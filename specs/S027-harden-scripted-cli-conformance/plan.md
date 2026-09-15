# Implementation Plan: Harden scripted CLI conformance

**Branch**: `codex/S027-harden-scripted-cli-conformance` | **Date**: 2026-09-15 | **Spec**: [spec.md](spec.md)

**Input**: S027 completes #62 and #63 with one shared CLI discovery, privacy, conformance and hardening proof.

## Summary

Audit and finish the existing CLI catalogue integration, add an optional counts-only scripted inspection summary, remove per-mutation filesystem conversion fuzz acquisition, exercise structured ownership/integrity refusals, and close truthful native development conformance evidence. Retain stable text/historical behavior and source truth.

## Technical Context

**Language/Version**: Pure Go, minimum Go 1.25 compatibility; existing current Go security tooling.

**Primary Dependencies**: Existing model/schema/source/codec/CLI/conversion boundaries and nested verifiers; no runtime dependency addition.

**Storage**: Immutable source envelope and in-memory fuzz mutations; existing atomic filesystem acquisition/publication remains CLI-owned.

**Testing**: Focused CLI/privacy/historical/source-target tests, governed conformance/native cycles, bounded fixed-work fuzzing, full root/nested tests and vet/static/security, six builds, site/docs/artifact and hosted non-publishing release proof.

**Target Platform**: Windows/macOS/Linux with CGO_ENABLED=0 builds. All console child launchers remain hidden and noninteractive on Windows.

**Project Type**: CLI and peer exact-version schema product; all Go domain packages remain internal.

**Performance Goals**: Preserve existing 64 MiB physical/target/source limits, 65,536 physical items, 1 MiB lines, 128 declaration fields, aggregate native ceilings and 8,192 diagnostic/loss ceiling; fuzz callbacks cap mutation input and allocations before parsing.

**Constraints**: No original-byte mutation, lossy normalization, raw source disclosure, silent truncation, partial target output, historical schema mutation, public error reclassification, or premature stable claims.

**Scale/Scope**: Two independently closeable issues; three user stories; all twelve conversion directions and four shell completions.

## Constitution Check

Pre-research PASS: source preservation/privacy, exact historical schema lockstep, native/common ownership, complete losses, focused tests alongside changes, bounded portable pure-Go behavior, installed Spec Kit and explicit publication/human-merge authority.

Post-design PASS: optional scripted report preserves old report shape; independent capture guards retain the private target boundary rather than broadening public validation; fuzz source construction remains test-only and preserves truthful source/encoding/integrity. Development proof excludes stable/release/production ownership. No constitution amendment.

## Project Structure

```text
specs/S027-harden-scripted-cli-conformance/
  spec.md, plan.md, research.md, research-cli.md, research-hardening.md, research-evidence.md
  data-model.md, contracts/cli-conformance.md, contracts/hardening.md, quickstart.md, tasks.md
  checklists/requirements.md, checklists/cli-safety.md, verification.md
internal/cli/
  surface.go, validate.go, inspect.go, scripted inspection/catalogue/privacy tests and goldens
internal/model/
  scripted capture source boundary, capture regressions, structured scripted fuzz
internal/convert/
  scripted_fuzz_test.go and test-only in-memory source construction
internal/testutil/
  conformance matrix decoder and trailing-data regressions
internal/conformance/
  named scripted native workflow/safety evidence
testdata/conformance-matrix.json
.github/workflows/ci.yml
docs/cli.md, docs/architecture.md, docs/roadmap.md, docs/Cueson-Project-Specification-v0.0.0.md
CHANGELOG.md
```

**Structure Decision**: CLI agent owns only internal/cli; hardening agent owns internal/model, internal/convert test/guard changes and internal/codec/scripted/fuzz_test.go seeds/invariants; evidence agent owns internal/testutil, internal/conformance tests, and matrix metadata. Root owns artifacts, docs, CI, integration, final verification and delivery. Agent test runners have separate temp jobs and no shared launcher-file mutation.

