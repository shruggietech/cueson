package convert

import (
	"context"
	"fmt"
	"strings"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
)

type targetRenderer func(context.Context, model.Document, string) ([]byte, []model.Diagnostic, error)

func renderProjected(ctx context.Context, document model.Document, targetFormat string) ([]byte, []model.Diagnostic, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	var bytes []byte
	var diagnostics []model.Diagnostic
	var err error
	switch targetFormat {
	case "subrip":
		bytes, err = subrip.Render(document.Cues)
	case "webvtt":
		var rendered webvtt.RenderResult
		rendered, err = webvtt.Render(document, webvtt.RenderOptions{Strict: true})
		bytes = rendered.Bytes
		for _, diagnostic := range rendered.Diagnostics {
			diagnostics = append(diagnostics, model.Diagnostic(diagnostic))
		}
	default:
		return nil, nil, fmt.Errorf("unsupported target format %q", targetFormat)
	}
	if err != nil {
		return nil, nil, err
	}
	if len(bytes) >= 3 && bytes[0] == 0xef && bytes[1] == 0xbb && bytes[2] == 0xbf || strings.ContainsRune(string(bytes), '\r') || !strings.HasSuffix(string(bytes), "\n") || strings.HasSuffix(string(bytes), "\n\n") {
		return nil, nil, fmt.Errorf("target renderer returned noncanonical text")
	}
	if err := validateTargetParser(bytes, targetFormat); err != nil {
		return nil, nil, fmt.Errorf("target parser rejected rendered output: %w", err)
	}
	if err := validateTargetSemantics(document, bytes, targetFormat); err != nil {
		return nil, nil, fmt.Errorf("target parser changed projected semantics: %w", err)
	}
	return bytes, diagnostics, nil
}

func validateTargetParser(bytes []byte, targetFormat string) error {
	switch targetFormat {
	case "subrip":
		_, err := subrip.Parse(string(bytes), subrip.Options{})
		return err
	case "webvtt":
		_, err := webvtt.Parse(string(bytes))
		return err
	default:
		return fmt.Errorf("unsupported target format %q", targetFormat)
	}
}

func validateTargetSemantics(document model.Document, bytes []byte, targetFormat string) error {
	var parsedCues []model.Cue
	switch targetFormat {
	case "subrip":
		parsed, err := subrip.Parse(string(bytes), subrip.Options{})
		if err != nil {
			return err
		}
		parsedCues = parsed.Cues
	case "webvtt":
		parsed, err := webvtt.Parse(string(bytes))
		if err != nil {
			return err
		}
		parsedCues = parsed.Cues
	default:
		return fmt.Errorf("unsupported target format %q", targetFormat)
	}
	if len(parsedCues) != len(document.Cues) {
		return fmt.Errorf("cue count changed from %d to %d", len(document.Cues), len(parsedCues))
	}
	for index := range parsedCues {
		want, got := document.Cues[index], parsedCues[index]
		if got.Timing != want.Timing {
			return fmt.Errorf("cue %d timing changed", index)
		}
		if got.Payload.PlainText != want.Payload.PlainText {
			return fmt.Errorf("cue %d readable payload changed", index)
		}
	}
	return nil
}

// AnalyzeSubRipRepresentability returns the reusable atomic target constraints
// shared by same-format rendering and cross-format conversion.
func AnalyzeSubRipRepresentability(document model.Document) (Report, error) {
	losses := make([]Loss, 0)
	appendMetadataLosses(&losses, document, "webvtt")
	for index := range document.Cues {
		cue := &document.Cues[index]
		if subrip.PayloadHasAmbiguousBoundary(cue.Payload.RawText) {
			losses = append(losses, cueLoss(document, "webvtt", *cue, LossCodeSubRipPayloadAmbiguous, KindAmbiguous, "payload has an ambiguous SubRip cue boundary", fmt.Sprintf("/cues/%d/payload/raw_text", index), nil))
		}
		appendCommonCueLosses(&losses, document, "webvtt", index, cue, false)
	}
	return NewReport(document, losses)
}

// SubRipRenderDiagnostics preserves the established public render warning
// vocabulary while keeping the underlying representability policy outside the
// CLI. Conversion itself uses the atomic Loss report above.
func SubRipRenderDiagnostics(document model.Document) []model.Diagnostic {
	diagnostics := make([]model.Diagnostic, 0)
	if document.Metadata.Title != nil || document.Metadata.Language != nil || document.Metadata.Kind != nil || document.Metadata.Description != nil {
		diagnostics = append(diagnostics, model.Diagnostic{Severity: "warning", Code: "subrip_render_metadata_unrepresented", Message: "document metadata has no canonical SubRip representation"})
	}
	for index := range document.Cues {
		cue := &document.Cues[index]
		order, id := cue.SourceOrder, cue.ID
		if subrip.PayloadHasAmbiguousBoundary(cue.Payload.RawText) {
			diagnostics = append(diagnostics, model.Diagnostic{Severity: "warning", Code: "subrip_render_payload_ambiguous", Message: "payload raw_text contains content that canonical SubRip reparses as a cue boundary", SourceOrder: &order, CueID: &id})
		}
		fields := make([]string, 0, 4)
		if len(cue.Speakers) > 0 {
			fields = append(fields, "speaker observations")
		}
		if len(cue.Tokens) > 0 {
			fields = append(fields, "token timing")
		}
		if len(cue.OCRObservations) > 0 {
			fields = append(fields, "OCR observations")
		}
		if cue.Placement != nil {
			fields = append(fields, "common placement")
		}
		if len(fields) > 0 {
			diagnostics = append(diagnostics, model.Diagnostic{Severity: "warning", Code: "subrip_render_fields_unrepresented", Message: strings.Join(fields, ", ") + " are not represented by canonical SubRip output", SourceOrder: &order, CueID: &id})
		}
	}
	return diagnostics
}

func subRipPayloadAmbiguous(payload string) bool {
	return subrip.PayloadHasAmbiguousBoundary(payload)
}

func subRipPlainText(payload string) string {
	return subrip.PlainText(payload)
}

func splitOnLF(value string) []string {
	return strings.Split(value, "\n")
}
