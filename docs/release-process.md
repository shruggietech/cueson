# Release process

**Current state:** v0.0.0 publicly released and independently verified

This document defines the protected release lifecycle around Cueson's repository-owned candidate proof and authorized public releases. It does not provide an executable publishing path or grant release authority by itself. The runnable, non-publishing artifact checks and v0.0.0 publication evidence are recorded in [release verification](release-verification.md).

## Authority split

Candidate construction and verification are ordinary reviewed repository work. Final pull-request merge, tag creation or movement, tag push, GitHub Release publication, release-asset publication, milestone closure, immutable release-schema admission, and production `cueson.io` changes are distinct governed actions. Each action must remain within an approved release specification and must receive any explicit operator authorization required by the [project constitution](../.specify/memory/constitution.md) and `AGENTS.md`.

An authorization applies only to the named action and state under review. It does not silently authorize a later retry after material changes, another pull request, another tag, another release, or production publication.

## v0.0.0 released state

The repository has a buildable `0.0.0` executable, an embedded `0.0.0` Cue JSON schema, a byte-identical immutable [versioned release schema](../schema/releases/v0.0.0/cueson.schema.json), published [v0.0.0 release notes](releases/v0.0.0.md), and a verified non-publishing six-target snapshot pipeline. The current CLI and schema boundary is documented in the [CLI contract](cli.md) and [schema baseline](schema.md).

