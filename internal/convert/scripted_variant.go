package convert

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

// projectScriptedVariant preserves ordered native owners and reports each
// dialect field whose role or presentation cannot survive the target profile.
func projectScriptedVariant(ctx context.Context, source model.Document, targetFormat string) (model.Document, []Loss, error) {
	if ctx == nil {
		return model.Document{}, nil, fmt.Errorf("invalid_context")
	}
	if (source.Format != "ass" && source.Format != "ssa") || (targetFormat != "ass" && targetFormat != "ssa") || source.Format == targetFormat {
		return model.Document{}, nil, fmt.Errorf("invalid_scripted_variant_pair")
	}
	if err := source.Validate(); err != nil {
		return model.Document{}, nil, err
	}
	for _, diagnostic := range source.Diagnostics {
		if diagnostic.Code == "unsupported_override" || diagnostic.Code == "malformed_override" || diagnostic.Code == "unsupported_karaoke_timing" || diagnostic.Code == "malformed_attachment" {
			return model.Document{}, nil, fmt.Errorf("unsupported_scripted_variant_content")
		}
	}
	original := source.FormatData.ASS
	if source.Format == "ssa" {
		original = source.FormatData.SSA
	}
	target := source
	target.Format = targetFormat
	target.FormatSupport = model.FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}
	target.Metadata = model.Metadata{}
	target.Source = cloneSourceEnvelope(source.Source)
	target.Cues = append([]model.Cue{}, source.Cues...)
	target.Diagnostics = append([]model.Diagnostic{}, source.Diagnostics...)
	native := *original
	native.Sections = append([]model.ScriptedSection{}, original.Sections...)
	native.Records = append([]model.ScriptedRecord{}, original.Records...)
	native.Styles = append([]model.ScriptedStyle{}, original.Styles...)
	native.Events = append([]model.ScriptedEvent{}, original.Events...)
	native.Attachments = append([]model.ScriptedAttachment{}, original.Attachments...)
	native.Dialect = "v4.00+"
	if targetFormat == "ssa" {
		native.Dialect = "v4.00"
		target.FormatData = model.DocumentFormatData{SSA: &native}
	} else {
		target.FormatData = model.DocumentFormatData{ASS: &native}
	}
	orders := map[string]int{}
	contexts := map[string]string{}
	for index, section := range native.Sections {
		name := strings.ToLower(section.Name)
		if name == "v4+ styles" || name == "v4 styles" {
			if targetFormat == "ssa" {
				native.Sections[index].Name = "V4 Styles"
			} else {
				native.Sections[index].Name = "V4+ Styles"
			}
			native.Sections[index].RawHeader = nil
			contexts[section.SectionID] = "style"
		} else if name == "events" {
			contexts[section.SectionID] = "event"
		}
	}
	for index, record := range native.Records {
		if err := ctx.Err(); err != nil {
			return model.Document{}, nil, err
		}
		orders[record.RecordID] = record.SourceOrder
		if record.Kind == "malformed" {
			return model.Document{}, nil, fmt.Errorf("malformed_scripted_variant_record")
		}
		if record.Kind == "metadata" {
			fields := append([]model.ScriptedField{}, record.Fields...)
			for fieldIndex, field := range fields {
				if strings.EqualFold(strings.TrimSpace(field.FieldName), "ScriptType") {
					fields[fieldIndex].RawValue = native.Dialect
					fields[fieldIndex].TypedValue = nil
					native.Records[index].RawLine = nil
				}
			}
			native.Records[index].Fields = fields
		}
		if record.Kind == "style" || record.Kind == "event" {
			native.Records[index].RawLine = nil
		}
		if record.Kind == "format_declaration" {
			kind := contexts[record.SectionID]
			extras := []string{}
			for _, field := range record.DeclarationFields {
				if _, known := model.ScriptedFieldIdentity(field.FieldName, source.Format, kind); !known {
					extras = append(extras, field.FieldName)
				}
			}
			native.Records[index].RawLine = nil
			native.Records[index].DeclarationFields = scriptedDeclaration(targetFormat, kind, extras)
		}
	}
	losses := []Loss{}
	if err := appendMetadataLosses(&losses, source, targetFormat); err != nil {
		return model.Document{}, nil, err
	}
	for styleIndex, style := range original.Styles {
		if err := ctx.Err(); err != nil {
			return model.Document{}, nil, err
		}
		if !style.Valid {
			return model.Document{}, nil, fmt.Errorf("malformed_scripted_variant_style")
		}
		mapped, err := scriptedVariantStyle(source, targetFormat, styleIndex, style, orders[style.RecordID], &losses)
		if err != nil {
			return model.Document{}, nil, err
		}
		native.Styles[styleIndex].Fields = mapped
	}
	cueByEvent := map[string]model.Cue{}
	for _, cue := range source.Cues {
		data := cue.FormatData.ASS
		if source.Format == "ssa" {
			data = cue.FormatData.SSA
		}
		cueByEvent[data.EventID] = cue
	}
	for eventIndex, event := range original.Events {
		if err := ctx.Err(); err != nil {
			return model.Document{}, nil, err
		}
		if !event.Valid {
			return model.Document{}, nil, fmt.Errorf("malformed_scripted_variant_event")
		}
		values := map[string]string{}
		extras := []model.ScriptedField{}
		for fieldIndex, field := range event.Fields {
			identity, known := model.ScriptedFieldIdentity(field.FieldName, source.Format, "event")
			if !known {
				copy := field
				copy.TypedValue = nil
				extras = append(extras, copy)
				continue
			}
			if identity == "layer" || identity == "marked" {
				order := orders[event.RecordID]
				var cueID *string
				if cue, ok := cueByEvent[event.EventID]; ok {
					id := cue.ID
					cueID = &id
				}
				loss := newLoss(source, targetFormat, LossCodeScriptedVariantFieldOmitted, KindOmitted, "dialect-specific event field has no equivalent target role", fmt.Sprintf("/format_data/%s/events/%d/fields/%d", source.Format, eventIndex, fieldIndex), &order, cueID, nil)
				if err := appendOneLoss(&losses, loss); err != nil {
					return model.Document{}, nil, err
				}
				continue
			}
			values[identity] = field.RawValue
		}
		values["layer"], values["marked"] = "0", "Marked=0"
		native.Events[eventIndex].Fields = scriptedVariantFields(targetFormat, "event", values, extras)
	}
	eventIndexes := map[string]int{}
	for index, event := range native.Events {
		eventIndexes[event.EventID] = index
	}
	textNames := model.ScriptedTextNames{ResetStyles: map[string]bool{}, FontFamilies: map[string]bool{}}
	for _, style := range native.Styles {
		textNames.ResetStyles[style.Name] = true
		for _, field := range style.Fields {
			if strings.EqualFold(field.FieldName, "Fontname") {
				textNames.FontFamilies[field.RawValue] = true
			}
		}
	}
	for index, cue := range source.Cues {
		if err := ctx.Err(); err != nil {
			return model.Document{}, nil, err
		}
		for occurrence := range cue.OCRObservations {
			loss := cueLoss(source, targetFormat, cue, LossCodeOCRObservationOmitted, KindOmitted, "OCR observation is omitted from subtitle output", fmt.Sprintf("/cues/%d/ocr_observations/%d", index, occurrence), []Attribute{{Name: "occurrence", Value: strconv.Itoa(occurrence)}})
			if err := appendOneLoss(&losses, loss); err != nil {
				return model.Document{}, nil, err
			}
		}
		if cue.Placement != nil {
			loss := cueLoss(source, targetFormat, cue, LossCodePlacementOmitted, KindOmitted, "common placement is omitted from target output", fmt.Sprintf("/cues/%d/placement", index), nil)
			if err := appendOneLoss(&losses, loss); err != nil {
				return model.Document{}, nil, err
			}
		}
		target.Cues[index].OCRObservations, target.Cues[index].Placement = []model.OCRObservation{}, nil
		if cue.SourceIdentifier != nil {
			loss := cueLoss(source, targetFormat, cue, LossCodeSourceIdentifierOmitted, KindOmitted, "source cue identifier has no target native representation", fmt.Sprintf("/cues/%d/source_identifier", index), nil)
			if err := appendOneLoss(&losses, loss); err != nil {
				return model.Document{}, nil, err
			}
			target.Cues[index].SourceIdentifier = nil
		}
		data := cue.FormatData.ASS
		if source.Format == "ssa" {
			data = cue.FormatData.SSA
		}
		copy := *data
		timing, err := quantizeScriptedTiming(cue.Timing)
		if err != nil {
			return model.Document{}, nil, projectionFailure(source.Format, targetFormat, fmt.Sprintf("/cues/%d/timing", index), err)
		}
		if err := appendScriptedPrecisionLosses(source, targetFormat, index, timing, &losses); err != nil {
			return model.Document{}, nil, err
		}
		target.Cues[index].Timing = timing
		event := &native.Events[eventIndexes[data.EventID]]
		facts, err := model.ProjectScriptedText(event.Text, copy.Projection.WrapStyle, timing.StartMilliseconds, timing.EndMilliseconds, textNames)
		if err != nil || len(facts.DiagnosticCodes) != 0 {
			return model.Document{}, nil, fmt.Errorf("unsupported_scripted_variant_projection")
		}
		event.Spans, event.Tags, event.Karaoke = facts.Spans, facts.Tags, facts.Karaoke
		target.Cues[index].Payload = model.Payload{RawText: event.Text, PlainText: strings.Join(facts.Lines, "\n"), Lines: facts.Lines}
		target.Cues[index].Tokens = facts.Tokens
		copy.StartTimestampRaw, copy.EndTimestampRaw = nil, nil
		copy.Projection.TextFieldName = "Text"
		if copy.Projection.SpeakerFieldName != nil {
			name := "Name"
			copy.Projection.SpeakerFieldName = &name
		}
		if targetFormat == "ssa" {
			target.Cues[index].FormatData = model.CueFormatData{SSA: &copy}
		} else {
			target.Cues[index].FormatData = model.CueFormatData{ASS: &copy}
		}
	}
	if len(target.Cues) != 0 {
		recomputeSummaries(&target)
		for _, cue := range target.Cues {
			if len(cue.Tokens) != 0 {
				target.Document.HasWordLevelTiming, target.Stats.HasWordLevelTiming = true, true
			}
		}
		target.Stats.DiagnosticCount, target.Stats.WarningCount, target.Stats.ErrorCount = source.Stats.DiagnosticCount, source.Stats.WarningCount, source.Stats.ErrorCount
	}
	if err := model.ValidateScriptedTarget(source, target); err != nil {
		return model.Document{}, nil, err
	}
	return target, losses, nil
}

