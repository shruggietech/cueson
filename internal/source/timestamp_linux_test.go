//go:build linux

package source

import (
	"errors"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/cueson/internal/model"
	"golang.org/x/sys/unix"
)

func TestOpenRegularNoFollowRejectsLinksAndDirectories(t *testing.T) {
	directory := t.TempDir()
	target := filepath.Join(directory, "target")
	if err := os.WriteFile(target, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(directory, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	if file, err := openRegularNoFollow(link, false); err == nil {
		_ = file.Close()
		t.Fatal("openRegularNoFollow accepted a symbolic link")
	}
	if file, err := openRegularNoFollow(directory, false); err == nil {
		_ = file.Close()
		t.Fatal("openRegularNoFollow accepted a directory")
	}

	file, err := openRegularNoFollow(target, false)
	if err != nil {
		t.Fatalf("open regular file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close regular file: %v", err)
	}
}

func TestCaptureNativeTimestampsUsesDescriptorAndNeverPromotesCTime(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	accessed := time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC)
	modified := time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)
	if err := unix.UtimesNano(path, []unix.Timespec{unix.NsecToTimespec(accessed.UnixNano()), unix.NsecToTimespec(modified.UnixNano())}); err != nil {
		t.Fatalf("set source timestamps: %v", err)
	}

	file, err := openRegularNoFollow(path, false)
	if err != nil {
		t.Fatalf("open source: %v", err)
	}
	defer file.Close()
	timestamps, err := captureNativeTimestamps(file)
	if err != nil {
		t.Fatalf("capture timestamps: %v", err)
	}
	if timestamps.Accessed == nil || timestamps.Accessed.UnixNS != accessed.UnixNano() {
		t.Fatalf("accessed = %#v, want %d", timestamps.Accessed, accessed.UnixNano())
	}
	if timestamps.Modified == nil || timestamps.Modified.UnixNS != modified.UnixNano() {
		t.Fatalf("modified = %#v, want %d", timestamps.Modified, modified.UnixNano())
	}
	switch timestamps.CreatedSource {
	case "birthtime":
		if timestamps.Created == nil {
			t.Fatal("birthtime provenance has no timestamp")
		}
	case "unavailable":
		if timestamps.Created != nil {
			t.Fatalf("unavailable provenance has timestamp %#v", timestamps.Created)
		}
	default:
		t.Fatalf("created_source = %q; Linux must not promote ctime", timestamps.CreatedSource)
	}
}

func TestApplyNativeTimestampsRestoresAccessAndModification(t *testing.T) {
	path := filepath.Join(t.TempDir(), "destination")
	if err := os.WriteFile(path, []byte("destination"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatalf("open destination: %v", err)
	}
	defer file.Close()

	accessed := time.Date(2022, 4, 5, 6, 7, 8, 0, time.UTC)
	modified := time.Date(2023, 5, 6, 7, 8, 9, 0, time.UTC)
	created := time.Date(2021, 1, 2, 3, 4, 5, 0, time.UTC)
	results := applyNativeTimestamps(file, model.Timestamps{
		Created:       linuxModelTimestamp(created),
		Modified:      linuxModelTimestamp(modified),
		Accessed:      linuxModelTimestamp(accessed),
		CreatedSource: "birthtime",
	})

	if len(results) != 3 {
		t.Fatalf("result count = %d, want 3", len(results))
	}
	if results[0].Kind != TimestampCreated || results[0].Status != TimestampUnsupported {
		t.Fatalf("creation result = %#v, want unsupported", results[0])
	}
	if results[1].Kind != TimestampModified || results[1].Status != TimestampRestored || results[1].EffectivePrecision != linuxTimestampPrecision {
		t.Fatalf("modification result = %#v, want restored", results[1])
	}
	if results[2].Kind != TimestampAccessed || results[2].Status != TimestampRestored || results[2].EffectivePrecision != linuxTimestampPrecision {
		t.Fatalf("access result = %#v, want restored", results[2])
	}
}

func TestApplyNativeTimestampsReportsUnavailableFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "destination")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	results := applyNativeTimestamps(file, model.Timestamps{CreatedSource: "unavailable"})
	for index, result := range results {
		if result.Status != TimestampUnavailable {
			t.Fatalf("results[%d] = %#v, want unavailable", index, result)
		}
	}
}

func TestVerifyLinuxTimestampClassifiesUnrepresentablePrecisionAsUnsupported(t *testing.T) {
	requested := &model.Timestamp{UnixNS: 1_704_067_200_123_456_789}
	actual := unix.NsecToTimespec(1_704_067_200_120_000_000)
	result := TimestampResult{Kind: TimestampModified, EffectivePrecision: linuxTimestampPrecision}

	verifyLinuxTimestamp(&result, requested, actual)

	if result.Status != TimestampUnsupported {
		t.Fatalf("status = %q, want %q", result.Status, TimestampUnsupported)
	}
	if result.EffectivePrecision != "" {
		t.Fatalf("effective precision = %q, want unknown", result.EffectivePrecision)
	}
	if !strings.Contains(result.Detail, "destination filesystem represented") {
		t.Fatalf("detail = %q, want destination precision explanation", result.Detail)
	}
}

func TestLinuxUnixNSChecksNativeRanges(t *testing.T) {
	tests := []struct {
		name        string
		seconds     int64
		nanoseconds int64
		want        int64
		wantError   bool
	}{
		{name: "epoch", want: 0},
		{name: "maximum", seconds: math.MaxInt64 / int64(time.Second), nanoseconds: math.MaxInt64 % int64(time.Second), want: math.MaxInt64},
		{name: "minimum", seconds: math.MinInt64/int64(time.Second) - 1, nanoseconds: int64(time.Second) + math.MinInt64%int64(time.Second), want: math.MinInt64},
		{name: "above maximum", seconds: math.MaxInt64/int64(time.Second) + 1, wantError: true},
		{name: "below minimum", seconds: math.MinInt64/int64(time.Second) - 2, wantError: true},
		{name: "negative nanoseconds", nanoseconds: -1, wantError: true},
		{name: "excess nanoseconds", nanoseconds: int64(time.Second), wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := linuxUnixNS(test.seconds, test.nanoseconds)
			if test.wantError {
				if err == nil {
					t.Fatalf("linuxUnixNS() = %d, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("linuxUnixNS(): %v", err)
			}
			if got != test.want {
				t.Fatalf("linuxUnixNS() = %d, want %d", got, test.want)
			}
		})
	}
}

func TestSetLinuxDescriptorTimesRejectsClosedFile(t *testing.T) {
	file, err := openRegularNoFollow(createLinuxTestFile(t), true)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	err = setLinuxDescriptorTimes(file, []unix.Timespec{{Sec: 1}, {Nsec: unix.UTIME_OMIT}})
	if err == nil {
		t.Fatal("setLinuxDescriptorTimes accepted a closed file")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("closed descriptor error was misclassified: %v", err)
	}
}

func createLinuxTestFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func linuxModelTimestamp(value time.Time) *model.Timestamp {
	return &model.Timestamp{ISO: value.Format(time.RFC3339Nano), UnixNS: value.UnixNano()}
}
