package convert

import (
	"context"
	"fmt"
	"strconv"

	"github.com/shruggietech/cueson/internal/model"
)

// Consumer annotations have no native subtitle representation. Account for
// each omitted occurrence using its location, never its identifier contents.
func appendConsumerAnnotationLosses(ctx context.Context, losses *[]Loss, document model.Document, targetFormat string) error {
	if ctx == nil {
		return fmt.Errorf("conversion consumer annotations: invalid_context")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if document.MediaTiming != nil {
		if err := appendOneLoss(losses, newLoss(document, targetFormat, LossCodeMediaTimingOmitted, KindOmitted, "consumer supplied media timing is omitted from subtitle output", "/media_timing", nil, nil, nil)); err != nil {
			return err
		}
	}
	for cueIndex, cue := range document.Cues {
		if err := ctx.Err(); err != nil {
			return err
		}
		for occurrence := range cue.SpeakerAttributions {
			loss := cueLoss(document, targetFormat, cue, LossCodeSpeakerAttributionOmitted, KindOmitted, "consumer speaker attribution is omitted from subtitle output", fmt.Sprintf("/cues/%d/speaker_attributions/%d", cueIndex, occurrence), []Attribute{{Name: "occurrence", Value: strconv.Itoa(occurrence)}})
			if err := appendOneLoss(losses, loss); err != nil {
				return err
			}
		}
	}
	return nil
}
