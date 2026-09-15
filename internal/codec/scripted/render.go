package scripted

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/source"
)

const maxRenderedBytes = 64 << 20

// RenderResult owns complete canonical output and ordered safe observations.
type RenderResult struct {
	Bytes       []byte
	Diagnostics []model.Diagnostic
}

// Render serializes editable native owners, independently of source restoration.
// No bytes are returned when any owner, framing, precision or safety check fails.
func Render(ctx context.Context, document model.Document, strict bool) (RenderResult, error) {
	if ctx == nil {
		return RenderResult{}, fmt.Errorf("render scripted: invalid_context")
	}
	if err := ctx.Err(); err != nil {
		return RenderResult{}, err
	}
	if document.Format != "ass" && document.Format != "ssa" {
		return RenderResult{}, fmt.Errorf("render scripted: mismatched_native_dialect")
	}
	if err := document.Validate(); err != nil {
		return RenderResult{}, err
	}
	if err := source.ValidateIntegrity(ctx, document); err != nil {
		if cancelled := ctx.Err(); cancelled != nil {
			return RenderResult{}, cancelled
		}
		// Envelope errors may include caller-controlled asset names. Keep this
		// publication boundary diagnostic independent of their unsafe contents.
		return RenderResult{}, fmt.Errorf("render scripted: invalid_source_integrity")
	}
	native := document.FormatData.ASS
	if document.Format == "ssa" {
		native = document.FormatData.SSA
	}
	for _, diagnostic := range document.Diagnostics {
		if slices.Contains([]string{"malformed_native_record", "malformed_attachment"}, diagnostic.Code) {
			return RenderResult{}, fmt.Errorf("render scripted: %s", diagnostic.Code)
		}
		if strict && diagnostic.Code != "format_extension_disagreement" {
			return RenderResult{}, fmt.Errorf("render scripted strict: %s", diagnostic.Code)
		}
	}
	sections := make(map[string]model.ScriptedSection, len(native.Sections))
	records := make(map[string]model.ScriptedRecord, len(native.Records))
	styles := make(map[string]model.ScriptedStyle, len(native.Styles))
	events := make(map[string]model.ScriptedEvent, len(native.Events))
	attachments := make(map[string]model.ScriptedAttachment, len(native.Attachments))
	cues := make(map[string]model.Cue, len(document.Cues))
	items := make([]string, len(native.Sections)+len(native.Records))
	for _, section := range native.Sections {
		sections[section.SectionID] = section
		items[section.SourceOrder] = "s:" + section.SectionID
	}
	for _, record := range native.Records {
		records[record.RecordID] = record
		items[record.SourceOrder] = "r:" + record.RecordID
	}
	styleNames := map[string]bool{}
	for _, style := range native.Styles {
		if !style.Valid {
			return RenderResult{}, fmt.Errorf("render scripted: malformed_native_record")
		}
		if strict && styleNames[style.Name] {
			return RenderResult{}, fmt.Errorf("render scripted strict: duplicate_style")
		}
		styleNames[style.Name] = true
		styles[style.StyleID] = style
	}
	for _, event := range native.Events {
		if !event.Valid || (event.EventType != "dialogue" && event.EventType != "comment") {
			return RenderResult{}, fmt.Errorf("render scripted: malformed_native_record")
		}
		events[event.EventID] = event
	}
	for _, attachment := range native.Attachments {
		attachments[attachment.AttachmentID] = attachment
	}
	for _, cue := range document.Cues {
		cues[cue.ID] = cue
	}
	var output strings.Builder
	emittedLines := 0
	writeLine := func(line string) error {
		if !utf8.ValidString(line) || strings.ContainsAny(line, "\r\n\x00") {
			return fmt.Errorf("render scripted: invalid_native_line")
		}
		if len(line) > model.MaxScriptedLineBytes || output.Len() > maxRenderedBytes-len(line)-1 || emittedLines >= model.MaxDocumentItems {
			return fmt.Errorf("render scripted: complexity_limit")
		}
		output.WriteString(line)
		output.WriteByte('\n')
		emittedLines++
		return nil
	}
	var active []string
	currentSection := ""
	knownNonconforming := len(document.Diagnostics) > 0
	for _, item := range items {
		if err := ctx.Err(); err != nil {
			return RenderResult{}, err
		}
		kind, id, _ := strings.Cut(item, ":")
		if kind == "s" {
			section := sections[id]
			currentSection = section.Name
			active = nil
			if err := writeLine("[" + section.Name + "]"); err != nil {
				return RenderResult{}, err
			}
			continue
		}
		record := records[id]
		line := ""
		switch record.Kind {
		case "metadata":
			if len(record.Fields) != 1 || strings.ContainsAny(record.Fields[0].FieldName, ":\r\n\x00") || strings.TrimSpace(record.Fields[0].FieldName) == "" {
				return RenderResult{}, fmt.Errorf("render scripted: invalid_native_fields")
			}
			field := record.Fields[0]
			value := field.RawValue
			if strings.EqualFold(strings.TrimSpace(field.FieldName), "Timer") {
				value = "100"
			}
			if strings.EqualFold(strings.TrimSpace(field.FieldName), "WrapStyle") {
				n, _ := strconv.Atoi(strings.TrimSpace(value))
				value = strconv.Itoa(n)
			}
			line = strings.TrimSpace(field.FieldName) + ": " + value
		case "format_declaration":
			active = make([]string, len(record.DeclarationFields))
			fieldContext := "style"
			if strings.EqualFold(currentSection, "Events") {
				fieldContext = "event"
			}
			for i, field := range record.DeclarationFields {
				identity, known := model.ScriptedFieldIdentity(field.FieldName, document.Format, fieldContext)
				if known != (field.RecognizedName != nil) || (known && *field.RecognizedName != identity) {
					return RenderResult{}, fmt.Errorf("render scripted: inconsistent_recognized_name")
				}
				active[i] = field.FieldName
			}
			if err := renderValidateDeclaration(active, document.Format, currentSection); err != nil {
				return RenderResult{}, err
			}
			line = "Format: " + strings.Join(active, ",")
		case "style", "event":
			var fields []model.ScriptedField
			var declaration *string
			prefix := "Style"
			var cue *model.Cue
			if record.Kind == "style" {
				owner := styles[*record.StyleID]
				fields, declaration = owner.Fields, owner.DeclarationID
			} else {
				owner := events[*record.EventID]
				fields, declaration = owner.Fields, owner.DeclarationID
				if owner.EventType == "dialogue" {
					prefix = "Dialogue"
					c := cues[*owner.CueID]
					cue = &c
				} else {
					prefix = "Comment"
				}
			}
			if declaration == nil {
				fields = renderConstructedFields(fields, document.Format, record.Kind)
			}
			names := make([]string, len(fields))
			for i, field := range fields {
				names[i] = field.FieldName
			}
			if !slices.Equal(names, active) {
				if err := writeLine("Format: " + strings.Join(names, ",")); err != nil {
					return RenderResult{}, err
				}
				active = names
			}
			values := make([]string, len(fields))
			for i, field := range fields {
				value, err := renderScalar(field, document.Format, record.Kind)
				if err != nil {
					return RenderResult{}, err
				}
				n := renderName(field.FieldName)
				if record.Kind == "event" && (n == "start" || n == "end") {
					var milliseconds int64
					if cue != nil {
						milliseconds = cue.Timing.StartMilliseconds
						if n == "end" {
							milliseconds = cue.Timing.EndMilliseconds
						}
					} else {
						milliseconds, err = renderParseTimestamp(field.RawValue)
						if err != nil {
							return RenderResult{}, err
						}
					}
					value, err = renderTimestamp(milliseconds)
					if err != nil {
						return RenderResult{}, err
					}
				}
				values[i] = value
			}
			line = prefix + ": " + strings.Join(values, ",")
		case "attachment_header":
			attachment := attachments[*record.AttachmentID]
			prefix := "fontname"
			if attachment.AttachmentType == "graphic" {
				prefix = "filename"
			}
			line = prefix + ": " + attachment.Name
		case "malformed":
			return RenderResult{}, fmt.Errorf("render scripted: malformed_native_record")
		case "blank", "comment", "unknown", "attachment_data":
			line = *record.RawLine
			trim := strings.TrimSpace(line)
			if (record.Kind == "blank" && trim != "") || (record.Kind == "comment" && !strings.HasPrefix(trim, ";")) || (record.Kind == "unknown" && (trim == "" || strings.HasPrefix(trim, ";") || strings.Contains(trim, ":") || (strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]")))) {
				return RenderResult{}, fmt.Errorf("render scripted: inconsistent_native_framing")
			}
			if record.Kind == "unknown" {
				knownNonconforming = true
				if strict {
					return RenderResult{}, fmt.Errorf("render scripted strict: unknown_native_record")
				}
			}
		default:
			return RenderResult{}, fmt.Errorf("render scripted: unsupported_native_record")
		}
		if err := writeLine(line); err != nil {
			return RenderResult{}, err
		}
	}
	// Reparse the complete bounded candidate before exposing it. This checks
	// constructed framing and declarations, and independently derives strict
	// conformance observations even when an edited model omitted diagnostics.
	candidate := []byte(output.String())
	parsed, err := Parse(ctx, candidate, document.Format)
	if err != nil {
		return RenderResult{}, err
	}
	if strict && len(parsed.Diagnostics) > 0 {
		return RenderResult{}, fmt.Errorf("render scripted strict: %s", parsed.Diagnostics[0].Code)
	}
	knownNonconforming = knownNonconforming || len(parsed.Diagnostics) > 0
	result := RenderResult{Bytes: candidate, Diagnostics: []model.Diagnostic{}}
	if knownNonconforming {
		result.Diagnostics = append(result.Diagnostics, model.Diagnostic{Severity: "warning", Code: "scripted_render_preserved_nonconforming", Message: "permissive rendering retained diagnosed native content"})
	}
	return result, nil
}

func renderName(name string) string {
	b := []byte(strings.Trim(name, " "))
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func renderCanonicalNames(format, kind string) []string {
	return model.ScriptedCanonicalFields(format, kind)
}

func renderValidateDeclaration(names []string, format, section string) error {
	kind := "style"
	if strings.EqualFold(section, "Events") {
		kind = "event"
	} else if !strings.EqualFold(section, "V4+ Styles") && !strings.EqualFold(section, "V4 Styles") {
		return fmt.Errorf("render scripted: invalid_declaration")
	}
	required := map[string]bool{}
	for _, name := range renderCanonicalNames(format, kind) {
		required[renderName(name)] = false
	}
	for i, name := range names {
		n := renderName(name)
		if kind == "event" && n == "actor" {
			n = "name"
		}
		if n == "" || strings.ContainsAny(name, ",:\r\n\x00") {
			return fmt.Errorf("render scripted: invalid_declaration")
		}
		if seen, exists := required[n]; exists {
			if seen {
				return fmt.Errorf("render scripted: invalid_declaration")
			}
			required[n] = true
		}
		if n == "text" && i != len(names)-1 {
			return fmt.Errorf("render scripted: nonfinal_native_text")
		}
	}
	for _, seen := range required {
		if !seen {
			return fmt.Errorf("render scripted: invalid_declaration")
		}
	}
	return nil
}

func renderConstructedFields(fields []model.ScriptedField, format, kind string) []model.ScriptedField {
	canonical := renderCanonicalNames(format, kind)
	byName := map[string]model.ScriptedField{}
	known := map[string]bool{}
	for _, name := range canonical {
		known[renderName(name)] = true
	}
	for _, field := range fields {
		n := renderName(field.FieldName)
		if kind == "event" && n == "actor" {
			n = "name"
		}
		byName[n] = field
	}
	result := make([]model.ScriptedField, 0, len(fields))
	for _, name := range canonical {
		if name == "Text" {
			continue
		}
		field := byName[renderName(name)]
		field.FieldName = name
		result = append(result, field)
	}
	for _, field := range fields {
		n := renderName(field.FieldName)
		if kind == "event" && n == "actor" {
			n = "name"
		}
		if !known[n] {
			result = append(result, field)
		}
	}
	if kind == "event" {
		field := byName["text"]
		field.FieldName = "Text"
		result = append(result, field)
	}
	return result
}

func renderTimestamp(milliseconds int64) (string, error) {
	if milliseconds < 0 || milliseconds%10 != 0 {
		return "", fmt.Errorf("render scripted: unrepresentable_centiseconds")
	}
	return fmt.Sprintf("%d:%02d:%02d.%02d", milliseconds/3600000, milliseconds/60000%60, milliseconds/1000%60, milliseconds/10%100), nil
}

func renderParseTimestamp(raw string) (int64, error) {
	return model.ScriptedMilliseconds(raw)
}

func renderScalar(field model.ScriptedField, format, kind string) (string, error) {
	n, _ := model.ScriptedFieldIdentity(field.FieldName, format, kind)
	// Time lexical values remain capture observations. Common timing is applied
	// by the event owner after this call, so their optional typed view is not
	// reinterpreted as an editable scalar.
	if n == "start" || n == "end" {
		return field.RawValue, nil
	}
	v, err := model.ScriptedScalarValue(field, format, kind, 0)
	if err != nil {
		return "", err
	}
	switch v.Kind {
	case "boolean":
		if *v.Boolean {
			return "-1", nil
		}
		return "0", nil
	case "decimal":
		return strconv.FormatFloat(*v.Decimal, 'f', -1, 64), nil
	case "integer":
		value := strconv.FormatInt(*v.Integer, 10)
		if n == "marked" {
			value = "Marked=" + value
		}
		return value, nil
	case "color":
		c := v.Color
		bits := uint32(c.Alpha)<<24 | uint32(c.Blue)<<16 | uint32(c.Green)<<8 | uint32(c.Red)
		if format == "ass" {
			return fmt.Sprintf("&H%08X", bits), nil
		}
		return strconv.FormatUint(uint64(bits), 10), nil
	case "string":
		return *v.String, nil
	}
	return "", fmt.Errorf("render scripted: inconsistent_typed_value")
}
