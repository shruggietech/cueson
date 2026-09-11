# Cueson Architecture

**Status:** v1.0.0 stable released architecture

**Ratified:** 2026-09-09 through Spec Kit slice `001-ratify-foundation-contracts`

This document is the architecture of record for the stable v1.0.0 release. The [project constitution](../.specify/memory/constitution.md) remains the highest repository authority. The [working project specification](Cueson-Project-Specification-v0.0.0.md) supplies broader context and roadmap detail when it does not conflict with ratified slices.

## Public and internal boundaries

Cueson's approved public v1 contracts are the `cueson` command-line interface and Cue JSON Schema. Go packages remain under `internal/` and make no public compatibility promise until a separate specification approves a library API. The maintained [compatibility contract](compatibility.md) defines version, format-state, platform, fidelity, and release boundaries for those surfaces.

The executable entry point under `cmd/cueson` is a minimal operating-system adapter. It passes context, arguments, standard input, standard output, and standard error to `internal/cli`, then exits with the returned status. Business and format packages do not import the CLI package.

## Responsibility map

| Boundary | Responsibility | Must not own |
|---|---|---|
| `cmd/cueson` | Process entry and exit | Argument policy, domain behavior, or format parsing |
| `internal/cli` | Command registration, argument parsing, validated-input classification, inspection reports, static completion, help, stream and color policy, typed-error mapping | Cue model rules or source restoration |
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

Current source implements `internal/codec` as a capability registry plus bounded detection/decoding and native SubRip and WebVTT parser/renderers. `internal/convert` owns target projection, compatibility analysis, deterministic runtime-only loss reports, and conversion rendering. `internal/ocr` remains reserved future ownership.

Validation, inspection, and conversion share one CLI-owned validated-input path so Cue JSON precedence, native content-first selection, schema and semantic checks, source integrity, and codec availability cannot drift among commands. Inspection projects that result into a fixed privacy-bounded report rather than exposing the Cue model directly. Help and four static shell completion definitions derive their public vocabulary from one ordered CLI surface catalogue while semantic option conflicts remain in the explicit parser.

## Source authority and restoration

Original source asset bytes are authoritative for exact restoration. Normalized cues, format-native parsed data, speaker observations, token timing, and OCR observations are derived or interpreted surfaces and never replace the source envelope.

Exact restoration is not a codec operation. It validates and writes assets directly from the source envelope, verifies byte length and SHA-256 before acceptance, and applies supported filesystem timestamps only after the final bytes and output name are in place. This separation allows a conforming Cue JSON document to restore assets even when the executable has no native codec for the document's format.

Restoration preflights the complete bundle, writes every asset to a same-directory private staging file, flushes and reopens each stage for integrity verification, and publishes only after every stage is ready. New destinations use exclusive hard-link publication so a raced-in path is never replaced. Forced regular-file destinations receive a same-filesystem hard-link rollback copy before atomic replacement. Controlled failures roll committed assets back in reverse order by recorded file identity; cleanup failure after the complete transaction is accepted is reported as a warning. This boundary does not claim crash atomicity or isolation from unrelated processes mutating the same paths.

Paths stored in Cue JSON are portable safe basenames only. Runtime input paths, directories, drive or mount information, hostnames, usernames, and other machine identifiers never enter the document. Restoration rejects path separators, traversal components, control characters, Windows-invalid punctuation and trailing characters, drive or URI prefixes, NTFS alternate-data-stream syntax, and case-insensitive Windows reserved device names even when they have extensions.

Before opening any output, restoration builds the complete destination plan and rejects duplicate destinations. Stored basenames in one source bundle must be unique under Unicode canonical caseless matching, defined as NFD normalization, default Unicode case folding, then NFD normalization again, so case-insensitive and normalization-insensitive filesystems cannot collapse distinct assets. Default destinations and paths beneath `--output-dir` derive from stored basenames. A single-asset `--output` is an explicit caller-selected runtime path, is validated separately, and may rename the restored file.

## Common model and native fidelity

The common model exposes timing and semantic content without requiring consumers to decode source bytes or parse a subtitle grammar. Format-native data remains adjacent wherever normalization would otherwise lose information. Current SubRip and WebVTT encode workflows construct both views from one bounded exact source acquisition.

Unknown, malformed, unsupported, or non-representable source information is preserved when practical and reported through deterministic diagnostics. Strict conversion refuses known loss. A successful operation never silently discards information covered by the active contract.

The dedicated [SubRip](formats/srt.md) and [WebVTT](formats/webvtt.md) pages define their stable v1 native grammars, diagnostics, fixtures, ingest, rendering, and conversion contracts. Conversion validates source-envelope integrity, projects only model data into a private target representation, computes its complete loss report before rendering, and keeps that target representation and report outside Cue JSON. Strict conversion rejects every known loss before output publication.

