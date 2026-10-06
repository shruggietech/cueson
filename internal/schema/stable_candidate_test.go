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

func TestCurrentStableCandidateAndHistoricalReleaseSchemaIdentity(t *testing.T) {
	t.Parallel()
	if ID() != "https://cueson.io/schema/v1.2.0/cueson.schema.json" || Version() != "1.2.0" {
		t.Fatalf("unexpected current contract: %s / %s", ID(), Version())
	}
	candidate, err := os.ReadFile(filepath.Join("..", "..", "schema", "releases", "v1.2.0", "cueson.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(candidate, Bytes()) {
		t.Fatal("canonical embedded and immutable candidate schema differ")
	}
	if got := sha256.Sum256(candidate); hex.EncodeToString(got[:]) != "f2661a3d52effbab4a82a4d47197b5c7fae58496dc30a397ea3f2f668358b654" {
		t.Fatalf("immutable candidate schema digest changed: %x", got)
	}
	if _, err := compileArtifact(candidate, ID(), Version()); err != nil {
		t.Fatalf("immutable candidate compilation: %v", err)
	}
	immutable, err := os.ReadFile(filepath.Join("..", "..", "schema", "releases", "v1.1.0", "cueson.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(immutable, historicalV11Bytes) {
		t.Fatal("historical embedded and immutable released schema differ")
	}
	digest := sha256.Sum256(immutable)
	if got := hex.EncodeToString(digest[:]); got != "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7" {
		t.Fatalf("released schema digest = %s", got)
	}
}

func TestRegistryRejectsFormerDevelopmentWithoutMutation(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"representative.cueson.json", "scripted-ass.cueson.json", "scripted-ssa.cueson.json"} {
		for _, development := range []string{"1.1.0-dev", "1.2.0-dev"} {
			t.Run(name+"/"+development, func(t *testing.T) {
				payload, err := os.ReadFile(filepath.Join("testdata", name))
				if err != nil {
					t.Fatal(err)
				}
				var instance map[string]any
				if err := json.Unmarshal(payload, &instance); err != nil {
					t.Fatal(err)
				}
				instance["$schema"] = "https://cueson.io/schema/v" + development + "/cueson.schema.json"
				instance["schema_version"] = development
				payload, err = json.Marshal(instance)
				if err != nil {
					t.Fatal(err)
				}
				before := bytes.Clone(payload)
				if _, err := Decode(payload); err == nil {
					t.Fatal("former development identity accepted")
				}
				if !bytes.Equal(before, payload) {
					t.Fatal("identity refusal rewrote original payload")
				}
			})
		}
	}
}

func TestStablePromotionPreservesReviewedNormativeContract(t *testing.T) {
	t.Parallel()
	artifact := annotationArtifact(t)
	stripAnnotations(artifact)
	// S034 changes exact current identity, with no change to S033 constraints.
	// Restoring its two instance constants and artifact identity must recover
	// the reviewed development normative snapshot exactly.
	artifact["$id"] = "https://cueson.io/schema/v1.2.0-dev/cueson.schema.json"
	properties := artifact["properties"].(map[string]any)
	properties["$schema"].(map[string]any)["const"] = "https://cueson.io/schema/v1.2.0-dev/cueson.schema.json"
	properties["schema_version"].(map[string]any)["const"] = "1.2.0-dev"
	encoded, err := json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	if got := sha256.Sum256(encoded); hex.EncodeToString(got[:]) != "0e69d9c884add390f345e0a914a80c8a8a315d78f1bcb6ae2b636a43816288cd" {
		t.Fatalf("stable promotion changed reviewed normative constraints: %x", got)
	}
}
