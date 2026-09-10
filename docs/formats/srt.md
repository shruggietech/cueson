# SubRip format contract

**Current status:** v0.0.0 `envelope_only`

**Stable target:** Planned for v1.0.0

This page defines the current Cueson boundary for SubRip (`.srt`) and the grammar and fidelity obligations planned for v1.0.0. It does not claim that a native SubRip codec exists in v0.0.0. The [Cue JSON schema](../schema.md), [CLI contract](../cli.md), and [architecture of record](../architecture.md) remain authoritative for implemented behavior.

## Current v0.0.0 capability

The v0.0.0 schema recognizes the canonical format key `subrip` and defines SubRip-specific document and cue fields. Schema recognition defines how a valid Cue JSON document represents the format. It does not mean the executable can parse or render a raw SubRip file.

| Capability | v0.0.0 state | Boundary |
|---|---|---|
| Cue JSON schema representation | Available | Valid documents use `format: "subrip"` and SubRip `format_data`. |
| Cue JSON structural and semantic validation | Available where consumed by shipped behavior | Validation checks the common model, format-specific shape, capability declaration, and source-envelope contract. There is no public `validate` command. |
| Exact source restoration | Available | `cueson restore` recreates verified source-envelope bytes without calling a SubRip codec. |
| Raw `.srt` detection and decoding | Unavailable | The executable does not accept a raw SubRip file for native ingest. |
| Semantic ingest | Unavailable | No parser currently derives cues, text, timing, coordinates, tags, or speakers from raw SubRip. |
| Model-driven rendering | Unavailable | The executable cannot render Cue JSON semantics as a new SubRip file. |
| Cross-format conversion | Unavailable | SRT-to-WebVTT and WebVTT-to-SRT commands are not shipped. |
| Stable SubRip support | Unavailable | Stable support remains a v1.0.0 acceptance gate. |

An externally authored Cue JSON document can contain a valid SubRip source envelope and common model. The current executable can restore that envelope exactly after validating the document and source integrity. It cannot prove how the semantic fields were derived because native ingest is not implemented.

## Source and model boundary

Original SubRip bytes in `source.assets[].data_base64`, together with their byte length and SHA-256 digest, are authoritative for exact restoration. Stored names are portable safe basenames, not original filesystem paths. The generic restore operation verifies the complete source bundle before publication and does not normalize line endings, timestamps, sequence lines, timecode punctuation, cue text, or tags.

The common cue payload and SubRip-specific fields are structured representations. When native ingest is implemented, `payload.raw_text`, `payload.plain_text`, `payload.lines`, normalized timing, speaker observations, and `format_data.subrip` will remain derived from the source and will not replace it. Model-driven rendering may produce canonical SubRip syntax, but only source restoration may claim byte identity.

## Planned v1 grammar and fidelity matrix

The following matrix records acceptance targets, not current functionality.

| Source construct | Planned v1 treatment |
|---|---|
| Integer cue sequence line | Parse and preserve the literal sequence line in format-specific data. |
| Missing or irregular sequence value | Preserve the source distinction and avoid destructive renumbering. |
| Comma millisecond separator | Accept as canonical SubRip timing syntax. |
| Period millisecond separator | Accept as a documented real-world variant and preserve the native timing line. |
| One-to-three digit millisecond field | Parse deterministically while retaining native syntax. |
| Multi-line cue text | Preserve ordered raw lines and expose normalized consumer text without destructive trimming. |
| CRLF, LF, lone-CR, and mixed line endings | Observe and preserve source bytes; derived text must not erase the original line-ending evidence. |
| UTF-8 with or without BOM | Detect and decode under the documented policy while retaining exact bytes. |
| Common legacy encodings | Support only through an explicit, pinned, tested decoding policy with an operator override for ambiguous input. |
| `X1`, `X2`, `Y1`, and `Y2` coordinate suffixes | Preserve coordinate values and their native timing-line representation without claiming equivalence to WebVTT placement. |
| Inline formatting tags | Preserve raw text and tags while deriving plain text separately. |
| Speaker prefix heuristic | Record an optional derived observation without removing or rewriting the source payload. |
| Malformed but preservable cue block | Retain the source information when practical and emit a deterministic diagnostic instead of silently dropping it. |

Because SubRip has no single universally authoritative grammar or encoding mandate, v1 completeness means complete coverage of Cueson's documented grammar and tolerated variants, backed by representative parser and renderer fixtures. The future parser must use an explicit deterministic grammar or state machine and must not depend on destructive whole-file blank-line normalization.

## Planned diagnostics and conversion

Future native operations must distinguish an unknown format from a schema-recognized format whose codec is unavailable. They must preserve unknown or malformed source fragments when practical and report every interpretation limitation deterministically.

Future cross-format conversion must account for every known non-representable semantic. In particular, SubRip pixel coordinates cannot be mapped safely to WebVTT placement without sufficient viewport semantics. Normal conversion may emit an explicit warning for a permitted adaptation or loss; strict conversion must reject the operation without creating output.

## Fixture boundary

The committed `testdata/fixtures/source-envelope/basic-lf/source/captions.srt` asset proves fixture provenance, source-envelope integrity, path hygiene, and byte-exact generic restoration. It is not a native SubRip parser, renderer, round-trip, encoding, speaker-detection, or conversion fixture.

Native SubRip work must add focused accepted, malformed, parser, renderer, and conversion fixtures before claiming any corresponding v1 grammar row. Every accepted source fixture must continue to prove byte-exact restoration independently from semantic parsing or rendering.

## Explicit v0.0.0 exclusions

The v0.0.0 release candidate does not provide raw SubRip ingest, format detection, text decoding, cue parsing, speaker extraction, model-driven SubRip rendering, native round trips, cross-format conversion, or a stable SubRip compatibility promise. The presence of `subrip` in the schema and the ability to restore source-envelope bytes must not be presented as any of those capabilities.
