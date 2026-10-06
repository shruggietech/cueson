# S034 candidate verification

**Preparation baseline:** `40473babe6fd649515f547115872b6316b641263`, the actual merged S033 main revision. **Candidate branch:** `codex/S034-prepare-v1-2-release-candidate`. **Date:** 2026-10-06. This baseline is not a publication target. The final PR body records the actual candidate head, hosted run URLs, retained artifacts and terminal reviews after they exist.

## Specification and integration

The installed Spec Kit workflow completed specify, clarify, checklist, plan, tasks, blocking read-only analysis and implementation. Analysis found complete coverage of fourteen requirements by twenty-three tasks, with no unresolved findings. Independent schema, documentation and release-proof reviews identified scoped proof and prose gaps. The corrected strict historical conversion, isolated selector probes, media-boundary checks and diagnostic privacy guards passed independent remediation review. The independent reviewer assessed and documented all eight release-contract requirements-quality criteria.

Current software and schema identify exact 1.2.0. The embedded and immutable candidate schema are byte-identical: 191170 bytes, SHA-256 `f2661a3d52effbab4a82a4d47197b5c7fae58496dc30a397ea3f2f668358b654`. Identity-only reversal recovers the reviewed S033 contract; consumer validation behavior is unchanged. Historical released schemas, fixtures, publication contracts, legal files and public site/download inventory remain unchanged. Former development identity is refused by schema and typed consumer validation.

## Completed local verification

- Product: uncached `go test -count=1 ./...`, vet, Staticcheck v0.7.0 and Govulncheck v1.8.0 passed. Govulncheck found no affected reachable product symbols; its informational imported-package finding is not a clean bill for every upstream package.
- Portability: six pure-Go, trimpath builds passed for Linux, Windows and macOS on amd64 and arm64. These builds do not claim foreign-host or arm64 execution.
- Robustness: all sixteen CI fuzz targets plus the two consumer identifier/media-arithmetic targets passed 1000 generated iterations each, with one worker per target.
- Nested modules: github-format, pr-policy, brand-verify, corpus-verify, docs-verify and release-verify tests, vet, Staticcheck and Govulncheck passed. The corpus module likewise reported an informational uncalled imported-package finding without affected reachable symbols.
- Repository: Go formatting, workflow lint (Actionlint v1.7.12), text formatting, UTF-8/no-BOM checks, diff whitespace and preserved-history checks passed. Brand verification accepted 265 entries and six references. Documentation verification accepted 27 documents, 145 local links, ten examples and 34 format rows.
- Published consumers: actual Windows integration authenticated downloaded published 1.0.0 and 1.1.0 executable bytes, legal/schema/build/version identity and positive historical acceptance before current-document refusal. Final proof counters are enforced in tests and the candidate evidence schema.
- Evidence contract: compiled with the existing JSON Schema validator; synthetic structural and three-host evidence accepted, while 36 malformed host, identity, size, count, consumer and privacy records were refused. This validates the contract, not the future hosted artifacts.
- Site: lint, generation/check, 88 unit tests, static build and artifact verification passed. Generated inventory retained twenty-two HTML routes and three published schemas. Browser/accessibility checks passed 81 tests with six established skips using an isolated port; the temporary configuration was removed. Wrangler deployment dry-run passed without deployment.

Verification tooling ran with redirected non-interactive I/O and hidden child-process creation on Windows. Tests were awaited to terminal completion. Tool-selected Go 1.26.8 supported analysis tools; the advertised compatibility floor remains Go 1.25 and receives hosted checks.

## Hosted gates and publication limits

Tasks T022 and T023 are runtime gates recorded on the final PR body. They require all final-head hosted checks, six packaged archives/SBOMs, the checksum manifest, same-bundle native Linux/Windows/macOS amd64 proof, every review finding addressed and at most one finding-driven second Codex request. Local source-built Windows integration does not substitute for those packaged gates.

The prepared public notes, intended thirteen-file inventory and publication decision procedure are reviewable in this slice. Actual public release approval requires the authorized merge's actual main revision, fresh successful main checks, retained accepted artifact IDs/expiry, exact asset lengths/digests and notes binding. No predicted squash revision or PR artifact is accepted as final main evidence. No tag, GitHub Release, asset upload, schema hosting or production deployment occurred during preparation; public v1.1.0 remains available.

Candidate issue #84 belongs to S034. Later public publication #85 depends on #84; production hosting #86 depends on #85. All three belong once to the Delivery Project and v1.2.0 milestone, using Stage and Slice while leaving default Status unused. No other active issues were found at kickoff or the pre-publication assessment.
