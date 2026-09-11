package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/testutil"
)

func TestCueJSONCommandsRejectBoundedPathPreconditionsWithoutOutputOrPathLeak(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	overPath := filepath.Join(directory, "PRIVATE-OVERSIZED.cueson.json")
	file, err := os.OpenFile(overPath, os.O_CREATE|os.O_WRONLY, 0o600)
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

	for _, command := range []struct {
		name string
		args func(string) []string
	}{
		{name: "restore", args: func(path string) []string { return []string{"restore", "--no-metadata", path} }},
		{name: "render", args: func(path string) []string { return []string{"render", path, "--to", "srt"} }},
	} {
		for _, path := range []string{directory, overPath, filepath.Join(directory, "PRIVATE-MISSING.cueson.json")} {
			status, stdout, stderr := runForTest(context.Background(), command.args(path))
			if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "path or size requirements") {
				t.Errorf("%s %q = (%d, %q, %q)", command.name, path, status, stdout, stderr)
			}
			if strings.Contains(stderr, path) || strings.Contains(stderr, directory) {
				t.Errorf("%s diagnostic leaked caller path: %q", command.name, stderr)
			}
		}
	}
}

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
		for _, want := range []string{"Usage:", "Commands:", "version", "Global options:", "Streams:", "Exit codes:", "Examples:"} {
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
		for _, implemented := range []string{"encode", "restore", "render", "convert", "validate", "inspect", "schema", "version", "completion"} {
			if !strings.Contains(stdout, "\n  "+implemented) {
				t.Errorf("Run(%q) help does not expose implemented command %q:\n%s", args, implemented, stdout)
			}
		}
	}
}

func TestParseConvertInvocationAndAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		wantFrom   string
		wantTarget string
	}{
		{name: "defaults", args: []string{"convert", "captions.srt", "--to", "vtt"}, wantFrom: "auto", wantTarget: "webvtt"},
		{name: "canonical Cue JSON", args: []string{"convert", "--from", "cueson", "document.json", "--to=subrip"}, wantFrom: "cueson", wantTarget: "subrip"},
		{name: "JSON alias", args: []string{"convert", "--from=json", "document.json", "--to", "srt"}, wantFrom: "cueson", wantTarget: "subrip"},
		{name: "Cue JSON alias", args: []string{"convert", "--from", "cue-json", "document.json", "--to", "webvtt"}, wantFrom: "cueson", wantTarget: "webvtt"},
		{name: "SubRip alias", args: []string{"convert", "--from=subrip", "captions.txt", "--to", "vtt"}, wantFrom: "subrip", wantTarget: "webvtt"},
		{name: "WebVTT alias", args: []string{"convert", "--from", "webvtt", "captions.txt", "--to=srt"}, wantFrom: "webvtt", wantTarget: "subrip"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := parseInvocation(tt.args)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.command != "convert" || parsed.convert.input == "" || parsed.convert.from != tt.wantFrom || parsed.convert.target != tt.wantTarget {
				t.Fatalf("convert invocation = %#v, want from=%q target=%q", parsed.convert, tt.wantFrom, tt.wantTarget)
			}
		})
	}

	parsed, err := parseInvocation([]string{"--quiet", "convert", "--encoding", "utf-8", "--output", "out.vtt", "--force", "--strict", "--no-speaker-detection", "in.srt", "--to", "vtt", "--silent", "--no-color"})
	if err != nil {
		t.Fatal(err)
	}
	if parsed.convert.encoding != "utf-8" || parsed.convert.output != "out.vtt" || !parsed.convert.outputSet || !parsed.convert.force || !parsed.convert.strict || !parsed.convert.noSpeakerDetection || !parsed.options.quiet || !parsed.options.silent || !parsed.options.noColor {
		t.Fatalf("convert options = %#v, globals = %#v", parsed.convert, parsed.options)
	}
}

