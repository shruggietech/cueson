# Conversion contract

**Current status:** Stable SubRip/WebVTT conversion published in v1.0.0; S026 adds the complete experimental four-format graph in unreleased 1.1.0-dev

**Historical release status:** v0.0.0 does not include conversion

`cueson convert` projects validated common-model subtitle semantics into the requested native format. Conversion is model-driven and never claims byte identity. Exact source restoration remains the responsibility of `cueson restore`.

## Operating rules

SubRip and WebVTT conversion preserves cue order, integer-millisecond timing, overlap, multiline payloads, and shared `b`, `i`, and `u` emphasis when the target grammar can represent them. The converter computes the complete ordered loss report before rendering. Permissive conversion emits one warning for every atomic loss, while `--strict` refuses all output when any known loss exists. Runtime loss reports are not stored in Cue JSON.

The loss report is deterministic and privacy-bounded. It contains stable codes, fixed messages, source and target formats, portable JSON Pointer paths, source-order references, cue IDs already present in the model, and small controlled context attributes. It never contains source bytes, caller paths, or machine identity. More than 8,192 atomic losses rejects the conversion rather than truncating the report.

All twelve ordered distinct pairs among SubRip, WebVTT, ASS v4+, and SSA v4 are available for their documented baselines. ASS/SSA remain experimental. Scripted conversion preserves the immutable original source envelope and constructs a private target model; target defaults and derived target owners do not become observations of original source provenance or a new Cue JSON envelope. The original document and all asset integrity checks complete before analysis.

Scripted sources preserve readable dialogue, integer-millisecond timing, and representable bold/italic/underline style and override emphasis when targeting text. Every omitted native metadata/style/event field, unsupported override, drawing span, attachment, and other retained native occurrence is accounted for independently. Source framing, ownership IDs, capture lexemes, and canonical declaration order are serialization choices rather than semantic losses. Unused styles still own observed presentation values and receive applicable field losses.

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
| `conversion_scripted_metadata_omitted` | ASS/SSA to text | Each omitted native script metadata field occurrence. |
| `conversion_scripted_style_field_omitted` | ASS/SSA to text | Each native style field occurrence not represented by the text target. |
| `conversion_scripted_event_field_omitted` | ASS/SSA to text | Each native dialogue field occurrence not represented by the text target. |
| `conversion_scripted_record_omitted` | ASS/SSA to text | Each omitted native non-dialogue event or physical record occurrence. |
| `conversion_scripted_section_omitted` | ASS/SSA to text | Each omitted source section whose content is not represented by target framing. |
| `conversion_scripted_attachment_omitted` | ASS/SSA to text | Each omitted native attachment, covering its header and encoded range once. |
| `conversion_scripted_override_omitted` | ASS/SSA to text | Each omitted unsupported native override occurrence. |
| `conversion_scripted_override_degraded` | ASS/SSA to text | Each native override whose supported target meaning is weaker, such as numeric bold weight. |
| `conversion_scripted_override_comment_omitted` | ASS/SSA to text | Each ignored non-tag portion of a native override span. |
| `conversion_scripted_drawing_omitted` | ASS/SSA to text | Each omitted drawing span in otherwise readable dialogue. |
| `conversion_scripted_centisecond_quantized` | Any scripted target | Each changed start or end endpoint under checked nearest-centisecond rounding. |
| `conversion_scripted_variant_field_omitted` | ASS/SSA variants | Each observed native field omitted by opposite-dialect construction. |
| `conversion_scripted_variant_field_degraded` | ASS/SSA variants | Each observed native field whose mapped dialect meaning is weakened. |
| `conversion_source_identifier_omitted` | Scripted sources | Each common source identifier annotation omitted from subtitle output. |

The established text-pair code meanings and ordering remain unchanged. Existing source-feature codes also apply to new directions only when that exact source feature is omitted or degraded. Scripted native pointers identify positional owner fields/tags/spans; contexts contain controlled feature names and occurrence numbers, never source values. One omitted attachment accounts for its complete encoded range without duplicating data-line losses.

Custom common source identifiers on scripted Cue JSON are annotations without native event columns and receive explicit omission losses. Native SubRip sequence framing and existing WebVTT identifier accounting retain their established policies. SSA nonzero per-color high-byte alpha observations receive field degradation losses when replaced by the dialect's effective AlphaLevel mapping; zero ignored high bytes are ordinary representation normalization.

