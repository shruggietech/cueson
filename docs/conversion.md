# Conversion contract

**Current status:** Published in v1.0.0

**Historical release status:** v0.0.0 does not include conversion

`cueson convert` projects validated common-model subtitle semantics into the requested native format. Conversion is model-driven and never claims byte identity. Exact source restoration remains the responsibility of `cueson restore`.

## Operating rules

SubRip and WebVTT conversion preserves cue order, integer-millisecond timing, overlap, multiline payloads, and shared `b`, `i`, and `u` emphasis when the target grammar can represent them. The converter computes the complete ordered loss report before rendering. Permissive conversion emits one warning for every atomic loss, while `--strict` refuses all output when any known loss exists. Runtime loss reports are not stored in Cue JSON.

The loss report is deterministic and privacy-bounded. It contains stable codes, fixed messages, source and target formats, portable JSON Pointer paths, source-order references, cue IDs already present in the model, and small controlled context attributes. It never contains source bytes, caller paths, or machine identity. More than 8,192 atomic losses rejects the conversion rather than truncating the report.

## Stable loss vocabulary

| Code | Direction | Meaning |
|---|---|---|
| `conversion_metadata_omitted` | Both | Common document metadata has no native target representation. |
| `conversion_payload_line_degraded` | Both | An empty payload line requires a non-visible target placeholder. |
| `conversion_ocr_observation_omitted` | Both | Derived OCR observations are not subtitle target syntax. |
| `conversion_speaker_observation_omitted` | Both | A common speaker observation has no independent target field. |
| `conversion_token_timing_omitted` | Both | Token timing has no equivalent target structure. |
| `conversion_placement_omitted` | Both | Common placement cannot be represented by the target evidence available. |
| `conversion_nul_degraded` | Both | NUL is replaced with the semantic replacement character. |
| `conversion_subrip_coordinates_omitted` | SRT to VTT | SubRip pixel coordinates lack viewport evidence needed for a truthful WebVTT placement. |
| `conversion_subrip_font_degraded` | SRT to VTT | SubRip `font` presentation is removed while text is retained. |
| `conversion_subrip_unrecognized_block_omitted` | SRT to VTT | An unrecognized SubRip source block remains in the source envelope but not target output. |
| `conversion_subrip_payload_ambiguous` | SRT to VTT | Payload structure would be reparsed as a different SubRip cue boundary. |
| `conversion_webvtt_description_omitted` | VTT to SRT | The WebVTT signature description has no SubRip representation. |
| `conversion_webvtt_metadata_omitted` | VTT to SRT | A WebVTT header metadata line has no SubRip representation. |
| `conversion_webvtt_block_omitted` | VTT to SRT | A NOTE, STYLE, REGION, or unrecognized body block has no SubRip representation. |
| `conversion_webvtt_identifier_omitted` | VTT to SRT | A native WebVTT cue identifier is not a SubRip sequence identity. |
| `conversion_webvtt_setting_omitted` | VTT to SRT | A valid effective WebVTT cue setting has no equivalent target field. |
| `conversion_webvtt_setting_occurrence_omitted` | VTT to SRT | A duplicate, unknown, or invalid retained setting occurrence is omitted. |
| `conversion_webvtt_voice_degraded` | VTT to SRT | Voice markup is removed while readable text is retained. |
| `conversion_webvtt_markup_degraded` | VTT to SRT | Native class, language, ruby, or other WebVTT markup meaning is weakened. |
| `conversion_webvtt_inline_timing_omitted` | VTT to SRT | Inline timestamp syntax and its token boundary are omitted. |
| `conversion_webvtt_entity_ambiguous` | VTT to SRT | A character reference stays encoded to avoid target markup reinterpretation. |

## Fatal incompatibilities

Fatal incompatibilities are not losses because no safe target document can be produced. Conversion rejects unsupported source-target pairs, structurally or semantically invalid Cue JSON, corrupt source-envelope integrity, a non-positive SubRip interval projected to WebVTT, carriage returns that cannot be represented as physical payload lines, target grammar ambiguity that would change readable text, and any complexity limit exceeded before a complete report can be produced.

Fatal rejection writes no target payload and leaves a selected destination absent or unchanged. Strict loss rejection follows the same publication rule. Diagnostics and warnings use stderr; successful target bytes use stdout unless a destination is selected.

## Compatibility summary

| Source | Target | Loss-free baseline | Common lossy boundaries | Fatal boundaries |
|---|---|---|---|---|
| SubRip | WebVTT | Ordered positive-duration cues, millisecond timing, overlap, multiline text, and shared emphasis. | Coordinates, `font`, speaker or token observations, OCR, placement, unrecognized blocks, NUL, empty lines, and ambiguous source payload boundaries. | Invalid model or envelope, unsupported pair, non-positive target interval, carriage return, report or item ceiling. |
| WebVTT | SubRip | Ordered positive-duration cues, millisecond timing, overlap, multiline text, and shared emphasis. | Signature description, header metadata, blocks, identifiers, settings, voice, class, language, ruby, inline timing, entities, speakers, OCR, placement, NUL, and empty lines. | Invalid model or envelope, unsupported pair, carriage return, target markup reinterpretation, report or item ceiling. |

Same-format transformation belongs to `cueson render`, not `cueson convert`. The accepted grammar, native-fidelity rules, and row-level conformance evidence remain in the [SubRip](formats/srt.md) and [WebVTT](formats/webvtt.md) contracts.
