# S035 Verification Record

## Scope and authority

Prepared on `codex/S035-publish-v1-2-public-schema-site` from actual S034 main `6a1fd7a541767c0b19dcc10022fafc0e46903cff`. Spec Kit specify → clarify → checklist → plan → tasks → blocking read-only analyze → implement ran in installed skills mode with installed PowerShell setup/prerequisite scripts. No extension hooks are installed. Analysis has twelve requirements/eighteen tasks with complete coverage and no unresolved findings. Built-in and independent requirements checklists each pass 8/8.

The kickoff authorizes push/official PR and at most two Codex rounds, stopping for final human review/merge. Subsequent owner steering authorizes the same S035 run to continue immediately after confirmed merge through exact-main production deployment and independent live verification, without another kickoff/approval. [Recorded authority on #86](https://github.com/shruggietech/cueson/issues/86#issuecomment-6028720348) was formatted and read back exactly. Preparation #88 closes through the PR; #86 and milestone remain open until T016–T018 actually pass.

## Published release binding

Fresh GitHub read-back confirms final public v1.2.0 release ID405257062 and unsigned annotated tag object `386087872854d77b7f130d3b24ec517f1a5d0aad`, peeling to `6a1fd7a541767c0b19dcc10022fafc0e46903cff`. All thirteen exact public asset records match names, positive sizes and SHA-256 from completed [#85 public proof](https://github.com/shruggietech/cueson/issues/85#issuecomment-6028397732). Independently downloaded public files retained from that publication still match every digest. All seven new primary URLs bind to those authentic public assets. No assets were rebuilt or replaced.

| Immutable version | Bytes | SHA-256 |
|---|---:|---|
| 0.0.0 | 22263 | `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975` |
| 1.0.0 | 61445 | `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541` |
| 1.1.0 | 185641 | `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7` |
| 1.2.0 | 191170 | `f2661a3d52effbab4a82a4d47197b5c7fae58496dc30a397ea3f2f668358b654` |

All four source/generated identities match. `git diff --exit-code origin/main -- internal schema testdata .goreleaser.yaml .github/workflows docs/releases/v1.1.0.md docs/releases/v1.0.0.md docs/releases/v0.0.0.md` passes. Runtime, schemas, historical release pages and workflows are unchanged.

## Local verification

All commands ran in the foreground through a verified `CREATE_NO_WINDOW` runner with redirected non-interactive I/O. Raw logs are ignored under `dist/S035`.

| Check | Result |
|---|---|
| `go test -count=1 ./...` and `go vet ./...` | PASS |
| `go -C scripts/docs-verify test -count=1 ./...` | PASS, including updated publication-marker rejection regressions |
| `go -C scripts/docs-verify vet ./...` and pinned Staticcheck v0.7.0 | PASS |
| Pinned Govulncheck v1.8.0 for docs verifier | PASS, no vulnerabilities found |
| `go run ./scripts/github-format/main.go .` and formatter module tests | PASS |
| `go -C scripts/docs-verify run . -repo ../..` | PASS:27 documents,151 local links,10 registered examples,34 format rows |
| Brand verifier | PASS:265 entries,6 references |
| Site lint/generate/generate:check | PASS |
| Site unit/negative checks | PASS:88 tests,0 skips |
| Next static build | PASS |
| Playwright with isolated port4174 preview | PASS:90 passed,6 established width-specific skips; new pages included at all three widths |
| Static artifact verifier | PASS:28 declared routes,24 HTML,four schemas,seven primary downloads |
| Wrangler deployment dry run | PASS; no production mutation |
| Independent implementation review | Clean; all new inventories, negative mutations, history and platform claims inspected |
| Git whitespace and UTF-8/no-BOM/mojibake checks | PASS |

Test-first proof: revised generator expectations failed against old authority (25 instead of28 routes and missing fourth schema/new release/download identity); updated documentation-marker regressions failed against old candidate-only rules. Both pass after authority/rule changes. The deliberate download mismatch now changes declared version to1.3.0 while retaining1.2.0 filenames, preserving a real rejection. Remote metadata mutations cover missing/altered new schema, missing consumer route, stale revision and historical primary download. Existing credential/revision/infrastructure and accessibility tests remain intact.

The local Playwright shell launcher could not start its `cmd.exe` subprocess. Browser checks therefore used a directly launched hidden Node preview and temporary equivalent configuration on isolated4174, avoiding stale4173 reuse. The temporary configuration and owned preview were removed after successful checks. Existing preview processes were retained. Pre-existing ignored release-proof Go helpers were isolated in an ignored scratch module so root `go test ./...` cannot discover them; no source/evidence was deleted.

## Runtime PR acceptance

First Codex review on PR #89 found two P2 issues: a scope-update bug substituted ASCII letters throughout the validation guide, and conversion/working-spec current-state prose still described unpublished v1.2.0. UTF-8 checks could not detect the ASCII corruption. The guide was restored, its actual workflow fence and handoff target now have a focused contract test, both stale current-state pages were reconciled, and documentation/generated-page regressions explicitly reject those old claims. All affected checks were rerun before updating the head and requesting the one authorized second review. Frozen history and runtime behavior remain unchanged.

The second and final Codex review found stale v1.1.0 current-state claims in the media guide and three native format contracts. All four were reconciled to published v1.2.0, preserving historical introduction evidence and the pending #86 live-acceptance boundary. Documentation negative regressions and generated-page checks now cover encode identity, native shape identity, all four schema versions and the current production acceptance owner. The finding is handled on a descendant remediation head without a third automated review.

T014/T015 exact pushed head, PR identity, hosted checks and external review outcomes are recorded in the final formatted/read-back PR body after they occur. They are not claimed complete by this pre-push record. No second Codex round is requested unless the automatic first round reports findings; no third round is authorized.

The single second-round request originally omitted its exact backticked head required by repository policy. Correcting that existing receipt after the completed review changed its evidence timestamp, so the trusted default-branch policy treats the result as stale despite all findings being fixed and resolved. No third request or operator waiver was issued. AGENTS.md now requires the exact head binding before publication and preserves the posted receipt. The remaining policy gate requires an explicit operator decision; successful product/site checks do not override it.

## Production before-state and remaining execution

Independent live HTTPS before-state confirms `https://cueson.io/deployment.json` returned200 with revision `ec8cab355f833efb5a632dcdc980881e294ee2f8`, release1.1.0 and three schemas. The v1.2.0 schema, release page and consumer-speaker route each returned404. This distinguishes already verified GitHub publication from website activation.

T016–T018 remain mandatory after the final owner merge: scoped housekeeping, fresh actual-main protected proof/deployment, infrastructure read-back and independent public verification, including explicit schema headers and separate IPv4/IPv6 HTTPS probes. Use [production-handoff.md](contracts/production-handoff.md). S035 is not complete until live acceptance passes; this preparation record does not assert deployment or defer it to a new slice.
