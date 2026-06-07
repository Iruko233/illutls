package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-65281-18-65037-27-45-11-17613-16-10-51-43-5-13-23-35,29-23-24,0
// JA3 HASH: 7d18416ecd7c4feb78b11ac783f6057a
// JA4: t13d1516h2_8daaf6152771_d8a2da3f94cd
// H2: 1:65536;2:0;4:6291456;6:262144|15663105|0|m,a,s,p
// H2 HASH: 52d84b11737d980aef856699f885ca86
// Status: Verified Clean (Simulated Android 16 UCBrowser 18.8)
// Notes: UCBrowser features randomized extension shuffling, ALPS (17613), and ECH (65037). Noticeably lacks ML-KEM768 despite modern features. Headers lack sec-ch-ua.
func init() {
	register(&illutls.BrowserProfile{
		Name:      "android_16_uc_18",
		UserAgent: "Mozilla/5.0 (Linux; U; Android 16; zh-CN; PKG110 Build/UKQ1.231108.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/123.0.6312.80 UCBrowser/18.8.4.1510 Mobile Safari/537.36",
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
				&utls.SNIExtension{},        // 0
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.SCTExtension{},                  // 18
				&utls.GenericExtension{Id: 65037, Data: []byte{}}, // 65037 ECH
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}}, // 11
				&utls.GenericExtension{Id: 17613, Data: []byte{0x00, 0x03, 0x02, 0x68, 0x32}}, // 17613 ALPS
				&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}}, // 16
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
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.StatusRequestExtension{},        // 5
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
				&utls.ExtendedMasterSecretExtension{}, // 23
				&utls.SessionTicketExtension{}, // 35
				&utls.UtlsGREASEExtension{}, // Last GREASE
			},
		},
		H2Settings: illutls.H2Settings{
			HeaderTableSize:      65536,
			EnablePush:           0,
			InitialWindowSize:    6291456,
			MaxHeaderListSize:    262144,
			SettingsOrder:        []uint16{1, 2, 4, 6},
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
			"sec-fetch-dest",
			"accept-encoding",
			"priority",
		},
		Headers: map[string]string{
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (Linux; U; Android 16; zh-CN; PKG110 Build/UKQ1.231108.001) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/123.0.6312.80 UCBrowser/18.8.4.1510 Mobile Safari/537.36",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"sec-fetch-site":            "cross-site",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-dest":            "iframe",
			"accept-encoding":           "gzip, deflate",
			"priority":                  "u=0, i",
		},
	})
}
