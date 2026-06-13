package illutls

import (
	"fmt"
	"math/rand"
	"strings"

	utls "github.com/refraction-networking/utls"
)

var (
	kGreaseyCharacters = []string{" ", "(", ":", "-", ".", "/", ")", ";", "=", "?", "_"}
	kGreasedVersions   = []string{"8", "99", "24"}

	// kCommonRawLanguages are typical raw base language strings without q-values.
	kCommonRawLanguages = []string{
		"en-US,en",
		"zh-CN,zh",
		"ja,en-US,en",
		"ru-RU,ru,en-US,en",
		"fr-FR,fr,en-US,en",
		"de-DE,de,en-US,en",
		"es-ES,es,en-US,en",
		"pt-BR,pt,en-US,en",
	}
)

// getGreasedBrandVersion mimics Chromium's GetGreasedUserAgentBrandVersion.
func getGreasedBrandVersion(seed int64) map[string]string {
	// Chrome uses seed mod length, not PRNG
	char1 := kGreaseyCharacters[seed%int64(len(kGreaseyCharacters))]
	char2 := kGreaseyCharacters[(seed+1)%int64(len(kGreaseyCharacters))]
	brand := fmt.Sprintf("Not%sA%sBrand", char1, char2)
	version := kGreasedVersions[seed%int64(len(kGreasedVersions))]

	return map[string]string{
		"brand":   brand,
		"version": version,
	}
}

// shuffleBrandList mimics Chromium's GenerateBrandVersionList shuffling.
func shuffleBrandList(brands []map[string]string, seed int64) []map[string]string {
	size := len(brands)
	if size == 3 {
		orders := [][]int{
			{0, 1, 2}, {0, 2, 1}, {1, 0, 2},
			{1, 2, 0}, {2, 0, 1}, {2, 1, 0},
		}
		order := orders[seed%6]
		return []map[string]string{
			brands[order[0]],
			brands[order[1]],
			brands[order[2]],
		}
	}
	// Fallback for other sizes, though Chrome standard is 3.
	r := rand.New(rand.NewSource(seed))
	shuffled := make([]map[string]string, len(brands))
	copy(shuffled, brands)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}


