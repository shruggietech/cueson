package convert

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestScriptedVariantAccountsForCommonAnnotations(t *testing.T) {
	for _, sourceFormat := range []string{"ass", "ssa"} {
		t.Run(sourceFormat, func(t *testing.T) {
			targetFormat := "ssa"
			if sourceFormat == "ssa" {
				targetFormat = "ass"
			}
			source := scriptedAnalysisDocument(t, sourceFormat, scriptedConversionSource(sourceFormat, `{\k10}Hello`))
			title, language, kind, description, line := "Shared title", "en", "captions", "Shared description", "10%"
			source.Metadata = model.Metadata{Title: &title, Language: &language, Kind: &kind, Description: &description}
			source.Cues[0].Placement = &model.Placement{Line: &line}
			source.Cues[0].OCRObservations = []model.OCRObservation{{ID: "ocr-001", Derived: true, Engine: "independent-test", Text: "observed", Lines: []string{"observed"}, SourceAssetID: source.Source.PrimaryAssetID}}
			before, err := json.Marshal(source)
			if err != nil {
				t.Fatal(err)
			}
			if err := schema.Validate(before); err != nil {
				t.Fatalf("source not valid CueJSON: %v", err)
			}
			if err := source.Validate(); err != nil {
				t.Fatal(err)
			}
			target, losses, err := projectScriptedVariant(context.Background(), source, targetFormat)
			if err != nil {
				t.Fatal(err)
			}
			report, err := NewReport(source, losses)
			if err != nil {
				t.Fatal(err)
			}
			counts := map[string]int{}
			for _, loss := range report.Losses {
				counts[loss.Code]++
			}
			if counts[LossCodeMetadataOmitted] != 4 || counts[LossCodeOCRObservationOmitted] != 1 || counts[LossCodePlacementOmitted] != 1 {
				t.Fatalf("unaccounted annotations: %+v", report)
			}
			if counts[LossCodeTokenTimingOmitted] != 0 || counts[LossCodeSpeakerObservationOmitted] != 0 {
				t.Fatalf("represented native facts reported omitted: %+v", report)
			}
			if target.Metadata != (model.Metadata{}) || len(target.Cues[0].OCRObservations) != 0 || target.Cues[0].Placement != nil {
				t.Fatal("target kept annotations it does not serialize")
			}
			result, err := scripted.RenderTarget(context.Background(), source, target)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := scripted.Parse(context.Background(), result.Bytes, targetFormat)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.Cues[0].Payload.PlainText != "Hello" || len(parsed.Cues[0].Tokens) != 1 || strings.Contains(string(result.Bytes), "Shared title") || strings.Contains(string(result.Bytes), "observed") {
				t.Fatalf("unexpected native result: %q", result.Bytes)
			}
			after, _ := json.Marshal(source)
			if string(before) != string(after) {
				t.Fatal("source annotations modified")
			}
		})
	}
}

func TestScriptedInboundBalancedEntityLiteralNeedsNoTextTargetWorkaround(t *testing.T) {
	source := scriptedTextSource(t, "webvtt", "WEBVTT\n\n00:01.000 --> 00:02.000\n&lt;b&gt;literal&lt;/b&gt; &amp; text\n")
	for _, targetFormat := range []string{"ass", "ssa"} {
		result, err := Convert(context.Background(), source, targetFormat, Options{Strict: true})
		if err != nil {
			t.Fatalf("representable angle literal rejected for %s: %v", targetFormat, err)
		}
		if result.LossReport.HasLosses() {
			t.Fatalf("unnecessary target loss: %+v", result.LossReport)
		}
		parsed, err := scripted.Parse(context.Background(), result.Bytes, targetFormat)
		if err != nil {
			t.Fatal(err)
		}
		if parsed.Cues[0].Payload.PlainText != source.Cues[0].Payload.PlainText || parsed.Cues[0].Payload.RawText != "<b>literal</b> & text" {
			t.Fatalf("literal target semantics=%+v", parsed.Cues[0].Payload)
		}
	}
}

func TestScriptedSourceIdentifierAnnotationHasExplicitOutcome(t *testing.T) {
	for _, sourceFormat := range []string{"ass", "ssa"} {
		source := scriptedAnalysisDocument(t, sourceFormat, scriptedConversionSource(sourceFormat, "Hello"))
		identifier := "custom-identifier"
		source.Cues[0].SourceIdentifier = &identifier
		encoded, _ := json.Marshal(source)
		if err := schema.Validate(encoded); err != nil {
			t.Fatalf("source annotation not valid CueJSON: %v", err)
		}
		if err := source.Validate(); err != nil {
			t.Fatal(err)
		}
		for _, targetFormat := range []string{"subrip", "webvtt", "ass", "ssa"} {
			if targetFormat == sourceFormat {
				continue
			}
			result, err := Convert(context.Background(), source, targetFormat, Options{})
			if err != nil {
				if len(result.Bytes) != 0 {
					t.Fatal("fatal annotation returned payload")
				}
				continue
			}
			accounted := false
			for _, loss := range result.LossReport.Losses {
				if loss.Path == "/cues/0/source_identifier" && loss.Code == LossCodeSourceIdentifierOmitted {
					accounted = true
				}
			}
			if !accounted {
				t.Fatalf("%s -> %s silently omitted valid identifier annotation", sourceFormat, targetFormat)
			}
			called := false
			renderer := func(context.Context, model.Document, model.Document, string) ([]byte, []model.Diagnostic, error) {
				called = true
				return []byte("unexpected"), nil, nil
			}
			strict, strictErr := convertWithRenderer(context.Background(), source, targetFormat, Options{Strict: true}, renderer)
			var strictLoss *StrictLossError
			if !errors.As(strictErr, &strictLoss) || called || len(strict.Bytes) != 0 || !reportsEqual(result.LossReport, strict.LossReport) {
				t.Fatalf("identifier strict refusal report mismatch %s -> %s: %+v %v", sourceFormat, targetFormat, strict, strictErr)
			}
		}
	}
}
