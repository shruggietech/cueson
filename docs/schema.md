# Cue JSON Schema Baseline

**Status:** v1.0.0 stable schema released and independently verified

**Ratified:** 2026-09-09 through Spec Kit slice `001-ratify-foundation-contracts`

This document defines the stable v1.0.0 schema, including the foundation realized by issues [#5](https://github.com/shruggietech/cueson/issues/5) and [#6](https://github.com/shruggietech/cueson/issues/6), native SubRip capability added by issues [#30](https://github.com/shruggietech/cueson/issues/30) and [#31](https://github.com/shruggietech/cueson/issues/31), native WebVTT capability added by issue [#32](https://github.com/shruggietech/cueson/issues/32), conversion behavior added by issue [#33](https://github.com/shruggietech/cueson/issues/33), and contract hardening completed by issues [#35](https://github.com/shruggietech/cueson/issues/35), [#36](https://github.com/shruggietech/cueson/issues/36), and [#41](https://github.com/shruggietech/cueson/issues/41). The [canonical schema artifact](../internal/schema/cueson.schema.json) is embedded in the executable and includes machine-readable annotations for schema-aware consumers.

## Dialect, identity, and version

The schema artifact uses JSON Schema Draft 2020-12. Three similar-looking fields have distinct meanings:

- The schema artifact's `$schema` keyword identifies the Draft 2020-12 metaschema.
- The stable schema artifact's `$id` is `https://cueson.io/schema/v1.0.0/cueson.schema.json`.
- A current-source Cue JSON instance uses the same canonical Cueson URI in `$schema` and uses `schema_version` value `1.0.0`.

The canonical Cueson URI remains an identifier until `cueson.io` separately serves public schema files. Consumers resolve 1.0.0 through the canonical repository schema, the [immutable v1 schema](../schema/releases/v1.0.0/cueson.schema.json), the schema packaged in the [v1.0.0 release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0), or executable embedding. Cueson never emits a mutable `latest` alias.

Users can discover the embedded contract version with `cueson schema --version` and emit the exact embedded schema with `cueson schema` or `cueson schema --output PATH`. A Cue JSON instance identifies its target contract through `$schema` and `schema_version`; consumers must evaluate both against an exact supported version rather than infer compatibility from the producer software version.

The v1.0.0 schema is released and immutable at [`schema/releases/v1.0.0/cueson.schema.json`](https://github.com/shruggietech/cueson/blob/v1.0.0/schema/releases/v1.0.0/cueson.schema.json), with SHA-256 `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541`. The historical v0.0.0 schema remains immutable with SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`. Public `cueson.io` schema hosting remains separately governed by the [release process](release-process.md).

Starting with v1.0.0, breaking changes require a major-version increase. Additive compatible changes may occur in a minor release, and compatible corrections may occur in a patch release. Official software and schema versions remain equal; a third-party producer version is independent from the schema version it targets.

S018 froze the implemented CLI and Cue JSON behavior, S019 promoted that reviewed contract to identity 1.0.0 with a byte-identical immutable repository copy, and S020 published and independently verified the exact tagged schema and release archives. Publication did not host the schema at `cueson.io` or add a production alias.

## Machine-readable annotations

The canonical schema gives every Cueson-owned root property and every consumer-facing property reachable through `$defs` a specific `description`. Public object, union, and constrained-scalar definitions have titles suitable for generated reference documentation. Examples cover every root semantic area, enumeration, format-native structure, and non-obvious constrained value such as identifiers, versions, safe basenames, digests, base64 payloads, timestamps, timings, and diagnostics.

`title`, `description`, and `examples` are non-normative JSON Schema annotations. They explain the contract but never broaden an `enum`, relax a pattern, satisfy a missing required property, override a conditional branch, replace model semantics, or weaken source-integrity checks. Validation keywords, the semantic model, and source-envelope integrity remain authoritative when an annotation is incomplete or a consumer does not process annotations.

Examples attached to a property or definition are valid for that exact schema fragment. Complete root examples additionally pass structural validation, semantic validation, and source-reference checks. Examples contain portable basenames and public identifiers only; they contain no original filesystem path, username, hostname, drive, mount, or other local machine identity. The examples under the safe-basename exclusion enum show values rejected by the enclosing basename contract, while examples on `safe_basename` itself show valid portable values.

Schema-aware tooling should present a property's own description together with the title and description of any referenced definition. It must not treat example values as defaults or assume an example exhausts allowed values. In particular, a producer `version` example identifies producer software and does not assert the Cue JSON schema version.

## Naming and root shape

Cueson-owned property names and enum values use lowercase `snake_case`. Standards-defined JSON Schema keywords such as `$schema`, `$id`, `$defs`, `oneOf`, and `additionalProperties` retain their standard spellings.

The stable v1 root contract contains these semantic areas:

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

The canonical schema and `format_data` keys are:

```text
subrip
webvtt
```

File extensions and CLI tokens `srt` and `vtt` are aliases that normalize to canonical keys and never appear as schema format values.

These keys and their format-native shapes identify format families. Capability fields separately declare whether the matching executable can ingest, render, or restore them. SubRip and WebVTT have completed the stable v1.0.0 release gate.

## Official format capability

`format_support` declares the official capabilities of the Cueson release associated with the document contract. It does not describe the capabilities of arbitrary third-party producer software.

The stable v1 SubRip contract uses:

```json
{
  "status": "stable",
  "ingest_supported": true,
  "render_supported": true,
  "restore_supported": true,
  "ocr_required_for_semantic_output": false
}
```

The stable v1 WebVTT contract uses the same capability values with status `stable`. Stable describes the reviewed and published v1.0.0 behavior and compatibility promise.

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

WebVTT document data preserves the raw signature line, optional description, ordered metadata lines, and non-cue `NOTE`, `STYLE`, `REGION`, or unrecognized blocks. Cue and non-cue `source_order` values form one unique contiguous sequence, so adjacent blocks and overlapping cues never depend on inferred ordering. Each block retains both its raw LF-joined body and physical lines; REGION blocks additionally retain ordered setting occurrences and the effective recognized setting map.

WebVTT cue data preserves the optional native identifier, raw timing line, raw settings text, ordered setting occurrences, effective recognized settings, and raw payload lines. Common payload lines must equal the native raw payload lines. Voice markup yields native speaker observations, valid inline timestamps yield timed text tokens, and neither derived view replaces the retained native payload.

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

Semantic validation resolves `primary_asset_id` and enforces asset identity, role, and cardinality rules. Structural validation checks the declared base64, byte-count, and SHA-256 shapes. Source preparation additionally requires canonical standard base64 and compares its actual decoded length and SHA-256 before any output is opened. Source bytes remain authoritative for exact restoration.

The document prohibits original filesystem paths, source directories, drive or mount data, working directories, hostnames, usernames, and other machine identifiers. A basename is non-empty, is neither `.` nor `..`, and rejects `/`, `\`, ASCII control characters, `<`, `>`, `:`, `"`, `|`, `?`, `*`, trailing spaces or periods, drive prefixes, URI prefixes, traversal components, and NTFS alternate-data-stream syntax. It also rejects the case-insensitive Windows device stems `CON`, `PRN`, `AUX`, `NUL`, `CLOCK$`, `CONIN$`, `CONOUT$`, `COM1` through `COM9`, and `LPT1` through `LPT9`, including the recognized superscript forms `COM¹` through `COM³` and `LPT¹` through `LPT³` and any of those stems followed by an extension.

Within one `source.assets` bundle, `file_name` values are unique by a portable collision key formed through Unicode canonical caseless matching: NFD normalization, default Unicode case folding, then NFD normalization again. Source preparation enforces that key and the complete restoration plan rejects destination collisions before output so distinct source assets cannot collapse onto one destination on case-insensitive or normalization-insensitive filesystems.

Timestamp metadata distinguishes true creation or birth time, platform creation time, fallback observations, and unavailable values. Unix `ctime` is never mislabeled as creation time. Captured metadata remains in the document even when the destination cannot restore it.

## Validation layers

Schema implementation distinguishes:

1. JSON parsing.
2. Draft 2020-12 structural validation.
3. Cueson semantic validation, including cross-field references and capability consistency.
4. Source-envelope integrity validation, including canonical base64, decoded length, SHA-256, safe basenames, and portable basename collisions; restoration separately validates its complete destination plan.
5. Official software/schema version equality.

Schema validation does not infer codec availability. Native codec presence remains an executable capability concern described by [the architecture](architecture.md) and [CLI contract](cli.md). Maintained documentation paths and local heading links are checked offline by the repository's standalone documentation verifier; that check does not publish the schema or validate external network availability.

Cross-format conversion loss reports and private target projections are runtime-only values. They are not Cue JSON properties, do not change schema version 1.0.0, and never replace the validated source document or its exact source envelope.

The versioned `inspect --json` report is likewise a separate CLI output contract rather than Cue JSON. It reports only bounded structural facts, excludes source bytes and content-bearing identifiers, and represents target-dependent conversion loss as not evaluated when no target was requested.
