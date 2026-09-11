package convert

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestNewReportProducesCanonicalAtomicOrder(t *testing.T) {
	t.Parallel()

	document := reportTestDocument(t)
	order, cueID := document.Cues[0].SourceOrder, document.Cues[0].ID
	losses := []Loss{
		newTestLoss(LossCodeOCRObservationOmitted, KindOmitted, "/cues/0/ocr_observations/0", &order, &cueID),
		newTestLoss(LossCodeMetadataOmitted, KindOmitted, "/metadata/title", nil, nil),
		newTestLoss(LossCodeSubRipCoordinatesOmitted, KindOmitted, "/cues/0/format_data/subrip/coordinates", &order, &cueID),
	}
	losses[0].Context = []Attribute{{Name: "observation_id", Value: "ocr-0"}, {Name: "field", Value: "text"}}

	report, err := NewReport(document, losses)
	if err != nil {
		t.Fatalf("NewReport() error = %v", err)
	}
	wantCodes := []string{LossCodeMetadataOmitted, LossCodeSubRipCoordinatesOmitted, LossCodeOCRObservationOmitted}
	gotCodes := make([]string, len(report.Losses))
	for index := range report.Losses {
		gotCodes[index] = report.Losses[index].Code
	}
	if !reflect.DeepEqual(gotCodes, wantCodes) {
		t.Fatalf("loss codes = %v, want %v", gotCodes, wantCodes)
	}
	if got := report.Losses[2].Context; !reflect.DeepEqual(got, []Attribute{{Name: "field", Value: "text"}, {Name: "observation_id", Value: "ocr-0"}}) {
		t.Fatalf("sorted context = %#v", got)
	}
	if !report.HasLosses() || (Report{}).HasLosses() {
		t.Fatal("HasLosses() returned an inconsistent result")
	}
	if err := report.Validate(document); err != nil {
		t.Fatalf("Report.Validate() error = %v", err)
	}
	if losses[0].Context[0].Name != "observation_id" {
		t.Fatal("NewReport() mutated caller-owned loss context")
	}
}

func TestNewReportUsesNumericPointerOrdering(t *testing.T) {
	t.Parallel()

	document := reportTestDocument(t)
	losses := make([]Loss, 11)
	for index := range losses {
		path := fmt.Sprintf("/metadata/items/%d", 10-index)
		losses[index] = newTestLoss(LossCodeMetadataOmitted, KindOmitted, path, nil, nil)
	}
	report, err := NewReport(document, losses)
	if err != nil {
		t.Fatal(err)
	}
	for index := range report.Losses {
		want := fmt.Sprintf("/metadata/items/%d", index)
		if report.Losses[index].Path != want {
			t.Fatalf("loss %d path = %q, want %q", index, report.Losses[index].Path, want)
		}
	}
}

