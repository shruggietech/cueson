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
		`version_template: "0.0.0"`, "release:", "disable: true", "CGO_ENABLED=0",
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
	const releaseSchemaSource = "- src: schema/releases/v0.0.0/cueson.schema.json"
	if count := strings.Count(config, releaseSchemaSource); count != 1 {
		t.Errorf(".goreleaser.yaml contains %d exact release-schema sources, want 1", count)
	}
	if strings.Contains(config, "- src: internal/schema/cueson.schema.json") {
		t.Error(".goreleaser.yaml still packages the mutable canonical schema")
	}

	for _, required := range []string{
		"pull_request:", "workflow_dispatch:", "contents: read", "persist-credentials: false",
		"ref: ${{ github.event.pull_request.head.sha || github.sha }}",
		"github.com/goreleaser/goreleaser/v2@v2.18.1", "github.com/anchore/syft/cmd/syft@v1.51.1",
		"GOTOOLCHAIN: auto", "goreleaser release --snapshot --clean --skip=publish", "retention-days: 3",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("release-proof.yml missing %q", required)
		}
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
	data := []byte(readPolicyFile(t, filepath.Join(repository, "specs", "S012-prepare-v0-release", "contracts", "release-evidence.schema.json")))
	var schema struct {
		Properties struct {
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
	want := []struct {
		name, goos, goarch, archive, binary, sbom string
	}{
		{"linux_amd64", "linux", "amd64", "cueson_0.0.0_linux_amd64.tar.gz", "cueson", "cueson_0.0.0_linux_amd64.tar.gz.sbom.json"},
		{"linux_arm64", "linux", "arm64", "cueson_0.0.0_linux_arm64.tar.gz", "cueson", "cueson_0.0.0_linux_arm64.tar.gz.sbom.json"},
		{"darwin_amd64", "darwin", "amd64", "cueson_0.0.0_darwin_amd64.tar.gz", "cueson", "cueson_0.0.0_darwin_amd64.tar.gz.sbom.json"},
		{"darwin_arm64", "darwin", "arm64", "cueson_0.0.0_darwin_arm64.tar.gz", "cueson", "cueson_0.0.0_darwin_arm64.tar.gz.sbom.json"},
		{"windows_amd64", "windows", "amd64", "cueson_0.0.0_windows_amd64.zip", "cueson.exe", "cueson_0.0.0_windows_amd64.zip.sbom.json"},
		{"windows_arm64", "windows", "arm64", "cueson_0.0.0_windows_arm64.zip", "cueson.exe", "cueson_0.0.0_windows_arm64.zip.sbom.json"},
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

func TestReleaseNotesFinalSuffix(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	releaseNotes := readPolicyFile(t, filepath.Join(repository, "docs", "releases", "v0.0.0.md"))
	const want = "Full changelog: https://github.com/shruggietech/cueson/blob/v0.0.0/CHANGELOG.md\n"
	if !strings.HasSuffix(releaseNotes, want) {
		t.Errorf("docs/releases/v0.0.0.md must end with exact suffix %q", want)
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
