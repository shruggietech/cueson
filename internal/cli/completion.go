package cli

import (
	"strings"
)

var completionShells = []string{"bash", "zsh", "fish", "powershell"}

func completionScript(shell string) ([]byte, bool) {
	var script string
	switch shell {
	case "bash":
		script = bashCompletionScript()
	case "zsh":
		script = zshCompletionScript()
	case "fish":
		script = fishCompletionScript()
	case "powershell":
		script = powershellCompletionScript()
	default:
		return nil, false
	}
	return []byte(script), true
}

func bashCompletionScript() string {
	var output strings.Builder
	output.WriteString("# Bash completion for cueson.\n")
	output.WriteString("_cueson_filter_candidates() {\n")
	output.WriteString("  local prefix=\"$1\" candidate\n")
	output.WriteString("  shift\n")
	output.WriteString("  COMPREPLY=()\n")
	output.WriteString("  for candidate in \"$@\"; do\n")
	output.WriteString("    case \"$candidate\" in \"$prefix\"*) COMPREPLY[${#COMPREPLY[@]}]=\"$candidate\" ;; esac\n")
	output.WriteString("  done\n")
	output.WriteString("}\n\n")
	output.WriteString("_cueson_completion() {\n")
	output.WriteString("  local current=\"${COMP_WORDS[COMP_CWORD]}\" previous=\"\" command=\"\" word candidate\n")
	output.WriteString("  local options_active=1 index\n")
	output.WriteString("  if (( COMP_CWORD > 0 )); then previous=\"${COMP_WORDS[COMP_CWORD-1]}\"; fi\n")
	output.WriteString("  for (( index=1; index<COMP_CWORD; index++ )); do\n")
	output.WriteString("    word=\"${COMP_WORDS[index]}\"\n")
	output.WriteString("    if [[ \"$word\" == \"--\" ]]; then options_active=0; continue; fi\n")
	output.WriteString("    if [[ -z \"$command\" ]]; then\n")
	output.WriteString("      case \"$word\" in\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("        ")
		output.WriteString(command.Name)
		output.WriteString(") command=\"")
		output.WriteString(command.Name)
		output.WriteString("\" ;;\n")
	}
	output.WriteString("      esac\n")
	output.WriteString("    fi\n")
	output.WriteString("  done\n")
	output.WriteString("  if (( ! options_active )); then COMPREPLY=(); return; fi\n\n")
	writeBashValueCases(&output)
	output.WriteString("  case \"$command\" in\n")
	output.WriteString("    \"\") _cueson_filter_candidates \"$current\"")
	for _, command := range orderedCommandSurface {
		output.WriteByte(' ')
		output.WriteString(bashQuote(command.Name))
	}
	for _, option := range globalOptionSurface {
		for _, spelling := range option.Spellings {
			output.WriteByte(' ')
			output.WriteString(bashQuote(spelling))
		}
	}
	output.WriteString(" ;;\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("    ")
		output.WriteString(command.Name)
		output.WriteString(") _cueson_filter_candidates \"$current\"")
		for _, value := range command.OperandValues {
			output.WriteByte(' ')
			output.WriteString(bashQuote(value))
		}
		for _, option := range appendOptions(globalOptionSurface, command.Options) {
			for _, spelling := range option.Spellings {
				output.WriteByte(' ')
				output.WriteString(bashQuote(spelling))
			}
		}
		output.WriteString(" ;;\n")
	}
	output.WriteString("    *) COMPREPLY=() ;;\n")
	output.WriteString("  esac\n")
	output.WriteString("}\n\n")
	output.WriteString("complete -F _cueson_completion cueson\n")
	return output.String()
}

