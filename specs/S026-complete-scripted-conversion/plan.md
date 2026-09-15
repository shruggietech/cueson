# Implementation Plan: Complete scripted conversion

**Branch**: `codex/S026-complete-scripted-conversion` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Input**: S026 closes #60 and #61 and verifies all twelve distinct four-format directions.

## Summary

Extend the existing side-effect-free conversion pipeline with native scripted loss analysis, deterministic scripted target construction, faithful bounded variant mapping, and a private source-aware rendering boundary. Preserve the existing two text directions and report vocabulary. Conversion command selection/help/completion and governed evidence complete the user journey; scripted codecs remain experimental.

## Technical Context

**Language/Version**: Pure Go, minimum compatibility Go 1.25.0, existing current toolchain for security tooling.

**Primary Dependencies**: Existing model, schema, source, registered codecs, conversion helpers, and CLI atomic publication; no new runtime dependency.

**Storage**: Immutable in-memory source envelope, constructed private targets, existing filesystem output transaction.

**Testing**: Go domain/conformance/CLI tests, synthetic governed fixtures, fixed-work fuzzing, nested verifier tests, vet/Staticcheck/govulncheck, six platform cross-builds, existing site test and release-proof CI.

**Target Platform**: Windows, macOS, Linux, CGO_ENABLED=0; hidden noninteractive Windows subprocess guarantees.

**Project Type**: CLI and peer schema product; internal Go packages remain private.

**Performance Goals**: Complete bounded conversion with existing 8,192 loss ceiling, 65,536 native item ceiling, 1 MiB native line ceiling, 64 MiB target ceiling, cancellation inside long traversals.

**Constraints**: No source mutation, private identity leakage, unsupported reinterpretation, silent timing change, partial target publication, schema release mutation, or unsupported stable claims.

**Scale/Scope**: Ten new directions, two established regressions, twelve verified distinct pairs; no release or production operation.

## Constitution Check

Pre-research PASS: preserve source truth; retain native owners and common semantics; compute complete losses; strict refusal; tests alongside format work; bounded portable pure-Go implementation; installed Spec Kit and blocking analysis; explicit push/PR authority and human final merge.

Post-design PASS: the constructed-target API avoids falsifying source-native identity while retaining public validation and capture checks. All loss/default/precision/fatal policies are explicit. Historical schema and released artifacts remain unchanged. No constitution exception is needed.

## Project Structure

```text
specs/S026-complete-scripted-conversion/
  spec.md, plan.md, research.md, research-outbound.md, research-targets.md, research-evidence.md
  data-model.md, quickstart.md, contracts/conversion.md, tasks.md
  checklists/requirements.md, checklists/conversion.md, verification.md
internal/convert/
  convert.go, matrix.go, project.go, render.go, report.go
  scripted_analysis.go, scripted_text.go, scripted_project.go, scripted_variant.go
  focused scripted tests and fuzz targets
internal/model/
  scripted.go, scripted_target.go, scripted_target_test.go
internal/codec/scripted/
  render.go, render_target.go, render_target_test.go
internal/cli/
  cli.go, surface.go, focused conversion tests, help/completion goldens
internal/conformance/scripted_conversion_test.go
testdata/
  fixtures/scripted-conversion/, manifest.json, conformance-matrix.json
docs/conversion.md, docs/formats/ass-ssa.md, docs/cli.md
CHANGELOG.md, .github/workflows/ci.yml
```

**Structure Decision**: Extend existing private domain boundaries; new files isolate scripted policies. The coordinating agent owns existing conversion routers/reports, CLI/docs/workflow, integration, and final verification. Outbound agent owns new analysis/text files and focused tests. Target agent owns new projection/variant files plus narrow model/codec boundary refactors/tests. Evidence agent owns new conformance/fixtures and manifest/matrix additions. No shared edit ownership.

## Phase 0: Research

Completed independent research of native loss categories, emphasis, private source validation, defaults, timing, variant semantics, and governed matrix evidence. All decisions and alternatives are consolidated in research.md and clarification entries in spec.md.

## Phase 1: Design

matrixAnalysis gains an optional private Target; existing text translations remain unchanged. Source-aware targetRenderer receives original and projected documents. For script-to-text, native analysis supplies one payload per source cue. For text-to-script, reuse established source-feature accounting while target construction avoids text-target-only blank-line or entity degradation. For variant conversion, one builder owns mapping and its atomic loss rules. Complete report validation and strict refusal precede all rendering.

ValidateScriptedTarget requires a validated original and exactly identical preserved source envelope; constructed target native validation retains all applicable invariants. RenderTarget validates original full integrity, validates constructed owners, shares the bounded owner serializer, and reparses before returning bytes. Public Validate/Render do not expose a bypass.

## Phase 2: Implementation and integration

Run installed tasks and blocking analyze first. Implement tests before or alongside policies. Work can proceed independently on disjoint files after the gate. Integrate four-format routing, extended loss source references, source integrity and nil/cancellation handling, complete semantic target validation, necessary CLI selection/help/completion, and truthful documentation.

## Phase 3: Verification and delivery

Run focused twelve-direction, native loss, precision, variant, strict/fatal atomic publication, privacy/bounds/historical tests. Run full CI parity and meaningful fixed-work conversion fuzzing; inspect convergence and all task coverage. Commit conventionally, push authorized S026 branch, publish formatted official PR closing both issues and verify exact body readback. Reconcile Project stages and unused default Status. Await all current-head CI and all external reviews. Address every finding, resolve handled threads, rerun checks, and request only one second Codex review if required. Record external completion without changing an already approved head; hand off final review and merge to the operator.

## Dependencies and exclusions

#60 depends on closed #59. #61's open #60 dependency remains native and is fulfilled through outbound analysis followed by target/variant integration in this same PR. #62/#63 and later stable/release/site issues remain open and excluded except necessary conversion surface integration. No tag/release/deployment command is authorized.
