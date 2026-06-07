package illutls

import (
	utls "github.com/refraction-networking/utls"
	"math/rand"
)

// greasePlaceholders are valid GREASE values per RFC 8701.
var greasePlaceholders = []uint16{
	0x0a0a, 0x1a1a, 0x2a2a, 0x3a3a, 0x4a4a, 0x5a5a,
	0x6a6a, 0x7a7a, 0x8a8a, 0x9a9a, 0xaaaa, 0xbaba,
	0xcaca, 0xdada, 0xeaea, 0xfafa,
}

// randomGREASE returns a random GREASE value.
func randomGREASE() uint16 {
	return greasePlaceholders[rand.Intn(len(greasePlaceholders))]
}

// CloneClientHelloSpec deep-copies a ClientHelloSpec so that GREASE
// randomization on one connection does not affect others.
func CloneClientHelloSpec(src *utls.ClientHelloSpec) *utls.ClientHelloSpec {
	if src == nil {
		return nil
	}
	dst := &utls.ClientHelloSpec{
		TLSVersMin:         src.TLSVersMin,
		TLSVersMax:         src.TLSVersMax,
		CompressionMethods: append([]byte(nil), src.CompressionMethods...),
		GetSessionID:       src.GetSessionID,
	}
	// Deep-copy cipher suites.
	dst.CipherSuites = make([]uint16, len(src.CipherSuites))
	copy(dst.CipherSuites, src.CipherSuites)
	
	// Deep-copy extensions.
	dst.Extensions = make([]utls.TLSExtension, len(src.Extensions))
	for i, ext := range src.Extensions {
		dst.Extensions[i] = cloneExtension(ext)
	}
	return dst
}

func cloneExtension(e utls.TLSExtension) utls.TLSExtension {
	switch ext := e.(type) {
	case *utls.SNIExtension:
		return &utls.SNIExtension{}
	case *utls.KeyShareExtension:
		ks := make([]utls.KeyShare, len(ext.KeyShares))
		for i, k := range ext.KeyShares {
			ks[i] = utls.KeyShare{Group: k.Group, Data: append([]byte(nil), k.Data...)}
		}
		return &utls.KeyShareExtension{KeyShares: ks}
	case *utls.SupportedCurvesExtension:
		curves := make([]utls.CurveID, len(ext.Curves))
		copy(curves, ext.Curves)
		return &utls.SupportedCurvesExtension{Curves: curves}
	case *utls.SupportedPointsExtension:
		points := make([]byte, len(ext.SupportedPoints))
		copy(points, ext.SupportedPoints)
		return &utls.SupportedPointsExtension{SupportedPoints: points}
	case *utls.SessionTicketExtension:
		return &utls.SessionTicketExtension{}
	case *utls.ALPNExtension:
		alpn := make([]string, len(ext.AlpnProtocols))
		copy(alpn, ext.AlpnProtocols)
		return &utls.ALPNExtension{AlpnProtocols: alpn}
	case *utls.SignatureAlgorithmsExtension:
		algs := make([]utls.SignatureScheme, len(ext.SupportedSignatureAlgorithms))
		copy(algs, ext.SupportedSignatureAlgorithms)
		return &utls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: algs}
	case *utls.SupportedVersionsExtension:
		versions := make([]uint16, len(ext.Versions))
		copy(versions, ext.Versions)
		return &utls.SupportedVersionsExtension{Versions: versions}
	case *utls.UtlsGREASEExtension:
		return &utls.UtlsGREASEExtension{Value: ext.Value, Body: append([]byte(nil), ext.Body...)}
	case *utls.UtlsCompressCertExtension:
		algos := make([]utls.CertCompressionAlgo, len(ext.Algorithms))
		copy(algos, ext.Algorithms)
		return &utls.UtlsCompressCertExtension{Algorithms: algos}
	case *utls.RenegotiationInfoExtension:
		return &utls.RenegotiationInfoExtension{Renegotiation: ext.Renegotiation}
	case *utls.PSKKeyExchangeModesExtension:
		modes := make([]uint8, len(ext.Modes))
		copy(modes, ext.Modes)
		return &utls.PSKKeyExchangeModesExtension{Modes: modes}
	case *utls.ExtendedMasterSecretExtension:
		return &utls.ExtendedMasterSecretExtension{}
	case *utls.SCTExtension:
		return &utls.SCTExtension{}
	case *utls.StatusRequestExtension:
		return &utls.StatusRequestExtension{}
	case *utls.GenericExtension:
		return &utls.GenericExtension{Id: ext.Id, Data: append([]byte(nil), ext.Data...)}
	case *utls.UtlsPaddingExtension:
		return &utls.UtlsPaddingExtension{GetPaddingLen: ext.GetPaddingLen, WillPad: ext.WillPad}
	case *utls.ApplicationSettingsExtension:
		alpn := make([]string, len(ext.SupportedProtocols))
		copy(alpn, ext.SupportedProtocols)
		return &utls.ApplicationSettingsExtension{SupportedProtocols: alpn}
	}
	return e // Fallback
}

// RandomizeGREASE replaces all GREASE placeholder cipher suites in the spec
// with fresh random values. This ensures each connection looks slightly
// different while maintaining the correct GREASE pattern.
func RandomizeGREASE(spec *utls.ClientHelloSpec) {
	if spec == nil {
		return
	}
	// Pick two independent GREASE values for variety.
	grease1 := randomGREASE()
	grease2 := randomGREASE()
	for grease2 == grease1 {
		grease2 = randomGREASE()
	}
	greaseIdx := 0
	for i, cs := range spec.CipherSuites {
		if isGREASE(cs) {
			if greaseIdx%2 == 0 {
				spec.CipherSuites[i] = grease1
			} else {
				spec.CipherSuites[i] = grease2
			}
			greaseIdx++
		}
	}
}

// isGREASE reports whether v is a GREASE value (0x?a?a pattern).
func isGREASE(v uint16) bool {
	return (v & 0x0f0f) == 0x0a0a
}
