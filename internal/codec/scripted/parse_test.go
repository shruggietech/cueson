package scripted_test

import (
	"context"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/model"
)

func nativeFixture(format, text string) string {
	dialect, section := "v4.00+", "V4+ Styles"
	style := "Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,-1,0,0,0,100,100,0,0,1,2,1,2,10,10,10,1"
	first := "0"
	if format == "ssa" {
		dialect = "v4.00"
		section = "V4 Styles"
		style = "Default,Arial,20,16777215,255,0,-2147483648,-1,0,1,2,1,2,10,10,10,0,1"
		first = "Marked=0"
	}
	return "[Script Info]\nScriptType: " + dialect + "\n[" + section + "]\nFormat: " + strings.Join(model.ScriptedCanonicalFields(format, "style"), ",") + "\nStyle: " + style + "\n[Events]\nFormat: " + strings.Join(model.ScriptedCanonicalFields(format, "event"), ",") + "\nDialogue: " + first + ",0:00:01.00,0:00:03.00,Default,Actor,0,0,0,," + text + "\n"
}

func TestParseNativeFactsAndCaptures(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			text := `{\k20}Hi,{\K30}there\Nfile C:\spoken.srt  `
			input := "\ufeff" + strings.ReplaceAll(nativeFixture(format, text), "\n", "\r\n")
			r, err := scripted.Parse(context.Background(), []byte(input), format)
			if err != nil {
				t.Fatal(err)
			}
			if len(r.Cues) != 1 || len(r.Native.Events) != 1 || len(r.Native.Sections) != 3 {
				t.Fatalf("unexpected native ownership: %+v", r)
			}
			c, e := r.Cues[0], r.Native.Events[0]
			if c.Payload.RawText != text || e.Text != text || c.Payload.PlainText != "Hi,there\nfile C:\\spoken.srt  " {
				t.Fatalf("text was normalized: %+v", c.Payload)
			}
			if len(c.Tokens) != 2 || c.Tokens[0].StartMilliseconds != 1000 || c.Tokens[0].EndMilliseconds != 1200 || c.Tokens[1].EndMilliseconds != 1500 {
				t.Fatalf("wrong karaoke: %+v", c.Tokens)
			}
			if len(c.Speakers) != 1 || c.Speakers[0].Name != "Actor" || c.Speakers[0].Origin != "native" {
				t.Fatalf("wrong speaker: %+v", c.Speakers)
			}
			if r.Encoding.BOM == nil || *r.Encoding.BOM != "utf-8" || r.Encoding.LineEndings != "crlf" {
				t.Fatalf("encoding observation %+v", r.Encoding)
			}
			if *r.Native.Records[len(r.Native.Records)-1].RawLine != strings.TrimSuffix(strings.Split(input, "\r\n")[7], "\r") {
				t.Fatal("raw event capture differs")
			}
			if c.SourceOrder != 7 || c.ID != "cue-000000" || e.EventID != "event-000007" {
				t.Fatal("position ownership differs")
			}
		})
	}
}

func TestDetectionCannotHijackSpokenSignatures(t *testing.T) {
	for _, input := range []string{"WEBVTT\n\n00:00.000 --> 00:01.000\n[Script Info]\n", "1\n00:00:00,000 --> 00:00:01,000\n[V4+ Styles]\n"} {
		if d := scripted.Detect([]byte(input)); d.Candidate {
			t.Fatalf("spoken signature claimed script: %+v", d)
		}
	}
	for _, input := range []string{"[Script Info]\nScriptType: v4.00++\n[V4+ Styles]\n", "[Script Info]\nScriptType: v4.00+\n[V4 Styles]\n", "[Script Info]\nScriptType: v4.00+\n[V4+ Styles]\n[V4 Styles]\n"} {
		d := scripted.Detect([]byte(input))
		if !d.Candidate || d.Err == nil {
			t.Fatalf("unsupported candidate not rejected: %+v", d)
		}
	}
}