## Version and schema relationship

The executable version has one build-version source under `internal/version`. The canonical schema has its own embedded identity and `schema_version`, and the build verifies equality rather than deriving one mutable value from the other. Official release success requires executable version, embedded schema version, and release tag version to agree.

Third-party producer software versions are independent from the Cue JSON schema version they target. The schema records that producer identity without weakening official Cueson release lockstep.

## CLI and domain error boundary

Domain packages return errors that retain enough type or classification for `internal/cli` to select a stable exit code and diagnostic. They do not print directly. The CLI owns help, quiet/silent filtering, color, stream selection, and final process status.

The CLI diagnostic writer retains output failures. A command that otherwise succeeded but could not deliver a promised success or warning diagnostic becomes a runtime failure; an existing invocation or runtime failure keeps its original status class.

Invalid invocation and missing pre-execution requirements use exit code 2. Failures after a valid shipped command is accepted, including missing runtime capability, parsing, validation, integrity, rendering, conversion, and I/O, use exit code 1. Success uses exit code 0.

## Platform boundary

The target remains pure Go with `CGO_ENABLED=0` unless a separately recorded decision proves a native dependency necessary. Timestamp capture and restoration use platform-selected internal adapters because portable APIs cannot make identical claims on Windows, macOS, and Linux. Native tests prove supported behavior; cross-compilation alone is insufficient.

Windows captures and attempts creation, modification, and access timestamps through no-follow handles at 100-nanosecond precision. Linux captures descriptor-visible modification and access timestamps, captures birth time when `statx` exposes it, and reports birth-time setting as unsupported. The macOS v0.0.0 adapter captures birth, modification, and access timestamps but restores only descriptor-based modification and access at its effective microsecond precision, reporting creation restoration as unsupported. Default restore warns for unsupported captured values, strict mode rolls back, and no-metadata mode makes no timestamp claim.

