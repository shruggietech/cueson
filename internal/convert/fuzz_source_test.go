package convert_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/version"
)

// conversionFuzzDocument constructs acquired source truth without touching the
// filesystem. Fixed absent timestamps are truthful for these in-memory assets.
func conversionFuzzDocument(raw []byte, format, encoding string) (model.Document, bool) {
	if len(raw) > 64<<10 {
		return model.Document{}, false
	}
	d := model.Document{
		Schema: schema.ID(), SchemaVersion: schema.Version(), Format: format,
		FormatSupport: model.FormatSupport{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
		Producer:      model.Producer{Name: "cueson", Version: version.String()},
		Source: model.SourceEnvelope{PrimaryAssetID: "asset-0", Assets: []model.SourceAsset{{
			ID: "asset-0", Role: "primary", FileName: "fuzz-input.bin",
			Size: model.AssetSize{Bytes: int64(len(raw))}, Hashes: model.AssetHashes{SHA256: fmt.Sprintf("%x", sha256.Sum256(raw))},
			Timestamps: model.Timestamps{CreatedSource: "unavailable"}, DataBase64: base64.StdEncoding.EncodeToString(raw),
		}}},
		Diagnostics: []model.Diagnostic{},
	}
	switch format {
	case "subrip":
		decoded, err := codec.DecodeText(raw, encoding)
		if err != nil {
			return model.Document{}, false
		}
		parsed, err := subrip.Parse(decoded.Text, subrip.Options{})
		if err != nil {
			return model.Document{}, false
		}
		d.Source.Assets[0].Encoding = &decoded.Observation
		d.Cues, d.Diagnostics = parsed.Cues, parsed.Diagnostics
		d.FormatData.SubRip = &model.SubRipDocumentData{Dialect: "subrip"}
	case "webvtt":
		decoded, err := webvtt.DecodeUTF8(raw, encoding)
		if err != nil {
			return model.Document{}, false
		}
		parsed, err := webvtt.Parse(decoded.Text)
		if err != nil {
			return model.Document{}, false
		}
		d.Source.Assets[0].Encoding = &decoded.Observation
		d.Cues, d.Diagnostics = parsed.Cues, parsed.Diagnostics
		for i := range d.Cues {
			d.Cues[i].Speakers = []model.Speaker{}
		}
		d.FormatData.WebVTT = &parsed.DocumentData
	case "ass", "ssa":
		parsed, err := scripted.Parse(context.Background(), raw, format)
		if err != nil {
			return model.Document{}, false
		}
		d.FormatSupport.Status = "experimental"
		d.Source.Assets[0].Encoding = &parsed.Encoding
		d.Cues, d.Diagnostics = parsed.Cues, parsed.Diagnostics
		if format == "ass" {
			d.FormatData.ASS = &parsed.Native
		} else {
			d.FormatData.SSA = &parsed.Native
		}
	default:
		return model.Document{}, false
	}
	d.Document.CueCount, d.Stats.CueCount = len(d.Cues), len(d.Cues)
	d.Stats.DiagnosticCount = len(d.Diagnostics)
	for _, diagnostic := range d.Diagnostics {
		if diagnostic.Severity == "warning" {
			d.Stats.WarningCount++
		}
		if diagnostic.Severity == "error" {
			d.Stats.ErrorCount++
		}
	}
	if len(d.Cues) > 0 {
		start, end := d.Cues[0].Timing.StartMilliseconds, d.Cues[0].Timing.EndMilliseconds
		for _, cue := range d.Cues {
			start, end = min(start, cue.Timing.StartMilliseconds), max(end, cue.Timing.EndMilliseconds)
			d.Document.HasWordLevelTiming = d.Document.HasWordLevelTiming || len(cue.Tokens) > 0
		}
		span := end - start
		d.Document.MediaStartMilliseconds, d.Document.MediaEndMilliseconds, d.Document.MediaSpanMilliseconds = &start, &end, &span
		d.Stats.MediaSpanMilliseconds = &span
	}
	d.Stats.HasWordLevelTiming = d.Document.HasWordLevelTiming
	if err := d.Validate(); err != nil {
		return model.Document{}, false
	}
	return d, true
}
