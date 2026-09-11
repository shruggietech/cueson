package cli

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/version"
)

const inspectReportVersion = "1"

type inspectOptions struct {
	input       string
	format      string
	encoding    string
	json        bool
	formatSet   bool
	encodingSet bool
	jsonSet     bool
}

type inspectionReport struct {
	InspectReportVersion string                        `json:"inspect_report_version"`
	Input                inspectionInputSummary        `json:"input"`
	Format               string                        `json:"format"`
	Schema               inspectionSchemaSummary       `json:"schema"`
	Capabilities         inspectionCapabilitySummary   `json:"capabilities"`
	Integrity            inspectionIntegritySummary    `json:"integrity"`
	Assets               []inspectionAssetSummary      `json:"assets"`
	Document             inspectionDocumentSummary     `json:"document"`
	Cues                 []inspectionCueSummary        `json:"cues"`
	Blocks               inspectionBlockSummary        `json:"blocks"`
	Diagnostics          []inspectionDiagnosticSummary `json:"diagnostics"`
	Loss                 inspectionLossState           `json:"loss"`
}

type inspectionInputSummary struct {
	Kind            string  `json:"kind"`
	SelectionBasis  string  `json:"selection_basis"`
	ContentFormat   *string `json:"content_format"`
	ExtensionFormat *string `json:"extension_format"`
}

type inspectionSchemaSummary struct {
	ID         string `json:"id"`
	Version    string `json:"version"`
	Compatible bool   `json:"compatible"`
}

type inspectionCapabilitySummary struct {
	Declared  inspectionDeclaredCapabilities  `json:"declared"`
	Installed inspectionInstalledCapabilities `json:"installed"`
}

type inspectionDeclaredCapabilities struct {
	Status                       string `json:"status"`
	IngestSupported              bool   `json:"ingest_supported"`
	RenderSupported              bool   `json:"render_supported"`
	RestoreSupported             bool   `json:"restore_supported"`
	OCRRequiredForSemanticOutput bool   `json:"ocr_required_for_semantic_output"`
}

type inspectionInstalledCapabilities struct {
	Ingest   bool `json:"ingest"`
	Render   bool `json:"render"`
	Restore  bool `json:"restore"`
	Validate bool `json:"validate"`
	Inspect  bool `json:"inspect"`
}

type inspectionIntegritySummary struct {
	Status     string `json:"status"`
	AssetCount int    `json:"asset_count"`
	TotalBytes int64  `json:"total_bytes"`
}

type inspectionAssetSummary struct {
	Index     int                        `json:"index"`
	Primary   bool                       `json:"primary"`
	Role      string                     `json:"role"`
	FileName  string                     `json:"file_name"`
	MediaType *string                    `json:"media_type"`
	Bytes     int64                      `json:"bytes"`
	Encoding  *inspectionEncodingSummary `json:"encoding"`
}

type inspectionEncodingSummary struct {
	BOM              *string  `json:"bom"`
	LineEndings      string   `json:"line_endings"`
	DetectedEncoding *string  `json:"detected_encoding"`
	Confidence       *float64 `json:"confidence"`
}

type inspectionDocumentSummary struct {
	CueCount               int    `json:"cue_count"`
	NonCueBlockCount       int    `json:"non_cue_block_count"`
	BodyItemCount          int    `json:"body_item_count"`
	MediaStartMilliseconds *int64 `json:"media_start_milliseconds"`
	MediaEndMilliseconds   *int64 `json:"media_end_milliseconds"`
	MediaSpanMilliseconds  *int64 `json:"media_span_milliseconds"`
	HasWordLevelTiming     bool   `json:"has_word_level_timing"`
	DiagnosticCount        int    `json:"diagnostic_count"`
	InfoCount              int    `json:"info_count"`
	WarningCount           int    `json:"warning_count"`
	ErrorCount             int    `json:"error_count"`
}

