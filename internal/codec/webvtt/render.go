package webvtt

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/model"
)

type renderItem struct {
	order int
	text  string
}

// Render serializes a WebVTT Cue JSON document as canonical LF WebVTT.
func Render(document model.Document, options RenderOptions) (RenderResult, error) {
	if document.Format != "webvtt" || document.FormatData.WebVTT == nil {
		return RenderResult{}, fmt.Errorf("render WebVTT: document format_data is not WebVTT")
	}
	if len(document.Cues) == 0 {
		return RenderResult{}, fmt.Errorf("render WebVTT: at least one cue is required")
	}
	knownNonconforming := false
	for _, diagnostic := range document.Diagnostics {
		if strings.HasPrefix(diagnostic.Code, "webvtt_") && diagnostic.Code != "webvtt_nul_replaced" {
			knownNonconforming = true
			if options.Strict {
				return RenderResult{}, fmt.Errorf("render WebVTT strict: preserved diagnostic %s prevents conforming output", diagnostic.Code)
			}
		}
	}
	data := document.FormatData.WebVTT
	header := "WEBVTT"
	if data.Description != nil {
		if strings.ContainsAny(*data.Description, "\r\n") {
			return RenderResult{}, fmt.Errorf("render WebVTT: description contains a line terminator")
		}
		header += " " + semanticText(*data.Description)
	}
	headerLines := []string{header}
	for index, line := range data.MetadataLines {
		if strings.ContainsAny(line, "\r\n") || strings.Contains(line, "-->") {
			return RenderResult{}, fmt.Errorf("render WebVTT: metadata line %d is unsafe", index)
		}
		headerLines = append(headerLines, semanticText(line))
	}

	firstCueOrder := len(document.Cues) + len(data.Blocks)
	for index := range document.Cues {
		if document.Cues[index].SourceOrder < firstCueOrder {
			firstCueOrder = document.Cues[index].SourceOrder
		}
	}
	items := make([]renderItem, 0, len(document.Cues)+len(data.Blocks))
	orders := make(map[int]bool, cap(items))
	for index := range data.Blocks {
		block := &data.Blocks[index]
		if block.SourceOrder < 0 || orders[block.SourceOrder] {
			return RenderResult{}, fmt.Errorf("render WebVTT: duplicate or negative source_order %d", block.SourceOrder)
		}
		orders[block.SourceOrder] = true
		if block.Raw != strings.Join(block.RawLines, "\n") || strings.ContainsRune(block.Raw, '\r') {
			return RenderResult{}, fmt.Errorf("render WebVTT block %d: raw and raw_lines are inconsistent", index)
		}
		if strings.Contains(block.Raw, "\n\n") {
			return RenderResult{}, fmt.Errorf("render WebVTT block %d: embedded blank line is unsafe", index)
		}
		blockNonconforming := block.Type == "unrecognized" || ((block.Type == "style" || block.Type == "region") && block.SourceOrder > firstCueOrder)
		if (block.Type == "note" || block.Type == "style") && strings.Contains(block.Raw, "-->") {
			blockNonconforming = true
		}
		text := semanticText(block.Raw)
		if block.Type == "region" {
			if block.Region == nil {
				return RenderResult{}, fmt.Errorf("render WebVTT block %d: REGION data is missing", index)
			}
			regionText, regionNonconforming, err := renderRegionBlock(*block)
			if err != nil {
				return RenderResult{}, fmt.Errorf("render WebVTT block %d: %w", index, err)
			}
			text = regionText
			blockNonconforming = blockNonconforming || regionNonconforming
		}
		if blockNonconforming {
			knownNonconforming = true
			if options.Strict {
				return RenderResult{}, fmt.Errorf("render WebVTT strict: block %d contains nonconforming native content", index)
			}
		}
		items = append(items, renderItem{order: block.SourceOrder, text: text})
	}

	seenIdentifiers := make(map[string]bool)
	var previousStart int64
	for index := range document.Cues {
		cue := &document.Cues[index]
		if cue.SourceOrder < 0 || orders[cue.SourceOrder] {
			return RenderResult{}, fmt.Errorf("render WebVTT: duplicate or negative source_order %d", cue.SourceOrder)
		}
		orders[cue.SourceOrder] = true
		if cue.FormatData.WebVTT == nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: native data is missing", index)
		}
		if cue.Timing.StartMilliseconds < 0 || cue.Timing.EndMilliseconds <= cue.Timing.StartMilliseconds || cue.Timing.DurationMilliseconds != cue.Timing.EndMilliseconds-cue.Timing.StartMilliseconds {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: invalid timing", index)
		}
		native := cue.FormatData.WebVTT
		if native.RawPayload != cue.Payload.RawText || native.RawPayload != strings.Join(native.RawPayloadLines, "\n") || !reflect.DeepEqual(native.RawPayloadLines, cue.Payload.Lines) {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: payload lexical and structured fields disagree", index)
		}
		if len(native.RawPayloadLines) == 0 {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: at least one payload line is required", index)
		}
		for lineIndex, line := range native.RawPayloadLines {
			if line == "" || strings.ContainsAny(line, "\r\n") {
				return RenderResult{}, fmt.Errorf("render WebVTT cue %d: payload line %d is an unsafe blank or non-physical line", index, lineIndex)
			}
		}
		start, err := formatTimestamp(cue.Timing.StartMilliseconds)
		if err != nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d start: %w", index, err)
		}
		end, err := formatTimestamp(cue.Timing.EndMilliseconds)
		if err != nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d end: %w", index, err)
		}
		settingsSuffix, settingsNonconforming, err := renderCueSettingsSuffix(*native, cue.SourceOrder, cue.ID)
		if err != nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: %w", index, err)
		}
		_, _, _, markupDiagnostics := scanPayload(native.RawPayload, cue.Timing.StartMilliseconds, cue.Timing.EndMilliseconds, cue.SourceOrder, cue.ID)
		cueNonconforming := settingsNonconforming || len(markupDiagnostics) > 0 || index > 0 && cue.Timing.StartMilliseconds < previousStart
		if native.IdentifierRaw != nil {
			semanticIdentifier := semanticText(*native.IdentifierRaw)
			if strings.ContainsAny(semanticIdentifier, "\r\n") || strings.Contains(semanticIdentifier, "-->") {
				return RenderResult{}, fmt.Errorf("render WebVTT cue %d: identifier is unsafe", index)
			}
			if seenIdentifiers[semanticIdentifier] {
				cueNonconforming = true
			}
			seenIdentifiers[semanticIdentifier] = true
		}
		if cueNonconforming {
			knownNonconforming = true
			if options.Strict {
				return RenderResult{}, fmt.Errorf("render WebVTT strict: cue %d contains nonconforming native content", index)
			}
		}
		lines := make([]string, 0, 3+len(native.RawPayloadLines))
		if native.IdentifierRaw != nil {
			lines = append(lines, semanticText(*native.IdentifierRaw))
		}
		lines = append(lines, start+" --> "+end+settingsSuffix)
		for _, line := range native.RawPayloadLines {
			lines = append(lines, semanticText(line))
		}
		items = append(items, renderItem{order: cue.SourceOrder, text: strings.Join(lines, "\n")})
		previousStart = cue.Timing.StartMilliseconds
	}
	for expected := 0; expected < len(items); expected++ {
		if !orders[expected] {
			return RenderResult{}, fmt.Errorf("render WebVTT: source_order is not contiguous at %d", expected)
		}
	}
	sort.Slice(items, func(left, right int) bool { return items[left].order < items[right].order })
	body := make([]string, len(items))
	for index := range items {
		body[index] = items[index].text
	}
	output := strings.Join(headerLines, "\n") + "\n\n" + strings.Join(body, "\n\n") + "\n"
	result := RenderResult{Bytes: []byte(output), Diagnostics: []Diagnostic{}}
	if knownNonconforming {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: "warning", Code: "webvtt_render_preserved_nonconforming", Message: "permissive rendering retained diagnosed nonconforming WebVTT content"})
	}
	return result, nil
}