func writeBashValueCases(output *strings.Builder) {
	output.WriteString("  case \"$command:$previous\" in\n")
	for _, command := range orderedCommandSurface {
		for _, option := range command.Options {
			if len(option.Values) == 0 {
				continue
			}
			for _, spelling := range option.Spellings {
				output.WriteString("    ")
				output.WriteString(command.Name)
				output.WriteByte(':')
				output.WriteString(spelling)
				output.WriteString(") _cueson_filter_candidates \"$current\"")
				for _, value := range option.Values {
					output.WriteByte(' ')
					output.WriteString(bashQuote(value))
				}
				output.WriteString("; return ;;\n")
			}
		}
	}
	output.WriteString("  esac\n")
	for _, command := range orderedCommandSurface {
		for _, option := range command.Options {
			if len(option.Values) == 0 {
				continue
			}
			for _, spelling := range option.Spellings {
				if !strings.HasPrefix(spelling, "--") {
					continue
				}
				output.WriteString("  if [[ \"$command\" == ")
				output.WriteString(bashQuote(command.Name))
				output.WriteString(" && \"$current\" == ")
				output.WriteString(bashQuote(spelling + "="))
				output.WriteString("* ]]; then\n")
				output.WriteString("    local value_prefix=\"${current#*=}\"\n")
				output.WriteString("    _cueson_filter_candidates \"$value_prefix\"")
				for _, value := range option.Values {
					output.WriteByte(' ')
					output.WriteString(bashQuote(value))
				}
				output.WriteString("\n")
				output.WriteString("    for (( index=0; index<${#COMPREPLY[@]}; index++ )); do COMPREPLY[index]=")
				output.WriteString(bashQuote(spelling + "="))
				output.WriteString("\"${COMPREPLY[index]}\"; done\n")
				output.WriteString("    return\n  fi\n")
			}
		}
	}
}

func zshCompletionScript() string {
	var output strings.Builder
	output.WriteString("#compdef cueson\n\n")
	output.WriteString("_cueson() {\n")
	output.WriteString("  local command='' previous='' word\n")
	output.WriteString("  local -i options_active=1 index\n")
	output.WriteString("  local -a candidates\n")
	output.WriteString("  if (( CURRENT > 1 )); then previous=$words[CURRENT-1]; fi\n")
	output.WriteString("  for (( index=2; index<CURRENT; index++ )); do\n")
	output.WriteString("    word=$words[index]\n")
	output.WriteString("    if [[ $word == '--' ]]; then options_active=0; continue; fi\n")
	output.WriteString("    if [[ -z $command ]]; then\n")
	output.WriteString("      case $word in\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("        ")
		output.WriteString(command.Name)
		output.WriteString(") command=")
		output.WriteString(zshQuote(command.Name))
		output.WriteString(" ;;\n")
	}
	output.WriteString("      esac\n")
	output.WriteString("    fi\n")
	output.WriteString("  done\n")
	output.WriteString("  if (( ! options_active )); then return 0; fi\n\n")
	writeZshValueCases(&output)
	output.WriteString("  case $command in\n")
	output.WriteString("    '') candidates=(")
	for _, command := range orderedCommandSurface {
		output.WriteByte(' ')
		output.WriteString(zshQuote(command.Name))
	}
	for _, option := range globalOptionSurface {
		for _, spelling := range option.Spellings {
			output.WriteByte(' ')
			output.WriteString(zshQuote(spelling))
		}
	}
	output.WriteString(" ) ;;\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("    ")
		output.WriteString(command.Name)
		output.WriteString(") candidates=(")
		for _, value := range command.OperandValues {
			output.WriteByte(' ')
			output.WriteString(zshQuote(value))
		}
		for _, option := range appendOptions(globalOptionSurface, command.Options) {
			for _, spelling := range option.Spellings {
				output.WriteByte(' ')
				output.WriteString(zshQuote(spelling))
			}
		}
		output.WriteString(" ) ;;\n")
	}
	output.WriteString("    *) candidates=() ;;\n")
	output.WriteString("  esac\n")
	output.WriteString("  compadd -Q -- $candidates\n")
	output.WriteString("}\n\n")
	output.WriteString("compdef _cueson cueson\n")
	return output.String()
}

func writeZshValueCases(output *strings.Builder) {
	output.WriteString("  case \"$command:$previous\" in\n")
	for _, command := range orderedCommandSurface {
		for _, option := range command.Options {
			if len(option.Values) == 0 {
				continue
			}
			for _, spelling := range option.Spellings {
				output.WriteString("    ")
				output.WriteString(command.Name)
				output.WriteByte(':')
				output.WriteString(spelling)
				output.WriteString(") candidates=(")
				for _, value := range option.Values {
					output.WriteByte(' ')
					output.WriteString(zshQuote(value))
				}
				output.WriteString(" ); compadd -Q -- $candidates; return 0 ;;\n")
			}
		}
	}
	output.WriteString("  esac\n")
	for _, command := range orderedCommandSurface {
		for _, option := range command.Options {
			if len(option.Values) == 0 {
				continue
			}
			for _, spelling := range option.Spellings {
				if !strings.HasPrefix(spelling, "--") {
					continue
				}
				output.WriteString("  if [[ $command == ")
				output.WriteString(zshQuote(command.Name))
				output.WriteString(" && $words[CURRENT] == ")
				output.WriteString(zshQuote(spelling + "="))
				output.WriteString("* ]]; then candidates=(")
				for _, value := range option.Values {
					output.WriteByte(' ')
					output.WriteString(zshQuote(spelling + "=" + value))
				}
				output.WriteString(" ); compadd -Q -- $candidates; return 0; fi\n")
			}
		}
	}
}

