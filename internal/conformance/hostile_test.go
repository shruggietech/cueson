package conformance_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/testutil"
)

func TestHostileNativeInputsFailWithoutPayloadOrPathDisclosure(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fixture string
		args    func(string) []string
	}{
		{name: "SubRip coordinate grammar", fixture: "malformed/subrip/incomplete-coordinates/input/incomplete.srt", args: func(path string) []string { return []string{"validate", "--format", "srt", path} }},
		{name: "WebVTT empty payload", fixture: "malformed/webvtt/empty-payload/input/empty-payload.vtt", args: func(path string) []string { return []string{"validate", "--format", "vtt", path} }},
		{name: "conversion rejects malformed source", fixture: "malformed/subrip/incomplete-coordinates/input/incomplete.srt", args: func(path string) []string { return []string{"convert", "--from", "srt", "--to", "vtt", path} }},
	}
	root := fixtureRoot(t)
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			sourcePath := filepath.Join(root, filepath.FromSlash(test.fixture))
			payload, err := os.ReadFile(sourcePath)
			if err != nil {
				t.Fatal(err)
			}
			directory := t.TempDir()
			inputPath := filepath.Join(directory, filepath.Base(sourcePath))
			if err := os.WriteFile(inputPath, payload, 0o600); err != nil {
				t.Fatal(err)
			}
			before, err := os.ReadDir(directory)
			if err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			status := cli.Run(context.Background(), test.args(inputPath), strings.NewReader(""), &stdout, &stderr)
			if status != cli.ExitRuntimeFailure || stdout.Len() != 0 {
				t.Fatalf("status=%d stdout=%q stderr=%q", status, stdout.String(), stderr.String())
			}
			after, err := os.ReadDir(directory)
			if err != nil {
				t.Fatal(err)
			}
			if len(after) != len(before) {
				t.Fatalf("workflow changed directory entries: before=%d after=%d", len(before), len(after))
			}
			if err := testutil.CheckNoForbidden("hostile/"+test.name, "diagnostics", stderr.Bytes(), directory, inputPath, string(payload)); err != nil {
				t.Fatal(err)
			}
		})
	}
}
