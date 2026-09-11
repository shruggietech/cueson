// Package convert owns cross-format projection and explicit loss accounting.
package convert

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/model"
)

const (
	// MaxLosses bounds report amplification before cloning, sorting, or indexing.
	MaxLosses              = 8_192
	maxLossMessageBytes    = 512
	maxLossPathBytes       = 1024
	maxContextAttributes   = 16
	maxAttributeNameBytes  = 64
	maxAttributeValueBytes = 128
)

// Severity classifies a recoverable conversion loss independently from strict
// policy. Fatal conversion conditions are returned as errors instead.
type Severity string

const SeverityWarning Severity = "warning"

// Kind identifies how target output changes one source semantic.
type Kind string

const (
	KindOmitted   Kind = "omitted"
	KindDegraded  Kind = "degraded"
	KindAmbiguous Kind = "ambiguous"
)

// Stable loss codes shared by conversion and representability analysis.
const (
	LossCodeMetadataOmitted                = "conversion_metadata_omitted"
	LossCodePayloadLineDegraded            = "conversion_payload_line_degraded"
	LossCodeOCRObservationOmitted          = "conversion_ocr_observation_omitted"
	LossCodeSpeakerObservationOmitted      = "conversion_speaker_observation_omitted"
	LossCodeTokenTimingOmitted             = "conversion_token_timing_omitted"
	LossCodePlacementOmitted               = "conversion_placement_omitted"
	LossCodeNULDegraded                    = "conversion_nul_degraded"
	LossCodeSubRipCoordinatesOmitted       = "conversion_subrip_coordinates_omitted"
	LossCodeSubRipFontDegraded             = "conversion_subrip_font_degraded"
	LossCodeSubRipUnrecognizedBlockOmitted = "conversion_subrip_unrecognized_block_omitted"
	LossCodeSubRipPayloadAmbiguous         = "conversion_subrip_payload_ambiguous"
	LossCodeWebVTTDescriptionOmitted       = "conversion_webvtt_description_omitted"
	LossCodeWebVTTMetadataOmitted          = "conversion_webvtt_metadata_omitted"
	LossCodeWebVTTBlockOmitted             = "conversion_webvtt_block_omitted"
	LossCodeWebVTTIdentifierOmitted        = "conversion_webvtt_identifier_omitted"
	LossCodeWebVTTSettingOmitted           = "conversion_webvtt_setting_omitted"
	LossCodeWebVTTSettingOccurrenceOmitted = "conversion_webvtt_setting_occurrence_omitted"
	LossCodeWebVTTVoiceDegraded            = "conversion_webvtt_voice_degraded"
	LossCodeWebVTTMarkupDegraded           = "conversion_webvtt_markup_degraded"
	LossCodeWebVTTInlineTimingOmitted      = "conversion_webvtt_inline_timing_omitted"
	LossCodeWebVTTEntityAmbiguous          = "conversion_webvtt_entity_ambiguous"
)

type lossCodeSpec struct {
	kind Kind
	rank int
}

var knownLossCodes = map[string]lossCodeSpec{
	LossCodeMetadataOmitted:                {kind: KindOmitted, rank: 10},
	LossCodeWebVTTDescriptionOmitted:       {kind: KindOmitted, rank: 20},
	LossCodeWebVTTMetadataOmitted:          {kind: KindOmitted, rank: 30},
	LossCodePayloadLineDegraded:            {kind: KindDegraded, rank: 100},
	LossCodeSubRipPayloadAmbiguous:         {kind: KindAmbiguous, rank: 110},
	LossCodeSubRipCoordinatesOmitted:       {kind: KindOmitted, rank: 120},
	LossCodeSubRipFontDegraded:             {kind: KindDegraded, rank: 130},
	LossCodeSpeakerObservationOmitted:      {kind: KindOmitted, rank: 140},
	LossCodeTokenTimingOmitted:             {kind: KindOmitted, rank: 150},
	LossCodeOCRObservationOmitted:          {kind: KindOmitted, rank: 160},
	LossCodePlacementOmitted:               {kind: KindOmitted, rank: 170},
	LossCodeNULDegraded:                    {kind: KindDegraded, rank: 180},
	LossCodeSubRipUnrecognizedBlockOmitted: {kind: KindOmitted, rank: 190},
	LossCodeWebVTTBlockOmitted:             {kind: KindOmitted, rank: 200},
	LossCodeWebVTTIdentifierOmitted:        {kind: KindOmitted, rank: 210},
	LossCodeWebVTTSettingOmitted:           {kind: KindOmitted, rank: 220},
	LossCodeWebVTTSettingOccurrenceOmitted: {kind: KindOmitted, rank: 230},
	LossCodeWebVTTVoiceDegraded:            {kind: KindDegraded, rank: 240},
	LossCodeWebVTTMarkupDegraded:           {kind: KindDegraded, rank: 250},
	LossCodeWebVTTInlineTimingOmitted:      {kind: KindOmitted, rank: 260},
	LossCodeWebVTTEntityAmbiguous:          {kind: KindAmbiguous, rank: 270},
}

