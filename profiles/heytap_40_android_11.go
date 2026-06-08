package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,65281-43-10-51-13-23-45-35-17513-0-18-27-5-11-16-21,29-23-24,0
// JA3 HASH: aac9c0ab80a72f3533ffafbea6ebf09a
// JA4: t13d1516h2_8daaf6152771_e5627efa2ab1
// H2: 1:65536;2:0;3:1000;4:6291456;6:262144|15663105|0|m,a,s,p
// H2 HASH: a345a694846ad9f6c97bcc3c75adbe26
// Status: Verified Clean (Simulated Android 11 HeyTapBrowser 40)
// Notes: OPPO HeyTapBrowser. Based on Chrome 115 but features randomized extension shuffling (with Padding strictly at the end). Uses ALPS 17513. Lacks ECH. Missing sec-ch-ua headers.
func init() {
	register(&illutls.BrowserProfile{
		Name:      "heytap-40-android-11",
		UserAgent: "Mozilla/5.0 (Linux; U; Android 11; zh-cn; PCAM00 Build/RKQ1.201217.002) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/115.0.5790.168 Mobile Safari/537.36 HeyTapBrowser/40.10.16.2",
		TLSSpec: &utls.ClientHelloSpec{
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
				utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256, // 52393
				utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,   // 52392
				utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_RSA_WITH_AES_128_CBC_SHA,
				utls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []byte{0x00},
			Extensions: []utls.TLSExtension{
				&utls.UtlsGREASEExtension{}, // First GREASE
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
				}},
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.X25519}, // 29
				}},
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
				&utls.ExtendedMasterSecretExtension{},                                         // 23
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}},           // 45
				&utls.SessionTicketExtension{},                                                // 35
				&utls.GenericExtension{Id: 17513, Data: []byte{0x00, 0x03, 0x02, 0x68, 0x32}}, // 17513 ALPS
				&utls.SNIExtension{},                                                          // 0
				&utls.SCTExtension{},                                                          // 18
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
				&utls.StatusRequestExtension{},                                     // 5
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},      // 11
				&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},     // 16
				&utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle}, // 21
				&utls.UtlsGREASEExtension{},                                        // Last GREASE
			},
		},
		H2Settings: illutls.H2Settings{
			HeaderTableSize:      65536,
			EnablePush:           0,
			MaxConcurrentStreams: 1000,
			InitialWindowSize:    6291456,
			MaxHeaderListSize:    262144,
			SettingsOrder:        []uint16{1, 2, 3, 4, 6}, // Notice 3: 1000 is present
		},
		H2WindowUpdate: 15663105,
		H2Priority: illutls.H2Priority{
			Weight:    0,
			DependsOn: 0,
			Exclusive: false,
		},
		PHeaderOrder: []string{
			":method",
			":authority",
			":scheme",
			":path",
		},
		HeaderOrder: []string{
			"upgrade-insecure-requests",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
		},
		Headers: map[string]string{
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (Linux; U; Android 11; zh-cn; PCAM00 Build/RKQ1.201217.002) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/115.0.5790.168 Mobile Safari/537.36 HeyTapBrowser/40.10.16.2",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"sec-fetch-site":            "none",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-dest":            "document",
			"sec-fetch-user":            "?1",
			"accept-encoding":           "gzip, deflate, br",
			"accept-language":           "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		},
	})
}
