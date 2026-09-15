package conformance_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func stableCandidateSource(t *testing.T, format string) []byte {
	t.Helper()
	switch format {
	case "srt":
		return []byte("1\n00:00:01,000 --> 00:00:02,000\nHello candidate\n")
	case "vtt":
		return []byte("WEBVTT\n\n00:01.000 --> 00:02.000\nHello candidate\n")
	default:
		root := fixtureRoot(t)
		fixture := mustFixture(t, verifiedManifest(t, root), "scripted/"+format+"-basic")
		return readArtifact(t, root, mustArtifact(t, fixture, "source"))
	}
}

func TestScriptedStableCandidateRenderConformance(t *testing.T) {
	t.Parallel()
	repository := filepath.Dir(fixtureRoot(t))
	canonical := readFile(t, filepath.Join(repository, "internal", "schema", "cueson.schema.json"))
	immutable := readFile(t, filepath.Join(repository, "schema", "releases", "v1.1.0", "cueson.schema.json"))
	emitted, diagnostics := runScriptedPlatformCLI(t, []string{"schema"}, cli.ExitSuccess)
	if schema.ID() != "https://cueson.io/schema/v1.1.0/cueson.schema.json" || schema.Version() != "1.1.0" || len(diagnostics) != 0 || !bytes.Equal(canonical, immutable) || !bytes.Equal(canonical, emitted) {
		t.Fatal("stable canonical/immutable/emitted contract differs")
	}
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			root := fixtureRoot(t)
			fixture := mustFixture(t, verifiedManifest(t, root), "scripted/"+format+"-basic")
			wantRender := readArtifact(t, root, mustArtifact(t, fixture, "expected_bytes"))
			directory := t.TempDir()
			input := filepath.Join(directory, "source."+format)
			if err := os.WriteFile(input, stableCandidateSource(t, format), 0o600); err != nil {
				t.Fatal(err)
			}
			encoded, _ := runScriptedPlatformCLI(t, []string{"encode", input, "--stdout"}, cli.ExitSuccess)
			document, err := schema.Decode(encoded)
			if err != nil {
				t.Fatal(err)
			}
			for _, status := range []string{"stable", "experimental", "schema_only"} {
				observed := document
				observed.FormatSupport = model.FormatSupport{Status: status, IngestSupported: status != "schema_only", RenderSupported: status != "schema_only", RestoreSupported: true}
				before, err := json.Marshal(observed)
				if err != nil {
					t.Fatal(err)
				}
				path := filepath.Join(directory, status+".cueson.json")
				if err = os.WriteFile(path, before, 0o600); err != nil {
					t.Fatal(err)
				}
				runScriptedPlatformCLI(t, []string{"validate", path}, cli.ExitSuccess)
				rendered, _ := runScriptedPlatformCLI(t, []string{"render", path, "--to", format, "--output", "-"}, cli.ExitSuccess)
				originalRender, _ := runScriptedPlatformCLI(t, []string{"render", path, "--to", format, "--output", "-"}, cli.ExitSuccess)
				if !bytes.Equal(rendered, wantRender) || !bytes.Equal(rendered, originalRender) || !bytes.Equal(readFile(t, path), before) {
					t.Fatal("stable rendering changed loaded capability observation or is nondeterministic")
				}
			}
		})
	}
}

// Hosted native tests execute the complete stable identity workflow on each governed platform.
func TestScriptedStableCandidateNativePlatformConformance(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"srt", "vtt", "ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			directory := t.TempDir()
			input := filepath.Join(directory, "source."+format)
			original := stableCandidateSource(t, format)
			if err := os.WriteFile(input, original, 0o600); err != nil {
				t.Fatal(err)
			}
			encoded, _ := runScriptedPlatformCLI(t, []string{"encode", input, "--stdout"}, cli.ExitSuccess)
			document, err := schema.Decode(encoded)
			if err != nil {
				t.Fatal(err)
			}
			if document.Schema != "https://cueson.io/schema/v1.1.0/cueson.schema.json" || document.SchemaVersion != "1.1.0" || document.Producer.Version != "1.1.0" || document.FormatSupport.Status != "stable" || !document.FormatSupport.IngestSupported || !document.FormatSupport.RenderSupported || !document.FormatSupport.RestoreSupported {
				t.Fatal("official candidate output has inconsistent identity or native support")
			}
			asset := document.Source.Assets[0]
			captured, err := base64.StdEncoding.DecodeString(asset.DataBase64)
			digest := sha256.Sum256(original)
			if err != nil || !bytes.Equal(captured, original) || asset.Size.Bytes != int64(len(original)) || asset.Hashes.SHA256 != hex.EncodeToString(digest[:]) {
				t.Fatal("stable encode changed source bytes or integrity")
			}
			path := filepath.Join(directory, "source.cueson.json")
			if err = os.WriteFile(path, encoded, 0o600); err != nil {
				t.Fatal(err)
			}
			runScriptedPlatformCLI(t, []string{"validate", path}, cli.ExitSuccess)
			inspected, _ := runScriptedPlatformCLI(t, []string{"inspect", path, "--json"}, cli.ExitSuccess)
			var report struct {
				Schema       struct{ Version string }               `json:"schema"`
				Capabilities struct{ Declared model.FormatSupport } `json:"capabilities"`
			}
			if err = json.Unmarshal(inspected, &report); err != nil || report.Schema.Version != "1.1.0" || report.Capabilities.Declared.Status != "stable" {
				t.Fatal("stable inspection did not report loaded identity/capabilities")
			}
			output := filepath.Join(directory, "restored."+format)
			runScriptedPlatformCLI(t, []string{"restore", path, "--output", output, "--no-metadata"}, cli.ExitSuccess)
			if !bytes.Equal(readFile(t, output), original) || !bytes.Equal(readFile(t, path), encoded) {
				t.Fatal("stable restoration changed source/model")
			}
			document.Schema = "https://cueson.io/schema/v1.1.0-dev/cueson.schema.json"
			document.SchemaVersion = "1.1.0-dev"
			bad, err := json.Marshal(document)
			if err != nil {
				t.Fatal(err)
			}
			badPath := filepath.Join(directory, "unsupported.cueson.json")
			if err = os.WriteFile(badPath, bad, 0o600); err != nil {
				t.Fatal(err)
			}
			stdout, _ := runScriptedPlatformCLI(t, []string{"restore", badPath, "--output", output, "--force", "--no-metadata"}, cli.ExitRuntimeFailure)
			if len(stdout) != 0 || !bytes.Equal(readFile(t, output), original) {
				t.Fatal("unsupported development identity replaced forced destination")
			}
		})
	}
}
