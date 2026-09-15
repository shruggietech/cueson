package convert

import (
	"context"
	"fmt"
	"html"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

// translateTextToScripted preserves shared emphasis without passing through a
// text-target placeholder policy or reinterpreting literals as native controls.
func translateTextToScripted(document model.Document, cueIndex int) (payloadTranslation, error) {
	if cueIndex < 0 || cueIndex >= len(document.Cues) {
		return payloadTranslation{}, fmt.Errorf("invalid_scripted_cue_index")
	}
	cue := document.Cues[cueIndex]
	raw := cue.Payload.RawText
	if err := validateTranslationComplexity(raw, "scripted"); err != nil {
		return payloadTranslation{}, err
	}
	var tokens []markupToken
	switch document.Format {
	case "subrip":
		tokens = scanSubRipMarkup(raw)
		pairSubRipMarkup(tokens)
	case "webvtt":
		tokens = scanWebVTTMarkup(raw, cue.Timing.StartMilliseconds, cue.Timing.EndMilliseconds)
		pairMarkup(tokens, true)
	default:
		return payloadTranslation{}, fmt.Errorf("invalid_scripted_text_source")
	}
	collector := payloadIssueCollector{issues: []payloadIssue{}}
	occurrences := map[string]int{}
	depths := map[string]int{}
	var output strings.Builder
	writeLiteral := func(value string, decode bool) error {
		if decode && document.Format == "webvtt" {
			value = decodeScriptedWebVTTLiteral(value)
		}
		if strings.ContainsRune(value, '\x00') {
			for _, r := range value {
				if r == 0 {
					appendPayloadIssue(&collector, occurrences, LossCodeNULDegraded, KindDegraded, "nul")
				}
			}
			value = strings.ReplaceAll(value, "\x00", "\ufffd")
		}
		encoded, err := encodeScriptedLiteral(value)
		if err != nil {
			return err
		}
		if output.Len() > model.MaxScriptedLineBytes-len(encoded) {
			return fmt.Errorf("scripted_payload_complexity_limit")
		}
		output.WriteString(encoded)
		return nil
	}
	cursor := 0
	for _, token := range tokens {
		if err := writeLiteral(raw[cursor:token.start], true); err != nil {
			return payloadTranslation{}, err
		}
		switch {
		case token.timestamp:
			appendPayloadIssue(&collector, occurrences, LossCodeWebVTTInlineTimingOmitted, KindOmitted, "inline_timestamp")
		case token.matched && (token.name == "b" || token.name == "i" || token.name == "u"):
			if document.Format == "webvtt" && token.classes && !token.closing {
				appendPayloadIssue(&collector, occurrences, LossCodeWebVTTMarkupDegraded, KindDegraded, "class")
			}
			before := depths[token.name] > 0
			if token.closing {
				depths[token.name]--
			} else {
				depths[token.name]++
			}
			after := depths[token.name] > 0
			if before != after {
				value := "0"
				if after {
					value = "1"
				}
				output.WriteString("{\\" + token.name + value + "}")
			}
		case token.matched && token.known:
			if !token.closing {
				if document.Format == "subrip" && token.name == "font" {
					appendPayloadIssue(&collector, occurrences, LossCodeSubRipFontDegraded, KindDegraded, "font")
				} else if document.Format == "webvtt" {
					code := LossCodeWebVTTMarkupDegraded
					if token.name == "v" {
						code = LossCodeWebVTTVoiceDegraded
					}
					appendPayloadIssue(&collector, occurrences, code, KindDegraded, token.name)
				}
			}
		default:
			if err := writeLiteral(token.raw, false); err != nil {
				return payloadTranslation{}, err
			}
		}
		cursor = token.end
	}
	if err := writeLiteral(raw[cursor:], true); err != nil {
		return payloadTranslation{}, err
	}
	if collector.overflow || output.Len() > model.MaxScriptedLineBytes {
		return payloadTranslation{}, fmt.Errorf("scripted_payload_complexity_limit")
	}
	return payloadTranslation{Text: output.String(), Issues: collector.issues}, nil
}

func encodeScriptedLiteral(text string) (string, error) {
	if !utf8.ValidString(text) || strings.ContainsAny(text, "{}\r\x00") || strings.Contains(text, `\N`) || strings.Contains(text, `\n`) || strings.Contains(text, `\h`) {
		return "", fmt.Errorf("unrepresentable_scripted_literal")
	}
	if len(text) > model.MaxScriptedLineBytes {
		return "", fmt.Errorf("scripted_payload_complexity_limit")
	}
	encoded := strings.NewReplacer("\n", `\N`, "\u00a0", `\h`).Replace(text)
	if len(encoded) > model.MaxScriptedLineBytes {
		return "", fmt.Errorf("scripted_payload_complexity_limit")
	}
	return encoded, nil
}

// Decode only terminated character references, matching the WebVTT codec.
// Scripted text has no SubRip angle-markup ambiguity policy.
func decodeScriptedWebVTTLiteral(text string) string {
	var output strings.Builder
	for position := 0; position < len(text); {
		relative := strings.IndexByte(text[position:], '&')
		if relative < 0 {
			output.WriteString(text[position:])
			break
		}
		start := position + relative
		output.WriteString(text[position:start])
		semicolon := strings.IndexByte(text[start+1:], ';')
		if semicolon < 0 {
			output.WriteString(text[start:])
			break
		}
		end := start + semicolon + 2
		output.WriteString(html.UnescapeString(text[start:end]))
		position = end
	}
	return output.String()
}

func quantizeScriptedTiming(timing model.Timing) (model.Timing, error) {
	round := func(value int64) (int64, error) {
		if value < 0 {
			return 0, fmt.Errorf("invalid_scripted_timing")
		}
		quotient := value / 10
		if value%10 >= 5 {
			quotient++
		}
		if quotient > math.MaxInt64/10 {
			return 0, fmt.Errorf("scripted_timing_overflow")
		}
		return quotient * 10, nil
	}
	start, err := round(timing.StartMilliseconds)
	if err != nil {
		return model.Timing{}, err
	}
	end, err := round(timing.EndMilliseconds)
	if err != nil {
		return model.Timing{}, err
	}
	if end <= start {
		return model.Timing{}, fmt.Errorf("collapsed_scripted_interval")
	}
	return model.Timing{StartMilliseconds: start, EndMilliseconds: end, DurationMilliseconds: end - start}, nil
}

func projectScriptedTarget(ctx context.Context, source model.Document, targetFormat string, translations []payloadTranslation) (model.Document, []Loss, error) {
	if ctx == nil {
		return model.Document{}, nil, fmt.Errorf("invalid_context")
	}
	if (targetFormat != "ass" && targetFormat != "ssa") || len(translations) != len(source.Cues) || len(source.Cues) > model.MaxDocumentItems-8 {
		return model.Document{}, nil, fmt.Errorf("invalid_scripted_target")
	}
	target := model.Document{Schema: schema.ID(), SchemaVersion: schema.Version(), Format: targetFormat, FormatSupport: model.FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}, Producer: source.Producer, Source: cloneSourceEnvelope(source.Source), Metadata: model.Metadata{}, Cues: []model.Cue{}, Diagnostics: []model.Diagnostic{}}
	native := &model.ScriptedDocumentData{Dialect: "v4.00+", Sections: []model.ScriptedSection{}, Records: []model.ScriptedRecord{}, Styles: []model.ScriptedStyle{}, Events: []model.ScriptedEvent{}, Attachments: []model.ScriptedAttachment{}}
	styleSection := "V4+ Styles"
	if targetFormat == "ssa" {
		native.Dialect, styleSection = "v4.00", "V4 Styles"
		target.FormatData.SSA = native
	} else {
		target.FormatData.ASS = native
	}
	native.Sections = []model.ScriptedSection{{SectionID: "section-000000", SourceOrder: 0, Name: "Script Info"}, {SectionID: "section-000003", SourceOrder: 3, Name: styleSection}, {SectionID: "section-000006", SourceOrder: 6, Name: "Events"}}
	styleID, styleDeclaration, eventDeclaration := "style-000005", "record-000004", "record-000007"
	native.Records = []model.ScriptedRecord{{RecordID: "record-000001", SourceOrder: 1, SectionID: "section-000000", Kind: "metadata", Fields: []model.ScriptedField{{FieldName: "ScriptType", RawValue: native.Dialect}}}, {RecordID: "record-000002", SourceOrder: 2, SectionID: "section-000000", Kind: "metadata", Fields: []model.ScriptedField{{FieldName: "WrapStyle", RawValue: "0"}}}, {RecordID: styleDeclaration, SourceOrder: 4, SectionID: "section-000003", Kind: "format_declaration", DeclarationFields: scriptedDeclaration(targetFormat, "style", nil)}, {RecordID: "record-000005", SourceOrder: 5, SectionID: "section-000003", Kind: "style", StyleID: &styleID}, {RecordID: eventDeclaration, SourceOrder: 7, SectionID: "section-000006", Kind: "format_declaration", DeclarationFields: scriptedDeclaration(targetFormat, "event", nil)}}
	native.Styles = []model.ScriptedStyle{{StyleID: styleID, RecordID: "record-000005", DeclarationID: &styleDeclaration, Name: "Default", Valid: true, Fields: scriptedDefaultStyle(targetFormat)}}
	losses := []Loss{}
	for index, sourceCue := range source.Cues {
		if err := ctx.Err(); err != nil {
			return model.Document{}, nil, err
		}
		timing, err := quantizeScriptedTiming(sourceCue.Timing)
		if err != nil {
			return model.Document{}, nil, projectionFailure(source.Format, targetFormat, fmt.Sprintf("/cues/%d/timing", index), err)
		}
		if err := appendScriptedPrecisionLosses(source, targetFormat, index, timing, &losses); err != nil {
			return model.Document{}, nil, err
		}
		text := translations[index].Text
		facts, err := model.ProjectScriptedText(text, 0, timing.StartMilliseconds, timing.EndMilliseconds)
		if err != nil || len(facts.DiagnosticCodes) != 0 || facts.DrawingExcluded {
			return model.Document{}, nil, fmt.Errorf("invalid_scripted_target_projection")
		}
		id, eventID, recordID := fmt.Sprintf("cue-%06d", index), fmt.Sprintf("event-%06d", index+8), fmt.Sprintf("record-%06d", index+8)
		fields := scriptedDefaultEvent(targetFormat, timing, text)
		native.Events = append(native.Events, model.ScriptedEvent{EventID: eventID, RecordID: recordID, DeclarationID: &eventDeclaration, EventType: "dialogue", CueID: &id, Valid: true, Fields: fields, Text: text, Spans: facts.Spans, Tags: facts.Tags, Karaoke: facts.Karaoke})
		native.Records = append(native.Records, model.ScriptedRecord{RecordID: recordID, SourceOrder: index + 8, SectionID: "section-000006", Kind: "event", EventID: &eventID})
		data := &model.ScriptedCueData{EventID: eventID, StyleID: styleID, Projection: model.ScriptedProjection{Origin: "native_derived", SourceEventID: eventID, TextFieldName: "Text", WrapStyle: 0}}
		cue := model.Cue{ID: id, Ordinal: index, SourceOrder: index + 8, Timing: timing, Payload: model.Payload{RawText: text, PlainText: strings.Join(facts.Lines, "\n"), Lines: facts.Lines}, Speakers: []model.Speaker{}, Tokens: facts.Tokens, OCRObservations: []model.OCRObservation{}}
		if targetFormat == "ass" {
			cue.FormatData.ASS = data
		} else {
			cue.FormatData.SSA = data
		}
		target.Cues = append(target.Cues, cue)
	}
	if len(target.Cues) != 0 {
		recomputeSummaries(&target)
	}
	if err := model.ValidateScriptedTarget(source, target); err != nil {
		return model.Document{}, nil, err
	}
	return target, losses, nil
}

