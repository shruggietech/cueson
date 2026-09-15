package scripted

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

// FuzzRenderParseCycle checks native semantic retention and canonical stability.
// Bounded fuzz inputs supplement the separate large-model/output regressions.
func FuzzRenderParseCycle(f *testing.F) {
	fixtures := map[string][]byte{}
	for _, format := range []string{"ass", "ssa"} {
		encoded, err := os.ReadFile("../../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			f.Fatal(err)
		}
		fixtures[format] = encoded
		var d model.Document
		if err = json.Unmarshal(encoded, &d); err != nil {
			f.Fatal(err)
		}
		raw, err := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(raw)
		f.Add(append(append([]byte{}, raw...), []byte("Inert extension\n")...))
		f.Add(append([]byte{0xef, 0xbb, 0xbf}, bytes.ReplaceAll(raw, []byte{'\n'}, []byte{'\r', '\n'})...))
		// Seed declaration ownership, contextual override grammar, and opaque
		// attachments independently so mutations reach their owning boundaries.
		first := "Layer"
		if format == "ssa" {
			first = "Marked"
		}
		reordered := bytes.Replace(raw, []byte("Format: "+first+", Start"), []byte("Format: Start, "+first), 1)
		eventFirst := "0"
		if format == "ssa" {
			eventFirst = "Marked=0"
		}
		reordered = bytes.Replace(reordered, []byte("Dialogue: "+eventFirst+",0:00:01.00"), []byte("Dialogue: 0:00:01.00,"+eventFirst), 1)
		f.Add(reordered)
		f.Add(bytes.ReplaceAll(raw, []byte("0:00:01.00"), []byte("000:00:01.00")))
		f.Add(bytes.Replace(raw, []byte(`Hello\Nworld`), []byte(`{\k10}Hi{\p1}m 0 0 l 1 1{\p0}{\t(0,100,\bord2)}{\unknown7}世界\Nworld`), 1))
		f.Add(append(append([]byte{}, raw...), []byte("[Fonts]\nfontname: example.ttf\n!!!!\n[Graphics]\nfilename: example.png\n!!!!\n")...))
		f.Add(append(append([]byte{}, raw...), []byte("[Events]\nFormat: Text\nComment: bounded inert malformed\n")...))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 64<<10 {
			t.Skip()
		}
		detection := Detect(raw)
		if detection.Err != nil || !detection.Candidate {
			return
		}
		parsed, err := Parse(context.Background(), raw, detection.Format)
		if err != nil {
			return
		}
		d := fuzzParsedDocument(t, fixtures[detection.Format], parsed, raw)
		captured, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		first, err := Render(context.Background(), d, false)
		if err != nil {
			return
		} // Malformed retained records are restoration-only.
		afterRender, err := json.Marshal(d)
		if err != nil || !bytes.Equal(captured, afterRender) {
			t.Fatal("renderer modified source truth or native owners")
		}
		nativeParsed, err := Parse(context.Background(), first.Bytes, detection.Format)
		if err != nil {
			t.Fatalf("renderer emitted invalid grammar: %v", err)
		}
		reparsed := fuzzParsedDocument(t, fixtures[detection.Format], nativeParsed, first.Bytes)
		second, err := Render(context.Background(), reparsed, false)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(first.Bytes, second.Bytes) {
			t.Fatal("canonical render/parse bytes drift")
		}
		if len(d.Cues) != len(reparsed.Cues) {
			t.Fatal("dialogue occurrence omitted")
		}
		for i := range d.Cues {
			a, b := d.Cues[i], reparsed.Cues[i]
			if a.Timing != b.Timing || !reflect.DeepEqual(a.Payload, b.Payload) || !reflect.DeepEqual(a.Speakers, b.Speakers) || !reflect.DeepEqual(a.Tokens, b.Tokens) {
				t.Fatal("common semantic data lost")
			}
		}
		a, b := renderNative(&d), renderNative(&reparsed)
		if len(a.Sections) != len(b.Sections) || len(a.Records) != len(b.Records) || len(a.Styles) != len(b.Styles) || len(a.Events) != len(b.Events) || len(a.Attachments) != len(b.Attachments) {
			t.Fatal("ordered native occurrence omitted")
		}
		for i := range a.Sections {
			if a.Sections[i].Name != b.Sections[i].Name {
				t.Fatal("section identity lost")
			}
		}
		for i := range a.Events {
			if a.Events[i].EventType != b.Events[i].EventType || a.Events[i].Text != b.Events[i].Text {
				t.Fatal("native event Text or type changed")
			}
			if !reflect.DeepEqual(a.Events[i].Spans, b.Events[i].Spans) || !reflect.DeepEqual(a.Events[i].Tags, b.Events[i].Tags) || !reflect.DeepEqual(a.Events[i].Karaoke, b.Events[i].Karaoke) {
				t.Fatal("native override or karaoke observations changed")
			}
			fuzzRetainedFields(t, a.Events[i].Fields, b.Events[i].Fields, detection.Format, "event")
		}
		for i := range a.Styles {
			if a.Styles[i].Name != b.Styles[i].Name || a.Styles[i].Valid != b.Styles[i].Valid {
				t.Fatal("native style identity changed")
			}
			fuzzRetainedFields(t, a.Styles[i].Fields, b.Styles[i].Fields, detection.Format, "style")
		}
		for i := range a.Attachments {
			if a.Attachments[i].Name != b.Attachments[i].Name || a.Attachments[i].AttachmentType != b.Attachments[i].AttachmentType || a.Attachments[i].DataRecordCount != b.Attachments[i].DataRecordCount {
				t.Fatal("embedded content owner changed")
			}
		}
		// Canonical output may change framing and numeric lexemes. Opaque
		// attachment data and inert unknown/comment content must remain exact.
		for i := range a.Records {
			if a.Records[i].Kind != b.Records[i].Kind {
				t.Fatal("native record classification changed")
			}
			switch a.Records[i].Kind {
			case "attachment_data", "unknown", "comment":
				if !reflect.DeepEqual(a.Records[i].RawLine, b.Records[i].RawLine) {
					t.Fatal("opaque native record changed")
				}
			}
		}
	})
}

