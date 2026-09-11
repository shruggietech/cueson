# Tasks: Deliver Usable SubRip Workflows

## Phase 1: Shared foundation

- [x] T001 Add failing source-size and exact-byte tests in `internal/source/source_test.go`, then implement bounded exact capture in `internal/source/source.go`.
- [x] T002 [P] Add failing registry and detection tests in `internal/codec/registry_test.go`, then implement `internal/codec/registry.go` and `internal/codec/detect.go`.
- [x] T003 [P] Add failing Unicode and legacy-override tests in `internal/codec/text_test.go`, then implement `internal/codec/text.go`.
- [x] T004 Add failing capability-profile tests in `internal/model/model_test.go`, then update `internal/model/model.go`.
- [x] T005 Add failing `0.1.0` schema tests in `internal/schema/schema_test.go`, then evolve `internal/schema/cueson.schema.json`, `internal/schema/schema.go`, and `internal/version/version.go` without editing the immutable v0.0.0 schema.

## Phase 2: User Story 1 - Encode SubRip

- [x] T006 [US1] Add provenance-recorded grammar, encoding, line-ending, tag, coordinate, and malformed fixtures under `testdata/subrip/` and update `testdata/manifest.json`.
- [x] T007 [US1] Add failing physical-line, timecode, tag, speaker, parser, malformed, and fuzz-seed tests in `internal/codec/subrip/*_test.go`.
- [x] T008 [US1] Implement SubRip scanning, parsing, diagnostics, and model construction in `internal/codec/subrip/`.
- [x] T009 [US1] Add failing encode invocation, stdout, output-safety, path-leakage, schema, and restore tests in `internal/cli/cli_test.go`.
- [x] T010 [US1] Implement `cueson encode` and safe generalized publication in `internal/cli/`.

## Phase 3: User Story 2 - Render SubRip

- [x] T011 [US2] Add failing canonical render, coordinates, multiline, unsupported renderer, stdout, force, and round-trip tests in `internal/codec/subrip/render_test.go` and `internal/cli/cli_test.go`.
- [x] T012 [US2] Implement canonical rendering in `internal/codec/subrip/render.go` and `cueson render` in `internal/cli/`.

## Phase 4: User Story 3 - Deterministic diagnosis

- [x] T013 [US3] Add and pass diagnostic-order, extension-disagreement, ambiguity, missing-capability, malformed-input, and no-output regression tests across `internal/codec/` and `internal/cli/`.
- [x] T014 [US3] Complete stable diagnostic codes and help/error wording in `internal/codec/`, `internal/codec/subrip/`, and `internal/cli/`.

## Phase 5: Integration and delivery

- [x] T015 Update `docs/architecture.md`, `docs/cli.md`, `docs/schema.md`, `docs/formats/srt.md`, `README.md`, and `CHANGELOG.md`.
- [x] T016 Adapt `.goreleaser.yaml`, `.github/workflows/release-proof.yml`, `scripts/release-verify/`, and tests so development proof stays green while immutable v0.0.0 assertions remain intact.
- [x] T017 Update conformance expectations, then run formatting, mojibake, docs, focused, full, and quickstart verification.
- [x] T018 Run Spec Kit convergence and analysis, complete discovered tasks, and prove the immutable v0.0.0 schema digest is unchanged.
- [ ] T019 Reconcile GitHub metadata, publish a formatted PR closing #30 and #31, verify read-back, and move both issues to `PR review`.
- [ ] T020 Watch every current-head check and review, resolve all findings, request at most one second Codex round, and stop only after a green fully reviewed head.
