package model

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

func addScriptedAttachment(d *Document, encoded string) {
	n := d.FormatData.ASS
	if n == nil {
		n = d.FormatData.SSA
	}
	sectionOrder := len(n.Sections) + len(n.Records)
	n.Sections = append(n.Sections, ScriptedSection{SectionID: "section-fonts", SourceOrder: sectionOrder, Name: "Fonts"})
	aid := "attachment-font"
	header := "fontname: sample.ttf"
	n.Records = append(n.Records, ScriptedRecord{RecordID: "record-font-header", SourceOrder: sectionOrder + 1, SectionID: "section-fonts", Kind: "attachment_header", RawLine: &header, AttachmentID: &aid}, ScriptedRecord{RecordID: "record-font-data", SourceOrder: sectionOrder + 2, SectionID: "section-fonts", Kind: "attachment_data", RawLine: &encoded})
	n.Attachments = append(n.Attachments, ScriptedAttachment{AttachmentID: aid, HeaderRecordID: "record-font-header", DataStartRecordID: "record-font-data", DataRecordCount: 1, AttachmentType: "font", Name: "sample.ttf"})
	syncScriptedAttachmentSource(d)
}
func syncScriptedAttachmentSource(d *Document) {
	b, _ := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
	prefix, _, _ := strings.Cut(string(b), "[Fonts]")
	n := d.FormatData.ASS
	if n == nil {
		n = d.FormatData.SSA
	}
	header := *n.Records[len(n.Records)-2].RawLine
	data := *n.Records[len(n.Records)-1].RawLine
	setScriptedSourceBytes(d, []byte(prefix+"[Fonts]\n"+header+"\n"+data+"\n"))
}
func TestScriptedAttachmentRangeAndEncodedBounds(t *testing.T) {
	for _, s := range []string{"!!", "!!!", "!!!!", strings.Repeat("!", 80)} {
		d := scriptedExample(t, "ass")
		addScriptedAttachment(&d, s)
		if e := d.Validate(); e != nil {
			t.Fatalf("accepted encoded line %q: %v", s, e)
		}
	}
	cases := []struct {
		name, want string
		change     func(*Document)
	}{
		{"encoded line ceiling", "complexity_limit", func(d *Document) { s := strings.Repeat("!", 81); d.FormatData.ASS.Records[6].RawLine = &s }},
		{"unsafe basename", "attachment", func(d *Document) { d.FormatData.ASS.Attachments[0].Name = "../sample.ttf" }},
		{"range into event", "attachment", func(d *Document) { d.FormatData.ASS.Attachments[0].DataStartRecordID = "record-event" }},
		{"range beyond section", "range", func(d *Document) { d.FormatData.ASS.Attachments[0].DataRecordCount = 2 }},
		{"bad alphabet undiagnosed", "diagnostic", func(d *Document) { s := "~!"; d.FormatData.ASS.Records[6].RawLine = &s }},
		{"bad final bits undiagnosed", "diagnostic", func(d *Document) { s := "!\""; d.FormatData.ASS.Records[6].RawLine = &s }},
		{"duplicated ownership", "attachment", func(d *Document) {
			d.FormatData.ASS.Attachments = append(d.FormatData.ASS.Attachments, d.FormatData.ASS.Attachments[0])
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := scriptedExample(t, "ass")
			addScriptedAttachment(&d, "!!")
			c.change(&d)
			syncScriptedAttachmentSource(&d)
			e := d.Validate()
			if e == nil || !strings.Contains(e.Error(), c.want) {
				t.Fatalf("%v, want %s", e, c.want)
			}
		})
	}
	d := scriptedExample(t, "ass")
	addScriptedAttachment(&d, "~!")
	order := 9
	d.Diagnostics = []Diagnostic{{Severity: "warning", Code: "malformed_attachment", Message: "Malformed embedded attachment retained.", SourceOrder: &order}}
	d.Stats.DiagnosticCount = 1
	d.Stats.WarningCount = 1
	if e := d.Validate(); e != nil {
		t.Fatalf("diagnosed malformed preservation model: %v", e)
	}
	// The physical record ceiling (65,536 *80 encoded bytes) is tighter than the
	// 16/32 MiB decoded attachment ceilings. Larger bundles reject that earlier bound.
	if MaxDocumentItems*80/4*3 >= 16<<20 {
		t.Fatal("documented attachment ceiling hierarchy changed")
	}
}
func TestScriptedPhysicalAndRepeatedCollectionBounds(t *testing.T) {
	cases := []struct {
		name   string
		change func(*Document)
	}{
		{"physical items", func(d *Document) { d.FormatData.ASS.Records = make([]ScriptedRecord, MaxDocumentItems) }},
		{"field declaration", func(d *Document) {
			d.FormatData.ASS.Records[1].DeclarationFields = make([]ScriptedDeclarationField, MaxScriptedDeclarationFields+1)
		}},
		{"item fields", func(d *Document) { d.FormatData.ASS.Records[0].Fields = make([]ScriptedField, MaxItemOccurrences+1) }},
		{"item spans", func(d *Document) { d.FormatData.ASS.Events[0].Spans = make([]ScriptedSpan, MaxItemOccurrences+1) }},
		{"item tags", func(d *Document) { d.FormatData.ASS.Events[0].Tags = make([]ScriptedTag, MaxItemOccurrences+1) }},
		{"attachments", func(d *Document) { d.FormatData.ASS.Attachments = make([]ScriptedAttachment, MaxItemOccurrences+1) }},
		{"diagnostics", func(d *Document) { d.Diagnostics = make([]Diagnostic, MaxDiagnostics+1) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := scriptedExample(t, "ass")
			c.change(&d)
			e := d.Validate()
			if e == nil || (!strings.Contains(e.Error(), "complexity_limit") && !strings.Contains(e.Error(), "maximum")) {
				t.Fatalf("ceiling not enforced: %v", e)
			}
		})
	}
	if !scriptedPhysical(strings.Repeat("x", MaxScriptedLineBytes)) || scriptedPhysical(strings.Repeat("x", MaxScriptedLineBytes+1)) {
		t.Fatal("physical line byte boundary")
	}
	for _, depth := range []int{MaxScriptedDepth, MaxScriptedDepth + 1} {
		text := "{\\t" + strings.Repeat("(", depth) + "1" + strings.Repeat(")", depth) + "}x"
		_, e := ProjectScriptedText(text, 0, 0, 1000)
		if (depth == MaxScriptedDepth) != (e == nil) {
			t.Fatalf("depth %d: %v", depth, e)
		}
	}
	for _, count := range []int{MaxItemOccurrences, MaxItemOccurrences + 1} {
		_, e := ProjectScriptedText("{"+strings.Repeat("\\b1", count)+"}x", 0, 0, 1000)
		if (count == MaxItemOccurrences) != (e == nil) {
			t.Fatalf("tag count %d: %v", count, e)
		}
	}
	if _, e := ProjectScriptedText(strings.Repeat(`\N`, MaxItemOccurrences), 0, 0, 1000); e == nil {
		t.Fatal("derived line occurrence ceiling")
	}
}
func TestScriptedScalarAndTimestampBoundaries(t *testing.T) {
	for _, s := range []string{"0:00:00.00", "123:59:59.99"} {
		if _, e := scriptedMilliseconds(s); e != nil {
			t.Fatal(e)
		}
	}
	for _, s := range []string{"0:60:00.00", "0:00:60.00", "-1:00:00.00", "0:00:00.000", "9999999999999999999:00:00.00"} {
		if _, e := scriptedMilliseconds(s); e == nil {
			t.Fatalf("bad native time accepted: %s", s)
		}
	}
	fields := []ScriptedField{{FieldName: "Fontsize", RawValue: "NaN"}, {FieldName: "Fontsize", RawValue: "+Inf"}, {FieldName: "MarginL", RawValue: "9223372036854775808"}, {FieldName: "Name", RawValue: "unsafe,comma"}, {FieldName: "Name", RawValue: "unsafe\nline"}, {FieldName: "PrimaryColour", RawValue: "&H100000000"}}
	for _, f := range fields {
		if e := validateScriptedScalar(f, "ass", 0); e == nil {
			t.Fatalf("bad native scalar accepted: %#v", f)
		}
	}
	value := int64(10)
	f := ScriptedField{FieldName: "MarginL", RawValue: "010", TypedValue: &ScriptedValue{Kind: "integer", Integer: &value}}
	if e := validateScriptedScalar(f, "ass", 0); e != nil {
		t.Fatal(e)
	}
	f.RawValue = "11"
	if e := validateScriptedScalar(f, "ass", 0); e == nil {
		t.Fatal("typed/lexical disagreement accepted")
	}
	facts, e := ProjectScriptedText("😀{\\i1}é", 0, 0, 1000)
	if e != nil || len(facts.Spans) != 3 || facts.Spans[1].StartScalar != 1 || facts.Spans[1].EndScalar != 6 || facts.Tags[0].StartScalar != 2 {
		t.Fatalf("scalar offsets: %s %v", fmt.Sprint(facts), e)
	}
}

