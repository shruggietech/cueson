package testutil

import (
	"strings"
	"testing"
)

type comparisonModel struct {
	Name  string   `json:"name"`
	Lines []string `json:"lines"`
}

type comparisonDiagnostic struct {
	Code     string `json:"code"`
	Severity string `json:"severity"`
}

func TestCompareJSONUsesSemanticEquality(t *testing.T) {
	t.Parallel()

	want := comparisonModel{Name: "example", Lines: []string{"first", "second"}}
	got := map[string]any{"lines": []any{"first", "second"}, "name": "example"}
	if err := CompareJSON("fixture-json", "model", want, got); err != nil {
		t.Fatalf("CompareJSON() error = %v", err)
	}
	if err := CompareJSON("fixture-json", "model", []byte(`{"lines":[],"name":"example"}`), map[string]any{"name": "example", "lines": nil}); err == nil || err.Error() != "fixture-json model: semantic JSON mismatch" {
		t.Fatalf("CompareJSON(null versus empty) error = %v", err)
	}
	if err := CompareJSON("fixture-json", "model", []string{"first", "second"}, []string{"second", "first"}); err == nil || err.Error() != "fixture-json model: semantic JSON mismatch" {
		t.Fatalf("CompareJSON(ordered array) error = %v", err)
	}
}

func TestCompareDiagnosticsUsesExactOrder(t *testing.T) {
	t.Parallel()

	want := []comparisonDiagnostic{{Code: "first", Severity: "warning"}, {Code: "second", Severity: "error"}}
	if err := CompareDiagnostics("fixture-diagnostics", "diagnostics", want, append([]comparisonDiagnostic(nil), want...)); err != nil {
		t.Fatalf("CompareDiagnostics() error = %v", err)
	}
	got := []comparisonDiagnostic{want[1], want[0]}
	if err := CompareDiagnostics("fixture-diagnostics", "diagnostics", want, got); err == nil || err.Error() != "fixture-diagnostics diagnostics: diagnostic mismatch at index 0 (expected_length=2 observed_length=2)" {
		t.Fatalf("CompareDiagnostics(order) error = %v", err)
	}
	if err := CompareDiagnostics("fixture-diagnostics", "diagnostics", nil, []comparisonDiagnostic{}); err == nil || err.Error() != "fixture-diagnostics diagnostics: diagnostic sequence mismatch (expected_null=true observed_null=false)" {
		t.Fatalf("CompareDiagnostics(nil versus empty) error = %v", err)
	}
}

func TestCompareBytesReportsStableFirstOffset(t *testing.T) {
	t.Parallel()

	if err := CompareBytes("fixture-bytes", "source_bytes", []byte{0x00, 0xff, '\r', '\n'}, []byte{0x00, 0xff, '\r', '\n'}); err != nil {
		t.Fatalf("CompareBytes() error = %v", err)
	}
	if err := CompareBytes("fixture-bytes", "rendered_bytes", []byte("abc"), []byte("axc")); err == nil || err.Error() != "fixture-bytes rendered_bytes: byte mismatch at offset 1 (expected_length=3 observed_length=3)" {
		t.Fatalf("CompareBytes(content) error = %v", err)
	}
	if err := CompareBytes("fixture-bytes", "source_bytes", []byte("abc"), []byte("ab")); err == nil || err.Error() != "fixture-bytes source_bytes: byte mismatch at offset 2 (expected_length=3 observed_length=2)" {
		t.Fatalf("CompareBytes(length) error = %v", err)
	}
}

