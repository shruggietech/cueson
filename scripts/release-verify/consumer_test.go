package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStableNativeDispatchIncludesCurrentCandidate(t *testing.T) {
	for _, version := range []string{"1.1.0", "1.2.0"} {
		if !requiresStableNativeProof(version, false) || requiresStableNativeProof(version, true) {
			t.Fatalf("incorrect stable dispatch for %s", version)
		}
	}
	if requiresStableNativeProof("1.2.0-dev", false) || requiresStableNativeProof("1.0.0", false) {
		t.Fatal("historical or development identity enters current stable proof")
	}
}

func TestConsumerMutationPreservesIntegerSourceTruth(t *testing.T) {
	base := []byte(`{"source":{"unix_ns":1788912000000000001},"cues":[{"timing":{"start_milliseconds":1250,"end_milliseconds":4200}}]}`)
	annotated, err := consumerPayload(base, "valid")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(annotated, []byte("1788912000000000001")) {
		t.Fatal("consumer proof rounded source nanoseconds")
	}
	var document map[string]json.RawMessage
	if err := json.Unmarshal(annotated, &document); err != nil {
		t.Fatal(err)
	}
	if string(document["source"]) != `{"unix_ns":1788912000000000001}` {
		t.Fatal("consumer proof altered source metadata")
	}
}

func TestMinimalIdentityProbeIsExactlyTwoFields(t *testing.T) {
	payload := []byte(`{"$schema":"https://cueson.io/schema/v1.2.0/cueson.schema.json","schema_version":"1.2.0","source":{"private":"omitted"},"cues":[]}`)
	minimal, err := minimalIdentityPayload(payload)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(minimal, &fields); err != nil {
		t.Fatal(err)
	}
	if len(fields) != 2 || string(fields["schema_version"]) != `"1.2.0"` || string(fields["$schema"]) != `"https://cueson.io/schema/v1.2.0/cueson.schema.json"` {
		t.Fatal("identity probe includes source fields or relabels exact identity")
	}
}

func TestMediaBoundaryInputKeepsAttributionInsideCue(t *testing.T) {
	payload, err := consumerPayload([]byte(`{"cues":[{"timing":{"start_milliseconds":1250,"end_milliseconds":4200}}]}`), "invalid_media_boundary")
	if err != nil {
		t.Fatal(err)
	}
	var result struct {
		Media struct {
			Duration int64 `json:"duration_milliseconds"`
		} `json:"media_timing"`
		Cues []struct {
			Attributions []struct {
				Start int64 `json:"start_milliseconds"`
				End   int64 `json:"end_milliseconds"`
			} `json:"speaker_attributions"`
		} `json:"cues"`
	}
	if err := json.Unmarshal(payload, &result); err != nil {
		t.Fatal(err)
	}
	interval := result.Cues[0].Attributions[0]
	if interval.Start != 1250 || interval.End != 4200 || result.Media.Duration != 4199 {
		t.Fatal("media boundary probe also violates cue containment")
	}
}

func TestSuccessfulConsumerReadoutRejectsBothOutputLeaks(t *testing.T) {
	for _, outputs := range [][2][]byte{{[]byte(proofSpeakerID), nil}, {nil, []byte(proofSpeakerID)}, {[]byte("safe"), []byte("safe")}} {
		err := verifyConsumerReadoutPrivacy(outputs[0], outputs[1])
		wantError := bytes.Contains(outputs[0], []byte(proofSpeakerID)) || bytes.Contains(outputs[1], []byte(proofSpeakerID))
		if (err != nil) != wantError {
			t.Fatal("successful consumer readout privacy gate differs")
		}
	}
}

func TestRefusalInventoryDetectsEqualCountStagingReplacement(t *testing.T) {
	directory := t.TempDir()
	first := filepath.Join(directory, "original")
	if err := os.WriteFile(first, []byte("sentinel"), 0600); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(first, filepath.Join(directory, "staging")); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	if sameNativeDirectoryEntries(before, after) {
		t.Fatal("equal-count staging replacement accepted")
	}
}

