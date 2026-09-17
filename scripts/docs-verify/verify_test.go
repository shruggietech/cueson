package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunCLI(t *testing.T) {
	repo := newRepository(t)
	var stdout, stderr bytes.Buffer
	if code := runCLI([]string{"-repo", repo}, &stdout, &stderr); code != 0 {
		t.Fatalf("runCLI() = %d, stderr = %q", code, stderr.String())
	}
	if got := stdout.String(); !strings.Contains(got, "checked 24 documents") || !strings.Contains(got, "0 local links") || !strings.Contains(got, "10 registered examples") || !strings.Contains(got, "2 format rows") {
		t.Fatalf("unexpected success output: %q", got)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runCLI([]string{"-repo", repo, "extra"}, &stdout, &stderr); code != 2 {
		t.Fatalf("runCLI(extra) = %d, want 2", code)
	}

	stdout.Reset()
	stderr.Reset()
	if err := os.Remove(filepath.Join(repo, filepath.FromSlash("docs/formats/srt.md"))); err != nil {
		t.Fatal(err)
	}
	if code := runCLI([]string{"-repo", repo}, &stdout, &stderr); code != 1 {
		t.Fatalf("runCLI(violations) = %d, want 1", code)
	}

	stdout.Reset()
	stderr.Reset()
	if code := runCLI([]string{"-repo", filepath.Join(repo, "absent")}, &stdout, &stderr); code != 2 {
		t.Fatalf("runCLI(absent root) = %d, want 2", code)
	}
}

func TestRequiredDocuments(t *testing.T) {
	repo := newRepository(t)
	if err := os.Remove(filepath.Join(repo, filepath.FromSlash("docs/formats/srt.md"))); err != nil {
		t.Fatal(err)
	}
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/formats/srt.md: required document is missing")
}

func TestRegisteredExamplesAreCompleteAndUnique(t *testing.T) {
	repo := newRepository(t)
	readme := strings.Replace(readFileForVerifierTest(t, repo, "README.md"), "<!-- docs-verify:example encode -->\n", "", 1)
	writeFile(t, repo, "README.md", readme+"<!-- docs-verify:example source-run -->\n<!-- docs-verify:example unknown -->\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "README.md: registered example is missing: encode")
	assertViolation(t, result.violations, "README.md: registered example appears more than once: source-run")
	assertViolation(t, result.violations, "README.md: unknown registered example: unknown")
}

func TestReferenceMarkersAndStaleClaims(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/schema.md", "# Schema\n\n$id schema_version format_support format_data source v0.0.0 v1.0.0\n")
	writeFile(t, repo, "CONTRIBUTING.md", "# Contributing\n\nCueson is an unreleased v0.0.0 foundation candidate.\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/schema.md: required contract marker is missing: non-normative")
	assertViolation(t, result.violations, "CONTRIBUTING.md: stale capability or release claim remains: Cueson is an unreleased v0.0.0 foundation candidate")
}

func TestPublishedV1DocsRejectCandidateClaims(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/formats/srt.md", readFileForVerifierTest(t, repo, "docs/formats/srt.md")+"\nThat evidence supports the stable candidate declaration; tag and release publication remain separate.\n")
	writeFile(t, repo, "docs/cueson-media-format-guide.html", "<p>Current v1.0.0 candidate boundary: The candidate is not yet published. SRT is a v1.0.0 release candidate.</p>\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/formats/srt.md: stale capability or release claim remains: stable candidate declaration; tag and release publication remain separate")
	assertViolation(t, result.violations, "docs/cueson-media-format-guide.html: stale capability or release claim remains: Current v1.0.0 candidate boundary:")
	assertViolation(t, result.violations, "docs/cueson-media-format-guide.html: stale capability or release claim remains: The candidate is not yet published")
	assertViolation(t, result.violations, "docs/cueson-media-format-guide.html: stale capability or release claim remains: v1.0.0 release candidate")
}

func TestFormatMatrixRowsMustAppearInMatchingGuide(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/formats/srt.md", "# SRT\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/formats/srt.md: format guide lacks matrix row: srt-example")
}

func TestReleaseNotesAreRequired(t *testing.T) {
	for _, name := range []string{"docs/releases/v0.0.0.md", "docs/releases/v1.0.0.md", "docs/releases/v1.1.0.md"} {
		t.Run(name, func(t *testing.T) {
			repo := newRepository(t)
			if err := os.Remove(filepath.Join(repo, filepath.FromSlash(name))); err != nil {
				t.Fatal(err)
			}
			result, err := verifyRepository(repo)
			if err != nil {
				t.Fatal(err)
			}
			assertViolation(t, result.violations, name+": required document is missing")
		})
	}
}

func TestV1ReleaseNotesRequireExactTaggedChangelogSuffix(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/releases/v1.0.0.md", "# Cueson v1.0.0\n\nstable Cue JSON; unsigned and unattested\n\nFull changelog: https://github.com/shruggietech/cueson/blob/main/CHANGELOG.md\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/releases/v1.0.0.md: required final suffix is missing: Full changelog: https://github.com/shruggietech/cueson/blob/v1.0.0/CHANGELOG.md")
}

func TestV1ReleaseNotesRejectPrematurePublicationClaim(t *testing.T) {
	repo := newRepository(t)
	notes := strings.Replace(readFileForVerifierTest(t, repo, "docs/releases/v1.0.0.md"), "Full changelog:", "Cueson v1.0.0 is now published.\n\nFull changelog:", 1)
	writeFile(t, repo, "docs/releases/v1.0.0.md", notes)
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/releases/v1.0.0.md: stale capability or release claim remains: is now published")
}

func TestScriptedFreezeRejectsStaleCapabilityAndPrematureRelease(t *testing.T) {
	repo := newRepository(t)
	native := readFileForVerifierTest(t, repo, "docs/formats/ass-ssa.md")
	writeFile(t, repo, "docs/formats/ass-ssa.md", native+"\nUnavailable for scripted formats.\n")
	notes := readFileForVerifierTest(t, repo, "docs/releases/v1.1.0.md")
	writeFile(t, repo, "docs/releases/v1.1.0.md", strings.Replace(notes, "Full changelog:", "Download Cueson v1.1.0.\n\nFull changelog:", 1))
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/formats/ass-ssa.md: stale capability or release claim remains: Unavailable for scripted formats")
	assertViolation(t, result.violations, "docs/releases/v1.1.0.md: stale capability or release claim remains: Download Cueson v1.1.0")
}

func TestCandidateRequiresExactHistoricalRefusalAndProspectiveNotes(t *testing.T) {
	repo := newRepository(t)
	readme := readFileForVerifierTest(t, repo, "README.md")
	writeFile(t, repo, "README.md", strings.ReplaceAll(readme, "published v1.0.0 executable rejects new 1.1.0 output", "older consumers may accept new output"))
	notes := readFileForVerifierTest(t, repo, "docs/releases/v1.1.0.md")
	writeFile(t, repo, "docs/releases/v1.1.0.md", strings.ReplaceAll(notes, "blob/v1.1.0/CHANGELOG.md", "blob/main/CHANGELOG.md"))
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "README.md: required contract marker is missing: published v1.0.0 executable rejects new 1.1.0 output")
	assertViolation(t, result.violations, "docs/releases/v1.1.0.md: required final suffix is missing: Full changelog: https://github.com/shruggietech/cueson/blob/v1.1.0/CHANGELOG.md")
}

func TestV1ChangelogRequiresUnreleasedBeforeReleaseSection(t *testing.T) {
	repo := newRepository(t)
	changelog := strings.Replace(readFileForVerifierTest(t, repo, "CHANGELOG.md"), "## [Unreleased]\n\n### Added\n\n- Post-release work.\n\n", "", 1)
	writeFile(t, repo, "CHANGELOG.md", changelog)
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "CHANGELOG.md: prepared sections must begin Unreleased, dated 1.1.0, then historical 1.0.0")
}

