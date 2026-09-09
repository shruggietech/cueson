# Cueson Project Specification

**Project:** `cueson`\
**Domain:** `https://cueson.io`\
**Owner:** ShruggieTech\
**Implementation language:** Go\
**License:** Apache License 2.0 (`Apache-2.0`)\
**Initial software version:** `0.0.0`\
**Initial schema version:** `0.0.0`\
**Document status:** Working pre-release architecture and repository draft\
**Primary audience:** Maintainers, AI coding agents, release automation, downstream integrators

## Draft authority and lifecycle

This document is the project's initial working specification. During repository bootstrap and pre-release design, maintainers MAY revise its implementation details, sequencing, examples, and provisional schema shapes as contradictions or better approaches emerge. Explicit operator decisions and the ratified project constitution control when they conflict with this draft. Core product invariants remain binding unless deliberately amended and documented. Only tagged release schemas and released public contracts are immutable.

The concise ratified implementation baselines are [architecture.md](architecture.md), [schema.md](schema.md), and [cli.md](cli.md). Those topic documents control implementation when this broader roadmap retains an older provisional alternative. Changes to a ratified baseline require a later Spec Kit decision and changelog entry.

## Executive summary

Cueson is an official ShruggieTech utility for converting subtitle and caption formats into and out of a canonical, versioned JSON representation called Cue JSON. The first stable release targets complete SubRip (`.srt`) and WebVTT (`.vtt`) support, and the project is deliberately designed so additional text, XML, and bitmap subtitle families can be added without replacing the core document model.

The project has two equally important products: the `cueson` executable and the Cueson JSON Schema. They are versioned in lock step. Every released binary emits the schema version matching its own software version, embeds that schema in the executable, and validates generated Cue JSON against it before reporting success.

Cueson MUST distinguish normalized data from source-preserving data. Normalized fields exist for search, analysis, format-independent processing, cross-format conversion, and downstream application development. Source-preserving fields exist so an encoded Cue JSON document can reproduce the original source asset byte for byte. No parser normalization, heuristic extraction, whitespace cleanup, timestamp canonicalization, or future OCR operation may destroy the original representation.

The canonical Cue JSON document therefore contains both a common semantic subtitle model and a lossless source-asset envelope. The common model is not optional ornamentation. It is the part of Cue JSON that downstream systems actually consume. It MUST expose cue timing, cue text, line structure, speaker observations, token timing where available, and format-native detail without requiring consumers to decode base64 or re-parse the source format. The envelope stores the original asset name, content hash, encoding observations, filesystem timestamp metadata, and original bytes encoded as base64. Original filesystem paths are prohibited.

Cueson is implemented in pure Go wherever practical and released as native standalone binaries for Windows, macOS, and Linux. The initial release family targets SRT and WebVTT. Future schema and software versions are expected to add TTML/IMSC, SMPTE-TT, ASS/SSA, EBU-STL, EBU-TT, PGS/SUP, VobSub IDX/SUB, and SAMI. The source-asset model is multi-asset from the beginning so future paired and binary formats do not require a foundational schema redesign.

## Product identity

The canonical project name is `cueson`, pronounced as "cues on" and derived from "Cue JSON."

The canonical project domain is `cueson.io`.

The domain has already been acquired and is managed through Cloudflare. Pre-v1 development MUST NOT modify production DNS, hosting, redirects, TLS, or other domain configuration. Public domain activation is intentionally deferred until after the v1 software release has completed successfully.

The canonical repository SHOULD be `shruggietech/cueson`.

The project is licensed under the Apache License, Version 2.0. The repository MUST use SPDX identifier `Apache-2.0` wherever a machine-readable license identifier is required.

The binary name is:

```text
cueson
```

Windows release artifacts use `cueson.exe`.

The canonical JSON media type is initially:

```text
application/vnd.cueson+json
```

The project MAY pursue an IANA media-type registration later. Until then, the vendor-tree media type is the project-defined interchange identifier.

The default encoded-document suffix is:

```text
.cueson.json
```

Given `captions.srt`, the default encoded output is `captions.srt.cueson.json`. Given `captions.vtt`, the default encoded output is `captions.vtt.cueson.json`. Keeping the source extension in the generated file name is intentional because it makes the source format immediately visible without opening the JSON.

## Goals

Cueson MUST:

- provide one canonical Cue JSON document model for supported subtitle and caption formats;
- provide first-class SRT and WebVTT parsing and encoding;
- guarantee exact restoration of original source assets from generated Cue JSON;
- expose normalized cue data without requiring consumers to understand the original subtitle format;
- make subtitle text and timing first-class, explicit, stable fields in the common cue model;
- preserve format-specific information that cannot be represented in the normalized common model;
- detect source format from content when possible and use the extension only as supporting evidence;
- store the original file name but never store the original filesystem path;
- preserve original source bytes and their SHA-256 digest;
- support multiple source assets in one document model so future paired formats such as VobSub do not require a breaking root-schema redesign;
- allow schema-recognized future formats to appear in the canonical schema before the CLI can ingest them natively;
- distinguish clearly between a format being recognized by the schema and being supported by the implementation;
- validate all emitted Cue JSON against the embedded canonical schema;
- support conversion from Cueson JSON back to the original source format;
- support cross-format SRT-to-WebVTT and WebVTT-to-SRT conversion with explicit loss accounting;
- expose a stable, scriptable CLI with clean stdout and diagnostic stderr;
- build as portable native executables for Windows, macOS, and Linux;
- keep the canonical versioned schema in the same repository as the implementation;
- publish immutable released schemas through repository and release artifacts, then through `cueson.io` after public-domain activation;
- integrate GitHub Spec Kit into the repository;
- require the ShruggieTech `shruggie-speckit` skill for full feature-slice autopilot work when that skill is available;
- use GitHub-native project management as the authoritative planning and delivery system;
- treat v1 delivery as an urgent project objective while preserving all quality and release gates;
- enable safe parallel work by coding sub-agents when independent tasks can be isolated and verified;
- maintain a detailed Keep a Changelog-compatible `CHANGELOG.md`;
- publish concise GitHub release notes containing highlights only and a link to the full changelog;
- start at software and schema version `0.0.0`;
- make `v1.0.0` the first stable contract with full SRT and WebVTT coverage.

## Non-goals for v1

Cueson v1 does not need to:

- edit subtitles interactively;
- provide a GUI;
- host a network service;
- become a subtitle playback engine;
- render subtitles onto video;
- infer semantic meaning from subtitle content beyond explicitly documented derived metadata;
- silently repair malformed source files;
- implement OCR before a bitmap subtitle codec enters active development;
- claim semantic support for bitmap subtitle formats before OCR extraction is implemented for those formats;
- support TTML, IMSC, SMPTE-TT, ASS, SSA, EBU-STL, EBU-TT, PGS, SUP, VobSub, or SAMI in the v1.0.0 release;
- guarantee lossless conversion between two different subtitle formats when the destination format cannot express source semantics;
- require every schema-recognized future format to have immediate native ingest support in the CLI;
- expose a stable public Go library API in v1.

The CLI and JSON schema are the v1 public contracts. Go packages SHOULD remain under `internal/` until a separate public library contract is deliberately specified.

## Governing invariants

### The schema is a first-class product

The schema is not generated documentation and is not subordinate to the executable. The executable and schema are peers.

Every release MUST contain:

- a `cueson` binary version;
- a schema with the exact same semantic version;
- an embedded copy of that schema in the binary;
- a release artifact containing the schema;
- an immutable canonical schema identifier under `cueson.io`, with repository, embedded, and release-artifact availability before the public domain is activated.

The executable MUST refuse to report a successful release build when its compiled version and embedded schema version differ.

### Software and schema versions move in lock step

If the software version is `1.3.2`, the emitted `schema_version` MUST be `1.3.2`.

The `$schema` URI in generated documents MUST point to the exact immutable release schema:

```text
https://cueson.io/schema/v1.3.2/cueson.schema.json
```

A convenience alias MAY exist at:

```text
https://cueson.io/schema/latest/cueson.schema.json
```

The `latest` alias is never canonical and MUST NOT be emitted into Cue JSON documents.

A released schema is immutable.

At and after v1.0.0, breaking changes to released software or schema contracts are permitted only at a major-version increase. Before v1.0.0, documented breaking contract changes may occur in a minor release while patch releases remain non-breaking. Additive changes may occur in minor versions. Corrections that do not change the accepted or emitted contract may occur in patch versions.

The repository begins at `0.0.0`. Until the first public release tag is created, the in-development `0.0.0` schema may be refined. After `v0.0.0` is released, its schema becomes immutable and the compatibility rules above apply.

### The common cue model is a first-class contract

Cue JSON is not merely a transport envelope for original subtitle bytes.

The canonical document MUST make timed subtitle content directly accessible through stable structured fields. For supported text-based formats this means Cue JSON MUST contain explicit cue timing and explicit cue text. For future bitmap formats it MUST provide a structured place for OCR-derived text and related provenance once that support arrives.

The common cue model exists so downstream systems can perform search, indexing, NLP, graph extraction, semantic comparison, quality checks, and cross-format processing without first re-parsing SRT, WebVTT, TTML, or future format grammars.

The common model therefore needs to be intentionally durable. Fields such as cue timing, `raw_text`, `plain_text`, `lines`, speakers, tokens, and placement are not incidental convenience fields. They are part of the long-haul public contract.

### Schema-recognized formats may precede codec support

The schema and the CLI are version-locked peers, but the set of schema-recognized formats does not need to equal the set of formats the current executable can ingest.

Cueson MUST distinguish at least these concepts:

- **schema-recognized format:** a format family named and structurally represented in the canonical schema;
- **implementation-supported format:** a format family for which the current executable ships a native codec or other documented ingest path;
- **support level:** the declared maturity of a schema-recognized format.

This distinction allows the project and community to debate and refine long-term document structure for future formats before the corresponding codec is implemented.

A future format may therefore appear in the schema in a reserved or experimental state, with the CLI able to validate a Cue JSON document using that structure even though `cueson encode` cannot yet ingest the original source asset directly.

The implementation MUST never pretend that a schema-recognized format is automatically ingest-supported. Commands that require a missing codec MUST fail clearly and truthfully.

### Exact restoration is distinct from normalized rendering

Cueson MUST support two reverse directions:

1. **Restore** reproduces the original source asset bytes exactly from the source envelope.
2. **Render** serializes the structured Cue JSON model into a subtitle format.

For an untouched Cueson document created by `encode`, `restore` MUST reproduce the original input byte for byte and MUST produce the same SHA-256 digest.

`render` is model-driven. It may produce canonical syntax rather than identical source bytes when the structured model has been edited.

Cross-format conversion uses rendering, not restoration.

This distinction prevents structured edits from being ignored merely because original bytes are available.

### OCR is a first-class derived representation

OCR is permanently in scope for Cueson. It is not required for the initial v1.0.0 SRT and WebVTT milestone because those formats already contain textual cue payloads, but the architecture MUST support OCR as a first-class derived representation from the beginning.

For any image-based subtitle family, source preservation alone is insufficient to claim useful semantic support. A format such as PGS/SUP or VobSub MAY be introduced experimentally before OCR is complete, but it MUST NOT be described as fully supported until Cueson can expose useful timed textual content derived from the subtitle images.

OCR output is derived data, never source truth. The original bitmap or binary source remains authoritative for exact restoration, while OCR supplies the searchable and processable text representation required by downstream consumers.

The model MUST permit multiple OCR observations for the same subtitle image or event so a future Cueson version can retain results from different engines, models, languages, or processing passes without overwriting prior observations.

Each OCR observation SHOULD be traceable to the source image or subtitle event from which it was derived and SHOULD support provenance fields including engine, engine version, model identifier when applicable, language, confidence, processing options, and source-image identity.

### Original bytes are authoritative for lossless provenance

Cueson MUST preserve original source bytes as base64 in the source envelope beginning with v0.0.0. Native SRT and WebVTT ingest remain deferred, so a v0.0.0 source envelope may be externally authored while still supporting generic integrity validation and exact restoration.

This requirement is intentionally stronger than the existing PowerShell converters. SRT has no universal encoding requirement, mixed line endings can occur within one file, malformed real-world inputs exist, and text decoding can be lossy. A normalized Unicode representation plus a single encoding or line-ending label cannot prove byte-exact reversibility in all cases.

Subtitle files are small enough that base64's approximate 33 percent storage overhead is acceptable in exchange for a truthful exact-restoration guarantee.

Future binary formats reuse the same source-asset mechanism instead of introducing a second binary-preservation model.

### Original filesystem paths are prohibited

Input filesystem paths are runtime-only information.

The schema MUST NOT contain:

- absolute source paths;
- relative source paths;
- source directories;
- working directories;
- home-directory values;
- drive letters derived from the source location;
- mount paths.

The source asset stores only the original base file name.

File names stored in the schema MUST be portable safe basenames. They MUST reject path separators, `.` and `..`, ASCII control characters, Windows-invalid punctuation and trailing characters, drive or URI prefixes, traversal constructs, NTFS alternate-data-stream syntax, and case-insensitive Windows reserved device names even when a device stem has an extension. Within one source bundle, basenames MUST be unique under Unicode canonical caseless matching, defined as NFD normalization, default Unicode case folding, then NFD normalization again. Semantic validation completes these checks before restoration opens any output.

Additional metadata is encouraged when it describes the media, subtitle track, source encoding, source format, language, tool version, or conversion process. Metadata MUST NOT be used to smuggle filesystem paths into the document.

### Parsers preserve, derived processors derive

Parsing MUST NOT rewrite source meaning.

