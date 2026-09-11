# Quickstart: Verify the v0.0.0 Publication

## 1. Verify the authorized inputs

Confirm the accepted run, artifact, source revision, schema digest, reviewed release notes, and thirteen-file contract before mutation. Confirm `v0.0.0` and its GitHub Release do not already exist.

Expected: every expected file exists with the contract digest, the release body formatter produces no change, the accepted evidence identifies `b294a6952c8bd041d852c502f5d7206c0b58edd6`, and no conflicting public state exists.

## 2. Publish the immutable tag

Create an unsigned annotated `v0.0.0` tag at the authorized revision, prove the local peeled target, push only that reference, and read back both the remote tag object and peeled target.

Expected: `refs/tags/v0.0.0^{}` equals the authorized 40-character commit. Do not move or replace the tag after push.

## 3. Publish and read back the GitHub Release

Create `Cueson v0.0.0` as a public, non-draft, non-prerelease release using `docs/releases/v0.0.0.md` and the thirteen explicit accepted files. Immediately read the release record back from GitHub.

Expected: release state, body, and exact asset inventory match `contracts/release-publication-contract.json`. A partial or mismatched result blocks dependent work.

## 4. Verify public downloads independently

Download the release into a newly created clean temporary directory. Compare all thirteen file digests with the contract, apply the six-entry checksum manifest, and execute the Windows amd64 binary's `version`, `schema --version`, and schema-emission commands.

Expected: every public byte matches accepted evidence, all six archive checksums pass, compatible commands report `0.0.0`, and emitted schema bytes match the immutable tagged schema. Because the public archives and SBOMs are byte-identical to the accepted proof files, the repository-owned six-target structural and semantic verification applies to the public set.

## 5. Update and verify repository records

Update release-facing status prose only after public verification. Keep `docs/releases/v0.0.0.md` and `schema/releases/v0.0.0/cueson.schema.json` unchanged, preserve the `envelope_only` limitation, and record the post-tag documentation work under `[Unreleased]`.

Run the complete local gate documented in `docs/release-verification.md`, including formatter, product tests, race detection, vet, standalone tool tests, documentation and brand verification, actionlint, vulnerability analysis, release snapshot, repository-owned release verification, and `git diff --check`.

Expected: all checks pass, repository text contains no current unreleased status claim, no native-codec overclaim appears, and no tagged release record changed.

## 6. Audit protected boundaries and GitHub delivery

Verify milestone v0.0.0 remains open, production `cueson.io` and public schema hosting remain unchanged, issue #27 appears exactly once in `cueson Delivery`, Slice is `S013`, governed Stage matches the lifecycle, and default Status is empty.

Expected: publication is verified without milestone closure, production changes, signatures, attestations, or schema hosting.

## 7. Validate the official pull request

Push `S013-publish-v0-release`, open the official pull request with `Closes #27`, read its body back, and wait for every current-head CI, CodeQL, release-proof, policy, Codex, security, and review result. Address every finding and request exactly one second Codex review only when round one reports findings.

Expected: the pull request is green, fully reviewed, conflict-free, and ready for the operator's final merge ritual. Do not merge or auto-merge it.
