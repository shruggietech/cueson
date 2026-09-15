package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestScriptedInspectionCountsRetainedNativeOwners(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"ass", "ssa"} {
		for _, sample := range []struct {
			name string
			want inspectionScriptedSummary
		}{
			{"zero-dialogue", inspectionScriptedSummary{SectionCount: 3, RecordCount: 5, FormatDeclarationCount: 2, StyleCount: 1, EventCount: 1, CommentEventCount: 1}},
			{"empty-overlap-drawing", inspectionScriptedSummary{SectionCount: 3, RecordCount: 7, FormatDeclarationCount: 2, StyleCount: 1, EventCount: 3, DialogueEventCount: 3, OverrideTagCount: 2}},
			{"attachments", inspectionScriptedSummary{SectionCount: 5, RecordCount: 9, FormatDeclarationCount: 2, StyleCount: 1, EventCount: 1, DialogueEventCount: 1, AttachmentCount: 2}},
			{"malformed-comment", inspectionScriptedSummary{SectionCount: 3, RecordCount: 6, FormatDeclarationCount: 2, StyleCount: 1, EventCount: 2, DialogueEventCount: 1, CommentEventCount: 1, InvalidEventCount: 1}},
			{"unused-invalid-style", inspectionScriptedSummary{SectionCount: 3, RecordCount: 6, FormatDeclarationCount: 2, StyleCount: 2, InvalidStyleCount: 1, EventCount: 1, DialogueEventCount: 1}},
			{"native-order", inspectionScriptedSummary{SectionCount: 5, RecordCount: 14, FormatDeclarationCount: 3, StyleCount: 1, EventCount: 2, DialogueEventCount: 1, CommentEventCount: 1, UnknownRecordCount: 1}},
			{"karaoke-speaker", inspectionScriptedSummary{SectionCount: 3, RecordCount: 5, FormatDeclarationCount: 2, StyleCount: 1, EventCount: 1, DialogueEventCount: 1, OverrideTagCount: 4, KaraokeSpanCount: 4}},
			{"diagnosed-text", inspectionScriptedSummary{SectionCount: 3, RecordCount: 8, FormatDeclarationCount: 2, StyleCount: 1, EventCount: 4, DialogueEventCount: 4, OverrideTagCount: 4, KaraokeSpanCount: 1, UnsupportedKaraokeSpanCount: 1}},
		} {
			t.Run(format+"/"+sample.name, func(t *testing.T) {
				path := filepath.Join("..", "..", "testdata", "fixtures", "scripted", format+"-"+sample.name, "source", "captions."+format)
				status, encoded, stderr := runForTest(context.Background(), []string{"encode", path, "--stdout"})
				if status != ExitSuccess {
					t.Fatalf("encode fixture = %d %s", status, stderr)
				}
				cueJSON := filepath.Join(t.TempDir(), "misleading.srt")
				if err := os.WriteFile(cueJSON, []byte(encoded), 0600); err != nil {
					t.Fatal(err)
				}
				for _, input := range []string{path, cueJSON} {
					status, output, diagnostics := runForTest(context.Background(), []string{"inspect", input, "--json"})
					if status != ExitSuccess {
						t.Fatalf("inspect fixture = %d %s", status, diagnostics)
					}
					var report inspectionReport
					if err := json.Unmarshal([]byte(output), &report); err != nil {
						t.Fatal(err)
					}
					if report.Scripted == nil || *report.Scripted != sample.want {
						t.Fatalf("%s native owner counts = %+v, want %+v", input, report.Scripted, sample.want)
					}
					if report.Blocks.Total != 0 || report.Document.NonCueBlockCount != 0 || report.Document.BodyItemCount != report.Document.CueCount || report.Scripted.DialogueEventCount != report.Document.CueCount {
						t.Fatalf("native records redefined common/WebVTT counts: %+v", report)
					}
					if report.Capabilities.Declared.Status != "stable" || !report.Capabilities.Installed.Ingest || !report.Capabilities.Installed.Render {
						t.Fatalf("native capabilities = %+v", report.Capabilities)
					}
					assertInspectionJSONKeys(t, []byte(output))
				}
			})
		}
	}
}

