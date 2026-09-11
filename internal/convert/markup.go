package convert

import (
	"fmt"
	"html"
	"strings"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
)

type payloadTranslation struct {
	Text   string
	Issues []payloadIssue
}

type payloadIssue struct {
	Code       string
	Kind       Kind
	Feature    string
	Occurrence int
}

type payloadIssueCollector struct {
	issues   []payloadIssue
	overflow bool
}

func (collector *payloadIssueCollector) append(issue payloadIssue) {
	if len(collector.issues) >= MaxLosses {
		collector.overflow = true
		return
	}
	collector.issues = append(collector.issues, issue)
}

type markupToken struct {
	start      int
	end        int
	raw        string
	name       string
	closing    bool
	known      bool
	matched    bool
	timestamp  bool
	classes    bool
	annotation string
}

func translateSubRipPayload(raw string) (payloadTranslation, error) {
	if strings.ContainsRune(raw, '\r') {
		return payloadTranslation{}, fmt.Errorf("SubRip payload contains a non-physical carriage return")
	}
	if err := validateTranslationComplexity(raw, "SubRip"); err != nil {
		return payloadTranslation{}, err
	}
	tokens := scanSubRipMarkup(raw)
	pairSubRipMarkup(tokens)
	collector := payloadIssueCollector{issues: make([]payloadIssue, 0)}
	var output strings.Builder
	cursor := 0
	fontOccurrence := 0
	nulOccurrence := 0
	for index := range tokens {
		token := &tokens[index]
		writeWebVTTText(&output, raw[cursor:token.start], &collector, &nulOccurrence)
		switch {
		case token.matched && (token.name == "b" || token.name == "i" || token.name == "u"):
			if token.closing {
				output.WriteString("</" + token.name + ">")
			} else {
				output.WriteString("<" + token.name + ">")
			}
		case token.matched && token.name == "font":
			if !token.closing {
				collector.append(payloadIssue{Code: LossCodeSubRipFontDegraded, Kind: KindDegraded, Occurrence: fontOccurrence})
				fontOccurrence++
			}
		default:
			writeWebVTTText(&output, token.raw, &collector, &nulOccurrence)
		}
		cursor = token.end
	}
	writeWebVTTText(&output, raw[cursor:], &collector, &nulOccurrence)
	translated := output.String()
	translated = replaceEmptyLines(translated, raw, "<i></i>", &collector)
	if collector.overflow {
		return payloadTranslation{}, fmt.Errorf("SubRip payload produces more than %d conversion losses", MaxLosses)
	}
	return payloadTranslation{Text: translated, Issues: collector.issues}, nil
}

