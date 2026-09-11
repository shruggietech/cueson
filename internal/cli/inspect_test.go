package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/source"
)

func TestBuildInspectionReportProjectsSafeDeterministicStructure(t *testing.T) {
	t.Parallel()

	input := inspectTestWebVTTInput(t)
	report, err := buildInspectionReport(input)
	if err != nil {
		t.Fatal(err)
	}

	if report.InspectReportVersion != "1" || report.Input.Kind != "native_subtitle" || report.Input.SelectionBasis != "content" {
		t.Fatalf("report identity = %#v", report)
	}
	if report.Format != "webvtt" || !report.Schema.Compatible {
		t.Fatalf("format/schema = (%q, %#v)", report.Format, report.Schema)
	}
	if !report.Capabilities.Installed.Ingest || !report.Capabilities.Installed.Render || !report.Capabilities.Installed.Restore || !report.Capabilities.Installed.Validate || !report.Capabilities.Installed.Inspect {
		t.Fatalf("installed capabilities = %#v", report.Capabilities.Installed)
	}
	if report.Integrity.Status != "verified" || report.Integrity.AssetCount != 1 || report.Integrity.TotalBytes <= 0 {
		t.Fatalf("integrity = %#v", report.Integrity)
	}
	if len(report.Assets) != 1 || !report.Assets[0].Primary || report.Assets[0].FileName != "captions.vtt" || report.Assets[0].Encoding == nil {
		t.Fatalf("assets = %#v", report.Assets)
	}
	if report.Document.CueCount != 1 || report.Document.NonCueBlockCount != 2 || report.Document.BodyItemCount != 3 || report.Document.WarningCount != 1 {
		t.Fatalf("document = %#v", report.Document)
	}
	if len(report.Cues) != 1 || report.Cues[0].PayloadLineCount != 1 || report.Cues[0].SpeakerCount != 1 || report.Cues[0].NativeSettingCount != 1 || report.Cues[0].NativeSettingOccurrenceCount != 2 || report.Cues[0].InvalidNativeSettingCount != 1 {
		t.Fatalf("cues = %#v", report.Cues)
	}
	if report.Blocks.Total != 2 || report.Blocks.NoteCount != 1 || report.Blocks.RegionCount != 1 || len(report.Blocks.Entries) != 2 || report.Blocks.Entries[1].NativeSettingOccurrenceCount != 2 || report.Blocks.Entries[1].InvalidNativeSettingCount != 1 {
		t.Fatalf("blocks = %#v", report.Blocks)
	}
	if len(report.Diagnostics) != 1 || report.Diagnostics[0].Code != "privacy_sentinel" || report.Diagnostics[0].CueOrdinal == nil || *report.Diagnostics[0].CueOrdinal != 0 {
		t.Fatalf("diagnostics = %#v", report.Diagnostics)
	}
	if report.Loss.Status != "not_evaluated" || report.Loss.Reason != "target_format_required" {
		t.Fatalf("loss = %#v", report.Loss)
	}

	first, err := marshalInspectionReport(report)
	if err != nil {
		t.Fatal(err)
	}
	second, err := marshalInspectionReport(report)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) || len(first) == 0 || first[len(first)-1] != '\n' || bytes.Count(first, []byte("\n")) != 1 {
		t.Fatalf("JSON output is not one repeated deterministic LF record: %q", first)
	}
	previous := -1
	for _, key := range []string{"inspect_report_version", "input", "format", "schema", "capabilities", "integrity", "assets", "document", "cues", "blocks", "diagnostics", "loss"} {
		position := bytes.Index(first, []byte(`"`+key+`":`))
		if position <= previous {
			t.Fatalf("root key %q is missing or out of contract order in %s", key, first)
		}
		previous = position
	}
	assertInspectionJSONKeys(t, first)
	assertInspectionPrivacy(t, first)
}