S004 supplied native-selectable adapter tests, ran Windows behavior natively in the development environment, and proved Linux and macOS build selection through CGO-disabled cross-compilation. S006 now supplies the hosted three-operating-system execution described in [Hosted delivery automation](#hosted-delivery-automation), after the source and test foundations removed the earlier sequencing dependency.

Windows process launches for console applications must use `CREATE_NO_WINDOW` or an equivalent hidden-process guarantee and disable interactive prompts.

## Brand source boundary

The operator-designated ShruggieTech Cueson brand-kit archive is retained byte for byte under `brand/cueson/1.0.0/archive`, and its complete safe extraction is retained under the adjacent `kit` directory. The Cueson-owned import manifest records the acquired archive identity, every extracted regular file, and every repository document that directly consumes a retained asset. The standalone `scripts/brand-verify` module independently proves archive, catalog, extraction, and consumer-reference identity without network access or repository writes.

Imported brand content is immutable source evidence. Repository formatting and repair tooling excludes protected brand paths, and later kit revisions use new versioned directories rather than mutating or mixing snapshots. Repository documentation consumes selected assets directly from the retained kit so there is no duplicate consumer-copy authority.

Brand assets remain repository source material rather than product runtime or release-archive members. Their import does not authorize changes to `cueson.io`, production metadata, repository social-preview settings, tags, releases, schemas, or other production surfaces. The detailed provenance, usage, licensing, verification, and update contract is recorded in the [brand guide](brand.md).

## Fixture and conformance foundation

The root `testdata/manifest.json` is the only inventory authority for committed payloads beneath root `testdata/fixtures`, `testdata/malformed`, and the reserved `testdata/fuzz` tree. Each payload has one portable identity, explicit purpose and origin, redistribution approval, license and attribution, NOTICE decision, intentional byte characteristics, exact byte length, and lowercase SHA-256. Verification rejects unlisted or missing files, unsafe or colliding paths, links, non-regular files, incomplete provenance, disallowed redistribution, and byte drift before a conformance assertion begins.

Repository-authored fixture metadata and normalized JSON expectations remain UTF-8 without BOM and LF-normalized. Authoritative source, expected-byte, malformed-input, and future fuzz-input paths bypass Git text and whitespace normalization even when their current bytes happen to be valid text. Verification reads these artifacts as bytes and never normalizes or repairs them.

`internal/testutil` remains domain-neutral to avoid import cycles and incompatible format-specific assertion dialects. It compares semantic JSON, ordered diagnostics, exact bytes, integrity values, and explicit timestamp outcomes using fixture IDs and logical surface names. Cross-package tests under `internal/conformance` own projections from the Cueson model and source reports, prove exact restoration from the accepted seed, exercise parse, structure, semantics, and integrity rejection stages, and reject explicit native, slash, backslash, and JSON-escaped local identifiers.

The embedded `internal/schema/testdata/representative.cueson.json` remains the canonical schema contract example owned by `internal/schema`; it is intentionally outside the root manifest rather than an undeclared corpus payload. Generic fuzz boundaries exercise complete Cue JSON decoding, canonical source-envelope base64 integrity, and safe basenames with callback guards and no restoration writes. Native SubRip and WebVTT packages add format-specific parsing, rendering, and parser-cycle fuzz boundaries. S006 supplies the hosted platform automation described below.

## Hosted delivery automation

The `CI` workflow provides stable, independently visible gates for Go formatting, repository-authored text and publication formatting, vet, schema and conformance tests, Staticcheck, govulncheck, native tests, race detection, and pure-Go cross-builds. It runs on every pull request to `main` and every `main` update with read-only repository permission. Superseded pull-request runs cancel only within the same workflow and pull request.

Native root-module tests run on explicit Linux, Windows, and macOS runner generations using the latest patch in the Go 1.25 compatibility line. These jobs own platform-behavior evidence. S006 raises the minimum from Go 1.24 because reachable standard-library vulnerabilities and the fixed `golang.org/x/text` release cannot retain that obsolete compatibility floor. A separate six-target matrix builds the executable for Windows, macOS, and Linux on amd64 and arm64 with `CGO_ENABLED=0`; those temporary binaries prove buildability only and are never uploaded or published.

Staticcheck, govulncheck, and actionlint use explicit Go module versions. External actions use immutable full commit identifiers with readable release comments, and checkout credential persistence is disabled before repository-controlled code runs. The independent `CodeQL` workflow receives only read access plus the narrow `security-events: write` permission required to publish analysis, and ordinary CI receives no write authority or secrets.

The stable check names are recorded in `specs/S006-establish-ci-gates/contracts/check-contract.md` and the 17 CI and CodeQL contexts are now active required checks on `main`. S006 did not create empty gates for native SRT/WebVTT parsing, model-driven rendering, or cross-format conversion. Those names appear only after their owning implementation exists.

## Pull-request policy automation

The standalone `scripts/pr-policy` module owns deterministic issue-link and Codex-review decisions without entering the shipped product module. GitHub's fully paginated closing-issue references are authoritative for link syntax and target resolution. Normal pull requests require at least one resolved closing issue or an exact pull-request-specific exception from the configured human operator. Dependabot receives an explicit issue-link exception and remains excluded from automatic Codex review.

The policy publishes `Cueson PR policy / Issue link` and `Cueson PR policy / Codex review` commit statuses against the pull request's current head. Draft and acknowledged first reviews remain pending. A current-head terminal Codex result passes only when its thumbs-up and summary are attributable to the configured integration and no actionable finding remains unresolved. Unknown, failed, incomplete, stale, contradictory, or partially paginated evidence fails closed with a corrective description.

Round one remains the native Codex integration's responsibility. Automation may issue one marked `@codex review` request only after first-round findings are resolved. When remediation changes the head, the reviewed finding commit must be an ancestor of the newer current head and every configured S006 CI gate must pass on that head. Only an exact GitHub Actions-authored marker or configured-operator request establishes the second-round boundary; marker-shaped comments from other actors fail closed. An operator request must cite one backticked commit prefix that resolves uniquely in the pull-request head history, use its latest edit time as the boundary, and acquire an authenticated durable reservation when first observed, so it cannot be silently rebound or forgotten after a later edit, deletion, or push. An authenticated marked or operator-issued second request permanently consumes the automated allowance. A resolved second-round finding may complete on a proven descendant only when every second-round thread is resolved and all configured S006 gates pass on that descendant. Failed, unresolved, ancestry-ambiguous, or CI-incomplete second-round evidence remains blocking and never creates a third automatic request.

The `Pull request policy` workflow runs only trusted `main` code through `pull_request_target`, pull-request conversation comments, and a non-hourly recovery schedule. It never checks out a pull-request head or merge ref. Its token grants only repository, pull-request, and check-result read access plus issue-comment and commit-status write authority; checkout credentials and dependency caches are disabled. One non-canceling global concurrency group serializes mutations, while complete evidence refetches and mutation read-back make retries idempotent. Reaction-only completion and review-thread resolution converge through scheduled recovery, and the configured operator can request immediate recovery with `/cueson reconcile` on the pull request.

Subsequent hosted pull requests proved current-head publication by the GitHub Actions identity for both policy statuses. The permitted GitHub Actions-authored second-round comment path returned HTTP 403 during S009, so the operator supplied the single manual second-round request and recorded a waiver after remediation. The automatic comment path therefore remains unproven, and neither policy status is protection-required.

## Repository protection

Repository protection is cumulative. The existing organization-owned default-branch ruleset remains an immutable S008 input, while Cueson-specific resolved-conversation and required-check policy belongs to a repository-owned ruleset targeting only `main`. Classic branch protection is not treated as the sole authority when rulesets supply the effective policy.

Required checks enter protection only after the exact context, pull-request head, successful terminal result, and provider identity are read back from GitHub. The repository-owned ruleset requires the 17 verified S006 CI and CodeQL contexts on the current base. Both S007 commit-status contexts remain outside protection because the GitHub Actions-authored second-round request path lacks complete activation proof.

Repository Actions default to read permission, pull-request approval remains disabled, action sources are limited to GitHub-owned actions, and immutable action references are required. Versioned workflows continue to declare their own minimum permissions. The versioned CodeQL workflow remains authoritative rather than enabling a competing default setup.

One organization-administrator bypass on the repository-owned ruleset preserves recovery from a misconfigured or renamed gate. It is not ordinary delivery authority: pull requests still follow the bounded review protocol, and final merge remains a human decision.

## Release proof and publication boundary

The repository-root GoReleaser v2 configuration is intentionally snapshot-only. It builds `cueson` with `CGO_ENABLED=0` for Windows, macOS, and Linux on amd64 and arm64, injects the `internal/version` release override with a verifier-readable marker consumed by the public version surface, uses commit-derived timestamps and trimmed build paths, packages the byte-identical immutable v1 schema plus legal files, emits one SHA-256 archive manifest, and asks a pinned Syft command to generate one target-bound SPDX JSON SBOM from each packaged binary. Each SBOM is named for its corresponding archive. Generating from the binary avoids leaking Syft's temporary archive-extraction paths. `release.disable: true` prevents the same configuration from becoming a publishing path when snapshot mode is omitted.

The standalone `scripts/release-verify` module is the release-candidate acceptance authority. It remains outside the shipped product dependency graph and uses only the Go standard library. It opens every ZIP and tar.gz archive, enforces the exact six-target and four-member contracts, verifies checksum bijection, schema bytes, and the release-version marker in every binary, inspects Go build information for target, pure-Go, source revision, clean state, trimmed paths, and local identifiers, validates SBOM source identity, and executes only the current host's compatible packaged binary. It scans semantic metadata rather than arbitrary compressed or executable bytes so binary data cannot create false path matches. It writes deterministic semantic evidence only after every check passes.

The `Release proof` workflow runs repository code through an ordinary unprivileged pull-request event with read-only repository permission, exact build-tool versions, no secrets, and no publication token. It builds and structurally verifies one accepted six-target bundle, then executes that same bundle on hosted Linux, Windows, and macOS runners. Its short-lived GitHub Actions artifact is review evidence rather than a GitHub Release asset. This is distinct from S006's temporary cross-build binaries and does not change the current required-check ruleset automatically.

S009 established candidate packaging, and S012 admitted the immutable schema and bound the accepted post-squash candidate to `b294a6952c8bd041d852c502f5d7206c0b58edd6`. S013 used explicit operator authority to create annotated tag [`v0.0.0`](https://github.com/shruggietech/cueson/tree/v0.0.0), publish the [thirteen-asset GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0), and verify every public file against accepted evidence. S019 promoted the frozen contracts to identity 1.0.0, admitted the immutable v1 schema, and proved exact schema and legal-file bytes plus native packaged-binary execution at `2cad4c816340404289b4d1d87179a4071713bb46`. S020 created annotated tag [`v1.0.0`](https://github.com/shruggietech/cueson/tree/v1.0.0), published the exact [thirteen-asset GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0), and independently verified every public download against accepted run [34621429626](https://github.com/shruggietech/cueson/actions/runs/34621429626). Syft SBOMs are checked for stable meaning and source binding; they are not claimed byte-for-byte reproducible across rebuilds while upstream output includes variable timestamps and document identifiers.

The runnable artifact and public-byte evidence remain in [release verification](release-verification.md). The broader [release process](release-process.md) separates candidate preparation and authorized GitHub publication from milestone closure, signatures, attestations, public schema hosting, and production actions.

## Maintained documentation verification

The standalone `scripts/docs-verify` module verifies the canonical documentation inventory, registered README workflow set, required schema and compatibility markers, stale capability claims, conformance-row linkage, and maintained local file, directory, image, and heading-fragment links without network access. It remains outside the shipped product module and performs no repair. Root-module CLI tests execute the registered workflows with governed fixtures, temporary destinations, and captured streams. `scripts/github-format` separately owns repository text encoding, line endings, mojibake rejection, Markdown source layout, and GitHub publication formatting.

The `Repository text` CI job runs both standalone modules. Historical Spec Kit artifacts remain formatter-governed delivery evidence but are not part of the maintained-document link graph.

## Work ownership after ratification

| Issue | Implementation ownership |
|---|---|
| [#4](https://github.com/shruggietech/cueson/issues/4) | Root Go module, process entry, CLI foundation, executable version source, root help, and `version` |
| [#5](https://github.com/shruggietech/cueson/issues/5) | Canonical v0.0.0 schema, representative document, embedding, structural and semantic validation foundation, lockstep tests, and `schema` |
| [#6](https://github.com/shruggietech/cueson/issues/6) | Source integrity, capture-before-read metadata, safe generic exact restoration, public `restore` command, and platform timestamp adapters/results |
| [#7](https://github.com/shruggietech/cueson/issues/7) | Fixture provenance, golden helpers, conformance infrastructure, malformed corpus conventions, and fuzz boundaries |
| [#8](https://github.com/shruggietech/cueson/issues/8) | Stable CI, native Windows/macOS/Linux execution, pure-Go cross-build proof, pinned analysis, vulnerability scanning, and independent CodeQL |
| [#9](https://github.com/shruggietech/cueson/issues/9) | Issue-linked pull-request policy, native Codex review reconciliation, one bounded second-round request, and trusted recovery automation |
| [#10](https://github.com/shruggietech/cueson/issues/10) | Active repository-owned `main` rules, verified required checks, Actions defaults, action-source restrictions, and recovery bypass |
| [#11](https://github.com/shruggietech/cueson/issues/11) | Non-publishing six-target candidate packaging, checksums, SBOM generation, and standalone artifact verification |
| [#12](https://github.com/shruggietech/cueson/issues/12) | Canonical documentation completion, offline link verification, and v0.0.0 milestone-readiness evidence |
| [#23](https://github.com/shruggietech/cueson/issues/23) | Complete official brand-kit retention, offline integrity verification, and repository brand integration |
| [#25](https://github.com/shruggietech/cueson/issues/25) | Immutable v0.0.0 schema admission, dated release records, and exact post-squash default-branch candidate proof |
| [#27](https://github.com/shruggietech/cueson/issues/27) | Authorized v0.0.0 tag and GitHub Release publication, independent public-download verification, and released-state reconciliation |
| [#30](https://github.com/shruggietech/cueson/issues/30) | Capability-based codec registry, bounded source acquisition, detection, decoding, diagnostics, and development schema/model transition |
| [#31](https://github.com/shruggietech/cueson/issues/31) | Native SubRip parsing, exact-envelope encode, canonical rendering, CLI workflows, fixtures, and round-trip verification |
| [#32](https://github.com/shruggietech/cueson/issues/32) | Native WebVTT parsing, exact-envelope encode, canonical rendering, CLI workflows, fixtures, diagnostics, and round-trip verification |
| [#33](https://github.com/shruggietech/cueson/issues/33) | Bidirectional SubRip and WebVTT conversion, deterministic loss accounting, strict rejection, CLI workflows, and cross-format verification |
| [#34](https://github.com/shruggietech/cueson/issues/34) | Complete validation, privacy-bounded inspection, static completion, generated help, and shared validated-input workflows |
| [#35](https://github.com/shruggietech/cueson/issues/35) | Whole-system conformance, hostile-input bounds, corpus and fuzz evidence, platform proof, and release-readiness hardening |
| [#36](https://github.com/shruggietech/cueson/issues/36) | Executable installation, CLI, format, conversion, compatibility, security, and release-transition documentation |
| [#37](https://github.com/shruggietech/cueson/issues/37) | Stable v1.0.0 identity, immutable v1 schema, release records, candidate artifacts, and native packaged-binary proof |
| [#38](https://github.com/shruggietech/cueson/issues/38) | Authorized v1.0.0 tag and GitHub Release publication, independent public-download verification, and released-state reconciliation |
| [#41](https://github.com/shruggietech/cueson/issues/41) | Consumer-facing canonical-schema descriptions, titles, examples, and automated annotation coverage |

The SubRip, WebVTT, and cross-format contracts frozen by S018 carry stable published identity 1.0.0 through S019 and S020. This document records the verified release but does not authorize public schema hosting, milestone closure, signatures, attestations, or production changes.
