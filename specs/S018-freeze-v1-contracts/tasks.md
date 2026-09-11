# Tasks: Freeze and Prove the v1 Contract

**Input**: Design documents from `specs/S018-freeze-v1-contracts/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: S018 explicitly requires test-first hostile-input, schema-annotation, conformance, fuzz, documentation, and release-readiness evidence.

**Organization**: Tasks are grouped by user story so schema, conformance, and documentation work can proceed with isolated file ownership after the shared safety foundation.

## Phase 1: Setup and Governance

**Purpose**: Establish the governed slice and its traceability before implementation.

- [x] T001 Register issues #35, #36, and #41 as slice S018 at the Specced stage in the `cueson Delivery` Project and retain their native parent and dependency relationships
- [x] T002 Validate the S018 specification checklist and persist active feature discovery in `.specify/feature.json` and `specs/S018-freeze-v1-contracts/checklists/requirements.md`

---

## Phase 2: Foundational Safety Boundaries

**Purpose**: Close shared untrusted-input amplification gaps before story-specific proof.

**CRITICAL**: User story implementation depends on these common limits and acquisition rules.

- [x] T003 Add failing exact-limit, limit-plus-one, symlink, non-regular, cancellation, and no-output tests for shared Cue JSON acquisition in `internal/source/source_test.go`, `internal/cli/cli_test.go`, and `internal/cli/input_test.go`
- [x] T004 Implement one context-aware bounded regular-file no-follow reader and route restore, render, validate, inspect, and conversion Cue JSON acquisition through it in `internal/source/source.go`, `internal/cli/cli.go`, `internal/cli/input.go`, and `internal/cli/workflows.go`
- [x] T005 Add failing amplification-limit tests for model collections, native body items, repeated settings and observations, diagnostics, and conversion losses in `internal/model/model_test.go`, `internal/codec/subrip/parse_test.go`, `internal/codec/webvtt/parse_test.go`, and `internal/convert/report_test.go`
- [x] T006 Implement named deterministic complexity ceilings with typed rejection and no truncation in `internal/model/model.go`, `internal/codec/subrip/parse.go`, `internal/codec/webvtt/parse.go`, `internal/codec/webvtt/settings.go`, `internal/codec/webvtt/markup.go`, and `internal/convert/report.go`

**Checkpoint**: Every public file input and amplification surface has a deterministic outer bound.

---

## Phase 3: User Story 1 - Trust the complete v1 workflow (Priority: P1) MVP

**Goal**: Prove complete format, source-fidelity, conversion, hostile-input, fuzz, platform, and security behavior.

**Independent Test**: Run the conformance matrix, governed corpus, hostile regressions, fixed-work fuzz targets, optional corpus verifier tests, native jobs, static builds, and all quality and security modules.

### Tests for User Story 1

- [x] T007 [P] [US1] Add the versioned machine-readable format evidence index and failing matrix validation in `testdata/conformance-matrix.json` and `internal/conformance/matrix_test.go`
- [x] T008 [P] [US1] Add representative governed SubRip and WebVTT malformed, boundary, and fuzz regression payloads with complete provenance and hashes in `testdata/manifest.json`, `testdata/malformed/`, and `testdata/fuzz/`
- [x] T009 [P] [US1] Add hostile size, encoding, base64, timestamp, path, markup, diagnostic, and output-privacy integration regressions in `internal/conformance/hostile_test.go`
- [x] T010 [P] [US1] Add all-accepted-SubRip CLI encode and exact-restore byte-identity coverage in `internal/cli/cli_test.go`
- [x] T011 [P] [US1] Strengthen actual detection, SubRip block and timing, WebVTT block and settings, markup, source-envelope, and parser-renderer-cycle fuzz invariants in `internal/codec/fuzz_test.go`, `internal/codec/subrip/fuzz_test.go`, `internal/codec/webvtt/fuzz_test.go`, `internal/source/fuzz_test.go`, and `internal/convert/fuzz_test.go`
- [x] T012 [P] [US1] Add failing read-only, no-follow, bounded, portable-output, and no-write tests for the optional maintainer verifier in `scripts/corpus-verify/verify_test.go`

### Implementation for User Story 1

- [x] T013 [US1] Implement and validate conformance matrix row, fixture, test, fuzz, applicability, and documentation traceability in `internal/conformance/matrix_test.go`, `testdata/conformance-matrix.json`, `docs/formats/srt.md`, and `docs/formats/webvtt.md`
- [x] T014 [US1] Implement the opt-in external corpus verifier and usage contract in `scripts/corpus-verify/go.mod`, `scripts/corpus-verify/main.go`, `scripts/corpus-verify/verify.go`, in-process CLI integration, and `testdata/README.md`
- [x] T015 [US1] Document the complete bidirectional loss, fatal-boundary, and strict-mode matrix and enforce equality with runtime codes in `docs/conversion.md` and `internal/convert/documentation_test.go`
- [x] T016 [US1] Run every required fuzz surface at a fixed work count and verify all repository modules under existing quality and security contexts in `.github/workflows/ci.yml`

**Checkpoint**: User Story 1 independently proves the complete implemented v1 workflow under supported and hostile conditions.

---

## Phase 4: User Story 2 - Discover the schema contract (Priority: P2)

**Goal**: Make every consumer-facing Cue JSON semantic discoverable and keep annotation coverage executable.

**Independent Test**: Recursively validate description, title, and example coverage; compile fragment examples; validate complete examples semantically; compare embedded and emitted bytes; and verify the released schema digest.

### Tests for User Story 2

- [x] T017 [P] [US2] Add failing recursive annotation coverage, title-policy, fragment-example, complete-example, unsafe-value, stale-version, and immutable-digest tests in `internal/schema/annotations_test.go`

### Implementation for User Story 2

- [x] T018 [US2] Enrich every public root property, reachable definition, direct property, enumeration, and constrained-value category without altering validation semantics in `internal/schema/cueson.schema.json`
- [x] T019 [US2] Document the non-normative annotation policy, v1 contract surface, native extensions, support states, source authority, compatibility, and immutable-copy rules in `docs/schema.md`

**Checkpoint**: User Story 2 independently provides a complete, valid, embedded machine-readable schema reference.

---

## Phase 5: User Story 3 - Install and use documented v1 behavior (Priority: P3)

**Goal**: Make installation, workflows, CLI behavior, compatibility, and release transition accurate and executable.

**Independent Test**: Execute every registered workflow in temporary directories, compare CLI vocabulary and loss codes, validate maintained links and matrices, and reject stale claims or encoding corruption.

### Tests for User Story 3

- [x] T020 [P] [US3] Add failing isolated encode, restore, render, convert, validate, inspect, and completion documentation scenarios in `internal/cli/documentation_test.go`
- [x] T021 [P] [US3] Add failing maintained-document inventory, registered-example, schema-marker, compatibility, stale-claim, and format-matrix checks in `scripts/docs-verify/verify_test.go`
- [x] T022 [P] [US3] Add failing help and documentation vocabulary equality checks against the ordered CLI surface in `internal/cli/documentation_test.go`

### Implementation for User Story 3

- [x] T023 [US3] Replace placeholder commands with isolated governed-fixture quick starts and separate published v0.0.0 installation from current source and future v1 artifacts in `README.md`
- [x] T024 [US3] Complete the command, option, alias, stream, status, overwrite, strict, completion, and example reference in `docs/cli.md`
- [x] T025 [US3] Add the v1 public compatibility promise and reconcile architecture, security, contribution, fixture, and release-transition guidance in `docs/compatibility.md`, `docs/architecture.md`, `SECURITY.md`, `CONTRIBUTING.md`, `testdata/README.md`, and `docs/release-process.md`
- [x] T026 [US3] Implement maintained-document, example-registration, schema/reference, matrix-linkage, and stale-claim verification in `scripts/docs-verify/verify.go` and `scripts/docs-verify/main.go`

**Checkpoint**: User Story 3 independently lets a user install, understand, and exercise every current v1-bound workflow accurately.

---

## Phase 6: Polish and Cross-Cutting Verification

**Purpose**: Reconcile the three stories and prove release-candidate readiness without publishing.

- [x] T027 Record S018 safety, conformance, annotation, compatibility, and executable-documentation decisions under `[Unreleased]` in `CHANGELOG.md`
- [x] T028 Reconcile generated help, schema bytes, format row IDs, loss codes, documentation examples, compatibility claims, source fidelity, and privacy guarantees across all S018 files
- [x] T029 Run formatter, documentation, brand, Project policy, release-verifier, full tests, race, vet, Staticcheck, vulnerability, fixed-work fuzz, six-target static build, UTF-8, mojibake, and `git diff --check` checks from `specs/S018-freeze-v1-contracts/quickstart.md`
- [x] T030 Build and independently verify a clean non-publishing `0.1.0` six-target snapshot at the full S018 commit without admitting an immutable v1 schema or publishing protected state
- [x] T031 Run Spec Kit convergence against all 33 functional requirements, eight success criteria, three user stories, three contracts, and seven constitutional principles, appending and completing any remediation tasks
- [ ] T032 Publish the official pull request with complete closing references for #35, #36, and #41, verify its rendered body, set Project stages to PR review, and retain the unused default Status field blank
- [ ] T033 Resolve every external review finding, request at most one second Codex review only when round one has findings, and wait for all current-head CI, CodeQL, release-proof, security, review, and policy gates to become green

---

## Dependencies and Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately.
- **Foundational safety (Phase 2)**: Depends on setup and blocks all user stories.
- **User Story 1, 2, and 3 (Phases 3 through 5)**: Start after Phase 2 and proceed in parallel with isolated ownership.
- **Polish (Phase 6)**: Depends on all three user stories.

### User Story Dependencies

- **User Story 1**: No dependency on the schema or documentation implementation after shared safety work.
- **User Story 2**: No dependency on conformance implementation after shared safety work; it owns `docs/schema.md` to avoid conflict.
- **User Story 3**: May reference completed matrix and schema anchors during integration but owns distinct documentation and test files.

### Parallel Ownership

- **Conformance owner**: T007 through T016, excluding shared documentation files after handoff.
- **Schema owner**: T017 through T019.
- **Documentation owner**: T020 through T026.
- **Coordinator**: T001 through T006 and T027 through T033, plus integration and final verification.

## Parallel Example

```text
Agent A: T007-T016 in conformance, corpus, fuzz, external verifier, format and conversion files
Agent B: T017-T019 in canonical schema, annotation tests, and schema reference
Agent C: T020-T026 in README, CLI and compatibility docs, executable scenarios, and docs verifier
Coordinator: integrate shared boundaries, resolve overlaps, run complete verification, converge, publish, and review
```

## Implementation Strategy

### MVP first

Complete safety boundaries and User Story 1 first. This produces the release-critical hostile-input and whole-system evidence even if schema or prose work needs further iteration.

### Aggressive integrated delivery

Run all three story workstreams concurrently after the shared foundation, integrate once, and immediately run the full release-readiness gate. Keep release-candidate versioning and publication out of S018 so a correction cannot invalidate protected candidate evidence.

## Notes

- Every completed task is marked `[x]` only after its stated files and verification are complete.
- Tests precede or accompany their implementation.
- Unknown source content is preserved or the operation rejects; no limit silently truncates it.
- Repository-authored Markdown uses one source line per paragraph and list item.
- Push and official pull-request creation are authorized for S018. Merge, tag, release, schema publication, and production mutation are not authorized.
