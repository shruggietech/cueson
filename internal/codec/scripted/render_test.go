package scripted

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func renderFixture(t *testing.T, format string) model.Document {
	t.Helper()
	raw, err := os.ReadFile("../../schema/testdata/scripted-" + format + ".cueson.json")
	if err != nil {
		t.Fatal(err)
	}
	var document model.Document
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatal(err)
	}
	return document
}

func renderNative(document *model.Document) *model.ScriptedDocumentData {
	if document.Format == "ssa" {
		return document.FormatData.SSA
	}
	return document.FormatData.ASS
}

func renderSetTimes(document *model.Document, start, end int64) {
	document.Cues[0].Timing = model.Timing{StartMilliseconds: start, EndMilliseconds: end, DurationMilliseconds: end - start}
	span := end - start
	document.Document.MediaStartMilliseconds = &start
	document.Document.MediaEndMilliseconds = &end
	document.Document.MediaSpanMilliseconds = &span
	document.Stats.MediaSpanMilliseconds = &span
}

func TestRenderCanonicalBothDialects(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			document := renderFixture(t, format)
			original, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
			if err != nil {
				t.Fatal(err)
			}
			result, err := Render(context.Background(), document, true)
			if err != nil {
				t.Fatal(err)
			}
			expected := string(original)
			lines := strings.Split(expected, "\n")
			for i, line := range lines {
				if strings.HasPrefix(line, "Format: ") {
					lines[i] = strings.ReplaceAll(line, ", ", ",")
				}
			}
			expected = strings.Join(lines, "\n")
			if format == "ssa" {
				expected = strings.NewReplacer("&H00FFFFFF", "16777215", "&H000000FF", "255", "&H00000000", "0").Replace(expected)
			}
			if string(result.Bytes) != expected {
				t.Fatalf("canonical output\ngot %q\nwant %q", result.Bytes, expected)
			}
			second, err := Render(context.Background(), document, true)
			if err != nil || !bytes.Equal(result.Bytes, second.Bytes) {
				t.Fatalf("determinism: %v", err)
			}
			if len(result.Diagnostics) != 0 {
				t.Fatalf("unexpected warnings: %v", result.Diagnostics)
			}
			if document.Source.Assets[0].DataBase64 != base64.StdEncoding.EncodeToString(original) {
				t.Fatal("source capture changed")
			}
		})
	}
}

func TestRenderTimingAndTextOwners(t *testing.T) {
	document := renderFixture(t, "ass")
	native := renderNative(&document)
	before := *native.Records[len(native.Records)-1].RawLine
	renderSetTimes(&document, 1230, 4560)
	event := &native.Events[0]
	event.Text = "Edited, comma  "
	for i := range event.Fields {
		if renderName(event.Fields[i].FieldName) == "text" {
			event.Fields[i].RawValue = event.Text
			event.Fields[i].TypedValue = nil
		}
	}
	facts, err := model.ProjectScriptedText(event.Text, 0, 1230, 4560)
	if err != nil {
		t.Fatal(err)
	}
	event.Spans, event.Tags, event.Karaoke = facts.Spans, facts.Tags, facts.Karaoke
	document.Cues[0].Payload = model.Payload{RawText: event.Text, PlainText: event.Text, Lines: facts.Lines}
	result, err := Render(context.Background(), document, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result.Bytes), "Dialogue: 0,0:00:01.23,0:00:04.56,Default,Narrator,0,0,0,,Edited, comma  \n") {
		t.Fatalf("owner edits ignored: %s", result.Bytes)
	}
	if *native.Records[len(native.Records)-1].RawLine != before {
		t.Fatal("capture line rewritten")
	}
	if *document.Cues[0].FormatData.ASS.StartTimestampRaw != "0:00:01.00" {
		t.Fatal("capture timestamp rewritten")
	}
}

