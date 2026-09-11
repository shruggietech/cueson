package codec

import (
	"strings"
	"testing"
	"unicode/utf16"

	"golang.org/x/text/encoding/charmap"
)

func TestDecodeTextHonorsUnicodeBOMsAndStrictUTF8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		payload  []byte
		encoding string
		bom      string
		text     string
	}{
		{name: "utf8", payload: []byte("Hello\n"), encoding: "utf-8", text: "Hello\n"},
		{name: "utf8 bom", payload: append([]byte{0xef, 0xbb, 0xbf}, []byte("Hello\r\n")...), encoding: "utf-8", bom: "utf-8", text: "Hello\r\n"},
		{name: "utf16 little endian", payload: append([]byte{0xff, 0xfe}, encodeUTF16("Hello\r", true)...), encoding: "utf-16le", bom: "utf-16le", text: "Hello\r"},
		{name: "utf16 big endian", payload: append([]byte{0xfe, 0xff}, encodeUTF16("Hello", false)...), encoding: "utf-16be", bom: "utf-16be", text: "Hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			decoded, err := DecodeText(tt.payload, "")
			if err != nil {
				t.Fatal(err)
			}
			if decoded.Text != tt.text || stringValue(decoded.Observation.DetectedEncoding) != tt.encoding || stringValue(decoded.Observation.BOM) != tt.bom {
				t.Fatalf("DecodeText() = %#v", decoded)
			}
		})
	}
}

func TestDecodeTextRecognizesStrongBOMlessUTF16Evidence(t *testing.T) {
	t.Parallel()

	for _, littleEndian := range []bool{true, false} {
		payload := encodeUTF16("1\n00:00:00,000 --> 00:00:01,000\nHello\n", littleEndian)
		decoded, err := DecodeText(payload, "")
		if err != nil {
			t.Fatal(err)
		}
		want := "utf-16be"
		if littleEndian {
			want = "utf-16le"
		}
		if stringValue(decoded.Observation.DetectedEncoding) != want {
			t.Fatalf("detected encoding = %q, want %q", stringValue(decoded.Observation.DetectedEncoding), want)
		}
	}
}

func TestDecodeTextRequiresExplicitLegacyEncoding(t *testing.T) {
	t.Parallel()

	payload := []byte{0x93, 'H', 'i', 0x94}
	if _, err := DecodeText(payload, ""); err == nil || !strings.Contains(err.Error(), "--encoding") {
		t.Fatalf("DecodeText(auto) error = %v", err)
	}
	decoded, err := DecodeText(payload, "cp1252")
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Text != "“Hi”" || stringValue(decoded.Observation.DetectedEncoding) != "windows-1252" {
		t.Fatalf("DecodeText(cp1252) = %#v", decoded)
	}

	latinPayload, err := charmap.ISO8859_1.NewEncoder().Bytes([]byte("olá"))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err = DecodeText(latinPayload, "latin1")
	if err != nil || decoded.Text != "olá" || stringValue(decoded.Observation.DetectedEncoding) != "iso-8859-1" {
		t.Fatalf("DecodeText(latin1) = (%#v, %v)", decoded, err)
	}
}

func TestDecodeTextRejectsBOMConflictsAndMalformedUnicode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		payload  []byte
		override string
	}{
		{name: "bom conflict", payload: append([]byte{0xff, 0xfe}, encodeUTF16("Hello", true)...), override: "utf-8"},
		{name: "utf8 bom required", payload: []byte("Hello"), override: "utf-8-bom"},
		{name: "odd utf16", payload: []byte{0xff, 0xfe, 'H'}, override: ""},
		{name: "unpaired surrogate", payload: []byte{0xff, 0xfe, 0x00, 0xd8, 'A', 0x00}, override: ""},
		{name: "invalid utf8", payload: []byte{0xc3, 0x28}, override: "utf-8"},
		{name: "unknown override", payload: []byte("Hello"), override: "guess-me"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := DecodeText(tt.payload, tt.override); err == nil {
				t.Fatal("DecodeText() unexpectedly succeeded")
			}
		})
	}
}

func TestDecodeTextClassifiesLineEndings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		text string
		want string
	}{
		{text: "none", want: "none"},
		{text: "one\n", want: "lf"},
		{text: "one\r\n", want: "crlf"},
		{text: "one\r", want: "cr"},
		{text: "one\r\ntwo\nthree\r", want: "mixed"},
	}
	for _, tt := range tests {
		decoded, err := DecodeText([]byte(tt.text), "")
		if err != nil {
			t.Fatal(err)
		}
		if decoded.Observation.LineEndings != tt.want {
			t.Fatalf("line endings for %q = %q, want %q", tt.text, decoded.Observation.LineEndings, tt.want)
		}
	}
}

func encodeUTF16(text string, littleEndian bool) []byte {
	words := utf16.Encode([]rune(text))
	payload := make([]byte, 0, len(words)*2)
	for _, word := range words {
		first, second := byte(word), byte(word>>8)
		if !littleEndian {
			first, second = second, first
		}
		payload = append(payload, first, second)
	}
	return payload
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func TestDecodeTextDoesNotAliasInput(t *testing.T) {
	t.Parallel()

	payload := []byte("hello")
	decoded, err := DecodeText(payload, "")
	if err != nil {
		t.Fatal(err)
	}
	payload[0] = 'x'
	if decoded.Text != "hello" {
		t.Fatalf("decoded text changed to %q", decoded.Text)
	}
}
