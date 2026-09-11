# Tasks: Complete CLI Workflows

## Phase 1: Shared validation and stream foundation

- [x] T001 Add failing Cue JSON precedence, explicit selector, native evidence, encoding, source-integrity, cancellation, and deterministic input-classification tests in `internal/cli/input_test.go` for FR-002 through FR-007, FR-010, and SC-001 through SC-002.
- [x] T002 Implement one read-only validated-input result and loader in `internal/cli/input.go`, then replace conversion-local classification in `internal/cli/workflows.go` without changing conversion behavior for FR-002 through FR-007, FR-010, and FR-027 through FR-028.
- [x] T003 Add failing diagnostic-stream write tests and implement retained diagnostic write errors plus status finalization in `internal/cli/cli_test.go` and `internal/cli/cli.go` for FR-009, FR-016, FR-023, and FR-026 through FR-029.

## Phase 2: User Story 1 - Validate without producing output

**Goal**: Assert complete Cue JSON or native input validity with ordered diagnostics and no payload or file output.

**Independent Test**: Validate accepted, warning-bearing, malformed, corrupt, ambiguous, unsupported, oversized, non-regular, and canceled inputs and compare exact statuses and empty stdout.

- [x] T004 [US1] Add failing `validate` parser, alias, conflict, help, stream, filter, no-output, validation-stage, and status tests in `internal/cli/validate_test.go` for FR-001 through FR-009, FR-026 through FR-029, SC-001, and SC-002.
- [x] T005 [US1] Implement `validate` options, parsing, help target, dispatch, assertion-only workflow, warnings, success diagnostic, and typed errors in `internal/cli/validate.go` and `internal/cli/cli.go` for FR-001 through FR-009 and FR-026 through FR-028.
- [x] T006 [US1] Add cross-package validation fixtures and no-output, path-leak, and generated-envelope integrity assertions in `internal/conformance/conformance_test.go` for FR-006 through FR-009, FR-019, FR-029, SC-001, and SC-002.

## Phase 3: User Story 2 - Inspect safely in human or JSON form

**Goal**: Produce deterministic privacy-bounded structural inspection reports from the same validated-input authority.

**Independent Test**: Compare complete human and JSON goldens for each input class, validate the fixed snake_case shape, and prove prohibited source or machine values never appear.

- [x] T007 [P] [US2] Add failing report construction, aggregate, cue, block, diagnostic-location, ordering, overflow, repeated-determinism, JSON-shape, and privacy-sentinel tests in `internal/cli/inspect_test.go` for FR-010 through FR-019 and SC-003 through SC-004.
- [x] T008 [US2] Implement the fixed version-1 inspection types, safe report projection, installed capability lookup, integrity summary, and human and compact JSON renderers in `internal/cli/inspect.go` for FR-010 through FR-019.
- [x] T009 [US2] Add failing `inspect` parser, aliases, help, JSON purity, default report, warning, filter, classification-failure, stream-failure, and status tests in `internal/cli/inspect_test.go` for FR-001, FR-004 through FR-005, FR-010, FR-016 through FR-019, and FR-026 through FR-029.
- [x] T010 [US2] Implement `inspect` options, parsing, help target, dispatch, shared validation, and stdout behavior in `internal/cli/inspect.go` and `internal/cli/cli.go` for FR-001, FR-004 through FR-005, FR-010, FR-016 through FR-019, and FR-026 through FR-028.

## Phase 4: User Story 3 - Configure shells and discover the complete CLI

**Goal**: Generate safe static completion and expose an exact complete help contract for all shipped commands.

**Independent Test**: Compare all completion and help bytes with goldens, prove catalogue/parser agreement, reject unsupported invocations, and execute an actual-binary stream and status matrix.

