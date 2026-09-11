# WebVTT format contract

**Current status:** v0.0.0 `envelope_only`

**Stable target:** Planned for v1.0.0

This page defines the current Cueson boundary for WebVTT (`.vtt`) and the structure and fidelity obligations planned for v1.0.0. It does not claim that a native WebVTT codec exists in v0.0.0. The [Cue JSON schema](../schema.md), [CLI contract](../cli.md), and [architecture of record](../architecture.md) remain authoritative for implemented behavior.

## Current v0.0.0 capability

The v0.0.0 schema recognizes the canonical format key `webvtt` and defines WebVTT-specific document, block, and cue fields. Schema recognition defines how a valid Cue JSON document represents the format. It does not mean the executable can parse or render a raw WebVTT file.

| Capability | v0.0.0 state | Boundary |
|---|---|---|
| Cue JSON schema representation | Available | Valid documents use `format: "webvtt"` and WebVTT `format_data`. |
| Cue JSON structural and semantic validation | Available where consumed by shipped behavior | Validation checks the common model, ordered WebVTT block shape, capability declaration, and source-envelope contract. There is no public `validate` command. |
| Exact source restoration | Available | `cueson restore` recreates verified source-envelope bytes without calling a WebVTT codec. |
| Raw `.vtt` detection and UTF-8 decoding | Unavailable | The executable does not accept a raw WebVTT file for native ingest. |
| Semantic ingest | Unavailable | No parser currently derives cues, blocks, settings, markup, speakers, tokens, or placement from raw WebVTT. |
| Model-driven rendering | Unavailable | The executable cannot render Cue JSON semantics as a new WebVTT file. |
| Cross-format conversion | Unavailable | WebVTT-to-SRT and SRT-to-WebVTT commands are not shipped. |
| Stable WebVTT support | Unavailable | Stable support remains a v1.0.0 acceptance gate. |

An externally authored Cue JSON document can contain a valid WebVTT source envelope and common model. The current executable can restore that envelope exactly after validating the document and source integrity. It cannot prove how the semantic fields were derived because native ingest is not implemented.

## Source and model boundary

Original WebVTT bytes in `source.assets[].data_base64`, together with their byte length and SHA-256 digest, are authoritative for exact restoration. Stored names are portable safe basenames, not original filesystem paths. The generic restore operation verifies the complete source bundle before publication and does not normalize the signature, header, block ordering, timing lines, cue settings, payload markup, whitespace, BOM, or line endings.

The common cue payload and WebVTT-specific fields are structured representations. When native ingest is implemented, decoded cue text, plain text, lines, normalized timing, speaker and token observations, placement, and `format_data.webvtt` will remain derived from the source and will not replace it. Model-driven rendering may produce canonical WebVTT syntax, but only source restoration may claim byte identity.

## Planned v1 structure and fidelity matrix

The following matrix records acceptance targets, not current functionality.

| Source construct | Planned v1 treatment |
|---|---|
| UTF-8 BOM | Detect it, preserve its source-byte presence, and decode according to WebVTT's UTF-8 requirement. |
| `WEBVTT` signature | Require it after permissible BOM handling and preserve the native signature line. |
| Header description and metadata lines | Preserve their text and order independently from normalized document metadata. |
| `NOTE` block | Preserve raw content and total source order even when no common-model equivalent exists. |
| `STYLE` block | Preserve raw content and total source order without claiming that every style can be converted. |
| `REGION` block | Preserve raw content, parsed region semantics when implemented, and total source order. |
| Cue identifier | Preserve the native identifier independently from Cueson's document-local cue ID. |
| Hours-optional timestamp | Parse deterministically, expose normalized milliseconds, and preserve the native timing line. |
| Cue timing and positioning settings | Preserve both raw settings and parsed values without forcing them into unsafe SubRip coordinate equivalence. |
| Voice, class, language, and ruby markup | Preserve native markup while deriving plain text and speaker observations separately. |
| HTML entity | Decode only for derived consumer views while retaining the original cue payload. |
| Inline timestamp | Preserve native syntax and derive word or phrase token timing when valid. |
| Whitespace-only cue payload line | Preserve it as source content instead of trimming or discarding it. |
| Overlapping rolling caption cue | Keep each cue one-to-one in source order without deduplication. |
| Adjacent non-cue blocks | Give every block an explicit source order so relative placement never depends on inference. |
| Unrecognized or malformed but preservable block | Retain raw content and emit a deterministic diagnostic instead of silently dropping it. |

The future parser must preserve raw payload text and raw structural lines independently of normalized consumer fields. Cue and non-cue block order must remain total across the whole document, including multiple adjacent blocks and overlapping rolling captions.

## Planned diagnostics and conversion

Future native operations must distinguish an unknown format from a schema-recognized format whose codec is unavailable. They must preserve unknown settings, markup, and blocks when practical and report interpretation limits deterministically.

Future WebVTT-to-SubRip conversion must account for information without a safe SubRip representation, including `STYLE`, `REGION`, and `NOTE` blocks, cue settings, inline word timing, and WebVTT-specific markup. Normal conversion may write representable content only while reporting each accepted loss. Strict conversion must reject any known loss without creating output.

## Fixture boundary

The v0.0.0 repository has no native WebVTT grammar fixture corpus. The schema's WebVTT shapes and generic source-envelope restoration contract are not parser, renderer, round-trip, rolling-caption, markup, or conversion evidence.

Native WebVTT work must add focused accepted, malformed, parser, renderer, ordering, rolling-caption, and conversion fixtures before claiming any corresponding v1 structure row. Every accepted source fixture must also prove byte-exact restoration independently from semantic parsing or rendering.

## Explicit v0.0.0 exclusions

The v0.0.0 release does not provide raw WebVTT ingest, signature detection, UTF-8 decoding, cue or block parsing, markup interpretation, speaker or token extraction, model-driven WebVTT rendering, native round trips, cross-format conversion, or a stable WebVTT compatibility promise. The presence of `webvtt` in the schema and the ability to restore source-envelope bytes must not be presented as any of those capabilities.
