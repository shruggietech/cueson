package model

import (
	"strings"
	"testing"
)

func TestScriptedCaptureIndependentComplexityGuards(t *testing.T) {
	t.Run("encoded source preflight", func(t *testing.T) {
		// Invalid base64 proves size refusal precedes decoding and its allocation.
		d := scriptedExample(t, "ass")
		d.Source.Assets[0].DataBase64 = strings.Repeat("!", ((64<<20)+2)/3*4+1)
		if _, err := inspectScriptedCaptureSource(d.Source); err == nil || !strings.Contains(err.Error(), "complexity_limit: original capture source") {
			t.Fatalf("independent encoded size preflight: %v", err)
		}
	})
	for _, count := range []int{MaxScriptedDeclarationFields, MaxScriptedDeclarationFields + 1, MaxScriptedLineBytes / 2} {
		d := scriptedExample(t, "ass")
		raw := []byte("[Events]\nFormat: " + strings.Repeat("X,", count-1) + "Text\n")
		setScriptedSourceBytes(&d, raw)
		capture, err := inspectScriptedCaptureSource(d.Source)
		if count == MaxScriptedDeclarationFields {
			if err != nil || !capture.matches(1, string(raw[len("[Events]\n"):len(raw)-1])) {
				t.Fatalf("exact declaration capture limit: %v", err)
			}
		} else if err == nil || !strings.Contains(err.Error(), "complexity_limit: original declaration") {
			t.Fatalf("%d-field declaration refusal: %v", count, err)
		}
	}
}
