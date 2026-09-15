package convert

import (
	"context"
	"fmt"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/source"
)

// Convert projects one validated source document into canonical target bytes.
// It has no filesystem or diagnostic-stream side effects.
func Convert(ctx context.Context, document model.Document, targetFormat string, options Options) (Result, error) {
	return convertWithRenderer(ctx, document, targetFormat, options, renderProjected)
}

func convertWithRenderer(ctx context.Context, document model.Document, targetFormat string, options Options, renderer targetRenderer) (Result, error) {
	if ctx == nil {
		return Result{}, fmt.Errorf("convert: invalid_context")
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if targetFormat == document.Format {
		return Result{}, &SameFormatError{Format: document.Format}
	}
	if !supportedPair(document.Format, targetFormat) {
		return Result{}, &UnsupportedPairError{SourceFormat: document.Format, TargetFormat: targetFormat}
	}
	if err := document.Validate(); err != nil {
		return Result{}, fmt.Errorf("validate conversion source: %w", err)
	}
	if err := source.ValidateIntegrity(ctx, document); err != nil {
		if cancelled := ctx.Err(); cancelled != nil {
			return Result{}, cancelled
		}
		return Result{}, fmt.Errorf("convert: invalid_source_integrity")
	}
	analysis, err := analyzeCompatibilityContext(ctx, document, targetFormat)
	if err != nil {
		return Result{}, err
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		return Result{}, fmt.Errorf("build conversion loss report: %w", err)
	}
	if options.Strict && report.HasLosses() {
		return Result{LossReport: report}, &StrictLossError{Report: report}
	}
	var target model.Document
	if analysis.Target != nil {
		target = *analysis.Target
	} else {
		target, err = projectDocument(document, targetFormat, analysis.Translations)
	}
	if err != nil {
		return Result{LossReport: report}, projectionFailure(document.Format, targetFormat, "format_data", err)
	}
	bytes, diagnostics, err := renderer(ctx, document, target, targetFormat)
	if err != nil {
		return Result{LossReport: report}, projectionFailure(document.Format, targetFormat, "format_data", err)
	}
	return Result{Bytes: bytes, LossReport: report, Diagnostics: diagnostics}, nil
}

func supportedPair(sourceFormat, targetFormat string) bool {
	return supportedFormat(sourceFormat) && supportedFormat(targetFormat) && sourceFormat != targetFormat
}
