package testutil

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestVerifyFixtures(t *testing.T) {
	t.Parallel()

	root := createCorpus(t, []byte("hello\n"))
	manifest, err := VerifyFixtures(root)
	if err != nil {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
	if len(manifest.Fixtures) != 1 || manifest.Fixtures[0].ID != "accepted/basic" {
		t.Fatalf("VerifyFixtures() returned unexpected manifest: %#v", manifest)
	}
}

func TestVerifyFixturesRejectsInvalidCorpus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*testing.T, string, *Manifest)
		want   string
	}{
		{name: "unknown manifest field", mutate: addUnknownManifestField, want: "decode manifest"},
		{name: "duplicate fixture id", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures = append(manifest.Fixtures, manifest.Fixtures[0])
		}, want: "fixture identifier is duplicated"},
		{name: "duplicate path", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			duplicate := manifest.Fixtures[0]
			duplicate.ID = "accepted/other"
			manifest.Fixtures = append(manifest.Fixtures, duplicate)
		}, want: "artifact path is duplicated"},
		{name: "portable path collision", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			duplicate := manifest.Fixtures[0]
			duplicate.ID = "accepted/other"
			duplicate.Artifacts[0].Path = "fixtures/basic/source/INPUT.TXT"
			manifest.Fixtures = append(manifest.Fixtures, duplicate)
		}, want: "artifact path is duplicated"},
		{name: "unsafe path", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Artifacts[0].Path = "../outside"
		}, want: "artifact source path is not portable"},
		{name: "declared directory", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Artifacts[0].Path = "fixtures/basic/source"
		}, want: "declared payload is not a regular file"},
		{name: "unknown redistribution", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Redistribution.Status = "unknown"
		}, want: "redistribution status must be approved"},
		{name: "missing notice decision", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Redistribution.NoticeRequired = nil
		}, want: "NOTICE decision is required"},
		{name: "digest mismatch", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Artifacts[0].SHA256 = strings.Repeat("0", 64)
		}, want: "SHA-256 does not match"},
		{name: "size mismatch", mutate: func(_ *testing.T, _ string, manifest *Manifest) { (*manifest.Fixtures[0].Artifacts[0].SizeBytes)++ }, want: "byte length does not match"},
		{name: "missing size", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Artifacts[0].SizeBytes = nil
		}, want: "size_bytes is required"},
		{name: "byte contract mismatch", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Artifacts[0].ByteContract.LineEndings = "crlf"
		}, want: "line endings do not match"},
		{name: "orphan payload", mutate: func(t *testing.T, root string, _ *Manifest) {
			writeFile(t, filepath.Join(root, "fixtures", "orphan.bin"), []byte{1})
		}, want: "payload is not declared"},
		{name: "missing payload", mutate: func(t *testing.T, root string, _ *Manifest) {
			if err := os.Remove(filepath.Join(root, filepath.FromSlash("fixtures/basic/source/input.txt"))); err != nil {
				t.Fatal(err)
			}
		}, want: "declared payload is missing"},
		{name: "local origin path", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Origin.Source = `C:\Users\person\fixture.txt`
		}, want: "origin source must not be a local path"},
		{name: "embedded origin path", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Origin.Source = "Captured from file:///tmp/private.txt"
		}, want: "origin source must not be a local path"},
		{name: "embedded Unix origin path", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Origin.Source = "Captured from /mnt/private/fixture.txt"
		}, want: "origin source must not be a local path"},
		{name: "local recipe path", mutate: func(_ *testing.T, _ string, manifest *Manifest) {
			manifest.Fixtures[0].Origin.Recipe = stringPointer(`Copy C:\Users\person\fixture.txt.`)
		}, want: "origin recipe must not contain a local path"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			root := createCorpus(t, []byte("hello\n"))
			manifest := readManifest(t, root)
			test.mutate(t, root, &manifest)
			if test.name != "unknown manifest field" {
				writeManifest(t, root, manifest)
			}
			_, err := VerifyFixtures(root)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("VerifyFixtures() error = %v, want fragment %q", err, test.want)
			}
			if strings.Contains(err.Error(), root) {
				t.Fatalf("error leaked fixture root: %v", err)
			}
		})
	}
}

