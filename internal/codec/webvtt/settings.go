package webvtt

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

var (
	signedNumber      = regexp.MustCompile(`^-?[0-9]+$`)
	percentagePattern = regexp.MustCompile(`^[0-9]+(?:\.[0-9]+)?%$`)
)

var cueSettingOrder = []string{"region", "vertical", "line", "position", "size", "align"}
var regionSettingOrder = []string{"id", "width", "lines", "regionanchor", "viewportanchor", "scroll"}

func parseCueSettings(raw string, sourceOrder int, cueID string) ([]model.WebVTTSettingOccurrence, map[string]string, []Diagnostic) {
	return parseSettings(raw, cueSettingOrder, validCueSetting, sourceOrder, &cueID)
}

func parseRegionSettings(raw string, sourceOrder int) ([]model.WebVTTSettingOccurrence, map[string]string, []Diagnostic) {
	return parseSettings(raw, regionSettingOrder, validRegionSetting, sourceOrder, nil)
}

func parseSettings(raw string, recognized []string, validate func(string, string) bool, sourceOrder int, cueID *string) ([]model.WebVTTSettingOccurrence, map[string]string, []Diagnostic) {
	recognizedSet := make(map[string]struct{}, len(recognized))
	for _, name := range recognized {
		recognizedSet[name] = struct{}{}
	}
	occurrences := make([]model.WebVTTSettingOccurrence, 0)
	effective := make(map[string]string)
	diagnostics := make([]Diagnostic, 0)
	seen := make(map[string]bool)
	for _, token := range splitASCIIWhitespace(raw) {
		name, value, found := strings.Cut(token, ":")
		semanticName, semanticValue := semanticText(name), semanticText(value)
		_, known := recognizedSet[semanticName]
		valid := found && semanticName != "" && semanticValue != "" && known && validate(semanticName, semanticValue)
		occurrences = append(occurrences, model.WebVTTSettingOccurrence{Raw: token, Name: name, Value: value, Recognized: known, Valid: valid})
		switch {
		case !known:
			diagnostics = append(diagnostics, settingDiagnostic("webvtt_setting_unknown", "unknown WebVTT setting was preserved", sourceOrder, cueID))
		case !valid:
			diagnostics = append(diagnostics, settingDiagnostic("webvtt_setting_invalid", "invalid WebVTT setting was preserved", sourceOrder, cueID))
		default:
			if seen[semanticName] {
				diagnostics = append(diagnostics, settingDiagnostic("webvtt_setting_duplicate", "duplicate WebVTT setting was preserved; the last valid value is effective", sourceOrder, cueID))
			}
			seen[semanticName] = true
			effective[semanticName] = semanticValue
		}
	}
	return occurrences, effective, diagnostics
}

func validCueSetting(name, value string) bool {
	switch name {
	case "region":
		return value != "" && !strings.ContainsAny(value, " \t\r\n") && !strings.Contains(value, "-->")
	case "vertical":
		return value == "rl" || value == "lr"
	case "line":
		position, alignment, hasAlignment := strings.Cut(value, ",")
		if hasAlignment && alignment != "start" && alignment != "center" && alignment != "end" {
			return false
		}
		if strings.HasSuffix(position, "%") {
			return validPercentage(position)
		}
		return signedNumber.MatchString(position)
	case "position":
		position, alignment, hasAlignment := strings.Cut(value, ",")
		if hasAlignment && alignment != "line-left" && alignment != "center" && alignment != "line-right" {
			return false
		}
		return validPercentage(position)
	case "size":
		return validPercentage(value)
	case "align":
		return value == "start" || value == "center" || value == "end" || value == "left" || value == "right"
	default:
		return false
	}
}

func validRegionSetting(name, value string) bool {
	switch name {
	case "id":
		return value != "" && !strings.ContainsAny(value, " \t\r\n") && !strings.Contains(value, "-->")
	case "width":
		return validPercentage(value)
	case "lines":
		return asciiDigits(value)
	case "regionanchor", "viewportanchor":
		left, right, found := strings.Cut(value, ",")
		return found && validPercentage(left) && validPercentage(right)
	case "scroll":
		return value == "up"
	default:
		return false
	}
}

func validPercentage(value string) bool {
	if !percentagePattern.MatchString(value) {
		return false
	}
	number := strings.TrimSuffix(value, "%")
	parsed, err := strconv.ParseFloat(number, 64)
	return err == nil && parsed >= 0 && parsed <= 100
}

func splitASCIIWhitespace(value string) []string {
	fields := make([]string, 0)
	start := -1
	for index := 0; index <= len(value); index++ {
		separator := index == len(value) || value[index] == ' ' || value[index] == '\t' || value[index] == '\n' || value[index] == '\r'
		if separator {
			if start >= 0 {
				fields = append(fields, value[start:index])
				start = -1
			}
			continue
		}
		if start < 0 {
			start = index
		}
	}
	return fields
}

func countASCIIWhitespaceFields(value string) int {
	count := 0
	inField := false
	for index := 0; index < len(value); index++ {
		separator := value[index] == ' ' || value[index] == '\t' || value[index] == '\n' || value[index] == '\r'
		if separator {
			inField = false
			continue
		}
		if !inField {
			count++
			inField = true
		}
	}
	return count
}

func settingDiagnostic(code, message string, sourceOrder int, cueID *string) Diagnostic {
	order := sourceOrder
	diagnostic := Diagnostic{Severity: "warning", Code: code, Message: message, SourceOrder: &order}
	if cueID != nil {
		id := *cueID
		diagnostic.CueID = &id
	}
	return diagnostic
}