func TestRunConvertHelpAndInvalidInvocation(t *testing.T) {
	t.Parallel()

	for _, args := range [][]string{{"convert", "--help"}, {"convert", "-h"}} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess || stderr != "" {
			t.Fatalf("Run(%q) = (%d, %q, %q)", args, status, stdout, stderr)
		}
		for _, want := range []string{"cueson [global options] convert", "--from", "--to", "--encoding", "--output", "--force", "--strict", "--no-speaker-detection"} {
			if !strings.Contains(stdout, want) {
				t.Errorf("Run(%q) help lacks %q:\n%s", args, want, stdout)
			}
		}
	}

	tests := []struct {
		args []string
		want string
	}{
		{args: []string{"convert"}, want: "requires one INPUT"},
		{args: []string{"convert", "in.srt"}, want: "requires --to"},
		{args: []string{"convert", "in.srt", "--to", "unknown"}, want: "--to must be srt or vtt"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--from", "unknown"}, want: "--from must be auto, cueson, srt, or vtt"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--from", "cueson", "--encoding", "utf-8"}, want: "--encoding cannot be used with Cue JSON"},
		{args: []string{"convert", "in.vtt", "--to", "srt", "--from", "vtt", "--encoding", "windows-1252"}, want: "WebVTT requires UTF-8"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--encoding", "unknown"}, want: "encoding \"unknown\" is not supported"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--force"}, want: "--force requires a filesystem output"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--output", "-", "--force"}, want: "--force requires a filesystem output"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--output", "a", "--output", "b"}, want: "--output may be specified only once"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--to", "srt"}, want: "--to may be specified only once"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--from", "srt", "--from", "srt"}, want: "--from may be specified only once"},
		{args: []string{"convert", "in.srt", "--to", "vtt", "--encoding", "utf-8", "--encoding", "utf-8"}, want: "--encoding may be specified only once"},
		{args: []string{"convert", "a", "b", "--to", "vtt"}, want: "no additional arguments"},
		{args: []string{"convert", "--bogus", "in.srt", "--to", "vtt"}, want: "unknown option"},
	}
	for _, tt := range tests {
		status, stdout, stderr := runForTest(context.Background(), tt.args)
		if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, tt.want) || !strings.Contains(stderr, "cueson [global options] convert") {
			t.Errorf("Run(%q) = (%d, %q, %q), want invocation error containing %q", tt.args, status, stdout, stderr, tt.want)
		}
	}
}

func TestConvertNativeAndCueJSONInputsHaveParity(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	sourceBytes := []byte("1\n00:00:01,000 --> 00:00:02,500\n<i>Hello</i> world\n")
	if err := os.WriteFile(sourcePath, sourceBytes, 0o600); err != nil {
		t.Fatal(err)
	}

	status, direct, stderr := runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt"})
	if status != ExitSuccess || stderr != "" || !strings.HasPrefix(direct, "WEBVTT\n") {
		t.Fatalf("direct conversion = (%d, %q, %q)", status, direct, stderr)
	}

	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("encode = (%d, %q, %q)", status, encoded, stderr)
	}
	misleadingJSONPath := filepath.Join(directory, "document.vtt")
	if err := os.WriteFile(misleadingJSONPath, []byte(encoded), 0o600); err != nil {
		t.Fatal(err)
	}
	status, fromJSON, stderr := runForTest(context.Background(), []string{"convert", misleadingJSONPath, "--to", "webvtt"})
	if status != ExitSuccess || stderr != "" || fromJSON != direct {
		t.Fatalf("Cue JSON conversion = (%d, %q, %q), want direct bytes %q", status, fromJSON, stderr, direct)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"convert", misleadingJSONPath, "--to", "webvtt", "--encoding", "utf-8"})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "--encoding cannot be used with automatically detected Cue JSON input") || !strings.Contains(stderr, "cueson [global options] convert") {
		t.Fatalf("auto-detected Cue JSON encoding conflict = (%d, %q, %q)", status, stdout, stderr)
	}

	renderedPath := filepath.Join(directory, "converted.vtt")
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", "--from", "cueson", "--to", "vtt", "--output", renderedPath, misleadingJSONPath})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("file conversion = (%d, %q, %q)", status, stdout, stderr)
	}
	before := readFileForCLI(t, renderedPath)
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", "--from", "cueson", "--to", "vtt", "--output", renderedPath, misleadingJSONPath})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "already exists") || !bytes.Equal(readFileForCLI(t, renderedPath), before) {
		t.Fatalf("refused conversion replacement = (%d, %q, %q)", status, stdout, stderr)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", "--from", "cueson", "--to", "vtt", "--output", renderedPath, "--force", misleadingJSONPath})
	if status != ExitSuccess || stdout != "" || stderr != "" || !bytes.Equal(readFileForCLI(t, renderedPath), before) {
		t.Fatalf("forced conversion replacement = (%d, %q, %q)", status, stdout, stderr)
	}
	status, reencoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", renderedPath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("target re-encode = (%d, %q, %q)", status, reencoded, stderr)
	}
	if document, err := schema.Decode([]byte(reencoded)); err != nil || document.Format != "webvtt" {
		t.Fatalf("target re-encoded document format = %q, error = %v", document.Format, err)
	}
}

