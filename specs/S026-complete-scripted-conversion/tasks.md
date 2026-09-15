# Tasks: Complete scripted conversion

**Input**: Design documents from specs/S026-complete-scripted-conversion.
**Prerequisites**: spec.md, plan.md, research.md, data-model.md, contracts/conversion.md, quickstart.md.
**Tests**: Required by the feature specification and constitution.
**Organization**: User story phases preserve independently testable outcomes.

## Phase 1: Setup

- [x] T001 Inspect current constitution, architecture, native GitHub #60/#61 dependencies, active issues, installed skills, and explicit delivery authority in specs/S026-complete-scripted-conversion/spec.md.
- [x] T002 Execute installed specify, clarify, checklist, plan research, and tasks workflows and generate specs/S026-complete-scripted-conversion design artifacts.

## Phase 2: Foundational gates

- [x] T003 Execute blocking installed analyze and record coverage/gate evidence in specs/S026-complete-scripted-conversion/verification.md.
- [x] T004 Extend four-format routing, nil context/full integrity checks, optional private target analysis, and source-aware renderer in internal/convert/convert.go, matrix.go, project.go, and render.go.
- [x] T005 Append fourteen stable loss codes for scripted observations and custom source identifiers, and native source-order validation with focused report tests in internal/convert/report.go, report_test.go, and scripted_report_test.go.
- [x] T006 [P] Add original-plus-target private validation/rendering boundaries without weakening public source validation in internal/model/scripted.go, scripted_target.go, scripted_target_test.go, internal/codec/scripted/render.go, render_target.go, and render_target_test.go.

## Phase 3: User Story 1 (P1), scripted outbound

**Goal**: Close #60 through four usable outbound directions with complete losses.
**Independent Test**: Both dialects to both text targets have expected text/timing/emphasis, complete native loss entries, and strict/fatal refusal.

- [x] T007 [P] [US1] Add native source/payload loss and shared emphasis/reset/drawing regression tests in internal/convert/scripted_analysis_test.go and scripted_review_test.go.
- [x] T008 [US1] Implement bounded native owner loss traversal and explicit fatal drawing/empty/malformed edges in internal/convert/scripted_analysis.go.
- [x] T009 [US1] Implement ordered balanced emphasis, native breaks, safe literal target text, and atomic override/span accounting in internal/convert/scripted_text.go.
- [x] T010 [P] [US1] Author outbound scripted baselines/native-loss expectations and semantic/loss conformance tests in testdata/fixtures/scripted-conversion and internal/conformance/scripted_conversion_test.go.
- [x] T011 [US1] Integrate outbound target semantic validation and verify established text-target rendering in internal/convert/render.go and scripted_analysis_test.go.

## Phase 4: User Story 2 (P1), scripted targets and variants

**Goal**: Close #61 through four text-source scripted targets and both variant directions.
**Independent Test**: Six directions have deterministic output, documented defaults, exact/lossy timing, faithful variant mapping, target reparse, and safe refusal.

- [x] T012 [P] [US2] Add deterministic default, nested emphasis, literal controls, timing ties/neighbors/collapse/overflow, and historical target tests in internal/convert/scripted_project_test.go.
- [x] T013 [US2] Implement context-aware deterministic target owners and checked centisecond endpoint losses in internal/convert/scripted_project.go.
- [x] T014 [P] [US2] Add all alignment, style/color/alpha/Layer/Marked, safe capture/comment/attachment, and malformed variant regressions in internal/convert/scripted_project_test.go.
- [x] T015 [US2] Implement deep-copied native dialect mapping and complete per-field/override variant losses in internal/convert/scripted_variant.go.
- [x] T016 [US2] Integrate source-specific text loss accounting and native target semantic reparse without false blank-line/entity losses in internal/convert/matrix.go and render.go.
- [x] T017 [P] [US2] Author six scripted-target goldens, precision/default/variant expectations, and conformance coverage in testdata/fixtures/scripted-conversion and internal/conformance/scripted_conversion_test.go.

## Phase 5: User Story 3 (P1), complete matrix and publication safety

**Goal**: Prove all twelve directions, old compatibility, privacy, and atomic CLI behavior.
**Independent Test**: Matrix/domain/CLI checks pass and every strict/fatal case publishes no payload or destination change.

- [x] T018 [P] [US3] Implement twelve-direction exact golden/report/semantic/determinism/immutability coverage and append manifest/matrix evidence in internal/conformance/scripted_conversion_test.go, testdata/manifest.json, and testdata/conformance-matrix.json.
- [x] T019 [US3] Integrate conversion command selection/help/completion and update governed goldens in internal/cli/cli.go, surface.go, and testdata/help and testdata/completion.
- [x] T020 [P] [US3] Add strict/fatal stdout/new/forced-destination, source integrity/privacy, cancellation/bounds, and historical CLI conversion cases in internal/cli/scripted_conversion_test.go and internal/conformance/scripted_conversion_test.go.
- [x] T021 [US3] Add meaningful new conversion fuzz target(s), preserve old pair goldens/codes, and govern fixed-work runs in internal/convert/scripted_fuzz_test.go and .github/workflows/ci.yml.

## Phase 6: Polish, verification, and delivery

- [x] T022 Update truthful twelve-direction/default/precision/loss/scope documentation and dated architecture decisions in docs/conversion.md, docs/formats/ass-ssa.md, docs/cli.md, testdata/README.md, and CHANGELOG.md.
- [x] T023 Run focused and full foreground CI-parity verification, convergence, schema integrity, encoding/whitespace checks, and record evidence in specs/S026-complete-scripted-conversion/verification.md.
- [x] T024 Commit and explicitly authorized push; publish formatted official PR closing #60/#61, verify body readback, and reconcile native Project Stage/Slice and unused Status in GitHub.
- [ ] T025 Await all current-head CI and every external review; handle findings/resolve threads, request at most one second Codex review, and record final human merge handoff evidence on the official GitHub PR.

## Dependencies and parallel execution

T001 -> T002 -> T003 blocks all implementation. T004/T005 are coordinating integration; T006 is target agent-owned. After the blocking gate, outbound T007/T008/T009, target T006/T012/T013/T014/T015, and evidence T010/T017/T018 can proceed on disjoint files. Evidence tasks share one owner and are sequential within that task. Native outbound analysis establishes #60's shared loss boundary before #61's integrated six-direction acceptance. US3 joins both stories; final verification is coordinator-owned.

US1 parallel example: outbound agent owns analysis/text tests and policies while evidence agent authors outbound expectations; coordinator integrates report/router changes. US2 parallel example: target agent owns boundary/projection/variant files while coordinator integrates text-source accounting and source-aware target rendering. US3 parallel example: evidence agent owns matrix/publication conformance while coordinator owns CLI surface/goldens, docs, and meaningful workflow fuzz integration.

## Implementation Strategy

Deliver US1's four-direction domain outcome first, then validate six US2 directions through the shared report boundary, and integrate all twelve through US3. Tests precede or accompany policy implementation. No release-wide scope is added. T025 is an external current-head delivery gate; final review evidence can be recorded on GitHub without a post-approval administrative commit that changes the reviewed head.
