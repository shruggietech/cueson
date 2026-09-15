# v1.1.0 publication decision contract

**State:** Preparation only. This contract does not select a release commit, populate approval, or authorize a tag, GitHub Release, release assets, public schema, production deployment, milestone closure, or final pull-request merge.

**Prepared history date:** 2026-09-15. **Intended version:** `1.1.0`. **Intended tag:** `v1.1.0`. **Release title:** `Cueson v1.1.0`. The intended public release is final, non-draft and non-prerelease.

## Actual source and accepted proof

The release decision package MUST be populated only after the human has merged the S029 preparation pull request and the resulting `main` push has completed fresh verification. Record the actual full lowercase 40-character post-merge `main` revision, prove a clean non-divergent default-branch checkout, and confirm that it is still the current default-branch revision when approval is requested and when each authorized publication action begins. Neither the preparation baseline, pull-request head, predicted squash result nor proof from another revision is the publication target.

Record successful exact-revision CI, CodeQL, Site and Release proof run identities and URLs, terminal review/security results, and the accepted candidate artifact ID, exact artifact name, source revision, byte digest and expiry. The artifact name MUST be `cueson-1.1.0-candidate-` followed by the actual accepted source revision. The Release proof MUST build one bundle with the existing pinned GoReleaser v2.18.1 and Syft v1.51.1, and its recorded Cueson toolchain/dependency identity MUST agree with the accepted source. A local or pull-request rebuild is review evidence and MUST NOT substitute for the accepted default-branch artifact.

The accepted structural `release-evidence.json` MUST bind version `1.1.0`, intended tag `v1.1.0`, the actual source revision, six archives, six SBOMs, six archive checksum entries and `published: false`. Retain the evidence bytes and their SHA-256 digest. Canonical, immutable repository, embedded, emitted and packaged schema bytes MUST agree with `release_schema_sha256` equal to `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7`. Record the accepted `LICENSE` and `NOTICE` byte lengths and lowercase SHA-256 digests and prove every archive contains those exact bytes.

Retain three native evidence records for `linux/amd64`, `windows/amd64` and `darwin/amd64`, each downloaded from the same accepted bundle without rebuilding. Each record MUST agree on source revision, schema/legal identity and all target archive/SBOM digests, declare native formats `srt`, `vtt`, `ass`, `ssa`, accept the two historical 1.0.0 inputs, and record all 32 old-consumer refusal paths plus all 32 identity probes using the verified published 1.0.0 consumer. Each record MUST prove fresh current identity/capability, native validation/inspection/render/restoration/conversion and refusal/output-preservation safety. Arm64 packages receive structural and build-information proof; no arm64 native execution is claimed.

## Exact public inventory

The populated decision package MUST record every following filename exactly once, with its kind, target where applicable, positive byte length and lowercase 64-character SHA-256 digest taken from the accepted bundle. Do not predict sizes or digests. Public files MUST be copied byte for byte from that bundle; no rebuild, re-archive, SBOM regeneration or checksum rewrite is allowed after acceptance.

| Filename | Kind | Target |
|---|---|---|
| `cueson_1.1.0_linux_amd64.tar.gz` | Archive | `linux/amd64` |
| `cueson_1.1.0_linux_amd64.tar.gz.sbom.json` | SPDX JSON SBOM | `linux/amd64` |
| `cueson_1.1.0_linux_arm64.tar.gz` | Archive | `linux/arm64` |
| `cueson_1.1.0_linux_arm64.tar.gz.sbom.json` | SPDX JSON SBOM | `linux/arm64` |
| `cueson_1.1.0_darwin_amd64.tar.gz` | Archive | `darwin/amd64` |
| `cueson_1.1.0_darwin_amd64.tar.gz.sbom.json` | SPDX JSON SBOM | `darwin/amd64` |
| `cueson_1.1.0_darwin_arm64.tar.gz` | Archive | `darwin/arm64` |
| `cueson_1.1.0_darwin_arm64.tar.gz.sbom.json` | SPDX JSON SBOM | `darwin/arm64` |
| `cueson_1.1.0_windows_amd64.zip` | Archive | `windows/amd64` |
| `cueson_1.1.0_windows_amd64.zip.sbom.json` | SPDX JSON SBOM | `windows/amd64` |
| `cueson_1.1.0_windows_arm64.zip` | Archive | `windows/arm64` |
| `cueson_1.1.0_windows_arm64.zip.sbom.json` | SPDX JSON SBOM | `windows/arm64` |
| `cueson_1.1.0_checksums.txt` | Checksum manifest | Six archives |

