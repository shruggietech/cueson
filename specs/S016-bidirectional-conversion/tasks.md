# Tasks: Deliver Bidirectional Subtitle Conversion

## Phase 1: Foundational validation and loss contracts

- [x] T001 Add failing read-only source-envelope integrity tests in `internal/source/integrity_test.go`, then expose context-aware validation in `internal/source/integrity.go` for FR-004 and FR-030.
- [x] T002 [P] Add failing atomic loss, report validation, deterministic ordering, duplicate, bounded-context, path-leak, and typed strict-error tests in `internal/convert/report_test.go` for FR-010 through FR-014 and FR-030.
- [x] T003 Implement the stable runtime loss types, known-code registry, report validator and sorter, and typed conversion errors in `internal/convert/report.go` for FR-010 through FR-014 and FR-023.
- [x] T004 Move CLI-local SubRip representability aggregation into atomic reusable analysis with regression tests in `internal/convert/render.go`, `internal/convert/render_test.go`, and `internal/cli/workflows.go` for FR-010, FR-011, and FR-033.

## Phase 2: User Story 1 - Convert between SubRip and WebVTT

**Goal**: Produce deterministic parser-valid target bytes from native source or Cue JSON while preserving every representable common semantic.

**Independent Test**: Convert loss-free sources and documents in both directions, parse the targets, and compare cue order, timing, payload meaning, shared markup, and repeated output bytes.

- [x] T005 [P] [US1] Add provenance-recorded loss-free bidirectional sources and expected target bytes under `testdata/fixtures/conversion/`, then update `testdata/manifest.json` for FR-007 through FR-009, FR-021 through FR-022, SC-001, SC-005, and SC-006.
- [x] T006 [P] [US1] Add failing target-aware SubRip and WebVTT payload translation tests for shared markup, entities, literal angle syntax, cue-boundary text, empty lines, multiline text, and injection safety in `internal/convert/markup_test.go` for FR-019 and FR-021.
- [x] T007 [US1] Implement iterative target-aware payload scanning, proven-equivalent bold/italic/underline mapping, literal escaping, and target-safe empty-line handling in `internal/convert/markup.go` for FR-019 through FR-021.
- [x] T008 [US1] Add failing loss-free projection, input-immutability, target-native model, parser-cycle, overlap, duplicate, zero-duration fatal, and repeated-determinism tests in `internal/convert/convert_test.go` for FR-005 through FR-009, FR-021 through FR-023, SC-001, SC-005, and SC-006.
- [x] T009 [US1] Implement private SubRip-to-WebVTT and WebVTT-to-SubRip model projection, canonical rendering, summary recomputation, same-format and fatal-boundary errors, and immutable input handling in `internal/convert/convert.go` and `internal/convert/project.go` for FR-005 through FR-009 and FR-021 through FR-023.

## Phase 3: User Story 2 - Understand and prevent conversion loss

**Goal**: Account for every documented native omission or degradation and block all loss before rendering under strict policy.

**Independent Test**: Exercise every compatibility row, compare complete ordered loss goldens, and prove strict mode returns the entire report without target bytes.

- [x] T010 [P] [US2] Add provenance-recorded comprehensive lossy and fatal conversion sources under `testdata/fixtures/conversion/`, `testdata/malformed/conversion/`, and `testdata/fuzz/conversion/`, then update `testdata/manifest.json` for FR-016 through FR-020 and FR-031 through FR-032.
- [x] T011 [US2] Add failing bidirectional matrix tests for SubRip coordinates, font, heuristic speakers, tokens, OCR and unrecognized blocks plus WebVTT metadata, NOTE, STYLE, REGION, unknown blocks, identifiers, settings, placement, voice, class, language, ruby, entities, inline timing and tokens in `internal/convert/matrix_test.go` for FR-009 through FR-020, SC-002, and SC-003.
- [x] T012 [US2] Implement fixed-rank bidirectional compatibility analysis and atomic path-addressed loss emission in `internal/convert/matrix.go` for FR-009 through FR-020 and FR-030.
- [x] T013 [US2] Add failing full-report strict tests and implement pre-render strict rejection with complete typed report retention in `internal/convert/convert_test.go` and `internal/convert/convert.go` for FR-013 through FR-015 and SC-004.
- [x] T014 [US2] Add bounded deterministic conversion, loss-classification, markup, strict-equivalence, and renderer/parser-cycle fuzz targets in `internal/convert/fuzz_test.go` for FR-031, FR-032, and SC-007.

