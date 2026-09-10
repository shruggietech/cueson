# Research: v0.0.0 Documentation and Milestone Verification

## Decision 1: Use a separate offline documentation verifier

**Decision**: Add `scripts/docs-verify` as a dependency-free Go 1.25 module that validates the canonical documentation inventory and local link targets without changing `scripts/github-format`.

**Rationale**: `github-format` owns encoding, line endings, Markdown source layout, mojibake detection, and GitHub-bound body formatting. Required-document and link resolution are different contract concerns. A small separate verifier can be tested independently, invoked locally and in CI, and remain outside the shipped product dependency graph.

**Alternatives considered**: Extending `github-format` would blur its publication-tool responsibility. A third-party Markdown link checker would add a supply-chain dependency, frequently fetch the network, and make deterministic local validation harder. A shell script would be less portable and harder to test consistently across Windows, macOS, and Linux.

## Decision 2: Audit maintained release documentation, not historical Spec Kit artifacts

**Decision**: The verifier scans root documentation (`README.md`, `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md`, and `AGENTS.md`), `.specify/memory/constitution.md`, and all Markdown and HTML beneath `docs/`. It enforces a twelve-document canonical required set. Historical slice artifacts beneath `specs/` are excluded from the release-document link contract.

**Rationale**: Root documents, the constitution, and `docs/` are maintained current authorities. Spec Kit directories preserve time-bound decision evidence whose old links and state descriptions are not current product documentation. S010's own artifacts remain governed by the repository formatter and explicit Spec Kit validation.

**Alternatives considered**: Scanning every Markdown file would convert immutable historical evidence into a continuously maintained documentation surface. Scanning only `docs/` would miss the README, contribution/security entry points, changelog, agent contract, and constitution.

## Decision 3: Resolve local links with GitHub-compatible heading identities

**Decision**: Accept inline Markdown links and images, reference definitions and uses, autolinks, and HTML `href`/`src` attributes. Resolve relative paths from the containing document, leading-slash paths from repository root, fragment-only references against the current document, URL-decoded paths and fragments, explicit HTML anchors, and duplicate heading suffixes. Ignore fenced code, inline code, HTML comments, CSS imports, and network resolution for external schemes. Reject repository escapes, path-case mismatches, symlinks, controls, backslashes, malformed escapes, unsupported schemes, and fragments on non-heading targets.

**Rationale**: These forms cover the repository's actual documentation patterns, including the HTML badge block and heading fragments. Offline resolution makes failures deterministic and keeps CI independent from remote availability or rate limiting.

**Alternatives considered**: Checking only inline Markdown would miss the README badge links. Fetching external links would confuse remote availability with repository correctness. Requiring a reduced authoring syntax would create needless documentation churn.

## Decision 4: Preserve truthful unpublished changelog state

**Decision**: Keep accumulated foundation changes under `[Unreleased]`, identify them as the v0.0.0 release candidate, and remove the premature empty `[0.0.0]` section and nonexistent release/tag links until an operator authorizes publication.

**Rationale**: The changelog must be complete before release, but Keep a Changelog entries become released only when an operator authorizes and performs the actual version publication. S010 has no tag or release authority.

**Alternatives considered**: Moving entries into a dated `0.0.0` section would falsely imply publication. Retaining an empty `0.0.0` release section and tag link would preserve draft prose at the cost of contradicting current GitHub state and Keep a Changelog semantics.

## Decision 5: Separate release process from release verification

**Decision**: Keep `docs/release-verification.md` as the executable non-publishing candidate proof and add `docs/release-process.md` as the broader gated sequence for version preparation, operator authorization, tag and GitHub Release publication, immutable schema handling, verification, and later production-domain work.

**Rationale**: Candidate construction is already implemented and safe on pull requests. Public release actions carry distinct authority and lifecycle rules. Separate documents prevent an ordinary verifier command from being mistaken for publication authorization.

**Alternatives considered**: Expanding the existing verifier guide would mix runnable unprivileged checks with protected human actions. Replacing it would discard the precise artifact proof completed by S009.

## Decision 6: Close issue #12 and epic #2 together only through merge

**Decision**: The official pull request will include `Closes #12` and `Closes #2`. During review both issues remain open, use Slice `S010`, and use governed Stage `PR review`; the default Status field remains empty. The v0.0.0 milestone remains open for a separately authorized release decision.

**Rationale**: Issue #12 is the epic's only open child. The same merged documentation/evidence change completes both outcomes, while pre-merge closure would claim default-branch evidence that does not yet exist. Closing the milestone is a release lifecycle action that S010 does not own.

**Alternatives considered**: Closing the epic in a later paperwork-only slice would duplicate the same evidence boundary. Closing either issue before merge would violate the repository's closure-evidence rule. Closing the milestone with the documentation PR would conflate release readiness with release completion.

## Decision 7: Reconcile durable issue evidence before declaring milestone readiness

**Decision**: Update stale acceptance and verification records on closed issues #6, #8, #9, #10, and #11 with their merged pull request, default-branch commit, hosted-check, and final-state evidence. Preserve issue history and change only statements that are now false or incomplete.

**Rationale**: Issue #12 requires every foundation issue to have closure evidence. Those five issue bodies still contain unchecked acceptance boxes, pre-merge instructions, or superseded activation claims even though their implementation is merged. GitHub-native issue state is the planning authority, so the evidence must be durable and readable there rather than inferred from a later documentation summary.

**Alternatives considered**: A new summary table alone would duplicate issue-owned facts and leave their authoritative bodies stale. Reopening completed issues would misrepresent implementation state. Editing unrelated historical narrative would create noise without improving acceptance evidence.

## Decision 8: Treat the working specification as a correctable draft

**Decision**: Mark only proven v0.0.0 foundation criteria complete and revise the draft changelog criterion so a candidate remains under `[Unreleased]` until publication. Leave all v1, release, immutable-schema, milestone-closure, and production requirements incomplete.

**Rationale**: The operator explicitly designated the specification as an evolving pre-release draft. Literal preservation of its premature `[0.0.0]` heading requirement would conflict with truthful publication state, the constitution, and current GitHub evidence.

**Alternatives considered**: Checking every release-related box would falsely claim protected actions. Leaving the entire foundation list unchecked would make the source-of-truth document deny completed implementation and hosted proof.
