package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCompletionScriptsMatchExactGoldens(t *testing.T) {
	t.Parallel()

	for _, shell := range completionShells {
		shell := shell
		t.Run(shell, func(t *testing.T) {
			t.Parallel()
			first, ok := completionScript(shell)
			if !ok {
				t.Fatalf("completionScript(%q) was rejected", shell)
			}
			second, _ := completionScript(shell)
			if !bytes.Equal(first, second) {
				t.Fatalf("completionScript(%q) is not deterministic", shell)
			}
			want, err := os.ReadFile(filepath.Join("testdata", "completion", shell+".golden"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(first, want) {
				t.Fatalf("%s completion differs from its exact golden", shell)
			}
			text := string(first)
			if !utf8.Valid(first) || bytes.HasPrefix(first, []byte("\xEF\xBB\xBF")) || strings.Contains(text, "\r") || !strings.HasSuffix(text, "\n") {
				t.Fatalf("%s completion is not UTF-8, BOM-free, LF-only text with a final LF", shell)
			}
			assertSafeStaticCompletion(t, text)
			assertCompleteVocabulary(t, shell, text)
		})
	}
}

func TestCompletionShellSelectionIsExact(t *testing.T) {
	t.Parallel()

	for _, shell := range []string{"", "Bash", " bash", "bash ", "cmd", "pwsh"} {
		if payload, ok := completionScript(shell); ok || payload != nil {
			t.Errorf("completionScript(%q) = (%q, %t), want (nil, false)", shell, payload, ok)
		}
	}
}

func assertSafeStaticCompletion(t *testing.T, script string) {
	t.Helper()
	for _, forbidden := range []string{"Invoke-Expression", "eval ", "curl ", "wget ", "http://", "https://", "cueson __", "cueson completion", "Get-ChildItem", "CompleteFilename", "$(", "`cueson"} {
		if strings.Contains(script, forbidden) {
			t.Errorf("completion contains forbidden behavior %q", forbidden)
		}
	}
}

func assertCompleteVocabulary(t *testing.T, shell, script string) {
	t.Helper()
	for _, command := range orderedCommandSurface {
		if !strings.Contains(script, command.Name) {
			t.Errorf("completion omits command %q", command.Name)
		}
		for _, option := range append(append([]cliOptionSpec{}, globalOptionSurface...), command.Options...) {
			for _, spelling := range option.Spellings {
				needle := spelling
				if shell == "fish" {
					switch {
					case spelling == "--":
						needle = "-a '--'"
					case strings.HasPrefix(spelling, "--"):
						needle = "-l " + strings.TrimPrefix(spelling, "--")
					case strings.HasPrefix(spelling, "-"):
						needle = "-s " + strings.TrimPrefix(spelling, "-")
					}
				}
				if !strings.Contains(script, needle) {
					t.Errorf("completion omits %s option %q", command.Name, spelling)
				}
			}
		}
	}
}
