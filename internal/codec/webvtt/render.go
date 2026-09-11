package webvtt

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

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
		header += " " + *data.Description
	}
	headerLines := []string{header}
	for index, line := range data.MetadataLines {
		if strings.ContainsAny(line, "\r\n") || strings.Contains(line, "-->") {
			return RenderResult{}, fmt.Errorf("render WebVTT: metadata line %d is unsafe", index)
		}
		headerLines = append(headerLines, line)
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
		if block.Type == "unrecognized" {
			knownNonconforming = true
			if options.Strict {
				return RenderResult{}, fmt.Errorf("render WebVTT strict: unrecognized block %d", index)
			}
		}
		if block.Type == "region" {
			if block.Region == nil {
				return RenderResult{}, fmt.Errorf("render WebVTT block %d: REGION data is missing", index)
			}
			occurrences, effective, _ := parseRegionSettings(block.Region.SettingsRaw, block.SourceOrder)
			if !reflect.DeepEqual(occurrences, block.Region.SettingOccurrences) || !reflect.DeepEqual(effective, block.Region.Settings) {
				return RenderResult{}, fmt.Errorf("render WebVTT block %d: REGION lexical and structured settings disagree", index)
			}
		}
		items = append(items, renderItem{order: block.SourceOrder, text: block.Raw})
	}
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
		if strings.ContainsRune(native.RawPayload, '\r') || strings.Contains(native.RawPayload, "\n\n") {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: payload contains an unsafe separator", index)
		}
		start, err := formatTimestamp(cue.Timing.StartMilliseconds)
		if err != nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d start: %w", index, err)
		}
		end, err := formatTimestamp(cue.Timing.EndMilliseconds)
		if err != nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d end: %w", index, err)
		}
		settings, err := renderCueSettings(*native, cue.SourceOrder, cue.ID)
		if err != nil {
			return RenderResult{}, fmt.Errorf("render WebVTT cue %d: %w", index, err)
		}
		lines := make([]string, 0, 3+len(native.RawPayloadLines))
		if native.IdentifierRaw != nil {
			if strings.ContainsAny(*native.IdentifierRaw, "\r\n") || strings.Contains(*native.IdentifierRaw, "-->") {
				return RenderResult{}, fmt.Errorf("render WebVTT cue %d: identifier is unsafe", index)
			}
			lines = append(lines, *native.IdentifierRaw)
		}
		timingLine := start + " --> " + end
		if settings != "" {
			timingLine += " " + settings
		}
		lines = append(lines, timingLine)
		lines = append(lines, native.RawPayloadLines...)
		items = append(items, renderItem{order: cue.SourceOrder, text: strings.Join(lines, "\n")})
	}
	if len(orders) != len(items) {
		return RenderResult{}, fmt.Errorf("render WebVTT: source order is incomplete")
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

func renderCueSettings(native model.WebVTTCueData, sourceOrder int, cueID string) (string, error) {
	if strings.ContainsAny(native.SettingsRaw, "\r\n") {
		return "", fmt.Errorf("settings_raw contains a line terminator")
	}
	occurrences, effective, _ := parseCueSettings(native.SettingsRaw, sourceOrder, cueID)
	if reflect.DeepEqual(occurrences, native.SettingOccurrences) && reflect.DeepEqual(effective, native.Settings) {
		return native.SettingsRaw, nil
	}
	parts := make([]string, 0, len(native.Settings)+len(native.SettingOccurrences))
	for _, name := range cueSettingOrder {
		if value, exists := native.Settings[name]; exists {
			if !validCueSetting(name, value) {
				return "", fmt.Errorf("effective setting %s is invalid", name)
			}
			parts = append(parts, name+":"+value)
		}
	}
	for _, occurrence := range native.SettingOccurrences {
		if !occurrence.Recognized || !occurrence.Valid {
			parts = append(parts, occurrence.Raw)
		}
	}
	return strings.Join(parts, " "), nil
}
