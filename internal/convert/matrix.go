package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

var webVTTSettingNames = []string{"region", "vertical", "line", "position", "size", "align"}

type matrixAnalysis struct {
	Translations []payloadTranslation
	Losses       []Loss
}

func analyzeCompatibility(document model.Document, targetFormat string) (matrixAnalysis, error) {
	analysis := matrixAnalysis{Translations: make([]payloadTranslation, len(document.Cues)), Losses: []Loss{}}
	if err := appendMetadataLosses(&analysis.Losses, document, targetFormat); err != nil {
		return matrixAnalysis{}, err
	}
	switch document.Format {
	case "subrip":
		if targetFormat != "webvtt" {
			return matrixAnalysis{}, &UnsupportedPairError{SourceFormat: document.Format, TargetFormat: targetFormat}
		}
		for index, diagnostic := range document.Diagnostics {
			if diagnostic.Code != "subrip_block_unrecognized" {
				continue
			}
			if err := appendOneLoss(&analysis.Losses, newLoss(document, targetFormat, LossCodeSubRipUnrecognizedBlockOmitted, KindOmitted, "unrecognized SubRip block is omitted", fmt.Sprintf("/diagnostics/%d", index), diagnostic.SourceOrder, nil, []Attribute{{Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
				return matrixAnalysis{}, err
			}
		}
		for index := range document.Cues {
			cue := &document.Cues[index]
			if cue.Timing.EndMilliseconds == cue.Timing.StartMilliseconds {
				return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/timing", index), fmt.Errorf("WebVTT requires positive cue duration"))
			}
			translation, err := translateSubRipPayload(cue.Payload.RawText)
			if err != nil {
				return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/payload/raw_text", index), err)
			}
			analysis.Translations[index] = translation
			if err := appendSubRipCueLosses(&analysis.Losses, document, targetFormat, index, cue, translation); err != nil {
				return matrixAnalysis{}, err
			}
		}
	case "webvtt":
		if targetFormat != "subrip" {
			return matrixAnalysis{}, &UnsupportedPairError{SourceFormat: document.Format, TargetFormat: targetFormat}
		}
		if err := appendWebVTTDocumentLosses(&analysis.Losses, document, targetFormat); err != nil {
			return matrixAnalysis{}, err
		}
		for index := range document.Cues {
			cue := &document.Cues[index]
			translation, err := translateWebVTTPayload(cue.Payload.RawText, cue.Timing.StartMilliseconds, cue.Timing.EndMilliseconds)
			if err != nil {
				return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/payload/raw_text", index), err)
			}
			if subRipPayloadAmbiguous(translation.Text) {
				return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/payload/raw_text", index), fmt.Errorf("payload cannot retain its cue boundary in SubRip"))
			}
			analysis.Translations[index] = translation
			if err := appendWebVTTCueLosses(&analysis.Losses, document, targetFormat, index, cue, translation); err != nil {
				return matrixAnalysis{}, err
			}
		}
	default:
		return matrixAnalysis{}, &UnsupportedPairError{SourceFormat: document.Format, TargetFormat: targetFormat}
	}
	return analysis, nil
}

func appendMetadataLosses(losses *[]Loss, document model.Document, targetFormat string) error {
	fields := []struct {
		name  string
		value *string
	}{
		{name: "title", value: document.Metadata.Title},
		{name: "language", value: document.Metadata.Language},
		{name: "kind", value: document.Metadata.Kind},
		{name: "description", value: document.Metadata.Description},
	}
	for _, field := range fields {
		if field.value == nil {
			continue
		}
		if err := appendOneLoss(losses, newLoss(document, targetFormat, LossCodeMetadataOmitted, KindOmitted, "document metadata is omitted from subtitle output", "/metadata/"+field.name, nil, nil, []Attribute{{Name: "field", Value: field.name}})); err != nil {
			return err
		}
	}
	return nil
}

func appendSubRipCueLosses(losses *[]Loss, document model.Document, targetFormat string, cueIndex int, cue *model.Cue, translation payloadTranslation) error {
	base := fmt.Sprintf("/cues/%d", cueIndex)
	if cue.FormatData.SubRip != nil && cue.FormatData.SubRip.Coordinates != nil {
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeSubRipCoordinatesOmitted, KindOmitted, "SubRip pixel coordinates are omitted without viewport evidence", base+"/format_data/subrip/coordinates", nil)); err != nil {
			return err
		}
	}
	if err := appendCommonCueLosses(losses, document, targetFormat, cueIndex, cue, false); err != nil {
		return err
	}
	return appendPayloadLosses(losses, document, targetFormat, cueIndex, cue, translation)
}

func appendWebVTTDocumentLosses(losses *[]Loss, document model.Document, targetFormat string) error {
	data := document.FormatData.WebVTT
	if data == nil {
		return nil
	}
	if data.Description != nil {
		if err := appendOneLoss(losses, newLoss(document, targetFormat, LossCodeWebVTTDescriptionOmitted, KindOmitted, "WebVTT signature description is omitted", "/format_data/webvtt/description", nil, nil, nil)); err != nil {
			return err
		}
	}
	for index := range data.MetadataLines {
		if err := appendOneLoss(losses, newLoss(document, targetFormat, LossCodeWebVTTMetadataOmitted, KindOmitted, "WebVTT header metadata line is omitted", fmt.Sprintf("/format_data/webvtt/metadata_lines/%d", index), nil, nil, []Attribute{{Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
			return err
		}
	}
	for index := range data.Blocks {
		block := &data.Blocks[index]
		order := block.SourceOrder
		if err := appendOneLoss(losses, newLoss(document, targetFormat, LossCodeWebVTTBlockOmitted, KindOmitted, "WebVTT non-cue block is omitted", fmt.Sprintf("/format_data/webvtt/blocks/%d", index), &order, nil, []Attribute{{Name: "block_type", Value: block.Type}, {Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
			return err
		}
	}
	return nil
}

func appendWebVTTCueLosses(losses *[]Loss, document model.Document, targetFormat string, cueIndex int, cue *model.Cue, translation payloadTranslation) error {
	base := fmt.Sprintf("/cues/%d", cueIndex)
	native := cue.FormatData.WebVTT
	if cue.SourceIdentifier != nil {
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeWebVTTIdentifierOmitted, KindOmitted, "WebVTT cue identifier is omitted", base+"/source_identifier", nil)); err != nil {
			return err
		}
	}
	settingPaths := make(map[string]bool)
	if native != nil {
		for _, name := range webVTTSettingNames {
			if _, exists := native.Settings[name]; !exists {
				continue
			}
			settingPaths[name] = true
			if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeWebVTTSettingOmitted, KindOmitted, "WebVTT cue setting is omitted", base+"/format_data/webvtt/settings/"+escapePointer(name), []Attribute{{Name: "setting", Value: name}})); err != nil {
				return err
			}
		}
		seen := make(map[string]bool)
		for occurrenceIndex, occurrence := range native.SettingOccurrences {
			duplicate := occurrence.Recognized && occurrence.Valid && seen[occurrence.Name]
			if occurrence.Recognized && occurrence.Valid {
				seen[occurrence.Name] = true
			}
			if occurrence.Recognized && occurrence.Valid && !duplicate {
				continue
			}
			if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeWebVTTSettingOccurrenceOmitted, KindOmitted, "WebVTT setting occurrence is omitted", fmt.Sprintf("%s/format_data/webvtt/setting_occurrences/%d", base, occurrenceIndex), []Attribute{{Name: "occurrence", Value: strconv.Itoa(occurrenceIndex)}})); err != nil {
				return err
			}
		}
	}
	if err := appendCommonCueLosses(losses, document, targetFormat, cueIndex, cue, len(settingPaths) > 0); err != nil {
		return err
	}
	return appendPayloadLosses(losses, document, targetFormat, cueIndex, cue, translation)
}

func appendCommonCueLosses(losses *[]Loss, document model.Document, targetFormat string, cueIndex int, cue *model.Cue, placementAccounted bool) error {
	base := fmt.Sprintf("/cues/%d", cueIndex)
	for index := range cue.Speakers {
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeSpeakerObservationOmitted, KindOmitted, "speaker observation is omitted from target structure", fmt.Sprintf("%s/speakers/%d", base, index), []Attribute{{Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
			return err
		}
	}
	for index := range cue.Tokens {
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeTokenTimingOmitted, KindOmitted, "token timing is omitted from target structure", fmt.Sprintf("%s/tokens/%d", base, index), []Attribute{{Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
			return err
		}
	}
	for index := range cue.OCRObservations {
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodeOCRObservationOmitted, KindOmitted, "OCR observation is omitted from subtitle output", fmt.Sprintf("%s/ocr_observations/%d", base, index), []Attribute{{Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
			return err
		}
	}
	if cue.Placement != nil && !placementAccounted {
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, LossCodePlacementOmitted, KindOmitted, "common placement is omitted from target output", base+"/placement", nil)); err != nil {
			return err
		}
	}
	return nil
}

func appendPayloadLosses(losses *[]Loss, document model.Document, targetFormat string, cueIndex int, cue *model.Cue, translation payloadTranslation) error {
	path := fmt.Sprintf("/cues/%d/payload/raw_text", cueIndex)
	for _, issue := range translation.Issues {
		context := []Attribute{{Name: "occurrence", Value: occurrenceValue(issue.Occurrence)}}
		if issue.Feature != "" {
			context = append(context, Attribute{Name: "syntax", Value: issue.Feature})
		}
		if err := appendOneLoss(losses, cueLoss(document, targetFormat, *cue, issue.Code, issue.Kind, payloadLossMessage(issue.Code), path, context)); err != nil {
			return err
		}
	}
	return nil
}

func appendOneLoss(losses *[]Loss, loss Loss) error {
	if len(*losses) >= MaxLosses {
		return fmt.Errorf("conversion produces more than %d losses", MaxLosses)
	}
	*losses = append(*losses, loss)
	return nil
}

func payloadLossMessage(code string) string {
	switch code {
	case LossCodePayloadLineDegraded:
		return "empty payload line uses a non-visible target placeholder"
	case LossCodeSubRipFontDegraded:
		return "SubRip font presentation is removed while text is retained"
	case LossCodeNULDegraded:
		return "NUL is replaced by the target semantic replacement character"
	case LossCodeWebVTTVoiceDegraded:
		return "WebVTT voice annotation is removed while text is retained"
	case LossCodeWebVTTMarkupDegraded:
		return "WebVTT native markup meaning is weakened while text is retained"
	case LossCodeWebVTTInlineTimingOmitted:
		return "WebVTT inline timestamp is omitted while text is retained"
	case LossCodeWebVTTEntityAmbiguous:
		return "WebVTT character reference remains encoded to avoid target markup ambiguity"
	default:
		return "source payload semantic is changed for the target grammar"
	}
}

func newLoss(document model.Document, targetFormat, code string, kind Kind, message, path string, sourceOrder *int, cueID *string, context []Attribute) Loss {
	return Loss{Code: code, Severity: SeverityWarning, Kind: kind, Message: message, SourceFormat: document.Format, TargetFormat: targetFormat, Path: path, SourceOrder: sourceOrder, CueID: cueID, Context: context}
}

func cueLoss(document model.Document, targetFormat string, cue model.Cue, code string, kind Kind, message, path string, context []Attribute) Loss {
	order, cueID := cue.SourceOrder, cue.ID
	return newLoss(document, targetFormat, code, kind, message, path, &order, &cueID, context)
}

func escapePointer(value string) string {
	return strings.ReplaceAll(strings.ReplaceAll(value, "~", "~0"), "/", "~1")
}

func occurrenceValue(index int) string {
	return fmt.Sprintf("%09d", index)
}

func projectionFailure(sourceFormat, targetFormat, path string, err error) error {
	return &ProjectionError{SourceFormat: sourceFormat, TargetFormat: targetFormat, Path: path, Err: err}
}
