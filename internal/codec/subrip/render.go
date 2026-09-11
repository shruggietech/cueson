package subrip

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

// Render serializes common-model cues as canonical LF SubRip bytes.
func Render(cues []model.Cue) ([]byte, error) {
	if len(cues) == 0 {
		return nil, fmt.Errorf("render SubRip: at least one cue is required")
	}
	blocks := make([]string, 0, len(cues))
	for index := range cues {
		cue := &cues[index]
		if cue.Timing.StartMilliseconds < 0 || cue.Timing.EndMilliseconds < cue.Timing.StartMilliseconds {
			return nil, fmt.Errorf("render SubRip cue %d: invalid timing range", index)
		}
		if cue.Timing.DurationMilliseconds != cue.Timing.EndMilliseconds-cue.Timing.StartMilliseconds {
			return nil, fmt.Errorf("render SubRip cue %d: inconsistent duration", index)
		}
		if strings.ContainsRune(cue.Payload.RawText, '\r') {
			return nil, fmt.Errorf("render SubRip cue %d: raw_text contains a non-canonical carriage return", index)
		}
		start, err := formatTimecode(cue.Timing.StartMilliseconds)
		if err != nil {
			return nil, fmt.Errorf("render SubRip cue %d start: %w", index, err)
		}
		end, err := formatTimecode(cue.Timing.EndMilliseconds)
		if err != nil {
			return nil, fmt.Errorf("render SubRip cue %d end: %w", index, err)
		}
		timingLine := start + " --> " + end
		if cue.FormatData.SubRip != nil && cue.FormatData.SubRip.Coordinates != nil {
			coordinates := cue.FormatData.SubRip.Coordinates
			if coordinates.X1 < 0 || coordinates.X2 < 0 || coordinates.Y1 < 0 || coordinates.Y2 < 0 {
				return nil, fmt.Errorf("render SubRip cue %d: coordinates must be non-negative", index)
			}
			timingLine += fmt.Sprintf(" X1:%d X2:%d Y1:%d Y2:%d", coordinates.X1, coordinates.X2, coordinates.Y1, coordinates.Y2)
		}
		block := strconv.Itoa(index+1) + "\n" + timingLine + "\n" + cue.Payload.RawText
		blocks = append(blocks, block)
	}
	return []byte(strings.Join(blocks, "\n\n") + "\n"), nil
}

// PayloadHasAmbiguousBoundary reports whether canonical emission would let the
// SubRip parser reinterpret part of raw_text as another cue or a separator.
func PayloadHasAmbiguousBoundary(rawText string) bool {
	lines := splitPhysicalLines(rawText)
	if strings.HasSuffix(rawText, "\n") {
		lines = append(lines, physicalLine{text: "", number: len(lines) + 1})
	}
	for position := range lines {
		if isBlank(lines[position].text) && separatorBoundary(lines, position) {
			return true
		}
		if position > 0 && validCueStart(lines, position) {
			return true
		}
	}
	return false
}