type inspectionCueSummary struct {
	Ordinal                      int   `json:"ordinal"`
	SourceOrder                  int   `json:"source_order"`
	StartMilliseconds            int64 `json:"start_milliseconds"`
	EndMilliseconds              int64 `json:"end_milliseconds"`
	DurationMilliseconds         int64 `json:"duration_milliseconds"`
	PayloadLineCount             int   `json:"payload_line_count"`
	SpeakerCount                 int   `json:"speaker_count"`
	TokenCount                   int   `json:"token_count"`
	OCRObservationCount          int   `json:"ocr_observation_count"`
	HasSourceIdentifier          bool  `json:"has_source_identifier"`
	HasPlacement                 bool  `json:"has_placement"`
	NativeSettingCount           int   `json:"native_setting_count"`
	NativeSettingOccurrenceCount int   `json:"native_setting_occurrence_count"`
	InvalidNativeSettingCount    int   `json:"invalid_native_setting_count"`
}

type inspectionBlockSummary struct {
	Total        int                           `json:"total"`
	NoteCount    int                           `json:"note_count"`
	StyleCount   int                           `json:"style_count"`
	RegionCount  int                           `json:"region_count"`
	UnknownCount int                           `json:"unknown_count"`
	Entries      []inspectionBlockEntrySummary `json:"entries"`
}

type inspectionBlockEntrySummary struct {
	SourceOrder                  int    `json:"source_order"`
	Type                         string `json:"type"`
	LineCount                    int    `json:"line_count"`
	HasParsedRegion              bool   `json:"has_parsed_region"`
	NativeSettingCount           int    `json:"native_setting_count"`
	NativeSettingOccurrenceCount int    `json:"native_setting_occurrence_count"`
	InvalidNativeSettingCount    int    `json:"invalid_native_setting_count"`
}

type inspectionDiagnosticSummary struct {
	Severity    string `json:"severity"`
	Code        string `json:"code"`
	SourceOrder *int   `json:"source_order"`
	CueOrdinal  *int   `json:"cue_ordinal"`
}

type inspectionLossState struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func setInspectValue(options *inspectOptions, option, value string) error {
	if options == nil {
		return fmt.Errorf("inspect options are unavailable")
	}
	switch option {
	case "--format":
		if options.formatSet {
			return fmt.Errorf("--format may be specified only once")
		}
		options.format, options.formatSet = value, true
	case "--encoding":
		if options.encodingSet {
			return fmt.Errorf("--encoding may be specified only once")
		}
		options.encoding, options.encodingSet = value, true
	default:
		return fmt.Errorf("unknown inspect option %q", option)
	}
	return nil
}

func setInspectJSON(options *inspectOptions) error {
	if options == nil {
		return fmt.Errorf("inspect options are unavailable")
	}
	if options.jsonSet {
		return fmt.Errorf("--json may be specified only once")
	}
	options.json, options.jsonSet = true, true
	return nil
}

func setInspectOperand(options *inspectOptions, value string) error {
	if options == nil {
		return fmt.Errorf("inspect options are unavailable")
	}
	if options.input != "" {
		return fmt.Errorf("inspect accepts no additional arguments")
	}
	options.input = value
	return nil
}

func finalizeInspectOptions(options *inspectOptions) error {
	if options == nil {
		return fmt.Errorf("inspect options are unavailable")
	}
	if options.input == "" {
		return fmt.Errorf("inspect requires one INPUT path")
	}
	if options.formatSet && options.format == "" {
		return fmt.Errorf("--format requires a value")
	}
	if options.encodingSet && options.encoding == "" {
		return fmt.Errorf("--encoding requires a value")
	}
	format, known := normalizeInputFormat(options.format)
	if !known {
		return fmt.Errorf("--format must be auto, cueson, srt, or vtt")
	}
	options.format = format
	if options.encoding != "" {
		if _, known := codec.NormalizeEncoding(options.encoding); !known {
			return fmt.Errorf("encoding %q is not supported", options.encoding)
		}
	}
	if options.format == "cueson" && options.encoding != "" {
		return fmt.Errorf("--encoding cannot be used with Cue JSON input")
	}
	if options.format == string(codec.FormatWebVTT) {
		if err := validateWebVTTInputEncoding(options.encoding); err != nil {
			return err
		}
	}
	return nil
}

