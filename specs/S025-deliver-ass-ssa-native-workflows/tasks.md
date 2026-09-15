# Tasks: S025 ASS/SSA native workflows

**Input**: [spec.md](spec.md), [plan.md](plan.md), [research.md](research.md), [data-model.md](data-model.md), [contract](contracts/native-workflows.md).

**Tests**: Required by the ratified constitution and FR-020; meaningful format/security tests precede or accompany implementation.

## Phase 1: Setup

- [x] T001 Confirm current authority, active issues, scope and hidden foreground verification in specs/S025-deliver-ass-ssa-native-workflows/plan.md.
- [x] T002 Execute installed specify/clarify/checklist/plan/tasks and blocking analysis; record coverage in specs/S025-deliver-ass-ssa-native-workflows/verification.md.

## Phase 2: Foundation

- [x] T003 Define and test narrow authoritative model field/declaration/timestamp/scalar/attachment helper interfaces in internal/model/scripted_native.go and internal/model/scripted_native_test.go (FR-005/FR-008/FR-016).
- [x] T004 Test and permit truthful experimental native capability observations while retaining schema_only/historical input in internal/model/scripted.go and internal/schema/cueson.schema.json (FR-007/FR-018).
- [x] T005 Test and implement explicit rejected scripted-candidate selection without extension fallback in internal/codec/detect.go and internal/codec/scripted_selection_test.go (FR-003).

## Phase 3: US1 fidelity corpus (P1)

Independent test: shared inventory/provenance/hash/privacy/encoding and required row resolution, without marking downstream assertions passed.

- [x] T006 [P] [US1] Author paired accepted and malformed scripted source fixtures with redistribution/byte contracts in testdata/fixtures/scripted/ass, testdata/fixtures/scripted/ssa and testdata/malformed/scripted (FR-001).
- [x] T007 [US1] Add unique payload/projection/diagnostic/canonical-byte records to testdata/manifest.json (FR-001/FR-002).
- [x] T008 [US1] Extend testdata/conformance-matrix.json and internal/conformance/matrix_test.go with shared scripted rows and explicit deferred downstream evidence (FR-002).
- [x] T009 [US1] Implement independent native/common/provenance/diagnostic and exact restoration assertions in internal/conformance/scripted_test.go (FR-002/FR-010/FR-014).

## Phase 4: US2 native ingest (P1)

Independent test: both variants match corpus projections, validate and restore exactly, with tested selection and safe malformed/unsafe boundaries.

- [x] T010 [P] [US2] Write focused dialect/encoding/declaration/order/unknown-field tests in internal/codec/scripted/parse_test.go (FR-003/FR-005).
- [x] T011 [US2] Implement bounded native framing, declaration ownership and truthful captures in internal/codec/scripted/detect.go and internal/codec/scripted/parse.go (FR-003/FR-004/FR-005/FR-016).
- [x] T012 [US2] Derive common text/lines, drawing exclusion, native speaker and eligible karaoke observations in internal/codec/scripted/projection.go (FR-006).
- [x] T013 [US2] Test privacy/active-content, malformed preservation-only content, safe diagnostics, and hostile ceilings in internal/codec/scripted/parse_test.go (FR-008/FR-009/FR-016).
- [x] T014 [US2] Integrate official native document construction and registry adapters in internal/cli/workflows.go and internal/cli/cli.go (FR-004/FR-007/FR-017).
- [x] T015 [US2] Prove source integrity/exact restore, zero-cue scripts, aliases/extension conflicts and established CLI error/publication rules in internal/cli/scripted_native_test.go (FR-010/FR-017/FR-018).

## Phase 5: US3 model rendering (P1)

Independent test: both variants deterministically render/reparse, coherent edits/construction appear, original restore remains exact and malformed/unsafe/precision/strict outcomes refuse before publication.

- [x] T016 [P] [US3] Write model-edit/construction/declaration-transition/canonical-order tests in internal/codec/scripted/render_test.go (FR-011/FR-012/FR-013).
- [x] T017 [US3] Implement shared bounded owner-based native serialization in internal/codec/scripted/render.go (FR-011/FR-012/FR-013/FR-016).
- [x] T018 [US3] Test malformed retained records/assets, framing contradictions, privacy, centisecond precision and strict refusal in internal/codec/scripted/render_test.go (FR-008/FR-015/FR-016).
- [x] T019 [US3] Prove native/common edit and parse/render/reparse equivalence with independent expected bytes in internal/conformance/scripted_test.go and internal/cli/scripted_native_test.go (FR-014/FR-015/FR-017).
- [x] T020 [US3] Add meaningful parser/render-cycle fuzz targets and integrate fixed-work execution in internal/codec/scripted/fuzz_test.go and .github/workflows/ci.yml (FR-016/FR-020).

## Phase 6: Integration and delivery

- [x] T021 Extend optional ASS/SSA corpus counts/native cycles in scripts/corpus-verify and retain existing format expectations (FR-018).
- [x] T022 Update truthful experimental capability, ownership, roadmap status, examples/annotations and dated decisions in docs/formats/ass-ssa.md, docs/architecture.md, docs/schema.md, docs/cli.md, docs/compatibility.md, docs/project-management.md, docs/roadmap.md, README.md and CHANGELOG.md (FR-019).
- [x] T023 Run foreground full root/nested tests/vet/static/vulnerability/fuzz/build/text/docs/brand/site gates and record results in specs/S025-deliver-ass-ssa-native-workflows/verification.md (FR-018/FR-020).
- [ ] T024 Commit, push and publish the formatted official closing PR; reconcile current-head CI, every review finding and at most one second Codex request, then record final delivery evidence in specs/S025-deliver-ass-ssa-native-workflows/verification.md and GitHub (FR-020).

## Dependencies and parallel work

T001/T002 block implementation. T003/T004/T005 establish model/signature/capability boundaries. After settled helper/parser signatures, exclusive parser tests/implementation, renderer tests/implementation and corpus authorship can proceed in parallel. T009/T015/T019 integration requires completed native operations; #57 acceptance precedes #58 and #59 within the one slice. Preserve native dependency traceability without reciprocal blockers.

Corpus agent example: T006/T007/T008 in testdata and conformance files. Parser agent example: T003/T010-T013 in model wrappers and scripted parser files. Renderer agent example: T016-T018 in scripted render files. Coordinator integrates T004/T005/T014/T015/T020-T024 and owns all final verification.

## Implementation strategy

First audit authored inventory/expectations (US1), then native encode/restore (US2), then model-render/edit/construction (US3). Focused tests establish each independently closeable outcome before full parity and publication. No downstream conversion or stable assertion is closed by this slice.
