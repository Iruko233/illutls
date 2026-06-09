package illutls

import (
	"crypto/tls"
	"time"
)

// Options holds all configuration for a Client.
type Options struct {
	// ProfileName selects a registered BrowserProfile by name.
	ProfileName string
	// Profile is a user-supplied profile (takes precedence over ProfileName).
	Profile *BrowserProfile
	// Timeout is the total request timeout (default: 30s).
	Timeout time.Duration
	// TLSHandshakeTimeout limits the TLS handshake duration (default: 10s).
	TLSHandshakeTimeout time.Duration
	// ProxyURL is an optional HTTP/SOCKS5 proxy (e.g. "socks5://127.0.0.1:1080").
	ProxyURL string
	// TLSConfig provides additional tls.Config fields (e.g. InsecureSkipVerify).
	// ServerName is automatically set per-request and should not be set here.
	TLSConfig *tls.Config
	// FollowRedirects controls whether the client follows HTTP redirects (default: true).
	FollowRedirects bool
	// MaxRedirects limits redirect hops (default: 10).
	MaxRedirects int
	// DisableKeepAlives disables HTTP keep-alive connections.
	DisableKeepAlives bool
	// ShuffleExtensions dynamically randomizes the order of TLS extensions per connection.
	ShuffleExtensions bool
}

// Option is a functional option for configuring a Client.
type Option func(*Options)

// defaultOptions returns sane production defaults.
func defaultOptions() *Options {
	return &Options{
		Timeout:             30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		FollowRedirects:     true,
		MaxRedirects:        10,
	}
}

// WithProfile selects a registered browser profile by name.
func WithProfile(name string) Option {
	return func(o *Options) {
		o.ProfileName = name
	}
}

// WithCustomProfile provides a user-defined BrowserProfile directly.
func WithCustomProfile(p *BrowserProfile) Option {
	return func(o *Options) {
		o.Profile = p
	}
}

// WithTimeout sets the total request timeout.
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}

// WithTLSHandshakeTimeout sets the TLS handshake timeout.
func WithTLSHandshakeTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.TLSHandshakeTimeout = d
	}
}

// WithProxy sets an HTTP or SOCKS5 proxy URL.
func WithProxy(proxyURL string) Option {
	return func(o *Options) {
		o.ProxyURL = proxyURL
	}
}

// WithTLSConfig provides additional TLS configuration.
func WithTLSConfig(c *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = c
	}
}

// WithoutRedirects disables following HTTP redirects.
func WithoutRedirects() Option {
	return func(o *Options) {
		o.FollowRedirects = false
	}
}

// WithMaxRedirects sets the maximum number of redirect hops.
func WithMaxRedirects(n int) Option {
	return func(o *Options) {
		o.MaxRedirects = n
	}
}

// WithDisableKeepAlives disables HTTP keep-alive.
func WithDisableKeepAlives() Option {
	return func(o *Options) {
		o.DisableKeepAlives = true
	}
}

// WithShuffleExtensions dynamically randomizes the order of TLS extensions per connection
// to mimic Chrome's Extension Shuffling behavior and avoid TLS parroting.
func WithShuffleExtensions(b bool) Option {
	return func(o *Options) {
		o.ShuffleExtensions = b
	}
}
