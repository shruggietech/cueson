package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/shruggietech/cueson/internal/model"
	"os"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/source"
)

func TestScriptedCurrentAnnotatedExamples(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			b, e := os.ReadFile("testdata/scripted-" + format + ".cueson.json")
			if e != nil {
				t.Fatal(e)
			}
			d, e := Decode(b)
			if e != nil {
				t.Fatal(e)
			}
			if e = source.ValidateIntegrity(context.Background(), d); e != nil {
				t.Fatal(e)
			}
			again, e := json.Marshal(d)
			if e != nil {
				t.Fatal(e)
			}
			if e = Validate(again); e != nil {
				t.Fatal(e)
			}
			var before, after any
			if e = json.Unmarshal(b, &before); e != nil {
				t.Fatal(e)
			}
			if e = json.Unmarshal(again, &after); e != nil {
				t.Fatal(e)
			}
			canonicalBefore, _ := json.Marshal(before)
			canonicalAfter, _ := json.Marshal(after)
			if !bytes.Equal(canonicalBefore, canonicalAfter) {
				t.Fatal("typed representation dropped accepted scripted fields")
			}
		})
	}
}

func TestScriptedDrawingKaraokeAndMalformedSchemaModels(t *testing.T) {
	for _, text := range []string{`{\p1}m 0 0 l 1 1{\p0}`, `{\k20}one{\K30}two`, `literal{broken`, `a\N\Ntail, `} {
		b, e := os.ReadFile("testdata/scripted-ass.cueson.json")
		if e != nil {
			t.Fatal(e)
		}
		d, e := Decode(b)
		if e != nil {
			t.Fatal(e)
		}
		facts, e := model.ProjectScriptedText(text, 0, 1000, 2000)
		if e != nil {
			t.Fatal(e)
		}
		event := &d.FormatData.ASS.Events[0]
		event.Text = text
		event.Fields[9].RawValue = text
		event.Spans = facts.Spans
		event.Tags = facts.Tags
		event.Karaoke = facts.Karaoke
		cue := &d.Cues[0]
		cue.Payload = model.Payload{RawText: text, PlainText: strings.Join(facts.Lines, "\n"), Lines: facts.Lines}
		cue.Tokens = facts.Tokens
		cue.FormatData.ASS.Projection.DrawingExcluded = facts.DrawingExcluded
		d.Document.HasWordLevelTiming = len(facts.Tokens) > 0
		d.Stats.HasWordLevelTiming = d.Document.HasWordLevelTiming
		for _, code := range facts.DiagnosticCodes {
			order := 7
			d.Diagnostics = append(d.Diagnostics, model.Diagnostic{Severity: "warning", Code: code, Message: "Native ambiguity retained.", SourceOrder: &order})
		}
		d.Stats.DiagnosticCount = len(d.Diagnostics)
		d.Stats.WarningCount = len(d.Diagnostics)
		encoded, e := json.Marshal(d)
		if e != nil {
			t.Fatal(e)
		}
		if e = Validate(encoded); e != nil {
			t.Fatalf("%q: %v", text, e)
		}
	}
}
func TestScriptedSchemaRejectsUnknownAndMixedBranches(t *testing.T) {
	b, e := os.ReadFile("testdata/scripted-ass.cueson.json")
	if e != nil {
		t.Fatal(e)
	}
	for _, mutate := range []func(map[string]any){func(d map[string]any) { d["extra"] = true }, func(d map[string]any) { n := d["format_data"].(map[string]any); n["ssa"] = n["ass"] }, func(d map[string]any) {
		n := d["format_data"].(map[string]any)["ass"].(map[string]any)
		n["unratified_map"] = map[string]any{}
	}, func(d map[string]any) { d["format_support"].(map[string]any)["status"] = "stable" }} {
		var d map[string]any
		if e = json.Unmarshal(b, &d); e != nil {
			t.Fatal(e)
		}
		mutate(d)
		bad, _ := json.Marshal(d)
		if e = Validate(bad); e == nil {
			t.Fatal("unknown field/branch/capability accepted")
		}
	}
}
