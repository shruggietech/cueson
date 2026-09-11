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
	if got := stdout.String(); !strings.Contains(got, "checked 19 documents") || !strings.Contains(got, "0 local links") || !strings.Contains(got, "10 registered examples") || !strings.Contains(got, "2 format rows") {
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
	for _, name := range []string{"docs/releases/v0.0.0.md", "docs/releases/v1.0.0.md"} {
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

func TestV1ChangelogRequiresUnreleasedBeforeReleaseSection(t *testing.T) {
	repo := newRepository(t)
	changelog := strings.Replace(readFileForVerifierTest(t, repo, "CHANGELOG.md"), "## [Unreleased]\n\n### Added\n\n- Post-release work.\n\n", "", 1)
	writeFile(t, repo, "CHANGELOG.md", changelog)
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "CHANGELOG.md: Unreleased must precede the dated 1.0.0 section")
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
	writeFile(t, repo, "CHANGELOG.md", "# Changelog\n\n## [Unreleased]\n\n### Added\n\n- Post-release work.\n\n## [1.0.0] - 2026-09-11\n\n[Unreleased]: https://github.com/shruggietech/cueson/compare/v1.0.0...HEAD\n[1.0.0]: https://github.com/shruggietech/cueson/compare/v0.0.0...v1.0.0\n")
	writeFile(t, repo, "README.md", readme.String()+"\nv1.0.0 released and independently verified; v0.0.0 remains historical; v1.0.0 GitHub Release.\n")
	writeFile(t, repo, "docs/schema.md", "# Schema\n\n$id schema_version format_support format_data source v0.0.0 v1.0.0 non-normative stable schema released and independently verified\n")
	writeFile(t, repo, "docs/compatibility.md", "# Compatibility\n\nCLI Cue JSON Schema internal/ v0.0.0 v1.0.0 Windows macOS Linux production stable release published and independently verified\n")
	writeFile(t, repo, "docs/Cueson-Project-Specification-v0.0.0.md", "# Working specification\n\ncanonical, immutable repository, embedded, emitted, and packaged v1.0.0 schema copies match byte-for-byte\n\nv1 release-candidate verification issue is complete\n")
	writeFile(t, repo, "docs/releases/v1.0.0.md", "# Cueson v1.0.0\n\nstable Cue JSON; unsigned and unattested\n\nFull changelog: https://github.com/shruggietech/cueson/blob/v1.0.0/CHANGELOG.md\n")
	writeFile(t, repo, "docs/formats/srt.md", "# SRT\n\n`srt-example`\n")
	writeFile(t, repo, "docs/formats/webvtt.md", "# WebVTT\n\n`vtt-example`\n")
	writeFile(t, repo, "testdata/conformance-matrix.json", `{"formats":[{"format":"srt","rows":[{"row_id":"srt-example"}]},{"format":"vtt","rows":[{"row_id":"vtt-example"}]}]}`)
	writeFile(t, repo, ".specify/memory/constitution.md", "# Constitution\n")
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