func renderRegionBlock(block model.WebVTTBlock) (string, bool, error) {
	region := block.Region
	if region == nil {
		return "", false, fmt.Errorf("REGION data is missing")
	}
	if err := validateRawSettingOccurrences(region.SettingOccurrences, regionSettingOrder, validRegionSetting); err != nil {
		return "", false, err
	}
	occurrences, effective, _ := parseRegionSettings(region.SettingsRaw, block.SourceOrder)
	rawSettings := ""
	if len(block.RawLines) > 1 {
		rawSettings = strings.Join(block.RawLines[1:], "\n")
	}
	lexicalConsistent := rawSettings == region.SettingsRaw && reflect.DeepEqual(occurrences, region.SettingOccurrences) && reflect.DeepEqual(effective, region.Settings)
	nonconforming := settingOccurrencesNonconforming(region.SettingOccurrences)
	if _, exists := region.Settings["id"]; !exists {
		nonconforming = true
	}
	if lexicalConsistent {
		return semanticText(block.Raw), nonconforming, nil
	}
	settings, err := canonicalSettings(regionSettingOrder, region.Settings, region.SettingOccurrences, validRegionSetting)
	if err != nil {
		return "", false, err
	}
	if settings == "" {
		return "REGION", nonconforming, nil
	}
	return "REGION\n" + settings, nonconforming, nil
}