func TestConvertWebVTTSourceToSubRip(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.vtt")
	if err := os.WriteFile(sourcePath, []byte("WEBVTT\n\n00:00.500 --> 00:02.000\n<b>Hello</b> world\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	status, output, stderr := runForTest(context.Background(), []string{"convert", "--from", "webvtt", sourcePath, "--to", "subrip", "--output=-"})
	if status != ExitSuccess || stderr != "" || !strings.Contains(output, "00:00:00,500 --> 00:00:02,000") || !strings.Contains(output, "<b>Hello</b> world") {
		t.Fatalf("WebVTT to SubRip = (%d, %q, %q)", status, output, stderr)
	}
	convertedPath := filepath.Join(directory, "converted.srt")
	if err := os.WriteFile(convertedPath, []byte(output), 0o600); err != nil {
		t.Fatal(err)
	}
	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", convertedPath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("target re-encode = (%d, %q, %q)", status, encoded, stderr)
	}
	if document, err := schema.Decode([]byte(encoded)); err != nil || document.Format != "subrip" || document.Cues[0].Payload.PlainText != "Hello world" {
		t.Fatalf("target re-encoded document = %#v, error = %v", document, err)
	}
}

func TestConvertCueJSONClassificationAndIntegrityFailures(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	nativeBytes := []byte("1\n00:00:00,000 --> 00:00:01,000\nHello\n")
	jsonNamedNative := filepath.Join(directory, "captions.json")
	if err := os.WriteFile(jsonNamedNative, nativeBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"convert", jsonNamedNative, "--to", "vtt"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "Cue JSON") {
		t.Fatalf("JSON-named native auto classification = (%d, %q, %q)", status, stdout, stderr)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", "--from", "srt", jsonNamedNative, "--to", "vtt"})
	if status != ExitSuccess || stdout == "" || stderr != "" {
		t.Fatalf("explicit native classification = (%d, %q, %q)", status, stdout, stderr)
	}
	jsonLookingNative := filepath.Join(directory, "captions.srt")
	if err := os.WriteFile(jsonLookingNative, []byte("{not JSON or subtitles}"), 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", jsonLookingNative, "--to", "vtt"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "Cue JSON") {
		t.Fatalf("JSON-looking native auto classification = (%d, %q, %q)", status, stdout, stderr)
	}
	bomJSONPath := filepath.Join(directory, "document.txt")
	bomJSON := append([]byte{0xef, 0xbb, 0xbf}, []byte("{not valid Cue JSON}")...)
	if err := os.WriteFile(bomJSONPath, bomJSON, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", bomJSONPath, "--to", "vtt"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "Cue JSON") {
		t.Fatalf("BOM JSON-looking classification = (%d, %q, %q)", status, stdout, stderr)
	}

	sourcePath := filepath.Join(directory, "source.srt")
	if err := os.WriteFile(sourcePath, nativeBytes, 0o600); err != nil {
		t.Fatal(err)
	}
	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatal(stderr)
	}
	document, err := schema.Decode([]byte(encoded))
	if err != nil {
		t.Fatal(err)
	}
	document.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64)
	corrupt, err := marshalDocument(document, false)
	if err != nil {
		t.Fatal(err)
	}
	corruptPath := filepath.Join(directory, "corrupt.cueson.json")
	if err := os.WriteFile(corruptPath, corrupt, 0o600); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(directory, "must-not-exist.vtt")
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", corruptPath, "--to", "vtt", "--output", output})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "SHA-256") {
		t.Fatalf("integrity failure = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, statErr := os.Lstat(output); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("integrity failure left output: %v", statErr)
	}
}