func fishCompletionScript() string {
	var output strings.Builder
	output.WriteString("# Fish completion for cueson.\n")
	output.WriteString("function __fish_cueson_options_active\n")
	output.WriteString("  set -l tokens (commandline -opc)\n")
	output.WriteString("  not contains -- -- $tokens[2..-1]\n")
	output.WriteString("end\n\n")
	output.WriteString("function __fish_cueson_command_is\n")
	output.WriteString("  set -l tokens (commandline -opc)\n")
	output.WriteString("  for token in $tokens[2..-1]\n")
	output.WriteString("    switch $token\n")
	output.WriteString("      case")
	for _, command := range orderedCommandSurface {
		output.WriteByte(' ')
		output.WriteString(fishQuote(command.Name))
	}
	output.WriteString("\n")
	output.WriteString("        test \"$token\" = \"$argv[1]\"\n")
	output.WriteString("        return\n")
	output.WriteString("    end\n")
	output.WriteString("  end\n")
	output.WriteString("  return 1\n")
	output.WriteString("end\n\n")
	output.WriteString("complete -c cueson -f\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("complete -c cueson -f -n '__fish_cueson_options_active; and __fish_use_subcommand' -a ")
		output.WriteString(fishQuote(command.Name))
		output.WriteString(" -d ")
		output.WriteString(fishQuote(command.Summary))
		output.WriteByte('\n')
	}
	for _, option := range globalOptionSurface {
		writeFishOption(&output, "", option)
	}
	for _, command := range orderedCommandSurface {
		for _, option := range command.Options {
			writeFishOption(&output, command.Name, option)
		}
		if len(command.OperandValues) > 0 {
			output.WriteString("complete -c cueson -f -n ")
			output.WriteString(fishQuote("__fish_cueson_options_active; and __fish_cueson_command_is " + command.Name))
			output.WriteString(" -a ")
			output.WriteString(fishQuote(strings.Join(command.OperandValues, " ")))
			output.WriteString(" -d 'Supported shell'\n")
		}
	}
	return output.String()
}

func writeFishOption(output *strings.Builder, command string, option cliOptionSpec) {
	if len(option.Spellings) == 1 && option.Spellings[0] == "--" {
		output.WriteString("complete -c cueson -f -n '__fish_cueson_options_active' -a ")
		output.WriteString(fishQuote("--"))
		output.WriteString(" -d ")
		output.WriteString(fishQuote(option.Description))
		output.WriteByte('\n')
		return
	}
	output.WriteString("complete -c cueson -f -n ")
	if command != "" {
		output.WriteString(fishQuote("__fish_cueson_options_active; and __fish_cueson_command_is " + command))
	} else {
		output.WriteString(fishQuote("__fish_cueson_options_active"))
	}
	for _, spelling := range option.Spellings {
		switch {
		case strings.HasPrefix(spelling, "--"):
			output.WriteString(" -l ")
			output.WriteString(strings.TrimPrefix(spelling, "--"))
		case strings.HasPrefix(spelling, "-"):
			output.WriteString(" -s ")
			output.WriteString(strings.TrimPrefix(spelling, "-"))
		}
	}
	if option.ValueName != "" {
		output.WriteString(" -x")
	}
	if len(option.Values) > 0 {
		output.WriteString(" -a ")
		output.WriteString(fishQuote(strings.Join(option.Values, " ")))
	}
	output.WriteString(" -d ")
	output.WriteString(fishQuote(option.Description))
	output.WriteByte('\n')
}

