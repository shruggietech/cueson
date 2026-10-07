package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestS035ValidationGuideHasRunnableCommandsAndHandoff(t *testing.T) {
	repo, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	base := "specs/S035-publish-v1-2-public-schema-site/"
	guide := readFileForVerifierTest(t, repo, base+"quickstart.md")
	start := strings.Index(guide, "```text\n")
	if start < 0 {
		t.Fatal("validation guide requires a text command fence")
	}
	commands := guide[start+len("```text\n"):]
	end := strings.Index(commands, "\n```")
	if end < 0 {
		t.Fatal("validation command fence is unterminated")
	}
	want := []string{
		"go -C scripts/docs-verify test -count=1 ./...",
		"go -C scripts/docs-verify run . -repo ../..",
		"go run ./scripts/github-format/main.go",
		"cd site",
		"pnpm lint", "pnpm generate", "pnpm generate:check", "pnpm test:unit",
		"pnpm build", "pnpm test:browser", "pnpm verify:artifact", "pnpm deploy:dry-run",
	}
	if got := commands[:end]; got != strings.Join(want, "\n") {
		t.Fatalf("validation guide commands differ from the executable workflow:\n%s", got)
	}
	if !strings.Contains(guide, "[production-handoff.md](contracts/production-handoff.md)") {
		t.Fatal("validation guide must link its production contract")
	}
	if _, err := os.Stat(filepath.Join(repo, filepath.FromSlash(base+"contracts/production-handoff.md"))); err != nil {
		t.Fatal(err)
	}
}

func TestS035GeneratedSourcePagesRejectStalePublicationClaims(t *testing.T) {
	for _, tc := range []struct{ path, claim string }{
		{"docs/conversion.md", "S034 prepares this behavior in the exact 1.2.0 candidate without changing published downloads"},
		{"docs/Cueson-Project-Specification-v0.0.0.md", "public availability remains v1.1.0"},
		{"docs/formats/srt.md", "New encode output always uses exact 1.1.0"},
		{"docs/formats/webvtt.md", "New encode output always uses exact 1.1.0"},
		{"docs/formats/ass-ssa.md", "exact current `1.1.0` identity"},
		{"docs/formats/ass-ssa.md", "#67 governs production activation and live verification"},
		{"docs/cueson-media-format-guide.html", "Current v1.1.0 release boundary:"},
		{"docs/cueson-media-format-guide.html", "issue #67 governs exact-main production activation"},
		{"docs/cueson-media-format-guide.html", "Stable 1.1.0"},
		{"docs/release-process.md", "Current public software availability is v1.1.0"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			repo := newRepository(t)
			content := readFileForVerifierTest(t, repo, tc.path)
			writeFile(t, repo, tc.path, content+"\n"+tc.claim+"\n")
			result, err := verifyRepository(repo)
			if err != nil {
				t.Fatal(err)
			}
			assertViolation(t, result.violations, tc.path+": stale capability or release claim remains: "+tc.claim)
		})
	}
}