- [x] T011 [P] [US3] Add failing ordered catalogue invariants and exact root plus nine-command help goldens in `internal/cli/surface_test.go`, `internal/cli/help_test.go`, and `internal/cli/testdata/help/` for FR-001, FR-022, and FR-024 through FR-027, SC-005, and SC-006.
- [x] T012 [US3] Implement the ordered command/option catalogue and deterministic complete help and short-usage rendering in `internal/cli/surface.go`, `internal/cli/help.go`, and `internal/cli/cli.go` for FR-001, FR-022, and FR-024 through FR-027.
- [x] T013 [P] [US3] Add failing shell parser, exact four-script golden, repeatability, LF/UTF-8, vocabulary completeness, and forbidden-behavior tests in `internal/cli/completion_test.go` and `internal/cli/testdata/completion/` for FR-020 through FR-023 and SC-005.
- [x] T014 [US3] Implement strict completion invocation and deterministic static Bash, Zsh, Fish, and PowerShell generators in `internal/cli/completion.go` and `internal/cli/cli.go` for FR-020 through FR-023.
- [x] T015 [US3] Add optional available-interpreter syntax smoke tests plus one-build actual executable help, completion, validation, and inspection process tests with hidden Windows child-process configuration in `internal/cli/completion_syntax_test.go`, `internal/cli/process_test.go`, `internal/cli/process_windows_test.go`, and `internal/cli/process_other_test.go` for FR-021 through FR-030, SC-005 through SC-007.

## Phase 5: Documentation and delivery

- [x] T016 Update the complete implemented command, validation, inspection-report, completion, stream, exit, privacy, and architecture contracts plus runnable development examples in `docs/cli.md`, `docs/architecture.md`, `docs/schema.md`, `docs/formats/srt.md`, `docs/formats/webvtt.md`, `README.md`, and `CHANGELOG.md` for FR-031 through FR-032.
- [x] T017 Run formatting, UTF-8 and mojibake checks, exact goldens, focused and full tests, race tests, bounded fuzz seeds, static analysis, vulnerability scanning, pure-Go cross-builds, documentation verification, quickstart scenarios, and development release proof for FR-029 through FR-032 and SC-001 through SC-007.
- [x] T018 Run Spec Kit convergence and the blocking analysis gate, append and complete any traceable remediation tasks, and prove every requirement and success criterion has implementation evidence in `specs/S017-complete-cli-workflows/tasks.md`.
- [ ] T019 Reconcile issue #34 and the `cueson Delivery` Project, publish a formatter-verified pull request closing #34, read the body back, and move the issue to `PR review` under the user's push and PR authorization.
- [ ] T020 Watch every current-head CI and external review result, address every finding, request at most one second Codex review with a commit-bound `@codex review` only if round one has findings, and stop only after the pull request is green and fully reviewed for FR-033 and SC-008.

## Dependencies and execution order

- Phase 1 blocks validation and inspection because both require one classification and validation authority and reliable diagnostic finalization.
- User Story 1 proves the assertion-only pipeline before User Story 2 projects it into public report shapes.
- User Story 2 can build report types after the Phase 1 loader stabilizes while User Story 3 develops independent surface and completion files.
- Catalogue integration must use the final validation and inspection option contracts before its exact goldens are accepted.
- Documentation and delivery depend on all user stories and the final catalogue.

## Parallel opportunities

- T003 can proceed independently from T001 and T002 until final CLI integration.
- T007 and T011/T013 use separate inspection and surface/completion files after the shared foundation.
- Completion generator work can proceed alongside inspection implementation once final command options are fixed.
- Documentation can begin after help and JSON output stabilize while process and interpreter smoke evidence runs under separate file ownership.

## Implementation strategy

Build the shared validated-input and stream foundation first, then ship validation as the first independently usable workflow. Add inspection as a privacy-bounded projection of that exact pipeline. Complete the public surface with a small ordered catalogue and four static completion scripts, keeping semantic parsing manual and tested rather than replacing the CLI framework. Finish with full foreground verification, convergence, publication, two-round review handling, and a human-only merge handoff.