func TestConvertLossWarningsFiltersAndStrictOutputSafety(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "coordinates.srt")
	input := []byte("1\n00:00:00,000 --> 00:00:01,000 X1:10 X2:20 Y1:30 Y2:40\nHello\n")
	if err := os.WriteFile(sourcePath, input, 0o600); err != nil {
		t.Fatal(err)
	}

	status, converted, stderr := runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt"})
	if status != ExitSuccess || converted == "" || !strings.Contains(strings.ToLower(stderr), "coordinate") || strings.Contains(stderr, directory) {
		t.Fatalf("permissive loss = (%d, %q, %q)", status, converted, stderr)
	}
	status, quietOutput, quietDiagnostics := runForTest(context.Background(), []string{"--quiet", "convert", sourcePath, "--to", "vtt"})
	if status != ExitSuccess || quietOutput != converted || quietDiagnostics != stderr {
		t.Fatalf("quiet loss = (%d, %q, %q), want warnings retained", status, quietOutput, quietDiagnostics)
	}
	status, silentOutput, silentDiagnostics := runForTest(context.Background(), []string{"--silent", "convert", sourcePath, "--to", "vtt"})
	if status != ExitSuccess || silentOutput != converted || silentDiagnostics != "" {
		t.Fatalf("silent loss = (%d, %q, %q)", status, silentOutput, silentDiagnostics)
	}

	newOutput := filepath.Join(directory, "strict-new.vtt")
	status, stdout, stderr := runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt", "--strict", "--output", newOutput})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(strings.ToLower(stderr), "strict") {
		t.Fatalf("strict new output = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, statErr := os.Lstat(newOutput); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("strict conversion left new output: %v", statErr)
	}
	existingOutput := filepath.Join(directory, "strict-existing.vtt")
	before := []byte("preserve me")
	if err := os.WriteFile(existingOutput, before, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt", "--strict", "--output", existingOutput, "--force"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(strings.ToLower(stderr), "strict") {
		t.Fatalf("strict forced output = (%d, %q, %q)", status, stdout, stderr)
	}
	if got := readFileForCLI(t, existingOutput); !bytes.Equal(got, before) {
		t.Fatalf("strict conversion changed destination: %q", got)
	}
}

func TestConvertOrdersSourceDiagnosticsBeforeLossesAndRejectsBadDestinations(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	input := []byte("00:00:00,000 --> 00:00:01,000 X1:10 X2:20 Y1:30 Y2:40\nHello\n")
	if err := os.WriteFile(sourcePath, input, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt"})
	sourceIndex := strings.Index(stderr, "subrip_sequence_missing")
	lossIndex := strings.Index(stderr, "conversion_subrip_coordinates_omitted")
	if status != ExitSuccess || stdout == "" || sourceIndex < 0 || lossIndex <= sourceIndex {
		t.Fatalf("ordered diagnostics = (%d, %q, %q), indexes = (%d, %d)", status, stdout, stderr, sourceIndex, lossIndex)
	}

	missingParent := filepath.Join(directory, "missing", "captions.vtt")
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt", "--output", missingParent})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "does not exist") || !strings.Contains(stderr, "cueson [global options] convert") {
		t.Fatalf("missing output parent = (%d, %q, %q)", status, stdout, stderr)
	}
	directoryOutput := filepath.Join(directory, "destination")
	if err := os.Mkdir(directoryOutput, 0o700); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt", "--output", directoryOutput, "--force"})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "not a regular file") {
		t.Fatalf("directory output = (%d, %q, %q)", status, stdout, stderr)
	}
}

func TestConvertLossyFixtureWarningsMatchCanonicalGoldens(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		target string
		want   string
	}{
		{name: "SubRip to WebVTT", input: "../../testdata/fixtures/conversion/srt-lossy/source/input.srt", target: "vtt", want: "186060b924cc054b5b4192c4688f3434a5c5a28b1dceb6566c82ac449fd12494"},
		{name: "WebVTT to SubRip", input: "../../testdata/fixtures/conversion/webvtt-lossy/source/input.vtt", target: "srt", want: "dc930384b305233f2ebd27fac948ed146cd6eba8a50aea9df79b00255a8b27ad"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, stdout, stderr := runForTest(context.Background(), []string{"convert", test.input, "--to", test.target})
			if status != ExitSuccess || stdout == "" {
				t.Fatalf("conversion = (%d, %q, %q)", status, stdout, stderr)
			}
			got := fmt.Sprintf("%x", sha256.Sum256([]byte(stderr)))
			if got != test.want {
				t.Fatalf("canonical warning SHA-256 = %s, want %s; warnings = %q", got, test.want, stderr)
			}
		})
	}
}

