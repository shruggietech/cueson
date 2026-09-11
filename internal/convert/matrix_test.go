package convert

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestSubRipCompatibilityMatrixAccountsForEveryLossyFixtureSemantic(t *testing.T) {
	document := conversionTestDocument(t, readTestFile(t, "../../testdata/fixtures/conversion/srt-lossy/source/input.srt"), "subrip")
	analysis, err := analyzeCompatibility(document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{LossCodeSubRipCoordinatesOmitted, LossCodeSubRipFontDegraded, LossCodeSpeakerObservationOmitted}
	if got := lossCodes(report); !reflect.DeepEqual(got, want) {
		t.Fatalf("loss codes = %v, want %v", got, want)
	}
}

func TestWebVTTCompatibilityMatrixAccountsForDocumentBlockCueAndPayloadLosses(t *testing.T) {
	document := conversionTestDocument(t, readTestFile(t, "../../testdata/fixtures/conversion/webvtt-lossy/source/input.vtt"), "webvtt")
	analysis, err := analyzeCompatibility(document, "subrip")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		t.Fatal(err)
	}
	counts := make(map[string]int)
	for _, loss := range report.Losses {
		counts[loss.Code]++
		if loss.SourceFormat != "webvtt" || loss.TargetFormat != "subrip" || loss.Path == "" || loss.Message == "" {
			t.Fatalf("invalid loss = %#v", loss)
		}
	}
	want := map[string]int{
		LossCodeWebVTTDescriptionOmitted:  1,
		LossCodeWebVTTMetadataOmitted:     2,
		LossCodeWebVTTBlockOmitted:        4,
		LossCodeWebVTTIdentifierOmitted:   1,
		LossCodeWebVTTSettingOmitted:      6,
		LossCodeSpeakerObservationOmitted: 1,
		LossCodeTokenTimingOmitted:        2,
		LossCodeWebVTTVoiceDegraded:       1,
		LossCodeWebVTTMarkupDegraded:      4,
		LossCodeWebVTTInlineTimingOmitted: 1,
	}
	if !reflect.DeepEqual(counts, want) {
		t.Fatalf("loss counts = %#v, want %#v", counts, want)
	}
	settings := make([]string, 0, 6)
	for _, loss := range report.Losses {
		if loss.Code != LossCodeWebVTTSettingOmitted {
			continue
		}
		for _, attribute := range loss.Context {
			if attribute.Name == "setting" {
				settings = append(settings, attribute.Value)
			}
		}
	}
	if wantSettings := []string{"region", "vertical", "line", "position", "size", "align"}; !reflect.DeepEqual(settings, wantSettings) {
		t.Fatalf("setting loss order = %v, want %v", settings, wantSettings)
	}
	for index := 1; index < len(report.Losses); index++ {
		left, right := report.Losses[index-1], report.Losses[index]
		if left.SourceOrder != nil && right.SourceOrder != nil && *left.SourceOrder > *right.SourceOrder {
			t.Fatalf("body losses are not source ordered: %#v then %#v", left, right)
		}
	}
}

func TestCompatibilityMatrixAccountsForUnknownDuplicateSettingsAndCommonFields(t *testing.T) {
	document := conversionTestDocument(t, "WEBVTT\n\n00:00.000 --> 00:01.000 line:10% mystery:x line:20% size:101%\ntext\n", "webvtt")
	analysis, err := analyzeCompatibility(document, "subrip")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		t.Fatal(err)
	}
	counts := make(map[string]int)
	for _, loss := range report.Losses {
		counts[loss.Code]++
	}
	if counts[LossCodeWebVTTSettingOmitted] != 1 || counts[LossCodeWebVTTSettingOccurrenceOmitted] != 3 {
		t.Fatalf("setting loss counts = %#v", counts)
	}
}

