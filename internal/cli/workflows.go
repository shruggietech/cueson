package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/version"
)

func runEncode(ctx context.Context, options encodeOptions, stdout io.Writer, stderr io.Writer, diagnostics diagnosticWriter) int {
	captured, err := source.CaptureContext(ctx, options.input, source.CaptureOptions{MediaType: stringPointer("application/x-subrip")})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, fmt.Sprintf("input %q does not exist", options.input))
			writeUsage(stderr, encodeHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}

	registry, err := workflowRegistry(options.encoding)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}
	selection, err := registry.Select(captured.Bytes, captured.Asset.FileName, options.format)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}
	decoder, err := registry.RequireDecoder(selection.Format)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}
	document, err := decoder(ctx, captured, codec.DecodeOptions{Encoding: options.encoding, DisableSpeakerDetection: options.noSpeakerDetection})
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}
	combined := make([]model.Diagnostic, 0, len(selection.Diagnostics)+len(document.Diagnostics))
	combined = append(combined, selection.Diagnostics...)
	document.Diagnostics = append(combined, document.Diagnostics...)
	updateDiagnosticStats(&document)
	if err := document.Validate(); err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode model: %v", err))
		return ExitRuntimeFailure
	}
	payload, err := marshalDocument(document, options.pretty)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode Cue JSON: %v", err))
		return ExitRuntimeFailure
	}
	if err := schema.Validate(payload); err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode schema: %v", err))
		return ExitRuntimeFailure
	}
	for _, observation := range document.Diagnostics {
		if observation.Severity == "warning" {
			diagnostics.write(diagnosticWarning, observation.Code+": "+observation.Message)
		}
	}
	if options.stdout || (options.outputSet && options.output == "-") {
		return writeStdout(stdout, diagnostics, payload)
	}
	destination := options.output
	if !options.outputSet {
		destination = options.input + ".cueson.json"
	}
	return handleWorkflowFileResult("encode", encodeHelp, writeSchemaFile(destination, payload, options.force), stderr, diagnostics)
}

func runRender(ctx context.Context, options renderOptions, stdout io.Writer, stderr io.Writer, diagnostics diagnosticWriter) int {
	if err := ctx.Err(); err != nil {
		diagnostics.write(diagnosticError, "operation canceled")
		return ExitRuntimeFailure
	}
	payload, err := os.ReadFile(options.input)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, fmt.Sprintf("input %q does not exist", options.input))
			writeUsage(stderr, renderHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("render: read Cue JSON: %v", err))
		return ExitRuntimeFailure
	}
	document, err := schema.Decode(payload)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("render: %v", err))
		return ExitRuntimeFailure
	}
	registry, err := workflowRegistry("")
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("render: %v", err))
		return ExitRuntimeFailure
	}
	registration, ok := registry.Lookup(options.target)
	if !ok {
		diagnostics.write(diagnosticError, fmt.Sprintf("render: format %q is unknown", options.target))
		return ExitRuntimeFailure
	}
	renderer, err := registry.RequireRenderer(registration.Format)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("render: %v", err))
		return ExitRuntimeFailure
	}
	result, err := renderer(ctx, document, codec.RenderOptions{Strict: options.strict})
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("render: %v", err))
		return ExitRuntimeFailure
	}
	for _, observation := range result.Diagnostics {
		if observation.Severity == "warning" {
			diagnostics.write(diagnosticWarning, observation.Code+": "+observation.Message)
		}
	}
	if !options.outputSet || options.output == "-" {
		return writeStdout(stdout, diagnostics, result.Bytes)
	}
	return handleWorkflowFileResult("render", renderHelp, writeSchemaFile(options.output, result.Bytes, options.force), stderr, diagnostics)
}

func workflowRegistry(encoding string) (*codec.Registry, error) {
	decodeSubRip := func(ctx context.Context, captured source.Captured, options codec.DecodeOptions) (model.Document, error) {
		if err := ctx.Err(); err != nil {
			return model.Document{}, err
		}
		decoded, err := codec.DecodeText(captured.Bytes, options.Encoding)
		if err != nil {
			return model.Document{}, err
		}
		captured.Asset.Encoding = &decoded.Observation
		parsed, err := subrip.Parse(decoded.Text, subrip.Options{DetectSpeakers: !options.DisableSpeakerDetection})
		if err != nil {
			return model.Document{}, err
		}
		return newSubRipDocument(captured.Asset, parsed), nil
	}
	renderSubRip := func(ctx context.Context, document model.Document, options codec.RenderOptions) (codec.RenderResult, error) {
		if err := ctx.Err(); err != nil {
			return codec.RenderResult{}, err
		}
		if document.Format != string(codec.FormatSubRip) {
			return codec.RenderResult{}, fmt.Errorf("document format %q cannot be rendered as SubRip", document.Format)
		}
		losses := subRipRenderLosses(document)
		if options.Strict && len(losses) > 0 {
			return codec.RenderResult{}, fmt.Errorf("strict SubRip render blocked %d non-representable structured field set(s); first: %s", len(losses), losses[0].Message)
		}
		bytes, err := subrip.Render(document.Cues)
		return codec.RenderResult{Bytes: bytes, Diagnostics: losses}, err
	}
	detectSubRip := func(data []byte) codec.Evidence {
		decoded, err := codec.DecodeText(data, encoding)
		if err != nil {
			return codec.Evidence{}
		}
		if _, err := subrip.Parse(decoded.Text, subrip.Options{}); err != nil {
			return codec.Evidence{}
		}
		return codec.Evidence{Matched: true, Confidence: 90, Reason: "SubRip timing grammar"}
	}
	detectWebVTT := func(data []byte) codec.Evidence {
		decoded, err := codec.DecodeText(data, encoding)
		if err != nil {
			return codec.Evidence{}
		}
		if strings.HasPrefix(strings.TrimPrefix(decoded.Text, "\ufeff"), "WEBVTT") {
			return codec.Evidence{Matched: true, Confidence: 100, Reason: "WebVTT signature"}
		}
		return codec.Evidence{}
	}
	return codec.NewRegistry(
		codec.Registration{Format: codec.FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{"srt"}, Detect: detectSubRip, Decode: decodeSubRip, Render: renderSubRip},
		codec.Registration{Format: codec.FormatWebVTT, Aliases: []string{"vtt"}, Extensions: []string{"vtt"}, Detect: detectWebVTT},
	)
}

