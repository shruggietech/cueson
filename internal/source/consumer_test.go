package source

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestConsumerDataAndMediaConflictsPreserveFourFormatRestoration(t *testing.T) {
	for _, format := range []string{"subrip", "webvtt", "ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			payload, err := os.ReadFile(filepath.Join("..", "schema", "testdata", "historical-v1.1.0-"+format+".json"))
			if err != nil {
				t.Fatal(err)
			}
			document, err := schema.Decode(payload)
			if err != nil {
				t.Fatal(err)
			}
			// Explicitly selecting the current contract permits consumer additions.
			document.Schema, document.SchemaVersion = schema.ID(), schema.Version()
			document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "__proto__"}, {SpeakerID: "__proto__"}}
			document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 0}
			before, _ := json.Marshal(document)
			if err := document.Validate(); err != nil {
				t.Fatal(err)
			}
			warnings, err := document.ConsumerTimingDiagnostics()
			if err != nil || len(warnings) != 1 {
				t.Fatalf("runtime warnings=%v error=%v", warnings, err)
			}
			want, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(t.TempDir(), "restored."+format)
			if _, err := Restore(context.Background(), document, RestoreOptions{Output: path, Metadata: MetadataNone}); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(got, want) {
				t.Fatalf("source changed: %v", err)
			}
			after, _ := json.Marshal(document)
			if !bytes.Equal(before, after) {
				t.Fatal("restoration or runtime diagnostics changed persisted document")
			}
		})
	}
}
