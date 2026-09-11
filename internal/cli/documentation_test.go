package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestMaintainedCLIReferenceCoversExecutableSurface(t *testing.T) {
	t.Parallel()

	repo := documentationRepositoryRoot(t)
	reference := readDocumentationFile(t, filepath.Join(repo, "docs", "cli.md"))
	for _, command := range orderedCommandSurface {
		section := markdownSection(reference, "## `"+command.Name+"`")
		if section == "" {
			t.Errorf("CLI reference lacks %s section", command.Name)
			continue
		}
		for _, required := range []string{command.Usage} {
			if !strings.Contains(section, required) {
				t.Errorf("%s reference lacks %q", command.Name, required)
			}
		}
		for _, option := range command.Options {
			for _, spelling := range option.Spellings {
				if !strings.Contains(section, "`"+spelling+"`") {
					t.Errorf("%s reference lacks option %s", command.Name, spelling)
				}
			}
			for _, value := range option.Values {
				if !strings.Contains(section, "`"+value+"`") {
					t.Errorf("%s reference lacks option value %s", command.Name, value)
				}
			}
		}
		for _, selector := range command.OperandValues {
			if !strings.Contains(section, "`"+selector+"`") {
				t.Errorf("%s reference lacks selector %s", command.Name, selector)
			}
		}
		allowedOptions := make(map[string]bool)
		for _, option := range append(append([]cliOptionSpec{}, globalOptionSurface...), command.Options...) {
			for _, spelling := range option.Spellings {
				allowedOptions[spelling] = true
			}
		}
		for _, spelling := range regexp.MustCompile("`(--?[a-z][a-z-]*)`").FindAllStringSubmatch(section, -1) {
			if !allowedOptions[spelling[1]] {
				t.Errorf("%s reference documents unknown option %s", command.Name, spelling[1])
			}
		}
	}

	for _, required := range []string{"`-q`", "`--quiet`", "`--silent`", "`--no-color`", "`-h`", "`--help`", "`--`", "stdout", "stderr", "Exit code `0`", "Exit code `1`", "Exit code `2`", "`NO_COLOR`"} {
		if !strings.Contains(reference, required) {
			t.Errorf("CLI reference lacks global contract %q", required)
		}
	}
}

func TestREADMERegistersExecutableWorkflows(t *testing.T) {
	t.Parallel()

	readme := readDocumentationFile(t, filepath.Join(documentationRepositoryRoot(t), "README.md"))
	for _, command := range []string{
		"go run ./cmd/cueson version",
		"go run ./cmd/cueson encode --pretty --output quickstart/document.cueson.json testdata/fixtures/conversion/srt-loss-free/source/input.srt",
		"go run ./cmd/cueson restore --no-metadata --output quickstart/restored.srt quickstart/document.cueson.json",
		"go run ./cmd/cueson render --to srt --output quickstart/rendered.srt quickstart/document.cueson.json",
		"go run ./cmd/cueson convert --strict --no-speaker-detection --to vtt --output quickstart/converted.vtt testdata/fixtures/conversion/srt-loss-free/source/input.srt",
		"go run ./cmd/cueson convert --strict --to srt --output quickstart/converted.srt testdata/fixtures/conversion/webvtt-loss-free/source/input.vtt",
		"go run ./cmd/cueson validate testdata/fixtures/webvtt/minimal/source/minimal.vtt",
		"go run ./cmd/cueson inspect testdata/fixtures/webvtt/minimal/source/minimal.vtt",
		"go run ./cmd/cueson inspect --json testdata/fixtures/webvtt/minimal/source/minimal.vtt",
		"go run ./cmd/cueson completion bash > quickstart/cueson.bash",
	} {
		if !strings.Contains(readme, command) {
			t.Errorf("README lacks registered command %q", command)
		}
	}
}