func TestConvertSameFormatEncodingAndStdoutFailures(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.vtt")
	if err := os.WriteFile(sourcePath, []byte("WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"convert", sourcePath, "--to", "vtt"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "render") {
		t.Fatalf("same-format conversion = (%d, %q, %q)", status, stdout, stderr)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"convert", sourcePath, "--to", "srt", "--encoding", "windows-1252"})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "WebVTT requires UTF-8") {
		t.Fatalf("auto WebVTT encoding failure = (%d, %q, %q)", status, stdout, stderr)
	}

	var diagnostics bytes.Buffer
	status = Run(context.Background(), []string{"convert", sourcePath, "--to", "srt"}, strings.NewReader(""), errorWriter{err: errors.New("output unavailable")}, &diagnostics)
	if status != ExitRuntimeFailure || !strings.Contains(diagnostics.String(), "write stdout: output unavailable") {
		t.Fatalf("stdout failure = (%d, %q)", status, diagnostics.String())
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
		for _, want := range []string{"Usage:", "cueson version", "Print the Cueson executable version", "Examples:"} {
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

func TestWriteSchemaFileWriteFailureLeavesNewDestinationAbsent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		write func(*os.File, []byte) error
	}{
		{
			name: "short write",
			write: func(file *os.File, payload []byte) error {
				if _, err := file.Write(payload[:len(payload)/2]); err != nil {
					return err
				}
				return io.ErrShortWrite
			},
		},
		{name: "failed write", write: func(*os.File, []byte) error { return errors.New("write failed") }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			directory := t.TempDir()
			output := filepath.Join(directory, "new.json")
			err := writeSchemaFileUsing(output, schema.Bytes(), false, publicationHooks{write: tt.write})
			if err == nil {
				t.Fatal("writeSchemaFileUsing() error = nil")
			}
			if _, statErr := os.Lstat(output); !errors.Is(statErr, fs.ErrNotExist) {
				t.Fatalf("failed publication left destination: %v", statErr)
			}
			entries, readErr := os.ReadDir(directory)
			if readErr != nil {
				t.Fatal(readErr)
			}
			if len(entries) != 0 {
				t.Fatalf("failed publication left staging entries: %v", entries)
			}
		})
	}
}

func TestWriteSchemaFileForcedFailurePreservesDestination(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "existing.json")
	before := []byte("preserve me")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}
	fault := errors.New("replace failed")
	err := writeSchemaFileUsing(output, schema.Bytes(), true, publicationHooks{
		replace: func(string, string) error { return fault },
	})
	if !errors.Is(err, fault) {
		t.Fatalf("writeSchemaFileUsing() error = %v, want %v", err, fault)
	}
	if got, readErr := os.ReadFile(output); readErr != nil || !bytes.Equal(got, before) {
		t.Fatalf("failed forced replacement changed destination: bytes = %q, error = %v", got, readErr)
	}
	assertOnlyNamedEntry(t, directory, filepath.Base(output))
}

func TestWriteSchemaFileVerificationFailureRollsBackForcedReplacement(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "existing.json")
	before := []byte("preserve me")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}
	fault := errors.New("final verification failed")
	verification := 0
	err := writeSchemaFileUsing(output, schema.Bytes(), true, publicationHooks{
		verify: func(path string, size int64, digest string) error {
			verification++
			if verification == 2 {
				return fault
			}
			return verifyPublishedFile(path, size, digest)
		},
	})
	if !errors.Is(err, fault) {
		t.Fatalf("writeSchemaFileUsing() error = %v, want %v", err, fault)
	}
	if got, readErr := os.ReadFile(output); readErr != nil || !bytes.Equal(got, before) {
		t.Fatalf("verification rollback did not preserve destination: bytes = %q, error = %v", got, readErr)
	}
	assertOnlyNamedEntry(t, directory, filepath.Base(output))
}

