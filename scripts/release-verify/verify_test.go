package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"debug/buildinfo"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

const testCommit = "0123456789abcdef0123456789abcdef01234567"

func TestExpectedTargets(t *testing.T) {
	t.Parallel()
	targets := expectedTargets("1.0.0")
	if len(targets) != 6 {
		t.Fatalf("target count = %d, want 6", len(targets))
	}
	want := []string{
		"cueson_1.0.0_linux_amd64.tar.gz", "cueson_1.0.0_linux_arm64.tar.gz",
		"cueson_1.0.0_darwin_amd64.tar.gz", "cueson_1.0.0_darwin_arm64.tar.gz",
		"cueson_1.0.0_windows_amd64.zip", "cueson_1.0.0_windows_arm64.zip",
	}
	for index := range want {
		if targets[index].Archive != want[index] {
			t.Fatalf("target %d archive = %q, want %q", index, targets[index].Archive, want[index])
		}
		if targets[index].SBOM != want[index]+".sbom.json" {
			t.Fatalf("target %d SBOM = %q", index, targets[index].SBOM)
		}
	}
}

func TestSafeMemberName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want bool
	}{
		{"cueson", true}, {"cueson.exe", true}, {"LICENSE", true},
		{"", false}, {".", false}, {"..", false}, {"../cueson", false},
		{"dir/cueson", false}, {`dir\cueson`, false}, {`C:\cueson`, false},
		{"/cueson", false}, {"cue\x00son", false}, {"cue..son", false},
	}
	for _, test := range tests {
		if got := safeMemberName(test.name); got != test.want {
			t.Errorf("safeMemberName(%q) = %v, want %v", test.name, got, test.want)
		}
	}
}

func TestReadArchiveRejectsUnsafeAndNonRegularMembers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		data []byte
	}{
		{"unsafe.zip", makeZIP(t, map[string][]byte{"../cueson": []byte("bad")})},
		{"unsafe.tar.gz", makeTar(t, []tarEntry{{name: "cueson", typeflag: tar.TypeSymlink}})},
		{"duplicate.tar.gz", makeTar(t, []tarEntry{{name: "cueson"}, {name: "cueson"}})},
	}
	for _, test := range tests {
		if _, err := readArchive(test.name, test.data); err == nil {
			t.Errorf("readArchive(%s) accepted invalid archive", test.name)
		}
	}
}

func TestVerifyMembers(t *testing.T) {
	t.Parallel()
	target := expectedTargets("0.0.0")[0]
	schema := []byte(`{"schema_version":"0.0.0"}`)
	valid := map[string]archiveMember{
		"cueson":             {name: "cueson", mode: 0o755, data: append(append([]byte("binary"), schema...), []byte(binaryVersionMarkerPrefix+"0.0.0")...)},
		"cueson.schema.json": {name: "cueson.schema.json", mode: 0o644, data: schema},
		"LICENSE":            {name: "LICENSE", mode: 0o644, data: []byte("license")},
		"NOTICE":             {name: "NOTICE", mode: 0o644, data: []byte("notice")},
	}
	if _, err := verifyMembers(target, valid, schema, []byte("license"), []byte("notice"), "0.0.0"); err != nil {
		t.Fatalf("valid members: %v", err)
	}
	tests := []struct {
		name   string
		mutate func(map[string]archiveMember)
	}{
		{"missing", func(items map[string]archiveMember) { delete(items, "NOTICE") }},
		{"extra", func(items map[string]archiveMember) { items["README"] = archiveMember{} }},
		{"schema drift", func(items map[string]archiveMember) {
			items["cueson.schema.json"] = archiveMember{data: []byte("drift")}
		}},
		{"license drift", func(items map[string]archiveMember) {
			items["LICENSE"] = archiveMember{mode: 0o644, data: []byte("drift")}
		}},
		{"notice drift", func(items map[string]archiveMember) {
			items["NOTICE"] = archiveMember{mode: 0o644, data: []byte("drift")}
		}},
		{"license executable", func(items map[string]archiveMember) {
			items["LICENSE"] = archiveMember{mode: 0o755, data: []byte("license")}
		}},
		{"notice executable", func(items map[string]archiveMember) {
			items["NOTICE"] = archiveMember{mode: 0o755, data: []byte("notice")}
		}},
		{"license special mode", func(items map[string]archiveMember) {
			items["LICENSE"] = archiveMember{mode: os.ModeSetuid | 0o644, data: []byte("license")}
		}},
		{"not executable", func(items map[string]archiveMember) {
			items["cueson"] = archiveMember{mode: 0o644, data: valid["cueson"].data}
		}},
		{"schema not embedded", func(items map[string]archiveMember) {
			items["cueson"] = archiveMember{mode: 0o755, data: []byte("binary" + binaryVersionMarkerPrefix + "0.0.0")}
		}},
		{"version marker missing", func(items map[string]archiveMember) {
			items["cueson"] = archiveMember{mode: 0o755, data: append([]byte("binary"), schema...)}
		}},
		{"version marker wrong", func(items map[string]archiveMember) {
			items["cueson"] = archiveMember{mode: 0o755, data: append(append([]byte("binary"), schema...), []byte(binaryVersionMarkerPrefix+"0.0.1")...)}
		}},
	}
	for _, test := range tests {
		items := cloneMembers(valid)
		test.mutate(items)
		if _, err := verifyMembers(target, items, schema, []byte("license"), []byte("notice"), "0.0.0"); err == nil {
			t.Errorf("%s accepted", test.name)
		}
	}
}

