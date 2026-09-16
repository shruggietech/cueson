# Verification: Prepare v1.1.0 public schema and site

## Governed baseline

- Started 2026-09-16 from clean local/default-remote main `7ff45c1d8cd8df377e1fb568b9785286b649fd7c` on branch `codex/S030-publish-v1-1-public-schema-site`.
- Issue #66 is closed with independently downloaded public-byte, exact package/schema/legal identity and Windows execution evidence for final GitHub Release `v1.1.0`. Annotated tag object `e595ac75c04fa96f5373bf451c4b5cd87fde9471` peels to the exact baseline. Exactly six archives, six SPDX JSON SBOMs and one checksum manifest remain published.
- Issue #76 is open, assigned, milestone v1.1.0, child of #67, historically blocked by closed #66 and currently blocking #67. It is the only new issue since the prior roadmap assessment. Open set is #51/#67/#76; no unassessed external issue arrived.
- Project `cueson Delivery` contains every in-scope issue once. #51 is In progress, #66 is Done/S029, #67 is Backlog/S030 and #76 is Ready/S030 at kickoff; default Status is cleared. #67, #51 and milestone v1.1.0 remain open.
- User authorization covers exact v1.1.0 tag/release/assets already completed, and branch push plus official S030 PR publication. No production cueson.io mutation or final PR merge is authorized.

## Specification and clarification

- Installed individual Spec Kit skills were used because the repository-mandated `shruggie-speckit` orchestration skill is not installed in this session.
- `.specify/extensions.yml` is absent, so no before/after hooks apply.
- Specification and built-in requirements checklist passed in one iteration with 17 functional requirements, seven measurable outcomes, three independently testable stories and no unresolved clarification marker.
- Clarification fixed the #76/#67 split, exact released identity, three-version immutable route set and separate production authority.

## Planning and requirements-quality review

- Planning produced research, data model, quickstart and two contracts. Constitution gates I-VII pass with no exception.
- Independent reviewer completed all 27 custom site-publication requirements-quality criteria with no gap, count contradiction, encoding defect or out-of-scope edit.

## Blocking analysis

- Analysis covered 17 functional requirements, seven buildable success criteria, three user stories, 30 tasks, two contracts and seven constitutional principles.
- Requirement/task coverage is 24/24 (100%). No ambiguity marker, duplicate requirement, constitutional conflict, uncovered buildable requirement or unmapped task was found.
- Finding A1 (high): T002/T004/T005 were incorrectly marked parallel while sharing `site/tests/generator.test.mjs`. Remediation removed parallel markers from T004/T005 and made their ordering explicit.
- Finding U1 (medium): T005 named an existing workflow-policy test without an exact path. Remediation identified `site/tests/generator.test.mjs` explicitly.
- Post-remediation analysis has zero critical, high, medium or material findings. Issue #76 moved to Specced/S030 and default Status was cleared after read-back.

## Test-first evidence

- Site generator focused run produced three passes and four intended failures: v1.1.0 schema output absent, maintained release record absent, generated route count 23 instead of 25 and deployment workflow still accepted any main ancestor.
- Production-verification helper test failed at module load because `assertDeploymentMatches` and `assertContentManifestDigest` do not yet exist/export. The test covers exact ordered local-vs-remote metadata, missing/extra/reordered inventory, changed release identity, missing schema and public manifest digest mismatch.
- Worker immutable-header/absent-alias coverage passed 5/5, proving existing generic behavior already accommodates a declared v1.1.0 route and unknown aliases.
- Documentation-verifier focused tests failed for the intended missing current-state rules: no required `v1.1.0 released and independently verified` marker and no stale-claim rejection for `Published downloads below remain v1.0.0`. The frozen S028/S029 evidence scope guard passed.
- Red-test changes pass formatting, UTF-8/no-BOM, mojibake and `git diff --check` sanity checks.

## Implementation evidence

