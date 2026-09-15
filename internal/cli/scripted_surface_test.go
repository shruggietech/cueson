package cli

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestScriptedInputCatalogueHelpAndCommandSpecificCompletionValues(t *testing.T) {
	t.Parallel()
	want := []string{"auto", "cueson", "srt", "vtt", "ass", "ssa"}
	for _, name := range []string{"validate", "inspect"} {
		command, _ := lookupCommandSurface(name)
		var values []string
		for _, option := range command.Options {
			if slices.Contains(option.Spellings, "--format") {
				values = option.Values
			}
		}
		if !slices.Equal(values, want) {
			t.Fatalf("%s format vocabulary = %q, want %q", name, values, want)
		}
		for _, value := range values {
			if _, known := normalizeInputFormat(value); !known {
				t.Errorf("%s advertises unaccepted selector %q", name, value)
			}
		}
		help, _ := fullHelpText(name)
		if !strings.Contains(help, "ass, or ssa") || !strings.Contains(help, "experimental ASS/SSA") || !strings.Contains(help, "ASS/SSA accept UTF-8 only") {
			t.Errorf("%s help omits scripted vocabulary/status/profile", name)
		}
		for _, shell := range completionShells {
			payload, _ := completionScript(shell)
			text := string(payload)
			// Check the command-specific value line, not a global 'ass' already
			// present for encode/render/convert. Also cover joined --format=.
			var markers []string
			switch shell {
			case "bash":
				markers = []string{name + ":--format)", "\"$command\" == '" + name + "' && \"$current\" == '--format='"}
			case "zsh":
				markers = []string{name + ":--format)", "$command == '" + name + "' &&"}
			case "fish":
				markers = []string{"__fish_cueson_command_is " + name + "' -l format"}
			case "powershell":
				markers = []string{"$command -ceq '" + name + "' -and $previous -ceq '--format'", "$command -ceq '" + name + "' -and $wordToComplete.StartsWith('--format='"}
			}
			for _, marker := range markers {
				var matchingLine string
				lines := strings.Split(text, "\n")
				for index, line := range lines {
					if strings.Contains(line, marker) {
						matchingLine = line
						if shell == "bash" && strings.Contains(marker, "&&") {
							matchingLine = strings.Join(lines[index:min(index+6, len(lines))], "\n")
							if !strings.Contains(matchingLine, "COMPREPLY[index]='--format='") {
								t.Errorf("bash %s joined completion omits the literal option prefix", name)
							}
						}
						break
					}
				}
				if matchingLine == "" {
					t.Errorf("%s %s completion missing context %q", shell, name, marker)
					continue
				}
				for _, value := range want {
					candidate := "'" + value + "'"
					if shell == "fish" {
						candidate = value
					} else if strings.Contains(marker, "StartsWith") || (shell != "bash" && strings.Contains(marker, "&&")) {
						candidate = "'--format=" + value + "'"
					}
					if !strings.Contains(matchingLine, candidate) {
						t.Errorf("%s %s context %q omits %q", shell, name, marker, candidate)
					}
				}
			}
		}
		status, stdout, stderr := runForTest(context.Background(), []string{name, "missing.input", "--format", "unsupported"})
		if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "--format must be auto, cueson, srt, vtt, ass, or ssa") {
			t.Errorf("%s selector failure = %d %q %q", name, status, stdout, stderr)
		}
	}
}

func TestSelectorDiagnosticDerivesValuesWithoutMutatingCatalogue(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"validate", "inspect"} {
		command, _ := lookupCommandSurface(name)
		for _, option := range command.Options {
			if !slices.Contains(option.Spellings, "--format") {
				continue
			}
			before := slices.Clone(option.Values)
			want := fmt.Sprintf("--format must be %s, or %s", strings.Join(before[:len(before)-1], ", "), before[len(before)-1])
			if got := invalidSelectorError(name, "--format").Error(); got != want || !slices.Equal(before, option.Values) {
				t.Fatalf("%s catalogue diagnostic = %q, want %q without mutation", name, got, want)
			}
		}
	}
	if got := invalidSelectorError("unknown", "--format").Error(); got != "--format has no supported selector values" {
		t.Fatalf("missing catalogue fallback = %q", got)
	}
}