func TestVerifyReleaseMetadata(t *testing.T) {
	t.Parallel()
	valid := releaseMetadata{ProjectName: "cueson", Version: "0.0.0", Commit: testCommit}
	if err := verifyReleaseMetadata(valid, "0.0.0", testCommit); err != nil {
		t.Fatal(err)
	}
	for _, invalid := range []releaseMetadata{
		{ProjectName: "other", Version: "0.0.0", Commit: testCommit},
		{ProjectName: "cueson", Version: "0.0.1", Commit: testCommit},
		{ProjectName: "cueson", Version: "0.0.0", Commit: strings.Repeat("f", 40)},
	} {
		if err := verifyReleaseMetadata(invalid, "0.0.0", testCommit); err == nil {
			t.Errorf("accepted metadata %#v", invalid)
		}
	}
}

func TestVerifyBuildSettings(t *testing.T) {
	t.Parallel()
	target := expectedTargets("0.0.0")[0]
	valid := &buildinfo.BuildInfo{Settings: []debug.BuildSetting{
		{Key: "GOOS", Value: "linux"}, {Key: "GOARCH", Value: "amd64"},
		{Key: "CGO_ENABLED", Value: "0"}, {Key: "vcs.revision", Value: testCommit},
		{Key: "vcs.modified", Value: "false"}, {Key: "-trimpath", Value: "true"},
	}}
	if err := verifyBuildSettings(valid, target, testCommit, nil); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"GOOS", "GOARCH", "CGO_ENABLED", "vcs.revision", "vcs.modified", "-trimpath"} {
		invalid := *valid
		invalid.Settings = append([]debug.BuildSetting(nil), valid.Settings...)
		for index := range invalid.Settings {
			if invalid.Settings[index].Key == key {
				invalid.Settings[index].Value = "wrong"
			}
		}
		if err := verifyBuildSettings(&invalid, target, testCommit, nil); err == nil {
			t.Errorf("accepted invalid %s setting", key)
		}
	}
	leaking := *valid
	leaking.Settings = append(append([]debug.BuildSetting(nil), valid.Settings...), debug.BuildSetting{Key: "local.path", Value: "/home/alice/cueson"})
	if err := verifyBuildSettings(&leaking, target, testCommit, nil); err == nil {
		t.Error("accepted local path in Go build information")
	}
}