func runInspect(ctx context.Context, options inspectOptions, stdout, stderr io.Writer, diagnostics diagnosticWriter, usage helpTarget) int {
	if err := finalizeInspectOptions(&options); err != nil {
		diagnostics.write(diagnosticError, err.Error())
		writeUsage(stderr, usage)
		return ExitInvocation
	}
	captured, err := source.CaptureContext(ctx, options.input, source.CaptureOptions{})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, "inspect: input does not exist")
			writeUsage(stderr, usage)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, "inspect: input acquisition failed")
		return ExitRuntimeFailure
	}

	input, err := loadValidatedInput(ctx, captured, inputOptions{format: options.format, encoding: options.encoding})
	if err != nil {
		var constraint *inputConstraintError
		if errors.As(err, &constraint) {
			diagnostics.write(diagnosticError, constraint.Error())
			writeUsage(stderr, usage)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, "inspect: input validation failed")
		return ExitRuntimeFailure
	}
	report, err := buildInspectionReport(input)
	if err != nil {
		diagnostics.write(diagnosticError, "inspect: report construction failed")
		return ExitRuntimeFailure
	}
	for _, observation := range input.diagnostics {
		if observation.Severity == "warning" {
			// Report diagnostics deliberately exclude free-form messages. Use the
			// same safe stable code on stderr so Cue JSON cannot inject content.
			diagnostics.write(diagnosticWarning, observation.Code)
		}
	}

	var payload []byte
	if options.json {
		payload, err = marshalInspectionReport(report)
	} else {
		payload, err = renderHumanInspection(report)
	}
	if err != nil {
		diagnostics.write(diagnosticError, "inspect: report rendering failed")
		return ExitRuntimeFailure
	}
	return writeInspectionStdout(stdout, diagnostics, payload)
}

func writeInspectionStdout(stdout io.Writer, diagnostics diagnosticWriter, payload []byte) int {
	written, err := stdout.Write(payload)
	if err == nil && written != len(payload) {
		err = io.ErrShortWrite
	}
	if err != nil {
		diagnostics.write(diagnosticError, "inspect: stdout write failed")
		return ExitRuntimeFailure
	}
	return ExitSuccess
}

