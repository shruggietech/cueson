# S026 outbound conversion research

**Date:** 2026-09-15

**Authority inspected:** Root `AGENTS.md`, constitution 0.1.0, architecture of record, the S026 specification, GitHub #60/#61, the ratified ASS/SSA selected profile, and the current model, codecs, and conversion implementation. This artifact records research only and does not authorize implementation before the blocking analysis gate.

## Findings and decisions recommended

Scripted conversion must traverse ordered native owners in addition to common cues. `model.ProjectScriptedText` deliberately supplies readable lines, drawing exclusion, tags, spans, and karaoke facts; it does not project style emphasis. Using `cue.payload.plain_text` alone would silently discard style emphasis, resets, unknown native records, inline comments, presentation, attachments, and event fields. Preserve representable readable text and basic bold/italic/underline with an ordered style and override state machine, and account separately for every known target-incompatible native occurrence.

Recommend fatal rejection in both modes for an empty script, an entirely unreadable dialogue cue, or drawing-only dialogue when targeting SubRip/WebVTT. This avoids invented placeholder dialogue, silently dropping a source cue, and invasive source-to-target index remapping. Mixed drawing and readable dialogue may succeed permissively, with an omission entry for each drawing span; strict mode rejects those losses. Consecutive or trailing readable line breaks may use the existing non-visible target placeholder policy with the existing atomic empty-line degradation code. The explicit policy should distinguish an empty dialogue from an empty line inside otherwise readable dialogue.

Treat malformed preservation-only native records and malformed overrides as fatal rather than making an interpretation guess. Unknown but structurally valid overrides may be omitted with complete losses. Drawing omission never claims OCR or reconstruction of drawing text. Non-dialogue native events are omitted with explicit native record losses; their payload is not promoted to dialogue.

Keep the existing two-format translator and report ordering unchanged. Extend the source/target format vocabulary and append new loss-code ranks after the established ranks. Existing metadata, derived speaker/token/OCR/placement, SubRip coordinate/font, WebVTT metadata/block/identifier/setting/voice/inline timing, and empty-line codes retain their meanings when those exact features are omitted on a new direction. Native source features should use native pointers; distinct derived observation omissions may use their established common-model pointers.

## Shared text and emphasis interpretation

Use the owning `ScriptedCueData.StyleID` to obtain initial style bold, italic, and underline. For every override span, process its tags in source order before subsequent literal spans. Exact `b0`/`b1`, `i0`/`i1`, and `u0`/`u1` are representable; empty parameters reset that property to the currently applicable style default. A full reset returns to the owning style, and a declared named reset selects the unique named style. An invalid or ambiguous reset is fatal. Positive explicit bold weights may degrade to boolean bold with an atomic degradation entry; alternatively reject them explicitly if the contract chooses a smaller accepted projection. Never silently interpret malformed parameters as valid emphasis.

Emit balanced, canonically nested target `b`, `i`, and `u` tags at text-run boundaries, not a direct stream of independent opening/closing tags that can cross. An override that changes state while a drawing span is active still changes the state used by later readable text. Preserve the ratified explicit-break semantics: `N` is a hard break, `n` is a break only under wrapping mode 2, and `h` is a non-breaking space. Script-level wrapping and `q` overrides affect this projection but any unsupported automatic layout behavior still receives presentation loss accounting.

Escape WebVTT literal text using its native entity rules before surrounding it with target emphasis. For SubRip, verify that the rendered payload reparses to the intended readable string. SubRip has no universal literal-markup escape: scripted readable text that becomes recognized target markup or an ambiguous cue boundary must fail both modes. Preserve source literals rather than deleting or reinterpreting them. Existing runtime target reparse currently checks timing and plain text; add focused emphasis state assertions in conversion tests because readable equality alone cannot prove shared emphasis preservation.

