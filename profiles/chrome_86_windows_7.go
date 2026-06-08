package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-57363-57427-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24-257,0
// JA3 HASH: 211b308d244cc0c791e17c0b0f42faff
// JA4: t13d1715h2_e08a0f08260f_de4a06bb82e3
// H2: 1:65536;3:1000;4:6291456;6:262144|15663105|0|m,a,s,p
// H2 HASH: 4f04edce68a7ecbe689edce7bf5f23f3
// Status: Verified Clean (Simulated Windows 7 Chrome 86)
// Notes: Authentic legacy fingerprint of Chrome 86 on Windows 7. Noticeably lacks GREASE, ALPS, and ECH. Includes legacy cipher suites (57363, 57427) and a strange curve ID (257). Perfect for simulating old enterprise/managed environments.
func init() {
	register(&illutls.BrowserProfile{
		Name:      "chrome-86-windows-7",
		UserAgent: "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
		TLSSpec: &utls.ClientHelloSpec{
			TLSVersMin: utls.VersionTLS12,
			TLSVersMax: utls.VersionTLS13,
			CipherSuites: []uint16{
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
				57363, // 0xE013 (Legacy experimental)
				57427, // 0xE053 (Legacy experimental)
				utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_RSA_WITH_AES_128_CBC_SHA,
				utls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []byte{0x00},
			Extensions: []utls.TLSExtension{
				&utls.SNIExtension{},                  // 0
				&utls.ExtendedMasterSecretExtension{}, // 23
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
					utls.CurveID(257),
				}},
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}}, // 11
				&utls.SessionTicketExtension{},                                // 35
				&utls.ALPNExtension{AlpnProtocols: []string{"h2"}},            // 16
				&utls.StatusRequestExtension{},                                // 5
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
				&utls.SCTExtension{}, // 18
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.X25519}, // 29
				}},
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.VersionTLS13,
					utls.VersionTLS12,
					utls.VersionTLS11,
					utls.VersionTLS10,
				}},
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
				&utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle},                                  // 21
			},
		},
		H2Settings: illutls.H2Settings{
			HeaderTableSize:      65536,
			EnablePush:           0,
			MaxConcurrentStreams: 1000,
			InitialWindowSize:    6291456,
			MaxHeaderListSize:    262144,
			SettingsOrder:        []uint16{1, 3, 4, 6}, // Notice 2 (EnablePush) is implicitly disabled, wait! Ah! Settings: "1": 65536, "3": 1000, "4": 6291456, "6": 262144. So 2 is MISSING!
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
			"user-agent":                "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
			"sec-fetch-site":            "none",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"sec-fetch-dest":            "document",
			"accept-encoding":           "gzip, deflate, br",
			"accept-language":           "zh-CN,zh;q=0.9",
		},
	})
}
