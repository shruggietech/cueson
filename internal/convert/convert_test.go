package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestConvertLossFreeFixturesAndTargetParserCycles(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		wantPath   string
		target     string
	}{
		{name: "SubRip to WebVTT", sourcePath: "../../testdata/fixtures/conversion/srt-loss-free/source/input.srt", wantPath: "../../testdata/fixtures/conversion/srt-loss-free/expected/bytes/output.vtt", target: "webvtt"},
		{name: "WebVTT to SubRip", sourcePath: "../../testdata/fixtures/conversion/webvtt-loss-free/source/input.vtt", wantPath: "../../testdata/fixtures/conversion/webvtt-loss-free/expected/bytes/output.srt", target: "subrip"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			source := readTestFile(t, test.sourcePath)
			document := conversionTestDocument(t, source, oppositeFormat(test.target))
			before, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			result, err := Convert(context.Background(), document, test.target, Options{Strict: true})
			if err != nil {
				t.Fatal(err)
			}
			want := readTestFile(t, test.wantPath)
			if string(result.Bytes) != want || result.LossReport.HasLosses() || len(result.Diagnostics) != 0 {
				t.Fatalf("result bytes = %q, losses = %#v, diagnostics = %#v", result.Bytes, result.LossReport, result.Diagnostics)
			}
			after, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, before) {
				t.Fatal("Convert() mutated its input document")
			}
			if err := validateTargetParser(result.Bytes, test.target); err != nil {
				t.Fatalf("target parser cycle: %v", err)
			}
			repeated, err := Convert(context.Background(), document, test.target, Options{})
			if err != nil || !bytes.Equal(result.Bytes, repeated.Bytes) || !reportsEqual(result.LossReport, repeated.LossReport) {
				t.Fatalf("repeated conversion = %#v, %v", repeated, err)
			}
		})
	}
}

func TestConvertStrictReturnsCompleteReportBeforeRendering(t *testing.T) {
	document := conversionTestDocument(t, readTestFile(t, "../../testdata/fixtures/conversion/srt-lossy/source/input.srt"), "subrip")
	called := false
	renderer := func(context.Context, model.Document, string) ([]byte, []model.Diagnostic, error) {
		called = true
		return []byte("unexpected"), nil, nil
	}
	result, err := convertWithRenderer(context.Background(), document, "webvtt", Options{Strict: true}, renderer)
	var strict *StrictLossError
	if !errors.As(err, &strict) || called || len(result.Bytes) != 0 || !result.LossReport.HasLosses() {
		t.Fatalf("result=%#v err=%v renderer_called=%t", result, err, called)
	}
	if !reportsEqual(result.LossReport, strict.Report) || len(strict.Report.Losses) < 3 {
		t.Fatalf("strict report = %#v", strict.Report)
	}
}

func TestConvertStrictRejectsWebVTTNULReplacement(t *testing.T) {
	document := conversionTestDocument(t, "WEBVTT\n\n00:00.000 --> 00:01.000\nA\x00B\n", "webvtt")
	result, err := Convert(context.Background(), document, "subrip", Options{Strict: true})
	var strict *StrictLossError
	if !errors.As(err, &strict) || len(result.Bytes) != 0 || len(result.LossReport.Losses) != 1 || result.LossReport.Losses[0].Code != LossCodeNULDegraded {
		t.Fatalf("strict NUL result = %#v, error = %v", result, err)
	}
}

func TestConvertRejectsFatalAndSameFormatBoundaries(t *testing.T) {
	document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:01,000\nx\n", "subrip")
	result, err := Convert(context.Background(), document, "webvtt", Options{})
	var projection *ProjectionError
	if !errors.As(err, &projection) || len(result.Bytes) != 0 || !strings.Contains(err.Error(), "positive cue duration") {
		t.Fatalf("zero-duration result=%#v err=%v", result, err)
	}
	_, err = Convert(context.Background(), document, "subrip", Options{})
	var same *SameFormatError
	if !errors.As(err, &same) || !strings.Contains(err.Error(), "render") {
		t.Fatalf("same-format error = %v", err)
	}
}

func conversionTestDocument(t *testing.T, input, format string) model.Document {
	t.Helper()
	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	document.Format = format
	document.FormatSupport = model.FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}
	document.Metadata = model.Metadata{}
	document.Diagnostics = []model.Diagnostic{}
	if format == "subrip" {
		parsed, parseErr := subrip.Parse(input, subrip.Options{DetectSpeakers: true})
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		document.Cues = parsed.Cues
		document.FormatData = model.DocumentFormatData{SubRip: &model.SubRipDocumentData{Dialect: "subrip"}}
		for _, diagnostic := range parsed.Diagnostics {
			document.Diagnostics = append(document.Diagnostics, model.Diagnostic(diagnostic))
		}
	} else {
		parsed, parseErr := webvtt.Parse(input)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		document.Cues = parsed.Cues
		document.FormatData = model.DocumentFormatData{WebVTT: &parsed.DocumentData}
		for _, diagnostic := range parsed.Diagnostics {
			document.Diagnostics = append(document.Diagnostics, model.Diagnostic(diagnostic))
		}
	}
	recomputeSummaries(&document)
	document.Stats.DiagnosticCount = len(document.Diagnostics)
	for _, diagnostic := range document.Diagnostics {
		if diagnostic.Severity == "warning" {
			document.Stats.WarningCount++
		}
		if diagnostic.Severity == "error" {
			document.Stats.ErrorCount++
		}
	}
	for _, cue := range document.Cues {
		if len(cue.Tokens) > 0 {
			document.Document.HasWordLevelTiming = true
			document.Stats.HasWordLevelTiming = true
		}
	}
	if err := document.Validate(); err != nil {
		t.Fatalf("test document invalid: %v", err)
	}
	return document
}

func oppositeFormat(target string) string {
	if target == "subrip" {
		return "webvtt"
	}
	return "subrip"
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func reportsEqual(left, right Report) bool {
	leftJSON, leftErr := json.Marshal(left)
	rightJSON, rightErr := json.Marshal(right)
	return leftErr == nil && rightErr == nil && bytes.Equal(leftJSON, rightJSON)
}
