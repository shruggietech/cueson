# S035 Validation Guide

Use the repository's verified hidden Windows process runner or equivalent headless environment; retain foreground verification.

```text
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
go run ./scripts/github-format/main.go
cd site
pnpm lint
pnpm generate
pnpm generate:check
pnpm test:unit
pnpm build
pnpm test:browser
pnpm verify:artifact
pnpm deploy:dry-run
```

Use a temporary isolated browser configuration on port 4174 if 4173 is already serving an older artifact. Expect 28 routes, four schemas and seven v1.2 downloads; all four released schema hashes must match. Negative tests reject altered or missing new schema, consumer route, release/download metadata and stale proof. Retain existing credential/revision/Worker checks.

Full hosted CI, CodeQL, Site, Release proof and Codex review apply to the exact PR head. Production entry/completion is in [production-handoff.md](contracts/production-handoff.md). Production deployment and live verification continue within S035 immediately after the final owner merge under the recorded steering authorization.