func TestRenderConstructedCanonicalDeclaration(t *testing.T) {
	document := renderFixture(t, "ass")
	native := renderNative(&document)
	for i := range native.Sections {
		native.Sections[i].RawHeader = nil
	}
	for i := range native.Records {
		native.Records[i].RawLine = nil
	}
	native.Styles[0].DeclarationID = nil
	native.Events[0].DeclarationID = nil
	// Remove both explicit declarations and compact the physical sequence.
	records := native.Records[:0]
	for _, record := range native.Records {
		if record.Kind != "format_declaration" {
			records = append(records, record)
		}
	}
	native.Records = records
	native.Sections[0].SourceOrder = 0
	native.Records[0].SourceOrder = 1
	native.Sections[1].SourceOrder = 2
	native.Records[1].SourceOrder = 3
	native.Sections[2].SourceOrder = 4
	native.Records[2].SourceOrder = 5
	document.Cues[0].SourceOrder = 5
	document.Cues[0].FormatData.ASS.StartTimestampRaw = nil
	document.Cues[0].FormatData.ASS.EndTimestampRaw = nil
	result, err := Render(context.Background(), document, true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(result.Bytes), "Format: ") != 2 || !strings.Contains(string(result.Bytes), "Format: Layer,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text\nDialogue: ") {
		t.Fatalf("constructed declarations absent: %s", result.Bytes)
	}
}

func TestRenderRejectsInconsistentUnsafeAndUnrepresentableEdits(t *testing.T) {
	cases := []struct {
		name, code string
		edit       func(*model.Document)
	}{
		{"precision", "unrepresentable_centiseconds", func(d *model.Document) { renderSetTimes(d, 1011, 2000) }},
		{"derived-only", "inconsistent_projection", func(d *model.Document) { d.Cues[0].Payload.PlainText = "fabricated" }},
		{"typed-disagreement", "inconsistent_typed_value", func(d *model.Document) {
			value := 40.0
			renderNative(d).Styles[0].Fields[2].TypedValue.Decimal = &value
		}},
		{"metadata-path", "unsafe_source_metadata", func(d *model.Document) {
			renderNative(d).Records[0].Fields = append(renderNative(d).Records[0].Fields, model.ScriptedField{FieldName: "Title", RawValue: "relative/path"})
		}},
		{"style-comma", "invalid_native_field_value", func(d *model.Document) {
			renderNative(d).Styles[0].Fields[1].RawValue = "Arial,Other"
			renderNative(d).Styles[0].Fields[1].TypedValue = nil
		}},
		{"text-newline", "complexity_limit", func(d *model.Document) { renderNative(d).Events[0].Text = "injected\nDialogue: bad" }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := renderFixture(t, "ass")
			tc.edit(&d)
			result, err := Render(context.Background(), d, false)
			if err == nil || !strings.Contains(err.Error(), tc.code) {
				t.Fatalf("want %s, got %v", tc.code, err)
			}
			if len(result.Bytes) != 0 {
				t.Fatal("partial output")
			}
		})
	}
}

func TestRenderCancelledAndWrongFormat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if result, err := Render(ctx, renderFixture(t, "ass"), false); err == nil || len(result.Bytes) != 0 {
		t.Fatal("cancelled render published bytes")
	}
	if result, err := Render(context.Background(), model.Document{Format: "subrip"}, false); err == nil || len(result.Bytes) != 0 {
		t.Fatal("wrong format rendered")
	}
}

func TestRenderScalarCanonicalization(t *testing.T) {
	cases := []struct{ name, raw, format, want string }{
		{"Bold", "+1", "ass", "-1"}, {"Fontsize", "00020.000", "ass", "20"}, {"MarginL", "+0010", "ass", "10"},
		{"PrimaryColour", "-1", "ass", "&HFFFFFFFF"}, {"TertiaryColour", "-2147483648", "ssa", "2147483648"},
	}
	for _, tc := range cases {
		result, err := renderScalar(model.ScriptedField{FieldName: tc.name, RawValue: tc.raw}, tc.format, "style")
		if err != nil || result != tc.want {
			t.Fatalf("%s: got %q/%v want %q", tc.name, result, err, tc.want)
		}
	}
}