func scriptedVariantStyle(source model.Document, targetFormat string, styleIndex int, style model.ScriptedStyle, order int, losses *[]Loss) ([]model.ScriptedField, error) {
	values := map[string]string{}
	extras := []model.ScriptedField{}
	fieldIndexes := map[string]int{}
	for index, field := range style.Fields {
		identity, known := model.ScriptedFieldIdentity(field.FieldName, source.Format, "style")
		if !known {
			copy := field
			copy.TypedValue = nil
			extras = append(extras, copy)
			continue
		}
		fieldIndexes[identity] = index
		values[identity] = field.RawValue
	}
	add := func(identity, code string, kind Kind, message string) error {
		index, ok := fieldIndexes[identity]
		if !ok {
			return fmt.Errorf("missing_scripted_variant_field")
		}
		return appendOneLoss(losses, newLoss(source, targetFormat, code, kind, message, fmt.Sprintf("/format_data/%s/styles/%d/fields/%d", source.Format, styleIndex, index), &order, nil, nil))
	}
	alignment, err := strconv.ParseInt(strings.TrimSpace(values["alignment"]), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("unsupported_scripted_variant_alignment")
	}
	mapping := map[int64]int64{1: 1, 2: 2, 3: 3, 5: 7, 6: 8, 7: 9, 9: 4, 10: 5, 11: 6}
	if source.Format == "ass" {
		reverse := map[int64]int64{}
		for ssa, ass := range mapping {
			reverse[ass] = ssa
		}
		mapping = reverse
	}
	mapped, ok := mapping[alignment]
	if !ok {
		return nil, fmt.Errorf("unsupported_scripted_variant_alignment")
	}
	values["alignment"] = strconv.FormatInt(mapped, 10)
	colors := map[string]model.ScriptedColor{}
	for _, identity := range []string{"primarycolour", "secondarycolour", "backcolour"} {
		field := style.Fields[fieldIndexes[identity]]
		value, err := model.ScriptedScalarValue(field, source.Format, "style", order)
		if err != nil {
			return nil, err
		}
		colors[identity] = *value.Color
	}
	if source.Format == "ass" {
		for _, identity := range []string{"underline", "strikeout", "scalex", "scaley", "spacing", "angle"} {
			if err := add(identity, LossCodeScriptedVariantFieldOmitted, KindOmitted, "style field is not represented in the target dialect"); err != nil {
				return nil, err
			}
			delete(values, identity)
		}
		if err := add("outlinecolour", LossCodeScriptedVariantFieldDegraded, KindDegraded, "independent outline color is replaced by the target shared color role"); err != nil {
			return nil, err
		}
		primaryAlpha := colors["primarycolour"].Alpha
		if colors["secondarycolour"].Alpha != primaryAlpha {
			if err := add("secondarycolour", LossCodeScriptedVariantFieldDegraded, KindDegraded, "independent color alpha is replaced by the target global alpha"); err != nil {
				return nil, err
			}
		}
		if colors["backcolour"].Alpha != 128 {
			if err := add("backcolour", LossCodeScriptedVariantFieldDegraded, KindDegraded, "shadow color alpha is replaced by the target fixed alpha"); err != nil {
				return nil, err
			}
		}
		values["alphalevel"] = strconv.Itoa(int(primaryAlpha))
		values["tertiarycolour"] = "0"
		for identity, color := range colors {
			color.Alpha = 0
			values[identity] = strconv.FormatUint(uint64(scriptedColorBits(color)), 10)
		}
	} else {
		if err := add("tertiarycolour", LossCodeScriptedVariantFieldOmitted, KindOmitted, "legacy tertiary color observation has no target rendering role"); err != nil {
			return nil, err
		}
		alpha, err := strconv.ParseInt(strings.TrimSpace(values["alphalevel"]), 10, 64)
		if err != nil || alpha < 0 || alpha > 255 {
			return nil, fmt.Errorf("unsupported_scripted_variant_alpha")
		}
		for _, identity := range []string{"primarycolour", "secondarycolour"} {
			color := colors[identity]
			if color.Alpha != 0 && color.Alpha != uint8(alpha) {
				if err := add(identity, LossCodeScriptedVariantFieldDegraded, KindDegraded, "retained native color alpha observation is replaced by the active global alpha"); err != nil {
					return nil, err
				}
			}
			color.Alpha = uint8(alpha)
			values[identity] = fmt.Sprintf("&H%08X", scriptedColorBits(color))
		}
		outline := colors["backcolour"]
		outline.Alpha = uint8(alpha)
		shadow := colors["backcolour"]
		if shadow.Alpha != 0 && shadow.Alpha != 128 {
			if err := add("backcolour", LossCodeScriptedVariantFieldDegraded, KindDegraded, "retained native color alpha observation is replaced by the active fixed shadow alpha"); err != nil {
				return nil, err
			}
		}
		shadow.Alpha = 128
		values["outlinecolour"], values["backcolour"] = fmt.Sprintf("&H%08X", scriptedColorBits(outline)), fmt.Sprintf("&H%08X", scriptedColorBits(shadow))
		delete(values, "tertiarycolour")
		delete(values, "alphalevel")
		values["underline"], values["strikeout"], values["scalex"], values["scaley"], values["spacing"], values["angle"] = "0", "0", "100", "100", "0", "0"
	}
	return scriptedVariantFields(targetFormat, "style", values, extras), nil
}

func scriptedVariantFields(format, kind string, values map[string]string, extras []model.ScriptedField) []model.ScriptedField {
	fields := []model.ScriptedField{}
	for _, name := range model.ScriptedCanonicalFields(format, kind) {
		identity, _ := model.ScriptedFieldIdentity(name, format, kind)
		if kind == "event" && identity == "text" {
			fields = append(fields, extras...)
		}
		fields = append(fields, model.ScriptedField{FieldName: name, RawValue: values[identity]})
	}
	if kind == "style" {
		fields = append(fields, extras...)
	}
	return fields
}

func scriptedColorBits(color model.ScriptedColor) uint32 {
	return uint32(color.Alpha)<<24 | uint32(color.Blue)<<16 | uint32(color.Green)<<8 | uint32(color.Red)
}
