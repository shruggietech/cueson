package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
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
		for _, deferred := range []string{"schema", "restore", "encode", "render", "convert", "validate", "inspect", "completion"} {
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
		{name: "unknown command", args: []string{"schema"}, wantError: `unknown command "schema"`, wantUsage: "cueson [global options] <command>"},
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
