// Package cli owns Cueson's command-line contract.
package cli

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shruggietech/cueson/internal/codec"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/version"
)

// Process exit statuses defined by the Cueson CLI contract.
const (
	ExitSuccess        = 0
	ExitRuntimeFailure = 1
	ExitInvocation     = 2
)

type globalOptions struct {
	quiet   bool
	silent  bool
	noColor bool
}

type schemaOptions struct {
	version   bool
	output    string
	outputSet bool
	force     bool
}

type restoreOptions struct {
	input          string
	output         string
	outputSet      bool
	outputDir      string
	outputDirSet   bool
	force          bool
	strictMetadata bool
	noMetadata     bool
}

type encodeOptions struct {
	input              string
	output             string
	outputSet          bool
	force              bool
	format             string
	encoding           string
	pretty             bool
	stdout             bool
	noSpeakerDetection bool
}

type renderOptions struct {
	input     string
	target    string
	output    string
	outputSet bool
	force     bool
	strict    bool
}

type convertOptions struct {
	input              string
	from               string
	target             string
	encoding           string
	output             string
	outputSet          bool
	force              bool
	strict             bool
	noSpeakerDetection bool
}

type helpTarget uint8

const (
	noHelp helpTarget = iota
	rootHelp
	versionHelp
	schemaHelp
	restoreHelp
	encodeHelp
	renderHelp
	convertHelp
	validateHelp
	inspectHelp
	completionHelp
)

type invocation struct {
	command         string
	help            helpTarget
	options         globalOptions
	schema          schemaOptions
	restore         restoreOptions
	encode          encodeOptions
	render          renderOptions
	convert         convertOptions
	validate        validateOptions
	inspect         inspectOptions
	completionShell string
}

type invocationError struct {
	message string
	usage   helpTarget
}

func (err *invocationError) Error() string {
	return err.message
}

type outputPreconditionError struct {
	message string
}

func (err *outputPreconditionError) Error() string {
	return err.message
}

type diagnosticLevel uint8

const (
	diagnosticSuccess diagnosticLevel = iota
	diagnosticInfo
	diagnosticWarning
	diagnosticError
)

type diagnosticWriter struct {
	writer io.Writer
	quiet  bool
	silent bool
	color  bool
	state  *diagnosticState
}

type diagnosticState struct {
	err error
}

// Run executes one Cueson command against the supplied process streams.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	_ = stdin

	parsed, parseErr := parseInvocation(args)
	diagnostics := newDiagnosticWriter(stderr, parsed.options)
	if parseErr != nil {
		diagnostics.write(diagnosticError, parseErr.Error())
		writeUsage(stderr, parseErr.usage)
		return diagnostics.finalize(ExitInvocation)
	}

	if parsed.help != noHelp {
		helpText, ok := fullHelpText(commandForHelp(parsed.help))
		if !ok {
			diagnostics.write(diagnosticError, "internal help dispatch failure")
			return diagnostics.finalize(ExitRuntimeFailure)
		}
		return diagnostics.finalize(writeStdout(stdout, diagnostics, []byte(helpText)))
	}

	select {
	case <-ctx.Done():
		diagnostics.write(diagnosticError, "operation canceled")
		return diagnostics.finalize(ExitRuntimeFailure)
	default:
	}

	status := ExitRuntimeFailure
	switch parsed.command {
	case "version":
		status = writeStdout(stdout, diagnostics, []byte(version.String()+"\n"))
	case "schema":
		status = runSchema(parsed.schema, stdout, stderr, diagnostics)
	case "restore":
		status = runRestore(ctx, parsed.restore, stderr, diagnostics)
	case "encode":
		status = runEncode(ctx, parsed.encode, stdout, stderr, diagnostics)
	case "render":
		status = runRender(ctx, parsed.render, stdout, stderr, diagnostics)
	case "convert":
		status = runConvert(ctx, parsed.convert, stdout, stderr, diagnostics)
	case "validate":
		status = runValidate(ctx, parsed.validate, stderr, diagnostics, validateHelp)
	case "inspect":
		status = runInspect(ctx, parsed.inspect, stdout, stderr, diagnostics, inspectHelp)
	case "completion":
		status = runCompletion(parsed.completionShell, stdout, stderr, diagnostics)
	default:
		diagnostics.write(diagnosticError, "internal command dispatch failure")
	}
	return diagnostics.finalize(status)
}

