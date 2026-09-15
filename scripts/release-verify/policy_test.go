package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestRepositoryReleasePolicy(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	config := readPolicyFile(t, filepath.Join(repository, ".goreleaser.yaml"))
	workflow := readPolicyFile(t, filepath.Join(repository, ".github", "workflows", "release-proof.yml"))

	for _, required := range []string{
		`version_template: "1.1.0"`, "release:", "disable: true", "CGO_ENABLED=0",
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
	const stableSchemaSource = "- src: schema/releases/v1.1.0/cueson.schema.json"
	if count := strings.Count(config, stableSchemaSource); count != 1 {
		t.Errorf(".goreleaser.yaml contains %d immutable candidate schema sources, want 1", count)
	}
	if strings.Contains(config, "- src: internal/schema/") {
		t.Error("stable candidate packages the evolving schema instead of its immutable copy")
	}

	for _, required := range []string{
		"pull_request:", "workflow_dispatch:", "contents: read", "persist-credentials: false",
		"SOURCE_COMMIT: ${{ github.event.pull_request.head.sha || github.sha }}",
		"ref: ${{ env.SOURCE_COMMIT }}", "name: cueson-1.1.0-candidate-${{ env.SOURCE_COMMIT }}",
		"github.com/goreleaser/goreleaser/v2@v2.18.1", "github.com/anchore/syft/cmd/syft@v1.51.1",
		"GOTOOLCHAIN: auto", "goreleaser release --snapshot --clean --skip=publish", "-version 1.1.0", "-execute-host", "-evidence", "retention-days: 3",
		"needs: release-proof", "actions/download-artifact@", "ubuntu-24.04", "windows-2025", "macos-15-intel",
	} {
		if !strings.Contains(workflow, required) {
			t.Errorf("release-proof.yml missing %q", required)
		}
	}
	if count := strings.Count(workflow, "-version 1.1.0"); count != 2 {
		t.Errorf("release-proof.yml has %d stable verifications, want bundle and native proof", count)
	}
	if strings.Contains(workflow, "-development") {
		t.Error("stable proof uses development schema verification")
	}
	if count := strings.Count(workflow, "goreleaser release --snapshot --clean --skip=publish"); count != 1 {
		t.Errorf("release-proof.yml builds %d candidate bundles, want exactly 1", count)
	}
	if count := strings.Count(workflow, "name: cueson-1.1.0-candidate-${{ env.SOURCE_COMMIT }}"); count != 2 {
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
	canonical, schemaDigest, err := loadRepositorySchemas(repository, "1.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(canonical) == 0 || !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(schemaDigest) {
		t.Fatalf("invalid canonical schema identity: bytes=%d digest=%q", len(canonical), schemaDigest)
	}
	if schemaDigest != "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7" {
		t.Fatalf("reviewed stable schema digest changed: %s", schemaDigest)
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
	for _, version := range []string{"0.0.0", "1.0.0", "1.1.0"} {
		releaseNotes := readPolicyFile(t, filepath.Join(repository, "docs", "releases", "v"+version+".md"))
		want := "Full changelog: https://github.com/shruggietech/cueson/blob/v" + version + "/CHANGELOG.md\n"
		if !strings.HasSuffix(releaseNotes, want) {
			t.Errorf("docs/releases/v%s.md must end with exact suffix %q", version, want)
		}
	}
}

func TestStableReleaseEvidenceContract(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	data := readPolicyFile(t, filepath.Join(repository, "specs", "S028-freeze-v1-1-release-candidate", "contracts", "release-evidence-contract.schema.json"))
	var schema struct {
		Properties  map[string]json.RawMessage `json:"properties"`
		Definitions map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
			Required   []string                   `json:"required"`
			AllOf      []struct {
				Properties map[string]struct {
					Const string `json:"const"`
				} `json:"properties"`
			} `json:"allOf"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal([]byte(data), &schema); err != nil {
		t.Fatal(err)
	}
	for field, expected := range map[string]string{"version": "1.1.0", "intended_tag": "v1.1.0", "release_schema_sha256": "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7"} {
		var value struct {
			Const string `json:"const"`
		}
		if err := json.Unmarshal(schema.Properties[field], &value); err != nil {
			t.Fatal(err)
		}
		if value.Const != expected {
			t.Fatalf("stable %s contract = %q", field, value.Const)
		}
	}
	for _, target := range expectedTargets("1.1.0") {
		definition := schema.Definitions[target.GOOS+"_"+target.GOARCH]
		if len(definition.AllOf) != 2 {
			t.Fatalf("missing exact target %s/%s", target.GOOS, target.GOARCH)
		}
		for field, want := range map[string]string{"goos": target.GOOS, "goarch": target.GOARCH, "archive": target.Archive, "binary": target.Binary, "sbom": target.SBOM} {
			if definition.AllOf[1].Properties[field].Const != want {
				t.Fatalf("target field %s is not bound", field)
			}
		}
	}
	var targets struct {
		PrefixItems []struct {
			Ref string `json:"$ref"`
		} `json:"prefixItems"`
		Items   bool `json:"items"`
		Minimum int  `json:"minItems"`
		Maximum int  `json:"maxItems"`
	}
	if err := json.Unmarshal(schema.Properties["targets"], &targets); err != nil {
		t.Fatal(err)
	}
	if targets.Items || len(targets.PrefixItems) != 6 || targets.Minimum != 6 || targets.Maximum != 6 {
		t.Fatal("stable evidence does not bind exactly six ordered target tuples")
	}
	for index, target := range expectedTargets("1.1.0") {
		if targets.PrefixItems[index].Ref != "#/$defs/"+target.GOOS+"_"+target.GOARCH {
			t.Fatal("stable evidence target order differs")
		}
	}
	proof := schema.Definitions["native_proof"]
	for _, field := range []string{"formats", "historical_input_count", "old_consumer_version", "old_consumer_revision", "old_consumer_archive", "old_consumer_archive_sha256", "old_consumer_refusal_count", "old_consumer_identity_probe_count"} {
		if !containsString(proof.Required, field) {
			t.Fatalf("native proof does not require %s", field)
		}
	}
	var count struct {
		Const int `json:"const"`
	}
	if err := json.Unmarshal(proof.Properties["old_consumer_refusal_count"], &count); err != nil {
		t.Fatal(err)
	}
	if count.Const != 32 {
		t.Fatal("native proof does not require all old-consumer safety probes")
	}
}

func TestPreparedV110HistoryPreservesReleasedBytes(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	changelog := readPolicyFile(t, filepath.Join(repository, "CHANGELOG.md"))
	const preparedHeading = "## [1.1.0] - 2026-09-15\n"
	if strings.Count(changelog, "## [Unreleased]\n") != 1 || strings.Count(changelog, preparedHeading) != 1 {
		t.Fatal("preparation requires one fresh Unreleased and one dated v1.1.0 section")
	}
	prepared := strings.Index(changelog, preparedHeading)
	historical := strings.Index(changelog, "## [1.0.0] - 2026-09-11\n")
	if strings.Index(changelog, "## [Unreleased]\n") >= prepared || historical <= prepared {
		t.Fatal("release sections are not in newest-first order")
	}
	section := changelog[prepared:historical]
	seenCategories := make(map[string]bool)
	var previousDecision time.Time
	for _, line := range strings.Split(section, "\n") {
		if strings.HasPrefix(line, "### ") {
			category := strings.TrimPrefix(line, "### ")
			if seenCategories[category] {
				t.Fatalf("prepared release duplicates category %q", category)
			}
			seenCategories[category] = true
		}
		if len(line) >= 14 && line[0:2] == "- " && line[12:14] == ": " {
			date, err := time.Parse("2006-01-02", line[2:12])
			if err != nil || (!previousDecision.IsZero() && date.Before(previousDecision)) {
				t.Fatalf("invalid or nonchronological prepared decision: %s", line)
			}
			previousDecision = date
		}
	}
	for _, category := range []string{"Added", "Changed", "Decisions", "Fixed"} {
		if !seenCategories[category] {
			t.Errorf("prepared history omits %s", category)
		}
	}
	historicalEnd := strings.Index(changelog[historical:], "\n[Unreleased]:")
	if historicalEnd < 0 {
		t.Fatal("changelog comparison footer is absent")
	}
	const releasedHistorySHA256 = "b262a0e199e6f54451db3211f5c2fbccf7a959e7fbfe2e1e7653db20c170f417"
	if got := fmt.Sprintf("%x", sha256.Sum256([]byte(changelog[historical:historical+historicalEnd]))); got != releasedHistorySHA256 {
		t.Fatal("earlier released changelog bytes changed during preparation")
	}
}

func TestPreparedV110PublicNotesBoundaries(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	notes := readPolicyFile(t, filepath.Join(repository, "specs", "S029-publish-verify-v1-1", "contracts", "release-notes.md"))
	if !utf8.ValidString(notes) || strings.HasPrefix(notes, "\ufeff") || strings.Contains(notes, "\r") {
		t.Fatal("public notes are not UTF-8 without BOM and LF")
	}
	if !strings.HasPrefix(notes, "# Cueson v1.1.0\n\n") || strings.Count(notes, "\n# ") != 0 {
		t.Fatal("public notes do not have the exact single release title")
	}
	const suffix = "Full changelog: https://github.com/shruggietech/cueson/blob/v1.1.0/CHANGELOG.md\n"
	if !strings.HasSuffix(notes, suffix) {
		t.Fatal("public notes omit the exact final tagged changelog link")
	}
	for _, required := range []string{"UTF-8 ASS v4+/SSA v4", "all twelve conversion directions", "strict/fatal refusal", "Exact historical 1.0.0", "Every new encode uses exact 1.1.0", "v1.0.0 consumers reject", "no historical-output selector", "schema-only", "experimental input observations", "original bytes", "byte or pixel equivalence", "external resources", "legacy ASS/SSA codepages", "unsigned and unattested", "Public schema/site hosting remains a separately authorized outcome"} {
		if !strings.Contains(notes, required) {
			t.Errorf("public notes omit reviewed disclosure %q", required)
		}
	}
	for _, forbidden := range []string{"**Status:**", "unpublished", "No tag", "No GitHub Release", "downloads remain", "Publication-ready"} {
		if strings.Contains(notes, forbidden) {
			t.Errorf("public body contains candidate-only availability claim %q", forbidden)
		}
	}
	changelog := readPolicyFile(t, filepath.Join(repository, "CHANGELOG.md"))
	if len(notes)*2 >= len(changelog) {
		t.Fatal("public highlights are not substantially shorter than the detailed changelog")
	}
}

func TestPublicationDecisionContractBindsPreparationOnly(t *testing.T) {
	t.Parallel()
	repository := filepath.Clean(filepath.Join("..", ".."))
	contract := readPolicyFile(t, filepath.Join(repository, "specs", "S029-publish-verify-v1-1", "contracts", "publication-decision.md"))
	if regexp.MustCompile(`\b[0-9a-f]{40}\b`).MatchString(contract) {
		t.Fatal("preparation contract prematurely selects a publication revision")
	}
	for _, required := range []string{"Preparation only", "actual full lowercase 40-character post-merge", "current default-branch revision", "predicted squash result", "fresh verification", "GoReleaser v2.18.1", "Syft v1.51.1", "artifact ID", "expiry", "published: false", "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7", "LICENSE", "NOTICE", "linux/amd64", "windows/amd64", "darwin/amd64", "same accepted bundle without rebuilding", "32 old-consumer refusal paths", "32 identity probes", "no arm64 native execution", "positive byte length", "SHA-256 digest", "2026-09-15", "github-format/main.go -stdin", "without implicit banner removal", "Create tag", "Push tag", "Publish GitHub Release", "Publish release assets", "Independently read back", "expired artifact", "invalidates the pending decision", "newly reviewed decision package", "upload is partial", "stop all dependent publication", "operator recovery judgment", "Close #66 and unblock #67 only after"} {
		if !strings.Contains(contract, required) {
			t.Errorf("decision contract omits governed obligation %q", required)
		}
	}
	want := make(map[string]string, 13)
	for _, target := range expectedTargets("1.1.0") {
		want[target.Archive] = "Archive|" + target.GOOS + "/" + target.GOARCH
		want[target.SBOM] = "SPDX JSON SBOM|" + target.GOOS + "/" + target.GOARCH
	}
	want["cueson_1.1.0_checksums.txt"] = "Checksum manifest|Six archives"
	seen := make(map[string]bool, len(want))
	for _, line := range strings.Split(contract, "\n") {
		if !strings.HasPrefix(line, "| `cueson_") {
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) != 5 {
			t.Fatalf("invalid publication inventory row %q", line)
		}
		name := strings.Trim(strings.TrimSpace(fields[1]), "`")
		kind := strings.TrimSpace(fields[2])
		target := strings.Trim(strings.TrimSpace(fields[3]), "`")
		if seen[name] || want[name] != kind+"|"+target {
			t.Fatalf("unexpected, duplicate or wrong-target publication asset %q", name)
		}
		seen[name] = true
	}
	if len(seen) != len(want) {
		t.Fatalf("publication contract lists %d assets, want exactly thirteen", len(seen))
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
