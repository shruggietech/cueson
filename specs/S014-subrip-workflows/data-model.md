# S014 Data Model

## Capability profiles

- `subrip`: `experimental`, ingest/render/restore true, OCR false.
- `webvtt`: `envelope_only`, ingest/render false, restore true, OCR false.

## Source capture

One bounded acquisition returns `SourceAsset` plus the exact bytes used for detection and decoding. Length and digest must agree; paths and machine identifiers never enter the document.

## Codec registration

A registration contains canonical format, aliases, extensions, and optional detector, decoder, and renderer functions. Lookup normalizes `srt` to `subrip` and `vtt` to `webvtt`.

## SubRip cue

Each cue contains source-order ID, normalized timing, raw/plain text, payload lines, optional heuristic speaker, optional coordinates, exact decoded native sequence/timing lines, and associated diagnostics. Rendering derives sequence numbers without rewriting native observations.

## Validation transitions

1. Bounded bytes become a validated source asset.
2. Detection selects an installed decoder or fails.
3. Decoding produces Unicode text and observations or fails.
4. Parsing produces common/native data and diagnostics.
5. Model and embedded-schema validation complete before publication.
6. Rendering accepts only valid documents with an installed target renderer.
