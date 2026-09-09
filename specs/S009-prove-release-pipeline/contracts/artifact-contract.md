# Release Artifact Contract

## Matrix and names

| Target | Archive | Executable |
|---|---|---|
| `linux/amd64` | `cueson_0.0.0_linux_amd64.tar.gz` | `cueson` |
| `linux/arm64` | `cueson_0.0.0_linux_arm64.tar.gz` | `cueson` |
| `darwin/amd64` | `cueson_0.0.0_darwin_amd64.tar.gz` | `cueson` |
| `darwin/arm64` | `cueson_0.0.0_darwin_arm64.tar.gz` | `cueson` |
| `windows/amd64` | `cueson_0.0.0_windows_amd64.zip` | `cueson.exe` |
| `windows/arm64` | `cueson_0.0.0_windows_arm64.zip` | `cueson.exe` |

Every archive contains exactly four root members:

```text
cueson or cueson.exe
cueson.schema.json
LICENSE
NOTICE
```

Directories, absolute paths, traversal, duplicate canonical names, links, device entries, and additional members are rejected.

## Checksum manifest

`cueson_0.0.0_checksums.txt` contains one lowercase SHA-256 entry for each archive in deterministic filename order. It contains six entries and no entries for SBOMs, metadata, evidence, unpacked binaries, or itself.

## Software bills of materials

Each target binary is the Syft catalog input for one sibling SPDX JSON document named `<archive>.sbom.json` for its corresponding archive. Every document parses as a JSON object, declares an SPDX version and document namespace, and identifies its target plus `0.0.0+<full-commit>` source version. Temporary archive-extraction paths are therefore outside the generation path and remain forbidden in accepted output.

S009 does not sign SBOMs or publish an attestation. Absence of those claims is part of acceptance.

## Metadata and evidence

GoReleaser emits `metadata.json` and `artifacts.json`. Repository verification emits `release-evidence.json` only after all assertions pass. The evidence summary contains the version, full source revision, expected counts, target identities, archive and SBOM digests, host-execution result, and `published: false`.

No inspected JSON, archive name, archive member, Go build setting, or executable string may contain supplied checkout-root, username, hostname, drive-root, or other local-identifier needles.

## Version and source identity

The following values must all equal `0.0.0`:

- requested snapshot version;
- packaged executable `cueson version` output on a compatible host;
- packaged executable `cueson schema --version` output on a compatible host;
- release override marker consumed by `internal/version.String`;
- canonical packaged schema `schema_version`;
- canonical packaged schema version segment in `$id`.

Every binary's Go build information must identify its expected target and the full source revision when the Go toolchain emits VCS settings. The revision must equal the verified checkout revision and `vcs.modified` must be false.

Every target binary must contain `cueson-release-version:0.0.0`, the release override marker consumed by the public version surface. This makes the injected version independently verifiable for foreign targets that the current host cannot execute.

## Failure behavior

The verifier exits nonzero with a stable diagnostic for every missing, duplicate, unknown, unsafe, mismatched, stale, or path-leaking artifact. It does not repair output or delete evidence needed to diagnose the failed run.