func TestCompareIntegrityChecksLengthAndLowercaseSHA256(t *testing.T) {
	t.Parallel()

	const digest = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if err := CompareIntegrity("fixture-integrity", "source_integrity", []byte("abc"), 3, digest); err != nil {
		t.Fatalf("CompareIntegrity() error = %v", err)
	}
	if err := CompareIntegrity("fixture-integrity", "source_integrity", []byte("abc"), 2, digest); err == nil || err.Error() != "fixture-integrity source_integrity: byte length mismatch (expected=2 observed=3)" {
		t.Fatalf("CompareIntegrity(length) error = %v", err)
	}
	if err := CompareIntegrity("fixture-integrity", "source_integrity", []byte("abc"), 3, strings.ToUpper(digest)); err == nil || err.Error() != "fixture-integrity source_integrity: expected SHA-256 must be 64 lowercase hexadecimal characters" {
		t.Fatalf("CompareIntegrity(uppercase digest) error = %v", err)
	}
	const wrongDigest = "aa7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if err := CompareIntegrity("fixture-integrity", "source_integrity", []byte("abc"), 3, wrongDigest); err == nil || err.Error() != "fixture-integrity source_integrity: SHA-256 mismatch (expected=aa7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad observed=ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad)" {
		t.Fatalf("CompareIntegrity(digest) error = %v", err)
	}
}

func TestCompareTimestampsDistinguishesExplicitOutcomes(t *testing.T) {
	t.Parallel()

	want := []TimestampOutcome{
		{Kind: "modified", Status: TimestampRestored, EffectivePrecision: "100ns"},
		{Kind: "created", Status: TimestampUnsupported},
		{Kind: "accessed", Status: TimestampUnavailable},
		{Kind: "extension", Status: TimestampFailed},
	}
	got := append([]TimestampOutcome(nil), want...)
	got[1].EffectivePrecision = "ignored for unsupported outcomes"
	if err := CompareTimestamps("fixture-timestamps", "timestamps", want, got); err != nil {
		t.Fatalf("CompareTimestamps() error = %v", err)
	}

	tests := []struct {
		name string
		got  []TimestampOutcome
		want string
	}{
		{name: "kind", got: replaceTimestamp(want, 0, TimestampOutcome{Kind: "created", Status: TimestampRestored, EffectivePrecision: "100ns"}), want: "fixture-timestamps timestamps: timestamp kind mismatch at index 0"},
		{name: "status", got: replaceTimestamp(want, 1, TimestampOutcome{Kind: "created", Status: TimestampFailed}), want: "fixture-timestamps timestamps: timestamp status mismatch at index 1 (expected=unsupported observed=failed)"},
		{name: "precision", got: replaceTimestamp(want, 0, TimestampOutcome{Kind: "modified", Status: TimestampRestored, EffectivePrecision: "1us"}), want: "fixture-timestamps timestamps: restored timestamp precision mismatch at index 0"},
		{name: "missing precision", got: replaceTimestamp(want, 0, TimestampOutcome{Kind: "modified", Status: TimestampRestored}), want: "fixture-timestamps timestamps: restored timestamp precision is not declared at index 0"},
		{name: "unknown", got: replaceTimestamp(want, 2, TimestampOutcome{Kind: "accessed", Status: "skipped"}), want: "fixture-timestamps timestamps: unknown observed timestamp status at index 2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := CompareTimestamps("fixture-timestamps", "timestamps", want, test.got); err == nil || err.Error() != test.want {
				t.Fatalf("CompareTimestamps() error = %v", err)
			}
		})
	}
}

func TestComparisonErrorsDoNotIncludeComparedPathValues(t *testing.T) {
	t.Parallel()

	const localPath = `C:\Users\operator\private\captions.srt`
	err := CompareJSON("fixture-path-free", "model", map[string]string{"value": "portable"}, map[string]string{"value": localPath})
	if err == nil {
		t.Fatal("CompareJSON() error = nil, want mismatch")
	}
	if strings.Contains(err.Error(), localPath) || err.Error() != "fixture-path-free model: semantic JSON mismatch" {
		t.Fatalf("CompareJSON() error = %q, want stable path-free evidence", err)
	}
}

func replaceTimestamp(values []TimestampOutcome, index int, value TimestampOutcome) []TimestampOutcome {
	result := append([]TimestampOutcome(nil), values...)
	result[index] = value
	return result
}
