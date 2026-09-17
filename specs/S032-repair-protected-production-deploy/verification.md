# Verification: Repair protected production deployment

## Baseline

- Branch `codex/S032-repair-production-deploy` started clean at exact main `926b062f7417580b8fac83d906b0fd34dfda3223`.
- Production was healthy at v1.1.0 with 25 routes, seven downloads, and three immutable schemas.
- Issue #78 was the only open repository issue and remained the operational production-validation outcome.
- The temporary GitHub production secret and remote recovery branch were absent. The local recovery branch remained retained because prior safe deletion checks refused removal.

## Requirements and analysis

- Issue #78 was repaired from collapsed Markdown and read back with all six governed sections and eleven post-merge acceptance criteria.
- Child issue #79 owns the independently closeable workflow, test, and credential-contract result. GitHub records #79 as a native sub-issue of #78, and #79 appears exactly once in `cueson Delivery` with Stage Specced, Slice S032, and empty default Status.
- The requirements checklist passed. No clarification marker remained.
- Blocking analysis found one material lifecycle defect in the initial authority: closing #78 from the source PR would claim production proof before the merged workflow could run. The #79/#78 split resolves it.
- Cross-artifact review found no remaining critical, high, or material inconsistency across the specification, plan, tasks, credential research, deployment contract, constitution, architecture, S030 contract, S031 recovery evidence, or issue acceptance.

## Red evidence

- `node --test tests/workflow-command.test.mjs tests/generator.test.mjs` ran 12 tests with 10 passing and two intended failures after the Windows test harness was corrected to invoke the pnpm JavaScript CLI directly without a shell.
- Static workflow policy failed because the workflow lacked main ref/SHA binding and still used job-scoped credentials plus literal package-script separators.
- The production verifier command-boundary test failed because its parser silently discarded the standalone separator and proceeded to network access instead of rejecting malformed arguments.
- Both correct named-argument package-script probes passed their parsers. The Cloudflare probe reached the expected missing-credential validation, and the public verifier reached its downstream fetch boundary.

## Focused verification

- `node --test tests/workflow-command.test.mjs tests/generator.test.mjs`: PASS, 12 tests.
- Correct Cloudflare named arguments reached the explicit missing-credential gate; a standalone separator failed as `invalid argument: --`.
- Correct public-verifier named arguments passed parsing and reached the downstream fetch boundary; a standalone separator failed as `invalid argument: --`.
- Static policy proved manual-only invocation, full lowercase revision validation, `refs/heads/main` execution, `GITHUB_SHA == REVISION`, three fresh `origin/main` equalities, no ancestor fallback, direct verifier arguments, the public deployment-record environment URL, no job-scoped token, and exactly three step-scoped token references.
- The verifier call surface uses zone identity, DNS records, Worker scripts, Worker Custom Domains, zone ruleset inventory, and redirect ruleset detail. The documented Account Workers Scripts Write plus Zone Zone Read, DNS Read, and Transform Rules Read contract covers deployment and those calls without unrelated administrative permissions.

## Complete verification

- Complete site pipeline: PASS. ESLint, deterministic generation/check, 87 unit tests, Next production build, 81 Playwright tests with six intentional project-scoped skips, exact artifact verification for 25 routes and three immutable schemas, and Wrangler dry run all succeeded.
- Root Go gates: PASS. `go test -count=1 ./...`, `go vet ./...`, and `go test -race -count=1 ./...` succeeded.
- All six standalone Go modules passed tests. Documentation verification checked 25 documents, 131 local links, 10 registered examples, and 34 format rows. Brand verification checked 265 entries and six references.
- Repository text, UTF-8/no-BOM, mojibake, Markdown layout, whitespace, and `git diff --check`: PASS through `scripts/github-format` and Git.
- GitHub workflow lint: PASS with actionlint v1.7.12.
- Root and all standalone modules passed vet and Staticcheck v0.7.0.
- Root and all standalone modules passed govulncheck v1.8.0 with zero called vulnerabilities. The root and corpus verifier reported one informational vulnerability in imported packages that no code path calls.
- No schema, release asset, product runtime, public content, Cloudflare state, GitHub secret, tag, or release changed during pull-request implementation.

## Convergence

- All 12 functional requirements, five success criteria, three stories, ten acceptance scenarios, nine child-issue acceptance criteria, the protected deployment contract, and seven constitutional principles have reviewable implementation or an explicit post-merge validation task.
- Revision identity is covered at workflow execution, checkout, three fresh-main gates, generated metadata, the GitHub deployment record contract, and the public deployment verifier.
- Credential coverage matches the full existing Cloudflare call surface. No unrelated permission or job-level secret exposure remains.
- One low robustness finding was corrected: the subprocess tests now fail on launcher errors and require the correct production invocation to reach its downstream fetch boundary, preventing a missing executable from masquerading as parser success.
- Final convergence found no unresolved critical, high, medium, or low actionable issue. Production mutation remains deferred until the human merge, so #78 stays open.
- Prepublication Project read-back confirmed #78 and #79 are both open, In progress, Slice S032, and have empty default Status. GitHub records #79 as a native child of #78.

## Hosted checks and review

- The authorized branch was pushed and official pull request #80 was published with `Closes #79` and `Refs #78`.
- The PR body passed github-format before publication. Immediate read-back confirmed the intended headings, paragraphs, complete command fence, checklist, closing reference, production exclusion, and post-merge continuation.
- Project read-back confirms #78 and #79 are both PR review/S032, linked to #80, and retain empty default Status.
- Round-one Codex review reported one P2 Windows portability finding: standard Corepack `pnpm.cmd` shims may target `pnpm.js`, and a leading separator after `%~dp0` must remain relative to the shim directory. The resolver now accepts `.js`, `.cjs`, and `.mjs`, strips only leading separators from the suffix, resolves with Windows path semantics, and has an explicit standard-Corepack-shim regression test.
- Post-fix focused verification passed 13/13 tests. The complete site pipeline passed again with 88 unit tests, 81 browser tests plus six intentional project-scoped skips, exact 25-route/three-schema artifact verification, and Wrangler dry run.
- Final exact-head checks and second-round review are pending.
