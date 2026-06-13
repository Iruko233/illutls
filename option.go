package illutls

import (
	"crypto/tls"
	"hash/fnv"
	"strings"
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
	// GeoIPDBPath is the local path to the GeoLite2-Country.mmdb file.
	GeoIPDBPath string
	// RawLanguageList is the explicit raw language string to be parsed via Chrome's algorithm.
	RawLanguageList string
	// ForceRandomizeJA4 is a highly dangerous feature that persistently modifies the cipher suite
	// list to generate a randomized JA4 hash, useful only against primitive blocklists.
	ForceRandomizeJA4 bool
	// Dynamic profile generation parameters (evaluated in New)
	DynamicProfileSeed     *int64
	DynamicProfilePlatform string
	DynamicProfileVersion  int
}

// Option is a functional option for configuring a Client.
type Option func(*Options)

// WithDangerousJA4Randomization enables persistent randomization of the JA4 hash by deterministically
// mutating the cipher suites based on the profile seed.
// WARNING: This is a highly dangerous feature. Modifying the cipher suite list breaks the authentic
// Chrome fingerprint. Advanced WAFs (Cloudflare, Akamai) may flag this as an anomaly.
// Use ONLY against primitive self-built WAFs that blindly block JA4 hashes based on frequency.
func WithDangerousJA4Randomization() Option {
	return func(o *Options) {
		o.ForceRandomizeJA4 = true
	}
}

// defaultOptions returns sane production defaults.
func defaultOptions() *Options {
	return &Options{
		Timeout:             30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		FollowRedirects:     true,
		MaxRedirects:        10,
		ShuffleExtensions:   true,
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

// WithDynamicProfile generates a deterministic browser profile based on a seed.
// Seed can be an int64, int, or a string (which will be hashed to an int64).
// Usage 1: WithDynamicProfile("proxy_url") - Platform and version are deterministically picked based on the seed.
// Usage 2: WithDynamicProfile("proxy_url", "windows", 145) - Explicitly force platform and version.
func WithDynamicProfile(seed any, params ...any) Option {
	return func(o *Options) {
		var finalSeed int64
		switch v := seed.(type) {
		case string:
			h := fnv.New64a()
			h.Write([]byte(v))
			finalSeed = int64(h.Sum64() & 0x7FFFFFFFFFFFFFFF)
		case int64:
			if v < 0 {
				finalSeed = -v
			} else {
				finalSeed = v
			}
		case int:
			if v < 0 {
				finalSeed = int64(-v)
			} else {
				finalSeed = int64(v)
			}
		}

		platform := ""
		majorVersion := 0

		if len(params) >= 2 {
			if p, ok := params[0].(string); ok {
				platform = p
			}
			if m, ok := params[1].(int); ok {
				majorVersion = m
			}
		}

		if platform == "" || majorVersion == 0 {
			platforms := []string{"windows", "mac", "linux", "android", "ios"}
			idx := uint64(finalSeed) % uint64(len(platforms))
			platform = platforms[idx]

			versionRange := uint64(155 - 140 + 1)
			vIdx := (uint64(finalSeed) ^ 0x9e3779b97f4a7c15) % versionRange
			majorVersion = 140 + int(vIdx)
		}

		o.DynamicProfileSeed = &finalSeed
		o.DynamicProfilePlatform = platform
		o.DynamicProfileVersion = majorVersion
	}
}

// WithLanguage handles Accept-Language generation with extreme flexibility:
// 1. "auto": Automatically infers language from Proxy IP using the embedded GeoIP MMDB.
// 2. "/path/to/geo.mmdb": Automatically infers language using the provided external MMDB.
// 3. "FR": Uses a 2-letter ISO country code to dynamically generate Chrome locale strings.
// 4. "ja,en-US,en": Uses the explicit raw language string.
// All generated strings are passed through Chrome's exact mathematical q-value algorithm.
func WithLanguage(lang string) Option {
	return func(o *Options) {
		if lang == "auto" {
			o.GeoIPDBPath = "embedded"
			return
		}
		if strings.HasSuffix(lang, ".mmdb") {
			o.GeoIPDBPath = lang
			return
		}
		o.RawLanguageList = lang
	}
}