func translateWebVTTPayload(raw string, cueStart, cueEnd int64) (payloadTranslation, error) {
	if strings.ContainsRune(raw, '\r') {
		return payloadTranslation{}, fmt.Errorf("WebVTT payload contains a non-physical carriage return")
	}
	if err := validateTranslationComplexity(raw, "WebVTT"); err != nil {
		return payloadTranslation{}, err
	}
	tokens := scanWebVTTMarkup(raw, cueStart, cueEnd)
	pairMarkup(tokens, true)
	collector := payloadIssueCollector{issues: make([]payloadIssue, 0)}
	occurrences := make(map[string]int)
	var output strings.Builder
	var plain strings.Builder
	cursor := 0
	for index := range tokens {
		token := &tokens[index]
		text := replaceWebVTTNUL(translateWebVTTText(raw[cursor:token.start], &collector, occurrences), &collector, occurrences)
		output.WriteString(text)
		plain.WriteString(text)
		switch {
		case token.timestamp:
			appendPayloadIssue(&collector, occurrences, LossCodeWebVTTInlineTimingOmitted, KindOmitted, "inline_timestamp")
		case token.matched && (token.name == "b" || token.name == "i" || token.name == "u"):
			if token.classes && !token.closing {
				appendPayloadIssue(&collector, occurrences, LossCodeWebVTTMarkupDegraded, KindDegraded, "class")
			}
			if token.closing {
				output.WriteString("</" + token.name + ">")
			} else {
				output.WriteString("<" + token.name + ">")
			}
		case token.matched && token.known:
			if !token.closing {
				code := map[string]string{
					"c": LossCodeWebVTTMarkupDegraded, "v": LossCodeWebVTTVoiceDegraded,
					"lang": LossCodeWebVTTMarkupDegraded, "ruby": LossCodeWebVTTMarkupDegraded, "rt": LossCodeWebVTTMarkupDegraded,
				}[token.name]
				if code != "" {
					kind := KindDegraded
					appendPayloadIssue(&collector, occurrences, code, kind, token.name)
				}
			}
		default:
			literal := replaceWebVTTNUL(token.raw, &collector, occurrences)
			output.WriteString(literal)
			plain.WriteString(literal)
		}
		cursor = token.end
	}
	tail := replaceWebVTTNUL(translateWebVTTText(raw[cursor:], &collector, occurrences), &collector, occurrences)
	output.WriteString(tail)
	plain.WriteString(tail)
	translated := output.String()
	translated = replaceEmptyLines(translated, plain.String(), "<u></u>", &collector)
	if collector.overflow {
		return payloadTranslation{}, fmt.Errorf("WebVTT payload produces more than %d conversion losses", MaxLosses)
	}
	if subrip.PlainText(translated) != plain.String() {
		return payloadTranslation{}, fmt.Errorf("WebVTT payload cannot be represented without SubRip markup reinterpretation")
	}
	return payloadTranslation{Text: translated, Issues: collector.issues}, nil
}

func validateTranslationComplexity(raw, format string) error {
	if strings.Count(raw, "\n") >= model.MaxItemOccurrences {
		return fmt.Errorf("%s payload contains more than %d physical lines", format, model.MaxItemOccurrences)
	}
	if strings.Count(raw, "<") > model.MaxItemOccurrences || strings.Count(raw, "&") > model.MaxItemOccurrences {
		return fmt.Errorf("%s payload contains more than %d markup or entity occurrences", format, model.MaxItemOccurrences)
	}
	return nil
}

func scanSubRipMarkup(raw string) []markupToken {
	tokens := make([]markupToken, 0)
	for position := 0; position < len(raw); {
		startRelative := strings.IndexByte(raw[position:], '<')
		if startRelative < 0 {
			break
		}
		start := position + startRelative
		endRelative := strings.IndexByte(raw[start+1:], '>')
		if endRelative < 0 {
			break
		}
		end := start + endRelative + 2
		inner := strings.TrimSpace(raw[start+1 : end-1])
		if inner == "" {
			tokens = append(tokens, markupToken{start: start, end: end, raw: raw[start:end]})
			position = end
			continue
		}
		normalized := strings.ToLower(inner)
		closing := strings.HasPrefix(normalized, "/")
		if closing {
			normalized = strings.TrimPrefix(normalized, "/")
		}
		name := normalized
		known := false
		if closing {
			known = name == "b" || name == "i" || name == "u" || name == "font"
		} else {
			switch {
			case name == "b" || name == "i" || name == "u" || name == "font":
				known = true
			case strings.HasPrefix(name, "font ") && !strings.ContainsAny(name, "<>"):
				name = "font"
				known = true
			}
		}
		tokens = append(tokens, markupToken{start: start, end: end, raw: raw[start:end], name: name, closing: closing, known: known})
		position = end
	}
	return tokens
}

