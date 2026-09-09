# Research: Ratify Foundation Contracts

## Schema identifier and lifecycle

**Decision**: Use JSON Schema Draft 2020-12. The schema artifact uses the Draft 2020-12 metaschema URI in its `$schema` keyword and `https://cueson.io/schema/v0.0.0/cueson.schema.json` as its `$id`. Cue JSON instances use the same canonical Cueson URI in their project-defined `$schema` member and use `schema_version` value `0.0.0`. The canonical URI is an identifier before domain activation. The unreleased schema may be refined; after a `v0.0.0` release it is immutable. Before v1, breaking changes require a documented minor release and patch releases remain non-breaking.

**Rationale**: This distinguishes schema dialect from project identity, follows the existing versioned URI pattern, preserves schema/software lockstep, and avoids making the future documentation domain a bootstrap dependency.

**Alternatives considered**: A repository URL was rejected because it would make the canonical identifier depend on source hosting. A `latest` URI was rejected because it is mutable. Treating unreleased v0.0.0 as immutable was rejected because the canonical schema has not yet been implemented or released.

## Initial format capability state

**Decision**: The completed v0.0.0 milestone recognizes `subrip` and `webvtt` with `format_support.status` set to `envelope_only`, `ingest_supported` and `render_supported` false, `restore_supported` true for valid source envelopes, and `ocr_required_for_semantic_output` false. Source-foundation issue #6 owns the generic `restore` command required to make that public capability claim truthful.

**Rationale**: The milestone establishes the lossless source envelope and public exact restoration before native subtitle codecs. Capability fields describe the official Cueson release, not arbitrary producer software. CLI tokens or extensions `srt` and `vtt` may later be explicit aliases but are not schema keys.

**Alternatives considered**: `reserved` was rejected for the completed milestone because source preservation and public restoration are now assigned foundation outcomes. `stable` and `beta` were rejected because native codecs are explicitly deferred. Keeping `restore_supported` true while shipping no restore command was rejected as an untruthful executable capability claim.

## OCR observation cardinality

**Decision**: Every cue contains an `ocr_observations` array with zero or more observations. Each non-empty observation is derived, independently provenanced, identifies its engine and source reference, and cannot replace native cue text or authoritative source assets. If a `derived` field is retained, it is constrained to true.

**Rationale**: A stable array shape avoids later scalar-to-array migration and allows multiple engines, models, languages, and processing passes to coexist.

**Alternatives considered**: Omitting the property for text formats was rejected because it makes the common cue shape less predictable. A nullable scalar or single object was rejected because it cannot represent independent observations safely.

## CLI command availability and failures

**Decision**: Help and command listings contain only commands implemented by that executable. The first CLI slice exposes help and `version`; the schema slice adds `schema`; the source-foundation slice adds generic `restore`. An unregistered command returns invocation exit code 2. A shipped command that cannot complete because a required codec is absent returns runtime exit code 1 and distinguishes that condition from an unknown format and strict loss rejection.

**Rationale**: This keeps help truthful while preserving explicit runtime capability diagnostics once relevant commands exist.

**Alternatives considered**: Registering the entire v1 command surface as placeholders was rejected because it advertises unavailable behavior. Treating a missing codec as invocation misuse was rejected because the input may be valid and the limitation belongs to the executable's runtime capability.

## Architecture boundaries

**Decision**: Keep command wiring, version identity, common model, schema, source restoration, codecs, conversion, and future OCR providers as separate internal responsibilities. Exact restoration reads the source envelope directly and does not call a codec. No package is a public Go API in v0.0.0.

**Rationale**: The boundaries align with the constitution, keep source truth independent from derived interpretation, and let future formats grow without coupling OCR or restoration to each parser.

**Alternatives considered**: Putting restoration on the codec interface was rejected because exact source bytes do not require parsing. Publishing Go packages at v0.0.0 was rejected because only the CLI and schema are approved public contracts.

## README badge state

**Decision**: Use static planned/unreleased badges while the CI workflow and release do not exist, linking them to issue #8 and milestone #1. Replace them with live badges only after those dynamic resources exist and are verified.

**Rationale**: A static badge renders immediately and describes present state without treating a permanent 404 or no-release response as propagation delay.

**Alternatives considered**: Leaving the current badges was rejected because the CI target does not exist. Removing all badges was rejected because the repository intentionally follows the ShruggieTech header pattern. Creating a placeholder CI workflow was rejected as out of scope for issue #8.

## GitHub bookkeeping

**Decision**: Record completed acceptance criteria and verification evidence on issue #1, but leave its README criterion, closure, and Done transition pending until the badge correction reaches the default branch. Set issue #3 to In progress with this slice identifier and keep it open until the eventual pull request merges. Amend issue #6 to own the public generic restore boundary. Do not claim a remote branch or pull request before either exists.

**Rationale**: This keeps GitHub as the delivery authority without falsely closing bootstrap work or crossing push, PR, or merge boundaries.

**Alternatives considered**: Closing issue #1 before its badge correction merged was rejected because default-branch README state would remain false. Closing issue #3 before merge was rejected because its ratified documents are not yet on the default branch.
