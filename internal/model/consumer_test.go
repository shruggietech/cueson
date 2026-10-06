package model

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
)

func consumerExample() Document {
	doc := representativeDocument()
	doc.Schema = "https://cueson.io/schema/v1.2.0/cueson.schema.json"
	doc.SchemaVersion = "1.2.0"
	return doc
}

func consumerInt(value int64) *int64 { return &value }

func TestConsumerStableContractRejectsFormerDevelopment(t *testing.T) {
	doc := consumerExample()
	doc.Schema = "https://cueson.io/schema/v1.2.0-dev/cueson.schema.json"
	doc.SchemaVersion = "1.2.0-dev"
	doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "consumer-assigned"}}
	if err := doc.ValidateConsumerAnnotations(); err == nil {
		t.Fatal("former development contract accepted consumer attribution in typed revalidation")
	}
	doc.Cues[0].SpeakerAttributions = nil
	doc.MediaTiming = &MediaTiming{DurationMilliseconds: 5000}
	if err := doc.ValidateConsumerAnnotations(); err == nil {
		t.Fatal("former development contract accepted media timing in typed revalidation")
	}
}

func TestConsumerIdentifierPreservationAndLimits(t *testing.T) {
	for _, id := range []string{"Alex", "__proto__", "https://example.test/id", "<script>inert</script>", "YWJjZA==", "مريم", "李 四", "👩‍💻", "e\u0301", "é", strings.Repeat("🗣", 256)} {
		doc := consumerExample()
		doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: id}, {SpeakerID: id}}
		if err := doc.Validate(); err != nil {
			t.Fatalf("accepted identifier %q: %v", id, err)
		}
		encoded, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		var decoded Document
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			t.Fatal(err)
		}
		if len(decoded.Cues[0].SpeakerAttributions) != 2 || decoded.Cues[0].SpeakerAttributions[0].SpeakerID != id {
			t.Fatal("identifier/order did not survive serialization")
		}
	}
	invalid := []string{"", " left", "right ", "\u00a0x", "x\u3000", strings.Repeat("🗣", 257), string([]byte{0xff})}
	for _, r := range []rune{0, 9, 10, 31, 127, 159, 0x2028, 0x2029, 0x061c, 0x200e, 0x200f, 0x202a, 0x202b, 0x202c, 0x202d, 0x202e, 0x2066, 0x2067, 0x2068, 0x2069, 0xfeff} {
		invalid = append(invalid, "private"+string(r)+"value")
	}
	for _, id := range invalid {
		doc := consumerExample()
		doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: id}}
		if err := doc.Validate(); err == nil {
			t.Fatalf("invalid identifier accepted: %q", id)
		} else if strings.Contains(err.Error(), "private") {
			t.Fatalf("identifier leaked: %v", err)
		}
	}
}

func TestConsumerAttributionIntervals(t *testing.T) {
	for _, test := range []struct {
		name       string
		start, end *int64
		valid      bool
	}{
		{"untimed", nil, nil, true}, {"whole cue", consumerInt(1250), consumerInt(4200), true},
		{"inside", consumerInt(2000), consumerInt(3000), true}, {"missing start", nil, consumerInt(3000), false},
		{"missing end", consumerInt(2000), nil, false}, {"negative", consumerInt(-1), consumerInt(3000), false},
		{"empty", consumerInt(2000), consumerInt(2000), false}, {"reversed", consumerInt(3000), consumerInt(2000), false},
		{"before cue", consumerInt(1249), consumerInt(3000), false}, {"after cue", consumerInt(2000), consumerInt(4201), false},
		{"huge", consumerInt(2000), consumerInt(math.MaxInt64), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			doc := consumerExample()
			doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id", StartMilliseconds: test.start, EndMilliseconds: test.end}}
			err := doc.Validate()
			if (err == nil) != test.valid {
				t.Fatalf("Validate=%v valid=%v", err, test.valid)
			}
		})
	}
	doc := consumerExample()
	doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id", StartMilliseconds: consumerInt(2000), EndMilliseconds: consumerInt(3500)}, {SpeakerID: "id", StartMilliseconds: consumerInt(2500), EndMilliseconds: consumerInt(4000)}}
	if err := doc.Validate(); err != nil {
		t.Fatal("overlap/repetition rejected", err)
	}
}

