// Package model defines Cue JSON types and format-neutral semantic invariants.
package model

import (
	"fmt"
	"time"
)

// Document is the fixed Cue JSON root contract.
type Document struct {
	Schema        string             `json:"$schema"`
	SchemaVersion string             `json:"schema_version"`
	Format        string             `json:"format"`
	FormatSupport FormatSupport      `json:"format_support"`
	Producer      Producer           `json:"producer"`
	Source        SourceEnvelope     `json:"source"`
	Metadata      Metadata           `json:"metadata"`
	Document      DocumentSummary    `json:"document"`
	Cues          []Cue              `json:"cues"`
	FormatData    DocumentFormatData `json:"format_data"`
	Diagnostics   []Diagnostic       `json:"diagnostics"`
	Stats         Stats              `json:"stats"`
}

type FormatSupport struct {
	Status                       string `json:"status"`
	IngestSupported              bool   `json:"ingest_supported"`
	RenderSupported              bool   `json:"render_supported"`
	RestoreSupported             bool   `json:"restore_supported"`
	OCRRequiredForSemanticOutput bool   `json:"ocr_required_for_semantic_output"`
}

type Producer struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SourceEnvelope struct {
	PrimaryAssetID string        `json:"primary_asset_id"`
	Assets         []SourceAsset `json:"assets"`
}

type SourceAsset struct {
	ID         string               `json:"id"`
	Role       string               `json:"role"`
	FileName   string               `json:"file_name"`
	MediaType  *string              `json:"media_type"`
	Size       AssetSize            `json:"size"`
	Hashes     AssetHashes          `json:"hashes"`
	Timestamps Timestamps           `json:"timestamps"`
	Encoding   *EncodingObservation `json:"encoding"`
	DataBase64 string               `json:"data_base64"`
}

type AssetSize struct {
	Bytes int64   `json:"bytes"`
	Text  *string `json:"text"`
}

type AssetHashes struct {
	SHA256 string `json:"sha256"`
}

type Timestamps struct {
	Created       *Timestamp `json:"created"`
	Modified      *Timestamp `json:"modified"`
	Accessed      *Timestamp `json:"accessed"`
	CreatedSource string     `json:"created_source"`
}

type Timestamp struct {
	ISO    string `json:"iso"`
	UnixNS int64  `json:"unix_ns"`
}

type EncodingObservation struct {
	BOM              *string  `json:"bom"`
	LineEndings      string   `json:"line_endings"`
	DetectedEncoding *string  `json:"detected_encoding"`
	Confidence       *float64 `json:"confidence"`
}

type Metadata struct {
	Title       *string `json:"title"`
	Language    *string `json:"language"`
	Kind        *string `json:"kind"`
	Description *string `json:"description"`
}

type DocumentSummary struct {
	CueCount               int    `json:"cue_count"`
	MediaStartMilliseconds *int64 `json:"media_start_milliseconds"`
	MediaEndMilliseconds   *int64 `json:"media_end_milliseconds"`
	MediaSpanMilliseconds  *int64 `json:"media_span_milliseconds"`
	HasWordLevelTiming     bool   `json:"has_word_level_timing"`
}

type Cue struct {
	ID               string           `json:"id"`
	Ordinal          int              `json:"ordinal"`
	SourceOrder      int              `json:"source_order"`
	SourceIdentifier *string          `json:"source_identifier"`
	Timing           Timing           `json:"timing"`
	Payload          Payload          `json:"payload"`
	Speakers         []Speaker        `json:"speakers"`
	Tokens           []Token          `json:"tokens"`
	OCRObservations  []OCRObservation `json:"ocr_observations"`
	Placement        *Placement       `json:"placement"`
	FormatData       CueFormatData    `json:"format_data"`
}

type Timing struct {
	StartMilliseconds    int64 `json:"start_milliseconds"`
	EndMilliseconds      int64 `json:"end_milliseconds"`
	DurationMilliseconds int64 `json:"duration_milliseconds"`
}

type Payload struct {
	RawText   string   `json:"raw_text"`
	PlainText string   `json:"plain_text"`
	Lines     []string `json:"lines"`
}

type Speaker struct {
	Name   string `json:"name"`
	Origin string `json:"origin"`
}

type Token struct {
	Text              string  `json:"text"`
	StartMilliseconds int64   `json:"start_milliseconds"`
	EndMilliseconds   int64   `json:"end_milliseconds"`
	NativeTiming      *string `json:"native_timing"`
}

