# Cue JSON Schema Baseline

**Status:** Ratified v0.0.0 implementation baseline

**Ratified:** 2026-09-09 through Spec Kit slice `001-ratify-foundation-contracts`

This document defines the schema decisions realized by implementation issue [#5](https://github.com/shruggietech/cueson/issues/5). The [canonical schema artifact](../internal/schema/cueson.schema.json) is embedded in the executable, and its structural and semantic validation foundation is available internally; source-integrity execution and exact restoration remain assigned to issue [#6](https://github.com/shruggietech/cueson/issues/6).

## Dialect, identity, and version

The schema artifact uses JSON Schema Draft 2020-12. Three similar-looking fields have distinct meanings:

- The schema artifact's `$schema` keyword identifies the Draft 2020-12 metaschema.
- The schema artifact's `$id` is `https://cueson.io/schema/v0.0.0/cueson.schema.json`.
- A Cue JSON instance uses the same canonical Cueson URI in its project-defined `$schema` member and uses `schema_version` value `0.0.0`.

The canonical Cueson URI is an identifier before `cueson.io` serves public schema files. Before domain activation, consumers resolve the exact schema through the repository, executable embedding, or release artifact. Cueson never emits a mutable `latest` alias into a document.

The unreleased development copy of v0.0.0 may be refined. Once v0.0.0 is released, the repository release copy, embedded copy, and release artifact are byte-identical and immutable.

Before v1.0.0, a documented breaking contract change requires a minor-version increase and patch releases remain non-breaking. Additive compatible changes may occur in a minor release. At and after v1.0.0, breaking changes require a major-version increase. Official software and schema versions remain equal; a third-party producer version is independent from the schema version it targets.

## Naming and root shape

Cueson-owned property names and enum values use lowercase `snake_case`. Standards-defined JSON Schema keywords such as `$schema`, `$id`, `$defs`, `oneOf`, and `additionalProperties` retain their standard spellings.

The v0.0.0 root contract contains these semantic areas:

```text
$schema
schema_version
format
format_support
producer
source
metadata
document
cues
format_data
diagnostics
stats
```

The exact JSON Schema defines required versus optional properties. The common cue model and source envelope are both first-class; neither is a substitute for the other.

## Initial format keys

The initial canonical schema and `format_data` keys are:

```text
subrip
webvtt
```

File extensions and future CLI tokens `srt` and `vtt` may be accepted as explicit aliases, but they normalize to the canonical keys and never appear as schema format values.

These keys identify format families; they do not claim stable codec support. Stable SRT and WebVTT support remains a v1.0.0 gate.

## Official format capability

`format_support` declares the official capabilities of the Cueson release associated with the document contract. It does not describe the capabilities of arbitrary third-party producer software.

At the current v0.0.0 schema-foundation milestone, both `subrip` and `webvtt` use:

```json
{
  "status": "envelope_only",
  "ingest_supported": false,
  "render_supported": false,
  "restore_supported": false,
  "ocr_required_for_semantic_output": false
}
```

`envelope_only` means the release can represent an already-valid source envelope intended for later exact restoration, but it does not claim native semantic ingest, model-driven rendering, or restoration. The generic public `restore` command belongs to source-foundation issue [#6](https://github.com/shruggietech/cueson/issues/6). Without that public command, `restore_supported` remains false.

Schema recognition, structural validity, native ingest, model-driven render, exact restoration, and OCR dependency are separate facts. Implementations and documentation must not infer one from another.

## Common cue contract

Every timed semantic unit maps to a cue when such a mapping is meaningful. A cue contains:

```text
id
ordinal
source_order
source_identifier
timing
payload
speakers
tokens
ocr_observations
placement
format_data
```

Normalized timing uses integer milliseconds in `start_milliseconds`, `end_milliseconds`, and `duration_milliseconds`. Native timing syntax remains in format-specific data.

The payload exposes `raw_text`, `plain_text`, and ordered logical `lines`. `raw_text` preserves decoded native textual content without destructive semantic normalization; it is not a substitute for original bytes. Speaker and token observations never destructively alter native payload content.

Each cue contains an `ocr_observations` array with cardinality zero or more. Text-native cues use an empty array when no observation exists. Every non-empty observation is independently identified and provenanced, includes OCR engine identity and a resolvable source reference, and may include engine version, model, language, lines, confidence, alternatives, regions, and processing options. If the schema retains `derived`, its only valid value is true.

OCR text never replaces native cue text, source images, or authoritative source assets.

## Multi-asset source envelope

`source` is required from v0.0.0 and contains `primary_asset_id` plus an ordered `assets` array. Multi-asset structure is foundational even when a format commonly uses one file.

Each asset records:

- a document-local identifier and primary or companion role;
- a safe original basename without path or URI components;
- media type when known;
- exact decoded byte length;
- lowercase SHA-256 identity;
- original bytes in standard base64;
- truthful timestamp values and provenance;
- text-encoding observations when applicable, which may be null for binary assets.

S003 semantic validation resolves `primary_asset_id` and enforces asset identity, role, and cardinality rules. Structural validation checks the declared base64, byte-count, and SHA-256 shapes. Decoding source bytes and comparing their actual length and SHA-256 remain assigned to source-foundation issue [#6](https://github.com/shruggietech/cueson/issues/6). Source bytes remain authoritative for exact restoration.

The document prohibits original filesystem paths, source directories, drive or mount data, working directories, hostnames, usernames, and other machine identifiers. A basename is non-empty, is neither `.` nor `..`, and rejects `/`, `\`, ASCII control characters, `<`, `>`, `:`, `"`, `|`, `?`, `*`, trailing spaces or periods, drive prefixes, URI prefixes, traversal components, and NTFS alternate-data-stream syntax. It also rejects the case-insensitive Windows device stems `CON`, `PRN`, `AUX`, `NUL`, `CLOCK$`, `CONIN$`, `CONOUT$`, `COM1` through `COM9`, and `LPT1` through `LPT9`, including those stems followed by an extension.

Within one `source.assets` bundle, `file_name` values must ultimately be unique by a portable collision key formed through Unicode canonical caseless matching: NFD normalization, default Unicode case folding, then NFD normalization again. S003 enforces each basename's structural safety but does not yet calculate this cross-asset collision key. Source-foundation issue [#6](https://github.com/shruggietech/cueson/issues/6) owns that enforcement before restoration opens any output so distinct source assets cannot collapse onto one destination on case-insensitive or normalization-insensitive filesystems.

Timestamp metadata distinguishes true creation or birth time, platform creation time, fallback observations, and unavailable values. Unix `ctime` is never mislabeled as creation time. Captured metadata remains in the document even when the destination cannot restore it.

## Validation layers

Schema implementation distinguishes:

1. JSON parsing.
2. Draft 2020-12 structural validation.
3. Cueson semantic validation, including cross-field references and capability consistency.
4. Source-envelope integrity validation, including decoded base64 length and SHA-256 comparison (deferred to issue #6).
5. Official software/schema version equality.

Schema validation does not infer codec availability. Native codec presence remains an executable capability concern described by [the architecture](architecture.md) and [CLI contract](cli.md).