func TestCompatibilityMatrixAccountsForCommonCueObservationsAtomically(t *testing.T) {
	document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:03,000\ntext\n", "subrip")
	document.Cues[0].Tokens = []model.Token{{Text: "text", StartMilliseconds: 1000, EndMilliseconds: 3000}}
	document.Cues[0].OCRObservations = []model.OCRObservation{{ID: "ocr-0", Derived: true, Text: "text", Lines: []string{"text"}, SourceAssetID: document.Source.PrimaryAssetID}}
	line := "10%"
	document.Cues[0].Placement = &model.Placement{Line: &line}
	analysis, err := analyzeCompatibility(document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{LossCodeTokenTimingOmitted: 1, LossCodeOCRObservationOmitted: 1, LossCodePlacementOmitted: 1}
	counts := make(map[string]int)
	for _, loss := range report.Losses {
		if _, tracked := want[loss.Code]; tracked {
			counts[loss.Code]++
		}
	}
	if !reflect.DeepEqual(counts, want) {
		t.Fatalf("common-field losses = %#v, want %#v", counts, want)
	}
}

func TestCompatibilityMatrixAccountsForSubRipUnrecognizedBlockDiagnostic(t *testing.T) {
	document := conversionTestDocument(t, "unrecognized block\n\n1\n00:00:01,000 --> 00:00:02,000\ntext\n", "subrip")
	analysis, err := analyzeCompatibility(document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		t.Fatal(err)
	}
	var found *Loss
	for index := range report.Losses {
		if report.Losses[index].Code == LossCodeSubRipUnrecognizedBlockOmitted {
			found = &report.Losses[index]
			break
		}
	}
	if found == nil || found.Path != "/diagnostics/0" || found.SourceOrder == nil || *found.SourceOrder != 0 || found.CueID != nil {
		t.Fatalf("unrecognized-block loss = %#v", found)
	}
}

func TestLossyFixtureReportsMatchCanonicalGoldens(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		format string
		target string
		want   string
	}{
		{name: "SubRip to WebVTT", input: "../../testdata/fixtures/conversion/srt-lossy/source/input.srt", format: "subrip", target: "webvtt", want: "df0c15af8de4929f79b185d9b2c2eb4b3ea75dc0b03166fc300e95b3a9910ccd"},
		{name: "WebVTT to SubRip", input: "../../testdata/fixtures/conversion/webvtt-lossy/source/input.vtt", format: "webvtt", target: "subrip", want: "4df221847040790c3709ac6b10156022691c3d5478c5dccb5989ae179c62fd56"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := conversionTestDocument(t, readTestFile(t, test.input), test.format)
			analysis, err := analyzeCompatibility(document, test.target)
			if err != nil {
				t.Fatal(err)
			}
			report, err := NewReport(document, analysis.Losses)
			if err != nil {
				t.Fatal(err)
			}
			canonical, err := json.Marshal(report)
			if err != nil {
				t.Fatal(err)
			}
			got := fmt.Sprintf("%x", sha256.Sum256(canonical))
			if got != test.want {
				t.Fatalf("canonical report SHA-256 = %s, want %s; report = %s", got, test.want, canonical)
			}
		})
	}
}

func TestCompatibilityAnalysisRejectsLossAmplificationBeforeBuildingAnOversizedReport(t *testing.T) {
	t.Parallel()

	document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\nx\n", "subrip")
	base := document.Cues[0]
	document.Cues = make([]model.Cue, 5)
	for index := range document.Cues {
		cue := base
		cue.ID = fmt.Sprintf("cue-%06d", index)
		cue.Ordinal = index
		cue.SourceOrder = index
		cue.Speakers = make([]model.Speaker, model.MaxItemOccurrences)
		cue.Tokens = make([]model.Token, model.MaxItemOccurrences)
		for occurrence := 0; occurrence < model.MaxItemOccurrences; occurrence++ {
			cue.Speakers[occurrence] = model.Speaker{Name: "speaker", Origin: "heuristic"}
			cue.Tokens[occurrence] = model.Token{Text: "x", StartMilliseconds: cue.Timing.StartMilliseconds, EndMilliseconds: cue.Timing.EndMilliseconds}
		}
		document.Cues[index] = cue
	}
	document.Document.CueCount = len(document.Cues)
	document.Document.HasWordLevelTiming = true
	document.Stats.CueCount = len(document.Cues)
	document.Stats.HasWordLevelTiming = true
	result, err := Convert(context.Background(), document, "webvtt", Options{})
	if err == nil || !strings.Contains(err.Error(), "more than 8192 losses") || len(result.Bytes) != 0 || len(result.LossReport.Losses) != 0 {
		t.Fatalf("conversion = %#v, %v", result, err)
	}
}

func lossCodes(report Report) []string {
	codes := make([]string, len(report.Losses))
	for index := range report.Losses {
		codes[index] = report.Losses[index].Code
	}
	return codes
}
