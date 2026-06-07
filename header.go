package illutls

import (
	http "github.com/bogdanfinn/fhttp"
)

// BuildHeaders constructs the full HTTP header set for a request according
// to the profile's template and ordering rules.
//
// Dynamic headers (:authority, :path, host, etc.) are filled from the
// request itself; static headers (user-agent, sec-ch-ua, accept-encoding,
// etc.) come from the profile definition.
func BuildHeaders(profile *BrowserProfile, req *http.Request) http.Header {
	h := make(http.Header)
	// 1. Apply default headers from the profile.
	for k, v := range profile.Headers {
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
