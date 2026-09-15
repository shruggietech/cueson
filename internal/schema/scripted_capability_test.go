package schema

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestScriptedSchemaCapabilityAlternatives(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		b, err := os.ReadFile("testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			t.Fatal(err)
		}
		d, err := Decode(b)
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range []struct {
			support model.FormatSupport
			valid   bool
		}{
			{model.FormatSupport{Status: "schema_only", RestoreSupported: true}, true},
			{model.FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true}, true},
			{model.FormatSupport{Status: "experimental", IngestSupported: true, RestoreSupported: true}, false},
			{model.FormatSupport{Status: "schema_only", IngestSupported: true, RenderSupported: true, RestoreSupported: true}, false},
			{model.FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}, false},
			{model.FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true, OCRRequiredForSemanticOutput: true}, false},
		} {
			d.FormatSupport = test.support
			payload, err := json.Marshal(d)
			if err != nil {
				t.Fatal(err)
			}
			if err := Validate(payload); (err == nil) != test.valid {
				t.Errorf("%s %#v valid=%v: %v", format, test.support, test.valid, err)
			}
		}
	}
}
