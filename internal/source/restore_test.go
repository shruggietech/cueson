package source

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/cueson/internal/model"
)

func TestRestoreExactBytesAndOverwritePolicy(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	want, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
	if err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(t.TempDir(), "restored.srt")
	report, err := Restore(context.Background(), document, RestoreOptions{Output: output, Metadata: MetadataNone})
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}
	got, err := os.ReadFile(output)
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("restored bytes = %q, error = %v", got, err)
	}
	if len(report.Assets) != 1 || report.Assets[0].Bytes != int64(len(want)) || len(report.Assets[0].Timestamps) != 0 {
		t.Fatalf("Restore() report = %#v", report)
	}
	before := []byte("preserve")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Restore(context.Background(), document, RestoreOptions{Output: output, Metadata: MetadataNone}); err == nil {
		t.Fatal("Restore() accepted existing output without force")
	}
	if got, _ := os.ReadFile(output); !bytes.Equal(got, before) {
		t.Fatalf("refused restoration changed output to %q", got)
	}
	if _, err := Restore(context.Background(), document, RestoreOptions{Output: output, Force: true, Metadata: MetadataNone}); err != nil {
		t.Fatalf("forced Restore() error = %v", err)
	}
	if got, _ := os.ReadFile(output); !bytes.Equal(got, want) {
		t.Fatalf("forced restored bytes = %q", got)
	}
	assertNoTransactionFiles(t, filepath.Dir(output))
}

func TestRestorePreflightsWholeBundle(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	companion := document.Source.Assets[0]
	companion.ID = "asset-1"
	companion.Role = "companion"
	companion.FileName = "other.srt"
	companion.Hashes.SHA256 = "0000000000000000000000000000000000000000000000000000000000000000"
	document.Source.Assets = append(document.Source.Assets, companion)
	directory := t.TempDir()
	if _, err := Restore(context.Background(), document, RestoreOptions{OutputDir: directory, Metadata: MetadataNone}); err == nil {
		t.Fatal("Restore() accepted corrupt companion")
	}
	entries, err := os.ReadDir(directory)
	if err != nil || len(entries) != 0 {
		t.Fatalf("failed preflight left entries %v, error = %v", entries, err)
	}
}

func TestRestoreMultipleAssetsToDirectory(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	companion := document.Source.Assets[0]
	companion.ID = "asset-1"
	companion.Role = "companion"
	companion.FileName = "other.srt"
	document.Source.Assets = append(document.Source.Assets, companion)
	directory := t.TempDir()
	report, err := Restore(context.Background(), document, RestoreOptions{OutputDir: directory, Metadata: MetadataNone})
	if err != nil {
		t.Fatalf("Restore() error = %v", err)
	}
	if len(report.Assets) != 2 {
		t.Fatalf("Restore() asset count = %d, want 2", len(report.Assets))
	}
	for _, name := range []string{"captions.srt", "other.srt"} {
		if _, err := os.Stat(filepath.Join(directory, name)); err != nil {
			t.Errorf("restored %q: %v", name, err)
		}
	}
	assertNoTransactionFiles(t, directory)
}

