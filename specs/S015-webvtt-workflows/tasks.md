# Tasks: Deliver Native WebVTT Workflows

## Phase 1: Foundational model and output safety

- [x] T001 Add failing WebVTT capability, source-order, raw-line, setting-occurrence, region, payload, speaker, token, and SubRip-regression tests in `internal/model/model_test.go` and `internal/schema/schema_test.go` for FR-005, FR-006, FR-010 through FR-020.
- [x] T002 Evolve `internal/model/model.go` and `internal/schema/cueson.schema.json` to the reviewed WebVTT data model while preserving version `0.1.0`, SubRip behavior, and immutable `schema/releases/v0.0.0/cueson.schema.json`.
- [x] T003 Add failing short-write, failed-write, new-destination, and forced-replacement transaction tests in `internal/cli/cli_test.go`, then repair the shared schema, encode, and render publisher in `internal/cli/cli.go` for FR-029.

## Phase 2: User Story 1 - Encode WebVTT as Cue JSON

- [x] T004 [US1] Add provenance-recorded accepted WebVTT sources for signature, BOM, CRLF, LF, lone CR, headers, all native block kinds, identifiers, timestamps, settings, markup, entities, inline timestamps, whitespace-only payloads, adjacent blocks, and overlapping cues under `testdata/fixtures/webvtt/`, then update `testdata/manifest.json` for FR-030 and FR-031.
- [x] T005 [P] [US1] Add failing signature, UTF-8, NUL, physical-line, timestamp, and block-state tests in `internal/codec/webvtt/parse_test.go` and `internal/codec/webvtt/time_test.go` for FR-002 through FR-009.
- [x] T006 [P] [US1] Add failing cue-setting and REGION-setting occurrence tests in `internal/codec/webvtt/settings_test.go` for FR-006, FR-010, and FR-011.
- [x] T007 [P] [US1] Add failing markup, entity, voice, language, ruby, and inline-token tests in `internal/codec/webvtt/markup_test.go` for FR-012 through FR-016.
- [x] T008 [US1] Implement bounded UTF-8 physical-line scanning, signature/header handling, strict timestamps, setting parsing, and iterative body classification in `internal/codec/webvtt/lines.go`, `internal/codec/webvtt/time.go`, `internal/codec/webvtt/settings.go`, and `internal/codec/webvtt/parse.go`.
- [x] T009 [US1] Implement iterative cue-text scanning and derived plain text, speaker observations, and token timing in `internal/codec/webvtt/markup.go`.
- [x] T010 [US1] Add failing automatic and explicit WebVTT encode, extension-disagreement, incompatible-encoding, JSON-only stdout, exact-restore, schema-validation, path-leak, output, force, and no-partial-output tests in `internal/cli/cli_test.go` and `internal/conformance/conformance_test.go` for FR-017 through FR-022 and FR-029.
- [x] T011 [US1] Register native WebVTT detection and decode through `internal/codec/webvtt/codec.go` and `internal/cli/workflows.go`, remove the SubRip-only capture media type assumption, and update `internal/cli/cli.go` help and option validation.

## Phase 3: User Story 2 - Render Cue JSON as WebVTT

- [x] T012 [US2] Add deterministic header, source-order merge, timestamp, setting, payload, non-cue block, strict/permissive, malformed-model, and parser-cycle renderer tests in `internal/codec/webvtt/render_test.go` for FR-023 through FR-027.
- [x] T013 [US2] Implement canonical model-driven WebVTT serialization and lexical consistency checks in `internal/codec/webvtt/render.go`.
- [x] T014 [US2] Add failing stdout, output, force, strict, rollback, and encode-render-reencode CLI tests in `internal/cli/cli_test.go`, then register the WebVTT renderer in `internal/cli/workflows.go` for FR-023 through FR-029.

## Phase 4: User Story 3 - Preserve and explain native fidelity