func runRestore(ctx context.Context, options restoreOptions, stderr io.Writer, diagnostics diagnosticWriter) int {
	payload, err := source.ReadFileContext(ctx, options.input)
	if err != nil {
		if source.IsCapturePrecondition(err) {
			diagnostics.write(diagnosticError, "restore: input does not satisfy path or size requirements")
			writeUsage(stderr, restoreHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, "restore: Cue JSON acquisition failed")
		return ExitRuntimeFailure
	}
	document, err := schema.Decode(payload)
	if err != nil {
		diagnostics.write(diagnosticError, err.Error())
		return ExitRuntimeFailure
	}
	metadataMode := source.MetadataDefault
	if options.strictMetadata {
		metadataMode = source.MetadataStrict
	} else if options.noMetadata {
		metadataMode = source.MetadataNone
	}
	report, err := source.Restore(ctx, document, source.RestoreOptions{Output: options.output, OutputDir: options.outputDir, Force: options.force, Metadata: metadataMode})
	if err != nil {
		var precondition *source.PreconditionError
		if errors.As(err, &precondition) {
			diagnostics.write(diagnosticError, precondition.Error())
			writeUsage(stderr, restoreHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("restore: %v", err))
		return ExitRuntimeFailure
	}
	for _, warning := range report.Warnings {
		diagnostics.write(diagnosticWarning, warning)
	}
	return ExitSuccess
}

func runSchema(options schemaOptions, stdout, stderr io.Writer, diagnostics diagnosticWriter) int {
	if err := schema.CheckLockstep(version.String()); err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("schema version lockstep: %v", err))
		return ExitRuntimeFailure
	}
	if options.version {
		return writeStdout(stdout, diagnostics, []byte(schema.Version()+"\n"))
	}
	if !options.outputSet {
		return writeStdout(stdout, diagnostics, schema.Bytes())
	}
	return handleSchemaFileResult(writeSchemaFile(options.output, schema.Bytes(), options.force), stderr, diagnostics)
}

func handleSchemaFileResult(err error, stderr io.Writer, diagnostics diagnosticWriter) int {
	if err != nil {
		var precondition *outputPreconditionError
		if errors.As(err, &precondition) {
			diagnostics.write(diagnosticError, precondition.Error())
			writeUsage(stderr, schemaHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("write schema: %v", err))
		return ExitRuntimeFailure
	}
	return ExitSuccess
}

func writeStdout(stdout io.Writer, diagnostics diagnosticWriter, payload []byte) int {
	written, err := stdout.Write(payload)
	if err == nil && written != len(payload) {
		err = io.ErrShortWrite
	}
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("write stdout: %v", err))
		return ExitRuntimeFailure
	}
	return ExitSuccess
}