func buildInspectionReport(input validatedInput) (inspectionReport, error) {
	document := input.document
	if err := schema.CheckLockstep(version.String()); err != nil {
		return inspectionReport{}, fmt.Errorf("schema version lockstep: %w", err)
	}
	registry, err := workflowRegistry("")
	if err != nil {
		return inspectionReport{}, fmt.Errorf("load installed capabilities: %w", err)
	}
	registration, installed := registry.Lookup(document.Format)
	if !installed {
		return inspectionReport{}, fmt.Errorf("format %q has no installed registry entry", document.Format)
	}

	report := inspectionReport{
		InspectReportVersion: inspectReportVersion,
		Input: inspectionInputSummary{
			Kind:            string(input.inputKind),
			SelectionBasis:  string(input.selectionBasis),
			ContentFormat:   inspectionFormatPointer(input.contentFormat),
			ExtensionFormat: inspectionFormatPointer(input.extensionFormat),
		},
		Format: document.Format,
		Schema: inspectionSchemaSummary{ID: schema.ID(), Version: schema.Version(), Compatible: true},
		Capabilities: inspectionCapabilitySummary{
			Declared: inspectionDeclaredCapabilities{
				Status:                       document.FormatSupport.Status,
				IngestSupported:              document.FormatSupport.IngestSupported,
				RenderSupported:              document.FormatSupport.RenderSupported,
				RestoreSupported:             document.FormatSupport.RestoreSupported,
				OCRRequiredForSemanticOutput: document.FormatSupport.OCRRequiredForSemanticOutput,
			},
			Installed: inspectionInstalledCapabilities{
				Ingest:   registration.Decode != nil,
				Render:   registration.Render != nil,
				Restore:  true,
				Validate: true,
				Inspect:  true,
			},
		},
		Assets:      make([]inspectionAssetSummary, len(document.Source.Assets)),
		Cues:        make([]inspectionCueSummary, len(document.Cues)),
		Diagnostics: make([]inspectionDiagnosticSummary, len(document.Diagnostics)),
		Loss:        inspectionLossState{Status: "not_evaluated", Reason: "target_format_required"},
	}

	var totalBytes int64
	for index := range document.Source.Assets {
		asset := document.Source.Assets[index]
		if asset.Size.Bytes < 0 || totalBytes > math.MaxInt64-asset.Size.Bytes {
			return inspectionReport{}, fmt.Errorf("integrity total_bytes overflows int64")
		}
		totalBytes += asset.Size.Bytes
		report.Assets[index] = inspectionAssetSummary{
			Index:     index,
			Primary:   asset.ID == document.Source.PrimaryAssetID,
			Role:      asset.Role,
			FileName:  asset.FileName,
			MediaType: cloneStringPointer(asset.MediaType),
			Bytes:     asset.Size.Bytes,
			Encoding:  projectInspectionEncoding(asset.Encoding),
		}
	}
	report.Integrity = inspectionIntegritySummary{Status: "verified", AssetCount: len(report.Assets), TotalBytes: totalBytes}

	for index := range document.Cues {
		cue := document.Cues[index]
		settingCount, occurrenceCount, invalidCount := cueSettingCounts(cue)
		report.Cues[index] = inspectionCueSummary{
			Ordinal:                      cue.Ordinal,
			SourceOrder:                  cue.SourceOrder,
			StartMilliseconds:            cue.Timing.StartMilliseconds,
			EndMilliseconds:              cue.Timing.EndMilliseconds,
			DurationMilliseconds:         cue.Timing.DurationMilliseconds,
			PayloadLineCount:             len(cue.Payload.Lines),
			SpeakerCount:                 len(cue.Speakers),
			TokenCount:                   len(cue.Tokens),
			OCRObservationCount:          len(cue.OCRObservations),
			HasSourceIdentifier:          cue.SourceIdentifier != nil,
			HasPlacement:                 cue.Placement != nil,
			NativeSettingCount:           settingCount,
			NativeSettingOccurrenceCount: occurrenceCount,
			InvalidNativeSettingCount:    invalidCount,
		}
	}
	slices.SortFunc(report.Cues, func(left, right inspectionCueSummary) int {
		return cmp.Compare(left.Ordinal, right.Ordinal)
	})

	report.Blocks = projectInspectionBlocks(document)
	if len(report.Cues) > math.MaxInt-report.Blocks.Total {
		return inspectionReport{}, fmt.Errorf("document body_item_count overflows int")
	}
	infoCount := 0
	for _, diagnostic := range document.Diagnostics {
		if diagnostic.Severity == "info" {
			infoCount++
		}
	}
	report.Document = inspectionDocumentSummary{
		CueCount:               len(report.Cues),
		NonCueBlockCount:       report.Blocks.Total,
		BodyItemCount:          len(report.Cues) + report.Blocks.Total,
		MediaStartMilliseconds: cloneInt64Pointer(document.Document.MediaStartMilliseconds),
		MediaEndMilliseconds:   cloneInt64Pointer(document.Document.MediaEndMilliseconds),
		MediaSpanMilliseconds:  cloneInt64Pointer(document.Document.MediaSpanMilliseconds),
		HasWordLevelTiming:     document.Document.HasWordLevelTiming,
		DiagnosticCount:        len(document.Diagnostics),
		InfoCount:              infoCount,
		WarningCount:           document.Stats.WarningCount,
		ErrorCount:             document.Stats.ErrorCount,
	}

	cueOrdinals := make(map[string]int, len(document.Cues))
	for _, cue := range document.Cues {
		cueOrdinals[cue.ID] = cue.Ordinal
	}
	for index := range document.Diagnostics {
		diagnostic := document.Diagnostics[index]
		var ordinal *int
		if diagnostic.CueID != nil {
			value, exists := cueOrdinals[*diagnostic.CueID]
			if !exists {
				return inspectionReport{}, fmt.Errorf("diagnostic %d cue reference does not resolve", index)
			}
			ordinal = &value
		}
		report.Diagnostics[index] = inspectionDiagnosticSummary{
			Severity:    diagnostic.Severity,
			Code:        diagnostic.Code,
			SourceOrder: cloneIntPointer(diagnostic.SourceOrder),
			CueOrdinal:  ordinal,
		}
	}
	return report, nil
}