func TestVerifyArtifactCatalog(t *testing.T) {
	t.Parallel()
	targets := expectedTargets("0.0.0")
	checksum := "cueson_0.0.0_checksums.txt"
	valid := make([]artifactRecord, 0, len(targets)*2+1)
	for _, target := range targets {
		valid = append(valid,
			artifactRecord{Name: target.Archive, Path: "dist/" + target.Archive, GOOS: target.GOOS, GOARCH: target.GOARCH, Type: "Archive"},
			artifactRecord{Name: target.SBOM, Path: "dist/" + target.SBOM, Type: "SBOM"},
		)
	}
	valid = append(valid, artifactRecord{Name: checksum, Path: "dist/" + checksum, Type: "Checksum"})
	if err := verifyArtifactCatalog(valid, targets, checksum); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		mutate func([]artifactRecord) []artifactRecord
	}{
		{"missing", func(items []artifactRecord) []artifactRecord { return items[1:] }},
		{"duplicate", func(items []artifactRecord) []artifactRecord { return append(items, items[0]) }},
		{"unknown archive", func(items []artifactRecord) []artifactRecord {
			return append(items, artifactRecord{Name: "cueson_0.0.0_plan9_amd64.tar.gz", Path: "dist/other.tar.gz", Type: "Archive"})
		}},
		{"wrong target", func(items []artifactRecord) []artifactRecord { items[0].GOOS = "darwin"; return items }},
		{"wrong type", func(items []artifactRecord) []artifactRecord { items[0].Type = "Binary"; return items }},
		{"absolute path", func(items []artifactRecord) []artifactRecord { items[0].Path = `C:\\src\\archive.zip`; return items }},
		{"unix absolute path", func(items []artifactRecord) []artifactRecord { items[0].Path = "/src/archive.zip"; return items }},
		{"UNC path", func(items []artifactRecord) []artifactRecord {
			items[0].Path = `\\server\share\archive.zip`
			return items
		}},
		{"traversing path", func(items []artifactRecord) []artifactRecord { items[0].Path = "../archive.zip"; return items }},
	}
	for _, test := range tests {
		items := append([]artifactRecord(nil), valid...)
		if err := verifyArtifactCatalog(test.mutate(items), targets, checksum); err == nil {
			t.Errorf("%s accepted", test.name)
		}
	}
}

func TestVerifyDistCatalogRejectsUnknownReleaseOutput(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, "unexpected.tar.gz"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := verifyDistCatalog(directory, expectedTargets("0.0.0"), "cueson_0.0.0_checksums.txt"); err == nil {
		t.Fatal("accepted unknown release artifact")
	}
}

func TestParseChecksums(t *testing.T) {
	t.Parallel()
	digest := strings.Repeat("a", 64)
	valid := []byte(digest + "  cueson_0.0.0_linux_amd64.tar.gz\n")
	checksums, err := parseChecksums(valid)
	if err != nil || checksums["cueson_0.0.0_linux_amd64.tar.gz"] != digest {
		t.Fatalf("valid checksum: %v %#v", err, checksums)
	}
	invalid := [][]byte{
		[]byte("malformed"),
		[]byte(strings.ToUpper(digest) + "  file.tar.gz\n"),
		[]byte(digest + "  ../file.tar.gz\n"),
		[]byte(digest + "  file.tar.gz\n" + digest + "  file.tar.gz\n"),
	}
	for _, data := range invalid {
		if _, err := parseChecksums(data); err == nil {
			t.Errorf("accepted invalid checksum %q", data)
		}
	}
}

