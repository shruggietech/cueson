// Package source validates, captures, and restores lossless source envelopes.
package source

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/shruggietech/cueson/internal/model"
)

// MaxCaptureBytes is the maximum source asset size accepted by native ingest.
const MaxCaptureBytes int64 = 64 << 20

// MetadataMode controls timestamp restoration policy.
type MetadataMode uint8

const (
	MetadataDefault MetadataMode = iota
	MetadataStrict
	MetadataNone
)

// RestoreOptions selects destinations, overwrite policy, and metadata policy.
type RestoreOptions struct {
	Output    string
	OutputDir string
	Force     bool
	Metadata  MetadataMode
}

// CaptureOptions supplies envelope fields that cannot be inferred from a file.
type CaptureOptions struct {
	ID        string
	Role      string
	MediaType *string
}

// Captured contains one lossless source asset and the exact bytes from which a
// codec must derive its model. Callers must not reopen the source path.
type Captured struct {
	Asset model.SourceAsset
	Bytes []byte
}

// CaptureContext acquires a regular source file once, records metadata before
// reading content, and returns the same bounded bytes stored in the envelope.
func CaptureContext(ctx context.Context, path string, options CaptureOptions) (Captured, error) {
	if err := ctx.Err(); err != nil {
		return Captured{}, fmt.Errorf("capture source canceled: %w", err)
	}
	if path == "" {
		return Captured{}, fmt.Errorf("source path must not be empty")
	}
	name := filepath.Base(path)
	if err := validateSafeBasename(name); err != nil {
		return Captured{}, err
	}
	file, err := openRegularNoFollow(path, false)
	if err != nil {
		return Captured{}, fmt.Errorf("open source %q: %w", path, err)
	}
	defer file.Close()
	return captureContextFromOpenFile(ctx, file, name, options, MaxCaptureBytes, captureNativeTimestamps)
}

func captureContextFromOpenFile(ctx context.Context, file *os.File, name string, options CaptureOptions, maxBytes int64, capture func(*os.File) (model.Timestamps, error)) (Captured, error) {
	if maxBytes < 0 {
		return Captured{}, fmt.Errorf("source byte limit must not be negative")
	}
	if err := ctx.Err(); err != nil {
		return Captured{}, fmt.Errorf("capture source canceled: %w", err)
	}
	timestamps, err := capture(file)
	if err != nil {
		return Captured{}, fmt.Errorf("capture source timestamps: %w", err)
	}
	before, err := file.Stat()
	if err != nil {
		return Captured{}, fmt.Errorf("inspect source before read: %w", err)
	}
	if !before.Mode().IsRegular() {
		return Captured{}, fmt.Errorf("source is not a regular file")
	}
	if before.Size() > maxBytes {
		return Captured{}, fmt.Errorf("source size %d exceeds the %d-byte limit", before.Size(), maxBytes)
	}

	var exact bytes.Buffer
	hash := sha256.New()
	limited := io.LimitReader(contextReader{ctx: ctx, reader: file}, maxBytes+1)
	count, err := io.Copy(io.MultiWriter(&exact, hash), limited)
	if err != nil {
		return Captured{}, fmt.Errorf("read source bytes: %w", err)
	}
	if count > maxBytes {
		return Captured{}, fmt.Errorf("source size exceeds the %d-byte limit", maxBytes)
	}
	after, err := file.Stat()
	if err != nil {
		return Captured{}, fmt.Errorf("inspect source after read: %w", err)
	}
	if !os.SameFile(before, after) || after.Size() != count || before.Size() != count || !before.ModTime().Equal(after.ModTime()) {
		return Captured{}, fmt.Errorf("source changed while it was being captured")
	}

	payload := exact.Bytes()
	id := options.ID
	if id == "" {
		id = "asset-0"
	}
	role := options.Role
	if role == "" {
		role = "primary"
	}
	asset := model.SourceAsset{
		ID:         id,
		Role:       role,
		FileName:   name,
		MediaType:  options.MediaType,
		Size:       model.AssetSize{Bytes: count},
		Hashes:     model.AssetHashes{SHA256: fmt.Sprintf("%x", hash.Sum(nil))},
		Timestamps: timestamps,
		DataBase64: base64.StdEncoding.EncodeToString(payload),
	}
	return Captured{Asset: asset, Bytes: append([]byte(nil), payload...)}, nil
}

// TimestampKind identifies one filesystem timestamp.
type TimestampKind string

const (
	TimestampCreated  TimestampKind = "created"
	TimestampModified TimestampKind = "modified"
	TimestampAccessed TimestampKind = "accessed"
)

// TimestampStatus is a truthful terminal metadata outcome.
type TimestampStatus string

const (
	TimestampRestored    TimestampStatus = "restored"
	TimestampUnsupported TimestampStatus = "unsupported"
	TimestampUnavailable TimestampStatus = "unavailable"
	TimestampFailed      TimestampStatus = "failed"
)

// TimestampResult reports one timestamp independently from byte restoration.
type TimestampResult struct {
	Kind               TimestampKind
	Status             TimestampStatus
	EffectivePrecision string
	Detail             string
}

// AssetResult reports verified output for one restored asset.
type AssetResult struct {
	AssetID     string
	Destination string
	Bytes       int64
	SHA256      string
	Timestamps  []TimestampResult
}

// Report is returned only after the complete bundle is accepted.
type Report struct {
	Assets   []AssetResult
	Warnings []string
}

// PreconditionError identifies deterministic destination-mode failures.
type PreconditionError struct {
	Message string
}

func (err *PreconditionError) Error() string {
	return err.Message
}

func preconditionf(format string, arguments ...any) error {
	return &PreconditionError{Message: fmt.Sprintf(format, arguments...)}
}

type preparedAsset struct {
	asset model.SourceAsset
}
