package scripted

import (
	"bytes"
	"context"
	"strings"
)

// Detect distinguishes an accepted dialect and rejected scripted candidates.
// Rejection evidence prevents filename or another codec from hiding ambiguity.
func Detect(data []byte) Detection {
	// ASCII signatures also recognize malformed encodings as script candidates.
	candidate := false
	first := true
	for rest := data; len(rest) > 0; {
		line, remaining, _ := bytes.Cut(rest, []byte{'\n'})
		rest = remaining
		line = bytes.TrimSpace(bytes.TrimPrefix(line, []byte{0xef, 0xbb, 0xbf}))
		if len(line) == 0 {
			continue
		}
		if first {
			first = false
			if line[0] != '[' || line[len(line)-1] != ']' {
				return Detection{}
			}
		}
		if len(line) >= 2 && line[0] == '[' && line[len(line)-1] == ']' {
			name := bytes.TrimSpace(line[1 : len(line)-1])
			if bytes.EqualFold(name, []byte("script info")) || (len(name) >= 2 && bytes.EqualFold(name[:2], []byte("v4"))) {
				candidate = true
				break
			}
		}
	}
	if !candidate {
		return Detection{}
	}
	lines, _, err := physicalLines(context.Background(), data)
	if err != nil {
		return Detection{Candidate: true, Err: err}
	}
	return detectLines(lines)
}

func detectLines(lines []string) Detection {
	first, candidate := true, false
	section := ""
	dialect := ""
	seenInfo, assStyles, ssaStyles := false, false, false
	for order, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if first {
			first = false
			if _, ok := sectionHeader(line); !ok {
				return Detection{}
			}
		}
		if name, ok := sectionHeader(line); ok {
			section = asciiName(name)
			if section == "script info" || strings.HasPrefix(section, "v4") {
				candidate = true
			}
			switch section {
			case "script info":
				seenInfo = true
			case "v4+ styles":
				assStyles = true
			case "v4 styles":
				ssaStyles = true
			}
			continue
		}
		key, value, ok := strings.Cut(strings.TrimLeft(line, " "), ":")
		if section == "script info" && ok && asciiName(key) == "scripttype" {
			value = strings.TrimSpace(value)
			if value != "v4.00+" && value != "v4.00" {
				return Detection{Candidate: true, Err: failure("invalid_scripted_dialect", order)}
			}
			if dialect != "" && dialect != value {
				return Detection{Candidate: true, Err: failure("invalid_scripted_dialect", order)}
			}
			dialect = value
		}
	}
	if !candidate {
		return Detection{}
	}
	if !seenInfo || dialect == "" || assStyles == ssaStyles || (dialect == "v4.00+" && !assStyles) || (dialect == "v4.00" && !ssaStyles) {
		return Detection{Candidate: true, Err: failure("invalid_scripted_dialect", 0)}
	}
	format := "ass"
	if dialect == "v4.00" {
		format = "ssa"
	}
	return Detection{Candidate: true, Format: format}
}
