# Compatibility Contract: SubRip and WebVTT Conversion

## Classification

- **Represented**: Meaning survives exactly; canonical lexical changes do not create a loss.
- **Transformed**: Readable meaning survives through an explicitly weaker or altered representation and produces a `degraded` or `ambiguous` loss.
- **Lost**: The semantic is omitted and produces an `omitted` loss.
- **Fatal**: No deterministic parser-valid projection preserves defensible readable meaning; conversion fails in both modes.

## Common and bidirectional semantics

| Semantic | SubRip to WebVTT | WebVTT to SubRip |
|---|---|---|
| Cue array order and overlap | Represented with contiguous target source order | Represented with canonical numeric sequence order |
| Positive integer-millisecond timing | Represented with canonical target punctuation | Represented with canonical target punctuation |
| Zero-duration timing | Fatal because WebVTT requires positive duration | Not accepted by the WebVTT source model |
| Multiline readable payload | Represented with target-aware escaping | Represented with target-aware escaping |
| Empty logical payload line | Transformed to a non-visible target-safe lexical placeholder when possible | Transformed to a non-visible target-safe lexical placeholder when possible |
| Balanced bold, italic, and underline | Represented canonically | Represented canonically |
| OCR observations | Lost per observation | Lost per observation |
| Source envelope and transport bytes | Outside conversion semantics; exact restoration remains available separately | Outside conversion semantics; exact restoration remains available separately |

## SubRip source semantics targeting WebVTT

| Semantic | Treatment |
|---|---|
| Sequence line and exact sequence spelling | Canonicalized into cue order without repurposing it as a WebVTT identifier; irregular source observations remain source diagnostics rather than target semantics |
| Pixel coordinates | Lost per cue because no viewport dimensions establish an equivalent WebVTT placement |
| Heuristic speaker observation | Lost as structured speaker data; visible source prefix remains payload text and is not promoted to native voice truth |
| Token timing without anchored payload offsets | Lost per token because target insertion points cannot be proven |
| Balanced `font` presentation | Content retained, presentation transformed, one loss per affected span |
| Unknown or malformed angle text | Preserved as literal text through WebVTT-safe escaping |
| NUL semantic content | Transformed to the target-safe replacement character with an explicit loss |
| Unrecognized source block recorded by diagnostics | Lost per block occurrence rather than silently ignored |

## WebVTT source semantics targeting SubRip

| Semantic | Treatment |
|---|---|
| Signature keyword | Canonical target syntax change, no loss |
| Signature description and header metadata | Lost per field or line |
| NOTE, STYLE, REGION, and unrecognized blocks | Lost per ordered block; a REGION definition and cue region reference are separate losses |
| Cue identifier | Lost per cue because canonical SubRip sequence numbers are ordering syntax, not identifiers |
| `region`, `vertical`, `line`, `position`, `size`, and `align` settings | Lost per effective setting in fixed canonical setting order |
| Unknown, invalid, or duplicate setting occurrence | Lost per occurrence after its source diagnostic remains separately available |
| Common placement derived from settings | Accounted by the corresponding setting losses without duplicate entries |
| Cue-wide or partial voice annotation and speaker observation | Content retained, speaker meaning lost; no visible speaker prefix is invented |
| Class, language, ruby, and ruby-text annotation | Content linearized, semantic or presentation meaning transformed per occurrence |
| Inline timestamp and token timing | Marker removed, readable content retained, timing lost per occurrence |
| Valid character reference | Decoded then escaped only as needed for safe SubRip text; no semantic loss |
| Malformed or unknown reference | Preserved literally when target-safe; otherwise transformed with an explicit ambiguity loss |

## Fatal target-grammar boundaries

Conversion fails without bytes when timing cannot satisfy the target grammar or payload content cannot be made parser-safe without inventing visible meaning or dropping a cue. Target parser acceptance is mandatory in normal and strict modes.
