package webvtt

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestRenderReportsConsumerAnnotationsWithoutChangingNativeVoices(t *testing.T) {
	parsed, err := Parse("WEBVTT\n\n00:01.000 --> 00:03.000\n<v Native speaker>Hello</v>\n")
	if err != nil {
		t.Fatal(err)
	}
	document := model.Document{Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}, Diagnostics: parsed.Diagnostics}
	document.Schema, document.SchemaVersion = schema.ID(), schema.Version()
	baseline, err := Render(document, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "PRIVATE-IDENTIFIER"}}
	document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 5000}
	before, _ := json.Marshal(document)
	result, err := Render(document, RenderOptions{})
	if err != nil || !bytes.Equal(baseline.Bytes, result.Bytes) {
		t.Fatalf("native voice representation changed: %q, err = %v", result.Bytes, err)
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
	strict, err := Render(document, RenderOptions{Strict: true})
	if err == nil || len(strict.Bytes) != 0 || !strings.Contains(err.Error(), "consumer_") {
		t.Fatalf("strict render = %#v, err = %v", strict, err)
	}
	after, _ := json.Marshal(document)
	if !bytes.Equal(before, after) {
		t.Fatal("render changed consumer annotations or native speaker observations")
	}
}

func TestRenderRejectsInvalidEditedOrHistoricalConsumerAnnotations(t *testing.T) {
	start, end := int64(1000), int64(4000)
	for _, test := range []struct {
		name        string
		attribution model.SpeakerAttribution
		historical  bool
	}{
		{name: "control", attribution: model.SpeakerAttribution{SpeakerID: "bad\nidentifier"}},
		{name: "partial timing", attribution: model.SpeakerAttribution{SpeakerID: "speaker", StartMilliseconds: &start}},
		{name: "outside cue", attribution: model.SpeakerAttribution{SpeakerID: "speaker", StartMilliseconds: &start, EndMilliseconds: &end}},
		{name: "historical contract", attribution: model.SpeakerAttribution{SpeakerID: "speaker"}, historical: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := Parse("WEBVTT\n\n00:01.000 --> 00:03.000\nHello\n")
			if err != nil {
				t.Fatal(err)
			}
			document := model.Document{Schema: schema.ID(), SchemaVersion: schema.Version(), Format: "webvtt", Cues: parsed.Cues, FormatData: model.DocumentFormatData{WebVTT: &parsed.DocumentData}}
			document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{test.attribution}
			if test.historical {
				document.Schema, document.SchemaVersion = "https://cueson.io/schema/v1.1.0/cueson.schema.json", "1.1.0"
			}
			for _, strict := range []bool{false, true} {
				result, err := Render(document, RenderOptions{Strict: strict})
				if err == nil || len(result.Bytes) != 0 {
					t.Fatalf("invalid typed annotations produced bytes: strict=%t result=%#v err=%v", strict, result, err)
				}
			}
		})
	}
}
