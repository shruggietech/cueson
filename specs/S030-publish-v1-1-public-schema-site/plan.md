# Implementation Plan: Prepare v1.1.0 public schema and site

**Branch**: `codex/S030-publish-v1-1-public-schema-site` | **Date**: 2026-09-16 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/S030-publish-v1-1-public-schema-site/spec.md`

## Summary

Prepare the reviewed cueson.io artifact for the independently verified v1.1.0 release. Establish one validated release record in the site content authority, add the v1.1.0 release document and exact immutable schema, move seven primary downloads to exact official assets, update current-state source documentation, and extend generated inventories and tests while preserving all historical routes and bytes. Harden later production proof by comparing public metadata with the locally reviewed declaration, hashing the public content manifest, and rejecting any deployment revision other than current `origin/main`. Close preparation #76 through the official pull request; leave production #67 open and perform no production mutation.

## Technical Context

**Language/Version**: TypeScript 5.9, ECMAScript modules on Node.js 24+, Next.js 16.3 static export, repository Go 1.25 tooling

**Primary Dependencies**: React 19.3, Fumadocs, generated MDX content, Playwright 1.63, axe-core, Cloudflare Wrangler 4.131; Go standard-library documentation and publication verifiers

**Storage**: UTF-8/no-BOM JSON and Markdown authorities, immutable schema bytes, deterministic generated MDX/static files and ignored build artifacts

**Testing**: Node built-in test runner, ESLint, deterministic generation/drift check, Next build, Playwright browser/accessibility/responsive/link/metadata checks, artifact verifier, Wrangler dry run, Go tests/vet/staticcheck/govulncheck and repository text checks

**Target Platform**: Static cueson.io artifact served by Cloudflare Workers and Assets; desktop/mobile evergreen browsers; GitHub-hosted Linux validation; later owner-controlled Cloudflare production deployment

**Project Type**: Documentation-driven static web application plus repository verification and deployment automation

**Performance Goals**: Preserve the current static-only request model; generate one deterministic 25-route artifact without network access during build; keep browser checks bounded to the declared route set

**Constraints**: Existing immutable released schema bytes; no mutable schema alias; seven exact approved download links; no pull-request credentials or production mutation; no new tag/release/asset changes; exact current-main selection for later deployment; UTF-8 without BOM and repository line-ending policy

**Scale/Scope**: 22 HTML routes, three immutable schema routes, seven primary downloads, one current release record, maintained current-state documentation, one official PR and one later separately authorized production continuation

## Constitution Check

*GATE: Passed before Phase 0 research and re-checked after Phase 1 design.*

- **I. Lossless Source Preservation**: PASS. No Cue JSON/source-byte behavior changes. Immutable released schema copies are copied byte for byte and verified by exact length/hash.
- **II. Schema and Official Software Discipline**: PASS. Exact released v1.1.0 software/tag/schema identity is already independently verified. The new route uses the immutable release copy without regeneration or normalization; older schemas remain unchanged.
- **III. Common Model Plus Native Fidelity**: PASS. Current public documentation reports stable native/common-model behavior and its limits without changing runtime ownership.
- **IV. No Silent Loss**: PASS. Generated inventories are complete and bijective; missing routes/downloads/schemas, stale adaptations, remote subset declarations and unexpected aliases fail explicitly.
- **V. Test-First Format Work**: PASS. No format runtime change. Site and documentation-verifier expectations are updated first, with positive and negative release/schema/inventory tests before implementation.
- **VI. Portable and Secure Operation**: PASS. Pull requests remain credential-free and non-deploying. Child-process code retains hidden non-interactive execution. Later deployment requires exact current main and preserves Cloudflare state checks.
- **VII. Documentation and Delivery Authority**: PASS. Documentation, generated artifact, tests and release facts converge through Spec Kit, blocking analysis, official CI/reviews and native issue/Project lifecycle. #76 closes independently; #67 and protected production/final-merge authority remain open.

No constitutional exception is required.

## Project Structure

### Documentation (this feature)

```text
specs/S030-publish-v1-1-public-schema-site/
├── checklists/
│   ├── requirements.md
│   └── site-publication.md
├── contracts/
│   ├── public-site-update.md
│   └── production-verification.md
├── data-model.md
├── plan.md
├── quickstart.md
├── research.md
├── spec.md
├── tasks.md
└── verification.md
```

### Source Code (repository root)

```text
README.md
CHANGELOG.md
.github/workflows/site-deploy.yml
docs/
├── architecture.md
├── cli.md
├── compatibility.md
├── conversion.md
├── cueson-media-format-guide.html
├── formats/{srt,webvtt,ass-ssa}.md
├── project-management.md
├── release-process.md
├── release-verification.md
├── releases/v1.1.0.md
├── roadmap.md
├── schema.md
└── Cueson-Project-Specification-v0.0.0.md
schema/releases/
├── v0.0.0/cueson.schema.json
├── v1.0.0/cueson.schema.json
└── v1.1.0/cueson.schema.json
scripts/docs-verify/
├── verify.go
└── verify_test.go
site/
├── app/(home)/page.tsx
├── content-map.json
├── lib/site.ts
├── scripts/{generate,verify-production}.mjs
├── tests/{generator,production-verification}.test.mjs
├── tests/site.spec.ts
└── worker/index.test.ts
```

**Structure Decision**: Extend the established source-to-generated site architecture. `site/content-map.json` remains the maintained inventory authority and gains a validated current-release record. Generated files stay derived and are checked for drift. Immutable schema sources remain under versioned repository paths. Production proof stays in the existing manual workflow and verification script, with pure helper coverage added beside existing site tests.

## Phase 0: Research Conclusions

1. Split repository preparation from production. #76 owns reviewable source/artifact work; #67 retains post-merge Cloudflare mutation, live proof and milestone closure.
2. Treat the public release as an observed fact, not a future candidate. GitHub Release v1.1.0 is final at `7ff45c1d8cd8df377e1fb568b9785286b649fd7c`; issue #66 records independent public-byte and execution verification.
3. Centralize current release identity. A small validated `release` record supplies version, tag and URL; all seven downloads must use the same exact tag/version and approved filenames.
4. Extend rather than replace. Add one release document and one schema while retaining every v0.0.0/v1.0.0 route and historical statement.
5. Keep immutable source bytes authoritative. S028 already admitted `schema/releases/v1.1.0/cueson.schema.json`; activate that existing exact 185,641-byte/digest identity without transformation.
6. Compare public proof to reviewed local truth. The verifier must reject a remote deployment that omits a locally declared route/download/schema even when its own metadata is internally consistent, and must hash the fetched content manifest against the reviewed deployment record.
7. Enforce exact current main. Replace ancestor acceptance with equality between selected revision and freshly fetched `origin/main` immediately before verification/deployment.
8. Preserve frozen preparation contracts. S028/S029 historical Spec Kit artifacts and immutable release notes/decision contracts describe their time-bound candidate workflow and are not rewritten.

## Phase 1: Design Decisions

1. Add `release` to `site/content-map.json` with exact `version`, `tag` and official `url`. Reject malformed versions/tags/URLs, mismatched download versions, duplicate downloads and unapproved filenames during load.
2. Generate release identity into content/deployment metadata and export it to the home page through the existing site data layer. Do not duplicate a mutable version constant in page code.
3. Add `docs/releases/v1.1.0.md` before v1.0.0 in release navigation and add the v1.1.0 immutable schema record after older versions. The declared inventory becomes 22 HTML routes plus three schema routes.
4. Update tests first to expect the final release, seven exact downloads, three exact schemas, retained history, absent latest alias and no current candidate-availability language.
5. Add pure production-verification comparison helpers with focused tests. Verify exact local-vs-remote deployment records, fetch and hash `content-manifest.json`, and drive route/download/schema probes from the locally reviewed inventory.
6. Update `site-deploy.yml` policy to fetch main and require exact equality before artifact proof/deployment. Keep manual dispatch, production environment, Cloudflare preservation/read-back and non-interactive behavior unchanged.
7. Update maintained current-state docs proportionally. Preserve historical release paragraphs and Spec Kit evidence; change only prose whose present-tense candidate/publication claims are now false.
8. Run complete local verification, blocking cross-artifact analysis, independent convergence, official PR exact-head CI and all review rounds before the human merge handoff.

## Verification Strategy

1. Red phase: update generator, site, documentation-verifier, workflow-policy and production-verification tests; demonstrate candidate-era implementation fails relevant expectations.
2. Green phase: implement centralized release metadata, release page/schema/inventory generation, landing/current-state docs and verification hardening.
3. Focused gates: site unit tests, documentation-verifier tests, workflow policy tests, generation drift, schema byte/hash checks and production-verification helper negatives.
4. Full site gates: lint, generation/check, unit, static build, browser/accessibility/responsive/link/metadata checks, artifact proof and Wrangler dry run.
5. Full repository gates: Go test/build/vet/race/staticcheck/govulncheck modules as applicable, docs/github-format policy checks, UTF-8/mojibake/whitespace validation and clean generated state.
6. Spec Kit gates: analyze after tasks; converge after implementation with independent artifact/authority review; record results in `verification.md`.
7. Hosted gates: push authorized branch, open formatted/read-back PR closing #76 and referencing #67, wait for exact-head CI/CodeQL/Site/Release-proof and security/Codex review, resolve every finding, request at most one second Codex round if round one finds issues.

## Lifecycle and Authority

1. #66 remains closed Done with its publication evidence; #51 stays In progress.
2. #76 moves Ready to Specced after analysis, In progress during implementation, PR review after official PR publication, and Done only after merge.
3. #67 stays Backlog through the preparation PR. After human merge it can move to Release verification for a separately authorized exact-main deployment.
4. This slice instruction authorizes branch push and official PR publication. It does not authorize production credentials/mutation or the final merge.

## Complexity Tracking

No constitutional violations or exceptional complexity are introduced.
