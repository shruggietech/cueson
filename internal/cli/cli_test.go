package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/schema"
)

func TestRunVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "plain", args: []string{"version"}},
		{name: "quiet before", args: []string{"--quiet", "version"}},
		{name: "quiet short after", args: []string{"version", "-q"}},
		{name: "silent after", args: []string{"version", "--silent"}},
		{name: "no color before", args: []string{"--no-color", "version"}},
		{name: "repeated globals", args: []string{"--quiet", "version", "--quiet", "--silent", "--no-color"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			status, stdout, stderr := runForTest(context.Background(), tt.args)
			if status != ExitSuccess {
				t.Fatalf("Run() status = %d, want %d; stderr = %q", status, ExitSuccess, stderr)
			}
			if stdout != "0.1.0\n" {
				t.Errorf("Run() stdout = %q, want %q", stdout, "0.1.0\\n")
			}
			if stderr != "" {
				t.Errorf("Run() stderr = %q, want empty", stderr)
			}
		})
	}
}

func TestRunCanceledVersion(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	status, stdout, stderr := runForTest(ctx, []string{"version", "--silent"})
	if status != ExitRuntimeFailure {
		t.Errorf("Run() status = %d, want %d", status, ExitRuntimeFailure)
	}
	if stdout != "" {
		t.Errorf("Run() stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "canceled") {
		t.Errorf("Run() stderr = %q, want cancellation diagnostic", stderr)
	}
}

func TestRunRootHelp(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{nil, {"-h"}, {"--help"}, {"--quiet", "--help"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess {
			t.Errorf("Run(%q) status = %d, want %d; stderr = %q", args, status, ExitSuccess, stderr)
		}
		if stderr != "" {
			t.Errorf("Run(%q) stderr = %q, want empty", args, stderr)
		}
		for _, want := range []string{"Usage:", "Commands:", "version", "Global options:", "Examples:", "cueson version"} {
			if !strings.Contains(stdout, want) {
				t.Errorf("Run(%q) stdout does not contain %q:\n%s", args, want, stdout)
			}
		}
		if !strings.Contains(stdout, "\n  schema") {
			t.Errorf("Run(%q) help does not expose schema command:\n%s", args, stdout)
		}
		if !strings.Contains(stdout, "\n  restore") {
			t.Errorf("Run(%q) help does not expose restore command:\n%s", args, stdout)
		}
		for _, deferred := range []string{"encode", "render", "convert", "validate", "inspect", "completion"} {
			if strings.Contains(stdout, "\n  "+deferred) {
				t.Errorf("Run(%q) help exposes deferred command %q:\n%s", args, deferred, stdout)
			}
		}
	}
}

func TestRunRestore(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	input := filepath.Join(directory, "document.cueson.json")
	output := filepath.Join(directory, "restored.srt")
	if err := os.WriteFile(input, schema.Representative(), 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", output, input})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("restore = (%d, %q, %q), want silent success", status, stdout, stderr)
	}
	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	want, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(output); err != nil || !bytes.Equal(got, want) {
		t.Fatalf("restored bytes = %q, error = %v", got, err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"--silent", "restore", "--no-metadata", "--output", output, input})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "already exists") {
		t.Fatalf("refused restore = (%d, %q, %q)", status, stdout, stderr)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", output, "--force", input})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("forced restore = (%d, %q, %q)", status, stdout, stderr)
	}
}

func TestRunRestoreHelpAndInvalidInvocation(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"restore", "--help"}, {"restore", "-h"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess || stderr != "" || !strings.Contains(stdout, "--strict-metadata") || !strings.Contains(stdout, "--output-dir") {
			t.Errorf("Run(%q) = (%d, %q, %q), want restore help", args, status, stdout, stderr)
		}
	}
	tests := []struct {
		args []string
		want string
	}{
		{args: []string{"restore"}, want: "requires one INPUT"},
		{args: []string{"restore", "--output", "x", "--output-dir", "y", "in"}, want: "mutually exclusive"},
		{args: []string{"restore", "--strict-metadata", "--no-metadata", "in"}, want: "mutually exclusive"},
		{args: []string{"restore", "a", "b"}, want: "no additional arguments"},
		{args: []string{"restore", "--bogus", "in"}, want: "unknown option"},
	}
	for _, tt := range tests {
		status, stdout, stderr := runForTest(context.Background(), tt.args)
		if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, tt.want) || !strings.Contains(stderr, "cueson [global options] restore") {
			t.Errorf("Run(%q) = (%d, %q, %q), want invocation error containing %q", tt.args, status, stdout, stderr, tt.want)
		}
	}
}