func scanWebVTTMarkup(raw string, cueStart, cueEnd int64) []markupToken {
	tokens := make([]markupToken, 0)
	previousTimestamp := cueStart
	for position := 0; position < len(raw); {
		startRelative := strings.IndexByte(raw[position:], '<')
		if startRelative < 0 {
			break
		}
		start := position + startRelative
		endRelative := strings.IndexByte(raw[start+1:], '>')
		if endRelative < 0 {
			break
		}
		end := start + endRelative + 2
		inner := raw[start+1 : end-1]
		if timestamp, err := webvtt.ParseTimestamp(inner); err == nil && timestamp > previousTimestamp && timestamp < cueEnd {
			tokens = append(tokens, markupToken{start: start, end: end, raw: raw[start:end], timestamp: true})
			previousTimestamp = timestamp
			position = end
			continue
		}
		name, closing, known, classes, annotation := parseWebVTTTag(inner)
		tokens = append(tokens, markupToken{start: start, end: end, raw: raw[start:end], name: name, closing: closing, known: known, classes: classes, annotation: annotation})
		position = end
	}
	return tokens
}

func parseWebVTTTag(inner string) (string, bool, bool, bool, string) {
	if inner == "" {
		return "", false, false, false, ""
	}
	closing := inner[0] == '/'
	if closing {
		inner = inner[1:]
	}
	boundary := len(inner)
	for index := range inner {
		if inner[index] == '.' || inner[index] == ' ' || inner[index] == '\t' {
			boundary = index
			break
		}
	}
	name := inner[:boundary]
	known := name == "c" || name == "i" || name == "b" || name == "u" || name == "ruby" || name == "rt" || name == "v" || name == "lang"
	if closing && boundary != len(inner) {
		known = false
	}
	annotation := ""
	if annotationStart := strings.IndexAny(inner, " \t"); annotationStart >= 0 {
		annotation = strings.TrimSpace(inner[annotationStart:])
	}
	if !closing && (name == "v" || name == "lang") && normalizeASCIIWhitespace(html.UnescapeString(annotation)) == "" {
		known = false
	}
	return name, closing, known, strings.Contains(inner[:max(boundary, 0)], ".") || strings.Contains(inner, "."), annotation
}

func pairMarkup(tokens []markupToken, allowSoleVoice bool) {
	stack := make([]int, 0)
	for index := range tokens {
		if !tokens[index].known || tokens[index].timestamp {
			continue
		}
		if !tokens[index].closing {
			stack = append(stack, index)
			continue
		}
		if len(stack) == 0 {
			continue
		}
		opening := stack[len(stack)-1]
		if tokens[opening].name != tokens[index].name {
			continue
		}
		stack = stack[:len(stack)-1]
		tokens[opening].matched = true
		tokens[index].matched = true
	}
	if allowSoleVoice && len(stack) == 1 {
		opening := stack[0]
		if tokens[opening].name == "v" && tokens[opening].start == 0 {
			tokens[opening].matched = true
		}
	}
}

func pairSubRipMarkup(tokens []markupToken) {
	stack := make([]int, 0)
	for index := range tokens {
		if !tokens[index].known {
			continue
		}
		if !tokens[index].closing {
			stack = append(stack, index)
			continue
		}
		if len(stack) == 0 {
			continue
		}
		opening := stack[len(stack)-1]
		if tokens[opening].name != tokens[index].name {
			stack = stack[:0]
			continue
		}
		stack = stack[:len(stack)-1]
		tokens[opening].matched = true
		tokens[index].matched = true
	}
}

func writeWebVTTText(output *strings.Builder, value string, collector *payloadIssueCollector, nulOccurrence *int) {
	for _, character := range value {
		switch character {
		case '&':
			output.WriteString("&amp;")
		case '<':
			output.WriteString("&lt;")
		case '>':
			output.WriteString("&gt;")
		case '\x00':
			output.WriteRune('\ufffd')
			collector.append(payloadIssue{Code: LossCodeNULDegraded, Kind: KindDegraded, Occurrence: *nulOccurrence})
			*nulOccurrence++
		default:
			output.WriteRune(character)
		}
	}
}