func subRipRenderLosses(document model.Document) []model.Diagnostic {
	diagnostics := make([]model.Diagnostic, 0)
	if document.Metadata.Title != nil || document.Metadata.Language != nil || document.Metadata.Kind != nil || document.Metadata.Description != nil {
		diagnostics = append(diagnostics, model.Diagnostic{Severity: "warning", Code: "subrip_render_metadata_unrepresented", Message: "document metadata has no canonical SubRip representation"})
	}
	for index := range document.Cues {
		cue := &document.Cues[index]
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
		if len(fields) == 0 {
			continue
		}
		order, id := cue.SourceOrder, cue.ID
		diagnostics = append(diagnostics, model.Diagnostic{Severity: "warning", Code: "subrip_render_fields_unrepresented", Message: strings.Join(fields, ", ") + " are not represented by canonical SubRip output", SourceOrder: &order, CueID: &id})
	}
	return diagnostics
}

func newSubRipDocument(asset model.SourceAsset, parsed subrip.Result) model.Document {
	minimum := parsed.Cues[0].Timing.StartMilliseconds
	maximum := parsed.Cues[0].Timing.EndMilliseconds
	for index := range parsed.Cues {
		if parsed.Cues[index].Timing.StartMilliseconds < minimum {
			minimum = parsed.Cues[index].Timing.StartMilliseconds
		}
		if parsed.Cues[index].Timing.EndMilliseconds > maximum {
			maximum = parsed.Cues[index].Timing.EndMilliseconds
		}
	}
	span := maximum - minimum
	document := model.Document{
		Schema: schema.ID(), SchemaVersion: schema.Version(), Format: string(codec.FormatSubRip),
		FormatSupport: model.FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
		Producer:      model.Producer{Name: "cueson", Version: version.String()},
		Source:        model.SourceEnvelope{PrimaryAssetID: asset.ID, Assets: []model.SourceAsset{asset}},
		Metadata:      model.Metadata{},
		Document:      model.DocumentSummary{CueCount: len(parsed.Cues), MediaStartMilliseconds: int64Pointer(minimum), MediaEndMilliseconds: int64Pointer(maximum), MediaSpanMilliseconds: int64Pointer(span)},
		Cues:          parsed.Cues,
		FormatData:    model.DocumentFormatData{SubRip: &model.SubRipDocumentData{Dialect: "subrip"}},
		Diagnostics:   parsed.Diagnostics,
		Stats:         model.Stats{CueCount: len(parsed.Cues), MediaSpanMilliseconds: int64Pointer(span)},
	}
	updateDiagnosticStats(&document)
	return document
}

func updateDiagnosticStats(document *model.Document) {
	document.Stats.DiagnosticCount = len(document.Diagnostics)
	document.Stats.WarningCount = 0
	document.Stats.ErrorCount = 0
	for _, observation := range document.Diagnostics {
		switch observation.Severity {
		case "warning":
			document.Stats.WarningCount++
		case "error":
			document.Stats.ErrorCount++
		}
	}
}

func marshalDocument(document model.Document, pretty bool) ([]byte, error) {
	var payload []byte
	var err error
	if pretty {
		payload, err = json.MarshalIndent(document, "", "  ")
	} else {
		payload, err = json.Marshal(document)
	}
	if err != nil {
		return nil, err
	}
	return append(payload, '\n'), nil
}

func handleWorkflowFileResult(operation string, usage helpTarget, err error, stderr io.Writer, diagnostics diagnosticWriter) int {
	if err == nil {
		return ExitSuccess
	}
	var precondition *outputPreconditionError
	if errors.As(err, &precondition) {
		diagnostics.write(diagnosticError, precondition.Error())
		writeUsage(stderr, usage)
		return ExitInvocation
	}
	diagnostics.write(diagnosticError, fmt.Sprintf("%s: write output: %v", operation, err))
	return ExitRuntimeFailure
}

func stringPointer(value string) *string { return &value }
func int64Pointer(value int64) *int64    { return &value }
