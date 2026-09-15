# S022 Validation Guide

## Prerequisites

Use the existing Node.js 24, pinned pnpm, compatible Go, and installed Playwright Chromium baseline. Windows requires verified headless launchers with redirected non-interactive I/O; direct Git and gh remain allowed.

## Deterministic DNS Regression

From `site/`, run:

```text
node --test tests/production-dns.test.mjs
```

Expect both resolver matrices to cover IPv4-only, IPv6-only, dual-stack, partial failure, no addresses, malformed/wrong-family/alias-only evidence, deterministic inventories, and complete failure diagnostics. See [the DNS contract](contracts/dns-verification.md).

## Full Site and Repository Checks

From `site/`, run:

```text
corepack pnpm test
```

This includes lint, generation/drift, unit tests, static build, browser/accessibility checks, artifact verification, and Wrangler dry run. From repository root, run:

```text
go run ./scripts/github-format/main.go .
go -C scripts/docs-verify run . -repo ../..
go test -count=1 ./...
git diff --check
```

## Existing Production Acceptance

From `site/`, run:

```text
corepack pnpm verify:production -- --expected-commit 46838fd5cc888b299a09b89da676db0005b00f16
```

Expect both hostnames, trusted TLS, www path/query redirect, declared routes/downloads, deployment revision, immutable schema hashes, and absent latest alias to pass. Do not bypass network identity or deploy production to obtain acceptance.

The local S022 session found that this machine returns ENOTFOUND for dns.google. Local public verification therefore uses a temporary process-only hostname bootstrap to Google's independently confirmed 8.8.8.8/8.8.4.4 service addresses while retaining the original HTTPS hostname, normal certificate validation, and every verifier probe. This environment accommodation is not committed and changes no system or production DNS; environments resolving dns.google normally run the command above unchanged.

## Delivery

Commit and push S022 under explicit kickoff authorization. Format PR Markdown through the repository publication formatter, read it back, and include `Closes #49`. Wait for hosted checks and external reviews, handle every finding, use at most one second Codex request, and return the satisfied PR for human review and merge.
