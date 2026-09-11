# Research: Native WebVTT Workflows

## Decision: Standards authority and tolerance boundary

Use the current [W3C WebVTT specification](https://www.w3.org/TR/webvtt1/) as the normative grammar authority. Distinguish conforming syntax from the specification's user-agent recovery algorithm: Cueson may retain recoverable invalid constructs, but it records a deterministic diagnostic and does not silently relabel them as conforming.

**Rationale**: Cueson is both an ingest tool and a future authoring/conversion tool. Real-world compatibility needs bounded recovery, while trustworthy interchange requires conformance errors to remain visible.

**Alternatives considered**: Reject every nonconforming construct, which would exclude files ordinary players parse; copy browser recovery without diagnostics, which would violate no-silent-loss requirements.

## Decision: UTF-8, BOM, NUL, and physical lines

Accept strict UTF-8 with an optional leading UTF-8 BOM. Reject UTF-16, legacy encodings, malformed UTF-8, and non-UTF-8 overrides. Replace embedded NUL with U+FFFD only in the decoded semantic view, emit `webvtt_nul_replaced`, and retain the original byte unchanged in the source envelope. Scan CRLF, LF, and lone CR iteratively into physical lines while retaining exact source bytes and a line-ending observation.

**Rationale**: WebVTT is UTF-8-only, and the standard parser replaces NUL. This preserves browser-compatible semantics without weakening source fidelity.

**Alternatives considered**: Reuse every generic SubRip decoder, which would falsely accept non-WebVTT encodings; reject NUL outright, which is stricter than useful parser recovery and unnecessary when exact bytes remain authoritative.

## Decision: Total source order without duplicated cue storage

Retain common cues in `cues` and non-cue WebVTT units in `format_data.webvtt.blocks`; use the union of their unique, contiguous `source_order` values as the complete body order. Do not duplicate cue objects in the native block array.

**Rationale**: The existing model already separates common cue semantics from format-native non-cue structure. A validated contiguous union proves that no body item silently disappeared while avoiding two mutable copies of each cue.

**Alternatives considered**: Add cue-reference blocks, which adds synchronization and validation complexity without more ordering information; store only a cue list, which loses NOTE, STYLE, REGION, and unknown block placement.

## Decision: Ordered lexical setting fidelity plus effective semantics

Keep `settings_raw` and a recognized effective settings map, and add ordered setting occurrences containing raw, name, value, recognition, and validity. Apply the same representation to parsed REGION settings. Process occurrences in source order, diagnose duplicates or invalid values, and retain them even when only the last valid recognized value becomes effective.

**Rationale**: A map alone loses duplicates, unknown names, whitespace, lexical order, and invalid-but-preservable data. The effective map remains convenient for consumers and canonical output.

**Alternatives considered**: Replace settings with a map only, which violates native fidelity; store raw text only, which forces every consumer to reparse WebVTT.

## Decision: Iterative parser and tokenizer

Implement an iterative document state machine and a bounded cue-text tokenizer. The document parser validates signature and header, collects empty-line-delimited units, recognizes exact NOTE, STYLE, REGION, or cue grammar, and preserves unknown units. The tokenizer recognizes WebVTT tags, class suffixes, voice and language annotations, the WebVTT character-reference subset, and inline timestamps while keeping malformed or unknown text literal and diagnosed.

**Rationale**: Iterative scanning keeps memory and stack behavior proportional to the existing 64 MiB source ceiling and makes line, block, and diagnostic order deterministic.

**Alternatives considered**: Regular-expression-only parsing, which is fragile for ordered recovery and nested markup; a browser DOM/CSS dependency, which expands the product boundary and conflicts with pure-Go portability.

## Decision: Speaker and token derivation

Treat `<v annotation>` as a native speaker observation. Convert valid, strictly increasing inline timestamps inside the cue interval into ordered text-span token timing. Keep invalid, non-increasing, or out-of-range timestamp syntax literal and diagnose it. Never alter raw payload text or deduplicate overlapping cues.

**Rationale**: These are WebVTT-native semantics useful to Cue JSON consumers, but they remain derived observations rather than source truth.

**Alternatives considered**: Drop markup timing, which makes the common model incomplete; create word timing by guessing, which invents precision absent from the source.

## Decision: Canonical model-driven rendering

Render UTF-8 without BOM using LF, `WEBVTT`, preserved header data, and the validated source-order union. Emit canonical `HH:MM:SS.mmm` timestamps and fixed recognized-setting order. Reuse raw non-cue bodies and raw cue payload because they are the native representation, but regenerate timing and settings whenever structured semantics differ. Permissive rendering may retain diagnosed native constructs; strict rendering rejects known nonconforming or unsafe content.

**Rationale**: The output represents the current model, while `restore` remains the only byte-identity operation.

**Alternatives considered**: Return source bytes from render, which ignores model edits; canonicalize all STYLE, NOTE, markup, and unknown content, which would destroy information with no complete structured equivalent.

## Decision: Shared output transaction repair

Repair the existing CLI file publisher so every schema, encode, and render write stages, flushes, closes, reopens, and verifies content before atomic publication, with rollback for forced replacement. Add regression tests for short or failed writes.

**Rationale**: S014's direct new-file write can expose a partial destination after an I/O failure, contradicting the existing no-partial-output contract inherited by S015. Fixing the shared boundary once is proportional and prevents copying defective behavior.

**Alternatives considered**: Limit the repair to WebVTT, which leaves identical existing commands unsafe; defer it, which would knowingly ship S015 against its own acceptance criteria.