## Phase 4: User Story 3 - Automate conversion safely

**Goal**: Ship the complete `convert` command with unambiguous input classification, stable streams and statuses, and transactional file output.

**Independent Test**: Run native and Cue JSON inputs through every format, encoding, stream, filter, output, replacement, strict, failure, and help path and compare exact statuses, streams, and destination bytes.

- [x] T015 [P] [US3] Add failing command-registration, help, option, alias, conflict, same-format, format-capability, encoding, source-classification, stream, warning-order, quiet, silent, and exit-code tests in `internal/cli/cli_test.go` for FR-001 through FR-006 and FR-024 through FR-029.
- [x] T016 [US3] Add `convert` invocation parsing, help, `--from` aliases, target validation, option constraints, and command dispatch in `internal/cli/cli.go` for FR-001, FR-024, FR-027, and FR-028.
- [x] T017 [US3] Refactor bounded native source loading and installed rendering into shared helpers without changing encode or render behavior in `internal/cli/workflows.go` and `internal/cli/cli_test.go` for FR-003, FR-004, and FR-033.
- [x] T018 [US3] Implement Cue JSON precedence, source integrity validation, conversion execution, ordered loss projection, stdout purity, and transactional file publication in `internal/cli/workflows.go` for FR-002 through FR-006 and FR-024 through FR-030.
- [x] T019 [US3] Add strict new-file absence, forced-destination preservation, short-write, failed-write, input-path leakage, direct-source/Cue-JSON parity, and target re-encode tests in `internal/cli/cli_test.go` and `internal/conformance/conformance_test.go` for FR-015, FR-025 through FR-033, SC-001, SC-004, SC-006, and SC-007.

## Phase 5: Documentation and delivery

- [x] T020 Update the conversion contract, compatibility and loss matrices, architecture ownership, schema runtime-only loss boundary, runnable examples, and changelog in `docs/architecture.md`, `docs/cli.md`, `docs/schema.md`, `docs/formats/srt.md`, `docs/formats/webvtt.md`, `README.md`, and `CHANGELOG.md` for FR-034 through FR-036.
- [x] T021 Run formatting, UTF-8 and mojibake checks, fixture-manifest verification, focused and full tests, race tests, fuzz seeds, pure-Go cross-builds, quickstart scenarios, and development release proof for FR-031 through FR-035 and SC-001 through SC-007.
- [x] T022 Run Spec Kit convergence and the blocking analysis gate, append and complete any traceable remediation tasks, and prove every requirement and success criterion has implementation evidence in `specs/S016-bidirectional-conversion/tasks.md`.
- [x] T023 Reconcile issue #33 and the `cueson Delivery` Project, publish a formatter-verified pull request closing #33, read the body back, and move the issue to `PR review` under the user's push and PR authorization.
- [ ] T024 Watch every current-head CI and external review result, address every finding, request at most one second Codex review with `@codex review` only if round one has findings, and stop only after the pull request is green and fully reviewed for FR-037 and SC-008.

## Dependencies and execution order

- Phase 1 blocks all user stories because conversion needs one validated loss vocabulary, source-integrity entry point, and shared representability authority.
- User Story 1 establishes the loss-free projection and rendering path required by later loss policy and CLI integration.
- User Story 2 depends on User Story 1's projection and payload translation and completes the compatibility matrix plus strict behavior.
- User Story 3 depends on both library stories and then integrates them through the public CLI and publisher.
- Documentation and delivery depend on all user stories.

## Parallel opportunities

- T001 and T002 use independent source and conversion files before T003.
- T005 and T006 can proceed in parallel after the foundational report contract.
- T010 can prepare governed corpus data while T011 defines failing matrix tests.
- T015 can define the CLI contract while library implementation is confined to `internal/convert`.
- Documentation can begin after runtime behavior stabilizes while final conformance and transaction evidence runs under separate file ownership.

## Implementation strategy

Deliver the runtime loss contract and source-integrity boundary first. Build and verify loss-free conversion before layering every lossy compatibility row and strict policy, then expose the complete operation through the CLI. Keep loss reports out of Cue JSON, do not invent timing or semantic metadata, and finish with full foreground verification before publication.

## Phase 6: Convergence

- [x] T025 Add complete canonical loss-report and CLI warning goldens for both lossy directions and compare them under native cross-platform tests per SC-002 and SC-003 (partial).
- [x] T026 Add bounded adversarial structured-document fuzz coverage for conversion validation, loss classification, path hygiene, and strict equivalence per FR-032 and SC-007 (partial).
