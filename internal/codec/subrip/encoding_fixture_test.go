package subrip

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/source"
)

func TestEncodingAndLineEndingFixturesPreserveExactCapture(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..", "..", "testdata", "fixtures", "subrip", "encoding-line-endings", "source")
	tests := []struct {
		name           string
		override       string
		encoding       string
		bom            string
		lineEndings    string
		decodedRawText string
	}{
		{name: "utf8-bom-crlf.srt", encoding: codec.EncodingUTF8, bom: codec.EncodingUTF8, lineEndings: "crlf", decodedRawText: "Café déjà vu."},
		{name: "utf16le-bom.srt", encoding: codec.EncodingUTF16LE, bom: codec.EncodingUTF16LE, lineEndings: "lf", decodedRawText: "Café déjà vu."},
		{name: "utf16be-bom.srt", encoding: codec.EncodingUTF16BE, bom: codec.EncodingUTF16BE, lineEndings: "lf", decodedRawText: "Café déjà vu."},
		{name: "windows-1252.srt", override: codec.EncodingWindows1252, encoding: codec.EncodingWindows1252, lineEndings: "lf", decodedRawText: "“Café”"},
		{name: "iso-8859-1.srt", override: codec.EncodingISO88591, encoding: codec.EncodingISO88591, lineEndings: "lf", decodedRawText: "Café déjà vu."},
		{name: "lone-cr.srt", encoding: codec.EncodingUTF8, lineEndings: "cr", decodedRawText: "lone CR"},
		{name: "mixed.srt", encoding: codec.EncodingUTF8, lineEndings: "mixed", decodedRawText: "mixed"},
	}
	mediaType := "application/x-subrip"
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(root, test.name)
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			captured, err := source.CaptureContext(context.Background(), path, source.CaptureOptions{MediaType: &mediaType})
			if err != nil {
				t.Fatalf("CaptureContext() error = %v", err)
			}
			if !bytes.Equal(captured.Bytes, original) || captured.Asset.Size.Bytes != int64(len(original)) {
				t.Fatalf("captured bytes differ: size=%d want=%d", captured.Asset.Size.Bytes, len(original))
			}
			enveloped, err := base64.StdEncoding.DecodeString(captured.Asset.DataBase64)
			if err != nil || !bytes.Equal(enveloped, original) {
				t.Fatalf("source envelope bytes differ: error=%v", err)
			}
			wantHash := fmt.Sprintf("%x", sha256.Sum256(original))
			if captured.Asset.Hashes.SHA256 != wantHash {
				t.Fatalf("SHA-256 = %s, want %s", captured.Asset.Hashes.SHA256, wantHash)
			}
			if captured.Asset.FileName != test.name {
				t.Fatalf("file_name = %q, want portable basename %q", captured.Asset.FileName, test.name)
			}

			decoded, err := codec.DecodeText(captured.Bytes, test.override)
			if err != nil {
				t.Fatalf("DecodeText() error = %v", err)
			}
			if decoded.Observation.DetectedEncoding == nil || *decoded.Observation.DetectedEncoding != test.encoding {
				t.Fatalf("encoding = %#v, want %q", decoded.Observation.DetectedEncoding, test.encoding)
			}
			if test.bom == "" {
				if decoded.Observation.BOM != nil {
					t.Fatalf("BOM = %#v, want nil", decoded.Observation.BOM)
				}
			} else if decoded.Observation.BOM == nil || *decoded.Observation.BOM != test.bom {
				t.Fatalf("BOM = %#v, want %q", decoded.Observation.BOM, test.bom)
			}
			if decoded.Observation.LineEndings != test.lineEndings {
				t.Fatalf("line endings = %q, want %q", decoded.Observation.LineEndings, test.lineEndings)
			}
			parsed, err := Parse(decoded.Text, Options{})
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if len(parsed.Cues) != 1 || parsed.Cues[0].Payload.RawText != test.decodedRawText {
				t.Fatalf("decoded cue = %#v, want raw text %q", parsed.Cues, test.decodedRawText)
			}
		})
	}
}