func TestScriptedAggregateFieldAmplification(t *testing.T) {
	for _, count := range []int{MaxScriptedFields / 128, MaxScriptedFields/128 + 1} {
		d := scriptedExample(t, "ass")
		n := d.FormatData.ASS
		template := n.Styles[0].Fields
		n.Sections = []ScriptedSection{{SectionID: "section-info", SourceOrder: 0, Name: "Script Info"}, {SectionID: "section-styles", SourceOrder: 2, Name: "V4+ Styles"}, {SectionID: "section-events", SourceOrder: 3 + count, Name: "Events"}}
		n.Records = n.Records[:1]
		n.Styles = make([]ScriptedStyle, 0, count)
		n.Events = []ScriptedEvent{}
		d.Cues = []Cue{}
		d.Document = DocumentSummary{}
		d.Stats = Stats{}
		for i := 0; i < count; i++ {
			styleID := fmt.Sprintf("style-%d", i)
			recordID := fmt.Sprintf("record-%d", i)
			name := fmt.Sprintf("Style %d", i)
			fields := make([]ScriptedField, 128)
			copy(fields, template)
			fields[0].RawValue = name
			for j := len(template); j < len(fields); j++ {
				fields[j] = ScriptedField{FieldName: "X-Inert", RawValue: "opaque native content"}
			}
			n.Styles = append(n.Styles, ScriptedStyle{StyleID: styleID, RecordID: recordID, Name: name, Valid: true, Fields: fields})
			n.Records = append(n.Records, ScriptedRecord{RecordID: recordID, SourceOrder: 3 + i, SectionID: "section-styles", Kind: "style", StyleID: &styleID})
		}
		e := d.Validate()
		if count*128+1 <= MaxScriptedFields {
			if e != nil {
				t.Fatalf("within aggregate field ceiling: %v", e)
			}
		} else if e == nil || !strings.Contains(e.Error(), "aggregate occurrences") {
			t.Fatalf("aggregate field ceiling: %v", e)
		}
	}
}