- [x] T015 [US3] Add provenance-recorded tolerated and fatal WebVTT sources for missing separators, decreasing starts, duplicate identifiers, invalid and duplicate settings, misplaced blocks, unknown blocks, malformed markup and entities, NUL, invalid Unicode, invalid signature, and invalid timing under `testdata/malformed/webvtt/` and `testdata/fuzz/webvtt/`, then update `testdata/manifest.json`.
- [x] T016 [US3] Add and pass deterministic diagnostic-order, raw-preservation, contiguous-order, rolling-caption, no-deduplication, and strict-render regression tests across `internal/codec/webvtt/`, `internal/model/model_test.go`, and `internal/conformance/conformance_test.go` for FR-005 through FR-016 and FR-026.
- [x] T017 [US3] Add fuzz targets for signature/header, block parsing, timestamps, cue and region settings, markup/entities, source integration, and renderer/parser cycles in `internal/codec/webvtt/fuzz_test.go` for FR-032.

## Phase 5: Documentation and delivery

- [x] T018 Update implemented behavior and runnable examples in `docs/architecture.md`, `docs/cli.md`, `docs/schema.md`, `docs/formats/webvtt.md`, `README.md`, and `CHANGELOG.md` for FR-034 while preserving S015 exclusions in FR-035.
- [x] T019 Update development release-proof and conformance expectations where required in `.github/workflows/release-proof.yml`, `.goreleaser.yaml`, `scripts/release-verify/`, and `internal/conformance/` without modifying immutable v0.0.0 material.
- [x] T020 Run formatting, UTF-8 and mojibake checks, fixture-manifest verification, focused tests, full tests, race tests, bounded fuzz seeds, pure-Go cross-builds, quickstart scenarios, and development release proof for FR-030 through FR-033.
- [x] T021 Run Spec Kit convergence and the blocking analysis gate, append and complete any traceable remediation tasks, and verify every specification requirement and success criterion has implementation evidence in `specs/S015-webvtt-workflows/tasks.md`.
- [x] T022 Reconcile issue #32 and the `cueson Delivery` Project, publish a formatter-verified pull request closing #32, read the body back, and move the issue to `PR review` under the user's push and PR authorization.
- [ ] T023 Watch every current-head CI and external review result, address every finding, request at most one second Codex review with `@codex review` only if round one has findings, and stop only after the pull request is green and fully reviewed for FR-036 and SC-008.

## Dependencies and execution order

- Phase 1 blocks all user stories because the parser and renderer require the reviewed WebVTT model and safe publisher.
- User Story 1 is the MVP and blocks renderer-cycle integration in User Story 2.
- User Story 2 depends on the parser and setting semantics from User Story 1 but remains independently testable from a valid Cue JSON document.
- User Story 3 extends the accepted parser and renderer with recovery, fidelity, and hostile-input evidence and depends on both earlier stories.
- Documentation and delivery depend on all user stories.

## Parallel opportunities

- T005, T006, and T007 are independent test-first seams after T001 and T002.
- Model/schema ownership, the new WebVTT codec tree, and CLI transaction repair use separate files during Phase 1.
- Documentation and release-proof inspection can proceed in parallel after runtime behavior stabilizes, with one owner per file.

## Implementation strategy

Complete the model/schema foundation and shared output repair first. Deliver native encode as the usable MVP, add renderer semantics and parser cycles, then close tolerance/fuzz coverage before updating public documentation and publishing. Do not start cross-format conversion or stable-support work inside S015.

## Convergence remediation

- [x] T024 Add governed mixed-line-ending and combined tolerated-input WebVTT sources, require every accepted WebVTT source to encode and restore byte-exactly, and verify each promised recovery class has corpus evidence for SC-002, SC-003, and SC-006.
- [x] T025 Align the markup scanner with the normative sole-component voice-span end-tag omission, add regression coverage, and remove false nonconformance diagnostics from conforming voice cues for FR-014, FR-015, and SC-006.
- [x] T026 Resolve every first-round Codex finding with regression coverage for raw setting whitespace, raw NUL fidelity, REGION regeneration, native strict-mode recomputation, hostile setting occurrences, nonempty payload lines, and malformed empty voice annotations; align character-reference decoding with the current W3C WebVTT algorithm for FR-002, FR-006, FR-010 through FR-015, FR-023 through FR-026, and SC-005 through SC-008.
