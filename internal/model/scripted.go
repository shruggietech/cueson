package model

import (
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ScriptedDocumentData retains a complete ordered ASS v4+ or SSA v4 document.
// Recognition is distinct from the availability of a native codec.
type ScriptedDocumentData struct {
	Dialect     string               `json:"dialect"`
	Sections    []ScriptedSection    `json:"sections"`
	Records     []ScriptedRecord     `json:"records"`
	Styles      []ScriptedStyle      `json:"styles"`
	Events      []ScriptedEvent      `json:"events"`
	Attachments []ScriptedAttachment `json:"attachments"`
}
type ScriptedSection struct {
	SectionID   string  `json:"section_id"`
	SourceOrder int     `json:"source_order"`
	Name        string  `json:"name"`
	RawHeader   *string `json:"raw_header,omitempty"`
}
type ScriptedRecord struct {
	RecordID          string                     `json:"record_id"`
	SourceOrder       int                        `json:"source_order"`
	SectionID         string                     `json:"section_id"`
	Kind              string                     `json:"kind"`
	RawLine           *string                    `json:"raw_line,omitempty"`
	Fields            []ScriptedField            `json:"fields,omitempty"`
	DeclarationFields []ScriptedDeclarationField `json:"declaration_fields,omitempty"`
	StyleID           *string                    `json:"style_id,omitempty"`
	EventID           *string                    `json:"event_id,omitempty"`
	AttachmentID      *string                    `json:"attachment_id,omitempty"`
}
type ScriptedDeclarationField struct {
	FieldName      string  `json:"field_name"`
	RecognizedName *string `json:"recognized_name,omitempty"`
}
type ScriptedField struct {
	FieldName  string         `json:"field_name"`
	RawValue   string         `json:"raw_value"`
	TypedValue *ScriptedValue `json:"typed_value,omitempty"`
}
type ScriptedValue struct {
	Kind    string         `json:"kind"`
	Boolean *bool          `json:"boolean,omitempty"`
	Integer *int64         `json:"integer,omitempty"`
	Decimal *float64       `json:"decimal,omitempty"`
	Color   *ScriptedColor `json:"color,omitempty"`
	String  *string        `json:"string,omitempty"`
}
type ScriptedColor struct {
	Alpha uint8 `json:"alpha"`
	Blue  uint8 `json:"blue"`
	Green uint8 `json:"green"`
	Red   uint8 `json:"red"`
}
type ScriptedStyle struct {
	StyleID       string          `json:"style_id"`
	RecordID      string          `json:"record_id"`
	DeclarationID *string         `json:"declaration_id,omitempty"`
	Name          string          `json:"name"`
	Valid         bool            `json:"valid"`
	Fields        []ScriptedField `json:"fields"`
}
type ScriptedEvent struct {
	EventID       string            `json:"event_id"`
	RecordID      string            `json:"record_id"`
	DeclarationID *string           `json:"declaration_id,omitempty"`
	EventType     string            `json:"event_type"`
	CueID         *string           `json:"cue_id,omitempty"`
	Valid         bool              `json:"valid"`
	Fields        []ScriptedField   `json:"fields"`
	Text          string            `json:"text"`
	Spans         []ScriptedSpan    `json:"spans"`
	Tags          []ScriptedTag     `json:"tags"`
	Karaoke       []ScriptedKaraoke `json:"karaoke"`
}
type ScriptedSpan struct {
	Kind        string `json:"kind"`
	StartScalar int    `json:"start_scalar"`
	EndScalar   int    `json:"end_scalar"`
	Raw         string `json:"raw"`
}
type ScriptedTag struct {
	Name        string `json:"name"`
	Parameter   string `json:"parameter"`
	StartScalar int    `json:"start_scalar"`
	EndScalar   int    `json:"end_scalar"`
	Raw         string `json:"raw"`
}
type ScriptedKaraoke struct {
	Variant              string `json:"variant"`
	StartScalar          int    `json:"start_scalar"`
	EndScalar            int    `json:"end_scalar"`
	DurationCentiseconds *int64 `json:"duration_centiseconds,omitempty"`
	Supported            bool   `json:"supported"`
}
type ScriptedAttachment struct {
	AttachmentID      string `json:"attachment_id"`
	HeaderRecordID    string `json:"header_record_id"`
	DataStartRecordID string `json:"data_start_record_id"`
	DataRecordCount   int    `json:"data_record_count"`
	AttachmentType    string `json:"attachment_type"`
	Name              string `json:"name"`
}
type ScriptedCueData struct {
	EventID           string             `json:"event_id"`
	StyleID           string             `json:"style_id"`
	StartTimestampRaw *string            `json:"start_timestamp_raw,omitempty"`
	EndTimestampRaw   *string            `json:"end_timestamp_raw,omitempty"`
	Projection        ScriptedProjection `json:"projection"`
}
type ScriptedProjection struct {
	Origin           string  `json:"origin"`
	SourceEventID    string  `json:"source_event_id"`
	TextFieldName    string  `json:"text_field_name"`
	SpeakerFieldName *string `json:"speaker_field_name,omitempty"`
	WrapStyle        int     `json:"wrap_style"`
	DrawingExcluded  bool    `json:"drawing_excluded"`
}

const (
	MaxScriptedLineBytes         = 1 << 20
	MaxScriptedFields            = 1_000_000
	MaxScriptedSpans             = 1_000_000
	MaxScriptedDeclarationFields = 128
	MaxScriptedDepth             = 32
)

var scriptedID = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
var scriptedTime = regexp.MustCompile(`^([0-9]+):([0-5][0-9]):([0-5][0-9])\.([0-9]{2})$`)
var metadataIdentity = regexp.MustCompile(`(?i)(?:[a-z]:[\\/]|\\\\|file:|https?://|/(?:home|users|tmp|var|etc)/|localhost|(?:host|machine|user)(?:name|id)\s*[:=])`)
var assStyleNames = strings.Split("Name,Fontname,Fontsize,PrimaryColour,SecondaryColour,OutlineColour,BackColour,Bold,Italic,Underline,StrikeOut,ScaleX,ScaleY,Spacing,Angle,BorderStyle,Outline,Shadow,Alignment,MarginL,MarginR,MarginV,Encoding", ",")
var ssaStyleNames = strings.Split("Name,Fontname,Fontsize,PrimaryColour,SecondaryColour,TertiaryColour,BackColour,Bold,Italic,BorderStyle,Outline,Shadow,Alignment,MarginL,MarginR,MarginV,AlphaLevel,Encoding", ",")
var assEventNames = strings.Split("Layer,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text", ",")
var ssaEventNames = strings.Split("Marked,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text", ",")

func nativeName(s string) string {
	s = asciiLower(strings.Trim(s, " "))
	if s == "actor" {
		return "name"
	}
	return s
}
func asciiLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + ('a' - 'A')
		}
	}
	return string(b)
}
func scriptedError(code string, order int) error {
	return fmt.Errorf("%s at scripted source_order %d", code, order)
}
func scriptedPhysical(s string) bool {
	return utf8.ValidString(s) && len(s) <= MaxScriptedLineBytes && !strings.ContainsAny(s, "\r\n\x00")
}
func nativeID(s string) bool { return scriptedID.MatchString(s) }
func nativeSafeName(s string) bool {
	if s == "" || s == "." || s == ".." || strings.ContainsAny(s, `\/:*?"<>|`) || strings.TrimRight(s, " .") != s {
		return false
	}
	for _, r := range s {
		if r < 32 || r == 127 {
			return false
		}
	}
	stem := strings.ToUpper(strings.SplitN(s, ".", 2)[0])
	if slices.Contains([]string{"CON", "PRN", "AUX", "NUL", "CONIN$", "CONOUT$"}, stem) {
		return false
	}
	return !(len(stem) == 4 && (strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT")) && stem[3] >= '1' && stem[3] <= '9')
}

func validateScriptedDocument(doc Document) error {
	if err := validateCollectionLimits(doc); err != nil {
		return err
	}
	want := FormatSupport{Status: "schema_only", RestoreSupported: true}
	if doc.FormatSupport != want {
		return fmt.Errorf("format_support must declare schema_only scripted recognition")
	}
	native := doc.FormatData.ASS
	dialect := "v4.00+"
	if doc.Format == "ssa" {
		native = doc.FormatData.SSA
		dialect = "v4.00"
	}
	if native == nil || doc.FormatData.SubRip != nil || doc.FormatData.WebVTT != nil || (doc.Format == "ass" && doc.FormatData.SSA != nil) || (doc.Format == "ssa" && doc.FormatData.ASS != nil) {
		return fmt.Errorf("format_data must contain only matching scripted branch")
	}
	if native.Dialect != dialect {
		return fmt.Errorf("scripted dialect does not match format")
	}
	if native.Sections == nil || native.Records == nil || native.Styles == nil || native.Events == nil || native.Attachments == nil {
		return fmt.Errorf("scripted collections must not be null")
	}
	if len(native.Sections) == 0 || len(native.Sections)+len(native.Records) > MaxDocumentItems || len(native.Attachments) > MaxItemOccurrences {
		return fmt.Errorf("complexity_limit: scripted document items")
	}
	assets, err := validateAssets(doc.Source)
	if err != nil {
		return err
	}
	if err = scriptedSourcePrivacy(doc.Source, doc.Format); err != nil {
		return err
	}
	diagnosticIndex := map[string]bool{}
	for _, diagnostic := range doc.Diagnostics {
		if diagnostic.SourceOrder != nil {
			diagnosticIndex[diagnostic.Code+":"+strconv.Itoa(*diagnostic.SourceOrder)] = true
		}
	}
	hasDiag := func(code string, order int) bool { return diagnosticIndex[code+":"+strconv.Itoa(order)] }
	wrapStyle := 0
	for _, value := range []*string{doc.Metadata.Title, doc.Metadata.Description, doc.Metadata.Language, doc.Metadata.Kind} {
		if value != nil && metadataIdentity.MatchString(*value) {
			return fmt.Errorf("unsafe_source_metadata in document metadata")
		}
	}
	sections := map[string]ScriptedSection{}
	records := map[string]ScriptedRecord{}
	orders := make([]string, len(native.Sections)+len(native.Records))
	ids := map[string]bool{}
	add := func(id string, order int, kind string) error {
		if !nativeID(id) || ids[id] || order < 0 || order >= len(orders) || orders[order] != "" {
			return scriptedError("invalid_native_ownership", order)
		}
		ids[id] = true
		orders[order] = kind + ":" + id
		return nil
	}
	last := -1
	for _, s := range native.Sections {
		if err = add(s.SectionID, s.SourceOrder, "section"); err != nil {
			return err
		}
		if s.SourceOrder <= last || !scriptedPhysical(s.Name) || strings.ContainsAny(s.Name, "[]") || strings.TrimSpace(s.Name) == "" {
			return scriptedError("invalid_native_section", s.SourceOrder)
		}
		last = s.SourceOrder
		if metadataIdentity.MatchString(s.Name) {
			return scriptedError("unsafe_source_metadata", s.SourceOrder)
		}
		if s.RawHeader != nil && (!scriptedPhysical(*s.RawHeader) || metadataIdentity.MatchString(*s.RawHeader)) {
			return scriptedError("unsafe_source_metadata", s.SourceOrder)
		}
		sections[s.SectionID] = s
	}
	for _, section := range native.Sections {
		if (doc.Format == "ass" && strings.EqualFold(section.Name, "V4 Styles")) || (doc.Format == "ssa" && strings.EqualFold(section.Name, "V4+ Styles")) {
			return scriptedError("mismatched_native_dialect", section.SourceOrder)
		}
	}
	last = -1
	totalFields, totalSpans := 0, 0
	seenScriptType := false
	kinds := []string{"metadata", "format_declaration", "style", "event", "attachment_header", "attachment_data", "blank", "comment", "unknown", "malformed"}
	for _, r := range native.Records {
		if err = add(r.RecordID, r.SourceOrder, "record"); err != nil {
			return err
		}
		if r.SourceOrder <= last {
			return scriptedError("invalid_native_order", r.SourceOrder)
		}
		last = r.SourceOrder
		if _, ok := sections[r.SectionID]; !ok || !slices.Contains(kinds, r.Kind) {
			return scriptedError("invalid_native_ownership", r.SourceOrder)
		}
		if r.RawLine != nil && !scriptedPhysical(*r.RawLine) {
			return scriptedError("invalid_native_line", r.SourceOrder)
		}
		if slices.Contains([]string{"blank", "comment", "unknown", "malformed", "attachment_data"}, r.Kind) && r.RawLine == nil {
			return scriptedError("missing_capture_content", r.SourceOrder)
		}
		if len(r.Fields) > MaxItemOccurrences || len(r.DeclarationFields) > MaxScriptedDeclarationFields {
			return scriptedError("complexity_limit", r.SourceOrder)
		}
		totalFields += len(r.Fields) + len(r.DeclarationFields)
		if r.Kind != "format_declaration" && len(r.DeclarationFields) > 0 {
			return scriptedError("invalid_declaration_ownership", r.SourceOrder)
		}
		if r.Kind == "format_declaration" && len(r.DeclarationFields) == 0 {
			return scriptedError("invalid_declaration", r.SourceOrder)
		}
		if (r.Kind == "style") != (r.StyleID != nil) || (r.Kind == "event") != (r.EventID != nil) || (r.Kind == "attachment_header") != (r.AttachmentID != nil) {
			return scriptedError("invalid_native_ownership", r.SourceOrder)
		}
		if r.Kind != "metadata" && len(r.Fields) > 0 {
			return scriptedError("competing_field_ownership", r.SourceOrder)
		}
		for _, field := range r.Fields {
			if strings.EqualFold(strings.TrimSpace(field.FieldName), "ScriptType") {
				if !strings.EqualFold(sections[r.SectionID].Name, "Script Info") {
					return scriptedError("mismatched_native_dialect", r.SourceOrder)
				}
				seenScriptType = true
			}
			if strings.EqualFold(strings.TrimSpace(field.FieldName), "ScriptType") && strings.TrimSpace(field.RawValue) != dialect {
				return scriptedError("mismatched_native_dialect", r.SourceOrder)
			}
			if err = validateScriptedScalar(field, doc.Format, r.SourceOrder); err != nil {
				return err
			}
			if strings.EqualFold(strings.TrimSpace(field.FieldName), "WrapStyle") {
				wrapStyle, _ = strconv.Atoi(strings.TrimSpace(field.RawValue))
			}
		}
		if err = scriptedRecordPrivacy(r, sections[r.SectionID].Name); err != nil {
			return err
		}
		if r.Kind == "malformed" && !hasDiag("malformed_native_record", r.SourceOrder) {
			return scriptedError("missing_native_diagnostic", r.SourceOrder)
		}
		records[r.RecordID] = r
	}
	if !seenScriptType {
		return fmt.Errorf("mismatched_native_dialect: native Script Info requires ScriptType")
	}
	activeSection := ""
	activeDeclaration := ""
	declarationAtRecord := map[string]string{}
	for order, item := range orders {
		if item == "" {
			return scriptedError("invalid_native_order", order)
		}
		kind, id, _ := strings.Cut(item, ":")
		if kind == "section" {
			activeSection = id
			activeDeclaration = ""
			continue
		}
		r := records[id]
		if r.SectionID != activeSection {
			return scriptedError("cross_section_ownership", order)
		}
		if r.Kind == "format_declaration" {
			activeDeclaration = r.RecordID
		}
		declarationAtRecord[r.RecordID] = activeDeclaration
	}
	styleNames := assStyleNames
	eventNames := assEventNames
	styleSection := "v4+ styles"
	if doc.Format == "ssa" {
		styleNames = ssaStyleNames
		eventNames = ssaEventNames
		styleSection = "v4 styles"
	}
	validDeclaration := func(record ScriptedRecord, ref *string, names []string, fields []ScriptedField) error {
		if ref == nil {
			if record.RawLine != nil {
				return scriptedError("missing_declaration", record.SourceOrder)
			}
			return validateScriptedFields(fields, names, nil, record.SourceOrder)
		}
		decl, ok := records[*ref]
		if !ok || decl.Kind != "format_declaration" || decl.SectionID != record.SectionID || decl.SourceOrder >= record.SourceOrder {
			return scriptedError("invalid_declaration_reference", record.SourceOrder)
		}
		if declarationAtRecord[record.RecordID] != *ref {
			return scriptedError("stale_declaration_reference", record.SourceOrder)
		}
		return validateScriptedFields(fields, names, decl.DeclarationFields, record.SourceOrder)
	}
	styles := map[string]ScriptedStyle{}
	styleByName := map[string][]string{}
	last = -1
	for _, s := range native.Styles {
		r, ok := records[s.RecordID]
		if !ok || r.Kind != "style" || r.StyleID == nil || *r.StyleID != s.StyleID || ids[s.StyleID] || !nativeID(s.StyleID) || r.SourceOrder <= last || strings.ToLower(sections[r.SectionID].Name) != styleSection {
			return scriptedError("invalid_style_ownership", r.SourceOrder)
		}
		ids[s.StyleID] = true
		last = r.SourceOrder
		if s.Fields == nil || len(s.Fields) > MaxItemOccurrences {
			return scriptedError("complexity_limit", r.SourceOrder)
		}
		totalFields += len(s.Fields)
		if s.Valid {
			if err = validDeclaration(r, s.DeclarationID, styleNames, s.Fields); err != nil {
				return err
			}
			if fieldRaw(s.Fields, "name") != s.Name {
				return scriptedError("inconsistent_style_name", r.SourceOrder)
			}
		} else if r.RawLine == nil || !hasDiag("malformed_native_record", r.SourceOrder) {
			return scriptedError("missing_native_diagnostic", r.SourceOrder)
		}
		for _, f := range s.Fields {
			if s.Valid || f.TypedValue != nil {
				err = validateScriptedScalar(f, doc.Format, r.SourceOrder)
			} else {
				err = validateRetainedScriptedField(f, r.SourceOrder)
			}
			if err != nil {
				return err
			}
		}
		styles[s.StyleID] = s
		styleByName[s.Name] = append(styleByName[s.Name], s.StyleID)
	}
	events := map[string]ScriptedEvent{}
	last = -1
	for _, e := range native.Events {
		r, ok := records[e.RecordID]
		if !ok || r.Kind != "event" || r.EventID == nil || *r.EventID != e.EventID || ids[e.EventID] || !nativeID(e.EventID) || r.SourceOrder <= last || strings.ToLower(sections[r.SectionID].Name) != "events" {
			return scriptedError("invalid_event_ownership", r.SourceOrder)
		}
		ids[e.EventID] = true
		last = r.SourceOrder
		if !scriptedPhysical(e.Text) || e.Fields == nil || e.Spans == nil || e.Tags == nil || e.Karaoke == nil || len(e.Fields) > MaxItemOccurrences || len(e.Spans) > MaxItemOccurrences || len(e.Tags) > MaxItemOccurrences || len(e.Karaoke) > MaxItemOccurrences {
			return scriptedError("complexity_limit", r.SourceOrder)
		}
		totalFields += len(e.Fields)
		totalSpans += len(e.Spans) + len(e.Tags)
		if strings.EqualFold(e.EventType, "command") || strings.Contains(strings.ToLower(e.EventType), "template") {
			return scriptedError("unsafe_active_content", r.SourceOrder)
		}
		if e.EventType == "dialogue" || e.EventType == "comment" {
			if e.Valid {
				if err = validDeclaration(r, e.DeclarationID, eventNames, e.Fields); err != nil {
					return err
				}
				if fieldRaw(e.Fields, "text") != e.Text {
					return scriptedError("inconsistent_projection", r.SourceOrder)
				}
			} else if e.EventType == "dialogue" {
				return scriptedError("malformed_native_record", r.SourceOrder)
			}
		}
		if !e.Valid && (r.RawLine == nil || !hasDiag("malformed_native_record", r.SourceOrder)) {
			return scriptedError("missing_native_diagnostic", r.SourceOrder)
		}
		if (e.EventType == "dialogue") != (e.CueID != nil) {
			return scriptedError("invalid_dialogue_ownership", r.SourceOrder)
		}
		for _, f := range e.Fields {
			if !e.Valid && f.TypedValue == nil {
				if err = validateRetainedScriptedField(f, r.SourceOrder); err != nil {
					return err
				}
				continue
			}
			if nativeName(f.FieldName) != "start" && nativeName(f.FieldName) != "end" {
				if err = validateScriptedScalar(f, doc.Format, r.SourceOrder); err != nil {
					return err
				}
			} else if _, err = scriptedMilliseconds(f.RawValue); err != nil {
				return scriptedError("invalid_native_timestamp", r.SourceOrder)
			} else if f.TypedValue != nil && (f.TypedValue.Kind != "string" || f.TypedValue.String == nil || *f.TypedValue.String != f.RawValue || f.TypedValue.Boolean != nil || f.TypedValue.Integer != nil || f.TypedValue.Decimal != nil || f.TypedValue.Color != nil) {
				return scriptedError("inconsistent_typed_value", r.SourceOrder)
			}
		}
		events[e.EventID] = e
	}
	if len(styles) != countScriptedKind(native.Records, "style") || len(events) != countScriptedKind(native.Records, "event") {
		return fmt.Errorf("orphaned native record owner")
	}
	if totalFields > MaxScriptedFields || totalSpans > MaxScriptedSpans {
		return fmt.Errorf("complexity_limit: scripted aggregate occurrences")
	}
	for _, event := range native.Events {
		if err = validateScriptedNativeOffsets(event, records[event.RecordID].SourceOrder); err != nil {
			return err
		}
		if event.EventType != "dialogue" && event.Valid && event.EventType == "comment" {
			start, _ := scriptedMilliseconds(fieldRaw(event.Fields, "start"))
			end, _ := scriptedMilliseconds(fieldRaw(event.Fields, "end"))
			facts, projectionErr := ProjectScriptedText(event.Text, wrapStyle, start, end)
			order := records[event.RecordID].SourceOrder
			if projectionErr != nil || !slices.Equal(event.Spans, facts.Spans) || !slices.Equal(event.Tags, facts.Tags) || !equalScriptedKaraoke(event.Karaoke, facts.Karaoke) {
				return scriptedError("inconsistent_projection", order)
			}
			for _, code := range facts.DiagnosticCodes {
				if !hasDiag(code, order) {
					return scriptedError("missing_native_diagnostic", order)
				}
			}
		}
	}
	for _, diagnostic := range doc.Diagnostics {
		if diagnostic.SourceOrder != nil && (*diagnostic.SourceOrder < 0 || *diagnostic.SourceOrder >= len(orders)) {
			return fmt.Errorf("invalid scripted diagnostic source_order")
		}
		if metadataIdentity.MatchString(diagnostic.Message) {
			return fmt.Errorf("unsafe_source_metadata in diagnostic message")
		}
	}
	if err = validateScriptedAttachments(native, records, ids, diagnosticIndex); err != nil {
		return err
	}
	seenEvents := map[string]bool{}
	for i := range doc.Cues {
		c := &doc.Cues[i]
		n := c.FormatData.ASS
		if doc.Format == "ssa" {
			n = c.FormatData.SSA
		}
		if n == nil || c.FormatData.SubRip != nil || c.FormatData.WebVTT != nil || (doc.Format == "ass" && c.FormatData.SSA != nil) || (doc.Format == "ssa" && c.FormatData.ASS != nil) {
			return fmt.Errorf("cues[%d].format_data does not match scripted format", i)
		}
		e, ok := events[n.EventID]
		r := records[e.RecordID]
		if !ok || e.EventType != "dialogue" || e.CueID == nil || *e.CueID != c.ID || seenEvents[e.EventID] || c.SourceOrder != r.SourceOrder {
			return scriptedError("invalid_dialogue_ownership", c.SourceOrder)
		}
		seenEvents[e.EventID] = true
		s, ok := styles[n.StyleID]
		if !ok || !s.Valid || s.Name != fieldRaw(e.Fields, "style") || len(styleByName[s.Name]) != 1 {
			return scriptedError("unresolved_style", c.SourceOrder)
		}
		if c.Timing.StartMilliseconds < 0 || c.Timing.EndMilliseconds <= c.Timing.StartMilliseconds {
			return scriptedError("invalid_native_timestamp", c.SourceOrder)
		}
		for _, raw := range []*string{n.StartTimestampRaw, n.EndTimestampRaw} {
			if raw != nil {
				if _, err = scriptedMilliseconds(*raw); err != nil {
					return scriptedError("invalid_native_timestamp", c.SourceOrder)
				}
			}
		}
		if err = validateScriptedProjection(c, e, n, wrapStyle, diagnosticIndex); err != nil {
			return err
		}
	}
	if len(seenEvents) != countScriptedDialogue(native.Events) {
		return fmt.Errorf("dialogue-to-cue ownership must be one-to-one")
	}
	if err = validateCues(doc.Cues, assets, DocumentFormatData{}); err != nil {
		return err
	}
	if len(doc.Cues) == 0 {
		if doc.Document.CueCount != 0 || doc.Stats.CueCount != 0 || doc.Document.MediaStartMilliseconds != nil || doc.Document.MediaEndMilliseconds != nil || doc.Document.MediaSpanMilliseconds != nil || doc.Stats.MediaSpanMilliseconds != nil || doc.Document.HasWordLevelTiming || doc.Stats.HasWordLevelTiming {
			return fmt.Errorf("empty scripted summaries must be null/zero")
		}
	} else if err = validateSummaries(doc); err != nil {
		return err
	}
	return validateDiagnostics(doc.Diagnostics, doc.Cues, doc.Stats)
}

func validateScriptedNativeOffsets(event ScriptedEvent, order int) error {
	r := []rune(event.Text)
	previousEnd := 0
	for _, span := range event.Spans {
		if span.StartScalar < previousEnd || span.EndScalar < span.StartScalar || span.EndScalar > len(r) || span.Raw != string(r[span.StartScalar:span.EndScalar]) {
			return scriptedError("inconsistent_projection", order)
		}
		previousEnd = span.EndScalar
	}
	previousStart := -1
	tags := map[string]ScriptedTag{}
	for _, tag := range event.Tags {
		if tag.StartScalar < 0 || tag.StartScalar < previousStart || tag.EndScalar < tag.StartScalar || tag.EndScalar > len(r) || tag.Raw != string(r[tag.StartScalar:tag.EndScalar]) || tag.Raw != "\\"+tag.Name+tag.Parameter {
			return scriptedError("inconsistent_projection", order)
		}
		previousStart = tag.StartScalar
		tags[strconv.Itoa(tag.StartScalar)+":"+strconv.Itoa(tag.EndScalar)] = tag
	}
	previousStart = -1
	for _, karaoke := range event.Karaoke {
		tag, ok := tags[strconv.Itoa(karaoke.StartScalar)+":"+strconv.Itoa(karaoke.EndScalar)]
		if !ok || tag.Name != karaoke.Variant || karaoke.StartScalar < previousStart || !slices.Contains([]string{"k", "K", "kf", "ko", "kt"}, karaoke.Variant) {
			return scriptedError("inconsistent_projection", order)
		}
		previousStart = karaoke.StartScalar
		duration, parseErr := strconv.ParseInt(strings.TrimSpace(tag.Parameter), 10, 64)
		if karaoke.DurationCentiseconds != nil && (parseErr != nil || *karaoke.DurationCentiseconds != duration) {
			return scriptedError("inconsistent_projection", order)
		}
		if karaoke.Supported && (parseErr != nil || duration < 0 || karaoke.Variant == "kt" || karaoke.DurationCentiseconds == nil) {
			return scriptedError("inconsistent_projection", order)
		}
	}
	return nil
}
func countScriptedKind(records []ScriptedRecord, kind string) int {
	n := 0
	for _, r := range records {
		if r.Kind == kind {
			n++
		}
	}
	return n
}
func countScriptedDialogue(events []ScriptedEvent) int {
	n := 0
	for _, e := range events {
		if e.EventType == "dialogue" {
			n++
		}
	}
	return n
}
func fieldRaw(fields []ScriptedField, name string) string {
	for _, f := range fields {
		if nativeName(f.FieldName) == name {
			return f.RawValue
		}
	}
	return ""
}
func validateScriptedFields(fields []ScriptedField, names []string, decl []ScriptedDeclarationField, order int) error {
	if len(fields) < len(names) || len(fields) > MaxScriptedDeclarationFields {
		return scriptedError("invalid_native_fields", order)
	}
	required := map[string]bool{}
	for _, n := range names {
		required[nativeName(n)] = false
	}
	if decl != nil && len(fields) != len(decl) {
		return scriptedError("invalid_native_fields", order)
	}
	for i, f := range fields {
		n := nativeName(f.FieldName)
		if n == "" || !scriptedPhysical(f.FieldName) || strings.Contains(f.FieldName, ",") {
			return scriptedError("invalid_native_fields", order)
		}
		if decl != nil {
			if f.FieldName != decl[i].FieldName {
				return scriptedError("inconsistent_declaration_fields", order)
			}
			_, known := required[n]
			if known != (decl[i].RecognizedName != nil) || (known && *decl[i].RecognizedName != n) {
				return scriptedError("inconsistent_recognized_name", order)
			}
		}
		if seen, known := required[n]; known {
			if seen {
				return scriptedError("duplicate_native_field", order)
			}
			required[n] = true
		}
		if n == "text" && i != len(fields)-1 {
			return scriptedError("nonfinal_native_text", order)
		}
	}
	for _, present := range required {
		if !present {
			return scriptedError("missing_native_field", order)
		}
	}
	return nil
}
func scriptedMilliseconds(s string) (int64, error) {
	m := scriptedTime.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid timestamp")
	}
	h, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil || h > (math.MaxInt64-3_599_990)/3_600_000 {
		return 0, fmt.Errorf("timestamp overflow")
	}
	mm, _ := strconv.ParseInt(m[2], 10, 64)
	ss, _ := strconv.ParseInt(m[3], 10, 64)
	cc, _ := strconv.ParseInt(m[4], 10, 64)
	return h*3_600_000 + mm*60_000 + ss*1000 + cc*10, nil
}
func validateScriptedScalar(f ScriptedField, format string, order int) error {
	n := nativeName(f.FieldName)
	if err := scriptedContentFieldPrivacy(f.FieldName, f.RawValue, order); err != nil {
		return err
	}
	if !scriptedPhysical(f.RawValue) || (n != "text" && strings.Contains(f.RawValue, ",")) {
		return scriptedError("invalid_native_field_value", order)
	}
	if n == "text" {
		if f.TypedValue != nil && (f.TypedValue.Kind != "string" || f.TypedValue.String == nil || *f.TypedValue.String != f.RawValue) {
			return scriptedError("inconsistent_typed_value", order)
		}
		return nil
	}
	raw := strings.TrimSpace(f.RawValue)
	kind := "string"
	var iv int64
	var dv float64
	var bv bool
	var cv ScriptedColor
	var err error
	switch n {
	case "bold", "italic", "underline", "strikeout":
		kind = "boolean"
		iv, err = strconv.ParseInt(raw, 10, 64)
		if iv != 0 && iv != 1 && iv != -1 {
			err = fmt.Errorf("invalid boolean")
		}
		bv = iv != 0
	case "fontsize", "scalex", "scaley", "spacing", "angle", "outline", "shadow":
		kind = "decimal"
		dv, err = strconv.ParseFloat(raw, 64)
		if math.IsNaN(dv) || math.IsInf(dv, 0) {
			err = fmt.Errorf("invalid decimal")
		}
	case "layer", "marked", "borderstyle", "alignment", "marginl", "marginr", "marginv", "encoding", "alphalevel":
		kind = "integer"
		if n == "marked" {
			raw = strings.TrimPrefix(raw, "Marked=")
		}
		iv, err = strconv.ParseInt(raw, 10, 64)
	case "primarycolour", "secondarycolour", "outlinecolour", "tertiarycolour", "backcolour":
		kind = "color"
		var value uint64
		if strings.HasPrefix(strings.ToUpper(raw), "&H") {
			value, err = strconv.ParseUint(strings.TrimSuffix(raw[2:], "&"), 16, 32)
		} else if strings.HasPrefix(raw, "-") {
			// SSA signed colors represent the same 32 native ABGR bits. Reject
			// overflow before explicitly reinterpreting that equal-width value.
			var signed int64
			signed, err = strconv.ParseInt(raw, 10, 32)
			if err != nil {
				return scriptedError("invalid_native_field_value", order)
			}
			value = uint64(uint32(int32(signed)))
		} else {
			value, err = strconv.ParseUint(strings.TrimPrefix(raw, "+"), 10, 32)
		}
		if err != nil {
			return scriptedError("invalid_native_field_value", order)
		}
		cv = ScriptedColor{uint8(value >> 24), uint8(value >> 16), uint8(value >> 8), uint8(value)}
	}
	if err != nil {
		return scriptedError("invalid_native_field_value", order)
	}
	if f.TypedValue == nil {
		return nil
	}
	v := f.TypedValue
	count := 0
	for _, present := range []bool{v.Boolean != nil, v.Integer != nil, v.Decimal != nil, v.Color != nil, v.String != nil} {
		if present {
			count++
		}
	}
	if count != 1 || v.Kind != kind {
		return scriptedError("inconsistent_typed_value", order)
	}
	equal := false
	switch kind {
	case "boolean":
		equal = v.Boolean != nil && *v.Boolean == bv
	case "integer":
		equal = v.Integer != nil && *v.Integer == iv
	case "decimal":
		equal = v.Decimal != nil && *v.Decimal == dv
	case "color":
		equal = v.Color != nil && *v.Color == cv
	case "string":
		equal = v.String != nil && *v.String == f.RawValue
	}
	if !equal {
		return scriptedError("inconsistent_typed_value", order)
	}
	return nil
}

func validateRetainedScriptedField(f ScriptedField, order int) error {
	if err := scriptedContentFieldPrivacy(f.FieldName, f.RawValue, order); err != nil {
		return err
	}
	if !scriptedPhysical(f.FieldName) || !scriptedPhysical(f.RawValue) || (nativeName(f.FieldName) != "text" && strings.Contains(f.RawValue, ",")) {
		return scriptedError("invalid_native_field_value", order)
	}
	return nil
}

func scriptedContentFieldPrivacy(name, value string, order int) error {
	n := nativeName(name)
	if strings.Contains(n, "automation") || strings.Contains(n, "script execution") || (n == "effect" && (strings.Contains(strings.ToLower(value), "template") || strings.Contains(strings.ToLower(value), "!code"))) {
		if strings.TrimSpace(value) != "" {
			return scriptedError("unsafe_active_content", order)
		}
		return nil
	}
	if slices.Contains([]string{"audio uri", "audio file", "video file", "timecodes file", "keyframes file", "last style storage", "computer", "hostname", "machine", "username", "user id", "machine id"}, n) && strings.TrimSpace(value) != "" {
		return scriptedError("unsafe_source_metadata", order)
	}
	return nil
}

func scriptedMetadata(name, value string, order int) error {
	key := strings.ToLower(strings.TrimSpace(name))
	if strings.Contains(key, "automation") || strings.Contains(key, "script execution") {
		if strings.TrimSpace(value) != "" {
			return scriptedError("unsafe_active_content", order)
		}
		return nil
	}
	resource := []string{"audio uri", "audio file", "video file", "timecodes file", "keyframes file", "last style storage", "computer", "hostname", "machine", "username", "user id", "machine id"}
	if slices.Contains(resource, key) {
		if strings.TrimSpace(value) != "" {
			return scriptedError("unsafe_source_metadata", order)
		}
		return nil
	}
	allowed := []string{"title", "scripttype", "collisions", "playresx", "playresy", "playdepth", "timer", "wrapstyle", "scaledborderandshadow", "ycbcr matrix", "synch point"}
	if !slices.Contains(allowed, key) || metadataIdentity.MatchString(value) {
		return scriptedError("unsafe_source_metadata", order)
	}
	if key == "timer" {
		v, e := strconv.ParseFloat(strings.TrimSpace(value), 64)
		if e != nil || v != 100 {
			return scriptedError("unsupported_timer", order)
		}
	}
	if key == "wrapstyle" {
		v, e := strconv.Atoi(strings.TrimSpace(value))
		if e != nil || v < 0 || v > 3 {
			return scriptedError("inconsistent_projection", order)
		}
	}
	return nil
}
func scriptedRecordPrivacy(r ScriptedRecord, section string) error {
	s := strings.ToLower(section)
	content := s == "events" || s == "v4 styles" || s == "v4+ styles" || s == "fonts" || s == "graphics"
	if r.Kind == "metadata" {
		if len(r.Fields) == 0 {
			return scriptedError("unsafe_source_metadata", r.SourceOrder)
		}
		for _, f := range r.Fields {
			if err := scriptedMetadata(f.FieldName, f.RawValue, r.SourceOrder); err != nil {
				return err
			}
		}
	}
	if r.RawLine == nil {
		return nil
	}
	raw := *r.RawLine
	prefix, _, hasPrefix := strings.Cut(strings.TrimSpace(raw), ":")
	if hasPrefix && (strings.EqualFold(prefix, "Command") || strings.Contains(strings.ToLower(prefix), "automation") || strings.Contains(strings.ToLower(prefix), "template")) {
		return scriptedError("unsafe_active_content", r.SourceOrder)
	}
	trimmed := strings.TrimSpace(raw)
	if content && (trimmed == "" || strings.HasPrefix(trimmed, ";")) {
		return nil
	}
	if content && (r.Kind == "attachment_data" || (r.Kind == "style" && hasPrefix && strings.EqualFold(prefix, "Style")) || (r.Kind == "event" && hasPrefix && (strings.EqualFold(prefix, "Dialogue") || strings.EqualFold(prefix, "Comment"))) || (r.Kind == "format_declaration" && hasPrefix && strings.EqualFold(prefix, "Format"))) {
		return nil
	}
	if r.Kind == "attachment_header" {
		_, name, ok := strings.Cut(raw, ":")
		if !ok || !nativeSafeName(strings.TrimSpace(name)) {
			return scriptedError("unsafe_source_metadata", r.SourceOrder)
		}
		return nil
	}
	if metadataIdentity.MatchString(raw) {
		return scriptedError("unsafe_source_metadata", r.SourceOrder)
	}
	if trimmed == "" || strings.HasPrefix(trimmed, ";") {
		return nil
	}
	if s == "script info" || s == "aegisub project garbage" {
		name, value, ok := strings.Cut(raw, ":")
		if !ok {
			return scriptedError("unsafe_source_metadata", r.SourceOrder)
		}
		return scriptedMetadata(name, value, r.SourceOrder)
	}
	if !content || hasPrefix {
		return scriptedError("unsafe_source_metadata", r.SourceOrder)
	}
	return nil
}
func scriptedSourcePrivacy(source SourceEnvelope, format string) error {
	for _, a := range source.Assets {
		if a.ID != source.PrimaryAssetID {
			continue
		}
		if len(a.DataBase64) > ((64<<20)+2)/3*4 {
			return fmt.Errorf("complexity_limit: native source")
		}
		raw, err := base64.StdEncoding.Strict().DecodeString(a.DataBase64)
		if err != nil {
			return fmt.Errorf("invalid scripted source base64")
		}
		if len(raw) == 0 || !utf8.Valid(raw) || strings.ContainsRune(string(raw), '\x00') {
			return fmt.Errorf("invalid scripted source text")
		}
		text := strings.TrimPrefix(string(raw), "\ufeff")
		text = strings.ReplaceAll(text, "\r\n", "\n")
		if strings.Contains(text, "\r") {
			return fmt.Errorf("invalid scripted source line endings")
		}
		section := ""
		dialect := "v4.00+"
		styleSection := "V4+ Styles"
		if format == "ssa" {
			dialect = "v4.00"
			styleSection = "V4 Styles"
		}
		seenType, seenStyles := false, false
		var declaration []string
		for order := 0; text != ""; order++ {
			if order >= MaxDocumentItems {
				return scriptedError("complexity_limit", order)
			}
			line, remaining, _ := strings.Cut(text, "\n")
			text = remaining
			if !scriptedPhysical(line) {
				return scriptedError("complexity_limit", order)
			}
			trim := strings.TrimSpace(line)
			if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
				section = trim[1 : len(trim)-1]
				declaration = nil
				if (format == "ass" && strings.EqualFold(section, "V4 Styles")) || (format == "ssa" && strings.EqualFold(section, "V4+ Styles")) {
					return scriptedError("mismatched_native_dialect", order)
				}
				if strings.EqualFold(section, styleSection) {
					seenStyles = true
				}
				if metadataIdentity.MatchString(section) {
					return scriptedError("unsafe_source_metadata", order)
				}
				continue
			}
			kind := "unknown"
			prefix, value, hasValue := strings.Cut(trim, ":")
			if strings.EqualFold(section, "Script Info") && hasValue && strings.EqualFold(prefix, "ScriptType") {
				if strings.TrimSpace(value) != dialect {
					return scriptedError("mismatched_native_dialect", order)
				}
				seenType = true
			}
			if strings.EqualFold(section, "fonts") || strings.EqualFold(section, "graphics") {
				kind = "attachment_data"
				if strings.HasPrefix(strings.ToLower(trim), "fontname:") || strings.HasPrefix(strings.ToLower(trim), "filename:") {
					kind = "attachment_header"
				}
			} else if strings.EqualFold(section, "events") {
				kind = "event"
				prefix, _, _ := strings.Cut(trim, ":")
				if strings.EqualFold(prefix, "Format") {
					kind = "format_declaration"
				}
				if strings.EqualFold(prefix, "command") || strings.Contains(strings.ToLower(prefix), "template") {
					return scriptedError("unsafe_active_content", order)
				}
			} else if strings.EqualFold(section, "v4 styles") || strings.EqualFold(section, "v4+ styles") {
				kind = "style"
				prefix, _, _ := strings.Cut(trim, ":")
				if strings.EqualFold(prefix, "Format") {
					kind = "format_declaration"
				}
			}
			if kind == "format_declaration" {
				declaration = strings.Split(value, ",")
				if len(declaration) > MaxScriptedDeclarationFields {
					return scriptedError("complexity_limit", order)
				}
			}
			if (kind == "event" || kind == "style") && hasValue && len(declaration) > 0 && (strings.EqualFold(prefix, "Dialogue") || strings.EqualFold(prefix, "Comment") || strings.EqualFold(prefix, "Style")) {
				values := strings.SplitN(value, ",", len(declaration))
				if len(values) == len(declaration) {
					for i, name := range declaration {
						if err = scriptedContentFieldPrivacy(name, values[i], order); err != nil {
							return err
						}
					}
				}
			}
			r := ScriptedRecord{SourceOrder: order, Kind: kind, RawLine: &line}
			if err = scriptedRecordPrivacy(r, section); err != nil {
				return err
			}
		}
		if !seenType || !seenStyles {
			return fmt.Errorf("mismatched_native_dialect: original script requires matching ScriptType and style section")
		}
	}
	return nil
}
func validateScriptedAttachments(n *ScriptedDocumentData, records map[string]ScriptedRecord, ids map[string]bool, diagnosticIndex map[string]bool) error {
	occupied := map[string]bool{}
	recordByOrder := map[int]ScriptedRecord{}
	sectionByID := map[string]string{}
	for _, section := range n.Sections {
		sectionByID[section.SectionID] = strings.ToLower(section.Name)
	}
	for _, r := range n.Records {
		recordByOrder[r.SourceOrder] = r
	}
	total := 0
	last := -1
	for _, a := range n.Attachments {
		h, ok := records[a.HeaderRecordID]
		start, found := records[a.DataStartRecordID]
		if !ok || !found || h.Kind != "attachment_header" || h.AttachmentID == nil || *h.AttachmentID != a.AttachmentID || !nativeID(a.AttachmentID) || ids[a.AttachmentID] || h.SourceOrder <= last || a.DataRecordCount <= 0 || start.SourceOrder != h.SourceOrder+1 || !nativeSafeName(a.Name) || !slices.Contains([]string{"font", "graphic"}, a.AttachmentType) {
			return scriptedError("invalid_attachment_ownership", h.SourceOrder)
		}
		ids[a.AttachmentID] = true
		last = h.SourceOrder
		sectionName := sectionByID[h.SectionID]
		if (a.AttachmentType == "font" && sectionName != "fonts") || (a.AttachmentType == "graphic" && sectionName != "graphics") {
			return scriptedError("invalid_attachment_section", h.SourceOrder)
		}
		if h.RawLine != nil {
			prefix, _, found := strings.Cut(*h.RawLine, ":")
			want := "fontname"
			if a.AttachmentType == "graphic" {
				want = "filename"
			}
			if !found || !strings.EqualFold(strings.TrimSpace(prefix), want) {
				return scriptedError("invalid_attachment_header", h.SourceOrder)
			}
		}
		var encodedBuilder strings.Builder
		count := 0
		malformed := false
		if a.DataRecordCount > len(n.Records) {
			return scriptedError("invalid_attachment_range", h.SourceOrder)
		}
		for offset := 0; offset < a.DataRecordCount; offset++ {
			r, exists := recordByOrder[start.SourceOrder+offset]
			if !exists {
				return scriptedError("invalid_attachment_range", h.SourceOrder)
			}
			if r.SectionID != h.SectionID || r.Kind != "attachment_data" || r.RawLine == nil || occupied[r.RecordID] {
				return scriptedError("invalid_attachment_range", r.SourceOrder)
			}
			occupied[r.RecordID] = true
			count++
			line := *r.RawLine
			if len(line) > 80 {
				return scriptedError("complexity_limit", r.SourceOrder)
			}
			for _, b := range []byte(line) {
				if b < 33 || b > 96 {
					malformed = true
				}
			}
			encodedBuilder.WriteString(line)
		}
		if count != a.DataRecordCount {
			return scriptedError("invalid_attachment_range", h.SourceOrder)
		}
		encoded := encodedBuilder.String()
		length := len(encoded) / 4 * 3
		rem := len(encoded) % 4
		if rem == 1 {
			malformed = true
		}
		if rem == 2 {
			length++
		}
		if rem == 3 {
			length += 2
		}
		if rem == 2 && len(encoded) > 0 && ((encoded[len(encoded)-1]-33)&15) != 0 {
			malformed = true
		}
		if rem == 3 && len(encoded) > 0 && ((encoded[len(encoded)-1]-33)&3) != 0 {
			malformed = true
		}
		if length > 16<<20 {
			return scriptedError("complexity_limit", h.SourceOrder)
		}
		total += length
		if total > 32<<20 {
			return scriptedError("complexity_limit", h.SourceOrder)
		}
		if malformed && !diagnosticIndex["malformed_attachment:"+strconv.Itoa(h.SourceOrder)] {
			return scriptedError("missing_native_diagnostic", h.SourceOrder)
		}
	}
	if len(n.Attachments) != countScriptedKind(n.Records, "attachment_header") || len(occupied) != countScriptedKind(n.Records, "attachment_data") {
		return fmt.Errorf("orphaned attachment record")
	}
	return nil
}
