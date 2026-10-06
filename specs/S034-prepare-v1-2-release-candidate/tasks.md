# Tasks: Prepare the v1.2.0 release candidate

**Input**: [spec.md](spec.md), [plan.md](plan.md), research, data model and contracts.

## Phase 1: Setup

- [x] T001 Establish exact S034 branch/directory, current authority and fresh issue inventory in specs/S034-prepare-v1-2-release-candidate/spec.md. [FR-012]
- [x] T002 Establish milestone, candidate/publication/hosting issues, native dependencies and unique Project stages in specs/S034-prepare-v1-2-release-candidate/issue-map.json. [FR-012]

## Phase 2: Foundation

- [x] T003 Complete specify/clarify/checklist/plan/research/design in specs/S034-prepare-v1-2-release-candidate/. [FR-013]
- [x] T004 Run blocking read-only analysis of spec.md, plan.md and tasks.md; resolve all critical/high findings before implementation. [FR-013]

## Phase 3: US1 - Exact stable contract

**Independent test**: Matching executable/schema/immutable bytes, preserved historical inputs and former development refusal.

- [x] T005 [P] [US1] Add candidate equality/history/dev-refusal tests in internal/schema/stable_candidate_test.go and internal/conformance/scripted_stable_test.go. [FR-001, FR-002, FR-003]
- [x] T006 [US1] Promote only current identity in internal/schema/, internal/version/, internal/model/consumer.go and current expectations/fixtures; add schema/releases/v1.2.0/cueson.schema.json. [FR-001, FR-002, FR-003]
- [x] T007 [US1] Run focused schema/model/conformance/CLI identity and historical tests; preserve released bytes. [FR-001, FR-002, FR-003]

## Phase 4: US2 - Same-bundle package proof

**Independent test**: Six packages and three native hosts exercise current/historical behavior and authenticated public old-consumer refusal.

- [x] T008 [P] [US2] Add stable 1.2.0 dispatch, old-consumer authentication and consumer-proof regression tests in scripts/release-verify/native_test.go and policy_test.go. [FR-004, FR-005, FR-006, FR-007, FR-008]
- [x] T009 [US2] Freeze checksummed published 1.1.0 host asset contract in specs/S034-prepare-v1-2-release-candidate/contracts/published-v1.1.0-consumer.json; preserve existing 1.0.0 proof. [FR-006]
- [x] T010 [US2] Implement stable 1.2.0 proof selection, same-bundle annotations, six historical paths, version-aware published 1.1.0 positive/refusal probes and evidence in scripts/release-verify/native.go and verify.go. [FR-004, FR-005, FR-006]
- [x] T011 [US2] Restore immutable stable schema packaging and exact 1.2.0 nonpublishing workflow in .goreleaser.yaml and .github/workflows/release-proof.yml. [FR-002, FR-007, FR-008, FR-014]
- [x] T012 [US2] Add new exact candidate evidence schema in specs/S034-prepare-v1-2-release-candidate/contracts/release-evidence-contract.schema.json and validate policy tests. [FR-007, FR-008, FR-011]
- [x] T013 [US2] Run focused release-verifier tests/vet/static/security analysis with unchanged historical proof contracts. [FR-004, FR-005, FR-006, FR-007, FR-008]

## Phase 5: US3 - Publication decision

**Independent test**: Prepared current documentation and immutable public history agree; actual publication requires fresh main evidence and concrete authority.

- [x] T014 [P] [US3] Add prepared 1.2.0 changelog/release-note/current-identity validation tests in scripts/docs-verify/prepared_changelog_test.go and verify_test.go. [FR-009, FR-010]
- [x] T015 [P] [US3] Finalize dated CHANGELOG.md, docs/releases/v1.2.0.md, current contract/examples and chronological roadmap/release status; preserve public site/download inventory and dated historical records. [FR-009, FR-010, FR-014]
- [x] T016 [US3] Implement prepared 1.2.0 documentation verification in scripts/docs-verify/prepared_changelog.go and verify.go. [FR-009, FR-010]
- [x] T017 [US3] Prepare exact intended inventory, evidence/expiry/changed-state/recovery and authority contract in specs/S034-prepare-v1-2-release-candidate/contracts/publication-decision.md. [FR-011, FR-014]
- [x] T018 [US3] Verify notes/body formatting, docs links/examples/status and public version invariants through scripts/docs-verify and scripts/github-format. [FR-009, FR-010, FR-011, FR-014]

## Phase 6: Integration and delivery

- [x] T019 Run foreground required product and six nested-module tests/vet/static/security, repository text/brand, meaningful fuzz, workflow lint and six pure-Go build checks; record results in verification.md. [FR-013]
- [x] T020 Run site lint/generation/tests/build/artifact/browser/accessibility checks with public v1.1 inventory retained; preserve existing unrelated preview state. [FR-010, FR-014]
- [x] T021 Independently inspect integrated diff, historical hashes and quickstart; correct scoped findings and commit conventionally on the exact S034 branch. [FR-001, FR-002, FR-003, FR-013]
- [ ] T022 Publish the authorized official candidate PR with Closes #84, formatted/read-back body and correct Project stage; complete final-head hosted checks and same-bundle native proof. [FR-007, FR-008, FR-012, FR-013]
- [ ] T023 Resolve every review finding, at most one finding-driven second Codex request, and present the actual reviewed candidate/publication boundary with executed evidence on the PR. [FR-011, FR-013, FR-014]

## Dependencies and parallel execution

T001-T004 precede product edits. US1 identity and US2 verifier tests can proceed in parallel with separate file ownership; final US2 execution requires US1. US3 documentation can proceed alongside both, but final identities/digests depend on the immutable candidate bytes. Root owns docs verifier to avoid cross-agent conflicts. T019-T023 follow integrated story completion.

The smallest usable increment is US1, but it does not complete candidate acceptance without all three stories. Parallel examples: schema agent handles T005-T007; proof agent handles T008-T013; documentation agent handles T015; root handles T014/T016-T018 and final verification. No agent commits/pushes independently.

## Runtime evidence and checkbox lifecycle

Mark executed local tasks when complete. Final-head hosted/review tasks T022/T023 are runtime gates: their authoritative completion is recorded on the final PR body after results arrive, avoiding a new evidence-only source commit that would invalidate the exact head and restart review. Unchecked runtime markers are not a claim of completion or permission to skip gates.
