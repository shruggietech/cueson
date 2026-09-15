package scripted

import (
	"context"
	"github.com/shruggietech/cueson/internal/model"
	"strings"
	"unicode/utf8"
)

func physicalLines(ctx context.Context, data []byte) ([]string, model.EncodingObservation, error) {
	o := model.EncodingObservation{LineEndings: "none", DetectedEncoding: ptr("utf-8"), Confidence: ptr(1.0)}
	if len(data) > 64<<20 {
		return nil, o, failure("complexity_limit", 0)
	}
	if len(data) == 0 || !utf8.Valid(data) {
		return nil, o, failure("invalid_scripted_encoding", 0)
	}
	text := string(data)
	if strings.HasPrefix(text, "\ufeff") {
		text = text[len("\ufeff"):]
		o.BOM = ptr("utf-8")
	}
	lines := make([]string, 0)
	hasLF, hasCRLF := false, false
	for text != "" {
		if err := ctx.Err(); err != nil {
			return nil, o, err
		}
		if len(lines) >= model.MaxDocumentItems {
			return nil, o, failure("complexity_limit", len(lines))
		}
		line, rest, ended := strings.Cut(text, "\n")
		text = rest
		if ended {
			if strings.HasSuffix(line, "\r") {
				hasCRLF = true
				line = strings.TrimSuffix(line, "\r")
			} else {
				hasLF = true
			}
		}
		if len(line) > model.MaxScriptedLineBytes {
			return nil, o, failure("complexity_limit", len(lines))
		}
		if strings.ContainsAny(line, "\r\x00") {
			return nil, o, failure("invalid_scripted_encoding", len(lines))
		}
		lines = append(lines, line)
	}
	if hasLF && hasCRLF {
		o.LineEndings = "mixed"
	} else if hasLF {
		o.LineEndings = "lf"
	} else if hasCRLF {
		o.LineEndings = "crlf"
	}
	return lines, o, nil
}

func asciiName(s string) string {
	b := []byte(strings.Trim(s, " "))
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}
func sectionHeader(line string) (string, bool) {
	v := strings.TrimSpace(line)
	if len(v) >= 2 && v[0] == '[' && v[len(v)-1] == ']' {
		return v[1 : len(v)-1], true
	}
	return "", false
}