Heuristic speaker extraction, markup stripping, entity decoding, word-token extraction, and similar convenience operations produce derived fields. They MUST NOT mutate the source-preserving payload.

For example, a source SRT cue containing:

```text
William: This is a line.
```

may produce a derived speaker value of `William`, but the source payload remains exactly `William: This is a line.`

### Unknown input is never silently discarded

When Cueson can preserve an unrecognized or nonconforming source fragment but cannot interpret it, the parser MUST preserve that fragment and record a diagnostic.

A warning is preferable to data loss.

A parser MUST NOT silently skip an unknown WebVTT block, malformed SRT block, unsupported tag, or unrecognized format-specific setting if the bytes can be preserved.

### Normalized data is convenient, raw data is authoritative for reconstruction

Common fields exist to make applications easy to write. They are not a replacement for format-specific source fidelity.

If a normalized value and source-preserving value disagree because a source file is malformed, the discrepancy MUST remain visible through diagnostics. Cueson MUST NOT silently rewrite the source envelope to match the normalized model.

## Existing converter baseline

### Existing SRT converter

The existing `ConvertFrom-Srt.ps1` establishes useful behavior that Cueson SHOULD preserve semantically:

- cue sequence numbers;
- normalized start and end times;
- derived start, end, and duration milliseconds;
- optional SRT coordinate values;
- cue content;
- optional heuristic speaker detection;
- tolerance for a UTF-8 BOM;
- tolerance for CRLF, LF, and lone-CR inputs;
- tolerance for `.` in place of `,` before milliseconds;
- tolerance for one-to-three digit millisecond fields;
- preservation of inline SRT styling tags;
- pipeline-friendly JSON output;
- explicit force behavior for overwrite protection;
- a documented `0`, `1`, `2` exit-code contract.

The existing SRT JSON contract is a top-level array of cues.

Cueson MUST NOT preserve the existing converter's destructive normalization behavior as the source-of-truth representation. The current script normalizes line endings, trims trailing whitespace, normalizes timestamps, and may remove a detected speaker prefix from `content`. Those transformations are acceptable as derived views but not as the only stored representation.

### Existing WebVTT converter

The existing `ConvertFrom-WebVtt.ps1` establishes a richer whole-document model that Cueson SHOULD preserve semantically:

- WebVTT signature and header description;
- header metadata;
- NOTE blocks;
- STYLE blocks;
- REGION blocks;
- cue identifiers;
- cue timing;
- cue settings and raw setting strings;
- raw cue payloads;
- plain-text payloads;
- voice extraction;
- inline timestamp token extraction;
- preservation of overlapping rolling cues;
- summary statistics;
- source BOM and line-ending observations.

Cueson MUST remove the current absolute `filePath` field from the unified contract.

The current WebVTT model stores non-cue blocks in separate arrays with positional hints such as `afterCue`. Cueson SHOULD preserve explicit source order as a first-class integer so multiple adjacent non-cue blocks remain totally ordered without inference.

### Unification requirements

Cueson MUST resolve the following differences rather than expose them as two unrelated schemas:

| Concern | Existing SRT | Existing WebVTT | Cueson |
|---|---|---|---|
| Top level | Cue array | Whole document | Whole document |
| Naming | Mostly snake_case | camelCase | Lowercase snake_case |
| Cue ordinal | Source index | Zero-based ordinal | Separate normalized ordinal and source identifier |
| Time convenience | Milliseconds | Seconds | Milliseconds in common model, native raw timing in format data |
| Raw source bytes | No | No | Required base64 source assets |
| Source filename | Not in output | Yes | Required |
| Source path | Not in output | Absolute path | Prohibited |
| Non-cue structure | Not applicable | Separate arrays | Ordered format structure |
| Speaker data | Heuristic and destructive to content | Native voice extraction | Derived and non-destructive |
| Exact restoration | No | No | Required |
| Cross-format conversion | No | No | Required |

## Canonical Cue JSON architecture

### Schema key and object style

All Cueson-owned JSON property names MUST use lowercase `snake_case`.

JSON Schema keywords such as `$schema`, `$id`, `$defs`, `oneOf`, and `additionalProperties` retain their standards-defined spellings and are exempt from this rule.

Schema objects SHOULD favor small reusable semantic types rather than long flattened property names. Repeated concepts should have one consistent shape wherever they appear.

Representative style rules include:

- timestamps are structured values rather than timestamp strings scattered across unrelated objects;
- file size uses a grouped size object containing a human-readable representation and exact bytes;
- digests are grouped under a `hashes` object;
- encoding observations are grouped under an `encoding` object;
- provenance belongs next to the value whose provenance it qualifies;
- machine-precise values and human-readable values may coexist when each has a distinct operational purpose;
- boolean names use an `is_`, `has_`, or similarly explicit predicate form when doing so removes ambiguity;
- abbreviations such as `id`, `sha256`, `srt`, `vtt`, and `ocr` remain lowercase;
- project-owned enum values use lowercase identifiers, with underscores when multiple words are required.

The schema SHOULD avoid encoding units only in a long property name when a reusable typed object communicates the unit and semantics more cleanly. Where precision or interoperability makes a unit-bearing key preferable, the unit suffix MUST be explicit and stable.

### Root object

The root object MUST use JSON Schema Draft 2020-12 and MUST have a fixed documented shape.

A representative v1 document is:

```json
{
  "$schema": "https://cueson.io/schema/v1.0.0/cueson.schema.json",
  "schema_version": "1.0.0",
  "format": "subrip",
  "format_support": {
    "status": "stable",
    "ingest_supported": true,
    "render_supported": true,
    "restore_supported": true,
    "ocr_required_for_semantic_output": false
  },
  "producer": {
    "name": "cueson",
    "version": "1.0.0"
  },
  "source": {
    "primary_asset_id": "asset-0",
    "assets": [
      {
        "id": "asset-0",
        "role": "primary",
        "file_name": "captions.srt",
        "media_type": "application/x-subrip",
        "size": {
          "text": "53 bytes",
          "bytes": 53
        },
        "hashes": {
          "sha256": "955096d8a9c0df1a51bc836d0fd7cfc2d06c75807c5908b5017791d436673eea"
        },
        "timestamps": {
          "created": {
            "iso": "2026-09-09T00:00:00.000000000Z",
            "unix_ns": 1788912000000000000
          },
          "modified": {
            "iso": "2026-09-09T00:00:00.000000000Z",
            "unix_ns": 1788912000000000000
          },
          "accessed": {
            "iso": "2026-09-09T00:00:00.000000000Z",
            "unix_ns": 1788912000000000000
          },
          "created_source": "windows_creation_time"
        },
        "encoding": {
          "bom": null,
          "line_endings": "lf",
          "detected_encoding": "utf-8",
          "confidence": 1.0
        },
        "data_base64": "MQowMDowMDowMSwyNTAgLS0+IDAwOjAwOjA0LDIwMAo8aT5IZWxsbywgd29ybGQuPC9pPgo="
      }
    ]
  },
  "metadata": {
    "title": null,
    "language": "en",
    "kind": "subtitles",
    "description": null
  },
  "document": {
    "cue_count": 1,
    "media_start_milliseconds": 1250,
    "media_end_milliseconds": 4200,
    "media_span_milliseconds": 2950,
    "has_word_level_timing": false
  },
  "cues": [
    {
      "id": "cue-000000",
      "ordinal": 0,
      "source_order": 0,
      "source_identifier": "1",
      "timing": {
        "start_milliseconds": 1250,
        "end_milliseconds": 4200,
        "duration_milliseconds": 2950
      },
      "payload": {
        "raw_text": "<i>Hello, world.</i>",
        "plain_text": "Hello, world.",
        "lines": [
          "Hello, world."
        ]
      },
      "speakers": [],
      "tokens": [],
      "ocr_observations": [],
      "placement": null,
      "format_data": {
        "subrip": {
          "sequence_line_raw": "1",
          "timing_line_raw": "00:00:01,250 --> 00:00:04,200"
        }
      }
    }
  ],
  "format_data": {
    "subrip": {
      "dialect": "subrip"
    }
  },
  "diagnostics": [],
  "stats": {
    "cue_count": 1,
    "diagnostic_count": 0,
    "warning_count": 0,
    "error_count": 0,
    "has_word_level_timing": false,
    "media_span_milliseconds": 2950
  }
}
```

The exact schema is developed through Spec Kit and becomes authoritative. The example above defines the intended architecture, not a substitute for the canonical JSON Schema.

### `$schema`

`$schema` identifies the immutable released Cueson schema.

It MUST use the exact software/schema version.

This instance member is distinct from the schema artifact's standards-defined `$schema` keyword, which identifies the JSON Schema Draft 2020-12 dialect. The schema artifact's `$id` and the Cue JSON instance's project-defined `$schema` member use the canonical versioned Cueson URI.

### `schema_version`

`schema_version` is the semantic version of the document contract.

Documents emitted by the official Cueson executable MUST use the schema version embedded in that executable. A third-party producer's own software version is independent and MAY differ from the `schema_version` of the Cue JSON contract it targets.

### `format`

`format` identifies the source subtitle family.

Initial canonical format keys:

```text
subrip
webvtt
```

These keys identify format families and do not claim stable codec support. Stable SRT and WebVTT support remains a v1.0.0 gate. Additional schema-recognized values may be declared before native ingest support arrives. Such values are still valid in the canonical schema when their document shape has been defined and the corresponding `format_support.status` truthfully reflects their maturity.

Future values may include:

```text
ttml
imsc
smpte_tt
ass
ssa
ebu_stl
ebu_tt
pgs
vobsub
sami
```

Adding a new format under an existing compatible root model is an additive change unless the common contract itself must break.

The schema MAY define document structures for additional formats ahead of implementation support. When it does so, those structures MUST be clearly marked through `format_support.status` and through the corresponding documentation pages.


### `format_support`

`format_support` describes the maturity and official capability state assigned to the document's `format` by the Cueson release associated with the document contract. It does not describe arbitrary capabilities of a third-party producer.

For the completed v0.0.0 milestone, `subrip` and `webvtt` are `envelope_only`: native ingest and model-driven render are false, generic exact restoration is true for a valid source envelope, and OCR is not required for semantic output. Source-foundation issue #6 owns the public generic `restore` command required for that capability claim.

Representative fields are:

```text
status
ingest_supported
render_supported
restore_supported
ocr_required_for_semantic_output
```

`status` expresses support maturity, not merely parser presence. The long-haul baseline SHOULD define at least these values:

```text
reserved
experimental
envelope_only
beta
stable
deprecated
```

Suggested semantics:

- `reserved`: the schema names the format and may define baseline structural placeholders, but no official ingest contract exists yet;
- `experimental`: a provisional structure exists and may change in compatible or pre-stable ways according to the active version policy;
- `envelope_only`: the project can preserve source assets and validate baseline document structure, but native semantic ingest is not yet claimed;
- `beta`: native ingest exists but the feature is not yet considered stable;
- `stable`: the format is part of the supported contract for the current major version;
- `deprecated`: still recognized, but scheduled for removal or replacement in a future major version.

`ingest_supported` states whether the current executable can ingest original source assets for this format.

`render_supported` states whether the current executable can render this document back into the indicated source format from the structured model.

`restore_supported` states whether exact asset restoration from `source.assets` is supported. This is normally true whenever the source envelope is valid.

`ocr_required_for_semantic_output` declares whether usable semantic text for this format depends on OCR or a similar derived extraction stage.

The schema SHOULD require `format_support` for emitted documents so downstream tools can distinguish mature content from forward-declared structures.

### `producer`

`producer` identifies the software that created the document.

Required fields:

```text
name
version
```

Optional non-sensitive build metadata MAY include commit identity or build target. Build metadata MUST NOT include filesystem paths.

### `source`

`source` is the exact-restoration envelope.

It MUST be multi-asset from the first schema even though SRT and WebVTT commonly use one source file.

This avoids a future breaking redesign for formats such as VobSub, where an `.idx` text index and `.sub` binary payload are a logical pair.

`source.primary_asset_id` identifies the main source asset.

`source.assets` is an ordered array.

Each asset MUST contain or explicitly null the stable fields required by the canonical schema. A representative source asset is:

```json
{
  "id": "asset-0",
  "role": "primary",
  "file_name": "captions.srt",
  "media_type": "application/x-subrip",
  "size": {
    "text": "18.42 KB",
    "bytes": 18420
  },
  "hashes": {
    "sha256": "..."
  },
  "timestamps": {
    "created": {
      "iso": "2026-09-08T18:02:03.123456700Z",
      "unix_ns": 1788890523123456700
    },
    "modified": {
      "iso": "2026-09-08T18:12:19.987654300Z",
      "unix_ns": 1788891139987654300
    },
    "accessed": {
      "iso": "2026-09-08T18:15:00.000000000Z",
      "unix_ns": 1788891300000000000
    },
    "created_source": "birthtime"
  },
  "encoding": {
    "bom": null,
    "line_endings": "crlf",
    "detected_encoding": "utf-8",
    "confidence": 1.0
  },
  "data_base64": "..."
}
```

`id` is a document-local source-asset identifier.

`role` initially supports `primary` and `companion`.

`file_name` is the original source basename. It MUST satisfy the portable safe-basename contract, including rejection of directory paths, control characters, Windows-invalid punctuation and trailing characters, drive or URI prefixes, traversal components, NTFS alternate-data-stream syntax, and case-insensitive Windows reserved device names even when followed by an extension. Its portable collision key MUST be unique within `source.assets`.

`media_type` describes the source asset media type when known.

