# Release verification

**Status:** Non-publishing v0.0.0 release-candidate proof

This guide defines the repository-owned release proof completed by issue [#11](https://github.com/shruggietech/cueson/issues/11). It produces and verifies candidate artifacts only. It does not create a tag, GitHub Release, release asset, signature, attestation, or production `cueson.io` state.

## Supported matrix

The snapshot builds with `CGO_ENABLED=0` for:

| Target | Archive | Executable |
|---|---|---|
| `linux/amd64` | `cueson_0.0.0_linux_amd64.tar.gz` | `cueson` |
| `linux/arm64` | `cueson_0.0.0_linux_arm64.tar.gz` | `cueson` |
| `darwin/amd64` | `cueson_0.0.0_darwin_amd64.tar.gz` | `cueson` |
| `darwin/arm64` | `cueson_0.0.0_darwin_arm64.tar.gz` | `cueson` |
| `windows/amd64` | `cueson_0.0.0_windows_amd64.zip` | `cueson.exe` |
| `windows/arm64` | `cueson_0.0.0_windows_arm64.zip` | `cueson.exe` |

Every archive contains exactly the target executable, `cueson.schema.json`, `LICENSE`, and `NOTICE` at its root. `cueson_0.0.0_checksums.txt` covers exactly those six archives. Each target binary produces one target-bound SPDX JSON SBOM named for the corresponding archive, which avoids introducing temporary archive-extraction paths into the document.

## Exact tools

S009 pins GoReleaser v2.18.1 and Syft v1.51.1. They are build-only commands and are not dependencies of the Cueson product module.

Use a dedicated temporary command directory. In PowerShell:

```powershell
$env:GOBIN = Join-Path $env:TEMP "cueson-release-tools"
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
go install github.com/goreleaser/goreleaser/v2@v2.18.1
go install github.com/anchore/syft/cmd/syft@v1.51.1
export PATH="${GOBIN}:${PATH}"
goreleaser --version
syft version
```

The GoReleaser version command must report v2.18.1. A Syft binary built through the exact `go install` command can report `[not provided]` because the Go build path does not inject Syft's presentation version; the versioned module command is the reviewed identity in that case. A later stable tool is not an equivalent S009 proof until the versioned configuration and verification are reviewed together.

## Foreground snapshot

Work from a clean checkout with complete Git history, then run:

```text
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

The `.goreleaser.yaml` file fixes the snapshot identity to `0.0.0`, disables release publishing in configuration, closes the build matrix to six targets, normalizes artifact timestamps to the source commit, strips build paths, and injects the existing `internal/version.current` variable. GoReleaser writes only beneath ignored `dist/`.

GoReleaser snapshot mode does not upload artifacts. `release.disable: true` is a second boundary so the checked-in configuration cannot publish even if a caller omits snapshot mode. A future public release requires a separate specification and operator authorization.

## Repository-owned verification

Obtain the full current commit identifier with `git rev-parse HEAD`, then run the standalone standard-library verifier:

```text
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.0.0 -commit <full-commit> -execute-host
```

The verifier:

- requires the exact six archive names and target tuples;
- accepts only flat regular archive members and rejects absolute paths, traversal, links, devices, duplicates, and unexpected content;
- compares every packaged schema byte-for-byte with `internal/schema/cueson.schema.json`;
- requires a one-to-one lowercase SHA-256 checksum mapping for the six archives;
- validates target, `CGO_ENABLED=0`, trimmed paths, source revision, clean VCS state, and injected version from Go build information;
- validates one binary-derived SPDX JSON SBOM per archive with target and source-revision identity;
- scans names, members, binaries, GoReleaser metadata, and SBOM JSON for structural or supplied local identifiers;
- executes only the host-compatible packaged binary and requires exact `version`, `schema --version`, and emitted-schema output;
- writes `dist/release-evidence.json` only after every assertion passes.

Pass additional local values with repeated `-forbid` flags when a machine-specific identifier is not already derived from the repository root, user profile, temporary directory, hostname, or environment.

The verifier structurally inspects all targets but never pretends the current host can execute foreign binaries. Hosted Ubuntu exercises Linux amd64; a native Windows foreground run exercises Windows amd64. Native macOS execution remains covered by its existing product tests until a later release workflow deliberately adds a macOS snapshot job.

## Determinism and provenance boundary

Release binaries and archives use fixed tools, stable names, source-commit timestamps, `-trimpath`, explicit VCS metadata, and no wall-clock linker value. Release evidence omits a verification timestamp so equivalent accepted inputs produce stable semantic evidence.

Syft currently includes variable timestamps and document identifiers in SPDX output. S009 verifies SBOM semantics, target association, and source binding but does not claim byte-for-byte reproducible SBOM documents. S009 also does not claim cryptographic provenance, signature, or GitHub artifact attestation. Those require a separately authorized release design.

## Hosted proof

The `Release proof` workflow runs on ordinary pull requests to `main` and explicit manual dispatch. Its `Non-publishing snapshot` job uses `contents: read`, persists no checkout credentials, references no secrets, passes no `GITHUB_TOKEN` to release tools, installs the exact Go tool versions, runs the same snapshot and verifier, and uploads `dist/` only after acceptance for short-lived review.

The retained workflow artifact is CI evidence, not a GitHub Release asset. A failed verifier produces a failed check and uploads no candidate bundle. The workflow has no tag or release trigger, write permission, identity-token permission, signing step, or production deployment step.

The repository ruleset does not automatically require this new check. S009 treats a green release-proof check on its official pull request as mandatory operational evidence; changing the protected required-check set remains separately governed repository-control work.

## Complete local gate

Before publication or review response, run:

```powershell
go run ./scripts/github-format/main.go .
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/github-format test -count=1 ./...
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/release-verify test -count=1 ./...
goreleaser check
goreleaser release --snapshot --clean --skip=publish
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.0.0 -commit <full-commit> -execute-host
git diff --check
```

Success proves a non-published candidate. It does not authorize the final merge, a tag, release publication, schema publication, milestone closure, or production mutation.