func TestCurrentV1CandidateAndWorkingSpecificationMarkers(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "README.md", readFileForVerifierTest(t, repo, "README.md")+"v0.1.0 development\n")
	writeFile(t, repo, "docs/Cueson-Project-Specification-v0.0.0.md", "# Working specification\n\n- [ ] public v1.0.0 schema matches the repository artifact exactly;\n- [ ] release verification issue is complete;\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "README.md: stale capability or release claim remains: v0.1.0 development")
	assertViolation(t, result.violations, "docs/Cueson-Project-Specification-v0.0.0.md: required contract marker is missing: canonical, immutable repository, embedded, emitted, and packaged v1.0.0 schema copies match byte-for-byte")
	assertViolation(t, result.violations, "docs/Cueson-Project-Specification-v0.0.0.md: stale capability or release claim remains: public v1.0.0 schema matches the repository artifact exactly")
}

func TestPublishedV110CurrentStateMarkersAreRequired(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		marker string
	}{
		{name: "readme release", path: "README.md", marker: "v1.1.0 released and independently verified"},
		{name: "readme GitHub release", path: "README.md", marker: "v1.1.0 GitHub Release"},
		{name: "schema release", path: "docs/schema.md", marker: "v1.1.0 schema released and independently verified"},
		{name: "compatibility release", path: "docs/compatibility.md", marker: "v1.1.0 stable release published and independently verified"},
		{name: "immutable schema copies", path: "docs/Cueson-Project-Specification-v0.0.0.md", marker: "canonical, immutable repository, embedded, emitted, and packaged v1.1.0 schema copies match byte-for-byte"},
		{name: "release notes publication", path: "docs/releases/v1.1.0.md", marker: "public v1.1.0 GitHub Release"},
		{name: "release checksum", path: "docs/releases/v1.1.0.md", marker: "cueson_1.1.0_checksums.txt"},
		{name: "scripted format publication", path: "docs/formats/ass-ssa.md", marker: "published v1.1.0 release"},
		{name: "roadmap issue split", path: "docs/roadmap.md", marker: "S030 child [#76]"},
		{name: "project issue split", path: "docs/project-management.md", marker: "S030 child #76 owns the reviewed 1.1.0 schema/site artifact"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newRepository(t)
			content := readFileForVerifierTest(t, repo, test.path)
			writeFile(t, repo, test.path, strings.ReplaceAll(content, test.marker, "removed published-state marker"))
			result, err := verifyRepository(repo)
			if err != nil {
				t.Fatal(err)
			}
			assertViolation(t, result.violations, test.path+": required contract marker is missing: "+test.marker)
		})
	}
}