`size` groups exact byte size and an optional human-readable representation. `size.bytes` is authoritative for integrity validation. `size.text` is derived convenience metadata.

`hashes` groups integrity digests. SHA-256 is required beginning with v0.0.0. Digest values MUST use lowercase hexadecimal consistently in the schema and implementation.

`timestamps` preserves the source filesystem's observable date information and is discussed separately below.

`encoding` contains textual encoding observations. It SHOULD use stable fields such as `bom`, `line_endings`, `detected_encoding`, and `confidence` rather than independent top-level encoding keys.

`data_base64` is the original asset bytes encoded with standard base64 and is authoritative for byte-exact restoration.

#### Source filesystem timestamps

Each source asset MUST attempt to capture the original file's filesystem timestamps before Cueson opens or reads the file content.

The initial timestamp set is:

```text
created
modified
accessed
created_source
```

`created`, `modified`, and `accessed` are timestamp objects.

A timestamp object contains:

```text
iso
unix_ns
```

`iso` is a human-readable RFC 3339 timestamp with the greatest precision faithfully available from the source API. Cueson SHOULD serialize it in UTC so the representation is deterministic across machines. Filesystems generally store an instant rather than an originating timezone, so Cueson MUST NOT invent a source timezone that the filesystem did not actually retain.

`unix_ns` is the same instant represented as signed Unix epoch nanoseconds. It is the machine-precise restoration value and MUST be derived directly from the native filesystem timestamp without first round-tripping through the ISO string.

Nanosecond representation is required even when a filesystem has coarser precision. For example, Windows FILETIME values naturally land on 100-nanosecond boundaries. The stored integer therefore preserves the source API value without reducing it to milliseconds.

`created_source` records how the creation value was obtained. The canonical enum is finalized during schema work but MUST distinguish a true creation or birth time from any fallback approximation. Representative values include:

```text
birthtime
windows_creation_time
ctime_fallback
unavailable
```

When a platform exposes additional date metadata that is useful for forensic preservation but is not safely restorable across platforms, Cueson MAY add it as an optional observational field in a compatible schema version. Such a value MUST identify its semantics accurately. In particular, Unix `ctime` MUST NOT be mislabeled as creation time.

#### Timestamp capture ordering

Timestamp capture is part of source preservation, not an afterthought.

Encoding MUST conceptually perform:

1. stat or otherwise query the source file using the best native API available;
2. capture filesystem timestamps and their provenance;
3. read the source bytes;
4. compute byte hashes and continue parsing.

This order prevents Cueson's own source read from becoming the `accessed` value stored in the document on filesystems that update access time.

#### Timestamp restoration

`restore` MUST attempt to reproduce captured filesystem timestamps in addition to restoring the original bytes and original basename.

Timestamp restoration is platform- and filesystem-aware. A portable Go implementation MUST use native operating-system facilities where the generic standard library cannot provide the required fidelity.

The intended capability model is:

- Windows: use native file-time APIs to restore creation, modification, and access times when the destination filesystem supports them;
- macOS: use native APIs to restore modification and access times and to restore creation/birth time where the destination filesystem and API permit it;
- Linux: restore modification and access times, capture birth time through `statx` when available, and report honestly that Linux generally does not provide a normal user-space mechanism for setting filesystem birth time after file creation;
- other future targets: document exact capture and restoration capabilities before claiming support.

Cueson MUST distinguish **metadata preservation** from **metadata restoration capability**. The Cue JSON document must retain every captured timestamp even when the current destination platform cannot reproduce one of them.

Restoration MUST apply file times only after:

1. source bytes have been written;
2. byte length has been verified;
3. SHA-256 has been verified;
4. the final output name is in place.

Timestamp application is the final mutating step so Cueson's own writes and integrity reads do not immediately replace the values being restored.

After setting timestamps, Cueson MUST stat the restored file and compare every timestamp that the destination claims to support. A metadata restoration report MUST identify each timestamp as restored, unsupported, unavailable in source metadata, or failed.

The restore CLI MUST provide a strict metadata mode:

```text
--strict-metadata
```

When enabled, inability to reproduce any captured timestamp that is part of the source restoration request is a runtime failure with exit code `1`. The implementation SHOULD stage output so a strict metadata failure does not leave a partially accepted restored artifact.

The restore CLI SHOULD also provide:

```text
--no-metadata
```

for callers who intentionally want byte restoration without filesystem metadata restoration.

Without `--strict-metadata`, an unsupported timestamp MUST produce a warning rather than being silently ignored.

Exact-byte restoration and exact-filesystem-metadata restoration MUST be reported separately. Cueson MUST never claim that a destination file has identical filesystem dates when the operating system or destination filesystem prevented one or more captured values from being applied.

### `metadata`

`metadata` contains normalized document-level metadata useful across formats.

The initial schema SHOULD provide typed fields for common metadata such as:

```text
title
language
kind
description
```

Fields without values SHOULD be null or empty according to the fixed schema contract.

Arbitrary source-derived metadata MAY be represented in a dedicated namespaced extension map. The schema and application MUST prevent filesystem provenance from being added through this mechanism.

### `document`

`document` contains normalized document-level structure that is meaningful across formats.

Initial fields SHOULD include:

```text
cue_count
media_start_milliseconds
media_end_milliseconds
media_span_milliseconds
has_word_level_timing
```

Format-specific structures do not belong here unless their semantics genuinely apply across formats.

### `cues`

Every supported timed-text format maps timed presentation units into the common `cues` array when that mapping is meaningful.

For schema-recognized formats that do not yet have native semantic ingest support, the schema may still define how a future or externally-authored document represents cues, payloads, OCR observations, placement, and format-native data. This lets the contract mature before the codec does.


The `cues` array is the primary semantic surface of Cue JSON. Consumers should be able to obtain subtitle timing and subtitle text from this array without re-parsing the original format.

A common cue SHOULD contain:

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

#### `id`

`id` is a deterministic document-local identifier such as `cue-000000`.

It is not a UUID and does not claim identity across independently produced files.

#### `ordinal`

`ordinal` is the zero-based cue ordinal in normalized cue order.

It does not replace an SRT sequence number or WebVTT cue identifier.

#### `source_order`

`source_order` is the zero-based position of the source construct among all ordered source constructs represented by the parser.

This lets cues, notes, styles, regions, and future block types share a total source order.

#### `source_identifier`

`source_identifier` preserves a format's source-level cue identity when one exists.

For SRT this may preserve the literal sequence-line value.

For WebVTT it may preserve the cue identifier.

The normalized cue ordinal remains independent.

#### `timing`

The v1 common timing model uses integer milliseconds:

```text
start_milliseconds
end_milliseconds
duration_milliseconds
```

Format-specific raw timing syntax remains in `format_data`.

Future formats that require greater precision MAY add a higher-precision normalized representation additively when possible. Native timing expressions and original bytes remain authoritative for exact source restoration.

#### `payload`

The common payload SHOULD contain:

```text
raw_text
plain_text
lines
```

These fields are part of the project's durable baseline and are expected to remain useful across the long haul.

`raw_text` is the decoded cue payload without destructive semantic normalization. For text-based formats it is the first-class representation of what the cue says in that format's textual surface.

`plain_text` is a derived consumer-friendly text representation with format markup removed where a deterministic stripping rule exists. This is the baseline field for indexing, search, NLP, and downstream graph extraction.

`lines` preserves the cue's logical text line structure as an ordered array. This avoids forcing consumers to reverse-engineer newline handling from a single text blob.

The schema SHOULD also allow future additive payload fields such as `normalized_text`, `language`, or `ruby_text` when a stable cross-format abstraction proves useful.

`raw_text` is not the byte-level source-of-truth. Exact bytes remain in `source.assets[].data_base64`.

#### `speakers`

`speakers` is an array of derived or source-native speaker observations.

Each observation SHOULD identify its origin:

```text
native
heuristic
```

A WebVTT `<v Speaker>` voice annotation is `native`.

An SRT `Name: dialogue` pattern is `heuristic`.

Speaker extraction MUST NOT mutate `payload.raw_text`.

#### `tokens`

`tokens` contains word- or phrase-level timing when the source provides it.

WebVTT inline timestamps map naturally here.

SRT normally emits an empty array.

A token SHOULD include normalized timing plus raw/native timing context when available.

#### `ocr_observations`

The common cue model SHOULD reserve space for OCR-derived observations so bitmap subtitle formats can mature into first-class semantic citizens rather than remaining blind envelopes.

`ocr_observations` is a required ordered array so results from multiple engines, models, languages, or processing passes can coexist without overwriting one another. For text-native formats this field is an empty array when no observation exists.

Each OCR observation SHOULD be able to represent:

```text
id
derived
engine
engine_version
model
language
text
lines
confidence
alternatives
source_asset_id
source_event_id
source_image_id
source_regions
processing_options
```

`derived` is a boolean indicating that the text content was produced through OCR rather than being native subtitle text.

`alternatives` may preserve multiple hypotheses produced by one observation when that materially improves downstream quality. It does not replace the observation array.

`source_regions` may describe the image region or timing segment from which the OCR result was derived.

OCR-derived text MUST NOT replace the authoritative source envelope. It is a semantic convenience layer with provenance.

#### `placement`

`placement` contains normalized positioning information that genuinely maps across formats.

The initial model SHOULD be conservative. It SHOULD NOT pretend that SRT pixel coordinates and WebVTT percentage/line/region settings are semantically identical.

When no safe normalized mapping exists, `placement` remains null or partial and the complete native information remains in `format_data`.

#### Cue `format_data`

Each cue carries exactly one format-specific object keyed by the source format.

For example:

```json
{
  "subrip": {
    "sequence_line_raw": "17",
    "timing_line_raw": "00:00:14,100 --> 00:00:17,240 X1:40 X2:600 Y1:20 Y2:50",
    "coordinates": {
      "x1": 40,
      "x2": 600,
      "y1": 20,
      "y2": 50
    }
  }
}
```

or:

```json
{
  "webvtt": {
    "identifier_raw": "intro",
    "timing_line_raw": "00:00:14.100 --> 00:00:17.240 align:start position:0%",
    "settings_raw": "align:start position:0%",
    "settings": {
      "align": "start",
      "position": "0%"
    },
    "raw_payload": "<v Narrator>Hello.</v>"
  }
}
```

Format-specific data SHOULD prefer preserving information over prematurely forcing it into shared abstractions.

### Root `format_data`

The root `format_data` object contains document-level format-specific structure.

For SRT it MAY contain dialect and decoding observations.

For WebVTT it MUST represent:

- header signature;
- header description;
- parsed and raw header metadata;
- NOTE blocks;
- STYLE blocks;
- REGION blocks;
- any preserved unrecognized blocks;
- total `source_order` for every non-cue construct.

Each non-cue block MUST carry `source_order`.

The schema MUST preserve multiple non-cue blocks in exact source order without relying solely on `afterCue`.

### `diagnostics`

Cueson MUST make tolerated defects and lossy interpretations visible.

Each diagnostic SHOULD contain:

```text
severity
code
message
source_order
cue_id
```

`source_order` and `cue_id` may be null when the diagnostic applies to the whole document.

Severity values SHOULD initially be:

```text
info
warning
error
```

Diagnostics are data about parsing or conversion. They are not log messages and therefore belong in the JSON when relevant.

### `stats`

`stats` contains derived counters and summary values.

Initial fields SHOULD include:

```text
cue_count
diagnostic_count
warning_count
error_count
has_word_level_timing
media_span_milliseconds
```

Format-specific counts MAY live under format-specific stats.

## Schema mechanics

The canonical schema MUST use JSON Schema Draft 2020-12.

The root schema SHOULD use `$defs` for reusable objects and `oneOf` or conditional `if`/`then` rules keyed by `format` for format-specific requirements.

The schema SHOULD set `additionalProperties: false` at all stable fixed object boundaries.

Explicit extension maps MAY allow namespaced keys where extensibility is intentional.

The schema MUST validate:

- version values;
- supported source format values;
- asset file-name safety;
- SHA-256 shape;
- base64 syntax;
- required common cue fields;
- format-specific cue objects;
- format-specific root structures.

Application-level semantic validation MUST supplement JSON Schema where the constraint cannot be expressed reliably in JSON Schema, including:

- `hashes.sha256` matches decoded `data_base64`;
- `size.bytes` matches decoded bytes;
- `primary_asset_id` resolves to an asset;
- one and only one primary asset exists for single-file formats;
- cue duration arithmetic is internally consistent;
- software/schema compatibility rules;
- source-order uniqueness and ordering;
- format-specific timing constraints;
- exact-restore digest verification;
- timestamp object consistency between `iso` and `unix_ns`;
- timestamp provenance consistency;
- source timestamp capture occurs before source-byte reads;
- restored timestamp verification on each supported operating system.

## Format detection

The CLI MUST support explicit format selection and automatic detection.

Automatic detection SHOULD use evidence in this order:

1. content signature or grammar;
2. source structure;
3. file extension as supporting evidence.

WebVTT detection requires the `WEBVTT` signature after permissible BOM handling.

SRT detection uses the de facto cue grammar, including a timecode arrow with SubRip-style time expressions and cue-block structure.

An extension disagreement SHOULD produce a warning rather than silently override stronger content evidence.

Ambiguous input MUST fail with an actionable diagnostic unless the operator supplies `--format`.

## SRT v1 coverage

Cueson v1.0.0 MUST cover the de facto SubRip grammar and the real-world variations already supported by the ShruggieTech SRT converter.

Coverage MUST include:

- integer cue sequence lines;
- missing or irregular sequence values without destructive renumbering;
- canonical comma millisecond separator;
- period millisecond separator as a tolerated variant;
- one-to-three digit millisecond fields;
- multi-line cue text;
- CRLF, LF, lone-CR, and mixed line endings;
- UTF-8 with and without BOM;
- common legacy text encodings under an explicit decoding policy;
- `--encoding` override when automatic decoding is ambiguous;
- optional `X1`, `X2`, `Y1`, `Y2` coordinate suffixes;
- inline formatting tags preserved in raw text;
- malformed but preservable blocks surfaced through diagnostics;
- optional derived speaker detection without source mutation.

Because SRT has no single universally authoritative specification or encoding mandate, "full SRT coverage" means complete coverage of the documented Cueson SRT grammar and tolerated-dialect contract, backed by corpus fixtures and tests. The project MUST document that grammar.

The parser MUST use a deterministic state machine or equivalent explicit parser. It MUST NOT depend on splitting the entire file solely on normalized blank lines in a way that destroys source distinctions.

## WebVTT v1 coverage

Cueson v1.0.0 MUST cover conforming WebVTT document structure and the real-world YouTube rolling-caption patterns already handled by the ShruggieTech converter.

Coverage MUST include:

- UTF-8 BOM detection;
- `WEBVTT` signature;
- optional header description;
- header metadata lines;
- NOTE blocks;
- STYLE blocks;
- REGION blocks;
- cue identifiers;
- hours-optional timestamps;
- cue timing settings;
- cue positioning settings;
- inline markup;
- voice tags;
- class tags;
- language tags;
- ruby-related markup;
- HTML entities;
- inline timestamps;
- word/phrase token extraction;
- whitespace-only cue payload lines;
- overlapping rolling cues without deduplication;
- source-order preservation for all blocks;
- preserved unrecognized blocks and diagnostics instead of silent dropping.

The parser MUST preserve raw payload text and raw structural lines independently of normalized consumer fields.

## Encoding policy

WebVTT is decoded according to its UTF-8 requirement.

SRT requires a broader policy because legacy files may use multiple encodings.

The v1 SRT decoder SHOULD:

1. honor a Unicode BOM when present;
2. attempt strict UTF-8;
3. detect UTF-16 LE/BE when evidence supports it;
4. support common legacy encodings through a pinned Go text-encoding dependency;
5. accept an explicit `--encoding` override;
6. fail processing with a clear diagnostic when text cannot be decoded safely.

Even when semantic decoding fails, the source bytes can still be preserved. The application MAY provide a future raw-envelope mode for unsupported textual encodings, but v1 success requires successful parsing.

The exact list of supported named encodings MUST be documented and tested before v1.0.0.

## Restore, render, and convert behavior

### Exact restore

```text
cueson restore INPUT.cueson.json
```

Restores the original source asset or asset set from `source.assets[].data_base64`.

With no destination option, a single asset uses its stored portable safe basename in the current directory. `--output-dir` restores one or more assets beneath the selected directory using their stored basenames. A single-asset `--output` supplies a separately validated literal runtime destination and may rename the file. `--output` and `--output-dir` are mutually exclusive.

Before opening any output, Cueson MUST validate all stored basenames, compute the complete destination plan, and reject duplicate portable basename keys or destination paths. `--force` MUST NOT permit one source asset to replace another asset from the same bundle.

The command MUST refuse to overwrite existing files unless `--force` is supplied.

After writing, the command MUST recompute SHA-256 and compare it with the stored digest before reporting byte-restoration success.

After byte verification, the command MUST attempt filesystem timestamp restoration unless `--no-metadata` is supplied. Timestamp verification occurs after the values are applied.

For multi-asset documents, all assets are restored to the selected output directory and each asset receives its own preserved timestamp metadata.

`--strict-metadata` requires every captured source timestamp requested for restoration to be reproduced and verified. `--no-metadata` intentionally skips filesystem timestamp restoration. The two options are mutually exclusive.

### Render

```text
cueson render INPUT.cueson.json --to srt
```

or:

```text
cueson render INPUT.cueson.json --to vtt
```

Renders the structured model into the requested subtitle format.

When rendering to the original format, Cueson SHOULD reuse original lexical details when the relevant structured values have not changed and the reuse is safe. It MAY otherwise emit canonical syntax.

`render` does not claim byte identity. `restore` does.

### Cross-format convert

```text
cueson convert INPUT.vtt --to srt
```

`convert` accepts a supported subtitle source or Cue JSON and writes the target subtitle format.

Cross-format conversion MAY be lossy because source formats do not have equivalent feature sets.

Every conversion MUST create an in-memory loss report.

By default, representable content is written and losses are emitted as warnings to stderr.

With:

```text
--strict
```

any non-representable source semantic MUST fail the conversion with exit code `1` and MUST NOT create the target file.

Examples of potential WebVTT-to-SRT loss include:

- STYLE blocks;
- REGION semantics;
- NOTE blocks;
- cue settings without an SRT equivalent;
- inline word timing;
- WebVTT-specific markup.

Examples of potential SRT-to-WebVTT adaptation include pixel coordinate semantics that cannot be mapped safely without viewport information.

Cueson MUST never imply that cross-format conversion is lossless merely because the original Cueson source envelope remains lossless.

## Behavior for schema-recognized but implementation-unsupported formats

The CLI MUST behave truthfully when encountering a schema-recognized format for which no native codec is present.

Expected baseline behavior:

- `cueson validate` SHOULD validate Cue JSON documents for any schema-recognized format whose structure is defined by the current canonical schema;
- `cueson inspect` SHOULD report the format's declared `format_support` state;
- `cueson restore` SHOULD work whenever the source envelope is valid, regardless of whether native ingest is currently supported;
- `cueson encode` MUST fail clearly when asked to ingest a source format without an implemented codec;
- `cueson render` and `cueson convert` MUST fail clearly when asked to produce a format without an implemented renderer;
- failure messages SHOULD differentiate "unknown format", "schema-recognized but codec missing", and "recognized but lossy conversion blocked by strict mode".

This separation is important to the long-haul contract. A future document using a reserved or experimental format should still be a valid Cue JSON document if it conforms to the current schema, even when the present executable cannot author it from raw source files.

## CLI contract

### General principles

The CLI MUST generally mirror the operator behavior established by ShruggieTech Bash and PowerShell utilities while remaining idiomatic Go.

Required conventions:

- `-h` and `--help` are reserved exclusively for help;
- `-q` and `--quiet` suppress informational and success diagnostics but retain warnings and errors;
- `--silent` suppresses non-error diagnostics;
- `--no-color` disables terminal color;
- the `NO_COLOR` environment convention is honored;
- color is automatically disabled when the diagnostic stream is not a TTY;
- stdout contains payload data only;
- diagnostics and progress go to stderr;
- no emojis appear in CLI output;
- overwrite is never implicit;
- `--force` is required to replace an existing output;
- `--` ends option processing;
- file paths supplied by the user are treated literally;
- command help contains meaningful descriptions and multiple worked examples;
- text output uses UTF-8 without BOM and LF unless the output is an exact restored source;
- structured JSON emitted to stdout contains no decorative text.

### Exit codes

The baseline contract is:

| Code | Meaning |
|---:|---|
| `0` | Success |
| `1` | Runtime, parse, validation, rendering, conversion, integrity, or assertion failure |
| `2` | Environment or invocation precondition failure |

Higher codes MUST NOT be introduced without updating the project constitution, CLI documentation, and changelog.

### Commands

The v1 target command surface SHOULD include:

```text
cueson encode
cueson restore
cueson render
cueson convert
cueson validate
cueson inspect
cueson schema
cueson version
cueson completion
```

Release help and command listings MUST expose only implemented commands. The v0.0.0 executable foundation begins with help and `version`; the schema slice adds `schema`; the source-foundation slice adds generic `restore`. Other target commands remain absent until their implementation slices provide truthful behavior. The release-specific authority is [cli.md](cli.md).

#### `encode`

```text
cueson encode INPUT [options]
```

Parses a subtitle source and writes Cue JSON.

Core options:

```text
-o, --output PATH
-f, --force
--format auto|srt|vtt
--encoding NAME
--pretty
--stdout
-q, --quiet
--silent
--no-color
-h, --help
```

`--stdout` and `-o -` MAY be equivalent.

When writing to stdout, informational output MUST be suppressed automatically so the stream contains only Cue JSON.

#### `restore`

```text
cueson restore INPUT.cueson.json [options]
```

Restores exact original asset bytes.

Core options:

```text
-o, --output PATH
--output-dir DIRECTORY
-f, --force
--strict-metadata
--no-metadata
-q, --quiet
--silent
--no-color
-h, --help
```

For a multi-asset source, `--output-dir` is the natural destination.

#### `render`

```text
cueson render INPUT.cueson.json --to FORMAT [options]
```

Serializes the structured model.

Core options:

```text
--to srt|vtt
-o, --output PATH
-f, --force
--strict
-q, --quiet
--silent
--no-color
-h, --help
```

#### `convert`

```text
cueson convert INPUT --to FORMAT [options]
```

Performs source-to-source or Cue-JSON-to-source conversion through the common model.

It MUST use the same loss-accounting rules as `render`.

#### `validate`

```text
cueson validate INPUT
```

For Cue JSON, validation includes:

- JSON parsing;
- JSON Schema validation;
- semantic validation;
- source-envelope integrity validation;
- software/schema compatibility checks.

For SRT or WebVTT, validation checks source grammar and reports diagnostics without writing Cue JSON unless explicitly requested.

#### `inspect`

```text
cueson inspect INPUT
```

Prints a concise human-readable summary by default.

`--json` emits a machine-readable summary to stdout.

No source bytes are dumped unless explicitly requested.

#### `schema`

```text
cueson schema
cueson schema --version
cueson schema --output PATH
```

Prints or writes the embedded canonical schema.

`cueson schema --version` prints only the schema version.

#### `version`

```text
cueson version
```

Prints the application version.

A verbose version option MAY include commit and target information without exposing local build paths.

#### `completion`

Generates shell completion for supported shells when the selected CLI framework provides it safely.

## Go architecture

Cueson SHOULD remain pure Go and build with `CGO_ENABLED=0` unless a future format introduces a justified native dependency.

The parser architecture SHOULD use a codec registry.

A conceptual codec contract is:

```go
type Codec interface {
    Format() Format
    Detect(sample []byte, file_name string) Detection
    Decode(ctx context.Context, source SourceBundle, opts DecodeOptions) (*Document, error)
    Render(ctx context.Context, doc *Document, opts RenderOptions) (*RenderResult, error)
}
```

Exact restoration is intentionally not a codec operation because it restores source assets directly from the source envelope.

Initial codecs:

```text
internal/codec/srt
internal/codec/webvtt
```

Future bitmap codecs MUST integrate through an OCR abstraction rather than embedding OCR logic directly into each parser. The anticipated boundary is an internal provider interface under `internal/ocr`, with codec-specific image extraction feeding one or more OCR providers and normalized OCR observations returning to the common model. The provider contract is intentionally deferred until the first bitmap-format feature slice, but the repository architecture MUST reserve this separation.

The common model lives under:

```text
internal/model
```

Schema handling lives under:

```text
internal/schema
```

Source asset integrity, filesystem metadata capture, native timestamp restoration, and exact source restoration live under:

```text
internal/source
```

Operating-system-specific timestamp adapters SHOULD be isolated behind an internal interface and implemented in platform-selected Go files such as `_windows.go`, `_darwin.go`, and `_linux.go`. Platform behavior MUST be covered by native CI tests rather than inferred from successful cross-compilation alone.

Loss reporting and cross-format orchestration live under:

```text
internal/convert
```

CLI wiring lives under:

```text
cmd/cueson
internal/cli
```

No package under `internal/` is a public compatibility promise.

### Schema embedding

The current canonical schema MUST be embedded with Go's `embed` facility.

The release build MUST derive the software version from the release tag or an equivalent single build-version source.

CI MUST verify:

```text
binary version == embedded schema_version == release tag version
```

A mismatch fails the build.

### Dependency policy

Prefer the Go standard library.

Third-party dependencies are acceptable where they materially reduce correctness risk, particularly for:

- CLI command parsing;
- JSON Schema Draft 2020-12 validation;
- legacy text encodings;
- terminal detection if the CLI framework does not provide it.

Every direct dependency MUST be pinned through `go.mod`, justified by actual use, and covered by automated vulnerability checks.

Avoid dependencies for functionality that is small, security-sensitive, or central to source fidelity when implementing the behavior directly is clearer.

## Repository scaffold

The repository SHOULD begin with this shape:

```text
cueson/
├── .editorconfig
├── .gitattributes
├── .gitignore
├── .github/
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug.yml
│   │   ├── feature.yml
│   │   └── config.yml
│   ├── workflows/
│   │   ├── ci.yml
│   │   ├── codeql.yml
│   │   ├── codex-review-gate.yml
│   │   ├── github-format.yml
│   │   ├── issue-project-reconcile.yml
│   │   ├── pr-issue-link.yml
│   │   └── release.yml
│   ├── dependabot.yml
│   └── pull_request_template.md
├── .specify/
│   ├── memory/
│   │   └── constitution.md
│   └── ...
├── cmd/
│   └── cueson/
│       └── main.go
├── docs/
│   ├── architecture.md
│   ├── cli.md
│   ├── formats/
│   │   ├── srt.md
│   │   └── webvtt.md
│   ├── project-management.md
│   ├── release-process.md
│   └── schema.md
├── internal/
│   ├── cli/
│   ├── codec/
│   │   ├── srt/
│   │   └── webvtt/
│   ├── convert/
│   ├── model/
│   ├── schema/
│   ├── source/
│   └── version/
├── schema/
│   ├── cueson.schema.json
│   └── releases/
│       └── v0.0.0/
│           └── cueson.schema.json
├── scripts/
│   └── github-format/
│       ├── go.mod
│       ├── main.go
│       └── main_test.go
├── specs/
├── testdata/
│   ├── srt/
│   ├── webvtt/
│   ├── roundtrip/
│   └── malformed/
├── AGENTS.md
├── CHANGELOG.md
├── CONTRIBUTING.md
├── LICENSE
├── NOTICE
├── README.md
├── SECURITY.md
├── go.mod
├── go.sum
└── .goreleaser.yml
```