func TestScriptedInspectionAllFifteenCountsAndPrivacyWithoutMutation(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"ass", "ssa"} {
		fixture, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			t.Fatal(err)
		}
		document, err := schema.Decode(fixture)
		if err != nil {
			t.Fatal(err)
		}
		native := document.FormatData.ASS
		if format == "ssa" {
			native = document.FormatData.SSA
		}
		// Test the projection independently of source grammar so every preserved
		// owner kind and prohibited content location is exercised precisely.
		const private = "PRIVATE-NATIVE-CONTENT"
		native.Sections[0].Name = private
		native.Records = append(native.Records, model.ScriptedRecord{Kind: "unknown", RawLine: pointer(private)}, model.ScriptedRecord{Kind: "malformed", RawLine: pointer(private)})
		native.Styles[0].Name = private
		native.Styles = append(native.Styles, model.ScriptedStyle{Valid: false, Name: private})
		native.Events[0].Text = private
		native.Events[0].Fields[0].RawValue = private
		native.Events[0].Tags = []model.ScriptedTag{{Name: private, Parameter: private, Raw: private}, {Name: private}}
		native.Events[0].Karaoke = []model.ScriptedKaraoke{{Variant: private, Supported: true}, {Variant: private, Supported: false}}
		native.Events = append(native.Events, model.ScriptedEvent{EventType: "comment", Valid: false, Text: private, Tags: []model.ScriptedTag{{Raw: private}}, Karaoke: []model.ScriptedKaraoke{{Supported: false}}})
		native.Attachments = []model.ScriptedAttachment{{Name: private}}
		before, _ := json.Marshal(document)
		report, err := buildInspectionReport(validatedInput{document: document, inputKind: inputKindCueJSON, selectionBasis: selectionBasisCueJSON})
		if err != nil {
			t.Fatal(err)
		}
		want := inspectionScriptedSummary{SectionCount: 3, RecordCount: 7, FormatDeclarationCount: 2, StyleCount: 2, InvalidStyleCount: 1, EventCount: 2, DialogueEventCount: 1, CommentEventCount: 1, InvalidEventCount: 1, AttachmentCount: 1, OverrideTagCount: 3, KaraokeSpanCount: 3, UnsupportedKaraokeSpanCount: 2, UnknownRecordCount: 1, MalformedRecordCount: 1}
		if report.Scripted == nil || *report.Scripted != want {
			t.Fatalf("all native aggregate counts = %+v, want %+v", report.Scripted, want)
		}
		if report.Capabilities.Declared.Status != "schema_only" || report.Capabilities.Declared.IngestSupported || !report.Capabilities.Installed.Ingest || !report.Capabilities.Installed.Render {
			t.Fatalf("declaration/installation conflated: %+v", report.Capabilities)
		}
		first, err := marshalInspectionReport(report)
		if err != nil {
			t.Fatal(err)
		}
		second, _ := marshalInspectionReport(report)
		human, _ := renderHumanInspection(report)
		after, _ := json.Marshal(document)
		if !bytes.Equal(before, after) || !bytes.Equal(first, second) || bytes.Count(first, []byte("\n")) != 1 || !bytes.HasSuffix(first, []byte("\n")) {
			t.Fatal("inspection mutated input or lost deterministic one-LF JSON framing")
		}
		for _, output := range [][]byte{first, human} {
			for _, prohibited := range []string{private, document.Source.Assets[0].DataBase64, document.Source.Assets[0].Hashes.SHA256} {
				if bytes.Contains(output, []byte(prohibited)) {
					t.Fatalf("scripted inspection exposed prohibited value %q", prohibited)
				}
			}
		}
		if !strings.Contains(string(human), "Scripted: sections=3 records=7 format_declarations=2") || !strings.Contains(string(human), "unsupported_karaoke_spans=2 unknown_records=1 malformed_records=1") {
			t.Fatalf("human scripted summary = %s", human)
		}
		var decoded map[string]json.RawMessage
		if err := json.Unmarshal(first, &decoded); err != nil {
			t.Fatal(err)
		}
		var counts map[string]int
		if err := json.Unmarshal(decoded["scripted"], &counts); err != nil || len(counts) != 15 {
			t.Fatalf("scripted JSON must contain exactly fifteen integer counts: %s, %v", decoded["scripted"], err)
		}
		assertInspectionJSONKeys(t, first)
	}
}

