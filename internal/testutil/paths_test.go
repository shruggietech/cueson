package testutil

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func TestPortablePathKey(t *testing.T) {
	t.Parallel()

	composed, err := PortablePathKey("fixtures/café/input.txt")
	if err != nil {
		t.Fatal(err)
	}
	decomposed, err := PortablePathKey("fixtures/cafe\u0301/input.txt")
	if err != nil {
		t.Fatal(err)
	}
	if composed != decomposed {
		t.Fatalf("canonical keys differ: %q != %q", composed, decomposed)
	}
}

func TestPortablePathKeyRejectsUnsafePaths(t *testing.T) {
	t.Parallel()

	paths := []string{
		"",
		"/absolute/input.txt",
		`C:\absolute\input.txt`,
		`fixtures\backslash.txt`,
		"fixtures/../outside.txt",
		"fixtures/./input.txt",
		"fixtures//input.txt",
		"fixtures/CON.txt",
		"fixtures/CLOCK$",
		"fixtures/conin$.txt",
		"fixtures/CONOUT$.log",
		"fixtures/trailing. ",
		"fixtures/question?.txt",
		"fixtures/" + strings.Repeat("a", 256),
	}
	for _, candidate := range paths {
		candidate := candidate
		t.Run(candidate, func(t *testing.T) {
			t.Parallel()
			if _, err := PortablePathKey(candidate); err == nil {
				t.Fatalf("PortablePathKey(%q) succeeded", candidate)
			}
		})
	}
}

func TestCheckNoForbiddenDetectsPortableVariants(t *testing.T) {
	t.Parallel()

	root := filepath.Join(`C:\Users`, "operator", "fixture root")
	variants := []string{
		root,
		filepath.ToSlash(root),
		strings.ReplaceAll(root, `\`, `/`),
		strings.ReplaceAll(root, `/`, `\`),
	}
	escaped, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	variants = append(variants, string(escaped[1:len(escaped)-1]))
	for _, variant := range variants {
		if err := CheckNoForbidden("fixture-path", "portable_output", []byte("prefix "+variant+" suffix"), root); err == nil || err.Error() != "fixture-path portable_output: forbidden local identifier detected" {
			t.Fatalf("CheckNoForbidden(%q) error = %v", variant, err)
		}
	}
	if err := CheckNoForbidden("fixture-path", "portable_output", []byte("https://example.test/C:/dialogue"), root); err != nil {
		t.Fatalf("CheckNoForbidden(generic path-like text) error = %v", err)
	}
	slashRoot := "C:/Users/operator/fixture root"
	if err := CheckNoForbidden("fixture-path", "portable_output", []byte(`{"path":"C:\\Users\\operator\\fixture root"}`), slashRoot); err == nil {
		t.Fatal("CheckNoForbidden() missed JSON-escaped backslash variant")
	}
}

func TestPortableOutputIsRootIndependent(t *testing.T) {
	t.Parallel()

	leftRoot := t.TempDir()
	rightRoot := t.TempDir()
	left := []byte(`{"fixture_id":"source-envelope/basic-lf","surface":"model"}`)
	right := append([]byte(nil), left...)
	if string(left) != string(right) {
		t.Fatal("equivalent portable outputs differ")
	}
	if err := CheckNoForbidden("source-envelope/basic-lf", "portable_output", left, leftRoot, rightRoot); err != nil {
		t.Fatal(err)
	}
}
