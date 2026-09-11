package cli

import (
	"slices"
	"testing"
)

func TestOrderedSurfaceContainsEveryShippedCommand(t *testing.T) {
	t.Parallel()

	want := []string{"encode", "restore", "render", "convert", "validate", "inspect", "schema", "version", "completion"}
	got := make([]string, 0, len(orderedCommandSurface))
	for _, command := range orderedCommandSurface {
		got = append(got, command.Name)
		if command.Summary == "" || command.Description == "" || command.Usage == "" || len(command.Examples) == 0 {
			t.Errorf("command %q has an incomplete surface definition", command.Name)
		}
	}
	if !slices.Equal(got, want) {
		t.Fatalf("command order = %q, want %q", got, want)
	}
}

func TestSurfaceOptionVocabularyIsUniqueAndComplete(t *testing.T) {
	t.Parallel()

	wantGlobals := []string{"-q", "--quiet", "--silent", "--no-color", "-h", "--help", "--"}
	if got := optionSpellings(globalOptionSurface); !slices.Equal(got, wantGlobals) {
		t.Fatalf("global options = %q, want %q", got, wantGlobals)
	}

	wantLocal := map[string][]string{
		"encode":     {"-o", "--output", "-f", "--force", "--format", "--encoding", "--pretty", "--stdout", "--no-speaker-detection"},
		"restore":    {"-o", "--output", "--output-dir", "-f", "--force", "--strict-metadata", "--no-metadata"},
		"render":     {"--to", "-o", "--output", "-f", "--force", "--strict"},
		"convert":    {"--to", "--from", "--encoding", "-o", "--output", "-f", "--force", "--strict", "--no-speaker-detection"},
		"validate":   {"--format", "--encoding"},
		"inspect":    {"--format", "--encoding", "--json"},
		"schema":     {"--version", "-o", "--output", "-f", "--force"},
		"version":    {},
		"completion": {},
	}
	for command, want := range wantLocal {
		surface, ok := lookupCommandSurface(command)
		if !ok {
			t.Fatalf("lookupCommandSurface(%q) failed", command)
		}
		if got := optionSpellings(surface.Options); !slices.Equal(got, want) {
			t.Errorf("%s options = %q, want %q", command, got, want)
		}
	}
}

func optionSpellings(options []cliOptionSpec) []string {
	var spellings []string
	seen := make(map[string]struct{})
	for _, option := range options {
		if len(option.Spellings) == 0 || option.Description == "" {
			continue
		}
		for _, spelling := range option.Spellings {
			if _, duplicate := seen[spelling]; duplicate {
				spellings = append(spellings, "duplicate:"+spelling)
				continue
			}
			seen[spelling] = struct{}{}
			spellings = append(spellings, spelling)
		}
	}
	return spellings
}
