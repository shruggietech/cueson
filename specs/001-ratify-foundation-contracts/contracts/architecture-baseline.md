# Architecture Baseline Contract

The published architecture document must define the following implementation boundaries without shipping their implementation in this slice.

## Required responsibilities

- CLI wiring owns command parsing and stream behavior, not domain models.
- Version handling owns the single executable build-version source.
- The common model owns normalized cross-format semantics.
- Schema handling owns the canonical contract, validation, and schema/software version equality.
- Source handling owns asset integrity, metadata capture, safe restoration, and byte authority.
- Codecs own detection, native parsing, and model-driven rendering for their format.
- Conversion owns loss accounting between model and target format.
- Future OCR providers own derived recognition; codecs may extract source images but do not embed OCR engines.

## Required separations

- Exact restoration is independent from codecs and model rendering.
- Restoration validates portable safe basenames before constructing output paths and rejects Windows device names and alternate-data-stream syntax on every platform.
- Restoration builds and validates the complete output plan before opening any destination and rejects basename or destination collisions.
- Native format data complements rather than replaces the common model.
- OCR observations never replace source assets or native cue content.
- Public compatibility is limited to the CLI and Cue JSON Schema; Go packages remain internal.

## Deferred ownership

- Issue #4 implements CLI and version foundations.
- Issue #5 implements the canonical schema and embedding.
- Issue #6 implements source integrity, generic exact restoration, its CLI command, and platform metadata boundaries.
- Issue #7 implements fixtures and conformance infrastructure.
- Later codec slices implement SRT and WebVTT parsing and rendering.
