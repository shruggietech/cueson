// Package source validates, captures, and restores lossless source envelopes.
package source

import (
	"fmt"

	"github.com/shruggietech/cueson/internal/model"
)

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