func TestParseDeclaredOrderAndSectionReset(t *testing.T) {
	input := nativeFixture("ass", "unused")
	input = strings.Replace(input, "Format: Layer,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text", "Format: End,Layer,Start,Style,Actor,MarginL,MarginR,MarginV,Effect,Extra,Extra,Text", 1)
	input = strings.Replace(input, "Dialogue: 0,0:00:01.00,0:00:03.00,Default,Actor,0,0,0,,unused", "Dialogue: 0:00:03.00,0,0:00:01.00,Default,Native Actor,0,0,0,,first,second,Text, with commas  ", 1)
	r, err := scripted.Parse(context.Background(), []byte(input), "ass")
	if err != nil {
		t.Fatal(err)
	}
	e := r.Native.Events[0]
	if len(e.Fields) != 12 || e.Fields[9].RawValue != "first" || e.Fields[10].RawValue != "second" || e.Text != "Text, with commas  " {
		t.Fatalf("declared order/suffix lost: %+v", e)
	}
	if r.Cues[0].Speakers[0].Name != "Native Actor" || *r.Cues[0].FormatData.ASS.Projection.SpeakerFieldName != "Actor" {
		t.Fatal("Actor alias provenance differs")
	}
	if _, err = scripted.Parse(context.Background(), []byte(input+"[Events]\nDialogue: 0,0:00:04.00,0:00:05.00,Default,,0,0,0,,bad\n"), "ass"); err == nil {
		t.Fatal("declaration leaked across section occurrence")
	}
}

func TestParseEmptyDrawingAndUnsupportedKaraoke(t *testing.T) {
	for _, text := range []string{"", `{\p1}m 0 0 l 10 10{\p0}`} {
		r, err := scripted.Parse(context.Background(), []byte(nativeFixture("ass", text)), "ass")
		if err != nil {
			t.Fatal(err)
		}
		if len(r.Cues) != 1 || r.Cues[0].Payload.PlainText != "" || len(r.Cues[0].Payload.Lines) != 1 || len(r.Cues[0].OCRObservations) != 0 {
			t.Fatalf("empty timed native unit lost: %+v", r.Cues)
		}
		if text != "" && !r.Cues[0].FormatData.ASS.Projection.DrawingExcluded {
			t.Fatal("drawing lacks exclusion provenance")
		}
	}
	r, err := scripted.Parse(context.Background(), []byte(nativeFixture("ssa", `{\kt20}bad{\k300}long`)), "ssa")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Cues[0].Tokens) != 0 || len(r.Native.Events[0].Karaoke) != 2 || len(r.Diagnostics) != 1 || r.Diagnostics[0].Code != "unsupported_karaoke_timing" {
		t.Fatalf("fabricated karaoke timing: %+v", r)
	}
}

func TestParseFatalProfileAndPrivacy(t *testing.T) {
	base := nativeFixture("ass", "Hello")
	cases := map[string]string{
		"leading blank":                       "\n" + base,
		"mismatched explicit":                 nativeFixture("ssa", "Hello"),
		"NUL":                                 strings.Replace(base, "Hello", "\x00", 1),
		"bare CR":                             strings.Replace(base, "Hello", "Hello\rbroken", 1),
		"invalid UTF8":                        base + string([]byte{0xff}),
		"timer":                               strings.Replace(base, "ScriptType: v4.00+", "ScriptType: v4.00+\nTimer: 99", 1),
		"unknown metadata":                    strings.Replace(base, "ScriptType: v4.00+", "ScriptType: v4.00+\nUnknown: innocent", 1),
		"relative path":                       strings.Replace(base, "ScriptType: v4.00+", "ScriptType: v4.00+\nTitle: secret/thing", 1),
		"resource name":                       strings.Replace(base, "ScriptType: v4.00+", "ScriptType: v4.00+\nVideo File: portable.mp4", 1),
		"active effect":                       strings.Replace(base, ",Actor,0,0,0,,Hello", ",Actor,0,0,0,code once,Hello", 1),
		"bad interval":                        strings.Replace(base, "0:00:03.00", "0:00:01.00", 1),
		"overflow":                            strings.Replace(base, "0:00:03.00", "9223372036854775807:00:00.00", 1),
		"style unresolved":                    strings.Replace(base, ",Default,Actor,", ",Missing,Actor,", 1),
		"ambiguous Text":                      strings.Replace(base, "Effect,Text", "Text,Effect", 1),
		"malformed Comment content exemption": strings.Replace(base, "Dialogue: 0,0:00:01.00,0:00:03.00,Default,Actor,0,0,0,,Hello", "Comment: secret/thing", 1),
		"unknown colon content":               base + "Other: data\n",
		"unsafe declaration name":             strings.Replace(base, "Effect,Text", "Effect,secret/thing,Text", 1),
		"unsafe section identity":             base + "[localhost]\n",
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			r, err := scripted.Parse(context.Background(), []byte(input), "ass")
			if err == nil {
				t.Fatalf("accepted unsafe/fatal input: %+v", r)
			}
			if strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "portable.mp4") {
				t.Fatalf("unsafe diagnostic echo: %v", err)
			}
		})
	}
}

