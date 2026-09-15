package scripted

import (
	"context"
	"github.com/shruggietech/cueson/internal/model"
	"sort"
	"strconv"
	"strings"
)

type parser struct {
	result          ParseResult
	format          string
	section         model.ScriptedSection
	declaration     *model.ScriptedRecord
	wrap            int
	fields          int
	spans           int
	attachmentBytes int
	seenMetadata    map[string]bool
	seenStyles      map[string]bool
	seenAttachments map[string]bool
}

// Parse retains all accepted native occurrences and derives dialogue semantics.
// It operates on already acquired bytes and never opens or discovers resources.
func Parse(ctx context.Context, data []byte, format string) (ParseResult, error) {
	if ctx == nil {
		return ParseResult{}, failure("invalid_context", 0)
	}
	if err := ctx.Err(); err != nil {
		return ParseResult{}, err
	}
	lines, encoding, err := physicalLines(ctx, data)
	if err != nil {
		return ParseResult{}, err
	}
	detect := detectLines(lines)
	if detect.Err != nil {
		return ParseResult{}, detect.Err
	}
	if !detect.Candidate || detect.Format != format {
		return ParseResult{}, failure("invalid_scripted_dialect", 0)
	}
	p := parser{format: format, result: ParseResult{Encoding: encoding, Cues: []model.Cue{}, Diagnostics: []model.Diagnostic{}, Native: model.ScriptedDocumentData{Dialect: "v4.00+", Sections: []model.ScriptedSection{}, Records: []model.ScriptedRecord{}, Styles: []model.ScriptedStyle{}, Events: []model.ScriptedEvent{}, Attachments: []model.ScriptedAttachment{}}}}
	p.seenMetadata = map[string]bool{}
	p.seenStyles = map[string]bool{}
	p.seenAttachments = map[string]bool{}
	if format == "ssa" {
		p.result.Native.Dialect = "v4.00"
	}
	for order := 0; order < len(lines); order++ {
		if err = ctx.Err(); err != nil {
			return ParseResult{}, err
		}
		line := lines[order]
		if name, ok := sectionHeader(line); ok {
			p.section = model.ScriptedSection{SectionID: occurrence("section", order), SourceOrder: order, Name: name, RawHeader: ptr(line)}
			if err = model.ValidateScriptedNativeSection(name, order); err != nil {
				return ParseResult{}, err
			}
			p.result.Native.Sections = append(p.result.Native.Sections, p.section)
			p.declaration = nil
			switch asciiName(name) {
			case "script info", "v4+ styles", "v4 styles", "events", "fonts", "graphics", "aegisub project garbage":
			default:
				if err = p.warn("unknown_native_section", order); err != nil {
					return ParseResult{}, err
				}
			}
			continue
		}
		if p.section.SectionID == "" {
			return ParseResult{}, failure("unsupported_preamble", order)
		}
		r := model.ScriptedRecord{RecordID: occurrence("record", order), SourceOrder: order, SectionID: p.section.SectionID, Kind: "unknown", RawLine: ptr(line)}
		trim := strings.TrimSpace(line)
		prefix, value, hasPrefix := strings.Cut(strings.TrimLeft(line, " "), ":")
		prefix = asciiName(prefix)
		section := asciiName(p.section.Name)
		if trim == "" {
			r.Kind = "blank"
		} else if strings.HasPrefix(trim, ";") {
			r.Kind = "comment"
		} else if (section == "fonts" && prefix == "fontname") || (section == "graphics" && prefix == "filename") {
			var last int
			last, err = p.attachment(lines, order, &r, value)
			if err != nil {
				return ParseResult{}, err
			}
			order = last
			continue
		} else if (section == "events" || section == "v4+ styles" || section == "v4 styles") && hasPrefix && prefix == "format" {
			r.Kind = "format_declaration"
			names := strings.SplitN(strings.TrimPrefix(value, " "), ",", model.MaxScriptedDeclarationFields+1)
			if len(names) > model.MaxScriptedDeclarationFields {
				return ParseResult{}, failure("complexity_limit", order)
			}
			context := "style"
			if section == "events" {
				context = "event"
			}
			for _, name := range names {
				if nameErr := model.ValidateScriptedNativeField(model.ScriptedField{FieldName: name}, format, context, order); nameErr != nil && strings.Contains(nameErr.Error(), "unsafe_") {
					return ParseResult{}, nameErr
				}
				identity, known := model.ScriptedFieldIdentity(name, format, context)
				f := model.ScriptedDeclarationField{FieldName: name}
				if known {
					f.RecognizedName = ptr(identity)
				}
				r.DeclarationFields = append(r.DeclarationFields, f)
			}
			p.fields += len(names)
			// Declaration syntax must frame every canonical field before it can
			// grant content roles to any subsequent recognized record.
			dummy := make([]string, len(names))
			if context == "event" {
				for i := range dummy {
					dummy[i] = "0"
				}
			}
			_, valid := model.ScriptedSplitFields(strings.Join(dummy, ","), names, format, context)
			if !valid {
				r.Kind = "malformed"
				r.DeclarationFields = nil
				p.declaration = nil
				if err = p.warn("malformed_native_record", order); err != nil {
					return ParseResult{}, err
				}
			} else {
				p.declaration = &r
			}
		} else if hasPrefix && ((section == "events" && (prefix == "dialogue" || prefix == "comment")) || ((section == "v4+ styles" || section == "v4 styles") && prefix == "style")) {
			if err = p.nativeRecord(&r, prefix, value); err != nil {
				return ParseResult{}, err
			}
		} else if hasPrefix && (section == "script info" || section == "aegisub project garbage") {
			r.Kind = "metadata"
			f := model.ScriptedField{FieldName: strings.TrimSpace(strings.SplitN(strings.TrimLeft(line, " "), ":", 2)[0]), RawValue: strings.TrimLeft(value, " ")}
			if err = model.ValidateScriptedNativeField(f, format, "metadata", order); err != nil {
				return ParseResult{}, err
			}
			r.Fields = []model.ScriptedField{f}
			p.fields++
			if prefix == "wrapstyle" {
				p.wrap, _ = strconv.Atoi(strings.TrimSpace(value))
			}
			if prefix == "scripttype" || prefix == "timer" {
				if p.seenMetadata[prefix] {
					if err = p.warn(map[string]string{"scripttype": "duplicate_script_type", "timer": "duplicate_timer"}[prefix], order); err != nil {
						return ParseResult{}, err
					}
				}
				p.seenMetadata[prefix] = true
			}
		} else {
			if err = p.warn("unknown_native_record", order); err != nil {
				return ParseResult{}, err
			}
		}
		if err = model.ValidateScriptedNativeRecord(r, p.section.Name); err != nil {
			return ParseResult{}, err
		}
		p.result.Native.Records = append(p.result.Native.Records, r)
		if p.fields > model.MaxScriptedFields {
			return ParseResult{}, failure("complexity_limit", order)
		}
	}
	if err = p.project(ctx); err != nil {
		return ParseResult{}, err
	}
	sort.SliceStable(p.result.Diagnostics, func(i, j int) bool {
		return *p.result.Diagnostics[i].SourceOrder < *p.result.Diagnostics[j].SourceOrder
	})
	return p.result, nil
}

