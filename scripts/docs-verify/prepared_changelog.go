package main

import (
	"regexp"
	"strings"
	"time"
)

var (
	preparedReleaseHeading = regexp.MustCompile(`^## \[([^\]]+)\](?: - (.*))?$`)
	preparedCompareLink    = regexp.MustCompile(`^\[([^\]]+)\]:\s*(.*)$`)
)

type preparedSection struct {
	name string
}

func verifyPreparedChangelog(content string) []violation {
	const name = "CHANGELOG.md"
	var violations []violation
	var sections []preparedSection
	counts := make(map[string]int)
	linkCounts := make(map[string]int)
	links := map[string]string{
		"Unreleased": "https://github.com/shruggietech/cueson/compare/v1.1.0...HEAD",
		"1.1.0":      "https://github.com/shruggietech/cueson/compare/v1.0.0...v1.1.0",
		"1.0.0":      "https://github.com/shruggietech/cueson/compare/v0.0.0...v1.0.0",
	}
	categories := make(map[string]bool)
	current := ""
	fence := ""
	allowedCategories := map[string]bool{
		"Added": true, "Changed": true, "Deprecated": true, "Removed": true,
		"Fixed": true, "Security": true, "Decisions": true,
	}
	for index, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			marker := trimmed[:3]
			if fence == "" {
				fence = marker
			} else if fence == marker {
				fence = ""
			}
			continue
		}
		if fence != "" {
			continue
		}
		if match := preparedCompareLink.FindStringSubmatch(line); len(match) != 0 {
			if expected, ok := links[match[1]]; ok {
				linkCounts[match[1]]++
				if match[2] != expected {
					violations = append(violations, violation{path: name, line: index + 1, message: "prepared comparison link must use the exact planned tags: " + match[1]})
				}
			}
		}
		if strings.HasPrefix(line, "## ") {
			current = ""
			categories = make(map[string]bool)
			match := preparedReleaseHeading.FindStringSubmatch(line)
			if len(match) == 0 {
				if strings.HasPrefix(line, "## [Unreleased]") || strings.HasPrefix(line, "## [1.1.0]") {
					violations = append(violations, violation{path: name, line: index + 1, message: "malformed prepared release heading"})
				}
				continue
			}
			current = match[1]
			sections = append(sections, preparedSection{name: current})
			counts[current]++
			if current == "Unreleased" && match[2] != "" {
				violations = append(violations, violation{path: name, line: index + 1, message: "Unreleased must not have a release date"})
			}
			if current == "1.1.0" {
				if _, err := time.Parse("2006-01-02", match[2]); err != nil {
					violations = append(violations, violation{path: name, line: index + 1, message: "prepared 1.1.0 date must be a valid YYYY-MM-DD"})
				}
			}
			continue
		}
		if (current == "Unreleased" || current == "1.1.0") && strings.HasPrefix(line, "### ") {
			category := strings.TrimPrefix(line, "### ")
			if !allowedCategories[category] {
				violations = append(violations, violation{path: name, line: index + 1, message: "invalid prepared changelog category: " + category})
			} else if categories[category] {
				violations = append(violations, violation{path: name, line: index + 1, message: "duplicate prepared changelog category: " + category})
			}
			categories[category] = true
		}
	}
	for _, version := range []string{"Unreleased", "1.1.0", "1.0.0"} {
		if counts[version] != 1 {
			violations = append(violations, violation{path: name, message: "prepared history requires exactly one " + version + " section"})
		}
		if linkCounts[version] != 1 {
			violations = append(violations, violation{path: name, message: "prepared history requires exactly one comparison link: " + version})
		}
	}
	if len(sections) < 3 || sections[0].name != "Unreleased" || sections[1].name != "1.1.0" || sections[2].name != "1.0.0" {
		violations = append(violations, violation{path: name, message: "prepared sections must begin Unreleased, dated 1.1.0, then historical 1.0.0"})
	}
	return violations
}
