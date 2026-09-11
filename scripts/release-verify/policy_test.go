package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestRepositoryReleasePolicy(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	config := readPolicyFile(t, filepath.Join(repository, ".goreleaser.yaml"))
	workflow := readPolicyFile(t, filepath.Join(repository, ".github", "workflows", "release-proof.yml"))

	for _, required := range []string{
		`version_template: "1.0.0"`, "release:", "disable: true", "CGO_ENABLED=0",
		"github.com/shruggietech/cueson/internal/version.releaseOverride=cueson-release-version:{{ .Version }}",
		"cueson_{{ .Version }}_checksums.txt", "artifacts: binary", "spdx-json=$document",
	} {
		if !strings.Contains(config, required) {
			t.Errorf(".goreleaser.yaml missing %q", required)
		}
	}
	for _, forbidden := range []string{"signs:", "attestations:", "dockers:", "brews:", "announce:", "blobs:"} {
		if strings.Contains(config, forbidden) {
			t.Errorf(".goreleaser.yaml contains publishing surface %q", forbidden)
		}
	}
	const immutableSchemaSource = "- src: schema/releases/v1.0.0/cueson.schema.json"
	if count := strings.Count(config, immutableSchemaSource); count != 1 {
		t.Errorf(".goreleaser.yaml contains %d immutable v1 schema sources, want 1", count)
	}
	if strings.Contains(config, "- src: internal/schema/cueson.schema.json") || strings.Contains(config, "- src: schema/releases/v0.0.0/cueson.schema.json") {
		t.Error("candidate packages a non-v1 schema source")
	}

	for _, required := range []string{
		"pull_request:", "workflow_dispatch:", "contents: read", "persist-credentials: false",
		"SOURCE_COMMIT: ${{ github.event.pull_request.head.sha || github.sha }}",
		"ref: ${{ env.SOURCE_COMMIT }}", "name: cueson-1.0.0-candidate-${{ env.SOURCE_COMMIT }}",
		"github.com/goreleaser/goreleaser/v2@v2.18.1", "github.com/anchore/syft/cmd/syft@v1.51.1",
		"GOTOOLCHAIN: auto", "goreleaser release --snapshot --clean --skip=publish", "-version 1.0.0", "-execute-host", "-evidence", "retention-days: 3",
		"needs: release-proof", "actions/download-artifact@", "ubuntu-24.04", "windows-2025", "macos-15-intel",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("release-proof.yml missing %q", required)
		}
	}
	if strings.Contains(workflow, "-development") {
		t.Error("release-proof.yml still enables development verification")
	}
	if count := strings.Count(workflow, "goreleaser release --snapshot --clean --skip=publish"); count != 1 {
		t.Errorf("release-proof.yml builds %d candidate bundles, want exactly 1", count)
	}
	if count := strings.Count(workflow, "name: cueson-1.0.0-candidate-${{ env.SOURCE_COMMIT }}"); count != 2 {
		t.Errorf("release-proof.yml names the accepted bundle %d times, want one upload and one download", count)
	}
	pushPattern := regexp.MustCompile(`(?m)^  push:\n((?: {4,}.*\n)*)`)
	pushMatch := pushPattern.FindStringSubmatch(workflow)
	if len(pushMatch) != 2 || pushMatch[1] != "    branches:\n      - main\n" {
		t.Errorf("release-proof.yml push trigger must contain only main, got %q", pushMatch)
	}
	for _, forbidden := range []string{
		"pull_request_target:", "secrets.", "GITHUB_TOKEN", "github.token", "persist-credentials: true", "cosign", "gh release", "action-gh-release",
	} {
		if strings.Contains(workflow, forbidden) {
			t.Errorf("release-proof.yml contains forbidden authority %q", forbidden)
		}
	}
	forbiddenAuthorityKey := regexp.MustCompile(`(?m)^\s*(?:id-token|token|secrets|tags|tags-ignore|environment|deployment|deploy|signing|signs|sign|attestations|attest|publication|publish|release):`)
	if match := forbiddenAuthorityKey.FindString(workflow); match != "" {
		t.Errorf("release-proof.yml contains forbidden authority key %q", strings.TrimSpace(match))
	}
	writePermission := regexp.MustCompile(`(?m)^\s+[a-zA-Z0-9_-]+:\s*write\s*$`)
	if writePermission.MatchString(workflow) {
		t.Error("release-proof.yml contains a write permission")
	}
	for _, line := range strings.Split(workflow, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "run: goreleaser release") && trimmed != "run: goreleaser release --snapshot --clean --skip=publish" {
			t.Errorf("release-proof.yml contains a publishing-capable GoReleaser command %q", trimmed)
		}
		if strings.HasPrefix(trimmed, "goreleaser release") && trimmed != "goreleaser release --snapshot --clean --skip=publish" {
			t.Errorf("release-proof.yml contains a publishing-capable GoReleaser command %q", trimmed)
		}
	}
	actionPattern := regexp.MustCompile(`uses:\s+([^\s@]+)@([^\s]+)`)
	for _, match := range actionPattern.FindAllStringSubmatch(workflow, -1) {
		if !strings.HasPrefix(match[1], "actions/") {
			t.Errorf("workflow action %s is not GitHub-owned", match[1])
		}
		if !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(match[2]) {
			t.Errorf("workflow action %s is not pinned to a full commit", match[1])
		}
	}
}

