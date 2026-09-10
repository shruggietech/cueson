# Data Model: Non-Publishing Release Proof

## Release Candidate

Represents one complete non-published v0.0.0 build from one source revision.

| Field | Rule |
|---|---|
| `version` | Exactly `0.0.0` for S009 |
| `source_revision` | Full lowercase Git commit identifier for the checkout |
| `snapshot` | Always true |
| `dirty` | Must be false for final evidence |
| `targets` | Exactly the six approved targets |
| `archives` | Exactly one archive per target |
| `checksum_manifest` | Exactly one SHA-256 manifest |
| `sboms` | Exactly one SPDX JSON document per archive |
| `published` | Always false |

## Target

Represents one operating-system and architecture combination.

| Field | Rule |
|---|---|
| `goos` | `windows`, `darwin`, or `linux` |
| `goarch` | `amd64` or `arm64` |
| `binary_name` | `cueson.exe` for Windows, otherwise `cueson` |
| `archive_format` | ZIP for Windows, tar.gz otherwise |
| `archive_name` | `cueson_0.0.0_<goos>_<goarch>.<format>` |

Target identity is unique by the ordered pair `(goos, goarch)`. All three operating systems have both architectures.

## Distributable Archive

Represents one target package.

| Field | Rule |
|---|---|
| `target` | Resolves to exactly one Target |
| `path` | Safe filename directly beneath `dist/` |
| `members` | Exactly binary, `cueson.schema.json`, `LICENSE`, and `NOTICE` |
| `sha256` | Lowercase 64-character digest of exact archive bytes |
| `schema_sha256` | Equal to the canonical repository schema digest |
| `source_revision` | Equal to the Release Candidate source revision through Go build information |

Unsafe, absolute, traversing, duplicate, directory-only, or unexpected member names invalidate the archive.

## Checksum Entry

Represents one line in `cueson_0.0.0_checksums.txt`.

| Field | Rule |
|---|---|
| `sha256` | Lowercase 64-character digest |
| `filename` | Exact safe distributable archive filename |

The manifest contains exactly six unique entries. It excludes itself, SBOMs, metadata, unpacked binaries, and internal build directories.

## Software Bill of Materials

Represents the SPDX JSON catalog generated from one target binary and associated with its corresponding archive.

| Field | Rule |
|---|---|
| `filename` | Exact archive filename plus `.sbom.json` |
| `spdx_version` | A supported SPDX JSON version identifier |
| `document_namespace` | Present and non-empty |
| `archive` | Resolves to exactly one Distributable Archive whose target binary was cataloged |
| `source_revision` | Encoded in source identity and matched to Release Candidate metadata |
| `local_path_hits` | Empty |

## GoReleaser Metadata

Represents `metadata.json` and `artifacts.json` emitted by the snapshot orchestrator. Required evidence includes project name `cueson`, version `0.0.0`, full source revision, snapshot state, six target archives, six archive SBOMs, and one checksum manifest. Any absolute workspace path or supplied machine identifier invalidates the proof.

## Release Evidence

Represents the verifier's deterministic JSON summary written only after all checks pass.

| Field | Rule |
|---|---|
| `version` | `0.0.0` |
| `source_revision` | Full expected commit |
| `archive_count` | 6 |
| `sbom_count` | 6 |
| `checksum_count` | 6 |
| `host_executed` | Target identifier when compatible execution ran, otherwise an explicit truthful null/result |
| `published` | false |
| `verified_at` | Omitted to keep evidence deterministic |

## State Transitions

```text
source checkout
  -> tools pinned
  -> snapshot built
  -> archives inspected
  -> checksums verified
  -> schema and version lockstep verified
  -> SBOM and metadata verified
  -> host-compatible CLI executed
  -> release evidence accepted
```

Any failed transition terminates the proof. There is no transition from S009 evidence to tag, signed artifact, attestation, GitHub Release, or production schema publication.