type OCRObservation struct {
	ID                string         `json:"id"`
	Derived           bool           `json:"derived"`
	Engine            string         `json:"engine"`
	EngineVersion     *string        `json:"engine_version,omitempty"`
	Model             *string        `json:"model,omitempty"`
	Language          *string        `json:"language,omitempty"`
	Text              string         `json:"text"`
	Lines             []string       `json:"lines"`
	Confidence        *float64       `json:"confidence,omitempty"`
	Alternatives      []string       `json:"alternatives,omitempty"`
	SourceAssetID     string         `json:"source_asset_id"`
	SourceEventID     *string        `json:"source_event_id,omitempty"`
	SourceImageID     *string        `json:"source_image_id,omitempty"`
	SourceRegions     []SourceRegion `json:"source_regions,omitempty"`
	ProcessingOptions map[string]any `json:"processing_options,omitempty"`
}

type SourceRegion struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Placement struct {
	Line     *string `json:"line"`
	Position *string `json:"position"`
	Size     *string `json:"size"`
	Align    *string `json:"align"`
}

type CueFormatData struct {
	SubRip *SubRipCueData `json:"subrip,omitempty"`
	WebVTT *WebVTTCueData `json:"webvtt,omitempty"`
}

type SubRipCueData struct {
	SequenceLineRaw string       `json:"sequence_line_raw"`
	TimingLineRaw   string       `json:"timing_line_raw"`
	Coordinates     *Coordinates `json:"coordinates,omitempty"`
}

type Coordinates struct {
	X1 int `json:"x1"`
	X2 int `json:"x2"`
	Y1 int `json:"y1"`
	Y2 int `json:"y2"`
}

type WebVTTCueData struct {
	IdentifierRaw *string           `json:"identifier_raw"`
	TimingLineRaw string            `json:"timing_line_raw"`
	SettingsRaw   string            `json:"settings_raw"`
	Settings      map[string]string `json:"settings"`
	RawPayload    string            `json:"raw_payload"`
}

type DocumentFormatData struct {
	SubRip *SubRipDocumentData `json:"subrip,omitempty"`
	WebVTT *WebVTTDocumentData `json:"webvtt,omitempty"`
}

type SubRipDocumentData struct {
	Dialect       string  `json:"dialect"`
	DecodingNotes *string `json:"decoding_notes"`
}

type WebVTTDocumentData struct {
	Signature     string        `json:"signature"`
	Description   *string       `json:"description"`
	MetadataLines []string      `json:"metadata_lines"`
	Blocks        []WebVTTBlock `json:"blocks"`
}

type WebVTTBlock struct {
	Type        string `json:"type"`
	SourceOrder int    `json:"source_order"`
	Raw         string `json:"raw"`
}

type Diagnostic struct {
	Severity    string  `json:"severity"`
	Code        string  `json:"code"`
	Message     string  `json:"message"`
	SourceOrder *int    `json:"source_order"`
	CueID       *string `json:"cue_id"`
}

type Stats struct {
	CueCount              int    `json:"cue_count"`
	DiagnosticCount       int    `json:"diagnostic_count"`
	WarningCount          int    `json:"warning_count"`
	ErrorCount            int    `json:"error_count"`
	HasWordLevelTiming    bool   `json:"has_word_level_timing"`
	MediaSpanMilliseconds *int64 `json:"media_span_milliseconds"`
}

// Validate checks invariants that are intentionally clearer outside JSON Schema.
func (document Document) Validate() error {
	if err := validateCapability(document.FormatSupport); err != nil {
		return err
	}
	assetIDs, err := validateAssets(document.Source)
	if err != nil {
		return err
	}
	if err := validateFormatData(document.Format, document.FormatData, document.Cues); err != nil {
		return err
	}
	if err := validateCues(document.Cues, assetIDs, document.FormatData); err != nil {
		return err
	}
	if err := validateSummaries(document); err != nil {
		return err
	}
	return validateDiagnostics(document.Diagnostics, document.Cues, document.Stats)
}

func validateCapability(support FormatSupport) error {
	if support.Status != "envelope_only" {
		return fmt.Errorf("format_support.status must be envelope_only")
	}
	if support.IngestSupported {
		return fmt.Errorf("format_support.ingest_supported must remain false")
	}
	if support.RenderSupported {
		return fmt.Errorf("format_support.render_supported must remain false")
	}
	if support.RestoreSupported {
		return fmt.Errorf("format_support.restore_supported must remain false until restore ships")
	}
	if support.OCRRequiredForSemanticOutput {
		return fmt.Errorf("format_support.ocr_required_for_semantic_output must be false")
	}
	return nil
}