func assertOnlyNamedEntry(t *testing.T, directory, name string) {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != name {
		t.Fatalf("directory entries = %v, want only %q", entries, name)
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

func TestDiagnosticWriteFailureFinalizesSuccessfulStatus(t *testing.T) {
	t.Parallel()

	diagnostics := newDiagnosticWriter(errorWriter{err: errors.New("diagnostic unavailable")}, globalOptions{})
	diagnostics.write(diagnosticSuccess, "valid input")
	if got := diagnostics.finalize(ExitSuccess); got != ExitRuntimeFailure {
		t.Fatalf("finalize(success) = %d, want %d", got, ExitRuntimeFailure)
	}
	if got := diagnostics.finalize(ExitInvocation); got != ExitInvocation {
		t.Fatalf("finalize(invocation) = %d, want %d", got, ExitInvocation)
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

func TestEncodeDefaultOutputAndWebVTTCodec(t *testing.T) {
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
	status, stdout, stderr = runForTest(context.Background(), []string{"encode", vttPath})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("default WebVTT encode = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, err := os.Stat(vttPath + ".cueson.json"); err != nil {
		t.Fatal(err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"encode", "--output=-", vttPath})
	if status != ExitSuccess || stdout == "" || stderr != "" {
		t.Fatalf("WebVTT encode = (%d, %q, %q)", status, stdout, stderr)
	}
	document, err := schema.Decode([]byte(stdout))
	if err != nil {
		t.Fatal(err)
	}
	if document.Format != "webvtt" || document.FormatData.WebVTT == nil || document.Source.Assets[0].MediaType == nil || *document.Source.Assets[0].MediaType != "text/vtt" {
		t.Fatalf("WebVTT document = %#v", document)
	}
}

func TestWebVTTEncodeRenderRestoreWorkflow(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "private-location.vtt")
	sourceBytes := append([]byte{0xef, 0xbb, 0xbf}, []byte("WEBVTT sample\r\n\r\nNOTE retained\r\n\r\nid-one\r\n00:01.000 --> 00:02.500 align:start\r\n<v Ada>Hello &amp; welcome</v>\r\n")...)
	if err := os.WriteFile(sourcePath, sourceBytes, 0o644); err != nil {
		t.Fatal(err)
	}

	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", "--pretty", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("encode = (%d, stderr %q)", status, stderr)
	}
	if strings.Contains(encoded, directory) || strings.Contains(encoded, sourcePath) {
		t.Fatalf("Cue JSON leaked input path: %q", encoded)
	}
	document, err := schema.Decode([]byte(encoded))
	if err != nil {
		t.Fatalf("schema.Decode(encoded) error = %v", err)
	}
	if document.Format != "webvtt" || document.FormatSupport.Status != "experimental" || !document.FormatSupport.IngestSupported || !document.FormatSupport.RenderSupported || !document.FormatSupport.RestoreSupported {
		t.Fatalf("format support = (%q, %#v)", document.Format, document.FormatSupport)
	}
	if document.FormatData.WebVTT == nil || len(document.FormatData.WebVTT.Blocks) != 1 || len(document.Cues) != 1 || document.Cues[0].Payload.PlainText != "Hello & welcome" {
		t.Fatalf("WebVTT model = %#v", document)
	}
	if len(document.Cues[0].Speakers) != 1 || document.Cues[0].Speakers[0].Name != "Ada" {
		t.Fatalf("WebVTT speaker observations = %#v", document.Cues[0].Speakers)
	}
	status, withoutSpeakers, stderr := runForTest(context.Background(), []string{"encode", "--stdout", "--no-speaker-detection", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("encode without speaker detection = (%d, stderr %q)", status, stderr)
	}
	withoutSpeakerDocument, err := schema.Decode([]byte(withoutSpeakers))
	if err != nil || len(withoutSpeakerDocument.Cues[0].Speakers) != 0 {
		t.Fatalf("disabled speaker observations = %#v, decode error = %v", withoutSpeakerDocument.Cues[0].Speakers, err)
	}

	documentPath := filepath.Join(directory, "document.cueson.json")
	if err := os.WriteFile(documentPath, []byte(encoded), 0o644); err != nil {
		t.Fatal(err)
	}
	status, rendered, stderr := runForTest(context.Background(), []string{"render", documentPath, "--to", "webvtt"})
	if status != ExitSuccess || stderr != "" || !strings.HasPrefix(rendered, "WEBVTT sample\n") || strings.Contains(rendered, "\r") {
		t.Fatalf("render = (%d, %q, %q)", status, rendered, stderr)
	}
	status, dashRendered, stderr := runForTest(context.Background(), []string{"render", documentPath, "--to", "vtt", "--output=-"})
	if status != ExitSuccess || dashRendered != rendered || stderr != "" {
		t.Fatalf("render --output=- = (%d, %q, %q), want the default stdout rendering", status, dashRendered, stderr)
	}

	reencodedPath := filepath.Join(directory, "reencoded.json")
	renderedPath := filepath.Join(directory, "rendered.vtt")
	if err := os.WriteFile(renderedPath, []byte(rendered), 0o644); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"encode", "--format", "vtt", "--output", reencodedPath, renderedPath})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("re-encode = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, err := schema.Decode(readFileForCLI(t, reencodedPath)); err != nil {
		t.Fatalf("schema.Decode(re-encoded) error = %v", err)
	}

	restoredPath := filepath.Join(directory, "restored.vtt")
	status, stdout, stderr = runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", restoredPath, documentPath})
	if status != ExitSuccess || stdout != "" || stderr != "" {
		t.Fatalf("restore = (%d, %q, %q)", status, stdout, stderr)
	}
	if restored := readFileForCLI(t, restoredPath); !bytes.Equal(restored, sourceBytes) {
		t.Fatalf("restored bytes = %q, want %q", restored, sourceBytes)
	}
}

func TestWebVTTEncodeSelectionAndEncodingFailures(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	disguised := filepath.Join(directory, "captions.srt")
	if err := os.WriteFile(disguised, []byte("WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", disguised})
	if status != ExitSuccess || encoded == "" || !strings.Contains(stderr, "format_extension_disagreement") {
		t.Fatalf("extension disagreement = (%d, %q, %q)", status, encoded, stderr)
	}
	document, err := schema.Decode([]byte(encoded))
	if err != nil || document.Format != "webvtt" {
		t.Fatalf("selected format = %q, decode error = %v", document.Format, err)
	}

	output := filepath.Join(directory, "must-not-exist.json")
	status, stdout, stderr := runForTest(context.Background(), []string{"encode", "--format", "vtt", "--encoding", "windows-1252", "--output", output, disguised})
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "WebVTT requires UTF-8") {
		t.Fatalf("explicit incompatible encoding = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, statErr := os.Lstat(output); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("explicit incompatible encoding left output: %v", statErr)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"encode", "--encoding", "windows-1252", "--output", output, disguised})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "WebVTT requires UTF-8") {
		t.Fatalf("automatic incompatible encoding = (%d, %q, %q)", status, stdout, stderr)
	}
	if _, statErr := os.Lstat(output); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("automatic incompatible encoding left output: %v", statErr)
	}
}