func TestPublishedV110DocsRejectCandidateEraClaims(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		claim string
	}{
		{name: "readme unpublished candidate", path: "README.md", claim: "unpublished 1.1.0 stable candidate"},
		{name: "readme release absent", path: "README.md", claim: "The candidate is not tagged or publicly released"},
		{name: "readme old downloads", path: "README.md", claim: "Published downloads below remain v1.0.0"},
		{name: "release notes unpublished candidate", path: "docs/releases/v1.1.0.md", claim: "Publication-ready notes for the unpublished 1.1.0 stable candidate"},
		{name: "release notes release absent", path: "docs/releases/v1.1.0.md", claim: "No tag or GitHub Release has been published"},
		{name: "release notes old downloads", path: "docs/releases/v1.1.0.md", claim: "public downloads remain v1.0.0"},
		{name: "scripted format unpublished", path: "docs/formats/ass-ssa.md", claim: "unpublished 1.1.0 candidate"},
		{name: "schema unpublished", path: "docs/schema.md", claim: "The 1.1.0 candidate is unpublished"},
		{name: "compatibility candidate", path: "docs/compatibility.md", claim: "Frozen 1.1.0 stable candidate"},
		{name: "release process candidate", path: "docs/release-process.md", claim: "## v1.1.0 stable candidate preparation"},
		{name: "release issue open", path: "docs/release-process.md", claim: "#66 remains open until independent public verification"},
		{name: "roadmap direct production ownership", path: "docs/roadmap.md", claim: "Public schema/site hosting remains S030 issue #67"},
		{name: "roadmap stale preparation", path: "docs/roadmap.md", claim: "S029 preparation is active"},
		{name: "roadmap direct deployment", path: "docs/roadmap.md", claim: "Reviewed artifact and explicitly authorized exact-main deployment"},
		{name: "project direct production ownership", path: "docs/project-management.md", claim: "S030 issue #67 owns the public 1.1.0 schema/site continuation"},
		{name: "project omitted preparation child", path: "docs/project-management.md", claim: "S030 owns the remaining public schema/site outcome #67"},
		{name: "project production snapshot", path: "docs/project-management.md", claim: "Production remains the reviewed S021 revision"},
		{name: "compatibility temporal absence", path: "docs/compatibility.md", claim: "The 1.1.0 schema URI is not yet publicly hosted"},
		{name: "CLI temporal absence", path: "docs/cli.md", claim: "still-pending public schema route"},
		{name: "schema temporal preparation", path: "docs/schema.md", claim: "issue #67 prepares the production route"},
		{name: "architecture temporal availability", path: "docs/architecture.md", claim: "Public schema/site availability remains #67"},
		{name: "conversion temporal hosting", path: "docs/conversion.md", claim: "Public schema/site hosting remains the separate #67 outcome"},
		{name: "scripted format temporal hosting", path: "docs/formats/ass-ssa.md", claim: "public 1.1.0 schema/site availability remains #67"},
		{name: "release process temporal hosting", path: "docs/release-process.md", claim: "Production schema hosting remains pending under #67"},
		{name: "roadmap production snapshot", path: "docs/roadmap.md", claim: "The production site remains deployed from reviewed S021 revision"},
		{name: "media guide production snapshot", path: "docs/cueson-media-format-guide.html", claim: "The current production site still serves the v1.0.0-era pages"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newRepository(t)
			content := readFileForVerifierTest(t, repo, test.path)
			if suffix, ok := requiredSuffixes[test.path]; ok {
				content = strings.TrimSuffix(content, suffix) + test.claim + "\n\n" + suffix
			} else {
				content += "\n" + test.claim + "\n"
			}
			writeFile(t, repo, test.path, content)
			result, err := verifyRepository(repo)
			if err != nil {
				t.Fatal(err)
			}
			assertViolation(t, result.violations, test.path+": stale capability or release claim remains: "+test.claim)
		})
	}
}

