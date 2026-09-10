package main

import (
	"bufio"
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maximumDocumentSize = 4 << 20

var requiredDocuments = []string{
	"README.md",
	"CONTRIBUTING.md",
	"SECURITY.md",
	"AGENTS.md",
	"CHANGELOG.md",
	"docs/architecture.md",
	"docs/schema.md",
	"docs/cli.md",
	"docs/formats/srt.md",
	"docs/formats/webvtt.md",
	"docs/project-management.md",
	"docs/release-process.md",
}

var (
	markdownHeadingPattern = regexp.MustCompile(`^\s{0,3}(#{1,6})\s+(.+?)\s*$`)
	setextHeadingPattern   = regexp.MustCompile(`^\s{0,3}(?:=+|-+)\s*$`)
	referenceDefinition    = regexp.MustCompile(`^\s{0,3}\[([^]]+)\]:\s*(.*)$`)
	referenceUse           = regexp.MustCompile(`\[([^]]+)\]\[([^]]*)\]`)
	shortcutReference      = regexp.MustCompile(`\[([^][\r\n]+)\]`)
	autolinkPattern        = regexp.MustCompile(`<((?:https?://|mailto:)[^ <>]+)>`)
	htmlAttributePattern   = regexp.MustCompile(`(?i)\b(?:href|src)\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s>]+))`)
	htmlAnchorPattern      = regexp.MustCompile(`(?i)\b(?:id|name)\s*=\s*(?:"([^"]+)"|'([^']+)'|([^\s>]+))`)
	htmlTagPattern         = regexp.MustCompile(`<[^>]*>`)
	markdownMarkupPattern  = regexp.MustCompile("[`*_~]")
	schemePattern          = regexp.MustCompile(`^([A-Za-z][A-Za-z0-9+.-]*):`)
)

type violation struct {
	path    string
	line    int
	message string
}

type verificationResult struct {
	documents  int
	localLinks int
	violations []violation
}

type repositoryEntry struct {
	mode os.FileMode
}

type documentReference struct {
	line   int
	target string
}

func verifyRepository(repo string) (verificationResult, error) {
	root, err := filepath.Abs(repo)
	if err != nil {
		return verificationResult{}, fmt.Errorf("resolve repository root: %w", err)
	}
	info, err := os.Lstat(root)
	if err != nil {
		return verificationResult{}, fmt.Errorf("read repository root: %w", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return verificationResult{}, fmt.Errorf("repository root is not a regular directory")
	}

	index, folded, err := repositoryIndex(root)
	if err != nil {
		return verificationResult{}, err
	}
	result := verificationResult{}
	for _, required := range requiredDocuments {
		entry, ok := index[required]
		if !ok {
			result.violations = append(result.violations, violation{path: required, message: "required document is missing"})
			continue
		}
		if !entry.mode.IsRegular() {
			result.violations = append(result.violations, violation{path: required, message: "required document is not a regular file"})
		}
	}

	documents := auditedDocuments(index)
	result.documents = len(documents)
	anchorCache := make(map[string]map[string]struct{}, len(documents))
	for _, name := range documents {
		content, readErr := readDocument(root, name)
		if readErr != nil {
			result.violations = append(result.violations, violation{path: name, message: readErr.Error()})
			continue
		}
		if isAnchorDocument(name) {
			anchorCache[name] = documentAnchors(name, content)
		}
		for _, ref := range documentReferences(name, content) {
			local, issue := resolveReference(name, ref, index, folded, anchorCache, root)
			if local {
				result.localLinks++
			}
			if issue != "" {
				result.violations = append(result.violations, violation{path: name, line: ref.line, message: issue})
			}
		}
	}
	sort.Slice(result.violations, func(i, j int) bool {
		a, b := result.violations[i], result.violations[j]
		if a.path != b.path {
			return a.path < b.path
		}
		if a.line != b.line {
			return a.line < b.line
		}
		return a.message < b.message
	})
	return result, nil
}

func repositoryIndex(root string) (map[string]repositoryEntry, map[string][]string, error) {
	index := make(map[string]repositoryEntry)
	folded := make(map[string][]string)
	err := filepath.WalkDir(root, func(name string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if entry.IsDir() && (rel == ".git" || strings.HasPrefix(rel, ".git/")) {
			return filepath.SkipDir
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		index[rel] = repositoryEntry{mode: info.Mode()}
		key := strings.ToLower(rel)
		folded[key] = append(folded[key], rel)
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("index repository: %w", err)
	}
	for key := range folded {
		sort.Strings(folded[key])
	}
	return index, folded, nil
}

func auditedDocuments(index map[string]repositoryEntry) []string {
	rootDocuments := map[string]bool{
		"README.md": true, "CHANGELOG.md": true, "CONTRIBUTING.md": true, "SECURITY.md": true, "AGENTS.md": true,
		".specify/memory/constitution.md": true,
	}
	var names []string
	for name, entry := range index {
		if !entry.mode.IsRegular() {
			continue
		}
		extension := strings.ToLower(path.Ext(name))
		if rootDocuments[name] || (strings.HasPrefix(name, "docs/") && (extension == ".md" || extension == ".html" || extension == ".htm")) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func readDocument(root, name string) (string, error) {
	full := filepath.Join(root, filepath.FromSlash(name))
	info, err := os.Stat(full)
	if err != nil {
		return "", fmt.Errorf("read document: %w", err)
	}
	if info.Size() > maximumDocumentSize {
		return "", fmt.Errorf("document exceeds %d bytes", maximumDocumentSize)
	}
	data, err := os.ReadFile(full)
	if err != nil {
		return "", fmt.Errorf("read document: %w", err)
	}
	if !utf8.Valid(data) {
		return "", fmt.Errorf("document is not valid UTF-8")
	}
	return string(data), nil
}

func documentReferences(name, content string) []documentReference {
	markdown := strings.HasSuffix(strings.ToLower(name), ".md")
	lines := sanitizedLines(content, markdown, markdown)
	definitions := make(map[string]documentReference)
	for index, line := range lines {
		if match := referenceDefinition.FindStringSubmatch(line); match != nil {
			definitions[normalizeReference(match[1])] = documentReference{line: index + 1, target: linkDestination(match[2])}
		}
	}
	var references []documentReference
	for index, line := range lines {
		lineNumber := index + 1
		if match := referenceDefinition.FindStringSubmatch(line); match != nil {
			references = append(references, documentReference{line: lineNumber, target: linkDestination(match[2])})
			continue
		}
		for _, target := range inlineMarkdownDestinations(line) {
			references = append(references, documentReference{line: lineNumber, target: target})
		}
		for _, match := range referenceUse.FindAllStringSubmatch(line, -1) {
			label := match[2]
			if label == "" {
				label = match[1]
			}
			if definition, ok := definitions[normalizeReference(label)]; ok {
				references = append(references, documentReference{line: lineNumber, target: definition.target})
			}
		}
		for _, match := range shortcutReference.FindAllStringSubmatchIndex(line, -1) {
			start, end := match[0], match[1]
			if (end < len(line) && (line[end] == '(' || line[end] == '[')) || (start > 0 && line[start-1] == ']') {
				continue
			}
			label := normalizeReference(line[match[2]:match[3]])
			if definition, ok := definitions[label]; ok {
				references = append(references, documentReference{line: lineNumber, target: definition.target})
			}
		}
		for _, match := range autolinkPattern.FindAllStringSubmatch(line, -1) {
			references = append(references, documentReference{line: lineNumber, target: match[1]})
		}
		for _, match := range htmlAttributePattern.FindAllStringSubmatch(line, -1) {
			references = append(references, documentReference{line: lineNumber, target: firstNonempty(match[1:]...)})
		}
	}
	return references
}

func sanitizedLines(content string, markdown, removeInlineCode bool) []string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 4096), maximumDocumentSize)
	var lines []string
	inComment := false
	inFence := false
	fenceMarker := ""
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if markdown && !inComment && (strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")) {
			marker := trimmed[:3]
			if !inFence {
				inFence, fenceMarker = true, marker
			} else if marker == fenceMarker {
				inFence, fenceMarker = false, ""
			}
			lines = append(lines, "")
			continue
		}
		if inFence {
			lines = append(lines, "")
			continue
		}
		line = stripHTMLComments(line, &inComment)
		if removeInlineCode {
			line = stripInlineCode(line)
		}
		lines = append(lines, line)
	}
	return lines
}

func stripHTMLComments(line string, inComment *bool) string {
	var output strings.Builder
	for len(line) > 0 {
		if *inComment {
			end := strings.Index(line, "-->")
			if end < 0 {
				return output.String()
			}
			line = line[end+3:]
			*inComment = false
			continue
		}
		start := strings.Index(line, "<!--")
		if start < 0 {
			output.WriteString(line)
			break
		}
		output.WriteString(line[:start])
		line = line[start+4:]
		*inComment = true
	}
	return output.String()
}

func stripInlineCode(line string) string {
	var output strings.Builder
	for index := 0; index < len(line); {
		if line[index] != '`' {
			output.WriteByte(line[index])
			index++
			continue
		}
		run := 1
		for index+run < len(line) && line[index+run] == '`' {
			run++
		}
		closing := strings.Index(line[index+run:], strings.Repeat("`", run))
		if closing < 0 {
			output.WriteString(line[index:])
			break
		}
		index += run + closing + run
	}
	return output.String()
}

func inlineMarkdownDestinations(line string) []string {
	var destinations []string
	for search := 0; search < len(line); {
		offset := strings.Index(line[search:], "](")
		if offset < 0 {
			break
		}
		start := search + offset + 2
		depth := 1
		end := start
		for ; end < len(line); end++ {
			switch line[end] {
			case '\\':
				end++
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					destinations = append(destinations, linkDestination(line[start:end]))
					end++
					break
				}
			}
			if depth == 0 {
				break
			}
		}
		if depth != 0 {
			destinations = append(destinations, strings.TrimSpace(line[start:]))
			break
		}
		search = end
	}
	return destinations
}

func linkDestination(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "<") {
		if end := strings.Index(value, ">"); end >= 0 {
			return html.UnescapeString(value[1:end])
		}
	}
	for index, r := range value {
		if unicode.IsSpace(r) {
			return html.UnescapeString(value[:index])
		}
	}
	return html.UnescapeString(value)
}

