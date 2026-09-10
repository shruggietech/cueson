package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatPublicationNormalizesAndUnwraps(t *testing.T) {
	input := []byte("\ufeff# Title\r\n\r\nA paragraph that was\r\nwrapped.\r\n\r\n- A list item that was\r\n  wrapped.\r\n")
	want := "# Title\n\nA paragraph that was wrapped.\n\n- A list item that was wrapped.\n"
	got, err := formatPublication(input)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("formatPublication() = %q, want %q", got, want)
	}
}

func TestFormatPublicationRejectsInvalidUTF8AndMojibake(t *testing.T) {
	for _, input := range [][]byte{{0xff}, []byte("broken " + string(rune(0x00c3)) + " text")} {
		if _, err := formatPublication(input); err == nil {
			t.Fatalf("formatPublication(%q) returned nil error", input)
		}
	}
}

func TestFormatMarkdownPreservesLiteralStructures(t *testing.T) {
	input := "---\ntitle: Example\n---\n\nA | B\n--- | ---\n1 | 2\n\nSetext heading\n==============\n\n```text\nline one — literal\nline two\n```\n"
	if got := formatMarkdown(input); got != input {
		t.Fatalf("literal structures changed:\n%s", got)
	}
}

func TestFormatMarkdownReplacesEmDashesOnlyInProse(t *testing.T) {
	input := "Prose — changes; `code — stays`; [link](https://example.test/a_(b)—c); HTTPS://example.test/a—b.\n\n```json\n{\"value\": \"literal — stays\"}\n```\n"
	want := "Prose, changes; `code — stays`; [link](https://example.test/a_(b)—c); HTTPS://example.test/a—b.\n\n```json\n{\"value\": \"literal — stays\"}\n```\n"
	if got := formatMarkdown(input); got != want {
		t.Fatalf("formatMarkdown() = %q, want %q", got, want)
	}
}

func TestRepositoryCheckFindsBOMAndWrongLineEnding(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "bad.md"), []byte("\ufeff# Title\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	problems, err := checkRepository(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "UTF-8 BOM") || !strings.Contains(problems[0], "non-LF") {
		t.Fatalf("problems = %v", problems)
	}
}

func TestRepositoryFixDoesNotHideMojibake(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "bad.md")
	content := "broken " + string(rune(0x00c3)) + " text\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	problems, err := checkRepository(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "probable mojibake") {
		t.Fatalf("problems = %v", problems)
	}
}

func TestRepositoryChecksWindowsScripts(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"bad.ps1", "bad.cmd", "bad.bat"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("echo hello\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	problems, err := checkRepository(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 3 {
		t.Fatalf("problems = %v, want one for every Windows script", problems)
	}
}

func TestRepositorySkipsProtectedBrandTrees(t *testing.T) {
	root := t.TempDir()
	protected := []string{
		"brand/source.md",
		"docs/assets/brand/source.md",
		"docs/assets/favicons/source.md",
		"docs/assets/fonts/source.md",
	}
	content := []byte("\ufeff# Brand\r\n\r\nPreserve these bytes.\r\n")
	for _, name := range protected {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, content, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	authored := filepath.Join(root, "authored.md")
	if err := os.WriteFile(authored, []byte("broken "+string(rune(0x00c3))+" text\r\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, fix := range []bool{false, true} {
		problems, err := checkRepository(root, fix)
		if err != nil {
			t.Fatal(err)
		}
		if len(problems) != 1 || !strings.HasPrefix(problems[0], "authored.md:") {
			t.Fatalf("checkRepository(fix=%t) problems = %v, want only authored.md", fix, problems)
		}
		for _, name := range protected {
			path := filepath.Join(root, filepath.FromSlash(name))
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(content) {
				t.Fatalf("checkRepository(fix=%t) changed protected file %s: %q", fix, name, got)
			}
		}
	}
}

func TestRepositoryRejectsSymlinkedTextFile(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(target, []byte("outside — unchanged\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "linked.md")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	problems, err := checkRepository(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], "symlinked text file") {
		t.Fatalf("problems = %v", problems)
	}
	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "outside — unchanged\n" {
		t.Fatalf("external target changed: %q", data)
	}
}

func TestRepositorySourcesConform(t *testing.T) {
	problems, err := checkRepository("../..", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 0 {
		t.Fatalf("run go run ./scripts/github-format/main.go -fix; problems: %v", problems)
	}
}
