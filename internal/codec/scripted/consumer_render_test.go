package scripted

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestRenderReportsConsumerAnnotationsForBothScriptedDialects(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			document := renderFixture(t, format)
			document.Schema, document.SchemaVersion = schema.ID(), schema.Version()
			baseline, err := Render(context.Background(), document, false)
			if err != nil {
				t.Fatal(err)
			}
			document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "PRIVATE-IDENTIFIER"}}
			document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 5000}
			before, _ := json.Marshal(document)
			result, err := Render(context.Background(), document, false)
			if err != nil || !bytes.Equal(baseline.Bytes, result.Bytes) {
				t.Fatalf("consumer annotations changed native bytes: %q, err = %v", result.Bytes, err)
			}
			codes := map[string]bool{}
			for _, diagnostic := range result.Diagnostics {
				codes[diagnostic.Code] = true
				if strings.Contains(diagnostic.Message, "PRIVATE-IDENTIFIER") {
					t.Fatal("consumer identifier leaked into render warning")
				}
			}
			if !codes["consumer_speaker_attribution_omitted"] || !codes["consumer_media_timing_omitted"] {
				t.Fatalf("consumer omission diagnostics = %#v", result.Diagnostics)
			}
			strict, err := Render(context.Background(), document, true)
			if err == nil || len(strict.Bytes) != 0 || !strings.Contains(err.Error(), "consumer_") {
				t.Fatalf("strict render = %#v, err = %v", strict, err)
			}
			after, _ := json.Marshal(document)
			if !bytes.Equal(before, after) {
				t.Fatal("render changed consumer annotations or native source fields")
			}
		})
	}
}