func TestWebVTTOutputForceAndStrictRender(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.vtt")
	if err := os.WriteFile(sourcePath, []byte("WEBVTT\n\n00:00.000 --> 00:01.000 mystery:x\nHello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	encodedPath := filepath.Join(directory, "document.json")
	status, stdout, stderr := runForTest(context.Background(), []string{"encode", "--output", encodedPath, sourcePath})
	if status != ExitSuccess || stdout != "" || !strings.Contains(stderr, "webvtt_setting_unknown") {
		t.Fatalf("encode output = (%d, %q, %q)", status, stdout, stderr)
	}

	before := []byte("preserve render destination")
	renderedPath := filepath.Join(directory, "captions.rendered.vtt")
	if err := os.WriteFile(renderedPath, before, 0o600); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"render", "--strict", "--to", "vtt", "--output", renderedPath, "--force", encodedPath})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "webvtt_setting_unknown") {
		t.Fatalf("strict render = (%d, %q, %q)", status, stdout, stderr)
	}
	if got := readFileForCLI(t, renderedPath); !bytes.Equal(got, before) {
		t.Fatalf("strict render changed output: %q", got)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"render", "--to", "vtt", "--output", renderedPath, "--force", encodedPath})
	if status != ExitSuccess || stdout != "" || !strings.Contains(stderr, "webvtt_render_preserved_nonconforming") {
		t.Fatalf("permissive forced render = (%d, %q, %q)", status, stdout, stderr)
	}
	if got := readFileForCLI(t, renderedPath); !bytes.HasPrefix(got, []byte("WEBVTT\n")) {
		t.Fatalf("rendered file = %q", got)
	}
}