func fuzzRetainedFields(t *testing.T, a, b []model.ScriptedField, format, kind string) {
	t.Helper()
	if len(a) != len(b) {
		t.Fatal("native field occurrence omitted")
	}
	for i := range a {
		if a[i].FieldName != b[i].FieldName {
			t.Fatal("native field declaration identity changed")
		}
		identity, _ := model.ScriptedFieldIdentity(a[i].FieldName, format, kind)
		if identity == "start" || identity == "end" {
			left, err := model.ScriptedMilliseconds(a[i].RawValue)
			if err != nil {
				t.Fatal(err)
			}
			right, err := model.ScriptedMilliseconds(b[i].RawValue)
			if err != nil || left != right {
				t.Fatal("native interval field semantics changed")
			}
			continue
		}
		left, err := model.ScriptedScalarValue(a[i], format, kind, 0)
		if err != nil {
			t.Fatal(err)
		}
		right, err := model.ScriptedScalarValue(b[i], format, kind, 0)
		if err != nil || !reflect.DeepEqual(left, right) {
			t.Fatalf("native field semantics changed: %s: %v", a[i].FieldName, err)
		}
	}
}

func fuzzParsedDocument(t *testing.T, fixture []byte, parsed ParseResult, raw []byte) model.Document {
	t.Helper()
	var d model.Document
	if err := json.Unmarshal(fixture, &d); err != nil {
		t.Fatal(err)
	}
	if d.Format == "ass" {
		d.FormatData.ASS = &parsed.Native
	} else {
		d.FormatData.SSA = &parsed.Native
	}
	d.Cues, d.Diagnostics = parsed.Cues, parsed.Diagnostics
	a := &d.Source.Assets[0]
	a.DataBase64, a.Size.Bytes, a.Size.Text = base64.StdEncoding.EncodeToString(raw), int64(len(raw)), nil
	a.Hashes.SHA256, a.Encoding = fmt.Sprintf("%x", sha256.Sum256(raw)), &parsed.Encoding
	d.Document, d.Stats = model.DocumentSummary{CueCount: len(d.Cues)}, model.Stats{CueCount: len(d.Cues), DiagnosticCount: len(d.Diagnostics)}
	for _, diagnostic := range d.Diagnostics {
		if diagnostic.Severity == "warning" {
			d.Stats.WarningCount++
		}
		if diagnostic.Severity == "error" {
			d.Stats.ErrorCount++
		}
	}
	if len(d.Cues) > 0 {
		start, end := d.Cues[0].Timing.StartMilliseconds, d.Cues[0].Timing.EndMilliseconds
		for _, cue := range d.Cues {
			start, end = min(start, cue.Timing.StartMilliseconds), max(end, cue.Timing.EndMilliseconds)
			d.Document.HasWordLevelTiming = d.Document.HasWordLevelTiming || len(cue.Tokens) > 0
		}
		span := end - start
		d.Document.MediaStartMilliseconds, d.Document.MediaEndMilliseconds, d.Document.MediaSpanMilliseconds = &start, &end, &span
		d.Stats.MediaSpanMilliseconds = &span
	}
	d.Stats.HasWordLevelTiming = d.Document.HasWordLevelTiming
	if err := d.Validate(); err != nil {
		t.Fatalf("parsed fuzz model validation: %v", err)
	}
	return d
}
