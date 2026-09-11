# Quickstart: Publish and Verify v1.0.0

## 1. Verify the frozen candidate

Confirm source revision `2cad4c816340404289b4d1d87179a4071713bb46`, run `34621429626`, artifact `10273380044`, artifact expiry, schema digest, release notes, and the thirteen-file contract. Confirm local and remote `v1.0.0` tags and the GitHub Release do not exist.

Expected: every candidate assertion and file digest matches, the release body formatter produces no change, and no conflicting public state exists.

## 2. Halt for exact publication authority

Present the operator with the exact commit, run, artifact, tag, release name, and thirteen-file inventory. Do not create or push the tag or publish the release until the operator explicitly authorizes those actions.

Expected: the authorization names or unambiguously accepts this frozen transaction. If authority is absent or the artifact expires, stop.

## 3. Publish the immutable tag

After authorization, create unsigned annotated tag `v1.0.0` at the frozen revision, prove the local peeled target, push only that reference, and read back both the remote tag object and peeled target.

Expected: `refs/tags/v1.0.0^{}` equals the authorized 40-character commit. Never move or replace the tag.

## 4. Publish and read back the GitHub Release

Create `Cueson v1.0.0` as a public, non-draft, non-prerelease release using unchanged `docs/releases/v1.0.0.md` and the thirteen explicit accepted files. Immediately read the release record back from GitHub.

Expected: release state, body, target, and exact asset inventory match `contracts/release-publication-contract.json`. A partial or mismatched result blocks dependent work.

## 5. Verify public downloads independently

Download the release into a newly created clean temporary directory. Compare all thirteen file sizes and digests with the contract, apply the six-entry checksum manifest, inspect archives and SBOM semantics, and execute the Windows amd64 binary's version and schema commands.

Expected: every public byte matches accepted evidence, all archive checksums pass, the compatible commands report `1.0.0`, emitted schema bytes match the immutable tagged schema, and stable platform and revision semantics remain intact.

## 6. Update and verify repository records

Only after public verification, update current release identity, download, status, process, architecture, schema, and verification prose. Keep `docs/releases/v1.0.0.md` and `schema/releases/v1.0.0/cueson.schema.json` unchanged and record post-tag documentation under `[Unreleased]`.

Expected: the repository truthfully describes the public release and deferred boundaries without changing tagged history.

## 7. Complete local and hosted delivery gates

Run the complete local gate documented in `docs/release-verification.md`, including formatter, product tests, race detection, vet, standalone tool tests, documentation and brand verification, actionlint, vulnerability analysis, release snapshot, repository-owned release verification, encoding, mojibake, and `git diff --check` checks. After separate push and pull-request authority, publish an official pull request with `Closes #38`, read its body back, and wait for every current-head CI, CodeQL, release-proof, policy, Codex, security, and review result.

Expected: every check is green, every finding is resolved, no more than two Codex rounds occur, and the pull request is ready for the operator's final merge ritual.

## 8. Reconcile delivery only within authority

Verify issue #38 appears once in `cueson Delivery`, Slice is `S020`, governed Stage matches lifecycle, and default Status is empty. After merge, reconcile issue #38 and move it to Done. Close epic #29 and milestone v1.0.0 only if all work is complete and the operator separately authorizes those closures.

Expected: GitHub planning agrees with verified delivery without an unauthorized release, production, milestone, epic, or merge mutation.
