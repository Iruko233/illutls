package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4866-4867-4865-49196-49195-52393-49200-49199-52392-49162-49161-49172-49171-157-156-53-47-49160-49170-10,0-23-65281-10-11-16-5-13-18-51-45-43-27,4588-29-23-24-25,0
// JA3 HASH: ecdf4f49dd59effc439639da29186671
// JA4: t13d2013h2_a09f3c656075_7f0f34a4126d
// H2: 2:0;3:100;4:2097152;9:1|10420225|0|m,s,a,p
// H2 HASH: c52879e43202aeb92740be6e8c86ea96
// Status: Verified Clean (Simulated Safari 26)
// Notes: Anomalous Apple Network.framework fingerprint with a spoofed "Version/26.1 Safari" User-Agent. Uniquely starts TLS 1.3 ciphers with AES_256_GCM (4866) and adds ML-KEM768 (4588).
func init() {
	register(&illutls.BrowserProfile{
		Name:      "safari-26-ios-26",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.1 Mobile/15E148 Safari/604.1",
		TLSSpec: &utls.ClientHelloSpec{
			TLSVersMin: utls.VersionTLS12,
			TLSVersMax: utls.VersionTLS13,
			CipherSuites: []uint16{
				utls.GREASE_PLACEHOLDER,
				utls.TLS_AES_256_GCM_SHA384, // 4866
				utls.TLS_CHACHA20_POLY1305_SHA256,
				utls.TLS_AES_128_GCM_SHA256,
				utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				utls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_RSA_WITH_AES_256_CBC_SHA,
				utls.TLS_RSA_WITH_AES_128_CBC_SHA,
				49160,                              // utls.TLS_ECDHE_ECDSA_WITH_3DES_EDE_CBC_SHA
				49170,                              // utls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
				utls.TLS_RSA_WITH_3DES_EDE_CBC_SHA, // 10
			},
			CompressionMethods: []byte{0x00},
			Extensions: []utls.TLSExtension{
				&utls.UtlsGREASEExtension{},           // First GREASE
				&utls.SNIExtension{},                  // 0
				&utls.ExtendedMasterSecretExtension{}, // 23
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.X25519MLKEM768, // 4588
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
					utls.CurveP521,
				}},
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},  // 11
				&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}}, // 16
				&utls.StatusRequestExtension{},                                 // 5
				&utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []utls.SignatureScheme{ // 13
					utls.ECDSAWithP256AndSHA256,
					utls.PSSWithSHA256,
					utls.PKCS1WithSHA256,
					utls.ECDSAWithP384AndSHA384,
					utls.PSSWithSHA384, // 2053
					utls.PSSWithSHA384, // 2053 (Duplicated)
					utls.PKCS1WithSHA384,
					utls.PSSWithSHA512,
					utls.PKCS1WithSHA512,
					utls.PKCS1WithSHA1, // 513
				}},
				&utls.SCTExtension{}, // 18
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.X25519MLKEM768}, // 4588
					{Group: utls.X25519},         // 29
				}},
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionZlib}}, // 27
				&utls.UtlsGREASEExtension{}, // Last GREASE
			},
		},
		H2Settings: illutls.H2Settings{
			EnablePush:           0,       // 2
			MaxConcurrentStreams: 100,     // 3
			InitialWindowSize:    2097152, // 4
			HeaderTableSize:      0,
			MaxHeaderListSize:    0,
			NoRFC7540Priorities:  1,                    // 9
			SettingsOrder:        []uint16{2, 3, 4, 9}, // Wait, in JSON: 2:0;3:100;4:2097152;9:1
		},
		H2WindowUpdate: 10420225, // From JSON
		H2Priority: illutls.H2Priority{
			Weight:    0,
			DependsOn: 0,
			Exclusive: false,
		},
		PHeaderOrder: []string{
			":method",
			":scheme",
			":authority",
			":path",
		},
		HeaderOrder: []string{
			"sec-fetch-dest",
			"user-agent",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-user",
			"accept-language",
			"priority",
			"accept-encoding",
		},
		Headers: map[string]string{
			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"accept-encoding": "gzip, deflate, br",
			"accept-language": "zh-CN,zh-Hans;q=0.9",
			"priority":        "u=0, i",
			"sec-fetch-dest":  "document",
			"sec-fetch-mode":  "navigate",
			"sec-fetch-user":  "?1",
			"sec-fetch-site":  "none",
			"user-agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.1 Mobile/15E148 Safari/604.1",
		},
	})
}