func TestEveryAcceptedWebVTTFixtureEncodesAndRestoresExactly(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..", "testdata")
	manifest, err := testutil.VerifyFixtures(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range manifest.Fixtures {
		if fixture.Expectation.Result != "accepted" || !strings.HasPrefix(fixture.ID, "webvtt/") {
			continue
		}
		fixture := fixture
		t.Run(strings.TrimPrefix(fixture.ID, "webvtt/"), func(t *testing.T) {
			t.Parallel()
			var artifact *testutil.Artifact
			for index := range fixture.Artifacts {
				if fixture.Artifacts[index].Role == "source" {
					artifact = &fixture.Artifacts[index]
					break
				}
			}
			if artifact == nil {
				t.Fatal("accepted WebVTT fixture has no source artifact")
			}
			sourcePath := testutil.ArtifactFile(root, *artifact)
			sourceBytes := readFileForCLI(t, sourcePath)
			status, encoded, _ := runForTest(context.Background(), []string{"encode", "--stdout", sourcePath})
			if status != ExitSuccess {
				t.Fatalf("encode status = %d", status)
			}
			if _, err := schema.Decode([]byte(encoded)); err != nil {
				t.Fatalf("schema.Decode(encoded) error = %v", err)
			}
			directory := t.TempDir()
			documentPath := filepath.Join(directory, "document.cueson.json")
			if err := os.WriteFile(documentPath, []byte(encoded), 0o600); err != nil {
				t.Fatal(err)
			}
			restoredPath := filepath.Join(directory, "restored.vtt")
			status, stdout, stderr := runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", restoredPath, documentPath})
			if status != ExitSuccess || stdout != "" || stderr != "" {
				t.Fatalf("restore = (%d, %q, %q)", status, stdout, stderr)
			}
			if restored := readFileForCLI(t, restoredPath); !bytes.Equal(restored, sourceBytes) {
				t.Fatal("restored bytes differ from governed source")
			}
		})
	}
}

func TestEveryAcceptedSubRipFixtureEncodesAndRestoresExactly(t *testing.T) {
	t.Parallel()

	root := filepath.Join("..", "..", "testdata")
	manifest, err := testutil.VerifyFixtures(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, fixture := range manifest.Fixtures {
		if fixture.Expectation.Result != "accepted" || !strings.HasPrefix(fixture.ID, "subrip/") {
			continue
		}
		for _, artifact := range fixture.Artifacts {
			if artifact.Role != "source" {
				continue
			}
			fixture, artifact := fixture, artifact
			t.Run(strings.TrimPrefix(fixture.ID, "subrip/")+"/"+artifact.ID, func(t *testing.T) {
				t.Parallel()
				sourcePath := testutil.ArtifactFile(root, artifact)
				sourceBytes := readFileForCLI(t, sourcePath)
				args := []string{"encode", "--stdout"}
				switch artifact.ID {
				case "windows_1252":
					args = append(args, "--encoding", "windows-1252")
				case "iso_8859_1":
					args = append(args, "--encoding", "iso-8859-1")
				}
				args = append(args, sourcePath)
				status, encoded, stderr := runForTest(context.Background(), args)
				if status != ExitSuccess {
					t.Fatalf("encode = (%d, %q, %q)", status, encoded, stderr)
				}
				if _, err := schema.Decode([]byte(encoded)); err != nil {
					t.Fatalf("schema.Decode(encoded) error = %v", err)
				}
				directory := t.TempDir()
				documentPath := filepath.Join(directory, "document.cueson.json")
				if err := os.WriteFile(documentPath, []byte(encoded), 0o600); err != nil {
					t.Fatal(err)
				}
				restoredPath := filepath.Join(directory, "restored.srt")
				status, stdout, stderr := runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", restoredPath, documentPath})
				if status != ExitSuccess || stdout != "" || stderr != "" {
					t.Fatalf("restore = (%d, %q, %q)", status, stdout, stderr)
				}
				if restored := readFileForCLI(t, restoredPath); !bytes.Equal(restored, sourceBytes) {
					t.Fatal("restored bytes differ from governed source")
				}
			})
		}
	}
}

func readFileForCLI(t *testing.T, path string) []byte {
	t.Helper()
	payload, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func TestEncodeDoesNotTreatWebVTTPrefixAsSignature(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	input := []byte("WEBVTTfoo\n\n1\n00:00:00,000 --> 00:00:01,000\nHello\n")
	if err := os.WriteFile(sourcePath, input, 0o644); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"encode", "--stdout", sourcePath})
	if status != ExitSuccess || stdout == "" || strings.Contains(stderr, "WebVTT") {
		t.Fatalf("SubRip encode with WEBVTT prefix = (%d, %q, %q)", status, stdout, stderr)
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

func TestRenderReportsPayloadThatReparsesAsCueBoundary(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	sourcePath := filepath.Join(directory, "captions.srt")
	if err := os.WriteFile(sourcePath, []byte("1\n00:00:00,000 --> 00:00:01,000\ncaption\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, encoded, stderr := runForTest(context.Background(), []string{"encode", "--stdout", "--no-speaker-detection", sourcePath})
	if status != ExitSuccess || stderr != "" {
		t.Fatalf("encode = (%d, %q)", status, stderr)
	}
	document, err := schema.Decode([]byte(encoded))
	if err != nil {
		t.Fatal(err)
	}
	document.Cues[0].Payload.RawText = "caption\n2\n00:00:02,000 --> 00:00:03,000\ntail"
	document.Cues[0].Payload.PlainText = document.Cues[0].Payload.RawText
	document.Cues[0].Payload.Lines = strings.Split(document.Cues[0].Payload.RawText, "\n")
	payload, err := marshalDocument(document, false)
	if err != nil {
		t.Fatal(err)
	}
	documentPath := filepath.Join(directory, "document.json")
	if err := os.WriteFile(documentPath, payload, 0o644); err != nil {
		t.Fatal(err)
	}
	status, stdout, stderr := runForTest(context.Background(), []string{"render", "--strict", "--to", "srt", documentPath})
	if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "payload raw_text") {
		t.Fatalf("strict render = (%d, %q, %q)", status, stdout, stderr)
	}
	status, stdout, stderr = runForTest(context.Background(), []string{"render", "--to", "srt", documentPath})
	if status != ExitSuccess || stdout == "" || !strings.Contains(stderr, "subrip_render_payload_ambiguous") {
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
