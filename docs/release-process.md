# Release process

**Current state:** v0.0.0 release candidate, not publicly released

This document defines the protected release lifecycle around Cueson's repository-owned candidate proof. It does not provide an executable publishing path and does not authorize any release action. The runnable, non-publishing artifact checks remain in [release verification](release-verification.md).

## Authority split

Candidate construction and verification are ordinary reviewed repository work. Final pull-request merge, tag creation or movement, tag push, GitHub Release publication, release-asset publication, milestone closure, immutable release-schema admission, and production `cueson.io` changes are distinct governed actions. Each action must remain within an approved release specification and must receive any explicit operator authorization required by the [project constitution](../.specify/memory/constitution.md) and `AGENTS.md`.

An authorization applies only to the named action and state under review. It does not silently authorize a later retry after material changes, another pull request, another tag, another release, or production publication.

## v0.0.0 candidate state

The repository currently has a buildable `0.0.0` executable, an embedded `0.0.0` Cue JSON schema, and a verified non-publishing six-target snapshot pipeline. The current CLI and schema boundary is documented in the [CLI contract](cli.md) and [schema baseline](schema.md).

The candidate is not a public release. There is no release tag, GitHub Release, immutable repository release-schema copy, published schema endpoint, signature, attestation, or production-domain activation. Candidate artifacts retained by GitHub Actions are review evidence, not release assets.

S010 prepares documentation and milestone evidence only. Its pull request is expected to close the documentation issue and the v0.0.0 foundation epic when the operator merges it. That merge does not publish v0.0.0 and does not close the v0.0.0 milestone.

## Candidate readiness

A release decision starts only from a clean, reviewed default-branch commit. Before asking for publication authority, maintainers must verify:

1. Every issue committed to the candidate is closed or truthfully moved, and the Project and milestone agree with that state.
2. The executable version, embedded schema version, and intended tag version are identical.
3. The detailed candidate history is complete under `[Unreleased]` in [CHANGELOG.md](../CHANGELOG.md).
4. Repository formatting, tests, race detection, vet, vulnerability analysis, CodeQL, platform-native tests, and pure-Go target builds are green.
5. The non-publishing snapshot and repository-owned verifier pass for the exact candidate commit.
6. Documentation describes shipped commands and `envelope_only` format support without claiming native SRT or WebVTT codecs.
7. No unresolved review or security finding remains.
8. The proposed release does not depend on an unapproved production-domain change.

Candidate readiness proves that a release decision can be made. It does not make the decision and does not grant permission to publish.

## Release preparation lifecycle

A later release-specific slice must prepare the publication state without weakening the existing non-publishing controls.

1. Select and record the exact default-branch commit proposed for release.
2. Reconfirm the intended semantic version and verify software/schema lockstep.
3. Convert the accumulated `[Unreleased]` entries into a dated version section only as part of the reviewed release change. Until that change is authorized and published, the candidate remains unreleased.
4. Create `schema/releases/v0.0.0/cueson.schema.json` as the immutable repository copy for this version and prove it is byte-identical to the canonical embedded schema source. S010 deliberately does not create this path.
5. Prepare concise GitHub release notes containing highlights only. For v0.0.0, the final line must be exactly `Full changelog: https://github.com/shruggietech/cueson/blob/v0.0.0/CHANGELOG.md`.
6. Rebuild and verify the complete release candidate from the exact proposed commit using the process in [release verification](release-verification.md).
7. Present the commit, version, changelog, immutable schema digest, artifact inventory, checksums, SBOM evidence, release notes, checks, and unresolved limitations to the operator.
8. Stop before every protected publication action until the operator grants the required specific authorization.

The checked-in GoReleaser configuration remains snapshot-only and has publication disabled. A release slice must not repurpose that reviewed configuration into a hidden publication mechanism.

## Protected publication actions

| Action | Required boundary |
|---|---|
| Merge the release-preparation pull request | Human operator approval for that specific pull request unless a single-use override explicitly identifies it. |
| Create, move, or push `v0.0.0` | Explicit authorization for the exact tag and target commit. Moving an existing tag requires fresh judgment and must never be inferred from creation authority. |
| Publish a GitHub Release or release assets | Explicit authorization for the exact release, notes, and verified artifact set. |
| Admit the immutable release-schema copy | Reviewed release change plus byte-identity proof; released copies are never mutated in place. |
| Close the v0.0.0 milestone | Separate lifecycle decision after publication and post-publication verification, or an explicit operator decision to retire the milestone without publishing. |
| Publish a schema or change production `cueson.io` | Separate post-v1 production specification and explicit production authorization. It is not part of the v0.0.0 transaction. |

Authority for one row does not authorize another row. In particular, permission to merge a release-preparation pull request is not permission to create a tag or GitHub Release.

## Tag and GitHub Release boundary

After release preparation reaches the verified default branch, the operator may authorize tag and release publication for the exact recorded commit and artifact set. The publishing actor must read back the resulting tag target, release notes, release state, and asset inventory before treating an API success as completion.

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
- every packaged executable and embedded schema reports the release version;
- every packaged schema matches the immutable tagged repository schema byte-for-byte;
- SBOMs remain associated with the correct target and source revision;
- release notes are the reviewed highlights and end with the exact tagged changelog link;
- no production-domain state changed as a side effect.

If any comparison fails, stop dependent actions, retain the evidence, and ask the operator to choose a recovery path. Do not move a published tag, overwrite an immutable schema, replace assets, or close the milestone by inference.

## Milestone lifecycle

The v0.0.0 milestone remains open when S010 finishes because repository-foundation readiness and public release completion are different states. After an authorized v0.0.0 publication passes post-publication verification, the milestone can be closed through its separately governed lifecycle action.

If the operator decides not to publish v0.0.0, the milestone may be retired or its remaining commitments moved only through an explicit recorded decision. The repository must continue to say that no v0.0.0 release occurred.

## After a successful release

After publication and milestone reconciliation, maintainers create a fresh `[Unreleased]` section for later work, preserve the released changelog and schema, complete post-merge housekeeping, and record any release follow-up as new issues rather than editing historical evidence.

Signatures, attestations, public schema hosting, DNS, redirects, TLS, documentation hosting, and other production `cueson.io` work remain separate future outcomes. None is implied by the existing snapshot verifier or by a successful v0.0.0 release.