func parseInvocation(args []string) (invocation, *invocationError) {
	var parsed invocation
	optionsActive := true

	for index := 0; index < len(args); index++ {
		arg := args[index]
		if optionsActive {
			switch arg {
			case "--":
				optionsActive = false
				continue
			case "-q", "--quiet":
				parsed.options.quiet = true
				continue
			case "--silent":
				parsed.options.silent = true
				continue
			case "--no-color":
				parsed.options.noColor = true
				continue
			case "-h", "--help":
				parsed.help = usageForCommand(parsed.command)
				return parsed, nil
			}

			if parsed.command == "schema" {
				switch arg {
				case "--version":
					parsed.schema.version = true
					continue
				case "-f", "--force":
					parsed.schema.force = true
					continue
				case "-o", "--output":
					if parsed.schema.outputSet {
						return parsed, schemaInvocationError("--output may be specified only once")
					}
					if index+1 >= len(args) {
						return parsed, schemaInvocationError("--output requires a path")
					}
					index++
					parsed.schema.output = args[index]
					parsed.schema.outputSet = true
					if parsed.schema.output == "" {
						return parsed, schemaInvocationError("--output requires a path")
					}
					continue
				}
				if strings.HasPrefix(arg, "--output=") {
					if parsed.schema.outputSet {
						return parsed, schemaInvocationError("--output may be specified only once")
					}
					parsed.schema.output = strings.TrimPrefix(arg, "--output=")
					parsed.schema.outputSet = true
					if parsed.schema.output == "" {
						return parsed, schemaInvocationError("--output requires a path")
					}
					continue
				}
			}
			if parsed.command == "restore" {
				switch arg {
				case "-f", "--force":
					parsed.restore.force = true
					continue
				case "--strict-metadata":
					parsed.restore.strictMetadata = true
					continue
				case "--no-metadata":
					parsed.restore.noMetadata = true
					continue
				case "-o", "--output":
					if parsed.restore.outputSet {
						return parsed, restoreInvocationError("--output may be specified only once")
					}
					if index+1 >= len(args) {
						return parsed, restoreInvocationError("--output requires a path")
					}
					index++
					parsed.restore.output, parsed.restore.outputSet = args[index], true
					if parsed.restore.output == "" {
						return parsed, restoreInvocationError("--output requires a path")
					}
					continue
				case "--output-dir":
					if parsed.restore.outputDirSet {
						return parsed, restoreInvocationError("--output-dir may be specified only once")
					}
					if index+1 >= len(args) {
						return parsed, restoreInvocationError("--output-dir requires a path")
					}
					index++
					parsed.restore.outputDir, parsed.restore.outputDirSet = args[index], true
					if parsed.restore.outputDir == "" {
						return parsed, restoreInvocationError("--output-dir requires a path")
					}
					continue
				}
				if strings.HasPrefix(arg, "--output=") {
					if parsed.restore.outputSet {
						return parsed, restoreInvocationError("--output may be specified only once")
					}
					parsed.restore.output, parsed.restore.outputSet = strings.TrimPrefix(arg, "--output="), true
					if parsed.restore.output == "" {
						return parsed, restoreInvocationError("--output requires a path")
					}
					continue
				}
				if strings.HasPrefix(arg, "--output-dir=") {
					if parsed.restore.outputDirSet {
						return parsed, restoreInvocationError("--output-dir may be specified only once")
					}
					parsed.restore.outputDir, parsed.restore.outputDirSet = strings.TrimPrefix(arg, "--output-dir="), true
					if parsed.restore.outputDir == "" {
						return parsed, restoreInvocationError("--output-dir requires a path")
					}
					continue
				}
			}
			if parsed.command == "encode" {
				switch arg {
				case "-f", "--force":
					parsed.encode.force = true
					continue
				case "--pretty":
					parsed.encode.pretty = true
					continue
				case "--stdout":
					parsed.encode.stdout = true
					continue
				case "--no-speaker-detection":
					parsed.encode.noSpeakerDetection = true
					continue
				case "-o", "--output", "--format", "--encoding":
					if index+1 >= len(args) {
						return parsed, encodeInvocationError(arg + " requires a value")
					}
					index++
					if args[index] == "" {
						return parsed, encodeInvocationError(arg + " requires a value")
					}
					switch arg {
					case "-o", "--output":
						if parsed.encode.outputSet {
							return parsed, encodeInvocationError("--output may be specified only once")
						}
						parsed.encode.output, parsed.encode.outputSet = args[index], true
					case "--format":
						if parsed.encode.format != "" {
							return parsed, encodeInvocationError("--format may be specified only once")
						}
						parsed.encode.format = args[index]
					case "--encoding":
						if parsed.encode.encoding != "" {
							return parsed, encodeInvocationError("--encoding may be specified only once")
						}
						parsed.encode.encoding = args[index]
					}
					continue
				}
				matched := false
				for _, option := range []string{"--output=", "--format=", "--encoding="} {
					if strings.HasPrefix(arg, option) {
						value := strings.TrimPrefix(arg, option)
						if value == "" {
							return parsed, encodeInvocationError(strings.TrimSuffix(option, "=") + " requires a value")
						}
						switch option {
						case "--output=":
							if parsed.encode.outputSet {
								return parsed, encodeInvocationError("--output may be specified only once")
							}
							parsed.encode.output, parsed.encode.outputSet = value, true
						case "--format=":
							if parsed.encode.format != "" {
								return parsed, encodeInvocationError("--format may be specified only once")
							}
							parsed.encode.format = value
						case "--encoding=":
							if parsed.encode.encoding != "" {
								return parsed, encodeInvocationError("--encoding may be specified only once")
							}
							parsed.encode.encoding = value
						}
						matched = true
						break
					}
				}
				if matched {
					continue
				}
			}
			if parsed.command == "render" {
				switch arg {
				case "-f", "--force":
					parsed.render.force = true
					continue
				case "--strict":
					parsed.render.strict = true
					continue
				case "-o", "--output", "--to":
					if index+1 >= len(args) {
						return parsed, renderInvocationError(arg + " requires a value")
					}
					index++
					if args[index] == "" {
						return parsed, renderInvocationError(arg + " requires a value")
					}
					if arg == "--to" {
						if parsed.render.target != "" {
							return parsed, renderInvocationError("--to may be specified only once")
						}
						parsed.render.target = args[index]
					} else {
						if parsed.render.outputSet {
							return parsed, renderInvocationError("--output may be specified only once")
						}
						parsed.render.output, parsed.render.outputSet = args[index], true
					}
					continue
				}
				if strings.HasPrefix(arg, "--output=") || strings.HasPrefix(arg, "--to=") {
					parts := strings.SplitN(arg, "=", 2)
					if parts[1] == "" {
						return parsed, renderInvocationError(parts[0] + " requires a value")
					}
					if parts[0] == "--to" {
						if parsed.render.target != "" {
							return parsed, renderInvocationError("--to may be specified only once")
						}
						parsed.render.target = parts[1]
					} else {
						if parsed.render.outputSet {
							return parsed, renderInvocationError("--output may be specified only once")
						}
						parsed.render.output, parsed.render.outputSet = parts[1], true
					}
					continue
				}
			}
			if parsed.command == "convert" {
				switch arg {
				case "-f", "--force":
					parsed.convert.force = true
					continue
				case "--strict":
					parsed.convert.strict = true
					continue
				case "--no-speaker-detection":
					parsed.convert.noSpeakerDetection = true
					continue
				case "-o", "--output", "--from", "--to", "--encoding":
					if index+1 >= len(args) {
						return parsed, convertInvocationError(arg + " requires a value")
					}
					index++
					if args[index] == "" {
						return parsed, convertInvocationError(arg + " requires a value")
					}
					if err := setConvertValue(&parsed.convert, arg, args[index]); err != nil {
						return parsed, err
					}
					continue
				}
				matched := false
				for _, option := range []string{"--output=", "--from=", "--to=", "--encoding="} {
					if !strings.HasPrefix(arg, option) {
						continue
					}
					value := strings.TrimPrefix(arg, option)
					if value == "" {
						return parsed, convertInvocationError(strings.TrimSuffix(option, "=") + " requires a value")
					}
					if err := setConvertValue(&parsed.convert, strings.TrimSuffix(option, "="), value); err != nil {
						return parsed, err
					}
					matched = true
					break
				}
				if matched {
					continue
				}
			}
			if parsed.command == "validate" {
				switch arg {
				case "--format", "--encoding":
					if index+1 >= len(args) {
						return parsed, validateInvocationError(arg + " requires a value")
					}
					index++
					if args[index] == "" {
						return parsed, validateInvocationError(arg + " requires a value")
					}
					if err := setValidateValue(&parsed.validate, arg, args[index]); err != nil {
						return parsed, validateInvocationError(err.Error())
					}
					continue
				}
				for _, option := range []string{"--format=", "--encoding="} {
					if !strings.HasPrefix(arg, option) {
						continue
					}
					value := strings.TrimPrefix(arg, option)
					if value == "" {
						return parsed, validateInvocationError(strings.TrimSuffix(option, "=") + " requires a value")
					}
					if err := setValidateValue(&parsed.validate, strings.TrimSuffix(option, "="), value); err != nil {
						return parsed, validateInvocationError(err.Error())
					}
					continue
				}
				if strings.HasPrefix(arg, "--format=") || strings.HasPrefix(arg, "--encoding=") {
					continue
				}
			}
			if parsed.command == "inspect" {
				switch arg {
				case "--json":
					if err := setInspectJSON(&parsed.inspect); err != nil {
						return parsed, inspectInvocationError(err.Error())
					}
					continue
				case "--format", "--encoding":
					if index+1 >= len(args) {
						return parsed, inspectInvocationError(arg + " requires a value")
					}
					index++
					if args[index] == "" {
						return parsed, inspectInvocationError(arg + " requires a value")
					}
					if err := setInspectValue(&parsed.inspect, arg, args[index]); err != nil {
						return parsed, inspectInvocationError(err.Error())
					}
					continue
				}
				for _, option := range []string{"--format=", "--encoding="} {
					if !strings.HasPrefix(arg, option) {
						continue
					}
					value := strings.TrimPrefix(arg, option)
					if value == "" {
						return parsed, inspectInvocationError(strings.TrimSuffix(option, "=") + " requires a value")
					}
					if err := setInspectValue(&parsed.inspect, strings.TrimSuffix(option, "="), value); err != nil {
						return parsed, inspectInvocationError(err.Error())
					}
					continue
				}
				if strings.HasPrefix(arg, "--format=") || strings.HasPrefix(arg, "--encoding=") {
					continue
				}
			}

			if strings.HasPrefix(arg, "-") {
				return parsed, &invocationError{message: fmt.Sprintf("unknown option %q", arg), usage: usageForCommand(parsed.command)}
			}
		}

		if parsed.command == "" {
			parsed.command = arg
			if _, exists := lookupCommandSurface(parsed.command); !exists {
				return parsed, &invocationError{message: fmt.Sprintf("unknown command %q", parsed.command), usage: rootHelp}
			}
			continue
		}

		if parsed.command == "restore" && parsed.restore.input == "" {
			parsed.restore.input = arg
			continue
		}
		if parsed.command == "encode" && parsed.encode.input == "" {
			parsed.encode.input = arg
			continue
		}
		if parsed.command == "render" && parsed.render.input == "" {
			parsed.render.input = arg
			continue
		}
		if parsed.command == "convert" && parsed.convert.input == "" {
			parsed.convert.input = arg
			continue
		}
		if parsed.command == "validate" {
			if err := setValidateOperand(&parsed.validate, arg); err != nil {
				return parsed, validateInvocationError(err.Error())
			}
			continue
		}
		if parsed.command == "inspect" {
			if err := setInspectOperand(&parsed.inspect, arg); err != nil {
				return parsed, inspectInvocationError(err.Error())
			}
			continue
		}
		if parsed.command == "completion" && parsed.completionShell == "" {
			parsed.completionShell = arg
			continue
		}
		if parsed.command != "restore" && parsed.command != "encode" && parsed.command != "render" && parsed.command != "convert" && parsed.command != "validate" && parsed.command != "inspect" && parsed.command != "completion" {
			return parsed, &invocationError{message: fmt.Sprintf("%s accepts no arguments", parsed.command), usage: usageForCommand(parsed.command)}
		}
		return parsed, &invocationError{message: fmt.Sprintf("%s accepts no additional arguments", parsed.command), usage: usageForCommand(parsed.command)}
	}

	if parsed.command == "" {
		parsed.help = rootHelp
		return parsed, nil
	}
	if parsed.command == "schema" {
		if parsed.schema.version && (parsed.schema.outputSet || parsed.schema.force) {
			return parsed, schemaInvocationError("--version cannot be combined with --output or --force")
		}
		if parsed.schema.force && !parsed.schema.outputSet {
			return parsed, schemaInvocationError("--force requires --output")
		}
	}
	if parsed.command == "restore" {
		if parsed.restore.input == "" {
			return parsed, restoreInvocationError("restore requires one INPUT path")
		}
		if parsed.restore.outputSet && parsed.restore.outputDirSet {
			return parsed, restoreInvocationError("--output and --output-dir are mutually exclusive")
		}
		if parsed.restore.strictMetadata && parsed.restore.noMetadata {
			return parsed, restoreInvocationError("--strict-metadata and --no-metadata are mutually exclusive")
		}
	}
	if parsed.command == "encode" {
		if parsed.encode.input == "" {
			return parsed, encodeInvocationError("encode requires one INPUT path")
		}
		if parsed.encode.format == "" {
			parsed.encode.format = "auto"
		}
		if parsed.encode.format != "auto" && parsed.encode.format != "srt" && parsed.encode.format != "vtt" && parsed.encode.format != "subrip" && parsed.encode.format != "webvtt" {
			return parsed, encodeInvocationError("--format must be auto, srt, or vtt")
		}
		if (parsed.encode.format == "vtt" || parsed.encode.format == "webvtt") && parsed.encode.encoding != "" {
			encoding, known := codec.NormalizeEncoding(parsed.encode.encoding)
			if !known || (encoding != codec.EncodingUTF8 && encoding != codec.EncodingUTF8BOM) {
				return parsed, encodeInvocationError("WebVTT requires UTF-8; --encoding is incompatible with --format vtt")
			}
		}
		stdoutSelected := parsed.encode.stdout || (parsed.encode.outputSet && parsed.encode.output == "-")
		if parsed.encode.stdout && parsed.encode.outputSet && parsed.encode.output != "-" {
			return parsed, encodeInvocationError("--stdout and --output are mutually exclusive")
		}
		if parsed.encode.force && stdoutSelected {
			return parsed, encodeInvocationError("--force requires a filesystem output")
		}
	}
	if parsed.command == "render" {
		if parsed.render.input == "" {
			return parsed, renderInvocationError("render requires one INPUT path")
		}
		if parsed.render.target == "" {
			return parsed, renderInvocationError("render requires --to FORMAT")
		}
		if parsed.render.target != "srt" && parsed.render.target != "subrip" && parsed.render.target != "vtt" && parsed.render.target != "webvtt" {
			return parsed, renderInvocationError("--to must be srt or vtt")
		}
		if parsed.render.force && (!parsed.render.outputSet || parsed.render.output == "-") {
			return parsed, renderInvocationError("--force requires a filesystem output")
		}
	}
	if parsed.command == "convert" {
		if parsed.convert.input == "" {
			return parsed, convertInvocationError("convert requires one INPUT path")
		}
		if parsed.convert.from == "" {
			parsed.convert.from = "auto"
		}
		from, ok := normalizeConvertSource(parsed.convert.from)
		if !ok {
			return parsed, convertInvocationError("--from must be auto, cueson, srt, or vtt")
		}
		parsed.convert.from = from
		if parsed.convert.target == "" {
			return parsed, convertInvocationError("convert requires --to FORMAT")
		}
		target, ok := normalizeSubtitleFormat(parsed.convert.target)
		if !ok {
			return parsed, convertInvocationError("--to must be srt or vtt")
		}
		parsed.convert.target = target
		canonicalEncoding := ""
		if parsed.convert.encoding != "" {
			var known bool
			canonicalEncoding, known = codec.NormalizeEncoding(parsed.convert.encoding)
			if !known {
				return parsed, convertInvocationError(fmt.Sprintf("encoding %q is not supported", parsed.convert.encoding))
			}
		}
		if parsed.convert.from == "cueson" && parsed.convert.encoding != "" {
			return parsed, convertInvocationError("--encoding cannot be used with Cue JSON input")
		}
		if parsed.convert.from == string(codec.FormatWebVTT) && parsed.convert.encoding != "" {
			if canonicalEncoding != codec.EncodingUTF8 && canonicalEncoding != codec.EncodingUTF8BOM {
				return parsed, convertInvocationError("WebVTT requires UTF-8; --encoding is incompatible with --from vtt")
			}
		}
		if parsed.convert.force && (!parsed.convert.outputSet || parsed.convert.output == "-") {
			return parsed, convertInvocationError("--force requires a filesystem output")
		}
	}
	if parsed.command == "validate" {
		if err := finalizeValidateOptions(&parsed.validate); err != nil {
			return parsed, validateInvocationError(err.Error())
		}
	}
	if parsed.command == "inspect" {
		if err := finalizeInspectOptions(&parsed.inspect); err != nil {
			return parsed, inspectInvocationError(err.Error())
		}
	}
	if parsed.command == "completion" {
		if parsed.completionShell == "" {
			return parsed, completionInvocationError("completion requires one SHELL selector")
		}
		if _, ok := completionScript(parsed.completionShell); !ok {
			return parsed, completionInvocationError("completion SHELL must be bash, zsh, fish, or powershell")
		}
	}
	return parsed, nil
}

