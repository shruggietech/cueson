package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const proofSpeakerID = "candidate-speaker-opaque-7deae601"

type ConsumerAnnotationProof struct {
	MatchingRenderCount                  int `json:"matching_render_count"`
	ConversionCount                      int `json:"conversion_count"`
	StrictRefusalCount                   int `json:"strict_refusal_count"`
	InvalidTimingRefusalCount            int `json:"invalid_timing_refusal_count"`
	InvalidIdentifierRefusalCount        int `json:"invalid_identifier_refusal_count"`
	MediaBoundaryRefusalCount            int `json:"media_boundary_refusal_count"`
	MediaBoundaryDestinationRefusalCount int `json:"media_boundary_destination_refusal_count"`
	CueMediaConflictCount                int `json:"cue_media_conflict_count"`
	RestoreCount                         int `json:"restore_count"`
}

func requiresStableNativeProof(version string, development bool) bool {
	return !development && (version == "1.1.0" || version == "1.2.0")
}

// RawMessage preserves nanosecond integer timestamps while adding consumer data.
func consumerPayload(payload []byte, variant string) ([]byte, error) {
	var document map[string]json.RawMessage
	if err := json.Unmarshal(payload, &document); err != nil {
		return nil, err
	}
	var cues []map[string]json.RawMessage
	if err := json.Unmarshal(document["cues"], &cues); err != nil || len(cues) == 0 {
		return nil, fmt.Errorf("consumer proof requires cues")
	}
	var first struct {
		Start int64 `json:"start_milliseconds"`
		End   int64 `json:"end_milliseconds"`
	}
	if err := json.Unmarshal(cues[0]["timing"], &first); err != nil || first.End <= first.Start {
		return nil, fmt.Errorf("consumer proof requires a positive cue interval")
	}
	maximum := first.End
	for _, cue := range cues {
		var timing struct {
			End int64 `json:"end_milliseconds"`
		}
		if err := json.Unmarshal(cue["timing"], &timing); err != nil {
			return nil, err
		}
		if timing.End > maximum {
			maximum = timing.End
		}
	}
	identifier := proofSpeakerID
	end := first.End
	if variant == "invalid_identifier" {
		identifier += "\x00"
	}
	if variant == "invalid_timing" {
		end++
	}
	attributions := []map[string]any{{"speaker_id": identifier, "start_milliseconds": first.Start, "end_milliseconds": end}, {"speaker_id": proofSpeakerID + "-untimed"}}
	if variant == "conflict" {
		attributions = []map[string]any{{"speaker_id": proofSpeakerID}, {"speaker_id": proofSpeakerID + "-untimed"}}
	}
	encoded, err := json.Marshal(attributions)
	if err != nil {
		return nil, err
	}
	cues[0]["speaker_attributions"] = encoded
	document["cues"], err = json.Marshal(cues)
	if err != nil {
		return nil, err
	}
	if variant == "unavailable" {
		delete(document, "media_timing")
	} else {
		duration := maximum + 1000
		if variant == "conflict" {
			duration = 0
		}
		if variant == "invalid_media_boundary" {
			duration = first.End - 1
		}
		document["media_timing"], err = json.Marshal(map[string]any{"duration_milliseconds": duration})
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(document)
}

func verifyConsumerReadoutPrivacy(stdout, stderr []byte) error {
	if bytes.Contains(stdout, []byte(proofSpeakerID)) || bytes.Contains(stderr, []byte(proofSpeakerID)) {
		return fmt.Errorf("consumer readout disclosed a speaker identifier")
	}
	return nil
}

func verifyConsumerInspection(payload []byte, state string, media bool, conflicts int) error {
	if bytes.Contains(payload, []byte(proofSpeakerID)) {
		return fmt.Errorf("consumer inspection disclosed a speaker identifier")
	}
	var report struct {
		Consumer struct {
			Count     int    `json:"attribution_count"`
			Timed     int    `json:"timed_attribution_count"`
			Untimed   int    `json:"untimed_attribution_count"`
			Media     bool   `json:"media_timing_present"`
			State     string `json:"media_boundary_check"`
			Conflicts int    `json:"cue_media_conflict_count"`
		} `json:"consumer_annotations"`
	}
	if err := json.Unmarshal(payload, &report); err != nil {
		return err
	}
	wantTimed := 1
	if state == "not_evaluated" {
		wantTimed = 0
	}
	if report.Consumer.Count != 2 || report.Consumer.Timed != wantTimed || report.Consumer.Untimed != 2-wantTimed || report.Consumer.Media != media || report.Consumer.State != state || report.Consumer.Conflicts != conflicts {
		return fmt.Errorf("consumer inspection counts or media boundary state differ")
	}
	return nil
}

func verifyConsumerOmissions(stderr []byte, conversion bool) error {
	if bytes.Contains(stderr, []byte(proofSpeakerID)) {
		return fmt.Errorf("omission diagnostics disclosed a speaker identifier")
	}
	prefix := "consumer_"
	if conversion {
		prefix = "conversion_"
	}
	if strings.Count(string(stderr), prefix+"speaker_attribution_omitted") != 2 || strings.Count(string(stderr), prefix+"media_timing_omitted") != 1 {
		return fmt.Errorf("consumer omissions do not report each attribution and the media declaration exactly once")
	}
	return nil
}

// Both occupied and absent destinations must survive refusal without staging.
func verifyNativeDestinationRefusal(ctx context.Context, binary, directory, input, action, target string, occupied bool, extra ...string) error {
	beforePayload, err := os.ReadFile(input)
	if err != nil {
		return err
	}
	destination := filepath.Join(directory, "consumer-refusal-output")
	if occupied {
		if err := os.WriteFile(destination, []byte("sentinel"), 0600); err != nil {
			return err
		}
	} else {
		if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	before, err := os.ReadDir(directory)
	if err != nil {
		return err
	}
	args := []string{action, input, "--output", destination, "--force"}
	if action == "restore" {
		args = append(args, "--no-metadata")
	} else {
		args = append(args, "--to", target)
	}
	args = append(args, extra...)
	status, stdout, stderr, err := invokeNative(ctx, binary, args...)
	if err != nil {
		return err
	}
	checkDestination := ""
	if occupied {
		checkDestination = destination
	}
	if err := verifyRefusal(status, stdout, stderr, checkDestination, append(deriveForbidden(directory, nil), proofSpeakerID)); err != nil {
		return err
	}
	if !occupied {
		if _, err := os.Lstat(destination); !os.IsNotExist(err) {
			return fmt.Errorf("refusal published an absent destination")
		}
	}
	after, err := os.ReadDir(directory)
	if err != nil || !sameNativeDirectoryEntries(before, after) {
		return fmt.Errorf("refusal left staging artifacts")
	}
	afterPayload, err := os.ReadFile(input)
	if err != nil || !bytes.Equal(beforePayload, afterPayload) {
		return fmt.Errorf("refusal mutated candidate input")
	}
	return nil
}

func sameNativeDirectoryEntries(before, after []os.DirEntry) bool {
	if len(before) != len(after) {
		return false
	}
	for index, entry := range before {
		if entry.Name() != after[index].Name() || entry.Type() != after[index].Type() {
			return false
		}
	}
	return true
}

func verifyConsumerNativeWorkflows(ctx context.Context, binary, repository, directory string) (ConsumerAnnotationProof, error) {
	var proof ConsumerAnnotationProof
	formats := []string{"srt", "vtt", "ass", "ssa"}
	for _, format := range formats {
		basePath := filepath.Join(directory, "source."+format+".cueson.json")
		base, err := os.ReadFile(basePath)
		if err != nil {
			return proof, err
		}
		annotated, err := consumerPayload(base, "valid")
		if err != nil {
			return proof, err
		}
		input := filepath.Join(directory, "annotated-"+format+".json")
		if err := os.WriteFile(input, annotated, 0600); err != nil {
			return proof, err
		}
		if _, err := nativeSuccess(ctx, binary, "validate", input); err != nil {
			return proof, err
		}
		inspection, err := nativeSuccess(ctx, binary, "inspect", input, "--json")
		if err != nil {
			return proof, err
		}
		if err := verifyConsumerInspection(inspection, "checked", true, 0); err != nil {
			return proof, err
		}
		textInspection, err := nativeSuccess(ctx, binary, "inspect", input)
		if err != nil {
			return proof, err
		}
		if bytes.Contains(textInspection, []byte(proofSpeakerID)) || !bytes.Contains(textInspection, []byte("attributions=2 timed=1 untimed=1")) {
			return proof, fmt.Errorf("human inspection disclosed identifiers or lost counts")
		}
		for _, target := range formats {
			action := "convert"
			if target == format {
				action = "render"
			}
			status, output, diagnostics, err := invokeNative(ctx, binary, action, input, "--to", target, "--output", "-")
			if err != nil || status != 0 || len(output) == 0 {
				return proof, fmt.Errorf("consumer %s %s to %s failed", action, format, target)
			}
			if bytes.Contains(output, []byte(proofSpeakerID)) {
				return proof, fmt.Errorf("consumer identifier escaped native output")
			}
			if err := verifyDiagnosticPrivacy(diagnostics, deriveForbidden(repository, []string{directory})); err != nil {
				return proof, err
			}
			if err := verifyConsumerOmissions(diagnostics, action == "convert"); err != nil {
				return proof, fmt.Errorf("%s %s to %s: %w", action, format, target, err)
			}
			resultPath := filepath.Join(directory, "consumer-"+format+"-"+target+"."+target)
			if err := os.WriteFile(resultPath, output, 0600); err != nil {
				return proof, err
			}
			if _, err := nativeSuccess(ctx, binary, "validate", resultPath); err != nil {
				return proof, err
			}
			if action == "render" {
				proof.MatchingRenderCount++
			} else {
				proof.ConversionCount++
			}
			for _, occupied := range []bool{true, false} {
				if err := verifyNativeDestinationRefusal(ctx, binary, directory, input, action, target, occupied, "--strict"); err != nil {
					return proof, err
				}
				proof.StrictRefusalCount++
			}
		}
		for _, variant := range []string{"invalid_timing", "invalid_identifier", "invalid_media_boundary", "unavailable", "conflict"} {
			payload, err := consumerPayload(base, variant)
			if err != nil {
				return proof, err
			}
			path := filepath.Join(directory, variant+"-"+format+".json")
			if err := os.WriteFile(path, payload, 0600); err != nil {
				return proof, err
			}
			status, stdout, stderr, err := invokeNative(ctx, binary, "validate", path)
			if err != nil {
				return proof, err
			}
			if strings.HasPrefix(variant, "invalid_") {
				if err := verifyRefusal(status, stdout, stderr, "", append(deriveForbidden(directory, nil), proofSpeakerID)); err != nil {
					return proof, err
				}
				if variant == "invalid_media_boundary" {
					if !bytes.Contains(stderr, []byte("outside declared media_timing")) {
						return proof, fmt.Errorf("media boundary input did not fail its declared media check")
					}
					proof.MediaBoundaryRefusalCount++
					for _, occupied := range []bool{true, false} {
						if err := verifyNativeDestinationRefusal(ctx, binary, directory, path, "render", format, occupied); err != nil {
							return proof, err
						}
						proof.MediaBoundaryDestinationRefusalCount++
					}
				} else if variant == "invalid_timing" {
					proof.InvalidTimingRefusalCount++
				} else {
					proof.InvalidIdentifierRefusalCount++
				}
				continue
			}
			if status != 0 || len(stdout) != 0 {
				return proof, fmt.Errorf("valid consumer variant %s refused", variant)
			}
			if err := verifyConsumerReadoutPrivacy(stdout, stderr); err != nil {
				return proof, err
			}
			inspection, err := nativeSuccess(ctx, binary, "inspect", path, "--json")
			if err != nil {
				return proof, err
			}
			state := "unavailable"
			conflicts := 0
			if variant == "conflict" {
				state = "not_evaluated"
				var document struct {
					Cues []json.RawMessage `json:"cues"`
				}
				if err := json.Unmarshal(base, &document); err != nil {
					return proof, err
				}
				conflicts = len(document.Cues)
				if strings.Count(string(stderr), "consumer_cue_media_conflict") != 1 {
					return proof, fmt.Errorf("source/media conflict warning is missing or unbounded")
				}
				proof.CueMediaConflictCount++
			}
			if err := verifyConsumerInspection(inspection, state, variant == "conflict", conflicts); err != nil {
				return proof, err
			}
			if variant == "conflict" {
				destination := filepath.Join(directory, "consumer-restored."+format)
				if _, err := nativeSuccess(ctx, binary, "restore", path, "--output", destination, "--no-metadata"); err != nil {
					return proof, err
				}
				original, err := os.ReadFile(filepath.Join(directory, "source."+format))
				if err != nil {
					return proof, err
				}
				restored, err := os.ReadFile(destination)
				if err != nil || !bytes.Equal(original, restored) {
					return proof, fmt.Errorf("consumer source/media conflict altered restoration")
				}
				proof.RestoreCount++
			}
		}
		after, err := os.ReadFile(input)
		if err != nil || !bytes.Equal(after, annotated) {
			return proof, fmt.Errorf("consumer workflows mutated annotation input")
		}
		after, err = os.ReadFile(basePath)
		if err != nil || !bytes.Equal(after, base) {
			return proof, fmt.Errorf("consumer workflows mutated base input")
		}
	}
	return proof, nil
}