func TestConsumerDeclaredMediaBounds(t *testing.T) {
	for _, test := range []struct {
		name       string
		duration   int64
		offset     *int64
		start, end *int64
		valid      bool
	}{
		{"default alignment", 4200, nil, consumerInt(1250), consumerInt(4200), true},
		{"positive alignment", 2950, consumerInt(1250), consumerInt(1250), consumerInt(4200), true},
		{"negative alignment", 5200, consumerInt(-1000), consumerInt(1250), consumerInt(4200), true},
		{"beyond media", 4199, nil, consumerInt(1250), consumerInt(4200), false},
		{"before media", 3000, consumerInt(1251), consumerInt(1250), consumerInt(4200), false},
		{"zero untimed", 0, nil, nil, nil, true}, {"zero timed", 0, nil, consumerInt(1250), consumerInt(4200), false},
		{"negative duration", -1, nil, nil, nil, false},
		{"endpoint overflow", 1, consumerInt(math.MaxInt64), nil, nil, false},
		{"minimum offset", math.MaxInt64, consumerInt(math.MinInt64), nil, nil, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			doc := consumerExample()
			doc.MediaTiming = &MediaTiming{DurationMilliseconds: test.duration, TimelineStartMilliseconds: test.offset}
			doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id", StartMilliseconds: test.start, EndMilliseconds: test.end}}
			err := doc.Validate()
			if (err == nil) != test.valid {
				t.Fatalf("Validate=%v valid=%v", err, test.valid)
			}
		})
	}
}

func TestConsumerAssessmentAndOmissions(t *testing.T) {
	doc := consumerExample()
	if doc.ConsumerAnnotationsSummary() != nil {
		t.Fatal("unannotated summary must be absent")
	}
	doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "private-id"}}
	if got := doc.ConsumerAnnotationsSummary(); got.MediaBoundaryCheck != "unavailable" || got.AttributionCount != 1 || got.UntimedAttributionCount != 1 {
		t.Fatalf("summary %+v", got)
	}
	doc.MediaTiming = &MediaTiming{DurationMilliseconds: 0}
	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
	if got := doc.ConsumerAnnotationsSummary(); got.MediaBoundaryCheck != "not_evaluated" || got.CueMediaConflictCount != 1 {
		t.Fatalf("summary %+v", got)
	}
	warnings, err := doc.ConsumerTimingDiagnostics()
	if err != nil || len(warnings) != 1 || warnings[0].Code != "consumer_cue_media_conflict" {
		t.Fatalf("warnings=%+v err=%v", warnings, err)
	}
	if doc.Stats.WarningCount != 0 || len(doc.Diagnostics) != 0 {
		t.Fatal("runtime warnings altered source diagnostics")
	}
	before, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(before), "timeline_start_milliseconds") {
		t.Fatal("omitted alignment was materialized")
	}
	losses, err := doc.ConsumerAnnotationLossDiagnostics()
	if err != nil || len(losses) != 2 {
		t.Fatalf("losses=%+v err=%v", losses, err)
	}
	for _, loss := range losses {
		if strings.Contains(loss.Message, "private-id") {
			t.Fatal("identifier leaked")
		}
	}
	after, _ := json.Marshal(doc)
	if string(before) != string(after) {
		t.Fatal("helper mutated document")
	}
	doc.MediaTiming.DurationMilliseconds = 4200
	doc.Cues[0].SpeakerAttributions[0].StartMilliseconds = consumerInt(1250)
	doc.Cues[0].SpeakerAttributions[0].EndMilliseconds = consumerInt(4200)
	if got := doc.ConsumerAnnotationsSummary(); got.MediaBoundaryCheck != "checked" || got.TimedAttributionCount != 1 || got.CueMediaConflictCount != 0 {
		t.Fatalf("summary %+v", got)
	}
}

func TestConsumerLimitsAndHistoricalContracts(t *testing.T) {
	doc := consumerExample()
	doc.Cues[0].SpeakerAttributions = make([]SpeakerAttribution, MaxItemOccurrences)
	for i := range doc.Cues[0].SpeakerAttributions {
		doc.Cues[0].SpeakerAttributions[i].SpeakerID = "id"
	}
	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
	doc.Cues[0].SpeakerAttributions = append(doc.Cues[0].SpeakerAttributions, SpeakerAttribution{SpeakerID: "id"})
	if err := doc.Validate(); err == nil {
		t.Fatal("attribution limit not enforced")
	}
	for _, version := range []string{"1.0.0", "1.1.0"} {
		doc = consumerExample()
		doc.SchemaVersion = version
		doc.Schema = "https://cueson.io/schema/v" + version + "/cueson.schema.json"
		doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{}
		if err := doc.Validate(); err == nil {
			t.Fatal("historical typed model accepted consumer field")
		}
		doc.Cues[0].SpeakerAttributions = nil
		doc.MediaTiming = &MediaTiming{}
		if err := doc.Validate(); err == nil {
			t.Fatal("historical typed model accepted media declaration")
		}
	}
	doc = consumerExample()
	doc.Cues = make([]Cue, 9)
	for i := range doc.Cues {
		doc.Cues[i].SpeakerAttributions = make([]SpeakerAttribution, MaxItemOccurrences)
	}
	if _, err := doc.ConsumerAnnotationLossDiagnostics(); err == nil {
		t.Fatal("loss amplification not bounded")
	}
}

