package conformance_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/testutil"
)

func TestGovernedCorpus(t *testing.T) {
	t.Parallel()

	root := fixtureRoot(t)
	manifest, err := testutil.VerifyFixtures(root)
	if err != nil {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
	if len(manifest.Fixtures) < 16 {
		t.Fatalf("fixture count = %d, want at least 16", len(manifest.Fixtures))
	}
}

func TestAcceptedSourceEnvelopeConformance(t *testing.T) {
	t.Parallel()

	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	fixture := mustFixture(t, manifest, "source-envelope/basic-lf")
	sourceArtifact := mustArtifact(t, fixture, "source")
	modelArtifact := mustArtifact(t, fixture, "expected_model")
	diagnosticsArtifact := mustArtifact(t, fixture, "expected_diagnostics")

	sourceBytes := readArtifact(t, root, sourceArtifact)
	modelBytes := readArtifact(t, root, modelArtifact)
	diagnosticBytes := readArtifact(t, root, diagnosticsArtifact)
	document, err := schema.Decode(modelBytes)
	if err != nil {
		t.Fatalf("schema.Decode() error = %v", err)
	}
	if err := testutil.CompareJSON(fixture.ID, "model", modelBytes, document); err != nil {
		t.Fatal(err)
	}
	var expectedDiagnostics []model.Diagnostic
	if err := json.Unmarshal(diagnosticBytes, &expectedDiagnostics); err != nil {
		t.Fatalf("decode expected diagnostics: %v", err)
	}
	if err := testutil.CompareDiagnostics(fixture.ID, "diagnostics", expectedDiagnostics, document.Diagnostics); err != nil {
		t.Fatal(err)
	}

	decodedSource, err := base64.StdEncoding.Strict().DecodeString(document.Source.Assets[0].DataBase64)
	if err != nil {
		t.Fatalf("decode source envelope: %v", err)
	}
	if err := testutil.CompareBytes(fixture.ID, "source_bytes", sourceBytes, decodedSource); err != nil {
		t.Fatal(err)
	}
	if err := testutil.CompareIntegrity(fixture.ID, "source_integrity", decodedSource, *sourceArtifact.SizeBytes, sourceArtifact.SHA256); err != nil {
		t.Fatal(err)
	}

	output := filepath.Join(t.TempDir(), "restored.srt")
	report, err := source.Restore(context.Background(), document, source.RestoreOptions{Output: output})
	if err != nil {
		t.Fatalf("source.Restore() error = %v", err)
	}
	restored := readFile(t, output)
	if err := testutil.CompareBytes(fixture.ID, "restored_bytes", sourceBytes, restored); err != nil {
		t.Fatal(err)
	}
	if len(report.Assets) != 1 {
		t.Fatalf("restoration asset count = %d, want 1", len(report.Assets))
	}
	if err := testutil.CompareIntegrity(fixture.ID, "restored_integrity", restored, report.Assets[0].Bytes, report.Assets[0].SHA256); err != nil {
		t.Fatal(err)
	}
	timestampOutcomes := portableTimestamps(report.Assets[0].Timestamps)
	if err := testutil.CompareTimestamps(fixture.ID, "timestamps", expectedTimestamps(), timestampOutcomes); err != nil {
		t.Fatal(err)
	}

	portable, err := json.Marshal(struct {
		FixtureID  string                      `json:"fixture_id"`
		Model      model.Document              `json:"model"`
		Timestamps []testutil.TimestampOutcome `json:"timestamps"`
	}{FixtureID: fixture.ID, Model: document, Timestamps: timestampOutcomes})
	if err != nil {
		t.Fatal(err)
	}
	home, _ := os.UserHomeDir()
	host, _ := os.Hostname()
	currentUser, err := user.Current()
	if err != nil {
		t.Fatalf("resolve current user: %v", err)
	}
	userSentinels := []string{currentUser.Username, filepath.Base(home)}
	if separator := strings.LastIndexAny(currentUser.Username, `\\/`); separator >= 0 {
		userSentinels = append(userSentinels, currentUser.Username[separator+1:])
	}
	if err := testutil.CheckNoForbidden(fixture.ID, "portable_expectation", portable, append([]string{root, filepath.Dir(output), home, host}, userSentinels...)...); err != nil {
		t.Fatal(err)
	}
	proveEquivalentRoots(t, fixture.ID, modelBytes)
}

func TestMalformedRegressionsAreDeterministic(t *testing.T) {
	t.Parallel()

	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, fixture := range manifest.Fixtures {
		if fixture.Expectation.Result != "rejected" {
			continue
		}
		fixture := fixture
		t.Run(fixture.ID, func(t *testing.T) {
			t.Parallel()
			artifact := mustArtifact(t, fixture, "input")
			input := readArtifact(t, root, artifact)
			original := append([]byte(nil), input...)
			var prior string
			for run := 0; run < 2; run++ {
				stage, diagnostic := rejectAtDeclaredBoundary(t, fixture, input, run)
				if stage != *fixture.Expectation.Stage {
					t.Fatalf("rejection stage = %q, want %q", stage, *fixture.Expectation.Stage)
				}
				if !strings.Contains(diagnostic, *fixture.Expectation.DiagnosticContains) {
					t.Fatalf("diagnostic = %q, want fragment %q", diagnostic, *fixture.Expectation.DiagnosticContains)
				}
				if run > 0 && diagnostic != prior {
					t.Fatalf("diagnostic changed between runs:\nfirst: %s\nsecond: %s", prior, diagnostic)
				}
				prior = diagnostic
			}
			if !bytes.Equal(input, original) {
				t.Fatal("malformed regression bytes changed during validation")
			}
		})
	}
}

func rejectAtDeclaredBoundary(t *testing.T, fixture testutil.Fixture, input []byte, run int) (string, string) {
	t.Helper()
	if fixture.ID == "subrip/reversed-time" {
		_, err := subrip.Parse(string(input), subrip.Options{DetectSpeakers: true})
		if err == nil {
			t.Fatal("subrip.Parse() accepted malformed fixture")
		}
		return "parse", err.Error()
	}
	if strings.HasPrefix(fixture.ID, "webvtt/") {
		decoded, err := webvtt.DecodeUTF8(input, "")
		if err == nil {
			_, err = webvtt.Parse(decoded.Text)
		}
		if err == nil {
			t.Fatal("WebVTT codec accepted malformed fixture")
		}
		diagnostic := err.Error()
		if strings.Contains(diagnostic, "webvtt_signature_invalid") {
			diagnostic = "WebVTT signature: " + diagnostic
		}
		return "parse", diagnostic
	}
	document, err := schema.Decode(input)
	if *fixture.Expectation.Stage != "integrity" {
		if err == nil {
			t.Fatal("schema.Decode() accepted malformed fixture")
		}
		return schemaErrorStage(err.Error()), err.Error()
	}
	if err != nil {
		t.Fatalf("schema.Decode() rejected integrity fixture before integrity stage: %v", err)
	}
	output := filepath.Join(t.TempDir(), "integrity-run-"+string(rune('0'+run))+".srt")
	_, err = source.Restore(context.Background(), document, source.RestoreOptions{Output: output, Metadata: source.MetadataNone})
	if err == nil {
		t.Fatal("source.Restore() accepted corrupt envelope")
	}
	if _, statErr := os.Lstat(output); !os.IsNotExist(statErr) {
		t.Fatalf("integrity failure left output: %v", statErr)
	}
	return "integrity", err.Error()
}

func schemaErrorStage(diagnostic string) string {
	switch {
	case strings.HasPrefix(diagnostic, "parse Cue JSON"):
		return "parse"
	case strings.HasPrefix(diagnostic, "validate Cue JSON structure"):
		return "structure"
	case strings.HasPrefix(diagnostic, "validate Cue JSON semantics"):
		return "semantics"
	default:
		return "unknown"
	}
}

func proveEquivalentRoots(t *testing.T, fixtureID string, modelBytes []byte) {
	t.Helper()
	roots := []string{t.TempDir(), t.TempDir()}
	outputs := make([][]byte, 0, len(roots))
	for _, root := range roots {
		path := filepath.Join(root, "model.json")
		if err := os.WriteFile(path, modelBytes, 0o600); err != nil {
			t.Fatal(err)
		}
		document, err := schema.Decode(readFile(t, path))
		if err != nil {
			t.Fatal(err)
		}
		portable, err := json.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		outputs = append(outputs, portable)
	}
	if err := testutil.CompareBytes(fixtureID, "equivalent_roots", outputs[0], outputs[1]); err != nil {
		t.Fatal(err)
	}
	if err := testutil.CheckNoForbidden(fixtureID, "equivalent_roots", outputs[0], roots...); err != nil {
		t.Fatal(err)
	}
}

func portableTimestamps(results []source.TimestampResult) []testutil.TimestampOutcome {
	if results == nil {
		return nil
	}
	outcomes := make([]testutil.TimestampOutcome, len(results))
	for index, result := range results {
		outcomes[index] = testutil.TimestampOutcome{Kind: string(result.Kind), Status: string(result.Status), EffectivePrecision: result.EffectivePrecision}
	}
	return outcomes
}

func expectedTimestamps() []testutil.TimestampOutcome {
	switch runtime.GOOS {
	case "windows":
		return []testutil.TimestampOutcome{
			{Kind: string(source.TimestampCreated), Status: testutil.TimestampRestored, EffectivePrecision: "100ns"},
			{Kind: string(source.TimestampModified), Status: testutil.TimestampRestored, EffectivePrecision: "100ns"},
			{Kind: string(source.TimestampAccessed), Status: testutil.TimestampRestored, EffectivePrecision: "100ns"},
		}
	case "linux":
		return []testutil.TimestampOutcome{
			{Kind: string(source.TimestampCreated), Status: testutil.TimestampUnsupported},
			{Kind: string(source.TimestampModified), Status: testutil.TimestampRestored, EffectivePrecision: "1ns"},
			{Kind: string(source.TimestampAccessed), Status: testutil.TimestampRestored, EffectivePrecision: "1ns"},
		}
	case "darwin":
		return []testutil.TimestampOutcome{
			{Kind: string(source.TimestampCreated), Status: testutil.TimestampUnsupported},
			{Kind: string(source.TimestampModified), Status: testutil.TimestampRestored, EffectivePrecision: "1us"},
			{Kind: string(source.TimestampAccessed), Status: testutil.TimestampRestored, EffectivePrecision: "1us"},
		}
	default:
		return []testutil.TimestampOutcome{
			{Kind: string(source.TimestampCreated), Status: testutil.TimestampUnsupported},
			{Kind: string(source.TimestampModified), Status: testutil.TimestampUnsupported},
			{Kind: string(source.TimestampAccessed), Status: testutil.TimestampUnsupported},
		}
	}
}

func fixtureRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate conformance test")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "testdata"))
}

func verifiedManifest(t *testing.T, root string) testutil.Manifest {
	t.Helper()
	manifest, err := testutil.VerifyFixtures(root)
	if err != nil {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
	return manifest
}

func mustFixture(t *testing.T, manifest testutil.Manifest, id string) testutil.Fixture {
	t.Helper()
	fixture, ok := manifest.FixtureByID(id)
	if !ok {
		t.Fatalf("fixture %q is missing", id)
	}
	return fixture
}

func mustArtifact(t *testing.T, fixture testutil.Fixture, id string) testutil.Artifact {
	t.Helper()
	artifact, ok := fixture.ArtifactByID(id)
	if !ok {
		t.Fatalf("fixture %q artifact %q is missing", fixture.ID, id)
	}
	return artifact
}

func readArtifact(t *testing.T, root string, artifact testutil.Artifact) []byte {
	t.Helper()
	return readFile(t, testutil.ArtifactFile(root, artifact))
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
