package scripted

import (
	"context"
	"fmt"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/source"
)

// RenderTarget serializes a private conversion target beside validated source
// truth. Generated native bytes never replace its original source envelope.
func RenderTarget(ctx context.Context, original, target model.Document) (RenderResult, error) {
	if ctx == nil {
		return RenderResult{}, fmt.Errorf("render scripted target: invalid_context")
	}
	if err := ctx.Err(); err != nil {
		return RenderResult{}, err
	}
	if err := model.ValidateScriptedTarget(original, target); err != nil {
		return RenderResult{}, err
	}
	if err := source.ValidateIntegrity(ctx, original); err != nil {
		if cancelled := ctx.Err(); cancelled != nil {
			return RenderResult{}, cancelled
		}
		return RenderResult{}, fmt.Errorf("render scripted target: invalid_source_integrity")
	}
	return renderValidated(ctx, target, false)
}