- The maintained authority now declares final release `v1.1.0`, seven exact primary downloads and three immutable schema identities. Generation validates release/tag/URL/filename consistency, unique ordered inventories and exact schema length/hash/source/public records.
- The existing S028-admitted v1.1.0 schema was intentionally activated without recopying or transformation. Both `internal/schema/cueson.schema.json` and `schema/releases/v1.1.0/cueson.schema.json` remain exactly 185,641 bytes with SHA-256 `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7`.
- Present-tense documentation now records the independently verified public v1.1.0 release at `7ff45c1d8cd8df377e1fb568b9785286b649fd7c`, four stable formats, twelve conversions, exact release assets, compatibility boundary and pending production work in #67. Frozen S028/S029 evidence remains historical.
- Documentation verification passed: 25 documents, 131 local links, 10 registered examples and 34 format rows. The documentation verifier rejects stale candidate and v1.0-current claims while preserving historical evidence.
- Site unit verification passed 83/83, including exact inventories, release/version asset consistency, negative drift, production-authority comparison, manifest-digest mismatch and immutable Worker behavior.
- The complete site pipeline passed lint, deterministic generation/check, unit tests, Next production build, Playwright at three viewports (81 passed, six intentionally project-scoped skips), artifact verification and Wrangler dry run. The deployable record contains exactly 25 routes: 22 HTML routes and three immutable schemas, plus seven downloads.
- Production verification now compares the complete public deployment metadata to the local reviewed authority, verifies public manifest bytes by digest and drives all public probes from the local inventory. The deployment workflow accepts only the freshly fetched exact `origin/main` revision. No Cloudflare or cueson.io mutation occurred.

## Convergence and publication evidence

- Cross-cutting verification passed root and standalone-module tests, root vet and race tests, all sixteen fixed-work fuzz targets at 1,000 iterations each, six pure-Go target cross-builds, standalone vet/staticcheck/govulncheck gates, documentation and brand verification, GitHub workflow lint, github-format, UTF-8/no-BOM/mojibake checks and `git diff --check`. Govulncheck found no called vulnerability; its informational unreachable imported-module report does not affect the built program.
- Generated-content and artifact checks are current. The final site run passed 83/83 unit tests, 81 browser tests with six intentionally project-scoped skips, exact artifact verification for 25 declared routes (22 HTML plus three schemas), seven release downloads and a Wrangler dry run. A transient Windows `ENOTEMPTY` during generated-directory replacement prompted bounded retry hardening; the complete rerun passed.
- Historical `v0.0.0` and `v1.0.0` schema bytes, all release evidence outside the maintained v1.1.0 page, and S028/S029 Spec Kit evidence remain unchanged. No production credential, Cloudflare state or cueson.io content changed.
- Independent convergence reviewed all 17 functional requirements, seven success criteria, three user stories with nine scenarios, six issue acceptance criteria, two contracts and seven constitutional principles. Review findings corrected durable release/governance wording, the #76/#67 ownership split, version-derived asset names, deployment metadata alignment, and exact-main checks immediately before proof and again after Cloudflare preflight immediately before mutation. The final audit reports zero critical, high, medium or low actionable findings.
- Prepublication Project read-back confirmed #76 was In progress/S030 with default Status unused, while #67 remained Backlog/S030 and production-authorized separately.
- The authorized branch was pushed and official pull request #77 was published with `Closes #76` and `Refs #67`. Its UTF-8 Markdown body passed github-format and immediate read-back confirmed 49 intentionally structured lines after correcting a PowerShell capture that had collapsed the first publication attempt. Project read-back confirms #76 is PR review/S030, the linked pull request is #77 and default Status remains unused.
- Round-one Codex review found one P1 credential-ordering gap: checkout-controlled package installation preceded the first exact-main equality check while the protected job token was already scoped to the job. Remediation moved a fresh exact-main check to immediately after checkout, retained the pre-proof and pre-deploy checks, and added ordering regression coverage plus matching architecture/contract language.