func projectInspectionBlocks(document model.Document) inspectionBlockSummary {
	summary := inspectionBlockSummary{Entries: []inspectionBlockEntrySummary{}}
	if document.FormatData.WebVTT == nil {
		return summary
	}
	summary.Entries = make([]inspectionBlockEntrySummary, len(document.FormatData.WebVTT.Blocks))
	for index := range document.FormatData.WebVTT.Blocks {
		block := document.FormatData.WebVTT.Blocks[index]
		settingCount, occurrenceCount, invalidCount := blockSettingCounts(block)
		summary.Entries[index] = inspectionBlockEntrySummary{
			SourceOrder:                  block.SourceOrder,
			Type:                         block.Type,
			LineCount:                    len(block.RawLines),
			HasParsedRegion:              block.Region != nil,
			NativeSettingCount:           settingCount,
			NativeSettingOccurrenceCount: occurrenceCount,
			InvalidNativeSettingCount:    invalidCount,
		}
		switch block.Type {
		case "note":
			summary.NoteCount++
		case "style":
			summary.StyleCount++
		case "region":
			summary.RegionCount++
		case "unrecognized":
			summary.UnknownCount++
		}
	}
	slices.SortFunc(summary.Entries, func(left, right inspectionBlockEntrySummary) int {
		if left.SourceOrder != right.SourceOrder {
			return cmp.Compare(left.SourceOrder, right.SourceOrder)
		}
		return strings.Compare(left.Type, right.Type)
	})
	summary.Total = len(summary.Entries)
	return summary
}

func cueSettingCounts(cue model.Cue) (int, int, int) {
	if cue.FormatData.WebVTT == nil {
		return 0, 0, 0
	}
	return len(cue.FormatData.WebVTT.Settings), len(cue.FormatData.WebVTT.SettingOccurrences), countInvalidSettings(cue.FormatData.WebVTT.SettingOccurrences)
}

func blockSettingCounts(block model.WebVTTBlock) (int, int, int) {
	if block.Region == nil {
		return 0, 0, 0
	}
	return len(block.Region.Settings), len(block.Region.SettingOccurrences), countInvalidSettings(block.Region.SettingOccurrences)
}

func countInvalidSettings(occurrences []model.WebVTTSettingOccurrence) int {
	count := 0
	for _, occurrence := range occurrences {
		if !occurrence.Valid {
			count++
		}
	}
	return count
}

func projectInspectionEncoding(observation *model.EncodingObservation) *inspectionEncodingSummary {
	if observation == nil {
		return nil
	}
	return &inspectionEncodingSummary{
		BOM:              cloneStringPointer(observation.BOM),
		LineEndings:      observation.LineEndings,
		DetectedEncoding: cloneStringPointer(observation.DetectedEncoding),
		Confidence:       cloneFloat64Pointer(observation.Confidence),
	}
}

func inspectionFormatPointer(format *codec.Format) *string {
	if format == nil {
		return nil
	}
	value := string(*format)
	return &value
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func cloneFloat64Pointer(value *float64) *float64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func marshalInspectionReport(report inspectionReport) ([]byte, error) {
	payload, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("marshal inspection report: %w", err)
	}
	return append(payload, '\n'), nil
}