func scriptedDeclaration(format, kind string, extras []string) []model.ScriptedDeclarationField {
	names := model.ScriptedCanonicalFields(format, kind)
	if kind == "event" {
		names = append(append(names[:len(names)-1], extras...), "Text")
	} else {
		names = append(names, extras...)
	}
	fields := make([]model.ScriptedDeclarationField, len(names))
	for index, name := range names {
		identity, known := model.ScriptedFieldIdentity(name, format, kind)
		fields[index].FieldName = name
		if known {
			fields[index].RecognizedName = &identity
		}
	}
	return fields
}

func scriptedDefaultStyle(format string) []model.ScriptedField {
	values := map[string]string{"Name": "Default", "Fontname": "Arial", "Fontsize": "20", "PrimaryColour": "&H00FFFFFF", "SecondaryColour": "&H000000FF", "OutlineColour": "&H00000000", "TertiaryColour": "0", "BackColour": "&H00000000", "Bold": "0", "Italic": "0", "Underline": "0", "StrikeOut": "0", "ScaleX": "100", "ScaleY": "100", "Spacing": "0", "Angle": "0", "BorderStyle": "1", "Outline": "2", "Shadow": "0", "Alignment": "2", "MarginL": "10", "MarginR": "10", "MarginV": "10", "AlphaLevel": "0", "Encoding": "1"}
	if format == "ssa" {
		values["PrimaryColour"], values["SecondaryColour"], values["BackColour"] = "16777215", "255", "0"
	}
	fields := []model.ScriptedField{}
	for _, name := range model.ScriptedCanonicalFields(format, "style") {
		fields = append(fields, model.ScriptedField{FieldName: name, RawValue: values[name]})
	}
	return fields
}

