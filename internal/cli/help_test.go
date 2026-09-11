package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestHelpMatchesExactGoldens(t *testing.T) {
	t.Parallel()

	commands := append([]string{"root"}, commandNames()...)
	for _, name := range commands {
		name := name
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			command := name
			if command == "root" {
				command = ""
			}
			got, ok := fullHelpText(command)
			if !ok {
				t.Fatalf("fullHelpText(%q) did not find a surface", command)
			}
			want, err := os.ReadFile(filepath.Join("testdata", "help", name+".golden"))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal([]byte(got), want) {
				t.Fatalf("%s help differs from its exact golden", name)
			}
			if !utf8.ValidString(got) || strings.HasPrefix(got, "\uFEFF") || strings.Contains(got, "\r") || !strings.HasSuffix(got, "\n") {
				t.Fatalf("%s help is not UTF-8, BOM-free, LF-only text with a final LF", name)
			}
			for _, heading := range []string{"Global options:", "Streams:", "Exit codes:", "Examples:"} {
				if !strings.Contains(got, heading) {
					t.Errorf("%s help lacks %q", name, heading)
				}
			}
		})
	}
}

func TestShortUsageIsBoundedAndDeterministic(t *testing.T) {
	t.Parallel()

	for _, command := range append([]string{""}, commandNames()...) {
		first, ok := shortUsageText(command)
		if !ok {
			t.Fatalf("shortUsageText(%q) did not find a surface", command)
		}
		second, _ := shortUsageText(command)
		if first != second || !strings.HasPrefix(first, "Usage:\n  cueson ") || strings.Count(first, "\n") != 2 {
			t.Errorf("shortUsageText(%q) = %q", command, first)
		}
	}
	if _, ok := fullHelpText("missing"); ok {
		t.Error("fullHelpText accepted an unknown command")
	}
}
