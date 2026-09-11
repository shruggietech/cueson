package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/convert"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/version"
)

func runEncode(ctx context.Context, options encodeOptions, stdout io.Writer, stderr io.Writer, diagnostics diagnosticWriter) int {
	captured, err := source.CaptureContext(ctx, options.input, source.CaptureOptions{})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, fmt.Sprintf("input %q does not exist", options.input))
			writeUsage(stderr, encodeHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}

	document, err := decodeCapturedSource(ctx, captured, options.format, options.encoding, options.noSpeakerDetection)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("encode: %v", err))
		return ExitRuntimeFailure
	}
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
	writeModelWarnings(diagnostics, document.Diagnostics)
	if options.stdout || (options.outputSet && options.output == "-") {
		return writeStdout(stdout, diagnostics, payload)
	}
	destination := options.output
	if !options.outputSet {
		destination = options.input + ".cueson.json"
	}
	return handleWorkflowFileResult("encode", encodeHelp, writeSchemaFile(destination, payload, options.force), stderr, diagnostics)
}

func decodeCapturedSource(ctx context.Context, captured source.Captured, requested, encoding string, disableSpeakerDetection bool) (model.Document, error) {
	loaded, err := loadNativeInput(ctx, captured, inputOptions{format: requested, encoding: encoding, disableSpeakerDetection: disableSpeakerDetection})
	if err != nil {
		return model.Document{}, err
	}
	return loaded.document, nil
}

func loadConversionDocument(ctx context.Context, captured source.Captured, options convertOptions) (model.Document, error) {
	loaded, err := loadValidatedInput(ctx, captured, inputOptions{format: options.from, encoding: options.encoding, disableSpeakerDetection: options.noSpeakerDetection})
	if err != nil {
		var constraint *inputConstraintError
		if errors.As(err, &constraint) {
			return model.Document{}, convertInvocationError(constraint.Error())
		}
		return model.Document{}, err
	}
	return loaded.document, nil
}

func runConvert(ctx context.Context, options convertOptions, stdout io.Writer, stderr io.Writer, diagnostics diagnosticWriter) int {
	if err := inspectConversionOutput(options); err != nil {
		var precondition *outputPreconditionError
		if errors.As(err, &precondition) {
			diagnostics.write(diagnosticError, precondition.Error())
			writeUsage(stderr, convertHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("convert: inspect output: %v", err))
		return ExitRuntimeFailure
	}

	captured, err := source.CaptureContext(ctx, options.input, source.CaptureOptions{})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, fmt.Sprintf("input %q does not exist", options.input))
			writeUsage(stderr, convertHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("convert: capture input: %v", err))
		return ExitRuntimeFailure
	}
	document, err := loadConversionDocument(ctx, captured, options)
	if err != nil {
		var invocation *invocationError
		if errors.As(err, &invocation) {
			diagnostics.write(diagnosticError, invocation.Error())
			writeUsage(stderr, invocation.usage)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("convert: %v", err))
		return ExitRuntimeFailure
	}

	writeModelWarnings(diagnostics, document.Diagnostics)
	result, err := convert.Convert(ctx, document, options.target, convert.Options{Strict: options.strict})
	writeConversionWarnings(diagnostics, result.LossReport)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("convert: %v", err))
		return ExitRuntimeFailure
	}
	writeModelWarnings(diagnostics, result.Diagnostics)

	if !options.outputSet || options.output == "-" {
		return writeStdout(stdout, diagnostics, result.Bytes)
	}
	return handleWorkflowFileResult("convert", convertHelp, writeSchemaFile(options.output, result.Bytes, options.force), stderr, diagnostics)
}

func inspectConversionOutput(options convertOptions) error {
	if !options.outputSet || options.output == "-" {
		return nil
	}
	parent := filepath.Dir(options.output)
	directory, err := os.Stat(parent)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return &outputPreconditionError{message: fmt.Sprintf("output directory %q does not exist", parent)}
	case err != nil:
		return err
	case !directory.IsDir():
		return &outputPreconditionError{message: fmt.Sprintf("output directory %q is not a directory", parent)}
	}
	_, err = inspectOutputDestination(options.output, options.force)
	return err
}

func writeModelWarnings(diagnostics diagnosticWriter, observations []model.Diagnostic) {
	for _, observation := range observations {
		if observation.Severity == "warning" {
			diagnostics.write(diagnosticWarning, observation.Code+": "+observation.Message)
		}
	}
}

func writeConversionWarnings(diagnostics diagnosticWriter, report convert.Report) {
	for _, loss := range report.Losses {
		diagnostics.write(diagnosticWarning, loss.Code+": "+loss.Message)
	}
}

