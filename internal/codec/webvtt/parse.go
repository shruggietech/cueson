package webvtt

import (
	"fmt"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

// Parse converts decoded WebVTT text into common and native model views.
func Parse(input string) (Result, error) {
	lexicalLines := splitPhysicalLines(input)
	if len(lexicalLines) == 0 {
		return Result{}, &ParseError{Code: "webvtt_signature_invalid", Line: 1, Text: "input is empty"}
	}
	lines := lexicalLines
	result := Result{Cues: []model.Cue{}, Diagnostics: []Diagnostic{}}
	if strings.ContainsRune(input, '\x00') {
		lines = splitPhysicalLines(strings.ReplaceAll(input, "\x00", "\ufffd"))
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: "warning", Code: "webvtt_nul_replaced", Message: "embedded NUL was replaced only in the decoded semantic view"})
	}
	description, ok := parseSignature(lines[0].text)
	if !ok {
		return Result{}, &ParseError{Code: "webvtt_signature_invalid", Line: 1, Text: "signature must be WEBVTT followed by a valid boundary"}
	}
	result.DocumentData = model.WebVTTDocumentData{
		Signature: "WEBVTT", SignatureLineRaw: lexicalLines[0].text, Description: description,
		MetadataLines: []string{}, Blocks: []model.WebVTTBlock{},
	}
	if len(lines) == 1 {
		return Result{}, &ParseError{Code: "webvtt_no_cues", Text: "input contains no valid WebVTT cues"}
	}

	position := 1
	for position < len(lines) && !blank(lines[position].text) {
		if cueStartAt(lines, position) || noteLine(lines[position].text) || keywordLine(lines[position].text, "STYLE") || keywordLine(lines[position].text, "REGION") {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: "warning", Code: "webvtt_header_separator_missing", Message: "body began without a blank line after the WebVTT header"})
			break
		}
		if strings.Contains(lines[position].text, "-->") {
			return Result{}, &ParseError{Code: "webvtt_header_invalid", Line: lines[position].number, Text: "header metadata must not contain -->"}
		}
		if len(result.DocumentData.MetadataLines) >= model.MaxItemOccurrences {
			return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[position].number, Text: fmt.Sprintf("header contains more than %d metadata lines", model.MaxItemOccurrences)}
		}
		result.DocumentData.MetadataLines = append(result.DocumentData.MetadataLines, lexicalLines[position].text)
		position++
	}
	for position < len(lines) && blank(lines[position].text) {
		position++
	}

	seenCue := false
	identifiers := make(map[string]bool)
	var previousStart int64
	for position < len(lines) {
		for position < len(lines) && blank(lines[position].text) {
			position++
		}
		if position >= len(lines) {
			break
		}
		sourceOrder := len(result.Cues) + len(result.DocumentData.Blocks)
		if sourceOrder >= model.MaxDocumentItems {
			return Result{}, &ParseError{Code: "webvtt_item_limit_exceeded", Line: lines[position].number, Text: fmt.Sprintf("document contains more than %d body items", model.MaxDocumentItems)}
		}
		line := lines[position].text
		if noteLine(line) || keywordLine(line, "STYLE") || keywordLine(line, "REGION") {
			end := position + 1
			for end < len(lines) && !blank(lines[end].text) {
				end++
			}
			rawLines := lineTexts(lexicalLines[position:end])
			if len(rawLines) > model.MaxItemOccurrences {
				return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[position].number, Text: fmt.Sprintf("block contains more than %d physical lines", model.MaxItemOccurrences)}
			}
			blockType := strings.ToLower(strings.Fields(line)[0])
			block := model.WebVTTBlock{Type: blockType, SourceOrder: sourceOrder, Raw: strings.Join(rawLines, "\n"), RawLines: rawLines}
			if (blockType == "style" || blockType == "region") && seenCue {
				result.Diagnostics = append(result.Diagnostics, blockDiagnostic("webvtt_"+blockType+"_after_cue", strings.ToUpper(blockType)+" block after the first cue was preserved", sourceOrder))
			}
			if blockType == "note" || blockType == "style" {
				for _, content := range rawLines[1:] {
					if strings.Contains(content, "-->") {
						result.Diagnostics = append(result.Diagnostics, blockDiagnostic("webvtt_"+blockType+"_invalid", strings.ToUpper(blockType)+" content containing --> was preserved", sourceOrder))
						break
					}
				}
			}
			if blockType == "region" {
				settingsRaw := strings.Join(rawLines[1:], "\n")
				if countASCIIWhitespaceFields(settingsRaw) > model.MaxItemOccurrences {
					return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[position].number, Text: fmt.Sprintf("REGION contains more than %d setting occurrences", model.MaxItemOccurrences)}
				}
				occurrences, settings, diagnostics := parseRegionSettings(settingsRaw, sourceOrder)
				block.Region = &model.WebVTTRegionData{SettingsRaw: settingsRaw, Settings: settings, SettingOccurrences: occurrences}
				result.Diagnostics = append(result.Diagnostics, diagnostics...)
				if _, exists := settings["id"]; !exists {
					result.Diagnostics = append(result.Diagnostics, blockDiagnostic("webvtt_region_id_missing", "REGION block has no valid id setting", sourceOrder))
				}
			}
			result.DocumentData.Blocks = append(result.DocumentData.Blocks, block)
			if len(result.Diagnostics) > model.MaxDiagnostics {
				return Result{}, &ParseError{Code: "webvtt_diagnostic_limit_exceeded", Line: lines[position].number, Text: fmt.Sprintf("document produces more than %d diagnostics", model.MaxDiagnostics)}
			}
			position = end
			continue
		}

		identifierIndex, timingIndex, timing, cueFound, timingErr := locateCue(lines, position)
		if timingErr != nil {
			return Result{}, timingErr
		}
		if !cueFound {
			end := position + 1
			for end < len(lines) && !blank(lines[end].text) && !cueStartAt(lines, end) {
				end++
			}
			rawLines := lineTexts(lexicalLines[position:end])
			if len(rawLines) > model.MaxItemOccurrences {
				return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[position].number, Text: fmt.Sprintf("block contains more than %d physical lines", model.MaxItemOccurrences)}
			}
			result.DocumentData.Blocks = append(result.DocumentData.Blocks, model.WebVTTBlock{Type: "unrecognized", SourceOrder: sourceOrder, Raw: strings.Join(rawLines, "\n"), RawLines: rawLines})
			result.Diagnostics = append(result.Diagnostics, blockDiagnostic("webvtt_block_unrecognized", fmt.Sprintf("unrecognized WebVTT block beginning at line %d was preserved", lines[position].number), sourceOrder))
			position = end
			continue
		}

		payloadStart := timingIndex + 1
		payloadEnd := payloadStart
		missingSeparator := false
		for payloadEnd < len(lines) && !blank(lines[payloadEnd].text) {
			if payloadEnd > payloadStart && cueStartAt(lines, payloadEnd) {
				missingSeparator = true
				break
			}
			payloadEnd++
		}
		if payloadEnd == payloadStart {
			return Result{}, &ParseError{Code: "webvtt_payload_empty", Line: lines[timingIndex].number, Text: "cue must contain at least one payload line"}
		}
		if payloadEnd-payloadStart > model.MaxItemOccurrences {
			return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[timingIndex].number, Text: fmt.Sprintf("cue contains more than %d payload lines", model.MaxItemOccurrences)}
		}
		ordinal := len(result.Cues)
		cueID := fmt.Sprintf("cue-%06d", ordinal)
		var identifierRaw *string
		if identifierIndex >= 0 {
			identifier := lexicalLines[identifierIndex].text
			identifierRaw = &identifier
			semanticIdentifier := lines[identifierIndex].text
			if identifiers[semanticIdentifier] {
				result.Diagnostics = append(result.Diagnostics, cueDiagnostic("webvtt_cue_identifier_duplicate", "duplicate WebVTT cue identifier was preserved", sourceOrder, cueID))
			}
			identifiers[semanticIdentifier] = true
		}
		if seenCue && timing.StartMilliseconds < previousStart {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("webvtt_start_order_invalid", "decreasing cue start time was preserved in source order", sourceOrder, cueID))
		}
		lexicalTiming, lexicalTimingErr := ParseTimingLine(lexicalLines[timingIndex].text)
		if lexicalTimingErr == nil {
			timing.SettingsRaw = lexicalTiming.SettingsRaw
		}
		if countASCIIWhitespaceFields(timing.SettingsRaw) > model.MaxItemOccurrences {
			return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[timingIndex].number, Text: fmt.Sprintf("cue contains more than %d setting occurrences", model.MaxItemOccurrences)}
		}
		occurrences, settings, settingDiagnostics := parseCueSettings(timing.SettingsRaw, sourceOrder, cueID)
		result.Diagnostics = append(result.Diagnostics, settingDiagnostics...)
		if missingSeparator {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("webvtt_separator_missing", "cue boundary was recovered without a blank separator", sourceOrder, cueID))
		}
		payloadLines := lineTexts(lexicalLines[payloadStart:payloadEnd])
		rawPayload := strings.Join(payloadLines, "\n")
		if strings.Count(rawPayload, "<") > model.MaxItemOccurrences || strings.Count(rawPayload, "&") > model.MaxItemOccurrences {
			return Result{}, &ParseError{Code: "webvtt_occurrence_limit_exceeded", Line: lines[timingIndex].number, Text: fmt.Sprintf("cue contains more than %d markup or entity occurrences", model.MaxItemOccurrences)}
		}
		plain, speakers, tokens, markupDiagnostics := scanPayload(rawPayload, timing.StartMilliseconds, timing.EndMilliseconds, sourceOrder, cueID)
		result.Diagnostics = append(result.Diagnostics, markupDiagnostics...)
		cue := model.Cue{
			ID: cueID, Ordinal: ordinal, SourceOrder: sourceOrder, SourceIdentifier: identifierRaw,
			Timing:   model.Timing{StartMilliseconds: timing.StartMilliseconds, EndMilliseconds: timing.EndMilliseconds, DurationMilliseconds: timing.EndMilliseconds - timing.StartMilliseconds},
			Payload:  model.Payload{RawText: rawPayload, PlainText: plain, Lines: payloadLines},
			Speakers: speakers, Tokens: tokens, OCRObservations: []model.OCRObservation{},
			FormatData: model.CueFormatData{WebVTT: &model.WebVTTCueData{IdentifierRaw: identifierRaw, TimingLineRaw: lexicalLines[timingIndex].text, SettingsRaw: timing.SettingsRaw, Settings: settings, SettingOccurrences: occurrences, RawPayload: rawPayload, RawPayloadLines: payloadLines}},
		}
		cue.Placement = placementFromSettings(settings)
		result.Cues = append(result.Cues, cue)
		if len(result.Diagnostics) > model.MaxDiagnostics {
			return Result{}, &ParseError{Code: "webvtt_diagnostic_limit_exceeded", Line: lines[timingIndex].number, Text: fmt.Sprintf("document produces more than %d diagnostics", model.MaxDiagnostics)}
		}
		seenCue = true
		previousStart = timing.StartMilliseconds
		position = payloadEnd
	}
	if len(result.Cues) == 0 {
		return Result{}, &ParseError{Code: "webvtt_no_cues", Text: "input contains no valid WebVTT cues"}
	}
	return result, nil
}

