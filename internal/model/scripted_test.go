package model

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"
)

func scriptedExample(t *testing.T, format string) Document {
	t.Helper()
	b, e := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
	if e != nil {
		t.Fatal(e)
	}
	var d Document
	if e = json.Unmarshal(b, &d); e != nil {
		t.Fatal(e)
	}
	return d
}
func setScriptedSourceBytes(d *Document, b []byte) {
	asset := &d.Source.Assets[0]
	asset.DataBase64 = base64.StdEncoding.EncodeToString(b)
	asset.Size.Bytes = int64(len(b))
	asset.Hashes.SHA256 = fmt.Sprintf("%x", sha256.Sum256(b))
	size := fmt.Sprintf("%d bytes", len(b))
	asset.Size.Text = &size
}
func TestScriptedModelCapturedAndConstructed(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			d := scriptedExample(t, format)
			if e := d.Validate(); e != nil {
				t.Fatal(e)
			}
			n := d.FormatData.ASS
			if n == nil {
				n = d.FormatData.SSA
			}
			for i := range n.Sections {
				n.Sections[i].RawHeader = nil
			}
			for i := range n.Records {
				n.Records[i].RawLine = nil
			}
			n.Styles[0].DeclarationID = nil
			n.Events[0].DeclarationID = nil
			if e := d.Validate(); e != nil {
				t.Fatalf("constructed: %v", e)
			}
			d.Cues[0].Timing.StartMilliseconds = 1011
			d.Cues[0].Timing.DurationMilliseconds = 989
			start, span := int64(1011), int64(989)
			d.Document.MediaStartMilliseconds = &start
			d.Document.MediaSpanMilliseconds = &span
			d.Stats.MediaSpanMilliseconds = &span
			if e := d.Validate(); e != nil {
				t.Fatalf("editable common timing: %v", e)
			}
		})
	}
}
func TestScriptedModelRejectsForgedOwnershipAndPrivacy(t *testing.T) {
	cases := []struct {
		name, want string
		change     func(*Document)
	}{
		{"branch mixing", "only matching", func(d *Document) { d.FormatData.SSA = d.FormatData.ASS }},
		{"historical branch", "current scripted", func(d *Document) {
			d.SchemaVersion = "1.0.0"
			d.Schema = "https://cueson.io/schema/v1.0.0/cueson.schema.json"
		}},
		{"capability lie", "schema_only", func(d *Document) { d.FormatSupport.IngestSupported = true }},
		{"record order gap", "invalid_native_ownership", func(d *Document) { d.FormatData.ASS.Records[0].SourceOrder = 8 }},
		{"section crossover", "cross_section", func(d *Document) { d.FormatData.ASS.Records[2].SectionID = "section-events" }},
		{"declaration crossover", "declaration_reference", func(d *Document) { s := "record-event-format"; d.FormatData.ASS.Styles[0].DeclarationID = &s }},
		{"orphan dialogue", "dialogue", func(d *Document) { d.Cues = []Cue{} }},
		{"ambiguous style", "style", func(d *Document) { d.FormatData.ASS.Styles[0].Name = "Missing" }},
		{"nonfinal Text", "nonfinal", func(d *Document) {
			e := &d.FormatData.ASS.Events[0]
			e.DeclarationID = nil
			d.FormatData.ASS.Records[4].RawLine = nil
			e.Fields[0], e.Fields[9] = e.Fields[9], e.Fields[0]
		}},
		{"typed disagreement", "typed", func(d *Document) { v := 21.0; d.FormatData.ASS.Styles[0].Fields[2].TypedValue.Decimal = &v }},
		{"unicode declaration lookalike", "recognized_name", func(d *Document) {
			d.FormatData.ASS.Styles[0].Fields[2].FieldName = "FONTSİZE"
			d.FormatData.ASS.Records[1].DeclarationFields[2].FieldName = "FONTSİZE"
		}},
		{"original wrong dialect", "mismatched_native_dialect", func(d *Document) {
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			d.Source.Assets[0].DataBase64 = base64.StdEncoding.EncodeToString([]byte(strings.Replace(string(b), "v4.00+", "v4.00", 1)))
		}},
		{"edited wrong dialect", "mismatched_native_dialect", func(d *Document) { d.FormatData.ASS.Records[0].Fields[0].RawValue = "v4.00++" }},
		{"original active effect", "unsafe_active_content", func(d *Document) {
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			d.Source.Assets[0].DataBase64 = base64.StdEncoding.EncodeToString([]byte(strings.Replace(string(b), ",0,0,0,,Hello", ",0,0,0,template,Hello", 1)))
		}},
		{"derived-only edit", "projection", func(d *Document) { d.Cues[0].Payload.PlainText = "invented" }},
		{"span forgery", "projection", func(d *Document) { d.FormatData.ASS.Events[0].Spans[0].EndScalar-- }},
		{"provenance forgery", "projection", func(d *Document) { d.Cues[0].FormatData.ASS.Projection.SourceEventID = "event-other" }},
		{"unsafe edited metadata", "unsafe_source_metadata", func(d *Document) {
			d.FormatData.ASS.Records[0].Fields = []ScriptedField{{FieldName: "Video File", RawValue: "relative.mp4"}}
		}},
		{"unclassified edited metadata", "unsafe_source_metadata", func(d *Document) {
			d.FormatData.ASS.Records[0].Fields = []ScriptedField{{FieldName: "Private authoring state", RawValue: "opaque"}}
		}},
		{"unsafe original metadata", "unsafe_source_metadata", func(d *Document) {
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = append([]byte("[Aegisub Project Garbage]\nVideo File: https://example.invalid/movie.mp4\n"), b...)
			d.Source.Assets[0].DataBase64 = base64.StdEncoding.EncodeToString(b)
		}},
		{"unsafe unknown original section", "unsafe_source_metadata", func(d *Document) {
			b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
			b = append(b, []byte("[Private Authoring]\nOpaque Identity: workstation\n")...)
			d.Source.Assets[0].DataBase64 = base64.StdEncoding.EncodeToString(b)
		}},
		{"active event", "unsafe_active_content", func(d *Document) { d.FormatData.ASS.Events[0].EventType = "command" }},
		{"complexity line", "complexity_limit", func(d *Document) { d.FormatData.ASS.Events[0].Text = strings.Repeat("x", MaxScriptedLineBytes+1) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := scriptedExample(t, "ass")
			c.change(&d)
			e := d.Validate()
			if e == nil || !strings.Contains(e.Error(), c.want) {
				t.Fatalf("got %v, want %s", e, c.want)
			}
			if strings.Contains(e.Error(), "example.invalid") {
				t.Fatal("unsafe diagnostic leaked source")
			}
		})
	}
}
func TestScriptedProjectionExplicitBreaksDrawingsAndKaraoke(t *testing.T) {
	cases := []struct {
		text    string
		wrap    int
		lines   []string
		drawing bool
	}{
		{`a\Nb\N`, 0, []string{"a", "b", ""}, false}, {`\Na\N\Nb`, 0, []string{"", "a", "", "b"}, false}, {`a\nb`, 0, []string{"a b"}, false}, {`a\nb`, 2, []string{"a", "b"}, false}, {`a{\q2}\nb`, 0, []string{"a", "b"}, false}, {`a\hb`, 0, []string{"a\u00a0b"}, false}, {`{\p1}m 0 0 l 1 1{\p0}`, 0, []string{""}, true}, {`{\b1}😀text{\b0}`, 0, []string{"😀text"}, false}, {`x{broken`, 0, []string{"x{broken"}, false},
	}
	for _, c := range cases {
		f, e := ProjectScriptedText(c.text, c.wrap, 1000, 2000)
		if e != nil || !slices.Equal(f.Lines, c.lines) || f.DrawingExcluded != c.drawing {
			t.Errorf("%q: %#v %v", c.text, f, e)
		}
	}
	f, e := ProjectScriptedText(`{\k20}one{\K30}two`, 0, 1000, 2000)
	if e != nil || len(f.Tokens) != 2 || f.Tokens[0].EndMilliseconds != 1200 || f.Tokens[1].StartMilliseconds != 1200 || f.Tokens[1].EndMilliseconds != 1500 || f.Tokens[1].Text != "two" {
		t.Fatalf("karaoke: %#v %v", f, e)
	}
	for _, s := range []string{`{\kt20}one`, `{\k999}one`, `{\k-1}one`} {
		f, e = ProjectScriptedText(s, 0, 1000, 2000)
		if e != nil || len(f.Tokens) != 0 || !slices.Contains(f.DiagnosticCodes, "unsupported_karaoke_timing") {
			t.Errorf("unsupported %q: %#v %v", s, f, e)
		}
	}
	for _, s := range []string{`{\p-1}x`, `{\q9}x`} {
		if _, e = ProjectScriptedText(s, 0, 1000, 2000); e == nil {
			t.Errorf("ambiguous state accepted: %q", s)
		}
	}
}
func TestScriptedZeroDialogueSummaries(t *testing.T) {
	d := scriptedExample(t, "ass")
	n := d.FormatData.ASS
	n.Records = n.Records[:4]
	n.Events = []ScriptedEvent{}
	d.Cues = []Cue{}
	d.Document = DocumentSummary{}
	d.Stats = Stats{}
	if e := d.Validate(); e != nil {
		t.Fatal(e)
	}
}

