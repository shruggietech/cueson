package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/source"
)

func consumerMap(t *testing.T) map[string]any {
	t.Helper()
	doc := representativeMap(t)
	firstCue(doc)["speaker_attributions"] = []any{map[string]any{"speaker_id": "Interviewer", "start_milliseconds": 1250, "end_milliseconds": 4200}}
	doc["media_timing"] = map[string]any{"duration_milliseconds": 5000}
	return doc
}

func TestConsumerSchemaExactTextAndOptionality(t *testing.T) {
	t.Parallel()
	for _, id := range []string{"Alex", "__proto__", "<script>alert(1)</script>", "https://example.invalid/a", "SGVsbG8=", "日本語", "العربية", "A\u200dB", "é", "e\u0301", "\U0001f642", strings.Repeat("\U0001f642", 256)} {
		t.Run(id, func(t *testing.T) {
			doc := consumerMap(t)
			firstCue(doc)["speaker_attributions"] = []any{map[string]any{"speaker_id": id}, map[string]any{"speaker_id": id}}
			payload, err := json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := Decode(payload)
			if err != nil {
				t.Fatal(err)
			}
			projection := mustJSONMap(t, decoded)
			if !reflect.DeepEqual(firstCue(doc)["speaker_attributions"], firstCue(projection)["speaker_attributions"]) {
				t.Fatal("consumer text/order changed")
			}
			if !reflect.DeepEqual(firstCue(doc)["speakers"], firstCue(projection)["speakers"]) || !reflect.DeepEqual(doc["source"], projection["source"]) {
				t.Fatal("native labels or source changed")
			}
			if err = source.ValidateIntegrity(context.Background(), decoded); err != nil {
				t.Fatal(err)
			}
		})
	}
	for _, assignments := range []any{[]any{}} {
		doc := representativeMap(t)
		firstCue(doc)["speaker_attributions"] = assignments
		assertMapValid(t, doc)
	}
	assertMapValid(t, representativeMap(t))
}

func TestConsumerSchemaRejectsInvalidStringsAndShapes(t *testing.T) {
	t.Parallel()
	invalid := []any{"", strings.Repeat("a", 257), 12, nil, []any{"id"}, map[string]any{"id": "x"}, " leading", "trailing ", "\u00a0id", "id\u3000"}
	for r := rune(0); r <= 0x1f; r++ {
		invalid = append(invalid, "a"+string(r)+"b")
	}
	for r := rune(0x7f); r <= 0x9f; r++ {
		invalid = append(invalid, "a"+string(r)+"b")
	}
	for _, r := range []rune{0x9, 0xa, 0xb, 0xc, 0xd, 0x20, 0x85, 0xa0, 0x1680, 0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005, 0x2006, 0x2007, 0x2008, 0x2009, 0x200a, 0x2028, 0x2029, 0x202f, 0x205f, 0x3000} {
		invalid = append(invalid, string(r)+"id", "id"+string(r))
	}
	for _, r := range []rune{0, 0x1f, 0x7f, 0x85, 0x9f, 0x2028, 0x2029, 0x061c, 0x200e, 0x200f, 0x202a, 0x202b, 0x202c, 0x202d, 0x202e, 0x2066, 0x2067, 0x2068, 0x2069, 0xfeff} {
		invalid = append(invalid, "a"+string(r)+"b")
	}
	for i, value := range invalid {
		doc := consumerMap(t)
		firstCue(doc)["speaker_attributions"] = []any{map[string]any{"speaker_id": value}}
		data, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		if err = Validate(data); err == nil {
			t.Fatalf("invalid speaker case %d accepted", i)
		}
	}
	for _, mutate := range []func(map[string]any){
		func(d map[string]any) { firstCue(d)["speaker_attributions"] = nil },
		func(d map[string]any) {
			firstCue(d)["speaker_attributions"] = []any{map[string]any{"speaker_id": "x", "extra": true}}
		},
		func(d map[string]any) {
			firstCue(d)["speaker_attributions"] = []any{map[string]any{"speaker_id": "x", "start_milliseconds": 1250}}
		},
		func(d map[string]any) {
			firstCue(d)["speaker_attributions"] = []any{map[string]any{"speaker_id": "x", "end_milliseconds": 4200}}
		},
		func(d map[string]any) { d["media_timing"] = nil },
		func(d map[string]any) { d["media_timing"] = map[string]any{} },
		func(d map[string]any) {
			d["media_timing"] = map[string]any{"duration_milliseconds": 5000, "extra": true}
		},
		func(d map[string]any) { d["media_timing"] = map[string]any{"duration_milliseconds": -1} },
		func(d map[string]any) { d["media_timing"] = map[string]any{"duration_milliseconds": 1.5} },
		func(d map[string]any) {
			d["media_timing"] = map[string]any{"duration_milliseconds": json.Number("9223372036854775808")}
		},
		func(d map[string]any) {
			d["media_timing"] = map[string]any{"duration_milliseconds": 5000, "timeline_start_milliseconds": json.Number("-9223372036854775809")}
		},
	} {
		doc := consumerMap(t)
		mutate(doc)
		data, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		if err = Validate(data); err == nil {
			t.Fatal("invalid consumer shape accepted")
		}
	}
}

