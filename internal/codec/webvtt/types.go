// Package webvtt parses and renders the Cueson WebVTT dialect.
package webvtt

import (
	"fmt"

	"github.com/shruggietech/cueson/internal/model"
)

// Diagnostic is an ordered, non-fatal interpretation observation.
type Diagnostic = model.Diagnostic

// Result contains every parsed cue, non-cue block, and diagnostic.
type Result struct {
	DocumentData model.WebVTTDocumentData
	Cues         []model.Cue
	Diagnostics  []Diagnostic
}

// DecodedText contains strict WebVTT Unicode text and its byte-derived observations.
type DecodedText struct {
	Text        string
	Observation model.EncodingObservation
}

// TimingLine contains the normalized semantics of one cue timing line.
type TimingLine struct {
	StartMilliseconds int64
	EndMilliseconds   int64
	SettingsRaw       string
}

// RenderOptions controls conformance handling during model-driven output.
type RenderOptions struct {
	Strict bool
}

// RenderResult contains canonical WebVTT bytes and render-time diagnostics.
type RenderResult struct {
	Bytes       []byte
	Diagnostics []Diagnostic
}

// ParseError identifies a deterministic fatal WebVTT boundary.
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
