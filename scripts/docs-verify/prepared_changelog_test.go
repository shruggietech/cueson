package main

import (
	"path/filepath"
	"strings"
	"testing"
)

const validPreparedHistory = "# Changelog\n\n## [Unreleased]\n\n## [1.1.0] - 2026-09-15\n\n### Added\n\n- Complete prepared history.\n\n### Decisions\n\n- 2026-09-15: Preserve publication authority.\n\n## [1.0.0] - 2026-09-11\n\n### Legacy category\n\n### Added\n\n### Added\n\nHistorical bytes are outside the new category gate.\n\n[Unreleased]: https://github.com/shruggietech/cueson/compare/v1.1.0...HEAD\n[1.1.0]: https://github.com/shruggietech/cueson/compare/v1.0.0...v1.1.0\n[1.0.0]: https://github.com/shruggietech/cueson/compare/v0.0.0...v1.0.0\n"

func TestPreparedChangelogMetadata(t *testing.T) {
	for _, test := range []struct {
		name    string
		content string
		want    string
	}{
		{"valid scoped history", validPreparedHistory, ""},
		{"valid leap date", strings.Replace(validPreparedHistory, "2026-09-15", "2024-02-29", 1), ""},
		{"missing current", strings.Replace(validPreparedHistory, "## [1.1.0] - 2026-09-15\n", "", 1), "exactly one 1.1.0 section"},
		{"duplicate current", strings.Replace(validPreparedHistory, "## [1.0.0]", "## [1.1.0] - 2026-09-15\n\n## [1.0.0]", 1), "exactly one 1.1.0 section"},
		{"missing Unreleased", strings.Replace(validPreparedHistory, "## [Unreleased]\n", "", 1), "exactly one Unreleased section"},
		{"duplicate Unreleased", strings.Replace(validPreparedHistory, "## [Unreleased]", "## [Unreleased]\n\n## [Unreleased]", 1), "exactly one Unreleased section"},
		{"date absent", strings.Replace(validPreparedHistory, "## [1.1.0] - 2026-09-15", "## [1.1.0]", 1), "valid YYYY-MM-DD"},
		{"invalid calendar date", strings.Replace(validPreparedHistory, "2026-09-15", "2026-02-29", 1), "valid YYYY-MM-DD"},
		{"noncanonical date", strings.Replace(validPreparedHistory, "2026-09-15", "2026-9-15", 1), "valid YYYY-MM-DD"},
		{"date trailing text", strings.Replace(validPreparedHistory, "2026-09-15\n", "2026-09-15 published\n", 1), "valid YYYY-MM-DD"},
		{"dated Unreleased", strings.Replace(validPreparedHistory, "## [Unreleased]", "## [Unreleased] - 2026-09-15", 1), "Unreleased must not have a release date"},
		{"reordered current", strings.ReplaceAll(strings.ReplaceAll(validPreparedHistory, "## [Unreleased]", "## [SWAP]"), "## [1.1.0] - 2026-09-15", "## [Unreleased]"), "must begin Unreleased"},
		{"invalid current category", strings.Replace(validPreparedHistory, "### Decisions", "### Features", 1), "invalid prepared changelog category"},
		{"duplicate current category", strings.Replace(validPreparedHistory, "### Decisions", "### Added", 1), "duplicate prepared changelog category"},
		{"duplicate future category", strings.Replace(validPreparedHistory, "## [Unreleased]\n", "## [Unreleased]\n\n### Fixed\n\n### Fixed\n", 1), "duplicate prepared changelog category"},
		{"future tag link", strings.Replace(validPreparedHistory, "compare/v1.0.0...v1.1.0", "compare/v1.0.0...v1.2.0", 1), "exact planned tags: 1.1.0"},
		{"duplicate comparison", validPreparedHistory + "[1.1.0]: https://github.com/shruggietech/cueson/compare/v1.0.0...v1.1.0\n", "exactly one comparison link: 1.1.0"},
		{"fenced fake headings", strings.Replace(validPreparedHistory, "- Complete prepared history.", "- Complete prepared history.\n\n```text\n## [1.1.0] - 2026-01-01\n### Added\n### Bogus\n```", 1), ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			violations := violationStrings(verifyPreparedChangelog(test.content))
			if test.want == "" {
				if len(violations) != 0 {
					t.Fatalf("unexpected violations: %v", violations)
				}
				return
			}
			if !strings.Contains(strings.Join(violations, "\n"), test.want) {
				t.Fatalf("want %q, got %v", test.want, violations)
			}
		})
	}
}

