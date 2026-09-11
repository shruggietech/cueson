package cli

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
)

func TestLoadValidatedInputClassifiesCueJSONBeforeNative(t *testing.T) {
	t.Parallel()

	captured := captureValidatedInputFixture(t, "misleading.srt", schema.Representative())
	loaded, err := loadValidatedInput(context.Background(), captured, inputOptions{format: "auto"})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.inputKind != inputKindCueJSON || loaded.selectionBasis != selectionBasisCueJSON || loaded.document.Format != "subrip" {
		t.Fatalf("loaded input = %#v", loaded)
	}
	if loaded.contentFormat != nil || loaded.extensionFormat != nil || !reflect.DeepEqual(loaded.diagnostics, loaded.document.Diagnostics) {
		t.Fatalf("Cue JSON evidence or diagnostics = %#v", loaded)
	}
}

func TestLoadValidatedInputPreservesNativeSelectionEvidenceAndDiagnostics(t *testing.T) {
	t.Parallel()

	captured := captureValidatedInputFixture(t, "misleading.vtt", []byte("00:00:00,000 --> 00:00:01,000\nAlice: Hello\n"))
	loaded, err := loadValidatedInput(context.Background(), captured, inputOptions{format: "auto", disableSpeakerDetection: true})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.inputKind != inputKindNativeSubtitle || loaded.selectionBasis != selectionBasisContent || loaded.document.Format != string(codec.FormatSubRip) {
		t.Fatalf("loaded input = %#v", loaded)
	}
	if loaded.contentFormat == nil || *loaded.contentFormat != codec.FormatSubRip || loaded.extensionFormat == nil || *loaded.extensionFormat != codec.FormatWebVTT {
		t.Fatalf("native evidence = content %#v, extension %#v", loaded.contentFormat, loaded.extensionFormat)
	}
	if len(loaded.diagnostics) == 0 || loaded.diagnostics[0].Code != codec.DiagnosticExtensionDisagreement {
		t.Fatalf("native diagnostics = %#v", loaded.diagnostics)
	}
	if len(loaded.document.Cues[0].Speakers) != 0 {
		t.Fatalf("validation derived speakers = %#v", loaded.document.Cues[0].Speakers)
	}
	if err := source.ValidateIntegrity(context.Background(), loaded.document); err != nil {
		t.Fatalf("generated source envelope integrity: %v", err)
	}
}

func TestLoadValidatedInputHonorsExplicitSelectionAndLegacyEncoding(t *testing.T) {
	t.Parallel()

	captured := captureValidatedInputFixture(t, "captions.bin", []byte("1\n00:00:00,000 --> 00:00:01,000\n\x93Hello\x94\n"))
	loaded, err := loadValidatedInput(context.Background(), captured, inputOptions{format: "srt", encoding: "cp1252", disableSpeakerDetection: true})
	if err != nil {
		t.Fatal(err)
	}
	if loaded.selectionBasis != selectionBasisExplicit || loaded.document.Format != string(codec.FormatSubRip) || loaded.document.Source.Assets[0].Encoding == nil || inputStringValue(loaded.document.Source.Assets[0].Encoding.DetectedEncoding) != codec.EncodingWindows1252 {
		t.Fatalf("explicit legacy input = %#v", loaded)
	}
}

func TestLoadValidatedInputRejectsCueJSONCandidatesWithoutFallback(t *testing.T) {
	t.Parallel()

	for _, fixture := range []struct {
		name    string
		payload []byte
	}{
		{name: "broken.srt", payload: []byte("  {not JSON or subtitles}")},
		{name: "broken.json", payload: []byte("1\n00:00:00,000 --> 00:00:01,000\nHello\n")},
		{name: "bom.srt", payload: append([]byte{0xef, 0xbb, 0xbf}, []byte(" [broken")...)},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			captured := captureValidatedInputFixture(t, fixture.name, fixture.payload)
			if _, err := loadValidatedInput(context.Background(), captured, inputOptions{format: "auto"}); err == nil || !strings.Contains(err.Error(), "Cue JSON") {
				t.Fatalf("loadValidatedInput() error = %v", err)
			}
		})
	}
}

func TestLoadValidatedInputRejectsCorruptIntegrityAndOptionConflicts(t *testing.T) {
	t.Parallel()

	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	document.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64)
	payload, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	corrupt := captureValidatedInputFixture(t, "corrupt.json", payload)
	if _, err := loadValidatedInput(context.Background(), corrupt, inputOptions{format: "cueson"}); err == nil || !strings.Contains(err.Error(), "SHA-256") {
		t.Fatalf("corrupt integrity error = %v", err)
	}

	validJSON := captureValidatedInputFixture(t, "document.json", schema.Representative())
	for _, options := range []inputOptions{{format: "cueson", encoding: "utf-8"}, {format: "auto", encoding: "utf-8"}} {
		_, err := loadValidatedInput(context.Background(), validJSON, options)
		var conflict *inputConstraintError
		if !errors.As(err, &conflict) || !strings.Contains(err.Error(), "Cue JSON") {
			t.Fatalf("Cue JSON encoding conflict = %T %v", err, err)
		}
	}

	webVTT := captureValidatedInputFixture(t, "captions.vtt", []byte("WEBVTT\n\n00:00.000 --> 00:01.000\nHello\n"))
	if _, err := loadValidatedInput(context.Background(), webVTT, inputOptions{format: "vtt", encoding: "windows-1252"}); err == nil || !strings.Contains(err.Error(), "WebVTT requires UTF-8") {
		t.Fatalf("WebVTT encoding conflict = %v", err)
	}
}

func TestLoadValidatedInputHonorsCancellation(t *testing.T) {
	t.Parallel()

	captured := captureValidatedInputFixture(t, "document.json", schema.Representative())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := loadValidatedInput(ctx, captured, inputOptions{format: "auto"}); !errors.Is(err, context.Canceled) {
		t.Fatalf("loadValidatedInput() error = %v, want context cancellation", err)
	}
	var nilContext context.Context
	if _, err := loadValidatedInput(nilContext, captured, inputOptions{format: "auto"}); err == nil || !strings.Contains(err.Error(), "nil context") {
		t.Fatalf("nil context error = %v", err)
	}
}

func captureValidatedInputFixture(t *testing.T, name string, payload []byte) source.Captured {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	captured, err := source.CaptureContext(context.Background(), path, source.CaptureOptions{})
	if err != nil {
		t.Fatal(err)
	}
	return captured
}

func inputStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
