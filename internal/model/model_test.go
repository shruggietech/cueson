package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDocumentValidateRepresentativeSemantics(t *testing.T) {
	t.Parallel()

	doc := representativeDocument()
	if err := doc.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDocumentValidateRejectsSemanticViolations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Document)
		want   string
	}{
		{name: "unresolved primary", mutate: func(doc *Document) { doc.Source.PrimaryAssetID = "missing" }, want: "primary_asset_id"},
		{name: "two primary roles", mutate: func(doc *Document) {
			doc.Source.Assets = append(doc.Source.Assets, companionAsset("asset-1"))
			doc.Source.Assets[1].Role = "primary"
		}, want: "primary"},
		{name: "duplicate asset id", mutate: func(doc *Document) { doc.Source.Assets = append(doc.Source.Assets, companionAsset("asset-0")) }, want: "asset id"},
		{name: "duration", mutate: func(doc *Document) { doc.Cues[0].Timing.DurationMilliseconds++ }, want: "duration_milliseconds"},
		{name: "ordinal", mutate: func(doc *Document) { doc.Cues[0].Ordinal = 1 }, want: "ordinal"},
		{name: "duplicate cue id", mutate: func(doc *Document) {
			duplicate := doc.Cues[0]
			duplicate.Ordinal = 1
			duplicate.SourceOrder = 1
			doc.Cues = append(doc.Cues, duplicate)
			syncCounts(doc)
		}, want: "cue id"},
		{name: "source order", mutate: func(doc *Document) {
			duplicate := doc.Cues[0]
			duplicate.ID = "cue-1"
			duplicate.Ordinal = 1
			doc.Cues = append(doc.Cues, duplicate)
			syncCounts(doc)
		}, want: "source_order"},
		{name: "cue count", mutate: func(doc *Document) { doc.Document.CueCount = 2 }, want: "cue_count"},
		{name: "timestamp instant", mutate: func(doc *Document) { doc.Source.Assets[0].Timestamps.Modified.UnixNS++ }, want: "timestamps.modified"},
		{name: "sub-nanosecond timestamp", mutate: func(doc *Document) { doc.Source.Assets[0].Timestamps.Modified.ISO = "2026-09-09T00:00:00.0000000001Z" }, want: "fractional precision"},
		{name: "creation provenance", mutate: func(doc *Document) { doc.Source.Assets[0].Timestamps.CreatedSource = "unavailable" }, want: "created_source"},
		{name: "format data", mutate: func(doc *Document) {
			doc.Cues[0].FormatData = CueFormatData{WebVTT: &WebVTTCueData{TimingLineRaw: "x"}}
		}, want: "format_data"},
		{name: "ocr source", mutate: func(doc *Document) {
			doc.Cues[0].OCRObservations = []OCRObservation{{ID: "ocr-0", Derived: true, Engine: "test", SourceAssetID: "missing", Text: "x", Lines: []string{"x"}}}
		}, want: "source_asset_id"},
		{name: "capability", mutate: func(doc *Document) { doc.FormatSupport.RestoreSupported = true }, want: "restore_supported"},
		{name: "diagnostic counts", mutate: func(doc *Document) { doc.Stats.WarningCount = 1 }, want: "warning_count"},
		{name: "media bounds", mutate: func(doc *Document) { value := int64(10); doc.Document.MediaStartMilliseconds = &value }, want: "media_start_milliseconds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := representativeDocument()
			tt.mutate(&doc)
			err := doc.Validate()
			if err == nil {
				t.Fatal("Validate() error = nil, want rejection")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Validate() error = %q, want mention of %q", err, tt.want)
			}
		})
	}
}

func TestDocumentJSONTagsStayCanonical(t *testing.T) {
	t.Parallel()

	encoded, err := json.Marshal(representativeDocument())
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{`"$schema"`, `"schema_version"`, `"format_support"`, `"ocr_observations"`, `"start_milliseconds"`, `"raw_text"`, `"plain_text"`} {
		if !strings.Contains(string(encoded), key) {
			t.Errorf("encoded model does not contain %s", key)
		}
	}
}

func representativeDocument() Document {
	start := int64(1250)
	end := int64(4200)
	span := int64(2950)
	newTimestamp := func() *Timestamp {
		return &Timestamp{ISO: "2026-09-09T00:00:00Z", UnixNS: 1788912000000000000}
	}
	return Document{
		Schema:        "https://cueson.io/schema/v0.0.0/cueson.schema.json",
		SchemaVersion: "0.0.0",
		Format:        "subrip",
		FormatSupport: FormatSupport{Status: "envelope_only"},
		Producer:      Producer{Name: "test", Version: "1.2.3"},
		Source: SourceEnvelope{PrimaryAssetID: "asset-0", Assets: []SourceAsset{{
			ID: "asset-0", Role: "primary", FileName: "captions.srt", MediaType: stringPointer("application/x-subrip"),
			Size: AssetSize{Bytes: 0, Text: stringPointer("0 bytes")}, Hashes: AssetHashes{SHA256: strings.Repeat("0", 64)},
			Timestamps: Timestamps{Created: newTimestamp(), Modified: newTimestamp(), Accessed: newTimestamp(), CreatedSource: "birthtime"}, DataBase64: "",
		}}},
		Metadata: Metadata{Title: stringPointer("Example"), Language: stringPointer("en"), Kind: stringPointer("subtitles")},
		Document: DocumentSummary{CueCount: 1, MediaStartMilliseconds: &start, MediaEndMilliseconds: &end, MediaSpanMilliseconds: &span},
		Cues: []Cue{{
			ID: "cue-0", Ordinal: 0, SourceOrder: 0, SourceIdentifier: stringPointer("1"),
			Timing:   Timing{StartMilliseconds: start, EndMilliseconds: end, DurationMilliseconds: span},
			Payload:  Payload{RawText: "Hello.", PlainText: "Hello.", Lines: []string{"Hello."}},
			Speakers: []Speaker{}, Tokens: []Token{}, OCRObservations: []OCRObservation{},
			FormatData: CueFormatData{SubRip: &SubRipCueData{SequenceLineRaw: "1", TimingLineRaw: "00:00:01,250 --> 00:00:04,200"}},
		}},
		FormatData:  DocumentFormatData{SubRip: &SubRipDocumentData{Dialect: "subrip"}},
		Diagnostics: []Diagnostic{},
		Stats:       Stats{CueCount: 1, MediaSpanMilliseconds: &span},
	}
}

func companionAsset(id string) SourceAsset {
	return SourceAsset{ID: id, Role: "companion", FileName: id + ".bin", Size: AssetSize{}, Hashes: AssetHashes{SHA256: strings.Repeat("0", 64)}, Timestamps: Timestamps{CreatedSource: "unavailable"}, DataBase64: ""}
}

func syncCounts(doc *Document) {
	doc.Document.CueCount = len(doc.Cues)
	doc.Stats.CueCount = len(doc.Cues)
}

func stringPointer(value string) *string {
	return &value
}
