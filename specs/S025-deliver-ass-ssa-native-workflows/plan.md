# Implementation Plan: S025 ASS/SSA native workflows

**Branch**: `codex/S025-deliver-ass-ssa-native-workflows` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

## Summary

Close #57/#58/#59 as one independently evidenced corpus, native ingest and model-render outcome. Reuse the typed scripted model and exact source boundary, add a shared bounded native codec, then integrate current registry/CLI operations and audited corpus evidence. Historical schema/semantics remain immutable; current identity stays 1.1.0-dev.

## Technical Context

**Language/Version**: Go 1.25 compatibility floor, pure Go.

**Primary Dependencies**: Existing model/schema/source/codec/CLI boundaries and manifest v1; no new runtime dependency.

**Storage**: Verified original assets in the existing source envelope; portable manifest/projection artifacts.

**Testing**: Go unit/conformance/CLI tests, independent byte/semantic goldens, fixed-work parser/render fuzz and native CI.

**Target Platform**: Windows/macOS/Linux native execution; amd64/arm64 six-target pure-Go builds.

**Project Type**: Internal CLI codecs and public versioned Cue JSON contracts.

**Performance Goals**: Bounded linear native framing/projection/serialization with checked aggregate allocation; no superlinear duplicate validation.

**Constraints**: Existing 64 MiB native acquisition, 1 GiB Cue JSON, 65,536 combined native occurrences, 1,024 repeated values/item, 8,192 diagnostics/losses; scripted 1 MiB line, 128 declarations, 1M fields/spans, depth 32, 16/32 MiB decoded attachment bounds and 64 MiB canonical output.

**Scale/Scope**: Both ratified ASS/SSA dialects; twenty requirements and three independently closeable issue outcomes. Cross-format conversion/stable/publication remain downstream.

## Constitution Check

Pre-research PASS: exact source fidelity, common/native ownership, no silent loss, test-first security/format work, pure-Go portability, schema lockstep, no runtime identity and explicit delivery authority are covered. Post-design PASS: independent expected projections and restoration, grammar-dependent privacy, declaration transitions, bounded serialization, historical regressions and no protected release transition are retained. No constitution exception or public stable contract relaxation.

## Research and design

[Research](research.md), [data model](data-model.md), [native workflow contract](contracts/native-workflows.md), and [quickstart](quickstart.md) resolve all design unknowns.

Material departures from S024: native codecs become available under experimental observations; valid schema_only inputs remain accepted. Explicit detector rejection closes filename-fallback ambiguity. Narrow model helper APIs reuse existing authoritative validation. The conformance matrix gains honest deferred evidence for downstream assertions. These decisions are dated in the changelog.

## Project Structure

```text
specs/S025-deliver-ass-ssa-native-workflows/
  spec.md, plan.md, research.md, data-model.md, quickstart.md, tasks.md
  contracts/native-workflows.md
  checklists/requirements.md, native-fidelity.md
  verification.md
internal/codec/scripted/
  detect.go, parse.go, projection.go, render.go, *_test.go
internal/model/scripted_native.go
internal/codec/detect.go
internal/cli/
  cli.go, workflows.go, scripted_native_test.go
internal/schema/cueson.schema.json
internal/conformance/scripted_test.go
testdata/
  manifest.json, conformance-matrix.json, fixtures/, malformed/
scripts/corpus-verify/
docs/formats/ass-ssa.md, docs/architecture.md, docs/roadmap.md
```

## Execution and ownership

Setup/analysis first; then bounded parallel work on exclusive files. Parser agent owns scripted detect/parse/projection and narrow model helper wrappers, coordinating signatures before renderer/corpus consumption. Renderer agent owns scripted render and renderer tests. Corpus agent owns testdata inventory/matrix and scripted conformance expectations/tests. Coordinator owns parent registry/detection rejection, current schema/capability integration, CLI integration, optional corpus verifier, docs/changelog and all final checks. Consumers may write tests against settled interfaces while parser implementation proceeds; integration tests require all owners.

## Verification plan

Run focused tests before/alongside each implementation. Verify independent corpus projections plus exact restoration and edit/construction render/reparse, all S024 security regressions and boundaries. Full foreground parity includes root/nested tests/vet/staticcheck/govulncheck, fixed-work old and new fuzz, six builds, actionlint, text/docs/brand checks, generation drift and site unit/browser/accessibility/build/artifact/Wrangler dry-run. Hosted CI adds race and native platforms plus exact candidate package smoke; wait for every current-head check.

## Authority and checklist decisions

Installed extension configuration is absent, so before/after hooks are skipped for each invoked skill. The checklist prerequisite requires a plan file despite autopilot checklist-before-plan ordering; setup-plan produced only its template before the checklist, then full design follows. Built-in requirements checklist passes 13/13. Custom twelve-item requirements checklist is generated unchecked; autopilot assesses all as clear and proceeds under existing full-slice authority without changing reviewer markers. Explicit S025 push/official-PR authority satisfies the autopilot pre-push authorization boundary. No tag/release/production/final merge is authorized.

## Complexity Tracking

No constitution violation requires an exception. One shared parser/renderer and narrow authority wrappers avoid independent variant stacks or copied privacy logic.