// GenerateProfile creates a deterministic BrowserProfile based on a seed.
// platform can be "windows", "mac", "linux", "android", "ios".
func GenerateProfile(seed int64, platform string, majorVersion int, forceRandomizeJA4 bool) *BrowserProfile {
	platform = strings.ToLower(platform)
	minorVersionSuffix := fmt.Sprintf("0.%d.%d", 4000+(seed%3000), 50+(seed%100))
	fullVersion := fmt.Sprintf("%d.%s", majorVersion, minorVersionSuffix)

	// 1. Generate Header Brands
	brands := []map[string]string{
		getGreasedBrandVersion(int64(majorVersion)),
		{"brand": "Chromium", "version": fmt.Sprintf("%d", majorVersion)},
		{"brand": "Google Chrome", "version": fmt.Sprintf("%d", majorVersion)},
	}
	shuffledBrands := shuffleBrandList(brands, int64(majorVersion))

	var secChUaParts []string
	var secChUaFullParts []string
	for _, b := range shuffledBrands {
		secChUaParts = append(secChUaParts, fmt.Sprintf(`"%s";v="%s"`, b["brand"], b["version"]))
		
		fullV := fullVersion
		if strings.Contains(b["brand"], "Not") && strings.Contains(b["brand"], "Brand") {
			fullV = b["version"] // GREASE version remains short
		}
		secChUaFullParts = append(secChUaFullParts, fmt.Sprintf(`"%s";v="%s"`, b["brand"], fullV))
	}
	secChUa := strings.Join(secChUaParts, ", ")
	secChUaFull := strings.Join(secChUaFullParts, ", ")

	// 2. Platform Specifics
	var uaOSString string
	var chPlatform string
	var chPlatformVersion string
	var chMobile string = "?0"
	var chArch string = `"x86"`
	var chBitness string = `"64"`
	var chModel string = `""`

	switch platform {
	case "windows", "win":
		uaOSString = "Windows NT 10.0; Win64; x64"
		chPlatform = `"Windows"`
		chPlatformVersion = `"10.0.0"`
	case "mac", "macos", "darwin":
		uaOSString = "Macintosh; Intel Mac OS X 10_15_7"
		chPlatform = `"macOS"`
		// pseudo-random deterministic macOS version
		chPlatformVersion = fmt.Sprintf(`"%d.%d.0"`, 13+(seed%3), seed%5)
	case "linux":
		uaOSString = "X11; Linux x86_64"
		chPlatform = `"Linux"`
		chPlatformVersion = `""`
	case "android":
		uaOSString = "Linux; Android 10; K"
		chPlatform = `"Android"`
		chPlatformVersion = `"10.0.0"`
		chMobile = "?1"
		chArch = `""`
		chBitness = `""`
		chModel = `"Pixel 5"` // Placeholder deterministic model
	case "ios":
		uaOSString = "iPhone; CPU iPhone OS 16_0 like Mac OS X"
		chPlatform = `"iOS"`
		chPlatformVersion = `"16.0"`
		chMobile = "?1"
		chArch = `"arm"`
		chModel = `"iPhone"`
	default:
		uaOSString = "Windows NT 10.0; Win64; x64"
		chPlatform = `"Windows"`
		chPlatformVersion = `"10.0.0"`
	}

	userAgent := fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", uaOSString, fullVersion)

	// 2.5 Deterministic Language
	rawLang := kCommonRawLanguages[seed%int64(len(kCommonRawLanguages))]
	acceptLanguage := GenerateAcceptLanguageHeader(rawLang)

	// 3. Assemble Headers
	headers := map[string]string{
		"sec-ch-ua":                   secChUa,
		"sec-ch-ua-mobile":            chMobile,
		"sec-ch-ua-platform":          chPlatform,
		"upgrade-insecure-requests":   "1",
		"user-agent":                  userAgent,
		"accept":                      "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"sec-fetch-site":              "none",
		"sec-fetch-mode":              "navigate",
		"sec-fetch-user":              "?1",
		"sec-fetch-dest":              "document",
		"sec-fetch-storage-access":    "active",
		"accept-encoding":             "gzip, deflate, br, zstd",
		"accept-language":             acceptLanguage,
		"priority":                    "u=0, i",
		"sec-ch-ua-platform-version":  chPlatformVersion,
		"sec-ch-ua-arch":              chArch,
		"sec-ch-ua-bitness":           chBitness,
		"sec-ch-ua-model":             chModel,
		"sec-ch-ua-full-version-list": secChUaFull,
	}

	headerOrder := []string{
		"sec-ch-ua",
		"sec-ch-ua-mobile",
		"sec-ch-ua-platform",
		"upgrade-insecure-requests",
		"user-agent",
		"accept",
		"sec-fetch-site",
		"sec-fetch-mode",
		"sec-fetch-user",
		"sec-fetch-dest",
		"sec-fetch-storage-access",
		"accept-encoding",
		"accept-language",
		"priority",
	}

	pHeaderOrder := []string{
		":method",
		":authority",
		":scheme",
		":path",
	}

	// 4. Base TLS Spec (Chrome 120+)
	tlsSpec := &utls.ClientHelloSpec{
		TLSVersMin: utls.VersionTLS12,
		TLSVersMax: utls.VersionTLS13,
		CipherSuites: []uint16{
			utls.GREASE_PLACEHOLDER,
			utls.TLS_AES_128_GCM_SHA256,
			utls.TLS_AES_256_GCM_SHA384,
			utls.TLS_CHACHA20_POLY1305_SHA256,
			utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			utls.TLS_RSA_WITH_AES_128_CBC_SHA,
			utls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		CompressionMethods: []byte{0x00},
		Extensions: []utls.TLSExtension{
			&utls.UtlsGREASEExtension{},           // First GREASE
			&utls.SNIExtension{},                  // 0
			&utls.ExtendedMasterSecretExtension{}, // 23
			&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
				{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
				{Group: utls.X25519MLKEM768}, // Post-Quantum Key Share
				{Group: utls.X25519},
			}},
			&utls.SessionTicketExtension{}, // 35
			&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
				utls.GREASE_PLACEHOLDER,
				utls.VersionTLS13,
				utls.VersionTLS12,
			}},
			&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}}, // 11
			&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
				utls.GREASE_PLACEHOLDER,
				utls.X25519MLKEM768, // Post-Quantum curve
				utls.X25519,
				utls.CurveP256,
				utls.CurveP384,
			}},
			&utls.SCTExtension{},                              // 18
			&utls.GenericExtension{Id: 65037, Data: []byte{
				0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x00, 0x20,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0x00, 0x40,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			}}, // ECH
			&utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []utls.SignatureScheme{ // 13
				utls.ECDSAWithP256AndSHA256,
				utls.PSSWithSHA256,
				utls.PKCS1WithSHA256,
				utls.ECDSAWithP384AndSHA384,
				utls.PSSWithSHA384,
				utls.PKCS1WithSHA384,
				utls.PSSWithSHA512,
				utls.PKCS1WithSHA512,
			}},
			&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient},                       // 65281
			&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}},                                 // 45
			&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},                                      // 16
			&utls.StatusRequestExtension{},                                                                      // 5
			&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
			&utls.GenericExtension{Id: 17613, Data: []byte{0x00, 0x03, 0x02, 0x68, 0x32}},                       // ALPS
			&utls.UtlsGREASEExtension{},                                                                         // Last GREASE
		},
	}

	if forceRandomizeJA4 {
		// Use a local random generator tied to the seed so the mutation is persistent for this profile
		rng := rand.New(rand.NewSource(seed))
		
		// The last 6 ciphers in the list are CBC/RSA variants that are safe to drop or shuffle
		// without breaking most modern handshakes.
		// Randomly drop between 0 and 3 of them.
		dropCount := rng.Intn(4) // 0 to 3
		if dropCount > 0 {
			origLength := len(tlsSpec.CipherSuites)
			// Remove 'dropCount' elements from the end of the slice
			tlsSpec.CipherSuites = tlsSpec.CipherSuites[:origLength-dropCount]
		}
		
		// We could also swap two ciphers to change the order, but JA4 *sorts* ciphers
		// so order swapping won't change the JA4 B-hash. Dropping them is the only way.
	}

	// 5. Build Chrome HTTP/2 Settings
	h2Settings := H2Settings{
		HeaderTableSize:   65536,
		EnablePush:        0,
		InitialWindowSize: 6291456,
		MaxHeaderListSize: 262144,
		SettingsOrder:     []uint16{1, 2, 4, 6},
	}
	
	h2WindowUpdate := uint32(15663105)

	profileName := fmt.Sprintf("dynamic-seed-%d-%s-%d", seed, platform, majorVersion)

	return &BrowserProfile{
		Name:           profileName,
		UserAgent:      userAgent,
		TLSSpec:        tlsSpec,
		H2Settings:     h2Settings,
		H2WindowUpdate: h2WindowUpdate,
		H2Priority: H2Priority{
			Weight:    0,
			DependsOn: 0,
			Exclusive: false,
		},
		HeaderOrder:  headerOrder,
		PHeaderOrder: pHeaderOrder,
		Headers:      headers,
	}
}
