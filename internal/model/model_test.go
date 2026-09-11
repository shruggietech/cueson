package model

import (
	"encoding/json"
	"fmt"
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

func TestDocumentValidateCapabilityProfiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		format  string
		support FormatSupport
		valid   bool
	}{
		{name: "subrip stable", format: "subrip", support: FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}, valid: true},
		{name: "subrip experimental", format: "subrip", support: FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true}},
		{name: "subrip envelope only", format: "subrip", support: FormatSupport{Status: "envelope_only", RestoreSupported: true}},
		{name: "webvtt stable", format: "webvtt", support: FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}, valid: true},
		{name: "webvtt experimental", format: "webvtt", support: FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true}},
		{name: "webvtt envelope only", format: "webvtt", support: FormatSupport{Status: "envelope_only", RestoreSupported: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := representativeDocument()
			if tt.format == "webvtt" {
				doc = representativeWebVTTDocument()
			}
			doc.Format = tt.format
			doc.FormatSupport = tt.support
			err := doc.Validate()
			if tt.valid && err != nil {
				t.Fatalf("Validate() error = %v", err)
			}
			if !tt.valid && err == nil {
				t.Fatal("Validate() error = nil, want capability rejection")
			}
		})
	}
}

func TestDocumentValidateWebVTTNativeModel(t *testing.T) {
	t.Parallel()

	if err := representativeWebVTTDocument().Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestDocumentValidateRejectsWebVTTNativeInconsistency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*Document)
		want   string
	}{
		{name: "source order duplicate", mutate: func(doc *Document) { doc.Cues[0].SourceOrder = 0 }, want: "source_order"},
		{name: "source order gap", mutate: func(doc *Document) { doc.Cues[0].SourceOrder = 2 }, want: "contiguous"},
		{name: "block raw lines", mutate: func(doc *Document) { doc.FormatData.WebVTT.Blocks[0].RawLines[1] = "id:other" }, want: "raw_lines"},
		{name: "block embedded line ending", mutate: func(doc *Document) {
			doc.FormatData.WebVTT.Blocks[0].RawLines = []string{"REGION", "id:fred\nwidth:40%"}
			doc.FormatData.WebVTT.Blocks[0].Raw = strings.Join(doc.FormatData.WebVTT.Blocks[0].RawLines, "\n")
		}, want: "raw_lines"},
		{name: "region missing", mutate: func(doc *Document) { doc.FormatData.WebVTT.Blocks[0].Region = nil }, want: "region"},
		{name: "region on note", mutate: func(doc *Document) { doc.FormatData.WebVTT.Blocks[0].Type = "note" }, want: "region"},
		{name: "payload raw", mutate: func(doc *Document) { doc.Cues[0].FormatData.WebVTT.RawPayload = "different" }, want: "raw_payload"},
		{name: "payload lines", mutate: func(doc *Document) { doc.Cues[0].FormatData.WebVTT.RawPayloadLines = []string{"different"} }, want: "raw_payload_lines"},
		{name: "payload embedded line ending", mutate: func(doc *Document) {
			lines := []string{"bad\nline"}
			doc.Cues[0].FormatData.WebVTT.RawPayloadLines = lines
			doc.Cues[0].FormatData.WebVTT.RawPayload = strings.Join(lines, "\n")
			doc.Cues[0].Payload.Lines = lines
			doc.Cues[0].Payload.RawText = strings.Join(lines, "\n")
		}, want: "raw_payload_lines"},
		{name: "payload lines empty", mutate: func(doc *Document) {
			doc.Cues[0].FormatData.WebVTT.RawPayloadLines = []string{}
			doc.Cues[0].FormatData.WebVTT.RawPayload = ""
			doc.Cues[0].Payload.Lines = []string{}
			doc.Cues[0].Payload.RawText = ""
		}, want: "raw_payload_lines"},
		{name: "payload line blank", mutate: func(doc *Document) {
			doc.Cues[0].FormatData.WebVTT.RawPayloadLines = []string{""}
			doc.Cues[0].FormatData.WebVTT.RawPayload = ""
			doc.Cues[0].Payload.Lines = []string{""}
			doc.Cues[0].Payload.RawText = ""
		}, want: "raw_payload_lines"},
		{name: "identifier", mutate: func(doc *Document) { doc.Cues[0].SourceIdentifier = stringPointer("other") }, want: "identifier_raw"},
		{name: "speaker origin", mutate: func(doc *Document) { doc.Cues[0].Speakers[0].Origin = "heuristic" }, want: "speakers"},
		{name: "token order", mutate: func(doc *Document) { doc.Cues[0].Tokens[1].StartMilliseconds = 1900 }, want: "tokens"},
		{name: "zero duration token", mutate: func(doc *Document) { doc.Cues[0].Tokens[0].EndMilliseconds = doc.Cues[0].Tokens[0].StartMilliseconds }, want: "tokens"},
		{name: "unrecognized valid occurrence", mutate: func(doc *Document) { doc.Cues[0].FormatData.WebVTT.SettingOccurrences[0].Recognized = false }, want: "recognized"},
		{name: "empty occurrence raw", mutate: func(doc *Document) { doc.Cues[0].FormatData.WebVTT.SettingOccurrences[0].Raw = "" }, want: "raw"},
		{name: "valid occurrence empty value", mutate: func(doc *Document) { doc.Cues[0].FormatData.WebVTT.SettingOccurrences[0].Value = "" }, want: "valid"},
		{name: "unknown effective setting", mutate: func(doc *Document) { doc.Cues[0].FormatData.WebVTT.Settings["unknown"] = "value" }, want: "settings"},
		{name: "unknown region setting", mutate: func(doc *Document) { doc.FormatData.WebVTT.Blocks[0].Region.Settings["unknown"] = "value" }, want: "settings"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := representativeWebVTTDocument()
			tt.mutate(&doc)
			err := doc.Validate()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Validate() error = %v, want rejection containing %q", err, tt.want)
			}
		})
	}
}

