package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"debug/buildinfo"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

const (
	evidenceFilename          = "release-evidence.json"
	binaryVersionMarkerPrefix = "cueson-release-version:"
)

var (
	commitPattern           = regexp.MustCompile(`^[0-9a-f]{40}$`)
	windowsDrivePathPattern = regexp.MustCompile(`^[A-Za-z]:/`)
)

type Config struct {
	DistDir     string
	RepoDir     string
	Version     string
	Commit      string
	ExecuteHost bool
	Forbidden   []string
}

type Target struct {
	GOOS       string `json:"goos"`
	GOARCH     string `json:"goarch"`
	Archive    string `json:"archive"`
	Binary     string `json:"binary"`
	SBOM       string `json:"sbom"`
	ArchiveSHA string `json:"archive_sha256"`
	SBOMSHA    string `json:"sbom_sha256"`
}

type ReleaseEvidence struct {
	Version        string   `json:"version"`
	SourceRevision string   `json:"source_revision"`
	ArchiveCount   int      `json:"archive_count"`
	SBOMCount      int      `json:"sbom_count"`
	ChecksumCount  int      `json:"checksum_count"`
	HostExecuted   *string  `json:"host_executed"`
	Targets        []Target `json:"targets"`
	Published      bool     `json:"published"`
}

type archiveMember struct {
	name string
	mode os.FileMode
	data []byte
}

type artifactRecord struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	GOOS   string `json:"goos"`
	GOARCH string `json:"goarch"`
	Type   string `json:"type"`
}

type releaseMetadata struct {
	ProjectName string `json:"project_name"`
	Version     string `json:"version"`
	Commit      string `json:"commit"`
}

type spdxDocument struct {
	SPDXVersion       string `json:"spdxVersion"`
	Name              string `json:"name"`
	DocumentNamespace string `json:"documentNamespace"`
	Packages          []struct {
		Name        string `json:"name"`
		VersionInfo string `json:"versionInfo"`
	} `json:"packages"`
}

func expectedTargets(version string) []Target {
	definitions := []struct {
		goos string
		arch string
		ext  string
		bin  string
	}{
		{"linux", "amd64", ".tar.gz", "cueson"},
		{"linux", "arm64", ".tar.gz", "cueson"},
		{"darwin", "amd64", ".tar.gz", "cueson"},
		{"darwin", "arm64", ".tar.gz", "cueson"},
		{"windows", "amd64", ".zip", "cueson.exe"},
		{"windows", "arm64", ".zip", "cueson.exe"},
	}
	targets := make([]Target, 0, len(definitions))
	for _, definition := range definitions {
		archive := fmt.Sprintf("cueson_%s_%s_%s%s", version, definition.goos, definition.arch, definition.ext)
		targets = append(targets, Target{
			GOOS: definition.goos, GOARCH: definition.arch,
			Archive: archive, Binary: definition.bin, SBOM: archive + ".sbom.json",
		})
	}
	return targets
}

