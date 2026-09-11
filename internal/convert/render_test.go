package convert

import (
	"context"
	"reflect"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
)

func TestProjectDocumentCreatesOnlyTargetNativeState(t *testing.T) {
	source := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\n<b>x</b>\n", "subrip")
	analysis, err := analyzeCompatibility(source, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	target, err := projectDocument(source, "webvtt", analysis.Translations)
	if err != nil {
		t.Fatal(err)
	}
	if target.FormatData.SubRip != nil || target.FormatData.WebVTT == nil || len(target.FormatData.WebVTT.Blocks) != 0 {
		t.Fatalf("target document data = %#v", target.FormatData)
	}
	if target.FormatSupport.Status != "stable" || !target.FormatSupport.IngestSupported || !target.FormatSupport.RenderSupported || !target.FormatSupport.RestoreSupported {
		t.Fatalf("target format support = %#v", target.FormatSupport)
	}
	if target.Cues[0].FormatData.SubRip != nil || target.Cues[0].FormatData.WebVTT == nil || target.Cues[0].SourceIdentifier != nil {
		t.Fatalf("target cue = %#v", target.Cues[0])
	}
	if target.Cues[0].SourceOrder != 0 || target.Cues[0].Payload.RawText != "<b>x</b>" {
		t.Fatalf("target cue projection = %#v", target.Cues[0])
	}
}

func TestRenderProjectedAlwaysProducesCanonicalParserValidText(t *testing.T) {
	source := conversionTestDocument(t, "WEBVTT\n\n00:00.000 --> 00:01.000\ntext\n", "webvtt")
	analysis, err := analyzeCompatibility(source, "subrip")
	if err != nil {
		t.Fatal(err)
	}
	target, err := projectDocument(source, "subrip", analysis.Translations)
	if err != nil {
		t.Fatal(err)
	}
	if target.FormatSupport.Status != "stable" || !target.FormatSupport.IngestSupported || !target.FormatSupport.RenderSupported || !target.FormatSupport.RestoreSupported {
		t.Fatalf("target format support = %#v", target.FormatSupport)
	}
	bytes, diagnostics, err := renderProjected(context.Background(), target, "subrip")
	if err != nil || len(diagnostics) != 0 {
		t.Fatalf("render = %q, %#v, %v", bytes, diagnostics, err)
	}
	if string(bytes) != "1\n00:00:00,000 --> 00:00:01,000\ntext\n" {
		t.Fatalf("bytes = %q", bytes)
	}
}

func TestSubRipRenderDiagnosticsPreserveEstablishedAggregation(t *testing.T) {
	document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\ntext\n", "subrip")
	title := "title"
	document.Metadata.Title = &title
	document.Cues[0].Payload.RawText = "text\n"
	document.Cues[0].Speakers = []model.Speaker{{Name: "Alice", Origin: "heuristic"}}
	document.Cues[0].Tokens = []model.Token{{Text: "text", StartMilliseconds: 1000, EndMilliseconds: 2000}}
	diagnostics, err := SubRipRenderDiagnostics(document)
	if err != nil {
		t.Fatal(err)
	}
	wantCodes := []string{"subrip_render_metadata_unrepresented", "subrip_render_payload_ambiguous", "subrip_render_fields_unrepresented"}
	gotCodes := make([]string, len(diagnostics))
	for index := range diagnostics {
		gotCodes[index] = diagnostics[index].Code
	}
	if !reflect.DeepEqual(gotCodes, wantCodes) {
		t.Fatalf("diagnostic codes = %v, want %v", gotCodes, wantCodes)
	}
}

func TestSubRipRenderDiagnosticsRejectsBeforeExceedingLimit(t *testing.T) {
	document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\ntext\n", "subrip")
	cue := document.Cues[0]
	cue.Speakers = []model.Speaker{{Name: "Alice", Origin: "heuristic"}}
	document.Cues = make([]model.Cue, model.MaxDiagnostics+1)
	for index := range document.Cues {
		document.Cues[index] = cue
	}

	diagnostics, err := SubRipRenderDiagnostics(document)
	if err == nil || len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %d, error = %v", len(diagnostics), err)
	}
}

func TestAnalyzeSubRipRepresentabilityReturnsAtomicAmbiguousPayloadLoss(t *testing.T) {
	document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\ntext\n", "subrip")
	document.Cues[0].Payload.RawText = "text\n"
	report, err := AnalyzeSubRipRepresentability(document)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Losses) != 1 || report.Losses[0].Code != LossCodeSubRipPayloadAmbiguous || report.Losses[0].Path != "/cues/0/payload/raw_text" {
		t.Fatalf("report = %#v", report)
	}
}

func TestTranslatedWebVTTTargetHasNoPreservedNonconformance(t *testing.T) {
	source := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\nliteral <v Bob>\n", "subrip")
	result, err := Convert(context.Background(), source, "webvtt", Options{Strict: true})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := webvtt.Parse(string(result.Bytes))
	if err != nil || len(parsed.Diagnostics) != 0 {
		t.Fatalf("parsed target = %#v, %v", parsed, err)
	}
}
