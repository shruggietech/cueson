// Package subrip parses and renders the documented Cueson SubRip dialect.
package subrip

import (
	"fmt"

	"github.com/shruggietech/cueson/internal/model"
)

// Coordinates is the complete optional SubRip pixel bounding box.
type Coordinates = model.Coordinates

// Diagnostic is an ordered, non-fatal interpretation observation.
type Diagnostic = model.Diagnostic

// Options controls optional semantic derivation.
type Options struct {
	DetectSpeakers bool
}

// Result contains common-model cues and every non-fatal parser diagnostic.
type Result struct {
	Cues        []model.Cue
	Diagnostics []Diagnostic
}

// TimecodeLine is one parsed SubRip timing line.
type TimecodeLine struct {
	StartMilliseconds   int64
	EndMilliseconds     int64
	Coordinates         *Coordinates
	Separator           byte
	EndSeparator        byte
	StartFractionDigits int
	EndFractionDigits   int
}

// ParseError identifies a deterministic fatal grammar boundary.
type ParseError struct {
	Code string
	Line int
	Text string
}

func (err *ParseError) Error() string {
	if err.Line > 0 {
		return fmt.Sprintf("%s at line %d: %s", err.Code, err.Line, err.Text)
	}
	return fmt.Sprintf("%s: %s", err.Code, err.Text)
}