func TestVerifyChecksumCatalogAndDigest(t *testing.T) {
	t.Parallel()
	targets := expectedTargets("0.0.0")
	checksums := make(map[string]string, len(targets))
	for _, target := range targets {
		checksums[target.Archive] = strings.Repeat("a", 64)
	}
	if err := verifyChecksumCatalog(checksums, targets); err != nil {
		t.Fatal(err)
	}
	missing := cloneChecksums(checksums)
	delete(missing, targets[0].Archive)
	if err := verifyChecksumCatalog(missing, targets); err == nil {
		t.Fatal("accepted missing checksum")
	}
	unknown := cloneChecksums(checksums)
	delete(unknown, targets[0].Archive)
	unknown["cueson_0.0.0_checksums.txt"] = strings.Repeat("b", 64)
	if err := verifyChecksumCatalog(unknown, targets); err == nil {
		t.Fatal("accepted self-referential or unknown checksum")
	}
	data := []byte("archive")
	digest := sha256.Sum256(data)
	valid := map[string]string{targets[0].Archive: fmt.Sprintf("%x", digest)}
	if _, err := verifyArchiveDigest(targets[0].Archive, data, valid); err != nil {
		t.Fatal(err)
	}
	valid[targets[0].Archive] = strings.Repeat("0", 64)
	if _, err := verifyArchiveDigest(targets[0].Archive, data, valid); err == nil {
		t.Fatal("accepted digest mismatch")
	}
}

func TestVerifySBOM(t *testing.T) {
	t.Parallel()
	target := expectedTargets("0.0.0")[0]
	document := map[string]any{
		"spdxVersion": "SPDX-2.3", "name": "cueson_linux_amd64",
		"documentNamespace": "https://example.test/spdx/1",
		"packages":          []map[string]string{{"name": "cueson_linux_amd64", "versionInfo": "0.0.0+" + testCommit}},
	}
	data, _ := json.Marshal(document)
	if err := verifySBOM(data, target, "0.0.0", testCommit); err != nil {
		t.Fatalf("valid SBOM: %v", err)
	}
	mutations := []func(map[string]any){
		func(value map[string]any) { value["spdxVersion"] = "SPDX-2.2" },
		func(value map[string]any) { value["name"] = "wrong" },
		func(value map[string]any) { value["packages"] = []any{} },
		func(value map[string]any) { value["published"] = true },
		func(value map[string]any) { value["attestation"] = "claimed" },
	}
	for index, mutate := range mutations {
		copyValue := deepCopy(t, document)
		mutate(copyValue)
		encoded, _ := json.Marshal(copyValue)
		if err := verifySBOM(encoded, target, "0.0.0", testCommit); err == nil {
			t.Errorf("mutation %d accepted", index)
		}
	}
}

func TestScanForbiddenRecognizesEscapedAndStructuralPaths(t *testing.T) {
	t.Parallel()
	forbidden := deriveForbidden(`C:\Users\alice\src\cueson`, []string{"alice-host"})
	invalid := [][]byte{
		[]byte(`C:\Users\alice\src\cueson`),
		[]byte(`{"path":"C:\\Users\\alice\\src\\cueson"}`),
		[]byte(`/home/alice/cueson`),
		[]byte(`{"path":"/opt/ci/cueson"}`),
		[]byte(`{"path":"D:\\build\\cueson"}`),
		[]byte(`{"path":"D:/build/cueson"}`),
		[]byte(`{"path":"\\\\server\\share\\cueson"}`),
		[]byte(`{"detail":"built from /workspace/cueson"}`),
		[]byte(`{"fileName":"/opt/ci/cueson"}`),
		[]byte(`{"sourceInfo":"acquired package info from go module information: /workspace/cueson"}`),
		[]byte(`alice-host`),
	}
	for _, data := range invalid {
		if err := scanForbidden("test", data, forbidden); err == nil {
			t.Errorf("accepted forbidden data %q", data)
		}
	}
	if err := scanForbidden("test", []byte(`{"path":"internal/schema/cueson.schema.json"}`), forbidden); err != nil {
		t.Fatalf("safe relative path rejected: %v", err)
	}
	valid := [][]byte{
		[]byte(`{"url":"https://cueson.io/schema/v0.0.0/cueson.schema.json"}`),
		[]byte(`{"purl":"pkg:golang/github.com/shruggietech/cueson@v0.0.0"}`),
	}
	for _, data := range valid {
		if err := scanForbidden("test", data, nil); err != nil {
			t.Errorf("safe metadata rejected: %v", err)
		}
	}
	virtualRoots := [][]byte{
		[]byte(`{"fileName":"\\cueson","sourceInfo":"acquired package info from go module information: \\cueson"}`),
		[]byte(`{"fileName":"/cueson","sourceInfo":"acquired package info from go module information: /cueson"}`),
		[]byte(`{"fileName":"/cueson.exe","sourceInfo":"acquired package info from the following paths: /cueson.exe"}`),
	}
	for _, data := range virtualRoots {
		if err := scanForbidden("test", data, nil, "cueson", "cueson.exe"); err != nil {
			t.Errorf("Syft virtual root rejected: %v", err)
		}
	}
}