// Attribute is one bounded deterministic context value. Context supplements a
// stable code and JSON Pointer; it never carries source excerpts or paths.
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Loss describes one atomic source semantic changed by a requested target.
type Loss struct {
	Code         string      `json:"code"`
	Severity     Severity    `json:"severity"`
	Kind         Kind        `json:"kind"`
	Message      string      `json:"message"`
	SourceFormat string      `json:"source_format"`
	TargetFormat string      `json:"target_format"`
	Path         string      `json:"path"`
	SourceOrder  *int        `json:"source_order"`
	CueID        *string     `json:"cue_id"`
	Context      []Attribute `json:"context"`
}

// Report contains the complete canonically ordered loss list.
type Report struct {
	Losses []Loss `json:"losses"`
}

// Result is the side-effect-free output of one conversion request.
type Result struct {
	Bytes       []byte
	LossReport  Report
	Diagnostics []model.Diagnostic
}

// Options controls conversion loss policy.
type Options struct {
	Strict bool
}

// NewReport clones, canonicalizes, and validates a loss list against its
// immutable source document.
func NewReport(document model.Document, losses []Loss) (Report, error) {
	if len(losses) > MaxLosses {
		return Report{}, fmt.Errorf("losses contains %d items, maximum is %d", len(losses), MaxLosses)
	}
	prepared := make([]Loss, len(losses))
	for index := range losses {
		prepared[index] = cloneLoss(losses[index])
		if prepared[index].Context == nil {
			prepared[index].Context = []Attribute{}
		}
		slices.SortFunc(prepared[index].Context, compareAttribute)
	}
	slices.SortFunc(prepared, compareLoss)
	report := Report{Losses: prepared}
	if err := report.Validate(document); err != nil {
		return Report{}, err
	}
	return report, nil
}

// Validate checks loss vocabulary, bounds, source associations, uniqueness,
// and canonical order without changing the report or source document.
func (report Report) Validate(document model.Document) error {
	if len(report.Losses) > MaxLosses {
		return fmt.Errorf("losses contains %d items, maximum is %d", len(report.Losses), MaxLosses)
	}
	cueByID := make(map[string]model.Cue, len(document.Cues))
	orders := make(map[int]struct{}, len(document.Cues)+len(document.Diagnostics))
	for _, cue := range document.Cues {
		cueByID[cue.ID] = cue
		orders[cue.SourceOrder] = struct{}{}
	}
	if document.FormatData.WebVTT != nil {
		for _, block := range document.FormatData.WebVTT.Blocks {
			orders[block.SourceOrder] = struct{}{}
		}
	}
	for _, diagnostic := range document.Diagnostics {
		if diagnostic.SourceOrder != nil {
			orders[*diagnostic.SourceOrder] = struct{}{}
		}
	}

	identities := make(map[string]struct{}, len(report.Losses))
	for index := range report.Losses {
		loss := &report.Losses[index]
		if err := validateLoss(*loss, document.Format, cueByID, orders); err != nil {
			return fmt.Errorf("losses[%d]: %w", index, err)
		}
		identity := lossIdentity(*loss)
		if _, duplicate := identities[identity]; duplicate {
			return fmt.Errorf("losses[%d]: duplicate loss identity", index)
		}
		identities[identity] = struct{}{}
		if index > 0 && compareLoss(report.Losses[index-1], *loss) > 0 {
			return fmt.Errorf("losses[%d]: loss report is not in canonical order", index)
		}
	}
	return nil
}

