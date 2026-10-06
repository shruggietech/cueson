package main

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

const validV120History = "# Changelog\n\n## [Unreleased]\n\n## [1.2.0] - 2026-10-06\n\n### Added\n\n- Consumer attribution.\n\n### Decisions\n\n- 2026-10-06: Candidate only.\n\n## [1.1.0] - 2026-09-15\n\n### Historical category\n\n### Added\n\n### Added\n\n## [1.0.0] - 2026-09-11\n\n[Unreleased]: https://github.com/shruggietech/cueson/compare/v1.2.0...HEAD\n[1.2.0]: https://github.com/shruggietech/cueson/compare/v1.1.0...v1.2.0\n[1.1.0]: https://github.com/shruggietech/cueson/compare/v1.0.0...v1.1.0\n[1.0.0]: https://github.com/shruggietech/cueson/compare/v0.0.0...v1.0.0\n"

func TestV120PreparedHistory(t *testing.T) {
	for _, test := range []struct{ name, content, want string }{
		{"valid with frozen historical categories", validV120History, ""},
		{"missing candidate", strings.Replace(validV120History, "## [1.2.0] - 2026-10-06\n", "", 1), "exactly one 1.2.0 section"},
		{"invalid candidate date", strings.Replace(validV120History, "2026-10-06", "2026-02-29", 1), "valid YYYY-MM-DD"},
		{"duplicate candidate category", strings.Replace(validV120History, "### Decisions", "### Added", 1), "duplicate prepared changelog category"},
		{"invalid current category", strings.Replace(validV120History, "### Decisions", "### Magic", 1), "invalid prepared changelog category"},
		{"missing historical 1.1", strings.Replace(validV120History, "## [1.1.0] - 2026-09-15\n", "", 1), "exactly one 1.1.0 section"},
		{"stale future compare", strings.Replace(validV120History, "compare/v1.2.0...HEAD", "compare/v1.1.0...HEAD", 1), "exact planned tags: Unreleased"},
		{"wrong current compare", strings.Replace(validV120History, "compare/v1.1.0...v1.2.0", "compare/v1.0.0...v1.2.0", 1), "exact planned tags: 1.2.0"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := violationStrings(verifyPreparedChangelog(test.content))
			if test.want == "" && len(got) != 0 {
				t.Fatalf("unexpected violations: %v", got)
			}
			if test.want != "" && !strings.Contains(strings.Join(got, "\n"), test.want) {
				t.Fatalf("want %q, got %v", test.want, got)
			}
		})
	}
}

func TestRepositoryV120DecisionRemainsPreparationOnly(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", ".."))
	base := "specs/S034-prepare-v1-2-release-candidate/contracts/"
	decision := readFileForVerifierTest(t, repo, base+"publication-decision.md")
	if regexp.MustCompile(`\b[0-9a-f]{40}\b`).MatchString(decision) {
		t.Fatal("preparation must not predict or select a publication source revision")
	}
	for _, requirement := range []string{"Preparation only", "actual full lowercase 40-character post-merge", "same accepted bundle without rebuilding", "expiry", "published: false", "f2661a3d52effbab4a82a4d47197b5c7fae58496dc30a397ea3f2f668358b654", "positive byte length", "Create tag", "Push tag", "Publish GitHub Release", "Publish release assets", "stop all dependent publication/hosting"} {
		if !strings.Contains(decision, requirement) {
			t.Errorf("decision omits required source, inventory or authority binding %q", requirement)
		}
	}
	want := make(map[string]bool)
	for _, platform := range []string{"linux", "darwin", "windows"} {
		for _, arch := range []string{"amd64", "arm64"} {
			extension := ".tar.gz"
			if platform == "windows" {
				extension = ".zip"
			}
			archive := "cueson_1.2.0_" + platform + "_" + arch + extension
			want[archive], want[archive+".sbom.json"] = true, true
		}
	}
	want["cueson_1.2.0_checksums.txt"] = true
	for _, line := range strings.Split(decision, "\n") {
		if !strings.HasPrefix(line, "| `cueson_") {
			continue
		}
		fields := strings.Split(line, "|")
		if len(fields) != 5 {
			t.Fatalf("malformed public inventory row %q", line)
		}
		name := strings.Trim(strings.TrimSpace(fields[1]), "`")
		if !want[name] {
			t.Fatalf("unexpected or duplicate public file %q", name)
		}
		delete(want, name)
	}
	if len(want) != 0 {
		t.Fatalf("public inventory omits %v", want)
	}
	notes := readFileForVerifierTest(t, repo, base+"release-notes.md")
	if !utf8.ValidString(notes) || strings.HasPrefix(notes, "\ufeff") || strings.Contains(notes, "\r") || !strings.HasPrefix(notes, "# Cueson v1.2.0\n\n") || !strings.HasSuffix(notes, "Full changelog: https://github.com/shruggietech/cueson/blob/v1.2.0/CHANGELOG.md\n") {
		t.Fatal("public notes must retain exact title, UTF-8/LF and tagged changelog suffix")
	}
	if strings.Contains(notes, "**Status:**") || strings.Contains(notes, "is now published") {
		t.Fatal("public notes must not carry a candidate banner or premature availability claim")
	}
}

func TestV120NotesRequireCandidateAndExactSuffix(t *testing.T) {
	repo := newRepository(t)
	writeFile(t, repo, "docs/releases/v1.2.0.md", "# Cueson v1.2.0\n\nCueson v1.2.0 is now published.\n\nFull changelog: https://github.com/shruggietech/cueson/blob/main/CHANGELOG.md\n")
	result, err := verifyRepository(repo)
	if err != nil {
		t.Fatal(err)
	}
	assertViolation(t, result.violations, "docs/releases/v1.2.0.md: required final suffix is missing: Full changelog: https://github.com/shruggietech/cueson/blob/v1.2.0/CHANGELOG.md")
	assertViolation(t, result.violations, "docs/releases/v1.2.0.md: stale capability or release claim remains: v1.2.0 is now published")
}
