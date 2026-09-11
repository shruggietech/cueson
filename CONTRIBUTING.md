# Contributing to Cueson

Cueson v1.0.0 is published with stable SubRip, WebVTT, conversion, validation, inspection, completion, schema, and restoration workflows. The immutable v0.0.0 envelope-only foundation remains available as historical evidence. Development and verification require Go 1.25.0 or newer.

Start with an issue so the intended outcome, dependencies, and verification can be agreed before implementation begins. Read the [architecture](docs/architecture.md), [CLI contract](docs/cli.md), [schema contract](docs/schema.md), and [compatibility contract](docs/compatibility.md) before changing public surfaces. Format or conversion work must also begin with the dedicated [SubRip](docs/formats/srt.md), [WebVTT](docs/formats/webvtt.md), and [conversion](docs/conversion.md) contracts. Read the [brand guide](docs/brand.md) before referencing, updating, or distributing Cueson identity assets.

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
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
```

The documentation verifier requires registered README workflows, checks maintained contract markers and stale claims, and links every machine-readable conformance row to its format guide. Root-module CLI tests execute the registered workflows with governed inputs, temporary destinations, and captured streams; do not point documentation tests at writable fixture destinations.

The official brand archive and extracted kit are immutable byte-protected inputs. Do not format, optimize, repair, rename, or edit imported files in place. A revised official kit requires a separately specified versioned acquisition and a regenerated Cueson import manifest.

The default branch accepts squash pull requests through active repository protection. Required CI and CodeQL checks must pass on the current head, and review conversations must be resolved. The pull-request policy reports issue-link and Codex-review state, but those two status contexts are not currently protection-required. The organization-administrator bypass is for recovery from broken controls rather than ordinary delivery.

An AI agent may prepare and verify a pull request, but the final merge remains a human decision unless the operator grants one explicit, pull-request-specific override.

Candidate packaging is non-publishing. Immutable schema admission occurs through reviewed candidate work; tag creation, GitHub Release and asset publication, milestone closure, public schema hosting, and production changes follow the separately authorized [release process](docs/release-process.md).
