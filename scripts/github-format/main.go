// Command github-format checks or repairs repository-authored text and formats Markdown for GitHub publication.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	alertMarker     = regexp.MustCompile(`^\[!(?:NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]$`)
	headingLine     = regexp.MustCompile(`^#{1,6}(?:\s|$)`)
	linkDefinition  = regexp.MustCompile(`^\[[^]]+\]:\s*`)
	listItem        = regexp.MustCompile(`^(?:(?:[-+*]\s+\[[ xX]\])|[-+*]|[0-9]+[.)])\s+`)
	metadataLine    = regexp.MustCompile(`^\*\*[^*]+\*\*:\s*`)
	setextUnderline = regexp.MustCompile(`^(?:=+|-+)\s*$`)
	tableDelimiter  = regexp.MustCompile(`^\s*:?-{3,}:?\s*(?:\|\s*:?-{3,}:?\s*)+\|?\s*$`)
	mojibake        = []string{string(rune(0xfffd)), string(rune(0x00c3)), string(rune(0x00c2)), string([]rune{0x00e2, 0x20ac}), string([]rune{0x00f0, 0x0178})}
)

func main() {
	fix := flag.Bool("fix", false, "rewrite repairable repository violations in place")
	stdin := flag.Bool("stdin", false, "format Markdown from standard input")
	flag.Parse()

	if *stdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fail(2, "read stdin: %v", err)
		}
		formatted, err := formatPublication(data)
		if err != nil {
			fail(1, "%v", err)
		}
		if _, err := os.Stdout.Write(formatted); err != nil {
			fail(2, "write stdout: %v", err)
		}
		return
	}

	root := "."
	if flag.NArg() == 1 {
		root = flag.Arg(0)
	} else if flag.NArg() > 1 {
		fail(2, "usage: github-format [-fix] [-stdin] [repository-root]")
	}

	problems, err := checkRepository(root, *fix)
	if err != nil {
		fail(2, "%v", err)
	}
	if len(problems) > 0 {
		for _, problem := range problems {
			fmt.Fprintln(os.Stderr, problem)
		}
		fail(1, "FAILED with %d file(s)", len(problems))
	}
	fmt.Println("github-format: OK")
}

func fail(code int, format string, args ...any) {
	fmt.Fprintf(os.Stderr, "github-format: "+format+"\n", args...)
	os.Exit(code)
}

func formatPublication(data []byte) ([]byte, error) {
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("stdin is not valid UTF-8")
	}
	data = bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf})
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	if bad := firstMojibake(text); bad != "" {
		return nil, fmt.Errorf("stdin contains probable mojibake marker %q", bad)
	}
	return []byte(formatMarkdown(text)), nil
}

func checkRepository(root string, fix bool) ([]string, error) {
	var problems []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if entry.IsDir() {
			if relative != "." && skipDirectory(relative, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(relative, ".specify/") && !strings.HasPrefix(relative, ".specify/memory/") {
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			if textFile(path) {
				problems = append(problems, relative+": symlinked text file is not allowed")
			}
			return nil
		}
		if !textFile(path) {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		formatted, reasons := inspectAndFormat(path, data)
		if len(reasons) == 0 {
			return nil
		}
		repairable := utf8.Valid(data) && !containsReason(reasons, "probable mojibake")
		if fix && repairable {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if err := os.WriteFile(path, formatted, info.Mode().Perm()); err != nil {
				return err
			}
			return nil
		}
		problems = append(problems, relative+": "+strings.Join(reasons, ", "))
		return nil
	})
	sort.Strings(problems)
	return problems, err
}

func containsReason(reasons []string, prefix string) bool {
	for _, reason := range reasons {
		if strings.HasPrefix(reason, prefix) {
			return true
		}
	}
	return false
}

func inspectAndFormat(path string, data []byte) ([]byte, []string) {
	if !utf8.Valid(data) {
		return data, []string{"invalid UTF-8"}
	}
	var reasons []string
	text := string(data)
	if strings.HasPrefix(text, "\ufeff") {
		reasons = append(reasons, "UTF-8 BOM")
		text = strings.TrimPrefix(text, "\ufeff")
	}
	if bad := firstMojibake(text); bad != "" {
		reasons = append(reasons, "probable mojibake "+fmt.Sprintf("%q", bad))
	}
	wantsCRLF := windowsScript(path)
	if wantsCRLF && strings.Contains(strings.ReplaceAll(text, "\r\n", ""), "\n") {
		reasons = append(reasons, "non-CRLF line ending")
	}
	if !wantsCRLF && strings.Contains(text, "\r") {
		reasons = append(reasons, "non-LF line ending")
	}

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	formatted := normalized
	if strings.EqualFold(filepath.Ext(path), ".md") {
		formatted = formatMarkdown(formatted)
	}
	if wantsCRLF {
		formatted = strings.ReplaceAll(formatted, "\n", "\r\n")
	}
	if formatted != string(data) && len(reasons) == 0 {
		reasons = append(reasons, "formatting")
	}
	return []byte(formatted), reasons
}

func firstMojibake(text string) string {
	for _, marker := range mojibake {
		if strings.Contains(text, marker) {
			return marker
		}
	}
	return ""
}

func skipDirectory(relative, name string) bool {
	if name == ".git" || name == ".idea" || name == ".vscode" || name == "dist" || name == "node_modules" || name == "testdata" || name == ".agents" {
		return true
	}
	switch relative {
	case ".specify/integrations", ".specify/scripts", ".specify/templates", ".specify/workflows",
		"brand", "docs/assets/brand", "docs/assets/favicons", "docs/assets/fonts",
		"site/.next", "site/.source", "site/.wrangler-dry-run", "site/content/generated", "site/out", "site/playwright-report", "site/public", "site/test-results":
		return true
	default:
		return false
	}
}

func textFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".bat", ".cmd", ".css", ".go", ".html", ".js", ".json", ".jsx", ".md", ".ps1", ".py", ".scss", ".sh", ".svg", ".toml", ".ts", ".tsx", ".txt", ".xml", ".yaml", ".yml":
		return true
	}
	switch filepath.Base(path) {
	case ".editorconfig", ".gitattributes", ".gitignore", "CODEOWNERS", "Dockerfile", "LICENSE", "Makefile", "NOTICE":
		return true
	default:
		return false
	}
}