func Verify(ctx context.Context, config Config) (ReleaseEvidence, error) {
	var evidence ReleaseEvidence
	if !commitPattern.MatchString(config.Commit) {
		return evidence, fmt.Errorf("expected commit must be a full lowercase 40-character SHA-1")
	}
	distDir, err := filepath.Abs(config.DistDir)
	if err != nil {
		return evidence, fmt.Errorf("resolve dist: %w", err)
	}
	repoDir, err := filepath.Abs(config.RepoDir)
	if err != nil {
		return evidence, fmt.Errorf("resolve repository: %w", err)
	}
	canonicalSchema, err := os.ReadFile(filepath.Join(repoDir, "internal", "schema", "cueson.schema.json"))
	if err != nil {
		return evidence, fmt.Errorf("read canonical schema: %w", err)
	}
	if err := verifySchemaIdentity(canonicalSchema, config.Version); err != nil {
		return evidence, fmt.Errorf("canonical schema: %w", err)
	}

	forbidden := deriveForbidden(repoDir, config.Forbidden)
	metadataRaw, err := readAndScan(filepath.Join(distDir, "metadata.json"), forbidden)
	if err != nil {
		return evidence, err
	}
	var metadata releaseMetadata
	if err := json.Unmarshal(metadataRaw, &metadata); err != nil {
		return evidence, fmt.Errorf("parse metadata.json: %w", err)
	}
	if err := verifyReleaseMetadata(metadata, config.Version, config.Commit); err != nil {
		return evidence, err
	}
	artifactsRaw, err := readAndScan(filepath.Join(distDir, "artifacts.json"), forbidden)
	if err != nil {
		return evidence, err
	}
	var artifacts []artifactRecord
	if err := json.Unmarshal(artifactsRaw, &artifacts); err != nil {
		return evidence, fmt.Errorf("parse artifacts.json: %w", err)
	}

	targets := expectedTargets(config.Version)
	checksumName := fmt.Sprintf("cueson_%s_checksums.txt", config.Version)
	if err := verifyArtifactCatalog(artifacts, targets, checksumName); err != nil {
		return evidence, err
	}
	if err := verifyDistCatalog(distDir, targets, checksumName); err != nil {
		return evidence, err
	}

	checksumRaw, err := readAndScan(filepath.Join(distDir, checksumName), forbidden)
	if err != nil {
		return evidence, err
	}
	checksums, err := parseChecksums(checksumRaw)
	if err != nil {
		return evidence, err
	}
	if err := verifyChecksumCatalog(checksums, targets); err != nil {
		return evidence, err
	}

	var hostExecuted *string
	for index := range targets {
		target := &targets[index]
		archivePath := filepath.Join(distDir, target.Archive)
		archiveBytes, err := os.ReadFile(archivePath)
		if err != nil {
			return evidence, fmt.Errorf("read %s: %w", target.Archive, err)
		}
		if err := scanForbidden(target.Archive, archiveBytes, forbidden); err != nil {
			return evidence, err
		}
		target.ArchiveSHA, err = verifyArchiveDigest(target.Archive, archiveBytes, checksums)
		if err != nil {
			return evidence, err
		}
		members, err := readArchive(target.Archive, archiveBytes)
		if err != nil {
			return evidence, fmt.Errorf("inspect %s: %w", target.Archive, err)
		}
		binaryBytes, err := verifyMembers(*target, members, canonicalSchema, config.Version)
		if err != nil {
			return evidence, fmt.Errorf("inspect %s: %w", target.Archive, err)
		}
		if err := scanForbidden(target.Archive+" binary", binaryBytes, forbidden); err != nil {
			return evidence, err
		}
		if err := verifyBuildInfo(binaryBytes, *target, config.Commit); err != nil {
			return evidence, fmt.Errorf("inspect %s build information: %w", target.Archive, err)
		}

		sbomRaw, err := readAndScan(filepath.Join(distDir, target.SBOM), forbidden)
		if err != nil {
			return evidence, err
		}
		if err := verifySBOM(sbomRaw, *target, config.Version, config.Commit); err != nil {
			return evidence, fmt.Errorf("inspect %s: %w", target.SBOM, err)
		}
		sbomDigest := sha256.Sum256(sbomRaw)
		target.SBOMSHA = hex.EncodeToString(sbomDigest[:])

		if config.ExecuteHost && target.GOOS == runtime.GOOS && target.GOARCH == runtime.GOARCH {
			if hostExecuted != nil {
				return evidence, fmt.Errorf("multiple host-compatible targets found")
			}
			if err := verifyHostCLI(ctx, binaryBytes, target.Binary, canonicalSchema, config.Version); err != nil {
				return evidence, fmt.Errorf("execute %s: %w", target.Archive, err)
			}
			identity := target.GOOS + "/" + target.GOARCH
			hostExecuted = &identity
		}
	}
	if config.ExecuteHost && hostExecuted == nil {
		return evidence, fmt.Errorf("no archive is executable on host %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	return ReleaseEvidence{
		Version: config.Version, SourceRevision: config.Commit,
		ArchiveCount: len(targets), SBOMCount: len(targets), ChecksumCount: len(checksums),
		HostExecuted: hostExecuted, Targets: targets, Published: false,
	}, nil
}

func readArchive(name string, data []byte) (map[string]archiveMember, error) {
	if strings.HasSuffix(name, ".zip") {
		reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			return nil, err
		}
		members := make(map[string]archiveMember, len(reader.File))
		for _, file := range reader.File {
			if !safeMemberName(file.Name) || file.FileInfo().IsDir() || file.Mode()&os.ModeType != 0 {
				return nil, fmt.Errorf("unsafe or non-regular ZIP member %q", file.Name)
			}
			if _, exists := members[file.Name]; exists {
				return nil, fmt.Errorf("duplicate member %q", file.Name)
			}
			opened, err := file.Open()
			if err != nil {
				return nil, err
			}
			contents, readErr := io.ReadAll(opened)
			closeErr := opened.Close()
			if readErr != nil {
				return nil, readErr
			}
			if closeErr != nil {
				return nil, closeErr
			}
			members[file.Name] = archiveMember{name: file.Name, mode: file.Mode(), data: contents}
		}
		return members, nil
	}
	if !strings.HasSuffix(name, ".tar.gz") {
		return nil, fmt.Errorf("unsupported archive format")
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	reader := tar.NewReader(gzipReader)
	members := make(map[string]archiveMember)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if !safeMemberName(header.Name) || (header.Typeflag != tar.TypeReg && header.Typeflag != tar.TypeRegA) {
			return nil, fmt.Errorf("unsafe or non-regular tar member %q", header.Name)
		}
		if _, exists := members[header.Name]; exists {
			return nil, fmt.Errorf("duplicate member %q", header.Name)
		}
		contents, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}
		members[header.Name] = archiveMember{name: header.Name, mode: os.FileMode(header.Mode), data: contents}
	}
	return members, nil
}