The exact Spec Kit-generated files beneath `.specify/` are owned by the installed Spec Kit version and MUST NOT be invented manually.

The repository MUST initialize Spec Kit with the official `specify` CLI and then commit its project artifacts.

The repository MUST include a project constitution before feature implementation begins.

The repository MUST also include a standalone Go-based `scripts/github-format` helper from the bootstrap phase. The helper is repository-owned publication tooling with its own module and tests, not an application runtime dependency. It MUST support repository checking and a `-stdin` mode suitable for normalizing GitHub-bound Markdown before an issue, pull request, review, or release-note body crosses a shell or API boundary. Before the product Go module exists, invoke it with `go run ./scripts/github-format/main.go`; later repository automation MAY use an equivalent verified invocation.

The initial scaffold does not include the post-v1 documentation website package. The `site/` package and its deployment workflow are added only after v1 has been released successfully.

## Spec Kit integration

Cueson MUST use GitHub Spec Kit for spec-driven feature slices.

The repository initialization workflow MUST:

1. install or invoke the official `specify` CLI;
2. initialize the current repository with the installed CLI's supported non-interactive command and the integration appropriate to the active coding-agent environment;
3. commit the generated `.specify/` project structure;
4. create `.specify/memory/constitution.md`;
5. capture the Cueson invariants from this specification in the constitution;
6. keep feature artifacts under `specs/`;
7. use the full Spec Kit quality path for production features.

The expected full sequence is conceptually:

```text
constitution
specify
clarify
plan
checklist
tasks
analyze
implement
converge
```

The exact command spelling MUST follow the installed integration rather than being hardcoded into the repository.

### `shruggie-speckit` requirement

When a coding environment has access to the ShruggieTech `shruggie-speckit` skill, any request to run a full feature slice, kickoff, or autopilot sequence MUST use that skill.

`AGENTS.md` MUST state this as a repository rule.

The skill MUST orchestrate the repository's installed Spec Kit commands rather than recreating Spec Kit behavior.

The `analyze` gate is blocking.

Verification must run to completion in the foreground.

The agent MUST halt before `git push`, tagging, or release publication unless the operator explicitly authorizes those actions.

Architecture-affecting decisions MUST be recorded in the slice artifacts and detailed changelog.

## v1 delivery posture

Reaching v1.0.0 is an urgent project objective.

Agents SHOULD exploit concurrency aggressively when work can be partitioned safely. Appropriate parallel work includes independent format fixtures, parser research, schema conformance tests, platform-specific timestamp verification, documentation updates, fuzz targets, and other tasks with clearly separated ownership.

Spec Kit task generation SHOULD identify safe parallel tasks explicitly. A coordinating agent MAY delegate those tasks to sub-agents and integrate their results.

Parallel execution MUST NOT weaken repository discipline:

- two agents must not edit the same files concurrently without an explicit merge plan;
- each sub-agent receives bounded scope, expected outputs, and verification requirements;
- shared schema and architecture changes remain centrally coordinated;
- parallel work must converge through the same `analyze`, formatting, test, review, and release gates as serial work;
- urgency never authorizes skipping tests, unresolved review findings, security checks, release verification, or the pre-push authorization halt.

When parallel work produces conflicting interpretations of the specification, the coordinating agent MUST resolve the conflict against the constitution and architecture of record before integration.

### Work-slice recommendation policy

When the operator asks the agent to recommend or assemble the next work slice, the agent MUST inspect the active GitHub planning surface rather than assuming the next open issue should become the next slice.

Unless the operator has directed otherwise, or prior decisions in the active conversation require a narrower boundary, the agent SHOULD attempt to retire as many compatible active GitHub Issues as can be completed coherently in one slice.

One GitHub Issue is not equivalent to one work slice.

The preferred slice is the largest coherent set of active issues that can share one implementation narrative, one review surface, and one verification story without becoming difficult to reason about.

Issues are good candidates for the same slice when several of these are true:

- they belong to the same release milestone or immediate delivery objective;
- they affect the same subsystem, schema boundary, codec, CLI surface, or release mechanism;
- they share prerequisites or naturally depend on the same foundational change;
- one implementation can satisfy multiple issue outcomes without artificial duplication;
- their acceptance criteria can be verified together without hiding individual failures;
- their combined change remains reviewable and has a bounded blast radius.

The agent SHOULD split or defer issues when grouping would create a confused scope, mix unrelated architectural concerns, combine incompatible release or verification lifecycles, create unsafe concurrent ownership, materially increase rollback difficulty, or produce a review surface too broad to validate confidently.

Atomic issue discipline remains in force. Every included issue retains its own acceptance criteria and closure evidence even when several issues share one slice.

A slice recommendation SHOULD identify:

- the proposed slice title and purpose;
- every issue proposed for inclusion;
- why those issues belong together;
- relevant active issues deliberately excluded and the reason for exclusion;
- known dependencies and blockers;
- safe opportunities for sub-agent parallelism;
- the expected verification boundary.

The agent MUST preserve explicit operator scoping decisions. A prior instruction to keep an issue separate, defer a topic, or limit the slice takes precedence over the default preference to clear more active issues.

## Project constitution requirements

`.specify/memory/constitution.md` MUST establish at least these principles:

1. **Lossless source preservation:** no supported encode operation may destroy the ability to restore original source bytes.
2. **Schema and software lockstep:** application and canonical schema versions are identical.
3. **No filesystem provenance:** source paths never enter Cue JSON.
4. **Common model plus native fidelity:** normalized cross-format fields never replace required native details.
5. **No silent loss:** ignored, unsupported, malformed, or non-representable information is preserved where possible and surfaced through diagnostics.
6. **Test-first format work:** parser and renderer changes require fixture and round-trip tests.
7. **Spec-driven development:** non-trivial features pass through Spec Kit artifacts and analysis.
8. **Portable releases:** supported release targets remain first-class and release failures on one supported platform block release.
9. **Documentation is contractual:** CLI, schema, format grammar, compatibility, and release behavior must match implementation.
10. **GitHub-native delivery:** issues, milestones, Project fields, pull requests, and release evidence follow the repository project-management contract.
11. **Schema naming consistency:** Cueson-owned JSON keys use lowercase `snake_case`, with reusable semantic objects for repeated data types.
12. **Filesystem metadata preservation:** source timestamp metadata is captured before source reads, preserved with machine precision and provenance, and restored when the destination platform permits it.
13. **Urgent but verified delivery:** reaching v1 is time-sensitive; safe parallelization is encouraged, but no quality, review, security, or release gate may be bypassed for speed.
14. **Coherent issue consolidation:** a work slice may and generally should close multiple compatible issues when doing so preserves a clear implementation and verification boundary; issue count alone never defines slice size.
15. **Human merge authority:** final pull-request merge authority belongs to the human operator unless one explicit, unambiguous, single-merge override is granted.

## `AGENTS.md` requirements

`AGENTS.md` is a binding operating contract for AI agents working in the repository.

It MUST include the following rules.

### Markdown line wrapping ban

Agents MUST NOT hard-wrap Markdown prose.

Markdown prose uses one logical source line per paragraph and one logical source line per list item regardless of length.

Agents MUST NOT wrap text merely to satisfy 80-column, 100-column, or 120-column aesthetics.

The rule applies broadly to:

- README files;
- specifications;
- plans;
- task files where the generating framework allows it;
- issue and pull-request templates;
- changelog entries;
- release documentation;
- architecture documents;
- contribution documentation.

Agents MUST NOT reflow existing unwrapped Markdown into wrapped lines.

Fenced code, tables, generated machine-owned files, and formats whose generators control line layout are handled according to their native rules. Agents must not hand-mangle such structures to imitate wrapped prose.

### Spec Kit and skill rule

Full feature-slice execution MUST use the repository's installed Spec Kit workflow.

When available, `shruggie-speckit` is required for full-slice autopilot.

Agents may not invent substitute `/speckit-*` behavior when the installed project commands are absent.

For v1-bound work, agents SHOULD parallelize independent sub-tasks through sub-agents whenever file ownership and verification can be isolated safely. The coordinating agent remains responsible for integration and final verification.

### Work-slice assembly rule

Before proposing a work slice, the agent MUST inspect active issues, milestones, dependencies, and prior operator decisions.

By default, the agent SHOULD bundle as many compatible active issues as can be completed and verified coherently. The agent MUST NOT assume one issue equals one slice.

A larger issue count is not itself a virtue. The grouping must preserve a clear purpose, bounded change surface, individual issue traceability, and a verification plan that can prove each included outcome separately.

When presenting a slice recommendation to the operator, the agent SHOULD explain both the issues included and the nearby active issues intentionally left out.

### Human-only merge rule

An AI agent MUST NOT perform the final pull-request merge unless the human operator gives explicit and unambiguous authorization to merge that specific pull request.

This prohibition includes:

- clicking or invoking a GitHub merge action;
- running `gh pr merge`;
- invoking a GitHub API merge endpoint;
- enabling auto-merge or a merge queue entry on behalf of the operator;
- using another automation surface whose effect is to complete the merge.

General instructions such as "autopilot this slice", "finish the work", "get it ready", or prior permission to merge another pull request do not grant merge authority.

A merge override is valid for exactly one specifically identified pull request. It does not become a standing repository permission and MUST NOT be reused for a later pull request.

If material pull-request state changes after authorization in a way that would require a new substantive judgment, the agent SHOULD treat the earlier merge authorization as exhausted and return to the operator rather than guessing.

Absent an explicit override, the agent's terminal state for implementation is a merge-ready pull request with completed verification and review evidence.

### Post-merge housekeeping

"Housekeeping" is the bounded repository reconciliation work that begins after the human operator confirms that a pull request has been successfully merged.

Operator confirmation is the trigger, but the agent MUST still verify the actual GitHub merge state before deleting branches, worktrees, or other references.

Housekeeping is not permission to perform new feature work, publish a release, create a tag, deploy production infrastructure, or make unrelated repository changes.

The normal housekeeping sequence is:

1. verify the pull request is merged, record its base branch, source branch, reviewed head revision, and linked issues, and stop if the observed GitHub state contradicts the operator's statement;
2. fetch from the canonical remote and prune stale remote-tracking references;
3. synchronize the local default branch with the remote using a non-destructive fast-forward path when the working tree is clean and local history is not divergent;
4. inspect local branches and worktrees created for the merged slice;
5. remove only branches and worktrees proven stale by the completed merge;
6. verify the repository's automatic remote head-branch deletion occurred, and delete a lingering same-repository source branch only when its identity and staleness are proven;
7. prune obsolete worktree metadata and remote-tracking references;
8. verify required post-merge CI on the default branch and report any regression;
9. reconcile linked issues, parent epics, milestones, dependencies, and GitHub Project `Stage` values against the merged result;
10. report exactly what was removed, retained, reconciled, and left pending.

A branch or worktree is considered stale only when the agent can tie it to completed work and prove that deleting it cannot discard unmerged work.

For local branches, a merged pull request is strong evidence, but the agent MUST also verify that the local branch has not advanced beyond the reviewed pull-request head. Squash merges require care because ordinary `git branch --merged` ancestry may not recognize the topic branch as merged.

For remote branches, the agent MUST NOT delete a protected branch, a reused branch, a branch from another contributor or fork, or a branch whose current head no longer matches the completed pull-request work.

Before removing a worktree, the agent MUST inspect it for tracked modifications, untracked files, and other local-only state. A dirty worktree is retained and reported instead of force-removed.

Housekeeping SHOULD use safe Git operations such as:

```text
git fetch --prune
git worktree list --porcelain
git worktree prune
git status
```

The agent MUST NOT use broad destructive cleanup commands such as `git clean -fdx`, `git reset --hard`, forced worktree removal, or indiscriminate branch deletion as routine housekeeping.

If the local default branch contains uncommitted changes or has diverged from the remote, the agent MUST preserve that state and report the obstacle instead of resetting it.

Issue and Project reconciliation is part of housekeeping. The agent SHOULD verify that issues named with closing references actually closed, intentionally referenced issues remain open when appropriate, parent-epic progress is truthful, and milestone state reflects the merged outcomes. A parent issue closes only when its own acceptance criteria are satisfied.

Post-merge CI is also part of housekeeping when the repository runs checks on the default branch. A red post-merge run is not ignored merely because pull-request CI was green. The agent SHOULD surface the failure and create or update the appropriate GitHub work item when project policy calls for it.

Release actions remain separately authorized. Housekeeping may identify that a release, deployment, domain change, or milestone ritual is now ready, but it MUST NOT perform those actions unless the operator separately authorizes them.

Housekeeping completes with a concise operator report covering:

- verified merge identity;
- default-branch synchronization result;
- local branches removed or retained;
- remote branches removed, automatically deleted, or retained;
- worktrees removed or retained;
- issue, epic, milestone, and Project reconciliation;
- post-merge CI status;
- remaining follow-up or release actions requiring authorization.

