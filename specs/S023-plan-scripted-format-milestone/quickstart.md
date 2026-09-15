# S023 validation guide

## Local review

Read [spec.md](spec.md), [plan.md](plan.md), [delivery contract](contracts/delivery.md), [version compatibility](contracts/version-compatibility.md), the maintained roadmap and ASS/SSA page. Confirm future capability labels and current v1.0.0 behavior remain distinct.

## Existing checks

Run commands through a verified hidden Windows launcher with redirected non-interactive I/O; direct git/gh are allowed. Verification remains foreground.

```text
go run ./scripts/github-format/main.go .
go -C scripts/github-format test -count=1 ./...
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
go test -count=1 ./...
corepack pnpm --dir site test
git diff --check
```

Expected: no formatting/encoding/links/example violations, unchanged old runtime/schema tests, complete generated site/build/browser/artifact/dry-run success and no whitespace errors. Native hosted CI supplies Linux/macOS, race, security and packaged proof beyond local Windows evidence. No production verifier/deployment or release command is needed for this slice.

## Native planning read-back

Inspect the v1.1.0 milestone and epic, all sixteen children and their parent/blocker APIs. Compare to [data-model.md](data-model.md) and actual issue-map snapshot. Verify forty-five Project items (twenty-eight existing plus seventeen new), unique issue numbers, all existing Done items unchanged, S023 children PR review, future children Backlog, child Slice text and unused Status. Inspect current issues again for unassessed arrivals.

## Final review handoff

Read back the official PR body, exact head check state, comments/reviews/reactions and unresolved threads. All required checks and reviews must be satisfied before asking for the human's final review and merge. Do not perform merge or post-merge housekeeping here.