func setConvertValue(options *convertOptions, option, value string) *invocationError {
	switch option {
	case "-o", "--output":
		if options.outputSet {
			return convertInvocationError("--output may be specified only once")
		}
		options.output, options.outputSet = value, true
	case "--from":
		if options.from != "" {
			return convertInvocationError("--from may be specified only once")
		}
		options.from = value
	case "--to":
		if options.target != "" {
			return convertInvocationError("--to may be specified only once")
		}
		options.target = value
	case "--encoding":
		if options.encoding != "" {
			return convertInvocationError("--encoding may be specified only once")
		}
		options.encoding = value
	}
	return nil
}

func normalizeConvertSource(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "auto":
		return "auto", true
	case "cueson", "json", "cue-json":
		return "cueson", true
	default:
		return normalizeSubtitleFormat(value)
	}
}

func normalizeSubtitleFormat(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "srt", "subrip":
		return string(codec.FormatSubRip), true
	case "vtt", "webvtt":
		return string(codec.FormatWebVTT), true
	default:
		return "", false
	}
}

func schemaInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: schemaHelp}
}

func restoreInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: restoreHelp}
}

func encodeInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: encodeHelp}
}

func renderInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: renderHelp}
}

func convertInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: convertHelp}
}

func validateInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: validateHelp}
}

func inspectInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: inspectHelp}
}

func completionInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: completionHelp}
}

func usageForCommand(command string) helpTarget {
	switch command {
	case "version":
		return versionHelp
	case "schema":
		return schemaHelp
	case "restore":
		return restoreHelp
	case "encode":
		return encodeHelp
	case "render":
		return renderHelp
	case "convert":
		return convertHelp
	case "validate":
		return validateHelp
	case "inspect":
		return inspectHelp
	case "completion":
		return completionHelp
	default:
		return rootHelp
	}
}

func writeUsage(writer io.Writer, target helpTarget) {
	usage, ok := shortUsageText(commandForHelp(target))
	if !ok {
		usage, _ = shortUsageText("")
	}
	fmt.Fprint(writer, usage)
}

func commandForHelp(target helpTarget) string {
	switch target {
	case versionHelp:
		return "version"
	case schemaHelp:
		return "schema"
	case restoreHelp:
		return "restore"
	case encodeHelp:
		return "encode"
	case renderHelp:
		return "render"
	case convertHelp:
		return "convert"
	case validateHelp:
		return "validate"
	case inspectHelp:
		return "inspect"
	case completionHelp:
		return "completion"
	default:
		return ""
	}
}

func runCompletion(shell string, stdout, stderr io.Writer, diagnostics diagnosticWriter) int {
	payload, ok := completionScript(shell)
	if !ok {
		diagnostics.write(diagnosticError, "completion SHELL must be bash, zsh, fish, or powershell")
		writeUsage(stderr, completionHelp)
		return ExitInvocation
	}
	return writeStdout(stdout, diagnostics, payload)
}