func TestVerifySchemaIdentity(t *testing.T) {
	t.Parallel()
	valid := []byte(`{"$id":"https://cueson.io/schema/v0.0.0/cueson.schema.json","properties":{"$schema":{"const":"https://cueson.io/schema/v0.0.0/cueson.schema.json"},"schema_version":{"const":"0.0.0"}}}`)
	if err := verifySchemaIdentity(valid, "0.0.0"); err != nil {
		t.Fatal(err)
	}
	if err := verifySchemaIdentity(valid, "0.1.0"); err == nil {
		t.Fatal("accepted mismatched version")
	}
	wrongSuffix := []byte(`{"$id":"https://cueson.io/schema/v0.0.0/cueson.schema.json?draft=1","properties":{"schema_version":{"const":"0.0.0"}}}`)
	if err := verifySchemaIdentity(wrongSuffix, "0.0.0"); err == nil {
		t.Fatal("accepted schema identity with content after the release path")
	}
	wrongAuthority := []byte(`{"$id":"https://example.test/schema/v0.0.0/cueson.schema.json","properties":{"$schema":{"const":"https://example.test/schema/v0.0.0/cueson.schema.json"},"schema_version":{"const":"0.0.0"}}}`)
	if err := verifySchemaIdentity(wrongAuthority, "0.0.0"); err == nil {
		t.Fatal("accepted schema identity from the wrong authority")
	}
	wrongInstanceSchema := []byte(`{"$id":"https://cueson.io/schema/v0.0.0/cueson.schema.json","properties":{"$schema":{"const":"https://example.test/schema/v0.0.0/cueson.schema.json"},"schema_version":{"const":"0.0.0"}}}`)
	if err := verifySchemaIdentity(wrongInstanceSchema, "0.0.0"); err == nil {
		t.Fatal("accepted a mismatched root instance $schema const")
	}
	if err := verifySchemaIdentity(append([]byte{0xef, 0xbb, 0xbf}, valid...), "0.0.0"); err == nil {
		t.Fatal("accepted UTF-8 BOM")
	}
	if err := verifySchemaIdentity([]byte{'{', '"', 'x', '"', ':', '"', 0xff, '"', '}'}, "0.0.0"); err == nil {
		t.Fatal("accepted invalid UTF-8")
	}
}

