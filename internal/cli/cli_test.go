package cli

import (
	"bytes"
	"context"
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
			if stdout != "0.0.0\n" {
				t.Errorf("Run() stdout = %q, want %q", stdout, "0.0.0\\n")
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
		for _, deferred := range []string{"restore", "encode", "render", "convert", "validate", "inspect", "completion"} {
			if strings.Contains(stdout, "\n  "+deferred) {
				t.Errorf("Run(%q) help exposes deferred command %q:\n%s", args, deferred, stdout)
			}
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
		if status != ExitSuccess || stdout != "0.0.0\n" || stderr != "" {
			t.Errorf("Run(%q) = (%d, %q, %q), want (0, %q, empty)", args, status, stdout, stderr, "0.0.0\\n")
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
