package source

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shruggietech/cueson/internal/model"
)

// Capture reads one regular file into a path-free lossless source asset.
func Capture(path string, options CaptureOptions) (model.SourceAsset, error) {
	if path == "" {
		return model.SourceAsset{}, fmt.Errorf("source path must not be empty")
	}
	name := filepath.Base(path)
	if err := validateSafeBasename(name); err != nil {
		return model.SourceAsset{}, err
	}
	file, err := openRegularNoFollow(path, false)
	if err != nil {
		return model.SourceAsset{}, fmt.Errorf("open source %q: %w", path, err)
	}
	defer file.Close()
	return captureFromOpenFile(file, name, options, captureNativeTimestamps)
}

func captureFromOpenFile(file *os.File, name string, options CaptureOptions, capture func(*os.File) (model.Timestamps, error)) (model.SourceAsset, error) {
	timestamps, err := capture(file)
	if err != nil {
		return model.SourceAsset{}, fmt.Errorf("capture source timestamps: %w", err)
	}
	hash := sha256.New()
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	count, err := io.Copy(io.MultiWriter(hash, encoder), file)
	if err != nil {
		return model.SourceAsset{}, fmt.Errorf("read source bytes: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return model.SourceAsset{}, fmt.Errorf("encode source bytes: %w", err)
	}
	id := options.ID
	if id == "" {
		id = "asset-0"
	}
	role := options.Role
	if role == "" {
		role = "primary"
	}
	return model.SourceAsset{
		ID: id, Role: role, FileName: name, MediaType: options.MediaType,
		Size: model.AssetSize{Bytes: count}, Hashes: model.AssetHashes{SHA256: fmt.Sprintf("%x", hash.Sum(nil))},
		Timestamps: timestamps, DataBase64: encoded.String(),
	}, nil
}