func TestScriptedSignedAndUnsignedColorBoundaries(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, c := range []struct {
			lexeme string
			color  ScriptedColor
		}{
			{"-1", ScriptedColor{255, 255, 255, 255}},
			{"-2147483648", ScriptedColor{128, 0, 0, 0}},
			{"4294967295", ScriptedColor{255, 255, 255, 255}},
			{"+4294967295", ScriptedColor{255, 255, 255, 255}},
			{"&HFFFFFFFF", ScriptedColor{255, 255, 255, 255}},
			{"&H80000000", ScriptedColor{128, 0, 0, 0}},
			{"0", ScriptedColor{}},
		} {
			t.Run(format+"/"+c.lexeme, func(t *testing.T) {
				field := ScriptedField{FieldName: "PrimaryColour", RawValue: c.lexeme, TypedValue: &ScriptedValue{Kind: "color", Color: &c.color}}
				if err := validateScriptedScalar(field, format, 0); err != nil {
					t.Fatalf("valid native color/equivalent typed bits: %v", err)
				}
				wrong := c.color
				wrong.Red ^= 1
				field.TypedValue.Color = &wrong
				if err := validateScriptedScalar(field, format, 0); err == nil || !strings.Contains(err.Error(), "inconsistent_typed_value") {
					t.Fatalf("lexical/typed disagreement: %v", err)
				}
			})
		}
		for _, lexeme := range []string{"-2147483649", "4294967296", "+4294967296", "-9223372036854775808", "18446744073709551615", "&H100000000"} {
			if err := validateScriptedScalar(ScriptedField{FieldName: "PrimaryColour", RawValue: lexeme}, format, 0); err == nil || !strings.Contains(err.Error(), "invalid_native_field_value") {
				t.Fatalf("%s out-of-range %q: %v", format, lexeme, err)
			}
		}
	}
}

func TestScriptedAggregateSpanAmplification(t *testing.T) {
	d := scriptedExample(t, "ass")
	n := d.FormatData.ASS
	template := n.Events[0]
	template.Text = "{" + strings.Repeat(`\b1`, 1000) + "}literal"
	template.Fields[9].RawValue = template.Text
	facts, err := ProjectScriptedText(template.Text, 0, 1000, 2000)
	if err != nil {
		t.Fatal(err)
	}
	n.Events = []ScriptedEvent{}
	n.Records = n.Records[:4]
	d.Cues = []Cue{}
	d.Document = DocumentSummary{}
	d.Stats = Stats{}
	for i := 0; i < 1000; i++ {
		eventID := fmt.Sprintf("event-%d", i)
		recordID := fmt.Sprintf("record-event-%d", i)
		event := template
		event.EventID = eventID
		event.RecordID = recordID
		event.EventType = "comment"
		event.CueID = nil
		event.Spans = facts.Spans
		event.Tags = facts.Tags
		event.Karaoke = facts.Karaoke
		n.Events = append(n.Events, event)
		n.Records = append(n.Records, ScriptedRecord{RecordID: recordID, SourceOrder: 7 + i, SectionID: "section-events", Kind: "event", EventID: &eventID})
	}
	if e := d.Validate(); e == nil || !strings.Contains(e.Error(), "aggregate occurrences") {
		t.Fatalf("aggregate native span/tag ceiling: %v", e)
	}
}
