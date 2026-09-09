# Research: Establish Fixture and Conformance-Test Infrastructure

## Decision: Govern source and expected artifacts through one ordered JSON manifest

Use `testdata/manifest.json` with `manifest_version: 1` and ordered fixture records. Each record contains a stable ID, purpose, classification, origin and reproducibility information, explicit redistribution approval, and an ordered artifact list. Each artifact declares a logical ID, role, portable slash-relative path, media type, byte length, lowercase SHA-256, and explicit encoding, BOM, line-ending, and final-newline characteristics.

**Rationale**: One machine-readable inventory lets tests reject missing, duplicate, orphaned, unsafe, unlicensed, or drifted corpus material. Governing expected artifacts as well as sources makes every golden change reviewable. The manifest itself cannot self-hash and remains governed repository text.

**Alternatives considered**: Per-directory metadata was rejected because completeness and duplicate detection become fragmented. Inferring license permission from project origin was rejected because contribution provenance and redistribution rights are distinct facts. Self-hashing the manifest was rejected as recursive and unverifiable.

## Decision: Preserve bytes by artifact class rather than file extension

Keep manifest, README, and normalized JSON expectations as UTF-8 LF repository text. Apply `-text -eol -whitespace` to authoritative `source`, exact expected-byte, malformed-input, and promoted fuzz-input subtrees, retaining extension rules as defense in depth. Replace the blanket `.editorconfig` testdata exemption with matching path-specific byte-area exemptions.

**Rationale**: Extension allowlists miss binary, JSON, text, or extensionless evidence, while exempting all testdata would also exempt repository-authored manifest and expectation text from encoding policy.

**Alternatives considered**: Marking all `testdata/**` as binary was rejected because it weakens text hygiene for metadata and normalized expectations. Extension-only rules were rejected because future formats and malformed cases can use arbitrary names.

## Decision: Keep reusable test utilities domain-neutral

Place manifest loading, artifact verification, semantic JSON comparison, exact byte and hash comparison, timestamp expectation comparison, and forbidden-token checks in `internal/testutil`. Return errors labeled by fixture and logical artifact IDs, never discovered absolute paths. Put model/schema/source projections and end-to-end assertions in `internal/conformance` tests.

**Rationale**: Domain-neutral utilities can be imported from schema, source, and future codec same-package tests without cycles. Separating portable values from runtime paths prevents accidental goldenization of `source.Report.Destination` or temporary directories.

**Alternatives considered**: Importing all Cueson domain packages from test utilities was rejected because same-package tests would create import cycles. A single lossy snapshot normalizer was rejected because it could erase nil-versus-empty, ordering, raw-byte, or timestamp distinctions. An automatic golden-update mode was rejected because it can silently bless path leakage or byte drift.

## Decision: Use a small synthetic Apache-2.0 seed corpus

Instantiate one accepted source-envelope case and four malformed validation cases using project-authored synthetic material. Reserve and document future SRT, WebVTT, render, conversion, and round-trip categories without creating grammar coverage evidence in S005. Treat the existing embedded schema representative as a separate official contract artifact.

**Rationale**: Synthetic seeds exercise the infrastructure without introducing third-party license uncertainty or implying native codec support. A small corpus keeps local and future CI verification deterministic and fast.

**Alternatives considered**: Importing real-world subtitle corpora was rejected until licensing and format-specific coverage are separately scoped. Duplicating the embedded representative into the manifest without a source artifact was rejected because the new case should demonstrate complete source-to-expectation governance.

## Decision: Separate structural, semantic, and integrity rejection stages

Malformed records declare one of `parse`, `structure`, `semantics`, or `integrity` plus a stable Cueson-owned diagnostic fragment. The conformance harness uses strict JSON decoding and the compiled schema for structural-only checks, typed model validation for semantic checks, and source restoration to a temporary destination for integrity rejection. It does not snapshot dependency-owned full error prose.