func TestScriptedInspectionEmptyOwnersAndMissingMatchingBranch(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"ass", "ssa"} {
		document := model.Document{Format: format}
		if _, err := projectInspectionScripted(document); err == nil {
			t.Fatal("missing native branch appeared to be a valid empty inventory")
		}
		native := &model.ScriptedDocumentData{}
		if format == "ass" {
			document.FormatData.ASS = native
		} else {
			document.FormatData.SSA = native
		}
		got, err := projectInspectionScripted(document)
		if err != nil || got == nil || *got != (inspectionScriptedSummary{}) {
			t.Fatalf("empty native owner aggregate = %+v, %v", got, err)
		}
	}
	for _, format := range []string{"subrip", "webvtt"} {
		got, err := projectInspectionScripted(model.Document{Format: format})
		if got != nil || err != nil {
			t.Fatalf("text format received scripted metadata: %+v, %v", got, err)
		}
	}
}

func TestTextInspectionBytesMatchPreS027NativeAndHistoricalSnapshots(t *testing.T) {
	t.Parallel()
	for _, sample := range []struct{ name, path string }{
		{"subrip", "../../testdata/fixtures/conversion/srt-loss-free/source/input.srt"},
		{"webvtt", "../../testdata/fixtures/conversion/webvtt-loss-free/source/input.vtt"},
		{"historical-subrip", "../schema/testdata/historical-v1.0.0-subrip.json"},
		{"historical-webvtt", "../schema/testdata/historical-v1.0.0-webvtt.json"},
	} {
		for _, mode := range []string{"json", "human"} {
			args := []string{"inspect", sample.path}
			if mode == "json" {
				args = append(args, "--json")
			}
			status, output, stderr := runForTest(context.Background(), args)
			if status != ExitSuccess {
				t.Fatalf("%s %s report = %d %s", sample.name, mode, status, stderr)
			}
			want, err := os.ReadFile(filepath.Join("testdata", "inspection", sample.name+"."+mode+".golden"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal([]byte(output), want) {
				t.Fatalf("%s %s inspection changed from exact S026 snapshot", sample.name, mode)
			}
			if strings.Contains(output, "\"scripted\"") || strings.Contains(output, "Scripted:") {
				t.Fatal("old text report received scripted field or human section")
			}
		}
	}
}

func TestScriptedEncodingRetainsOwningCommandExitClasses(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"ass", "ssa"} {
		path := "../../testdata/fixtures/scripted/" + format + "-basic/source/captions." + format
		for _, command := range []string{"encode", "validate", "inspect", "convert"} {
			args := []string{command, path, "--encoding", "windows-1252"}
			want := ExitRuntimeFailure
			if command == "convert" {
				args = append(args, "--from", format, "--to", "srt")
				want = ExitInvocation
			} else {
				args = append(args, "--format", format)
				if command == "encode" {
					args = append(args, "--stdout")
				}
			}
			status, stdout, stderr := runForTest(context.Background(), args)
			if status != want || stdout != "" || stderr == "" {
				t.Fatalf("%s %s rejected encoding changed class: %d %q %q, want %d", format, command, status, stdout, stderr, want)
			}
		}
	}
}

func pointer[T any](value T) *T { return &value }

// Assert the aggregate contains no strings, slices, pointers or native owners.
func TestScriptedInspectionAggregateTypeIsCountsOnly(t *testing.T) {
	t.Parallel()
	typeInfo := reflect.TypeFor[inspectionScriptedSummary]()
	if typeInfo.NumField() != 15 {
		t.Fatalf("summary has %d fields, want fifteen", typeInfo.NumField())
	}
	for index := range typeInfo.NumField() {
		if field := typeInfo.Field(index); field.Type.Kind() != reflect.Int {
			t.Errorf("inspection native field %s exposes non-count type %s", field.Name, field.Type)
		}
	}
}