## Scripted target defaults and timing

Text-source scripted targets contain matching ScriptType, WrapStyle 0, a canonical Default style, and ordered dialogue owners. Defaults use Arial 20; white primary, red secondary, black outline/shadow; neutral bold/italic/underline/strikeout; scales 100; spacing and angle 0; border style 1; outline 2; shadow 0; bottom-center alignment 2; margins 10; encoding 1; and SSA AlphaLevel 0. No video dimensions, installed-font availability, machine identity, or source-native authoring provenance are inferred. Target-only defaults are construction; replacing an observed source value requires loss accounting.

Every scripted target rounds each nonnegative common cue endpoint to the nearest centisecond, with ties upward. Quotient/remainder arithmetic avoids unchecked addition; each changed endpoint has its own degradation entry. Recompute duration from the rounded endpoints. Collapse (end no longer greater than start), negative timing, or arithmetic overflow is fatal in both modes. Do not silently truncate milliseconds or enlarge collapsed intervals. Native scripted source times are already exact centiseconds; valid edited scripted Cue JSON common timing can require the same explicit endpoint losses during dialect conversion.

Actual line feeds become native hard breaks and nonbreaking spaces become native hard-space controls. Literal braces and literal backslash followed by N, n, or h are fatal for text-source scripted targets because the selected native scanner cannot preserve them without acquiring control semantics. Other literal text must survive complete target reparsing. Shared emphasis becomes balanced native overrides with nested state; adjacent or trailing logical blank lines are representable without the text-target-only placeholder degradation.

## Scripted dialect conversion

ASS/SSA conversion deep-copies native owners and retains representable metadata, styles, dialogue/non-dialogue events, actors, safe effects, comments, unknown inert content, and attachments. Unchanged capture observations remain independently associated with original source positions; changed captures are removed rather than fabricated. Invalid preservation-only owners are fatal.

SSA style alignment maps to ASS as 1->1, 2->2, 3->3, 5->7, 6->8, 7->9, 9->4, 10->5, and 11->6; reverse those pairs for ASS to SSA. Unsupported values fail explicitly. Layer controls stacking while Marked is an authoring observation, so replacement by a zero target default is an accounted omission. Source-only style columns are accounted individually, including observed neutral values. Target-only neutral columns are defaults.

Color and alpha roles differ: SSA BackColour supplies both outline and shadow in the researched maintained renderer behavior, TertiaryColour is not equivalent to ASS OutlineColour, and AlphaLevel differs from ASS per-color alpha. Conversion maps representable values and reports each weakened or omitted observed field. It does not claim pixel equivalence, resolve fonts, load attachments, or execute native effects. Unknown or dialect-incompatible controls require explicit preservation, losses, or fatal refusal according to the selected contract.

## Fatal incompatibilities

Fatal incompatibilities are not losses because no safe target document can be produced. Conversion rejects unsupported source-target pairs, structurally or semantically invalid Cue JSON, corrupt source-envelope integrity, a non-positive SubRip interval projected to WebVTT, carriage returns that cannot be represented as physical payload lines, target grammar ambiguity that would change readable text, and any complexity limit exceeded before a complete report can be produced.

Fatal rejection writes no target payload and leaves a selected destination absent or unchanged. Strict loss rejection follows the same publication rule. Diagnostics and warnings use stderr; successful target bytes use stdout unless a destination is selected.

Drawing-only or entirely unreadable dialogue is fatal for scripted-to-text conversion, as is a zero-dialogue script targeting text. Mixed drawing/readable dialogue can succeed permissively with one omission loss per drawing span; no readable text is invented and no source cue is silently dropped. Valid empty scripts may remain valid across scripted variants. Preserved malformed native records, malformed override interpretation, target literal-markup reinterpretation, unsupported native alignment, quantized interval collapse, and native record/line/output ceilings are fatal edges. Every successful target reparses with the expected timing and readable semantics before any payload is returned.

## Compatibility summary