func powershellCompletionScript() string {
	var output strings.Builder
	output.WriteString("# PowerShell completion for cueson.\n")
	output.WriteString("Register-ArgumentCompleter -Native -CommandName cueson -ScriptBlock {\n")
	output.WriteString("  param($wordToComplete, $commandAst, $cursorPosition)\n")
	output.WriteString("  $tokens = @()\n")
	output.WriteString("  foreach ($element in $commandAst.CommandElements) {\n")
	output.WriteString("    if ($element -is [System.Management.Automation.Language.StringConstantExpressionAst]) { $tokens += $element.Value } else { $tokens += $element.Extent.Text }\n")
	output.WriteString("  }\n")
	output.WriteString("  $command = ''\n")
	output.WriteString("  $optionsActive = $true\n")
	output.WriteString("  $limit = $tokens.Count\n")
	output.WriteString("  if ($wordToComplete.Length -gt 0 -and $tokens.Count -gt 1 -and $tokens[$tokens.Count - 1] -ceq $wordToComplete) { $limit-- }\n")
	output.WriteString("  for ($index = 1; $index -lt $limit; $index++) {\n")
	output.WriteString("    $token = $tokens[$index]\n")
	output.WriteString("    if ($token -ceq '--') { $optionsActive = $false; continue }\n")
	output.WriteString("    if ($command.Length -eq 0) {\n")
	output.WriteString("      switch -CaseSensitive ($token) {\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("        ")
		output.WriteString(powershellQuote(command.Name))
		output.WriteString(" { $command = ")
		output.WriteString(powershellQuote(command.Name))
		output.WriteString("; break }\n")
	}
	output.WriteString("      }\n")
	output.WriteString("    }\n")
	output.WriteString("  }\n")
	output.WriteString("  if (-not $optionsActive) { return }\n")
	output.WriteString("  $previous = if ($limit -gt 1) { $tokens[$limit - 1] } else { '' }\n")
	output.WriteString("  $candidates = @()\n")
	writePowerShellValueCases(&output)
	output.WriteString("  if ($candidates.Count -eq 0) {\n")
	output.WriteString("    switch -CaseSensitive ($command) {\n")
	output.WriteString("      '' { $candidates = @(")
	writePowerShellCandidates(&output, rootCompletionCandidates())
	output.WriteString(") }\n")
	for _, command := range orderedCommandSurface {
		output.WriteString("      ")
		output.WriteString(powershellQuote(command.Name))
		output.WriteString(" { $candidates = @(")
		writePowerShellCandidates(&output, commandCompletionCandidates(command))
		output.WriteString(") }\n")
	}
	output.WriteString("    }\n")
	output.WriteString("  }\n")
	output.WriteString("  foreach ($candidate in $candidates) {\n")
	output.WriteString("    if ($candidate.StartsWith($wordToComplete, [System.StringComparison]::Ordinal)) {\n")
	output.WriteString("      [System.Management.Automation.CompletionResult]::new($candidate, $candidate, 'ParameterValue', $candidate)\n")
	output.WriteString("    }\n")
	output.WriteString("  }\n")
	output.WriteString("}\n")
	return output.String()
}

func writePowerShellValueCases(output *strings.Builder) {
	for _, command := range orderedCommandSurface {
		for _, option := range command.Options {
			if len(option.Values) == 0 {
				continue
			}
			for _, spelling := range option.Spellings {
				output.WriteString("  if ($command -ceq ")
				output.WriteString(powershellQuote(command.Name))
				output.WriteString(" -and $previous -ceq ")
				output.WriteString(powershellQuote(spelling))
				output.WriteString(") { $candidates = @(")
				writePowerShellCandidates(output, option.Values)
				output.WriteString(") }\n")
				if strings.HasPrefix(spelling, "--") {
					output.WriteString("  if ($command -ceq ")
					output.WriteString(powershellQuote(command.Name))
					output.WriteString(" -and $wordToComplete.StartsWith(")
					output.WriteString(powershellQuote(spelling + "="))
					output.WriteString(", [System.StringComparison]::Ordinal)) { $candidates = @(")
					prefixed := make([]string, 0, len(option.Values))
					for _, value := range option.Values {
						prefixed = append(prefixed, spelling+"="+value)
					}
					writePowerShellCandidates(output, prefixed)
					output.WriteString(") }\n")
				}
			}
		}
	}
}

func rootCompletionCandidates() []string {
	values := commandNames()
	for _, option := range globalOptionSurface {
		values = append(values, option.Spellings...)
	}
	return values
}

func commandCompletionCandidates(command cliCommandSpec) []string {
	values := append([]string{}, command.OperandValues...)
	for _, option := range appendOptions(globalOptionSurface, command.Options) {
		values = append(values, option.Spellings...)
	}
	return values
}

func appendOptions(first, second []cliOptionSpec) []cliOptionSpec {
	combined := make([]cliOptionSpec, 0, len(first)+len(second))
	combined = append(combined, first...)
	combined = append(combined, second...)
	return combined
}

func writePowerShellCandidates(output *strings.Builder, candidates []string) {
	for index, candidate := range candidates {
		if index > 0 {
			output.WriteString(", ")
		}
		output.WriteString(powershellQuote(candidate))
	}
}

func bashQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func zshQuote(value string) string {
	return bashQuote(value)
}

func fishQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "\\'") + "'"
}

func powershellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
