package cli

import (
	"bytes"
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestCompletionSyntaxWithAvailableInterpreters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		shell       string
		executables []string
		arguments   []string
	}{
		{shell: "bash", executables: []string{"bash"}, arguments: []string{"-n"}},
		{shell: "zsh", executables: []string{"zsh"}, arguments: []string{"-n"}},
		{shell: "fish", executables: []string{"fish"}, arguments: []string{"--no-execute"}},
		{shell: "powershell", executables: []string{"pwsh", "powershell"}, arguments: []string{"-NoLogo", "-NoProfile", "-NonInteractive", "-Command", "$source=[Console]::In.ReadToEnd(); $tokens=$null; $errors=$null; [System.Management.Automation.Language.Parser]::ParseInput($source, [ref]$tokens, [ref]$errors) > $null; if ($errors.Count -ne 0) { $errors | ForEach-Object { [Console]::Error.WriteLine($_.Message) }; exit 1 }"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.shell, func(t *testing.T) {
			t.Parallel()
			executable := firstAvailableExecutable(tt.executables)
			if executable == "" {
				t.Skip("supported shell interpreter is not installed on this host")
			}
			payload, ok := completionScript(tt.shell)
			if !ok {
				t.Fatalf("completionScript(%q) was rejected", tt.shell)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, executable, tt.arguments...)
			configureTestCommand(command)
			command.Stdin = bytes.NewReader(payload)
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("%s syntax check failed: %v\n%s", tt.shell, err, output)
			}
		})
	}
}

func firstAvailableExecutable(names []string) string {
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}