func writeSchemaFile(path string, payload []byte, force bool) error {
	return writeSchemaFileUsing(path, payload, force, publicationHooks{})
}

type publicationHooks struct {
	write   func(*os.File, []byte) error
	sync    func(*os.File) error
	close   func(*os.File) error
	verify  func(string, int64, string) error
	link    func(string, string) error
	replace func(string, string) error
}

func replaceSchemaFileUsing(path string, payload []byte, commit func(string, string) error) error {
	return writeSchemaFileUsing(path, payload, true, publicationHooks{replace: commit})
}

func writeSchemaFileUsing(path string, payload []byte, force bool, hooks publicationHooks) (result error) {
	existing, err := inspectOutputDestination(path, force)
	if err != nil {
		return err
	}

	directory := filepath.Dir(path)
	base := filepath.Base(path)
	temporary, err := os.CreateTemp(directory, "."+base+".tmp-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer func() {
		if temporaryPath != "" {
			_ = os.Remove(temporaryPath)
		}
	}()
	closed := false
	defer func() {
		if !closed {
			_ = temporary.Close()
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		closed = true
		return err
	}
	write := hooks.write
	if write == nil {
		write = func(file *os.File, data []byte) error { return writeAll(file, data) }
	}
	if err := write(temporary, payload); err != nil {
		_ = temporary.Close()
		closed = true
		return err
	}
	sync := hooks.sync
	if sync == nil {
		sync = func(file *os.File) error { return file.Sync() }
	}
	if err := sync(temporary); err != nil {
		_ = temporary.Close()
		closed = true
		return err
	}
	closeFile := hooks.close
	if closeFile == nil {
		closeFile = func(file *os.File) error { return file.Close() }
	}
	if err := closeFile(temporary); err != nil {
		_ = temporary.Close()
		closed = true
		return err
	}
	closed = true

	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	verify := hooks.verify
	if verify == nil {
		verify = verifyPublishedFile
	}
	if err := verify(temporaryPath, int64(len(payload)), digest); err != nil {
		return fmt.Errorf("verify staged output: %w", err)
	}
	stageInfo, err := os.Lstat(temporaryPath)
	if err != nil {
		return fmt.Errorf("inspect staged output: %w", err)
	}
	if !stageInfo.Mode().IsRegular() {
		return fmt.Errorf("staged output is not a regular file")
	}

	if existing == nil {
		link := hooks.link
		if link == nil {
			link = os.Link
		}
		if err := link(temporaryPath, path); err != nil {
			if errors.Is(err, fs.ErrExist) {
				return &outputPreconditionError{message: fmt.Sprintf("output %q appeared before publication; use --force to replace it", path)}
			}
			return fmt.Errorf("publish new output: %w", err)
		}
		committed, err := inspectCommittedOutput(path, stageInfo)
		if err != nil {
			return errors.Join(err, rollbackNewOutput(path, stageInfo))
		}
		if err := os.Remove(temporaryPath); err != nil {
			return errors.Join(fmt.Errorf("remove staging link: %w", err), rollbackNewOutput(path, committed))
		}
		temporaryPath = ""
		if err := verify(path, int64(len(payload)), digest); err != nil {
			return errors.Join(fmt.Errorf("verify published output: %w", err), rollbackNewOutput(path, committed))
		}
		return nil
	}

	current, err := os.Lstat(path)
	if err != nil || !current.Mode().IsRegular() || !os.SameFile(current, existing) {
		return fmt.Errorf("output %q changed before replacement", path)
	}
	backupPath, err := createOutputBackup(path)
	if err != nil {
		return fmt.Errorf("preserve output for rollback: %w", err)
	}
	removeBackup := true
	defer func() {
		if removeBackup {
			_ = os.Remove(backupPath)
		}
	}()
	current, err = os.Lstat(path)
	if err != nil || !current.Mode().IsRegular() || !os.SameFile(current, existing) {
		return fmt.Errorf("output %q changed while preparing replacement", path)
	}
	replace := hooks.replace
	if replace == nil {
		replace = os.Rename
	}
	if err := replace(temporaryPath, path); err != nil {
		return fmt.Errorf("replace output: %w", err)
	}
	temporaryPath = ""
	committed, err := inspectPublishedOutput(path)
	if err != nil {
		rollbackErr := rollbackReplacedOutput(path, stageInfo, backupPath, existing)
		removeBackup = rollbackErr == nil
		return errors.Join(err, rollbackErr)
	}
	if err := verify(path, int64(len(payload)), digest); err != nil {
		rollbackErr := rollbackReplacedOutput(path, committed, backupPath, existing)
		removeBackup = rollbackErr == nil
		return errors.Join(fmt.Errorf("verify published output: %w", err), rollbackErr)
	}
	if err := os.Remove(backupPath); err != nil {
		rollbackErr := rollbackReplacedOutput(path, committed, backupPath, existing)
		removeBackup = rollbackErr == nil
		return errors.Join(fmt.Errorf("remove rollback backup: %w", err), rollbackErr)
	}
	removeBackup = false
	return nil
}

func inspectOutputDestination(path string, force bool) (fs.FileInfo, error) {
	info, err := os.Lstat(path)
	switch {
	case err == nil:
		if info.Mode()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return nil, &outputPreconditionError{message: fmt.Sprintf("output %q is not a regular file", path)}
		}
		if !force {
			return nil, &outputPreconditionError{message: fmt.Sprintf("output %q already exists; use --force to replace it", path)}
		}
		return info, nil
	case errors.Is(err, fs.ErrNotExist):
		return nil, nil
	default:
		return nil, err
	}
}

