# Release verification

**Status:** Exact v1.1.0 release publicly published and independently verified; non-publishing proof retained

This guide defines the repository-owned artifact proof established by issue [#11](https://github.com/shruggietech/cueson/issues/11), records how accepted default-branch evidence was used to publish and independently verify [v0.0.0](https://github.com/shruggietech/cueson/releases/tag/v0.0.0), [v1.0.0](https://github.com/shruggietech/cueson/releases/tag/v1.0.0) and [v1.1.0](https://github.com/shruggietech/cueson/releases/tag/v1.1.0), and documents the retained non-publishing candidate proof. The checked-in workflow produces candidate evidence only; it does not create tags, GitHub Releases, release assets, signatures, attestations, or production `cueson.io` state.

## Current stable release proof

The checked-in workflow continues to build non-publishing exact `1.1.0` snapshots and uses stable verification of canonical/immutable/embedded/emitted/packaged schema byte identity. It retains exact head revision, schema/software identity, six archives, six target-bound SBOMs, checksum bijection and legal-file proof. Matching native Linux, Windows and macOS amd64 packages execute four-format workflows and historical/safety checks; verified published v1.0.0 bytes prove rejection of fresh current output. This workflow remains evidence-only even though the exact v1.1.0 release has now been published and independently verified.

After building a snapshot, the active local verifier command is:

```text
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.1.0 -commit <full-commit> -execute-host
```

For exact 1.1.0 with `-execute-host`, optional `native_proof` records the four tested formats, two historical input documents, and the verified old-consumer version/revision/archive/digest plus refusal count. The old consumer is the published 1.0.0 binary at `2cad4c816340404289b4d1d87179a4071713bb46`, authenticated against the frozen S020 public asset contract before execution. This verifier-only fetch requires network availability and fails closed on asset/hash/source mismatch; product schema/codec execution remains offline.

## Published v1.1.0 matrix

The snapshot builds with `CGO_ENABLED=0` for:

| Target | Archive | Executable |
|---|---|---|
| `linux/amd64` | `cueson_1.1.0_linux_amd64.tar.gz` | `cueson` |
| `linux/arm64` | `cueson_1.1.0_linux_arm64.tar.gz` | `cueson` |
| `darwin/amd64` | `cueson_1.1.0_darwin_amd64.tar.gz` | `cueson` |
| `darwin/arm64` | `cueson_1.1.0_darwin_arm64.tar.gz` | `cueson` |
| `windows/amd64` | `cueson_1.1.0_windows_amd64.zip` | `cueson.exe` |
| `windows/arm64` | `cueson_1.1.0_windows_arm64.zip` | `cueson.exe` |

Every archive contains exactly the target executable, `cueson.schema.json`, `LICENSE` and `NOTICE` at its root. Published v1.1.0 packages use `schema/releases/v1.1.0/cueson.schema.json` and stable verification against canonical/embedded/emitted/packaged bytes, with exact `1.1.0` identity. Historical published v1.0.0 schema and evidence bytes remain unchanged.

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

The `.goreleaser.yaml` file fixes the candidate snapshot identity to `1.1.0`, packages the immutable candidate schema, disables publishing in configuration, closes the build matrix to six pure-Go targets, normalizes artifact timestamps to the source commit, strips build paths, and injects the verifier-readable release-version marker. GoReleaser writes only beneath ignored `dist/`. Snapshot status describes the non-publishing packaging operation, not a development schema identity.

GoReleaser snapshot mode does not upload artifacts. `release.disable: true` is a second boundary so the checked-in configuration cannot publish even if a caller omits snapshot mode. Any later public release requires a separate specification and operator authorization.

## Repository-owned verification

Obtain the full current commit identifier with `git rev-parse HEAD`, then run the standalone standard-library verifier:

```text
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.1.0 -commit <full-commit> -execute-host
```

The verifier:

- requires regular canonical/immutable repository schemas with exact current URI, root instance `$schema` constant and `schema_version`, and validates their byte identity;
- requires the exact six archive names and target tuples;
- accepts only flat regular archive members and rejects absolute paths, traversal, links, devices, duplicates, and unexpected content;
- compares every packaged and emitted schema byte-for-byte with the canonical and immutable repository candidate bytes;
- requires packaged `LICENSE` and `NOTICE` bytes to match the repository sources, requires archive mode `0644`, and records both digests;
- requires a one-to-one lowercase SHA-256 checksum mapping for the six archives;
- validates target, `CGO_ENABLED=0`, trimmed paths, source revision, and clean VCS state from Go build information, then requires the release-version marker consumed by the public version surface in every target binary;
- validates one binary-derived SPDX JSON SBOM per archive with target and source-revision identity;
- scans Go build information, GoReleaser metadata, and SBOM JSON for structural or supplied local identifiers without treating arbitrary compressed or executable bytes as text;
- executes only the host-compatible packaged binary, requiring exact version/schema discovery, four-format native validation/inspection/render/restoration/conversion, historical acceptance and output-refusal safety;
- writes evidence only after every assertion passes, including intended tag `v1.1.0`, the lowercase schema and legal-file SHA-256 values, the exact source revision, `development` absent or false and `published: false`.

Accepted current candidate evidence records exact version `1.1.0`, intended tag `v1.1.0`, full source revision, schema/license/notice digests, archive/SBOM/checksum counts of six, ordered target identities and archive/SBOM digests, compatible-host workflow proof or `null`, and `published: false`. The intended tag is prospective evidence, not authorization or tag creation. Distinct evidence paths retain each structural/native proof; no historical accepted evidence is overwritten.

Pass additional local values with repeated `-forbid` flags when a machine-specific identifier is not already derived from the repository root, user profile, temporary directory, hostname, or environment.

The verifier structurally inspects all targets but never pretends the current host can execute foreign binaries. Candidate automation builds and structurally accepts one exact bundle on Ubuntu, then fans that same bundle out without rebuilding so hosted Linux, Windows, and macOS runners execute their matching amd64 packaged binaries. Arm64 artifacts receive structural and build-information proof but no execution claim because no governed arm64 runner is available.

## Determinism and provenance boundary

Release binaries and archives use fixed tools, stable names, source-commit timestamps, `-trimpath`, explicit VCS metadata, and no wall-clock linker value. Release evidence omits a verification timestamp so equivalent accepted inputs produce stable semantic evidence.

Syft currently includes variable timestamps and document identifiers in SPDX output. S009 verifies SBOM semantics, target association, and source binding but does not claim byte-for-byte reproducible SBOM documents. S009 also does not claim cryptographic provenance, signature, or GitHub artifact attestation. Those require a separately authorized release design.

## Hosted proof

The `Release proof` workflow runs on ordinary pull requests to `main`, pushes to `main`, and manual dispatch. Candidate jobs use read-only repository permission, no credentials/secrets/publication token and exact pinned tools. They build and stably verify the exact `1.1.0` six-target snapshot. Matching Linux, Windows and macOS amd64 smoke jobs download and execute that same bundle without rebuilding. Published old-consumer proof uses a source/digest-bound v1.0.0 executable and verifies current-output identity refusal before publication. Accepted short-lived artifacts are review evidence, never GitHub Release assets.

The retained workflow artifact is CI evidence, not a GitHub Release asset. A failed verifier produces a failed check and uploads no candidate bundle. The workflow has no tag or release trigger, write permission, identity-token permission, signing step, or production deployment step.

A pull-request run proves its exact reviewed head, not the final publication target. S029 received a fresh successful `main` push proof binding revision `7ff45c1d8cd8df377e1fb568b9785286b649fd7c` and artifact digests before #66 received exact tag/release authority. S030 child #76 owns the reviewed schema/site artifact, and #67 governs separately authorized exact-main production deployment and live proof. The historical default-branch evidence below records each independently verified published artifact set.

The repository ruleset does not require this check. S009 nevertheless required it as operational evidence on its official pull request; changing the protected required-check set remains separately governed repository-control work.

## Publication decision binding

The dated 1.1.0 history was prepared metadata before publication. The reviewed preparation pull request merged before final tag selection; fresh accepted main proof of the actual squash revision supplied the exact publication source and thirteen-asset inventory. Earlier S028 main evidence and preparation PR proof remain historical/review evidence and did not substitute for this source binding. The decision also bound the exact separately frozen formatted public-note digest.

The complete source/asset/schema/legal/native/notes package was presented before requesting exact tag creation/push and release/asset publication authority. Independently downloaded public bytes and GitHub read-back established the authorized result and completed #66. Public schema/site activation and live read-back are governed separately through #67.

## Delivered evidence

The `Non-publishing snapshot` job passed on the final head of pull request [#21](https://github.com/shruggietech/cueson/pull/21), together with all CI, CodeQL, and pull-request-policy gates. The operator then merged the pull request into `main` as `3da0a4b4eeae57024d837c5f46e9d62537ffab95` on 2026-09-10, closing issue [#11](https://github.com/shruggietech/cueson/issues/11). Post-merge [CI run 34421328445](https://github.com/shruggietech/cueson/actions/runs/34421328445) and [CodeQL run 34421328401](https://github.com/shruggietech/cueson/actions/runs/34421328401) also completed successfully.

S012 extended the workflow so its release-preparation squash-merge commit received a distinct proof rather than inheriting the pull-request head's evidence. For exact revision `b294a6952c8bd041d852c502f5d7206c0b58edd6`, [CI run 34543376819](https://github.com/shruggietech/cueson/actions/runs/34543376819), [CodeQL run 34543376786](https://github.com/shruggietech/cueson/actions/runs/34543376786), and [Release proof run 34543376814](https://github.com/shruggietech/cueson/actions/runs/34543376814) all completed successfully. Artifact `10178231596` recorded schema SHA-256 `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`, six archive digests, six SBOM digests, six checksum entries, Linux amd64 host execution, and `published: false` as the truthful state of the non-publishing verifier.

S013 created unsigned annotated tag [`v0.0.0`](https://github.com/shruggietech/cueson/tree/v0.0.0), proved its remote peeled target was the same exact revision, and published the final [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0) at `2026-09-11T00:20:19Z`. The public set contains exactly the accepted six archives, six matching SPDX JSON SBOMs, and `cueson_0.0.0_checksums.txt`; internal `release-evidence.json`, GoReleaser metadata and configuration, and extracted directories were not published.

Independent post-publication verification downloaded all thirteen assets into a new clean directory, matched every SHA-256 value against the accepted evidence, applied the exact six-entry checksum bijection, executed the public Windows amd64 binary through a hidden non-interactive process, required exact `0.0.0` output from `version` and `schema --version`, and compared emitted schema bytes with the tagged immutable schema. GitHub API read-back also proved the release body matched `docs/releases/v0.0.0.md`, every asset was uploaded, and the release was public, final, and non-prerelease. Milestone v0.0.0 remained open at the S013 boundary and was closed only through later authorized reconciliation; no signature or attestation asset was added, and the production schema hostname remained inactive.

S019 bound the stable v1.0.0 candidate to exact post-squash revision `2cad4c816340404289b4d1d87179a4071713bb46`. [CI run 34621429606](https://github.com/shruggietech/cueson/actions/runs/34621429606), [CodeQL run 34621429656](https://github.com/shruggietech/cueson/actions/runs/34621429656), and [Release proof run 34621429626](https://github.com/shruggietech/cueson/actions/runs/34621429626) completed successfully. Accepted artifact `10273380044` recorded schema SHA-256 `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541`, exact legal-file bytes and modes, six archive digests, six SBOM digests, six checksum entries, native Linux, Windows, and macOS amd64 execution, and `published: false` as the truthful state of the non-publishing verifier.

S020 created unsigned annotated tag [`v1.0.0`](https://github.com/shruggietech/cueson/tree/v1.0.0), proved its remote peeled target was the exact accepted revision, and published the final [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v1.0.0) at `2026-09-11T16:58:30Z`. The public set contains exactly the accepted six archives, six matching SPDX JSON SBOMs, and `cueson_1.0.0_checksums.txt`; internal evidence, GoReleaser metadata and configuration, and extracted directories were not published.

Independent v1.0.0 verification downloaded all thirteen public assets into a new clean directory, matched every filename, byte length, and SHA-256 value against `specs/S020-publish-v1-release/contracts/release-publication-contract.json`, applied the exact six-entry checksum bijection, structurally inspected all six archives, validated all six SBOMs, and executed the public Windows amd64 binary through the repository's hidden non-interactive verifier. Exact `1.0.0` output from `version` and `schema --version`, tagged schema byte identity, build target and revision, pure-Go state, legal-file bytes and modes, release marker, and local-identifier exclusions all passed. GitHub read-back proved the release body exactly matched `docs/releases/v1.0.0.md`, every asset was uploaded, and the release was public, final, and non-prerelease. The v1 epic and milestone remain open for separately authorized reconciliation; no signature, attestation, public schema endpoint, production-domain change, or pull-request merge occurred during publication.

S029 bound exact post-squash revision `7ff45c1d8cd8df377e1fb568b9785286b649fd7c`. [CI run 35040576852](https://github.com/shruggietech/cueson/actions/runs/35040576852), [CodeQL run 35040576902](https://github.com/shruggietech/cueson/actions/runs/35040576902), [Site run 35040576878](https://github.com/shruggietech/cueson/actions/runs/35040576878), and [Release proof run 35040576916](https://github.com/shruggietech/cueson/actions/runs/35040576916) completed successfully. Accepted artifact `10424912955` recorded the exact 1.1.0 schema/legal bytes, six archive digests, six SBOM digests, six checksum entries, and matching native Linux, Windows and macOS amd64 execution. Its `published: false` field truthfully describes the non-publishing verifier run used to make the later decision.

Issue #66 used explicit operator authority to create and push unsigned annotated tag [`v1.1.0`](https://github.com/shruggietech/cueson/tree/v1.1.0), prove its peeled target was the accepted exact revision, and publish the final [v1.1.0 GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v1.1.0) at `2026-09-16T22:47:59Z`. The public set contains exactly six archives, six matching SPDX JSON SBOMs, and `cueson_1.1.0_checksums.txt`; internal evidence, GoReleaser metadata and configuration, and extracted directories were not published.

Independent v1.1.0 verification downloaded all thirteen public assets into a new clean directory, matched every filename, byte length and SHA-256 value against the accepted exact-main decision package, applied the exact six-entry checksum bijection, structurally inspected all six archives, validated all six SBOMs, and executed the public Windows amd64 package. Exact version/schema discovery, tagged immutable schema identity, build target and revision, pure-Go state, legal-file bytes and modes, release marker, local-identifier exclusions and the four-format native/historical/safety proof all passed. GitHub read-back proved the release body exactly matched the reviewed formatted notes, every asset was uploaded, and the release was public, final and non-prerelease. That publication did not change production `cueson.io`; #76 owns the reviewed schema/site artifact and #67 governs its production activation and live verification.

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
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 1.1.0 -commit <full-commit> -execute-host
git diff --check
```

Success proves the non-publishing exact 1.1.0 stable snapshot at the current clean revision, including immutable-copy proof. The same checks continue to produce review evidence after publication; they do not merge, create or move a tag, publish or replace a release asset, serve the new schema route, close the milestone, add signatures/attestations or mutate production.
