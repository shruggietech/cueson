# Tasks: Prepare v1.2.0 public schema and site

**Input**: [spec.md](spec.md), [plan.md](plan.md), research, data model and production handoff.

**Tests**: Existing negative/security/accessibility checks remain required. Extend independently fixed inventories before changing authority.

## Phase 1: Setup

- [x] T001 Inspect constitution, architecture, active issues/milestone/native dependencies and verified public release; establish aligned branch/directory/issue mapping in `issue-map.json`.
- [x] T002 Execute installed specify, clarify, checklist, plan and tasks workflows; independently review `checklists/publication.md` and pass blocking read-only analysis.

## Phase 2: Foundational

- [x] T003 Bind public release/tag/schema/primary asset evidence in `verification.md`; preserve prior bytes and distinguish live before-state.

## Phase 3: US1 - Published release discovery

**Independent test**: Current metadata and both new public pages match shipped behavior and the consumer contract.

- [x] T004 [US1] Extend release/download/navigation/schema expectations in `site/tests/generator.test.mjs` and `site/tests/site.spec.ts`; keep version-mismatch mutation meaningful and observe expected failure before authority promotion.
- [x] T005 [US1] Promote v1.2.0 release/download authority and add release/speaker routes in `site/content-map.json`.
- [x] T006 [P] [US1] Reconcile current publication prose in `README.md`, `SECURITY.md`, `CONTRIBUTING.md`, `docs/architecture.md`, `docs/schema.md`, `docs/compatibility.md`, `docs/consumer-speakers.md`, `docs/releases/v1.2.0.md`.
- [x] T007 [P] [US1] Record public publication evidence and reconcile lifecycle prose in `docs/release-verification.md`, `docs/release-process.md`, `docs/roadmap.md`, `docs/project-management.md`, plus `CHANGELOG.md` Unreleased.

## Phase 4: US2 - Historical preservation

**Independent test**: Four fixed immutable identities and prior release routes/links all pass.

- [x] T008 [US2] Extend independent schema inventories and unchanged historical release expectations in `site/tests/generator.test.mjs`, `site/tests/site.spec.ts`, `site/worker/index.test.ts`.
- [x] T009 [US2] Compare all four source/generated schema sizes/hashes and verify historical release links through static/browser checks; record in `verification.md`.

## Phase 5: US3 - Exact production handoff

**Independent test**: Altered or omitted new metadata/routes/schemas reject; protected workflow boundary remains intact and handoff specifies actual-main/live proof.

- [x] T010 [US3] Update sample authority and add new-schema/consumer-route negative mutations in `site/tests/production-verification.test.mjs` while retaining all credential/revision/infrastructure checks.
- [x] T011 [US3] Finalize `contracts/production-handoff.md`, keep #86/milestone pending and native Project fields reconciled.

## Phase 6: Verification and delivery

- [x] T012 Run full site lint/generation/drift/unit/build/browser/accessibility/artifact/dry-run pipeline in foreground using isolated preview; record exact results in `verification.md`.
- [x] T013 Run applicable documentation/publication/module checks, UTF-8/no-BOM/mojibake and immutable/history scope checks; record in `verification.md`.
- [ ] T014 Obtain independent implementation review, address findings, commit and push authorized branch; publish formatted official PR closing #88 with runtime review/check evidence.
- [ ] T015 Monitor exact-head required CI and every external review finding; fix/verify/resolve handled threads and request at most one second Codex round if first round has findings; stop for final review/merge with production #86 pending.

## Phase 7: Continue immediately after owner merge

- [ ] T016 Verify owner-confirmed merge, perform scoped housekeeping and fresh actual-main proof; bind the recorded S035 production authorization in `verification.md`.
- [ ] T017 Dispatch protected `site-deploy.yml` for exact current main, watch proof/preflight/deploy/read-back to completion, and record actual production evidence in `verification.md` and #86.
- [ ] T018 Run independent public route/download/revision/content-manifest/four-schema verification; explicitly probe content type, immutable cache and nosniff headers, redirects and missing latest schema; demonstrate HTTPS independently over IPv4 and IPv6 using the S022 contract. Reconcile #86, Project and milestone only after all live results pass.

## Dependencies and parallel execution

T001 → T002 → T003 → T004 → T005. T006 and T007 have disjoint ownership and may run in parallel after the analysis gate. T008 follows test inventory integration; T010 is independent of prose edits. T009/T012/T013 require integrated authorities and prose. T014 follows local verification; T015 follows push/PR. US2 preserves historical compatibility alongside US1, and US3 governs later live completion.

The owner retains final merge at T015. Their subsequent steering authorizes T016–T018 immediately after that merge without another kickoff or approval. S035 remains open across this pause and is not complete until production verification passes.

## Implementation strategy

Complete all three reviewable preparation stories under the existing static architecture, then use local and exact-head hosted checks for acceptance. Do not substitute preparation for actual deployment or claim live completion. T014/T015 runtime outcomes are recorded in the final PR read-back; checking them in source before they occur would be premature.
