package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/source"
)

type validateOptions struct {
	input       string
	format      string
	encoding    string
	formatSet   bool
	encodingSet bool
}

func setValidateValue(options *validateOptions, option, value string) error {
	if options == nil {
		return fmt.Errorf("validate options are unavailable")
	}
	switch option {
	case "--format":
		if options.formatSet {
			return fmt.Errorf("--format may be specified only once")
		}
		options.format, options.formatSet = value, true
	case "--encoding":
		if options.encodingSet {
			return fmt.Errorf("--encoding may be specified only once")
		}
		options.encoding, options.encodingSet = value, true
	default:
		return fmt.Errorf("unknown validate option %q", option)
	}
	return nil
}

func setValidateOperand(options *validateOptions, value string) error {
	if options == nil {
		return fmt.Errorf("validate options are unavailable")
	}
	if options.input != "" {
		return fmt.Errorf("validate accepts no additional arguments")
	}
	options.input = value
	return nil
}

func finalizeValidateOptions(options *validateOptions) error {
	if options == nil {
		return fmt.Errorf("validate options are unavailable")
	}
	if options.input == "" {
		return fmt.Errorf("validate requires one INPUT path")
	}
	if options.formatSet && options.format == "" {
		return fmt.Errorf("--format requires a value")
	}
	if options.encodingSet && options.encoding == "" {
		return fmt.Errorf("--encoding requires a value")
	}
	format, known := normalizeInputFormat(options.format)
	if !known {
		return fmt.Errorf("--format must be auto, cueson, srt, or vtt")
	}
	options.format = format
	if options.encoding != "" {
		if _, known := codec.NormalizeEncoding(options.encoding); !known {
			return fmt.Errorf("encoding %q is not supported", options.encoding)
		}
	}
	if options.format == "cueson" && options.encoding != "" {
		return fmt.Errorf("--encoding cannot be used with Cue JSON input")
	}
	if options.format == string(codec.FormatWebVTT) {
		if err := validateWebVTTInputEncoding(options.encoding); err != nil {
			return err
		}
	}
	return nil
}

func runValidate(ctx context.Context, options validateOptions, stderr io.Writer, diagnostics diagnosticWriter, usage helpTarget) int {
	if err := finalizeValidateOptions(&options); err != nil {
		diagnostics.write(diagnosticError, err.Error())
		writeUsage(stderr, usage)
		return ExitInvocation
	}
	captured, err := source.CaptureContext(ctx, options.input, source.CaptureOptions{})
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, fmt.Sprintf("input %q does not exist", options.input))
			writeUsage(stderr, usage)
			return ExitInvocation
		}
		if source.IsCapturePrecondition(err) {
			diagnostics.write(diagnosticError, "validate: input path or size precondition failed")
			writeUsage(stderr, usage)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("validate: capture input: %v", err))
		return ExitRuntimeFailure
	}
	loaded, err := loadValidatedInput(ctx, captured, inputOptions{format: options.format, encoding: options.encoding, disableSpeakerDetection: true})
	if err != nil {
		var constraint *inputConstraintError
		if errors.As(err, &constraint) {
			diagnostics.write(diagnosticError, constraint.Error())
			writeUsage(stderr, usage)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("validate: %v", err))
		return ExitRuntimeFailure
	}
	writeModelWarnings(diagnostics, loaded.diagnostics)
	diagnostics.write(diagnosticSuccess, fmt.Sprintf("valid %s.%s", loaded.inputKind, loaded.document.Format))
	return ExitSuccess
}