func TestPublishedV110RulesPreserveFrozenSliceEvidence(t *testing.T) {
	repo := newRepository(t)
	historical := "unpublished 1.1.0 stable candidate\nNo tag or GitHub Release has been published\npublic downloads remain v1.0.0\n"
	writeFile(t, repo, "specs/S028-freeze-v1-1-release-candidate/verification.md", historical)
	writeFile(t, repo, "specs/S029-publish-verify-v1-1/verification.md", historical)
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.violations) != 0 {
		t.Fatalf("frozen S028/S029 evidence must retain its time-bound contract: %v", result.violations)
	}
}

func TestReleaseNotesLocalReference(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/releases/v0.0.0.md", "# Cueson v0.0.0\n\nSee the [schema contract](../schema.md).\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.violations) != 0 {
		t.Fatalf("unexpected violations: %v", result.violations)
	}
	if result.localLinks != 1 {
		t.Fatalf("localLinks = %d, want 1", result.localLinks)
	}
}

func TestRequiredDocumentMustBeRegular(t *testing.T) {
	repo := newRepository(t)
	name := filepath.Join(repo, filepath.FromSlash("docs/formats/srt.md"))
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(name, 0o755); err != nil {
		t.Fatal(err)
	}
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/formats/srt.md: required document is not a regular file")
}

func TestRequiredDocumentMustNotBeSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks normally requires an elevated Windows token")
	}
	repo := newRepository(t)
	name := filepath.Join(repo, filepath.FromSlash("docs/formats/srt.md"))
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(repo, "README.md"), name); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/formats/srt.md: required document is not a regular file")
}

