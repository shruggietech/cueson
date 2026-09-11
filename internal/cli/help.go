package cli

import (
	"fmt"
	"strings"
)

func fullHelpText(command string) (string, bool) {
	if command == "" {
		return rootSurfaceHelp(), true
	}
	surface, ok := lookupCommandSurface(command)
	if !ok {
		return "", false
	}
	return commandSurfaceHelp(surface), true
}

func shortUsageText(command string) (string, bool) {
	if command == "" {
		return "Usage:\n  cueson [global options] <command>\n", true
	}
	surface, ok := lookupCommandSurface(command)
	if !ok {
		return "", false
	}
	return "Usage:\n  " + surface.Usage + "\n", true
}

func rootSurfaceHelp() string {
	var output strings.Builder
	output.WriteString("Cueson command-line interface.\n\n")
	output.WriteString("Usage:\n  cueson [global options] <command>\n\n")
	output.WriteString("Commands:\n")
	width := 0
	for _, command := range orderedCommandSurface {
		if len(command.Name) > width {
			width = len(command.Name)
		}
	}
	for _, command := range orderedCommandSurface {
		fmt.Fprintf(&output, "  %s%s  %s\n", command.Name, strings.Repeat(" ", width-len(command.Name)), command.Summary)
	}
	output.WriteString("\nGlobal options:\n")
	writeOptionHelp(&output, globalOptionSurface)
	output.WriteString("\nStreams:\n")
	output.WriteString("  stdout  Payload data and explicit help only; never diagnostics or decoration.\n")
	output.WriteString("  stderr  Diagnostics, warnings, errors, and error-associated short usage.\n")
	output.WriteString("  Quiet and silent never suppress an explicitly requested stdout payload.\n")
	writeExitCodeHelp(&output)
	output.WriteString("\nExamples:\n")
	output.WriteString("  cueson --help\n")
	output.WriteString("  cueson validate captions.srt\n")
	output.WriteString("  cueson inspect --json document.cueson.json\n")
	output.WriteString("  cueson convert captions.vtt --to srt\n")
	output.WriteString("  cueson completion bash\n")
	return output.String()
}

func commandSurfaceHelp(command cliCommandSpec) string {
	var output strings.Builder
	output.WriteString(command.Description)
	output.WriteString("\n\nUsage:\n  ")
	output.WriteString(command.Usage)
	output.WriteString("\n\nOptions:\n")
	if len(command.Options) == 0 {
		output.WriteString("  None.\n")
	} else {
		writeOptionHelp(&output, command.Options)
	}
	output.WriteString("\nGlobal options:\n")
	writeOptionHelp(&output, globalOptionSurface)
	if len(command.Notes) > 0 {
		output.WriteString("\nNotes:\n")
		for _, note := range command.Notes {
			output.WriteString("  ")
			output.WriteString(note)
			output.WriteByte('\n')
		}
	}
	output.WriteString("\nStreams:\n")
	for _, stream := range command.Streams {
		output.WriteString("  ")
		output.WriteString(stream)
		output.WriteByte('\n')
	}
	writeExitCodeHelp(&output)
	output.WriteString("\nExamples:\n")
	for _, example := range command.Examples {
		output.WriteString("  ")
		output.WriteString(example)
		output.WriteByte('\n')
	}
	return output.String()
}

func writeOptionHelp(output *strings.Builder, options []cliOptionSpec) {
	labels := make([]string, len(options))
	width := 0
	for index, option := range options {
		labels[index] = strings.Join(option.Spellings, ", ")
		if option.ValueName != "" {
			labels[index] += " " + option.ValueName
		}
		if len(labels[index]) > width {
			width = len(labels[index])
		}
	}
	for index, option := range options {
		output.WriteString("  ")
		output.WriteString(labels[index])
		output.WriteString(strings.Repeat(" ", width-len(labels[index])))
		output.WriteString("  ")
		output.WriteString(option.Description)
		output.WriteByte('\n')
	}
}

func writeExitCodeHelp(output *strings.Builder) {
	output.WriteString("\nExit codes:\n")
	output.WriteString("  0  Command or explicit help completed successfully; warnings may be present.\n")
	output.WriteString("  1  An accepted operation failed during processing, cancellation, or runtime I/O.\n")
	output.WriteString("  2  Command syntax, an operand, an option, or a pre-execution requirement was invalid.\n")
}
