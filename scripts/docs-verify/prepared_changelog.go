package main

import (
	"regexp"
	"strings"
	"time"
)

var (
	preparedReleaseHeading = regexp.MustCompile(`^\[([^\]]+)\](?:[ \t]+-[ \t]+(.*))?$`)
	preparedATXHeading     = regexp.MustCompile(`^(#{1,6})(?:[ \t]+(.*)|[ \t]*)$`)
	preparedClosingHashes  = regexp.MustCompile(`[ \t]+#+[ \t]*$`)
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
	versions := map[string]string{"unreleased": "Unreleased", "1.1.0": "1.1.0", "1.0.0": "1.0.0"}
	categories := make(map[string]bool)
	current := ""
	fence := ""
	referencePrefix := ""
	allowedCategories := map[string]bool{
		"Added": true, "Changed": true, "Deprecated": true, "Removed": true,
		"Fixed": true, "Security": true, "Decisions": true,
	}
	for index, line := range strings.Split(content, "\n") {
		// CommonMark allows at most three spaces before block headings and
		// reference definitions. Four spaces or a leading tab are code.
		indent := len(line) - len(strings.TrimLeft(line, " "))
		isCode := indent > 3 || strings.HasPrefix(line[indent:], "\t")
		if isCode && referencePrefix == "" {
			continue
		}
		if !isCode {
			line = line[indent:]
		}
		trimmed := strings.TrimSpace(line)
		if !isCode && (strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")) {
			referencePrefix = ""
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
		reference := line
		if referencePrefix != "" {
			if trimmed != "" && !strings.Contains(line, "[") {
				reference = referencePrefix + "\n" + line
			}
			referencePrefix = ""
		}
		// Reference labels can span nonblank lines and collapse whitespace.
		// Keep the CommonMark 999-character ceiling and process headings
		// independently, so an unfinished label cannot hide metadata.
		if strings.HasPrefix(reference, "[") && !strings.Contains(reference, "]") && len(reference) <= 1000 {
			referencePrefix = reference
		}
		if match := preparedCompareLink.FindStringSubmatch(reference); len(match) != 0 && len(match[1]) <= 999 {
			if version, ok := versions[normalizeReference(match[1])]; ok {
				expected := links[version]
				linkCounts[version]++
				if match[2] != expected {
					violations = append(violations, violation{path: name, line: index + 1, message: "prepared comparison link must use the exact planned tags: " + version})
				}
			}
		}
		if isCode {
			continue
		}
		heading := preparedATXHeading.FindStringSubmatch(line)
		if len(heading) == 0 {
			continue
		}
		text := strings.TrimSpace(preparedClosingHashes.ReplaceAllString(heading[2], ""))
		if len(heading[1]) == 2 {
			current = ""
			categories = make(map[string]bool)
			match := preparedReleaseHeading.FindStringSubmatch(text)
			if len(match) == 0 {
				if strings.HasPrefix(text, "[Unreleased]") || strings.HasPrefix(text, "[1.1.0]") {
					violations = append(violations, violation{path: name, line: index + 1, message: "malformed prepared release heading"})
				}
				continue
			}
			current = match[1]
			if version, ok := versions[normalizeReference(current)]; ok {
				current = version
			}
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
		if (current == "Unreleased" || current == "1.1.0") && len(heading[1]) == 3 {
			category := text
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