### Codex pull-request review protocol

The repository MUST use Codex code review for pull requests except Dependabot-originated pull requests, unless the human operator explicitly overrides the rule for a specific pull request.

The normal review lifecycle is limited to at most two Codex review rounds initiated by project automation or an AI coding agent.

**Round one** occurs automatically when an eligible pull request becomes reviewable.

Codex may react to the pull request with the eyeballs reaction (`👀`). This means Codex has acknowledged the review request and has started or queued review work. It is not approval, completion, or evidence that the pull request is clean.

If Codex completes the first review with a thumbs-up reaction (`👍`) and no actionable review findings, the review is considered clean. The agent MUST NOT request a second Codex review.

If the first review produces findings, the authoring agent MUST:

1. inspect every Codex finding;
2. either implement the required correction or provide a technically justified response;
3. resolve the associated review threads only after the underlying concern has actually been addressed;
4. rerun the relevant verification;
5. update the pull-request head;
6. confirm that all first-round Codex findings are resolved.

Only after those conditions are true may the project request round two with exactly one `@codex review` invocation.

Round two is a final re-review, not the beginning of an unbounded loop. If round two finds additional problems, the agent MUST address them and rerun verification, but MUST NOT request a third Codex review on its own.

The human operator may explicitly waive a required review, request an additional review, or otherwise override this protocol for a specific pull request. An agent MUST NOT infer an override from silence.

Dependabot pull requests are excluded from automatic Codex review unless the operator explicitly requests one.

Repository automation SHOULD make this protocol idempotent. It MUST prevent duplicate second-round requests, distinguish the first and second review rounds, and never generate a third automated `@codex review`.

### GitHub body publication integrity

Structured Markdown MUST survive publication to GitHub with its required block structure intact.

Pull-request, issue, review, and release-note bodies commonly become corrupted when multi-line Markdown crosses a shell or API boundary as an array, an improperly escaped argument, or a command-substitution result. A successful `gh` or API response does not prove the body rendered correctly.

Before publication, agents MUST:

1. construct the body as a UTF-8 Markdown document with intentional blank lines;
2. pass it through `go run ./scripts/github-format/main.go -stdin`;
3. publish from a file whenever practical.

For pull requests, prefer:

```text
gh pr create --body-file <file>
```

or the corresponding file-based body option for an update.

For raw API calls, prefer a JSON request body read from a file rather than interpolating multi-line Markdown into one command-line argument.

When PowerShell must hold a multi-line body in memory and the value may be an array of strings, the agent MUST reconstruct a single scalar explicitly with:

```powershell
-join [Environment]::NewLine
```

After publishing, the agent MUST read the body back from GitHub immediately and verify:

- headings remain separated correctly;
- blank lines survived;
- lists and nested lists retain their intended structure;
- task checkboxes remain valid;
- fenced code blocks remain intact;
- tables remain valid;
- physical line count is plausible for the authored source.

Paragraph wrapping and Markdown structure are different concerns. Paragraph prose remains one physical source line under this repository's no-hard-wrap rule, while headings, separate list items, table rows, fenced blocks, and blank separators require real newline characters.

### Release note rule

GitHub release notes are a highlights summary, not a duplicate changelog.

Release notes MUST:

- contain only major user-visible highlights, compatibility notices, and critical upgrade information;
- omit exhaustive commit lists and exhaustive change-category duplication;
- remain substantially shorter than the corresponding changelog section;
- end with a link to the full changelog for that exact tag.

The final release-notes line MUST use the tagged changelog URL:

```text
Full changelog: https://github.com/shruggietech/cueson/blob/vX.Y.Z/CHANGELOG.md
```

### Changelog rule

`CHANGELOG.md` is the authoritative detailed human-readable release history.

Agents MUST update `[Unreleased]` during implementation.

Architecture-affecting decisions MUST be called out explicitly.

The changelog uses Keep a Changelog structure and Semantic Versioning.

### Push, merge, and release halt

Agents MUST NOT:

- push;
- perform a pull-request merge;
- enable auto-merge or place a pull request into a merge queue;
- create or move a release tag;
- publish a GitHub release;
- publish a schema version to the production domain

without explicit operator authorization for the specific action.

Merge authorization is always single-use and pull-request-specific. Authorization to merge one pull request does not authorize any later merge.

### Source-fidelity rule

Agents MUST NOT "clean up" parser behavior by dropping unknown source content, normalizing raw source fields, or changing raw-source bytes without a specification change.

### Encoding and text hygiene

Repository-authored text files MUST use UTF-8 without BOM. Line endings are governed by `.gitattributes`: text defaults to LF, while PowerShell, command, and batch scripts use CRLF. `.editorconfig` MUST mirror those rules for supported editors.

Intentional non-UTF-8 or byte-sensitive fixtures MUST be isolated under `testdata/`, excluded from Git text normalization through `.gitattributes`, and documented. Approved brand, font, and other integrity-tracked asset directories MUST likewise be protected from text normalization.

## GitHub-native project management

Cueson MUST use an organization-owned GitHub Project named:

```text
cueson Delivery
```

The project is linked to the canonical ShruggieTech repository.

GitHub is the project-management source of truth. No parallel ticketing or planning database is required.

### Repository settings and code-quality controls

Repository settings are part of the project scaffold and MUST be configured deliberately rather than left at platform defaults.

At minimum, the repository SHOULD enable every applicable GitHub code-quality and security feature available to the project plan and repository visibility, including:

- dependency graph;
- Dependabot alerts;
- Dependabot security updates;
- CodeQL code scanning;
- secret scanning;
- secret-scanning push protection;
- private vulnerability reporting when available and appropriate.

The default branch MUST be protected by a repository ruleset or equivalent branch protection. The baseline policy SHOULD:

- require pull requests for ordinary changes;
- require the defined CI and quality checks to pass;
- require review conversations to be resolved before merge;
- block force pushes and branch deletion on the protected default branch;
- preserve administrator or operator override capability for exceptional recovery.

Squash merge MUST be available as the normal pull-request merge path. The repository's native automatic head-branch deletion setting MUST be enabled so merged feature branches are removed after successful merge. A custom branch-deletion workflow SHOULD NOT be created when GitHub's native setting provides the behavior reliably.

GitHub Actions workflow permissions SHOULD default to the least privilege necessary. Workflows that need write permissions MUST declare those permissions explicitly at the narrowest practical scope.

Codex automatic review MUST be enabled for eligible non-Dependabot pull requests. If the native integration cannot enforce the project's two-round protocol by itself, repository automation MUST provide the missing gating logic without creating an unbounded review loop.

Repository settings that materially affect delivery or security SHOULD be verified through the GitHub API or CLI after mutation. Agents MUST NOT treat a successful settings command as sufficient evidence without reading the resulting state back.

### One authoritative location for each fact

The repository follows this mapping:

| Fact | Authoritative location |
|---|---|
| Outcome, scope, acceptance criteria, discussion | Issue |
| Responsible person or agent | Native assignee |
| Type, priority, effort, area, gates | Labels or approved organization issue types |
| Target release | Native milestone |
| Parent-child hierarchy | Native parent and sub-issue relationships |
| Blocking order | Native issue dependencies |
| Implementation evidence | Pull request with closing reference |
| Delivery lifecycle | Project `Stage` |
| Cross-issue implementation batch | Project `Slice` |
| Operating rules | Repository docs, `AGENTS.md`, Project README |

Native metadata MUST NOT be duplicated into custom Project fields.

### Atomic issues

Every actionable issue owns one independently closeable and independently testable outcome.

Atomicity defines the issue contract, not the work-slice size. Several atomic issues may share one work slice when their implementation and verification form one coherent delivery unit.

Implementation and later release/field verification MUST use separate issues when they have independent completion evidence.

Coordinator epics MAY group atomic children.

### Work-slice planning

When assembling the next slice, the agent SHOULD inspect all active issues relevant to the current milestone or delivery objective and attempt to include the largest coherent compatible set.

The default objective is to reduce active issue inventory efficiently without sacrificing clarity.

The agent MUST NOT mechanically assign one slice per issue.

The agent MUST preserve individual issue acceptance criteria inside a multi-issue slice and MUST NOT close an issue merely because neighboring issues in the same slice completed.

When a potentially compatible issue is excluded, the slice plan SHOULD record why, such as dependency ordering, architectural separation, excessive review scope, distinct platform verification, operator deferral, or different release timing.

### Issue contract

Actionable issue templates MUST contain:

```text
Outcome
Context
Scope
Acceptance criteria
Dependencies
Verification
```

Every issue also receives the governed type, priority, effort, area, and milestone metadata when evidence supports those values.

### Labels

If native organization issue types are unavailable, the label families SHOULD be:

```text
type:
priority:
effort:
area:
needs:
skip:
```

Exclusive families MUST remain mutually exclusive.

Workflow state MUST NOT be represented with labels such as `todo`, `doing`, `review`, or `done`.

Release identity MUST NOT be represented as a label. Milestones own releases.

### Milestones

Published versions and committed upcoming versions use native milestones such as:

```text
v0.0.0
v0.1.0
v1.0.0
```

Milestones represent delivery commitments, not workflow stages.

A milestone MUST NOT close until its release or verification requirements are complete.

### Project fields

The custom `Stage` field MUST use exactly:

```text
Backlog
Ready
Specced
In progress
PR review
Release verification
Done
```

The custom `Slice` text field is the only default additional planning field.

GitHub's default `Status` field is left unused rather than mirrored.

### Project invariant

Every in-scope repository issue appears exactly once in the Project.

Committed work uses repository issues, not draft Project items.

### Pull requests

Normal pull requests to the default branch MUST contain at least one complete closing reference:

```text
Closes #123
```

Cross-repository references use the full repository identifier.

CI MUST enforce the issue-link rule, with only documented `skip:` exceptions.

Eligible pull requests also pass the Codex review protocol defined by `AGENTS.md`. The review gate MUST distinguish a clean first-round thumbs-up from a first round with findings, allow at most one automated second-round request, and exclude Dependabot unless explicitly overridden.

The final merge is a human-controlled boundary. An agent may prepare, update, verify, review, and report a pull request, but MUST NOT merge it without a one-time explicit operator instruction naming or otherwise unambiguously identifying the pull request.

After the operator confirms a successful merge, the agent enters the post-merge housekeeping protocol defined in `AGENTS.md`.

### Reconciliation

Project automation MUST be evidence-based and idempotent.

Automation MAY safely:

- add missing repository issues;
- trigger or reconcile the first Codex review on eligible non-Dependabot pull requests;
- request exactly one second Codex review only after first-round findings are resolved;
- verify GitHub-published Markdown bodies after publication;
- initialize empty stages to `Backlog`;
- set closed issues to `Done`;
- set linked open pull requests to `PR review`;
- set open `needs: verification` issues to `Release verification`;
- resolve duplicate Project items;
- archive old completed items according to policy.

Automation MUST NOT infer priority, effort, milestone, readiness, slice, or verification completion without authoritative evidence.

Automated mutations MUST be read back and verified.

### Audit cadence

A recurring audit SHOULD verify:

- missing or duplicate Project items;
- invalid stages;
- open issues marked `Done`;
- closed issues not marked `Done`;
- `PR review` without an open pull request;
- `In progress` without current ownership or activity;
- gate labels on closed issues;
- multiple labels from exclusive families;
- unsupported milestone assignments;
- epic-child inconsistencies;
- inferred or unknown slices;
- merged pull requests whose source branches or slice worktrees remain stale;
- merged pull requests whose post-merge CI or parent-issue reconciliation remains incomplete.

## Changelog and release model

### Changelog

`CHANGELOG.md` uses a Keep a Changelog structure.

It begins with:

```markdown
# Changelog

## [Unreleased]

## [0.0.0]
```

The initial repository version is `0.0.0`.

Detailed changes belong in `CHANGELOG.md`, not GitHub release notes.

Recommended categories include:

```text
Added
Changed
Deprecated
Removed
Fixed
Security
Decisions
```

`Decisions` is a ShruggieTech extension for architecture-affecting choices that future agents need to understand.

### Release notes

Release notes summarize only the most important highlights. They are intentionally concise and MUST NOT become a second detailed changelog.

They MUST end with:

```text
Full changelog: https://github.com/shruggietech/cueson/blob/vX.Y.Z/CHANGELOG.md
```

The release pipeline SHOULD verify that the line exists and points to the release tag before publication.

### Semantic versioning

At and after v1.0.0, breaking changes require a major-version bump. Before v1.0.0, documented breaking changes may occur in a minor release while patch releases remain non-breaking.

Schema and software remain in lock step.

The v1.0.0 milestone is not complete until the SRT and WebVTT compatibility gates in this document are satisfied.

## Initial version roadmap

The pre-v1 roadmap is intentionally execution-focused. Agents should prefer coherent parallel workstreams over unnecessary serialization when tasks are independent, while preserving all integration and quality gates.

### Initial repository bootstrap

The first bootstrap session establishes enough trustworthy structure to plan product implementation. It does not need to deliver the complete `0.0.0` milestone in one change.

Required bootstrap outcomes:

- local Git repository and text hygiene rules, including `.gitattributes`;
- Apache 2.0 license and notice;
- a truthful pre-release README;
- official Spec Kit initialization and a ratified initial constitution;
- repository `AGENTS.md` and `[Unreleased]` changelog;
- standalone `scripts/github-format` publication tooling;
- GitHub issue and pull-request templates;
- the canonical GitHub repository, delivery Project, labels, milestones, and atomic foundation backlog.

