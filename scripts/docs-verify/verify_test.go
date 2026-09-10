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
	if got := stdout.String(); !strings.Contains(got, "checked 14 documents") || !strings.Contains(got, "0 local links") {
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
	writeFile(t, repo, ".specify/memory/constitution.md", "# Constitution\n")
	return repo
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