// renderCueSettingsSuffix returns the entire timing-line suffix, including its
// leading delimiter whitespace when settings are present.
func renderCueSettingsSuffix(native model.WebVTTCueData, sourceOrder int, cueID string) (string, bool, error) {
	if strings.ContainsAny(native.SettingsRaw, "\r\n") {
		return "", false, fmt.Errorf("settings_raw contains a line terminator")
	}
	if err := validateRawSettingOccurrences(native.SettingOccurrences, cueSettingOrder, validCueSetting); err != nil {
		return "", false, err
	}
	occurrences, effective, _ := parseCueSettings(native.SettingsRaw, sourceOrder, cueID)
	nonconforming := settingOccurrencesNonconforming(native.SettingOccurrences)
	if reflect.DeepEqual(occurrences, native.SettingOccurrences) && reflect.DeepEqual(effective, native.Settings) {
		if native.SettingsRaw == "" || isASCIISpace(native.SettingsRaw[0]) {
			return semanticText(native.SettingsRaw), nonconforming, nil
		}
		return " " + semanticText(native.SettingsRaw), nonconforming, nil
	}
	settings, err := canonicalSettings(cueSettingOrder, native.Settings, native.SettingOccurrences, validCueSetting)
	if err != nil {
		return "", false, err
	}
	if settings == "" {
		return "", nonconforming, nil
	}
	return " " + settings, nonconforming, nil
}

func canonicalSettings(order []string, effective map[string]string, occurrences []model.WebVTTSettingOccurrence, validate func(string, string) bool) (string, error) {
	parts := make([]string, 0, len(effective)+len(occurrences))
	known := make(map[string]bool, len(order))
	for _, name := range order {
		known[name] = true
		if value, exists := effective[name]; exists {
			value = semanticText(value)
			if !validate(name, value) {
				return "", fmt.Errorf("effective setting %s is invalid", name)
			}
			parts = append(parts, name+":"+value)
		}
	}
	for name := range effective {
		if !known[name] {
			return "", fmt.Errorf("effective setting %s is unknown", name)
		}
	}
	for _, occurrence := range occurrences {
		if !occurrence.Recognized || !occurrence.Valid {
			parts = append(parts, semanticText(occurrence.Raw))
		}
	}
	return strings.Join(parts, " "), nil
}

func validateRawSettingOccurrences(occurrences []model.WebVTTSettingOccurrence, recognizedNames []string, validate func(string, string) bool) error {
	recognized := make(map[string]bool, len(recognizedNames))
	for _, name := range recognizedNames {
		recognized[name] = true
	}
	for index, occurrence := range occurrences {
		if occurrence.Raw == "" || strings.Contains(occurrence.Raw, "-->") {
			return fmt.Errorf("setting occurrence %d has unsafe raw syntax", index)
		}
		for _, character := range occurrence.Raw {
			if character != '\x00' && (character <= ' ' || character == '\x7f') {
				return fmt.Errorf("setting occurrence %d has unsafe raw syntax", index)
			}
		}
		name, value, found := strings.Cut(occurrence.Raw, ":")
		if !found {
			value = ""
		}
		if name != occurrence.Name || value != occurrence.Value {
			return fmt.Errorf("setting occurrence %d raw syntax disagrees with name and value", index)
		}
		semanticName, semanticValue := semanticText(name), semanticText(value)
		wantRecognized := recognized[semanticName]
		wantValid := found && semanticName != "" && semanticValue != "" && wantRecognized && validate(semanticName, semanticValue)
		if occurrence.Recognized != wantRecognized || occurrence.Valid != wantValid {
			return fmt.Errorf("setting occurrence %d has forged recognition or validity flags", index)
		}
		if !utf8.ValidString(occurrence.Raw) {
			return fmt.Errorf("setting occurrence %d is not valid UTF-8", index)
		}
	}
	return nil
}

func settingOccurrencesNonconforming(occurrences []model.WebVTTSettingOccurrence) bool {
	seen := make(map[string]bool)
	for _, occurrence := range occurrences {
		if !occurrence.Recognized || !occurrence.Valid || seen[occurrence.Name] {
			return true
		}
		seen[occurrence.Name] = true
	}
	return false
}
