package webvtt

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ParseTimestamp parses one conforming WebVTT timestamp.
func ParseTimestamp(value string) (int64, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, fmt.Errorf("timestamp must use MM:SS.mmm or HH:MM:SS.mmm")
	}
	hoursText := "0"
	minutesText := parts[0]
	secondsText := parts[1]
	if len(parts) == 3 {
		hoursText, minutesText, secondsText = parts[0], parts[1], parts[2]
		if len(hoursText) < 2 || !asciiDigits(hoursText) {
			return 0, fmt.Errorf("hours must contain at least two ASCII digits")
		}
	}
	if len(minutesText) != 2 || !asciiDigits(minutesText) {
		return 0, fmt.Errorf("minutes must contain exactly two ASCII digits")
	}
	secondParts := strings.Split(secondsText, ".")
	if len(secondParts) != 2 || len(secondParts[0]) != 2 || len(secondParts[1]) != 3 || !asciiDigits(secondParts[0]) || !asciiDigits(secondParts[1]) {
		return 0, fmt.Errorf("seconds must use two digits and three millisecond digits")
	}
	hours, err := strconv.ParseInt(hoursText, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("hours overflow")
	}
	minutes, _ := strconv.ParseInt(minutesText, 10, 64)
	seconds, _ := strconv.ParseInt(secondParts[0], 10, 64)
	milliseconds, _ := strconv.ParseInt(secondParts[1], 10, 64)
	if minutes > 59 || seconds > 59 {
		return 0, fmt.Errorf("minutes and seconds must be less than 60")
	}
	base := minutes*60_000 + seconds*1_000 + milliseconds
	if hours > (math.MaxInt64-base)/3_600_000 {
		return 0, fmt.Errorf("timestamp exceeds millisecond range")
	}
	return hours*3_600_000 + base, nil
}

// ParseTimingLine parses one conforming WebVTT cue timing line.
func ParseTimingLine(line string) (TimingLine, error) {
	arrow := strings.Index(line, "-->")
	if arrow < 0 || strings.Contains(line[arrow+3:], "-->") {
		return TimingLine{}, fmt.Errorf("timing line must contain one --> delimiter")
	}
	left := line[:arrow]
	right := line[arrow+3:]
	if len(left) == 0 || len(right) == 0 || !isASCIISpace(left[len(left)-1]) || !isASCIISpace(right[0]) {
		return TimingLine{}, fmt.Errorf("timing delimiter requires ASCII whitespace on both sides")
	}
	startText := asciiTrim(left)
	right = strings.TrimLeft(right, " \t")
	endBoundary := strings.IndexAny(right, " \t")
	endText := right
	settingsRaw := ""
	if endBoundary >= 0 {
		endText = right[:endBoundary]
		settingsRaw = asciiTrim(right[endBoundary:])
	}
	start, err := ParseTimestamp(startText)
	if err != nil {
		return TimingLine{}, fmt.Errorf("start timestamp: %w", err)
	}
	end, err := ParseTimestamp(endText)
	if err != nil {
		return TimingLine{}, fmt.Errorf("end timestamp: %w", err)
	}
	if end <= start {
		return TimingLine{}, fmt.Errorf("end timestamp must be greater than start timestamp")
	}
	return TimingLine{StartMilliseconds: start, EndMilliseconds: end, SettingsRaw: settingsRaw}, nil
}

func formatTimestamp(milliseconds int64) (string, error) {
	if milliseconds < 0 {
		return "", fmt.Errorf("timestamp must be non-negative")
	}
	hours := milliseconds / 3_600_000
	remainder := milliseconds % 3_600_000
	minutes := remainder / 60_000
	remainder %= 60_000
	seconds := remainder / 1_000
	fraction := remainder % 1_000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, fraction), nil
}

func asciiDigits(value string) bool {
	if value == "" {
		return false
	}
	for index := range value {
		if value[index] < '0' || value[index] > '9' {
			return false
		}
	}
	return true
}

func isASCIISpace(value byte) bool {
	return value == ' ' || value == '\t'
}