// HasLosses reports whether strict conversion must reject the request.
func (report Report) HasLosses() bool {
	return len(report.Losses) > 0
}

// First returns the first canonical loss, or nil for an empty report.
func (report Report) First() *Loss {
	if len(report.Losses) == 0 {
		return nil
	}
	return &report.Losses[0]
}

// IsKnownLossCode reports whether code belongs to the stable S016 vocabulary.
func IsKnownLossCode(code string) bool {
	_, exists := knownLossCodes[code]
	return exists
}

// StrictLossError retains the complete report when strict policy rejects loss.
type StrictLossError struct {
	Report Report
}

func (err *StrictLossError) Error() string {
	if err == nil || len(err.Report.Losses) == 0 {
		return "strict conversion blocked by an invalid empty loss report"
	}
	first := err.Report.Losses[0]
	return fmt.Sprintf("strict conversion blocked %d loss(es); first: %s at %s: %s", len(err.Report.Losses), first.Code, first.Path, first.Message)
}

// SameFormatError directs callers to native model rendering.
type SameFormatError struct {
	Format string
}

func (err *SameFormatError) Error() string {
	if err == nil {
		return "same-format conversion is unsupported; use render"
	}
	return fmt.Sprintf("source and target format %q are the same; use render", err.Format)
}

// UnsupportedPairError identifies a source-target pair outside S016.
type UnsupportedPairError struct {
	SourceFormat string
	TargetFormat string
}

func (err *UnsupportedPairError) Error() string {
	if err == nil {
		return "unsupported conversion pair"
	}
	return fmt.Sprintf("unsupported conversion pair %q to %q", err.SourceFormat, err.TargetFormat)
}

// ProjectionError identifies a fatal field-level target projection boundary.
type ProjectionError struct {
	SourceFormat string
	TargetFormat string
	Path         string
	Err          error
}

func (err *ProjectionError) Error() string {
	if err == nil {
		return "conversion projection failed"
	}
	if err.Path == "" {
		return fmt.Sprintf("project %q to %q: %v", err.SourceFormat, err.TargetFormat, err.Err)
	}
	return fmt.Sprintf("project %q to %q at %s: %v", err.SourceFormat, err.TargetFormat, err.Path, err.Err)
}

func (err *ProjectionError) Unwrap() error {
	if err == nil {
		return nil
	}
	return err.Err
}

