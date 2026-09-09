package schema

import (
	"bytes"
	"encoding/json"
	"testing"
)

const maxCueJSONFuzzBytes = 64 << 10

func FuzzDecodeCueJSON(f *testing.F) {
	structurallyInvalid := fuzzRepresentativeMutation(f, func(document map[string]any) {
		document["source_path"] = `C:\Users\cueson-test\private\captions.srt`
	})
	semanticallyInvalid := fuzzRepresentativeMutation(f, func(document map[string]any) {
		cue := document["cues"].([]any)[0].(map[string]any)
		cue["timing"].(map[string]any)["duration_milliseconds"] = float64(1)
	})
	pairedSurrogate := bytes.Replace(Representative(), []byte("captions.srt"), []byte(`face\ud83d\ude00.srt`), 1)
	unpairedHighSurrogate := bytes.Replace(Representative(), []byte("captions.srt"), []byte(`bad\ud800.srt`), 1)
	unpairedLowSurrogate := bytes.Replace(Representative(), []byte("captions.srt"), []byte(`bad\udc00.srt`), 1)
	multipleValues := append(append([]byte(nil), Representative()...), []byte("{}")...)

	for _, seed := range [][]byte{
		Representative(),
		nil,
		[]byte("{"),
		[]byte{0xff},
		multipleValues,
		pairedSurrogate,
		unpairedHighSurrogate,
		unpairedLowSurrogate,
		structurallyInvalid,
		semanticallyInvalid,
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > maxCueJSONFuzzBytes {
			return
		}

		document, err := Decode(data)
		if err != nil {
			return
		}
		if err := document.Validate(); err != nil {
			t.Fatalf("Decode() returned a semantically invalid document: %v", err)
		}
		encoded, err := json.Marshal(document)
		if err != nil {
			t.Fatalf("json.Marshal(Decode()) error = %v", err)
		}
		if _, err := Decode(encoded); err != nil {
			t.Fatalf("Decode(json.Marshal(Decode())) error = %v", err)
		}
	})
}

func fuzzRepresentativeMutation(f *testing.F, mutate func(map[string]any)) []byte {
	f.Helper()

	var document map[string]any
	if err := json.Unmarshal(Representative(), &document); err != nil {
		f.Fatalf("json.Unmarshal(Representative()) error = %v", err)
	}
	mutate(document)
	encoded, err := json.Marshal(document)
	if err != nil {
		f.Fatalf("json.Marshal(mutated representative) error = %v", err)
	}
	return encoded
}
