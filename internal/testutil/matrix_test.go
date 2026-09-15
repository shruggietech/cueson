package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeConformanceMatrixRejectsTamperedMetadata(t *testing.T) {
	t.Parallel()
	valid := []byte(`{"version":1,"formats":[]}`)
	if _, err := DecodeConformanceMatrix(append(append([]byte(nil), valid...), []byte(" \n\t")...)); err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{
		"second value":             append(append([]byte(nil), valid...), []byte(" {}")...),
		"malformed trailing bytes": append(append([]byte(nil), valid...), []byte(" garbage{")...),
		"unknown property":         []byte(`{"version":1,"formats":[],"unapproved":true}`),
		"nested unknown property":  []byte(`{"version":1,"formats":[{"format":"x","rows":[],"unapproved":true}]}`),
		"BOM":                      append([]byte{0xef, 0xbb, 0xbf}, valid...),
		"invalid UTF-8":            append(append([]byte(nil), valid...), 0xff),
		"null matrix":              []byte("null"),
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := DecodeConformanceMatrix(data); err == nil {
				t.Fatal("tampered matrix accepted")
			}
		})
	}
}

func TestEvidenceFunctionRequiresGenuineKindCorrectPortableDeclaration(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "internal", "example")
	if err := os.MkdirAll(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	text := `package example
import alias "testing"
// func TestComment(t *alias.T) {}
const misleading = "func TestString(t *alias.T) {}"
func TestReal(t *alias.T) {}
func FuzzReal(f *alias.F) {}
func TestWrong(t *alias.F) {}
func TestMany(a,b *alias.T) {}
func TestReturning(t *alias.T) int { return 1 }
func Testlowercase(t *alias.T) {}
`
	if err := os.WriteFile(filepath.Join(directory, "example_test.go"), []byte(text), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		reference, kind string
		want            bool
	}{
		{"internal/example:TestReal", "render_test", true},
		{"internal/example:TestReal", "platform_test", true},
		{"internal/example:FuzzReal", "fuzz_target", true},
		{"internal/example:TestReal", "fuzz_target", false},
		{"internal/example:FuzzReal", "conversion_test", false},
		{"internal/example:TestComment", "render_test", false},
		{"internal/example:TestString", "render_test", false},
		{"internal/example:TestWrong", "render_test", false},
		{"internal/example:TestMany", "render_test", false},
		{"internal/example:TestReturning", "render_test", false},
		{"internal/example:Testlowercase", "render_test", false},
		{"../internal/example:TestReal", "render_test", false},
		{"internal/../internal/example:TestReal", "render_test", false},
		{`internal\example:TestReal`, "render_test", false},
		{"/internal/example:TestReal", "render_test", false},
		{"internal/example:TestMissing", "render_test", false},
	} {
		t.Run(test.reference+"/"+test.kind, func(t *testing.T) {
			got, err := EvidenceFunctionExists(root, test.reference, test.kind)
			if err != nil || got != test.want {
				t.Fatalf("exists=%t err=%v want=%t", got, err, test.want)
			}
		})
	}
}
