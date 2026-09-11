package subrip

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

const timeTokenPattern = `(\d+):(\d{1,2}):(\d{1,2})([,.])(\d{1,3})`

var (
	timingLinePattern = regexp.MustCompile(`^[ \t]*` + timeTokenPattern + `[ \t]*-->[ \t]*` + timeTokenPattern + `(?:[ \t]+(.+?))?[ \t]*$`)
	coordinatePattern = regexp.MustCompile(`(?i)(X1|X2|Y1|Y2):[ \t]*(\d+)`)
)

// ParseTimecodeLine parses a complete SubRip timecode line without normalizing its source text.
func ParseTimecodeLine(line string) (TimecodeLine, error) {
	match := timingLinePattern.FindStringSubmatch(line)
	if match == nil {
		return TimecodeLine{}, fmt.Errorf("timecode line does not match SubRip grammar")
	}
	start, err := timePartsToMilliseconds(match[1], match[2], match[3], match[5])
	if err != nil {
		return TimecodeLine{}, fmt.Errorf("start timecode: %w", err)
	}
	end, err := timePartsToMilliseconds(match[6], match[7], match[8], match[10])
	if err != nil {
		return TimecodeLine{}, fmt.Errorf("end timecode: %w", err)
	}
	if end < start {
		return TimecodeLine{}, fmt.Errorf("end timecode precedes start timecode")
	}
	var coordinates *Coordinates
	if match[11] != "" {
		coordinates, err = parseCoordinates(match[11])
		if err != nil {
			return TimecodeLine{}, err
		}
	}
	return TimecodeLine{
		StartMilliseconds:   start,
		EndMilliseconds:     end,
		Coordinates:         coordinates,
		Separator:           match[4][0],
		StartFractionDigits: len(match[5]),
		EndFractionDigits:   len(match[10]),
	}, nil
}

func timePartsToMilliseconds(hoursText, minutesText, secondsText, fractionText string) (int64, error) {
	hours, err := strconv.ParseInt(hoursText, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("hours overflow")
	}
	minutes, _ := strconv.ParseInt(minutesText, 10, 64)
	seconds, _ := strconv.ParseInt(secondsText, 10, 64)
	if minutes >= 60 || seconds >= 60 {
		return 0, fmt.Errorf("minutes and seconds must be less than 60")
	}
	fraction, _ := strconv.ParseInt(fractionText+strings.Repeat("0", 3-len(fractionText)), 10, 64)
	if hours > (math.MaxInt64-minutes*60_000-seconds*1_000-fraction)/3_600_000 {
		return 0, fmt.Errorf("timecode exceeds millisecond range")
	}
	return hours*3_600_000 + minutes*60_000 + seconds*1_000 + fraction, nil
}

func parseCoordinates(suffix string) (*Coordinates, error) {
	matches := coordinatePattern.FindAllStringSubmatchIndex(suffix, -1)
	if len(matches) != 4 {
		return nil, fmt.Errorf("coordinate suffix must contain X1, X2, Y1, and Y2 exactly once")
	}
	values := make(map[string]int, 4)
	cursor := 0
	for _, match := range matches {
		if strings.TrimSpace(suffix[cursor:match[0]]) != "" {
			return nil, fmt.Errorf("coordinate suffix contains unsupported text")
		}
		key := strings.ToUpper(suffix[match[2]:match[3]])
		if _, exists := values[key]; exists {
			return nil, fmt.Errorf("coordinate suffix repeats %s", key)
		}
		value, err := strconv.ParseInt(suffix[match[4]:match[5]], 10, strconv.IntSize)
		if err != nil {
			return nil, fmt.Errorf("coordinate %s exceeds integer range", key)
		}
		values[key] = int(value)
		cursor = match[1]
	}
	if strings.TrimSpace(suffix[cursor:]) != "" {
		return nil, fmt.Errorf("coordinate suffix contains unsupported text")
	}
	return &Coordinates{X1: values["X1"], X2: values["X2"], Y1: values["Y1"], Y2: values["Y2"]}, nil
}

func formatTimecode(milliseconds int64) (string, error) {
	if milliseconds < 0 {
		return "", fmt.Errorf("timestamp must be non-negative")
	}
	hours := milliseconds / 3_600_000
	remainder := milliseconds % 3_600_000
	minutes := remainder / 60_000
	remainder %= 60_000
	seconds := remainder / 1_000
	fraction := remainder % 1_000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, fraction), nil
}
