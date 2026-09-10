package main

import (
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
