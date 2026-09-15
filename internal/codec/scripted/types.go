// Package scripted implements the bounded ASS v4+/SSA v4 native text profile.
package scripted

import (
	"fmt"
	"github.com/shruggietech/cueson/internal/model"
)

type ParseResult struct {
	Native      model.ScriptedDocumentData
	Cues        []model.Cue
	Diagnostics []model.Diagnostic
	Encoding    model.EncodingObservation
}

// ParseError reports only a closed code and safe physical position.
type ParseError struct {
	Code        string
	SourceOrder int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s at scripted source_order %d", e.Code, e.SourceOrder)
}
func failure(code string, order int) error     { return &ParseError{Code: code, SourceOrder: order} }
func ptr[T any](v T) *T                        { return &v }
func occurrence(kind string, order int) string { return fmt.Sprintf("%s-%06d", kind, order) }

type Detection struct {
	Candidate bool
	Format    string
	Err       error
}
