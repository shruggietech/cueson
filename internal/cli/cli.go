// Package cli owns Cueson's command-line contract.
package cli

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

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

type helpTarget uint8

const (
	noHelp helpTarget = iota
	rootHelp
	versionHelp
)

type invocation struct {
	command string
	help    helpTarget
	options globalOptions
}

type invocationError struct {
	message string
	usage   helpTarget
}

func (err *invocationError) Error() string {
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

Global options:
  -q, --quiet  Suppress informational and success diagnostics.
  --silent     Suppress every non-error diagnostic.
  --no-color   Disable diagnostic color.
  -h, --help   Show help.
  --            End option processing.

Examples:
  cueson --help
  cueson version
`

const versionHelpText = `Print the Cueson executable version.

Usage:
  cueson [global options] version

Example:
  cueson version
`

const rootUsageText = `Usage:
  cueson [global options] <command>
`

const versionUsageText = `Usage:
  cueson [global options] version
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
		fmt.Fprint(stdout, rootHelpText)
		return ExitSuccess
	case versionHelp:
		fmt.Fprint(stdout, versionHelpText)
		return ExitSuccess
	}

	select {
	case <-ctx.Done():
		diagnostics.write(diagnosticError, "operation canceled")
		return ExitRuntimeFailure
	default:
	}

	fmt.Fprintln(stdout, version.String())
	return ExitSuccess
}

func parseInvocation(args []string) (invocation, *invocationError) {
	var parsed invocation
	optionsActive := true

	for _, arg := range args {
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
				if parsed.command == "version" {
					parsed.help = versionHelp
				} else {
					parsed.help = rootHelp
				}
				return parsed, nil
			}

			if strings.HasPrefix(arg, "-") {
				return parsed, &invocationError{
					message: fmt.Sprintf("unknown option %q", arg),
					usage:   usageForCommand(parsed.command),
				}
			}
		}

		if parsed.command == "" {
			parsed.command = arg
			if parsed.command != "version" {
				return parsed, &invocationError{
					message: fmt.Sprintf("unknown command %q", parsed.command),
					usage:   rootHelp,
				}
			}
			continue
		}

		return parsed, &invocationError{
			message: fmt.Sprintf("%s accepts no arguments", parsed.command),
			usage:   versionHelp,
		}
	}

	if parsed.command == "" {
		parsed.help = rootHelp
	}

	return parsed, nil
}

func usageForCommand(command string) helpTarget {
	if command == "version" {
		return versionHelp
	}
	return rootHelp
}

func writeUsage(writer io.Writer, target helpTarget) {
	if target == versionHelp {
		fmt.Fprint(writer, versionUsageText)
		return
	}
	fmt.Fprint(writer, rootUsageText)
}

func newDiagnosticWriter(writer io.Writer, options globalOptions) diagnosticWriter {
	_, noColorPresent := os.LookupEnv("NO_COLOR")
	return diagnosticWriter{
		writer: writer,
		quiet:  options.quiet,
		silent: options.silent,
		color:  diagnosticColorEnabled(options.noColor, noColorPresent, isTerminal(writer)),
	}
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
