package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49196-49195-52393-49200-49199-52392-49162-49161-49172-49171-157-156-53-47-49160-49170-10,0-23-65281-10-11-16-5-13-18-51-45-43-27-21,29-23-24-25,0
// JA3 HASH: 773906b0efdefa24a7f2b8eb6985bf37
// JA4: t13d2014h2_a09f3c656075_14788d8d241b
// H2: 2:0;4:2097152;3:100|10485760|0|m,s,p,a
// H2 HASH: ad8424af1cc590e09f7b0c499bf7fcdb
// Status: Verified Clean (Simulated iOS 17.6.1 Safari)
// Notes: Native iOS 17 Safari. Features the Apple-specific H2 settings order [2, 4, 3], zlib cert compression, and a duplicated signature algorithm (0x0805).
func init() {
	register(&illutls.BrowserProfile{
		Name:      "ios_17_safari_17",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 17_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6.1 Mobile/21G93 Safari/604.1",
		TLSSpec: &utls.ClientHelloSpec{
			TLSVersMin: utls.VersionTLS10,
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
				&utls.UtlsGREASEExtension{}, // First GREASE
				&utls.SNIExtension{},        // 0
				&utls.ExtendedMasterSecretExtension{}, // 23
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
					utls.CurveP521,
				}},
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}}, // 11
				&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}}, // 16
				&utls.StatusRequestExtension{},        // 5
				&utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []utls.SignatureScheme{ // 13
					utls.ECDSAWithP256AndSHA256,
					utls.PSSWithSHA256,
					utls.PKCS1WithSHA256,
					utls.ECDSAWithP384AndSHA384,
					utls.SignatureScheme(515), // ECDSAWithSHA1
					utls.PSSWithSHA384,        // 2053
					utls.PSSWithSHA384,        // 2053 (Duplicated by Apple)
					utls.PKCS1WithSHA384,
					utls.PSSWithSHA512,
					utls.PKCS1WithSHA512,
					utls.PKCS1WithSHA1, // 513
				}},
				&utls.SCTExtension{},                  // 18
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.X25519},
				}},
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
					utls.VersionTLS11,
					utls.VersionTLS10,
				}},
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionZlib}}, // 27
				&utls.UtlsGREASEExtension{}, // 56026
				&utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle}, // 21
			},
		},
		H2Settings: illutls.H2Settings{
			EnablePush:           0,
			MaxConcurrentStreams: 100,
			InitialWindowSize:    2097152,
			SettingsOrder:        []uint16{2, 4, 3}, // Apple Safari Specific Order
		},
		H2WindowUpdate: 10485760,
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
			"sec-fetch-site",
			"sec-fetch-dest",
			"accept-language",
			"sec-fetch-mode",
			"user-agent",
			"accept-encoding",
		},
		Headers: map[string]string{
			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"sec-fetch-site":  "cross-site",
			"sec-fetch-dest":  "iframe",
			"accept-language": "zh-CN,zh-Hans;q=0.9",
			"sec-fetch-mode":  "navigate",
			"user-agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 17_6_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.6.1 Mobile/21G93 Safari/604.1",
			"accept-encoding": "gzip, deflate, br",
		},
	})
}
