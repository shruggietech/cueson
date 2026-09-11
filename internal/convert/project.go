package convert

import (
	"fmt"
	"strconv"

	"github.com/shruggietech/cueson/internal/model"
)

func projectDocument(source model.Document, targetFormat string, translations []payloadTranslation) (model.Document, error) {
	if len(translations) != len(source.Cues) {
		return model.Document{}, fmt.Errorf("payload translation count does not match cue count")
	}
	target := model.Document{
		Schema: source.Schema, SchemaVersion: source.SchemaVersion, Format: targetFormat,
		FormatSupport: model.FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
		Producer:      source.Producer, Source: cloneSourceEnvelope(source.Source), Metadata: model.Metadata{},
		Cues: make([]model.Cue, len(source.Cues)), Diagnostics: []model.Diagnostic{},
	}
	for index := range source.Cues {
		sourceCue := &source.Cues[index]
		if targetFormat == "webvtt" && sourceCue.Timing.EndMilliseconds == sourceCue.Timing.StartMilliseconds {
			return model.Document{}, fmt.Errorf("WebVTT requires positive cue duration")
		}
		payload := translations[index].Text
		cue := model.Cue{
			ID: fmt.Sprintf("cue-%06d", index), Ordinal: index, SourceOrder: index,
			Timing:   sourceCue.Timing,
			Payload:  model.Payload{RawText: payload, PlainText: payload, Lines: splitPayloadLines(payload)},
			Speakers: []model.Speaker{}, Tokens: []model.Token{}, OCRObservations: []model.OCRObservation{},
		}
		switch targetFormat {
		case "subrip":
			cue.Payload.PlainText = subRipPlainText(payload)
			cue.FormatData.SubRip = &model.SubRipCueData{
				SequenceLineRaw: strconv.Itoa(index + 1),
				TimingLineRaw:   formatSubRipTimestamp(cue.Timing.StartMilliseconds) + " --> " + formatSubRipTimestamp(cue.Timing.EndMilliseconds),
			}
		case "webvtt":
			cue.Payload.PlainText = semanticText(sourceCue.Payload.PlainText)
			cue.FormatData.WebVTT = &model.WebVTTCueData{
				TimingLineRaw: formatWebVTTTimestamp(cue.Timing.StartMilliseconds) + " --> " + formatWebVTTTimestamp(cue.Timing.EndMilliseconds),
				SettingsRaw:   "", Settings: map[string]string{}, SettingOccurrences: []model.WebVTTSettingOccurrence{},
				RawPayload: payload, RawPayloadLines: splitPayloadLines(payload),
			}
		default:
			return model.Document{}, fmt.Errorf("unsupported target format %q", targetFormat)
		}
		target.Cues[index] = cue
	}
	if targetFormat == "subrip" {
		target.FormatData.SubRip = &model.SubRipDocumentData{Dialect: "subrip"}
	} else {
		target.FormatData.WebVTT = &model.WebVTTDocumentData{Signature: "WEBVTT", SignatureLineRaw: "WEBVTT", MetadataLines: []string{}, Blocks: []model.WebVTTBlock{}}
	}
	recomputeSummaries(&target)
	if err := target.Validate(); err != nil {
		return model.Document{}, fmt.Errorf("target document is invalid: %w", err)
	}
	return target, nil
}

func cloneSourceEnvelope(source model.SourceEnvelope) model.SourceEnvelope {
	cloned := source
	cloned.Assets = append([]model.SourceAsset(nil), source.Assets...)
	return cloned
}

func splitPayloadLines(payload string) []string {
	return splitOnLF(payload)
}

func recomputeSummaries(document *model.Document) {
	minimum := document.Cues[0].Timing.StartMilliseconds
	maximum := document.Cues[0].Timing.EndMilliseconds
	for index := range document.Cues {
		if document.Cues[index].Timing.StartMilliseconds < minimum {
			minimum = document.Cues[index].Timing.StartMilliseconds
		}
		if document.Cues[index].Timing.EndMilliseconds > maximum {
			maximum = document.Cues[index].Timing.EndMilliseconds
		}
	}
	span := maximum - minimum
	document.Document = model.DocumentSummary{CueCount: len(document.Cues), MediaStartMilliseconds: int64Pointer(minimum), MediaEndMilliseconds: int64Pointer(maximum), MediaSpanMilliseconds: int64Pointer(span)}
	document.Stats = model.Stats{CueCount: len(document.Cues), MediaSpanMilliseconds: int64Pointer(span)}
}

func formatSubRipTimestamp(milliseconds int64) string {
	hours, remainder := milliseconds/3_600_000, milliseconds%3_600_000
	minutes, remainder := remainder/60_000, remainder%60_000
	seconds, fraction := remainder/1_000, remainder%1_000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, fraction)
}

func formatWebVTTTimestamp(milliseconds int64) string {
	hours, remainder := milliseconds/3_600_000, milliseconds%3_600_000
	minutes, remainder := remainder/60_000, remainder%60_000
	seconds, fraction := remainder/1_000, remainder%1_000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, fraction)
}

func int64Pointer(value int64) *int64 { return &value }
