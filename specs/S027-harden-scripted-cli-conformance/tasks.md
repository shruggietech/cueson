# Tasks: Harden scripted CLI conformance

**Input**: [spec.md](spec.md), [plan.md](plan.md), research, data model, contracts and quickstart.

**Tests**: Required by the source-fidelity constitution and this slice's acceptance. Focused tests precede or accompany changes.

## Phase 1: Setup

- [x] T001 Reassess current authority, merged S026, active issues/dependencies and CI in specs/S027-harden-scripted-cli-conformance/spec.md.
- [x] T002 Run installed specify/clarify/checklist/plan/tasks workflows and record choices in specs/S027-harden-scripted-cli-conformance/plan.md.

## Phase 2: Foundational gate

- [x] T003 Run installed blocking analyze over specs/S027-harden-scripted-cli-conformance/spec.md, plan.md and tasks.md; resolve all material findings before product edits.
- [x] T004 Read checklists without changing custom markers, record authorized continuation, and verify existing .gitignore plus ownership boundaries in specs/S027-harden-scripted-cli-conformance/verification.md.

## Phase 3: User Story 1 - Consistent discovery and inspection (P1)

**Independent test**: ASS/SSA native and Cue JSON selection/capability/summary agreement, privacy, old text/historical/report/help/completion compatibility and failure classes.

- [x] T005 [P] [US1] Add focused catalogue option-value/help/completion and stale-diagnostic regressions in internal/cli scripted surface tests.
- [x] T006 [US1] Extend validate/inspect values and descriptions and catalogue-derived invalid-format diagnostics in internal/cli/surface.go, validate.go and inspect.go.
- [x] T007 [US1] Add precise optional fifteen-field scripted count, non-dialogue/empty/schema-only, privacy and old-output regressions in internal/cli/scripted_inspect_test.go.
- [x] T008 [US1] Implement optional counts-only scripted metadata and conditional human rendering in internal/cli/inspect.go, preserving all text fields and shape.
- [x] T009 [US1] Refresh deliberate help/four completion snapshots in internal/cli/testdata and run full old-document/shared-input/output-error matrix in internal/cli.

## Phase 4: User Story 2 - Owning-boundary safety (P1)

**Independent test**: Direct capture ceilings/ownership rejection, accepted source immutability and exact all-format conversion cycles, accepted controls and adversarial structured documents.

- [x] T010 [P] [US2] Add independent encoded-source/declaration capture-boundary regressions in internal/model scripted capture tests.
- [x] T011 [US2] Preflight encoded-size allocation and use bounded declaration split in internal/model independent scripted source-capture inspection.
- [x] T012 [US2] Replace per-callback CLI/temp file conversion source creation with bounded in-memory four-format codec/model construction in internal/convert/scripted_fuzz_test.go and focused builder tests.
- [x] T013 [US2] Retain strict/fatal/determinism/reparse/immutable-source twelve-direction assertions and supported encoding/BOM/empty/source metadata builder cases in internal/convert.
- [x] T014 [US2] Add bounded structured scripted owner/envelope/privacy fuzz mutations with genuine accepted controls in internal/model scripted fuzz tests and rich declaration/override/attachment seeds plus native cycle invariants in internal/codec/scripted/fuzz_test.go.
- [x] T015 [US2] Run focused model/convert/capture/parser/render/review regressions and fixed-work new targets, recording evidence in specs/S027-harden-scripted-cli-conformance/verification.md.

## Phase 5: User Story 3 - Executed development conformance (P2)

**Independent test**: Every required development row traces to runnable owning evidence, all 34 accepted scripted/native-conversion fixtures restore exactly, and native workflow refusals are safe.

- [x] T016 [P] [US3] Add malformed/trailing/unknown-field/BOM/UTF-8 conformance decoder regressions in internal/testutil and require io.EOF after one matrix document.
- [x] T017 [US3] Add required-row and kind-correct actual test/fuzz reference checks in internal/conformance matrix tests, preserving valid existing portable inapplicability.
- [x] T018 [US3] Prove all 34 accepted scripted and scripted-conversion fixtures restore exact original bytes in internal/conformance.
- [x] T019 [US3] Add named ASS/SSA native workflow/capability/privacy/strict/fatal safety tests in internal/conformance; integrate after US1.
- [x] T020 [US3] Replace completed development hardening deferral and narrow remaining stable/candidate gate to #64/#65 in testdata/conformance-matrix.json without changing old fixture records.
- [x] T021 [US3] Schedule structured scripted, structured conversion and safe-basename fixed-work fuzz targets in .github/workflows/ci.yml using existing bounded foreground policy.

## Phase 6: Polish and delivery

- [x] T022 Update exact inspection semantics, authority/chronological status and dated decisions in docs/cli.md, docs/architecture.md, docs/roadmap.md, docs/Cueson-Project-Specification-v0.0.0.md and CHANGELOG.md.
- [x] T023 Run full CI-parity verification, scheduled fuzzing, six pure-Go builds, site/docs/artifact checks and convergence; record honest local/native-hosted boundaries in specs/S027-harden-scripted-cli-conformance/verification.md.
- [ ] T024 Commit conventionally, push authorized branch, publish formatted/read-back official PR closing #62/#63, reconcile Project PR review/unused Status, await current-head hosted CI/release proof and every review, remediate findings within two rounds, and hand off human final merge via exact PR evidence.

## Dependencies and execution order

T001-T004 block all product edits. CLI T005-T009, hardening T010-T015 and decoder/corpus T016-T018 have disjoint file ownership and may run together after the gate. T019 requires integrated US1, T020 requires genuine evidence, and T021-T024 follow integrated owners. #63 keeps native #62 dependency until human merge closes both complete outcomes.

## Parallel examples and implementation strategy

CLI agent owns internal/cli. Hardening agent owns internal/model and internal/convert guard/fuzz tests. Evidence agent owns internal/testutil, internal/conformance and matrix metadata. Root owns docs/workflow/artifacts/integration/delivery. Independent US1 is the useful discovery MVP; complete all three stories in this PR without deployment or release publication.

## Coverage

FR-001 → T005/T006/T009/T019; FR-002 → T007/T008/T019; FR-003 → T005/T006/T009/T022; FR-004 → T013/T018/T019/T023; FR-005 → T010/T011/T014/T019; FR-006 → T010/T011/T014/T015/T019; FR-007 → T012-T015/T021/T023; FR-008 → T016-T020/T023; FR-009 → T021-T024; FR-010 → T001-T004/T022-T024. SC-001-SC-005 map to the same independent story gates and final exact-head delivery.