func normalizeReference(value string) string {
	return strings.Join(strings.Fields(strings.ToLower(value)), " ")
}

func resolveReference(source string, ref documentReference, index map[string]repositoryEntry, folded map[string][]string, anchors map[string]map[string]struct{}, root string) (bool, string) {
	raw := strings.TrimSpace(ref.target)
	if raw == "" {
		return true, "empty link target"
	}
	if strings.Contains(raw, "\\") {
		return true, "local target contains a backslash: " + raw
	}
	for _, r := range raw {
		if unicode.IsControl(r) {
			return true, "link target contains a control character"
		}
	}
	if strings.HasPrefix(raw, "//") {
		parsed, err := url.Parse(raw)
		if err != nil || parsed.Host == "" {
			return false, "malformed external URL: " + raw
		}
		return false, ""
	}
	if match := schemePattern.FindStringSubmatch(raw); match != nil {
		scheme := strings.ToLower(match[1])
		if scheme != "http" && scheme != "https" && scheme != "mailto" {
			return false, "unsupported URL scheme: " + scheme
		}
		parsed, err := url.Parse(raw)
		if err != nil || ((scheme == "http" || scheme == "https") && parsed.Host == "") {
			return false, "malformed external URL: " + raw
		}
		return false, ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return true, "malformed URL escape: " + raw
	}
	decodedPath, err := url.PathUnescape(parsed.EscapedPath())
	if err != nil {
		return true, "malformed URL escape: " + raw
	}
	fragment, err := url.PathUnescape(parsed.Fragment)
	if err != nil {
		return true, "malformed URL escape: " + raw
	}
	if decodedPath == "" && fragment == "" {
		return true, "empty link target"
	}

	target := source
	if decodedPath != "" {
		if strings.HasPrefix(decodedPath, "/") {
			target = path.Clean(strings.TrimPrefix(decodedPath, "/"))
		} else {
			target = path.Clean(path.Join(path.Dir(source), decodedPath))
		}
		if target == ".." || strings.HasPrefix(target, "../") || path.IsAbs(target) {
			return true, "local target escapes repository: " + decodedPath
		}
	}
	entry, ok := index[target]
	if !ok {
		if alternatives := folded[strings.ToLower(target)]; len(alternatives) != 0 {
			return true, "local target has incorrect case: " + decodedPath
		}
		return true, "local target does not exist: " + decodedPath
	}
	if entry.mode&os.ModeSymlink != 0 {
		return true, "local target is a symbolic link: " + decodedPath
	}
	if fragment == "" {
		return true, ""
	}
	if !isAnchorDocument(target) || entry.mode.IsDir() {
		return true, "fragment is not supported for target: " + decodedPath
	}
	targetAnchors, ok := anchors[target]
	if !ok {
		content, readErr := readDocument(root, target)
		if readErr != nil {
			return true, readErr.Error()
		}
		targetAnchors = documentAnchors(target, content)
		anchors[target] = targetAnchors
	}
	if _, ok := targetAnchors[fragment]; !ok {
		return true, "fragment does not exist: " + fragment
	}
	return true, ""
}

