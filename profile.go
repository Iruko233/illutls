package illutls

import (
	utls "github.com/refraction-networking/utls"
	"sync"
)

// BrowserProfile defines the complete fingerprint for a browser identity.
// It is immutable after creation 鈥?safe to share across goroutines.
type BrowserProfile struct {
	// Name is a human-readable identifier, e.g. "Chrome 141 Windows".
	Name string
	// UserAgent is the User-Agent header string.
	UserAgent string
	// TLSSpec defines the TLS ClientHello parameters for utls.
	// GREASE placeholders are randomized at connection time.
	TLSSpec *utls.ClientHelloSpec
	// H2Settings configures the HTTP/2 SETTINGS frame.
	H2Settings H2Settings
	// H2WindowUpdate is the WINDOW_UPDATE increment sent after SETTINGS.
	H2WindowUpdate uint32
	// H2Priority configures the PRIORITY field in HEADERS frames.
	H2Priority H2Priority
	// HeaderOrder defines the order of standard HTTP headers (lowercase).
	HeaderOrder []string
	// PHeaderOrder defines the order of HTTP/2 pseudo-headers.
	PHeaderOrder []string
	// Headers supplies default header values (sec-ch-ua, accept, etc.).
	// Keys are canonical HTTP header names.
	Headers map[string]string
}

// H2Settings corresponds to an HTTP/2 SETTINGS frame.
type H2Settings struct {
	HeaderTableSize      uint32
	EnablePush           uint32
	MaxConcurrentStreams uint32
	InitialWindowSize    uint32
	MaxFrameSize         uint32
	MaxHeaderListSize    uint32
	// NoRFC7540Priorities corresponds to SETTINGS_NO_RFC7540_PRIORITIES (0x09).
	NoRFC7540Priorities uint32

	// SettingsOrder explicitly defines the order of settings. If empty, a default order is used.
	SettingsOrder []uint16
}

// H2Priority corresponds to the PRIORITY field in an HTTP/2 HEADERS frame.
type H2Priority struct {
	Weight    uint8
	DependsOn uint32
	Exclusive bool
}

// ---------------------------------------------------------------------------
// Global profile registry
// ---------------------------------------------------------------------------
var (
	profilesMu sync.RWMutex
	profiles   = make(map[string]*BrowserProfile)
)

// RegisterProfile adds a browser profile to the global registry.
// If a profile with the same name already exists it is overwritten.
func RegisterProfile(p *BrowserProfile) {
	profilesMu.Lock()
	profiles[p.Name] = p
	profilesMu.Unlock()
}

// GetProfile returns a registered profile by name, or nil if not found.
func GetProfile(name string) *BrowserProfile {
	profilesMu.RLock()
	p := profiles[name]
	profilesMu.RUnlock()
	return p
}

// ListProfiles returns the names of all registered profiles.
func ListProfiles() []string {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	names := make([]string, 0, len(profiles))
	for k := range profiles {
		names = append(names, k)
	}
	return names
}
