# Cueson Compatibility Contract

**Status:** v1-bound contract frozen in v0.1.0 development source; v1.0.0 is not yet published

This document defines which Cueson surfaces receive a public compatibility promise, how format capability states are interpreted, and which release and platform claims are currently valid.

## Public v1 surfaces

The v1 compatibility promise covers exactly two public interfaces:

- The `cueson` CLI command names, options, aliases, streams, exit-code classes, overwrite rules, and strict-mode behavior documented in the [CLI contract](cli.md).
- The Cue JSON Schema identity, validation constraints, common model, native extensions, source envelope, and version rules documented in the [schema contract](schema.md).

Go packages remain under `internal/` and carry no public source-compatibility or import-compatibility promise. Repository scripts, fixtures, conformance indexes, inspection-report internals, conversion target projections, and diagnostic prose are verification or implementation surfaces unless another maintained contract explicitly says otherwise. Stable diagnostic and conversion-loss codes remain machine-consumable where their owning contract identifies them as stable.

## Version policy

The published v0.0.0 release is an immutable envelope-only foundation. Its executable, schema, archive set, and versioned schema copy remain unchanged and do not contain native SubRip or WebVTT ingest, rendering, conversion, validation, inspection, or completion.

Current development source uses executable and schema version v0.1.0. It implements the complete intended v1 command set and native SubRip and WebVTT workflows, but its format capability status remains `experimental` until the separately reviewed v1 candidate changes the release identity and capability declarations together.

At v1.0.0, incompatible changes to either public interface require a new major version. Compatible additions may use a minor version, and compatible corrections may use a patch version. Before v1.0.0, a documented breaking schema or CLI change requires at least a development minor-version change. Official Cueson executable and schema versions remain equal; a third-party producer's own version is independent from the Cue JSON schema version it targets.

No mutable `latest` schema identity is part of the contract. Every released schema copy is immutable. The v0.1.0 canonical schema is available from the repository and embedded executable; its `cueson.io` URI remains an identifier until production hosting is separately authorized.

## Format capability states

- `envelope_only`: the release can preserve and exactly restore a valid source envelope but does not claim native semantic ingest or model-driven rendering.
- `experimental`: the executable implements the documented native behavior, but the stable release gate has not yet been completed.
- `stable`: the format has passed the applicable v1 acceptance and release gates for that software/schema version.

Capability booleans are independent facts. Schema recognition does not imply an installed codec, model-driven rendering does not imply exact restoration, and exact restoration does not imply native parsing. Consumers must use `format_support` rather than infer capability from a format key or file extension.

## Fidelity and conversion

The source envelope is authoritative for exact restoration. Common cue fields and format-native fields are structured representations for consumption and editing. Rendering serializes that structured model and never claims byte identity with the captured source.

SubRip and WebVTT conversion preserves only the semantics identified as represented in the [conversion contract](conversion.md). Normal conversion emits complete deterministic loss warnings; `--strict` refuses publication when any known loss exists. Fatal target incompatibility fails in both modes. A successful exact restore is not evidence that a conversion or model-driven render is lossless.

## Platform support

Official release targets are Windows, macOS, and Linux on amd64 and arm64, built as pure Go with `CGO_ENABLED=0`. Cross-build success proves that a target binary can be produced. Native execution claims require the corresponding Windows, macOS, or Linux hosted test.

Core CLI, schema, parsing, rendering, conversion, validation, inspection, and completion behavior is portable. Filesystem timestamp capture and restoration legitimately differ by operating system: unsupported captured metadata warns by default, fails and rolls back with strict metadata, and is skipped with `--no-metadata`. The [architecture](architecture.md) records the native boundary.

## Installation and release boundary

Users who need the published foundation can download v0.0.0 from GitHub and verify it with the published checksum manifest. Users evaluating the v1-bound workflows must currently build or run v0.1.0 from source with Go 1.25.0 or newer.

A v1.0.0 binary, immutable v1 schema, release notes, tag, GitHub Release, public schema endpoint, and production documentation site do not exist yet. Contract hardening does not authorize those actions. Candidate preparation, publication, milestone closure, schema hosting, and any production `cueson.io` change remain distinct steps under the [release process](release-process.md).
