package schema

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/source"
)

func historicalRepresentative(t *testing.T) map[string]any {
	t.Helper()
	var artifact map[string]any
	if err := json.Unmarshal(historicalV1Bytes, &artifact); err != nil {
		t.Fatal(err)
	}
	examples, ok := artifact["examples"].([]any)
	if !ok || len(examples) == 0 {
		t.Fatal("historical schema has no representative example")
	}
	return examples[0].(map[string]any)
}

func TestHistoricalRegistryUsesImmutableLocalAuthority(t *testing.T) {
	t.Parallel()
	released, err := os.ReadFile(filepath.Join("..", "..", "schema", "releases", "v1.0.0", "cueson.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(released, historicalV1Bytes) {
		t.Fatal("embedded historical schema differs from immutable release")
	}
	digest := sha256.Sum256(historicalV1Bytes)
	if got := hex.EncodeToString(digest[:]); got != "1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541" {
		t.Fatalf("historical schema digest = %s", got)
	}
	doc := historicalRepresentative(t)
	doc["producer"] = map[string]any{"name": "independent-producer", "version": "42.7-custom"}
	payload, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	before := bytes.Clone(payload)
	decoded, err := Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.Schema != historicalV1ID || decoded.SchemaVersion != historicalV1Version || decoded.Producer.Name != "independent-producer" || decoded.Producer.Version != "42.7-custom" {
		t.Fatalf("historical identity or producer changed: %#v", decoded)
	}
	if !bytes.Equal(before, payload) {
		t.Fatal("Decode mutated original input bytes")
	}
	if !reflect.DeepEqual(doc["source"], mustJSONMap(t, decoded)["source"]) {
		t.Fatal("historical source envelope changed during typed decoding")
	}
}

func mustJSONMap(t *testing.T, value any) map[string]any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestRegistryRejectsIncompleteOrApproximateIdentity(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"missing schema", func(doc map[string]any) { delete(doc, "$schema") }},
		{"missing version", func(doc map[string]any) { delete(doc, "schema_version") }},
		{"wrong schema type", func(doc map[string]any) { doc["$schema"] = 10 }},
		{"wrong version type", func(doc map[string]any) { doc["schema_version"] = nil }},
		{"cross version pair", func(doc map[string]any) { doc["schema_version"] = Version() }},
		{"trailing slash", func(doc map[string]any) { doc["$schema"] = historicalV1ID + "/" }},
		{"case folding", func(doc map[string]any) { doc["$schema"] = strings.ToUpper(historicalV1ID) }},
		{"unversioned alias", func(doc map[string]any) { doc["$schema"] = "https://cueson.io/schema/cueson.schema.json" }},
		{"future identity", func(doc map[string]any) {
			doc["$schema"] = "https://cueson.io/schema/v9.0.0/cueson.schema.json"
			doc["schema_version"] = "9.0.0"
		}},
		{"unpromoted release identity", func(doc map[string]any) {
			doc["$schema"] = "https://cueson.io/schema/v1.1.0/cueson.schema.json"
			doc["schema_version"] = "1.1.0"
		}},
		{"unsupported envelope", func(doc map[string]any) {
			doc["$schema"] = "https://cueson.io/schema/v0.0.0/cueson.schema.json"
			doc["schema_version"] = "0.0.0"
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := historicalRepresentative(t)
			tc.mutate(doc)
			data, err := json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Decode(data); err == nil || !strings.Contains(err.Error(), "identity") {
				t.Fatalf("Decode error = %v, want explicit identity rejection", err)
			}
		})
	}
}

func TestHistoricalWebVTTNativeContentSurvivesRegistryDecode(t *testing.T) {
	t.Parallel()
	// The SubRip fixture is the immutable released schema's complete root
	// example. The WebVTT fixture was constructed from that historical root and
	// the tagged v1.0.0 native shapes, with actual UTF-8 source bytes, REGION,
	// NOTE, STYLE, duplicate/unknown settings, voice and timestamp observations.
	// They are frozen historical inputs, independent of current representatives.
	payload, err := os.ReadFile(filepath.Join("testdata", "historical-v1.0.0-webvtt.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(payload, &doc); err != nil {
		t.Fatal(err)
	}
	decoded, err := Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	projection := mustJSONMap(t, decoded)
	if !reflect.DeepEqual(doc["format_data"], projection["format_data"]) {
		t.Fatal("historical WebVTT document native content was dropped")
	}
	wantCue := firstCue(doc)
	gotCue := firstCue(projection)
	for _, key := range []string{"format_data", "payload", "source_identifier", "speakers", "tokens"} {
		if !reflect.DeepEqual(wantCue[key], gotCue[key]) {
			t.Fatalf("historical WebVTT cue %s changed", key)
		}
	}
	if err := decoded.Validate(); err != nil {
		t.Fatalf("downstream historical revalidation: %v", err)
	}
	if err := source.ValidateIntegrity(context.Background(), decoded); err != nil {
		t.Fatalf("frozen historical fixture integrity: %v", err)
	}
}

func TestFrozenHistoricalSubRipFixtureMatchesReleasedExample(t *testing.T) {
	t.Parallel()
	payload, err := os.ReadFile(filepath.Join("testdata", "historical-v1.0.0-subrip.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(payload, &doc); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(doc, historicalRepresentative(t)) {
		t.Fatal("frozen SubRip fixture differs from immutable released example")
	}
	decoded, err := Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	if err := source.ValidateIntegrity(context.Background(), decoded); err != nil {
		t.Fatalf("frozen historical fixture integrity: %v", err)
	}
}

func TestHistoricalRegistryRetainsStructuralAndSemanticFailures(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		mutate func(map[string]any)
		want   string
	}{
		{"historical provisional capabilities", func(doc map[string]any) { doc["format_support"].(map[string]any)["status"] = "experimental" }, "structure"},
		{"scripted historical branch", func(doc map[string]any) { doc["format_data"].(map[string]any)["ass"] = map[string]any{} }, "structure"},
		{"historical unknown field", func(doc map[string]any) { doc["source_path"] = "private" }, "structure"},
		{"historical summary disagreement", func(doc map[string]any) { doc["stats"].(map[string]any)["cue_count"] = 99 }, "semantics"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := historicalRepresentative(t)
			tc.mutate(doc)
			data, err := json.Marshal(doc)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := Decode(data); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("Decode error = %v, want %s rejection", err, tc.want)
			}
		})
	}
}

func TestRegistryArtifactIdentityAndExternalReferences(t *testing.T) {
	t.Parallel()
	var artifact map[string]any
	if err := json.Unmarshal(historicalV1Bytes, &artifact); err != nil {
		t.Fatal(err)
	}
	artifact["$schema"] = "http://json-schema.org/draft-07/schema#"
	data, err := json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := compileArtifact(data, historicalV1ID, historicalV1Version); err == nil || !strings.Contains(err.Error(), "dialect") {
		t.Fatalf("compile wrong dialect = %v", err)
	}
	artifact["$schema"] = "https://json-schema.org/draft/2020-12/schema"
	artifact["$ref"] = "https://example.invalid/unbundled.schema.json"
	data, err = json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := compileArtifact(data, historicalV1ID, historicalV1Version); err == nil || !strings.Contains(err.Error(), "unbundled schema resource") {
		t.Fatalf("compile external reference = %v, want local-only failure", err)
	}
}
