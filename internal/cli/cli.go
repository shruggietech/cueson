// Package cli owns Cueson's command-line contract.
package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

type helpTarget uint8

const (
	noHelp helpTarget = iota
	rootHelp
	versionHelp
	schemaHelp
	restoreHelp
	encodeHelp
	renderHelp
)

type invocation struct {
	command string
	help    helpTarget
	options globalOptions
	schema  schemaOptions
	restore restoreOptions
	encode  encodeOptions
	render  renderOptions
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
}

const rootHelpText = `Cueson command-line interface.

Usage:
  cueson [global options] <command>

Commands:
	encode   Encode a subtitle source as Cue JSON.
	render   Render Cue JSON from its structured model.
  version  Print the Cueson executable version.
  schema   Print or save the embedded Cue JSON schema.
  restore  Restore exact source-envelope bytes.

Global options:
  -q, --quiet  Suppress informational and success diagnostics.
  --silent     Suppress every non-error diagnostic.
  --no-color   Disable diagnostic color.
  -h, --help   Show help.
  --            End option processing.

Examples:
  cueson --help
  cueson version
  cueson schema --version
  cueson schema --output cueson.schema.json
	cueson encode captions.srt
	cueson render captions.srt.cueson.json --to srt
  cueson restore --output restored.srt document.cueson.json
`

const versionHelpText = `Print the Cueson executable version.

Usage:
  cueson [global options] version

Example:
  cueson version
`

const schemaHelpText = `Print or save the embedded canonical Cue JSON schema.

Usage:
  cueson [global options] schema [options]

Options:
  --version          Print the embedded schema version.
  -o, --output PATH  Write the schema to a literal path instead of stdout.
  -f, --force        Replace an existing regular output file.
  -h, --help         Show schema help.

Examples:
  cueson schema
  cueson schema --version
  cueson schema --output cueson.schema.json
  cueson schema --output cueson.schema.json --force
`

const restoreHelpText = `Restore exact source-envelope bytes without a format codec.

Usage:
  cueson [global options] restore [options] INPUT

Options:
  -o, --output PATH       Restore one asset to a literal path.
  --output-dir DIR        Restore all assets beneath an existing directory.
  -f, --force             Replace approved existing regular files.
  --strict-metadata       Require every captured timestamp to be restored.
  --no-metadata           Skip timestamp restoration.
  -h, --help              Show restore help.

Successful restoration writes no stdout payload. Warnings and errors use stderr.
Multi-asset documents require --output-dir. Metadata modes are mutually exclusive.

Example:
  cueson restore --output restored.srt document.cueson.json
`

const encodeHelpText = `Encode a subtitle source as Cue JSON while preserving its exact bytes.

Usage:
  cueson [global options] encode [options] INPUT

Options:
  -o, --output PATH       Write Cue JSON to PATH (default INPUT.cueson.json).
  -f, --force             Replace an existing regular output file.
  --format FORMAT         Select auto, srt, or vtt (default auto).
  --encoding NAME         Select utf-8, utf-16le, utf-16be, windows-1252, or iso-8859-1.
  --pretty                Indent Cue JSON output.
  --stdout                Write Cue JSON to stdout.
  --no-speaker-detection  Disable derived Name: speaker observations.
  -h, --help              Show encode help.

Encoding aliases: utf8; utf-8-bom, utf8-bom, utf-8-sig; utf16le, utf-16-le;
utf16be, utf-16-be; windows1252, cp1252; iso8859-1, latin1, latin-1.

Example:
  cueson encode --pretty captions.srt
`

const renderHelpText = `Render Cue JSON from its structured model.

Usage:
  cueson [global options] render [options] INPUT.cueson.json --to srt

Options:
  --to FORMAT        Select the target format.
  -o, --output PATH  Write output to PATH instead of stdout.
  -f, --force        Replace an existing regular output file.
  --strict           Reject known non-representable model content.
  -h, --help         Show render help.

Example:
  cueson render captions.srt.cueson.json --to srt --output captions.rendered.srt
`

const rootUsageText = `Usage:
  cueson [global options] <command>
`

const versionUsageText = `Usage:
  cueson [global options] version
`

const schemaUsageText = `Usage:
  cueson [global options] schema [--version | --output PATH [--force]]
`

const restoreUsageText = `Usage:
  cueson [global options] restore [--output PATH | --output-dir DIR] [--force] [--strict-metadata | --no-metadata] INPUT
`

const encodeUsageText = `Usage:
  cueson [global options] encode [--output PATH | --stdout] [--force] [--format FORMAT] [--encoding NAME] [--pretty] [--no-speaker-detection] INPUT
`

const renderUsageText = `Usage:
  cueson [global options] render [--output PATH] [--force] [--strict] INPUT --to FORMAT
`