func verifyPublishedFile(path string, expectedSize int64, expectedDigest string) error {
	named, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if named.Mode()&fs.ModeSymlink != 0 || !named.Mode().IsRegular() {
		return fmt.Errorf("output is not a regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	hash := sha256.New()
	count, readErr := io.Copy(hash, file)
	opened, statErr := file.Stat()
	closeErr := file.Close()
	if readErr != nil {
		return readErr
	}
	if statErr != nil {
		return statErr
	}
	if closeErr != nil {
		return closeErr
	}
	if !opened.Mode().IsRegular() || !os.SameFile(named, opened) {
		return fmt.Errorf("output changed while it was verified")
	}
	current, err := os.Lstat(path)
	if err != nil || !current.Mode().IsRegular() || !os.SameFile(opened, current) {
		return fmt.Errorf("output changed after verification")
	}
	actualDigest := fmt.Sprintf("%x", hash.Sum(nil))
	if count != expectedSize || actualDigest != expectedDigest {
		return fmt.Errorf("output integrity is size %d and SHA-256 %s, want size %d and SHA-256 %s", count, actualDigest, expectedSize, expectedDigest)
	}
	return nil
}

func inspectCommittedOutput(path string, expected fs.FileInfo) (fs.FileInfo, error) {
	current, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect published output: %w", err)
	}
	if !current.Mode().IsRegular() || !os.SameFile(current, expected) {
		return nil, fmt.Errorf("published output changed during commit")
	}
	return current, nil
}

func inspectPublishedOutput(path string) (fs.FileInfo, error) {
	current, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect published output: %w", err)
	}
	if current.Mode()&fs.ModeSymlink != 0 || !current.Mode().IsRegular() {
		return nil, fmt.Errorf("published output is not a regular file")
	}
	return current, nil
}