func TestReleaseEvidenceContractBindsExactTargetTuples(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	data := []byte(readPolicyFile(t, filepath.Join(repository, "specs", "S019-prepare-v1-release-candidate", "contracts", "release-evidence.schema.json")))
	var schema struct {
		Required   []string `json:"required"`
		Properties struct {
			Version struct {
				Const string `json:"const"`
			} `json:"version"`
			IntendedTag struct {
				Const string `json:"const"`
			} `json:"intended_tag"`
			Published struct {
				Const bool `json:"const"`
			} `json:"published"`
			Targets struct {
				PrefixItems []struct {
					Ref string `json:"$ref"`
				} `json:"prefixItems"`
				Items bool `json:"items"`
			} `json:"targets"`
		} `json:"properties"`
		Definitions map[string]struct {
			AllOf []struct {
				Properties map[string]struct {
					Const string `json:"const"`
				} `json:"properties"`
			} `json:"allOf"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	if schema.Properties.Version.Const != "1.0.0" || schema.Properties.IntendedTag.Const != "v1.0.0" || schema.Properties.Published.Const {
		t.Fatalf("release evidence contract does not bind the v1 non-publishing identity")
	}
	for _, required := range []string{"version", "intended_tag", "source_revision", "release_schema_sha256", "license_sha256", "notice_sha256", "archive_count", "sbom_count", "checksum_count", "host_executed", "targets", "published"} {
		if !containsString(schema.Required, required) {
			t.Errorf("release evidence contract does not require %q", required)
		}
	}
	want := []struct {
		name, goos, goarch, archive, binary, sbom string
	}{
		{"linux_amd64", "linux", "amd64", "cueson_1.0.0_linux_amd64.tar.gz", "cueson", "cueson_1.0.0_linux_amd64.tar.gz.sbom.json"},
		{"linux_arm64", "linux", "arm64", "cueson_1.0.0_linux_arm64.tar.gz", "cueson", "cueson_1.0.0_linux_arm64.tar.gz.sbom.json"},
		{"darwin_amd64", "darwin", "amd64", "cueson_1.0.0_darwin_amd64.tar.gz", "cueson", "cueson_1.0.0_darwin_amd64.tar.gz.sbom.json"},
		{"darwin_arm64", "darwin", "arm64", "cueson_1.0.0_darwin_arm64.tar.gz", "cueson", "cueson_1.0.0_darwin_arm64.tar.gz.sbom.json"},
		{"windows_amd64", "windows", "amd64", "cueson_1.0.0_windows_amd64.zip", "cueson.exe", "cueson_1.0.0_windows_amd64.zip.sbom.json"},
		{"windows_arm64", "windows", "arm64", "cueson_1.0.0_windows_arm64.zip", "cueson.exe", "cueson_1.0.0_windows_arm64.zip.sbom.json"},
	}
	if schema.Properties.Targets.Items {
		t.Fatal("release evidence contract permits target items after prefixItems")
	}
	if len(schema.Properties.Targets.PrefixItems) != len(want) {
		t.Fatalf("prefixItems = %d, want %d", len(schema.Properties.Targets.PrefixItems), len(want))
	}
	for index, expected := range want {
		if got := schema.Properties.Targets.PrefixItems[index].Ref; got != "#/$defs/"+expected.name {
			t.Errorf("prefixItems[%d] = %q, want #/$defs/%s", index, got, expected.name)
		}
		definition, ok := schema.Definitions[expected.name]
		if !ok || len(definition.AllOf) != 2 {
			t.Fatalf("definition %s missing exact target constraints", expected.name)
		}
		properties := definition.AllOf[1].Properties
		checks := map[string]string{"goos": expected.goos, "goarch": expected.goarch, "archive": expected.archive, "binary": expected.binary, "sbom": expected.sbom}
		for name, value := range checks {
			if got := properties[name].Const; got != value {
				t.Errorf("definition %s %s = %q, want %q", expected.name, name, got, value)
			}
		}
	}
}

func TestRepositoryCandidateSchemaAndLegalIdentity(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	canonical, schemaDigest, err := loadRepositorySchemas(repository, "1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(canonical) == 0 || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(schemaDigest) {
		t.Fatalf("invalid canonical schema identity: bytes=%d digest=%q", len(canonical), schemaDigest)
	}
	for _, name := range []string{"LICENSE", "NOTICE"} {
		file, err := loadRepositoryFile(repository, name)
		if err != nil {
			t.Fatalf("load %s: %v", name, err)
		}
		if len(file.data) == 0 || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(file.sha256) {
			t.Errorf("invalid %s identity: bytes=%d digest=%q", name, len(file.data), file.sha256)
		}
	}
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestReleaseNotesFinalSuffix(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	for _, version := range []string{"0.0.0", "1.0.0"} {
		releaseNotes := readPolicyFile(t, filepath.Join(repository, "docs", "releases", "v"+version+".md"))
		want := "Full changelog: https://github.com/shruggietech/cueson/blob/v" + version + "/CHANGELOG.md\n"
		if !strings.HasSuffix(releaseNotes, want) {
			t.Errorf("docs/releases/v%s.md must end with exact suffix %q", version, want)
		}
	}
}

func readPolicyFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
