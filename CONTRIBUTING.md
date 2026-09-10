# Contributing to Cueson

Cueson is an unreleased v0.0.0 foundation candidate. There is no supported installation package or public binary release, and native SubRip/WebVTT ingest, model-driven rendering, and conversion are not implemented. Development and verification use the Go 1.25 source tree.

Start with an issue so the intended outcome, dependencies, and verification can be agreed before implementation begins. Read the [architecture](docs/architecture.md), [CLI contract](docs/cli.md), and [schema baseline](docs/schema.md) before changing their surfaces. Format work must also begin with the dedicated [SubRip](docs/formats/srt.md) or [WebVTT](docs/formats/webvtt.md) boundary, which distinguishes current envelope-only behavior from planned v1 work.

## Development workflow

1. Select or create an atomic GitHub issue with explicit acceptance criteria and verification.
2. Group compatible issues into a coherent work slice when they share one implementation and review story.
3. Use the repository's Spec Kit workflow for non-trivial product work.
4. Update tests, documentation, and `CHANGELOG.md` alongside behavior changes.
5. Run the repository verification commands in the foreground.
6. Open a pull request containing a complete closing reference such as `Closes #123`.

Repository-authored Markdown uses one source line per paragraph and per list item. Before publishing GitHub text, pass it through the repository formatter and publish from a UTF-8 file when practical.

Run the repository text and maintained-document checks before submitting documentation changes:

```text
go run ./scripts/github-format/main.go .
go -C scripts/github-format test -count=1 ./...
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
```

The default branch accepts squash pull requests through active repository protection. Required CI and CodeQL checks must pass on the current head, and review conversations must be resolved. The pull-request policy reports issue-link and Codex-review state, but those two status contexts are not currently protection-required. The organization-administrator bypass is for recovery from broken controls rather than ordinary delivery.

An AI agent may prepare and verify a pull request, but the final merge remains a human decision unless the operator grants one explicit, pull-request-specific override.

Candidate packaging is non-publishing. Tag creation, GitHub Release publication, immutable release-schema copying, milestone closure, and production changes follow the separately authorized [release process](docs/release-process.md).