func TestEvidenceInventoryFieldsRemainAbsentForEarlierVersions(t *testing.T) {
	old, err := json.Marshal(ReleaseEvidence{Version: "1.1.0", Targets: expectedTargets("1.1.0")})
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"archive_size", "sbom_size", "checksum_manifest"} {
		if bytes.Contains(old, []byte(field)) {
			t.Fatalf("earlier evidence gained %s", field)
		}
	}
}

func TestPublishedV110IdentityGateRejectsUnrelatedFailure(t *testing.T) {
	for _, test := range []struct {
		text  string
		valid bool
	}{
		{"validate Cue JSON structure: select exact contract: unsupported or mismatched Cue JSON identity ($schema, schema_version)", true},
		{"source hash mismatch", false},
		{"unsupported or mismatched Cue JSON identity ($schema, schema_version)", false},
		{"select exact contract: unsupported source", false},
	} {
		if (verifyVersionedIdentityDiagnostic([]byte(test.text), "1.1.0") == nil) != test.valid {
			t.Fatalf("wrong identity gate for %q", test.text)
		}
	}
}

func TestPublishedV110ContractBindsThreeActualHostArchives(t *testing.T) {
	repository := filepath.Clean(filepath.Join("..", ".."))
	contract, err := loadPublishedV110Contract(repository)
	if err != nil {
		t.Fatal(err)
	}
	if len(contract.Assets) != 3 || contract.Version != "1.1.0" || contract.Revision != "7ff45c1d8cd8df377e1fb568b9785286b649fd7c" {
		t.Fatal("published contract identity differs")
	}
	for _, asset := range contract.Assets {
		if asset.GOARCH != "amd64" || asset.Size <= 0 || len(asset.SHA256) != 64 {
			t.Fatal("published host binding differs")
		}
	}
}

func TestPublishedV110ContractRejectsAuthenticationDrift(t *testing.T) {
	repository := filepath.Clean(filepath.Join("..", ".."))
	name := filepath.Join("specs", "S034-prepare-v1-2-release-candidate", "contracts", "published-v1.1.0-consumer.json")
	data, err := os.ReadFile(filepath.Join(repository, name))
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct{ old, new string }{
		{"7ff45c1d8cd8df377e1fb568b9785286b649fd7c", strings.Repeat("a", 40)},
		{"2306834", "2306835"},
		{"2aff230ca03127eeca49d4be447b5a49563026caab96b5c545189a4181e8e672", strings.Repeat("f", 64)},
		{"57089835e30184f9625d5b9b493ebb41d58ca5ac7ec76af84d0c72aa8ca5e908", strings.Repeat("b", 64)},
		{"\"version\": \"1.1.0\"", "\"version\": \"1.2.0\""},
		{"\"version\": \"1.1.0\"", "\"version\": \"1.1.0\", \"unexpected\": true"},
	} {
		directory := t.TempDir()
		path := filepath.Join(directory, name)
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, bytes.Replace(data, []byte(test.old), []byte(test.new), 1), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := loadPublishedV110Contract(directory); err == nil {
			t.Fatalf("published authentication drift accepted: %s", test.old)
		}
	}
}

func TestConsumerNativeProofRejectsOutputAndCountDrift(t *testing.T) {
	valid := []byte(`{"consumer_annotations":{"attribution_count":2,"timed_attribution_count":1,"untimed_attribution_count":1,"media_timing_present":true,"media_boundary_check":"checked","cue_media_conflict_count":0}}`)
	if err := verifyConsumerInspection(valid, "checked", true, 0); err != nil {
		t.Fatal(err)
	}
	var document map[string]any
	if err := json.Unmarshal(valid, &document); err != nil {
		t.Fatal(err)
	}
	document["leaked_speaker"] = proofSpeakerID
	leaked, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	if verifyConsumerInspection(leaked, "checked", true, 0) == nil {
		t.Fatal("speaker value escaped inspection")
	}
	if verifyConsumerInspection(valid, "unavailable", false, 0) == nil {
		t.Fatal("media state drift accepted")
	}
	if verifyConsumerOmissions([]byte("consumer_speaker_attribution_omitted consumer_media_timing_omitted"), false) == nil {
		t.Fatal("missing attribution occurrence accepted")
	}
}

