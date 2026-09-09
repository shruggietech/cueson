# Cueson Architecture

**Status:** Ratified v0.0.0 implementation baseline

**Ratified:** 2026-09-09 through Spec Kit slice `001-ratify-foundation-contracts`

This document is the architecture of record for the v0.0.0 foundation. The [project constitution](../.specify/memory/constitution.md) remains the highest repository authority. The [working project specification](Cueson-Project-Specification-v0.0.0.md) supplies broader context and roadmap detail when it does not conflict with this baseline.

## Public and internal boundaries

Cueson's approved public v1 contracts are the `cueson` command-line interface and Cue JSON Schema. Go packages remain under `internal/` and make no public compatibility promise until a separate specification approves a library API.

The executable entry point under `cmd/cueson` is a minimal operating-system adapter. It passes context, arguments, standard input, standard output, and standard error to `internal/cli`, then exits with the returned status. Business and format packages do not import the CLI package.

## Responsibility map

| Boundary | Responsibility | Must not own |
|---|---|---|
| `cmd/cueson` | Process entry and exit | Argument policy, domain behavior, or format parsing |
| `internal/cli` | Command registration, argument parsing, help, stream and color policy, typed-error mapping | Cue model rules or source restoration |
| `internal/version` | Single executable build-version source | Schema identity as a second mutable version source |
| `internal/model` | Cue JSON types and format-neutral semantic invariants | Filesystem I/O, command behavior, or codec selection |
| `internal/schema` | Embedded schema bytes, identity and version access, structural validation, and schema/software lockstep checks | Native codec capability inference |
| `internal/source` | Source bundles, safe basenames, byte length and SHA-256 verification, metadata capture ordering, exact restoration, and platform timestamp results | Native format parsing or model-driven rendering |
| `internal/codec` | Format detection plus native ingest and model-driven render implementations | Exact source restoration |
| `internal/convert` | Model-driven conversion and explicit loss accounting | Source-envelope restoration |
| `internal/ocr` | Future OCR provider boundary for derived recognition | Source truth or codec-specific parsing |

Dependencies point inward toward stable, dependency-light contracts. `internal/model` does not depend on CLI, source, codecs, conversion, or OCR. Codec availability is the authority for native ingest and render capability; schema recognition alone is not.

## Source authority and restoration

Original source asset bytes are authoritative for exact restoration. Normalized cues, format-native parsed data, speaker observations, token timing, and OCR observations are derived or interpreted surfaces and never replace the source envelope.

Exact restoration is not a codec operation. It validates and writes assets directly from the source envelope, verifies byte length and SHA-256 before acceptance, and applies supported filesystem timestamps only after the final bytes and output name are in place. This separation allows a conforming Cue JSON document to restore assets even when the executable has no native codec for the document's format.

Restoration preflights the complete bundle, writes every asset to a same-directory private staging file, flushes and reopens each stage for integrity verification, and publishes only after every stage is ready. New destinations use exclusive hard-link publication so a raced-in path is never replaced. Forced regular-file destinations receive a same-filesystem hard-link rollback copy before atomic replacement. Controlled failures roll committed assets back in reverse order by recorded file identity; cleanup failure after the complete transaction is accepted is reported as a warning. This boundary does not claim crash atomicity or isolation from unrelated processes mutating the same paths.

Paths stored in Cue JSON are portable safe basenames only. Runtime input paths, directories, drive or mount information, hostnames, usernames, and other machine identifiers never enter the document. Restoration rejects path separators, traversal components, control characters, Windows-invalid punctuation and trailing characters, drive or URI prefixes, NTFS alternate-data-stream syntax, and case-insensitive Windows reserved device names even when they have extensions.

Before opening any output, restoration builds the complete destination plan and rejects duplicate destinations. Stored basenames in one source bundle must be unique under Unicode canonical caseless matching, defined as NFD normalization, default Unicode case folding, then NFD normalization again, so case-insensitive and normalization-insensitive filesystems cannot collapse distinct assets. Default destinations and paths beneath `--output-dir` derive from stored basenames. A single-asset `--output` is an explicit caller-selected runtime path, is validated separately, and may rename the restored file.

## Common model and native fidelity

The common model exposes timing and semantic content without requiring consumers to decode source bytes or parse a subtitle grammar. Format-native data remains adjacent to the common model wherever normalization would otherwise lose information.

Unknown, malformed, unsupported, or non-representable source information is preserved when practical and reported through deterministic diagnostics. Strict conversion refuses known loss. A successful operation never silently discards information covered by the active contract.

## Version and schema relationship

The executable version has one build-version source under `internal/version`. The canonical schema has its own embedded identity and `schema_version`, and the build verifies equality rather than deriving one mutable value from the other. Official release success requires executable version, embedded schema version, and release tag version to agree.

Third-party producer software versions are independent from the Cue JSON schema version they target. The schema records that producer identity without weakening official Cueson release lockstep.

## CLI and domain error boundary

Domain packages return errors that retain enough type or classification for `internal/cli` to select a stable exit code and diagnostic. They do not print directly. The CLI owns help, quiet/silent filtering, color, stream selection, and final process status.

Invalid invocation and missing pre-execution requirements use exit code 2. Failures after a valid shipped command is accepted, including missing runtime capability, parsing, validation, integrity, rendering, conversion, and I/O, use exit code 1. Success uses exit code 0.

## Platform boundary

The target remains pure Go with `CGO_ENABLED=0` unless a separately recorded decision proves a native dependency necessary. Timestamp capture and restoration use platform-selected internal adapters because portable APIs cannot make identical claims on Windows, macOS, and Linux. Native tests prove supported behavior; cross-compilation alone is insufficient.

Windows captures and attempts creation, modification, and access timestamps through no-follow handles at 100-nanosecond precision. Linux captures descriptor-visible modification and access timestamps, captures birth time when `statx` exposes it, and reports birth-time setting as unsupported. The macOS v0.0.0 adapter captures birth, modification, and access timestamps but restores only descriptor-based modification and access at its effective microsecond precision, reporting creation restoration as unsupported. Default restore warns for unsupported captured values, strict mode rolls back, and no-metadata mode makes no timestamp claim.

S004 supplies native-selectable adapter tests, runs Windows behavior natively in the development environment, and proves Linux and macOS build selection through CGO-disabled cross-compilation. Hosted native execution on all three operating systems remains issue [#8](https://github.com/shruggietech/cueson/issues/8), which is intentionally downstream of the source and test foundations.

Windows process launches for console applications must use `CREATE_NO_WINDOW` or an equivalent hidden-process guarantee and disable interactive prompts.

## Work ownership after ratification

| Issue | Implementation ownership |
|---|---|
| [#4](https://github.com/shruggietech/cueson/issues/4) | Root Go module, process entry, CLI foundation, executable version source, root help, and `version` |
| [#5](https://github.com/shruggietech/cueson/issues/5) | Canonical v0.0.0 schema, representative document, embedding, structural and semantic validation foundation, lockstep tests, and `schema` |
| [#6](https://github.com/shruggietech/cueson/issues/6) | Source integrity, capture-before-read metadata, safe generic exact restoration, public `restore` command, and platform timestamp adapters/results |
| [#7](https://github.com/shruggietech/cueson/issues/7) | Fixture provenance, golden helpers, conformance infrastructure, malformed corpus conventions, and fuzz boundaries |

SRT and WebVTT codecs, model-driven render, and cross-format conversion are deliberately deferred to later implementation slices. This document does not authorize placeholder commands or premature capability claims.