func TestParseUnknownAndMalformedRetention(t *testing.T) {
	input := nativeFixture("ass", `{\unknown9}text`)
	input = strings.Replace(input, "[V4+ Styles]", "ScriptType: v4.00+\nTimer: 100\nTimer: 100\n[V4+ Styles]", 1)
	input = strings.Replace(input, "[Events]", "Style: Unused,Arial,bad,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,1,2,10,10,10,1\n[Events]", 1)
	input += "; harmless comment\nunknown inert line\n[Unknown Safe]\n\n[Fonts]\nfontname: unused.ttf\n"
	r, err := scripted.Parse(context.Background(), []byte(input), "ass")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Native.Styles) != 2 || r.Native.Styles[1].Valid || r.Native.Styles[1].Name != "Unused" {
		t.Fatalf("invalid style occurrence lost: %+v", r.Native.Styles)
	}
	want := []string{"duplicate_script_type", "duplicate_timer", "malformed_native_record", "unsupported_override", "unknown_native_record", "unknown_native_section", "malformed_native_record"}
	if len(r.Diagnostics) != len(want) {
		t.Fatalf("diagnostics %+v", r.Diagnostics)
	}
	for i, d := range r.Diagnostics {
		if d.Code != want[i] {
			t.Fatalf("diagnostic %d = %q want %q", i, d.Code, want[i])
		}
	}
	if r.Native.Records[len(r.Native.Records)-1].Kind != "malformed" {
		t.Fatal("empty attachment header lost")
	}
}

func TestParseBoundsAndCancellation(t *testing.T) {
	prefix := "[Script Info]\nScriptType: v4.00+\n[V4+ Styles]\n"
	for _, tc := range []struct {
		name, input string
		good        bool
	}{
		{"physical item exact", prefix + strings.Repeat("\n", model.MaxDocumentItems-3), true},
		{"physical item above", prefix + strings.Repeat("\n", model.MaxDocumentItems-2), false},
		{"physical line exact", prefix + "[Events]\n" + strings.Repeat("a", model.MaxScriptedLineBytes), true},
		{"physical line above", prefix + "[Events]\n" + strings.Repeat("a", model.MaxScriptedLineBytes+1), false},
		{"diagnostic exact", prefix + "[Events]\n" + strings.Repeat("inert\n", model.MaxDiagnostics), true},
		{"diagnostic above", prefix + "[Events]\n" + strings.Repeat("inert\n", model.MaxDiagnostics+1), false},
		{"depth above", nativeFixture("ass", "{\\t"+strings.Repeat("(", model.MaxScriptedDepth+1)+"x"+strings.Repeat(")", model.MaxScriptedDepth+1)+"}text"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := scripted.Parse(context.Background(), []byte(tc.input), "ass")
			if (err == nil) != tc.good {
				t.Fatalf("got %v, good=%v", err, tc.good)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := scripted.Parse(ctx, []byte(prefix), "ass"); err == nil {
		t.Fatal("canceled parse succeeded")
	}
}

func FuzzParse(f *testing.F) {
	f.Add([]byte(nativeFixture("ass", `{\k10}Hello\Nworld`)), "ass")
	f.Add([]byte(nativeFixture("ssa", "Hello, world")), "ssa")
	f.Add([]byte("[Script Info]\nScriptType: v4.00++\n"), "ass")
	f.Fuzz(func(t *testing.T, input []byte, format string) {
		if len(input) > 2<<20 {
			t.Skip()
		}
		r, err := scripted.Parse(context.Background(), input, format)
		if err == nil {
			if len(r.Native.Sections)+len(r.Native.Records) > model.MaxDocumentItems || len(r.Diagnostics) > model.MaxDiagnostics {
				t.Fatal("accepted above bounds")
			}
			if d := scripted.Detect(input); d.Err != nil || d.Format != format {
				t.Fatal("parse/detect disagreement")
			}
		}
	})
}
