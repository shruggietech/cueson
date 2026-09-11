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