func TestConsumerSchemaEscapedScalarValidation(t *testing.T) {
	t.Parallel()
	doc := consumerMap(t)
	firstCue(doc)["speaker_attributions"] = []any{map[string]any{"speaker_id": strings.Repeat("🙂", 256)}}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	escaped := bytes.ReplaceAll(data, []byte("🙂"), []byte(`\ud83d\ude42`))
	if err = Validate(escaped); err != nil {
		t.Fatalf("256 decoded supplementary scalars: %v", err)
	}
	for _, invalid := range [][]byte{[]byte(`\ud83d`), []byte(`\ude42`), []byte(`\ud83d\u0041`), []byte(`\u202e`), []byte(`\u0000`)} {
		bad := bytes.Replace(escaped, []byte(`\ud83d\ude42`), invalid, 1)
		if err = Validate(bad); err == nil {
			t.Fatal("invalid escaped consumer scalar accepted")
		}
	}
}

func TestConsumerSchemaErrorDoesNotExposeSpeakerText(t *testing.T) {
	t.Parallel()
	doc := consumerMap(t)
	firstCue(doc)["speaker_attributions"] = []any{map[string]any{"speaker_id": "private-speaker-marker\n<script>alert(1)</script>"}}
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err = Validate(data); err == nil || strings.Contains(err.Error(), "private-speaker-marker") || strings.Contains(err.Error(), "<script>") {
		t.Fatalf("consumer error leaked content or accepted input: %v", err)
	}
}

func TestConsumerSchemaIntervalsCollectionsAndMediaBounds(t *testing.T) {
	t.Parallel()
	for _, interval := range [][2]any{{1250, 1250}, {4200, 1250}, {-1, 4200}, {1250.5, 4200}, {1250, 4201}, {1249, 4200}, {1250, json.Number("9223372036854775808")}} {
		doc := consumerMap(t)
		firstCue(doc)["speaker_attributions"] = []any{map[string]any{"speaker_id": "x", "start_milliseconds": interval[0], "end_milliseconds": interval[1]}}
		data, _ := json.Marshal(doc)
		if err := Validate(data); err == nil {
			t.Fatal("invalid consumer interval accepted")
		}
	}
	for _, count := range []int{1024, 1025} {
		doc := consumerMap(t)
		entries := make([]any, count)
		for i := range entries {
			entries[i] = map[string]any{"speaker_id": "repeated"}
		}
		firstCue(doc)["speaker_attributions"] = entries
		data, _ := json.Marshal(doc)
		if err := Validate(data); (err == nil) != (count == 1024) {
			t.Fatalf("consumer count %d: %v", count, err)
		}
	}
	for _, tc := range []struct {
		duration, start int64
		valid           bool
	}{{2950, 1250, true}, {4200, 0, true}, {4199, 0, false}, {0, 0, false}, {5000, -500, true}, {math.MaxInt64, 1, false}} {
		doc := consumerMap(t)
		doc["media_timing"] = map[string]any{"duration_milliseconds": tc.duration, "timeline_start_milliseconds": tc.start}
		data, _ := json.Marshal(doc)
		if err := Validate(data); (err == nil) != tc.valid {
			t.Fatalf("media duration/start %d/%d: %v", tc.duration, tc.start, err)
		}
	}
}

func TestHistoricalV11RegistryPreservesFourFormats(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"subrip", "webvtt", "ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			data, err := os.ReadFile("testdata/historical-v1.1.0-" + format + ".json")
			if err != nil {
				t.Fatal(err)
			}
			before := bytes.Clone(data)
			doc, err := Decode(data)
			if err != nil {
				t.Fatal(err)
			}
			if doc.SchemaVersion != "1.1.0" || doc.Schema != historicalV11ID {
				t.Fatal("historical identity changed")
			}
			var original map[string]any
			if err = json.Unmarshal(data, &original); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(original, mustJSONMap(t, doc)) || !bytes.Equal(before, data) {
				t.Fatal("historical fields changed")
			}
			if err = doc.Validate(); err != nil {
				t.Fatal(err)
			}
			if err = source.ValidateIntegrity(context.Background(), doc); err != nil {
				t.Fatal(err)
			}
			firstCue(original)["speaker_attributions"] = []any{}
			payload, _ := json.Marshal(original)
			if err = Validate(payload); err == nil {
				t.Fatal("historical identity accepted new consumer field")
			}
			delete(firstCue(original), "speaker_attributions")
			original["media_timing"] = map[string]any{"duration_milliseconds": 5000}
			payload, _ = json.Marshal(original)
			if err = Validate(payload); err == nil {
				t.Fatal("historical identity accepted media declaration")
			}
		})
	}
}

func TestHistoricalV1RejectsConsumerAdditions(t *testing.T) {
	t.Parallel()
	for _, mutate := range []func(map[string]any){
		func(d map[string]any) { firstCue(d)["speaker_attributions"] = []any{} },
		func(d map[string]any) { d["media_timing"] = map[string]any{"duration_milliseconds": 5000} },
	} {
		doc := historicalRepresentative(t)
		mutate(doc)
		payload, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		if err = Validate(payload); err == nil {
			t.Fatal("historical1.0 accepted consumer additions")
		}
	}
}
