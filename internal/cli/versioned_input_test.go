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

// The fixture is the complete published v1.0.0 annotation example, not a
// current instance relabeled to bypass historical contract selection.
func historicalSubRipPayload(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("../../schema/releases/v1.0.0/cueson.schema.json")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Examples []json.RawMessage `json:"examples"`
	}
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatal(err)
	}
	if len(contract.Examples) != 1 {
		t.Fatal("historical example missing")
	}
	return append([]byte(nil), contract.Examples[0]...)
}

func TestHistoricalInputCommandsPreserveIdentityAndSource(t *testing.T) {
	t.Parallel()
	payload := historicalSubRipPayload(t)
	document, err := schema.Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	input := filepath.Join(directory, "misleading.vtt")
	if err := os.WriteFile(input, payload, 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{"validate", input},
		{"inspect", input, "--json"},
		{"render", input, "--to", "srt"},
		{"convert", input, "--to", "vtt"},
		{"restore", input, "--no-metadata", "--output", filepath.Join(directory, "restored.srt")},
	} {
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitSuccess {
			t.Fatalf("%v: status=%d stderr=%s", args, status, stderr)
		}
		if args[0] == "inspect" {
			var report inspectionReport
			if err := json.Unmarshal([]byte(stdout), &report); err != nil {
				t.Fatal(err)
			}
			if report.Schema.ID != document.Schema || report.Schema.Version != "1.0.0" {
				t.Fatalf("historical identity = %#v", report.Schema)
			}
		}
	}
	restored, err := os.ReadFile(filepath.Join(directory, "restored.srt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != "1\n00:00:01,250 --> 00:00:04,200\n<i>Hello, world.</i>\n" {
		t.Fatalf("restored source changed: %q", restored)
	}
	after, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(payload, after) {
		t.Fatal("historical input was mutated")
	}
}

func TestHistoricalIntegrityAndIdentityFailuresDoNotPublish(t *testing.T) {
	t.Parallel()
	payload := historicalSubRipPayload(t)
	for name, corrupt := range map[string][]byte{
		"hash":    bytes.Replace(payload, []byte("955096d8"), []byte("055096d8"), 1),
		"version": bytes.Replace(payload, []byte(`"schema_version": "1.0.0"`), []byte(`"schema_version": "1.1.0"`), 1),
		"remote":  bytes.Replace(payload, []byte("https://cueson.io/schema/v1.0.0/cueson.schema.json"), []byte("https://unbundled.invalid/schema.json"), 1),
	} {
		t.Run(name, func(t *testing.T) {
			directory := t.TempDir()
			input := filepath.Join(directory, "misleading.srt")
			output := filepath.Join(directory, "output")
			if err := os.WriteFile(input, corrupt, 0600); err != nil {
				t.Fatal(err)
			}
			for _, command := range []string{"render", "convert", "restore"} {
				if err := os.WriteFile(output, []byte("sentinel"), 0600); err != nil {
					t.Fatal(err)
				}
				args := []string{command, input, "--output", output, "--force"}
				if command == "restore" {
					args = append(args, "--no-metadata")
				} else {
					args = append(args, "--to", "srt")
				}
				status, stdout, stderr := runForTest(context.Background(), args)
				if status != ExitRuntimeFailure || stdout != "" || strings.Contains(stderr, directory) {
					t.Fatalf("%v = (%d, %q, %q)", args, status, stdout, stderr)
				}
				after, err := os.ReadFile(output)
				if err != nil {
					t.Fatal(err)
				}
				if string(after) != "sentinel" {
					t.Fatalf("%s published invalid input", command)
				}
			}
		})
	}
}

func TestFrozenHistoricalBothFormatCommandMatrix(t *testing.T) {
	t.Parallel()
	for _, format := range []struct{ name, target, other string }{
		{"subrip", "srt", "vtt"}, {"webvtt", "vtt", "srt"},
	} {
		t.Run(format.name, func(t *testing.T) {
			payload, err := os.ReadFile("../schema/testdata/historical-v1.0.0-" + format.name + ".json")
			if err != nil {
				t.Fatal(err)
			}
			document, err := schema.Decode(payload)
			if err != nil {
				t.Fatal(err)
			}
			expected, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
			if err != nil {
				t.Fatal(err)
			}
			directory := t.TempDir()
			input := filepath.Join(directory, "misleading."+format.other)
			output := filepath.Join(directory, "restored."+format.target)
			if err := os.WriteFile(input, payload, 0600); err != nil {
				t.Fatal(err)
			}
			for _, args := range [][]string{
				{"validate", input}, {"inspect", input, "--json"},
				{"render", input, "--to", format.target},
				{"convert", input, "--to", format.other},
				{"restore", input, "--no-metadata", "--output", output},
			} {
				status, stdout, stderr := runForTest(context.Background(), args)
				if status != ExitSuccess {
					t.Fatalf("%v: status=%d stderr=%s", args, status, stderr)
				}
				if args[0] == "inspect" {
					var report inspectionReport
					if err := json.Unmarshal([]byte(stdout), &report); err != nil {
						t.Fatal(err)
					}
					if report.Schema.Version != "1.0.0" || report.Schema.ID != document.Schema || report.Integrity.Status != "verified" {
						t.Fatalf("historical report = %#v", report)
					}
				}
			}
			restored, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(restored, expected) {
				t.Fatal("historical source bytes changed")
			}
			for _, command := range []string{"render", "convert"} {
				target := format.target
				if command == "convert" {
					target = format.other
				}
				status, stdout, _ := runForTest(context.Background(), []string{command, input, "--to", target, "--strict"})
				if status != ExitRuntimeFailure || stdout != "" {
					t.Fatalf("historical strict loss bypass: %s (%d, %q)", command, status, stdout)
				}
			}
			after, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(after, payload) {
				t.Fatal("historical identity/producer/native input changed")
			}
		})
	}
}
