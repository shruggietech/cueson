# Data Model: Native WebVTT Workflows

## WebVTT document data

`format_data.webvtt` retains normalized `signature`, lexical `signature_line_raw`, optional `description`, ordered `metadata_lines`, and non-cue `blocks`. The signature line's normalized identity remains `WEBVTT`; lexical description and spacing never replace that identity.

The union of cue and non-cue block `source_order` values is unique and contiguous from zero. Each individual array is strictly increasing. This union, rather than timing order, defines rendering order.

## WebVTT block

A non-cue block contains `type`, `source_order`, LF-joined `raw`, decoded `raw_lines`, and optional structured `region`. `raw` equals the LF join of `raw_lines`. Region data is required only when `type` is `region` and prohibited for other block types.

Supported block types are `note`, `style`, `region`, and `unrecognized`. Placement after the first cue, malformed content, and other conformance observations belong in diagnostics rather than changing source order.

## WebVTT region data

A parsed region contains `settings_raw`, ordered `setting_occurrences`, and a map of recognized effective `settings`. Recognized keys are `id`, `width`, `lines`, `regionanchor`, `viewportanchor`, and `scroll`. Raw lines remain authoritative for preserved lexical content; the map is the consumer and canonical-render view.

## WebVTT cue data

Each cue retains optional `identifier_raw`, `timing_line_raw`, `settings_raw`, recognized effective `settings`, ordered `setting_occurrences`, LF-joined `raw_payload`, and `raw_payload_lines`. Raw payload equality with common `payload.raw_text` and `payload.lines` is validated.

The recognized cue-setting keys are `region`, `vertical`, `line`, `position`, `size`, and `align`. The setting occurrence sequence preserves duplicates, unknown names, invalid values, and lexical spelling. The effective map includes only valid recognized values.

## Setting occurrence

A setting occurrence contains `raw`, `name`, `value`, `recognized`, and `valid`. Occurrence order is source order within its timing or REGION setting list. `raw` preserves the original `name:value` token; `name` and `value` are decoded lexical components without normalization.

## Common cue projections

WebVTT cue timing populates normalized millisecond start, end, and duration. Raw payload populates common raw text and lines. The cue-text scanner populates plain text, native speaker observations, and optional timed text spans without mutating native data.

Safely representable `line`, `position`, `size`, and `align` values may populate common placement. `vertical` and `region` remain native because the common placement object cannot express them completely.

## Diagnostics and state transitions

Document parsing transitions through signature, header, and body states. Body units transition to cue, NOTE, STYLE, REGION, unrecognized-preserved, or fatal. Diagnostics retain stable code, severity, message, and applicable source order and cue identity.

A document advances to successful encoded state only after exact capture, UTF-8 decoding, bounded parsing, model validation, schema validation, and output staging succeed. Rendering advances to published state only after model/schema validation, deterministic serialization, output verification, and atomic publication succeed.
