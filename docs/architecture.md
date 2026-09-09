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
| `internal/testutil` | Domain-neutral fixture loading, provenance and integrity verification, portable-path enforcement, generic golden comparisons, and explicit local-identifier rejection | Model, schema, source, codec, or CLI semantics |
| `internal/conformance` | Cross-package projections that prove schema, model, source, fixture, malformed, and path-leak behavior together | Production behavior or reusable domain APIs |
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

S004 supplied native-selectable adapter tests, ran Windows behavior natively in the development environment, and proved Linux and macOS build selection through CGO-disabled cross-compilation. S006 now supplies the hosted three-operating-system execution described in [Hosted delivery automation](#hosted-delivery-automation), after the source and test foundations removed the earlier sequencing dependency.

Windows process launches for console applications must use `CREATE_NO_WINDOW` or an equivalent hidden-process guarantee and disable interactive prompts.

## Fixture and conformance foundation

The root `testdata/manifest.json` is the only inventory authority for committed payloads beneath root `testdata/fixtures`, `testdata/malformed`, and the reserved `testdata/fuzz` tree. Each payload has one portable identity, explicit purpose and origin, redistribution approval, license and attribution, NOTICE decision, intentional byte characteristics, exact byte length, and lowercase SHA-256. Verification rejects unlisted or missing files, unsafe or colliding paths, links, non-regular files, incomplete provenance, disallowed redistribution, and byte drift before a conformance assertion begins.

Repository-authored fixture metadata and normalized JSON expectations remain UTF-8 without BOM and LF-normalized. Authoritative source, expected-byte, malformed-input, and future fuzz-input paths bypass Git text and whitespace normalization even when their current bytes happen to be valid text. Verification reads these artifacts as bytes and never normalizes or repairs them.

`internal/testutil` remains domain-neutral to avoid import cycles and incompatible format-specific assertion dialects. It compares semantic JSON, ordered diagnostics, exact bytes, integrity values, and explicit timestamp outcomes using fixture IDs and logical surface names. Cross-package tests under `internal/conformance` own projections from the Cueson model and source reports, prove exact restoration from the accepted seed, exercise parse, structure, semantics, and integrity rejection stages, and reject explicit native, slash, backslash, and JSON-escaped local identifiers.

The embedded `internal/schema/testdata/representative.cueson.json` remains the canonical schema contract example owned by `internal/schema`; it is intentionally outside the root manifest rather than an undeclared corpus payload. Initial fuzz boundaries exercise complete Cue JSON decoding, canonical source-envelope base64 integrity, and safe basenames with callback guards and no restoration writes. Native SRT and WebVTT grammar coverage remains downstream codec work. S006 supplies the hosted platform automation described below.

## Hosted delivery automation

The `CI` workflow provides stable, independently visible gates for Go formatting, repository-authored text and publication formatting, vet, schema and conformance tests, Staticcheck, govulncheck, native tests, race detection, and pure-Go cross-builds. It runs on every pull request to `main` and every `main` update with read-only repository permission. Superseded pull-request runs cancel only within the same workflow and pull request.

Native root-module tests run on explicit Linux, Windows, and macOS runner generations using the latest patch in the Go 1.25 compatibility line. These jobs own platform-behavior evidence. S006 raises the minimum from Go 1.24 because reachable standard-library vulnerabilities and the fixed `golang.org/x/text` release cannot retain that obsolete compatibility floor. A separate six-target matrix builds the executable for Windows, macOS, and Linux on amd64 and arm64 with `CGO_ENABLED=0`; those temporary binaries prove buildability only and are never uploaded or published.

Staticcheck, govulncheck, and actionlint use explicit Go module versions. External actions use immutable full commit identifiers with readable release comments, and checkout credential persistence is disabled before repository-controlled code runs. The independent `CodeQL` workflow receives only read access plus the narrow `security-events: write` permission required to publish analysis, and ordinary CI receives no write authority or secrets.

The stable check names are recorded in `specs/S006-establish-ci-gates/contracts/check-contract.md` for later repository protection work. S006 does not create empty gates for native SRT/WebVTT parsing, model-driven rendering, cross-format conversion, Codex-review automation, release packaging, or repository settings. Those names appear only after their owning implementation exists.

## Pull-request policy automation

The standalone `scripts/pr-policy` module owns deterministic issue-link and Codex-review decisions without entering the shipped product module. GitHub's fully paginated closing-issue references are authoritative for link syntax and target resolution. Normal pull requests require at least one resolved closing issue or an exact pull-request-specific exception from the configured human operator. Dependabot receives an explicit issue-link exception and remains excluded from automatic Codex review.

The policy publishes `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` commit statuses against the pull request's current head. Draft and acknowledged first reviews remain pending. A current-head terminal Codex result passes only when its thumbs-up and summary are attributable to the configured integration and no actionable finding remains unresolved. Unknown, failed, incomplete, stale, contradictory, or partially paginated evidence fails closed with a corrective description.

Round one remains the native Codex integration's responsibility. Automation may issue one marked `@codex review` request only after first-round findings are resolved. When remediation changes the head, the reviewed finding commit must be an ancestor of the newer current head and every configured S006 CI gate must pass on that head. Only an exact GitHub Actions-authored marker or configured-operator request establishes the second-round boundary; marker-shaped comments from other actors fail closed. An authenticated marked or operator-issued second request permanently consumes the automated allowance. Second-round findings, failure, or later head changes remain blocking and never create a third automatic request.

The `Pull request policy` workflow runs only trusted `main` code through `pull_request_target`, pull-request conversation comments, and a non-hourly recovery schedule. It never checks out a pull-request head or merge ref. Its token grants only repository, pull-request, and check-result read access plus issue-comment and commit-status write authority; checkout credentials and dependency caches are disabled. One non-canceling global concurrency group serializes mutations, while complete evidence refetches and mutation read-back make retries idempotent. Reaction-only completion and review-thread resolution converge through scheduled recovery, and the configured operator can request immediate recovery with `/cueson reconcile` on the pull request.

Because GitHub does not activate a newly introduced trusted-default-branch workflow for its own pull request, S007 verifies that pull request with a read-only live adapter audit and fixture-backed mutation tests. The first eligible pull request after merge must prove the hosted status source and Actions-bot comment behavior before issue #10 makes either policy context required.

## Work ownership after ratification

| Issue | Implementation ownership |
|---|---|
| [#4](https://github.com/shruggietech/cueson/issues/4) | Root Go module, process entry, CLI foundation, executable version source, root help, and `version` |
| [#5](https://github.com/shruggietech/cueson/issues/5) | Canonical v0.0.0 schema, representative document, embedding, structural and semantic validation foundation, lockstep tests, and `schema` |
| [#6](https://github.com/shruggietech/cueson/issues/6) | Source integrity, capture-before-read metadata, safe generic exact restoration, public `restore` command, and platform timestamp adapters/results |
| [#7](https://github.com/shruggietech/cueson/issues/7) | Fixture provenance, golden helpers, conformance infrastructure, malformed corpus conventions, and fuzz boundaries |
| [#8](https://github.com/shruggietech/cueson/issues/8) | Stable CI, native Windows/macOS/Linux execution, pure-Go cross-build proof, pinned analysis, vulnerability scanning, and independent CodeQL |
| [#9](https://github.com/shruggietech/cueson/issues/9) | Issue-linked pull-request policy, native Codex review reconciliation, one bounded second-round request, and trusted recovery automation |

SRT and WebVTT codecs, model-driven render, and cross-format conversion are deliberately deferred to later implementation slices. This document does not authorize placeholder commands or premature capability claims.