func TestVerifyFixturesRejectsCanonicalAliasInInventory(t *testing.T) {
	t.Parallel()

	root := createCorpus(t, []byte("hello\n"))
	manifest := readManifest(t, root)
	if runtime.GOOS == "windows" {
		manifest.Fixtures[0].Artifacts[0].Path = "fixtures/basic/source/INPUT.TXT"
		writeManifest(t, root, manifest)
	} else {
		writeFile(t, filepath.Join(root, filepath.FromSlash("fixtures/basic/source/INPUT.TXT")), []byte("hello\n"))
	}
	_, err := VerifyFixtures(root)
	if err == nil || !strings.Contains(err.Error(), "collides with the declared portable path") {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
}

func TestVerifyFixturesRequiresEmptyArtifactSize(t *testing.T) {
	t.Parallel()

	root := createCorpus(t, nil)
	manifest := readManifest(t, root)
	manifest.Fixtures[0].Artifacts[0].SizeBytes = nil
	writeManifest(t, root, manifest)
	_, err := VerifyFixtures(root)
	if err == nil || !strings.Contains(err.Error(), "size_bytes is required") {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
}

func TestVerifyFixturesInventoriesFuzzRegressions(t *testing.T) {
	t.Parallel()

	root := createCorpus(t, []byte("hello\n"))
	manifest := readManifest(t, root)
	data := []byte{0xff, 0x00}
	path := "fuzz/schema/minimal/input/document.bin"
	writeFile(t, filepath.Join(root, filepath.FromSlash(path)), data)
	sum := sha256.Sum256(data)
	notice := false
	manifest.Fixtures = append(manifest.Fixtures, Fixture{
		ID:      "fuzz/schema-minimal",
		Purpose: "Retain a minimized fuzz regression.",
		Class:   "fuzz_regression",
		Origin:  Origin{Kind: "synthetic", Source: "Cueson test suite", Recipe: stringPointer("Write bytes ff 00.")},
		Redistribution: Redistribution{
			Status: "approved", License: "MIT", Attribution: "Copyright ShruggieTech", NoticeRequired: &notice,
		},
		Artifacts: []Artifact{{
			ID: "input", Role: "malformed_input", Path: path, MediaType: "application/octet-stream", SizeBytes: int64Pointer(int64(len(data))), SHA256: hex.EncodeToString(sum[:]),
			ByteContract: ByteContract{Encoding: "binary", BOM: "not_applicable", LineEndings: "not_applicable", FinalNewline: "not_applicable"},
		}},
		Expectation: Expectation{Result: "rejected", Stage: stringPointer("parse"), DiagnosticContains: stringPointer("parse Cue JSON")},
	})
	writeManifest(t, root, manifest)
	if _, err := VerifyFixtures(root); err != nil {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
}

func TestVerifyFixturesRejectsSymlink(t *testing.T) {
	t.Parallel()

	root := createCorpus(t, []byte("hello\n"))
	target := filepath.Join(root, filepath.FromSlash("fixtures/basic/source/input.txt"))
	link := filepath.Join(root, filepath.FromSlash("fixtures/basic/source/link.txt"))
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symbolic links unavailable: %v", err)
	}
	manifest := readManifest(t, root)
	artifact := manifest.Fixtures[0].Artifacts[0]
	artifact.ID = "link"
	artifact.Path = "fixtures/basic/source/link.txt"
	manifest.Fixtures[0].Artifacts = append(manifest.Fixtures[0].Artifacts, artifact)
	writeManifest(t, root, manifest)

	_, err := VerifyFixtures(root)
	if err == nil || !strings.Contains(err.Error(), "symbolic links are not allowed") {
		t.Fatalf("VerifyFixtures() error = %v", err)
	}
}

func createCorpus(t *testing.T, data []byte) string {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash("fixtures/basic/source/input.txt"))
	writeFile(t, path, data)
	sum := sha256.Sum256(data)
	notice := false
	manifest := Manifest{
		ManifestVersion: 1,
		Fixtures: []Fixture{{
			ID:             "accepted/basic",
			Purpose:        "Exercise a basic accepted text fixture.",
			Class:          "accepted",
			Origin:         Origin{Kind: "synthetic", Source: "Cueson S005 test suite", Recipe: stringPointer("Write the ASCII bytes hello followed by LF.")},
			Redistribution: Redistribution{Status: "approved", License: "MIT", Attribution: "Copyright ShruggieTech", NoticeRequired: &notice},
			Artifacts:      []Artifact{{ID: "source", Role: "source", Path: "fixtures/basic/source/input.txt", MediaType: "text/plain", SizeBytes: int64Pointer(int64(len(data))), SHA256: hex.EncodeToString(sum[:]), ByteContract: byteContractForTestData(data)}},
			Expectation:    Expectation{Result: "accepted"},
		}},
	}
	writeManifest(t, root, manifest)
	return root
}

func addUnknownManifestField(t *testing.T, root string, _ *Manifest) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var value map[string]any
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatal(err)
	}
	value["unexpected"] = true
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "manifest.json"), encoded)
}

func readManifest(t *testing.T, root string) Manifest {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatal(err)
	}
	return manifest
}

func writeManifest(t *testing.T, root string, manifest Manifest) {
	t.Helper()
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "manifest.json"), append(data, '\n'))
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func stringPointer(value string) *string {
	return &value
}

func int64Pointer(value int64) *int64 {
	return &value
}

func byteContractForTestData(data []byte) ByteContract {
	contract := ByteContract{Encoding: "utf-8", BOM: "absent", LineEndings: "none", FinalNewline: "absent"}
	if len(data) > 0 && data[len(data)-1] == '\n' {
		contract.LineEndings = "lf"
		contract.FinalNewline = "present"
	}
	return contract
}
