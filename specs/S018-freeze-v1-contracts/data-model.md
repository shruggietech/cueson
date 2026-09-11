# Data Model: Freeze and Prove the v1 Contract

## Safety Limit Set

Represents the deterministic ceilings applied before untrusted input can amplify memory, work, diagnostics, or output.

Fields:

- `input_bytes`: 67,108,864 bytes accepted from one native or Cue JSON input.
- `document_items`: 65,536 common cues or native body items in one document.
- `item_occurrences`: 1,024 repeated settings or derived observations attached to one item.
- `diagnostics`: 8,192 ordered diagnostics accumulated for one operation.
- `losses`: 8,192 atomic conversion losses accumulated for one operation.
- `external_files`: 10,000 files admitted by one optional external-corpus run.
- `fuzz_input_bytes`: 65,536 bytes accepted by one mutation harness invocation.

Rules:

- Limits are named and shared from one owner where multiple packages rely on the same meaning.
- Exceeding a limit rejects the operation with one stable error.
- No accepted array or diagnostic list is silently truncated.
- These numeric limits are the S018 development contract and may change only through an explicit later contract decision.

## Conformance Matrix

Represents the complete mapping between documented format behavior and executable evidence.

Fields:

- `version`: Matrix contract version.
- `formats`: Ordered list of `srt` and `vtt` format entries.
- `row_id`: Stable lowercase identifier shared with maintained documentation.
- `description`: Concise behavior represented by the row.
- `evidence`: One or more evidence entries.
- `evidence.kind`: Accepted fixture, malformed fixture, render test, conversion test, fuzz target, platform test, or inapplicable.
- `evidence.reference`: Governed fixture ID or stable test reference.
- `evidence.reason`: Required only for inapplicable evidence.

Rules:

- Row identifiers are unique and appear in the matching format guide.
- Fixture references resolve exactly against `testdata/manifest.json`.
- Each row contains every evidence class required by its behavior or an explicit inapplicability reason.
- Absolute paths and local machine identifiers are forbidden.

## Schema Annotation Target

Represents one consumer-visible schema node covered by the annotation policy.

Fields:

- `pointer`: Canonical JSON Pointer from the schema root.
- `kind`: Root property, reusable definition, nested property, enumeration, or constrained scalar.
- `title_required`: Whether generated documentation benefits from a title.
- `description`: Non-empty, specific, non-tautological consumer meaning.
- `examples`: Zero or more applicable example instances.
- `validation_scope`: Exact schema fragment or complete root under which examples are validated.

Rules:

- Every root property and property directly declared by a reachable definition requires a description.
- Every enumeration and required constrained-value category has a valid example.
- Complete root examples also satisfy model semantics and source-integrity rules.
- Annotation values do not change validation semantics.

## Executable Documentation Scenario

Represents one advertised user workflow and its deterministic proof.

Fields:

- `id`: Stable scenario identity.
- `document`: Maintained document containing the command.
- `command`: Exact public command vocabulary.
- `input_fixture`: Governed fixture or generated safe temporary input.
- `expected_status`: Exit status.
- `stdout_contract`: Empty, Cue JSON, native subtitle, inspection JSON, help, or completion content.
- `stderr_contract`: Expected informational, warning, or error behavior.
- `artifact_contract`: Expected temporary output, exact restored bytes, or no file.

Rules:

- Scenarios never write into `testdata/` or another governed source directory.
- Every primary workflow is represented at least once.
- Stream assertions reject decorative or misplaced output.
- Temporary directories and machine-specific values never enter expected output.

## External Corpus Run

Represents one opt-in maintainer verification of private or non-redistributable inputs.

Fields:

- `root`: Caller-supplied runtime directory, never persisted or printed verbatim.
- `relative_identity`: Portable path relative to the supplied root.
- `format`: Detected supported format.
- `result`: Accepted, rejected, skipped unsupported, or verifier failure.
- `aggregate_counts`: Counts by format and result.

Rules:

- Links, non-regular files, traversal identities, oversized files, and excessive file counts are rejected.
- Source files are opened no-follow where supported and never modified.
- Reports contain portable relative identities only when needed for a failure.
- No external bytes or identifiers enter the governed manifest or CI artifacts.