func TestReportValidationRejectsInvalidOrUnsafeLosses(t *testing.T) {
	t.Parallel()

	document := reportTestDocument(t)
	order, cueID := document.Cues[0].SourceOrder, document.Cues[0].ID
	valid := newTestLoss(LossCodeSubRipCoordinatesOmitted, KindOmitted, "/cues/0/format_data/subrip/coordinates", &order, &cueID)
	tests := []struct {
		name   string
		mutate func(*Loss)
		want   string
	}{
		{name: "unknown code", mutate: func(loss *Loss) { loss.Code = "conversion_unknown" }, want: "unknown loss code"},
		{name: "severity", mutate: func(loss *Loss) { loss.Severity = "error" }, want: "severity"},
		{name: "kind", mutate: func(loss *Loss) { loss.Kind = KindDegraded }, want: "kind"},
		{name: "source format", mutate: func(loss *Loss) { loss.SourceFormat = "webvtt" }, want: "source format"},
		{name: "same pair", mutate: func(loss *Loss) { loss.TargetFormat = "subrip" }, want: "format pair"},
		{name: "message empty", mutate: func(loss *Loss) { loss.Message = "" }, want: "message"},
		{name: "message path leak", mutate: func(loss *Loss) { loss.Message = `read C:\Users\operator\source.srt` }, want: "path-like"},
		{name: "path relative", mutate: func(loss *Loss) { loss.Path = "cues/0" }, want: "JSON Pointer"},
		{name: "path escape", mutate: func(loss *Loss) { loss.Path = "/cues/~2bad" }, want: "JSON Pointer"},
		{name: "source order unresolved", mutate: func(loss *Loss) { value := 99; loss.SourceOrder = &value }, want: "source_order"},
		{name: "cue unresolved", mutate: func(loss *Loss) { value := "missing"; loss.CueID = &value }, want: "cue_id"},
		{name: "cue order mismatch", mutate: func(loss *Loss) { value := 1; loss.SourceOrder = &value }, want: "source_order"},
		{name: "context duplicate", mutate: func(loss *Loss) {
			loss.Context = []Attribute{{Name: "field", Value: "one"}, {Name: "field", Value: "two"}}
		}, want: "context"},
		{name: "context order", mutate: func(loss *Loss) {
			loss.Context = []Attribute{{Name: "zeta", Value: "one"}, {Name: "alpha", Value: "two"}}
		}, want: "context"},
		{name: "context bound", mutate: func(loss *Loss) {
			loss.Context = []Attribute{{Name: "field", Value: strings.Repeat("x", maxAttributeValueBytes+1)}}
		}, want: "context"},
		{name: "context path leak", mutate: func(loss *Loss) { loss.Context = []Attribute{{Name: "field", Value: `C:\Users\operator`}} }, want: "path-like"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			loss := cloneLoss(valid)
			test.mutate(&loss)
			err := (Report{Losses: []Loss{loss}}).Validate(document)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Report.Validate() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestReportValidationRejectsDuplicateAndNoncanonicalOrder(t *testing.T) {
	t.Parallel()

	document := reportTestDocument(t)
	first := newTestLoss(LossCodeMetadataOmitted, KindOmitted, "/metadata/title", nil, nil)
	second := newTestLoss(LossCodeMetadataOmitted, KindOmitted, "/metadata/language", nil, nil)

	if err := (Report{Losses: []Loss{first, first}}).Validate(document); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate validation error = %v", err)
	}
	canonical, err := NewReport(document, []Loss{first, second})
	if err != nil {
		t.Fatal(err)
	}
	canonical.Losses[0], canonical.Losses[1] = canonical.Losses[1], canonical.Losses[0]
	if err := canonical.Validate(document); err == nil || !strings.Contains(err.Error(), "canonical order") {
		t.Fatalf("order validation error = %v", err)
	}
}

func TestLossReportsRejectCollectionAmplificationBeforeCanonicalization(t *testing.T) {
	t.Parallel()

	document := reportTestDocument(t)
	losses := make([]Loss, MaxLosses+1)
	if _, err := NewReport(document, losses); err == nil || !strings.Contains(err.Error(), "maximum") {
		t.Fatalf("NewReport() error = %v, want bounded rejection", err)
	}
	if err := (Report{Losses: losses}).Validate(document); err == nil || !strings.Contains(err.Error(), "maximum") {
		t.Fatalf("Report.Validate() error = %v, want bounded rejection", err)
	}
}

func TestTypedConversionErrorsRetainStableContext(t *testing.T) {
	t.Parallel()

	document := reportTestDocument(t)
	loss := newTestLoss(LossCodeMetadataOmitted, KindOmitted, "/metadata/title", nil, nil)
	report, err := NewReport(document, []Loss{loss})
	if err != nil {
		t.Fatal(err)
	}

	strict := &StrictLossError{Report: report}
	var strictTarget *StrictLossError
	if !errors.As(strict, &strictTarget) || strictTarget.Report.Losses[0].Code != LossCodeMetadataOmitted || !strings.Contains(strict.Error(), LossCodeMetadataOmitted) || !strings.Contains(strict.Error(), "/metadata/title") {
		t.Fatalf("strict error = %#v (%q)", strictTarget, strict.Error())
	}

	same := &SameFormatError{Format: "subrip"}
	var sameTarget *SameFormatError
	if !errors.As(same, &sameTarget) || !strings.Contains(same.Error(), "render") {
		t.Fatalf("same-format error = %#v (%q)", sameTarget, same.Error())
	}

	unsupported := &UnsupportedPairError{SourceFormat: "subrip", TargetFormat: "unknown"}
	var unsupportedTarget *UnsupportedPairError
	if !errors.As(unsupported, &unsupportedTarget) || !strings.Contains(unsupported.Error(), "unsupported") {
		t.Fatalf("unsupported-pair error = %#v (%q)", unsupportedTarget, unsupported.Error())
	}

	cause := errors.New("unsafe target payload")
	projection := &ProjectionError{SourceFormat: "subrip", TargetFormat: "webvtt", Path: "/cues/0/payload", Err: cause}
	var projectionTarget *ProjectionError
	if !errors.As(projection, &projectionTarget) || !errors.Is(projection, cause) || !strings.Contains(projection.Error(), "/cues/0/payload") {
		t.Fatalf("projection error = %#v (%q)", projectionTarget, projection.Error())
	}

	result := Result{Bytes: []byte("target\n"), LossReport: report, Diagnostics: []model.Diagnostic{}}
	if len(result.Bytes) == 0 || len(result.LossReport.Losses) != 1 || result.Diagnostics == nil {
		t.Fatalf("result = %#v", result)
	}
	_ = Options{Strict: true}
}

func newTestLoss(code string, kind Kind, path string, sourceOrder *int, cueID *string) Loss {
	return Loss{
		Code: code, Severity: SeverityWarning, Kind: kind, Message: "a source semantic is not represented",
		SourceFormat: "subrip", TargetFormat: "webvtt", Path: path, SourceOrder: sourceOrder, CueID: cueID, Context: []Attribute{},
	}
}

func reportTestDocument(t *testing.T) model.Document {
	t.Helper()
	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	return document
}