func TestScriptedUnknownFieldsAliasesAndContentRoles(t *testing.T) {
	d := scriptedExample(t, "ass")
	n := d.FormatData.ASS
	extra := ScriptedField{FieldName: "X-Native", RawValue: "retained opaque content"}
	n.Styles[0].Fields = append(n.Styles[0].Fields, extra)
	n.Records[1].DeclarationFields = append(n.Records[1].DeclarationFields, ScriptedDeclarationField{FieldName: extra.FieldName})
	n.Styles[0].Fields[0], n.Styles[0].Fields[1] = n.Styles[0].Fields[1], n.Styles[0].Fields[0]
	n.Records[1].DeclarationFields[0], n.Records[1].DeclarationFields[1] = n.Records[1].DeclarationFields[1], n.Records[1].DeclarationFields[0]
	n.Events[0].Fields[4].FieldName = "Actor"
	n.Records[3].DeclarationFields[4].FieldName = "Actor"
	actor := "Actor"
	d.Cues[0].FormatData.ASS.Projection.SpeakerFieldName = &actor
	b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
	b = append(b, []byte("; content-role discussion of C:\\examples\\captions.ass\n")...)
	setScriptedSourceBytes(&d, b)
	if e := d.Validate(); e != nil {
		t.Fatalf("unknown/reordered native fields, Actor alias, inert content comment: %v", e)
	}
}