func validateAssets(source SourceEnvelope) (map[string]struct{}, error) {
	assetIDs := make(map[string]struct{}, len(source.Assets))
	primaryCount := 0
	for index := range source.Assets {
		asset := &source.Assets[index]
		if _, exists := assetIDs[asset.ID]; exists {
			return nil, fmt.Errorf("duplicate asset id %q", asset.ID)
		}
		assetIDs[asset.ID] = struct{}{}
		if asset.Role == "primary" {
			primaryCount++
			if asset.ID != source.PrimaryAssetID {
				return nil, fmt.Errorf("primary asset %q does not match primary_asset_id %q", asset.ID, source.PrimaryAssetID)
			}
		}
		if err := validateTimestamps(asset.Timestamps, index); err != nil {
			return nil, err
		}
	}
	if _, exists := assetIDs[source.PrimaryAssetID]; !exists {
		return nil, fmt.Errorf("primary_asset_id %q does not resolve", source.PrimaryAssetID)
	}
	if primaryCount != 1 {
		return nil, fmt.Errorf("source must contain exactly one primary asset, found %d", primaryCount)
	}
	return assetIDs, nil
}

func validateTimestamps(timestamps Timestamps, assetIndex int) error {
	ordered := []struct {
		name  string
		value *Timestamp
	}{
		{name: "created", value: timestamps.Created},
		{name: "modified", value: timestamps.Modified},
		{name: "accessed", value: timestamps.Accessed},
	}
	for _, observation := range ordered {
		if observation.value == nil {
			continue
		}
		parsed, err := time.Parse(time.RFC3339Nano, observation.value.ISO)
		if err != nil {
			return fmt.Errorf("source.assets[%d].timestamps.%s.iso: %w", assetIndex, observation.name, err)
		}
		if !parsed.Equal(time.Unix(0, observation.value.UnixNS)) {
			return fmt.Errorf("source.assets[%d].timestamps.%s iso and unix_ns identify different instants", assetIndex, observation.name)
		}
	}
	if timestamps.CreatedSource == "unavailable" && timestamps.Created != nil {
		return fmt.Errorf("source.assets[%d].timestamps.created_source unavailable requires a null created timestamp", assetIndex)
	}
	if timestamps.CreatedSource != "unavailable" && timestamps.Created == nil {
		return fmt.Errorf("source.assets[%d].timestamps.created_source %q requires a created timestamp", assetIndex, timestamps.CreatedSource)
	}
	return nil
}

func validateFormatData(format string, root DocumentFormatData, cues []Cue) error {
	switch format {
	case "subrip":
		if root.SubRip == nil || root.WebVTT != nil {
			return fmt.Errorf("format_data must contain only subrip for format subrip")
		}
		for index := range cues {
			if cues[index].FormatData.SubRip == nil || cues[index].FormatData.WebVTT != nil {
				return fmt.Errorf("cues[%d].format_data must contain only subrip", index)
			}
		}
	case "webvtt":
		if root.WebVTT == nil || root.SubRip != nil {
			return fmt.Errorf("format_data must contain only webvtt for format webvtt")
		}
		for index := range cues {
			if cues[index].FormatData.WebVTT == nil || cues[index].FormatData.SubRip != nil {
				return fmt.Errorf("cues[%d].format_data must contain only webvtt", index)
			}
		}
	default:
		return fmt.Errorf("format %q is not recognized", format)
	}
	return nil
}