func TestRestoreMetadataPolicyAndRollback(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	instant := time.Unix(0, 1704067200123456701).UTC()
	document.Source.Assets[0].Timestamps.Created = &model.Timestamp{ISO: instant.Format(time.RFC3339Nano), UnixNS: instant.UnixNano()}
	document.Source.Assets[0].Timestamps.Modified = nil
	document.Source.Assets[0].Timestamps.Accessed = nil
	document.Source.Assets[0].Timestamps.CreatedSource = "birthtime"
	directory := t.TempDir()
	output := filepath.Join(directory, "restored.srt")

	report, err := Restore(context.Background(), document, RestoreOptions{Output: output})
	if err != nil {
		t.Fatalf("default Restore() error = %v", err)
	}
	if len(report.Warnings) == 0 || len(report.Assets[0].Timestamps) != 3 || report.Assets[0].Timestamps[0].Status != TimestampUnsupported {
		t.Fatalf("default metadata report = %#v", report)
	}
	if err := os.Remove(output); err != nil {
		t.Fatal(err)
	}
	if _, err := Restore(context.Background(), document, RestoreOptions{Output: output, Metadata: MetadataStrict}); err == nil {
		t.Fatal("strict Restore() accepted unsupported creation timestamp")
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Fatalf("strict metadata failure left new output: %v", err)
	}

	before := []byte("original destination")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Restore(context.Background(), document, RestoreOptions{Output: output, Force: true, Metadata: MetadataStrict}); err == nil {
		t.Fatal("strict forced Restore() accepted unsupported creation timestamp")
	}
	if got, err := os.ReadFile(output); err != nil || !bytes.Equal(got, before) {
		t.Fatalf("strict rollback did not preserve destination: bytes = %q, error = %v", got, err)
	}
	assertNoTransactionFiles(t, directory)
}

func TestRestoreHonorsCancellationBeforeOutput(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	output := filepath.Join(t.TempDir(), "restored.srt")
	if _, err := Restore(ctx, testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}); err == nil {
		t.Fatal("Restore() accepted canceled context")
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Fatalf("canceled restore left output: %v", err)
	}
}

func TestRestoreCleansPartialStageOnStreamingCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	directory := t.TempDir()
	output := filepath.Join(directory, "restored.srt")
	hooks := transactionHooks{fail: func(point, _ string) error {
		if point == "stage-write" {
			cancel()
		}
		return nil
	}}
	if _, err := restoreWithHooks(ctx, testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}, hooks); !errors.Is(err, context.Canceled) {
		t.Fatalf("restoreWithHooks() error = %v, want context cancellation", err)
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Fatalf("streaming cancellation left output: %v", err)
	}
	assertNoTransactionFiles(t, directory)
}

func TestRestoreTransactionFaultBoundaries(t *testing.T) {
	t.Parallel()

	fault := errors.New("injected transaction fault")
	for _, point := range []string{"stage-write", "stage-sync", "stage-close", "publish-link", "final-verify"} {
		t.Run(point, func(t *testing.T) {
			directory := t.TempDir()
			output := filepath.Join(directory, "restored.srt")
			hooks := transactionHooks{fail: func(gotPoint, _ string) error {
				if gotPoint == point {
					return fault
				}
				return nil
			}}
			if _, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}, hooks); !errors.Is(err, fault) {
				t.Fatalf("restoreWithHooks() error = %v, want injected fault", err)
			}
			if _, err := os.Lstat(output); !os.IsNotExist(err) {
				t.Fatalf("%s fault left output: %v", point, err)
			}
			assertNoTransactionFiles(t, directory)
		})
	}
}

func TestRestoreReplacementFaultPreservesDestination(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "restored.srt")
	before := []byte("preserve replacement target")
	if err := os.WriteFile(output, before, 0o600); err != nil {
		t.Fatal(err)
	}
	fault := errors.New("rename unavailable")
	hooks := transactionHooks{fail: func(point, _ string) error {
		if point == "replace-rename" {
			return fault
		}
		return nil
	}}
	if _, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: output, Force: true, Metadata: MetadataNone}, hooks); !errors.Is(err, fault) {
		t.Fatalf("restoreWithHooks() error = %v, want rename fault", err)
	}
	if got, err := os.ReadFile(output); err != nil || !bytes.Equal(got, before) {
		t.Fatalf("replacement fault changed destination: bytes = %q, error = %v", got, err)
	}
	assertNoTransactionFiles(t, directory)
}

