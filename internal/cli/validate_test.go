package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
)

func TestValidateOptionHelpersCanonicalizeAndRejectConflicts(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		selector string
		want     string
	}{
		{selector: "auto", want: "auto"},
		{selector: "cueson", want: "cueson"},
		{selector: "json", want: "cueson"},
		{selector: "cue-json", want: "cueson"},
		{selector: "srt", want: "subrip"},
		{selector: "subrip", want: "subrip"},
		{selector: "vtt", want: "webvtt"},
		{selector: "webvtt", want: "webvtt"},
	} {
		options := validateOptions{input: "input", format: test.selector}
		if err := finalizeValidateOptions(&options); err != nil || options.format != test.want {
			t.Errorf("finalizeValidateOptions(%q) = %#v, %v", test.selector, options, err)
		}
	}

	for _, options := range []validateOptions{
		{},
		{input: "input", formatSet: true},
		{input: "input", encodingSet: true},
		{input: "input", format: "unknown"},
		{input: "input", format: "cueson", encoding: "utf-8"},
		{input: "input", format: "webvtt", encoding: "windows-1252"},
		{input: "input", encoding: "unknown"},
	} {
		if err := finalizeValidateOptions(&options); err == nil {
			t.Errorf("finalizeValidateOptions(%#v) error = nil", options)
		}
	}

	var options validateOptions
	if err := setValidateValue(&options, "--format", "srt"); err != nil {
		t.Fatal(err)
	}
	if err := setValidateValue(&options, "--format", "vtt"); err == nil || !strings.Contains(err.Error(), "only once") {
		t.Fatalf("duplicate --format error = %v", err)
	}
	if err := setValidateValue(&options, "--encoding", "utf-8"); err != nil {
		t.Fatal(err)
	}
	if err := setValidateValue(&options, "--encoding", "utf-8"); err == nil || !strings.Contains(err.Error(), "only once") {
		t.Fatalf("duplicate --encoding error = %v", err)
	}
	if err := setValidateOperand(&options, "first"); err != nil {
		t.Fatal(err)
	}
	if err := setValidateOperand(&options, "second"); err == nil || !strings.Contains(err.Error(), "additional") {
		t.Fatalf("extra operand error = %v", err)
	}
}

