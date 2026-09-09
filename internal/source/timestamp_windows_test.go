//go:build windows

package source

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shruggietech/cueson/internal/model"
	"golang.org/x/sys/windows"
)

func TestWindowsTimestampRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "timestamps.bin")
	if err := os.WriteFile(path, []byte("cueson"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	wanted := model.Timestamps{
		Created:       windowsTestTimestamp(1704067200123456700),
		Modified:      windowsTestTimestamp(1711929600234567800),
		Accessed:      windowsTestTimestamp(1722470400345678900),
		CreatedSource: "windows_creation_time",
	}
	results := applyNativeTimestamps(file, wanted)
	for _, result := range results {
		if result.Status != TimestampRestored {
			t.Fatalf("%s status = %s (%s), want restored", result.Kind, result.Status, result.Detail)
		}
		if result.EffectivePrecision != windowsPrecision {
			t.Fatalf("%s precision = %q, want %q", result.Kind, result.EffectivePrecision, windowsPrecision)
		}
	}

	captured, err := captureNativeTimestamps(file)
	if err != nil {
		t.Fatal(err)
	}
	if captured.CreatedSource != "windows_creation_time" {
		t.Fatalf("created_source = %q", captured.CreatedSource)
	}
	assertWindowsTimestamp(t, "created", captured.Created, wanted.Created)
	assertWindowsTimestamp(t, "modified", captured.Modified, wanted.Modified)
	assertWindowsTimestamp(t, "accessed", captured.Accessed, wanted.Accessed)
}

func TestWindowsTimestampUnsupportedPrecision(t *testing.T) {
	path := filepath.Join(t.TempDir(), "precision.bin")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	results := applyNativeTimestamps(file, model.Timestamps{
		Modified:      windowsTestTimestamp(1704067200123456701),
		CreatedSource: "unavailable",
	})
	if results[0].Status != TimestampUnavailable || results[1].Status != TimestampUnsupported || results[2].Status != TimestampUnavailable {
		t.Fatalf("unexpected results: %#v", results)
	}
	if !strings.Contains(results[1].Detail, "100ns") {
		t.Fatalf("detail = %q", results[1].Detail)
	}
}

func TestWindowsTimestampRejectsCTimeFallback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "fallback.bin")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, true)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	results := applyNativeTimestamps(file, model.Timestamps{
		Created:       windowsTestTimestamp(1704067200000000000),
		CreatedSource: "ctime_fallback",
	})
	if results[0].Status != TimestampUnsupported {
		t.Fatalf("creation status = %s (%s)", results[0].Status, results[0].Detail)
	}
}

func TestOpenRegularNoFollowRejectsNonRegularEntries(t *testing.T) {
	if file, err := openRegularNoFollow(t.TempDir(), false); err == nil {
		file.Close()
		t.Fatal("directory unexpectedly accepted")
	}

	target := filepath.Join(t.TempDir(), "target.bin")
	if err := os.WriteFile(target, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(filepath.Dir(target), "link.bin")
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("Windows host cannot create a test symlink: %v", err)
	}
	if file, err := openRegularNoFollow(link, false); err == nil {
		file.Close()
		t.Fatal("reparse point unexpectedly accepted")
	}
}

func TestWindowsFiletimeConversionBoundaries(t *testing.T) {
	const tickLimit = uint64(math.MaxInt64 / filetimeQuantum)
	for _, raw := range []uint64{filetimeEpochTicks - tickLimit, filetimeEpochTicks, filetimeEpochTicks + tickLimit} {
		unixNS, err := filetimeToUnixNS(raw)
		if err != nil {
			t.Fatalf("filetimeToUnixNS(%d): %v", raw, err)
		}
		converted, err := unixNSToFiletime(unixNS)
		if err != nil {
			t.Fatalf("unixNSToFiletime(%d): %v", unixNS, err)
		}
		if got := rawFiletime(converted); got != raw {
			t.Fatalf("round trip = %d, want %d", got, raw)
		}
	}
	for _, raw := range []uint64{filetimeEpochTicks - tickLimit - 1, filetimeEpochTicks + tickLimit + 1} {
		if _, err := filetimeToUnixNS(raw); err == nil {
			t.Fatalf("filetimeToUnixNS(%d) unexpectedly succeeded", raw)
		}
	}
	if timestamp, available, err := timestampFromFiletime(windows.Filetime{}); err != nil || available || timestamp != nil {
		t.Fatalf("zero FILETIME = (%#v, %t, %v), want unavailable", timestamp, available, err)
	}
}

func windowsTestTimestamp(unixNS int64) *model.Timestamp {
	return &model.Timestamp{ISO: time.Unix(0, unixNS).UTC().Format(time.RFC3339Nano), UnixNS: unixNS}
}

func assertWindowsTimestamp(t *testing.T, name string, got, want *model.Timestamp) {
	t.Helper()
	if got == nil || want == nil || got.UnixNS != want.UnixNS || got.ISO != want.ISO {
		t.Fatalf("%s = %#v, want %#v", name, got, want)
	}
}