func (p *parser) warn(code string, order int) error {
	if len(p.result.Diagnostics) >= model.MaxDiagnostics {
		return failure("complexity_limit", order)
	}
	p.result.Diagnostics = append(p.result.Diagnostics, model.Diagnostic{Severity: "warning", Code: code, Message: strings.ReplaceAll(code, "_", " "), SourceOrder: ptr(order)})
	return nil
}

func (p *parser) nativeRecord(r *model.ScriptedRecord, prefix, value string) error {
	context := "event"
	if prefix == "style" {
		context = "style"
	}
	var names []string
	if p.declaration != nil {
		for _, f := range p.declaration.DeclarationFields {
			names = append(names, f.FieldName)
		}
	}
	values, valid := model.ScriptedSplitFields(value, names, p.format, context)
	fields := []model.ScriptedField{}
	if valid {
		for i, name := range names {
			f := model.ScriptedField{FieldName: name, RawValue: values[i]}
			if err := model.ValidateScriptedNativeField(f, p.format, context, r.SourceOrder); err != nil {
				if strings.Contains(err.Error(), "unsafe_") {
					return err
				}
				valid = false
			}
			fields = append(fields, f)
		}
	}
	if valid && context == "event" {
		start, e1 := model.ScriptedMilliseconds(field(fields, "start"))
		end, e2 := model.ScriptedMilliseconds(field(fields, "end"))
		if e1 != nil || e2 != nil || end <= start {
			if prefix == "dialogue" {
				return failure("invalid_native_interval", r.SourceOrder)
			}
			valid = false
		}
	}
	if !valid {
		if prefix == "dialogue" {
			return failure("malformed_native_record", r.SourceOrder)
		}
		if err := model.ValidateScriptedMalformedValue(value, r.SourceOrder); err != nil {
			return err
		}
		if err := p.warn("malformed_native_record", r.SourceOrder); err != nil {
			return err
		}
	}
	p.fields += len(fields)
	if context == "style" {
		id := occurrence("style", r.SourceOrder)
		r.Kind = "style"
		r.StyleID = ptr(id)
		name := field(fields, "name")
		if p.seenStyles[name] {
			if err := p.warn("duplicate_style", r.SourceOrder); err != nil {
				return err
			}
		}
		p.seenStyles[name] = true
		var decl *string
		if p.declaration != nil {
			decl = ptr(p.declaration.RecordID)
		}
		p.result.Native.Styles = append(p.result.Native.Styles, model.ScriptedStyle{StyleID: id, RecordID: r.RecordID, DeclarationID: decl, Name: name, Valid: valid, Fields: fields})
	} else {
		id := occurrence("event", r.SourceOrder)
		r.Kind = "event"
		r.EventID = ptr(id)
		var decl *string
		if p.declaration != nil {
			decl = ptr(p.declaration.RecordID)
		}
		p.result.Native.Events = append(p.result.Native.Events, model.ScriptedEvent{EventID: id, RecordID: r.RecordID, DeclarationID: decl, EventType: prefix, Valid: valid, Fields: fields, Text: field(fields, "text"), Spans: []model.ScriptedSpan{}, Tags: []model.ScriptedTag{}, Karaoke: []model.ScriptedKaraoke{}})
	}
	return nil
}

