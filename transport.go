package illutls

import (
	"context"
	"crypto/tls"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	http2 "github.com/bogdanfinn/fhttp/http2"
	bogdanutls "github.com/bogdanfinn/utls"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
	"net"
	"net/url"
	"sync"
	"time"
)

// Transport implements http.RoundTripper with browser-identical TLS and
// HTTP/2 fingerprints.
type Transport struct {
	profile   *BrowserProfile
	dialer    *net.Dialer
	proxyFunc func() (*url.URL, error)
	tlsCfg    *tls.Config
	// h2Transport is the underlying fhttp HTTP/2 transport.
	h2Transport *http2.Transport
	mu          sync.Mutex
	cachedConns map[string]net.Conn // host 鈫?reusable conn
}

// NewTransport creates a Transport for the given profile.
func NewTransport(profile *BrowserProfile, opts *Options) (*Transport, error) {
	t := &Transport{
		profile: profile,
		dialer: &net.Dialer{
			Timeout:   opts.Timeout,
			KeepAlive: 30 * time.Second,
		},
		cachedConns: make(map[string]net.Conn),
	}
	// Proxy setup.
	if opts.ProxyURL != "" {
		u, err := url.Parse(opts.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidProxy, err)
		}
		t.proxyFunc = func() (*url.URL, error) { return u, nil }
	}
	// TLS config 鈥?copy user-supplied or create default.
	if opts.TLSConfig != nil {
		t.tlsCfg = opts.TLSConfig.Clone()
	} else {
		t.tlsCfg = &tls.Config{}
	}
	// Build the fhttp HTTP/2 transport with profile's H2 settings.
	t.h2Transport = &http2.Transport{
		DialTLS:            t.dialTLS,
		DisableCompression: false,
	}
	ApplyH2Settings(t.h2Transport, profile.H2Settings, profile.H2WindowUpdate, profile.H2Priority)
	return t, nil
}

// RoundTrip executes a single HTTP transaction, applying the profile's
// headers and TLS fingerprint.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Build profile headers and merge into the request.
	profileHeaders := BuildHeaders(t.profile, req)
	MergeHeaders(req.Header, profileHeaders)
	return t.h2Transport.RoundTrip(req)
}

// dialTLS performs a TCP dial (optionally through a proxy) and wraps the
// connection with utls using the profile's ClientHelloSpec.
func (t *Transport) dialTLS(network, addr string, cfg *bogdanutls.Config) (net.Conn, error) {
	ctx := context.Background()
	return t.dialTLSContext(ctx, network, addr)
}

// dialTLSContext is the context-aware TLS dial function.
func (t *Transport) dialTLSContext(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	// 1. Establish raw TCP connection (with optional proxy).
	rawConn, err := t.dialRaw(ctx, network, addr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err)
	}
	// 2. Clone and randomize the TLS spec for this connection.
	spec := CloneClientHelloSpec(t.profile.TLSSpec)
	RandomizeGREASE(spec)
	// 3. Wrap with utls.
	tlsCfg := &utls.Config{
		ServerName:         host,
		InsecureSkipVerify: t.tlsCfg.InsecureSkipVerify,
		NextProtos:         []string{"h2", "http/1.1"},
	}
	uConn := utls.UClient(rawConn, tlsCfg, utls.HelloCustom)
	if err := uConn.ApplyPreset(spec); err != nil {
		rawConn.Close()
		return nil, fmt.Errorf("%w: apply preset: %s", ErrTLSHandshake, err)
	}
	// 4. Handshake.
	if err := uConn.HandshakeContext(ctx); err != nil {
		rawConn.Close()
		return nil, fmt.Errorf("%w: %s", ErrTLSHandshake, err)
	}
	return uConn, nil
}

// dialRaw establishes a raw TCP connection, optionally through a proxy.
func (t *Transport) dialRaw(ctx context.Context, network, addr string) (net.Conn, error) {
	if t.proxyFunc != nil {
		proxyURL, err := t.proxyFunc()
		if err != nil {
			return nil, err
		}
		switch proxyURL.Scheme {
		case "socks5", "socks5h":
			auth := (*proxy.Auth)(nil)
			if proxyURL.User != nil {
				pass, _ := proxyURL.User.Password()
				auth = &proxy.Auth{
					User:     proxyURL.User.Username(),
					Password: pass,
				}
			}
			d, err := proxy.SOCKS5(network, proxyURL.Host, auth, t.dialer)
			if err != nil {
				return nil, err
			}
			return d.Dial(network, addr)
		case "http", "https":
			// For HTTP proxies, use CONNECT method via plain dial.
			return t.dialer.DialContext(ctx, network, proxyURL.Host)
		}
	}
	return t.dialer.DialContext(ctx, network, addr)
}

// Close shuts down idle connections.
func (t *Transport) Close() {
	t.mu.Lock()
	for _, c := range t.cachedConns {
		c.Close()
	}
	t.cachedConns = make(map[string]net.Conn)
	t.mu.Unlock()
}
