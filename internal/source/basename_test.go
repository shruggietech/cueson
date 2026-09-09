package source

import "testing"

func TestValidateSafeBasename(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"captions.srt", "résumé.vtt", "CAPTION_01.bin"} {
		if err := validateSafeBasename(name); err != nil {
			t.Errorf("validateSafeBasename(%q) error = %v", name, err)
		}
	}
	for _, name := range []string{"", ".", "..", "../x", `a\b`, "a:b", "x. ", "x.", "CON", "con.txt", "Lpt9.log", "bad\x00name"} {
		if err := validateSafeBasename(name); err == nil {
			t.Errorf("validateSafeBasename(%q) error = nil", name)
		}
	}
}

func TestPortableIdentityUsesCanonicalCaselessMatching(t *testing.T) {
	t.Parallel()

	if portableIdentity("RÉSUMÉ.SRT") != portableIdentity("re\u0301sume\u0301.srt") {
		t.Error("canonical-equivalent case variants do not collide")
	}
	document := testDocument(t)
	companion := document.Source.Assets[0]
	companion.ID = "asset-1"
	companion.Role = "companion"
	companion.FileName = "CAPTIONS.SRT"
	document.Source.Assets = append(document.Source.Assets, companion)
	if _, err := prepareAssets(document); err == nil {
		t.Fatal("prepareAssets() accepted portable basename collision")
	}
}
