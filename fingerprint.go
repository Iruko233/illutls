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
	// Deep-copy extensions (shallow element copy 鈥?extensions are typically
	// stateless or have internal cloning during ApplyPreset).
	dst.Extensions = make([]utls.TLSExtension, len(src.Extensions))
	copy(dst.Extensions, src.Extensions)
	return dst
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
