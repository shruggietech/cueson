package source

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestCaptureObtainsMetadataBeforeContentRead(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "capture.bin")
	want := []byte("source bytes")
	if err := os.WriteFile(path, want, 0o600); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	called := false
	asset, err := captureFromOpenFile(file, filepath.Base(path), CaptureOptions{}, func(file *os.File) (model.Timestamps, error) {
		called = true
		offset, seekErr := file.Seek(0, 1)
		if seekErr != nil || offset != 0 {
			t.Fatalf("capture callback offset = %d, error = %v", offset, seekErr)
		}
		return model.Timestamps{CreatedSource: "unavailable"}, nil
	})
	if err != nil || !called {
		t.Fatalf("captureFromOpenFile() = (%#v, %v), callback = %t", asset, err, called)
	}
	decoded, err := base64.StdEncoding.DecodeString(asset.DataBase64)
	if err != nil || !bytes.Equal(decoded, want) {
		t.Fatalf("captured bytes = %q, error = %v", decoded, err)
	}
	if asset.FileName != "capture.bin" || asset.ID != "asset-0" || asset.Role != "primary" {
		t.Fatalf("captured identity = %#v", asset)
	}
}

func TestCaptureRejectsUnsafeOrNonRegularSource(t *testing.T) {
	t.Parallel()

	if _, err := Capture(t.TempDir(), CaptureOptions{}); err == nil {
		t.Fatal("Capture() accepted directory")
	}
}