func TestExecutableDocumentationScenarios(t *testing.T) {
	repo := documentationRepositoryRoot(t)
	srtPath := filepath.Join(repo, "testdata", "fixtures", "conversion", "srt-loss-free", "source", "input.srt")
	vttPath := filepath.Join(repo, "testdata", "fixtures", "conversion", "webvtt-loss-free", "source", "input.vtt")
	minimalVTTPath := filepath.Join(repo, "testdata", "fixtures", "webvtt", "minimal", "source", "minimal.vtt")
	directory := t.TempDir()
	documentPath := filepath.Join(directory, "document.cueson.json")
	restoredPath := filepath.Join(directory, "restored.srt")
	renderedPath := filepath.Join(directory, "rendered.srt")
	convertedVTTPath := filepath.Join(directory, "converted.vtt")
	convertedSRTPath := filepath.Join(directory, "converted.srt")
	schemaPath := filepath.Join(directory, "cueson.schema.json")

	status, stdout, stderr := runForTest(context.Background(), []string{"version"})
	requireDocumentedResult(t, "source execution", status, stdout, stderr, ExitSuccess, "1.0.0\n", "")

	status, stdout, stderr = runForTest(context.Background(), []string{"encode", "--pretty", "--output", documentPath, srtPath})
	requireDocumentedResult(t, "encode", status, stdout, stderr, ExitSuccess, "", "")
	encoded := readDocumentationBytes(t, documentPath)
	if !json.Valid(encoded) {
		t.Fatal("encode output is not JSON")
	}
	if _, err := schema.Decode(encoded); err != nil {
		t.Fatalf("encode output does not satisfy Cue JSON: %v", err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"restore", "--no-metadata", "--output", restoredPath, documentPath})
	requireDocumentedResult(t, "restore", status, stdout, stderr, ExitSuccess, "", "")
	if !bytes.Equal(readDocumentationBytes(t, restoredPath), readDocumentationBytes(t, srtPath)) {
		t.Fatal("exact restoration differs from the governed source")
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"render", "--to", "srt", "--output", renderedPath, documentPath})
	requireDocumentedResult(t, "render", status, stdout, stderr, ExitSuccess, "", "")
	if _, err := subrip.Parse(string(readDocumentationBytes(t, renderedPath)), subrip.Options{}); err != nil {
		t.Fatalf("rendered SubRip does not reparse: %v", err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"convert", "--strict", "--no-speaker-detection", "--to", "vtt", "--output", convertedVTTPath, srtPath})
	requireDocumentedResult(t, "SubRip to WebVTT conversion", status, stdout, stderr, ExitSuccess, "", "")
	if _, err := webvtt.Parse(string(readDocumentationBytes(t, convertedVTTPath))); err != nil {
		t.Fatalf("converted WebVTT does not reparse: %v", err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"convert", "--strict", "--to", "srt", "--output", convertedSRTPath, vttPath})
	requireDocumentedResult(t, "WebVTT to SubRip conversion", status, stdout, stderr, ExitSuccess, "", "")
	if _, err := subrip.Parse(string(readDocumentationBytes(t, convertedSRTPath)), subrip.Options{}); err != nil {
		t.Fatalf("converted SubRip does not reparse: %v", err)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"validate", minimalVTTPath})
	if status != ExitSuccess || stdout != "" || !strings.Contains(stderr, "valid native_subtitle.webvtt") {
		t.Fatalf("validate = (%d, %q, %q), want status 0, empty stdout, and one success diagnostic", status, stdout, stderr)
	}

	status, stdout, stderr = runForTest(context.Background(), []string{"inspect", minimalVTTPath})
	if status != ExitSuccess || stdout == "" || stderr != "" || !strings.Contains(stdout, "Format: webvtt") {
		t.Fatalf("human inspect = (%d, %q, %q)", status, stdout, stderr)
	}
	assertDocumentedPrivacy(t, stdout, directory, minimalVTTPath)

	status, stdout, stderr = runForTest(context.Background(), []string{"inspect", "--json", minimalVTTPath})
	if status != ExitSuccess || stderr != "" || !json.Valid([]byte(stdout)) {
		t.Fatalf("JSON inspect = (%d, %q, %q)", status, stdout, stderr)
	}
	assertDocumentedPrivacy(t, stdout, directory, minimalVTTPath)

	status, stdout, stderr = runForTest(context.Background(), []string{"schema"})
	requireDocumentedResult(t, "schema stdout", status, stdout, stderr, ExitSuccess, string(schema.Bytes()), "")
	status, stdout, stderr = runForTest(context.Background(), []string{"schema", "--version"})
	requireDocumentedResult(t, "schema version", status, stdout, stderr, ExitSuccess, "1.0.0\n", "")
	status, stdout, stderr = runForTest(context.Background(), []string{"schema", "--output", schemaPath})
	requireDocumentedResult(t, "schema file", status, stdout, stderr, ExitSuccess, "", "")
	if !bytes.Equal(readDocumentationBytes(t, schemaPath), schema.Bytes()) {
		t.Fatal("schema file differs from embedded canonical bytes")
	}

	for _, shell := range []string{"bash", "powershell"} {
		status, stdout, stderr = runForTest(context.Background(), []string{"completion", shell})
		if status != ExitSuccess || stderr != "" || stdout == "" {
			t.Fatalf("%s completion = (%d, %q, %q)", shell, status, stdout, stderr)
		}
	}
}

func documentationRepositoryRoot(t *testing.T) string {
	t.Helper()
	_, name, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve documentation test source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(name), "..", ".."))
}

func readDocumentationFile(t *testing.T, name string) string {
	t.Helper()
	return string(readDocumentationBytes(t, name))
}

func readDocumentationBytes(t *testing.T, name string) []byte {
	t.Helper()
	payload, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func markdownSection(document, heading string) string {
	start := strings.Index(document, heading)
	if start < 0 {
		return ""
	}
	remainder := document[start+len(heading):]
	if end := strings.Index(remainder, "\n## "); end >= 0 {
		return remainder[:end]
	}
	return remainder
}

func requireDocumentedResult(t *testing.T, name string, status int, stdout, stderr string, wantStatus int, wantStdout, wantStderr string) {
	t.Helper()
	if status != wantStatus || stdout != wantStdout || stderr != wantStderr {
		t.Fatalf("%s = (%d, %q, %q), want (%d, %q, %q)", name, status, stdout, stderr, wantStatus, wantStdout, wantStderr)
	}
}

func assertDocumentedPrivacy(t *testing.T, output string, forbidden ...string) {
	t.Helper()
	for _, value := range forbidden {
		if value != "" && strings.Contains(output, value) {
			t.Fatalf("inspection output leaked local identifier %q", value)
		}
	}
}
