# Data Model: Prepare the v1.0.0 Release Candidate

## Stable Release Identity

Fields:

- `version`: exact semantic version `1.0.0`.
- `schema_id`: exact `https://cueson.io/schema/v1.0.0/cueson.schema.json` identifier.
- `intended_tag`: exact `v1.0.0` name; intent only until separately authorized.
- `format_status`: `stable` for the frozen SubRip and WebVTT contracts.

Validation rules:

- The executable, schema constants, schema bytes, emitted documents, candidate configuration, artifacts, evidence, changelog, and notes agree on the version.
- The identifier is an identity before production hosting and does not claim that the URL is live.
- Previously released versioned schema bytes never change.

## Immutable v1 Schema

Fields:

- `canonical_path`: `internal/schema/cueson.schema.json`.
- `release_path`: `schema/releases/v1.0.0/cueson.schema.json`.
- `sha256`: lowercase 64-character digest of the common bytes.
- `embedded_copies`: one copy in each of the six binaries.
- `packaged_copies`: one `cueson.schema.json` member in each of the six archives.

Validation rules:

- Both repository paths are regular UTF-8 files without a byte-order mark and are byte-identical.
- `$id`, the root instance `$schema` constant, and `schema_version` equal the Stable Release Identity.
- Every embedded, emitted, and packaged copy is byte-identical to the repository release copy.

## Candidate Revision

Fields:

- `source_revision`: exact 40-character lowercase Git commit.
- `clean`: build metadata reports an unmodified source tree.
- `review_context`: pull-request head or default-branch post-squash proof.

State transitions:

- `specified`: S019 artifacts pass analysis.
- `implemented`: version, schema, release records, and proof behavior pass focused verification.
- `review_candidate`: a clean committed pull-request head produces accepted non-publishing evidence.
- `merged_candidate`: an authorized squash-merge commit on `main` independently produces accepted evidence.
- `publication_eligible`: issue #38 may request exact tag and release authorization for the merged candidate.

## Target Artifact

Fields:

- `goos`: `linux`, `darwin`, or `windows`.
- `goarch`: `amd64` or `arm64`.
- `archive`: deterministic `cueson_1.0.0_<goos>_<goarch>` filename with platform extension.
- `binary`: `cueson` or `cueson.exe`.
- `schema`: root archive member `cueson.schema.json`.
- `license`: root archive member `LICENSE` whose bytes match the repository source and whose archive mode is exactly `0644`.
- `notice`: root archive member `NOTICE` whose bytes match the repository source and whose archive mode is exactly `0644`.
- `sbom`: target-bound SPDX JSON filename.
- `archive_sha256`: lowercase archive digest.
- `sbom_sha256`: lowercase software-bill-of-material document digest.

Validation rules:

- Exactly six unique operating-system and architecture pairs exist in canonical order.
- Every archive contains exactly four safe regular root members.
- Every binary is pure Go, path-trimmed, clean, source-bound, version-marked, and schema-bearing.
- Software-bill-of-material semantics bind the target, version, and source revision; upstream document bytes are not claimed reproducible.

## Candidate Evidence

Fields:

- `version`: `1.0.0`.
- `intended_tag`: `v1.0.0`.
- `source_revision`: Candidate Revision identity.
- `release_schema_sha256`: Immutable v1 Schema digest.
- `license_sha256`: repository and packaged `LICENSE` digest.
- `notice_sha256`: repository and packaged `NOTICE` digest.
- `archive_count`: exactly `6`.
- `sbom_count`: exactly `6`.
- `checksum_count`: exactly `6`.
- `host_executed`: compatible target identity or `null` when execution was not requested.
- `targets`: ordered collection of six Target Artifact evidence records.
- `published`: always `false`.

Validation rules:

- Candidate verification runs without development mode and requires the immutable release schema.
- Evidence is written only after every repository, archive, executable, checksum, schema, legal-file, and software-bill-of-material assertion passes.
- Existing evidence is never overwritten; each accepted proof uses a distinct output path or starts from a clean output directory.

## Candidate Release Notes

Fields:

- `title`: Cueson v1.0.0.
- `highlights`: concise stable user-facing capabilities.
- `boundaries`: truthful absence of signatures, attestations, public schema hosting, and production activation.
- `changelog_link`: exact final tagged URL.

Validation rules:

- The notes remain substantially shorter than the detailed changelog.
- The final non-empty line is exactly the required changelog sentence.
- The notes describe a candidate without claiming publication before issue #38 completes.

## Publication Decision Package

Fields:

- `candidate_revision`: exact post-merge default-branch commit.
- `release_record`: dated changelog, candidate notes, stable identity, and immutable schema digest.
- `artifact_record`: archives, checksums, software bills of materials, legal digests, native smoke evidence, and accepted verifier evidence.
- `quality_record`: green default-branch checks and resolved review and security findings.
- `boundaries`: unsigned, unattested, unhosted schema and unchanged production state unless separately authorized later.

State transitions:

- `prepared`: S019 pull request is green and fully reviewed.
- `candidate`: S019 is merged and exact `main` proof is green.
- `authorized`: the operator grants specific tag and release authority for the named candidate and artifact set.
- `published`: issue #38 completes the protected actions and public read-back.
- `verified`: independent post-publication comparison passes and lifecycle reconciliation can be separately decided.
