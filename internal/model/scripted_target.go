package model

import (
	"fmt"
	"reflect"
)

// ValidateScriptedTarget validates an internal conversion representation.
// It never authorizes a mismatched source envelope in published Cue JSON.
func ValidateScriptedTarget(source, target Document) error {
	if err := source.Validate(); err != nil {
		return fmt.Errorf("invalid_conversion_source")
	}
	if !reflect.DeepEqual(source.Source, target.Source) {
		return fmt.Errorf("conversion_source_envelope_changed")
	}
	if source.Producer != target.Producer {
		return fmt.Errorf("conversion_source_producer_changed")
	}
	if target.SchemaVersion != "1.1.0" || target.Schema != "https://cueson.io/schema/v1.1.0/cueson.schema.json" || (target.Format != "ass" && target.Format != "ssa") {
		return fmt.Errorf("invalid_scripted_target_identity")
	}
	return validateScriptedDocumentWithSourcePolicy(target, true)
}

func scriptedTargetHasCaptures(doc Document) bool {
	native := doc.FormatData.ASS
	if doc.Format == "ssa" {
		native = doc.FormatData.SSA
	}
	if native == nil {
		return false
	}
	for _, section := range native.Sections {
		if section.RawHeader != nil {
			return true
		}
	}
	for _, record := range native.Records {
		if record.RawLine != nil {
			return true
		}
	}
	for _, cue := range doc.Cues {
		data := cue.FormatData.ASS
		if doc.Format == "ssa" {
			data = cue.FormatData.SSA
		}
		if data != nil && (data.StartTimestampRaw != nil || data.EndTimestampRaw != nil) {
			return true
		}
	}
	return false
}
