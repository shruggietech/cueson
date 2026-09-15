package convert

import (
	"context"
	"fmt"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

// analyzeScriptedSource accounts for native owners before any target is rendered.
func analyzeScriptedSource(ctx context.Context, document model.Document, targetFormat string) (matrixAnalysis, error) {
	if ctx == nil {
		return matrixAnalysis{}, fmt.Errorf("invalid_conversion_context")
	}
	if err := ctx.Err(); err != nil {
		return matrixAnalysis{}, err
	}
	if (document.Format != "ass" && document.Format != "ssa") || (targetFormat != "subrip" && targetFormat != "webvtt") {
		return matrixAnalysis{}, &UnsupportedPairError{SourceFormat: document.Format, TargetFormat: targetFormat}
	}
	if len(document.Cues) == 0 {
		return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, "/cues", fmt.Errorf("zero_dialogue_text_target"))
	}
	for _, diagnostic := range document.Diagnostics {
		if diagnostic.Code == "malformed_attachment" {
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, "/format_data/"+document.Format+"/attachments", fmt.Errorf("malformed_scripted_attachment"))
		}
	}
	owners, err := indexedScriptedOwners(document)
	if err != nil {
		return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, "/format_data", err)
	}
	analysis := matrixAnalysis{Translations: make([]payloadTranslation, len(document.Cues)), Losses: []Loss{}}
	if err := appendMetadataLosses(&analysis.Losses, document, targetFormat); err != nil {
		return matrixAnalysis{}, err
	}
	base := "/format_data/" + document.Format
	appendNative := func(code string, path string, order int) error {
		return appendOneLoss(&analysis.Losses, newLoss(document, targetFormat, code, KindOmitted, scriptedLossMessage(code), path, &order, nil, nil))
	}
	malformedOrders := map[int]bool{}
	for _, diagnostic := range document.Diagnostics {
		if diagnostic.SourceOrder != nil && diagnostic.Code == "malformed_override" {
			malformedOrders[*diagnostic.SourceOrder] = true
		}
	}
	for index, section := range owners.native.Sections {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		name := strings.ToLower(section.Name)
		if name != "script info" && name != "v4+ styles" && name != "v4 styles" && name != "events" && name != "fonts" && name != "graphics" {
			if err := appendNative(LossCodeScriptedSectionOmitted, fmt.Sprintf("%s/sections/%d", base, index), section.SourceOrder); err != nil {
				return matrixAnalysis{}, err
			}
		}
	}
	for index, record := range owners.native.Records {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		switch record.Kind {
		case "malformed":
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("%s/records/%d", base, index), fmt.Errorf("malformed_scripted_record"))
		case "metadata":
			for fieldIndex, field := range record.Fields {
				if strings.EqualFold(strings.TrimSpace(field.FieldName), "ScriptType") {
					continue
				}
				if err := appendNative(LossCodeScriptedMetadataOmitted, fmt.Sprintf("%s/records/%d/fields/%d", base, index, fieldIndex), record.SourceOrder); err != nil {
					return matrixAnalysis{}, err
				}
			}
		case "blank", "comment", "unknown":
			if err := appendNative(LossCodeScriptedRecordOmitted, fmt.Sprintf("%s/records/%d", base, index), record.SourceOrder); err != nil {
				return matrixAnalysis{}, err
			}
		}
	}
	for index, style := range owners.native.Styles {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		if !style.Valid {
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("%s/styles/%d", base, index), fmt.Errorf("malformed_scripted_style"))
		}
	}
	for index, event := range owners.native.Events {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		if !event.Valid {
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("%s/events/%d", base, index), fmt.Errorf("malformed_scripted_event"))
		}
		if event.EventType != "dialogue" {
			if err := appendNative(LossCodeScriptedRecordOmitted, fmt.Sprintf("%s/events/%d", base, index), owners.recordOrders[event.RecordID]); err != nil {
				return matrixAnalysis{}, err
			}
		}
	}
	for index, attachment := range owners.native.Attachments {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		if err := appendNative(LossCodeScriptedAttachmentOmitted, fmt.Sprintf("%s/attachments/%d", base, index), owners.recordOrders[attachment.HeaderRecordID]); err != nil {
			return matrixAnalysis{}, err
		}
	}
	for cueIndex, cue := range document.Cues {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		if malformedOrders[cue.SourceOrder] {
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/payload/raw_text", cueIndex), fmt.Errorf("malformed_scripted_override"))
		}
		translation, nativeLosses, err := translateIndexedScriptedPayload(document, cueIndex, targetFormat, owners)
		if err != nil {
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/payload/raw_text", cueIndex), err)
		}
		analysis.Translations[cueIndex] = translation
		nativeCue := cue.FormatData.ASS
		if document.Format == "ssa" {
			nativeCue = cue.FormatData.SSA
		}
		eventIndex := owners.events[nativeCue.EventID]
		for fieldIndex, field := range owners.native.Events[eventIndex].Fields {
			name, _ := model.ScriptedFieldIdentity(field.FieldName, document.Format, "event")
			if name == "start" || name == "end" || name == "text" {
				continue
			}
			if err := appendOneLoss(&analysis.Losses, cueLoss(document, targetFormat, cue, LossCodeScriptedEventFieldOmitted, KindOmitted, scriptedLossMessage(LossCodeScriptedEventFieldOmitted), fmt.Sprintf("%s/events/%d/fields/%d", base, eventIndex, fieldIndex), nil)); err != nil {
				return matrixAnalysis{}, err
			}
		}
		for _, loss := range nativeLosses {
			if err := appendOneLoss(&analysis.Losses, cueLoss(document, targetFormat, cue, loss.code, loss.kind, scriptedLossMessage(loss.code), loss.path, nil)); err != nil {
				return matrixAnalysis{}, err
			}
		}
		if err := appendCommonCueLosses(&analysis.Losses, document, targetFormat, cueIndex, &cue, false); err != nil {
			return matrixAnalysis{}, err
		}
		if err := appendPayloadLosses(&analysis.Losses, document, targetFormat, cueIndex, &cue, translation); err != nil {
			return matrixAnalysis{}, err
		}
	}
	// The single translation pass has now recorded surviving dimensions across
	// all cues and initial or reset styles. Canonical reports restore native
	// source order after this deferred, bounded style-field accounting.
	for index, style := range owners.native.Styles {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		matched := owners.represented[index]
		for fieldIndex, field := range style.Fields {
			name := strings.ToLower(strings.TrimSpace(field.FieldName))
			if (name == "bold" && matched.bold) || (name == "italic" && matched.italic) || (name == "underline" && matched.underline) {
				continue
			}
			if err := appendNative(LossCodeScriptedStyleFieldOmitted, fmt.Sprintf("%s/styles/%d/fields/%d", base, index, fieldIndex), owners.recordOrders[style.RecordID]); err != nil {
				return matrixAnalysis{}, err
			}
		}
	}
	return analysis, nil
}

func scriptedLossMessage(code string) string {
	switch code {
	case LossCodeScriptedMetadataOmitted:
		return "scripted metadata field is omitted"
	case LossCodeScriptedStyleFieldOmitted:
		return "scripted style field is omitted"
	case LossCodeScriptedEventFieldOmitted:
		return "scripted event field is omitted"
	case LossCodeScriptedRecordOmitted:
		return "scripted non-dialogue record is omitted"
	case LossCodeScriptedSectionOmitted:
		return "scripted extension section is omitted"
	case LossCodeScriptedAttachmentOmitted:
		return "scripted embedded attachment is omitted"
	case LossCodeScriptedOverrideOmitted:
		return "scripted override semantic is omitted"
	case LossCodeScriptedOverrideDegraded:
		return "scripted override presentation is weakened"
	case LossCodeScriptedOverrideCommentOmitted:
		return "scripted inline override comment is omitted"
	case LossCodeScriptedDrawingOmitted:
		return "scripted drawing span is omitted"
	default:
		return "scripted source semantic is changed"
	}
}
