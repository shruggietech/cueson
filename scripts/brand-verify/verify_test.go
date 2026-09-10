package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestVerifyRepositoryHappyPath(t *testing.T) {
	fixture := newFixture(t, []zipFixtureEntry{{name: "logos/cueson.svg", content: "<svg/>"}})
	result, err := verifyRepository(fixture.root, defaultManifestPath)
	if err != nil {
		t.Fatalf("verifyRepository: %v", err)
	}
	if len(result.violations) != 0 {
		t.Fatalf("violations = %v", result.violations)
	}
	if result.entries != 1 || result.references != 1 {
		t.Fatalf("counts = %d entries, %d references", result.entries, result.references)
	}
}

func TestVerifyRepositoryRejectsUnknownAndDuplicateManifestFields(t *testing.T) {
	fixture := newFixture(t, []zipFixtureEntry{{name: "asset.svg", content: "x"}})
	manifestPath := filepath.Join(fixture.root, filepath.FromSlash(defaultManifestPath))
	original := mustRead(t, manifestPath)

	t.Run("unknown", func(t *testing.T) {
		changed := bytes.Replace(original, []byte(`"schema_version": 1`), []byte(`"unknown": true, "schema_version": 1`), 1)
		mustWrite(t, manifestPath, changed)
		result, err := verifyRepository(fixture.root, defaultManifestPath)
		if err != nil {
			t.Fatal(err)
		}
		assertViolationContains(t, result, "unknown field")
	})

	t.Run("duplicate", func(t *testing.T) {
		changed := bytes.Replace(original, []byte(`"schema_version": 1`), []byte(`"schema_version": 1, "schema_version": 1`), 1)
		mustWrite(t, manifestPath, changed)
		result, err := verifyRepository(fixture.root, defaultManifestPath)
		if err != nil {
			t.Fatal(err)
		}
		assertViolationContains(t, result, "duplicate field")
	})

	t.Run("wrong case", func(t *testing.T) {
		changed := bytes.Replace(original, []byte(`"schema_version": 1`), []byte(`"Schema_Version": 1`), 1)
		mustWrite(t, manifestPath, changed)
		result, err := verifyRepository(fixture.root, defaultManifestPath)
		if err != nil {
			t.Fatal(err)
		}
		assertViolationContains(t, result, "missing required field")
	})
}

func TestVerifyRepositoryRejectsMaliciousZIPPaths(t *testing.T) {
	tests := []struct {
		name string
		path string
		mode os.FileMode
		want string
	}{
		{name: "traversal", path: "../escape.svg", want: "parent components are not allowed"},
		{name: "absolute", path: "/absolute.svg", want: "absolute or drive-qualified"},
		{name: "backslash", path: `logos\evil.svg`, want: "backslashes are not allowed"},
		{name: "reserved", path: "logos/CON.svg", want: "Windows-reserved name"},
		{name: "symlink", path: "logos/link.svg", mode: os.ModeSymlink | 0o777, want: "only regular files"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t, []zipFixtureEntry{{name: test.path, content: "x", mode: test.mode}})
			result, err := verifyRepository(fixture.root, defaultManifestPath)
			if err != nil {
				t.Fatal(err)
			}
			assertViolationContains(t, result, test.want)
		})
	}
}

func TestVerifyRepositoryRejectsDuplicateAndCaseCollidingZIPEntries(t *testing.T) {
	tests := []struct {
		name    string
		entries []zipFixtureEntry
		want    string
	}{
		{name: "duplicate", entries: []zipFixtureEntry{{name: "a.svg", content: "one"}, {name: "a.svg", content: "two"}}, want: "duplicate path"},
		{name: "case collision", entries: []zipFixtureEntry{{name: "A.svg", content: "one"}, {name: "a.svg", content: "two"}}, want: "case-insensitive collision"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t, test.entries)
			result, err := verifyRepository(fixture.root, defaultManifestPath)
			if err != nil {
				t.Fatal(err)
			}
			assertViolationContains(t, result, test.want)
		})
	}
}