// Run executes one Cueson command against the supplied process streams.
func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	_ = stdin

	parsed, parseErr := parseInvocation(args)
	diagnostics := newDiagnosticWriter(stderr, parsed.options)
	if parseErr != nil {
		diagnostics.write(diagnosticError, parseErr.Error())
		writeUsage(stderr, parseErr.usage)
		return ExitInvocation
	}

	switch parsed.help {
	case rootHelp:
		return writeStdout(stdout, diagnostics, []byte(rootHelpText))
	case versionHelp:
		return writeStdout(stdout, diagnostics, []byte(versionHelpText))
	case schemaHelp:
		return writeStdout(stdout, diagnostics, []byte(schemaHelpText))
	case restoreHelp:
		return writeStdout(stdout, diagnostics, []byte(restoreHelpText))
	case encodeHelp:
		return writeStdout(stdout, diagnostics, []byte(encodeHelpText))
	case renderHelp:
		return writeStdout(stdout, diagnostics, []byte(renderHelpText))
	}

	select {
	case <-ctx.Done():
		diagnostics.write(diagnosticError, "operation canceled")
		return ExitRuntimeFailure
	default:
	}

	switch parsed.command {
	case "version":
		return writeStdout(stdout, diagnostics, []byte(version.String()+"\n"))
	case "schema":
		return runSchema(parsed.schema, stdout, stderr, diagnostics)
	case "restore":
		return runRestore(ctx, parsed.restore, stderr, diagnostics)
	case "encode":
		return runEncode(ctx, parsed.encode, stdout, stderr, diagnostics)
	case "render":
		return runRender(ctx, parsed.render, stdout, stderr, diagnostics)
	default:
		diagnostics.write(diagnosticError, "internal command dispatch failure")
		return ExitRuntimeFailure
	}
}

func runRestore(ctx context.Context, options restoreOptions, stderr io.Writer, diagnostics diagnosticWriter) int {
	inputInfo, err := os.Stat(options.input)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			diagnostics.write(diagnosticError, fmt.Sprintf("input %q does not exist", options.input))
			writeUsage(stderr, restoreHelp)
			return ExitInvocation
		}
		diagnostics.write(diagnosticError, fmt.Sprintf("read Cue JSON %q: %v", options.input, err))
		return ExitRuntimeFailure
	}
	if !inputInfo.Mode().IsRegular() {
		diagnostics.write(diagnosticError, fmt.Sprintf("input %q is not a regular file", options.input))
		writeUsage(stderr, restoreHelp)
		return ExitInvocation
	}
	payload, err := os.ReadFile(options.input)
	if err != nil {
		diagnostics.write(diagnosticError, fmt.Sprintf("read Cue JSON %q: %v", options.input, err))
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

			if strings.HasPrefix(arg, "-") {
				return parsed, &invocationError{message: fmt.Sprintf("unknown option %q", arg), usage: usageForCommand(parsed.command)}
			}
		}

		if parsed.command == "" {
			parsed.command = arg
			if parsed.command != "version" && parsed.command != "schema" && parsed.command != "restore" && parsed.command != "encode" && parsed.command != "render" {
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
		if parsed.command != "restore" && parsed.command != "encode" && parsed.command != "render" {
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
	return parsed, nil
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
	default:
		return rootHelp
	}
}

func writeUsage(writer io.Writer, target helpTarget) {
	switch target {
	case versionHelp:
		fmt.Fprint(writer, versionUsageText)
	case schemaHelp:
		fmt.Fprint(writer, schemaUsageText)
	case restoreHelp:
		fmt.Fprint(writer, restoreUsageText)
	case encodeHelp:
		fmt.Fprint(writer, encodeUsageText)
	case renderHelp:
		fmt.Fprint(writer, renderUsageText)
	default:
		fmt.Fprint(writer, rootUsageText)
	}
}

func writeSchemaFile(path string, payload []byte, force bool) error {
	info, err := os.Lstat(path)
	switch {
	case err == nil:
		if !info.Mode().IsRegular() {
			return &outputPreconditionError{message: fmt.Sprintf("output %q is not a regular file", path)}
		}
		if !force {
			return &outputPreconditionError{message: fmt.Sprintf("output %q already exists; use --force to replace it", path)}
		}
		return replaceSchemaFile(path, payload)
	case !errors.Is(err, fs.ErrNotExist):
		return err
	default:
		return createSchemaFile(path, payload)
	}
}

func createSchemaFile(path string, payload []byte) (result error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return &outputPreconditionError{message: fmt.Sprintf("output %q already exists; use --force to replace it", path)}
		}
		return err
	}
	defer func() {
		if result != nil {
			_ = os.Remove(path)
		}
	}()
	if err := writeAll(file, payload); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func replaceSchemaFile(path string, payload []byte) (result error) {
	return replaceSchemaFileUsing(path, payload, os.Rename)
}

func replaceSchemaFileUsing(path string, payload []byte, commit func(string, string) error) (result error) {
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
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := writeAll(temporary, payload); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}

	if err := commit(temporaryPath, path); err != nil {
		return err
	}
	temporaryPath = ""
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
	return diagnosticWriter{writer: writer, quiet: options.quiet, silent: options.silent, color: diagnosticColorEnabled(options.noColor, noColorPresent, isTerminal(writer))}
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
	if diagnostics.color {
		fmt.Fprintf(diagnostics.writer, "\x1b[%sm%s:\x1b[0m %s\n", colorCode, label, message)
		return
	}
	fmt.Fprintf(diagnostics.writer, "%s: %s\n", label, message)
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