func TestRepositoryChangelogIsPreparedForV110(t *testing.T) {
	repo := filepath.Clean(filepath.Join("..", ".."))
	content := readFileForVerifierTest(t, repo, "CHANGELOG.md")
	if got := strings.Count(content, "\n## [1.1.0] - "); got != 1 {
		t.Fatalf("prepared repository history has %d dated 1.1.0 headings, want exactly one before tag authorization", got)
	}
	if violations := verifyReleaseChangelog(repo); len(violations) != 0 {
		t.Fatalf("prepared repository metadata: %v", violationStrings(violations))
	}
}

func TestPreparedChangelogMarkdownEquivalence(t *testing.T) {
	for _, test := range []struct {
		name    string
		content string
		want    string
	}{
		{"indented valid headings", strings.ReplaceAll(validPreparedHistory, "\n##", "\n   ##"), ""},
		{"indented duplicate release", strings.Replace(validPreparedHistory, "## [1.0.0]", "   ## [1.1.0] - 2026-09-15\n\n## [1.0.0]", 1), "exactly one 1.1.0 section"},
		{"indented duplicate category", strings.Replace(validPreparedHistory, "### Decisions", "  ### Added", 1), "duplicate prepared changelog category"},
		{"indented invalid category", strings.Replace(validPreparedHistory, "### Decisions", " ### Features", 1), "invalid prepared changelog category"},
		{"indented reordered release", strings.Replace(validPreparedHistory, "## [Unreleased]", "   ## [1.0.0] - 2026-09-11\n\n## [Unreleased]", 1), "must begin Unreleased"},
		{"closing ATX hashes", strings.ReplaceAll(strings.ReplaceAll(validPreparedHistory, "## [1.1.0] - 2026-09-15", "  ## [1.1.0] - 2026-09-15 ###"), "### Decisions", " ### Decisions ###"), ""},
		{"closing-hash duplicate category", strings.Replace(validPreparedHistory, "### Decisions", "   ### Added ###", 1), "duplicate prepared changelog category"},
		{"tab-separated ATX", strings.ReplaceAll(validPreparedHistory, "## ", "##\t"), ""},
		{"four-space code", strings.Replace(validPreparedHistory, "- Complete prepared history.", "- Complete prepared history.\n\n    ## [1.1.0] - 2026-09-15\n    ### Added\n    ### Bogus", 1), ""},
		{"equivalent reference labels", strings.Replace(validPreparedHistory, "[Unreleased]:", "   [  UnReLeAsEd\t ]:", 1), ""},
		{"case-insensitive duplicate reference", validPreparedHistory + "[unreleased]: https://github.com/shruggietech/cueson/compare/v1.1.0...HEAD\n", "exactly one comparison link: Unreleased"},
		{"indented duplicate reference", validPreparedHistory + "   [  Unreleased  ]: https://github.com/shruggietech/cueson/compare/v1.1.0...HEAD\n", "exactly one comparison link: Unreleased"},
		{"earlier wrong case target", strings.Replace(validPreparedHistory, "[Unreleased]:", "[UNRELEASED]: https://github.com/shruggietech/cueson/compare/v1.0.0...HEAD\n[Unreleased]:", 1), "exact planned tags: Unreleased"},
		{"earlier wrong indented target", strings.Replace(validPreparedHistory, "[1.1.0]:", "  [ 1.1.0 ]: https://github.com/shruggietech/cueson/compare/v1.0.0...v1.2.0\n[1.1.0]:", 1), "exact planned tags: 1.1.0"},
		{"four-space reference code", validPreparedHistory + "    [unreleased]: https://github.com/shruggietech/cueson/compare/v1.0.0...HEAD\n", ""},
		{"leading-tab code", validPreparedHistory + "\n\t## [1.1.0] - 2026-09-15\n\t[unreleased]: https://github.com/shruggietech/cueson/compare/v1.0.0...HEAD\n", ""},
		{"multiline equivalent label", strings.Replace(validPreparedHistory, "[Unreleased]:", "[\n UnReLeAsEd\n ]:", 1), ""},
		{"multiline tab whitespace label", strings.Replace(validPreparedHistory, "[Unreleased]:", "[\n\tUnReLeAsEd\n ]:", 1), ""},
		{"multiline earlier wrong target", strings.Replace(validPreparedHistory, "[Unreleased]:", "[\n Unreleased\n ]: https://github.com/shruggietech/cueson/compare/v1.0.0...HEAD\n[Unreleased]:", 1), "exact planned tags: Unreleased"},
		{"unfinished label cannot hide heading", strings.Replace(validPreparedHistory, "## [1.1.0] - 2026-09-15", "[\n   ## [1.1.0] - 2026-09-15", 1), ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			violations := violationStrings(verifyPreparedChangelog(test.content))
			if test.want == "" {
				if len(violations) != 0 {
					t.Fatalf("unexpected violations: %v", violations)
				}
				return
			}
			if !strings.Contains(strings.Join(violations, "\n"), test.want) {
				t.Fatalf("want %q, got %v", test.want, violations)
			}
		})
	}
}