func TestRunRestorePreExecutionPathFailures(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	input := filepath.Join(directory, "document.cueson.json")
	if err := os.WriteFile(input, schema.Representative(), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"restore", filepath.Join(directory, "missing.cueson.json")},
		{"restore", directory},
		{"restore", "--no-metadata", "--output", "   ", input},
		{"restore", "--no-metadata", "--output-dir", "\t", input},
		{"restore", "--no-metadata", "--output", filepath.Join(directory, "missing", "out.srt"), input},
		{"restore", "--no-metadata", "--output-dir", filepath.Join(directory, "missing"), input},
	} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "cueson [global options] restore") {
			t.Errorf("Run(%q) = (%d, %q, %q), want pre-execution failure", args, status, stdout, stderr)
		}
	}
}

func TestRunVersionHelp(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"version", "-h"}, {"version", "--help"}, {"--quiet", "version", "--help"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess {
			t.Errorf("Run(%q) status = %d, want %d; stderr = %q", args, status, ExitSuccess, stderr)
		}
		if stderr != "" {
			t.Errorf("Run(%q) stderr = %q, want empty", args, stderr)
		}
		for _, want := range []string{"Usage:", "cueson version", "Print the Cueson executable version", "Example:"} {
			if !strings.Contains(stdout, want) {
				t.Errorf("Run(%q) stdout does not contain %q:\n%s", args, want, stdout)
			}
		}
	}
}

func TestRunInvalidInvocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		wantError string
		wantUsage string
	}{
		{name: "unknown command", args: []string{"missing"}, wantError: `unknown command "missing"`, wantUsage: "cueson [global options] <command>"},
		{name: "unknown root option", args: []string{"--force"}, wantError: `unknown option "--force"`, wantUsage: "cueson [global options] <command>"},
		{name: "unknown version option", args: []string{"version", "--bogus"}, wantError: `unknown option "--bogus"`, wantUsage: "cueson [global options] version"},
		{name: "version operand", args: []string{"version", "unexpected"}, wantError: "version accepts no arguments", wantUsage: "cueson [global options] version"},
		{name: "literal option as command", args: []string{"--", "--help"}, wantError: `unknown command "--help"`, wantUsage: "cueson [global options] <command>"},
		{name: "literal option as operand", args: []string{"version", "--", "--help"}, wantError: "version accepts no arguments", wantUsage: "cueson [global options] version"},
		{name: "literal tilde", args: []string{"~"}, wantError: `unknown command "~"`, wantUsage: "cueson [global options] <command>"},
		{name: "literal environment syntax", args: []string{"$HOME"}, wantError: `unknown command "$HOME"`, wantUsage: "cueson [global options] <command>"},
		{name: "silent retains error", args: []string{"--silent", "missing"}, wantError: `unknown command "missing"`, wantUsage: "cueson [global options] <command>"},
		{name: "quiet retains error", args: []string{"--quiet", "missing"}, wantError: `unknown command "missing"`, wantUsage: "cueson [global options] <command>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			status, stdout, stderr := runForTest(context.Background(), tt.args)
			if status != ExitInvocation {
				t.Errorf("Run() status = %d, want %d", status, ExitInvocation)
			}
			if stdout != "" {
				t.Errorf("Run() stdout = %q, want empty", stdout)
			}
			if !strings.Contains(stderr, tt.wantError) {
				t.Errorf("Run() stderr = %q, want error containing %q", stderr, tt.wantError)
			}
			if !strings.Contains(stderr, tt.wantUsage) {
				t.Errorf("Run() stderr = %q, want usage containing %q", stderr, tt.wantUsage)
			}
			if strings.Contains(stderr, "\x1b[") {
				t.Errorf("Run() non-terminal stderr contains terminal control sequence: %q", stderr)
			}
		})
	}
}