func field(fields []model.ScriptedField, name string) string {
	for _, f := range fields {
		identity, _ := model.ScriptedFieldIdentity(f.FieldName, "ass", "event")
		if identity == name {
			return f.RawValue
		}
	}
	return ""
}

func (p *parser) attachment(lines []string, order int, r *model.ScriptedRecord, value string) (int, error) {
	name := strings.TrimSpace(value)
	if !model.ScriptedNativeSafeName(name) {
		return order, failure("unsafe_source_metadata", order)
	}
	end := order + 1
	for end < len(lines) {
		if _, ok := sectionHeader(lines[end]); ok {
			break
		}
		prefix, _, ok := strings.Cut(strings.TrimSpace(lines[end]), ":")
		if strings.TrimSpace(lines[end]) == "" || (ok && (asciiName(prefix) == "fontname" || asciiName(prefix) == "filename")) {
			break
		}
		end++
	}
	if end == order+1 {
		r.Kind = "malformed"
		if err := model.ValidateScriptedNativeRecord(*r, p.section.Name); err != nil {
			return order, err
		}
		p.result.Native.Records = append(p.result.Native.Records, *r)
		return order, p.warn("malformed_native_record", order)
	}
	length, malformed, err := model.ScriptedAttachmentFacts(lines[order+1 : end])
	if err != nil {
		return order, err
	}
	p.attachmentBytes += length
	if p.attachmentBytes > 32<<20 || len(p.result.Native.Attachments) >= model.MaxItemOccurrences {
		return order, failure("complexity_limit", order)
	}
	id := occurrence("attachment", order)
	r.Kind = "attachment_header"
	r.AttachmentID = ptr(id)
	if p.seenAttachments[name] {
		if err = p.warn("duplicate_attachment_name", order); err != nil {
			return order, err
		}
	}
	p.seenAttachments[name] = true
	if malformed {
		if err = p.warn("malformed_attachment", order); err != nil {
			return order, err
		}
	}
	kind := "font"
	if asciiName(p.section.Name) == "graphics" {
		kind = "graphic"
	}
	p.result.Native.Attachments = append(p.result.Native.Attachments, model.ScriptedAttachment{AttachmentID: id, HeaderRecordID: r.RecordID, DataStartRecordID: occurrence("record", order+1), DataRecordCount: end - order - 1, AttachmentType: kind, Name: name})
	p.result.Native.Records = append(p.result.Native.Records, *r)
	for pos := order + 1; pos < end; pos++ {
		p.result.Native.Records = append(p.result.Native.Records, model.ScriptedRecord{RecordID: occurrence("record", pos), SourceOrder: pos, SectionID: p.section.SectionID, Kind: "attachment_data", RawLine: ptr(lines[pos])})
	}
	return end - 1, nil
}