func TestBuildInspectionReportOrdersCuesAndBlocksWithoutMutatingDocument(t *testing.T) {
	t.Parallel()

	input := inspectTestWebVTTInput(t)
	firstCue := input.document.Cues[0]
	secondCue := firstCue
	secondCue.ID = "cue-000001"
	secondCue.Ordinal = 1
	secondCue.SourceOrder = 3
	input.document.Cues = []model.Cue{secondCue, firstCue}
	firstBlock := input.document.FormatData.WebVTT.Blocks[0]
	secondBlock := input.document.FormatData.WebVTT.Blocks[1]
	input.document.FormatData.WebVTT.Blocks = []model.WebVTTBlock{secondBlock, firstBlock}

	report, err := buildInspectionReport(input)
	if err != nil {
		t.Fatal(err)
	}
	if report.Cues[0].Ordinal != 0 || report.Cues[1].Ordinal != 1 {
		t.Fatalf("cue order = %#v", report.Cues)
	}
	if report.Blocks.Entries[0].SourceOrder >= report.Blocks.Entries[1].SourceOrder {
		t.Fatalf("block order = %#v", report.Blocks.Entries)
	}
	if input.document.Cues[0].Ordinal != 1 || input.document.FormatData.WebVTT.Blocks[0].SourceOrder != secondBlock.SourceOrder {
		t.Fatal("report construction mutated the validated document")
	}
}

func TestBuildInspectionReportRetainsAssetAndDiagnosticOrder(t *testing.T) {
	t.Parallel()

	input := inspectTestWebVTTInput(t)
	companion := input.document.Source.Assets[0]
	companion.ID = "asset-1"
	companion.Role = "companion"
	companion.FileName = "first-in-stored-order.vtt"
	input.document.Source.Assets = []model.SourceAsset{companion, input.document.Source.Assets[0]}
	firstOrder := input.document.Cues[0].SourceOrder
	secondOrder := firstOrder
	input.document.Diagnostics = []model.Diagnostic{
		{Severity: "warning", Code: "z_code_stays_first", Message: "not reported", SourceOrder: &firstOrder},
		{Severity: "info", Code: "a_code_stays_second", Message: "not reported", SourceOrder: &secondOrder},
	}
	updateDiagnosticStats(&input.document)

	report, err := buildInspectionReport(input)
	if err != nil {
		t.Fatal(err)
	}
	if report.Assets[0].FileName != "first-in-stored-order.vtt" || report.Assets[0].Index != 0 || report.Assets[1].Index != 1 {
		t.Fatalf("asset order = %#v", report.Assets)
	}
	if report.Diagnostics[0].Code != "z_code_stays_first" || report.Diagnostics[1].Code != "a_code_stays_second" {
		t.Fatalf("diagnostic order = %#v", report.Diagnostics)
	}
	if report.Document.InfoCount != 1 || report.Document.WarningCount != 1 || report.Document.ErrorCount != 0 {
		t.Fatalf("diagnostic totals = %#v", report.Document)
	}
}

func TestBuildInspectionReportRejectsByteAndBodyCountOverflow(t *testing.T) {
	t.Parallel()

	input := inspectTestWebVTTInput(t)
	overflowAsset := input.document.Source.Assets[0]
	overflowAsset.FileName = "companion.vtt"
	overflowAsset.ID = "asset-1"
	input.document.Source.Assets[0].Size.Bytes = math.MaxInt64
	overflowAsset.Size.Bytes = 1
	input.document.Source.Assets = append(input.document.Source.Assets, overflowAsset)
	if _, err := buildInspectionReport(input); err == nil || !strings.Contains(err.Error(), "total_bytes") {
		t.Fatalf("buildInspectionReport() error = %v, want total_bytes overflow", err)
	}
}

func TestHumanInspectionUsesOnlyReportFacts(t *testing.T) {
	t.Parallel()

	report, err := buildInspectionReport(inspectTestWebVTTInput(t))
	if err != nil {
		t.Fatal(err)
	}
	output, err := renderHumanInspection(report)
	if err != nil {
		t.Fatal(err)
	}
	wantFragments := []string{
		"Inspection report 1\n",
		"Input: native_subtitle (content)",
		"Format: webvtt",
		"Integrity: verified (1 asset,",
		"Cues (1):",
		"Blocks (2): note=1 style=0 region=1 unknown=0",
		"Diagnostics (1): info=0 warning=1 error=0",
		"Loss: not_evaluated (target_format_required)",
	}
	for _, fragment := range wantFragments {
		if !bytes.Contains(output, []byte(fragment)) {
			t.Errorf("human output missing %q:\n%s", fragment, output)
		}
	}
	assertInspectionPrivacy(t, output)
}