func TestRunSchema(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"schema"}, {"--quiet", "schema"}, {"schema", "--silent"}, {"schema", "--no-color"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess {
			t.Errorf("Run(%q) status = %d, want %d; stderr = %q", args, status, ExitSuccess, stderr)
		}
		if stdout != string(schema.Bytes()) {
			t.Errorf("Run(%q) schema output does not match embedded bytes", args)
		}
		if stderr != "" {
			t.Errorf("Run(%q) stderr = %q, want empty", args, stderr)
		}
	}
}

func TestRunSchemaVersion(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"schema", "--version"}, {"--quiet", "schema", "--version"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess || stdout != "0.1.0\n" || stderr != "" {
			t.Errorf("Run(%q) = (%d, %q, %q), want (0, %q, empty)", args, status, stdout, stderr, "0.1.0\\n")
		}
	}
}

func TestRunSchemaHelp(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"schema", "-h"}, {"schema", "--help"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess || stderr != "" {
			t.Errorf("Run(%q) = (%d, _, %q), want success and empty stderr", args, status, stderr)
		}
		for _, want := range []string{"cueson schema", "--version", "--output PATH", "--force", "Examples:"} {
			if !strings.Contains(stdout, want) {
				t.Errorf("Run(%q) help does not contain %q:\n%s", args, want, stdout)
			}
		}
	}
}

func TestRunSchemaFileOutput(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "cueson.schema.json")
	status, stdout, stderr := runForTest(context.Background(), []string{"schema", "--output", output})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("initial output = (%d, %q, %q), want success with silent streams", status, stdout, stderr)
	}
	assertFileEqualsSchema(t, output)

	before := []byte("preserve me")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"schema", "-o", output})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "already exists") {
		t.Errorf("refusal = (%d, %q, %q), want invocation failure", status, stdout, stderr)
	}
	if got, err := os.ReadFile(output); err != nil || !bytes.Equal(got, before) {
		t.Fatalf("refused output changed file: bytes = %q, error = %v", got, err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"schema", "--output=" + output, "--force"})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("forced output = (%d, %q, %q), want success with silent streams", status, stdout, stderr)
	}
	assertFileEqualsSchema(t, output)
}

func TestRunSchemaInvalidInvocation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "missing output value", args: []string{"schema", "--output"}, want: "requires a path"},
		{name: "empty output value", args: []string{"schema", "--output="}, want: "requires a path"},
		{name: "force without output", args: []string{"schema", "--force"}, want: "requires --output"},
		{name: "version with output", args: []string{"schema", "--version", "--output", "x"}, want: "cannot be combined"},
		{name: "version with force", args: []string{"schema", "--version", "--force"}, want: "cannot be combined"},
		{name: "extra operand", args: []string{"schema", "extra"}, want: "accepts no arguments"},
		{name: "literal option", args: []string{"schema", "--", "--version"}, want: "accepts no arguments"},
		{name: "unknown option", args: []string{"schema", "--bogus"}, want: "unknown option"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, stdout, stderr := runForTest(context.Background(), tt.args)
			if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, tt.want) || !strings.Contains(stderr, "cueson [global options] schema") {
				t.Errorf("Run(%q) = (%d, %q, %q), want invocation failure containing %q and schema usage", tt.args, status, stdout, stderr, tt.want)
			}
		})
	}
}

func TestRunSchemaRuntimeFailures(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	status, stdout, stderr := runForTest(ctx, []string{"schema"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "canceled") {
		t.Errorf("canceled schema = (%d, %q, %q), want runtime cancellation", status, stdout, stderr)
	}

	missingParent := filepath.Join(t.TempDir(), "missing", "schema.json")
	status, stdout, stderr = runForTest(context.Background(), []string{"schema", "--output", missingParent})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "write schema") {
		t.Errorf("bad output = (%d, %q, %q), want runtime I/O failure", status, stdout, stderr)
	}
}

