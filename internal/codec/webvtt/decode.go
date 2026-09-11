package webvtt

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/model"
)

// Detect reports whether payload begins with a boundary-valid WebVTT signature.
func Detect(payload []byte) bool {
	if bytes.HasPrefix(payload, []byte{0xef, 0xbb, 0xbf}) {
		payload = payload[3:]
	}
	if !bytes.HasPrefix(payload, []byte("WEBVTT")) {
		return false
	}
	return len(payload) == 6 || payload[6] == ' ' || payload[6] == '\t' || payload[6] == '\r' || payload[6] == '\n'
}

// DecodeUTF8 applies WebVTT's strict UTF-8-only decoding boundary.
func DecodeUTF8(payload []byte, override string) (DecodedText, error) {
	normalized := strings.ToLower(strings.TrimSpace(override))
	requiresBOM := false
	if override != "" {
		switch normalized {
		case "utf-8", "utf8":
		case "utf-8-bom", "utf8-bom", "utf-8-sig":
			requiresBOM = true
		default:
			return DecodedText{}, fmt.Errorf("WebVTT requires UTF-8; encoding %q is incompatible", override)
		}
	}
	bom := false
	if bytes.HasPrefix(payload, []byte{0xef, 0xbb, 0xbf}) {
		bom = true
		payload = payload[3:]
	}
	if !utf8.Valid(payload) {
		return DecodedText{}, fmt.Errorf("decode WebVTT: malformed UTF-8")
	}
	if requiresBOM && !bom {
		return DecodedText{}, fmt.Errorf("encoding %q requires a UTF-8 BOM", override)
	}
	encoding := "utf-8"
	confidence := 1.0
	observation := model.EncodingObservation{LineEndings: classifyLineEndings(payload), DetectedEncoding: &encoding, Confidence: &confidence}
	if bom {
		observed := "utf-8"
		observation.BOM = &observed
	}
	return DecodedText{Text: string(payload), Observation: observation}, nil
}

func classifyLineEndings(payload []byte) string {
	crlf, lf, cr := false, false, false
	for index := 0; index < len(payload); index++ {
		switch payload[index] {
		case '\r':
			if index+1 < len(payload) && payload[index+1] == '\n' {
				crlf = true
				index++
			} else {
				cr = true
			}
		case '\n':
			lf = true
		}
	}
	count := 0
	for _, found := range []bool{crlf, lf, cr} {
		if found {
			count++
		}
	}
	if count == 0 {
		return "none"
	}
	if count > 1 {
		return "mixed"
	}
	if crlf {
		return "crlf"
	}
	if cr {
		return "cr"
	}
	return "lf"
}
