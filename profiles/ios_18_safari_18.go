package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// Target JA3: ecdf4f49dd59effc439639da29186671
// Target H2: 2:0;3:100;4:2097152;9:1|10420225|0|m,s,a,p
// Status: Verified Clean (No ECH, No PSK)
func init() {
	register(&illutls.BrowserProfile{
		Name:      "iOS 18 Safari 18",
		UserAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 18_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Mobile/15E148 Safari/604.1",
		TLSSpec: &utls.ClientHelloSpec{
			TLSVersMin: utls.VersionTLS12,
			TLSVersMax: utls.VersionTLS13,
			CipherSuites: []uint16{
				utls.GREASE_PLACEHOLDER,
				utls.TLS_AES_256_GCM_SHA384,
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
				0xc008, // TLS_ECDHE_ECDSA_WITH_3DES_EDE_CBC_SHA (49160)
				0xc012, // TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA (49170)
				0x000a, // TLS_RSA_WITH_3DES_EDE_CBC_SHA (10)
			},
			CompressionMethods: []byte{0x00},
			Extensions: []utls.TLSExtension{
				&utls.UtlsGREASEExtension{},           // 0x3a3a
				&utls.SNIExtension{},                  // 0
				&utls.ExtendedMasterSecretExtension{}, // 23
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.X25519MLKEM768,
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
					utls.PSSWithSHA384,
					utls.PSSWithSHA384, // Safari 18 known quirk: sends this twice
					utls.PKCS1WithSHA384,
					utls.PSSWithSHA512,
					utls.PKCS1WithSHA512,
					utls.PKCS1WithSHA1,
				}},
				&utls.SCTExtension{}, // 18
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.X25519MLKEM768},
					{Group: utls.X25519},
				}},
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.GenericExtension{Id: 27, Data: []byte{0x02, 0x00, 0x01}}, // 27 (compress_certificate with zlib)
				&utls.UtlsGREASEExtension{},                                    // 0x6a6a
			},
		},
		H2Settings: illutls.H2Settings{
			// Safari iOS omits HeaderTableSize and MaxHeaderListSize entirely, so we leave them 0
			EnablePush:           0,
			MaxConcurrentStreams: 100,
			InitialWindowSize:    2097152,
			NoRFC7540Priorities:  1,
		},
		H2WindowUpdate: 10420225,
		// PRIORITY frame is also sent according to raw JSON: priority: u=3, i
		// According to HTTP3 Ext Priority, this is u=3, i. But does Safari send HTTP/2 PRIORITY frame?
		// According to JA4, Safari iOS sends PRIORITY frame in HTTP2, but wait!
		// The json says priority header in the headers array. Let's see:
		// "priority: u=3, i"
		// If it's sent as a header, it's not a PRIORITY frame on stream 1.
		// Wait, the JSON says:
		// "priority": { "weight": 256, "depends_on": 0, "exclusive": 1 }
		// Wait! The JSON for Safari DOES NOT HAVE the "priority" object in sent_frames!
		// Let's re-read the JSON for Safari iOS ecdf4f49dd59effc439639da29186671.json.
		// "flags": ["EndStream (0x1)", "EndHeaders (0x4)"] -> It does NOT have Priority (0x20) flag!
		// So Safari does NOT send HTTP/2 PRIORITY frame! It only sends Ext-Priority header!
		// So I should not set H2Priority to anything. fhttp won't send PRIORITY if Weight == 0 and DependsOn == 0 and Exclusive == false and Priority param is not set. Wait, fhttp does send priority by default if not told otherwise?
		// In bogdanfinn/fhttp, StreamPriority is used. We will leave H2Priority empty to not set Priority flag.
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
			"referer",
			"sec-fetch-site",
			"sec-fetch-mode",
			"accept-language",
			"priority",
			"accept-encoding",
		},
		Headers: map[string]string{
			"Sec-Fetch-Dest":  "empty",
			"User-Agent":      "Mozilla/5.0 (iPhone; CPU iPhone OS 18_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Mobile/15E148 Safari/604.1",
			"Accept":          "*/*",
			"Sec-Fetch-Site":  "same-origin",
			"Sec-Fetch-Mode":  "cors",
			"Accept-Language": "zh-CN,zh-Hans;q=0.9",
			"Priority":        "u=3, i",
			"Accept-Encoding": "gzip, deflate, br",
		},
	})
}