func TestInspectionPrivacyExcludesContentAndUserControlledIdentifiers(t *testing.T) {
	t.Parallel()

	input := inspectTestWebVTTInput(t)
	const sentinel = "PRIVATE-CONTENT-C:\\Users\\operator\\captions"
	identifier := sentinel
	input.document.Cues[0].SourceIdentifier = &identifier
	input.document.Cues[0].FormatData.WebVTT.IdentifierRaw = &identifier
	input.document.Cues[0].Payload = model.Payload{RawText: sentinel, PlainText: sentinel, Lines: []string{sentinel}}
	input.document.Cues[0].FormatData.WebVTT.RawPayload = sentinel
	input.document.Cues[0].FormatData.WebVTT.RawPayloadLines = []string{sentinel}
	input.document.Cues[0].Speakers = []model.Speaker{{Name: sentinel, Origin: "native"}}
	input.document.Cues[0].Tokens = []model.Token{{Text: sentinel, StartMilliseconds: 0, EndMilliseconds: 1000}}
	input.document.Diagnostics[0].Message = sentinel

	report, err := buildInspectionReport(input)
	if err != nil {
		t.Fatal(err)
	}
	jsonOutput, err := marshalInspectionReport(report)
	if err != nil {
		t.Fatal(err)
	}
	humanOutput, err := renderHumanInspection(report)
	if err != nil {
		t.Fatal(err)
	}
	for _, output := range [][]byte{jsonOutput, humanOutput} {
		if bytes.Contains(output, []byte(sentinel)) || bytes.Contains(output, []byte(input.document.Source.Assets[0].DataBase64)) || bytes.Contains(output, []byte(input.document.Source.Assets[0].Hashes.SHA256)) {
			t.Fatalf("inspection leaked prohibited source content: %s", output)
		}
	}
}

