package illutls

import (
	"net/url"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

// calculateSecFetchSite determines the relationship between the referer and target URL.
// It roughly mimics Chromium's GetRelationOfURLChainToOrigin logic.
func calculateSecFetchSite(referer string, targetURL *url.URL) string {
	if referer == "" || targetURL == nil {
		return "none"
	}
	refURL, err := url.Parse(referer)
	if err != nil {
		return "cross-site"
	}

	if refURL.Hostname() == targetURL.Hostname() {
		return "same-origin"
	}

	// Simplistic same-site check (checking if they share the root domain)
	refParts := strings.Split(refURL.Hostname(), ".")
	targetParts := strings.Split(targetURL.Hostname(), ".")
	
	if len(refParts) >= 2 && len(targetParts) >= 2 {
		refRoot := refParts[len(refParts)-2] + "." + refParts[len(refParts)-1]
		targetRoot := targetParts[len(targetParts)-2] + "." + targetParts[len(targetParts)-1]
		if refRoot == targetRoot {
			return "same-site"
		}
	}
	return "cross-site"
}

// BuildHeaders constructs the full HTTP header set for a request according
// to the profile's template and ordering rules.
//
// Dynamic headers (:authority, :path, host, etc.) are filled from the
// request itself; static headers (user-agent, sec-ch-ua, accept-encoding,
// etc.) come from the profile definition.
func BuildHeaders(profile *BrowserProfile, req *http.Request) http.Header {
	h := make(http.Header)

	// Heuristic Lock-Free Context Inference:
	// If the user explicitly sets Accept: application/json or provides a Referer,
	// we infer this is an API/Fetch request, not a top-level document navigation.
	acceptHeader := strings.ToLower(req.Header.Get("Accept"))
	isAPIRequest := strings.Contains(acceptHeader, "json") || 
					strings.Contains(acceptHeader, "*/*") || 
					req.Header.Get("Referer") != "" ||
					req.Header.Get("X-Requested-With") == "XMLHttpRequest"

	// 1. Apply default headers from the profile, intercepting sec-fetch-* if it's an API request.
	for k, v := range profile.Headers {
		if isAPIRequest {
			if k == "sec-fetch-dest" && req.Header.Get("sec-fetch-dest") == "" {
				v = "empty"
			} else if k == "sec-fetch-mode" && req.Header.Get("sec-fetch-mode") == "" {
				v = "cors"
			} else if k == "sec-fetch-site" && req.Header.Get("sec-fetch-site") == "" {
				v = calculateSecFetchSite(req.Header.Get("Referer"), req.URL)
			} else if k == "sec-fetch-user" && req.Header.Get("sec-fetch-user") == "" {
				continue // Do not include sec-fetch-user for API requests
			}
		}
		h.Set(k, v)
	}
	// 2. Always set User-Agent from the profile.
	h.Set("User-Agent", profile.UserAgent)
	// 3. Set host-related headers from the request.
	if req.Host != "" {
		h.Set("Host", req.Host)
	} else if req.URL != nil {
		h.Set("Host", req.URL.Host)
	}
	// 4. Apply header ordering.
	if len(profile.HeaderOrder) > 0 {
		h[http.HeaderOrderKey] = profile.HeaderOrder
	}
	if len(profile.PHeaderOrder) > 0 {
		h[http.PHeaderOrderKey] = profile.PHeaderOrder
	}
	return h
}

// MergeHeaders copies src headers into dst without overwriting existing keys.
func MergeHeaders(dst, src http.Header) {
	for k, vs := range src {
		if k == http.HeaderOrderKey || k == http.PHeaderOrderKey {
			// Order keys are special 鈥?always take from src if present.
			dst[k] = vs
			continue
		}
		if _, exists := dst[k]; !exists {
			dst[k] = vs
		}
	}
}
