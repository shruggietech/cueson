# Research: Publish and Verify v1.0.0

## Decision 1: Publish only the exact accepted default-branch candidate

**Decision**: The only eligible tag target is `2cad4c816340404289b4d1d87179a4071713bb46`, backed by accepted release-proof run `34621429626` and artifact `10273380044`.

**Rationale**: This is the verified squash-merge commit on `main`, not a pull-request head or predicted merge result. Its default-branch CI, CodeQL, release proof, and native packaged smoke checks are green.

**Alternatives considered**: Tagging the S020 documentation branch would make reviewed follow-up prose part of the product release rather than the accepted candidate. Rebuilding would create a new unapproved input set. Tagging an earlier S019 head would not identify merged default-branch state.

## Decision 2: Separate kickoff authority from publication authority

**Decision**: S020 may be specified, planned, analyzed, and preflighted now. It must halt before tag creation, tag push, or GitHub Release publication until the operator explicitly authorizes those actions for the exact candidate and artifact.

**Rationale**: Repository policy treats tag and release publication as separately protected mutations. “Kick off S020” supplies work-slice intent but does not identify those mutations as authorized.

**Alternatives considered**: Inferring release authority from kickoff would violate the recorded boundary. Stopping before any work would waste the safe preparation window and increase artifact-expiry risk.

## Decision 3: Publish exactly thirteen user-facing files from one retained artifact

**Decision**: The public set is the six archives, six matching SPDX JSON SBOMs, and `cueson_1.0.0_checksums.txt` extracted from artifact `10273380044` without rebuild, rename, normalization, or substitution.

**Rationale**: The candidate verifier already bound these files to the exact source revision, schema, legal files, target matrix, and native smoke evidence. Publishing internal metadata would confuse users and expose files that are evidence rather than release deliverables.

**Alternatives considered**: Rebuilding from the tag would generate a different unapproved bundle. Publishing `release-evidence.json`, `artifacts.json`, `metadata.json`, `config.yaml`, or extracted directories would expand the public contract without user value.

## Decision 4: Freeze exact public-file digests before publication

**Decision**: Record the names, sizes, and SHA-256 digests of all thirteen public files in `contracts/release-publication-contract.json` before remote mutation.

**Rationale**: Names alone cannot detect substitution. A frozen contract makes pre-publication inputs and post-publication downloads directly comparable.

**Alternatives considered**: Relying only on the checksum manifest covers archives but not SBOMs or the manifest itself. Relying on GitHub-reported sizes cannot prove byte identity.

## Decision 5: Treat publication as a read-back transaction

**Decision**: After authorization, create and verify the annotated tag first, then create the release and upload assets, then immediately read the release and every public asset back before updating repository status claims.

**Rationale**: GitHub API success does not prove the intended reference, body, state, or complete asset set is public. Sequential read-back constrains partial failure and keeps evidence causally clear.

**Alternatives considered**: Parallel tag and release actions would make partial failure harder to reason about. Updating documentation before public verification would create false current-state claims.

## Decision 6: Independently verify public bytes and stable semantics

**Decision**: Download all thirteen public assets into a new clean directory, compare all digests with the frozen contract, apply the checksum manifest as a six-archive bijection, inspect archives and SBOMs, and execute the Windows amd64 package on the current host.

**Rationale**: This proves the actual user download surface. The accepted candidate evidence remains the detailed cross-platform proof, while public verification establishes that GitHub serves those exact accepted bytes.

**Alternatives considered**: Trusting upload inputs or GitHub asset names would not verify downloads. Rebuilding on each platform would test different bytes.

## Decision 7: Publish unchanged reviewed notes

**Decision**: Use `docs/releases/v1.0.0.md` unchanged as the release body and require the exact final line `Full changelog: https://github.com/shruggietech/cueson/blob/v1.0.0/CHANGELOG.md`.

**Rationale**: The release highlights and boundaries were reviewed with the candidate. Editing them during publication would create unreviewed release meaning.

**Alternatives considered**: Generated GitHub notes would not preserve the project’s reviewed narrative. Copying the full changelog would violate the short-notes policy.

## Decision 8: Keep the tag fixed while documenting publication afterward

**Decision**: Tag the accepted candidate commit, then update current-state repository documentation on the S020 branch after public verification. The follow-up documentation is intentionally outside the tag.

**Rationale**: The tag must identify the verified candidate, while post-publication facts cannot truthfully exist in that candidate before the release occurs. This is the same safe transaction shape used for v0.0.0.

**Alternatives considered**: Tagging a later documentation commit would discard the accepted candidate identity. Prewriting released-state claims would be false before publication.

## Decision 9: Preserve production, signature, and API boundaries

**Decision**: S020 does not activate `cueson.io`, host the schema, add signatures or attestations, or declare Go packages public and stable.

**Rationale**: Those are separate products or trust transitions that require their own designs and authority. The v1 CLI and Cue JSON release is complete without them.

**Alternatives considered**: Bundling production hosting into release publication would expand risk and delay usable software. Adding unsigned placeholder claims would be misleading.

## Decision 10: Reconcile planning only after evidence and explicit authority

**Decision**: Issue #38 closes through the reviewed S020 pull request. Epic #29 and milestone v1.0.0 may close only after publication evidence is complete, dependent work is closed, and the operator authorizes those protected planning mutations.

**Rationale**: Publication success and planning closure are related but distinct state changes. Sequencing them prevents a partially published release from appearing complete.

**Alternatives considered**: Closing all surfaces during upload would erase the distinction between attempted and verified publication. Leaving completed planning open forever would misstate delivery after authorized reconciliation.