func validateCues(cues []Cue, assetIDs map[string]struct{}, rootData DocumentFormatData) error {
	cueIDs := make(map[string]struct{}, len(cues))
	sourceOrders := make(map[int]string, len(cues))
	previousOrder := -1
	for index := range cues {
		cue := &cues[index]
		if _, exists := cueIDs[cue.ID]; exists {
			return fmt.Errorf("duplicate cue id %q", cue.ID)
		}
		cueIDs[cue.ID] = struct{}{}
		if cue.Ordinal != index {
			return fmt.Errorf("cues[%d].ordinal is %d, want %d", index, cue.Ordinal, index)
		}
		if cue.SourceOrder <= previousOrder {
			return fmt.Errorf("cues[%d].source_order must be strictly increasing", index)
		}
		previousOrder = cue.SourceOrder
		sourceOrders[cue.SourceOrder] = fmt.Sprintf("cue %q", cue.ID)
		if cue.Timing.EndMilliseconds < cue.Timing.StartMilliseconds {
			return fmt.Errorf("cues[%d].timing.end_milliseconds precedes start_milliseconds", index)
		}
		if cue.Timing.DurationMilliseconds != cue.Timing.EndMilliseconds-cue.Timing.StartMilliseconds {
			return fmt.Errorf("cues[%d].timing.duration_milliseconds is inconsistent", index)
		}
		for tokenIndex := range cue.Tokens {
			token := &cue.Tokens[tokenIndex]
			if token.EndMilliseconds < token.StartMilliseconds || token.StartMilliseconds < cue.Timing.StartMilliseconds || token.EndMilliseconds > cue.Timing.EndMilliseconds {
				return fmt.Errorf("cues[%d].tokens[%d] timing falls outside its cue", index, tokenIndex)
			}
		}
		observationIDs := make(map[string]struct{}, len(cue.OCRObservations))
		for observationIndex := range cue.OCRObservations {
			observation := &cue.OCRObservations[observationIndex]
			if !observation.Derived {
				return fmt.Errorf("cues[%d].ocr_observations[%d].derived must be true", index, observationIndex)
			}
			if _, exists := observationIDs[observation.ID]; exists {
				return fmt.Errorf("duplicate ocr observation id %q in cue %q", observation.ID, cue.ID)
			}
			observationIDs[observation.ID] = struct{}{}
			if _, exists := assetIDs[observation.SourceAssetID]; !exists {
				return fmt.Errorf("cues[%d].ocr_observations[%d].source_asset_id %q does not resolve", index, observationIndex, observation.SourceAssetID)
			}
		}
	}
	if rootData.WebVTT != nil {
		previousBlockOrder := -1
		for index, block := range rootData.WebVTT.Blocks {
			if block.SourceOrder <= previousBlockOrder {
				return fmt.Errorf("format_data.webvtt.blocks[%d].source_order must be strictly increasing", index)
			}
			previousBlockOrder = block.SourceOrder
			if owner, exists := sourceOrders[block.SourceOrder]; exists {
				return fmt.Errorf("format_data.webvtt.blocks[%d].source_order duplicates %s", index, owner)
			}
			sourceOrders[block.SourceOrder] = fmt.Sprintf("webvtt block %d", index)
		}
	}
	return nil
}

func validateSummaries(document Document) error {
	if document.Document.CueCount != len(document.Cues) || document.Stats.CueCount != len(document.Cues) {
		return fmt.Errorf("document and stats cue_count must equal %d", len(document.Cues))
	}
	if len(document.Cues) == 0 {
		return fmt.Errorf("cues must contain at least one cue")
	}
	minimumStart := document.Cues[0].Timing.StartMilliseconds
	maximumEnd := document.Cues[0].Timing.EndMilliseconds
	hasTokens := false
	for index := range document.Cues {
		cue := &document.Cues[index]
		if cue.Timing.StartMilliseconds < minimumStart {
			minimumStart = cue.Timing.StartMilliseconds
		}
		if cue.Timing.EndMilliseconds > maximumEnd {
			maximumEnd = cue.Timing.EndMilliseconds
		}
		hasTokens = hasTokens || len(cue.Tokens) > 0
	}
	span := maximumEnd - minimumStart
	if !equalsInt64(document.Document.MediaStartMilliseconds, minimumStart) {
		return fmt.Errorf("document.media_start_milliseconds must equal %d", minimumStart)
	}
	if !equalsInt64(document.Document.MediaEndMilliseconds, maximumEnd) {
		return fmt.Errorf("document.media_end_milliseconds must equal %d", maximumEnd)
	}
	if !equalsInt64(document.Document.MediaSpanMilliseconds, span) || !equalsInt64(document.Stats.MediaSpanMilliseconds, span) {
		return fmt.Errorf("document and stats media_span_milliseconds must equal %d", span)
	}
	if document.Document.HasWordLevelTiming != hasTokens || document.Stats.HasWordLevelTiming != hasTokens {
		return fmt.Errorf("has_word_level_timing must reflect token presence")
	}
	return nil
}

func validateDiagnostics(diagnostics []Diagnostic, cues []Cue, stats Stats) error {
	cueIDs := make(map[string]struct{}, len(cues))
	for index := range cues {
		cueIDs[cues[index].ID] = struct{}{}
	}
	warnings := 0
	errors := 0
	for index := range diagnostics {
		diagnostic := &diagnostics[index]
		switch diagnostic.Severity {
		case "warning":
			warnings++
		case "error":
			errors++
		}
		if diagnostic.CueID != nil {
			if _, exists := cueIDs[*diagnostic.CueID]; !exists {
				return fmt.Errorf("diagnostics[%d].cue_id %q does not resolve", index, *diagnostic.CueID)
			}
		}
	}
	if stats.DiagnosticCount != len(diagnostics) {
		return fmt.Errorf("stats.diagnostic_count must equal %d", len(diagnostics))
	}
	if stats.WarningCount != warnings {
		return fmt.Errorf("stats.warning_count must equal %d", warnings)
	}
	if stats.ErrorCount != errors {
		return fmt.Errorf("stats.error_count must equal %d", errors)
	}
	return nil
}

func equalsInt64(value *int64, want int64) bool {
	return value != nil && *value == want
}
