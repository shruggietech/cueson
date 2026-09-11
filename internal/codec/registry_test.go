package codec

import (
	"context"
	"errors"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/source"
)

func TestRegistryNormalizesAliasesAndKeepsCapabilitiesIndependent(t *testing.T) {
	t.Parallel()

	decode := func(context.Context, source.Captured, DecodeOptions) (model.Document, error) {
		return model.Document{}, nil
	}
	registry, err := NewRegistry(
		Registration{Format: FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{".srt"}, Decode: decode},
		Registration{Format: FormatWebVTT, Aliases: []string{"vtt"}, Extensions: []string{"vtt"}},
	)
	if err != nil {
		t.Fatal(err)
	}

	registration, ok := registry.Lookup("SRT")
	if !ok || registration.Format != FormatSubRip || registration.Decode == nil || registration.Render != nil {
		t.Fatalf("Lookup(srt) = (%#v, %t)", registration, ok)
	}
	registration, ok = registry.Lookup("webvtt")
	if !ok || registration.Format != FormatWebVTT || registration.Decode != nil || registration.Render != nil {
		t.Fatalf("Lookup(webvtt) = (%#v, %t)", registration, ok)
	}
	if _, err := registry.RequireDecoder(FormatWebVTT); err == nil {
		t.Fatal("RequireDecoder(webvtt) accepted a registration without a decoder")
	} else {
		var missing *MissingCapabilityError
		if !errors.As(err, &missing) || missing.Capability != CapabilityDecode {
			t.Fatalf("RequireDecoder(webvtt) error = %v", err)
		}
	}
}

func TestRegistryRejectsConflictingIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		registrations []Registration
	}{
		{name: "duplicate canonical", registrations: []Registration{{Format: FormatSubRip}, {Format: FormatSubRip}}},
		{name: "alias collision", registrations: []Registration{{Format: FormatSubRip, Aliases: []string{"caption"}}, {Format: FormatWebVTT, Aliases: []string{"caption"}}}},
		{name: "extension collision", registrations: []Registration{{Format: FormatSubRip, Extensions: []string{"sub"}}, {Format: FormatWebVTT, Extensions: []string{".SUB"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewRegistry(tt.registrations...); err == nil {
				t.Fatal("NewRegistry() accepted conflicting identity")
			}
		})
	}
}

func TestSelectPrefersContentAndDiagnosesExtensionDisagreement(t *testing.T) {
	t.Parallel()

	registry := detectionRegistry(t)
	selection, err := registry.Select([]byte("1\n00:00:00,000 --> 00:00:01,000\nHello\n"), "captions.vtt", "auto")
	if err != nil {
		t.Fatal(err)
	}
	if selection.Format != FormatSubRip || len(selection.Diagnostics) != 1 || selection.Diagnostics[0].Code != DiagnosticExtensionDisagreement {
		t.Fatalf("Select() = %#v", selection)
	}
	if selection.ContentFormat == nil || *selection.ContentFormat != FormatSubRip || selection.ExtensionFormat == nil || *selection.ExtensionFormat != FormatWebVTT {
		t.Fatalf("Select() evidence = %#v", selection)
	}
}

func TestSelectExplicitFormatWins(t *testing.T) {
	t.Parallel()

	registry := detectionRegistry(t)
	selection, err := registry.Select([]byte("WEBVTT\n"), "captions.vtt", "srt")
	if err != nil {
		t.Fatal(err)
	}
	if selection.Format != FormatSubRip || !selection.Explicit || len(selection.Diagnostics) != 0 {
		t.Fatalf("Select(explicit) = %#v", selection)
	}
}

func TestSelectUsesExtensionOnlyWhenContentHasNoWinner(t *testing.T) {
	t.Parallel()

	registry := detectionRegistry(t)
	selection, err := registry.Select([]byte("not enough evidence"), "captions.srt", "")
	if err != nil || selection.Format != FormatSubRip || selection.ContentFormat != nil {
		t.Fatalf("Select(extension) = (%#v, %v)", selection, err)
	}
}

func TestSelectRejectsAmbiguousOrUnknownInput(t *testing.T) {
	t.Parallel()

	registry, err := NewRegistry(
		Registration{Format: FormatSubRip, Aliases: []string{"srt"}, Detect: func([]byte) Evidence { return Evidence{Matched: true, Confidence: 50} }},
		Registration{Format: FormatWebVTT, Aliases: []string{"vtt"}, Detect: func([]byte) Evidence { return Evidence{Matched: true, Confidence: 50} }},
	)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Select([]byte("ambiguous"), "captions.bin", "auto"); err == nil {
		t.Fatal("Select() accepted tied content evidence")
	} else {
		var ambiguous *AmbiguousFormatError
		if !errors.As(err, &ambiguous) {
			t.Fatalf("Select() error = %v", err)
		}
	}
	if _, err := detectionRegistry(t).Select([]byte("unknown"), "captions.bin", "auto"); err == nil {
		t.Fatal("Select() accepted unknown input")
	} else {
		var unknown *UnknownFormatError
		if !errors.As(err, &unknown) {
			t.Fatalf("Select() error = %v", err)
		}
	}
}

func detectionRegistry(t *testing.T) *Registry {
	t.Helper()
	registry, err := NewRegistry(
		Registration{Format: FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{"srt"}, Detect: func(data []byte) Evidence {
			return Evidence{Matched: len(data) > 1 && data[0] == '1', Confidence: 90, Reason: "subrip timing grammar"}
		}},
		Registration{Format: FormatWebVTT, Aliases: []string{"vtt"}, Extensions: []string{"vtt"}, Detect: func(data []byte) Evidence {
			return Evidence{Matched: len(data) >= 6 && string(data[:6]) == "WEBVTT", Confidence: 100, Reason: "WebVTT signature"}
		}},
	)
	if err != nil {
		t.Fatal(err)
	}
	return registry
}
