# Research: Build Schema Foundation

## Decision: Use a maintained Draft 2020-12 validator

Use `github.com/santhosh-tekuri/jsonschema/v6` v6.0.3 and compile the embedded artifact with its canonical URI registered locally.

**Rationale**: The library explicitly supports Draft 2020-12, validates schemas against their metaschema, passes the JSON Schema test suite aside from documented optional cases, supports reusable compiled schemas, and remains pure Go. Registering the schema bytes under the canonical URI prevents runtime network dependence.

**Alternatives considered**: Hand-written structural validation was rejected because it would not establish standards conformance. Older major versions were rejected in favor of the current stable v6 release. External CLI validation was rejected because the executable must carry the contract offline.

**Primary references**: [jsonschema repository](https://github.com/santhosh-tekuri/jsonschema), [v6.0.3 release](https://github.com/santhosh-tekuri/jsonschema/releases/tag/v6.0.3)

## Decision: Keep one canonical artifact inside the embedding package

Store the canonical file at `internal/schema/cueson.schema.json` and embed that exact file.

**Rationale**: Go embedding cannot traverse parent directories. Keeping the reviewable artifact beside its internal owner avoids generated copies and byte-drift checks while honoring the no-public-Go-API boundary.

**Alternatives considered**: A root `schema/` Go package would expose an unintended public package. A canonical root file plus generated internal copy would create two mutable artifacts before release packaging exists.

## Decision: Separate structural validation from format-neutral semantics

Compile and apply JSON Schema in `internal/schema`, then decode structurally valid documents into `internal/model` and apply cross-field checks there.

**Rationale**: Required properties, types, enums, patterns, and format-conditioned shapes are declarative. Reference resolution, uniqueness, arithmetic, counts, ordering, timestamp equivalence, and release capability truth are clearer and more testable as model invariants. Neither layer performs filesystem work.

**Alternatives considered**: Encoding all cross-field rules into JSON Schema was rejected because several rules are inexpressible or obscure. Putting all checks in the schema package was rejected because it would blur the ratified model boundary.

## Decision: Defer executable source-integrity and Unicode collision enforcement

S003 defines size, SHA-256, base64, timestamps, and safe single basenames and validates their structural shapes. Issue #6 will decode bytes, compare length and digest, construct the Unicode canonical caseless collision key, capture metadata before reading, and restore exact files.

**Rationale**: The architecture assigns source bundles, hashes, collision planning, and restoration to `internal/source`. S003 must make those fields contractual without prematurely claiming the operation that proves or writes them.

**Alternatives considered**: Starting `internal/source` during S003 was rejected as cross-issue implementation. Omitting the fields was rejected because the multi-asset envelope is foundational.

## Decision: Keep restoration capability false during S003

Both recognized formats use `envelope_only` with all four capability booleans false in this slice.

**Rationale**: The schema baseline explicitly says `restore_supported` remains false until issue #6 provides the public restore command. This truth boundary overrides the broader completed-milestone example that shows restoration true.

**Alternatives considered**: Advertising the planned final v0.0.0 capability early was rejected because documentation and executable behavior would disagree.

## Decision: Stage schema file output before replacement

Write and close a same-directory temporary file before committing it to the selected destination. Refuse an existing path without `--force`; with force, replace the existing regular file with the completed temporary file in one same-directory rename operation. A failed commit leaves the destination in place, and the uncommitted temporary file is removed.

**Rationale**: The requested payload is fixed and small, but explicit overwrite must not turn a failed write into a destroyed prior file. A single replacement operation avoids the race introduced by first moving the prior destination to a disclosed or undisclosed backup path.

**Alternatives considered**: Direct truncating writes were rejected because they can corrupt an existing destination on failure. A backup-and-rollback sequence was rejected after review because another process can occupy the destination between the two renames and prevent restoration of the prior file.

## Decision: Treat output-precondition failures separately from runtime I/O

Missing `--output` values, invalid option combinations, and an already-existing destination without `--force` return exit code 2. Failures creating, writing, closing, or committing a valid requested output return exit code 1.

**Rationale**: This matches the ratified CLI distinction between invocation or deterministic pre-execution requirements and failures after a shipped operation is accepted.

**Alternatives considered**: Treating every filesystem failure as runtime was rejected because the existing-output rule is a documented precondition users can resolve before execution.
