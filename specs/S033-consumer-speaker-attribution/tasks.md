# Tasks: S033 Consumer speaker attribution and media timing

**Input**: [spec.md](spec.md), [plan.md](plan.md), [data-model.md](data-model.md), [contracts/consumer-annotations.md](contracts/consumer-annotations.md).

**Tests**: Required test-first security/Unicode/bounds/source/native coverage by constitution and issue #82.

## Phase 1: Setup

- [x] T001 Complete installed specify/clarify/checklist/plan artifacts in specs/S033-consumer-speaker-attribution/ (FR-020).
- [x] T002 Complete blocking read-only analysis and owner-authorized checklist assessment against specs/S033-consumer-speaker-attribution/tasks.md (FR-020).

## Phase 2: Foundational compatibility

- [x] T003 Add current/historical identity regression tests in internal/schema/consumer_test.go and internal/schema/stable_candidate_test.go (FR-016,017; SC-004).
- [x] T004 Stage1.2.0-dev in internal/schema/schema.go and internal/version/version.go; locally retain byte-exact historical/v1.1.0/cueson.schema.json (FR-016,017).
- [x] T005 Transition non-publishing development snapshot .goreleaser.yaml, .github/workflows/release-proof.yml and scripts/release-verify policy assertions without changing historical release authorities (FR-016,017,019).

## Phase 3: US1 Retain consumer results

**Independent test**: Unicode consumer values round-trip beside unchanged native labels; original source restores exactly.

- [x] T006 [P] [US1] Add current schema optional-field/null/type/Unicode/bounds tests in internal/schema/consumer_test.go (FR-001..005; SC-001,002).
- [x] T007 [P] [US1] Add exact-value/ordering/security/historical typed-validation tests in internal/model/consumer_test.go (FR-001..005,014,016).
- [x] T008 [US1] Implement current schema optional attribution/media definitions and annotations/examples in internal/schema/cueson.schema.json and focused fixtures (FR-001..005,008,018).
- [x] T009 [US1] Implement SpeakerAttribution/MediaTiming fields, frozen ID validation and collection checks in internal/model/model.go, consumer.go and limits.go (FR-001..005,014).
- [x] T010 [US1] Preserve scripted/history typed revalidation in internal/model/scripted_target.go and model.go (FR-016,017).

## Phase 4: US2 Validate cue/media timing

**Independent test**: Valid offset intervals pass, invalid assignments reject, source-cue conflicts warn and remain restorable.

- [x] T011 [US2] Add interval/media offset/overflow/edit/zero-duration tests in internal/model/consumer_test.go (FR-006..012; SC-002,005).
- [x] T012 [US2] Implement checked timing validation and immutable counts/diagnostics helpers in internal/model/consumer.go (FR-006..013).
- [x] T013 [P] [US2] Add command privacy/state/strict/no-publication tests in internal/cli/consumer_test.go and restoration tests in internal/source/consumer_test.go (FR-007,011..014; SC-001,002,005).
- [x] T014 [US2] Integrate optional inspect projection and validate/restore/runtime warning paths in internal/cli/input.go, inspect.go, validate.go and cli.go (FR-011..014).
- [x] T015 [US2] Add meaningful Unicode/timing/roundtrip fuzz coverage in internal/model/consumer_test.go and internal/schema/consumer_test.go (FR-019; SC-002).

## Phase 5: US3 Export with complete loss reporting

**Independent test**: Four matching render paths and twelve conversion edges report omissions; strict publishes nothing.

- [x] T016 [P] [US3] Add complete conversion and strict-before-render tests in internal/convert/consumer_test.go (FR-015; SC-003).
- [x] T017 [P] [US3] Add native renderer warning/strict tests in internal/codec/webvtt and internal/codec/scripted consumer tests (FR-015; SC-003).
- [x] T018 [US3] Add conversion-wide field-specific consumer loss pass in internal/convert/matrix.go/report.go and clear private variant fields in scripted_variant.go (FR-015,017).
- [x] T019 [US3] Integrate bounded neutral omission warnings in internal/convert/render.go and codec native renderers (FR-012,015).
- [x] T020 [US3] Prove historical native/restoration/current inspection regression in internal/cli and internal/conformance current identity expectations (FR-016,017,019; SC-004).

## Phase 6: Polish and delivery

- [x] T021 Update README.md, docs/schema.md, cli.md, conversion.md, compatibility.md, architecture.md and CHANGELOG.md with consumer meaning, duration/coverage distinction and process decisions (FR-018; SC-005).
- [x] T022 Update scripts/docs-verify current-source markers and Site-preview development prose while preserving immutable published routes/downloads (FR-016,018).
- [x] T023 Run focused/full foreground CI parity, fuzz, security, native/platform/build and formatting/docs checks; record actual results in specs/S033-consumer-speaker-attribution/verification.md (FR-019; SC-001..006).
- [x] T024 Commit, authorized push and official PR closing #82; verify formatted publication and Project stages (FR-020).
- [ ] T025 Handle all bot findings through at most two review rounds and green final-head CI; prepare human final merge handoff in verification.md (FR-020; SC-006).

## Dependencies and execution

T001/T002 block all implementation. T003/T004 version dispatch and T009 model helper contracts underpin integration. Test tasks precede or accompany their owning implementation. Model/schema/export ownership is disjoint; root integrates CLI/source/docs and final verification. T023 precedes T024; T025 is post-publication. No final merge/tag/release/production step is authorized.

## Parallel example

Model agent owns T007/T009..012/T015, schema agent owns T003..006/T008 and current identity snapshots, export agent owns T016..019. Root owns T013/T014/T020..025. Helper signatures are coordinated before dependent edits.

## Implementation strategy

Independently prove consumer storage, then timeline checks, then native omission accounting. Integrate the full issue outcome before publication; a partial story is not issue completion. Checklist quality assessment does not represent completed code tasks.
