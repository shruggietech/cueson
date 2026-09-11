# Research: Freeze and Prove the v1 Contract

## Decision 1: Combine hardening, schema annotations, and documentation

**Decision**: Close issues #35, #36, and #41 in one S018 pull request.

**Rationale**: Hardening can reveal final behavioral corrections, schema annotations describe the resulting public data contract, and documentation must match the same final behavior. One review surface eliminates serial correction loops while leaving release-candidate preparation independent.

**Alternatives considered**: Three issue-sized slices were rejected as unnecessary serialization. Combining release-candidate preparation was rejected because version lockstep, immutable v1 schema admission, and candidate proof require a distinct protected review and evidence package.

## Decision 2: Bound every file input before allocation

**Decision**: Reuse one context-aware, regular-file, no-follow reader with the existing 64 MiB ceiling for native inputs and a 1 GiB ceiling for `restore` and `render` Cue JSON. `encode` rejects a Cue JSON representation larger than that consumer boundary instead of publishing a document that Cueson cannot read back.

**Rationale**: Those two commands currently call an unbounded whole-file reader while newer validated-input commands use bounded source capture. A shared acquisition boundary prevents drift and preserves stable path, cancellation, and exit-status classification.

**Alternatives considered**: Retaining unbounded Cue JSON reads was rejected as a hostile-input gap. Streaming JSON directly was deferred because it would complicate exact validation without removing the need for a bounded complete document.

## Decision 3: Reject complexity overflow without truncation

**Decision**: Retain the 64 MiB native file-input ceiling, cap Cue JSON produced by encode and consumed by exact restore or model-driven render at 1 GiB, cap common cues and native body items at 65,536 per document, cap repeated settings or derived observations at 1,024 per item, cap diagnostics and conversion losses at 8,192 per operation, cap an optional external-corpus run at 10,000 files, and guard individual fuzz inputs at 64 KiB. Exceeding a ceiling returns one stable typed failure before further amplification; accepted data is never silently truncated.

**Rationale**: The byte limit alone still permits millions of tiny objects or diagnostics. Explicit ownership-local limits make resource behavior testable and preserve the constitution's no-silent-loss rule.

**Alternatives considered**: Diagnostic truncation was rejected because it would hide source information. Relying only on process memory was rejected as non-deterministic and unsafe.

## Decision 4: Make conformance traceability data-driven

**Decision**: Add `testdata/conformance-matrix.json` with stable format-contract row identifiers, evidence classes, governed fixture references, focused test references, and explicit inapplicability reasons.

**Rationale**: Prose and scattered test names cannot prove complete row coverage. A compact machine-readable index can be checked against both documentation row identifiers and the governed fixture manifest.

**Alternatives considered**: A Markdown-only matrix was rejected because omissions would not fail CI. Embedding the matrix only in Go was rejected because maintainers and users need a readable compatibility reference.

## Decision 5: Run fixed-work fuzzing in an existing required context

**Decision**: Strengthen the required fuzz harnesses with explicit input guards and invariants, then run each named surface with a deterministic work-count budget and serial worker setting inside an existing CI job.

**Rationale**: Ordinary tests execute seed corpora but do not perform mutation. A work-count budget is more comparable than wall-clock fuzz duration and avoids introducing a new unprotected status context.

**Alternatives considered**: Nightly-only fuzzing was rejected because S018 requires pull-request evidence. Long wall-clock sessions were rejected as noisy and slow for the v1 critical path.

## Decision 6: Keep external corpora opt-in and read-only

**Decision**: Add a standalone maintainer verifier that walks a supplied corpus without following links, enforces file-count and byte limits, runs real detection, parsing, and rendering cycles, writes no sibling output, and reports only aggregate or portable relative identities.

**Rationale**: Real-world corpora are valuable but may be private or non-redistributable. The repository must not ingest their bytes, absolute paths, or identities into fixtures, logs, or CI artifacts.

**Alternatives considered**: Committing third-party corpora was rejected for provenance and size reasons. Making the verifier mandatory in hosted CI was rejected because CI has no authorized corpus.

**Implementation note**: Keep one canonical `scripts/corpus-verify` module and invoke the real CLI in process. A duplicate `scripts/external-corpus-verify` authority and platform child-process adapters were intentionally omitted because they would add conflicting implementations and unnecessary process-launch risk without providing distinct behavior.

## Decision 7: Treat schema annotations as tested non-normative data

**Decision**: Add descriptions to every public root and reachable-definition property, titles where useful, and examples for every public area, enumeration, and constrained-value class. Validate examples explicitly against compiled fragments and complete documents.

**Rationale**: JSON Schema does not validate `examples` automatically. Recursive coverage and fragment compilation prevent documentation rot while retaining all existing validation keywords and accepted or rejected outcomes.

**Alternatives considered**: Root-only descriptions were rejected because generators expose nested fields. Maintaining examples only in prose was rejected because schema-aware consumers would still lack them.

## Decision 8: Make documentation examples executable without unsafe side effects

**Decision**: Use governed inputs, temporary destinations, and captured CLI streams in root-module tests. Add static registration and vocabulary checks in the documentation verifier and package-level equality tests for help and conversion-loss references.

**Rationale**: The README currently uses nonexistent placeholders and arbitrary output paths. Direct CLI invocation in isolated temporary directories proves the promised behavior without modifying fixtures or requiring visible child processes.

**Alternatives considered**: Shell transcript tests were rejected as platform-fragile. Link-only verification was retained but is insufficient on its own.

## Decision 9: Preserve release boundaries

**Decision**: Keep software and schema identity at `0.1.0`, pin the released v0.0.0 schema SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`, and leave v1 versioning, immutable schema admission, notes, tag, artifacts, publication, milestone closure, and production hosting to later slices.

**Rationale**: S018 proves and documents the contract. It does not nominate or publish the final candidate.

**Alternatives considered**: Folding release preparation into S018 was rejected because it would blur protected authority and make any late hardening correction invalidate candidate evidence.
