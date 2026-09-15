package codec

import (
	"errors"
	"testing"
)

func TestGuardedContentCannotEscapeViaExtensionOrExplicitSelector(t *testing.T) {
	rejected := errors.New("unsupported_scripted_dialect")
	registry, err := NewRegistry(
		Registration{Format: FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{"srt"}},
		Registration{Format: FormatASS, EnforceContentSelection: true, Detect: func(data []byte) Evidence {
			if string(data) == "bad" {
				return Evidence{Rejection: rejected}
			}
			if string(data) == "ass" {
				return Evidence{Matched: true, Confidence: 100}
			}
			return Evidence{}
		}},
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, selector := range []string{"auto", "srt", "ass"} {
		if _, err := registry.Select([]byte("bad"), "captions.srt", selector); !errors.Is(err, rejected) {
			t.Errorf("rejected candidate escaped via %q: %v", selector, err)
		}
	}
	if _, err := registry.Select([]byte("ass"), "captions.srt", "srt"); err == nil {
		t.Error("guarded ASS content escaped via explicit SubRip selector")
	}
	selected, err := registry.Select([]byte("ass"), "captions.srt", "ass")
	if err != nil || selected.Format != FormatASS || !selected.Explicit {
		t.Fatalf("matching explicit = %#v, %v", selected, err)
	}
	if _, err := registry.Select([]byte("ordinary srt"), "captions.srt", "srt"); err != nil {
		t.Fatalf("ordinary explicit format changed: %v", err)
	}
}