func parseSignature(line string) (*string, bool) {
	if strings.Contains(line, "-->") {
		return nil, false
	}
	if line == "WEBVTT" {
		return nil, true
	}
	if !strings.HasPrefix(line, "WEBVTT") || len(line) == 6 || !isASCIISpace(line[6]) {
		return nil, false
	}
	description := asciiTrim(line[6:])
	return &description, true
}

func locateCue(lines []physicalLine, position int) (int, int, TimingLine, bool, error) {
	if timing, err := ParseTimingLine(lines[position].text); err == nil {
		return -1, position, timing, true, nil
	} else if strings.Contains(lines[position].text, "-->") {
		return -1, position, TimingLine{}, false, &ParseError{Code: "webvtt_timing_invalid", Line: lines[position].number, Text: err.Error()}
	}
	if position+1 < len(lines) {
		if timing, err := ParseTimingLine(lines[position+1].text); err == nil {
			return position, position + 1, timing, true, nil
		} else if strings.Contains(lines[position+1].text, "-->") {
			return position, position + 1, TimingLine{}, false, &ParseError{Code: "webvtt_timing_invalid", Line: lines[position+1].number, Text: err.Error()}
		}
	}
	return -1, -1, TimingLine{}, false, nil
}

func cueStartAt(lines []physicalLine, position int) bool {
	if position >= len(lines) {
		return false
	}
	if _, err := ParseTimingLine(lines[position].text); err == nil {
		return true
	}
	return position+1 < len(lines) && !blank(lines[position].text) && func() bool { _, err := ParseTimingLine(lines[position+1].text); return err == nil }()
}

