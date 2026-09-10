# Release process

**Current state:** v0.0.0 release candidate, not publicly released

This document defines the protected release lifecycle around Cueson's repository-owned candidate proof. It does not provide an executable publishing path and does not authorize any release action. The runnable, non-publishing artifact checks remain in [release verification](release-verification.md).

## Authority split

Candidate construction and verification are ordinary reviewed repository work. Final pull-request merge, tag creation or movement, tag push, GitHub Release publication, release-asset publication, milestone closure, immutable release-schema admission, and production `cueson.io` changes are distinct governed actions. Each action must remain within an approved release specification and must receive any explicit operator authorization required by the [project constitution](../.specify/memory/constitution.md) and `AGENTS.md`.

An authorization applies only to the named action and state under review. It does not silently authorize a later retry after material changes, another pull request, another tag, another release, or production publication.

## v0.0.0 candidate state

The repository currently has a buildable `0.0.0` executable, an embedded `0.0.0` Cue JSON schema, a byte-identical [versioned release-schema candidate](../schema/releases/v0.0.0/cueson.schema.json), prepared [v0.0.0 release notes](releases/v0.0.0.md), and a verified non-publishing six-target snapshot pipeline. The current CLI and schema boundary is documented in the [CLI contract](cli.md) and [schema baseline](schema.md).

The candidate is not a public release. There is no release tag, GitHub Release, published schema endpoint, signature, attestation, or production-domain activation. The versioned schema remains a reviewed candidate until publication binds it to the authorized tag, and candidate artifacts retained by GitHub Actions are review evidence rather than release assets.

S012 prepares the permanent release record and extends the same non-publishing proof to pushes on `main`. Its pull request proves the reviewed head only. If the operator later squash-merges it, the resulting default-branch commit receives its own proof and becomes the proposed publication target. Neither merge nor successful default-branch proof publishes v0.0.0 or closes the v0.0.0 milestone.

## Candidate readiness

A release decision starts only from a clean, reviewed default-branch commit. Before asking for publication authority, maintainers must verify:

1. Every issue committed to the candidate is closed or truthfully moved, and the Project and milestone agree with that state.
2. The executable version, embedded schema version, and intended tag version are identical.
3. The detailed candidate history is complete under the dated `[0.0.0]` section in [CHANGELOG.md](../CHANGELOG.md), with a fresh empty `[Unreleased]` section retained for later work.
4. Repository formatting, tests, race detection, vet, vulnerability analysis, CodeQL, platform-native tests, and pure-Go target builds are green.
5. The non-publishing snapshot and repository-owned verifier pass for the exact default-branch candidate commit and record its full revision plus `release_schema_sha256` in deterministic evidence.
6. Documentation describes shipped commands and `envelope_only` format support without claiming native SRT or WebVTT codecs.
7. No unresolved review or security finding remains.
8. The proposed release does not depend on an unapproved production-domain change.

Candidate readiness proves that a release decision can be made. It does not make the decision and does not grant permission to publish.

## Release preparation lifecycle

S012 prepares the publication state without weakening the existing non-publishing controls.

1. Review the dated v0.0.0 changelog, concise release notes, and versioned release-schema candidate together on the S012 pull request.
2. Prove the pull-request head with the complete six-target snapshot and repository-owned verifier. Treat that run as review evidence only because squash merge creates a different commit.
3. Stop for the operator's final review and specific merge decision after the pull request is green and fully reviewed.
4. If the operator squash-merges S012, wait for the `main` push workflow to rebuild and verify the exact resulting commit. Do not substitute the pull-request head or pre-merge default-branch revision.
5. Inspect the default-branch `release-evidence.json` and confirm the full source revision, version, `release_schema_sha256`, target inventory, archive and SBOM digests, counts, compatible-host execution result, and `published: false`.
6. Present the complete decision package described below and stop before tag creation or release publication until the operator grants authority for those exact actions and artifacts.

The checked-in GoReleaser configuration remains snapshot-only and has publication disabled. A release slice must not repurpose that reviewed configuration into a hidden publication mechanism.

## Pull-request and default-branch candidate binding

Pull-request proof answers whether the reviewed change can produce the expected candidate. It does not identify the final publication commit. Repository policy uses squash merges, so the pull-request head is not the resulting `main` revision and must never be recorded as the release target by prediction.

The non-publishing workflow also runs on pushes to `main`. After an authorized S012 merge, only a successful run for the exact 40-character squash-merge revision can promote that revision from reviewed repository state to proposed publication target. If the default branch moves again before authorization, or if artifacts, tools, dependencies, notes, schema bytes, or checks change materially, the decision package must be regenerated and reviewed against the new state.

## Operator decision package

Before requesting tag or GitHub Release authority, present the operator with one bounded package containing:

- the exact 40-character post-squash `main` commit proposed for tag `v0.0.0`;
- version `0.0.0`, the dated [changelog](../CHANGELOG.md), and the reviewed [release notes](releases/v0.0.0.md);
- the lowercase `release_schema_sha256` for the byte-identical canonical, versioned, embedded, and packaged schema;
- the exact six-archive and six-SBOM inventory, the six archive checksum entries, and every recorded archive and SBOM digest;
- the accepted `release-evidence.json`, including target identities, compatible-host execution result, and `published: false`;
- the green default-branch CI, CodeQL, and release-proof results plus the resolved review and security record;
- the envelope-only limitation and the explicit absence of native SRT/WebVTT ingest, model-driven render, conversion, signatures, attestations, public schema hosting, and production activation.

The operator decision must name the target commit and intended actions. Approval to create or push the tag does not authorize GitHub Release or asset publication, and approval to publish the release does not authorize production-domain work or milestone closure.

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

After release preparation reaches the verified default branch, the operator may authorize tag and release publication for the exact commit and artifact set recorded by the successful `main` proof. The publishing actor must read back the resulting tag target, release notes, release state, and asset inventory before treating an API success as completion.

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

The v0.0.0 milestone remains open when S012 finishes because release preparation and public release completion are different states. After an authorized v0.0.0 publication passes post-publication verification, the milestone can be closed through its separately governed lifecycle action.

If the operator decides not to publish v0.0.0, the milestone may be retired or its remaining commitments moved only through an explicit recorded decision. The repository must continue to say that no v0.0.0 release occurred.

## After a successful release

After publication and milestone reconciliation, maintainers retain the fresh `[Unreleased]` section for later work, preserve the released changelog and schema, complete post-merge housekeeping, and record any release follow-up as new issues rather than editing historical evidence.

Signatures, attestations, public schema hosting, DNS, redirects, TLS, documentation hosting, and other production `cueson.io` work remain separate future outcomes. None is implied by the existing snapshot verifier or by a successful v0.0.0 release.
