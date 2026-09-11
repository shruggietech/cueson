# Research: Publish and Verify v0.0.0

## Decision 1: Publish the exact accepted default-branch candidate

**Decision**: Create an unsigned annotated tag `v0.0.0` that peels to `b294a6952c8bd041d852c502f5d7206c0b58edd6`, then publish the GitHub Release from the accepted run `34543376814` artifact only.

**Rationale**: The post-S012 default-branch run is the first proof that binds the final squash-merge revision itself. The operator's authorization identifies this previously proposed candidate and asset set, so no reconstruction or target movement is necessary.

**Alternatives considered**: Tagging the S013 branch would include post-release documentation rather than the accepted candidate. Rebuilding would produce a new unapproved input set. Using the S012 pull-request head would not identify the squash-merge commit.

## Decision 2: Publish thirteen user-facing files

**Decision**: Publish six archives, six corresponding SPDX JSON SBOMs, and `cueson_0.0.0_checksums.txt`. Exclude `release-evidence.json`, GoReleaser metadata and configuration, and extracted build directories.

**Rationale**: Archives, SBOMs, and checksums are useful public deliverables. The other retained files are internal build and decision evidence, and publishing the evidence file would expose a record whose `published: false` value correctly describes the non-publishing verifier but could be misread as current release state.

**Alternatives considered**: Publishing the entire workflow artifact would leak internal mechanics and create an ambiguous public record. Omitting SBOMs would discard verified supply-chain information. Publishing a separate checksum for SBOMs would require generating an unreviewed asset.

## Decision 3: Use one fail-closed publication transaction

**Decision**: Reconfirm that neither tag nor release exists, verify the exact local source files and reviewed notes, create and inspect the local annotated tag, push only that tag, verify the remote peeled target, and create the release with all thirteen assets in one explicit command. Read back state immediately after every remote mutation.

**Rationale**: This ordering catches conflicts before mutation and makes the tag immutable before release creation. A single release-create command minimizes partial state while read-back detects incomplete uploads.

**Alternatives considered**: Creating a draft would contradict the authorized final release. Uploading assets piecemeal expands partial-state risk. A publishing workflow would add persistent authority and unnecessary code for a one-time release.

## Decision 4: Prove public bytes by identity with accepted evidence

**Decision**: Download all public assets into a new clean directory and compare all thirteen SHA-256 digests with the pre-publication contract derived from accepted evidence. Apply the six-entry checksum manifest, execute the compatible Windows amd64 archive, compare its emitted schema with the tagged schema, and rely on exact byte identity with the already accepted verifier evidence for the five foreign executables and all archive/SBOM structural assertions.

**Rationale**: Cryptographic identity proves that GitHub serves the same files inspected by the repository-owned verifier. Repeating foreign-binary structural inspection with a new implementation would add no stronger evidence and could diverge from the accepted authority.

**Alternatives considered**: Trusting GitHub asset names alone cannot detect substitution. Rebuilding after publication would test different bytes. Adding another permanent verifier would duplicate the existing release authority and widen S013 unnecessarily.

## Decision 5: Preserve public notes exactly

**Decision**: Format-check `docs/releases/v0.0.0.md`, require that formatting makes no byte change, use it directly as the release body, then compare the GitHub read-back body with the source document and inspect rendered structure.

**Rationale**: These notes were reviewed with the candidate and already contain the required limitations and tagged changelog suffix. Editing them during publication would change the authorized package.

**Alternatives considered**: GitHub-generated notes would be unreviewed and overly broad. Copying the changelog would violate the highlights-only convention. A manually entered body risks whitespace corruption.

## Decision 6: Record publication without altering tagged history

**Decision**: Update current repository status documentation on the S013 branch after public verification and add a concise entry under the fresh `[Unreleased]` changelog section. Do not change the tagged `v0.0.0` release notes or versioned schema.

**Rationale**: The tag must remain the reviewed S012 candidate, while the default branch needs a truthful record of the subsequent release operation. The `[Unreleased]` entry accurately describes documentation and operational history added after the tag.

**Alternatives considered**: Amending the tagged changelog is impossible without moving the tag. Editing release notes after publication would break source-to-public identity. Leaving repository prose unchanged would falsely claim the project remains unreleased.

## Decision 7: Keep separately governed work separate

**Decision**: Leave milestone v0.0.0 open, do not publish the canonical schema endpoint, do not touch production `cueson.io`, and do not add signatures, attestations, native codecs, or auto-merge.

**Rationale**: The operator authorized tag and GitHub Release publication, not these distinct transitions. The release process explicitly separates them.

**Alternatives considered**: Treating release authorization as blanket production authority would violate the operating contract and erase useful lifecycle boundaries.
