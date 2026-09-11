package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestActualExecutableHelpCompletionAndInvocationStreams(t *testing.T) {
	if testing.Short() {
		t.Skip("actual executable process test is omitted in short mode")
	}

	binaryName := "cueson"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	binary := filepath.Join(t.TempDir(), binaryName)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, "go", "build", "-o", binary, "./cmd/cueson")
	build.Dir = filepath.Clean(filepath.Join("..", ".."))
	configureTestCommand(build)
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build actual executable: %v\n%s", err, output)
	}

	helpCases := append([]string{""}, commandNames()...)
	for _, commandName := range helpCases {
		args := []string{"--help"}
		if commandName != "" {
			args = []string{commandName, "--help"}
		}
		stdout, stderr, status := runActualProcess(t, binary, args)
		want, _ := fullHelpText(commandName)
		if status != ExitSuccess || stdout != want || stderr != "" {
			t.Errorf("%s help = status %d, stdout_match %t, stderr %q", commandName, status, stdout == want, stderr)
		}
	}

	for _, shell := range completionShells {
		stdout, stderr, status := runActualProcess(t, binary, []string{"completion", shell})
		want, _ := completionScript(shell)
		if status != ExitSuccess || !bytes.Equal([]byte(stdout), want) || stderr != "" {
			t.Errorf("completion %s = status %d, stdout_match %t, stderr %q", shell, status, bytes.Equal([]byte(stdout), want), stderr)
		}
	}

	for _, commandName := range []string{"validate", "inspect"} {
		stdout, stderr, status := runActualProcess(t, binary, []string{commandName})
		usage, _ := shortUsageText(commandName)
		if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, usage) {
			t.Errorf("missing %s input = status %d, stdout %q, stderr %q", commandName, status, stdout, stderr)
		}
	}

	stdout, stderr, status := runActualProcess(t, binary, []string{"completion", "cmd"})
	usage, _ := shortUsageText("completion")
	if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, usage) {
		t.Errorf("unsupported completion = status %d, stdout %q, stderr %q", status, stdout, stderr)
	}
}

func runActualProcess(t *testing.T, binary string, args []string) (string, string, int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := exec.CommandContext(ctx, binary, args...)
	configureTestCommand(command)
	command.Env = append(os.Environ(), "NO_COLOR=1")
	command.Stdin = nil
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err == nil {
		return stdout.String(), stderr.String(), ExitSuccess
	}
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return stdout.String(), stderr.String(), exitError.ExitCode()
	}
	t.Fatalf("execute %q: %v", append([]string{binary}, args...), err)
	return "", "", -1
}