func TestRenderPhysicalAndAggregateOutputBounds(t *testing.T) {
	for _, tc := range []struct {
		name              string
		count, titleBytes int
	}{{"physical-line", 1, model.MaxScriptedLineBytes}, {"aggregate-output", 65, model.MaxScriptedLineBytes - 8}} {
		t.Run(tc.name, func(t *testing.T) {
			d := renderFixture(t, "ass")
			n := renderNative(&d)
			n.Sections = []model.ScriptedSection{{SectionID: "section-info", SourceOrder: 0, Name: "Script Info"}, {SectionID: "section-styles", SourceOrder: tc.count + 2, Name: "V4+ Styles"}}
			n.Records = []model.ScriptedRecord{{RecordID: "record-type", SourceOrder: 1, SectionID: "section-info", Kind: "metadata", Fields: []model.ScriptedField{{FieldName: "ScriptType", RawValue: "v4.00+"}}}}
			value := strings.Repeat("a", tc.titleBytes)
			for i := 0; i < tc.count; i++ {
				n.Records = append(n.Records, model.ScriptedRecord{RecordID: "record-title-" + strconv.Itoa(i), SourceOrder: i + 2, SectionID: "section-info", Kind: "metadata", Fields: []model.ScriptedField{{FieldName: "Title", RawValue: value}}})
			}
			n.Styles = []model.ScriptedStyle{}
			n.Events = []model.ScriptedEvent{}
			d.Cues = []model.Cue{}
			d.Document = model.DocumentSummary{}
			d.Stats = model.Stats{}
			if err := d.Validate(); err != nil {
				t.Fatalf("accepted model before serialization: %v", err)
			}
			result, err := Render(context.Background(), d, false)
			if err == nil || !strings.Contains(err.Error(), "complexity_limit") || len(result.Bytes) != 0 {
				t.Fatalf("bound failure returned output: %v/%d", err, len(result.Bytes))
			}
		})
	}
}

func renderParsedFixture(t *testing.T, format string, raw []byte) model.Document {
	t.Helper()
	p, err := Parse(context.Background(), raw, format)
	if err != nil {
		t.Fatal(err)
	}
	d := renderFixture(t, format)
	n := renderNative(&d)
	*n = p.Native
	d.Cues = p.Cues
	d.Diagnostics = p.Diagnostics
	a := &d.Source.Assets[0]
	a.DataBase64 = base64.StdEncoding.EncodeToString(raw)
	a.Size.Bytes = int64(len(raw))
	a.Size.Text = nil
	a.Hashes.SHA256 = fmt.Sprintf("%x", sha256.Sum256(raw))
	a.Encoding = &p.Encoding
	d.Document = model.DocumentSummary{CueCount: len(d.Cues)}
	d.Stats = model.Stats{CueCount: len(d.Cues), DiagnosticCount: len(d.Diagnostics)}
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
		for _, c := range d.Cues {
			start = min(start, c.Timing.StartMilliseconds)
			end = max(end, c.Timing.EndMilliseconds)
			if len(c.Tokens) > 0 {
				d.Document.HasWordLevelTiming = true
				d.Stats.HasWordLevelTiming = true
			}
		}
		span := end - start
		d.Document.MediaStartMilliseconds = &start
		d.Document.MediaEndMilliseconds = &end
		d.Document.MediaSpanMilliseconds = &span
		d.Stats.MediaSpanMilliseconds = &span
	}
	if err := d.Validate(); err != nil {
		t.Fatalf("parsed fixture validation: %v", err)
	}
	return d
}

func TestRenderParsedDeclarationSpacingAndTransitions(t *testing.T) {
	d := renderFixture(t, "ass")
	raw, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
	lines := strings.Split(string(raw), "\n")
	last := lines[len(lines)-2]
	raw = append(raw, []byte(last+"\n"+last+"\n")...)
	d = renderParsedFixture(t, "ass", raw)
	n := renderNative(&d)
	// The middle occurrence has structured construction ownership. It needs a
	// canonical declaration, then the final captured occurrence needs its
	// original declaration restored without replaying either raw event line.
	n.Events[1].DeclarationID = nil
	for i := range n.Records {
		if n.Records[i].EventID != nil && *n.Records[i].EventID == n.Events[1].EventID {
			n.Records[i].RawLine = nil
		}
	}
	d.Cues[1].FormatData.ASS.StartTimestampRaw = nil
	d.Cues[1].FormatData.ASS.EndTimestampRaw = nil
	result, err := Render(context.Background(), d, true)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(result.Bytes), "Format: ") != 4 {
		t.Fatalf("declaration transition missing: %s", result.Bytes)
	}
	parsed := renderParsedFixture(t, "ass", result.Bytes)
	second, err := Render(context.Background(), parsed, true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(result.Bytes, second.Bytes) {
		t.Fatalf("declaration spacing drift\nfirst %q\nsecond %q", result.Bytes, second.Bytes)
	}
	for i := range parsed.Cues {
		if parsed.Cues[i].Payload.RawText != d.Cues[i].Payload.RawText || parsed.Cues[i].Timing != d.Cues[i].Timing {
			t.Fatal("event semantics drift")
		}
	}
}

