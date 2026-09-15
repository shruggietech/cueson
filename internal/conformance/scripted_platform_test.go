package conformance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/convert"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/testutil"
)

func TestScriptedAcceptedFixtureRestorationConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	accepted := 0
	for _, fixture := range manifest.Fixtures {
		if fixture.Expectation.Result != "accepted" || (!strings.HasPrefix(fixture.ID, "scripted/") && !strings.HasPrefix(fixture.ID, "scripted-conversion/")) {
			continue
		}
		accepted++
		t.Run(fixture.ID, func(t *testing.T) {
			document, _ := encodeScriptedFixture(t, root, fixture)
			before := mustMarshalConversion(t, document)
			if err := source.ValidateIntegrity(context.Background(), document); err != nil {
				t.Fatal(err)
			}
			artifact := mustArtifact(t, fixture, "source")
			original := readArtifact(t, root, artifact)
			directory := t.TempDir()
			output := filepath.Join(directory, "restored.subtitle")
			if _, err := source.Restore(context.Background(), document, source.RestoreOptions{Output: output, Metadata: source.MetadataNone}); err != nil {
				t.Fatal(err)
			}
			restored := readFile(t, output)
			if err := testutil.CompareBytes(fixture.ID, "restored_source", original, restored); err != nil {
				t.Fatal(err)
			}
			if err := testutil.CompareIntegrity(fixture.ID, "restored_source", restored, *artifact.SizeBytes, artifact.SHA256); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, mustMarshalConversion(t, document)) || !bytes.Equal(original, readArtifact(t, root, artifact)) {
				t.Fatal("restoration changed source truth or model")
			}
			entries, err := os.ReadDir(directory)
			if err != nil || len(entries) != 1 || entries[0].Name() != filepath.Base(output) {
				t.Fatal("successful restoration left staging artifacts")
			}
			if err := testutil.CheckNoForbidden(fixture.ID, "restoration_model", before, root, directory, output, testutil.ArtifactFile(root, artifact)); err != nil {
				t.Fatal(err)
			}
		})
	}
	if accepted < 34 {
		t.Fatalf("accepted scripted/conversion records=%d, want all existing thirty-four", accepted)
	}
}

// The existing native CI jobs execute this filesystem workflow on Windows, macOS, and Linux.
func TestScriptedNativePlatformWorkflowConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			fixture := mustFixture(t, manifest, "scripted/"+format+"-bom-crlf")
			original := readArtifact(t, root, mustArtifact(t, fixture, "source"))
			var expected scriptedExpectation
			decodeScriptedExpected(t, readArtifact(t, root, mustArtifact(t, fixture, "expected_model")), &expected)
			directory := t.TempDir()
			input := filepath.Join(directory, "source."+format)
			if err := os.WriteFile(input, original, 0o600); err != nil {
				t.Fatal(err)
			}
			encoded, _ := runScriptedPlatformCLI(t, []string{"encode", input, "--stdout"}, cli.ExitSuccess)
			document, err := schema.Decode(encoded)
			if err != nil {
				t.Fatal(err)
			}
			jsonPath := filepath.Join(directory, "source.cueson.json")
			if err = os.WriteFile(jsonPath, encoded, 0o600); err != nil {
				t.Fatal(err)
			}
			for _, path := range []string{input, jsonPath} {
				validated, _ := runScriptedPlatformCLI(t, []string{"validate", path}, cli.ExitSuccess)
				if len(validated) != 0 {
					t.Fatal("validation published payload")
				}
				inspected, diagnostics := runScriptedPlatformCLI(t, []string{"inspect", path, "--json"}, cli.ExitSuccess)
				var report struct {
					Format   string         `json:"format"`
					Scripted map[string]int `json:"scripted"`
				}
				if err = json.Unmarshal(inspected, &report); err != nil {
					t.Fatal(err)
				}
				if report.Format != format || report.Scripted["section_count"] != len(expected.SectionNames) || report.Scripted["record_count"] != expected.RecordCount || report.Scripted["style_count"] != len(expected.Styles) || report.Scripted["event_count"] != len(expected.Events) || report.Scripted["dialogue_event_count"] != len(expected.Cues) || report.Scripted["attachment_count"] != len(expected.AttachmentNames) {
					t.Fatalf("native inspection counts=%v", report.Scripted)
				}
				if err = testutil.CheckNoForbidden(fixture.ID, "inspection", append(inspected, diagnostics...), directory, input, jsonPath, "Narrator", "Arial", "Hello"); err != nil {
					t.Fatal(err)
				}
			}
			rendered, _ := runScriptedPlatformCLI(t, []string{"render", jsonPath, "--to", format, "--output", "-"}, cli.ExitSuccess)
			if err = testutil.CompareBytes(fixture.ID, "native_canonical_render", readArtifact(t, root, mustArtifact(t, fixture, "expected_bytes")), rendered); err != nil {
				t.Fatal(err)
			}
			parsed, err := scripted.Parse(context.Background(), rendered, format)
			if err != nil {
				t.Fatal(err)
			}
			assertScriptedProjection(t, expected, parsed.Native, parsed.Cues)
			output := filepath.Join(directory, "restored."+format)
			if payload, _ := runScriptedPlatformCLI(t, []string{"restore", jsonPath, "--output", output, "--no-metadata"}, cli.ExitSuccess); len(payload) != 0 {
				t.Fatal("file restoration published stdout payload")
			}
			if !bytes.Equal(readFile(t, output), original) || !bytes.Equal(readFile(t, input), original) {
				t.Fatal("native filesystem workflow changed original bytes")
			}
			converted, _ := runScriptedPlatformCLI(t, []string{"convert", jsonPath, "--to", "vtt"}, cli.ExitSuccess)
			cues := reparseConversion(t, converted, "webvtt")
			if len(cues) != len(document.Cues) {
				t.Fatal("native conversion lost dialogue")
			}
			for index, cue := range cues {
				if cue.Timing != document.Cues[index].Timing || cue.Payload.PlainText != document.Cues[index].Payload.PlainText {
					t.Fatal("native conversion changed readable baseline")
				}
			}
			assertConversionPublicationRefusal(t, original, format, "webvtt", true)
			entries, err := os.ReadDir(directory)
			if err != nil || len(entries) != 3 {
				t.Fatal("native workflow left staging artifacts")
			}
		})
	}
}

func TestScriptedNativePlatformSafetyConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, format := range []string{"ass", "ssa"} {
		fixture := mustFixture(t, manifest, "scripted/"+format+"-basic")
		original, _ := encodeScriptedFixture(t, root, fixture)
		originalBytes := mustMarshalConversion(t, original)
		for name, mutation := range map[string]func(*model.Document){
			"corrupt hash":       func(document *model.Document) { document.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64) },
			"corrupt length":     func(document *model.Document) { document.Source.Assets[0].Size.Bytes++ },
			"unsafe source name": func(document *model.Document) { document.Source.Assets[0].FileName = "../forbidden.subtitle" },
			"fabricated capture": func(document *model.Document) {
				native := document.FormatData.ASS
				if format == "ssa" {
					native = document.FormatData.SSA
				}
				wrong := "[Different]"
				native.Sections[0].RawHeader = &wrong
			},
			"missing owner reference": func(document *model.Document) {
				native := document.Cues[0].FormatData.ASS
				if format == "ssa" {
					native = document.Cues[0].FormatData.SSA
				}
				native.EventID = "event-missing"
			},
			"inconsistent derived text": func(document *model.Document) { document.Cues[0].Payload.PlainText = "different projection" },
		} {
			t.Run(format+"/"+name, func(t *testing.T) {
				var bad model.Document
				if err := json.Unmarshal(originalBytes, &bad); err != nil {
					t.Fatal(err)
				}
				mutation(&bad)
				directory := t.TempDir()
				output := filepath.Join(directory, "preserved.subtitle")
				preserved := []byte("preserve existing bytes\n")
				if err := os.WriteFile(output, preserved, 0o600); err != nil {
					t.Fatal(err)
				}
				if _, err := source.Restore(context.Background(), bad, source.RestoreOptions{Output: output, Force: true, Metadata: source.MetadataNone}); err == nil {
					t.Fatal("invalid source/model restored")
				}
				if !bytes.Equal(readFile(t, output), preserved) {
					t.Fatal("invalid restore changed forced destination")
				}
				if result, err := scripted.Render(context.Background(), bad, false); err == nil || len(result.Bytes) != 0 {
					t.Fatal("invalid native/model rendered")
				}
				if result, err := convert.Convert(context.Background(), bad, "webvtt", convert.Options{}); err == nil || len(result.Bytes) != 0 {
					t.Fatal("invalid native/model converted")
				}
				jsonPath := filepath.Join(directory, "invalid.cueson.json")
				if err := os.WriteFile(jsonPath, mustMarshalConversion(t, bad), 0o600); err != nil {
					t.Fatal(err)
				}
				for _, command := range []string{"restore", "render", "convert"} {
					args := []string{command, jsonPath, "--output", output, "--force"}
					if command == "convert" {
						args = append(args, "--to", "vtt")
					}
					if command == "render" {
						args = append(args, "--to", format)
					}
					if command == "restore" {
						args = append(args, "--no-metadata")
					}
					stdout, stderr := runScriptedPlatformCLI(t, args, cli.ExitRuntimeFailure)
					if len(stdout) != 0 || !bytes.Equal(readFile(t, output), preserved) {
						t.Fatal("invalid command published or replaced payload")
					}
					if err := testutil.CheckNoForbidden(fixture.ID, "rejection", stderr, directory, jsonPath, output); err != nil {
						t.Fatal(err)
					}
				}
				entries, err := os.ReadDir(directory)
				if err != nil || len(entries) != 2 {
					t.Fatal("invalid workflow left staging artifacts")
				}
				if !reflect.DeepEqual(originalBytes, mustMarshalConversion(t, original)) {
					t.Fatal("adversarial test mutated its accepted source")
				}
			})
		}
	}
}

func runScriptedPlatformCLI(t *testing.T, args []string, want int) ([]byte, []byte) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	status := cli.Run(context.Background(), args, strings.NewReader(""), &stdout, &stderr)
	if status != want {
		t.Fatalf("%s status=%d want=%d stdout=%q stderr=%q", args[0], status, want, stdout.Bytes(), stderr.Bytes())
	}
	return append([]byte(nil), stdout.Bytes()...), append([]byte(nil), stderr.Bytes()...)
}
