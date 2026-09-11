package codec

import "testing"

func FuzzDecodeText(f *testing.F) {
	f.Add([]byte("hello\n"), "")
	f.Add([]byte{0xff, 0xfe, 'h', 0}, "utf-16le")
	f.Add([]byte{0x93, 'h', 0x94}, "windows-1252")
	f.Fuzz(func(t *testing.T, payload []byte, encoding string) {
		_, _ = DecodeText(payload, encoding)
	})
}

func FuzzFormatSelection(f *testing.F) {
	registry, err := NewRegistry(
		Registration{Format: FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{"srt"}, Detect: func(data []byte) Evidence {
			return Evidence{Matched: len(data) >= 3 && string(data[:3]) == "1\n0", Confidence: 90}
		}},
		Registration{Format: FormatWebVTT, Aliases: []string{"vtt"}, Extensions: []string{"vtt"}, Detect: func(data []byte) Evidence {
			return Evidence{Matched: len(data) >= 6 && string(data[:6]) == "WEBVTT", Confidence: 100}
		}},
	)
	if err != nil {
		f.Fatal(err)
	}
	f.Add([]byte("1\n00:00:00,000 --> 00:00:01,000\ntext\n"), "captions.srt", "auto")
	f.Add([]byte("WEBVTT\n"), "captions.vtt", "")
	f.Fuzz(func(t *testing.T, payload []byte, name, requested string) {
		_, _ = registry.Select(payload, name, requested)
	})
}
