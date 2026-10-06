package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestConsumerInspectionCountsWithoutIdentityOrSourceMutation(t *testing.T) {
	input := inspectTestWebVTTInput(t)
	input.document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "consumer-identity-sentinel"}}
	input.document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 0}
	before, _ := json.Marshal(input.document)
	report, err := buildInspectionReport(input)
	if err != nil {
		t.Fatal(err)
	}
	if report.InspectReportVersion != "1" || report.ConsumerAnnotations == nil || report.ConsumerAnnotations.AttributionCount != 1 || report.ConsumerAnnotations.UntimedAttributionCount != 1 || report.ConsumerAnnotations.MediaBoundaryCheck != "not_evaluated" || report.ConsumerAnnotations.CueMediaConflictCount != 1 {
		t.Fatalf("summary = %#v", report.ConsumerAnnotations)
	}
	for _, render := range []func(inspectionReport) ([]byte, error){marshalInspectionReport, renderHumanInspection} {
		payload, err := render(report)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(payload, []byte("consumer-identity-sentinel")) {
			t.Fatal("inspection exposed consumer identity")
		}
	}
	after, _ := json.Marshal(input.document)
	if !bytes.Equal(before, after) {
		t.Fatal("inspection changed source document")
	}
}

func TestConsumerCommandsReportConflictAndRestoreExactBytes(t *testing.T) {
	input := inspectTestWebVTTInput(t)
	input.document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "opaque-consumer-id"}}
	input.document.MediaTiming = &model.MediaTiming{DurationMilliseconds: 0}
	payload, err := marshalDocument(input.document, false)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	path := filepath.Join(directory, "input.cueson.json")
	if err := os.WriteFile(path, payload, 0600); err != nil {
		t.Fatal(err)
	}
	restored := filepath.Join(directory, "restored.vtt")
	commands := [][]string{
		{"validate", path}, {"inspect", path, "--json"},
		{"restore", path, "--output", restored, "--no-metadata"},
		{"render", path, "--to", "vtt"}, {"convert", path, "--to", "srt"},
	}
	for _, args := range commands {
		var stdout, stderr bytes.Buffer
		status := Run(context.Background(), args, nil, &stdout, &stderr)
		if status != ExitSuccess {
			t.Fatalf("%v status=%d stderr=%s", args, status, &stderr)
		}
		if !strings.Contains(stderr.String(), "consumer_cue_media_conflict") {
			t.Fatalf("%v missed runtime warning: %s", args, &stderr)
		}
		if strings.Contains(stdout.String(), "opaque-consumer-id") || strings.Contains(stderr.String(), "opaque-consumer-id") {
			t.Fatalf("%v exposed consumer ID", args)
		}
	}
	want, _ := base64.StdEncoding.DecodeString(input.document.Source.Assets[0].DataBase64)
	got, err := os.ReadFile(restored)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("restore bytes changed: %v", err)
	}
	for _, args := range [][]string{{"render", path, "--to", "vtt", "--strict"}, {"convert", path, "--to", "srt", "--strict"}} {
		var stdout, stderr bytes.Buffer
		if status := Run(context.Background(), args, nil, &stdout, &stderr); status != ExitRuntimeFailure || stdout.Len() != 0 {
			t.Fatalf("strict %v status=%d output=%q", args, status, &stdout)
		}
	}
}

func TestConsumerValidationReportsUnavailableMediaBoundary(t *testing.T) {
	input := inspectTestWebVTTInput(t)
	input.document.Cues[0].SpeakerAttributions = []model.SpeakerAttribution{{SpeakerID: "caller-id"}}
	payload, err := marshalDocument(input.document, false)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "input.json")
	if err := os.WriteFile(path, payload, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if status := Run(context.Background(), []string{"validate", path}, nil, &stdout, &stderr); status != ExitSuccess || !strings.Contains(stderr.String(), "media boundary check: unavailable") {
		t.Fatalf("status=%d stderr=%s", status, &stderr)
	}
}