func verifyMembers(target Target, members map[string]archiveMember, canonicalSchema []byte, version string) ([]byte, error) {
	expected := map[string]bool{target.Binary: true, "cueson.schema.json": true, "LICENSE": true, "NOTICE": true}
	if len(members) != len(expected) {
		return nil, fmt.Errorf("archive has %d members, want %d", len(members), len(expected))
	}
	for name := range members {
		if !expected[name] {
			return nil, fmt.Errorf("unexpected member %q", name)
		}
	}
	for name := range expected {
		if _, ok := members[name]; !ok {
			return nil, fmt.Errorf("missing member %q", name)
		}
	}
	if !bytes.Equal(members["cueson.schema.json"].data, canonicalSchema) {
		return nil, fmt.Errorf("packaged schema differs from canonical schema")
	}
	if target.GOOS != "windows" && members[target.Binary].mode&0o111 == 0 {
		return nil, fmt.Errorf("binary member is not executable")
	}
	if !bytes.Contains(members[target.Binary].data, canonicalSchema) {
		return nil, fmt.Errorf("binary does not embed the canonical schema bytes")
	}
	if !bytes.Contains(members[target.Binary].data, []byte(binaryVersionMarkerPrefix+version)) {
		return nil, fmt.Errorf("binary does not contain the expected release-version marker")
	}
	return members[target.Binary].data, nil
}

func verifyBuildInfo(binary []byte, target Target, commit string) error {
	info, err := buildinfo.Read(bytes.NewReader(binary))
	if err != nil {
		return err
	}
	return verifyBuildSettings(info, target, commit)
}

func verifyBuildSettings(info *buildinfo.BuildInfo, target Target, commit string) error {
	settings := make(map[string]string, len(info.Settings))
	for _, setting := range info.Settings {
		settings[setting.Key] = setting.Value
	}
	required := map[string]string{
		"GOOS": target.GOOS, "GOARCH": target.GOARCH, "CGO_ENABLED": "0",
		"vcs.revision": commit, "vcs.modified": "false", "-trimpath": "true",
	}
	for key, want := range required {
		if got := settings[key]; got != want {
			return fmt.Errorf("setting %s=%q, want %q", key, got, want)
		}
	}
	return nil
}

