---

description: "Implementation tasks for the Cueson public site and domain launch"
---

# Tasks: Launch cueson.io

**Input**: Design documents from `/specs/S021-launch-cueson-site/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Generator, Worker, browser, artifact, and hosted verification are required by the feature specification.

**Organization**: Tasks are grouped by user story. Protected post-merge production execution is recorded as the continuation ritual after all implementation tasks complete.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it owns different files and has no dependency on an incomplete task
- **[Story]**: Maps the task to a user story in [spec.md](spec.md)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish the isolated site package and deterministic toolchain.

- [X] T001 Create the pinned pnpm package and lockfile in `site/package.json` and `site/pnpm-lock.yaml`
- [X] T002 [P] Configure TypeScript, Next.js static export, Fumadocs MDX, Tailwind CSS, and Playwright in `site/tsconfig.json`, `site/next.config.mjs`, `site/source.config.ts`, and `site/playwright.config.ts`
- [X] T003 [P] Add generated-site exclusions and byte-protected public assets to `.gitignore` and `.gitattributes`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Build the single deterministic input and route authority used by every user story.

**CRITICAL**: No user-story phase begins until this phase is complete.

- [X] T004 Write failing generator safety, mapping, determinism, and drift tests in `site/tests/generator.test.mjs`
- [X] T005 Define authoritative documents, routes, released schemas, and approved brand copies in `site/content-map.json`
- [X] T006 Implement bounded deterministic content, schema, brand, manifest, and revision generation in `site/scripts/generate.mjs`
- [X] T007 Implement the Fumadocs source loader and shared metadata/route helpers in `site/lib/source.ts`, `site/lib/site.ts`, and `site/components/mdx.tsx`
- [X] T008 Create the static preview server and common browser-test utilities in `site/scripts/serve.mjs` and `site/tests/helpers.ts`

**Checkpoint**: A clean checkout can generate and drift-check all disposable site inputs without modifying an authoritative source.

---

## Phase 3: User Story 1 - Learn and use Cueson from its public home (Priority: P1)

**Goal**: Deliver the branded landing page and complete navigable documentation experience.

**Independent Test**: Build and serve the static artifact, then verify every planned HTML route, navigation path, link, metadata field, accessibility rule, and representative viewport.

### Tests for User Story 1

- [X] T009 [US1] Write failing landing, documentation-route, link, metadata, accessibility, keyboard, and responsive browser tests in `site/tests/site.spec.ts`

### Implementation for User Story 1

- [X] T010 [P] [US1] Implement the global provider, fonts, metadata, navigation, and error surface in `site/app/layout.tsx`, `site/app/global.css`, `site/app/not-found.tsx`, and `site/lib/layout.shared.tsx`
- [X] T011 [P] [US1] Implement the product landing experience and release downloads in `site/app/(home)/layout.tsx` and `site/app/(home)/page.tsx`
- [X] T012 [US1] Implement the statically generated Fumadocs layout and pages in `site/app/docs/layout.tsx` and `site/app/docs/[[...slug]]/page.tsx`
- [X] T013 [P] [US1] Integrate approved logos, social preview, favicons, fonts, colors, spacing, focus, and reduced-motion behavior through `site/content-map.json` and `site/app/global.css`
- [X] T014 [P] [US1] Publish the authoritative media-format guide through the generator and link it from `site/app/(home)/page.tsx` and the shared navigation
- [X] T015 [US1] Make all User Story 1 browser checks pass against the clean static export in `site/tests/site.spec.ts`

**Checkpoint**: A visitor can learn, install, download, and navigate Cueson documentation from a responsive and accessible static artifact.

---

## Phase 4: User Story 2 - Resolve an immutable released schema (Priority: P1)

**Goal**: Serve exact v0.0.0 and v1.0.0 schema bytes at their canonical versioned paths without a mutable alias.

**Independent Test**: Compare generated and served schema response bytes, lengths, hashes, content types, cache policy, and missing-alias behavior with the immutable repository sources.

### Tests for User Story 2

- [X] T016 [US2] Write failing schema-copy, hash, route, content-type, cache, and no-latest-alias tests in `site/tests/generator.test.mjs`, `site/tests/site.spec.ts`, and `site/worker/index.test.ts`

### Implementation for User Story 2

- [X] T017 [US2] Generate only the two declared byte-exact versioned schemas and record their identities in `site/scripts/generate.mjs` and `site/content-map.json`
- [X] T018 [US2] Verify schema inventory, raw byte identity, and absence of mutable aliases in `site/scripts/verify-artifact.mjs`
- [X] T019 [US2] Make all User Story 2 generator, browser, Worker, and artifact checks pass in `site/tests/generator.test.mjs`, `site/tests/site.spec.ts`, `site/worker/index.test.ts`, and `site/scripts/verify-artifact.mjs`

**Checkpoint**: Both canonical released schema identifiers resolve from the built artifact with exact immutable bytes, and `latest` does not exist.

---

## Phase 5: User Story 3 - Publish and verify an owner-controlled site (Priority: P1)

**Goal**: Produce a reviewed Cloudflare deployment unit that cannot deploy from pull requests and can be independently verified after human merge.

**Independent Test**: Exercise Worker behavior locally, inspect workflow triggers and permissions, perform a Wrangler dry run, and run the production verifier against a controlled local origin before the post-merge Cloudflare continuation.

### Tests for User Story 3

- [X] T020 [US3] Write failing apex response-policy and path/query-preserving `www` redirect tests in `site/worker/index.test.ts`
- [X] T021 [US3] Write failing deployment-workflow authority and revision-binding assertions in `site/tests/generator.test.mjs`

### Implementation for User Story 3

- [X] T022 [US3] Implement the static-asset pass-through, response policy, schema caching, deployment-metadata caching, and `www` redirect in `site/worker/index.ts`
- [X] T023 [US3] Configure exact Workers Static Assets, dry-run, and Custom Domain behavior in `site/wrangler.jsonc`
- [X] T024 [P] [US3] Add read-only pull-request and default-branch site validation in `.github/workflows/site.yml`
- [X] T025 [P] [US3] Add manual exact-main-revision production deployment with no pull-request trigger in `.github/workflows/site-deploy.yml`
- [X] T026 [US3] Implement repeatable local and public DNS, TLS, HTTP, redirect, download, route, metadata, and schema verification in `site/scripts/verify-production.mjs`
- [X] T027 [US3] Make Worker, authority, production-verifier fixture, and Wrangler dry-run checks pass in `site/worker/index.test.ts`, `site/tests/generator.test.mjs`, `site/scripts/verify-production.mjs`, and `site/wrangler.jsonc` without mutating Cloudflare

**Checkpoint**: The repository contains a production-ready, owner-controlled deployment unit and an independent public verifier, while pull requests remain validation-only.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Reconcile documentation, governance, and the complete verification story.

- [X] T028 [P] Record the public-site architecture, generated-content boundary, Cloudflare target, deployment authority, schema publication, and verification process in `docs/architecture.md`, `docs/release-process.md`, and `docs/schema.md`
- [X] T029 [P] Record the S021 architecture and domain-activation decision under `[Unreleased]` in `CHANGELOG.md`
- [X] T030 Run the complete quickstart, repository CI-parity, brand, formatting, UTF-8, mojibake, whitespace, site, browser, artifact, and Wrangler dry-run verification from `specs/S021-launch-cueson-site/quickstart.md`

---

## Protected Post-Merge Continuation

After the human operator merges the green and fully reviewed pull request, the same S021 lifecycle selects the exact merged `main` commit, repeats T030 from a clean checkout, performs the authorized Cloudflare deployment described by [deployment.md](contracts/deployment.md), waits for DNS and certificate readiness, runs the independent public verifier, and reconciles issue #47 and its Project item. This protected continuation is deliberately outside the pull-request implementation checklist because it cannot occur before the human merge ritual.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Starts immediately.
- **Foundational (Phase 2)**: Depends on Phase 1 and blocks all user stories.
- **User Story 1 (Phase 3)**: Depends on Phase 2.
- **User Story 2 (Phase 4)**: Depends on Phase 2 and shares the generated artifact with User Story 1.
- **User Story 3 (Phase 5)**: Depends on the generated site and schema artifact from User Stories 1 and 2.
- **Polish (Phase 6)**: Depends on all three user stories.
- **Protected continuation**: Depends on a human merge of the fully reviewed, green pull request.

### User Story Dependencies

- **User Story 1**: Independently demonstrates the public product and documentation artifact after Phase 2.
- **User Story 2**: Independently demonstrates immutable schema publication after Phase 2.
- **User Story 3**: Uses the complete static artifact from User Stories 1 and 2 because deployment must publish both together.

### Parallel Opportunities

- T002 and T003 can proceed in parallel after T001 establishes package boundaries.
- T010 and T011 own separate layout and landing files after T009.
- T013 and T014 own independent brand and guide integration surfaces.
- T024 and T025 own separate GitHub workflows after the Worker and Wrangler contract is fixed.
- T028 and T029 own different documentation files after implementation stabilizes.

## Implementation Strategy

### MVP First

1. Complete Setup and Foundational phases.
2. Complete User Story 1 and validate the public product/documentation artifact.
3. Complete User Story 2 and prove immutable schema identity.
4. Complete User Story 3 and prove the deployment boundary without production mutation.
5. Complete documentation and full verification, publish the pull request, and resolve at most two Codex review rounds.
6. Stop for the human merge ritual, then resume the protected S021 continuation and activate the domain from merged `main`.

## Notes

- Every checklist task names its owned file or verification artifact.
- Test tasks precede their corresponding implementation tasks.
- Generated documentation, schemas, brand copies, and build output remain ignored and reproducible.
- Production success is not claimed until post-merge Cloudflare state and independent public probes pass.
