package profiles

import (
	"github.com/Iruko233/illutls"
	utls "github.com/refraction-networking/utls"
)

// JA3: 771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-45-23-18-5-43-17613-10-16-13-35-27-11-51-65037-65281,4588-29-23-24,0
// JA3 HASH: 93b91180f7396c2af5c9c738f8493fd5
// JA4: t13d1516h2_8daaf6152771_d8a2da3f94cd
// H2: 1:65536;2:0;4:6291456;6:262144|15663105|0|m,a,s,p
// H2 HASH: 52d84b11737d980aef856699f885ca86
// Status: Verified Clean (Simulated Android 12 WeChat 8.0)
// Notes: WeChat WebView based on Chrome 146. Features randomized extension shuffling, ALPS (17613), ECH, and ML-KEM768. Includes x-requested-with: com.tencent.mm.
func init() {
	register(&illutls.BrowserProfile{
		Name:      "wechat-8-android-12",
		UserAgent: "Mozilla/5.0 (Linux; Android 12; ANG-AN00 Build/HUAWEIANG-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/146.0.7680.178 Mobile Safari/537.36 XWEB/1460149 MMWEBSDK/20260502 MMWEBID/5226 REV/4843ebe11ecd759196310a3df45633a796c09f5b MicroMessenger/8.0.72.3100(0x28004850) WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64",
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
				&utls.PSKKeyExchangeModesExtension{Modes: []uint8{utls.PskModeDHE}}, // 45
				&utls.ExtendedMasterSecretExtension{},                               // 23
				&utls.SCTExtension{},                                                // 18
				&utls.StatusRequestExtension{},                                      // 5
				&utls.SupportedVersionsExtension{Versions: []uint16{ // 43
					utls.GREASE_PLACEHOLDER,
					utls.VersionTLS13,
					utls.VersionTLS12,
				}},
				&utls.GenericExtension{Id: 17613, Data: []byte{0x00, 0x03, 0x02, 0x68, 0x32}}, // 17613 ALPS
				&utls.SupportedCurvesExtension{Curves: []utls.CurveID{ // 10
					utls.GREASE_PLACEHOLDER,
					utls.CurveID(4588), // X25519MLKEM768
					utls.X25519,
					utls.CurveP256,
					utls.CurveP384,
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
				&utls.SessionTicketExtension{}, // 35
				&utls.UtlsCompressCertExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli}}, // 27
				&utls.SupportedPointsExtension{SupportedPoints: []byte{0x00}},                                       // 11
				&utls.KeyShareExtension{KeyShares: []utls.KeyShare{ // 51
					{Group: utls.GREASE_PLACEHOLDER, Data: []byte{0}},
					{Group: utls.CurveID(4588)}, // 4588
					{Group: utls.X25519},        // 29
				}},
				&utls.GenericExtension{Id: 65037, Data: []byte{
					0x00,
					0x00, 0x01, 0x00, 0x01,
					0x00,
					0x00, 0x20,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0x00, 0x40,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				}},                             // 65037 ECH
				&utls.RenegotiationInfoExtension{Renegotiation: utls.RenegotiateOnceAsClient}, // 65281
				&utls.UtlsGREASEExtension{},                                                   // Last GREASE
			},
		},
		H2Settings: illutls.H2Settings{
			HeaderTableSize:   65536,
			EnablePush:        0,
			InitialWindowSize: 6291456,
			MaxHeaderListSize: 262144,
			SettingsOrder:     []uint16{1, 2, 4, 6},
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
			"x-requested-with",
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
			"sec-ch-ua":                 `"Chromium";v="146", "Not-A.Brand";v="24", "Android WebView";v="146"`,
			"sec-ch-ua-mobile":          "?1",
			"sec-ch-ua-platform":        `"Android"`,
			"upgrade-insecure-requests": "1",
			"user-agent":                "Mozilla/5.0 (Linux; Android 12; ANG-AN00 Build/HUAWEIANG-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/146.0.7680.178 Mobile Safari/537.36 XWEB/1460149 MMWEBSDK/20260502 MMWEBID/5226 REV/4843ebe11ecd759196310a3df45633a796c09f5b MicroMessenger/8.0.72.3100(0x28004850) WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64",
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/wxpic,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"x-requested-with":          "com.tencent.mm",
			"sec-fetch-site":            "none",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"sec-fetch-dest":            "document",
			"sec-fetch-storage-access":  "active",
			"accept-encoding":           "gzip, deflate, br, zstd",
			"accept-language":           "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
			"priority":                  "u=0, i",
		},
	})
}