func TestVerifyRepositoryDetectsArchivePayloadAndReferenceDrift(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(t *testing.T, fixture testFixture)
		want   string
	}{
		{
			name: "archive",
			mutate: func(t *testing.T, fixture testFixture) {
				file, err := os.OpenFile(fixture.archivePath, os.O_APPEND|os.O_WRONLY, 0)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := file.Write([]byte("drift")); err != nil {
					t.Fatal(err)
				}
				if err := file.Close(); err != nil {
					t.Fatal(err)
				}
			},
			want: "archive: byte size",
		},
		{
			name: "payload content",
			mutate: func(t *testing.T, fixture testFixture) {
				mustWrite(t, fixture.payloadPath, []byte("changed"))
			},
			want: "payload entry",
		},
		{
			name: "additional payload",
			mutate: func(t *testing.T, fixture testFixture) {
				mustWrite(t, filepath.Join(filepath.Dir(fixture.payloadPath), "extra.svg"), []byte("extra"))
			},
			want: "payload: unrecorded entry",
		},
		{
			name: "missing reference",
			mutate: func(t *testing.T, fixture testFixture) {
				mustWrite(t, fixture.documentPath, []byte("reference removed"))
			},
			want: "exact reference text is absent",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t, []zipFixtureEntry{{name: "logos/cueson.svg", content: "<svg/>"}})
			test.mutate(t, fixture)
			result, err := verifyRepository(fixture.root, defaultManifestPath)
			if err != nil {
				t.Fatal(err)
			}
			assertViolationContains(t, result, test.want)
		})
	}
}