func validateLoss(loss Loss, documentFormat string, cues map[string]model.Cue, orders map[int]struct{}) error {
	spec, known := knownLossCodes[loss.Code]
	if !known {
		return fmt.Errorf("unknown loss code %q", loss.Code)
	}
	if loss.Severity != SeverityWarning {
		return fmt.Errorf("severity is %q, want %q", loss.Severity, SeverityWarning)
	}
	if loss.Kind != spec.kind {
		return fmt.Errorf("kind for %s is %q, want %q", loss.Code, loss.Kind, spec.kind)
	}
	if loss.SourceFormat != documentFormat {
		return fmt.Errorf("source format %q does not match document format %q", loss.SourceFormat, documentFormat)
	}
	if !supportedFormat(loss.SourceFormat) || !supportedFormat(loss.TargetFormat) || loss.SourceFormat == loss.TargetFormat {
		return fmt.Errorf("format pair %q to %q is not supported", loss.SourceFormat, loss.TargetFormat)
	}
	if err := validateSafeText("message", loss.Message, maxLossMessageBytes); err != nil {
		return err
	}
	if err := validateJSONPointer(loss.Path); err != nil {
		return err
	}
	if loss.CueID != nil {
		cue, exists := cues[*loss.CueID]
		if !exists {
			return fmt.Errorf("cue_id %q does not resolve", *loss.CueID)
		}
		if loss.SourceOrder == nil || *loss.SourceOrder != cue.SourceOrder {
			return fmt.Errorf("source_order does not match cue_id %q", *loss.CueID)
		}
	}
	if loss.SourceOrder != nil {
		if _, exists := orders[*loss.SourceOrder]; !exists {
			return fmt.Errorf("source_order %d does not resolve", *loss.SourceOrder)
		}
	} else if loss.CueID != nil {
		return fmt.Errorf("cue_id requires source_order")
	}
	if len(loss.Context) > maxContextAttributes {
		return fmt.Errorf("context contains %d attributes, maximum is %d", len(loss.Context), maxContextAttributes)
	}
	for index := range loss.Context {
		attribute := loss.Context[index]
		if !validAttributeName(attribute.Name) {
			return fmt.Errorf("context[%d].name %q is invalid", index, attribute.Name)
		}
		if err := validateSafeText(fmt.Sprintf("context[%d].value", index), attribute.Value, maxAttributeValueBytes); err != nil {
			return err
		}
		if index > 0 {
			comparison := compareAttribute(loss.Context[index-1], attribute)
			if comparison == 0 || loss.Context[index-1].Name == attribute.Name {
				return fmt.Errorf("context contains duplicate attribute name %q", attribute.Name)
			}
			if comparison > 0 {
				return fmt.Errorf("context is not in canonical order")
			}
		}
	}
	return nil
}

func supportedFormat(format string) bool {
	return format == "subrip" || format == "webvtt"
}

func validateSafeText(field, value string, maximum int) error {
	if value == "" {
		return fmt.Errorf("%s must not be empty", field)
	}
	if !utf8.ValidString(value) || len(value) > maximum {
		return fmt.Errorf("%s must be valid UTF-8 and at most %d bytes", field, maximum)
	}
	if containsControl(value) {
		return fmt.Errorf("%s contains a control character", field)
	}
	if strings.ContainsAny(value, `/\\`) || strings.Contains(value, "://") {
		return fmt.Errorf("%s contains path-like content", field)
	}
	return nil
}

func containsControl(value string) bool {
	for _, character := range value {
		if unicode.IsControl(character) {
			return true
		}
	}
	return false
}

func validateJSONPointer(pointer string) error {
	if pointer == "" || len(pointer) > maxLossPathBytes || !utf8.ValidString(pointer) || pointer[0] != '/' || containsControl(pointer) {
		return fmt.Errorf("path %q is not a bounded RFC 6901 JSON Pointer", pointer)
	}
	for index := 0; index < len(pointer); index++ {
		if pointer[index] != '~' {
			continue
		}
		if index+1 >= len(pointer) || pointer[index+1] != '0' && pointer[index+1] != '1' {
			return fmt.Errorf("path %q is not a bounded RFC 6901 JSON Pointer", pointer)
		}
		index++
	}
	root := pointerSegment(strings.Split(pointer[1:], "/")[0])
	switch root {
	case "$schema", "schema_version", "format", "format_support", "producer", "source", "metadata", "document", "cues", "format_data", "diagnostics", "stats":
		return nil
	default:
		return fmt.Errorf("path %q is not a Cue JSON Pointer", pointer)
	}
}

func pointerSegment(segment string) string {
	segment = strings.ReplaceAll(segment, "~1", "/")
	return strings.ReplaceAll(segment, "~0", "~")
}

func validAttributeName(value string) bool {
	if value == "" || len(value) > maxAttributeNameBytes || value[0] < 'a' || value[0] > 'z' {
		return false
	}
	for index := 1; index < len(value); index++ {
		character := value[index]
		if character != '_' && (character < 'a' || character > 'z') && (character < '0' || character > '9') {
			return false
		}
	}
	return true
}

