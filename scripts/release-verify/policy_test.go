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
		"github.com/shruggietech/cueson/internal/version.current={{ .Version }}",
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

	for _, required := range []string{
		"pull_request:", "workflow_dispatch:", "contents: read", "persist-credentials: false",
		"github.com/goreleaser/goreleaser/v2@v2.18.1", "github.com/anchore/syft/cmd/syft@v1.51.1",
		"GOTOOLCHAIN: auto", "goreleaser release --snapshot --clean --skip=publish", "retention-days: 3",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("release-proof.yml missing %q", required)
		}
	}
	for _, forbidden := range []string{"pull_request_target:", "contents: write", "id-token: write", "secrets.", "GITHUB_TOKEN", "tags:", "release:"} {
		if strings.Contains(workflow, forbidden) {
			t.Errorf("release-proof.yml contains forbidden authority %q", forbidden)
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

func readPolicyFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
