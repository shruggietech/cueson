# Documentation Contract

## Canonical required set

S010 requires these twelve enduring documentation paths:

1. `README.md`
2. `CONTRIBUTING.md`
3. `SECURITY.md`
4. `AGENTS.md`
5. `CHANGELOG.md`
6. `docs/architecture.md`
7. `docs/schema.md`
8. `docs/cli.md`
9. `docs/formats/srt.md`
10. `docs/formats/webvtt.md`
11. `docs/project-management.md`
12. `docs/release-process.md`

The maintained audit scope also includes `.specify/memory/constitution.md` and all other Markdown documents beneath `docs/`, including the working project specification, repository-control evidence, and non-publishing release-verification guide.

## Authority boundaries

- `README.md` is the product overview, current status, entry-point examples, and documentation index.
- `docs/architecture.md` is the architecture of record.
- `docs/schema.md` is the release-specific Cue JSON contract and compatibility summary.
- `docs/cli.md` is the complete currently shipped CLI contract.
- The two format pages describe current envelope-only capability and planned v1 grammar/fidelity coverage without claiming codec implementation.
- `docs/release-verification.md` owns the runnable non-publishing artifact proof.
- `docs/release-process.md` owns the broader protected release lifecycle.
- `docs/project-management.md` owns the current GitHub-native operating contract.
- The working project specification is the broader roadmap and completion-gate inventory where it does not conflict with ratified authorities.

## Claim vocabulary

- `implemented`: behavior exists and is covered by repository verification.
- `envelope_only`: the schema and source envelope support validation and exact generic restoration, but no native codec ingest or model-driven render is claimed.
- `planned`: future behavior with no present capability claim.
- `release-gated`: implemented preparation whose public effect requires a separately authorized release action.
- `post-release`: work intentionally deferred until after a successful public release.

## Publication boundary

S010 may document and verify the candidate state. It does not create a tag, release, immutable release-schema copy, public schema, milestone closure, or production-domain state.