func TestDocumentValidatePreservesSubRipContract(t *testing.T) {
	t.Parallel()

	doc := representativeDocument()
	doc.Cues[0].Speakers = []Speaker{{Name: "Narrator", Origin: "heuristic"}}
	if err := doc.Validate(); err != nil {
		t.Fatalf("Validate() rejected existing SubRip model: %v", err)
	}
}

func TestDocumentValidateRejectsCollectionAmplification(t *testing.T) {
	tests := []struct {
		name string
		doc  func() Document
		want string
	}{
		{name: "assets", doc: func() Document {
			doc := representativeDocument()
			doc.Source.Assets = make([]SourceAsset, MaxDocumentItems+1)
			return doc
		}, want: "source.assets"},
		{name: "cues", doc: func() Document {
			doc := representativeDocument()
			doc.Cues = make([]Cue, MaxDocumentItems+1)
			return doc
		}, want: "cues"},
		{name: "diagnostics", doc: func() Document {
			doc := representativeDocument()
			doc.Diagnostics = make([]Diagnostic, MaxDiagnostics+1)
			return doc
		}, want: "diagnostics"},
		{name: "payload lines", doc: func() Document {
			doc := representativeDocument()
			doc.Cues[0].Payload.Lines = make([]string, MaxItemOccurrences+1)
			return doc
		}, want: "payload.lines"},
		{name: "ocr nested alternatives", doc: func() Document {
			doc := representativeDocument()
			doc.Cues[0].OCRObservations = []OCRObservation{{Alternatives: make([]string, MaxItemOccurrences+1)}}
			return doc
		}, want: "alternatives"},
		{name: "webvtt blocks", doc: func() Document {
			doc := representativeWebVTTDocument()
			doc.FormatData.WebVTT.Blocks = make([]WebVTTBlock, MaxDocumentItems+1)
			return doc
		}, want: "format_data.webvtt.blocks"},
		{name: "webvtt settings", doc: func() Document {
			doc := representativeWebVTTDocument()
			doc.Cues[0].FormatData.WebVTT.Settings = make(map[string]string, MaxItemOccurrences+1)
			for index := 0; index <= MaxItemOccurrences; index++ {
				doc.Cues[0].FormatData.WebVTT.Settings[fmt.Sprintf("setting_%d", index)] = "x"
			}
			return doc
		}, want: "format_data.webvtt.settings"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc().Validate()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Validate() error = %v, want bounded rejection containing %q", err, tt.want)
			}
		})
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
		{name: "capability", mutate: func(doc *Document) { doc.FormatSupport.RestoreSupported = false }, want: "restore_supported"},
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
		Schema:        "https://cueson.io/schema/v1.0.0/cueson.schema.json",
		SchemaVersion: "1.0.0",
		Format:        "subrip",
		FormatSupport: FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
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