func compareLoss(left, right Loss) int {
	leftDocument, rightDocument := left.SourceOrder == nil, right.SourceOrder == nil
	if leftDocument != rightDocument {
		if leftDocument {
			return -1
		}
		return 1
	}
	if left.SourceOrder != nil {
		if comparison := cmp.Compare(*left.SourceOrder, *right.SourceOrder); comparison != 0 {
			return comparison
		}
	}
	if comparison := cmp.Compare(knownLossCodes[left.Code].rank, knownLossCodes[right.Code].rank); comparison != 0 {
		return comparison
	}
	if comparison := cmp.Compare(lossRuleRank(left), lossRuleRank(right)); comparison != 0 {
		return comparison
	}
	if comparison := compareJSONPointer(left.Path, right.Path); comparison != 0 {
		return comparison
	}
	if comparison := compareAttributes(left.Context, right.Context); comparison != 0 {
		return comparison
	}
	if comparison := cmp.Compare(left.Code, right.Code); comparison != 0 {
		return comparison
	}
	if comparison := compareOptionalInt(left.SourceOrder, right.SourceOrder); comparison != 0 {
		return comparison
	}
	return compareOptionalString(left.CueID, right.CueID)
}

func lossRuleRank(loss Loss) int {
	if loss.Code != LossCodeWebVTTSettingOmitted {
		return 0
	}
	for _, attribute := range loss.Context {
		if attribute.Name != "setting" {
			continue
		}
		switch attribute.Value {
		case "region":
			return 1
		case "vertical":
			return 2
		case "line":
			return 3
		case "position":
			return 4
		case "size":
			return 5
		case "align":
			return 6
		}
	}
	return 7
}

func compareJSONPointer(left, right string) int {
	leftParts, rightParts := strings.Split(left, "/"), strings.Split(right, "/")
	for index := 0; index < min(len(leftParts), len(rightParts)); index++ {
		leftPart, rightPart := pointerSegment(leftParts[index]), pointerSegment(rightParts[index])
		leftNumber, leftNumeric := pointerIndex(leftPart)
		rightNumber, rightNumeric := pointerIndex(rightPart)
		if leftNumeric && rightNumeric {
			if comparison := cmp.Compare(leftNumber, rightNumber); comparison != 0 {
				return comparison
			}
			continue
		}
		if comparison := cmp.Compare(leftPart, rightPart); comparison != 0 {
			return comparison
		}
	}
	return cmp.Compare(len(leftParts), len(rightParts))
}

func pointerIndex(value string) (uint64, bool) {
	if value == "" || len(value) > 1 && value[0] == '0' {
		return 0, false
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	return parsed, err == nil
}

func compareAttribute(left, right Attribute) int {
	if comparison := cmp.Compare(left.Name, right.Name); comparison != 0 {
		return comparison
	}
	return cmp.Compare(left.Value, right.Value)
}

func compareAttributes(left, right []Attribute) int {
	for index := 0; index < min(len(left), len(right)); index++ {
		if comparison := compareAttribute(left[index], right[index]); comparison != 0 {
			return comparison
		}
	}
	return cmp.Compare(len(left), len(right))
}

func compareOptionalInt(left, right *int) int {
	if left == nil && right == nil {
		return 0
	}
	if left == nil {
		return -1
	}
	if right == nil {
		return 1
	}
	return cmp.Compare(*left, *right)
}

func compareOptionalString(left, right *string) int {
	if left == nil && right == nil {
		return 0
	}
	if left == nil {
		return -1
	}
	if right == nil {
		return 1
	}
	return cmp.Compare(*left, *right)
}

func lossIdentity(loss Loss) string {
	var identity strings.Builder
	identity.WriteString(loss.Code)
	identity.WriteByte(0)
	identity.WriteString(loss.Path)
	for _, attribute := range loss.Context {
		identity.WriteByte(0)
		identity.WriteString(attribute.Name)
		identity.WriteByte('=')
		identity.WriteString(attribute.Value)
	}
	return identity.String()
}

func cloneLoss(loss Loss) Loss {
	cloned := loss
	if loss.SourceOrder != nil {
		value := *loss.SourceOrder
		cloned.SourceOrder = &value
	}
	if loss.CueID != nil {
		value := *loss.CueID
		cloned.CueID = &value
	}
	cloned.Context = slices.Clone(loss.Context)
	return cloned
}
