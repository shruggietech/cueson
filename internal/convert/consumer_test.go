package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func consumerConversionDocument(t *testing.T, format string) model.Document {
	t.Helper()
	switch format {
	case "subrip":
		return conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:03,000\nHello\n", format)
	case "webvtt":
		return conversionTestDocument(t, "WEBVTT\n\n00:01.000 --> 00:03.000\nHello\n", format)
	default:
		return scriptedAnalysisDocument(t, format, scriptedConversionSource(format, "Hello"))
	}
}

func TestConsumerAnnotationsAccountedForAcrossAllTwelveConversionDirections(t *testing.T) {
	formats := []string{"subrip", "webvtt", "ass", "ssa"}
	for _, sourceFormat := range formats {
		for _, targetFormat := range formats {
			if sourceFormat == targetFormat {
				continue
			}
			t.Run(sourceFormat+"_"+targetFormat, func(t *testing.T) {
				document := consumerConversionDocument(t, sourceFormat)
				baseline, err := Convert(context.Background(), document, targetFormat, Options{})
				if err != nil {
					t.Fatal(err)
				}
				document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 5000}
				start, end := int64(1000), int64(2000)
				document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "PRIVATE-IDENTIFIER-one"}, {SpeakerID: "PRIVATE-IDENTIFIER-two", StartMilliseconds: &start, EndMilliseconds: &end}}
				before, err := json.Marshal(document)
				if err != nil {
					t.Fatal(err)
				}
				result, err := Convert(context.Background(), document, targetFormat, Options{})
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(baseline.Bytes, result.Bytes) {
					t.Fatal("consumer declarations changed the native subtitle representation")
				}
				paths := map[string][]string{}
				for _, loss := range result.LossReport.Losses {
					if loss.Code == LossCodeSpeakerAttributionOmitted || loss.Code == LossCodeMediaTimingOmitted {
						paths[loss.Code] = append(paths[loss.Code], loss.Path)
					}
				}
				if !reflect.DeepEqual(paths[LossCodeSpeakerAttributionOmitted], []string{"/cues/0/speaker_attributions/0", "/cues/0/speaker_attributions/1"}) || !reflect.DeepEqual(paths[LossCodeMediaTimingOmitted], []string{"/media_timing"}) {
					t.Fatalf("consumer omission paths = %#v", paths)
				}
				if err := result.LossReport.Validate(document); err != nil {
					t.Fatal(err)
				}
				reportJSON, _ := json.Marshal(result.LossReport)
				if strings.Contains(string(reportJSON), "PRIVATE-IDENTIFIER") || bytes.Contains(result.Bytes, []byte("PRIVATE-IDENTIFIER")) {
					t.Fatal("consumer identifiers leaked into native output or loss metadata")
				}
				if err := validateTargetParser(result.Bytes, targetFormat); err != nil {
					t.Fatal(err)
				}
				repeated, err := Convert(context.Background(), document, targetFormat, Options{})
				if err != nil || !bytes.Equal(result.Bytes, repeated.Bytes) || !reportsEqual(result.LossReport, repeated.LossReport) {
					t.Fatalf("nondeterministic conversion: %v", err)
				}
				called := false
				renderer := func(context.Context, model.Document, model.Document, string) ([]byte, []model.Diagnostic, error) {
					called = true
					return []byte("unexpected"), nil, nil
				}
				strictResult, err := convertWithRenderer(context.Background(), document, targetFormat, Options{Strict: true}, renderer)
				var strict *StrictLossError
				if !errors.As(err, &strict) || called || len(strictResult.Bytes) != 0 || !reportsEqual(result.LossReport, strictResult.LossReport) {
					t.Fatalf("strict conversion = %#v, err = %v, renderer_called = %t", strictResult, err, called)
				}
				after, _ := json.Marshal(document)
				if !bytes.Equal(before, after) {
					t.Fatal("conversion modified consumer annotations or source truth")
				}
			})
		}
	}
}

func TestConsumerConversionOmissionReportBound(t *testing.T) {
	document := consumerConversionDocument(t, "subrip")
	for index := 0; index < 9; index++ {
		cue := document.Cues[0]
		cue.SpeakerAttributions = make([]model.SpeakerAttribution, model.MaxItemOccurrences)
		for occurrence := range cue.SpeakerAttributions {
			cue.SpeakerAttributions[occurrence].SpeakerID = "speaker"
		}
		if index == 0 {
			document.Cues[0] = cue
		} else {
			document.Cues = append(document.Cues, cue)
		}
	}
	losses := []Loss{}
	if err := appendConsumerAnnotationLosses(context.Background(), &losses, document, "webvtt"); err == nil || len(losses) > MaxLosses {
		t.Fatalf("unbounded omission report: losses=%d error=%v", len(losses), err)
	}
}

func TestSubRipRenderReportsConsumerAnnotationsSeparately(t *testing.T) {
	document := consumerConversionDocument(t, "subrip")
	document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 5000}
	document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "PRIVATE-IDENTIFIER"}}
	diagnostics, err := SubRipRenderDiagnostics(document)
	if err != nil {
		t.Fatal(err)
	}
	codes := map[string]bool{}
	for _, diagnostic := range diagnostics {
		codes[diagnostic.Code] = true
		if strings.Contains(diagnostic.Message, "PRIVATE-IDENTIFIER") {
			t.Fatal("speaker identifier leaked into render diagnostic")
		}
	}
	if !codes["consumer_speaker_attribution_omitted"] || !codes["consumer_media_timing_omitted"] {
		t.Fatalf("consumer render codes = %#v", codes)
	}
	report, err := AnalyzeSubRipRepresentability(document)
	if err != nil || len(report.Losses) != 2 {
		t.Fatalf("consumer representability report = %#v, error = %v", report, err)
	}
}

func TestSubRipRepresentabilityRejectsInvalidTypedConsumerAnnotations(t *testing.T) {
	start, end := int64(1000), int64(4000)
	for _, test := range []struct {
		name        string
		attribution model.SpeakerAttribution
		historical  bool
	}{
		{name: "control", attribution: model.SpeakerAttribution{SpeakerID: "bad\nidentifier"}},
		{name: "partial timing", attribution: model.SpeakerAttribution{SpeakerID: "speaker", StartMilliseconds: &start}},
		{name: "outside cue", attribution: model.SpeakerAttribution{SpeakerID: "speaker", StartMilliseconds: &start, EndMilliseconds: &end}},
		{name: "historical contract", attribution: model.SpeakerAttribution{SpeakerID: "speaker"}, historical: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			document := consumerConversionDocument(t, "subrip")
			document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{test.attribution}
			if test.historical {
				document.Schema, document.SchemaVersion = "https://cueson.io/schema/v1.1.0/cueson.schema.json", "1.1.0"
			}
			if diagnostics, err := SubRipRenderDiagnostics(document); err == nil || len(diagnostics) != 0 {
				t.Fatalf("invalid annotations produced diagnostics: %#v, err = %v", diagnostics, err)
			}
			if report, err := AnalyzeSubRipRepresentability(document); err == nil || len(report.Losses) != 0 {
				t.Fatalf("invalid annotations produced report: %#v, err = %v", report, err)
			}
		})
	}
}
