# Quickstart: Verify the Non-Publishing Release Proof

## Prerequisites

- Work from a clean `S009-prove-release-pipeline` checkout with complete Git history.
- Use Go 1.25.0 or the latest compatible Go 1.25 patch.
- Install GoReleaser v2.18.1 and Syft v1.51.1 from their exact Go module versions.
- Keep commands foreground and non-interactive.

## 1. Validate repository artifacts

```powershell
go run ./scripts/github-format/main.go .
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go -C scripts/github-format test -count=1 ./...
go -C scripts/pr-policy test -count=1 ./...
go -C scripts/release-verify test -count=1 ./...
git diff --check
```

Expected: every command succeeds and repository-authored files satisfy encoding and line-ending policy.

## 2. Install exact release tools

Use a dedicated temporary `GOBIN`, then install:

```text
go install github.com/goreleaser/goreleaser/v2@v2.18.1
go install github.com/anchore/syft/cmd/syft@v1.51.1
```

Expected: `goreleaser --version` reports v2.18.1 and `syft version` reports v1.51.1. The tools remain outside the product module and release archives.

## 3. Validate configuration and build the snapshot

```text
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

Expected: GoReleaser produces six archives, six SPDX JSON SBOMs, one checksum manifest, and metadata beneath ignored `dist/`; it does not upload or publish anything.

## 4. Verify the complete release candidate

```text
go -C scripts/release-verify run . -dist ../../dist -repo ../.. -version 0.0.0 -commit <full-HEAD-sha> -execute-host
```

Expected: the verifier opens every archive, validates its four members, verifies all six checksums and SBOMs, compares schema bytes and version identities, inspects build and release metadata, executes the compatible packaged binary when available, finds no supplied local identifiers, and writes `dist/release-evidence.json` with `published: false`.

## 5. Inspect non-publication evidence

Verify that no `v0.0.0` tag or GitHub Release was created, no release asset was uploaded, and no production domain configuration changed. Inspect the workflow definition to confirm `contents: read`, no secret references, no signing or attestation permission, no tag or release trigger, `release.disable: true`, and the mandatory `--snapshot` argument.

## 6. Validate hosted behavior

After the authorized branch push, wait for `Release proof / Non-publishing snapshot` on the pull-request head. Download the retained workflow artifact only after the check succeeds and compare its file inventory with [contracts/artifact-contract.md](contracts/artifact-contract.md).

Expected: hosted automation applies the same verifier and produces no publication side effect. S009 then proceeds through at most two Codex rounds and stops for the operator's final merge ritual.