func TestRunValidateCueJSONAndNativeSuccessStreams(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	jsonPath := filepath.Join(directory, "document.json")
	if err := os.WriteFile(jsonPath, schema.Representative(), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name    string
		options validateOptions
		want    string
	}{
		{name: "auto Cue JSON", options: validateOptions{input: jsonPath, format: "auto"}, want: "valid cue_json.subrip"},
		{name: "explicit Cue JSON", options: validateOptions{input: jsonPath, format: "cueson"}, want: "valid cue_json.subrip"},
	} {
		t.Run(test.name, func(t *testing.T) {
			status, stderr := runValidateForTest(context.Background(), test.options, globalOptions{})
			if status != ExitSuccess || stderr != "success: "+test.want+"\n" {
				t.Fatalf("runValidate() = (%d, %q)", status, stderr)
			}
		})
	}

	nativePath := filepath.Join(directory, "captions.vtt")
	if err := os.WriteFile(nativePath, []byte("WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	status, stderr := runValidateForTest(context.Background(), validateOptions{input: nativePath, format: "auto"}, globalOptions{})
	if status != ExitSuccess || stderr != "success: valid native_subtitle.webvtt\n" {
		t.Fatalf("native runValidate() = (%d, %q)", status, stderr)
	}
	if _, err := os.Lstat(nativePath + ".cueson.json"); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("validation created output: %v", err)
	}
}

func TestRunValidateWarningFiltersAndNoOutput(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	path := filepath.Join(directory, "misleading.vtt")
	input := []byte("00:00:00,000 --> 00:00:01,000\nHello\n")
	if err := os.WriteFile(path, input, 0o600); err != nil {
		t.Fatal(err)
	}
	status, standard := runValidateForTest(context.Background(), validateOptions{input: path, format: "auto"}, globalOptions{})
	if status != ExitSuccess || !strings.Contains(standard, "format_extension_disagreement") || !strings.Contains(standard, "subrip_sequence_missing") || !strings.HasSuffix(standard, "success: valid native_subtitle.subrip\n") {
		t.Fatalf("standard validation = (%d, %q)", status, standard)
	}
	status, quiet := runValidateForTest(context.Background(), validateOptions{input: path, format: "auto"}, globalOptions{quiet: true})
	if status != ExitSuccess || strings.Contains(quiet, "success:") || !strings.Contains(quiet, "warning:") {
		t.Fatalf("quiet validation = (%d, %q)", status, quiet)
	}
	status, silent := runValidateForTest(context.Background(), validateOptions{input: path, format: "auto"}, globalOptions{silent: true})
	if status != ExitSuccess || silent != "" {
		t.Fatalf("silent validation = (%d, %q)", status, silent)
	}
	if got, err := os.ReadFile(path); err != nil || !bytes.Equal(got, input) {
		t.Fatalf("validation changed input: %q, %v", got, err)
	}
}

func TestRunValidateFailureStagesAndStatuses(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	missing := filepath.Join(directory, "missing.json")
	status, stderr := runValidateForTest(context.Background(), validateOptions{input: missing, format: "auto"}, globalOptions{})
	if status != ExitInvocation || !strings.Contains(stderr, "does not exist") {
		t.Fatalf("missing validation = (%d, %q)", status, stderr)
	}

	malformed := filepath.Join(directory, "malformed.json")
	if err := os.WriteFile(malformed, []byte("{"), 0o600); err != nil {
		t.Fatal(err)
	}
	status, stderr = runValidateForTest(context.Background(), validateOptions{input: malformed, format: "auto"}, globalOptions{})
	if status != ExitRuntimeFailure || !strings.Contains(stderr, "parse Cue JSON") {
		t.Fatalf("malformed validation = (%d, %q)", status, stderr)
	}

	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	document.Source.Assets[0].Size.Bytes++
	payload, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	corrupt := filepath.Join(directory, "corrupt.json")
	if err := os.WriteFile(corrupt, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stderr = runValidateForTest(context.Background(), validateOptions{input: corrupt, format: "auto"}, globalOptions{})
	if status != ExitRuntimeFailure || !strings.Contains(stderr, "decoded length") {
		t.Fatalf("integrity validation = (%d, %q)", status, stderr)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	status, stderr = runValidateForTest(ctx, validateOptions{input: malformed, format: "auto"}, globalOptions{})
	if status != ExitRuntimeFailure || !strings.Contains(stderr, "canceled") {
		t.Fatalf("canceled validation = (%d, %q)", status, stderr)
	}
}

func TestRunValidateRejectsMalformedUnsupportedOversizedAndNonRegularInputs(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	malformedNative := filepath.Join(directory, "malformed.srt")
	if err := os.WriteFile(malformedNative, []byte("not a subtitle\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	unsupported := filepath.Join(directory, "unsupported.bin")
	if err := os.WriteFile(unsupported, []byte("not a recognized format\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	oversized := filepath.Join(directory, "oversized.srt")
	file, err := os.OpenFile(oversized, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(source.MaxCaptureBytes + 1); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name       string
		options    validateOptions
		wantStatus int
		want       string
	}{
		{name: "explicit malformed native", options: validateOptions{input: malformedNative, format: "srt"}, wantStatus: ExitRuntimeFailure, want: "SubRip"},
		{name: "unsupported auto input", options: validateOptions{input: unsupported, format: "auto"}, wantStatus: ExitRuntimeFailure, want: "format is unknown"},
		{name: "oversized input", options: validateOptions{input: oversized, format: "auto"}, wantStatus: ExitInvocation, want: "input path or size precondition failed"},
		{name: "non-regular input", options: validateOptions{input: directory, format: "auto"}, wantStatus: ExitInvocation, want: "input path or size precondition failed"},
		{name: "unsafe basename", options: validateOptions{input: filepath.Join(directory, "PRIVATE:BAD.srt"), format: "auto"}, wantStatus: ExitInvocation, want: "input path or size precondition failed"},
	} {
		t.Run(test.name, func(t *testing.T) {
			status, stderr := runValidateForTest(context.Background(), test.options, globalOptions{})
			if status != test.wantStatus || !strings.Contains(stderr, test.want) {
				t.Fatalf("runValidate() = (%d, %q), want status %d containing %q", status, stderr, test.wantStatus, test.want)
			}
			if test.wantStatus == ExitInvocation && !strings.Contains(stderr, "Usage:") {
				t.Fatalf("runValidate() invocation stderr = %q, want usage", stderr)
			}
		})
	}
}

func runValidateForTest(ctx context.Context, options validateOptions, globals globalOptions) (int, string) {
	var stderr bytes.Buffer
	diagnostics := newDiagnosticWriter(&stderr, globals)
	return runValidate(ctx, options, &stderr, diagnostics, rootHelp), stderr.String()
}
