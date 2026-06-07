// package illutls provides an HTTP client that simulates
// real browser TLS fingerprints and HTTP headers.
//
// Illutls combines utls (for TLS ClientHello spoofing) with fhttp (for
// HTTP/2 frame-level fingerprint control) to present a complete, authentic
// browser identity on every connection.
//
// # Quick Start
//
//	client, err := illutls.New(illutls.WithProfile("Chrome 141 Windows"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	resp, err := client.Get("https://example.com")
//
// # Available Profiles
//
// Call ListProfiles() to see all registered browser profiles. Out of the box,
// Illutls ships with 10 real-world fingerprints covering Chrome, Firefox,
// Edge, and Safari across Windows, macOS, Linux, Android, and iOS.
//
// # Concurrency
//
// Client is safe for concurrent use by multiple goroutines. Each TLS
// connection gets independently randomized GREASE values, making repeat
// connections appear as natural browser traffic.
//
// # Custom Profiles
//
// Use WithCustomProfile() to register your own BrowserProfile with a custom
// TLS spec, HTTP/2 settings, and header ordering.
package illutls
