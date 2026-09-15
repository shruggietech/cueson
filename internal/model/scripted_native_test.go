package model

import "testing"

func TestScriptedNativeAdapterAuthority(t *testing.T) {
	fields := ScriptedCanonicalFields("ass", "style")
	fields[0] = "changed"
	if ScriptedCanonicalFields("ass", "style")[0] != "Name" {
		t.Fatal("mutable canonical authority")
	}
	if n, known := ScriptedFieldIdentity(" Actor ", "ass", "event"); n != "name" || !known {
		t.Fatal("Actor event alias absent")
	}
	if _, known := ScriptedFieldIdentity("Actor", "ass", "style"); known {
		t.Fatal("Actor alias leaked into style")
	}
	for _, tc := range []struct{ name, raw, kind string }{{"Bold", "-1", "boolean"}, {"Fontsize", "20.5", "decimal"}, {"Alignment", "2", "integer"}, {"PrimaryColour", "-2147483648", "color"}, {"Fontname", "A Font", "string"}} {
		v, err := ScriptedScalarValue(ScriptedField{FieldName: tc.name, RawValue: tc.raw}, "ssa", "style", 0)
		if err != nil || v.Kind != tc.kind {
			t.Fatalf("%s: %+v %v", tc.name, v, err)
		}
	}
	if _, err := ScriptedScalarValue(ScriptedField{FieldName: "PrimaryColour", RawValue: "-2147483649"}, "ssa", "style", 0); err == nil {
		t.Fatal("signed color overflow accepted")
	}
	if _, err := ScriptedMilliseconds("9223372036854775807:00:00.00"); err == nil {
		t.Fatal("timestamp overflow accepted")
	}
	for _, tc := range []struct {
		line string
		n    int
		bad  bool
	}{{"!!!!", 3, false}, {"!!", 1, false}, {"!\"", 1, true}, {"!!!", 2, false}, {"!", 0, true}, {"z", 0, true}} {
		n, bad, err := ScriptedAttachmentFacts([]string{tc.line})
		if err != nil || n != tc.n || bad != tc.bad {
			t.Fatalf("%q: %d %v %v", tc.line, n, bad, err)
		}
	}
}
