package cli

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/version"
)

type inputKind string

const (
	inputKindCueJSON        inputKind = "cue_json"
	inputKindNativeSubtitle inputKind = "native_subtitle"
)

type selectionBasis string

const (
	selectionBasisCueJSON   selectionBasis = "cue_json"
	selectionBasisExplicit  selectionBasis = "explicit"
	selectionBasisContent   selectionBasis = "content"
	selectionBasisExtension selectionBasis = "extension"
)

type inputOptions struct {
	format                  string
	encoding                string
	disableSpeakerDetection bool
}

type validatedInput struct {
	document        model.Document
	inputKind       inputKind
	selectionBasis  selectionBasis
	contentFormat   *codec.Format
	extensionFormat *codec.Format
	diagnostics     []model.Diagnostic
}

type inputConstraintError struct {
	message string
}

func (err *inputConstraintError) Error() string {
	if err == nil {
		return "invalid input option combination"
	}
	return err.message
}

func loadValidatedInput(ctx context.Context, captured source.Captured, options inputOptions) (validatedInput, error) {
	if ctx == nil {
		return validatedInput{}, fmt.Errorf("validate input: nil context")
	}
	if err := ctx.Err(); err != nil {
		return validatedInput{}, fmt.Errorf("validate input canceled: %w", err)
	}
	format, known := normalizeInputFormat(options.format)
	if !known {
		return validatedInput{}, &codec.UnknownFormatError{Selector: options.format}
	}
	options.format = format

	if options.format == "cueson" {
		if options.encoding != "" {
			return validatedInput{}, &inputConstraintError{message: "--encoding cannot be used with Cue JSON input"}
		}
		return loadCueJSONInput(ctx, captured.Bytes, selectionBasisExplicit)
	}
	if options.format == string(codec.FormatWebVTT) {
		if err := validateWebVTTInputEncoding(options.encoding); err != nil {
			return validatedInput{}, err
		}
	}
	if options.format != string(codec.FormatAuto) {
		return loadNativeInput(ctx, captured, options)
	}

	document, cueJSONErr := schema.Decode(captured.Bytes)
	if cueJSONErr == nil {
		if options.encoding != "" {
			return validatedInput{}, &inputConstraintError{message: "--encoding cannot be used with automatically detected Cue JSON input"}
		}
		return finishCueJSONInput(ctx, document, selectionBasisCueJSON)
	}
	if cueJSONCandidate(captured.Bytes, captured.Asset.FileName) {
		return validatedInput{}, fmt.Errorf("decode Cue JSON input: %w", cueJSONErr)
	}
	return loadNativeInput(ctx, captured, options)
}

func loadCueJSONInput(ctx context.Context, payload []byte, basis selectionBasis) (validatedInput, error) {
	document, err := schema.Decode(payload)
	if err != nil {
		return validatedInput{}, fmt.Errorf("decode Cue JSON input: %w", err)
	}
	return finishCueJSONInput(ctx, document, basis)
}

func finishCueJSONInput(ctx context.Context, document model.Document, basis selectionBasis) (validatedInput, error) {
	if err := schema.CheckLockstep(version.String()); err != nil {
		return validatedInput{}, fmt.Errorf("schema version lockstep: %w", err)
	}
	if err := source.ValidateIntegrity(ctx, document); err != nil {
		return validatedInput{}, fmt.Errorf("validate Cue JSON source integrity: %w", err)
	}
	return validatedInput{
		document:       document,
		inputKind:      inputKindCueJSON,
		selectionBasis: basis,
		diagnostics:    cloneDiagnostics(document.Diagnostics),
	}, nil
}

func loadNativeInput(ctx context.Context, captured source.Captured, options inputOptions) (validatedInput, error) {
	registry, err := workflowRegistry(options.encoding)
	if err != nil {
		return validatedInput{}, err
	}
	selection, err := registry.Select(captured.Bytes, captured.Asset.FileName, options.format)
	if err != nil {
		return validatedInput{}, err
	}
	captured.Asset.MediaType = mediaTypeForFormat(selection.Format)
	decoder, err := registry.RequireDecoder(selection.Format)
	if err != nil {
		return validatedInput{}, err
	}
	document, err := decoder(ctx, captured, codec.DecodeOptions{Encoding: options.encoding, DisableSpeakerDetection: options.disableSpeakerDetection})
	if err != nil {
		return validatedInput{}, err
	}
	combined := make([]model.Diagnostic, 0, len(selection.Diagnostics)+len(document.Diagnostics))
	combined = append(combined, selection.Diagnostics...)
	document.Diagnostics = append(combined, document.Diagnostics...)
	updateDiagnosticStats(&document)
	if err := document.Validate(); err != nil {
		return validatedInput{}, fmt.Errorf("validate source model: %w", err)
	}
	if err := source.ValidateIntegrity(ctx, document); err != nil {
		return validatedInput{}, fmt.Errorf("validate generated source integrity: %w", err)
	}
	basis := selectionBasisExplicit
	if !selection.Explicit {
		basis = selectionBasisExtension
		if selection.ContentFormat != nil {
			basis = selectionBasisContent
		}
	}
	return validatedInput{
		document:        document,
		inputKind:       inputKindNativeSubtitle,
		selectionBasis:  basis,
		contentFormat:   cloneFormat(selection.ContentFormat),
		extensionFormat: cloneFormat(selection.ExtensionFormat),
		diagnostics:     cloneDiagnostics(document.Diagnostics),
	}, nil
}

func normalizeInputFormat(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return string(codec.FormatAuto), true
	case "cueson", "json", "cue-json":
		return "cueson", true
	case "srt", "subrip":
		return string(codec.FormatSubRip), true
	case "vtt", "webvtt":
		return string(codec.FormatWebVTT), true
	default:
		return "", false
	}
}

func validateWebVTTInputEncoding(encoding string) error {
	if encoding == "" {
		return nil
	}
	canonical, known := codec.NormalizeEncoding(encoding)
	if !known || canonical != codec.EncodingUTF8 && canonical != codec.EncodingUTF8BOM {
		return &inputConstraintError{message: "WebVTT requires UTF-8; --encoding is incompatible with --format vtt"}
	}
	return nil
}

func cueJSONCandidate(payload []byte, fileName string) bool {
	trimmed := bytes.TrimSpace(payload)
	trimmed = bytes.TrimPrefix(trimmed, []byte{0xef, 0xbb, 0xbf})
	trimmed = bytes.TrimSpace(trimmed)
	return (len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')) || strings.EqualFold(filepath.Ext(fileName), ".json")
}

func cloneFormat(format *codec.Format) *codec.Format {
	if format == nil {
		return nil
	}
	cloned := *format
	return &cloned
}

func cloneDiagnostics(diagnostics []model.Diagnostic) []model.Diagnostic {
	if diagnostics == nil {
		return nil
	}
	return append([]model.Diagnostic{}, diagnostics...)
}
