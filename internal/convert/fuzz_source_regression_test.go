package convert_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/convert"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
)

func TestConversionFuzzSourceExactInMemory(t *testing.T) {
	tests := []struct {
		name, format, encoding string
		raw                    []byte
	}{
		{"subrip inferred actor disabled", "subrip", "", []byte("1\r\n00:00:01,000 --> 00:00:02,000\r\nAlice: Hello\n")},
		{"webvtt inferred voice disabled", "webvtt", "", []byte("\xef\xbb\xbfWEBVTT\r\n\r\n00:01.000 --> 00:02.000\r\n<v Alice>Hello\n")},
		{"legacy text explicit", "subrip", "windows-1252", []byte("1\n00:00:01,000 --> 00:00:02,000\ncaf\xe9\n")},
		{"empty ass", "ass", "", []byte("[Script Info]\nScriptType: v4.00+\n[V4+ Styles]\n")},
		{"empty ssa", "ssa", "", []byte("[Script Info]\nScriptType: v4.00\n[V4 Styles]\n")},
	}
	for _, format := range []string{"ass", "ssa"} {
		encoded, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			t.Fatal(err)
		}
		d, err := schema.Decode(encoded)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		if err != nil {
			t.Fatal(err)
		}
		raw = append([]byte{0xef, 0xbb, 0xbf}, bytes.ReplaceAll(raw, []byte("\n"), []byte("\r\n"))...)
		tests = append(tests, struct {
			name, format, encoding string
			raw                    []byte
		}{format + " native actor and BOM", format, "", raw})
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, accepted := conversionFuzzDocument(tc.raw, tc.format, tc.encoding)
			if !accepted {
				t.Fatal("valid acquired source rejected")
			}
			asset := d.Source.Assets[0]
			raw, err := base64.StdEncoding.Strict().DecodeString(asset.DataBase64)
			if err != nil || !bytes.Equal(raw, tc.raw) {
				t.Fatal("source bytes changed")
			}
			if err := source.ValidateIntegrity(context.Background(), d); err != nil {
				t.Fatal(err)
			}
			if asset.FileName != "fuzz-input.bin" || asset.Timestamps.CreatedSource != "unavailable" || asset.Timestamps.Created != nil || asset.Timestamps.Modified != nil || asset.Timestamps.Accessed != nil {
				t.Fatal("untruthful in-memory metadata")
			}
			encoded, err := json.Marshal(d)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := schema.Decode(encoded); err != nil {
				t.Fatalf("constructed source violates schema: %v", err)
			}
			if len(d.Cues) > 0 {
				if tc.format == "ass" || tc.format == "ssa" {
					if len(d.Cues[0].Speakers) != 1 || d.Cues[0].Speakers[0].Origin != "native" {
						t.Fatal("native actor disappeared")
					}
				} else if len(d.Cues[0].Speakers) != 0 {
					t.Fatal("inferred actor enabled")
				}
			}
			if strings.Contains(tc.name, "empty") && (len(d.Cues) != 0 || d.Document.MediaStartMilliseconds != nil || d.Stats.MediaSpanMilliseconds != nil) {
				t.Fatal("empty native summary invented")
			}
		})
	}
	for _, format := range []string{"subrip", "webvtt", "ass", "ssa", "unknown"} {
		if _, accepted := conversionFuzzDocument([]byte("malformed\xff"), format, ""); accepted {
			t.Fatalf("%s accepted malformed bytes", format)
		}
		if _, accepted := conversionFuzzDocument(bytes.Repeat([]byte{'x'}, (64<<10)+1), format, ""); accepted {
			t.Fatalf("%s accepted oversized fuzz setup", format)
		}
	}
}

func TestScriptedPrivateEnvelopeRefusalDoesNotPublishIdentity(t *testing.T) {
	const sentinel = "C:/private-envelope-sentinel/source.ass"
	for _, format := range []string{"ass", "ssa"} {
		encoded, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			t.Fatal(err)
		}
		d, err := schema.Decode(encoded)
		if err != nil {
			t.Fatal(err)
		}
		// Filename safety belongs to the independent source boundary. Integrity
		// rejection at render/conversion publication must hide this private value.
		d.Source.Assets[0].FileName = sentinel
		for _, target := range []string{"subrip", "webvtt", "ass", "ssa"} {
			if target == format {
				continue
			}
			result, err := convert.Convert(context.Background(), d, target, convert.Options{})
			if err == nil || len(result.Bytes) != 0 || strings.Contains(err.Error(), sentinel) {
				t.Fatalf("%s -> %s private envelope refusal: %v", format, target, err)
			}
		}
		rendered, err := scripted.Render(context.Background(), d, false)
		if err == nil || len(rendered.Bytes) != 0 || strings.Contains(err.Error(), sentinel) {
			t.Fatalf("%s private native render refusal: %v", format, err)
		}
	}
}