func verifyReleaseMetadata(metadata releaseMetadata, version, commit string) error {
	if metadata.ProjectName != "cueson" || metadata.Version != version || metadata.Commit != commit {
		return fmt.Errorf("metadata identity is project=%q version=%q commit=%q", metadata.ProjectName, metadata.Version, metadata.Commit)
	}
	return nil
}

func verifyArtifactCatalog(artifacts []artifactRecord, targets []Target, checksumName string) error {
	type expectation struct {
		artifactType string
		goos         string
		goarch       string
	}
	expected := make(map[string]expectation, len(targets)*2+1)
	for _, target := range targets {
		expected[target.Archive] = expectation{artifactType: "Archive", goos: target.GOOS, goarch: target.GOARCH}
		expected[target.SBOM] = expectation{artifactType: "SBOM"}
	}
	expected[checksumName] = expectation{artifactType: "Checksum"}
	seen := make(map[string]int, len(expected))
	for _, artifact := range artifacts {
		name := artifact.Name
		if name == "" {
			name = filepath.Base(filepath.FromSlash(artifact.Path))
		}
		if artifact.Path != "" && (filepath.IsAbs(artifact.Path) || !safeRelativeArtifactPath(artifact.Path)) {
			return fmt.Errorf("artifacts.json contains unsafe path %q", artifact.Path)
		}
		want, required := expected[name]
		if !required {
			if artifact.Type == "Archive" || artifact.Type == "SBOM" || artifact.Type == "Checksum" {
				return fmt.Errorf("artifacts.json identifies unexpected %s %s", artifact.Type, name)
			}
			continue
		}
		if artifact.Type != want.artifactType {
			return fmt.Errorf("artifacts.json identifies %s as %s, want %s", name, artifact.Type, want.artifactType)
		}
		if want.artifactType == "Archive" && (artifact.GOOS != want.goos || artifact.GOARCH != want.goarch) {
			return fmt.Errorf("artifacts.json identifies %s as %s/%s, want %s/%s", name, artifact.GOOS, artifact.GOARCH, want.goos, want.goarch)
		}
		seen[name]++
	}
	for name := range expected {
		if seen[name] != 1 {
			return fmt.Errorf("artifacts.json identifies %s %d times, want 1", name, seen[name])
		}
	}
	return nil
}

func verifyDistCatalog(distDir string, targets []Target, checksumName string) error {
	expected := map[string]bool{checksumName: true, "metadata.json": true, "artifacts.json": true, "config.yaml": true}
	for _, target := range targets {
		expected[target.Archive] = true
		expected[target.SBOM] = true
	}
	entries, err := os.ReadDir(distDir)
	if err != nil {
		return fmt.Errorf("read dist catalog: %w", err)
	}
	for _, entry := range entries {
		if entry.IsDir() || expected[entry.Name()] || entry.Name() == evidenceFilename {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".zip") || strings.HasSuffix(entry.Name(), ".tar.gz") || strings.HasSuffix(entry.Name(), ".sbom.json") || strings.Contains(entry.Name(), "checksum") {
			return fmt.Errorf("dist contains unexpected release artifact %s", entry.Name())
		}
	}
	return nil
}

func parseChecksums(data []byte) (map[string]string, error) {
	checksums := make(map[string]string)
	for lineNumber, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) != 2 {
			return nil, fmt.Errorf("checksum line %d is malformed", lineNumber+1)
		}
		digest, name := fields[0], strings.TrimPrefix(fields[1], "*")
		if len(digest) != sha256.Size*2 || strings.ToLower(digest) != digest {
			return nil, fmt.Errorf("checksum line %d is not lowercase SHA-256", lineNumber+1)
		}
		if _, err := hex.DecodeString(digest); err != nil {
			return nil, fmt.Errorf("checksum line %d is not hexadecimal", lineNumber+1)
		}
		if !safeMemberName(name) {
			return nil, fmt.Errorf("checksum line %d has unsafe name %q", lineNumber+1, name)
		}
		if _, exists := checksums[name]; exists {
			return nil, fmt.Errorf("checksum manifest duplicates %s", name)
		}
		checksums[name] = digest
	}
	return checksums, nil
}