## Phase 0: Research and decisions

Existing validate/inspect already accept scripted input through shared selection, but their discovery catalogue and invalid-token prose lag. Preserve the shared authority and extend the catalogue instead of adding another parser. Scripted inspection exposes explicit native counts without redefining WebVTT block/body fields or disclosing actors/styles/text/resource names.

Normal model validation already bounds source and declarations. Narrow guards in independent conversion-target capture inspection preflight encoded length before allocation and limit declaration splitting before field-limit refusal. This is defensive consistency, not a claimed native ingest exploit.

Replace conversion fuzz CLI/disk acquisition with test-only direct codec/model construction using exact bytes, truthful unavailable timestamps, encoding observations, current identity and valid summaries. Keep accepted controls and all strict/fatal/determinism/immutable-source/reparse assertions. Add structured scripted owner/envelope mutations and schedule existing structured conversion/safe basename boundaries.

Matrix decoding must require io.EOF after exactly one document. Replace only remaining development platform deferrals with genuine named native safety test references; keep portable grammar platform inapplicability and stable/release #64/#65 deferral. Do not change accepted fixture bytes or historical records.

Explicit alternatives: a new general inspection report version or changed old text body counts would harm compatibility without value, so retain report version 1 and omit scripted metadata outside ASS/SSA. Uniformly reclassifying all incompatible encoding failures would change established public behavior, so preserve the existing command classes. Broad source/cancellation redesign duplicates bounded proven owners and is excluded.

## Phase 1: Design

See data-model.md and both contracts. Scripted summary fields count native owners and source observations only. Shared catalogue produces help/completion tokens and invalid-format allowed values. Counts traversal is bounded by validated model limits and never copies raw source fields.

Independent capture validation checks encoded length before base64 decoding, then existing decoded length/hash/capture ownership; bounded declaration split still detects excess fields. In-memory fuzz builders parse exact bytes and validate their resulting source/model before conversion. Structured mutation clones baseline state and distinguishes accepted controls from known-invalid ownership/integrity/privacy mutations.

Conformance test performs native encode/validate/inspect/render/reparse/restore and strict/fatal output safety using governed formats; filesystem activity belongs to this explicit native integration test, not fuzz callbacks. Testutil rejects both valid second documents and malformed trailing bytes.

## Phase 2: Tasks and blocking analysis

Run installed setup-tasks and analyze prerequisites. Map all ten FRs and both issues' acceptance criteria to ordered story tasks; resolve analysis findings before product edits. US1 discovery/inspection integrates before final #63 evidence acceptance; hardening/evidence tasks with disjoint ownership may proceed together after the gate.

Checklist protocol deviation: installed checklist prerequisite requires plan.md, so setup-plan created a template scaffold before checklist generation; completed plan design follows it. Reviewer-owned checklist markers remain unchecked. Their ten criteria are covered by the recorded requirements audit; the explicitly authorized full autopilot kickoff permits continuation without rewriting markers or an inter-step halt.

## Phase 3: Verification and delivery

Run focused agents' checks then root-owned integrated full CI parity, all scheduled fixed-work fuzz budgets, six pure-Go builds and site/docs/artifact proof. Inspect convergence against every FR/story/issue criterion and all source/privacy/security invariants. Record exact outcomes in verification.md and conventional/co-author commit.

Push the authorized codex/S027 branch and publish a formatted/read-back official PR with Closes #62 and Closes #63. Move both Project stages to PR review and leave default Status unused. Await current-head checks and external reviews; address every finding and resolve only handled concerns. Request one second Codex review only if round one finds issues and do not issue a third. Human final merge is the end-of-session handoff.

## Dependencies and exclusions

#62's #61 prerequisite is closed. #63 retains its native #62 blocker and completes after integrated CLI acceptance in this slice. Epic/milestone stay open. Stable documentation/support and exact 1.1.0 candidate remain #64/#65; tag/release and production routes/deployment remain #66/#67 and require separate authority.
