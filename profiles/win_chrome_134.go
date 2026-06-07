package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// Target JA3: 340b5b91cdefd018d1ab75fb30fbdd43
// Target H2: 1:65536;2:0;4:6291456;6:262144|15663105|0|m,a,s,p
// Status: Verified Clean (No ECH, No PSK)
func init() {
	register(&illutls.BrowserProfile{
		Name:      "Win Chrome 134",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
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
				&utls.UtlsGREASEExtension{}, // 0xeaea
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.X25519MLKEM768,
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
				}},
				&utls.SessionTicketExtension{}, // 35
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}},           // 45
				&utls.ExtendedMasterSecretExtension{},                                         // 23
				&utls.SCTExtension{},                                                          // 18
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.GenericExtension{Id: 17513, Data: []byte{0x00, 0x03, 0x02, 0x68, 0x32}}, // 17513 (ALPS)
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},                 // 11
				&utls.StatusRequestExtension{},                                                // 5
				&utls.SNIExtension{},                                                          // 0
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.X25519MLKEM768},
					{Group: utls.X25519},
				}},
				&utls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}}, // 16
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
				&utls.UtlsGREASEExtension{}, // 0x3a3a
			},
		},
		H2Settings: illutls.H2Settings{
			HeaderTableSize:   65536,
			EnablePush:        0,
			InitialWindowSize: 6291456,
			MaxHeaderListSize: 262144,
		},
		H2WindowUpdate: 15663105,
		H2Priority: illutls.H2Priority{
			Weight:    220,
			DependsOn: 0,
			Exclusive: true,
		},
		PHeaderOrder: []string{
			":method",
			":authority",
			":scheme",
			":path",
		},
		HeaderOrder: []string{
			"sec-ch-ua-platform",
			"user-agent",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"sec-fetch-dest",
			"accept-encoding",
			"accept-language",
			"priority",
		},
		Headers: map[string]string{
			"Sec-Ch-Ua-Platform": `"Windows"`,
			"Sec-Ch-Ua":          `"Chromium";v="134", "Not:A-Brand";v="24", "Google Chrome";v="134"`,
			"Sec-Ch-Ua-Mobile":   "?0",
			"Accept":             "*/*",
			"Sec-Fetch-Site":     "same-site",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Dest":     "empty",
			"Accept-Encoding":    "gzip, deflate, br, zstd",
			"Accept-Language":    "zh-CN,zh;q=0.9",
			"Priority":           "u=1, i",
		},
	})
}