func TestAcceptedReferences(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/target file.md", "# One Heading\n\n## Repeat\n\n## Repeat\n\nSetext heading\n--------------\n\n<a id=\"explicit-anchor\"></a>\n")
	writeFile(t, repo, "docs/assets/pixel.png", "not really an image")
	writeFile(t, repo, "docs/link-source.md", strings.Join([]string{
		"# Links",
		"",
		"[relative](target%20file.md#one-heading)",
		"[parent](../README.md#repository)",
		"[root](/docs/target%20file.md#repeat-1)",
		"[fragment](#links)",
		"![image](assets/pixel.png)",
		"[reference][target]",
		"[shortcut]",
		"[target]: target%20file.md#setext-heading",
		"[shortcut]: target%20file.md#repeat",
		"<https://example.com/path>",
		"<mailto:maintainer@example.com>",
		"<a href=\"//example.com/asset\">network path</a>",
		"<a href=\"target%20file.md#explicit-anchor\"><img src='assets/pixel.png'></a>",
		"[directory](assets/)",
		"`[ignored](missing-inline.md)`",
		"\\[literal](missing-literal.md)",
		"[escaped close\\](missing-close.md)",
		"```markdown",
		"[ignored](missing-fenced.md)",
		"```",
		"<!-- [ignored](missing-comment.md) -->",
		"<style>@import url(\"missing.css\");</style>",
		"",
	}, "\n"))
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.violations) != 0 {
		t.Fatalf("unexpected violations: %v", result.violations)
	}
	if result.localLinks != 12 {
		t.Fatalf("localLinks = %d, want 12", result.localLinks)
	}
}

func TestRejectedReferencesAndDeterministicOrder(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/existing.md", "# Existing\n")
	writeFile(t, repo, "docs/rejected.md", strings.Join([]string{
		"# Rejected",
		"",
		"[missing](z-missing.md)",
		"[anchor](existing.md#absent)",
		"[case](Existing.md)",
		"[escape](../../outside.md)",
		"[backslash](existing\\bad.md)",
		"[bad escape](bad%ZZ.md)",
		"[scheme](ftp://example.com/file)",
		"[empty]()",
		"[control](bad\x01name.md)",
		"",
	}, "\n"))
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"docs/rejected.md:3: local target does not exist: z-missing.md",
		"docs/rejected.md:4: fragment does not exist: absent",
		"docs/rejected.md:5: local target has incorrect case: Existing.md",
		"docs/rejected.md:6: local target escapes repository: ../../outside.md",
		"docs/rejected.md:7: local target contains a backslash: existing\\bad.md",
		"docs/rejected.md:8: malformed URL escape: bad%ZZ.md",
		"docs/rejected.md:9: unsupported URL scheme: ftp",
		"docs/rejected.md:10: empty link target",
		"docs/rejected.md:11: link target contains a control character",
	}
	if got := violationStrings(result.violations); strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("violations:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestHTMLAnchorsAndNonMarkdownFragments(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/page.html", "<!doctype html><h1 id=\"named\">Page</h1><a name='legacy'></a>\n")
	writeFile(t, repo, "docs/data.txt", "data\n")
	writeFile(t, repo, "docs/html-links.md", "[id](page.html#named) [name](page.html#legacy) [bad](data.txt#fragment)\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/html-links.md:1: fragment is not supported for target: data.txt")
}

func TestSymlinkTargetRejected(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("creating symlinks normally requires an elevated Windows token")
	}
	repo := newRepository(t)
	target := filepath.Join(repo, "docs", "real.md")
	writeFile(t, repo, "docs/real.md", "# Real\n")
	if err := os.Symlink(target, filepath.Join(repo, "docs", "linked.md")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	writeFile(t, repo, "docs/symlink-source.md", "[linked](linked.md)\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/symlink-source.md:1: local target is a symbolic link: linked.md")
}

func TestDuplicateAndSetextSlugs(t *testing.T) {
	anchors := markdownAnchors("# Héllo, World!\n\n# Héllo, World!\n\nA Title\n=======\n\n## `version`\n\n## [Linked heading](https://example.com)\n\n## `[Literal](https://example.com/path)`\n")
	for _, want := range []string{"héllo-world", "héllo-world-1", "a-title", "version", "linked-heading", "literalhttpsexamplecompath"} {
		if _, ok := anchors[want]; !ok {
			t.Errorf("missing anchor %q in %v", want, anchors)
		}
	}
}

func TestCodeHTMLAnchorsIgnored(t *testing.T) {
	anchors := documentAnchors("docs/example.md", "```html\n<a id=\"fenced-example\"></a>\n```\n\n`<a id=\"inline-example\"></a>`\n\n<a id=\"real-anchor\"></a>\n")
	if _, ok := anchors["fenced-example"]; ok {
		t.Errorf("fenced example became an anchor: %v", anchors)
	}
	if _, ok := anchors["inline-example"]; ok {
		t.Errorf("inline example became an anchor: %v", anchors)
	}
	if _, ok := anchors["real-anchor"]; !ok {
		t.Errorf("missing real explicit anchor: %v", anchors)
	}
}

func newRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	for _, name := range requiredDocuments {
		heading := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
		if name == "README.md" {
			heading = "Repository"
		}
		writeFile(t, repo, name, "# "+heading+"\n")
	}
	var readme strings.Builder
	readme.WriteString("# Repository\n\n")
	for _, id := range requiredExampleIDs {
		readme.WriteString("<!-- docs-verify:example " + id + " -->\n")
	}
	writeFile(t, repo, "README.md", readme.String())
	writeFile(t, repo, "docs/schema.md", "# Schema\n\n$id schema_version format_support format_data source v0.0.0 v1.0.0 non-normative\n")
	writeFile(t, repo, "CHANGELOG.md", "# Changelog\n\n## [Unreleased]\n\n### Added\n\n- Post-release work.\n\n## [1.1.0] - 2026-09-15\n\n### Added\n\n- Prepared candidate history.\n\n## [1.0.0] - 2026-09-11\n\n[Unreleased]: https://github.com/shruggietech/cueson/compare/v1.1.0...HEAD\n[1.1.0]: https://github.com/shruggietech/cueson/compare/v1.0.0...v1.1.0\n[1.0.0]: https://github.com/shruggietech/cueson/compare/v0.0.0...v1.0.0\n")
	writeFile(t, repo, "README.md", readme.String()+"\nv1.0.0 released and independently verified; v0.0.0 remains historical; v1.0.0 GitHub Release.\n")
	writeFile(t, repo, "docs/schema.md", "# Schema\n\n$id schema_version format_support format_data source v0.0.0 v1.0.0 non-normative stable schema released and independently verified\n")
	writeFile(t, repo, "docs/compatibility.md", "# Compatibility\n\nCLI Cue JSON Schema internal/ v0.0.0 v1.0.0 Windows macOS Linux production stable release published and independently verified\n")
	writeFile(t, repo, "docs/Cueson-Project-Specification-v0.0.0.md", "# Working specification\n\ncanonical, immutable repository, embedded, emitted, and packaged v1.0.0 schema copies match byte-for-byte\n\nv1 release-candidate verification issue is complete\n")
	writeFile(t, repo, "docs/releases/v1.0.0.md", "# Cueson v1.0.0\n\nstable Cue JSON; unsigned and unattested\n\nFull changelog: https://github.com/shruggietech/cueson/blob/v1.0.0/CHANGELOG.md\n")
	writeFile(t, repo, "docs/formats/srt.md", "# SRT\n\n`srt-example`\n")
	writeFile(t, repo, "docs/formats/webvtt.md", "# WebVTT\n\n`vtt-example`\n")
	writeFile(t, repo, "testdata/conformance-matrix.json", `{"formats":[{"format":"srt","rows":[{"row_id":"srt-example"}]},{"format":"vtt","rows":[{"row_id":"vtt-example"}]}]}`)
	writeFile(t, repo, ".specify/memory/constitution.md", "# Constitution\n")
	for name, markers := range requiredReferenceMarkers {
		content := readFileForVerifierTest(t, repo, name)
		for _, marker := range markers {
			if !strings.Contains(content, marker) {
				content += "\n" + marker + "\n"
			}
		}
		writeFile(t, repo, name, content)
	}
	for name, suffix := range requiredSuffixes {
		content := readFileForVerifierTest(t, repo, name)
		if !strings.HasSuffix(content, suffix) {
			writeFile(t, repo, name, content+"\n"+suffix)
		}
	}
	return repo
}

func readFileForVerifierTest(t *testing.T, repo, name string) string {
	t.Helper()
	payload, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return string(payload)
}

func writeFile(t *testing.T, repo, name, content string) {
	t.Helper()
	path := filepath.Join(repo, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func assertViolation(t *testing.T, violations []violation, want string) {
	t.Helper()
	for _, got := range violationStrings(violations) {
		if got == want {
			return
		}
	}
	t.Fatalf("missing violation %q in %v", want, violations)
}
