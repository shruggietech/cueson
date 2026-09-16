# Tasks: Prepare v1.1.0 public schema and site

**Input**: Design documents from `specs/S030-publish-v1-1-public-schema-site/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Test-first coverage is required by the specification for release/schema inventories, current documentation, production verification and deployment policy.

**Organization**: Tasks are grouped by user story so each outcome retains an independently reviewable verification surface. Production mutation is outside this task list.

## Phase 1: Setup and red evidence

**Purpose**: Freeze the governed baseline and make the required behavior fail before implementation.

- [x] T001 Reconfirm exact release #66 evidence, issue #76 parent/dependencies/Project state, open-issue inventory, baseline revision and protected authority boundary; record the result in `specs/S030-publish-v1-1-public-schema-site/verification.md`.
- [x] T002 [P] [US1] Update release/document/download expectations in `site/tests/generator.test.mjs` and `site/tests/site.spec.ts` for v1.1.0, retained history, exact seven downloads and absence of current candidate language.
- [x] T003 [P] [US1] Add published-state and negative drift expectations to `scripts/docs-verify/verify_test.go` while preserving frozen S028/S029 historical contracts.
- [x] T004 [US2] Update schema inventory/byte/hash/alias expectations in `site/tests/generator.test.mjs`, `site/tests/site.spec.ts` and `site/worker/index.test.ts` for all three immutable schemas after T002 establishes shared release expectations.
- [x] T005 [US3] Add focused local-authority, incomplete-remote-inventory, manifest-digest and exact-main workflow-policy expectations in `site/tests/production-verification.test.mjs` and `site/tests/generator.test.mjs` after T002/T004 complete their shared-file edits.
- [x] T006 Run the focused tests, prove the candidate-era implementation fails the new expectations for the intended reasons, and record concise red evidence in `specs/S030-publish-v1-1-public-schema-site/verification.md`.

---

## Phase 2: Foundational requirements gate

**Purpose**: Ensure the written requirements and task coverage are complete before implementation.

- [x] T007 Complete independent review of `specs/S030-publish-v1-1-public-schema-site/checklists/site-publication.md`; resolve every requirements-quality gap without treating implementation as checklist completion.
- [x] T008 Run the Spec Kit blocking analysis across `spec.md`, `plan.md`, `tasks.md`, checklists and contracts; resolve every critical/high or material coverage/consistency finding before source implementation.
- [x] T009 Move issue #76 from Ready to Specced, retain Slice S030 and clear default Status after successful analysis; read back Project and native relationships.

**Checkpoint**: Requirements, contracts and task traceability are complete; implementation may begin.

---

## Phase 3: User Story 1 - Find and use the current release (Priority: P1)

**Goal**: Present independently verified v1.1.0 consistently across landing, generated documentation and exact primary downloads while preserving history.

**Independent Test**: Generate/build the site and inspect release identity, four-format claims, v1.1.0 release route and all seven exact download targets while older release routes remain present.

- [x] T010 [US1] Add and validate one current release record plus seven exact v1.1.0 primary downloads in `site/content-map.json` and `site/scripts/generate.mjs`, including uniqueness and version/tag/filename consistency failures.
- [x] T011 [US1] Add `docs/releases/v1.1.0.md` to the generated release inventory and update `site/lib/site.ts` plus `site/app/(home)/page.tsx` to derive current release identity and four-format content from the maintained authority.
- [x] T012 [P] [US1] Reconcile present-tense release/capability/install/schema/compatibility/conversion/format statements in `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, `CHANGELOG.md`, `docs/architecture.md`, `docs/cli.md`, `docs/schema.md`, `docs/compatibility.md`, `docs/conversion.md`, `docs/formats/srt.md`, `docs/formats/webvtt.md`, `docs/formats/ass-ssa.md`, `docs/cueson-media-format-guide.html`, `docs/roadmap.md`, `docs/release-process.md`, `docs/release-verification.md`, `docs/project-management.md` and `docs/Cueson-Project-Specification-v0.0.0.md` without rewriting historical release/Spec Kit evidence.
- [x] T013 [US1] Update `scripts/docs-verify/verify.go` to require durable published/verified v1.1.0 markers across README, contributing, security and maintained contract surfaces, exact current release/download facts and retained historical compatibility while removing obsolete candidate-only assertions.
- [x] T014 [US1] Run focused generator, documentation-verifier and browser tests; resolve failures and record green US1 evidence in `specs/S030-publish-v1-1-public-schema-site/verification.md`.

**Checkpoint**: Current release identity and user-facing release/download documentation are independently complete.

---

## Phase 4: User Story 2 - Resolve every immutable schema (Priority: P1)

**Goal**: Add the exact immutable v1.1.0 schema route while preserving exact v0.0.0/v1.0.0 bytes and rejecting mutable aliases.

**Independent Test**: Generate the artifact and prove three exact versioned schema files by byte length/hash plus absent latest/unknown fallbacks.

