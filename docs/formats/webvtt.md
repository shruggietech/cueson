# WebVTT format contract

**Current status:** v0.1.0 `experimental`

**Stable target:** v1.0.0

This page defines the native WebVTT (`.vtt`) capability implemented in current development source. The [Cue JSON schema](../schema.md), [CLI contract](../cli.md), and [architecture of record](../architecture.md) remain authoritative for shared behavior.

## Current capability

| Capability | Development state | Boundary |
|---|---|---|
| Cue JSON schema representation | Available | Valid documents use `format: "webvtt"` and WebVTT-native `format_data`. |
| Raw detection and UTF-8 decoding | Experimental | A boundary-valid `WEBVTT` signature is detected after an optional UTF-8 BOM; other encodings are rejected. |
| Semantic ingest | Experimental | The native parser derives common cues while retaining signature, header, block, cue, setting, markup, and timing structure. |
| Model-driven rendering | Experimental | `cueson render --to vtt` writes deterministic LF WebVTT from validated structured data. |
| Exact source restoration | Available | `cueson restore` verifies and recreates source-envelope bytes without invoking the WebVTT codec. |
| Cross-format conversion | Unavailable | WebVTT-to-SubRip and SubRip-to-WebVTT conversion remain owned by issue #33. |
| Stable WebVTT support | Unavailable | Stable support remains a v1.0.0 acceptance gate. |

## Detection, Unicode, and source authority

Content detection accepts only an optional leading UTF-8 BOM followed by `WEBVTT` at a valid signature boundary. Content evidence wins over a disagreeing file extension and the disagreement is reported. Input is acquired through the shared 64 MiB bounded capture path before parsing.

WebVTT decoding accepts valid UTF-8 only. The optional BOM is observed but does not enter decoded text. Malformed UTF-8 and non-UTF-8 selections fail. An embedded NUL remains unchanged in the source envelope, becomes U+FFFD only in the semantic text view, and produces a deterministic diagnostic.

Original bytes, byte length, and SHA-256 in `source.assets` remain the authority for exact restoration. Stored names are portable safe basenames and Cue JSON never records the caller's path or another local machine identifier. Rendering never claims byte identity.

## Document and source order

The parser requires the `WEBVTT` signature and retains its complete raw line, optional description, and ordered header metadata. Header content remains distinct from body blocks.

Every cue and non-cue body block receives one `source_order` value. The union is unique and contiguous across the complete document. This keeps adjacent NOTE, STYLE, REGION, unrecognized blocks, and overlapping or rolling cues in their original relative order without deduplication.

NOTE, STYLE, and unrecognized blocks retain complete raw LF-joined content and physical lines. REGION blocks retain those raw views plus the complete ordered setting occurrences and valid effective values for `id`, `width`, `lines`, `regionanchor`, `viewportanchor`, and `scroll`. Misplaced, unknown, invalid, or duplicate but safely bounded content is retained with ordered diagnostics; fatal signature, timing, Unicode, or structural errors are rejected.

## Cue fidelity

Each WebVTT cue retains its optional native identifier independently from Cueson's document-local cue ID. Duplicate native identifiers are preserved and diagnosed. Timing accepts valid hours-optional WebVTT timestamps, stores normalized integer milliseconds, and retains the native timing line. Equal, reversed, malformed, or overflowing intervals fail; decreasing cue starts remain in source order with a diagnostic.

Cue settings preserve the complete raw setting string and every occurrence in lexical order. Valid effective values are exposed for `vertical`, `line`, `position`, `size`, `align`, and `region`; unknown, invalid, and duplicate occurrences remain available with diagnostics.

Raw payload text is the LF join of retained decoded payload lines. Leading, trailing, and whitespace-only content is not trimmed. The iterative markup scanner separately derives readable plain text, native voice speakers, and valid inline-timestamp token spans while leaving markup, entities, language, class, ruby, and source bytes unchanged. Malformed or unsupported markup and entities are preserved and diagnosed rather than silently discarded.

## Canonical rendering

`cueson render INPUT.cueson.json --to vtt` validates the complete model and writes a `WEBVTT` signature, retained header metadata, body items in contiguous source order, normalized dot-millisecond timestamps, consistent cue settings, raw native blocks and payload markup, one blank separator between body items, and a final LF. Output is UTF-8 without a BOM.

Native lexical content is reused only when it remains consistent with structured fields. Structured timing is always formatted canonically. In permissive mode, safely retained diagnosed content may be emitted with a warning. `--strict` rejects known preserved conformance errors or unsafe ambiguity before a destination is published.

Rendering and restoration are intentionally different operations. Editing structured cue timing changes rendered WebVTT, while restoring the same Cue JSON still reproduces the original captured bytes.

## Diagnostics and verification

Diagnostics are deterministic and preserve parser encounter order. They distinguish recoverable header separators, decreasing starts, duplicate identifiers, invalid or duplicate settings, misplaced blocks, unknown blocks, malformed markup or entities, and semantic NUL replacement. Diagnostics go to stderr in CLI workflows; structured stdout remains only Cue JSON or WebVTT bytes.

The governed corpus under `testdata/fixtures/webvtt`, `testdata/malformed/webvtt`, and `testdata/fuzz/webvtt` covers accepted signatures, line endings, headers, native block kinds, identifiers, timestamps, settings, markup, entities, inline timestamps, whitespace, adjacency, rolling captions, recovery, fatal input, and bounded fuzz seeds. Accepted fixtures separately prove exact source restoration and model-driven parser-renderer cycles.

## Explicit exclusions

Current development source does not provide cross-format conversion, a public `validate` or `inspect` command, shell completion, stable v1 support declarations, schema publication to `cueson.io`, or any claim that rendered output is byte-identical to the captured source.