func TestScriptedNonDialogueOffsetsAndProjectionCannotBeForged(t *testing.T) {
	for _, valid := range []bool{true, false} {
		for _, observation := range []string{"span", "tag", "karaoke"} {
			t.Run(fmt.Sprintf("valid-%t-%s", valid, observation), func(t *testing.T) {
				d := scriptedExample(t, "ass")
				event := &d.FormatData.ASS.Events[0]
				event.EventType = "comment"
				event.CueID = nil
				event.Valid = valid
				d.Cues = []Cue{}
				d.Document = DocumentSummary{}
				d.Stats = Stats{}
				if !valid {
					order := 7
					d.Diagnostics = []Diagnostic{{Severity: "warning", Code: "malformed_native_record", Message: "Malformed non-dialogue retained.", SourceOrder: &order}}
					d.Stats.DiagnosticCount = 1
					d.Stats.WarningCount = 1
				}
				switch observation {
				case "span":
					event.Spans[0].EndScalar++
				case "tag":
					event.Tags = []ScriptedTag{{Name: "b", Parameter: "1", StartScalar: 0, EndScalar: 3, Raw: `\b1`}}
				case "karaoke":
					event.Karaoke = []ScriptedKaraoke{{Variant: "k", StartScalar: 0, EndScalar: 3, Supported: true}}
				}
				if e := d.Validate(); e == nil || !strings.Contains(e.Error(), "inconsistent_projection") {
					t.Fatalf("forged non-dialogue observation accepted: %v", e)
				}
			})
		}
	}
}
