package illutls

import (
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"strings"
)

// Client is the main entry point for Illutls. It provides browser-identical
// HTTP requests with automatic TLS and header fingerprinting.
//
// Client is safe for concurrent use by multiple goroutines.
type Client struct {
	httpClient *http.Client
	profile    *BrowserProfile
	transport  *Transport
	opts       *Options
}

// New creates an Illutls Client configured with the given options.
//
// Example:
//
//	client, err := illutls.New(illutls.WithProfile("Chrome 141 Windows"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	resp, err := client.Get("https://example.com")
func New(opts ...Option) (*Client, error) {
	o := defaultOptions()
	for _, fn := range opts {
		fn(o)
	}
	// Resolve the profile.
	profile := o.Profile
	if o.DynamicProfileSeed != nil {
		profile = GenerateProfile(*o.DynamicProfileSeed, o.DynamicProfilePlatform, o.DynamicProfileVersion, o.ForceRandomizeJA4)
	}

	if profile == nil {
		if o.ProfileName == "" {
			// Default to Chrome 141 Windows if nothing specified.
			o.ProfileName = "Chrome 141 Windows"
		}
		profile = GetProfile(o.ProfileName)
		if profile == nil {
			return nil, fmt.Errorf("%w: %q (available: %s)",
				ErrProfileNotFound, o.ProfileName, strings.Join(ListProfiles(), ", "))
		}
	}

	// Determine Accept-Language based on user options.
	rawLanguage := "en-US,en" // Default fallback

	if o.RawLanguageList != "" {
		// 1. Explicit override via WithLanguage
		if len(o.RawLanguageList) == 2 {
			isoCode := strings.ToUpper(o.RawLanguageList)
			langCode, exists := countryToLangCode[isoCode]
			if !exists {
				langCode = "en"
			}
			rawLanguage = BuildRawLanguageList(langCode, isoCode)
		} else {
			rawLanguage = o.RawLanguageList
		}
	} else if o.GeoIPDBPath != "" && o.ProxyURL != "" {
		// 2. Proxy Auto-Adaptation (MMDB Lookup)
		proxyLang, err := resolveProxyLanguage(o.ProxyURL, o.GeoIPDBPath)
		if err == nil && proxyLang != "" {
			// resolveProxyLanguage already runs the Chrome algorithm
			profile.Headers["accept-language"] = proxyLang
			rawLanguage = "" // Prevent re-running algorithm below
		}
	} else if o.ProfileName == "" && o.Profile != nil {
		// 3. Dynamic Profile Seed extraction (pseudo-random language based on seed)
		// We can infer this is a dynamic profile if they used WithDynamicProfile but no language.
		// Let's assume if it's a dynamic profile, it should have a seed-based language.
		// Wait, we need the seed to do pseudo-random. Since the profile is already generated, 
		// we should actually inject it inside GenerateProfile in generator.go.
	}

	if rawLanguage != "" {
		profile.Headers["accept-language"] = GenerateAcceptLanguageHeader(rawLanguage)
	}

	// Build the transport.
	transport, err := NewTransport(profile, o)
	if err != nil {
		return nil, err
	}
	// Build the HTTP client.
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   o.Timeout,
	}
	if !o.FollowRedirects {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if o.MaxRedirects > 0 {
		max := o.MaxRedirects
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= max {
				return fmt.Errorf("stopped after %d redirects", max)
			}
			return nil
		}
	}
	return &Client{
		httpClient: httpClient,
		profile:    profile,
		transport:  transport,
		opts:       o,
	}, nil
}

// Do sends an HTTP request and returns an HTTP response.
// The caller is responsible for closing the response body.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}

// Get issues a GET request to the specified URL.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post issues a POST request to the specified URL.
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return c.Do(req)
}

// Head issues a HEAD request to the specified URL.
func (c *Client) Head(url string) (*http.Response, error) {
	req, err := c.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// NewRequest creates a new HTTP request with the profile's headers
// pre-populated. The caller can override any header before calling Do().
func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// Pre-populate profile headers.
	profileHeaders := BuildHeaders(c.profile, req)
	MergeHeaders(req.Header, profileHeaders)
	return req, nil
}

// Profile returns the active BrowserProfile.
func (c *Client) Profile() *BrowserProfile {
	return c.profile
}

// Close releases resources held by the Client.
func (c *Client) Close() {
	c.transport.Close()
}
