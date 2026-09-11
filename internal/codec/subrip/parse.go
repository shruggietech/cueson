package subrip

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

var integerSequencePattern = regexp.MustCompile(`^[ \t]*\d+[ \t]*$`)

type physicalLine struct {
	text   string
	number int
}

// Parse converts decoded SubRip text into ordered common-model cues.
func Parse(input string, options Options) (Result, error) {
	lines := splitPhysicalLines(input)
	result := Result{Cues: make([]model.Cue, 0), Diagnostics: make([]Diagnostic, 0)}
	position := 0
	sourceOrder := 0
	expectedSequence := int64(1)

	for position < len(lines) {
		for position < len(lines) && isBlank(lines[position].text) {
			position++
		}
		if position >= len(lines) {
			break
		}

		sequenceIndex, timingIndex, found := cueStart(lines, position)
		if !found {
			end := position + 1
			for end < len(lines) && !isBlank(lines[end].text) && !validCueStart(lines, end) {
				end++
			}
			order := sourceOrder
			result.Diagnostics = append(result.Diagnostics, Diagnostic{
				Severity: "warning", Code: "subrip_block_unrecognized",
				Message:     fmt.Sprintf("unrecognized SubRip content beginning at line %d remains available in the source envelope", lines[position].number),
				SourceOrder: &order,
			})
			sourceOrder++
			position = end
			continue
		}

		timing, err := ParseTimecodeLine(lines[timingIndex].text)
		if err != nil {
			return Result{}, &ParseError{Code: "subrip_timing_invalid", Line: lines[timingIndex].number, Text: err.Error()}
		}

		payloadStart := timingIndex + 1
		payloadEnd := payloadStart
		nextPosition := len(lines)
		missingSeparator := false
		for payloadEnd < len(lines) {
			if isBlank(lines[payloadEnd].text) && separatorBoundary(lines, payloadEnd) {
				nextPosition = payloadEnd
				for nextPosition < len(lines) && isBlank(lines[nextPosition].text) {
					nextPosition++
				}
				break
			}
			if payloadEnd > payloadStart && validCueStart(lines, payloadEnd) {
				nextPosition = payloadEnd
				missingSeparator = true
				break
			}
			payloadEnd++
		}
		payloadLines := make([]string, 0, payloadEnd-payloadStart)
		for index := payloadStart; index < payloadEnd; index++ {
			payloadLines = append(payloadLines, lines[index].text)
		}

		ordinal := len(result.Cues)
		cueID := fmt.Sprintf("cue-%06d", ordinal)
		order := sourceOrder
		sequenceRaw := ""
		var sourceIdentifier *string
		if sequenceIndex >= 0 {
			sequenceRaw = lines[sequenceIndex].text
			identifier := strings.TrimSpace(sequenceRaw)
			sourceIdentifier = &identifier
			if integerSequencePattern.MatchString(sequenceRaw) {
				value, parseErr := strconv.ParseInt(identifier, 10, 64)
				if parseErr != nil {
					result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_sequence_irregular", "sequence value exceeds the supported integer range and was preserved literally", order, cueID))
				} else {
					if value != expectedSequence {
						result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_sequence_irregular", fmt.Sprintf("sequence value %d was preserved; expected %d", value, expectedSequence), order, cueID))
					}
					if value < int64(^uint64(0)>>1) {
						expectedSequence = value + 1
					}
				}
			} else {
				result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_sequence_non_integer", "non-integer sequence line was preserved literally", order, cueID))
			}
		} else {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_sequence_missing", "cue has no sequence line; the missing distinction was preserved", order, cueID))
			expectedSequence++
		}
		if missingSeparator {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_separator_missing", "cue boundary was recovered without a blank separator", order, cueID))
		}
		if timing.Separator == '.' || timing.EndSeparator == '.' {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_timing_period_separator", "period millisecond separators were accepted as a tolerated variant", order, cueID))
		}
		if timing.StartFractionDigits < 3 || timing.EndFractionDigits < 3 {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_timing_short_fraction", "short millisecond fields were right-padded to three digits", order, cueID))
		}
		if len(payloadLines) == 0 {
			result.Diagnostics = append(result.Diagnostics, cueDiagnostic("subrip_payload_empty", "cue has no payload lines", order, cueID))
		}

		rawText := strings.Join(payloadLines, "\n")
		cue := model.Cue{
			ID: cueID, Ordinal: ordinal, SourceOrder: order, SourceIdentifier: sourceIdentifier,
			Timing:   model.Timing{StartMilliseconds: timing.StartMilliseconds, EndMilliseconds: timing.EndMilliseconds, DurationMilliseconds: timing.EndMilliseconds - timing.StartMilliseconds},
			Payload:  model.Payload{RawText: rawText, PlainText: PlainText(rawText), Lines: payloadLines},
			Speakers: []model.Speaker{}, Tokens: []model.Token{}, OCRObservations: []model.OCRObservation{},
			FormatData: model.CueFormatData{SubRip: &model.SubRipCueData{SequenceLineRaw: sequenceRaw, TimingLineRaw: lines[timingIndex].text, Coordinates: timing.Coordinates}},
		}
		if options.DetectSpeakers {
			if speaker := DetectSpeaker(payloadLines); speaker != nil {
				cue.Speakers = append(cue.Speakers, model.Speaker{Name: *speaker, Origin: "heuristic"})
			}
		}
		result.Cues = append(result.Cues, cue)
		sourceOrder++
		position = nextPosition
	}

	if len(result.Cues) == 0 {
		return Result{}, &ParseError{Code: "subrip_no_cues", Text: "input contains no valid SubRip cues"}
	}
	return result, nil
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

func cueStart(lines []physicalLine, position int) (int, int, bool) {
	if _, err := ParseTimecodeLine(lines[position].text); err == nil {
		return -1, position, true
	}
	if strings.Contains(lines[position].text, "-->") {
		return -1, position, true
	}
	if position+1 < len(lines) {
		if _, err := ParseTimecodeLine(lines[position+1].text); err == nil || strings.Contains(lines[position+1].text, "-->") {
			return position, position + 1, true
		}
	}
	return -1, -1, false
}

func validCueStart(lines []physicalLine, position int) bool {
	if position >= len(lines) {
		return false
	}
	if _, err := ParseTimecodeLine(lines[position].text); err == nil {
		return true
	}
	if position+1 < len(lines) {
		_, err := ParseTimecodeLine(lines[position+1].text)
		return err == nil
	}
	return false
}

func separatorBoundary(lines []physicalLine, position int) bool {
	for position < len(lines) && isBlank(lines[position].text) {
		position++
	}
	return position >= len(lines) || validCueStart(lines, position)
}

func isBlank(line string) bool {
	return strings.TrimSpace(line) == ""
}

func cueDiagnostic(code, message string, sourceOrder int, cueID string) Diagnostic {
	order := sourceOrder
	id := cueID
	return Diagnostic{Severity: "warning", Code: code, Message: message, SourceOrder: &order, CueID: &id}
}