func isAnchorDocument(name string) bool {
	extension := strings.ToLower(path.Ext(name))
	return extension == ".md" || extension == ".markdown" || extension == ".html" || extension == ".htm"
}

func documentAnchors(name, content string) map[string]struct{} {
	anchors := make(map[string]struct{})
	if strings.HasSuffix(strings.ToLower(name), ".md") || strings.HasSuffix(strings.ToLower(name), ".markdown") {
		for anchor := range markdownAnchors(content) {
			anchors[anchor] = struct{}{}
		}
	}
	for _, line := range sanitizedLines(content, false, false) {
		for _, match := range htmlAnchorPattern.FindAllStringSubmatch(line, -1) {
			anchors[html.UnescapeString(firstNonempty(match[1:]...))] = struct{}{}
		}
	}
	return anchors
}

func markdownAnchors(content string) map[string]struct{} {
	lines := sanitizedLines(content, true, false)
	anchors := make(map[string]struct{})
	counts := make(map[string]int)
	add := func(text string) {
		base := githubSlug(text)
		if base == "" {
			return
		}
		anchor := base
		if count := counts[base]; count != 0 {
			anchor = fmt.Sprintf("%s-%d", base, count)
		}
		counts[base]++
		anchors[anchor] = struct{}{}
	}
	for index, line := range lines {
		if match := markdownHeadingPattern.FindStringSubmatch(line); match != nil {
			text := strings.TrimSpace(match[2])
			text = strings.TrimSpace(strings.TrimRight(text, "#"))
			add(text)
			continue
		}
		if index > 0 && setextHeadingPattern.MatchString(line) && strings.TrimSpace(lines[index-1]) != "" {
			add(strings.TrimSpace(lines[index-1]))
		}
	}
	return anchors
}

func githubSlug(value string) string {
	value = html.UnescapeString(htmlTagPattern.ReplaceAllString(value, ""))
	value = markdownMarkupPattern.ReplaceAllString(value, "")
	var slug strings.Builder
	for _, r := range strings.ToLower(value) {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '-', r == '_':
			slug.WriteRune(r)
		case unicode.IsSpace(r):
			slug.WriteByte('-')
		}
	}
	return slug.String()
}

func firstNonempty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func violationStrings(violations []violation) []string {
	items := make([]string, len(violations))
	for index, item := range violations {
		if item.line > 0 {
			items[index] = fmt.Sprintf("%s:%d: %s", item.path, item.line, item.message)
		} else {
			items[index] = fmt.Sprintf("%s: %s", item.path, item.message)
		}
	}
	return items
}