func TestConsumerCurrentScriptedContracts(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		doc := scriptedExample(t, format)
		doc.SchemaVersion = "1.2.0"
		doc.Schema = "https://cueson.io/schema/v1.2.0/cueson.schema.json"
		doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "independent"}}
		if err := doc.Validate(); err != nil {
			t.Fatalf("%s: %v", format, err)
		}
	}
}

func TestConsumerCueConflictsRemainNonfatalAndBounded(t *testing.T) {
	doc := consumerExample()
	doc.MediaTiming = &MediaTiming{DurationMilliseconds: 3000}
	doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id", StartMilliseconds: consumerInt(1500), EndMilliseconds: consumerInt(2000)}}
	if err := doc.Validate(); err != nil {
		t.Fatal("outlying source cue rejected despite valid assignment", err)
	}
	doc.Cues = make([]Cue, MaxDocumentItems)
	for index := range doc.Cues {
		doc.Cues[index].Timing = Timing{StartMilliseconds: 1250, EndMilliseconds: 4200}
	}
	warnings, err := doc.ConsumerTimingDiagnostics()
	if err != nil || len(warnings) != 1 || !strings.Contains(warnings[0].Message, "65536") {
		t.Fatalf("unbounded or inaccurate warning: %+v %v", warnings, err)
	}
}

func TestConsumerZeroCueUntimedAndLossCeiling(t *testing.T) {
	doc := consumerExample()
	doc.Cues[0].Timing.EndMilliseconds = 1250
	doc.Cues[0].Timing.DurationMilliseconds = 0
	doc.Document.MediaEndMilliseconds = consumerInt(1250)
	doc.Document.MediaSpanMilliseconds = consumerInt(0)
	doc.Stats.MediaSpanMilliseconds = consumerInt(0)
	doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id"}}
	if err := doc.Validate(); err != nil {
		t.Fatal("untimed zero-duration cue rejected", err)
	}
	doc.Cues = make([]Cue, 8)
	for index := range doc.Cues {
		doc.Cues[index].SpeakerAttributions = make([]SpeakerAttribution, MaxItemOccurrences)
	}
	losses, err := doc.ConsumerAnnotationLossDiagnostics()
	if err != nil || len(losses) != MaxDiagnostics {
		t.Fatalf("exact omission ceiling rejected: %d %v", len(losses), err)
	}
	doc.MediaTiming = &MediaTiming{}
	if _, err := doc.ConsumerAnnotationLossDiagnostics(); err == nil {
		t.Fatal("media omission bypassed diagnostic ceiling")
	}
}

func TestConsumerPrivateScriptedTargetRevalidatesAnnotations(t *testing.T) {
	source := scriptedExample(t, "ssa")
	target := captureFreeScriptedExample(t, "ass")
	target.Source = source.Source
	target.Schema = "https://cueson.io/schema/v1.2.0/cueson.schema.json"
	target.SchemaVersion = "1.2.0"
	target.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id"}}
	if err := ValidateScriptedTarget(source, target); err != nil {
		t.Fatal("current private target rejected", err)
	}
	target.Cues[0].SpeakerAttributions[0].StartMilliseconds = consumerInt(0)
	if err := ValidateScriptedTarget(source, target); err == nil {
		t.Fatal("private target skipped attribution validation")
	}
}

func FuzzConsumerIdentifier(f *testing.F) {
	for _, seed := range []string{"Alex", "مريم", "👩‍💻", "x\u202ey", strings.Repeat("🗣", 257)} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, id string) {
		doc := consumerExample()
		doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: id}}
		if err := doc.Validate(); err != nil {
			return
		}
		data, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		var decoded Document
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatal(err)
		}
		if decoded.Cues[0].SpeakerAttributions[0].SpeakerID != id {
			t.Fatal("accepted ID changed")
		}
	})
}

func FuzzConsumerMediaArithmetic(f *testing.F) {
	f.Add(int64(0), int64(4200), int64(1250), int64(4200))
	f.Add(int64(math.MaxInt64), int64(1), int64(1250), int64(4200))
	f.Add(int64(math.MinInt64), int64(math.MaxInt64), int64(0), int64(1))
	f.Fuzz(func(t *testing.T, offset, duration, start, end int64) {
		doc := consumerExample()
		doc.MediaTiming = &MediaTiming{DurationMilliseconds: duration, TimelineStartMilliseconds: &offset}
		doc.Cues[0].SpeakerAttributions = []SpeakerAttribution{{SpeakerID: "id", StartMilliseconds: &start, EndMilliseconds: &end}}
		if err := doc.Validate(); err != nil {
			return
		}
		if duration < 0 || offset > math.MaxInt64-duration || start < offset || end > offset+duration || start < 1250 || end > 4200 || end <= start {
			t.Fatal("accepted invalid checked interval")
		}
		if summary := doc.ConsumerAnnotationsSummary(); summary.MediaBoundaryCheck != "checked" || summary.TimedAttributionCount != 1 {
			t.Fatalf("inaccurate assessment %+v", summary)
		}
	})
}
