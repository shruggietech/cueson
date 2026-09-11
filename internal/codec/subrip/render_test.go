package subrip

import (
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestRenderCanonicalSubRip(t *testing.T) {
	t.Parallel()
	cues := []model.Cue{
		{Timing: model.Timing{StartMilliseconds: 1250, EndMilliseconds: 4200, DurationMilliseconds: 2950}, Payload: model.Payload{RawText: "<i>Hello</i>\nworld"}, FormatData: model.CueFormatData{SubRip: &model.SubRipCueData{Coordinates: &model.Coordinates{X1: 1, X2: 2, Y1: 3, Y2: 4}}}},
		{Timing: model.Timing{StartMilliseconds: 3723004, EndMilliseconds: 3723500, DurationMilliseconds: 496}, Payload: model.Payload{RawText: "second"}, FormatData: model.CueFormatData{SubRip: &model.SubRipCueData{}}},
	}
	got, err := Render(cues)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want := "1\n00:00:01,250 --> 00:00:04,200 X1:1 X2:2 Y1:3 Y2:4\n<i>Hello</i>\nworld\n\n2\n01:02:03,004 --> 01:02:03,500\nsecond\n"
	if string(got) != want {
		t.Fatalf("Render() = %q, want %q", got, want)
	}
	parsed, err := Parse(string(got), Options{})
	if err != nil || len(parsed.Cues) != 2 {
		t.Fatalf("rendered parse = %#v, %v", parsed, err)
	}
}

func TestRenderRejectsInvalidStructuredCue(t *testing.T) {
	t.Parallel()
	for _, cues := range [][]model.Cue{
		nil,
		{{Timing: model.Timing{StartMilliseconds: 2, EndMilliseconds: 1, DurationMilliseconds: -1}, Payload: model.Payload{RawText: "bad"}}},
		{{Timing: model.Timing{StartMilliseconds: 0, EndMilliseconds: 1, DurationMilliseconds: 1}, Payload: model.Payload{RawText: "bad\rtext"}}},
	} {
		if _, err := Render(cues); err == nil {
			t.Errorf("Render(%#v) unexpectedly succeeded", cues)
		}
	}
}

func TestPayloadHasAmbiguousBoundary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		payload string
		want    bool
	}{
		{name: "ordinary multiline", payload: "first\nsecond", want: false},
		{name: "internal blank", payload: "first\n\nthird", want: false},
		{name: "sequence and timing suffix", payload: "caption\n2\n00:00:02,000 --> 00:00:03,000\ntail", want: true},
		{name: "timing suffix", payload: "caption\n00:00:02,000 --> 00:00:03,000\ntail", want: true},
		{name: "trailing empty line", payload: "caption\n", want: true},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := PayloadHasAmbiguousBoundary(test.payload); got != test.want {
				t.Fatalf("PayloadHasAmbiguousBoundary(%q) = %t, want %t", test.payload, got, test.want)
			}
		})
	}
}