func TestCurrentCandidateEvidenceContract(t *testing.T) {
	repository := filepath.Clean(filepath.Join("..", ".."))
	data, err := os.ReadFile(filepath.Join(repository, "specs", "S034-prepare-v1-2-release-candidate", "contracts", "release-evidence-contract.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema struct {
		Properties  map[string]json.RawMessage `json:"properties"`
		Definitions map[string]struct {
			Properties map[string]json.RawMessage `json:"properties"`
			Required   []string                   `json:"required"`
		} `json:"$defs"`
	}
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	var identity struct {
		Const string `json:"const"`
	}
	if err := json.Unmarshal(schema.Properties["version"], &identity); err != nil || identity.Const != "1.2.0" {
		t.Fatal("candidate schema identity differs")
	}
	for field, want := range map[string]int{"matching_render_count": 4, "conversion_count": 12, "strict_refusal_count": 32, "invalid_timing_refusal_count": 4, "invalid_identifier_refusal_count": 4, "media_boundary_refusal_count": 4, "media_boundary_destination_refusal_count": 8, "cue_media_conflict_count": 4, "restore_count": 4} {
		var count struct {
			Const int `json:"const"`
		}
		if err := json.Unmarshal(schema.Definitions["consumer_annotations"].Properties[field], &count); err != nil || count.Const != want {
			t.Fatalf("candidate proof count %s differs", field)
		}
	}
	for field, want := range map[string]string{"published_consumer_v100": "1.0.0", "published_consumer_v110": "1.1.0"} {
		definition := schema.Definitions[field]
		if err := json.Unmarshal(definition.Properties["version"], &identity); err != nil || identity.Const != want {
			t.Fatalf("consumer evidence identity %s differs", field)
		}
		for _, required := range []string{"version", "source_revision", "archive", "archive_sha256", "positive_control_count", "refusal_count", "identity_probe_count", "selector_identity_probe_count"} {
			if !containsString(definition.Required, required) {
				t.Fatalf("consumer evidence is missing required %s", required)
			}
		}
		positive, refusals := 2, 32
		if want == "1.1.0" {
			positive, refusals = 4, 64
		}
		for key, expected := range map[string]int{"positive_control_count": positive, "refusal_count": refusals, "identity_probe_count": refusals, "selector_identity_probe_count": refusals} {
			var count struct {
				Const int `json:"const"`
			}
			if err := json.Unmarshal(definition.Properties[key], &count); err != nil || count.Const != expected {
				t.Fatalf("consumer evidence %s %s differs", want, key)
			}
		}
	}
	for _, field := range []string{"archive_size", "sbom_size"} {
		if !containsString(schema.Definitions["target"].Required, field) {
			t.Fatalf("candidate inventory omits %s", field)
		}
	}
	for _, field := range []string{"consumer_annotations", "published_consumers"} {
		if !containsString(schema.Definitions["native_proof"].Required, field) {
			t.Fatalf("candidate native proof omits %s", field)
		}
	}
	var strict struct {
		Const int `json:"const"`
	}
	if err := json.Unmarshal(schema.Definitions["native_proof"].Properties["historical_strict_refusal_count"], &strict); err != nil || strict.Const != 12 || !containsString(schema.Definitions["native_proof"].Required, "historical_strict_refusal_count") {
		t.Fatal("historical strict destination proof differs")
	}
}

// The executable proof belongs to the existing test binary so no extra build is required.
func checkConsumerWorkflows(t *testing.T, binary, repository, directory string) {
	t.Helper()
	proof, err := verifyConsumerNativeWorkflows(context.Background(), binary, repository, directory)
	if err != nil {
		t.Fatal(err)
	}
	if proof.MatchingRenderCount != 4 || proof.ConversionCount != 12 || proof.StrictRefusalCount != 32 || proof.InvalidTimingRefusalCount != 4 || proof.InvalidIdentifierRefusalCount != 4 || proof.MediaBoundaryRefusalCount != 4 || proof.MediaBoundaryDestinationRefusalCount != 8 || proof.CueMediaConflictCount != 4 || proof.RestoreCount != 4 {
		t.Fatalf("incomplete consumer proof: %+v", proof)
	}
}