func runRender(ctx context.Context, options renderOptions, stdout io.Writer, stderr io.Writer, diagnostics diagnosticWriter) int {
	if err := ctx.Err(); err != nil {
		diagnostics.write(diagnosticError, "operation canceled")
		return ExitRuntimeFailure
	}
	payload, err := source.ReadCueJSONContext(ctx, options.input)
	if err != nil {
		if source.IsCapturePrecondition(err) {
			diagnostics.write(diagnosticError, "render: input does not satisfy path or size requirements")
			writeUsage(stderr, renderHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, "render: Cue JSON acquisition failed")
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
	writeModelWarnings(diagnostics, result.Diagnostics)
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
		observations, err := convert.SubRipRenderDiagnostics(document)
		if err != nil {
			return codec.RenderResult{}, err
		}
		if options.Strict && len(observations) > 0 {
			return codec.RenderResult{}, fmt.Errorf("strict SubRip render blocked %d non-representable structured field set(s); first: %s", len(observations), observations[0].Message)
		}
		bytes, renderErr := subrip.Render(document.Cues)
		return codec.RenderResult{Bytes: bytes, Diagnostics: observations}, renderErr
	}
	decodeWebVTT := func(ctx context.Context, captured source.Captured, options codec.DecodeOptions) (model.Document, error) {
		if err := ctx.Err(); err != nil {
			return model.Document{}, err
		}
		decoded, err := webvtt.DecodeUTF8(captured.Bytes, options.Encoding)
		if err != nil {
			return model.Document{}, err
		}
		captured.Asset.Encoding = &decoded.Observation
		parsed, err := webvtt.Parse(decoded.Text)
		if err != nil {
			return model.Document{}, err
		}
		if options.DisableSpeakerDetection {
			for index := range parsed.Cues {
				parsed.Cues[index].Speakers = []model.Speaker{}
			}
		}
		return newWebVTTDocument(captured.Asset, parsed), nil
	}
	renderWebVTT := func(ctx context.Context, document model.Document, options codec.RenderOptions) (codec.RenderResult, error) {
		if err := ctx.Err(); err != nil {
			return codec.RenderResult{}, err
		}
		result, err := webvtt.Render(document, webvtt.RenderOptions{Strict: options.Strict})
		return codec.RenderResult{Bytes: result.Bytes, Diagnostics: result.Diagnostics}, err
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
		if webvtt.Detect(data) {
			return codec.Evidence{Matched: true, Confidence: 100, Reason: "WebVTT signature"}
		}
		return codec.Evidence{}
	}
	return codec.NewRegistry(
		codec.Registration{Format: codec.FormatSubRip, Aliases: []string{"srt"}, Extensions: []string{"srt"}, Detect: detectSubRip, Decode: decodeSubRip, Render: renderSubRip},
		codec.Registration{Format: codec.FormatWebVTT, Aliases: []string{"vtt"}, Extensions: []string{"vtt"}, Detect: detectWebVTT, Decode: decodeWebVTT, Render: renderWebVTT},
	)
}

func mediaTypeForFormat(format codec.Format) *string {
	mediaType := ""
	switch format {
	case codec.FormatSubRip:
		mediaType = "application/x-subrip"
	case codec.FormatWebVTT:
		mediaType = "text/vtt"
	default:
		return nil
	}
	return &mediaType
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

func newWebVTTDocument(asset model.SourceAsset, parsed webvtt.Result) model.Document {
	minimum := parsed.Cues[0].Timing.StartMilliseconds
	maximum := parsed.Cues[0].Timing.EndMilliseconds
	hasTokens := false
	for index := range parsed.Cues {
		if parsed.Cues[index].Timing.StartMilliseconds < minimum {
			minimum = parsed.Cues[index].Timing.StartMilliseconds
		}
		if parsed.Cues[index].Timing.EndMilliseconds > maximum {
			maximum = parsed.Cues[index].Timing.EndMilliseconds
		}
		hasTokens = hasTokens || len(parsed.Cues[index].Tokens) > 0
	}
	span := maximum - minimum
	document := model.Document{
		Schema: schema.ID(), SchemaVersion: schema.Version(), Format: string(codec.FormatWebVTT),
		FormatSupport: model.FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
		Producer:      model.Producer{Name: "cueson", Version: version.String()},
		Source:        model.SourceEnvelope{PrimaryAssetID: asset.ID, Assets: []model.SourceAsset{asset}},
		Metadata:      model.Metadata{},
		Document:      model.DocumentSummary{CueCount: len(parsed.Cues), MediaStartMilliseconds: int64Pointer(minimum), MediaEndMilliseconds: int64Pointer(maximum), MediaSpanMilliseconds: int64Pointer(span), HasWordLevelTiming: hasTokens},
		Cues:          parsed.Cues,
		FormatData:    model.DocumentFormatData{WebVTT: &parsed.DocumentData},
		Diagnostics:   parsed.Diagnostics,
		Stats:         model.Stats{CueCount: len(parsed.Cues), HasWordLevelTiming: hasTokens, MediaSpanMilliseconds: int64Pointer(span)},
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

func int64Pointer(value int64) *int64 { return &value }
