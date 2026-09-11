# Release verification

**Status:** Published v0.0.0 verified; v1.0.0 stable release candidate verified without publication

This guide defines the repository-owned artifact proof established by issue [#11](https://github.com/shruggietech/cueson/issues/11), records how accepted default-branch evidence was used to publish and independently verify [v0.0.0](https://github.com/shruggietech/cueson/releases/tag/v0.0.0), and documents the non-publishing stable candidate proof required by issue [#37](https://github.com/shruggietech/cueson/issues/37). The checked-in workflow produces candidate evidence only; it does not create tags, GitHub Releases, release assets, signatures, attestations, or production `cueson.io` state.

## Supported matrix

The snapshot builds with `CGO_ENABLED=0` for:

| Target | Archive | Executable |
|---|---|---|
| `linux/amd64` | `cueson_1.0.0_linux_amd64.tar.gz` | `cueson` |
| `linux/arm64` | `cueson_1.0.0_linux_arm64.tar.gz` | `cueson` |
| `darwin/amd64` | `cueson_1.0.0_darwin_amd64.tar.gz` | `cueson` |
| `darwin/arm64` | `cueson_1.0.0_darwin_arm64.tar.gz` | `cueson` |
| `windows/amd64` | `cueson_1.0.0_windows_amd64.zip` | `cueson.exe` |
| `windows/arm64` | `cueson_1.0.0_windows_arm64.zip` | `cueson.exe` |

Every archive contains exactly the target executable, `cueson.schema.json`, `LICENSE`, and `NOTICE` at its root. Candidate construction packages the immutable [versioned v1.0.0 schema](../schema/releases/v1.0.0/cueson.schema.json), verifies its exact 1.0.0 identity, and compares it byte-for-byte with the canonical repository schema and the schema embedded in every binary. Packaged `LICENSE` and `NOTICE` bytes must match their repository sources, and each legal-file archive member must use the portable approved non-executable mode `0644`. `cueson_1.0.0_checksums.txt` covers exactly the six archives. Each target binary produces one target-bound SPDX JSON SBOM named for the corresponding archive. The immutable [v0.0.0 release schema](../schema/releases/v0.0.0/cueson.schema.json) remains unchanged.

## Exact tools

S009 pins GoReleaser v2.18.1 and Syft v1.51.1. They are build-only commands and are not dependencies of the Cueson product module.

Use a dedicated temporary command directory. In PowerShell:

```powershell
$env:GOBIN = Join-Path $env:TEMP "cueson-release-tools"
$env:GOTOOLCHAIN = "auto"
New-Item -ItemType Directory -Force -Path $env:GOBIN | Out-Null
go install github.com/goreleaser/goreleaser/v2@v2.18.1
go install github.com/anchore/syft/cmd/syft@v1.51.1
$env:Path = "$env:GOBIN;$env:Path"
goreleaser --version
syft version
```

In a POSIX shell:

```bash
export GOBIN="$(mktemp -d)"
export GOTOOLCHAIN="auto"
go install github.com/goreleaser/goreleaser/v2@v2.18.1
go install github.com/anchore/syft/cmd/syft@v1.51.1
export PATH="${GOBIN}:${PATH}"
goreleaser --version
syft version
```

Go's automatic toolchain selection is required while compiling the pinned release tools because their own modules require newer Go compilers than Cueson's Go 1.25 compatibility floor. The Cueson build still follows the repository module's selected compatible toolchain. The GoReleaser version command must report v2.18.1. A Syft binary built through the exact `go install` command can report `[not provided]` because the Go build path does not inject Syft's presentation version; the versioned module command is the reviewed identity in that case. A later stable tool is not an equivalent S009 proof until the versioned configuration and verification are reviewed together.

## Foreground snapshot

Work from a clean checkout with complete Git history, then run:

```text
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

The `.goreleaser.yaml` file fixes the candidate snapshot identity to `1.0.0`, packages the immutable v1 schema, disables release publishing in configuration, closes the build matrix to six targets, normalizes artifact timestamps to the source commit, strips build paths, and injects the `internal/version` override with a verifier-readable marker consumed by the public version surface. GoReleaser writes only beneath ignored `dist/`.

GoReleaser snapshot mode does not upload artifacts. `release.disable: true` is a second boundary so the checked-in configuration cannot publish even if a caller omits snapshot mode. Any later public release requires a separate specification and operator authorization.

## Repository-owned verification

Obtain the full current commit identifier with `git rev-parse HEAD`, then run the standalone standard-library verifier:

```text
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.0.0 -commit <full-commit> -execute-host
```

The verifier:

- requires regular canonical and immutable v1 repository schemas, validates the exact canonical origin, root instance `$schema` constant, and `schema_version`, and proves byte identity before artifact inspection;
- requires the exact six archive names and target tuples;
- accepts only flat regular archive members and rejects absolute paths, traversal, links, devices, duplicates, and unexpected content;
- compares every packaged and emitted schema byte-for-byte with the same immutable and canonical repository bytes;
- requires packaged `LICENSE` and `NOTICE` bytes to match the repository sources, requires archive mode `0644`, and records both digests;
- requires a one-to-one lowercase SHA-256 checksum mapping for the six archives;
- validates target, `CGO_ENABLED=0`, trimmed paths, source revision, and clean VCS state from Go build information, then requires the release-version marker consumed by the public version surface in every target binary;
- validates one binary-derived SPDX JSON SBOM per archive with target and source-revision identity;
- scans Go build information, GoReleaser metadata, and SBOM JSON for structural or supplied local identifiers without treating arbitrary compressed or executable bytes as text;
- executes only the host-compatible packaged binary and requires exact `version`, `schema --version`, and emitted-schema output;
- writes evidence only after every assertion passes, including intended tag `v1.0.0`, the lowercase schema and legal-file SHA-256 values, the exact source revision, and `published: false`; candidate proof runs without development mode.

Accepted candidate evidence records version `1.0.0`, intended tag `v1.0.0`, the exact full source revision, schema, license, and notice digests, archive, SBOM, and checksum counts of six, the ordered six-target archive and SBOM digests, the compatible-host execution result or `null`, and `published: false`. Development mode is absent or false. Existing evidence is not overwritten, so every accepted structural or native proof uses a distinct evidence path or a clean output directory.

Pass additional local values with repeated `-forbid` flags when a machine-specific identifier is not already derived from the repository root, user profile, temporary directory, hostname, or environment.

The verifier structurally inspects all targets but never pretends the current host can execute foreign binaries. Candidate automation builds and structurally accepts one exact bundle on Ubuntu, then fans that same bundle out without rebuilding so hosted Linux, Windows, and macOS runners execute their matching amd64 packaged binaries. Arm64 artifacts receive structural and build-information proof but no execution claim because no governed arm64 runner is available.

## Determinism and provenance boundary

Release binaries and archives use fixed tools, stable names, source-commit timestamps, `-trimpath`, explicit VCS metadata, and no wall-clock linker value. Release evidence omits a verification timestamp so equivalent accepted inputs produce stable semantic evidence.

Syft currently includes variable timestamps and document identifiers in SPDX output. S009 verifies SBOM semantics, target association, and source binding but does not claim byte-for-byte reproducible SBOM documents. S009 also does not claim cryptographic provenance, signature, or GitHub artifact attestation. Those require a separately authorized release design.

## Hosted proof

The `Release proof` workflow runs on ordinary pull requests to `main`, pushes to `main`, and explicit manual dispatch. Its candidate jobs use `contents: read`, persist no checkout credentials, reference no secrets, pass no `GITHUB_TOKEN` to release tools, install the exact Go tool versions, run the same 1.0.0 snapshot and candidate verifier without development mode, and retain only accepted short-lived review evidence. Linux, Windows, and macOS smoke jobs download the already accepted bundle and do not rebuild it.

The retained workflow artifact is CI evidence, not a GitHub Release asset. A failed verifier produces a failed check and uploads no candidate bundle. The workflow has no tag or release trigger, write permission, identity-token permission, signing step, or production deployment step.

A pull-request run proves the reviewed head but does not predict the publication target. After an operator-authorized squash merge, a `main` push run proves the resulting commit independently. Only that accepted default-branch evidence can nominate the exact revision and artifact set that issue #38 may present for separate tag and GitHub Release authority.

The repository ruleset does not require this check. S009 nevertheless required it as operational evidence on its official pull request; changing the protected required-check set remains separately governed repository-control work.

## Delivered evidence

The `Non-publishing snapshot` job passed on the final head of pull request [#21](https://github.com/shruggietech/cueson/pull/21), together with all CI, CodeQL, and pull-request-policy gates. The operator then merged the pull request into `main` as `3da0a4b4eeae57024d837c5f46e9d62537ffab95` on 2026-09-10, closing issue [#11](https://github.com/shruggietech/cueson/issues/11). Post-merge [CI run 34421328445](https://github.com/shruggietech/cueson/actions/runs/34421328445) and [CodeQL run 34421328401](https://github.com/shruggietech/cueson/actions/runs/34421328401) also completed successfully.

S012 extended the workflow so its release-preparation squash-merge commit received a distinct proof rather than inheriting the pull-request head's evidence. For exact revision `b294a6952c8bd041d852c502f5d7206c0b58edd6`, [CI run 34543376819](https://github.com/shruggietech/cueson/actions/runs/34543376819), [CodeQL run 34543376786](https://github.com/shruggietech/cueson/actions/runs/34543376786), and [Release proof run 34543376814](https://github.com/shruggietech/cueson/actions/runs/34543376814) all completed successfully. Artifact `10178231596` recorded schema SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`, six archive digests, six SBOM digests, six checksum entries, Linux amd64 host execution, and `published: false` as the truthful state of the non-publishing verifier.

S013 created unsigned annotated tag [`v0.0.0`](https://github.com/shruggietech/cueson/tree/v0.0.0), proved its remote peeled target was the same exact revision, and published the final [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0) at `2026-09-11T00:20:19Z`. The public set contains exactly the accepted six archives, six matching SPDX JSON SBOMs, and `cueson_0.0.0_checksums.txt`; internal `release-evidence.json`, GoReleaser metadata and configuration, and extracted directories were not published.

Independent post-publication verification downloaded all thirteen assets into a new clean directory, matched every SHA-256 value against the accepted evidence, applied the exact six-entry checksum bijection, executed the public Windows amd64 binary through a hidden non-interactive process, required exact `0.0.0` output from `version` and `schema --version`, and compared emitted schema bytes with the tagged immutable schema. GitHub API read-back also proved the release body matched `docs/releases/v0.0.0.md`, every asset was uploaded, and the release was public, final, and non-prerelease. Milestone v0.0.0 remained open at the S013 boundary and was closed only through later authorized reconciliation; no signature or attestation asset was added, and the production schema hostname remained inactive.

## Complete local gate

Before publication or review response, run:

```powershell
go run ./scripts/github-format/main.go .
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/github-format test -count=1 ./...
go -C scripts/docs-verify test -count=1 ./...
go -C scripts/docs-verify run . -repo ../..
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/release-verify test -count=1 ./...
go -C scripts/brand-verify test -count=1 ./...
go -C scripts/brand-verify run . -repo ../..
goreleaser check
goreleaser release --snapshot --clean --skip=publish
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.0.0 -commit <full-commit> -execute-host
git diff --check
```

Success proves a non-publishing v1.0.0 stable release candidate at the current revision. It does not merge the candidate pull request, create or move a tag, publish a GitHub Release or asset, serve the schema, close the milestone, add a signature or attestation, or mutate production state.
