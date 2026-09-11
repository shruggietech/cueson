# Native WebVTT codec contract

## Detection and decode

- Content detection accepts only a leading optional UTF-8 BOM followed by a boundary-valid `WEBVTT` signature.
- Decode accepts the existing bounded exact captured source and rejects non-UTF-8 encoding selections.
- Decode returns a complete Cue JSON document with experimental ingest, render, and restore capability, common cues, WebVTT-native data, ordered diagnostics, statistics, and the unchanged source asset.
- Fatal errors remain typed enough for the CLI to return exit code 1 without output publication.

## Parsing

- Signature and header parsing precede body parsing.
- Body order is represented by the contiguous union of cue and non-cue block `source_order` values.
- NOTE, STYLE, REGION, and unrecognized blocks retain decoded raw lines.
- Cue identifiers, timing lines, settings, and payload lines retain lexical representations alongside recognized semantics.
- Parsing is iterative, bounded by captured input size, and deterministic for the same bytes and options.

## Rendering

- Render requires a valid WebVTT Cue JSON document.
- Output is UTF-8 without BOM and uses deterministic LF lines and a final LF.
- Header and every body item are emitted in validated model order.
- Structured timing and effective recognized settings control canonical syntax; native bodies and payload markup remain available where no complete structured equivalent exists.
- Strict rendering rejects every known preserved conformance error or unsafe ambiguity; permissive rendering reports deterministic warnings.
- Render never claims source-byte identity.
