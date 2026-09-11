package webvtt

import (
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestRenderCanonicalOrderAndParserCycle(t *testing.T) {
	parsed, err := Parse("WEBVTT source\n\nNOTE keep\nthis\n\nid\n00:01.000 --> 00:02.000 align:start line:10%\n<b>Hello</b>\n")
	if err != nil {
		t.Fatal(err)
	}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}, Diagnostics: parsed.Diagnostics}
	rendered, err := Render(document, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	want := "WEBVTT source\n\nNOTE keep\nthis\n\nid\n00:00:01.000 --> 00:00:02.000 align:start line:10%\n<b>Hello</b>\n"
	if string(rendered.Bytes) != want {
		t.Fatalf("Render() = %q, want %q", rendered.Bytes, want)
	}
	if _, err := Parse(string(rendered.Bytes)); err != nil {
		t.Fatalf("Parse(Render()) = %v", err)
	}
}

func TestRenderStrictRejectsPreservedConformanceDiagnostic(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000 mystery:x\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}, Diagnostics: parsed.Diagnostics}
	if _, err := Render(document, RenderOptions{Strict: true}); err == nil || !strings.Contains(err.Error(), "webvtt_setting_unknown") {
		t.Fatalf("strict Render() error = %v", err)
	}
	permissive, err := Render(document, RenderOptions{})
	if err != nil || len(permissive.Diagnostics) != 1 {
		t.Fatalf("permissive Render() = (%#v, %v)", permissive, err)
	}
}

func TestRenderStrictRecomputesConformanceFromNativeContent(t *testing.T) {
	inputs := map[string]string{
		"unknown setting": "WEBVTT\n\n00:00.000 --> 00:01.000 mystery:x\nx\n",
		"misplaced style": "WEBVTT\n\n00:00.000 --> 00:01.000\nx\n\nSTYLE\n::cue { color: red }\n",
	}
	for name, input := range inputs {
		t.Run(name, func(t *testing.T) {
			parsed, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}, Diagnostics: nil}
			if _, err := Render(document, RenderOptions{Strict: true}); err == nil {
				t.Fatal("strict Render() trusted an omitted diagnostic array")
			}
		})
	}
}

func TestRenderRegeneratesEditedRegionSettingsAndKeepsSafeUnknownOccurrences(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\nREGION\nid:r width:80% mystery:x\n\n00:00.000 --> 00:01.000\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	parsed.DocumentData.Blocks[0].Region.Settings["width"] = "50%"
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}, Diagnostics: parsed.Diagnostics}
	rendered, err := Render(document, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rendered.Bytes), "REGION\nid:r width:50% mystery:x\n") {
		t.Fatalf("Render() = %q", rendered.Bytes)
	}
}

func TestRenderRejectsUnsafeRawSettingOccurrence(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000 mystery:x\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	native := parsed.Cues[0].FormatData.WebVTT
	native.SettingsRaw = ""
	native.SettingOccurrences[0].Raw = "mystery:x\n\nNOTE injected"
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	if _, err := Render(document, RenderOptions{}); err == nil {
		t.Fatal("Render() accepted an injected raw setting occurrence")
	}
}

func TestRenderRejectsEmptyPayloadLinesButAcceptsWhitespaceOnlyContent(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	for _, lines := range [][]string{{}, {""}, {"x", ""}} {
		mutated := document
		mutated.Cues = append([]model.Cue(nil), document.Cues...)
		native := *document.Cues[0].FormatData.WebVTT
		mutated.Cues[0].FormatData.WebVTT = &native
		native.RawPayloadLines = lines
		native.RawPayload = strings.Join(lines, "\n")
		mutated.Cues[0].Payload.Lines = lines
		mutated.Cues[0].Payload.RawText = native.RawPayload
		if _, err := Render(mutated, RenderOptions{}); err == nil {
			t.Fatalf("Render() accepted payload lines %#v", lines)
		}
	}
	native := document.Cues[0].FormatData.WebVTT
	native.RawPayloadLines = []string{" "}
	native.RawPayload = " "
	document.Cues[0].Payload.Lines = []string{" "}
	document.Cues[0].Payload.RawText = " "
	if _, err := Render(document, RenderOptions{}); err != nil {
		t.Fatalf("Render() rejected whitespace-only payload: %v", err)
	}
}

