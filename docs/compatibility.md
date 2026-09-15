# Cueson Compatibility Contract

**Status:** Frozen 1.1.0 stable candidate; v1.0.0 stable release published and independently verified

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

## Current candidate compatibility

The exact 1.1.0 stable candidate accepts historical 1.0.0 documents through the bundled immutable local schema and preserved version-specific semantics. Validate, inspect, restore, matching native render and established conversion paths preserve loaded identity, producer and source truth. Missing, approximate, mismatched, unknown, v0.0.0 and former 1.1.0-dev identities reject locally without network retrieval. Current discovery and every new encode, including SubRip/WebVTT, use exact 1.1.0. The published v1.0.0 executable rejects this new output until explicitly updated; historical input support on the candidate does not grant forward compatibility to old consumers. No historical-output selector or migration command is added.

The candidate freezes bounded ASS v4+/SSA v4 detection, strict UTF-8 ingest, common/native semantic ownership, textual model rendering, all twelve four-format conversion directions, validation, privacy-safe inspection and independently verified exact restoration. Official encode and private native targets declare `stable`; exact-current `schema_only` and complete `experimental` declarations remain valid observations unchanged. Stable identifies the verified bounded contract, not pixel rendering or public availability. The candidate is unpublished and its 1.1.0 schema URI is not yet publicly hosted.

## Format capability states

- `envelope_only`: the release can preserve and exactly restore a valid source envelope but does not claim native semantic ingest or model-driven rendering.
- `schema_only`: the declaration records schema/model recognition and generic restoration without claiming native ingest/render. Current source accepts this exact-current observation without promoting it to the running executable's installed capabilities.
- `experimental`: the executable implements documented native behavior whose compatibility contract has not completed its stable release gate.
- `stable`: the format has passed its applicable contract/conformance and candidate gates for that software/schema version; public release availability is stated separately.

Capability booleans are independent facts. Schema recognition does not imply an installed codec, model-driven rendering does not imply exact restoration, and exact restoration does not imply native parsing. Consumers use `format_support` for the document's declaration and the running executable's installed codec availability for an operation; neither a format key nor a file extension proves availability. Acceptance of a lower-capability declaration does not rewrite its recorded observations or imply a native operation is installed.

## Fidelity and conversion

The source envelope is authoritative for exact restoration. Common cue fields and format-native fields are structured representations for consumption and editing. Rendering serializes that structured model and never claims byte identity with the captured source.

Four-format conversion preserves only the semantics identified as represented in the [conversion contract](conversion.md). Normal conversion emits complete deterministic loss warnings; `--strict` refuses publication when any known loss exists. Fatal target incompatibility fails in both modes. A successful exact restore is not evidence that a conversion or model-driven render is lossless.

## Platform support

Official release targets are Windows, macOS, and Linux on amd64 and arm64, built as pure Go with `CGO_ENABLED=0`. Cross-build success proves that a target binary can be produced. Native execution claims require the corresponding Windows, macOS, or Linux hosted test.

Core CLI, schema, parsing, rendering, conversion, validation, inspection, and completion behavior is portable. Filesystem timestamp capture and restoration legitimately differ by operating system: unsupported captured metadata warns by default, fails and rolls back with strict metadata, and is skipped with `--no-metadata`. The [architecture](architecture.md) records the native boundary.

## Installation and release boundary

Users can download v1.0.0 from GitHub and verify the selected archive with the published `cueson_1.0.0_checksums.txt` manifest. Building or running from source requires Go 1.25.0 or newer.

The v1.0.0 executable, immutable schema, release notes, annotated tag, and thirteen-asset GitHub Release are published and independently verified. Its epic and milestone are closed. S021 completed public schema hosting and the production documentation site; S022 corrected independent address-family verification. Future release and production changes remain distinct explicitly authorized steps under the [release process](release-process.md).

## Minor-release compatibility boundary

The [S023 roadmap](roadmap.md) targets a compatible v1.1.0 ASS/SSA addition. The [ratified version contract](../specs/S023-plan-scripted-format-milestone/contracts/version-compatibility.md) requires the candidate executable to accept exact historical v1.0.0 documents using the released local schema and version-specific semantics, preserving their input identity, producer, source truth and promised CLI behavior. Released v1.0.0 and the 1.1.0 candidate reject v0.0.0 input identity; the historical v0.0.0 executable remains available for its own envelopes.

New native encode output targets exact 1.1.0 identity, with current software/schema lockstep. Old exact-version consumers, including the released v1.0.0 executable, reject these new documents until they explicitly support the new contract. Historical input support on the new executable does not provide forward compatibility to old consumers. No historical-output selector or migration command is introduced by this candidate. Any incompatible change to an established public guarantee blocks the minor target and requires a major-version decision; all compatibility claims must be proven by the documented command/version matrix and candidate evidence.

The [ASS/SSA contract](formats/ass-ssa.md) freezes the bounded native profile completed by S025-S027 and promoted in S028. Tags/releases and public schema/site hosting remain separate #66/#67 outcomes. Other format families remain explicitly deferred.
