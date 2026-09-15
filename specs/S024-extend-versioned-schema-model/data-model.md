# S024 data model

## Exact contract authority

The current1.1.0-dev pair and historical1.0.0 pair select immutable local resources independently. Registry entries carry expected artifact identity and version-specific semantics. No producer version, normalized URI, semver range or network reference participates.

## Scripted document and physical sequence

Matching ass/ssa branches own dialect, ordered sections/records/styles/events/attachments. Sections contain document-local section_id, source_order, exact name and optional raw_header capture. Records contain record_id, source_order, physically owning section_id, closed kind, optional raw_line capture and kind-specific fields/references. Section headers and records jointly own one contiguous physical sequence. Duplicate IDs, unowned records and sequence gaps reject. Nonnil capture observations must match their positions in the verified original asset; source-declared timestamp captures remain independent of edited common timing.

format_declaration records contain ordered exact field names and recognized field identities. Their record_id is declaration_id. Declared fields apply only within their section occurrence until replacement. Recognized canonical dialect fields are mandatory/unique; Name/Actor is the only alias, final Text owns comma suffix.

## Styles and events

Styles reference physical records and applicable declarations, retain exact names and ordered field occurrences with optional one typed value. Events reference records/declarations and explicit event_type/cue_id, own exact text and ordered native fields/spans/tags/karaoke observations. Constructed recognized styles/events without captured declarations may omit declaration_id. Typed non-time values agree with owning lexical fields; native time lexemes remain observations after common edits.

Each dialogue has exactly one common cue with matching physical source_order and event/style references. Non-dialogue content creates no artificial common cues. Common timing owns edits; raw_text equals event Text; readable logical lines join plain_text with LF. Source-relative span offsets use Unicode scalars, resolve into native Text and retain ordered interpretation/provenance. Zero dialogue uses zero cues and null media summaries for current scripted branches only.

## Attachments

attachment_id/header_record_id/data_start_record_id/data_record_count/type/safe name identify bounded contiguous encoded-data records in the same owning section, without overlap or duplicated arrays. Malformed bounded retained content requires deterministic diagnostics and remains preservation-only; generic restoration still verifies original source bytes.

## Safety and scale

Metadata context is determined by the closed accepted profile, not producer flags. Unsafe known/unknown filesystem/resource/identity/active or unclassifiable authoring metadata rejects original source and edited native views. Explicit dialogue, font-family and bounded embedded-content roles remain content. Declaration/field names and semicolon comments do not inherit those exemptions; comments require conservative inert-content inspection, replacing the earlier development allowance for path-looking comments.

Existing common collection limits plus ratified scripted physical/item/field/span/tag/line/depth/attachment limits apply before indexed reference traversal. No silent truncation or diagnostic leakage.

## Lifecycle

Current scripted support is schema_only with no installed native ingest/render. Native codecs/conversion/conformance and final stable1.1.0 promotion remain later native issues. Historical1.0.0 structure, capabilities and semantic constraints are unchanged.
