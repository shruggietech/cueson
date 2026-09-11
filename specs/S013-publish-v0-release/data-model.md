# Data Model: Publish and Verify v0.0.0

S013 changes no application data. These entities define the authorized release transaction, its public state, and its verification record.

## Authorized Candidate

Fields:

- `version`: `0.0.0`.
- `tag`: `v0.0.0`.
- `source_revision`: `b294a6952c8bd041d852c502f5d7206c0b58edd6`.
- `release_schema_sha256`: `d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975`.
- `proof_run`: GitHub Actions run `34543376814`.
- `proof_artifact`: artifact `10178231596`.

Validation rules:

- The revision is the current accepted S012 squash-merge commit on `main` at authorization time.
- CI, CodeQL, and release proof succeeded for that exact revision.
- The retained evidence identifies the same version, revision, schema digest, and six targets.
- The tag and release are absent before the authorized transaction begins.

## Release Tag

Fields:

- `name`: `v0.0.0`.
- `kind`: annotated Git tag.
- `signed`: false.
- `peeled_target`: Authorized Candidate source revision.

State transitions:

- `absent`: preflight confirms no local or remote release tag.
- `local`: one annotated tag object exists and peels to the authorized revision.
- `remote`: the exact tag object is pushed and its remote peeled reference identifies the authorized revision.
- `immutable`: after push, S013 never moves, replaces, or deletes it.

## Public Asset

Fields:

- `name`: exact filename in `release-publication-contract.json`.
- `kind`: archive, SPDX JSON SBOM, or checksum manifest.
- `sha256`: lowercase 64-character digest of the accepted file.
- `target`: operating-system and architecture tuple for archives and SBOMs, otherwise absent.

Validation rules:

- Exactly thirteen unique names exist.
- Exactly six are archives, six are matching SBOMs, and one is the checksum manifest.
- Source bytes come from the accepted workflow artifact without rebuild, regeneration, rename, or substitution.
- Internal evidence, build metadata, configuration, and extracted directories are absent from the public set.

## GitHub Release

Fields:

- `tag_name`: `v0.0.0`.
- `name`: `Cueson v0.0.0`.
- `draft`: false.
- `prerelease`: false.
- `body_source`: `docs/releases/v0.0.0.md` at the authorized candidate.
- `assets`: ordered or unordered set of thirteen Public Asset records.

Validation rules:

- The release is publicly visible and associated with the immutable Release Tag.
- The body read back from GitHub matches the reviewed source Markdown and ends with the tagged changelog link.
- Asset-set equality is name- and digest-based; order is not significant.

## Publication Verification

Fields:

- `remote_tag_target`: full peeled source revision.
- `release_state`: tag, name, draft, prerelease, body, URL, and asset metadata read from GitHub.
- `download_directory`: newly created clean temporary directory.
- `downloaded_assets`: thirteen name-and-digest pairs.
- `checksum_result`: exact six-entry archive bijection and six matching digest results.
- `host_execution`: Windows amd64 `version`, `schema --version`, and schema emission results.
- `protected_state`: milestone, production, signature, attestation, and schema-hosting observations.

Validation rules:

- Remote tag target, release record, body, asset names, and every downloaded digest match the authorized contract.
- Exact digest identity transfers the already accepted archive and SBOM structural proof to the publicly downloaded bytes.
- Compatible host execution reports `0.0.0`, and emitted schema bytes equal the tagged schema.
- Any mismatch blocks the verified state and all dependent documentation claims.

State transitions:

- `unpublished`: preflight passes and protected state is unchanged.
- `published_pending_verification`: remote tag and release exist, but no success claim is made.
- `verified`: every public and protected-state assertion passes.
- `exception`: any mismatch or partial state is retained for operator-directed recovery; no tag movement or asset replacement occurs by inference.

## Release Documentation Record

Fields:

- `public_release_url`: immutable GitHub Release URL.
- `tag_url`: immutable Git tag URL.
- `source_revision`: Authorized Candidate revision.
- `proof_run_url`: accepted pre-publication proof.
- `capability`: `envelope_only` for SubRip and WebVTT.
- `deferred_boundaries`: native codecs, public schema hosting, production domain, signatures, attestations, and milestone closure.

Validation rules:

- Current status documents contain no claim that the project or v0.0.0 is unreleased.
- Historical passages may describe the pre-publication process when tense and context make that history explicit.
- The tagged release notes and versioned schema remain unchanged.
- Links resolve, Markdown layout passes, and native format limitations remain explicit.
