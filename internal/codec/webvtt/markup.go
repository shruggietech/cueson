package webvtt

import (
	"html"
	"sort"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

type payloadLexeme struct {
	start, end int
	raw        string
	name       string
	annotation string
	kind       string
	closing    bool
	matched    bool
	timestamp  int64
}

type diagnosticAt struct {
	offset int
	diag   Diagnostic
}

func scanPayload(raw string, cueStart, cueEnd int64, sourceOrder int, cueID string) (string, []model.Speaker, []model.Token, []Diagnostic) {
	lexemes, pending := lexPayload(raw, sourceOrder, cueID)
	pairMarkup(lexemes)
	validTimestamp := make([]bool, len(lexemes))
	previousTimestamp := cueStart
	for index := range lexemes {
		lexeme := &lexemes[index]
		if lexeme.kind != "timestamp" {
			continue
		}
		if lexeme.timestamp <= previousTimestamp || lexeme.timestamp >= cueEnd {
			pending = append(pending, diagnosticAt{offset: lexeme.start, diag: cueDiagnostic("webvtt_inline_timestamp_invalid", "invalid or non-increasing inline timestamp was preserved literally", sourceOrder, cueID)})
			continue
		}
		validTimestamp[index] = true
		previousTimestamp = lexeme.timestamp
	}

	var plain strings.Builder
	var span strings.Builder
	speakers := make([]model.Speaker, 0)
	tokens := make([]model.Token, 0)
	spanStart := cueStart
	var nativeTiming *string
	cursor := 0
	for index := range lexemes {
		lexeme := &lexemes[index]
		text, entityDiagnostics := decodeEntities(raw[cursor:lexeme.start], sourceOrder, cueID, cursor)
		plain.WriteString(text)
		span.WriteString(text)
		pending = append(pending, entityDiagnostics...)
		switch lexeme.kind {
		case "timestamp":
			if validTimestamp[index] {
				appendSpanToken(&tokens, span.String(), spanStart, lexeme.timestamp, nativeTiming)
				span.Reset()
				spanStart = lexeme.timestamp
				native := lexeme.raw
				nativeTiming = &native
			} else {
				plain.WriteString(semanticText(lexeme.raw))
				span.WriteString(semanticText(lexeme.raw))
			}
		case "timestamp_invalid":
			plain.WriteString(semanticText(lexeme.raw))
			span.WriteString(semanticText(lexeme.raw))
			pending = append(pending, diagnosticAt{offset: lexeme.start, diag: cueDiagnostic("webvtt_inline_timestamp_invalid", "malformed inline timestamp was preserved literally", sourceOrder, cueID)})
		case "tag":
			if !lexeme.matched {
				plain.WriteString(semanticText(lexeme.raw))
				span.WriteString(semanticText(lexeme.raw))
				pending = append(pending, diagnosticAt{offset: lexeme.start, diag: cueDiagnostic("webvtt_markup_unbalanced", "unbalanced WebVTT markup was preserved literally", sourceOrder, cueID)})
			} else if !lexeme.closing && lexeme.name == "v" {
				name, diagnostics := decodeEntities(lexeme.annotation, sourceOrder, cueID, lexeme.start)
				pending = append(pending, diagnostics...)
				speakers = append(speakers, model.Speaker{Name: name, Origin: "native"})
			}
		default:
			plain.WriteString(semanticText(lexeme.raw))
			span.WriteString(semanticText(lexeme.raw))
			pending = append(pending, diagnosticAt{offset: lexeme.start, diag: cueDiagnostic("webvtt_markup_unknown", "unknown WebVTT markup was preserved literally", sourceOrder, cueID)})
		}
		cursor = lexeme.end
	}
	text, entityDiagnostics := decodeEntities(raw[cursor:], sourceOrder, cueID, cursor)
	plain.WriteString(text)
	span.WriteString(text)
	pending = append(pending, entityDiagnostics...)
	if nativeTiming != nil {
		appendSpanToken(&tokens, span.String(), spanStart, cueEnd, nativeTiming)
	}
	sort.SliceStable(pending, func(left, right int) bool { return pending[left].offset < pending[right].offset })
	diagnostics := make([]Diagnostic, len(pending))
	for index := range pending {
		diagnostics[index] = pending[index].diag
	}
	return plain.String(), speakers, tokens, diagnostics
}

func lexPayload(raw string, sourceOrder int, cueID string) ([]payloadLexeme, []diagnosticAt) {
	lexemes := make([]payloadLexeme, 0)
	diagnostics := make([]diagnosticAt, 0)
	for position := 0; position < len(raw); {
		relative := strings.IndexByte(raw[position:], '<')
		if relative < 0 {
			break
		}
		start := position + relative
		relativeEnd := strings.IndexByte(raw[start+1:], '>')
		if relativeEnd < 0 {
			diagnostics = append(diagnostics, diagnosticAt{offset: start, diag: cueDiagnostic("webvtt_markup_unbalanced", "unterminated WebVTT markup was preserved literally", sourceOrder, cueID)})
			break
		}
		end := start + 1 + relativeEnd + 1
		rawToken := raw[start:end]
		inner := raw[start+1 : end-1]
		if timestamp, err := ParseTimestamp(inner); err == nil {
			lexemes = append(lexemes, payloadLexeme{start: start, end: end, raw: rawToken, kind: "timestamp", timestamp: timestamp})
			position = end
			continue
		}
		if len(inner) > 0 && inner[0] >= '0' && inner[0] <= '9' && strings.ContainsRune(inner, ':') {
			lexemes = append(lexemes, payloadLexeme{start: start, end: end, raw: rawToken, kind: "timestamp_invalid"})
			position = end
			continue
		}
		name, annotation, closing, known := parseTag(inner)
		kind := "unknown"
		if known {
			kind = "tag"
		}
		lexemes = append(lexemes, payloadLexeme{start: start, end: end, raw: rawToken, name: name, annotation: annotation, closing: closing, kind: kind})
		position = end
	}
	return lexemes, diagnostics
}

func parseTag(inner string) (name, annotation string, closing, known bool) {
	if inner == "" {
		return "", "", false, false
	}
	closing = inner[0] == '/'
	if closing {
		inner = inner[1:]
	}
	boundary := len(inner)
	for index := range inner {
		if inner[index] == '.' || isASCIISpace(inner[index]) {
			boundary = index
			break
		}
	}
	name = inner[:boundary]
	if name != "c" && name != "i" && name != "b" && name != "u" && name != "ruby" && name != "rt" && name != "v" && name != "lang" {
		return name, "", closing, false
	}
	if closing {
		return name, "", true, boundary == len(inner)
	}
	annotationStart := strings.IndexAny(inner, " \t")
	if annotationStart >= 0 {
		annotation = normalizeASCIIWhitespace(inner[annotationStart:])
	}
	if name == "v" && annotation == "" {
		return name, "", false, false
	}
	return name, annotation, false, true
}

func pairMarkup(lexemes []payloadLexeme) {
	stack := make([]int, 0)
	for index := range lexemes {
		if lexemes[index].kind != "tag" {
			continue
		}
		if !lexemes[index].closing {
			stack = append(stack, index)
			continue
		}
		if len(stack) == 0 {
			continue
		}
		opening := stack[len(stack)-1]
		if lexemes[opening].name != lexemes[index].name {
			continue
		}
		stack = stack[:len(stack)-1]
		lexemes[opening].matched = true
		lexemes[index].matched = true
	}
	if len(stack) == 1 {
		opening := stack[0]
		if lexemes[opening].name == "v" && lexemes[opening].start == 0 {
			lexemes[opening].matched = true
		}
	}
}

func decodeEntities(text string, sourceOrder int, cueID string, baseOffset int) (string, []diagnosticAt) {
	var output strings.Builder
	diagnostics := make([]diagnosticAt, 0)
	for position := 0; position < len(text); {
		relative := strings.IndexByte(text[position:], '&')
		if relative < 0 {
			output.WriteString(semanticText(text[position:]))
			break
		}
		start := position + relative
		output.WriteString(semanticText(text[position:start]))
		semicolon := strings.IndexByte(text[start+1:], ';')
		if semicolon < 0 {
			output.WriteByte('&')
			diagnostics = append(diagnostics, diagnosticAt{offset: baseOffset + start, diag: cueDiagnostic("webvtt_entity_invalid", "unterminated character reference was preserved literally", sourceOrder, cueID)})
			position = start + 1
			continue
		}
		end := start + 1 + semicolon + 1
		entity := text[start:end]
		decoded, valid := decodeWebVTTEntity(entity)
		if !valid {
			output.WriteString(semanticText(entity))
			diagnostics = append(diagnostics, diagnosticAt{offset: baseOffset + start, diag: cueDiagnostic("webvtt_entity_invalid", "unknown or malformed character reference was preserved literally", sourceOrder, cueID)})
		} else {
			output.WriteString(decoded)
		}
		position = end
	}
	return output.String(), diagnostics
}

func decodeWebVTTEntity(entity string) (string, bool) {
	decoded := html.UnescapeString(entity)
	if decoded != entity {
		return decoded, true
	}
	return entity, false
}

func semanticText(value string) string {
	return strings.ReplaceAll(value, "\x00", "\ufffd")
}

func normalizeASCIIWhitespace(value string) string {
	return strings.Join(strings.FieldsFunc(value, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == '\f'
	}), " ")
}

func appendSpanToken(tokens *[]model.Token, text string, start, end int64, native *string) {
	if text == "" {
		return
	}
	var nativeCopy *string
	if native != nil {
		value := *native
		nativeCopy = &value
	}
	*tokens = append(*tokens, model.Token{Text: text, StartMilliseconds: start, EndMilliseconds: end, NativeTiming: nativeCopy})
}
