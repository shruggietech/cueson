# Quickstart: Prepare v1.1.0 public schema and site

## 1. Confirm governed state

Verify branch `codex/S030-publish-v1-1-public-schema-site`, `.specify/feature.json`, issue #76 as child/blocker of #67, closed #66 evidence, milestone v1.1.0 and Project Slice S030. Confirm main baseline is `7ff45c1d8cd8df377e1fb568b9785286b649fd7c`, the working tree is clean and production authority is absent.

## 2. Run tests first

Update focused expectations for current release identity, exact downloads, v1.1.0 release route, three immutable schemas, absent latest alias, complete metadata, exact-main workflow policy and local-authority production proof. Run the focused tests and record expected failures against candidate-era implementation.

```text
pnpm --dir site test:unit
go test ./scripts/docs-verify
go test ./scripts/site-policy
```

## 3. Implement the maintained authority and artifact

Add the validated current release record, exact downloads, v1.1.0 release page and immutable schema source. Update generated content and deployment metadata, landing content and maintained current-state documents. Preserve earlier routes, schemas, release pages and frozen Spec Kit evidence.

## 4. Harden later production proof

Compare public deployment metadata with the locally generated reviewed record, hash the public content manifest and probe local declared inventories. Require exact equality between deployment revision and freshly fetched `origin/main`. Do not invoke deployment.

## 5. Run focused and full verification

```text
pnpm --dir site test
go test ./...
go test ./scripts/...
go vet ./...
```

Also run standalone module tests, formatting, static analysis, vulnerability scanning, race/native/platform checks and repository text/UTF-8/mojibake/whitespace gates defined by current CI. Confirm generation leaves no diff and the deployable artifact declares exactly 22 HTML plus three schema routes and seven downloads.

## 6. Run Spec Kit gates

Run blocking analysis after tasks are generated. After implementation and local verification, run independent convergence across the specification, plan, tasks, contracts, issue acceptance and changed artifacts. Resolve every material finding and record evidence in `verification.md`.

## 7. Publish and converge the official pull request

Commit conventionally, push the authorized branch and publish a formatted/read-back PR that includes `Closes #76` and `Refs #67`. Move #76 to PR review, clear default Status and wait for all exact-head CI/security/Codex results. Address every comment and resolve threads after verification. If round one has findings, request exactly one second review with `@codex review`; never request a third automatically.

## 8. Stop for the human merge ritual

Report exact PR head, green checks, terminal review state, resolved findings and remaining boundary. Do not merge or deploy. After human merge, #67 will require the exact resulting main SHA, separate production authorization, exact-main rebuild/deployment and complete live verification before the epic and milestone can close.