func TestInspectOptionValidation(t *testing.T) {
	t.Parallel()

	options := inspectOptions{input: "captions.vtt"}
	if err := setInspectValue(&options, "--format", "webvtt"); err != nil {
		t.Fatal(err)
	}
	if err := setInspectValue(&options, "--encoding", "utf8"); err != nil {
		t.Fatal(err)
	}
	if err := finalizeInspectOptions(&options); err != nil {
		t.Fatal(err)
	}
	if options.format != "webvtt" {
		t.Fatalf("canonical format = %q, want webvtt", options.format)
	}
	if err := setInspectValue(&options, "--format", "srt"); err == nil || !strings.Contains(err.Error(), "only once") {
		t.Fatalf("duplicate format error = %v", err)
	}
	if err := setInspectJSON(&options); err != nil {
		t.Fatal(err)
	}
	if err := setInspectJSON(&options); err == nil || !strings.Contains(err.Error(), "only once") {
		t.Fatalf("duplicate JSON error = %v", err)
	}

	for _, test := range []struct {
		name    string
		options inspectOptions
		want    string
	}{
		{name: "missing input", options: inspectOptions{}, want: "requires one INPUT"},
		{name: "explicit empty format", options: inspectOptions{input: "x", formatSet: true}, want: "--format requires a value"},
		{name: "explicit empty encoding", options: inspectOptions{input: "x", encodingSet: true}, want: "--encoding requires a value"},
		{name: "unknown format", options: inspectOptions{input: "x", format: "mystery"}, want: "must be auto"},
		{name: "unknown encoding", options: inspectOptions{input: "x", encoding: "mystery"}, want: "not supported"},
		{name: "Cue JSON encoding", options: inspectOptions{input: "x", format: "json", encoding: "utf-8"}, want: "cannot be used"},
		{name: "WebVTT legacy encoding", options: inspectOptions{input: "x", format: "vtt", encoding: "windows-1252"}, want: "requires UTF-8"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := finalizeInspectOptions(&test.options); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("finalizeInspectOptions() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestRunInspectCueJSONReportsCueJSONClassification(t *testing.T) {
	t.Parallel()

	input := inspectTestWebVTTInput(t)
	payload, err := marshalDocument(input.document, false)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "document.cueson.json")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	diagnostics := newDiagnosticWriter(&stderr, globalOptions{})
	status := runInspect(context.Background(), inspectOptions{input: path, format: "auto", json: true}, &stdout, &stderr, diagnostics, noHelp)
	if status != ExitSuccess {
		t.Fatalf("runInspect() status = %d, stderr = %q", status, stderr.String())
	}
	var report inspectionReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.Input.Kind != "cue_json" || report.Input.SelectionBasis != "cue_json" || report.Input.ContentFormat != nil || report.Input.ExtensionFormat != nil {
		t.Fatalf("Cue JSON classification = %#v", report.Input)
	}
	if strings.Contains(stdout.String(), "PRIVATE-DIAGNOSTIC-MESSAGE") || strings.Contains(stderr.String(), "PRIVATE-DIAGNOSTIC-MESSAGE") {
		t.Fatalf("Cue JSON diagnostic message leaked: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "warning: privacy_sentinel") {
		t.Fatalf("stderr = %q, want stable warning code", stderr.String())
	}
}

func TestRunInspectJSONKeepsPayloadPureAndWarningsSafe(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "misnamed.srt")
	input := []byte("WEBVTT\n\n00:00.000 --> 00:01.000 mystery:value\nPRIVATE-PAYLOAD\n")
	if err := os.WriteFile(path, input, 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	diagnostics := newDiagnosticWriter(&stderr, globalOptions{})
	status := runInspect(context.Background(), inspectOptions{input: path, format: "auto", json: true}, &stdout, &stderr, diagnostics, noHelp)
	if status != ExitSuccess {
		t.Fatalf("runInspect() status = %d, stderr = %q", status, stderr.String())
	}
	var report inspectionReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatalf("stdout is not one JSON report: %v\n%s", err, stdout.String())
	}
	if report.Input.SelectionBasis != "content" || report.Format != "webvtt" {
		t.Fatalf("classification = %#v", report.Input)
	}
	if strings.Contains(stdout.String(), "PRIVATE-PAYLOAD") || strings.Contains(stderr.String(), "PRIVATE-PAYLOAD") {
		t.Fatalf("inspection leaked payload: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "warning: format_extension_disagreement") || !strings.Contains(stderr.String(), "warning: webvtt_setting_unknown") {
		t.Fatalf("stderr = %q, want stable warning codes", stderr.String())
	}
}

func TestRunInspectStdoutFailureIsRuntimeFailure(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "captions.srt")
	if err := os.WriteFile(path, []byte("1\n00:00:00,000 --> 00:00:01,000\ncaption\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	diagnostics := newDiagnosticWriter(&stderr, globalOptions{})
	status := runInspect(context.Background(), inspectOptions{input: path, format: "auto", json: true}, inspectErrorWriter{}, &stderr, diagnostics, noHelp)
	if status != ExitRuntimeFailure || !strings.Contains(stderr.String(), "inspect: stdout write failed") || strings.Contains(stderr.String(), os.ErrPermission.Error()) {
		t.Fatalf("runInspect() = %d, stderr = %q", status, stderr.String())
	}
}

func TestRunInspectSanitizesFailureDiagnostics(t *testing.T) {
	t.Parallel()

	t.Run("source integrity", func(t *testing.T) {
		t.Parallel()
		input := inspectTestWebVTTInput(t)
		privateAssetID := "PRIVATE-ASSET-ID"
		privateHash := strings.Repeat("a", 64)
		input.document.Source.Assets[0].ID = privateAssetID
		input.document.Source.PrimaryAssetID = privateAssetID
		input.document.Source.Assets[0].Hashes.SHA256 = privateHash
		payload, err := marshalDocument(input.document, false)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "PRIVATE-CALLER-PATH.cueson.json")
		if err := os.WriteFile(path, payload, 0o644); err != nil {
			t.Fatal(err)
		}

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		diagnostics := newDiagnosticWriter(&stderr, globalOptions{})
		status := runInspect(context.Background(), inspectOptions{input: path, format: "auto", json: true}, &stdout, &stderr, diagnostics, noHelp)
		if status != ExitRuntimeFailure || stdout.Len() != 0 || !strings.Contains(stderr.String(), "inspect: input validation failed") {
			t.Fatalf("runInspect() = %d, stdout = %q, stderr = %q", status, stdout.String(), stderr.String())
		}
		for _, prohibited := range []string{privateAssetID, privateHash, path, filepath.Base(path), "SHA-256"} {
			if strings.Contains(stderr.String(), prohibited) {
				t.Fatalf("inspection failure leaked %q in stderr %q", prohibited, stderr.String())
			}
		}
	})

	t.Run("missing caller path", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(t.TempDir(), "PRIVATE-MISSING-PATH.srt")
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		diagnostics := newDiagnosticWriter(&stderr, globalOptions{})
		status := runInspect(context.Background(), inspectOptions{input: path, format: "auto", json: true}, &stdout, &stderr, diagnostics, inspectHelp)
		if status != ExitInvocation || stdout.Len() != 0 || !strings.Contains(stderr.String(), "inspect: input does not exist") {
			t.Fatalf("runInspect() = %d, stdout = %q, stderr = %q", status, stdout.String(), stderr.String())
		}
		if strings.Contains(stderr.String(), path) || strings.Contains(stderr.String(), filepath.Base(path)) {
			t.Fatalf("missing-input diagnostic leaked caller path in stderr %q", stderr.String())
		}
	})
}

type inspectErrorWriter struct{}

func (inspectErrorWriter) Write([]byte) (int, error) {
	return 0, os.ErrPermission
}

func inspectTestWebVTTInput(t *testing.T) validatedInput {
	t.Helper()

	directory := t.TempDir()
	path := filepath.Join(directory, "captions.vtt")
	// Supply a valid cue after the two blocks while retaining one warning whose
	// free-form message must never enter report output.
	content := "WEBVTT\n\nNOTE safe structural count only\nline\n\nREGION\nid:region-1 mystery:value\n\n00:00.000 --> 00:01.000 align:start mystery:value\n<v PRIVATE-SPEAKER>hello</v>\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	captured, err := source.CaptureContext(context.Background(), path, source.CaptureOptions{})
	if err != nil {
		t.Fatal(err)
	}
	document, err := decodeCapturedSource(context.Background(), captured, "vtt", "", false)
	if err != nil {
		t.Fatal(err)
	}
	order := document.Cues[0].SourceOrder
	id := document.Cues[0].ID
	document.Diagnostics = []model.Diagnostic{{Severity: "warning", Code: "privacy_sentinel", Message: "PRIVATE-DIAGNOSTIC-MESSAGE", SourceOrder: &order, CueID: &id}}
	updateDiagnosticStats(&document)
	contentFormat := codec.FormatWebVTT
	extensionFormat := codec.FormatWebVTT
	return validatedInput{
		document:        document,
		inputKind:       inputKind("native_subtitle"),
		selectionBasis:  selectionBasis("content"),
		contentFormat:   &contentFormat,
		extensionFormat: &extensionFormat,
		diagnostics:     document.Diagnostics,
	}
}

func assertInspectionPrivacy(t *testing.T, output []byte) {
	t.Helper()
	for _, prohibited := range []string{
		"data_base64", "sha256", "PRIVATE-SPEAKER", "PRIVATE-DIAGNOSTIC-MESSAGE",
		"NOTE safe structural count only", "id:region-1", "region-1", "hello",
	} {
		if bytes.Contains(output, []byte(prohibited)) {
			t.Errorf("inspection output contains prohibited %q: %s", prohibited, output)
		}
	}
}

func assertInspectionJSONKeys(t *testing.T, payload []byte) {
	t.Helper()
	var decoded any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatal(err)
	}
	keyPattern := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	var visit func(any)
	visit = func(value any) {
		switch typed := value.(type) {
		case map[string]any:
			for key, child := range typed {
				if !keyPattern.MatchString(key) {
					t.Errorf("JSON key %q is not lowercase snake_case", key)
				}
				visit(child)
			}
		case []any:
			for _, child := range typed {
				visit(child)
			}
		}
	}
	visit(decoded)
}