func scriptedDefaultEvent(format string, timing model.Timing, text string) []model.ScriptedField {
	values := map[string]string{"Layer": "0", "Marked": "Marked=0", "Start": scriptedTargetTimestamp(timing.StartMilliseconds), "End": scriptedTargetTimestamp(timing.EndMilliseconds), "Style": "Default", "Name": "", "MarginL": "0", "MarginR": "0", "MarginV": "0", "Effect": "", "Text": text}
	fields := []model.ScriptedField{}
	for _, name := range model.ScriptedCanonicalFields(format, "event") {
		fields = append(fields, model.ScriptedField{FieldName: name, RawValue: values[name]})
	}
	return fields
}

func scriptedTargetTimestamp(value int64) string {
	return strconv.FormatInt(value/3_600_000, 10) + fmt.Sprintf(":%02d:%02d.%02d", value/60_000%60, value/1_000%60, value/10%100)
}

func appendScriptedPrecisionLosses(source model.Document, targetFormat string, index int, timing model.Timing, losses *[]Loss) error {
	cue := source.Cues[index]
	for _, endpoint := range []struct {
		name          string
		before, after int64
	}{{"start_milliseconds", cue.Timing.StartMilliseconds, timing.StartMilliseconds}, {"end_milliseconds", cue.Timing.EndMilliseconds, timing.EndMilliseconds}} {
		if endpoint.before != endpoint.after {
			loss := cueLoss(source, targetFormat, cue, LossCodeScriptedCentisecondQuantized, KindDegraded, "cue timing is rounded to the nearest centisecond", fmt.Sprintf("/cues/%d/timing/%s", index, endpoint.name), []Attribute{{Name: "endpoint", Value: endpoint.name}})
			if err := appendOneLoss(losses, loss); err != nil {
				return err
			}
		}
	}
	return nil
}