func TestRenderReusesRawSettingSuffixWithoutDuplicatingDelimiter(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000   align:start  \ntext\n")
	if err != nil {
		t.Fatal(err)
	}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	rendered, err := Render(document, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	wantTiming := "00:00:00.000 --> 00:00:01.000   align:start  \n"
	if !strings.Contains(string(rendered.Bytes), wantTiming) {
		t.Fatalf("Render() = %q, want timing suffix %q", rendered.Bytes, wantTiming)
	}
}

func TestRenderRejectsForgedSettingOccurrenceFlags(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000 align:start\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	parsed.Cues[0].FormatData.WebVTT.SettingOccurrences[0].Valid = false
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	if _, err := Render(document, RenderOptions{}); err == nil {
		t.Fatal("Render() accepted forged setting occurrence flags")
	}
}

func TestRenderRejectsRecognizedBlocksWhoseRawHeadersDisagree(t *testing.T) {
	inputs := map[string]string{
		"note":   "WEBVTT\n\nNOTE source\nkept\n\n00:00.000 --> 00:01.000\nx\n",
		"style":  "WEBVTT\n\nSTYLE\n::cue { color: red }\n\n00:00.000 --> 00:01.000\nx\n",
		"region": "WEBVTT\n\nREGION\nid:r\n\n00:00.000 --> 00:01.000\nx\n",
	}
	for name, input := range inputs {
		t.Run(name, func(t *testing.T) {
			parsed, err := Parse(input)
			if err != nil {
				t.Fatal(err)
			}
			block := &parsed.DocumentData.Blocks[0]
			block.RawLines[0] = "STYLE"
			if name == "style" {
				block.RawLines[0] = ""
			}
			block.Raw = strings.Join(block.RawLines, "\n")
			document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
			if _, err := Render(document, RenderOptions{}); err == nil {
				t.Fatal("Render() reused a block with a mismatched raw header")
			}
		})
	}
}

func TestRenderRejectsUnrecognizedBlocksThatReparseAsAnotherBodyType(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\nunknown\nbody\n\n00:00.000 --> 00:01.000\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	block := &parsed.DocumentData.Blocks[0]
	for _, rawLines := range [][]string{
		{""},
		{"NOTE source", "body"},
		{"id", "00:02.000 --> 00:03.000"},
		{"bad --> timing"},
		{"id", "bad --> timing"},
	} {
		block.RawLines = rawLines
		block.Raw = strings.Join(rawLines, "\n")
		document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
		if _, err := Render(document, RenderOptions{}); err == nil {
			t.Fatalf("Render() accepted reclassifying unrecognized block %#v", rawLines)
		}
	}
}

func TestRenderAllowsLineSafeArrowInRegionButNotCueSettingOccurrence(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\nREGION\nid:r mystery:-->\n\n00:00.000 --> 00:01.000\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	rendered, err := Render(document, RenderOptions{})
	if err != nil {
		t.Fatalf("Render() rejected a line-safe REGION arrow: %v", err)
	}
	if !strings.Contains(string(rendered.Bytes), "mystery:-->") {
		t.Fatalf("Render() lost REGION occurrence: %q", rendered.Bytes)
	}
	if _, err := Parse(string(rendered.Bytes)); err != nil {
		t.Fatalf("Parse(Render()) rejected line-safe REGION arrow: %v", err)
	}

	cueParsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000 mystery:x\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	native := cueParsed.Cues[0].FormatData.WebVTT
	native.SettingsRaw = " mystery:-->"
	native.SettingOccurrences[0].Raw = "mystery:-->"
	native.SettingOccurrences[0].Value = "-->"
	document = model.Document{Format: "webvtt", Cues: cueParsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &cueParsed.DocumentData}}
	if _, err := Render(document, RenderOptions{}); err == nil {
		t.Fatal("Render() accepted a structurally unsafe cue-setting arrow")
	}
}

func TestRenderCanonicalEditIgnoresStaleNonconformingSettingsRaw(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000 mystery:x\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	native := parsed.Cues[0].FormatData.WebVTT
	native.SettingOccurrences = []model.WebVTTSettingOccurrence{}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	rendered, err := Render(document, RenderOptions{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(rendered.Bytes), "mystery:x") {
		t.Fatalf("Render() retained removed stale setting: %q", rendered.Bytes)
	}
}

func TestRenderRejectsUnknownEffectiveSetting(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:00.000 --> 00:01.000\nx\n")
	if err != nil {
		t.Fatal(err)
	}
	parsed.Cues[0].FormatData.WebVTT.Settings["injected"] = "value"
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
	if _, err := Render(document, RenderOptions{}); err == nil {
		t.Fatal("Render() silently dropped an unknown effective setting")
	}
}
