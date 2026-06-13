package illutls

import (
	_ "embed"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

//go:embed GeoLite2-Country.mmdb
var embeddedMMDB []byte

// countryToLangCode maps an ISO 3166-1 alpha-2 country code to its primary language code.
var countryToLangCode = map[string]string{
	"US": "en", "GB": "en", "CA": "en", "AU": "en", "NZ": "en",
	"CN": "zh", "TW": "zh", "HK": "zh",
	"JP": "ja",
	"KR": "ko",
	"FR": "fr", "BE": "fr", "CH": "fr",
	"DE": "de", "AT": "de",
	"RU": "ru", "UA": "uk", "BY": "be",
	"ES": "es", "AR": "es", "MX": "es", "CO": "es",
	"IT": "it",
	"BR": "pt", "PT": "pt",
	"IN": "hi",
	"NL": "nl",
	"SE": "sv",
	"NO": "no",
	"FI": "fi",
	"DK": "da",
}

// resolveProxyLanguage looks up the proxy IP in the MMDB and returns the matching Accept-Language string.
// If mmdbPath is empty, it uses the embedded GeoLite2-Country.mmdb.
func resolveProxyLanguage(proxyURL string, mmdbPath string) (string, error) {
	pURL, err := url.Parse(proxyURL)
	if err != nil {
		return "", err
	}

	host := pURL.Hostname()
	ip := net.ParseIP(host)
	if ip == nil {
		// Try resolving the hostname (might block, usually proxies are passed as IPs in scraping)
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return "", fmt.Errorf("could not resolve proxy host: %s", host)
		}
		ip = ips[0]
	}

	var db *geoip2.Reader
	if mmdbPath != "" && mmdbPath != "embedded" {
		db, err = geoip2.Open(mmdbPath)
	} else {
		db, err = geoip2.FromBytes(embeddedMMDB)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to open mmdb: %v", err)
	}
	defer db.Close()

	record, err := db.Country(ip)
	if err != nil {
		return "", fmt.Errorf("failed to lookup ip in mmdb: %v", err)
	}

	isoCode := strings.ToUpper(record.Country.IsoCode)
	if isoCode == "" {
		return "", fmt.Errorf("no country found for ip")
	}

	langCode, exists := countryToLangCode[isoCode]
	if !exists {
		// Default fallback for unknown countries
		langCode = "en"
	}

	rawLangs := BuildRawLanguageList(langCode, isoCode)
	return GenerateAcceptLanguageHeader(rawLangs), nil
}
