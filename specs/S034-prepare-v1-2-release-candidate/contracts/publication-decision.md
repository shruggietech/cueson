# v1.2.0 publication decision contract

**State:** Preparation only. This document does not select or approve a release revision, create/push a tag, publish a Release/assets/schema, deploy production, close the milestone or grant a final PR merge. Candidate issue [#84](https://github.com/shruggietech/cueson/issues/84), public release [#85](https://github.com/shruggietech/cueson/issues/85) and production hosting [#86](https://github.com/shruggietech/cueson/issues/86) own separate outcomes.

**Prepared history date:** 2026-10-06. **Intended version:** `1.2.0`. **Intended tag:** `v1.2.0`. **Release title:** `Cueson v1.2.0`. The intended public release is final, non-draft and non-prerelease. Latest published downloads and production schemas remain v1.1.0 during preparation.

## Actual source and accepted proof

Populate a concrete decision only after the owner authorizes the specific reviewed preparation PR merge and fresh actual-main checks finish. Record the actual full lowercase 40-character post-merge main revision, a clean non-divergent checkout, current default-branch identity, successful exact-revision CI/CodeQL/Site/Release proof run IDs/URLs and terminal review/security results. A PR head, preparation baseline, predicted squash result or proof from another revision cannot be the public release target.

The accepted artifact is exactly `cueson-1.2.0-candidate-` followed by that source revision. Record artifact ID/name, source revision, byte digest, expiry and retained evidence SHA-256. It must come from fresh default-branch proof using pinned GoReleaser v2.18.1 and Syft v1.51.1, with accepted toolchain/dependency identity and `published: false`. A local or PR rebuild is review evidence only.

Canonical, immutable repository, embedded, emitted and packaged schema bytes must match SHA-256 `f2661a3d52effbab4a82a4d47197b5c7fae58496dc30a397ea3f2f668358b654` (191170 bytes). Record accepted LICENSE/NOTICE lengths and lowercase SHA-256 values and require identical copies in every archive. Historical released schemas and publication contracts remain immutable.

The structural release-evidence.json must satisfy [the exact candidate evidence contract](release-evidence-contract.schema.json), bind version 1.2.0, tag v1.2.0 and the full source revision and prove six archives, six target-bound SPDX JSON SBOMs and six archive checksum entries. Retain its exact bytes/digest.

Retain three native records for linux/amd64, windows/amd64 and darwin/amd64 downloaded from the same accepted bundle without rebuilding. Each agrees on source/schema/legal/target inventory and proves four native formats, twelve conversions, six historical inputs (1.0.0 SubRip/WebVTT and 1.1.0 SubRip/WebVTT/ASS/SSA), current consumer storage/validation/privacy/restore/native losses, and strict output-preservation safety. Historical conversion checks include twelve strict refusals against occupied and absent destinations for the six fixtures' documented native losses. Consumer timing checks independently prove four invalid cue intervals, four timed assignments contained by their cues but outside declared media, and eight guarded media-bound destination refusals. Missing optional annotations retain existing behavior; unavailable duration checks are not passed checks.

Published old-consumer records must independently identify authenticated public 1.0.0 and 1.1.0 binaries. Preserve 1.0.0's two positive controls, 32 actual refusal paths and 32 complete-document identity probes. Bind the 1.1.0 archive to the frozen [published consumer contract](published-v1.1.0-consumer.json), verify all four historical 1.1.0 positive controls, and require 64 actual refusals plus 64 complete-document identity probes covering current base and annotated documents on every applicable consuming path. Separately require 32 selector-only probes for 1.0.0 and 64 for 1.1.0, with exactly `$schema` and `schema_version`; these isolate version refusal from unrelated missing-payload validation. Occupied/absent destinations and input files stay unchanged. Native consumer annotation counters and invalid/conflict paths must match the exact evidence schema. An unavailable public download is incomplete proof. Arm64 packages receive structural/build proof; no arm64 native execution is claimed.

## Exact public inventory

Every filename below appears exactly once with kind, target, positive byte length and lowercase 64-character SHA-256 taken from accepted bytes. Do not predict size/digest. Copy files byte for byte; do not rebuild, re-archive, regenerate SBOMs or rewrite checksums after acceptance.

| Filename | Kind | Target |
|---|---|---|
| `cueson_1.2.0_linux_amd64.tar.gz` | Archive | linux/amd64 |
| `cueson_1.2.0_linux_amd64.tar.gz.sbom.json` | SPDX JSON SBOM | linux/amd64 |
| `cueson_1.2.0_linux_arm64.tar.gz` | Archive | linux/arm64 |
| `cueson_1.2.0_linux_arm64.tar.gz.sbom.json` | SPDX JSON SBOM | linux/arm64 |
| `cueson_1.2.0_darwin_amd64.tar.gz` | Archive | darwin/amd64 |
| `cueson_1.2.0_darwin_amd64.tar.gz.sbom.json` | SPDX JSON SBOM | darwin/amd64 |
| `cueson_1.2.0_darwin_arm64.tar.gz` | Archive | darwin/arm64 |
| `cueson_1.2.0_darwin_arm64.tar.gz.sbom.json` | SPDX JSON SBOM | darwin/arm64 |
| `cueson_1.2.0_windows_amd64.zip` | Archive | windows/amd64 |
| `cueson_1.2.0_windows_amd64.zip.sbom.json` | SPDX JSON SBOM | windows/amd64 |
| `cueson_1.2.0_windows_arm64.zip` | Archive | windows/arm64 |
| `cueson_1.2.0_windows_arm64.zip.sbom.json` | SPDX JSON SBOM | windows/arm64 |
| `cueson_1.2.0_checksums.txt` | Checksum manifest | Six archives |

The manifest contains one lowercase SHA-256 entry per archive, with no missing/duplicate/unexpected/SBOM/self entry. Internal evidence/metadata/configuration/extracted files/signatures/attestations are outside the thirteen-file public inventory.

## Exact notes and dated history

Use [release-notes.md](release-notes.md) as the exact public body, including title, intentional blank lines and final tagged full-changelog link. Read the accepted source bytes, pass through `go run ./scripts/github-format/main.go -stdin`, retain UTF-8/no-BOM output and bind positive length/digest before authorization. Do not remove a banner implicitly or add generated text. docs/releases/v1.2.0.md is separate candidate guidance until public verification.

The accepted source contains exactly one `## [1.2.0] - 2026-10-06` section, one fresh Unreleased, chronological Decisions and unchanged earlier release history. Prepared date does not assert publication. If a new date is needed, review updated metadata on main and repeat source/proof/notes binding before authority.

## Separate authority transitions

| Transition | Required approval and observation |
|---|---|
| Create tag | Explicit approval for unsigned annotated v1.2.0 at the named accepted full revision; check no conflicting local/remote tag or existing release. |
| Push tag | Explicit approval for that exact tag/target; immediately read back remote tag and peeled revision. |
| Publish GitHub Release | Explicit approval for tag, title,final/non-draft/non-prerelease state and exact formatted notes digest. |
| Publish release assets | Explicit approval for thirteen names, positive sizes and SHA-256 values from the same accepted bundle. |
| Independently verify | Read back tag/release/body/state and download all thirteen public files into a separate clean directory; verify inventory/digests/checksum bijection, six structural/SBOM/schema/legal identities and compatible public-package execution. |

The operator can authorize several named transitions together. One transition does not authorize another, tag movement, asset replacement, signatures/attestations, production changes, final merge or milestone closure. Record exact approval text and approved bindings. Preserve candidate evidence's published: false and create separate independent public verification evidence.

## Refresh, expiry and partial failure

Immediately before each approved action recheck current main, accepted checks/reviews, source/tools/dependencies/schema/legal/notes/date identity, asset bytes and artifact availability/expiry. Any changed state, new unresolved finding, unavailable/expired artifact invalidates the pending decision. Obtain fresh accepted proof and a newly reviewed decision package; do not silently substitute cached/rebuilt bytes or extend approval.

If read-back differs, upload is partial, verification fails or unexpected files appear, stop all dependent publication/hosting, preserve evidence and report exact partial state for operator recovery judgment. Do not move/delete/replace public state or retry a changed transaction by inference. Close #85 only after authorized independently verified publication; #86 then owns separately reviewed and authorized hosting. Milestone closure follows all committed outcomes.