func translateWebVTTText(value string, collector *payloadIssueCollector, occurrences map[string]int) string {
	type entitySpan struct {
		start   int
		end     int
		raw     string
		decoded string
	}

	var decodedText strings.Builder
	angleEntities := make([]entitySpan, 0)
	for position := 0; position < len(value); {
		startRelative := strings.IndexByte(value[position:], '&')
		if startRelative < 0 {
			decodedText.WriteString(value[position:])
			break
		}
		start := position + startRelative
		decodedText.WriteString(value[position:start])
		endRelative := strings.IndexByte(value[start+1:], ';')
		if endRelative < 0 {
			decodedText.WriteByte('&')
			position = start + 1
			continue
		}
		end := start + endRelative + 2
		entity := value[start:end]
		decoded := html.UnescapeString(entity)
		if decoded == entity {
			decodedText.WriteString(entity)
		} else {
			decodedStart := decodedText.Len()
			decodedText.WriteString(decoded)
			if strings.ContainsAny(decoded, "<>") {
				angleEntities = append(angleEntities, entitySpan{start: decodedStart, end: decodedText.Len(), raw: entity, decoded: decoded})
			}
		}
		position = end
	}

	decoded := decodedText.String()
	if len(angleEntities) == 0 {
		return decoded
	}
	markup := scanSubRipMarkup(decoded)
	pairSubRipMarkup(markup)
	var output strings.Builder
	cursor := 0
	markupIndex := 0
	for _, entity := range angleEntities {
		output.WriteString(decoded[cursor:entity.start])
		for markupIndex < len(markup) && (!markup[markupIndex].matched || markup[markupIndex].end <= entity.start) {
			markupIndex++
		}
		ambiguous := markupIndex < len(markup) && markup[markupIndex].matched && entity.start < markup[markupIndex].end && entity.end > markup[markupIndex].start
		if ambiguous {
			output.WriteString(entity.raw)
			appendPayloadIssue(collector, occurrences, LossCodeWebVTTEntityAmbiguous, KindAmbiguous, "character_reference")
		} else {
			output.WriteString(entity.decoded)
		}
		cursor = entity.end
	}
	output.WriteString(decoded[cursor:])
	return output.String()
}

func replaceWebVTTNUL(value string, collector *payloadIssueCollector, occurrences map[string]int) string {
	if !strings.ContainsRune(value, '\x00') {
		return value
	}
	var output strings.Builder
	for _, character := range value {
		if character != '\x00' {
			output.WriteRune(character)
			continue
		}
		output.WriteRune('\ufffd')
		appendPayloadIssue(collector, occurrences, LossCodeNULDegraded, KindDegraded, "nul")
	}
	return output.String()
}

func replaceEmptyLines(translated, semantic, placeholder string, collector *payloadIssueCollector) string {
	lines := strings.Split(translated, "\n")
	semanticLines := strings.Split(semantic, "\n")
	emptyOccurrence := 0
	for index := range lines {
		if lines[index] == "" && index < len(semanticLines) && semanticLines[index] == "" {
			lines[index] = placeholder
			collector.append(payloadIssue{Code: LossCodePayloadLineDegraded, Kind: KindDegraded, Occurrence: emptyOccurrence})
			emptyOccurrence++
		}
	}
	return strings.Join(lines, "\n")
}

func appendPayloadIssue(collector *payloadIssueCollector, occurrences map[string]int, code string, kind Kind, feature string) {
	occurrence := occurrences[code]
	collector.append(payloadIssue{Code: code, Kind: kind, Feature: feature, Occurrence: occurrence})
	occurrences[code] = occurrence + 1
}

func semanticText(value string) string {
	return strings.ReplaceAll(value, "\x00", "\ufffd")
}

func normalizeASCIIWhitespace(value string) string {
	return strings.Join(strings.FieldsFunc(value, func(character rune) bool {
		return character == ' ' || character == '\t' || character == '\r' || character == '\n' || character == '\f'
	}), " ")
}