func TestVerifyRepositoryRejectsPayloadSymlink(t *testing.T) {
	fixture := newFixture(t, []zipFixtureEntry{{name: "logos/cueson.svg", content: "<svg/>"}})
	if err := os.Remove(fixture.payloadPath); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(fixture.root, "outside.svg")
	mustWrite(t, outside, []byte("<svg/>"))
	if err := os.Symlink(outside, fixture.payloadPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	result, err := verifyRepository(fixture.root, defaultManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	assertViolationContains(t, result, "symbolic links are not allowed")
}

func TestVerifyRepositorySortsDiagnostics(t *testing.T) {
	fixture := newFixture(t, []zipFixtureEntry{{name: "logos/cueson.svg", content: "<svg/>"}})
	mustWrite(t, fixture.payloadPath, []byte("changed"))
	mustWrite(t, fixture.documentPath, []byte("missing"))
	result, err := verifyRepository(fixture.root, defaultManifestPath)
	if err != nil {
		t.Fatal(err)
	}
	for index := 1; index < len(result.violations); index++ {
		if result.violations[index-1] > result.violations[index] {
			t.Fatalf("diagnostics not sorted: %v", result.violations)
		}
	}
}

func TestRunCLIExitConventions(t *testing.T) {
	fixture := newFixture(t, []zipFixtureEntry{{name: "logos/cueson.svg", content: "<svg/>"}})
	var stdout, stderr bytes.Buffer
	if code := runCLI([]string{"-repo", fixture.root}, &stdout, &stderr); code != 0 {
		t.Fatalf("success code = %d, stderr = %s", code, stderr.String())
	}
	mustWrite(t, fixture.payloadPath, []byte("drift"))
	stdout.Reset()
	stderr.Reset()
	if code := runCLI([]string{"-repo", fixture.root}, &stdout, &stderr); code != 1 {
		t.Fatalf("violation code = %d, stderr = %s", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := runCLI([]string{"extra"}, &stdout, &stderr); code != 2 {
		t.Fatalf("usage code = %d, stderr = %s", code, stderr.String())
	}
}

type zipFixtureEntry struct {
	name    string
	content string
	mode    os.FileMode
}

type testFixture struct {
	root         string
	archivePath  string
	payloadPath  string
	documentPath string
}

func newFixture(t *testing.T, archiveEntries []zipFixtureEntry) testFixture {
	t.Helper()
	root := t.TempDir()
	manifestRelative := filepath.FromSlash(defaultManifestPath)
	base := filepath.Dir(filepath.Join(root, manifestRelative))
	archivePath := filepath.Join(base, "archive", "cueson-brand-1.0.0.zip")
	payloadRoot := filepath.Join(base, "kit")
	if err := os.MkdirAll(filepath.Dir(archivePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(payloadRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	writer := zip.NewWriter(file)
	manifestEntries := make([]manifestEntry, 0, len(archiveEntries))
	var payloadPath string
	var total int64
	for _, entry := range archiveEntries {
		header := &zip.FileHeader{Name: entry.name, Method: zip.Deflate}
		if entry.mode != 0 {
			header.SetMode(entry.mode)
		} else {
			header.SetMode(0o644)
		}
		header.Modified = time.Unix(0, 0).UTC()
		part, err := writer.CreateHeader(header)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.WriteString(part, entry.content); err != nil {
			t.Fatal(err)
		}
		digest := sha256.Sum256([]byte(entry.content))
		manifestEntries = append(manifestEntries, manifestEntry{Path: entry.name, Bytes: int64(len(entry.content)), SHA256: hex.EncodeToString(digest[:])})
		total += int64(len(entry.content))
		if payloadPath == "" {
			payloadPath = filepath.Join(payloadRoot, filepath.FromSlash(entry.name))
		}
		if entry.mode&os.ModeSymlink == 0 && validatePortablePath(entry.name) == nil {
			target := filepath.Join(payloadRoot, filepath.FromSlash(entry.name))
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				t.Fatal(err)
			}
			mustWrite(t, target, []byte(entry.content))
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	archiveInfo, err := os.Stat(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	archiveDigest, err := digestFile(archivePath)
	if err != nil {
		t.Fatal(err)
	}

	documentRelative := "README.md"
	documentPath := filepath.Join(root, documentRelative)
	referenceText := "brand/cueson/1.0.0/kit/" + archiveEntries[0].name
	mustWrite(t, documentPath, []byte("![Cueson]("+referenceText+")\n"))
	manifest := importManifest{
		SchemaVersion: 1,
		Brand:         "cueson",
		KitVersion:    "1.0.0",
		SourceURL:     "https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip",
		EffectiveURL:  "https://brand.shruggie.tech/cueson/downloads/cueson-brand-1.0.0.zip",
		AcquiredAt:    "2026-09-10T12:00:00Z",
		Archive: archiveIdentity{
			Path:              "brand/cueson/1.0.0/archive/cueson-brand-1.0.0.zip",
			Bytes:             archiveInfo.Size(),
			SHA256:            archiveDigest,
			EntryCount:        len(manifestEntries),
			UncompressedBytes: total,
		},
		PayloadRoot: "brand/cueson/1.0.0/kit",
		Entries:     manifestEntries,
		References: []repositoryReference{{
			Document:  documentRelative,
			Asset:     archiveEntries[0].name,
			Reference: referenceText,
			Purpose:   "test identity",
		}},
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	manifestPath := filepath.Join(root, manifestRelative)
	if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWrite(t, manifestPath, data)
	return testFixture{root: root, archivePath: archivePath, payloadPath: payloadPath, documentPath: documentPath}
}

func mustWrite(t *testing.T, filePath string, data []byte) {
	t.Helper()
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, filePath string) []byte {
	t.Helper()
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertViolationContains(t *testing.T, result verificationResult, want string) {
	t.Helper()
	for _, violation := range result.violations {
		if strings.Contains(violation, want) {
			return
		}
	}
	t.Fatalf("violations %q do not contain %q", result.violations, want)
}
