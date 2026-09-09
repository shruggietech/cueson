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

type helpTarget uint8

const (
	noHelp helpTarget = iota
	rootHelp
	versionHelp
	schemaHelp
)

type invocation struct {
	command string
	help    helpTarget
	options globalOptions
	schema  schemaOptions
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
  version  Print the Cueson executable version.
  schema   Print or save the embedded Cue JSON schema.

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

const rootUsageText = `Usage:
  cueson [global options] <command>
`

const versionUsageText = `Usage:
  cueson [global options] version
`

const schemaUsageText = `Usage:
  cueson [global options] schema [--version | --output PATH [--force]]
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
	default:
		diagnostics.write(diagnosticError, "internal command dispatch failure")
		return ExitRuntimeFailure
	}
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

			if strings.HasPrefix(arg, "-") {
				return parsed, &invocationError{message: fmt.Sprintf("unknown option %q", arg), usage: usageForCommand(parsed.command)}
			}
		}

		if parsed.command == "" {
			parsed.command = arg
			if parsed.command != "version" && parsed.command != "schema" {
				return parsed, &invocationError{message: fmt.Sprintf("unknown command %q", parsed.command), usage: rootHelp}
			}
			continue
		}

		return parsed, &invocationError{message: fmt.Sprintf("%s accepts no arguments", parsed.command), usage: usageForCommand(parsed.command)}
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
	return parsed, nil
}

func schemaInvocationError(message string) *invocationError {
	return &invocationError{message: message, usage: schemaHelp}
}

func usageForCommand(command string) helpTarget {
	switch command {
	case "version":
		return versionHelp
	case "schema":
		return schemaHelp
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
