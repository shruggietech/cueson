package webvtt

import "strings"

type physicalLine struct {
	text   string
	number int
}

func splitPhysicalLines(input string) []physicalLine {
	lines := make([]physicalLine, 0, strings.Count(input, "\n")+1)
	start := 0
	lineNumber := 1
	for index := 0; index < len(input); index++ {
		if input[index] != '\r' && input[index] != '\n' {
			continue
		}
		lines = append(lines, physicalLine{text: input[start:index], number: lineNumber})
		if input[index] == '\r' && index+1 < len(input) && input[index+1] == '\n' {
			index++
		}
		start = index + 1
		lineNumber++
	}
	if start < len(input) {
		lines = append(lines, physicalLine{text: input[start:], number: lineNumber})
	}
	return lines
}

func blank(line string) bool {
	return line == ""
}

func asciiTrim(value string) string {
	return strings.Trim(value, " \t")
}
