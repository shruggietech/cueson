package codec_test

import (
	"reflect"
	"testing"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
)

const maxFuzzInputBytes = 64 << 10

func FuzzDecodeText(f *testing.F) {
	f.Add([]byte("hello\n"), "")
	f.Add([]byte{0xff, 0xfe, 'h', 0}, "utf-16le")
	f.Add([]byte{0x93, 'h', 0x94}, "windows-1252")
	f.Fuzz(func(t *testing.T, payload []byte, encoding string) {
		if len(payload)+len(encoding) > maxFuzzInputBytes {
			return
		}
		first, firstErr := codec.DecodeText(payload, encoding)
		second, secondErr := codec.DecodeText(payload, encoding)
		if (firstErr == nil) != (secondErr == nil) || firstErr == nil && !reflect.DeepEqual(first, second) {
			t.Fatalf("DecodeText is nondeterministic: (%#v, %v), (%#v, %v)", first, firstErr, second, secondErr)
		}
	})
}

func FuzzFormatSelection(f *testing.F) {
	registry, err := codec.NewRegistry(
		codec.Registration{Format: codec.FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{"srt"}, Detect: func(data []byte) codec.Evidence {
			decoded, decodeErr := codec.DecodeText(data, "")
			if decodeErr != nil {
				return codec.Evidence{}
			}
			if _, parseErr := subrip.Parse(decoded.Text, subrip.Options{}); parseErr != nil {
				return codec.Evidence{}
			}
			return codec.Evidence{Matched: true, Confidence: 90, Reason: "SubRip timing grammar"}
		}},
		codec.Registration{Format: codec.FormatWebVTT, Aliases: []string{"vtt"}, Extensions: []string{"vtt"}, Detect: func(data []byte) codec.Evidence {
			if !webvtt.Detect(data) {
				return codec.Evidence{}
			}
			return codec.Evidence{Matched: true, Confidence: 100, Reason: "WebVTT signature"}
		}},
	)
	if err != nil {
		f.Fatal(err)
	}
	f.Add([]byte("1\n00:00:00,000 --> 00:00:01,000\ntext\n"), "captions.srt", "auto")
	f.Add([]byte("WEBVTT\n"), "captions.vtt", "")
	f.Fuzz(func(t *testing.T, payload []byte, name, requested string) {
		if len(payload)+len(name)+len(requested) > maxFuzzInputBytes {
			return
		}
		first, firstErr := registry.Select(payload, name, requested)
		second, secondErr := registry.Select(payload, name, requested)
		if (firstErr == nil) != (secondErr == nil) || firstErr == nil && !selectionEqual(first, second) {
			t.Fatalf("format selection is nondeterministic: (%#v, %v), (%#v, %v)", first, firstErr, second, secondErr)
		}
	})
}

func selectionEqual(left, right codec.Selection) bool {
	if left.Format != right.Format || left.Explicit != right.Explicit || !reflect.DeepEqual(left.Diagnostics, right.Diagnostics) {
		return false
	}
	return formatPointerEqual(left.ContentFormat, right.ContentFormat) && formatPointerEqual(left.ExtensionFormat, right.ExtensionFormat)
}

func formatPointerEqual(left, right *codec.Format) bool {
	return left == nil && right == nil || left != nil && right != nil && *left == *right
}
