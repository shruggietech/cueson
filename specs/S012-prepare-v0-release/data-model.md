# Data Model: Prepare the v0.0.0 Release

S012 changes no application data. These entities define the repository release record and its deterministic evidence.

## Release Schema Candidate

Fields:

- `version`: semantic version `0.0.0` declared by schema identity and `schema_version` constant.
- `canonical_path`: `internal/schema/cueson.schema.json`.
- `release_path`: `schema/releases/v0.0.0/cueson.schema.json`.
- `sha256`: lowercase 64-character digest of the exact shared bytes.

Validation rules:

- Both paths are regular repository files and valid UTF-8 JSON without a BOM.
- Both files are byte-identical after checkout.
- `$id` ends in `/v0.0.0/cueson.schema.json` and `schema_version.const` equals `0.0.0`.
- The packaged and embedded schema bytes equal the same content.

State transitions:

- Before S012 merge, the release path is a reviewed candidate.
- After the authorized release is published from the tagged commit, the versioned release path is immutable forever.

## Candidate Revision

Fields:

- `source_revision`: one lowercase 40-character full Git SHA-1.
- `version`: `0.0.0`.
- `clean`: build metadata value `true` represented by `vcs.modified=false`.
- `branch_role`: pull-request review evidence or proposed default-branch publication target.

Validation rules:

- GoReleaser metadata, every binary build record, every SBOM package version, and release evidence bind to the same revision.
- A pull-request head remains review evidence only.
- The proposed publication target must be the exact S012 squash-merge commit on `main` with successful default-branch proof.

## Target Evidence

Fields:

- `goos`: one of `linux`, `darwin`, or `windows`.
- `goarch`: one of `amd64` or `arm64`.
- `archive`: exact versioned archive filename.
- `binary`: `cueson` or `cueson.exe`.
- `sbom`: exact target-bound SPDX JSON filename.
- `archive_sha256`: lowercase archive digest from the checksum manifest.
- `sbom_sha256`: lowercase digest of the inspected SBOM document.

Validation rules:

- Exactly six unique operating-system and architecture pairs exist.
- Every archive contains exactly the target binary, `cueson.schema.json`, `LICENSE`, and `NOTICE`.
- Every target reports `CGO_ENABLED=0`, trimmed paths, the expected revision, and unmodified VCS state.
- SBOM semantics identify the target and version plus source revision without claiming stable upstream timestamps or namespace bytes.

## Release Evidence

Fields:

- `version`: `0.0.0`.
- `source_revision`: Candidate Revision identity.
- `release_schema_sha256`: Release Schema Candidate digest.
- `archive_count`: exactly `6`.
- `sbom_count`: exactly `6`.
- `checksum_count`: exactly `6`.
- `host_executed`: compatible target identity or `null` when execution was not requested.
- `targets`: ordered collection of six Target Evidence records.
- `published`: always `false` for this verifier.

Validation rules:

- Evidence is written only after every repository, archive, executable, checksum, and SBOM assertion passes.
- Evidence is deterministic for stable semantic inputs except for artifact digests that legitimately reflect upstream SBOM byte variability.
- Existing evidence is never overwritten; a clean snapshot directory is required.

## Release Notes

Fields:

- `title`: Cueson v0.0.0.
- `highlights`: concise shipped foundation outcomes.
- `limitations`: truthful absence of native SRT and WebVTT ingest, render, and conversion.
- `changelog_link`: exact final tagged URL.

Validation rules:

- The document is UTF-8 Markdown with one source line per paragraph or list item.
- The final non-empty line is exactly the required changelog sentence.
- The notes do not claim publication before the protected release action occurs.

## Publication Decision Package

Fields:

- `candidate_revision`: exact post-merge default-branch commit.
- `release_record`: changelog, release notes, and versioned schema digest.
- `artifact_record`: archives, checksums, SBOMs, and verifier evidence.
- `quality_record`: green default-branch checks and resolved findings.
- `limitations`: envelope-only format capability and omitted signatures, attestations, schema hosting, and production changes.

State transitions:

- `prepared`: S012 pull request is green and fully reviewed.
- `candidate`: S012 is merged and exact `main` proof is green.
- `authorized`: the operator grants specific tag and release authority for the named candidate and artifacts.
- `published`: protected actions complete and read-back matches the authorization.
- `verified`: post-publication comparisons pass; milestone closure can then be separately decided.