func renderHumanInspection(report inspectionReport) ([]byte, error) {
	var output strings.Builder
	fmt.Fprintf(&output, "Inspection report %s\n", report.InspectReportVersion)
	fmt.Fprintf(&output, "Input: %s (%s), content_format=%s, extension_format=%s\n", report.Input.Kind, report.Input.SelectionBasis, nullableString(report.Input.ContentFormat), nullableString(report.Input.ExtensionFormat))
	fmt.Fprintf(&output, "Format: %s\n", report.Format)
	fmt.Fprintf(&output, "Schema: %s (version=%s compatible=%t)\n", report.Schema.ID, report.Schema.Version, report.Schema.Compatible)
	fmt.Fprintf(&output, "Declared capabilities: status=%s ingest_supported=%t render_supported=%t restore_supported=%t ocr_required_for_semantic_output=%t\n", report.Capabilities.Declared.Status, report.Capabilities.Declared.IngestSupported, report.Capabilities.Declared.RenderSupported, report.Capabilities.Declared.RestoreSupported, report.Capabilities.Declared.OCRRequiredForSemanticOutput)
	fmt.Fprintf(&output, "Installed capabilities: ingest=%t render=%t restore=%t validate=%t inspect=%t\n", report.Capabilities.Installed.Ingest, report.Capabilities.Installed.Render, report.Capabilities.Installed.Restore, report.Capabilities.Installed.Validate, report.Capabilities.Installed.Inspect)
	fmt.Fprintf(&output, "Integrity: %s (%d %s, %d bytes)\n", report.Integrity.Status, report.Integrity.AssetCount, plural(report.Integrity.AssetCount, "asset", "assets"), report.Integrity.TotalBytes)
	fmt.Fprintf(&output, "Assets (%d):\n", len(report.Assets))
	for _, asset := range report.Assets {
		fmt.Fprintf(&output, "  [%d] primary=%t role=%s file_name=%s media_type=%s bytes=%d\n", asset.Index, asset.Primary, asset.Role, strconv.Quote(asset.FileName), nullableQuotedString(asset.MediaType), asset.Bytes)
		if asset.Encoding == nil {
			output.WriteString("      encoding: null\n")
		} else {
			fmt.Fprintf(&output, "      encoding: bom=%s line_endings=%s detected_encoding=%s confidence=%s\n", nullableQuotedString(asset.Encoding.BOM), asset.Encoding.LineEndings, nullableQuotedString(asset.Encoding.DetectedEncoding), nullableFloat64(asset.Encoding.Confidence))
		}
	}
	fmt.Fprintf(&output, "Document: cues=%d non_cue_blocks=%d body_items=%d media_start_milliseconds=%s media_end_milliseconds=%s media_span_milliseconds=%s has_word_level_timing=%t\n", report.Document.CueCount, report.Document.NonCueBlockCount, report.Document.BodyItemCount, nullableInt64(report.Document.MediaStartMilliseconds), nullableInt64(report.Document.MediaEndMilliseconds), nullableInt64(report.Document.MediaSpanMilliseconds), report.Document.HasWordLevelTiming)
	fmt.Fprintf(&output, "Cues (%d):\n", len(report.Cues))
	for _, cue := range report.Cues {
		fmt.Fprintf(&output, "  ordinal=%d source_order=%d start_milliseconds=%d end_milliseconds=%d duration_milliseconds=%d payload_lines=%d speakers=%d tokens=%d ocr_observations=%d has_source_identifier=%t has_placement=%t native_settings=%d native_setting_occurrences=%d invalid_native_settings=%d\n", cue.Ordinal, cue.SourceOrder, cue.StartMilliseconds, cue.EndMilliseconds, cue.DurationMilliseconds, cue.PayloadLineCount, cue.SpeakerCount, cue.TokenCount, cue.OCRObservationCount, cue.HasSourceIdentifier, cue.HasPlacement, cue.NativeSettingCount, cue.NativeSettingOccurrenceCount, cue.InvalidNativeSettingCount)
	}
	fmt.Fprintf(&output, "Blocks (%d): note=%d style=%d region=%d unknown=%d\n", report.Blocks.Total, report.Blocks.NoteCount, report.Blocks.StyleCount, report.Blocks.RegionCount, report.Blocks.UnknownCount)
	for _, block := range report.Blocks.Entries {
		fmt.Fprintf(&output, "  source_order=%d type=%s lines=%d has_parsed_region=%t native_settings=%d native_setting_occurrences=%d invalid_native_settings=%d\n", block.SourceOrder, block.Type, block.LineCount, block.HasParsedRegion, block.NativeSettingCount, block.NativeSettingOccurrenceCount, block.InvalidNativeSettingCount)
	}
	fmt.Fprintf(&output, "Diagnostics (%d): info=%d warning=%d error=%d\n", report.Document.DiagnosticCount, report.Document.InfoCount, report.Document.WarningCount, report.Document.ErrorCount)
	for _, diagnostic := range report.Diagnostics {
		fmt.Fprintf(&output, "  severity=%s code=%s source_order=%s cue_ordinal=%s\n", diagnostic.Severity, diagnostic.Code, nullableInt(diagnostic.SourceOrder), nullableInt(diagnostic.CueOrdinal))
	}
	fmt.Fprintf(&output, "Loss: %s (%s)\n", report.Loss.Status, report.Loss.Reason)
	return []byte(output.String()), nil
}

func nullableString(value *string) string {
	if value == nil {
		return "null"
	}
	return *value
}

func nullableQuotedString(value *string) string {
	if value == nil {
		return "null"
	}
	return strconv.Quote(*value)
}

func nullableFloat64(value *float64) string {
	if value == nil {
		return "null"
	}
	return strconv.FormatFloat(*value, 'g', -1, 64)
}

func nullableInt(value *int) string {
	if value == nil {
		return "null"
	}
	return strconv.Itoa(*value)
}

func nullableInt64(value *int64) string {
	if value == nil {
		return "null"
	}
	return strconv.FormatInt(*value, 10)
}

func plural(count int, singular, pluralValue string) string {
	if count == 1 {
		return singular
	}
	return pluralValue
}
