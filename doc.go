// package illutls provides an HTTP client that simulates
// real browser TLS fingerprints and HTTP headers.
//
// Illutls combines utls (for TLS ClientHello spoofing) with fhttp (for
// HTTP/2 frame-level fingerprint control) to present a complete, authentic
// browser identity on every connection.
//
// # Quick Start
//
//	client, err := illutls.New(illutls.WithProfile("chrome-149-windows-10"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	resp, err := client.Get("https://example.com")
//
// # Available Profiles
//
// Call ListProfiles() to see all registered browser profiles. Out of the box,
// Illutls ships with real-world fingerprints covering Chrome, Firefox,
// Edge, and Safari across Windows, macOS, Linux, Android, and iOS.
// If no profile is specified, it defaults to "chrome-149-windows-10".
//
// # Concurrency & Randomization
//
// Client is safe for concurrent use by multiple goroutines. Each TLS
// connection gets independently randomized GREASE values and extension shuffling,
// making repeat connections appear as natural browser traffic.
//
// # Features
//
//   - Extension Shuffling: Randomizes TLS extension order per connection (enabled by default).
//   - ECH Support: Injects compliant GREASE payloads into the ECH extension (65037).
//   - Proxy Support: Route traffic through HTTP/SOCKS5 proxies via WithProxy().
//
// # Custom Profiles
//
// Use WithCustomProfile() to register your own BrowserProfile with a custom
// TLS spec, HTTP/2 settings, and header ordering.
package illutls