The checksum manifest MUST contain exactly one lowercase SHA-256 entry for each of the six archives, with no missing, duplicate, unexpected, SBOM or self-referential entry. Match those values to the decision package and archived bytes. Six target-specific SBOMs MUST retain the accepted target and source binding. Internal evidence, metadata, configuration, extracted files, signatures and attestations are outside this thirteen-file public inventory.

## Exact notes and dated history

The selected public body is [release-notes.md](release-notes.md), including its title, intentional blank lines and exact final tagged full-changelog link. Read its bytes from the accepted source, pass it through `go run ./scripts/github-format/main.go -stdin`, retain the exact UTF-8/no-BOM formatted file, and record its positive byte length and lowercase SHA-256 digest in the decision package. Review that exact file before authorization and use it without implicit banner removal, text substitution or generated additions. The candidate guidance in `docs/releases/v1.1.0.md` remains a separate truthful unpublished-state document until independently verified publication.

The selected source MUST include exactly one prepared `## [1.1.0] - 2026-09-15` changelog section, one fresh `[Unreleased]` section, complete release history, chronological dated Decisions and unchanged earlier release sections. The dated section prepares tagged history and does not assert that release publication has occurred. If the intended publication date requires a different prepared changelog date, update and review the metadata on `main`, refresh this decision contract and repeat source/proof/notes binding before seeking authority.

## Separate authority transitions

Preparation #74 can close after its reviewed metadata/contract outcome merges. Publication #66, hosting #67, epic #51 and the v1.1.0 milestone remain open. A green preparation pull request, its merge approval, this contract and any preflight observation grant no release or production authority.

| Transition | Required exact approval and observation |
|---|---|
| Create tag | Explicit approval to create unsigned annotated `v1.1.0` at the named accepted full revision; verify no conflicting local/remote tag and no existing release before acting. |
| Push tag | Explicit approval to push that exact tag and target; immediately read back the remote tag and its peeled revision. Creation approval alone does not authorize push. |
| Publish GitHub Release | Explicit approval for the exact tag, title, final/non-draft/non-prerelease state and reviewed formatted notes digest. |
| Publish release assets | Explicit approval for the exact thirteen-file inventory, byte lengths and digests from the accepted bundle. Release-body approval alone does not authorize asset publication. |
| Verify public publication | Independently read back release/tag/body/state and download all thirteen public files into a separate clean directory; verify all names, lengths, digests, checksum bijection, six structural/SBOM identities, packaged schema/legal bytes and compatible-host execution. |

An operator may authorize several named transitions in one explicit decision, but authority for one transition MUST NOT be inferred from another. No transition authorizes tag movement, asset replacement, signatures, attestations, production/public-schema changes, final pull-request merge or milestone closure. Record the exact approval text, approved source/body/inventory identity and actions before acting. Keep `published: false` as the preserved candidate evidence state; create a separate publication-verification record only after independent public checks pass.

## Refresh, expiry and partial-failure handling

Immediately before each approved action, recheck current default-branch revision, accepted successful checks, unresolved reviews/security findings, approved notes/date/schema/legal/tool/dependency identity, accepted asset bytes and artifact availability/expiry. A moved default branch, changed source, tools, dependencies, schema, legal files, notes, prepared date, assets or checks, new unresolved finding, unavailable artifact or expired artifact invalidates the pending decision. Obtain fresh accepted proof and a newly reviewed decision package; do not substitute cached or rebuilt bytes by inference, silently extend approval or reuse authority after a material state change.

If a tag, release or asset read-back differs, an upload is partial, independent verification fails or unexpected public files appear, stop all dependent publication and hosting actions, retain observations and accepted evidence, and report the exact partial state for operator recovery judgment. Do not move the tag, delete or overwrite public files, replace assets, promote a draft, retry a changed transaction or claim completion by inference. Close #66 and unblock #67 only after authorized publication and independent public verification satisfy every acceptance criterion; lifecycle and production actions retain their own authorization.