func TestReplaceSchemaFileCommitFailurePreservesDestination(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "cueson.schema.json")
	before := []byte("preserve me")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}

	commitFailure := errors.New("commit failed")
	err := replaceSchemaFileUsing(output, schema.Bytes(), func(temporaryPath, destinationPath string) error {
		if destinationPath != output {
			t.Fatalf("commit destination = %q, want %q", destinationPath, output)
		}
		if temporaryInfo, statErr := os.Stat(temporaryPath); statErr != nil || !temporaryInfo.Mode().IsRegular() {
			t.Fatalf("temporary output is not a regular file: info = %v, error = %v", temporaryInfo, statErr)
		}
		got, readErr := os.ReadFile(destinationPath)
		if readErr != nil || !bytes.Equal(got, before) {
			t.Fatalf("destination changed before commit: bytes = %q, error = %v", got, readErr)
		}
		return commitFailure
	})
	if !errors.Is(err, commitFailure) {
		t.Fatalf("replaceSchemaFileUsing() error = %v, want commit failure", err)
	}
	if got, err := os.ReadFile(output); err != nil || !bytes.Equal(got, before) {
		t.Fatalf("failed replacement changed destination: bytes = %q, error = %v", got, err)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(output) {
		t.Fatalf("failed replacement left temporary entries: %v", entries)
	}
}

func assertFileEqualsSchema(t *testing.T, path string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, schema.Bytes()) {
		t.Error("schema file does not match embedded bytes")
	}
}

func TestDiagnosticVisibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		quiet  bool
		silent bool
		want   string
	}{
		{name: "normal", want: "success: done\ninfo: working\nwarning: careful\nerror: failed\n"},
		{name: "quiet", quiet: true, want: "warning: careful\nerror: failed\n"},
		{name: "silent", silent: true, want: "error: failed\n"},
		{name: "quiet and silent", quiet: true, silent: true, want: "error: failed\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var output bytes.Buffer
			diagnostics := diagnosticWriter{writer: &output, quiet: tt.quiet, silent: tt.silent}
			diagnostics.write(diagnosticSuccess, "done")
			diagnostics.write(diagnosticInfo, "working")
			diagnostics.write(diagnosticWarning, "careful")
			diagnostics.write(diagnosticError, "failed")
			if got := output.String(); got != tt.want {
				t.Errorf("diagnostics = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDiagnosticColorEligibility(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		explicitOff    bool
		noColorPresent bool
		terminal       bool
		want           bool
	}{
		{name: "terminal", terminal: true, want: true},
		{name: "explicit off", explicitOff: true, terminal: true},
		{name: "NO_COLOR present", noColorPresent: true, terminal: true},
		{name: "non-terminal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := diagnosticColorEnabled(tt.explicitOff, tt.noColorPresent, tt.terminal); got != tt.want {
				t.Errorf("diagnosticColorEnabled() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestDiagnosticColorRendering(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	diagnostics := diagnosticWriter{writer: &output, color: true}
	diagnostics.write(diagnosticError, "failed")
	if got := output.String(); !strings.Contains(got, "\x1b[31merror:\x1b[0m failed\n") {
		t.Errorf("colored diagnostic = %q, want red error label", got)
	}
}

func TestRunStdoutWriteFailure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "root help"},
		{name: "version", args: []string{"version"}},
		{name: "version help", args: []string{"version", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stderr bytes.Buffer
			status := Run(context.Background(), tt.args, strings.NewReader(""), errorWriter{err: errors.New("output unavailable")}, &stderr)
			if status != ExitRuntimeFailure {
				t.Errorf("Run() status = %d, want %d", status, ExitRuntimeFailure)
			}
			if !strings.Contains(stderr.String(), "write stdout: output unavailable") {
				t.Errorf("Run() stderr = %q, want stdout write diagnostic", stderr.String())
			}
		})
	}
}

func TestParseEncodeAndRenderInvocation(t *testing.T) {
	t.Parallel()

	encode, err := parseInvocation([]string{"encode", "--format", "srt", "--encoding=utf-8", "--pretty", "--no-speaker-detection", "captions.srt"})
	if err != nil {
		t.Fatal(err)
	}
	if encode.encode.input != "captions.srt" || encode.encode.format != "srt" || encode.encode.encoding != "utf-8" || !encode.encode.pretty || !encode.encode.noSpeakerDetection {
		t.Fatalf("encode invocation = %+v", encode.encode)
	}

	render, err := parseInvocation([]string{"render", "document.cueson.json", "--to=srt", "--output", "rendered.srt", "--force", "--strict"})
	if err != nil {
		t.Fatal(err)
	}
	if render.render.input != "document.cueson.json" || render.render.target != "srt" || render.render.output != "rendered.srt" || !render.render.force || !render.render.strict {
		t.Fatalf("render invocation = %+v", render.render)
	}
}

func TestEncodeRenderAndRestoreWorkflow(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	sourceBytes := []byte("1\r\n00:00:01,250 --> 00:00:04,200\r\n<i>Ada: Hello.</i>\r\n")
	if err := os.WriteFile(sourcePath, sourceBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", "--pretty", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("encode = (%d, stderr %q)", status, stderr)
	}
	document, err := schema.Decode([]byte(encoded))
	if err != nil {
		t.Fatal(err)
	}
	if document.Format != "subrip" || document.FormatSupport.Status != "experimental" || len(document.Cues) != 1 {
		t.Fatalf("encoded document = %#v", document)
	}
	if document.Cues[0].Payload.RawText != "<i>Ada: Hello.</i>" || document.Cues[0].Payload.PlainText != "Ada: Hello." {
		t.Fatalf("encoded payload = %#v", document.Cues[0].Payload)
	}

	documentPath := filepath.Join(directory, "captions.cueson.json")
	if err := os.WriteFile(documentPath, []byte(encoded), 0o644); err != nil {
		t.Fatal(err)
	}
	status, rendered, stderr := runForTest(context.Background(), []string{"render", documentPath, "--to", "srt"})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("render = (%d, stderr %q)", status, stderr)
	}
	wantRendered := "1\n00:00:01,250 --> 00:00:04,200\n<i>Ada: Hello.</i>\n"
	if rendered != wantRendered {
		t.Fatalf("rendered = %q, want %q", rendered, wantRendered)
	}

	restoredPath := filepath.Join(directory, "restored.srt")
	status, stdout, stderr := runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", restoredPath, documentPath})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("restore = (%d, %q, %q)", status, stdout, stderr)
	}
	restored, err := os.ReadFile(restoredPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(restored, sourceBytes) {
		t.Fatalf("restored bytes = %q, want %q", restored, sourceBytes)
	}
}

func TestEncodeDefaultOutputAndMissingWebVTTCodec(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	if err := os.WriteFile(sourcePath, []byte("1\n00:00:00,000 --> 00:00:01,000\nHello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"encode", sourcePath})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("default encode = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, err := os.Stat(sourcePath + ".cueson.json"); err != nil {
		t.Fatal(err)
	}

	vttPath := filepath.Join(directory, "captions.vtt")
	if err := os.WriteFile(vttPath, []byte("WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"encode", "--stdout", vttPath})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "recognized but its native decode capability is unavailable") {
		t.Fatalf("WebVTT encode = (%d, %q, %q)", status, stdout, stderr)
	}
}

func TestRenderStrictRejectsUnrepresentableStructuredFields(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	if err := os.WriteFile(sourcePath, []byte("1\n00:00:00,000 --> 00:00:01,000\nAlice: Hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("encode = (%d, %q)", status, stderr)
	}
	documentPath := filepath.Join(directory, "document.json")
	if err := os.WriteFile(documentPath, []byte(encoded), 0o644); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"render", "--strict", "--to", "srt", documentPath})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "strict SubRip render blocked") {
		t.Fatalf("strict render = (%d, %q, %q)", status, stdout, stderr)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"render", "--to", "srt", documentPath})
	if status != ExitSuccess || stdout == "" || !strings.Contains(stderr, "subrip_render_fields_unrepresented") {
		t.Fatalf("normal render = (%d, %q, %q)", status, stdout, stderr)
	}
}

func TestParseEncodeAndRenderInvocationRejectsConflicts(t *testing.T) {
	t.Parallel()

	tests := [][]string{
		{"encode", "captions.srt", "--stdout", "--output", "other.json"},
		{"encode", "captions.srt", "--stdout", "--force"},
		{"encode", "captions.srt", "--format", "unknown"},
		{"render", "document.json"},
		{"render", "document.json", "--to", "unknown"},
		{"render", "document.json", "--to", "srt", "--force"},
	}
	for _, args := range tests {
		if _, err := parseInvocation(args); err == nil {
			t.Errorf("parseInvocation(%q) error = nil", args)
		}
	}
}

type errorWriter struct {
	err error
}

func (writer errorWriter) Write([]byte) (int, error) {
	return 0, writer.err
}

func runForTest(ctx context.Context, args []string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	status := Run(ctx, args, strings.NewReader(""), &stdout, &stderr)
	return status, stdout.String(), stderr.String()
}