func windowsScript(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".ps1", ".cmd", ".bat":
		return true
	default:
		return false
	}
}

func replaceEmDashes(input string) string {
	input = strings.ReplaceAll(input, " \u2014 ", ", ")
	return strings.ReplaceAll(input, "\u2014", "-")
}

func replaceEmDashesInMarkdownProse(input string) string {
	var output strings.Builder
	for index := 0; index < len(input); {
		if input[index] == '`' {
			runLength := 1
			for index+runLength < len(input) && input[index+runLength] == '`' {
				runLength++
			}
			delimiter := strings.Repeat("`", runLength)
			end := strings.Index(input[index+runLength:], delimiter)
			if end >= 0 {
				end += index + runLength + runLength
				output.WriteString(input[index:end])
				index = end
				continue
			}
		}
		if strings.HasPrefix(input[index:], "](") {
			if end := markdownDestinationEnd(input, index); end >= 0 {
				output.WriteString(input[index:end])
				index = end
				continue
			}
		}
		if bareURLStart(input[index:]) {
			end := bareURLEnd(input, index)
			output.WriteString(input[index:end])
			index = end
			continue
		}
		if input[index] == '<' {
			end := strings.IndexByte(input[index:], '>')
			if end >= 0 {
				end += index + 1
				output.WriteString(input[index:end])
				index = end
				continue
			}
		}
		end := index
		for end < len(input) && input[end] != '`' && input[end] != '<' && !strings.HasPrefix(input[end:], "](") && !bareURLStart(input[end:]) {
			end++
		}
		if end == index {
			output.WriteByte(input[index])
			index++
			continue
		}
		output.WriteString(replaceEmDashes(input[index:end]))
		index = end
	}
	return output.String()
}

func markdownDestinationEnd(input string, start int) int {
	depth := 1
	for index := start + 2; index < len(input); index++ {
		if input[index] == '\\' && index+1 < len(input) {
			index++
			continue
		}
		switch input[index] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return index + 1
			}
		}
	}
	return -1
}

func bareURLStart(input string) bool {
	return hasPrefixFold(input, "https://") || hasPrefixFold(input, "http://") || hasPrefixFold(input, "www.")
}

func hasPrefixFold(input, prefix string) bool {
	return len(input) >= len(prefix) && strings.EqualFold(input[:len(prefix)], prefix)
}

func bareURLEnd(input string, start int) int {
	end := start
	for end < len(input) && !strings.ContainsRune(" \t\r\n<>", rune(input[end])) {
		end++
	}
	return end
}

