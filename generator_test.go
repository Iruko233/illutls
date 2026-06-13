package illutls

import (
	"fmt"
	"testing"
)

func TestGenerateProfileDeterministic(t *testing.T) {
	// Generate two profiles with the same seed, platform, and version
	p1 := GenerateProfile(12345, "windows", 122, false)
	p2 := GenerateProfile(12345, "windows", 122, false)

	// User-Agent and Headers should be identical
	if p1.UserAgent != p2.UserAgent {
		t.Errorf("UserAgent mismatch for same seed")
	}

	if p1.Headers["sec-ch-ua"] != p2.Headers["sec-ch-ua"] {
		t.Errorf("sec-ch-ua mismatch for same seed: %v vs %v", p1.Headers["sec-ch-ua"], p2.Headers["sec-ch-ua"])
	}

	// Extensions count and order should be identical
	if len(p1.TLSSpec.Extensions) != len(p2.TLSSpec.Extensions) {
		t.Errorf("TLS Extension count mismatch")
	}

	for i := range p1.TLSSpec.Extensions {
		// Type comparison
		if fmt.Sprintf("%T", p1.TLSSpec.Extensions[i]) != fmt.Sprintf("%T", p2.TLSSpec.Extensions[i]) {
			t.Errorf("Extension order mismatch at index %d", i)
		}
	}

	// Different seed should produce a different extension order or different GREASE brand
	p3 := GenerateProfile(54321, "windows", 122, false)
	if p1.Headers["sec-ch-ua"] == p3.Headers["sec-ch-ua"] {
		// Note: There's a 1/6 chance the order matches, but the GREASE char will likely differ
		t.Logf("Note: p1 and p3 generated identical sec-ch-ua, check if expected based on seed mod")
	}
}