func keywordLine(line, keyword string) bool {
	if line == keyword {
		return true
	}
	return strings.HasPrefix(line, keyword) && len(line) > len(keyword) && isASCIISpace(line[len(keyword)]) && asciiTrim(line[len(keyword):]) == ""
}

func noteLine(line string) bool {
	return line == "NOTE" || strings.HasPrefix(line, "NOTE ") || strings.HasPrefix(line, "NOTE\t")
}

func lineTexts(lines []physicalLine) []string {
	values := make([]string, len(lines))
	for index := range lines {
		values[index] = lines[index].text
	}
	return values
}

func placementFromSettings(settings map[string]string) *model.Placement {
	placement := &model.Placement{}
	count := 0
	for name, destination := range map[string]**string{"line": &placement.Line, "position": &placement.Position, "size": &placement.Size, "align": &placement.Align} {
		if value, exists := settings[name]; exists {
			copied := value
			*destination = &copied
			count++
		}
	}
	if count == 0 {
		return nil
	}
	return placement
}

func blockDiagnostic(code, message string, sourceOrder int) Diagnostic {
	order := sourceOrder
	return Diagnostic{Severity: "warning", Code: code, Message: message, SourceOrder: &order}
}

func cueDiagnostic(code, message string, sourceOrder int, cueID string) Diagnostic {
	order, id := sourceOrder, cueID
	return Diagnostic{Severity: "warning", Code: code, Message: message, SourceOrder: &order, CueID: &id}
}
