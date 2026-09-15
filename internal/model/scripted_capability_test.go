package model

import "testing"

func TestScriptedCapabilityObservationsRemainTruthful(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, support := range []FormatSupport{
			{Status: "schema_only", RestoreSupported: true},
			{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
			{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
		} {
			d := scriptedExample(t, format)
			d.FormatSupport = support
			if err := d.Validate(); err != nil {
				t.Errorf("%s %s valid observation rejected: %v", format, support.Status, err)
			}
		}
		for _, support := range []FormatSupport{
			{Status: "schema_only", IngestSupported: true, RestoreSupported: true},
			{Status: "experimental", IngestSupported: true, RestoreSupported: true},
			{Status: "experimental", RenderSupported: true, RestoreSupported: true},
			{Status: "experimental", IngestSupported: true, RenderSupported: true},
			{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true, OCRRequiredForSemanticOutput: true},
			{Status: "stable", IngestSupported: true, RestoreSupported: true},
		} {
			d := scriptedExample(t, format)
			d.FormatSupport = support
			if err := d.Validate(); err == nil {
				t.Errorf("%s accepted invalid or premature capability: %#v", format, support)
			}
		}
	}
}