| Source | Target | Represented baseline | Common lossy boundaries | Fatal boundaries |
|---|---|---|---|---|
| SubRip | WebVTT | Ordered positive-duration cues, millisecond timing, overlap, multiline text, and shared emphasis. | Coordinates, `font`, speaker or token observations, OCR, placement, unrecognized blocks, NUL, empty lines, and ambiguous source payload boundaries. | Invalid model or envelope, unsupported pair, non-positive target interval, carriage return, report or item ceiling. |
| WebVTT | SubRip | Ordered positive-duration cues, millisecond timing, overlap, multiline text, and shared emphasis. | Signature description, header metadata, blocks, identifiers, settings, voice, class, language, ruby, inline timing, entities, speakers, OCR, placement, NUL, and empty lines. | Invalid model or envelope, unsupported pair, carriage return, target markup reinterpretation, report or item ceiling. |
| ASS | SubRip | Readable ordered dialogue, exact timing, multiline text, and shared emphasis. | Native script/style/event presentation, non-dialogue records, unsupported overrides, mixed drawings, attachments, and derived observations. | Drawing-only/unreadable/zero dialogue, malformed owners, literal markup or cue-boundary reinterpretation, bounds. |
| ASS | WebVTT | Readable ordered dialogue, exact timing, multiline text, escaped literals, and shared emphasis. | Native script/style/event presentation, non-dialogue records, unsupported overrides, mixed drawings, attachments, and derived observations. | Drawing-only/unreadable/zero dialogue, malformed owners, nonphysical lines, bounds. |
| SSA | SubRip | Readable ordered dialogue, exact timing, multiline text, and shared emphasis. | Native script/style/event presentation, non-dialogue records, unsupported overrides, mixed drawings, attachments, and derived observations. | Drawing-only/unreadable/zero dialogue, malformed owners, literal markup or cue-boundary reinterpretation, bounds. |
| SSA | WebVTT | Readable ordered dialogue, exact timing, multiline text, escaped literals, and shared emphasis. | Native script/style/event presentation, non-dialogue records, unsupported overrides, mixed drawings, attachments, and derived observations. | Drawing-only/unreadable/zero dialogue, malformed owners, nonphysical lines, bounds. |
| SubRip | ASS | Positive-duration centisecond-exact dialogue and shared emphasis with deterministic defaults. | Coordinates, font, unrecognized blocks, common observations, NUL, and changed timing endpoints. | Literal braces/control pairs, collapsed/overflow timing, target grammar mismatch, bounds. |
| SubRip | SSA | Positive-duration centisecond-exact dialogue and shared emphasis with deterministic defaults. | Coordinates, font, unrecognized blocks, common observations, NUL, and changed timing endpoints. | Literal braces/control pairs, collapsed/overflow timing, target grammar mismatch, bounds. |
| WebVTT | ASS | Centisecond-exact dialogue and shared emphasis with deterministic defaults. | Metadata, blocks, identifiers/settings, voice/class/ruby/language/inline timing, common observations, NUL, and changed endpoints. | Literal braces/control pairs, collapsed/overflow timing, target grammar mismatch, bounds. |
| WebVTT | SSA | Centisecond-exact dialogue and shared emphasis with deterministic defaults. | Metadata, blocks, identifiers/settings, voice/class/ruby/language/inline timing, common observations, NUL, and changed endpoints. | Literal braces/control pairs, collapsed/overflow timing, target grammar mismatch, bounds. |
| ASS | SSA | Dialogue/native owners retain representable semantics; observed dialect-only columns mean a canonical baseline can be lossy. | Layer, ASS-only style columns, outline/back/alpha roles, and incompatible dialect fields or controls. | Malformed owners, unsupported alignment/scalars, unsafe captures, target reparse or bounds. |
| SSA | ASS | Dialogue/native owners retain representable semantics; observed dialect-only columns mean a canonical baseline can be lossy. | Marked, TertiaryColour, AlphaLevel/color roles, and incompatible dialect fields or controls. | Malformed owners, unsupported alignment/scalars, unsafe captures, target reparse or bounds. |

Same-format transformation belongs to `cueson render`, not `cueson convert`. The accepted grammar, native-fidelity rules, and row-level conformance evidence remain in the [SubRip](formats/srt.md), [WebVTT](formats/webvtt.md), and [ASS/SSA](formats/ass-ssa.md) contracts. Full release-wide CLI integration and hardening remain #62/#63; stable freeze and publication remain later separately governed gates.
