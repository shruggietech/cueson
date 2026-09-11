package codec

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/shruggietech/cueson/internal/model"
	"golang.org/x/text/encoding/charmap"
)

const (
	EncodingUTF8        = "utf-8"
	EncodingUTF8BOM     = "utf-8-bom"
	EncodingUTF16LE     = "utf-16le"
	EncodingUTF16BE     = "utf-16be"
	EncodingWindows1252 = "windows-1252"
	EncodingISO88591    = "iso-8859-1"
)

var encodingAliases = map[string]string{
	"utf-8":        EncodingUTF8,
	"utf8":         EncodingUTF8,
	"utf-8-bom":    EncodingUTF8BOM,
	"utf8-bom":     EncodingUTF8BOM,
	"utf-8-sig":    EncodingUTF8BOM,
	"utf-16le":     EncodingUTF16LE,
	"utf16le":      EncodingUTF16LE,
	"utf-16-le":    EncodingUTF16LE,
	"utf-16be":     EncodingUTF16BE,
	"utf16be":      EncodingUTF16BE,
	"utf-16-be":    EncodingUTF16BE,
	"windows-1252": EncodingWindows1252,
	"windows1252":  EncodingWindows1252,
	"cp1252":       EncodingWindows1252,
	"iso-8859-1":   EncodingISO88591,
	"iso8859-1":    EncodingISO88591,
	"latin1":       EncodingISO88591,
	"latin-1":      EncodingISO88591,
}

// DecodedText contains strict Unicode text and observations derived from the exact bytes.
type DecodedText struct {
	Text        string
	Observation model.EncodingObservation
}

// NormalizeEncoding resolves one documented encoding name or stable alias.
func NormalizeEncoding(name string) (string, bool) {
	canonical, exists := encodingAliases[normalizeToken(name)]
	return canonical, exists
}

// DecodeText decodes exact source bytes without guessing between legacy encodings.
func DecodeText(payload []byte, override string) (DecodedText, error) {
	bom, bomBytes := detectBOM(payload)
	selected := ""
	if override != "" {
		canonical, exists := NormalizeEncoding(override)
		if !exists {
			return DecodedText{}, fmt.Errorf("encoding %q is not supported", override)
		}
		selected = canonical
		if selected == EncodingUTF8BOM {
			if bom != EncodingUTF8 {
				return DecodedText{}, fmt.Errorf("encoding %q requires a UTF-8 BOM", override)
			}
			selected = EncodingUTF8
		} else if bom != "" && selected != bom {
			return DecodedText{}, fmt.Errorf("encoding %q conflicts with the detected %s BOM", override, bom)
		}
	} else if bom != "" {
		selected = bom
	} else if inferred := inferBOMlessUTF16(payload); inferred != "" {
		selected = inferred
	} else if utf8.Valid(payload) {
		selected = EncodingUTF8
	} else {
		return DecodedText{}, fmt.Errorf("text is neither strict UTF-8 nor strongly identified UTF-16; provide --encoding for legacy text")
	}

	content := payload[bomBytes:]
	var text string
	var err error
	switch selected {
	case EncodingUTF8:
		if !utf8.Valid(content) {
			return DecodedText{}, fmt.Errorf("decode %s: malformed UTF-8", selected)
		}
		text = string(content)
	case EncodingUTF16LE:
		text, err = decodeUTF16(content, binary.LittleEndian)
	case EncodingUTF16BE:
		text, err = decodeUTF16(content, binary.BigEndian)
	case EncodingWindows1252:
		text, err = decodeSingleByte(content, charmap.Windows1252)
	case EncodingISO88591:
		text, err = decodeSingleByte(content, charmap.ISO8859_1)
	default:
		return DecodedText{}, fmt.Errorf("internal unsupported encoding %q", selected)
	}
	if err != nil {
		return DecodedText{}, fmt.Errorf("decode %s: %w", selected, err)
	}

	confidence := 1.0
	observationEncoding := selected
	observation := model.EncodingObservation{LineEndings: classifyLineEndings(text), DetectedEncoding: &observationEncoding, Confidence: &confidence}
	if bom != "" {
		observedBOM := bom
		observation.BOM = &observedBOM
	}
	return DecodedText{Text: text, Observation: observation}, nil
}

func detectBOM(payload []byte) (string, int) {
	switch {
	case len(payload) >= 3 && payload[0] == 0xef && payload[1] == 0xbb && payload[2] == 0xbf:
		return EncodingUTF8, 3
	case len(payload) >= 2 && payload[0] == 0xff && payload[1] == 0xfe:
		return EncodingUTF16LE, 2
	case len(payload) >= 2 && payload[0] == 0xfe && payload[1] == 0xff:
		return EncodingUTF16BE, 2
	default:
		return "", 0
	}
}

func inferBOMlessUTF16(payload []byte) string {
	if len(payload) < 4 || len(payload)%2 != 0 {
		return ""
	}
	pairs := len(payload) / 2
	zeroEven := 0
	zeroOdd := 0
	for index := 0; index < len(payload); index += 2 {
		if payload[index] == 0 {
			zeroEven++
		}
		if payload[index+1] == 0 {
			zeroOdd++
		}
	}
	strong := func(primary, opposite int) bool {
		return primary*4 >= pairs*3 && opposite*8 <= pairs
	}
	if strong(zeroOdd, zeroEven) {
		return EncodingUTF16LE
	}
	if strong(zeroEven, zeroOdd) {
		return EncodingUTF16BE
	}
	return ""
}

func decodeUTF16(payload []byte, order binary.ByteOrder) (string, error) {
	if len(payload)%2 != 0 {
		return "", fmt.Errorf("odd trailing byte")
	}
	words := make([]uint16, len(payload)/2)
	for index := range words {
		words[index] = order.Uint16(payload[index*2 : index*2+2])
	}
	runes := make([]rune, 0, len(words))
	for index := 0; index < len(words); index++ {
		word := words[index]
		switch {
		case word >= 0xd800 && word <= 0xdbff:
			if index+1 >= len(words) || words[index+1] < 0xdc00 || words[index+1] > 0xdfff {
				return "", fmt.Errorf("unpaired high surrogate at code unit %d", index)
			}
			runes = append(runes, utf16.DecodeRune(rune(word), rune(words[index+1])))
			index++
		case word >= 0xdc00 && word <= 0xdfff:
			return "", fmt.Errorf("unpaired low surrogate at code unit %d", index)
		default:
			runes = append(runes, rune(word))
		}
	}
	return string(runes), nil
}

func decodeSingleByte(payload []byte, characterMap *charmap.Charmap) (string, error) {
	decoded, err := characterMap.NewDecoder().Bytes(payload)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func classifyLineEndings(text string) string {
	crlf, lf, cr := false, false, false
	for index := 0; index < len(text); index++ {
		switch text[index] {
		case '\r':
			if index+1 < len(text) && text[index+1] == '\n' {
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
	for _, present := range []bool{crlf, lf, cr} {
		if present {
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
	if lf {
		return "lf"
	}
	return "cr"
}
