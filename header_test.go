package illutls

import (
	http "github.com/bogdanfinn/fhttp"
	"net/url"
	"testing"
)

func TestBuildHeadersAutoAdaptive(t *testing.T) {
	profile := &BrowserProfile{
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		Headers: map[string]string{
			"sec-fetch-dest": "document",
			"sec-fetch-mode": "navigate",
			"sec-fetch-site": "none",
			"sec-fetch-user": "?1",
			"accept":         "text/html,application/xhtml+xml,application/xml",
		},
	}

	targetURL, _ := url.Parse("https://api.example.com/data.json")

	// 1. Navigation Request (Default)
	req1, _ := http.NewRequest("GET", targetURL.String(), nil)
	h1 := BuildHeaders(profile, req1)
	
	if h1.Get("sec-fetch-dest") != "document" {
		t.Errorf("Expected document, got %s", h1.Get("sec-fetch-dest"))
	}
	if h1.Get("sec-fetch-user") != "?1" {
		t.Errorf("Expected ?1, got %s", h1.Get("sec-fetch-user"))
	}

	// 2. API Request (via Accept JSON)
	req2, _ := http.NewRequest("GET", targetURL.String(), nil)
	req2.Header.Set("Accept", "application/json")
	h2 := BuildHeaders(profile, req2)
	
	if h2.Get("sec-fetch-dest") != "empty" {
		t.Errorf("Expected empty, got %s", h2.Get("sec-fetch-dest"))
	}
	if h2.Get("sec-fetch-mode") != "cors" {
		t.Errorf("Expected cors, got %s", h2.Get("sec-fetch-mode"))
	}
	if h2.Get("sec-fetch-user") != "" {
		t.Errorf("Expected empty sec-fetch-user, got %s", h2.Get("sec-fetch-user"))
	}

	// 3. API Request (via Referer, Same-Site)
	req3, _ := http.NewRequest("GET", targetURL.String(), nil)
	req3.Header.Set("Referer", "https://www.example.com/page")
	h3 := BuildHeaders(profile, req3)

	if h3.Get("sec-fetch-dest") != "empty" {
		t.Errorf("Expected empty, got %s", h3.Get("sec-fetch-dest"))
	}
	if h3.Get("sec-fetch-site") != "same-site" {
		t.Errorf("Expected same-site, got %s", h3.Get("sec-fetch-site"))
	}
}
