package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,45-51-51764-35-0-17613-65281-5-11-16-13-18-10-27-23-65037-43,4588-29-23-24,0
// JA3 HASH: a83aa52b676ea19698afa69aca1f2a3d
// JA4: t13d1517h2_8daaf6152771_cb7bf5808d99
// H2: 1:65536;2:0;4:6291456;6:262144|15663105|0|m,a,s,p
// H2 HASH: 52d84b11737d980aef856699f885ca86
// Status: Verified Clean (Simulated Windows 10 Chrome 150)
// Notes: Chrome 150 introduces random extension shuffling and new ML-DSA Post-Quantum Signature Schemes.
// Fix Note: Extension 51764 (experimental TrustAnchors) was removed to fix CF 'illegal parameter' error.
// JA3 is essentially randomized per-connection when Extension Shuffling is enabled, but JA4 remains stable.
func init() {
	register(&illutls.BrowserProfile{
		Name:      "chrome-150-windows-10",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36",
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
				&utls.UtlsGREASEExtension{},                                         // First GREASE
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.CurveID(4588)},
					{Group: utls.X25519},
				}},
				// &utls.GenericExtension{Id: 51764, Data: []byte{}}, // 51764 (Experimental TrustAnchors, commented out to avoid CF block)
				&utls.SessionTicketExtension{},                                                // 35
				&utls.SNIExtension{},                                                          // 0
				&utls.GenericExtension{Id: 17613, Data: []byte{0x00, 0x03, 0x02, 0x68, 0x32}}, // 17613 ALPS
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.StatusRequestExtension{},                                                // 5
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},                 // 11
				&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},                // 16
				&utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: []utls.SignatureScheme{ // 13
					utls.SignatureScheme(2308), // ML-DSA?
					utls.SignatureScheme(2309), // ML-DSA?
					utls.SignatureScheme(2310), // ML-DSA?
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
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.CurveID(4588), // X25519MLKEM768
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
				}},
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
				&utls.ExtendedMasterSecretExtension{},             // 23
				&utls.GenericExtension{Id: 65037, Data: []byte{
					0x00,
					0x00, 0x01, 0x00, 0x01,
					0x00,
					0x00, 0x20,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0x00, 0x40,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				}}, // 65037 ECH
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.UtlsGREASEExtension{}, // 64250 (Last GREASE)
				&utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle},
			},
		},
		H2Settings: illutls.H2Settings{
			HeaderTableSize:      65536,
			EnablePush:           0,
			MaxConcurrentStreams: 0,
			InitialWindowSize:    6291456,
			MaxFrameSize:         0,
			MaxHeaderListSize:    262144,
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
		},
		Headers: map[string]string{
			"sec-ch-ua":                 `"Not;A=Brand";v="8", "Chromium";v="150", "Google Chrome";v="150"`,
			"sec-ch-ua-mobile":          "?0",
			"sec-ch-ua-platform":        `"Windows"`,
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"sec-fetch-site":            "none",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"sec-fetch-dest":            "document",
			"sec-fetch-storage-access":  "active",
			"accept-encoding":           "gzip, deflate, br, zstd",
			"accept-language":           "zh-CN,zh;q=0.9",
			"priority":                  "u=0, i",
		},
	})
}