- [x] T015 [US2] Verify the existing S028-admitted `schema/releases/v1.1.0/cueson.schema.json` remains byte-identical to `internal/schema/cueson.schema.json` at exact 185,641-byte length and SHA-256 `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7`; do not recopy or transform it.
- [x] T016 [US2] Add v1.1.0 schema length/hash/source/public metadata to `site/content-map.json` and extend `site/scripts/generate.mjs` validation and generated content/deployment records while retaining both older schema entries.
- [x] T017 [US2] Extend artifact, Worker, browser and route checks in `site/tests/generator.test.mjs`, `site/tests/site.spec.ts`, `site/worker/index.test.ts` and `site/scripts/verify-artifact.mjs` as needed for exactly three immutable schemas and no latest alias.
- [x] T018 [US2] Run focused schema generation, unit, Worker, browser and artifact checks; prove all three byte identities and record green US2 evidence in `specs/S030-publish-v1-1-public-schema-site/verification.md`.

**Checkpoint**: All three immutable schema contracts are present, exact and independently testable.

---

## Phase 5: User Story 3 - Review an exact deployable update (Priority: P2)

**Goal**: Produce a reproducible reviewed artifact and harden later production selection/read-back without deploying it.

**Independent Test**: Build the artifact at the slice head, compare generated deployment/content metadata to declared inventories, exercise negative remote-subset/manifest/revision cases and complete a Wrangler dry run.

- [x] T019 [US3] Refactor `site/scripts/verify-production.mjs` to export pure comparison/hash helpers, load local deployment metadata as authority, compare exact remote metadata and hash public `content-manifest.json`; drive route/download/schema probes from local inventories.
- [x] T020 [US3] Complete positive/negative helper coverage in `site/tests/production-verification.test.mjs` for missing, extra, reordered or changed remote inventory and manifest digest mismatch.
- [x] T021 [US3] Change `.github/workflows/site-deploy.yml` to require selected revision equality with freshly fetched `origin/main`, update workflow policy coverage, and preserve manual/protected credential, Cloudflare preflight/read-back and no-PR-deploy boundaries.
- [x] T022 [US3] Run the complete site suite (`lint`, generation/check, unit, build, browser/accessibility/responsive/link/metadata, artifact and Wrangler dry run) and record exact route/download/schema counts in `specs/S030-publish-v1-1-public-schema-site/verification.md`.

**Checkpoint**: The reviewed artifact is ready for a later separately authorized exact-main production continuation.

---

## Phase 6: Cross-cutting verification and convergence

**Purpose**: Prove repository-wide quality, reconcile governance and publish one review-ready head.

- [x] T023 Run all applicable root and standalone Go module tests, build/vet/race/staticcheck/govulncheck gates, documentation/github-format policy tests and repository UTF-8/no-BOM/mojibake/whitespace checks; record commands and outcomes in `specs/S030-publish-v1-1-public-schema-site/verification.md`.
- [x] T024 Confirm generated files are current, immutable historical schema/release/Spec Kit bytes are unchanged, no production credential/state changed and `git diff --check` plus working-tree review are clean enough to commit.
- [x] T025 Run Spec Kit convergence with independent review of requirements, issue acceptance, route/download/schema identities, documentation truth, production boundary and changed artifacts; resolve every material finding and update `tasks.md` plus `verification.md`.
- [x] T026 Commit conventionally on `codex/S030-publish-v1-1-public-schema-site`, verify the exact committed tree locally, and ensure #76 remains In progress/S030 with default Status empty before publication.

---

## Phase 7: Official pull request and review convergence

**Purpose**: Publish the explicitly authorized official PR and reach a terminal reviewed exact head.

- [ ] T027 Push the authorized branch and publish a github-format/read-back official PR with `Closes #76` and `Refs #67`; move #76 to PR review and clear default Status.
- [ ] T028 Wait for all exact-head CI, CodeQL, Site, Release-proof, security and first-round Codex results; inspect every review/comment/reaction and address every actionable finding with focused plus complete re-verification.
- [ ] T029 If round one contains findings, request exactly one second Codex review with `@codex review`, then handle every second-round result without requesting a third automatic round; resolve threads only after concerns are satisfied.
- [ ] T030 Publish formatted/read-back completion evidence on the PR, confirm all required checks are green at the final head, terminal reviews have no actionable finding, #67/#51/milestone remain open and production is unchanged, then hand off for the human final review and merge ritual.

## Dependencies and Execution Order

- T001-T006 establish red evidence. T007-T009 are blocking requirements/analysis gates.
- US1 and US2 share `site/content-map.json` and `site/scripts/generate.mjs`, so T010-T018 execute in task order despite independent acceptance surfaces.
- US3 depends on the final local deployment record from US1/US2.
- T023-T026 require all three stories complete. T027-T030 require the committed converged head.
- Production deployment, live verification, #67/#51 closure and milestone closure are intentionally outside this task list.

## Parallel Opportunities

- T002-T005 use separate test owners and can proceed in parallel.
- Documentation reconciliation in T012 can be independently reviewed while T010-T011 establish release metadata, provided the coordinator integrates and reruns all checks.
- Requirements checklist review and independent convergence are suitable for separate read-only/reviewer ownership.

## Implementation Strategy

1. Freeze governance and make focused tests fail.
2. Pass analysis before implementation.
3. Deliver current release/documentation identity.
4. Add exact immutable schema hosting.
5. Harden exact-main production proof without mutation.
6. Run full verification and independent convergence.
7. Publish and converge the authorized official PR.
8. Stop before merge and production.