func rollbackNewOutput(path string, committed fs.FileInfo) error {
	current, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect new output for rollback: %w", err)
	}
	if !current.Mode().IsRegular() || !os.SameFile(current, committed) {
		return fmt.Errorf("new output changed before rollback and was retained")
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove failed new output: %w", err)
	}
	return nil
}

func createOutputBackup(path string) (string, error) {
	placeholder, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".backup-*")
	if err != nil {
		return "", err
	}
	backupPath := placeholder.Name()
	if err := placeholder.Close(); err != nil {
		_ = os.Remove(backupPath)
		return "", err
	}
	if err := os.Remove(backupPath); err != nil {
		return "", err
	}
	if err := os.Link(path, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}

func rollbackReplacedOutput(path string, committed fs.FileInfo, backupPath string, original fs.FileInfo) error {
	current, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("inspect replacement for rollback: %w; original retained at %q", err, backupPath)
	}
	if !current.Mode().IsRegular() || !os.SameFile(current, committed) {
		return fmt.Errorf("replacement changed before rollback; original retained at %q", backupPath)
	}
	if err := os.Rename(backupPath, path); err != nil {
		return fmt.Errorf("restore rollback backup %q: %w", backupPath, err)
	}
	restored, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("verify restored rollback output: %w", err)
	}
	if !restored.Mode().IsRegular() || !os.SameFile(restored, original) {
		return fmt.Errorf("restored rollback output does not match the original")
	}
	return nil
}

func writeAll(writer io.Writer, payload []byte) error {
	written, err := writer.Write(payload)
	if err == nil && written != len(payload) {
		err = io.ErrShortWrite
	}
	return err
}

func newDiagnosticWriter(writer io.Writer, options globalOptions) diagnosticWriter {
	_, noColorPresent := os.LookupEnv("NO_COLOR")
	return diagnosticWriter{writer: writer, quiet: options.quiet, silent: options.silent, color: diagnosticColorEnabled(options.noColor, noColorPresent, isTerminal(writer)), state: &diagnosticState{}}
}

func diagnosticColorEnabled(explicitOff, noColorPresent, terminal bool) bool {
	return !explicitOff && !noColorPresent && terminal
}

func isTerminal(writer io.Writer) bool {
	type statter interface {
		Stat() (fs.FileInfo, error)
	}
	stream, ok := writer.(statter)
	if !ok {
		return false
	}
	info, err := stream.Stat()
	return err == nil && info.Mode()&fs.ModeCharDevice != 0
}

func (diagnostics diagnosticWriter) write(level diagnosticLevel, message string) {
	if diagnostics.suppresses(level) {
		return
	}
	label, colorCode := diagnosticLabel(level)
	var err error
	if diagnostics.color {
		_, err = fmt.Fprintf(diagnostics.writer, "\x1b[%sm%s:\x1b[0m %s\n", colorCode, label, message)
	} else {
		_, err = fmt.Fprintf(diagnostics.writer, "%s: %s\n", label, message)
	}
	if err != nil && diagnostics.state != nil && diagnostics.state.err == nil {
		diagnostics.state.err = err
	}
}

func (diagnostics diagnosticWriter) finalize(status int) int {
	if status == ExitSuccess && diagnostics.state != nil && diagnostics.state.err != nil {
		return ExitRuntimeFailure
	}
	return status
}

func (diagnostics diagnosticWriter) suppresses(level diagnosticLevel) bool {
	if level == diagnosticError {
		return false
	}
	if diagnostics.silent {
		return true
	}
	return diagnostics.quiet && (level == diagnosticSuccess || level == diagnosticInfo)
}

func diagnosticLabel(level diagnosticLevel) (string, string) {
	switch level {
	case diagnosticSuccess:
		return "success", "32"
	case diagnosticInfo:
		return "info", "36"
	case diagnosticWarning:
		return "warning", "33"
	default:
		return "error", "31"
	}
}
