package convert

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shruggietech/cueson/internal/model"
)

// analyzeTextToScripted shares source-feature loss accounting with the
// established text conversion paths, while native text encoding is target-aware.
func analyzeTextToScripted(ctx context.Context, document model.Document, targetFormat string) (matrixAnalysis, error) {
	analysis := matrixAnalysis{Translations: make([]payloadTranslation, len(document.Cues)), Losses: []Loss{}}
	if err := appendMetadataLosses(&analysis.Losses, document, targetFormat); err != nil {
		return matrixAnalysis{}, err
	}
	switch document.Format {
	case "subrip":
		for index, diagnostic := range document.Diagnostics {
			if err := ctx.Err(); err != nil {
				return matrixAnalysis{}, err
			}
			if diagnostic.Code == "subrip_block_unrecognized" {
				if err := appendOneLoss(&analysis.Losses, newLoss(document, targetFormat, LossCodeSubRipUnrecognizedBlockOmitted, KindOmitted, "unrecognized SubRip block is omitted", fmt.Sprintf("/diagnostics/%d", index), diagnostic.SourceOrder, nil, []Attribute{{Name: "occurrence", Value: strconv.Itoa(index)}})); err != nil {
					return matrixAnalysis{}, err
				}
			}
		}
	case "webvtt":
		if err := appendWebVTTDocumentLosses(&analysis.Losses, document, targetFormat); err != nil {
			return matrixAnalysis{}, err
		}
	default:
		return matrixAnalysis{}, &UnsupportedPairError{SourceFormat: document.Format, TargetFormat: targetFormat}
	}
	for index := range document.Cues {
		if err := ctx.Err(); err != nil {
			return matrixAnalysis{}, err
		}
		translation, err := translateTextToScripted(document, index)
		if err != nil {
			return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, fmt.Sprintf("/cues/%d/payload/raw_text", index), err)
		}
		analysis.Translations[index] = translation
		cue := &document.Cues[index]
		if document.Format == "subrip" {
			err = appendSubRipCueLosses(&analysis.Losses, document, targetFormat, index, cue, translation)
		} else {
			err = appendWebVTTCueLosses(&analysis.Losses, document, targetFormat, index, cue, translation)
		}
		if err != nil {
			return matrixAnalysis{}, err
		}
	}
	target, losses, err := projectScriptedTarget(ctx, document, targetFormat, analysis.Translations)
	if err != nil {
		return matrixAnalysis{}, projectionFailure(document.Format, targetFormat, "/format_data", err)
	}
	for _, loss := range losses {
		if err := appendOneLoss(&analysis.Losses, loss); err != nil {
			return matrixAnalysis{}, err
		}
	}
	analysis.Target = &target
	return analysis, nil
}
