# Tasks: Freeze the v1.1.0 release candidate

**Input**: [plan.md](plan.md), [spec.md](spec.md), research/data-model/contracts and quickstart.

**Tests**: Required by FR-007 to FR-010 and the constitution. Focused identity/capability, matrix and package refusal tests precede implementation.

## Phase 1: Setup and blocking design

- [x] T001 Verify clean S027 main, active issues, native dependencies and project/milestone authority.
- [x] T002 Execute installed specify, clarify, domain checklist and plan workflows with requirements/decision records in specs/S028-freeze-v1-1-release-candidate/.
- [x] T003 Consolidate independent research, explicit file ownership and test strategy in research.md and plan.md.
- [x] T004 Complete installed tasks/analyze blocking gate with full requirement/story/task coverage before implementation.

## Phase 2: User Story 1 - Frozen public contract (P1)

- [x] T005 [P] [US1] Add stable candidate documentation/schema-example/release-note requirements checks in scripts/docs-verify and scripts/release-verify/policy_test.go.
- [x] T006 [P] [US1] Reconcile README.md, docs/architecture.md, docs/compatibility.md, docs/schema.md, docs/cli.md, docs/conversion.md, docs/formats/ass-ssa.md and maintained format/provenance guidance with completed bounded native and conversion/CLI/conformance scope.
- [x] T007 [US1] Update docs/release-process.md and docs/release-verification.md for stable non-publishing candidate proof and exact published/production distinctions.
- [x] T008 [US1] Prepare concise docs/releases/v1.1.0.md with compatibility disclosures and required changelog suffix.
- [x] T009 [P] [US1] Add authored ASS/SSA route, governed format navigation and four-format conversion description in site/content-map.json, site/scripts/generate.mjs and Site generator/browser tests, preserving published download/schema inventory.
- [x] T010 [US1] Validate authored/schema examples, documentation markers, generated-content drift and preview route/support consistency.

## Phase 3: User Story 2 - Stable coherent current contract (P1)

- [x] T011 [P] [US2] Add meaningful stable current identity/capability/historical-byte/invalid-development assertions in internal/schema, internal/model and internal/cli tests; retain frozen historical inputs/goldens.
- [x] T012 [US2] Promote exact runtime schema/model/version dispatch and official CLI/convert native capability construction to 1.1.0 stable; preserve precise schema_only/experimental observations in current structural/semantic acceptance.
- [x] T013 [US2] Update current canonical schema annotations/examples, identity-specific current fixtures/help/completion/inspection goldens, and add identical schema/releases/v1.1.0/cueson.schema.json without touching historical resources/source fixtures.
- [x] T014 [P] [US2] Add stable four-format native workflow/restoration/refusal evidence in internal/conformance and forbid all scripted matrix deferrals.
- [x] T015 [US2] Replace both scripted-stable-gate deferred render/platform rows in testdata/conformance-matrix.json with genuine stable evidence.
- [x] T016 [US2] Run focused current/historical/native/convert/model/conformance tests and verify original fixture/manifest and historical schema/golden byte preservation.

## Phase 4: User Story 3 - Exact non-publishing candidate (P2)

- [x] T017 [P] [US3] Add packaged four-format native/historical/safety smoke and published-old-consumer integrity/refusal tests in scripts/release-verify.
- [x] T018 [US3] Extend matching host package execution in scripts/release-verify using hidden noninteractive platform process launch and deterministic revision/digest-bound published v1.0.0 consumer proof.
- [x] T019 [US3] Switch .goreleaser.yaml and .github/workflows/release-proof.yml to stable 1.1.0 immutable schema/candidate/native evidence while preserving pinned tools, read-only same-bundle proof and publication disable.
- [x] T020 [US3] Add exact S028 six-target release evidence contract under specs/S028-freeze-v1-1-release-candidate/contracts and require it through candidate policy tests.
- [x] T021 [US3] Run standalone release policy/smoke/refusal tests; prove concise prepared notes and protected publication boundaries.

## Phase 5: Integration, convergence and delivery

- [x] T022 Update CHANGELOG.md chronologically with candidate additions/decisions and reconcile docs/roadmap.md and main working plan with merged S027 and S028's stable candidate outcome.
- [x] T023 Execute installed convergence and independent cross-owner audit, resolving every material finding before verification.
- [x] T024 Run complete root/nested tests/vet/static/vulnerability, six pure-Go builds, sixteen fixed-work fuzz targets, workflow/format/encoding/docs/brand gates and full Site unit/browser/build/artifact/dry-run.
- [x] T025 Record accurate pre-publication evidence in verification.md and commit conventionally on the slice branch with clean VCS before exact local/hosted package proof.
- [ ] T026 Execute exact clean-revision six-target candidate structural/native package proof and published-old-consumer rejection; record source-bound evidence without dirtying the reviewed head.
- [ ] T027 Push and publish official non-draft formatted/read-back PR closing #64/#65; reconcile native Project stages/Slice with unused default Status.
- [ ] T028 Wait exact-head hosted CI/security/Site/three-host package checks and round-one Codex; respond to every finding, fix/verify/update head and request at most one second round if necessary.
- [ ] T029 Record terminal exact-head external completion by formatted/read-back PR comment and hand off to operator final review/merge with #66/#67 and epic/milestone pending.

## Dependencies and execution order

T004 blocks all product implementation. US1 tests/docs/Site and US2 identity/capability work can proceed with isolated file ownership after design; candidate acceptance follows the integrated frozen contract. T014/T015 coordinator ownership avoids matrix conflicts. T020 candidate owner writes the exact evidence contract. T022 coordinator owns changelog/roadmap/main plan. Candidate owner owns all other internal runtime/CLI/convert/schema tests and release verifier/config/workflow; documentation owner owns other root docs, docs verifier and Site.

Tests precede the corresponding promotion/proof changes. T023-T024 integrate all owners; exact VCS package proof follows commit because dirty candidate metadata is invalid. T028-T029 remain external delivery gates recorded on the official PR after arrival rather than creating an evidence-only reviewed-head change.

## Coverage

- US1: FR-001 to FR-005 and FR-011 -> T005-T010, T022, T024.
- US2: FR-006 to FR-009 -> T011-T016, T024.
- US3: FR-008 to FR-011 and FR-013 -> T017-T021, T025-T029.
- Cross-cutting: FR-012 -> T001-T004, T023-T025; FR-014 -> T001, T027-T029; FR-015 -> T007-T009, T019, T021, T027-T029.
