package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestVerifyAcceptsSupportedInputsWithoutWriting(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "nested", "captions.srt"), "1\n00:00:00,000 --> 00:00:01,000\nhello\n")
	writeTestFile(t, filepath.Join(root, "captions.vtt"), "WEBVTT\n\n00:00.000 --> 00:01.000\nhello\n")
	writeTestFile(t, filepath.Join(root, "unknown.bin"), "not subtitles")
	before := treeSnapshot(t, root)
	report, err := Verify(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if report.Files != 3 || report.Accepted != 2 || report.Unsupported != 1 || report.Formats.SRT != 1 || report.Formats.VTT != 1 || len(report.Failures) != 0 {
		t.Fatalf("report = %#v", report)
	}
	if after := treeSnapshot(t, root); !reflect.DeepEqual(after, before) {
		t.Fatalf("corpus changed: before=%#v after=%#v", before, after)
	}
	payload, err := json.Marshal(report)
	if err != nil || strings.Contains(string(payload), filepath.ToSlash(root)) || strings.Contains(string(payload), filepath.Clean(root)) {
		t.Fatalf("report leaked root: %s, %v", payload, err)
	}
}

func TestVerifyRejectsLinksLimitsAndCancellationWithoutRootLeak(t *testing.T) {
	t.Parallel()
	t.Run("file limit", func(t *testing.T) {
		root := t.TempDir()
		writeTestFile(t, filepath.Join(root, "one.srt"), "x")
		writeTestFile(t, filepath.Join(root, "two.srt"), "x")
		_, err := verifyWithLimits(context.Background(), root, 1, 100)
		if err == nil || !strings.Contains(err.Error(), "more than 1 files") || strings.Contains(err.Error(), root) {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("byte limit", func(t *testing.T) {
		root := t.TempDir()
		writeTestFile(t, filepath.Join(root, "large.srt"), "12345")
		_, err := verifyWithLimits(context.Background(), root, 1, 4)
		if err == nil || !strings.Contains(err.Error(), "4-byte limit") || strings.Contains(err.Error(), root) {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("cancellation", func(t *testing.T) {
		root := t.TempDir()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := Verify(ctx, root)
		if err == nil || !strings.Contains(err.Error(), "canceled") || strings.Contains(err.Error(), root) {
			t.Fatalf("error = %v", err)
		}
	})
	t.Run("symbolic link", func(t *testing.T) {
		root := t.TempDir()
		target := filepath.Join(root, "target.srt")
		writeTestFile(t, target, "x")
		link := filepath.Join(root, "link.srt")
		if err := os.Symlink(target, link); err != nil {
			t.Skipf("symbolic links unavailable: %v", err)
		}
		_, err := Verify(context.Background(), root)
		if err == nil || !strings.Contains(err.Error(), "symbolic link") || strings.Contains(err.Error(), root) {
			t.Fatalf("error = %v", err)
		}
	})
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func treeSnapshot(t *testing.T, root string) map[string]int64 {
	t.Helper()
	result := make(map[string]int64)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		result[filepath.ToSlash(relative)] = info.Size()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}