func representativeWebVTTDocument() Document {
	doc := representativeDocument()
	doc.Format = "webvtt"
	doc.FormatSupport = FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true}
	doc.Source.Assets[0].FileName = "captions.vtt"
	doc.Cues[0].SourceOrder = 1
	doc.Cues[0].SourceIdentifier = stringPointer("intro")
	doc.Cues[0].Payload = Payload{RawText: "<v Narrator>Hello <00:00:02.000>world.</v>", PlainText: "Hello world.", Lines: []string{"<v Narrator>Hello <00:00:02.000>world.</v>"}}
	doc.Cues[0].Speakers = []Speaker{{Name: "Narrator", Origin: "native"}}
	doc.Cues[0].Tokens = []Token{
		{Text: "Hello ", StartMilliseconds: 1250, EndMilliseconds: 2000},
		{Text: "world.", StartMilliseconds: 2000, EndMilliseconds: 4200, NativeTiming: stringPointer("00:00:02.000")},
	}
	doc.Cues[0].FormatData = CueFormatData{WebVTT: &WebVTTCueData{
		IdentifierRaw: stringPointer("intro"), TimingLineRaw: "00:01.250 --> 00:04.200 align:start mystery:opaque", SettingsRaw: "align:start mystery:opaque",
		Settings: map[string]string{"align": "start"},
		SettingOccurrences: []WebVTTSettingOccurrence{
			{Raw: "align:start", Name: "align", Value: "start", Recognized: true, Valid: true},
			{Raw: "mystery:opaque", Name: "mystery", Value: "opaque", Recognized: false, Valid: false},
			{Raw: ":x", Name: "", Value: "x", Recognized: false, Valid: false},
			{Raw: "bare", Name: "bare", Value: "", Recognized: false, Valid: false},
		},
		RawPayload: "<v Narrator>Hello <00:00:02.000>world.</v>", RawPayloadLines: []string{"<v Narrator>Hello <00:00:02.000>world.</v>"},
	}}
	doc.FormatData = DocumentFormatData{WebVTT: &WebVTTDocumentData{
		Signature: "WEBVTT", SignatureLineRaw: "WEBVTT Example", Description: stringPointer("Example"), MetadataLines: []string{"Kind: captions"},
		Blocks: []WebVTTBlock{{
			Type: "region", SourceOrder: 0, Raw: "REGION\nid:fred\nwidth:40%", RawLines: []string{"REGION", "id:fred", "width:40%"},
			Region: &WebVTTRegionData{SettingsRaw: "id:fred width:40%", Settings: map[string]string{"id": "fred", "width": "40%"}, SettingOccurrences: []WebVTTSettingOccurrence{
				{Raw: "id:fred", Name: "id", Value: "fred", Recognized: true, Valid: true},
				{Raw: "width:40%", Name: "width", Value: "40%", Recognized: true, Valid: true},
			}},
		}},
	}}
	doc.Document.HasWordLevelTiming = true
	doc.Stats.HasWordLevelTiming = true
	return doc
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