func formatMarkdown(input string) string {
	normalized := strings.ReplaceAll(input, "\r\n", "\n")
	hasFinalNewline := strings.HasSuffix(normalized, "\n")
	lines := strings.Split(strings.TrimSuffix(normalized, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return input
	}

	out := make([]string, 0, len(lines))
	var fence byte
	var fenceLength int
	inFrontMatter := len(lines) > 0 && strings.TrimSpace(lines[0]) == "---"
	inHTMLComment := false
	for index := 0; index < len(lines); {
		line := lines[index]
		trimmed := strings.TrimSpace(line)
		if index == 0 && inFrontMatter {
			out = append(out, line)
			index++
			continue
		}
		if inFrontMatter {
			out = append(out, line)
			index++
			if trimmed == "---" {
				inFrontMatter = false
			}
			continue
		}
		if inHTMLComment {
			out = append(out, line)
			index++
			if strings.Contains(line, "-->") {
				inHTMLComment = false
			}
			continue
		}
		if strings.HasPrefix(trimmed, "<!--") {
			out = append(out, line)
			index++
			inHTMLComment = !strings.Contains(line, "-->")
			continue
		}
		if fenceLength > 0 {
			out = append(out, line)
			index++
			if closesFence(trimmed, fence, fenceLength) {
				fenceLength = 0
			}
			continue
		}
		if character, length, ok := fenceDelimiter(trimmed); ok {
			fence, fenceLength = character, length
			out = append(out, line)
			index++
			continue
		}
		if index+1 < len(lines) && setextUnderline.MatchString(strings.TrimSpace(lines[index+1])) {
			out = append(out, replaceEmDashesInMarkdownProse(line), lines[index+1])
			index += 2
			continue
		}
		if tableStart(lines, index) {
			for index < len(lines) && strings.Contains(lines[index], "|") {
				out = append(out, lines[index])
				index++
			}
			continue
		}
		if literalMarkdownLine(line) {
			out = append(out, line)
			index++
			continue
		}
		if content, ok := blockquoteProse(line); ok {
			joined := content
			index++
			for index < len(lines) {
				content, ok = blockquoteProse(lines[index])
				if !ok {
					break
				}
				joined += " " + content
				index++
			}
			out = append(out, "> "+replaceEmDashesInMarkdownProse(joined))
			continue
		}
		joined := strings.TrimRight(line, " \t")
		marker := listItem.FindString(trimmed)
		index++
		for index < len(lines) && proseContinuation(lines, index, marker != "", leadingIndent(line)+len(marker)) {
			joined += " " + strings.TrimSpace(lines[index])
			index++
		}
		out = append(out, replaceEmDashesInMarkdownProse(joined))
	}
	result := strings.Join(out, "\n")
	if hasFinalNewline {
		result += "\n"
	}
	return result
}

func fenceDelimiter(trimmed string) (byte, int, bool) {
	if trimmed == "" || (trimmed[0] != '`' && trimmed[0] != '~') {
		return 0, 0, false
	}
	character := trimmed[0]
	length := 0
	for length < len(trimmed) && trimmed[length] == character {
		length++
	}
	return character, length, length >= 3
}

func closesFence(trimmed string, character byte, minimumLength int) bool {
	found, length, ok := fenceDelimiter(trimmed)
	return ok && found == character && length >= minimumLength && strings.TrimSpace(trimmed[length:]) == ""
}

func proseContinuation(lines []string, index int, parentIsList bool, contentIndent int) bool {
	line := lines[index]
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || listItem.MatchString(trimmed) || markdownBlockStart(trimmed) || literalMarkdownLine(line) {
		return false
	}
	if tableStart(lines, index) || (index+1 < len(lines) && setextUnderline.MatchString(strings.TrimSpace(lines[index+1]))) {
		return false
	}
	if parentIsList {
		return leadingIndent(line) < contentIndent+4
	}
	return true
}

func tableStart(lines []string, index int) bool {
	return index+1 < len(lines) && strings.Contains(lines[index], "|") && tableDelimiter.MatchString(lines[index+1])
}

func leadingIndent(line string) int {
	indent := 0
	for _, character := range line {
		if character == ' ' {
			indent++
		} else if character == '\t' {
			indent += 4
		} else {
			break
		}
	}
	return indent
}

func blockquoteProse(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "> ") {
		return "", false
	}
	content := strings.TrimSpace(strings.TrimPrefix(trimmed, ">"))
	if content == "" || alertMarker.MatchString(content) || listItem.MatchString(content) || markdownBlockStart(content) {
		return "", false
	}
	return content, true
}

func markdownBlockStart(trimmed string) bool {
	_, _, isFence := fenceDelimiter(trimmed)
	return isFence || headingLine.MatchString(trimmed) || setextUnderline.MatchString(trimmed) || tableDelimiter.MatchString(trimmed) || strings.HasPrefix(trimmed, ">") || strings.HasPrefix(trimmed, "<") || strings.HasPrefix(trimmed, "|") || strings.HasPrefix(trimmed, "![") || strings.HasPrefix(trimmed, ":::") || strings.HasPrefix(trimmed, "$$")
}

func literalMarkdownLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") || strings.HasSuffix(line, "\\") || strings.HasSuffix(line, "  ") {
		return true
	}
	if markdownBlockStart(trimmed) || trimmed == "---" || trimmed == "***" || trimmed == "___" || linkDefinition.MatchString(trimmed) || metadataLine.MatchString(trimmed) {
		return true
	}
	return false
}
