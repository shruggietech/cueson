package model

import "fmt"

// Structural safety limits bound collection amplification after a Cue JSON
// document has passed the byte-level input boundary.
const (
	MaxDocumentItems   = 65_536
	MaxItemOccurrences = 1_024
	MaxDiagnostics     = 8_192
)

func validateCollectionLimits(document Document) error {
	if err := boundedLength("source.assets", len(document.Source.Assets), MaxDocumentItems); err != nil {
		return err
	}
	if err := boundedLength("cues", len(document.Cues), MaxDocumentItems); err != nil {
		return err
	}
	if err := boundedLength("diagnostics", len(document.Diagnostics), MaxDiagnostics); err != nil {
		return err
	}

	for cueIndex := range document.Cues {
		cue := &document.Cues[cueIndex]
		prefix := fmt.Sprintf("cues[%d]", cueIndex)
		for _, collection := range []struct {
			name   string
			length int
		}{
			{name: prefix + ".payload.lines", length: len(cue.Payload.Lines)},
			{name: prefix + ".speakers", length: len(cue.Speakers)},
			{name: prefix + ".tokens", length: len(cue.Tokens)},
			{name: prefix + ".ocr_observations", length: len(cue.OCRObservations)},
		} {
			if err := boundedLength(collection.name, collection.length, MaxItemOccurrences); err != nil {
				return err
			}
		}
		for observationIndex := range cue.OCRObservations {
			observation := &cue.OCRObservations[observationIndex]
			observationPrefix := fmt.Sprintf("%s.ocr_observations[%d]", prefix, observationIndex)
			for _, collection := range []struct {
				name   string
				length int
			}{
				{name: observationPrefix + ".lines", length: len(observation.Lines)},
				{name: observationPrefix + ".alternatives", length: len(observation.Alternatives)},
				{name: observationPrefix + ".source_regions", length: len(observation.SourceRegions)},
				{name: observationPrefix + ".processing_options", length: len(observation.ProcessingOptions)},
			} {
				if err := boundedLength(collection.name, collection.length, MaxItemOccurrences); err != nil {
					return err
				}
			}
		}
		if cue.FormatData.WebVTT != nil {
			native := cue.FormatData.WebVTT
			for _, collection := range []struct {
				name   string
				length int
			}{
				{name: prefix + ".format_data.webvtt.settings", length: len(native.Settings)},
				{name: prefix + ".format_data.webvtt.setting_occurrences", length: len(native.SettingOccurrences)},
				{name: prefix + ".format_data.webvtt.raw_payload_lines", length: len(native.RawPayloadLines)},
			} {
				if err := boundedLength(collection.name, collection.length, MaxItemOccurrences); err != nil {
					return err
				}
			}
		}
	}

	if document.FormatData.WebVTT == nil {
		return nil
	}
	webvtt := document.FormatData.WebVTT
	if err := boundedLength("format_data.webvtt.metadata_lines", len(webvtt.MetadataLines), MaxItemOccurrences); err != nil {
		return err
	}
	if err := boundedLength("format_data.webvtt.blocks", len(webvtt.Blocks), MaxDocumentItems); err != nil {
		return err
	}
	for blockIndex := range webvtt.Blocks {
		block := &webvtt.Blocks[blockIndex]
		prefix := fmt.Sprintf("format_data.webvtt.blocks[%d]", blockIndex)
		if err := boundedLength(prefix+".raw_lines", len(block.RawLines), MaxItemOccurrences); err != nil {
			return err
		}
		if block.Region != nil {
			if err := boundedLength(prefix+".region.settings", len(block.Region.Settings), MaxItemOccurrences); err != nil {
				return err
			}
			if err := boundedLength(prefix+".region.setting_occurrences", len(block.Region.SettingOccurrences), MaxItemOccurrences); err != nil {
				return err
			}
		}
	}
	return nil
}

func boundedLength(path string, length, maximum int) error {
	if length > maximum {
		return fmt.Errorf("%s contains %d items, maximum is %d", path, length, maximum)
	}
	return nil
}
