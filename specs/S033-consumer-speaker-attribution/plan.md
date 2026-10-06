# Implementation Plan: Consumer speaker attribution and media timing

**Branch**: `codex/S033-consumer-speaker-attribution` | **Date**:2026-10-06 | **Spec**: [spec.md](spec.md)

## Summary

Deliver complete issue #82 storage, bounds, reports and export behavior in one coherent schema/model/common CLI/conversion slice. Native annotations and source assets stay authoritative for their existing roles. Stage exact1.2.0-dev and retain locally bundled immutable1.1.0/1.0.0 contracts. No release or production publication.

## Technical Context

**Language/Version**: Existing pure Go1.25+ CLI, JSON Schema2020-12, Markdown and existing Node/pnpm Site tooling.

**Primary Dependencies**: Existing jsonschema/v6, model, source, codecs, converter, local registry, CLI inspection/diagnostic boundaries and verification scripts.

**Storage**: Optional cue speaker_attributions array; optional closed root media_timing; caller data is persistent while validation state, cue conflicts and losses are runtime projections.

**Testing**: Test-first Unicode/type/interval/overflow/collection/security, current/historical schema round trips, all twelve conversions, four native renderers, CLI privacy/state/restore and focused fuzzing. Root/nested Go test/vet, formatting/docs/brand/site and CI parity. Hosted race/native Linux/macOS/security and non-publishing snapshot proof supplement local Windows execution.

**Target Platform**: Windows/macOS/Linux; six pure-Go build targets.

**Performance Goals**: Linear bounded scans over at most65536 cues and1024 entries per cue, existing input/diagnostic/loss ceilings, no uncontrolled report amplification.

**Constraints**: Exact source bytes, consumer IDs never executed/resolved, metadata privacy, no release/production/final merge, hidden Windows subprocesses.

**Scale/Scope**:20FRs,3stories,22issue acceptance criteria; one issue and one PR.

## Constitution Check

Pre-research and post-design PASS. No exceptions: source unchanged; immutable historical schema bytes; explicit optional new contract; deterministic strict omission reports; test-first boundary/security work; bounded pure-Go/offline inputs; installed Spec Kit gate and authorized publication with human merge.

## Phase 0: Research

Two independent research agents reviewed schema/history and render/conversion integration, and a third reviewed model/Unicode/media checks. Findings consolidated in [research.md](research.md). No unresolved technical question requires scope narrowing.

## Phase 1: Design

[data-model.md](data-model.md), [contracts/consumer-annotations.md](contracts/consumer-annotations.md) and [quickstart.md](quickstart.md) define exact issue semantics and observable proof.

Current consumer fields are gated to exact1.2.0-dev. Historical typed validation rejects these fields independently of schema.Decode. Scripted model/target guards accept exact historical1.1.0 and current identity without widening1.0.0.

Compute bounded immutable annotation summaries and runtime cue/media conflict observations. inspect JSON v1 receives one optional counts/state summary only when consumer fields are present; absent annotations retain established output layout. Validate/restore/render/convert report cue conflicts through their owning runtime diagnostic surfaces. No persisted diagnostics/stats rewrite.

No native target faithfully preserves this independent consumer surface. Every cue attribution occurrence receives a deterministic omission observation and media_timing receives a document observation, within existing ceilings. Matching renders use model-owned neutral diagnostics; conversion uses common loss pass across all twelve edges. Private scripted targets clear omitted consumer fields after original validation.

## Project Structure and Ownership

- Model agent: internal/model/model.go, scripted_target.go, limits.go, new consumer.go and consumer tests/fuzz.
- Schema agent: internal/schema/**, internal/version/**, .goreleaser.yaml, .github/workflows/release-proof.yml, scripts/release-verify/**, necessary current identity/conformance expectations.
- Export agent: internal/convert/** and consumer behavior in internal/codec/webvtt/render.go, internal/codec/scripted/render.go with focused tests.
- Root: internal/cli/**, internal/source consumer restoration tests, docs/README/changelog/site docs-verifier contract updates, Spec Kit artifacts, integration, final verification/publication/reviews.
- Existing named files have a single owner; agents communicate helper signatures and do not commit/push.

## Decisions and Alternatives

- Stage1.2.0-dev instead of mutating released1.1.0 or prematurely claiming stable publication. Development snapshot CI uses current schema rather than immutable stable1.1.0 package; record changed pinned process artifacts in changelog.
- Preserve arbitrary ID values and apply only explicit issue character/size bounds, rather than restricting syntax through resource/SQL/script blacklists or document-local registries.
- Keep runtime cue/media conflict observations outside persisted Diagnostics to avoid corrupting archived counts/provenance.
- Extend inspect v1 optionally, preserving absent-field behavior, rather than changing every old report or exposing IDs.
- Loss per attribution occurrence rather than aggregated per cue preserves issue-required field-specific loss accounting; refuse before report ceiling amplification rather than silently truncate.
- Source cue conflicts remain warnings; invalid attribution intervals are errors. Native output and restore share full validation but source-only conflicts do not prevent restoration.
- Custom checklist markers stay reviewer-owned and unchecked; requirements-quality criteria were assessed in analysis. Existing explicit autopilot authorization permits implementation without another stop or misleading marker changes.

## Implementation Strategy and Dependencies

Complete specify/clarify/checklist/plan/tasks and blocking read-only analysis before code. Parallel agents then implement disjoint model/schema/export surfaces with root CLI/docs integration. All work depends on model helper contract; schema/version changes are coordinated before final integration tests.

## Delivery and Verification

Record red/green focused evidence and complete foreground required checks, CI parity and hosted final-head checks. Commit with repository conventional style; no invented co-author identity. Explicit current owner request authorizes branch push and official PR, so skill pre-push halt is already satisfied. At most two bot review rounds; respond to all findings and resolve actually handled threads. Stop before final merge/tag/release/production.