func verifyChecksumCatalog(checksums map[string]string, targets []Target) error {
	if len(checksums) != len(targets) {
		return fmt.Errorf("checksum manifest has %d entries, want %d", len(checksums), len(targets))
	}
	expected := make(map[string]bool, len(targets))
	for _, target := range targets {
		expected[target.Archive] = true
		if _, ok := checksums[target.Archive]; !ok {
			return fmt.Errorf("checksum manifest is missing %s", target.Archive)
		}
	}
	for name := range checksums {
		if !expected[name] {
			return fmt.Errorf("checksum manifest contains unexpected entry %s", name)
		}
	}
	return nil
}

func verifyArchiveDigest(name string, data []byte, checksums map[string]string) (string, error) {
	digest := sha256.Sum256(data)
	actual := hex.EncodeToString(digest[:])
	want, ok := checksums[name]
	if !ok {
		return "", fmt.Errorf("checksum manifest is missing %s", name)
	}
	if actual != want {
		return "", fmt.Errorf("checksum mismatch for %s", name)
	}
	return actual, nil
}

func verifySBOM(data []byte, target Target, version, commit string) error {
	var document spdxDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return fmt.Errorf("parse SPDX JSON: %w", err)
	}
	if document.SPDXVersion != "SPDX-2.3" || document.DocumentNamespace == "" {
		return fmt.Errorf("missing SPDX-2.3 identity")
	}
	wantName := "cueson_" + target.GOOS + "_" + target.GOARCH
	if document.Name != wantName {
		return fmt.Errorf("document name %q, want %q", document.Name, wantName)
	}
	wantVersion := version + "+" + commit
	foundRoot := false
	for _, item := range document.Packages {
		if item.Name == wantName && item.VersionInfo == wantVersion {
			foundRoot = true
		}
	}
	if !foundRoot {
		return fmt.Errorf("missing source package %s at %s", wantName, wantVersion)
	}
	var generic map[string]any
	if err := json.Unmarshal(data, &generic); err != nil {
		return err
	}
	if key, ok := prohibitedClaim(generic); ok {
		return fmt.Errorf("contains prohibited publication claim %q", key)
	}
	return nil
}

func prohibitedClaim(value any) (string, bool) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			normalized := strings.ToLower(strings.ReplaceAll(key, "_", ""))
			if normalized == "published" || normalized == "attestation" || normalized == "signature" || normalized == "releaseurl" {
				switch claim := child.(type) {
				case bool:
					if claim {
						return key, true
					}
				case string:
					if claim != "" {
						return key, true
					}
				default:
					if claim != nil {
						return key, true
					}
				}
			}
			if key, ok := prohibitedClaim(child); ok {
				return key, true
			}
		}
	case []any:
		for _, child := range typed {
			if key, ok := prohibitedClaim(child); ok {
				return key, true
			}
		}
	}
	return "", false
}

