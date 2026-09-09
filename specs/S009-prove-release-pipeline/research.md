# Research: Non-Publishing Release Proof

## Decision: Pin GoReleaser v2.18.1 and Syft v1.51.1 as build-only tools

**Rationale:** GoReleaser v2.18.1 and Syft v1.51.1 are the current stable upstream releases on 2026-09-09. Exact versions make local and hosted behavior reviewable and avoid floating tool resolution. Both remain release tooling and do not enter the product module or shipped binary.

**Alternatives considered:** Floating `latest` or `~> v2` selectors were rejected because the same source revision could produce structurally different evidence over time. Adding these large tool dependency graphs to the root module was rejected because they are not product dependencies.

## Decision: Install pinned commands instead of using third-party workflow actions

**Rationale:** Cueson's repository Actions policy allows GitHub-owned actions only and requires immutable action references. The GoReleaser and Syft actions are therefore unavailable even though they are upstream-supported. Hosted automation will use exact `go install module@version` commands after a pinned GitHub-owned Go setup action, while local maintainers may install the same exact versions.

**Alternatives considered:** Relaxing the repository action allowlist was rejected as an S008 policy regression outside issue #11. Downloading unchecked binaries with shell installers was rejected because it would weaken supply-chain provenance. Building a replacement release orchestrator was rejected because GoReleaser already expresses the approved archive contract.

## Decision: Use GoReleaser snapshot mode with a fixed v0.0.0 identity and disabled release publishing

**Rationale:** GoReleaser documents `--snapshot` for complete local or CI builds generated only into `dist/` and not uploaded. A `snapshot.version_template` of `0.0.0` proves the actual v0.0.0 filenames and injected executable version before an authorized `v0.0.0` tag exists. `release.disable: true` adds defense in depth if a future caller omits snapshot mode.

**Alternatives considered:** Building from a temporary local tag was rejected because tag creation is outside S009 authority and would make the safety proof depend on cleanup. Accepting a generated `0.0.0-SNAPSHOT-<commit>` identity was rejected because it would not prove the intended v0.0.0 artifact names or version lockstep.

## Decision: Package ZIP on Windows and tar.gz on Unix targets

**Rationale:** GoReleaser supports per-operating-system archive-format overrides. ZIP is native to Windows consumers, while tar.gz is conventional for macOS and Linux. All archives expose the same four root members: the target binary, `cueson.schema.json`, `LICENSE`, and `NOTICE`.

**Alternatives considered:** One format for every platform was rejected because it reduces usability without improving verification. Wrapping members in a variable directory was rejected because a stable flat four-member contract is easier to inspect and consume.

## Decision: Generate one target-bound SPDX JSON SBOM from each packaged binary with Syft

**Rationale:** GoReleaser's SBOM pipe invokes Syft for every packaged target binary and names the output for the corresponding archive. SPDX JSON is machine-readable, standardized, and supported by the toolchain. Cataloging the binary avoids Syft temporary archive-extraction paths in the document. Each SBOM receives a target-specific source name and a source version containing v0.0.0 plus the full commit, while GoReleaser metadata and release evidence bind all targets to that revision.

**Alternatives considered:** Omitting SBOMs was rejected because the working project specification and issue context explicitly expect them where supported. Claiming GitHub artifact attestations was rejected because S009 has no `id-token: write`, attestation action, signature, or publication authority.

## Decision: Use a dependency-free repository-owned verifier as the acceptance authority

**Rationale:** GoReleaser's successful exit proves pipeline completion but does not alone prove Cueson's exact member, checksum, schema, lockstep, path-hygiene, and publication assertions. A standalone Go module can parse ZIP, tar.gz, JSON, Go build information, and SHA-256 output consistently on every host without entering the shipped product dependency graph.

**Alternatives considered:** Platform-specific shell inspection was rejected because archive tools and quoting differ across Windows and Unix. Trusting filenames and metadata without opening archives was rejected because missing or unsafe members would remain undetected.

## Decision: Separate structural matrix checks from host-compatible execution

**Rationale:** One host cannot execute every cross-compiled binary. The verifier validates all six binaries structurally, requires the release-version marker consumed by the public version surface, and uses Go build information for target and revision evidence, while executing `cueson version` and `cueson schema --version` only for the current host's compatible archive. Hosted Ubuntu proves Linux amd64; the foreground Windows run proves Windows amd64.

**Alternatives considered:** Emulation was rejected as disproportionate and a new trusted dependency. Skipping all execution was rejected because public CLI surfaces must be proven from a packaged binary on compatible hosts.

## Decision: Keep the workflow read-only and artifact-only

**Rationale:** The workflow runs on pull requests and explicit dispatch with `contents: read`, no secrets, no `GITHUB_TOKEN` passed to release tooling, no signing or attestation permissions, and no release command lacking `--snapshot`. A pinned GitHub-owned artifact action may retain `dist/` for review, but it cannot create a GitHub Release.

**Alternatives considered:** `contents: write`, tag or release events, environment secrets, and a combined publish workflow were rejected because issue #11 and operator authority stop before publication.

## Decision: Normalize archive timestamps and strip build paths without claiming byte-identical SBOMs

**Rationale:** GoReleaser supports source-commit timestamps for build outputs, archives, and metadata. Go builds use `-trimpath`, explicit VCS metadata, and `CGO_ENABLED=0`. Repository verification scans archive members, Go build settings, GoReleaser metadata, and SBOM JSON for checkout paths, drive prefixes, usernames, and other supplied local identifiers. Syft currently emits time- and UUID-sensitive SPDX fields, so deterministic means pinned tools, stable semantics, and source binding rather than byte-for-byte SBOM reproducibility.

**Alternatives considered:** Claiming full cross-host bit-for-bit reproducibility was rejected because S009 is not a reproducible-build certification and upstream deterministic SBOM work remains open. Merely searching archive names was rejected because build metadata and SBOMs can leak paths independently.

## Sources

- [GoReleaser snapshots](https://goreleaser.com/customization/publish/snapshots/)
- [GoReleaser Go builder](https://goreleaser.com/customization/builds/builders/go/)
- [GoReleaser archives](https://goreleaser.com/customization/package/archives/)
- [GoReleaser checksums](https://goreleaser.com/customization/package/checksum/)
- [GoReleaser SBOMs](https://goreleaser.com/customization/sbom/)
- [GoReleaser metadata](https://goreleaser.com/customization/general/metadata/)
- [GoReleaser v2.18.1](https://github.com/goreleaser/goreleaser/releases/tag/v2.18.1)
- [Syft v1.51.1](https://github.com/anchore/syft/releases/tag/v1.51.1)
- [Open Syft deterministic-output work](https://github.com/anchore/syft/pull/3932)
- [Go reproducible builds](https://go.dev/blog/rebuild)
- [GitHub Actions secure use](https://docs.github.com/en/actions/reference/security/secure-use)
- [GitHub Actions workflow permissions](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#permissions)
- [GitHub artifact attestations](https://docs.github.com/en/actions/concepts/security/artifact-attestations)
