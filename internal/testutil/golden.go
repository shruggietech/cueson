package testutil

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

const (
	TimestampRestored    = "restored"
	TimestampUnsupported = "unsupported"
	TimestampUnavailable = "unavailable"
	TimestampFailed      = "failed"
)

// TimestampOutcome is the portable portion of one platform timestamp result.
// Runtime paths and diagnostic detail intentionally do not belong in a golden.
type TimestampOutcome struct {
	Kind               string `json:"kind"`
	Status             string `json:"status"`
	EffectivePrecision string `json:"effective_precision,omitempty"`
}

// CompareJSON compares JSON values semantically. Byte slices and RawMessages
// are decoded as JSON; other values are marshaled before comparison. Object key
// order is ignored while array order and null-versus-empty distinctions remain.
func CompareJSON(fixtureID, surface string, expected, observed any) error {
	expectedValue, err := semanticJSON(expected)
	if err != nil {
		return comparisonError(fixtureID, surface, "expected value is not valid JSON")
	}
	observedValue, err := semanticJSON(observed)
	if err != nil {
		return comparisonError(fixtureID, surface, "observed value is not valid JSON")
	}
	if !reflect.DeepEqual(expectedValue, observedValue) {
		return comparisonError(fixtureID, surface, "semantic JSON mismatch")
	}
	return nil
}

// CompareDiagnostics compares a diagnostic sequence exactly and in order.
func CompareDiagnostics[T any](fixtureID, surface string, expected, observed []T) error {
	if (expected == nil) != (observed == nil) {
		return comparisonError(
			fixtureID,
			surface,
			"diagnostic sequence mismatch (expected_null=%t observed_null=%t)",
			expected == nil,
			observed == nil,
		)
	}
	sharedLength := len(expected)
	if len(observed) < sharedLength {
		sharedLength = len(observed)
	}
	for index := 0; index < sharedLength; index++ {
		expectedValue, expectedErr := semanticJSON(expected[index])
		observedValue, observedErr := semanticJSON(observed[index])
		if expectedErr != nil || observedErr != nil || !reflect.DeepEqual(expectedValue, observedValue) {
			return comparisonError(
				fixtureID,
				surface,
				"diagnostic mismatch at index %d (expected_length=%d observed_length=%d)",
				index,
				len(expected),
				len(observed),
			)
		}
	}
	if len(expected) != len(observed) {
		return comparisonError(
			fixtureID,
			surface,
			"diagnostic mismatch at index %d (expected_length=%d observed_length=%d)",
			sharedLength,
			len(expected),
			len(observed),
		)
	}
	return nil
}

// CompareBytes compares authoritative bytes without normalization. Mismatch
// evidence reports only a stable offset and lengths, never payload content.
func CompareBytes(fixtureID, surface string, expected, observed []byte) error {
	if bytes.Equal(expected, observed) {
		return nil
	}
	offset := firstDifferentOffset(expected, observed)
	return comparisonError(
		fixtureID,
		surface,
		"byte mismatch at offset %d (expected_length=%d observed_length=%d)",
		offset,
		len(expected),
		len(observed),
	)
}

// CompareIntegrity verifies the declared byte length and lowercase SHA-256 of
// an observed byte sequence without decoding or rewriting it.
func CompareIntegrity(fixtureID, surface string, observed []byte, expectedLength int64, expectedSHA256 string) error {
	if int64(len(observed)) != expectedLength {
		return comparisonError(
			fixtureID,
			surface,
			"byte length mismatch (expected=%d observed=%d)",
			expectedLength,
			len(observed),
		)
	}
	if !isLowerSHA256(expectedSHA256) {
		return comparisonError(fixtureID, surface, "expected SHA-256 must be 64 lowercase hexadecimal characters")
	}
	observedDigest := sha256.Sum256(observed)
	observedSHA256 := hex.EncodeToString(observedDigest[:])
	if observedSHA256 != expectedSHA256 {
		return comparisonError(
			fixtureID,
			surface,
			"SHA-256 mismatch (expected=%s observed=%s)",
			expectedSHA256,
			observedSHA256,
		)
	}
	return nil
}

// CompareTimestamps compares timestamp kind and status in order. Effective
// precision participates only when both sides truthfully claim restoration.
func CompareTimestamps(fixtureID, surface string, expected, observed []TimestampOutcome) error {
	if len(expected) != len(observed) {
		return comparisonError(
			fixtureID,
			surface,
			"timestamp count mismatch (expected=%d observed=%d)",
			len(expected),
			len(observed),
		)
	}
	for index := range expected {
		if !knownTimestampStatus(expected[index].Status) {
			return comparisonError(fixtureID, surface, "unknown expected timestamp status at index %d", index)
		}
		if !knownTimestampStatus(observed[index].Status) {
			return comparisonError(fixtureID, surface, "unknown observed timestamp status at index %d", index)
		}
		if expected[index].Kind != observed[index].Kind {
			return comparisonError(fixtureID, surface, "timestamp kind mismatch at index %d", index)
		}
		if expected[index].Status != observed[index].Status {
			return comparisonError(
				fixtureID,
				surface,
				"timestamp status mismatch at index %d (expected=%s observed=%s)",
				index,
				expected[index].Status,
				observed[index].Status,
			)
		}
		if expected[index].Status == TimestampRestored {
			if expected[index].EffectivePrecision == "" || observed[index].EffectivePrecision == "" {
				return comparisonError(fixtureID, surface, "restored timestamp precision is not declared at index %d", index)
			}
			if expected[index].EffectivePrecision != observed[index].EffectivePrecision {
				return comparisonError(fixtureID, surface, "restored timestamp precision mismatch at index %d", index)
			}
		}
	}
	return nil
}

func semanticJSON(value any) (any, error) {
	var payload []byte
	switch typed := value.(type) {
	case []byte:
		payload = typed
	case json.RawMessage:
		payload = typed
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		payload = encoded
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	var decoded any
	if err := decoder.Decode(&decoded); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("multiple JSON values")
		}
		return nil, err
	}
	return decoded, nil
}

func firstDifferentOffset(expected, observed []byte) int {
	limit := len(expected)
	if len(observed) < limit {
		limit = len(observed)
	}
	for index := 0; index < limit; index++ {
		if expected[index] != observed[index] {
			return index
		}
	}
	return limit
}

func isLowerSHA256(value string) bool {
	if len(value) != sha256.Size*2 {
		return false
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

func knownTimestampStatus(status string) bool {
	switch status {
	case TimestampRestored, TimestampUnsupported, TimestampUnavailable, TimestampFailed:
		return true
	default:
		return false
	}
}

func comparisonError(fixtureID, surface, format string, arguments ...any) error {
	return fmt.Errorf("%s %s: %s", fixtureID, surface, fmt.Sprintf(format, arguments...))
}
