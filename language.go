package illutls

import (
	"fmt"
	"strings"
)

// GenerateAcceptLanguageHeader perfectly mimics Chromium's net::HttpUtil::GenerateAcceptLanguageHeader
// It takes a comma-separated list of languages (e.g. "ja,en-US,en") and applies descending q-values.
func GenerateAcceptLanguageHeader(rawLanguageList string) string {
	qvalue10 := 10
	tokens := strings.Split(rawLanguageList, ",")
	var langListWithQ []string

	for _, lang := range tokens {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			continue
		}
		if qvalue10 == 10 {
			// q=1.0 is implicit
			langListWithQ = append(langListWithQ, lang)
		} else {
			// e.g. ;q=0.9
			langListWithQ = append(langListWithQ, fmt.Sprintf("%s;q=0.%d", lang, qvalue10))
		}
		// It does not make sense to have 'q=0'.
		if qvalue10 > 1 {
			qvalue10 -= 1
		}
	}
	return strings.Join(langListWithQ, ",")
}

// BuildRawLanguageList perfectly mimics Chrome's OS locale fallback algorithm.
// When an OS is set to a specific locale (e.g., "ja-JP"), Chrome builds the Accept-Language
// by taking the specific locale, then the base language, and appending English fallbacks.
func BuildRawLanguageList(lang, country string) string {
	lang = strings.ToLower(lang)
	country = strings.ToUpper(country)

	if lang == "en" {
		if country == "US" {
			return "en-US,en"
		}
		return fmt.Sprintf("en-%s,en-US,en", country)
	}

	if lang == "zh" {
		if country == "CN" {
			return "zh-CN,zh"
		}
		return fmt.Sprintf("zh-%s,zh-CN,zh", country)
	}

	return fmt.Sprintf("%s-%s,%s,en-US,en", lang, country, lang)
}
