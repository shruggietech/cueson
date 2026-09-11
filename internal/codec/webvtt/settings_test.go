package webvtt

import "testing"

func TestCueSettingsRetainOccurrencesAndLastValidValue(t *testing.T) {
	occurrences, settings, diagnostics := parseCueSettings("line:10% mystery:x line:20%,center size:101%", 3, "cue-000000")
	if len(occurrences) != 4 || len(diagnostics) != 3 {
		t.Fatalf("occurrences=%#v diagnostics=%#v", occurrences, diagnostics)
	}
	if got := settings["line"]; got != "20%,center" {
		t.Fatalf("effective line = %q", got)
	}
	if occurrences[1].Recognized || occurrences[1].Valid || !occurrences[2].Recognized || !occurrences[2].Valid || occurrences[3].Valid {
		t.Fatalf("occurrence flags = %#v", occurrences)
	}
}

func TestRegionSettingsValidateCompleteKnownSurface(t *testing.T) {
	raw := "id:captions width:80% lines:3 regionanchor:0%,100% viewportanchor:10%,90% scroll:up"
	occurrences, settings, diagnostics := parseRegionSettings(raw, 0)
	if len(occurrences) != 6 || len(settings) != 6 || len(diagnostics) != 0 {
		t.Fatalf("occurrences=%#v settings=%#v diagnostics=%#v", occurrences, settings, diagnostics)
	}
}

func TestSettingsUseOnlyASCIIWhitespaceAndStrictPercentages(t *testing.T) {
	occurrences, settings, diagnostics := parseCueSettings("size:10%\u00a0align:start position:+2% line:1.5", 0, "cue-000000")
	if len(occurrences) != 3 || len(settings) != 0 || len(diagnostics) != 3 {
		t.Fatalf("occurrences=%#v settings=%#v diagnostics=%#v", occurrences, settings, diagnostics)
	}
}

func TestSettingsPreserveNULLexicallyAndReplaceItInEffectiveSemantics(t *testing.T) {
	occurrences, settings, diagnostics := parseRegionSettings("id:r\x00name", 0)
	if len(occurrences) != 1 || occurrences[0].Value != "r\x00name" || settings["id"] != "r\ufffdname" || len(diagnostics) != 0 {
		t.Fatalf("occurrences=%#v settings=%#v diagnostics=%#v", occurrences, settings, diagnostics)
	}
}