These semantics are informed by the [Aegisub override manual](https://aegisub.org/docs/latest/ass_tags/) and the [pinned libass parser](https://github.com/libass/libass/blob/bbb3c7f1570a4a021e52683f3fbdf74fe492ae84/libass/ass.c). The manual describes text-state overrides, breaks, resets, weights, and drawings; the pinned parser supplies concrete dialect field lists. Cueson's selected profile and explicit conversion policy remain authoritative, without a pixel-equivalence claim.

## Proposed private APIs

`analyzeScriptedSource(document model.Document, targetFormat string) (matrixAnalysis, error)` belongs in `internal/convert/scripted_analysis.go`. It populates the existing one-translation-per-source-cue contract for text targets, traverses native owners once, appends bounded losses through `appendOneLoss`, and rejects unsupported or unsafe target edges before projection. Root-owned dispatch chooses this function for ASS/SSA without modifying established SubRip/WebVTT branches.

`translateScriptedPayload(document model.Document, cueIndex int, targetFormat string) (payloadTranslation, error)` belongs in `internal/convert/scripted_text.go`. It reads validated native event/style ownership, derives balanced shared-emphasis output, and returns payload-local issues for safe target representation changes. Native tag/span losses should be appended by analysis using their native structural paths; do not mislabel them as old SubRip/WebVTT markup losses. A bounded internal run representation with readable text and three boolean emphasis properties can serve both text rendering and scripted-target construction.

`scriptedNative(document model.Document) *model.ScriptedDocumentData` is a small branch selector; maps from record/style/event IDs to their validated owners prevent repeated linear searches. New analysis must check context or use a context-aware entry point for its long loops, so cancellation is not delayed until rendering. Root can extend `matrixAnalysis` narrowly if variant preservation needs additional projection data; avoid inventing a second competing conversion report model.

## Proposed closed loss vocabulary and references

| Code | Kind | Atomic source reference and applicability |
|---|---|---|
| `conversion_scripted_metadata_omitted` | omitted | Each omitted `/format_data/{dialect}/records/N/fields/M`, associated with that metadata record order. |
| `conversion_scripted_style_field_omitted` | omitted | Each omitted `/format_data/{dialect}/styles/N/fields/M`, associated with the owning style record order. |
| `conversion_scripted_event_field_omitted` | omitted | Each omitted native dialogue field occurrence other than represented timing/text; associated with the owning cue and event record order. |
| `conversion_scripted_record_omitted` | omitted | Each omitted non-dialogue event, comment, unknown, or retained nonsemantic record; pointer to its native owner or record and that record order. |
| `conversion_scripted_section_omitted` | omitted | Each omitted noncanonical or empty source section occurrence not already explained by target format framing; associated with its header order. |
| `conversion_scripted_attachment_omitted` | omitted | Each attachment occurrence, pointer to `/format_data/{dialect}/attachments/N`, associated with the header record order; its encoded range is covered once. |
| `conversion_scripted_override_omitted` | omitted | Each unsupported top-level `/format_data/{dialect}/events/N/tags/M`; associated with the owning cue/order. |
| `conversion_scripted_override_degraded` | degraded | Each representable-but-weakened override, such as a bold weight reduced to boolean emphasis. |
| `conversion_scripted_override_comment_omitted` | omitted | Each ignored non-tag portion of an override span; pointer to `/format_data/{dialect}/events/N/spans/M` with controlled occurrence context. |
| `conversion_scripted_drawing_omitted` | omitted | Each mixed-content drawing span, using the native span pointer and owning cue/order. |
| `conversion_scripted_centisecond_quantized` | degraded | Each changed `/cues/N/timing/{start_milliseconds|end_milliseconds}` endpoint when constructing a scripted target. |
| `conversion_scripted_variant_field_omitted` | omitted | Each native field not representable in the opposite dialect, retaining its original style/event field pointer. |
| `conversion_scripted_variant_field_degraded` | degraded | Each mapped dialect field whose original semantics are weakened, retaining its original native field pointer. |

The final ratification can consolidate code names when one code still expresses one stable semantic rule. Source format framing, capture lexemes, IDs used only for ownership, and canonical declaration order are target construction rather than semantic omissions; document that distinction explicitly. Do not omit an unused style's field values merely because no current dialogue uses it. Shared emphasis fields represented through cue runs do not also need an omission entry for that represented property, but unused or unsupported style properties do. Default margins, layer, fonts, colors, and layout still express native presentation even when common defaults make their effects unobtrusive; do not silently dismiss them as irrelevant.

Every report message must be a constant safe sentence, never a native raw value or caught parser message. Context may contain decimal occurrence indices and closed normalized feature names, never arbitrary field names, style names, actor strings, attachment names, tag parameters, source excerpts, URLs, or caller paths. Ordinal native pointers are already sufficient to identify unknown fields safely. Tags nested inside a transform are one unsupported complex native tag unless the implementation explicitly models their children; do not claim a transform was represented by interpreting its parameters as ordinary top-level emphasis.

## Architectural and safety integration risks

`Report.Validate` currently recognizes source orders only from cues, WebVTT blocks, and diagnostics. Add every scripted section and record order before validating native style, metadata, comment, and attachment losses. Preserve uniqueness, canonical numeric pointer ordering, and the 8,192-loss ceiling. Run loss collection once with checked append operations; crossing the ceiling must fail without a truncated report or output bytes.

`projectDocument` currently indexes the first cue while recomputing summaries and only constructs SubRip/WebVTT. The explicit no-readable-target policy prevents an outbound empty-cue panic, but the helper still needs safe empty-summary handling for legitimate empty scripted targets. Conversion must validate the complete original source envelope, including secondary assets, before analysis. Nil context should produce a safe fatal error rather than panic.

Most critically, `validateScriptedDocument` unconditionally applies `scriptedSourcePrivacy` using the document's format, and that validator requires the preserved primary bytes to declare that same ScriptType/style dialect. A cloned SRT source cannot validate as an ASS document, and an unchanged ASS source cannot validate as SSA. Omitting raw capture fields does not remove this check. Do not weaken native Cue JSON privacy validation or rewrite the source envelope to make a private target document pass. Use a narrow constructed-target rendering boundary that validates the original source document and integrity independently, validates constructed owners, serializes them, and fully reparses the bounded target candidate before publication. The resulting target native bytes are output, not claimed original source provenance or a target Cue JSON envelope.

ASS/SSA variant conversion needs a separate fidelity plan, not a call to outbound flattening. Shared style fields, text, actor, timings, margins, metadata, comments, and safe attachments should survive when their target meaning is representable. ASS Layer versus SSA Marked, OutlineColour versus TertiaryColour, alpha handling, alignment numbering, and ASS-only style columns require explicit semantic mapping, per-field losses, or fatal boundaries. Unknown fields and dialect-specific overrides must not be silently reclassified as compatible. Opposite-dialect privacy and bounded target reparse remain blocking even when the visible dialogue is unchanged.

## Required focused evidence

Exercise both scripted dialects to both text targets with inherited and changed emphasis, resets, explicit breaks, mixed drawings, karaoke, actor, style presentation, non-dialogue/unknown records, duplicate declarations, attachments, and complete report snapshots. Add fatal cases for drawing-only/empty scripts, malformed retained owners, literal markup reinterpretation, target cue-boundary ambiguity, cancellation, source integrity, and loss amplification. Verify strict mode collects its complete report before any renderer call, returns no payload, and preserves both existing destinations and stdout. Retain unchanged old pair output and loss snapshots, and prove original source envelope equality before and after each new direction.