The immutable annotated tag [`v0.0.0`](https://github.com/shruggietech/cueson/tree/v0.0.0) peels to `b294a6952c8bd041d852c502f5d7206c0b58edd6`. The [GitHub Release](https://github.com/shruggietech/cueson/releases/tag/v0.0.0) is public, final, and contains exactly six archives, six target-bound SPDX JSON SBOMs, and one checksum manifest copied byte for byte from accepted default-branch proof run [34543376814](https://github.com/shruggietech/cueson/actions/runs/34543376814).

The release does not publish the canonical schema endpoint, add a signature or attestation, activate production `cueson.io`, implement native subtitle codecs, or close milestone v0.0.0. Those remain distinct governed outcomes. GitHub Actions artifacts from later snapshot runs remain review evidence rather than release assets unless another release is explicitly authorized.

## Candidate readiness

A release decision starts only from a clean, reviewed default-branch commit. Before asking for publication authority, maintainers must verify:

1. Every issue committed to the candidate is closed or truthfully moved, and the Project and milestone agree with that state.
2. The executable version, embedded schema version, and intended tag version are identical.
3. The detailed candidate history is complete under a dated section matching the intended version in [CHANGELOG.md](../CHANGELOG.md), with a fresh empty `[Unreleased]` section retained for later work.
4. Repository formatting, tests, race detection, vet, vulnerability analysis, CodeQL, platform-native tests, and pure-Go target builds are green.
5. The non-publishing snapshot and repository-owned verifier pass for the exact default-branch candidate commit and record its full revision plus `release_schema_sha256` in deterministic evidence.
6. Documentation describes the candidate's actual commands and format-support declarations, distinguishes exact restoration from rendering and conversion, and contains no capability inherited from a different release.
7. No unresolved review or security finding remains.
8. The proposed release does not depend on an unapproved production-domain change.

Candidate readiness proves that a release decision can be made. It does not make the decision and does not grant permission to publish.

## Release preparation lifecycle

For v0.0.0, S012 prepared the publication state without weakening the existing non-publishing controls.

1. The dated v0.0.0 changelog, concise release notes, and versioned release-schema candidate were reviewed together on the S012 pull request.
2. The complete six-target snapshot and repository-owned verifier proved the pull-request head as review evidence only because squash merge created a different commit.
3. The operator completed the final review and specific merge decision after the pull request was green and fully reviewed.
4. The `main` push workflow rebuilt and verified the exact squash-merge commit rather than substituting the pull-request head or pre-merge default-branch revision.
5. The default-branch `release-evidence.json` recorded the full source revision, version, `release_schema_sha256`, target inventory, archive and SBOM digests, counts, compatible-host execution result, and `published: false`.
6. S013 froze the complete decision package, received exact tag and release authority, published only the authorized state, and verified the public result independently.

The checked-in GoReleaser configuration remains snapshot-only and has publication disabled. A release slice must not repurpose that reviewed configuration into a hidden publication mechanism. For v1.0.0, candidate versioning, immutable-schema admission, dated changelog and concise release-note preparation, and exact candidate proof belong to the later release-candidate slice; S018 contract hardening does not perform those actions.

## Pull-request and default-branch candidate binding

Pull-request proof answers whether the reviewed change can produce the expected candidate. It does not identify the final publication commit. Repository policy uses squash merges, so the pull-request head is not the resulting `main` revision and must never be recorded as the release target by prediction.

The non-publishing workflow also runs on pushes to `main`. After the authorized S012 merge, successful run [34543376814](https://github.com/shruggietech/cueson/actions/runs/34543376814) promoted exact 40-character squash-merge revision `b294a6952c8bd041d852c502f5d7206c0b58edd6` from reviewed repository state to the authorized publication target. Future releases must repeat this binding; if the default branch moves before authorization, or artifacts, tools, dependencies, notes, schema bytes, or checks change materially, the decision package must be regenerated and reviewed against the new state.

## Operator decision package

Before requesting tag or GitHub Release authority, present the operator with one bounded package containing:

- the exact 40-character post-squash `main` commit proposed for the intended semantic-version tag;
- the intended version, dated [changelog](../CHANGELOG.md), and reviewed concise release notes;
- the lowercase `release_schema_sha256` for the byte-identical canonical, versioned, embedded, and packaged schema;
- the exact six-archive and six-SBOM inventory, the six archive checksum entries, and every recorded archive and SBOM digest;
- the accepted `release-evidence.json`, including target identities, compatible-host execution result, and `published: false`;
- the green default-branch CI, CodeQL, and release-proof results plus the resolved review and security record;
- every capability limitation that applies to that candidate, plus the explicit absence of signatures, attestations, public schema hosting, or production activation unless those items were separately specified, proved, and authorized.

The operator decision must name the target commit and intended actions. Approval to create or push the tag does not authorize GitHub Release or asset publication, and approval to publish the release does not authorize production-domain work or milestone closure.

## Protected publication actions

| Action | Required boundary |
|---|---|
| Merge the release-preparation pull request | Human operator approval for that specific pull request unless a single-use override explicitly identifies it. |
| Create, move, or push a version tag | Explicit authorization for the exact tag and target commit. Moving an existing tag requires fresh judgment and must never be inferred from creation authority. |
| Publish a GitHub Release or release assets | Explicit authorization for the exact release, notes, and verified artifact set. |
| Admit the immutable release-schema copy | Reviewed release change plus byte-identity proof; released copies are never mutated in place. |
| Close a release milestone | Separate lifecycle decision after publication and post-publication verification, or an explicit operator decision to retire the milestone without publishing. |
| Publish a schema or change production `cueson.io` | Separate post-v1 production specification and explicit production authorization. It is not implied by a GitHub release transaction. |

Authority for one row does not authorize another row. In particular, permission to merge a release-preparation pull request is not permission to create a tag or GitHub Release.

## Tag and GitHub Release boundary

After release preparation reaches the verified default branch, the operator may authorize tag and release publication for the exact commit and artifact set recorded by the successful `main` proof. For v0.0.0, the operator authorized the exact S012 commit and thirteen public files, and S013 read back the resulting tag target, release notes, release state, and asset inventory before treating publication as complete.

The release must use the verified artifacts associated with the approved source commit. A rebuild from different source, a changed dependency, a changed tool version, or a materially different artifact set requires verification and operator judgment again. Draft, partial, or failed publication is not silently promoted to success.

## Immutable schema boundary

The versioned repository schema, the schema embedded in every official binary, and the schema distributed in every release archive must be byte-identical. Their schema identity and version must match the release version.

After release, `schema/releases/v0.0.0/cueson.schema.json` is immutable. A correction that changes contract bytes requires a new software and schema version; it must not overwrite the released copy or move the existing tag to conceal the change.

The canonical `https://cueson.io/schema/v0.0.0/cueson.schema.json` value is an identifier before public-domain activation. Consumers may use the embedded, repository, or release-artifact copy. Publishing that URL is intentionally deferred until the separately specified post-v1 production phase.

## Post-publication verification

Publication is complete only after read-back and independent comparison prove:

- the release tag identifies the authorized commit;
- the GitHub Release identifies the intended tag and is in the intended publication state;
- the published asset inventory is complete and contains no unexpected artifact;
- archive checksums match the published archives;
- every packaged executable contains the verified release marker and expected build identity, the host-compatible executable reports the release version, and every embedded schema reports the schema version;
- every packaged schema matches the immutable tagged repository schema byte-for-byte;
- SBOMs remain associated with the correct target and source revision;
- release notes are the reviewed highlights and end with the exact tagged changelog link;
- no production-domain state changed as a side effect.

If any comparison fails, stop dependent actions, retain the evidence, and ask the operator to choose a recovery path. Do not move a published tag, overwrite an immutable schema, replace assets, or close the milestone by inference.

## Milestone lifecycle

The v0.0.0 milestone remains open after S013 publication verification because milestone closure is a separate lifecycle action that was not included in the operator's release authority. It can be closed only through a separately authorized reconciliation.

For a future abandoned release, its milestone may be retired or its remaining commitments moved only through an explicit recorded decision. Repository status must continue to reflect the actual publication outcome.

## After a successful release

After publication, maintainers retain the fresh `[Unreleased]` section for later work, preserve the released changelog and schema, complete post-merge housekeeping, and record any release follow-up as new issues rather than editing historical evidence. Milestone reconciliation follows only when separately authorized.

Signatures, attestations, public schema hosting, DNS, redirects, TLS, documentation hosting, and other production `cueson.io` work remain separate future outcomes. None is implied by the existing snapshot verifier or by a successful v0.0.0 release.
