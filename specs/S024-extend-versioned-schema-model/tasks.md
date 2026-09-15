# Tasks: Extend versioned schema and scripted model

**Input:** [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), and [input contract](contracts/input-and-capabilities.md).

**Tests:** Required behavior/security/limits tests precede or accompany implementation. No prose-mirroring tests.

## Phase 1: Setup and blocking foundation

- [x] T001 Establish approved branch/current authority and installed Spec Kit artifacts under specs/S024-extend-versioned-schema-model/. (FR-018–020)
- [x] T002 Complete research, five clarification decisions, requirements-quality assessment and blocking read-only analysis in plan.md/research.md before implementation. (FR-001–020)

## Phase 2: User Story 1 - Historical input

- [x] T003 [P] [US1] Write exact historical/current identity and local immutable authority tests in internal/schema/registry_test.go and schema_test.go. (FR-001–004,017)
- [x] T004 [US1] Embed byte-identical historical schema and implement selected local structure/identity dispatch in internal/schema/schema.go and historical/v1.0.0/cueson.schema.json. (FR-001–005,017)
- [x] T005 [P] [US1] Add historical semantic/capability/provenance regressions and isolate historical/current typed validation in internal/model/model.go and model tests. (FR-003–005,016)
- [x] T006 [US1] Write both-format five-command historical/third-party/invalid-identity/source-corruption and force/no-publication regressions in internal/cli/versioned_input_test.go. (FR-003–006)
- [x] T007 [US1] Integrate direct render/restore through shared Cue JSON finishing and report loaded schema identity in internal/cli/cli.go, workflows.go, input.go and inspect.go. (FR-005–007)

## Phase 3: User Story 2 - Scripted native model

- [x] T008 [P] [US2] Add positive/negative ordered ownership/declaration/capture/construction/attachment tests in internal/model/scripted_test.go and internal/schema/scripted_test.go. (FR-008–010)
- [x] T009 [US2] Implement typed native branches, local references/ordering/dialect/value semantics and bounded traversal in internal/model/model.go, scripted.go and limits helpers. (FR-008–010,012)
- [x] T010 [US2] Test and implement deterministic logical Text projection, drawing/override/speaker/token provenance and zero-dialogue summaries in internal/model/scripted*.go and tests. (FR-011,016)
- [x] T011 [US2] Test and enforce original/edited known/unknown metadata privacy, active-content refusal, and all ratified complexity boundaries in internal/model/scripted*.go and tests. (FR-006,009,012–013)
- [x] T012 [US2] Add annotated current schema shapes, valid portable examples and recursive coverage/native fidelity assertions in internal/schema/cueson.schema.json, annotations_test.go and testdata/. (FR-008–014,017)
- [x] T013 [US2] Add recognized nil-codec ASS/SSA lookup and truthful generic inspect/restore/unavailable-native behavior with tests in internal/codec/registry.go and internal/cli/. (FR-015–016)

## Phase 4: User Story 3 - Current discovery and development proof

- [x] T014 [US3] Stage exact1.1.0-dev software/current schema/representative output identity and discovery/encoding tests in internal/version/, internal/schema/ and internal/cli/. (FR-007,017)
- [x] T015 [US3] Reconcile only current candidate version/schema-source packaging assumptions and associated meaningful proof tests in .goreleaser.yaml, .github/workflows/release-proof.yml and internal/releasepolicy/ or scripts/release-verify/. (FR-007,017,019)
- [x] T016 [US3] Record S023 ratification/current development behavior, current/historical negotiation and dated architecture/process decisions in docs/roadmap.md, docs/schema.md, docs/compatibility.md, docs/architecture.md, docs/project-management.md, README.md and CHANGELOG.md. (FR-007,015,018)

## Phase 5: Integration and final verification

- [x] T017 Audit final issue acceptance/native state and spec/plan/tasks/contracts/code consistency; record change inventory and decisions in verification.md. (FR-001–020)
- [x] T018 Run complete appropriate foreground root/nested behavioral/security tests, fuzz budgets, vet/static/vulnerability, six builds, annotation/immutable integrity and formatter/docs/encoding/whitespace proof in verification.md. (FR-001–019)
- [x] T019 Run complete foreground site generation/lint/unit/build/browser/artifact/dry-run verification and document current/released separation in verification.md. (FR-014–019)
- [x] T020 Prepare conventional commit and formatted official PR body closing #55/#56 under the surrounding authorized delivery protocol; preserve human merge and maximum two review rounds in plan.md/verification.md. (FR-020)

## Dependencies and parallel execution

T001 precedes T002. All implementation begins after clean blocking analysis. Compatibility T003/T004 and native model tests T005/T008 can run with independent ownership. Root writes T006 in parallel with agents, then T007/T013/T014 integrate their actual interfaces. Native tests precede native implementation T009–T012. T015/T016 follow design identity and T017–T020 follow integration/passing gates.

Compatibility agent owns schema.go, schema_test.go and historical packaging/tests. Native agent owns model files/current JSON/annotations/representative/dedicated scripted tests. Root owns CLI/version/registry/process/docs and final proof. Same-file operations are sequential. Historical MVP is independently verifiable; full authorized slice includes all three stories.

Commit/push/official PR/CI/review execution surrounds these preparation tasks and is recorded as actual GitHub evidence. Final human merge and later release/production remain pending.
