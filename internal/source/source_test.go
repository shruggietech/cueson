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

	"github.com/shruggietech/cueson/internal/model"
)

func TestCaptureContextReturnsTheExactBytesUsedByTheEnvelope(t *testing.T) {
	t.Parallel()

	payload := []byte("1\r\n00:00:00,000 --> 00:00:01,000\r\nHello\r\n")
	path := filepath.Join(t.TempDir(), "captions.srt")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}

	captured, err := CaptureContext(context.Background(), path, CaptureOptions{})
	if err != nil {
		t.Fatalf("CaptureContext() error = %v", err)
	}
	if !bytes.Equal(captured.Bytes, payload) {
		t.Fatalf("CaptureContext() bytes = %q, want %q", captured.Bytes, payload)
	}
	decoded, err := base64.StdEncoding.DecodeString(captured.Asset.DataBase64)
	if err != nil {
		t.Fatalf("decode asset: %v", err)
	}
	if !bytes.Equal(decoded, captured.Bytes) {
		t.Fatal("source envelope bytes differ from the returned codec bytes")
	}
	if captured.Asset.FileName != "captions.srt" || captured.Asset.Size.Bytes != int64(len(payload)) {
		t.Fatalf("CaptureContext() asset = %#v", captured.Asset)
	}
}

func TestCaptureContextEnforcesTheByteLimitAfterCapturingMetadata(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "bounded.srt")
	if err := os.WriteFile(path, []byte("12345"), 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, false)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	metadataCaptured := false
	_, err = captureContextFromOpenFile(context.Background(), file, "bounded.srt", CaptureOptions{}, 4, func(open *os.File) (model.Timestamps, error) {
		metadataCaptured = true
		offset, seekErr := open.Seek(0, 1)
		if seekErr != nil || offset != 0 {
			t.Fatalf("metadata callback offset = %d, error = %v", offset, seekErr)
		}
		return model.Timestamps{CreatedSource: "unavailable"}, nil
	})
	if !metadataCaptured {
		t.Fatal("metadata was not captured before the size rejection")
	}
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("captureContextFromOpenFile() error = %v, want byte-limit rejection", err)
	}
}

func TestCaptureContextAcceptsTheExactByteLimit(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "boundary.srt")
	payload := []byte("1234")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := openRegularNoFollow(path, false)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	captured, err := captureContextFromOpenFile(context.Background(), file, "boundary.srt", CaptureOptions{}, int64(len(payload)), func(*os.File) (model.Timestamps, error) {
		return model.Timestamps{CreatedSource: "unavailable"}, nil
	})
	if err != nil {
		t.Fatalf("captureContextFromOpenFile() error = %v", err)
	}
	if !bytes.Equal(captured.Bytes, payload) {
		t.Fatalf("captured bytes = %q, want %q", captured.Bytes, payload)
	}
}

func TestCaptureContextHonorsCancellation(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "cancel.srt")
	if err := os.WriteFile(path, []byte("source"), 0o600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := CaptureContext(ctx, path, CaptureOptions{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("CaptureContext() error = %v, want context cancellation", err)
	}
}

func TestCaptureContextClassifiesDeterministicPreconditions(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	oversized := filepath.Join(directory, "oversized.srt")
	file, err := os.OpenFile(oversized, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(MaxCaptureBytes + 1); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	for _, path := range []string{
		filepath.Join(directory, "missing.srt"),
		directory,
		oversized,
		filepath.Join(directory, "PRIVATE:BAD.srt"),
	} {
		if _, err := CaptureContext(context.Background(), path, CaptureOptions{}); !IsCapturePrecondition(err) {
			t.Errorf("CaptureContext(%q) error = %v, want capture precondition", path, err)
		}
	}
}

func TestReadFileContextUsesBoundedNoFollowBoundary(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	exactPath := filepath.Join(directory, "exact.json")
	if err := os.WriteFile(exactPath, []byte("1234"), 0o600); err != nil {
		t.Fatal(err)
	}
	payload, err := readFileContextWithLimit(context.Background(), exactPath, 4)
	if err != nil || !bytes.Equal(payload, []byte("1234")) {
		t.Fatalf("exact boundary = (%q, %v)", payload, err)
	}

	overPath := filepath.Join(directory, "over.json")
	if err := os.WriteFile(overPath, []byte("12345"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readFileContextWithLimit(context.Background(), overPath, 4); !IsCapturePrecondition(err) {
		t.Fatalf("limit-plus-one error = %v, want precondition", err)
	}

	for _, path := range []string{filepath.Join(directory, "missing.json"), directory, filepath.Join(directory, "PRIVATE:BAD.json")} {
		if _, err := readFileContextWithLimit(context.Background(), path, 4); !IsCapturePrecondition(err) {
			t.Errorf("readFileContextWithLimit(%q) error = %v, want precondition", path, err)
		}
	}

	linkPath := filepath.Join(directory, "linked.json")
	if err := os.Symlink(exactPath, linkPath); err == nil {
		if _, err := readFileContextWithLimit(context.Background(), linkPath, 4); !IsCapturePrecondition(err) {
			t.Errorf("symlink error = %v, want precondition", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := readFileContextWithLimit(ctx, exactPath, 4); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled read error = %v", err)
	}
}

func TestReadCueJSONContextAllowsBoundedSourceExpansion(t *testing.T) {
	directory := t.TempDir()
	expandedPath := filepath.Join(directory, "expanded.cueson.json")
	file, err := os.OpenFile(expandedPath, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(MaxCaptureBytes + 1); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	payload, err := ReadCueJSONContext(context.Background(), expandedPath)
	if err != nil || int64(len(payload)) != MaxCaptureBytes+1 {
		t.Fatalf("expanded Cue JSON boundary = (%d bytes, %v)", len(payload), err)
	}

	tooLargePath := filepath.Join(directory, "too-large.cueson.json")
	file, err = os.OpenFile(tooLargePath, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Truncate(MaxCueJSONBytes + 1); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadCueJSONContext(context.Background(), tooLargePath); !IsCapturePrecondition(err) {
		t.Fatalf("Cue JSON limit-plus-one error = %v, want precondition", err)
	}
}
