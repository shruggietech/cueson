//go:build darwin

package source

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shruggietech/cueson/internal/model"
)

func TestDarwinOpenRegularNoFollowRejectsSymlinkAndDirectory(t *testing.T) {
	directory := t.TempDir()
	regularPath := filepath.Join(directory, "regular.srt")
	if err := os.WriteFile(regularPath, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	regular, err := openRegularNoFollow(regularPath, false)
	if err != nil {
		t.Fatalf("open regular file: %v", err)
	}
	if err := regular.Close(); err != nil {
		t.Fatalf("close regular file: %v", err)
	}

	symlinkPath := filepath.Join(directory, "source-link.srt")
	if err := os.Symlink(regularPath, symlinkPath); err != nil {
		t.Fatal(err)
	}
	if file, err := openRegularNoFollow(symlinkPath, false); err == nil {
		_ = file.Close()
		t.Fatal("openRegularNoFollow accepted a symbolic link")
	}
	if file, err := openRegularNoFollow(directory, false); err == nil {
		_ = file.Close()
		t.Fatal("openRegularNoFollow accepted a directory")
	}
}

func TestDarwinCaptureNativeTimestampsUsesTruthfulCreationProvenance(t *testing.T) {
	path := filepath.Join(t.TempDir(), "capture.srt")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, false)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	timestamps, err := captureNativeTimestamps(file)
	if err != nil {
		t.Fatalf("captureNativeTimestamps: %v", err)
	}
	if timestamps.Modified == nil || timestamps.Accessed == nil {
		t.Fatalf("capture = %#v, want modification and access timestamps", timestamps)
	}
	switch timestamps.CreatedSource {
	case "birthtime":
		if timestamps.Created == nil {
			t.Fatal("birthtime provenance has no created timestamp")
		}
	case "ctime_fallback":
		if timestamps.Created == nil {
			t.Fatal("ctime fallback provenance has no created timestamp")
		}
	case "unavailable":
		if timestamps.Created != nil {
			t.Fatal("unavailable creation provenance has a created timestamp")
		}
	default:
		t.Fatalf("created_source = %q, want truthful Darwin provenance", timestamps.CreatedSource)
	}
}

func TestDarwinApplyNativeTimestampsRestoresDescriptorTimes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "restore.srt")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	created := timestampForDarwinTest(time.Date(2026, 9, 9, 12, 0, 0, 123456000, time.UTC))
	modified := timestampForDarwinTest(time.Date(2026, 9, 9, 12, 1, 0, 654321000, time.UTC))
	accessed := timestampForDarwinTest(time.Date(2026, 9, 9, 12, 2, 0, 111222000, time.UTC))
	results := applyNativeTimestamps(file, model.Timestamps{
		Created:       created,
		Modified:      modified,
		Accessed:      accessed,
		CreatedSource: "birthtime",
	})

	if len(results) != 3 {
		t.Fatalf("result count = %d, want 3", len(results))
	}
	if results[0].Kind != TimestampCreated || results[0].Status != TimestampUnsupported {
		t.Fatalf("created result = %#v, want unsupported", results[0])
	}
	for index, kind := range []TimestampKind{TimestampModified, TimestampAccessed} {
		result := results[index+1]
		if result.Kind != kind || result.Status != TimestampRestored || result.EffectivePrecision != darwinTimestampPrecision {
			t.Fatalf("%s result = %#v, want restored at %s", kind, result, darwinTimestampPrecision)
		}
	}
}

func TestDarwinApplyNativeTimestampsReportsUnavailableValues(t *testing.T) {
	path := filepath.Join(t.TempDir(), "unavailable.srt")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	results := applyNativeTimestamps(file, model.Timestamps{CreatedSource: "unavailable"})
	if len(results) != 3 {
		t.Fatalf("result count = %d, want 3", len(results))
	}
	for _, result := range results {
		if result.Status != TimestampUnavailable {
			t.Fatalf("result = %#v, want unavailable", result)
		}
	}
}

func timestampForDarwinTest(instant time.Time) *model.Timestamp {
	instant = instant.UTC()
	return &model.Timestamp{ISO: instant.Format(time.RFC3339Nano), UnixNS: instant.UnixNano()}
}
