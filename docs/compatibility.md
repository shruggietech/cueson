# Cueson Compatibility Contract

**Status:** v1.0.0 stable release published and independently verified

This document defines which Cueson surfaces receive a public compatibility promise, how format capability states are interpreted, and which release and platform claims are currently valid.

## Public v1 surfaces

The v1 compatibility promise covers exactly two public interfaces:

- The `cueson` CLI command names, options, aliases, streams, exit-code classes, overwrite rules, and strict-mode behavior documented in the [CLI contract](cli.md).
- The Cue JSON Schema identity, validation constraints, common model, native extensions, source envelope, and version rules documented in the [schema contract](schema.md).

Go packages remain under `internal/` and carry no public source-compatibility or import-compatibility promise. Repository scripts, fixtures, conformance indexes, inspection-report internals, conversion target projections, and diagnostic prose are verification or implementation surfaces unless another maintained contract explicitly says otherwise. Stable diagnostic and conversion-loss codes remain machine-consumable where their owning contract identifies them as stable.

## Version policy

The published v0.0.0 release is an immutable envelope-only foundation. Its executable, schema, archive set, and versioned schema copy remain unchanged and do not contain native SubRip or WebVTT ingest, rendering, conversion, validation, inspection, or completion.

The published v1.0.0 executable and schema implement the complete stable v1 command set and native SubRip and WebVTT workflows, and both format capability declarations are `stable`. Immutable tag [`v1.0.0`](https://github.com/shruggietech/cueson/tree/v1.0.0) and the verified [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0) identify the exact stable contract.

Starting with v1.0.0, incompatible changes to either public interface require a new major version. Compatible additions may use a minor version, and compatible corrections may use a patch version. Official Cueson executable and schema versions remain equal; a third-party producer's own version is independent from the Cue JSON schema version it targets.

No mutable `latest` schema identity is part of the contract. Every released schema copy is immutable. The v1.0.0 canonical and immutable schemas are available from the repository, release archives, embedded executable, and byte-exact public `cueson.io` route. S021 completed separately authorized public hosting of both v0.0.0 and v1.0.0 immutable schemas.

## Current development compatibility

S024 current source uses unreleased `1.1.0-dev` and accepts exact historical 1.0.0 input with the bundled released schema and preserved historical semantics. Validate, inspect, restore, matching native render and established conversion paths retain historical identity, producer and source truth. Missing, approximate, mismatched, unsupported and unpromoted final identities reject locally. Current schema discovery and new encode output use the development identity; released exact-version consumers reject new output until explicitly updated. No historical-output selector is added.

Current ASS/SSA branches support schema/model recognition, generic validation, inspection and exact restoration. Native ingest, model rendering and conversion remain later outcomes. Development shape recognition provides no stable release claim and does not publish a schema URI.

## Format capability states

- `envelope_only`: the release can preserve and exactly restore a valid source envelope but does not claim native semantic ingest or model-driven rendering.
- `schema_only`: the current schema/model recognizes the native shape; generic validation, inspection and restoration are available while native ingest/render are absent.
- `experimental`: the executable implements documented native behavior whose compatibility contract has not completed its stable release gate.
- `stable`: the format has passed the applicable v1 acceptance and release gates for that software/schema version.

Capability booleans are independent facts. Schema recognition does not imply an installed codec, model-driven rendering does not imply exact restoration, and exact restoration does not imply native parsing. Consumers must use `format_support` rather than infer capability from a format key or file extension.

## Fidelity and conversion

The source envelope is authoritative for exact restoration. Common cue fields and format-native fields are structured representations for consumption and editing. Rendering serializes that structured model and never claims byte identity with the captured source.

SubRip and WebVTT conversion preserves only the semantics identified as represented in the [conversion contract](conversion.md). Normal conversion emits complete deterministic loss warnings; `--strict` refuses publication when any known loss exists. Fatal target incompatibility fails in both modes. A successful exact restore is not evidence that a conversion or model-driven render is lossless.

## Platform support

Official release targets are Windows, macOS, and Linux on amd64 and arm64, built as pure Go with `CGO_ENABLED=0`. Cross-build success proves that a target binary can be produced. Native execution claims require the corresponding Windows, macOS, or Linux hosted test.

Core CLI, schema, parsing, rendering, conversion, validation, inspection, and completion behavior is portable. Filesystem timestamp capture and restoration legitimately differ by operating system: unsupported captured metadata warns by default, fails and rolls back with strict metadata, and is skipped with `--no-metadata`. The [architecture](architecture.md) records the native boundary.

## Installation and release boundary

Users can download v1.0.0 from GitHub and verify the selected archive with the published `cueson_1.0.0_checksums.txt` manifest. Building or running from source requires Go 1.25.0 or newer.

The v1.0.0 executable, immutable schema, release notes, annotated tag, and thirteen-asset GitHub Release are published and independently verified. Its epic and milestone are closed. S021 completed public schema hosting and the production documentation site; S022 corrected independent address-family verification. Future release and production changes remain distinct explicitly authorized steps under the [release process](release-process.md).

## Planned minor-release evolution

The [S023 roadmap](roadmap.md) targets a compatible v1.1.0 ASS/SSA addition. The [future version contract](../specs/S023-plan-scripted-format-milestone/contracts/version-compatibility.md) requires the future executable to accept exact historical v1.0.0 documents using the released local schema and version-specific semantics, preserving their input identity, producer, source truth and promised CLI behavior. Current v1.0.0 and planned v1.1.0 reject v0.0.0 input identity; the historical v0.0.0 executable remains available for its own envelopes.

New native encode output in the future release targets exact 1.1.0 identity, with current software/schema lockstep. Old exact-version consumers, including the released v1.0.0 executable, reject these new documents until they explicitly support the new contract. Historical input support on the new executable does not provide forward compatibility to old consumers. No historical-output selector or migration command is introduced by this planning slice. Any incompatible change to an established public guarantee blocks the minor target and requires a major-version decision; all compatibility claims must be proven by the documented command/version matrix and candidate evidence.

The [ASS/SSA contract](formats/ass-ssa.md) is bounded: S023 ratified planning and S024 realizes schema/model recognition. Native interpretation, stable support and later release remain separate outcomes.
