# Data Model: Publish and Verify v1.0.0

## Authorized Candidate

Fields:

- `version`: exact semantic version `1.0.0`.
- `source_revision`: `2cad4c816340404289b4d1d87179a4071713bb46`.
- `proof_run_id`: `34621429626`.
- `proof_artifact_id`: `10273380044`.
- `artifact_name`: `cueson-1.0.0-candidate-2cad4c816340404289b4d1d87179a4071713bb46`.
- `release_schema_sha256`: `1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541`.
- `expires_at`: `2026-09-14T16:21:38Z`.

Validation rules:

- The revision is the exact S019 squash-merge commit on `main`.
- Default-branch CI, CodeQL, release proof, and Linux, Windows, and macOS packaged smoke checks are green for the revision.
- Candidate evidence reports `published: false`, version `1.0.0`, intended tag `v1.0.0`, and the exact source and schema identities.
- Unavailability or mismatch blocks publication; no rebuild or substitution is inferred.

## Release Tag

Fields:

- `name`: exact `v1.0.0`.
- `kind`: unsigned annotated tag.
- `target`: Authorized Candidate source revision.
- `annotation`: concise v1.0.0 release identity.

State transitions:

- `absent`: verified locally and remotely during preflight.
- `authorized`: operator specifically authorizes creation and publication at the frozen candidate.
- `local_verified`: local tag object is annotated and peels to the exact target.
- `remote_verified`: pushed reference reads back and peels to the exact target.
- `exception`: any conflicting or partial state blocks movement, replacement, or recreation.

## GitHub Release

Fields:

- `tag`: Release Tag name.
- `name`: `Cueson v1.0.0`.
- `body_source`: unchanged `docs/releases/v1.0.0.md`.
- `draft`: `false`.
- `prerelease`: `false`.
- `target`: Authorized Candidate source revision.
- `assets`: Public Asset Set.

Validation rules:

- The body passes `scripts/github-format` unchanged and reads back with intentional Markdown structure.
- The final non-empty body line is the tagged changelog link.
- The release is created only after remote tag verification and exact authorization.
- Partial upload or metadata mismatch enters exception state and blocks released-state documentation.

## Public Asset

Fields:

- `name`: exact release filename.
- `kind`: `archive`, `sbom`, or `checksums`.
- `size`: accepted byte length.
- `sha256`: accepted lowercase digest.
- `goos`: platform for archives and SBOMs.
- `goarch`: architecture for archives and SBOMs.

Validation rules:

- Exactly six archives, six SBOMs, and one checksum file exist.
- Names, sizes, and digests are unique where required and match `contracts/release-publication-contract.json`.
- Each archive has exactly four safe root members: executable, `cueson.schema.json`, `LICENSE`, and `NOTICE`.
- Each archive schema has the Authorized Candidate schema digest.
- Each SBOM identifies its platform, version, and source revision semantically.
- No internal evidence, metadata, configuration, or extracted directory is public.

## Publication Evidence

Fields:

- `tag_object`: remote tag object identifier.
- `peeled_revision`: exact Authorized Candidate revision.
- `release_url`: public GitHub Release URL.
- `body_verified`: structural and exact-source comparison result.
- `asset_inventory_verified`: exact thirteen-name comparison result.
- `asset_digests_verified`: thirteen-file SHA-256 comparison result.
- `checksums_verified`: exact six-archive bijection result.
- `archives_verified`: member, legal-file, schema, version, and target results.
- `sboms_verified`: platform, version, and revision semantic results.
- `host_execution_verified`: Windows amd64 runtime result.
- `protected_boundaries_verified`: production, schema-hosting, signature, attestation, milestone, epic, and merge audit.

State transitions:

- `unpublished`: preflight passes and protected public state is absent.
- `published_pending_verification`: remote tag and release exist, but no success claim is made.
- `verified`: every public and protected-state assertion passes.
- `exception`: mismatch or partial state is retained for operator-directed recovery; no tag movement or silent asset replacement occurs.

## Release Documentation Record

Fields:

- `public_release_url`: immutable GitHub Release URL.
- `tag_url`: immutable tag URL.
- `source_revision`: Authorized Candidate revision.
- `proof_run_url`: accepted pre-publication proof.
- `capabilities`: stable Cue JSON, CLI, SubRip, WebVTT, validation, inspection, completion, rendering, and bounded conversion behavior.
- `deferred_boundaries`: public schema hosting, production domain activation, signatures, attestations, and external Go API stability.

Validation rules:

- Current status documents contain no claim that v1.0.0 remains unpublished after verification.
- Historical passages may describe candidate preparation when tense and context make that history explicit.
- The tagged release notes and versioned schema remain unchanged.
- Links resolve, Markdown layout passes, and capability limitations remain explicit.

## Delivery Reconciliation

Fields:

- `issue`: #38, closed by the merged S020 pull request.
- `epic`: #29, eligible for closure only after all child outcomes are complete.
- `milestone`: v1.0.0, eligible for closure only after its open issue count is zero and authority is explicit.
- `project_item`: issue #38 represented exactly once with Slice `S020`.
- `stage`: governed lifecycle value.
- `default_status`: absent.

Validation rules:

- Specced and In progress state can be recorded during authorized local work.
- Release verification is used only after publication has occurred and before final repository reconciliation completes.
- Done is used only after the issue has closed through merge.
- Epic and milestone closure never occur by inference from release upload alone.
