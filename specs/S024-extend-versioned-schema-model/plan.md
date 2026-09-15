# Implementation Plan: Extend versioned schema and scripted model

**Branch:** `codex/S024-extend-versioned-schema-model` | **Date:** 2026-09-15 | **Spec:** [spec.md](spec.md)

## Summary

Deliver #55 and #56 as one complete schema/model transition: exact locally bundled historical 1.0.0 validation and current 1.1.0-dev scripted native shapes/semantics. Current output discovery remains separate from historical input identity. Root integrates shared CLI boundaries, meaningful old-input command tests, truthful nil-codec recognition, candidate packaging, maintained documentation and final foreground verification.

## Technical Context

**Language/Version:** Pure Go 1.25+; Markdown/JSON schema; Node/pnpm existing site verification.

**Primary Dependencies:** Existing jsonschema/v6 compiler, typed model, source acquisition/integrity/restore, codec registry, CLI catalogue, corpus/annotation/docs/publication verifiers, established CI/security/release proof.

**Storage:** Exact current canonical schema plus one embedded byte-identical historical 1.0.0 copy; immutable release copies remain unchanged.

**Testing:** Test-first exact identity dispatch, independent historical semantic/capability regression, positive/negative scripted examples/invariants, full old-input command matrix and no-publication/privacy cases. Foreground root/nested checks, existing fuzz budgets, native current platform, all six builds and complete site/artifact/dry-run proof. Hosted native/race/security/packaged proof before handoff.

**Target Platform:** Windows/macOS/Linux, CGO disabled; Windows child processes use existing verified hidden launch paths.

**Performance Goals:** Existing 64 MiB native/1 GiB Cue JSON acquisition bounds, 65,536 physical native occurrences, 1,024 repeated item values, 8,192 diagnostics/losses, ratified line/depth/field/tag/attachment bounds. Indexed reference validation and iterative native text checks avoid unbounded amplification.

**Scale/Scope:** Twenty requirements, three stories, two atomic GitHub outcomes. No full scripted codec/converter or stable publication.

## Constitution Check

Pre-research and post-design: PASS, no exception.

| Principle | Design application |
|---|---|
| Lossless source/privacy | Historical/source identity untouched; unknown accepted native content retained; unsafe metadata causes whole-operation rejection. |
| Schema/software discipline | Exact 1.1.0-dev lockstep; released 0.0.0/1.0.0 bytes immutable; #65 later promotes final1.1.0. |
| Common/native fidelity | Typed ordered structures and explicit derived projection, no source envelope substitute. |
| No silent loss | Declared relationships/capture construction rules; native operations unavailable until implemented. |
| Test-first behavior | Identity, model, commands, security/limits and output safety tests before/alongside code. |
| Portable bounded operation | Pure Go, maps/ordered passes, bounded decoding and no filesystem/native asset fetch/execution. |
| Documentation/delivery | Native metadata authority, installed Spec Kit analysis, current proof, explicit push/PR authority and human merge. |

## Phase 0: Research

Two bounded independent research agents inspected historical dispatch/CLI integration and native shape realization, respectively. Findings are consolidated in [research.md](research.md). No unresolved technology or interface question remains.

## Phase 1: Design

[data-model.md](data-model.md) and [contracts/input-and-capabilities.md](contracts/input-and-capabilities.md) specify exact registry authority, version semantic dispatch, ordered scripted relationships, privacy/content contexts and temporary support. [quickstart.md](quickstart.md) defines executable proof scenarios.

Current model validation isolates historical entry from scripted current rules because source/render/convert revalidate typed documents after decoding. Existing historical functions retain released behavior. New scripted native branches reuse common safety checks but have distinct native reference/projection/zero-dialogue rules.

Native Text projection, lexical scalar checks and bounded original metadata inspection are model-validation helpers reusable by later codecs; they do not constitute whole-document native ingest/render. Capture-only observations are optional for constructed recognized objects; retained uninterpreted data requires original captured content.

## Source Tree and Ownership

- Compatibility agent: internal/schema/schema.go, internal/schema/schema_test.go, internal/schema/historical/v1.0.0/cueson.schema.json and dedicated registry tests.
- Native agent: internal/model/model.go, internal/model/scripted*.go/tests, current schema JSON, representative current document, annotations tests and dedicated scripted schema tests.
- Root: internal/cli shared input/restore/render/inspect and command tests; codec nil registrations/capability integration; version and current candidate packaging/process assumptions; docs/site generation and all final verification.
- No concurrent ownership of the same file. Agents do not commit/push/publish; root integrates and owns final gates.

## Decisions and Alternatives

- Exact staged1.1.0-dev was chosen over prematurely final1.1.0 or mutating1.0.0; current candidate proof changes only necessary version/source assumptions. Dated changelog decision required for pinned packaging changes.
- Historical schema embedding was chosen over filesystem/network resolution and relaxed identity; local bytes are verified against immutable authority.
- Historical semantic dispatch was chosen over indiscriminate current validation and duplicated executables; historical native behavior remains frozen while current scripted rules evolve.
- Schema-only ASS/SSA declarations were chosen over fabricated stable/ingest/render flags or making model recognition inaccessible to generic inspect/restore.
- Nonempty scripted zero-dialogue models use null media summaries without weakening historical formats.
- Declaration record IDs own active declarations; constructed recognized objects can omit capture declarations and canonical serialization stays a later renderer concern.
- Custom requirements checklist markers remain reviewer owned. All criteria are assessed before implementation; existing explicit autopilot authorization permits proceeding after clean analysis without altering markers or treating them as code completion.

## Delivery and Verification

Model validation enforces physical-item, occurrence, field, line, depth, provenance and attachment bounds. The 64 MiB canonical native output ceiling remains a renderer publication gate owned by S025/#59 because S024 installs no scripted renderer; no fake serialization is added just to measure future output. The 16/32 MiB attachment decoded ceilings remain explicit, while the shared physical-record count and 80-byte encoded-line ceiling reject oversized bundles earlier. This verification hierarchy is recorded without weakening either ceiling.

Run installed tasks/analyze prerequisite gates before implementation. Integrate test-driven work, run CI-parity in the foreground, reconcile Issue/Project stages, and prepare a conventional commit. Explicit S024 kickoff authorizes push and official PR without another pre-push halt. Publish formatted bodies from files and verify GitHub readback. Handle every security/review finding and at most one second Codex review; final merge/release/production remain separate operator decisions.
