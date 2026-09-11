package subrip

import (
	"regexp"
	"strings"
)

var (
	sameLineSpeaker  = regexp.MustCompile(`^(\p{L}[^:\r\n]{0,48}):[ \t]+(\S.*)$`)
	labelLineSpeaker = regexp.MustCompile(`^(\p{L}[^:\r\n]{0,48}):[ \t]*$`)
)

// PlainText removes only documented SubRip presentation tags and keeps their content.
func PlainText(input string) string {
	type tag struct {
		start   int
		end     int
		name    string
		closing bool
	}
	tags := make([]tag, 0)
	for offset := 0; offset < len(input); {
		start := strings.IndexByte(input[offset:], '<')
		if start < 0 {
			break
		}
		start += offset
		end := strings.IndexByte(input[start+1:], '>')
		if end < 0 {
			break
		}
		end += start + 1
		if name, closing, ok := documentedTag(input[start+1 : end]); ok {
			tags = append(tags, tag{start: start, end: end + 1, name: name, closing: closing})
		}
		offset = end + 1
	}
	remove := make([]bool, len(tags))
	stack := make([]int, 0)
	for index := range tags {
		if !tags[index].closing {
			stack = append(stack, index)
			continue
		}
		if len(stack) == 0 {
			continue
		}
		if tags[stack[len(stack)-1]].name != tags[index].name {
			stack = stack[:0]
			continue
		}
		opening := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		remove[opening], remove[index] = true, true
	}
	var output strings.Builder
	offset := 0
	for index := range tags {
		if !remove[index] {
			continue
		}
		output.WriteString(input[offset:tags[index].start])
		offset = tags[index].end
	}
	output.WriteString(input[offset:])
	return output.String()
}

func documentedTag(inner string) (string, bool, bool) {
	normalized := strings.ToLower(strings.TrimSpace(inner))
	closing := strings.HasPrefix(normalized, "/")
	if closing {
		normalized = strings.TrimPrefix(normalized, "/")
	}
	if normalized == "b" || normalized == "i" || normalized == "u" || normalized == "font" {
		return normalized, closing, true
	}
	if !closing && strings.HasPrefix(normalized, "font ") && !strings.ContainsAny(normalized, "<>") {
		return "font", false, true
	}
	return "", false, false
}

// DetectSpeaker applies the documented transcript-prefix heuristic without changing the input lines.
func DetectSpeaker(lines []string) *string {
	if len(lines) == 0 {
		return nil
	}
	if match := sameLineSpeaker.FindStringSubmatch(lines[0]); match != nil && !strings.ContainsAny(match[1], "?!") {
		name := strings.TrimSpace(match[1])
		return &name
	}
	if len(lines) > 1 {
		if match := labelLineSpeaker.FindStringSubmatch(lines[0]); match != nil && !strings.ContainsAny(match[1], "?!") {
			name := strings.TrimSpace(match[1])
			return &name
		}
	}
	return nil
}