func TestRenderUnknownAndMalformedContentPolicies(t *testing.T) {
	base := renderFixture(t, "ass")
	raw, _ := base64.StdEncoding.DecodeString(base.Source.Assets[0].DataBase64)
	unknown := renderParsedFixture(t, "ass", append(append([]byte{}, raw...), []byte("Inert extension\n")...))
	if result, err := Render(context.Background(), unknown, false); err != nil || !strings.HasSuffix(string(result.Bytes), "Inert extension\n") || len(result.Diagnostics) == 0 {
		t.Fatalf("permissive retention: %v", err)
	}
	if result, err := Render(context.Background(), unknown, true); err == nil || len(result.Bytes) != 0 {
		t.Fatal("strict unknown content published")
	}
	for _, tail := range []string{"Comment: inert\n", "[Fonts]\nfontname: embedded.ttf\n!\n"} {
		d := renderParsedFixture(t, "ass", append(append([]byte{}, raw...), []byte(tail)...))
		for _, strict := range []bool{false, true} {
			if result, err := Render(context.Background(), d, strict); err == nil || len(result.Bytes) != 0 {
				t.Fatalf("malformed preservation rendered in strict=%v", strict)
			}
		}
	}
}

func TestRenderCompleteSourceIntegrityAndSafeFailure(t *testing.T) {
	d := renderFixture(t, "ass")
	extra := d.Source.Assets[0]
	extra.ID = "asset-extra"
	extra.Role = "related"
	extra.FileName = "related.ass"
	extra.Hashes.SHA256 = strings.Repeat("0", 64)
	d.Source.Assets = append(d.Source.Assets, extra)
	if result, err := Render(context.Background(), d, false); err == nil || !strings.Contains(err.Error(), "invalid_source_integrity") || len(result.Bytes) != 0 {
		t.Fatalf("secondary source integrity skipped: %v", err)
	}
	d = renderFixture(t, "ass")
	d.Source.Assets[0].FileName = "C:/private/machine/source.ass"
	if result, err := Render(context.Background(), d, false); err == nil || strings.Contains(err.Error(), "private") || len(result.Bytes) != 0 {
		t.Fatalf("unsafe asset name escaped into diagnostic/output: %v", err)
	}
}

func TestRenderRejectsUnusedDeclarationRecognitionForgery(t *testing.T) {
	d := renderFixture(t, "ass")
	n := renderNative(&d)
	// Append an unused declaration in the final Events occurrence. An unused
	// declaration still owns an explicit recognized-field observation.
	declaration := n.Records[len(n.Records)-2]
	declaration.RecordID = "record-extra-declaration"
	declaration.SourceOrder = len(n.Sections) + len(n.Records)
	declaration.RawLine = nil
	declaration.DeclarationFields = append([]model.ScriptedDeclarationField{}, declaration.DeclarationFields...)
	wrong := "end"
	declaration.DeclarationFields[0].RecognizedName = &wrong
	n.Records = append(n.Records, declaration)
	if err := d.Validate(); err != nil {
		t.Fatalf("unused declaration model fixture: %v", err)
	}
	if result, err := Render(context.Background(), d, false); err == nil || !strings.Contains(err.Error(), "inconsistent_recognized_name") || len(result.Bytes) != 0 {
		t.Fatalf("unused declaration silently repaired: %v", err)
	}
}