Product Go code, the root Go module, schema implementation, CI enforcement, repository rulesets, security automation, review automation, and release dry runs follow as tracked work. Required checks MUST NOT be configured before their workflows exist.

### `0.0.0` repository foundation

`0.0.0` establishes the project without claiming production format completeness.

Required outcomes:

- repository scaffold;
- Go module;
- CLI skeleton;
- version command;
- embedded schema plumbing;
- canonical `0.0.0` schema;
- multi-asset source envelope and integrity validation;
- public generic exact restoration from a valid source envelope;
- Spec Kit initialization;
- project constitution;
- `AGENTS.md`;
- `CHANGELOG.md`;
- GitHub issue and pull-request templates;
- GitHub Project contract;
- required GitHub repository quality/security settings;
- automatic eligible-PR Codex first-round review;
- two-round Codex review gate scaffolding;
- `scripts/github-format` publication helper;
- CI;
- cross-platform build proof;
- release pipeline dry run;
- schema/software version-equality test;
- initial test fixture structure.

### `0.x` implementation series

Pre-v1 minor releases may add compatible capabilities or introduce documented breaking contract changes. Patch releases remain non-breaking. Released schema artifacts themselves remain immutable.

Expected work includes:

- native format ingest into the source asset envelope;
- codec-integrated integrity validation;
- SRT parser;
- SRT renderer;
- WebVTT parser;
- WebVTT renderer;
- format detection;
- normalized cue model;
- diagnostics;
- codec-to-envelope-to-restore round trips;
- cross-format conversion;
- loss reporting;
- corpus hardening;
- fuzzing;
- schema release packaging;
- documentation;
- verification that the externally produced Cueson brand kit has been imported into the repository before v1.

The exact minor-version allocation is managed through milestones and Spec Kit slices rather than frozen in this architecture document.

The Cueson brand kit is produced outside the coding-agent workflow and will be made available through `https://brand.shruggie.tech`. Coding agents MUST NOT invent, redesign, or substitute temporary brand assets. Once the operator or external brand process adds the approved kit to the repository, agents may integrate and verify it. Absence of the approved brand kit is a v1 release blocker, not permission for an agent to create one.

### `1.0.0` stability gate

v1.0.0 requires:

- full documented SRT structural coverage;
- full documented WebVTT structural coverage;
- byte-exact restore for every accepted SRT and WebVTT fixture;
- no original filesystem path in generated Cue JSON;
- common cue model validated across both formats;
- format-specific information preserved;
- SRT-to-WebVTT conversion;
- WebVTT-to-SRT conversion;
- strict-mode loss prevention;
- canonical schema included in the tagged repository state and release artifacts, ready for post-v1 public-domain publication;
- approved externally produced brand kit present in the repository;
- schema and executable version lockstep;
- supported release binaries for Windows, macOS, and Linux;
- complete CLI reference;
- complete schema reference;
- complete format grammar documentation;
- green fuzz, corpus, integrity, and cross-platform CI gates.

## Schema-first future format declarations

Cueson SHOULD allow the schema to define structures for future formats before the implementation claims native ingest support for those formats.

The purpose is not to promise functionality prematurely. The purpose is to let maintainers and the broader community debate, test, and refine document structure in the open before codec work hardens around it.

The recommended process is:

1. introduce or refine the schema shape for a future format under a truthful `format_support.status` such as `reserved`, `experimental`, or `envelope_only`;
2. document what parts of the structure are normative and what parts remain provisional;
3. publish example Cue JSON fixtures using that structure;
4. solicit discussion and revise the schema as allowed by the active version policy;
5. implement codecs only after the structural contract is sufficiently mature.

This model lets Cueson become a long-haul interchange contract for subtitle data rather than merely a mirror of whatever codecs happen to exist this month.

Schema-recognized future formats MUST still satisfy the repository's ordinary invariants: no source paths, stable naming style, exact source-asset preservation, explicit diagnostics, and truthful support-level reporting.

## Future format roadmap

The v1 architecture MUST avoid assumptions that prevent these families from being added later. For image-based families, future format support is defined as both lossless source preservation and usable semantic extraction. OCR is part of that future support contract, not a separate unrelated product.

### XML and timed-text families

Planned candidates:

- TTML;
- IMSC;
- SMPTE-TT;
- EBU-TT.

These formats introduce namespaces, styling, layout, timing expressions, metadata, and document structures that are richer than SRT and often richer than WebVTT.

Their format-specific data SHOULD remain native rather than being flattened destructively into the common cue model.

### Scripted subtitle families

Planned candidates:

- ASS;
- SSA.

These formats introduce script sections, styles, events, overrides, positioning, and karaoke timing.

The common cue model may expose normalized dialogue events while format-specific data preserves script structure and override syntax.

### Legacy broadcast formats

Planned candidate:

- EBU-STL.

This introduces binary records, broadcast-specific character sets, timing conventions, and metadata.

The existing source-asset base64 envelope is directly reusable.

### Bitmap subtitle families

Planned candidates:

- PGS/SUP;
- VobSub IDX/SUB.

These formats are image-based rather than text-based.

The schema will require binary source assets and MUST support independently addressable subtitle images or events, timing metadata, palettes or rendering information where applicable, and derived OCR observations.

OCR is a required semantic capability for full Cueson support of bitmap subtitle families. A codec that can only preserve and restore PGS/SUP or VobSub bytes is an archival envelope, not a complete Cueson implementation of that format. Experimental envelope-only support MAY exist during development, but release documentation MUST label it explicitly as incomplete.

OCR remains derived. It MUST NOT replace, rewrite, or become the sole authoritative representation of the bitmap source. Exact restoration continues to use the original source assets.

The multi-asset source model is specifically intended to support paired VobSub files.

### SAMI

SAMI is HTML-like and may contain style, class, language, and timing structures.

Its source fidelity belongs in format-specific data while normalized cues expose supported timed text.

## Base64 and OCR

Base64 is part of the source envelope from v1 for exact restoration.

Future binary subtitle support therefore does not require a second transport mechanism.

Future schema versions MAY add per-event or per-frame binary payloads where independently addressable binary objects are useful. Such payloads SHOULD use explicit media types, hashes, and base64.

The common model MUST reserve a durable derived-data mechanism capable of carrying OCR observations without requiring a breaking redesign of the root document. The exact v1 shape may remain empty for SRT and WebVTT, but its extension path must be defined before v1.0.0 is frozen.

A future OCR observation SHOULD be able to describe:

```text
id
engine
engine_version
model
language
text
lines
confidence
source_asset_id
source_event_id
source_image_id
region
processing_options
```

Confidence SHOULD support both an overall observation value and, when the engine exposes them, finer-grained line, token, or region confidence values.

OCR implementations SHOULD permit engine abstraction rather than baking one OCR implementation permanently into the Cue JSON contract. The software may ship with a preferred built-in or bundled engine in a future release, but engine identity belongs in provenance so results remain interpretable and replaceable.

For bitmap subtitle formats, OCR text becomes the primary normalized textual representation for search, indexing, analysis, and cross-format conversion. It remains derived data and MUST never become the sole representation of the subtitle event or replace the original source bytes and images.

Cross-format conversion from a bitmap subtitle format into a text subtitle format MUST use explicit OCR-derived text and MUST report that OCR derivation in conversion provenance or diagnostics. Conversion MUST NOT pretend OCR output was native source text.

## Downstream consumer compatibility

Cueson is expected to become a foundational interchange layer for other subtitle- and transcription-oriented software.

A future downstream consumer MAY:

- inspect media containers for embedded subtitle tracks;
- extract embedded subtitle assets;
- invoke Cueson codecs for supported extracted formats;
- perform basic speech transcription when no usable subtitle track exists;
- present the resulting timed content as Cue JSON.

That downstream media inspection, demuxing, and speech-to-text workflow is not part of the Cueson core executable's v1 responsibility. Cueson SHOULD instead keep clean boundaries so such software can depend on the canonical schema and CLI without forking the document model.

Cue JSON authored by another compliant producer MUST identify that producer truthfully in the `producer` object. Derived transcription or OCR text MUST retain provenance indicating that the text was generated rather than native source text.

Cueson MUST NOT acquire a hard dependency on a media demuxer, video runtime, or speech-recognition model merely to support this downstream use case. If a stable public Go library becomes valuable for downstream consumers, that API requires its own explicit compatibility specification rather than leaking `internal/` packages into public use.

## Cross-platform release requirements

Official releases MUST provide native binaries for:

- Windows;
- macOS;
- Linux.

The initial release matrix SHOULD include:

```text
windows/amd64
windows/arm64
darwin/amd64
darwin/arm64
linux/amd64
linux/arm64
```

A target may be removed only through an explicit compatibility decision and changelog entry.

Release builds SHOULD use `CGO_ENABLED=0`.

Artifacts SHOULD include:

- platform binary archives;
- SHA-256 checksum manifest;
- canonical schema;
- license and notices as required;
- SBOM where supported by the release toolchain;
- provenance or artifact attestation where supported by GitHub.

GoReleaser is the preferred release orchestrator unless implementation planning identifies a concrete reason to use a different tool.

## CI quality gates

Pull-request CI MUST run the following foundation gates:

```text
gofmt verification
go run ./scripts/github-format/main.go
go vet ./...
go test ./...
go test with race detection on an appropriate CI target
static analysis
govulncheck
schema validation tests
schema key-style conformance tests
source-path prohibition tests
version-lockstep tests
pull-request body formatting tests
```

Capability-specific gates activate only after their owning implementation exists. Native format round trips and cross-format conversion tests begin with their codec and conversion slices. Codex review-gate automation tests begin with issue #9. CI MUST NOT publish empty or falsely passing placeholders for these deferred surfaces.

The foundation CI MUST run platform-selected tests natively on Windows, macOS, and Linux and MUST separately cross-build the supported Windows, macOS, and Linux amd64/arm64 matrix with `CGO_ENABLED=0`. Cross-build success is portability evidence and MUST NOT be presented as native behavioral proof.

A pinned `golangci-lint` configuration MAY provide the static-analysis umbrella.

CodeQL SHOULD run independently as a GitHub security workflow.

Tests MUST run in the foreground in agent workflows. Agents must not launch long test suites in the background and infer success from polling.

## Testing strategy

### Unit tests

Every parser and renderer component receives focused unit tests for its grammar and transformations.

### Golden fixtures

Every supported source dialect requires source fixtures and expected Cue JSON fixtures.

Golden tests MUST include:

- source input;
- expected normalized model;
- expected diagnostics;
- expected rendered output where canonical rendering is defined;
- expected exact-restoration SHA-256.

### Exact round-trip tests

For every accepted source fixture:

```text
source bytes
-> cueson encode
-> Cue JSON
-> cueson restore
-> restored bytes
```

MUST result in identical bytes and identical SHA-256.

Where the test platform and filesystem support restoring the captured timestamp set, the same test MUST also verify identical `created`, `modified`, and `accessed` instants at the source API's effective precision.

CI MUST include platform-native timestamp restoration tests on Windows, macOS, and Linux. Tests MUST verify both successful restoration and the documented unsupported-path behavior. Linux tests MUST NOT fake support for setting birth time when the underlying filesystem/API does not provide it.

This is a release-blocking invariant for byte restoration and a release-blocking capability-contract test for filesystem metadata restoration.

### Render round-trip tests

For supported, conforming inputs:

```text
source
-> encode
-> render to same format
-> re-encode
```

SHOULD produce semantically equivalent normalized models even when rendered lexical form becomes canonical.

This is distinct from exact restoration.

### Cross-format tests

SRT and WebVTT conversion tests MUST verify:

- representable content maps correctly;
- non-representable content appears in a loss report;
- `--strict` prevents lossy output;
- default mode warns but does not silently discard unsupported semantics.

### Malformed corpus

`testdata/malformed/` contains intentionally malformed and tolerated inputs.

Tests verify that parsers:

- preserve bytes;
- produce deterministic diagnostics;
- never panic;
- never silently discard preservable input.

### Fuzz testing

Go fuzz tests MUST target:

- format detection;
- SRT block parsing;
- SRT time parsing;
- WebVTT block parsing;
- WebVTT timing settings;
- markup/token scanning;
- base64/source-envelope validation;
- renderer/parser cycles.

Fuzz-found regressions MUST be committed as permanent fixtures.

### Corpus testing

A larger real-world corpus SHOULD be maintained outside the repository when licensing or size makes committing it inappropriate.

The CI-compatible fixture corpus remains small, redistributable, and deterministic.

A separate maintainers' corpus test command MAY run against locally supplied directories without encoding local paths into outputs.

## Security requirements

Subtitle files and Cue JSON are untrusted input.

Cueson MUST:

- never execute subtitle markup or scripts;
- never interpret embedded HTML-like content as executable content;
- enforce safe restoration basenames;
- prevent directory traversal on restore;
- cap or stream allocations where practical;
- avoid regex patterns vulnerable to catastrophic backtracking;
- handle malformed base64 without panic;
- validate claimed byte lengths and hashes;
- refuse unsafe overwrite without `--force`;
- avoid embedding source paths or local machine identifiers;
- validate timestamp ranges before passing them to native filesystem APIs;
- treat stored timestamps as untrusted input during restoration;
- fuzz parser boundaries;
- run dependency vulnerability scanning.

Future archive support, if added, requires a separate extraction-safety specification.

## Documentation requirements

The repository MUST maintain:

- `README.md` for product overview and quick start;
- `docs/architecture.md` for architecture of record;
- `docs/schema.md` for schema model and compatibility;
- `docs/cli.md` for complete CLI contract;
- `docs/formats/srt.md` for the supported SRT grammar and tolerated variants;
- `docs/formats/webvtt.md` for WebVTT coverage;
- `docs/project-management.md` for the GitHub-native operating contract;
- `docs/release-process.md` for versioning, schema publication, release verification, and notes;
- `CONTRIBUTING.md`;
- `SECURITY.md`;
- `AGENTS.md`;
- `CHANGELOG.md`.

Markdown in repository-authored documents MUST follow the no-hard-wrap rule defined by `AGENTS.md`.

Before v1, root `docs/` is the authoritative documentation source and is consumed directly from the repository. The public documentation website is a post-v1 deliverable and MUST NOT become a release blocker for the v1 software itself.

## README requirements

The bootstrap README SHOULD communicate:

- what Cue JSON is;
- why exact source restoration matters;
- the planned initial stable formats and current implementation status;
- the project roadmap and repository documentation;
- project domain;
- ShruggieTech attribution;
- contribution and security links.

It SHOULD follow the ShruggieTech repository header pattern with truthful CI, release, license, and documentation badges. Badges MAY point to planned repository surfaces before their workflows or releases exist only when the README clearly labels the project as pre-release.

Before v1.0.0, the README MUST add installation guidance, encode, restore, render, convert, and validate examples, the canonical schema identifier, and the final compatibility statement as those capabilities become real.

The README MUST NOT claim full SRT or WebVTT support before the v1.0.0 acceptance gate is met.

## License and legal scaffolding

Cueson is licensed under the Apache License, Version 2.0.

The repository MUST contain the complete Apache 2.0 license text in `LICENSE` from the initial scaffold and MUST use SPDX identifier:

```text
Apache-2.0
```

`NOTICE` MUST be present and maintained when required by Apache 2.0 attribution obligations or bundled third-party material.

Direct and bundled dependencies, fixtures, fonts, brand assets, and other redistributed material MUST be reviewed for license compatibility before release.

Third-party notices MUST be included when dependencies or fixture licensing require them.

Test fixtures MUST have a clear provenance and redistribution status.

## Release workflow

A release candidate must satisfy this sequence:

1. all milestone issues are resolved or truthfully moved;
2. all required verification issues are complete;
3. `CHANGELOG.md` has a complete version section;
4. schema version is updated in lock step with software version;
5. the immutable release schema is copied into `schema/releases/vX.Y.Z/`;
6. schema IDs reference the exact version URL;
7. format, lint, test, fuzz-regression, vulnerability, schema, and round-trip gates pass;
8. release binaries build for every supported target;
9. release artifacts include checksums and schema;
10. release notes are generated as highlights only;
11. release notes end with the exact tagged changelog link;
12. an agent or maintainer verifies artifacts after build;
13. operator authorization is obtained before tag push or publication;
14. the GitHub release is published;
15. the release artifact and tagged repository schema are verified against one another;
16. the milestone is closed only after v1 release verification succeeds.

Production `cueson.io` DNS, hosting, and schema publication are intentionally not part of the v1 release transaction. They occur in a separate post-v1 public-web activation phase after the software release has succeeded.

## Schema publication

The repository copy under:

```text
schema/releases/vX.Y.Z/cueson.schema.json
```

is the immutable release source of truth.

Each GitHub release MUST include the exact same schema artifact.

The schema continues to use its intended canonical `https://cueson.io/schema/vX.Y.Z/cueson.schema.json` identifier, but the public URL is not required to resolve until the post-v1 public-web activation phase. Consumers validating before public-domain activation can use the embedded schema, repository copy, or release artifact.

After the public domain is activated, released schemas are published beneath:

```text
https://cueson.io/schema/vX.Y.Z/cueson.schema.json
```

Publication automation MUST NOT mutate a released schema.

A post-publication verification job MUST fetch the public URL and compare its SHA-256 with the tagged repository and release artifact.

## Post-v1 public web and domain activation

Public web activation begins only after v1.0.0 has been released successfully and its release verification is complete.

The work is a separate Spec Kit slice or set of slices and MUST NOT be silently folded into the v1 release transaction.

### Domain configuration

`cueson.io` is already registered through Cloudflare.

The coding agent responsible for post-v1 web activation MUST use Cloudflare command-line or API-backed CLI tooling for necessary domain configuration. Dashboard-only manual mutation is not the expected implementation path.

The agent MUST:

- use a least-privilege Cloudflare API token rather than a global account key when automation credentials are required;
- inspect existing zone and DNS state before changing anything;
- make the smallest necessary DNS and domain changes for the selected static host;
- avoid deleting unrelated records;
- read the resulting Cloudflare state back after mutation;
- independently verify public DNS resolution and TLS behavior;
- record the final domain configuration in `docs/release-process.md` or the post-v1 site architecture document.

Cloudflare production mutations remain operator-authorized actions. The agent MUST halt before making them if the active authorization does not explicitly cover production domain configuration.

### Documentation site

After v1, the repository gains a static documentation/public-site package that serves authoritative content originating from root `docs/`.

The baseline house stack is:

```text
Next.js App Router
TypeScript
React
Fumadocs Core/UI and Fumadocs MDX
Tailwind CSS
next-themes
pnpm
Playwright
static export
```

Exact dependency versions are pinned during the post-v1 implementation slice according to the then-current approved house baseline.

The intended repository boundary is:

```text
docs/                    authored documentation source of truth
site/                    static Next.js/Fumadocs application
site/content/generated/  build-only adaptations of docs/
site/out/                build-only deployable static artifact
```

Generated website content MUST NOT become a second documentation authority. Root `docs/` remains the authored source.

The site SHOULD include:

- product landing content;
- documentation navigation;
- CLI documentation;
- schema documentation;
- versioned schema download links;
- supported-format documentation;
- release and changelog links;
- security and contribution links;
- Cueson brand assets already present in the repository.

The site build MUST include static-build verification, route checks, link checks, accessibility tests, responsive tests, metadata checks, and production-equivalent artifact verification.

Deployment MUST be owner-controlled and MUST NOT occur from arbitrary pull requests.

### Brand integration

The approved brand kit must already exist in the repository before v1 release. The post-v1 site consumes that approved kit.

Brand creation itself remains outside coding-agent scope. Agents MAY integrate, adapt through provided bindings, and verify approved assets but MUST NOT invent replacement branding.

## Implementation principles for AI agents

Agents working on Cueson MUST:

- inspect the architecture and constitution before implementation;
- use Spec Kit for non-trivial feature work;
- use `shruggie-speckit` for full-slice autopilot when available;
- parallelize independent v1 work through sub-agents when ownership boundaries are clear;
- treat the schema as code;
- add tests before or alongside parser changes;
- preserve unknown source data rather than deleting it;
- make the smallest compatible change;
- avoid broad refactors unrelated to the active atomic issue;
- update the detailed changelog for behavior or architecture changes;
- avoid hard-wrapping Markdown;
- run `go run ./scripts/github-format/main.go` before publishing repository-authored GitHub text;
- publish multi-line GitHub bodies from files whenever practical and read them back after publication;
- follow the two-round Codex review protocol for eligible pull requests;
- recommend multi-issue slices by default when active issues can be grouped coherently;
- never perform the final merge without a one-time explicit operator override;
- perform post-merge housekeeping only after operator confirmation and verified merge state;
- run CI-parity checks before claiming completion;
- read back GitHub mutations when performing project-management work;
- halt before push, tag, release, or production schema publication unless explicitly authorized.

## v0.0.0 foundation completion gate

The initial scaffolding effort is complete when:

- [ ] the `shruggietech/cueson` repository exists;
- [ ] the Go module builds;
- [ ] the project version is `0.0.0`;
- [ ] the schema version is `0.0.0`;
- [ ] the schema is embedded in the executable;
- [ ] the representative root document example includes explicit cue timing and subtitle text;
- [ ] the canonical schema defines the common cue payload baseline for subtitle text;
- [ ] the schema defines `format_support` and its baseline maturity values;
- [ ] `cueson version` reports `0.0.0`;
- [ ] `cueson schema --version` reports `0.0.0`;
- [ ] CI fails if those versions diverge;
- [ ] Spec Kit is initialized;
- [ ] the project constitution contains the required invariants;
- [ ] `AGENTS.md` contains the Markdown wrapping ban;
- [ ] `AGENTS.md` contains the `shruggie-speckit` requirement;
- [ ] `AGENTS.md` contains the release-notes highlights-only rule;
- [ ] `AGENTS.md` contains the multi-issue work-slice recommendation rule;
- [ ] `AGENTS.md` contains the human-only final merge rule and single-use override semantics;
- [ ] `AGENTS.md` defines the post-merge housekeeping protocol;
- [ ] `AGENTS.md` contains the pre-push/pre-release halt;
- [ ] `CHANGELOG.md` contains `[Unreleased]` and `[0.0.0]`;
- [ ] `LICENSE` contains Apache License 2.0 and repository metadata uses `Apache-2.0`;
- [ ] `NOTICE` exists;
- [ ] `scripts/github-format` exists and supports `-stdin`;
- [ ] issue templates follow the atomic outcome contract;
- [ ] the GitHub Project is configured with the required `Stage` and `Slice` fields;
- [ ] PR issue-link enforcement exists;
- [ ] eligible non-Dependabot pull requests receive automatic first-round Codex review;
- [ ] automation can request at most one second Codex review after first-round findings are resolved;
- [ ] GitHub body publication readback is part of the agent contract;
- [ ] GitHub code-quality and security settings are enabled as applicable;
- [ ] protected-default-branch ruleset is configured;
- [ ] automatic head-branch deletion after merge is enabled;
- [ ] cross-platform build jobs prove Windows, macOS, and Linux compilation;
- [ ] release configuration exists but does not auto-publish without authorization;
- [ ] repository-authored text files are UTF-8 without BOM, with line endings governed by `.gitattributes`;
- [ ] no generated Cue JSON field stores an original filesystem path;
- [ ] all Cueson-owned schema keys conform to the lowercase `snake_case` rule;
- [ ] source assets capture timestamp metadata before reading source bytes;
- [ ] platform timestamp capability tests exist for Windows, macOS, and Linux.

## v1.0.0 done gate

The v1 release MUST NOT ship until all of the following are true:

- [ ] every documented SRT grammar feature has parser fixtures;
- [ ] every documented SRT grammar feature has renderer fixtures where applicable;
- [ ] every documented WebVTT structure has parser fixtures;
- [ ] every documented WebVTT structure has renderer fixtures where applicable;
- [ ] all accepted SRT fixtures restore byte-for-byte;
- [ ] all accepted WebVTT fixtures restore byte-for-byte;
- [ ] source hashes are verified during restore;
- [ ] source `created`, `modified`, and `accessed` timestamps are captured before source reads;
- [ ] restorable timestamps are reapplied after byte verification;
- [ ] timestamp restoration is read back and verified;
- [ ] unsupported creation/birth-time restoration is reported explicitly rather than silently ignored;
- [ ] `--strict-metadata` enforces full timestamp reproduction when requested;
- [ ] SRT derived speaker detection never mutates raw cue text;
- [ ] WebVTT rolling cues are preserved one-to-one;
- [ ] non-cue WebVTT blocks retain total source order;
- [ ] no accepted source construct is silently discarded;
- [ ] generated documents validate against the embedded v1.0.0 schema;
- [ ] the common cue model exposes subtitle text through `payload.raw_text`, `payload.plain_text`, and `payload.lines` for SRT and WebVTT;
- [ ] text consumers do not need to decode `data_base64` to access cue text;
- [ ] `format_support` truthfully reports support maturity and capability state;
- [ ] public v1.0.0 schema matches the repository artifact exactly;
- [ ] SRT-to-WebVTT conversion passes the compatibility matrix;
- [ ] WebVTT-to-SRT conversion reports every known unrepresentable semantic;
- [ ] `--strict` blocks lossy conversion;
- [ ] Windows binaries pass smoke tests;
- [ ] macOS binaries pass smoke tests;
- [ ] Linux binaries pass smoke tests;
- [ ] amd64 and arm64 release artifacts are produced for the supported matrix;
- [ ] CLI help documents every command, option, and exit code;
- [ ] the approved externally produced Cueson brand kit is present in the repository;
- [ ] GitHub code-quality and security settings are verified;
- [ ] eligible release pull requests satisfied the two-round Codex review protocol or have an explicit operator override;
- [ ] no pull request was merged by an AI agent without a pull-request-specific single-use operator authorization;
- [ ] post-merge housekeeping evidence exists for v1-bound implementation pull requests;
- [ ] security and dependency scans are green;
- [ ] fuzz regression corpus is green;
- [ ] release notes contain highlights only and end with the tagged changelog link;
- [ ] release verification issue is complete;
- [ ] no pre-v1 or v1 release step has mutated production `cueson.io` Cloudflare configuration.

## References

- Existing ShruggieTech SRT converter and schema: https://gist.github.com/h8rt3rmin8r/2773b55fc12bb0a666f2f903195f8c57
- Existing ShruggieTech WebVTT converter: https://gist.github.com/h8rt3rmin8r/a05e8941432303ae6e8319208f49dae6
- Existing WebVTT schema referenced by the converter: https://schemas.shruggie.tech/data/webvtt-json.schema.json
- Cueson project domain: https://cueson.io
- ShruggieTech brand distribution: https://brand.shruggie.tech
- GitHub Spec Kit: https://github.com/github/spec-kit
- JSON Schema Draft 2020-12: https://json-schema.org/draft/2020-12
- Semantic Versioning: https://semver.org/
- Keep a Changelog: https://keepachangelog.com/
- WebVTT specification: https://www.w3.org/TR/webvtt1/