func TestRestoreReportsRollbackAndAcceptedCleanupFaults(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "rollback.srt")
	fault := errors.New("rollback removal unavailable")
	hooks := transactionHooks{fail: func(point, _ string) error {
		if point == "final-verify" || point == "rollback-remove-new" {
			return fault
		}
		return nil
	}}
	if _, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}, hooks); err == nil || !strings.Contains(err.Error(), "rollback") {
		t.Fatalf("rollback fault error = %v, want explicit rollback failure", err)
	}
	if err := os.Remove(output); err != nil {
		t.Fatal(err)
	}

	accepted := filepath.Join(directory, "accepted.srt")
	cleanupFault := errors.New("cleanup unavailable")
	report, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: accepted, Metadata: MetadataNone}, transactionHooks{fail: func(point, _ string) error {
		if point == "cleanup-stage" {
			return cleanupFault
		}
		return nil
	}})
	if err != nil || len(report.Warnings) == 0 {
		t.Fatalf("accepted cleanup fault = (%#v, %v), want warning success", report, err)
	}
	assertNoTransactionFiles(t, directory)
}

func TestRestoreReportsFailedTransactionCleanupFault(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "failed.srt")
	fault := errors.New("injected cleanup failure")
	hooks := transactionHooks{fail: func(point, _ string) error {
		if point == "final-verify" || point == "failure-cleanup-stage" {
			return fault
		}
		return nil
	}}
	_, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}, hooks)
	if err == nil || !strings.Contains(err.Error(), "cleanup failed staging files") {
		t.Fatalf("restoreWithHooks() error = %v, want cleanup failure", err)
	}
	if _, err := os.Lstat(output); !os.IsNotExist(err) {
		t.Fatalf("failed transaction left accepted output: %v", err)
	}
	entries, readErr := os.ReadDir(directory)
	if readErr != nil || len(entries) != 1 || !strings.Contains(entries[0].Name(), ".cueson-stage-") {
		t.Fatalf("cleanup fault entries = %v, error = %v", entries, readErr)
	}
	if err := os.Remove(filepath.Join(directory, entries[0].Name())); err != nil {
		t.Fatal(err)
	}
}

func TestRestoreReportsFailedPartialStageCleanupFault(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "failed-write.srt")
	fault := errors.New("injected partial-stage cleanup failure")
	hooks := transactionHooks{fail: func(point, _ string) error {
		if point == "stage-write" || point == "failure-cleanup-stage" {
			return fault
		}
		return nil
	}}
	_, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}, hooks)
	if err == nil || !strings.Contains(err.Error(), "remove failed staging file") {
		t.Fatalf("restoreWithHooks() error = %v, want partial-stage cleanup failure", err)
	}
	entries, readErr := os.ReadDir(directory)
	if readErr != nil || len(entries) != 1 || !strings.Contains(entries[0].Name(), ".cueson-stage-") {
		t.Fatalf("partial-stage cleanup entries = %v, error = %v", entries, readErr)
	}
	if err := os.Remove(filepath.Join(directory, entries[0].Name())); err != nil {
		t.Fatal(err)
	}
}

func TestRestorePostPlanningRaceIsRuntimeFailure(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	output := filepath.Join(directory, "raced.srt")
	racedBytes := []byte("external file")
	hooks := transactionHooks{fail: func(point, path string) error {
		if point == "publish-link" {
			return os.WriteFile(path, racedBytes, 0o600)
		}
		return nil
	}}
	_, err := restoreWithHooks(context.Background(), testDocument(t), RestoreOptions{Output: output, Metadata: MetadataNone}, hooks)
	if err == nil {
		t.Fatal("restoreWithHooks() accepted raced-in destination")
	}
	var precondition *PreconditionError
	if errors.As(err, &precondition) {
		t.Fatalf("post-planning race returned PreconditionError: %v", err)
	}
	if got, readErr := os.ReadFile(output); readErr != nil || !bytes.Equal(got, racedBytes) {
		t.Fatalf("raced-in destination = %q, error = %v", got, readErr)
	}
	assertNoTransactionFiles(t, directory)
}

func assertNoTransactionFiles(t *testing.T, directory string) {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if bytes.Contains([]byte(entry.Name()), []byte(".cueson-")) {
			t.Errorf("transaction artifact remains: %s", entry.Name())
		}
	}
}
