package schema

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestStableCandidateImmutableSchemaIdentity(t *testing.T) {
	t.Parallel()
	if ID() != "https://cueson.io/schema/v1.1.0/cueson.schema.json" || Version() != "1.1.0" {
		t.Fatalf("unexpected current contract: %s / %s", ID(), Version())
	}
	immutable, err := os.ReadFile(filepath.Join("..", "..", "schema", "releases", "v1.1.0", "cueson.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(immutable, Bytes()) {
		t.Fatal("canonical and immutable candidate schema differ")
	}
	digest := sha256.Sum256(immutable)
	if got := hex.EncodeToString(digest[:]); got != "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7" {
		t.Fatalf("candidate schema digest = %s", got)
	}
}

func TestStableCandidateRejectsUnreleasedDevelopmentWithoutMutation(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"representative.cueson.json", "scripted-ass.cueson.json", "scripted-ssa.cueson.json"} {
		t.Run(name, func(t *testing.T) {
			payload, err := os.ReadFile(filepath.Join("testdata", name))
			if err != nil {
				t.Fatal(err)
			}
			var instance map[string]any
			if err := json.Unmarshal(payload, &instance); err != nil {
				t.Fatal(err)
			}
			instance["$schema"] = "https://cueson.io/schema/v1.1.0-dev/cueson.schema.json"
			instance["schema_version"] = "1.1.0-dev"
			payload, err = json.Marshal(instance)
			if err != nil {
				t.Fatal(err)
			}
			before := bytes.Clone(payload)
			if _, err := Decode(payload); err == nil {
				t.Fatal("unreleased development identity accepted")
			}
			if !bytes.Equal(before, payload) {
				t.Fatal("identity refusal rewrote original payload")
			}
		})
	}
}
