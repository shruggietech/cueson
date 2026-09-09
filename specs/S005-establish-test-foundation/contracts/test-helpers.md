# Shared Test Helper Contract

## Boundary

`internal/testutil` is an internal test-support package with no dependency on `internal/model`, `internal/schema`, `internal/source`, codecs, or CLI behavior. Domain tests project their values into its portable inputs.

## Fixture verification

The helper loads exactly one manifest JSON value, rejects unknown fields, validates every record, inventories governed roots without following links, and returns ordered portable records only after provenance, redistribution, path, byte-contract, size, and digest checks pass. It never modifies corpus files.

## Golden comparisons

The helper exposes separate comparison functions for semantic JSON, ordered diagnostics, exact bytes, byte length and SHA-256, and timestamp expectations. It does not combine these into a normalization snapshot and has no update-goldens mode.

Mismatch errors begin with the caller-supplied fixture ID and logical surface. Byte mismatches identify the first differing offset and expected or observed lengths. JSON comparison preserves null-versus-empty and ordered arrays. Timestamp comparison requires exact kind and status; a restored result also compares effective precision.

## Path-leak detection

The helper accepts explicit forbidden sentinels and rejects native, slash-normalized, backslash-normalized, and JSON-escaped occurrences. It does not reject generic path-looking text, URLs, or drive-like dialogue that may be legitimate payload content.

Whole runtime restoration reports are not goldenized because they contain caller-selected destinations. Cross-package conformance tests project only portable result fields.