**Rationale**: Stable stage ownership makes regression intent clear while avoiding brittle assertions against third-party validator wording. The integrity stage must prove zero accepted output.

**Alternatives considered**: Treating every failure as generic validation was rejected because it cannot prove the separate schema, model, and source boundaries. Full error snapshots were rejected because dependency revisions can change non-contractual prose.

## Decision: Start with three pure bounded fuzz boundaries

Add fuzz targets for complete Cue JSON decoding, canonical base64 inspection, and portable safe-basename validation. Each callback rejects inputs above 64 KiB as a harness resource guard, calls no filesystem mutation, network, clock, external process, restoration, or codec behavior, and repeats accepted results where determinism is meaningful. Seed ordinary tests from representative and malformed cases, then run each mutation target for a fixed count in the foreground.

**Rationale**: These are the implemented untrusted-input boundaries with high signal and no platform side effects. A size guard protects developer runs without inventing a product input limit. Panics must surface naturally so the fuzz engine can minimize them.

**Alternatives considered**: Fuzzing restore transactions or native timestamps was rejected because filesystem state and platform variance reduce determinism. Fuzzing nonexistent SRT/WebVTT parsers was rejected as a false coverage claim. Wrapping callbacks in recovery was rejected because it hides the failure the harness must preserve.

## Decision: Promote fuzz discoveries into governed human-named regressions

Useful minimized failures move into root manifest-governed malformed or fuzz-regression cases with explicit provenance and expected rejection. Package-local automatic fuzz corpora remain temporary unless their native Go seed identity has unique value; duplicate retained copies require an explicit reason.

**Rationale**: Root cases are reviewable, portable, integrity-checked, and shared across deterministic tests. Go's automatic package-local corpus remains an execution mechanism rather than a second ungoverned source of truth.

**Alternatives considered**: Committing every automatic corpus file was rejected because it creates opaque duplicate evidence. Keeping discoveries only in test source literals was rejected because provenance and byte identity would bypass the manifest.

## Decision: Check explicit sentinels rather than banning path-looking content globally

Path-leak tests inject source root, temporary root, workspace root, username, hostname, and unique sentinel values, then reject their native, slash-normalized, backslash-normalized, and JSON-escaped forms in generated expectations. Legitimate payload text may contain URLs or path-looking dialogue and is not globally scrubbed.

**Rationale**: The product forbids local source identity, not textual discussion of paths. Explicit sentinels prove that runtime locations are absent without corrupting valid subtitle content.

**Alternatives considered**: Broad regular expressions for drive letters, slashes, or URLs were rejected because they would flag legitimate cue text. Normalizing forbidden paths to placeholders was rejected for Cue JSON model goldens because source paths must be absent, not concealed.

## Decision: Defer hosted portability proof and basename-contract reconciliation

S005 makes test infrastructure portable and runs current-host validation. Issue #8 owns hosted native execution. The existing schema counts its 255-name limit in Unicode characters while source validation uses 255 UTF-8 bytes; S005 fixture paths follow the stricter byte-safe rule but do not alter either product contract.

**Rationale**: Hosted workflows are downstream of this foundation. Changing source/schema acceptance is product behavior outside issue #7 and requires its own specification decision.

**Alternatives considered**: Adding CI now was rejected because issue #8 owns it. Encoding a fuzz differential assertion that one basename rule must equal the other was rejected because it would silently choose an unresolved product policy.

## Primary references

- Go fuzzing: https://go.dev/doc/security/fuzz/
- Go fuzz package testing: https://pkg.go.dev/testing#hdr-Fuzzing
- Git attributes: https://git-scm.com/docs/gitattributes
- SPDX license identifiers: https://spdx.org/licenses/
- Project constitution: `.specify/memory/constitution.md`
- Architecture of record: `docs/architecture.md`
- Working fixture and testing contract: `docs/Cueson-Project-Specification-v0.0.0.md`
