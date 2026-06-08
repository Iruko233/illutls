package illutls

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	http2 "github.com/bogdanfinn/fhttp/http2"
	bogdanutls "github.com/bogdanfinn/utls"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

var ErrHTTP11Negotiated = errors.New("http/1.1 negotiated")

type bufferedConn struct {
	net.Conn
	r *bufio.Reader
}

func (b *bufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

// Transport implements http.RoundTripper with browser-identical TLS and
// HTTP/2 fingerprints.
type Transport struct {
	profile   *BrowserProfile
	dialer    *net.Dialer
	proxyFunc func() (*url.URL, error)
	tlsCfg    *tls.Config
	// h2Transport is the underlying fhttp HTTP/2 transport.
	h2Transport *http2.Transport
	// h1Transport is the fallback HTTP/1.1 transport.
	h1Transport *http.Transport
	mu          sync.Mutex
	stashedConns map[string]net.Conn
}

// NewTransport creates a Transport for the given profile.
func NewTransport(profile *BrowserProfile, opts *Options) (*Transport, error) {
	t := &Transport{
		profile: profile,
		dialer: &net.Dialer{
			Timeout:   opts.Timeout,
			KeepAlive: 30 * time.Second,
		},
		stashedConns: make(map[string]net.Conn),
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
	
	// Build the fallback HTTP/1.1 transport.
	t.h1Transport = &http.Transport{
		DialContext:           t.dialRaw,
		DialTLSContext:        t.dialTLSForH1,
		ForceAttemptHTTP2:     false,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   opts.TLSHandshakeTimeout,
		DisableKeepAlives:     opts.DisableKeepAlives,
		ExpectContinueTimeout: 1 * time.Second,
	}
	
	return t, nil
}

// RoundTrip executes a single HTTP transaction, applying the profile's
// headers and TLS fingerprint.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to prevent data races and mutating the original.
	req = req.Clone(req.Context())

	// Build profile headers and merge into the request.
	profileHeaders := BuildHeaders(t.profile, req)
	MergeHeaders(req.Header, profileHeaders)
	
	// Plain HTTP does not support HTTP/2 over cleartext in this context.
	if req.URL.Scheme == "http" {
		return t.h1Transport.RoundTrip(req)
	}
	
	resp, err := t.h2Transport.RoundTrip(req)
	if err != nil && (errors.Is(err, ErrHTTP11Negotiated) || strings.Contains(err.Error(), ErrHTTP11Negotiated.Error())) {
		// ALPN negotiated HTTP/1.1. Hijack the connection and route to h1Transport.
		respH1, errH1 := t.h1Transport.RoundTrip(req)
		
		// Cleanup the stash in case h1Transport reused an idle connection and didn't consume the stash.
		addr := req.URL.Host
		if !strings.Contains(addr, ":") {
			if req.URL.Scheme == "https" {
				addr += ":443"
			} else {
				addr += ":80"
			}
		}
		t.mu.Lock()
		if conn, ok := t.stashedConns[addr]; ok {
			conn.Close()
			delete(t.stashedConns, addr)
		}
		t.mu.Unlock()
		
		return respH1, errH1
	}
	return resp, err
}

// dialTLS performs a TCP dial (optionally through a proxy) and wraps the
// connection with utls using the profile's ClientHelloSpec.
func (t *Transport) dialTLS(network, addr string, cfg *bogdanutls.Config) (net.Conn, error) {
	ctx := context.Background()
	uConn, err := t.dialTLSContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	
	if uConn.(*utls.UConn).ConnectionState().NegotiatedProtocol == "http/1.1" {
		t.mu.Lock()
		if old, ok := t.stashedConns[addr]; ok {
			old.Close()
		}
		t.stashedConns[addr] = uConn
		t.mu.Unlock()
		return nil, ErrHTTP11Negotiated
	}
	
	return uConn, nil
}

func (t *Transport) dialTLSForH1(ctx context.Context, network, addr string) (net.Conn, error) {
	t.mu.Lock()
	if conn, ok := t.stashedConns[addr]; ok {
		delete(t.stashedConns, addr)
		t.mu.Unlock()
		return conn, nil
	}
	t.mu.Unlock()
	
	// Fallback to dialing explicitly if stash is empty
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
	nextProtos := []string{"h2", "http/1.1"}
	if len(t.tlsCfg.NextProtos) > 0 {
		nextProtos = t.tlsCfg.NextProtos
	}
	// 3. Wrap with utls.
	tlsCfg := &utls.Config{
		ServerName:         host,
		InsecureSkipVerify: t.tlsCfg.InsecureSkipVerify,
		NextProtos:         nextProtos,
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
			conn, err := t.dialer.DialContext(ctx, network, proxyURL.Host)
			if err != nil {
				return nil, err
			}
			req := &http.Request{
				Method: "CONNECT",
				URL:    &url.URL{Opaque: addr},
				Host:   addr,
				Header: make(http.Header),
			}
			if proxyURL.User != nil {
				pass, _ := proxyURL.User.Password()
				auth := proxyURL.User.Username() + ":" + pass
				basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
				req.Header.Set("Proxy-Authorization", basicAuth)
			}
			if err := req.Write(conn); err != nil {
				conn.Close()
				return nil, fmt.Errorf("proxy write failed: %w", err)
			}
			br := bufio.NewReader(conn)
			resp, err := http.ReadResponse(br, req)
			if err != nil {
				conn.Close()
				return nil, fmt.Errorf("proxy read failed: %w", err)
			}
			if resp.StatusCode != 200 {
				conn.Close()
				return nil, fmt.Errorf("proxy error: %s", resp.Status)
			}
			if br.Buffered() > 0 {
				return &bufferedConn{Conn: conn, r: br}, nil
			}
			return conn, nil
		}
	}
	return t.dialer.DialContext(ctx, network, addr)
}

// Close shuts down idle connections.
func (t *Transport) Close() {
	t.mu.Lock()
	for _, c := range t.stashedConns {
		c.Close()
	}
	t.stashedConns = make(map[string]net.Conn)
	t.mu.Unlock()
	
	if t.h1Transport != nil {
		t.h1Transport.CloseIdleConnections()
	}
	if t.h2Transport != nil {
		t.h2Transport.CloseIdleConnections()
	}
}