func TestLoadRepositorySchemasRequiresExactVersionedCopy(t *testing.T) {
	t.Parallel()
	valid := []byte(`{"$id":"https://cueson.io/schema/v0.0.0/cueson.schema.json","properties":{"$schema":{"const":"https://cueson.io/schema/v0.0.0/cueson.schema.json"},"schema_version":{"const":"0.0.0"}}}` + "\n")
	wantDigest := sha256.Sum256(valid)

	t.Run("valid", func(t *testing.T) {
		repository := makeSchemaRepository(t, valid, valid)
		canonical, digest, err := loadRepositorySchemas(repository, "0.0.0")
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(canonical, valid) {
			t.Fatal("returned canonical schema differs")
		}
		if digest != fmt.Sprintf("%x", wantDigest) {
			t.Fatalf("release schema digest = %q, want %x", digest, wantDigest)
		}
	})

	t.Run("missing release schema", func(t *testing.T) {
		repository := makeSchemaRepository(t, valid, nil)
		if _, _, err := loadRepositorySchemas(repository, "0.0.0"); err == nil || !strings.Contains(err.Error(), "versioned release schema") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("byte drift", func(t *testing.T) {
		repository := makeSchemaRepository(t, valid, append(append([]byte(nil), valid...), '\n'))
		if _, _, err := loadRepositorySchemas(repository, "0.0.0"); err == nil || !strings.Contains(err.Error(), "differs byte-for-byte") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("canonical version mismatch", func(t *testing.T) {
		repository := makeSchemaRepository(t, valid, valid)
		if _, _, err := loadRepositorySchemas(repository, "0.0.1"); err == nil || !strings.Contains(err.Error(), "canonical schema") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("release version mismatch", func(t *testing.T) {
		wrong := []byte(`{"$id":"https://cueson.io/schema/v0.0.1/cueson.schema.json","properties":{"$schema":{"const":"https://cueson.io/schema/v0.0.1/cueson.schema.json"},"schema_version":{"const":"0.0.1"}}}` + "\n")
		repository := makeSchemaRepository(t, valid, wrong)
		if _, _, err := loadRepositorySchemas(repository, "0.0.0"); err == nil || !strings.Contains(err.Error(), "versioned release schema") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("unsafe version", func(t *testing.T) {
		repository := makeSchemaRepository(t, valid, valid)
		if _, _, err := loadRepositorySchemas(repository, "../../escape"); err == nil || !strings.Contains(err.Error(), "safe semantic version") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestLoadDevelopmentSchemaUsesCanonicalWithoutReleaseCopy(t *testing.T) {
	t.Parallel()
	valid := []byte(`{"$id":"https://cueson.io/schema/v0.1.0/cueson.schema.json","properties":{"$schema":{"const":"https://cueson.io/schema/v0.1.0/cueson.schema.json"},"schema_version":{"const":"0.1.0"}}}` + "\n")
	wantDigest := sha256.Sum256(valid)
	repository := makeSchemaRepository(t, valid, nil)

	canonical, digest, err := loadDevelopmentSchema(repository, "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(canonical, valid) {
		t.Fatal("returned development schema differs")
	}
	if digest != fmt.Sprintf("%x", wantDigest) {
		t.Fatalf("development schema digest = %q, want %x", digest, wantDigest)
	}
	if _, _, err := loadDevelopmentSchema(repository, "0.0.0"); err == nil || !strings.Contains(err.Error(), "canonical schema") {
		t.Fatalf("version mismatch error = %v", err)
	}
}

func TestVerifyProbeOutput(t *testing.T) {
	t.Parallel()
	want := []byte("0.0.0\n")
	if err := verifyProbeOutput([]string{"version"}, want, nil, nil, want); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		stdout []byte
		stderr []byte
		runErr error
	}{
		{"nonzero", nil, []byte("failed"), errors.New("exit status 1")},
		{"stderr", want, []byte("warning"), nil},
		{"extra stdout", []byte("0.0.0\nextra\n"), nil, nil},
		{"wrong version", []byte("0.0.1\n"), nil, nil},
		{"schema drift", []byte("{}\n"), nil, nil},
	}
	for _, test := range tests {
		if err := verifyProbeOutput([]string{"probe"}, test.stdout, test.stderr, test.runErr, want); err == nil {
			t.Errorf("%s accepted", test.name)
		}
	}
}

func TestWriteEvidenceIsExclusiveAndDeterministic(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	firstPath := filepath.Join(directory, "first.json")
	secondPath := filepath.Join(directory, "second.json")
	evidence := ReleaseEvidence{
		Version: "1.0.0", IntendedTag: "v1.0.0", SourceRevision: testCommit,
		ReleaseSchemaSHA256: strings.Repeat("a", 64), LicenseSHA256: strings.Repeat("b", 64), NoticeSHA256: strings.Repeat("c", 64),
		ArchiveCount: 6, SBOMCount: 6, ChecksumCount: 6,
	}
	if err := writeEvidence(firstPath, evidence); err != nil {
		t.Fatal(err)
	}
	if err := writeEvidence(secondPath, evidence); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("identical evidence inputs produced different bytes")
	}
	if err := writeEvidence(firstPath, evidence); err == nil {
		t.Fatal("second write replaced evidence")
	}
	if !bytes.Contains(first, []byte(`"release_schema_sha256": "`+strings.Repeat("a", 64)+`"`)) {
		t.Fatal("evidence omits release schema digest")
	}
	for _, field := range []string{`"intended_tag": "v1.0.0"`, `"license_sha256": "` + strings.Repeat("b", 64) + `"`, `"notice_sha256": "` + strings.Repeat("c", 64) + `"`} {
		if !bytes.Contains(first, []byte(field)) {
			t.Fatalf("evidence omits %s", field)
		}
	}
	if !bytes.Contains(first, []byte(`"published": false`)) {
		t.Fatal("evidence omits non-publication state")
	}
}

func TestEvidenceOutputPath(t *testing.T) {
	t.Parallel()
	dist := t.TempDir()
	if got, want := evidenceOutputPath(Config{DistDir: dist}), filepath.Join(dist, evidenceFilename); got != want {
		t.Fatalf("default evidence path = %q, want %q", got, want)
	}
	custom := filepath.Join(t.TempDir(), "native-smoke.json")
	if got := evidenceOutputPath(Config{DistDir: dist, EvidencePath: custom}); got != custom {
		t.Fatalf("custom evidence path = %q, want %q", got, custom)
	}
}

func TestLoadRepositoryFile(t *testing.T) {
	t.Parallel()
	repository := t.TempDir()
	want := []byte("governed material\n")
	if err := os.WriteFile(filepath.Join(repository, "LICENSE"), want, 0o644); err != nil {
		t.Fatal(err)
	}
	file, err := loadRepositoryFile(repository, "LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(file.data, want) {
		t.Fatal("loaded repository material differs")
	}
	digest := sha256.Sum256(want)
	if file.sha256 != fmt.Sprintf("%x", digest) {
		t.Fatalf("digest = %q, want %x", file.sha256, digest)
	}
}

func TestDevelopmentEvidenceIsExplicit(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "development.json")
	evidence := ReleaseEvidence{Version: "0.1.0", Development: true, SourceRevision: testCommit, ReleaseSchemaSHA256: strings.Repeat("a", 64)}
	if err := writeEvidence(path, evidence); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"development": true`)) {
		t.Fatal("development evidence does not identify its non-release proof mode")
	}
}

func makeSchemaRepository(t *testing.T, canonical, release []byte) string {
	t.Helper()
	repository := t.TempDir()
	canonicalPath := filepath.Join(repository, "internal", "schema", releaseSchemaFilename)
	if err := os.MkdirAll(filepath.Dir(canonicalPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(canonicalPath, canonical, 0o644); err != nil {
		t.Fatal(err)
	}
	if release != nil {
		releasePath := filepath.Join(repository, "schema", "releases", "v0.0.0", releaseSchemaFilename)
		if err := os.MkdirAll(filepath.Dir(releasePath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(releasePath, release, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return repository
}

type tarEntry struct {
	name     string
	typeflag byte
}

func makeTar(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	var output bytes.Buffer
	gzipWriter := gzip.NewWriter(&output)
	writer := tar.NewWriter(gzipWriter)
	for _, entry := range entries {
		typeflag := entry.typeflag
		if typeflag == 0 {
			typeflag = tar.TypeReg
		}
		if err := writer.WriteHeader(&tar.Header{Name: entry.name, Typeflag: typeflag, Mode: 0o755, Size: 0}); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func makeZIP(t *testing.T, entries map[string][]byte) []byte {
	t.Helper()
	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for name, data := range entries {
		item, err := writer.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := item.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func cloneMembers(source map[string]archiveMember) map[string]archiveMember {
	result := make(map[string]archiveMember, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func cloneChecksums(source map[string]string) map[string]string {
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func deepCopy(t *testing.T, source map[string]any) map[string]any {
	t.Helper()
	data, err := json.Marshal(source)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}
