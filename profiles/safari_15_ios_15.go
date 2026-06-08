package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49196-49195-52393-49200-49199-52392-49188-49187-49162-49161-49192-49191-49172-49171-157-156-61-60-53-47-49160-49170-10,0-23-65281-10-11-16-5-13-18-51-45-43-27-21,29-23-24-25,0
// JA3 HASH: c59b5aeb69936c251f090be89e1c4ca5
// JA4: t13d2614h2_2802a3db6c62_14788d8d241b
// H2: 4:2097152;3:100|10485760|0|m,s,p,a
// H2 HASH: d5fcbdc393757341115a861bf8d23265
// Status: Verified Clean (Simulated Safari 15)
// Notes: iOS 15.1 Safari. Standard cipher order. Includes padding extension (21) and duplicated 2053 signature algorithm.
func init() {
	register(&illutls.BrowserProfile{
		Name:      "safari-15-ios-15",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 15_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Mobile/15E148 Safari/604.1",
		TLSSpec: &utls.ClientHelloSpec{
			TLSVersMin: utls.VersionTLS12,
			TLSVersMax: utls.VersionTLS13,
			CipherSuites: []uint16{
				utls.GREASE_PLACEHOLDER,
				utls.TLS_AES_128_GCM_SHA256,
				utls.TLS_AES_256_GCM_SHA384,
				utls.TLS_CHACHA20_POLY1305_SHA256,
				utls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				utls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				utls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
				49188, // utls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA384
				49187, // utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256
				utls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				utls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				49192, // utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384
				49191, // utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
				utls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				utls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				utls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				utls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				61, // utls.TLS_RSA_WITH_AES_256_CBC_SHA256
				60, // utls.TLS_RSA_WITH_AES_128_CBC_SHA256
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
					utls.ECDSAWithSHA1, // 515
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
					{Group: utls.X25519}, // 29
				}},
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionZlib}}, // 27
				&utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle},                                // 21
				&utls.UtlsGREASEExtension{}, // Last GREASE
			},
		},
		H2Settings: illutls.H2Settings{
			MaxConcurrentStreams: 100,     // 3
			InitialWindowSize:    2097152, // 4
			SettingsOrder:        []uint16{4, 3},
		},
		H2WindowUpdate: 10485760, // From JSON
		H2Priority: illutls.H2Priority{
			Weight:    0,
			DependsOn: 0,
			Exclusive: false,
		},
		PHeaderOrder: []string{
			":method",
			":scheme",
			":path",
			":authority",
		},
		HeaderOrder: []string{
			"accept",
			"accept-encoding",
			"user-agent",
			"accept-language",
		},
		Headers: map[string]string{
			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"accept-encoding": "gzip, deflate, br",
			"accept-language": "zh-CN,zh-Hans;q=0.9",
			"user-agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 15_1_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Mobile/15E148 Safari/604.1",
		},
	})
}