func verifySchemaIdentity(data []byte, version string) error {
	var schema struct {
		ID         string `json:"$id"`
		Properties map[string]struct {
			Const string `json:"const"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		return err
	}
	if !strings.Contains(schema.ID, "/v"+version+"/cueson.schema.json") {
		return fmt.Errorf("schema id %q does not identify v%s", schema.ID, version)
	}
	if got := schema.Properties["schema_version"].Const; got != version {
		return fmt.Errorf("schema_version const %q, want %q", got, version)
	}
	return nil
}

func verifyHostCLI(ctx context.Context, binary []byte, name string, schema []byte, version string) error {
	directory, err := os.MkdirTemp("", "cueson-release-verify-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(directory)
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, binary, 0o700); err != nil {
		return err
	}
	probes := []struct {
		args []string
		want []byte
	}{
		{[]string{"version"}, []byte(version + "\n")},
		{[]string{"schema", "--version"}, []byte(version + "\n")},
		{[]string{"schema"}, schema},
	}
	for _, probe := range probes {
		command := exec.CommandContext(ctx, path, probe.args...)
		configureProcess(command)
		var stdout, stderr bytes.Buffer
		command.Stdout = &stdout
		command.Stderr = &stderr
		runErr := command.Run()
		if err := verifyProbeOutput(probe.args, stdout.Bytes(), stderr.Bytes(), runErr, probe.want); err != nil {
			return err
		}
	}
	return nil
}

func verifyProbeOutput(args []string, stdout, stderr []byte, runErr error, want []byte) error {
	label := strings.Join(args, " ")
	if runErr != nil {
		return fmt.Errorf("%s: %w: %s", label, runErr, strings.TrimSpace(string(stderr)))
	}
	if len(stderr) != 0 {
		return fmt.Errorf("%s wrote stderr: %s", label, strings.TrimSpace(string(stderr)))
	}
	if !bytes.Equal(stdout, want) {
		return fmt.Errorf("%s output differs from contract", label)
	}
	return nil
}

func readAndScan(path string, forbidden []string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filepath.Base(path), err)
	}
	if err := scanForbidden(filepath.Base(path), data, forbidden); err != nil {
		return nil, err
	}
	return data, nil
}

func deriveForbidden(repoDir string, supplied []string) []string {
	values := append([]string{repoDir, filepath.ToSlash(repoDir), os.Getenv("USERPROFILE"), os.Getenv("HOME"), os.TempDir()}, supplied...)
	for _, localIdentifier := range []string{os.Getenv("USERNAME"), os.Getenv("USER"), os.Getenv("COMPUTERNAME"), os.Getenv("HOSTNAME")} {
		localIdentifier = strings.TrimSpace(localIdentifier)
		if len(localIdentifier) >= 8 {
			values = append(values, localIdentifier)
		}
	}
	unique := make(map[string]struct{})
	result := make([]string, 0, len(values)*2)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if len(value) < 3 {
			continue
		}
		for _, candidate := range []string{value, strings.ReplaceAll(value, `\`, `\\`)} {
			candidate = strings.ToLower(candidate)
			if _, exists := unique[candidate]; !exists {
				unique[candidate] = struct{}{}
				result = append(result, candidate)
			}
		}
	}
	sort.Strings(result)
	return result
}

func scanForbidden(label string, data []byte, forbidden []string) error {
	text := strings.ToLower(string(data))
	for _, needle := range forbidden {
		if strings.Contains(text, needle) {
			return fmt.Errorf("%s contains forbidden local identifier", label)
		}
	}
	structural := []string{`c:\\users\\`, `c:\users\`, `/home/`, `/users/`, `/tmp/`, `/private/var/folders/`}
	for _, needle := range structural {
		if strings.Contains(text, needle) {
			return fmt.Errorf("%s contains structural absolute path %q", label, needle)
		}
	}
	return nil
}

func safeMemberName(name string) bool {
	if name == "" || name == "." || name == ".." || filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return false
	}
	if strings.ContainsAny(name, `/\\`) || strings.ContainsRune(name, 0) || strings.Contains(name, "..") {
		return false
	}
	for _, character := range name {
		if character < 0x20 || character == 0x7f {
			return false
		}
	}
	return true
}

func safeRelativeArtifactPath(path string) bool {
	normalized := strings.ReplaceAll(path, `\`, "/")
	if strings.HasPrefix(normalized, "/") || windowsDrivePathPattern.MatchString(normalized) {
		return false
	}
	for _, component := range strings.Split(normalized, "/") {
		if component == ".." || component == "" {
			return false
		}
	}
	return filepath.VolumeName(path) == ""
}
