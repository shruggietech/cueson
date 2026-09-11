# SubRip format contract

**Current source status:** v0.1.0 development `experimental`

**Published release status:** v0.0.0 `envelope_only`

**Stable target:** v1.0.0

Current source implements native SubRip detection, decoding, semantic ingest, exact source restoration, and deterministic model-driven rendering. This contract defines accepted grammar and tolerated variants; it does not claim that the published v0.0.0 executable contains the codec or that v1 stability gates are complete.

## Capability matrix

| Capability | Current source | Boundary |
|---|---|---|
| Cue JSON representation | Available | Canonical `format: "subrip"` with common and native fields. |
| Raw `.srt` detection and decoding | Experimental | Content evidence outranks extension evidence. |
| Semantic ingest | Experimental | `cueson encode` derives cues while retaining exact bytes. |
| Exact restoration | Available | `cueson restore` verifies and recreates the source envelope without the codec. |
| Model-driven rendering | Experimental | `cueson render --to srt` emits canonical SubRip from structured fields. |
| Cross-format conversion | Experimental | `cueson convert INPUT --to vtt` projects common cue semantics and reports every known incompatible SubRip feature. |
| Stable SubRip support | Unavailable | Stable support remains a v1.0.0 gate. |

## Source and decoding

Encode opens one regular no-follow source once, captures timestamps before reading, limits retained input to 64 MiB, and uses the same exact bytes for hashing, envelope storage, detection, and decoding. Cue JSON stores a portable basename, never the original path or machine identity.

Automatic decoding accepts UTF-8 with or without BOM and UTF-16 LE/BE with BOM. BOM-less UTF-16 is accepted only when a strong tested NUL-byte pattern identifies byte order. Invalid Unicode fails. A conflicting BOM and explicit encoding fails. Windows-1252 and ISO-8859-1 require explicit `--encoding` selection because byte-only guessing cannot distinguish them reliably.

Canonical encoding names are `utf-8`, `utf-16le`, `utf-16be`, `windows-1252`, and `iso-8859-1`. The CLI also accepts common punctuation-insensitive aliases documented in command help. The selected canonical name, BOM observation, confidence, and source line-ending observation are retained in the source asset.

## Grammar and tolerated variants

The parser scans physical CRLF, LF, and lone-CR lines and uses an explicit block state model. Exact source bytes always remain authoritative.

| Construct | Treatment |
|---|---|
| Integer sequence line | Accepted and preserved in `sequence_line_raw`. |
| Missing sequence line | Accepted when a timing line begins the block; diagnosed and preserved as an empty native sequence observation. |
| Duplicate, decreasing, or irregular integer | Cue order remains source order; a deterministic warning records the irregularity. |
| Non-integer line before timing | Rejected or diagnosed as malformed rather than silently treated as a sequence number. |
| Timing arrow | Requires `-->` between complete start and end times. |
| Timestamp | Non-negative hours fitting the integer-millisecond range, one-to-two minute/second digits below 60, comma or tolerated period separator, and one-to-three millisecond digits right-padded to milliseconds. |
| Reversed, negative, or overflowing time | Fatal parse error; endpoints are never silently swapped. |
| Coordinates | Optional complete `X1`, `X2`, `Y1`, `Y2` non-negative integer suffix with each key exactly once. Partial, duplicate, or malformed groups fail; canonical rendering orders the keys. |
| Payload | One or more ordered decoded lines; empty internal lines and trailing characters are retained. |
| Line endings | CRLF, LF, and lone CR are accepted; mixed input is observed while exact bytes remain unchanged. |
| Formatting | Well-formed `b`, `i`, `u`, and `font` tags remain in raw text and are removed only from derived plain text while retaining contents. Unknown or malformed angle text stays literal. |
| Speaker prefix | A conservative `Name:` prefix may add one `heuristic` speaker observation; it never changes payload text and can be disabled. |
| Malformed but preservable content | Exact bytes remain restorable and every known interpretation gap receives an ordered diagnostic; input with no safely recognized cue fails. |

`payload.lines` contains decoded payload lines without physical terminators. `payload.raw_text` joins them with LF without stripping content. `payload.plain_text` applies only the tag rule above. Native timing syntax and coordinates remain in `format_data.subrip` alongside normalized integer-millisecond timing.

## Detection and diagnostics

Explicit `--format srt` selects SubRip. In auto mode, valid content grammar outranks `.srt` extension evidence. A content/extension disagreement emits a warning. Ambiguous input, unknown format, invalid encoding, and malformed source fail without output. Selecting `vtt` identifies a schema-recognized format whose native decoder is unavailable rather than calling it unknown.

Diagnostics are ordered deterministically by source position and code. Stdout mode contains only Cue JSON or rendered SubRip; diagnostics use stderr and follow global quiet/silent rules.

## Rendering

Model-driven rendering is distinct from exact restoration. Canonical output uses cue order, sequence values `1` through `N`, timestamps formatted as `HH:MM:SS,mmm`, complete coordinates when present, structured `payload.raw_text`, LF line endings, one blank line between cues, and one final LF. Rendered output is accepted by the same parser and re-encoding preserves all representable normalized semantics.

A payload is not representable without ambiguity when its line structure would be reparsed as a cue boundary, including an internal sequence-and-timing suffix or a trailing empty payload line. Normal rendering emits `subrip_render_payload_ambiguous`; strict rendering refuses the output.

Existing destinations are never replaced without `--force`. Output is staged and committed only after successful validation/rendering, so failures do not leave partial files.

## Conversion to WebVTT

SubRip-to-WebVTT conversion preserves cue order, integer-millisecond timing, overlap, multiline payloads, and shared `b`, `i`, and `u` emphasis. Sequence lexemes are canonicalized into cue order rather than promoted to WebVTT identifiers. Coordinates, `font` presentation, heuristic speaker observations, unanchored token timing, OCR observations, and unrecognized native data are omitted or degraded with one deterministic runtime loss observation per occurrence. A cue whose end does not follow its start is rejected because Cueson does not invent target timing. `--strict` rejects any known loss before WebVTT is rendered or published.

## Fixture and stability boundary

The manifest-governed `testdata/subrip/` and conversion corpus covers canonical and tolerated timing, sequences, coordinates, tags, speakers, encodings, BOMs, line endings, malformed input, golden rendering, exact restoration, loss-free and lossy conversion, strict rejection, and fuzz seeds. Stable support still requires completion of the broader v1 milestone.
