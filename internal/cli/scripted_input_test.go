package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/schema"
)

func TestScriptedGenericCommandsAndUnavailableNativePublication(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			payload, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
			if err != nil {
				t.Fatal(err)
			}
			document, err := schema.Decode(payload)
			if err != nil {
				t.Fatal(err)
			}
			directory := t.TempDir()
			input := filepath.Join(directory, "misleading.srt")
			output := filepath.Join(directory, "output."+format)
			if err := os.WriteFile(input, payload, 0600); err != nil {
				t.Fatal(err)
			}
			for _, command := range []string{"validate", "inspect"} {
				args := []string{command, input}
				if command == "inspect" {
					args = append(args, "--json")
				}
				status, stdout, stderr := runForTest(context.Background(), args)
				if status != ExitSuccess {
					t.Fatalf("%v = (%d, %q, %q)", args, status, stdout, stderr)
				}
				if command == "inspect" {
					var report inspectionReport
					if err := json.Unmarshal([]byte(stdout), &report); err != nil {
						t.Fatal(err)
					}
					capabilities := report.Capabilities
					if report.Format != format || report.Schema.Version != "1.1.0-dev" || capabilities.Declared.Status != "schema_only" || capabilities.Installed.Ingest || capabilities.Installed.Render || !capabilities.Installed.Restore || !capabilities.Installed.Validate || !capabilities.Installed.Inspect {
						t.Fatalf("scripted report = %#v", report)
					}
				}
			}
			status, stdout, stderr := runForTest(context.Background(), []string{"restore", input, "--no-metadata", "--output", output})
			if status != ExitSuccess || stdout != "" {
				t.Fatalf("restore = (%d, %q, %q)", status, stdout, stderr)
			}
			expected, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
			if err != nil {
				t.Fatal(err)
			}
			restored, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(expected, restored) {
				t.Fatal("scripted restore changed source bytes")
			}
			for _, args := range [][]string{
				{"render", input, "--to", format, "--output", output, "--force"},
				{"convert", input, "--to", "srt", "--output", output, "--force"},
			} {
				status, stdout, stderr := runForTest(context.Background(), args)
				if status != ExitRuntimeFailure || stdout != "" || strings.Contains(stderr, directory) {
					t.Fatalf("%v = (%d, %q, %q)", args, status, stdout, stderr)
				}
				if args[0] == "render" && !strings.Contains(stderr, "recognized but its native render capability is unavailable") {
					t.Fatalf("untruthful capability diagnostic: %s", stderr)
				}
				after, err := os.ReadFile(output)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(expected, after) {
					t.Fatal("unavailable codec published output")
				}
			}
		})
	}
}
