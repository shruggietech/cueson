# Research: Prepare the v1.0.0 Release Candidate

## Decision 1: Promote identity and support maturity in one transaction

**Decision**: Change the executable, canonical schema, schema constants, representative documents, governed fixtures, model validation, producers, tests, and current-state documentation together from `0.1.0` and `experimental` to `1.0.0` and `stable` for SubRip and WebVTT.

**Rationale**: The architecture deliberately held the complete v1-bound behavior at a development identity until S018 froze its contracts. Issue #37 owns the coordinated stable transition; changing only version strings or only capability status would create a false and internally inconsistent release record.

**Alternatives considered**: Keep `experimental` in the v1 schema, which contradicts the recorded release gate; change schema status without model and producer behavior, which breaks validation; retain `0.1.0` while calling artifacts v1, which violates software/schema lockstep.

## Decision 2: Finalize canonical bytes before admitting the immutable copy

**Decision**: Complete every `1.0.0` identity, annotation, example, and stable-support edit in `internal/schema/cueson.schema.json`, verify the resulting normative and full digests, then create `schema/releases/v1.0.0/cueson.schema.json` as an exact byte copy and pin that identity in tests.

**Rationale**: Released schema copies are immutable. Copying early and editing two files invites drift, while generating the copy from a finalized reviewed canonical source gives one clear admission event and preserves the v0.0.0 bytes.

**Alternatives considered**: Edit both copies in parallel, which creates duplicate authorities; generate the immutable file only during packaging, which omits a reviewed repository contract; overwrite v0.0.0, which violates immutability.

## Decision 3: Package the immutable schema and run candidate-mode proof

**Decision**: Fix the GoReleaser snapshot identity to `1.0.0`, package `schema/releases/v1.0.0/cueson.schema.json`, retain `release.disable: true` and `--snapshot --skip=publish`, and run `scripts/release-verify` without `-development` in local and hosted proof.

**Rationale**: The existing verifier's default path already requires canonical/versioned byte identity and performs the six-target artifact, build-information, checksum, software-bill-of-material, embedded-schema, and host execution checks needed for a candidate. Development mode intentionally skips immutable-schema admission and therefore cannot prove v1 readiness.

**Alternatives considered**: Add a publication-capable GoReleaser path, which crosses an excluded authority boundary; keep `-development`, which fails issue #37; build ad hoc archives outside the reviewed pipeline, which duplicates release logic.

## Decision 4: Strengthen exact schema and legal-file verification

**Decision**: Require the canonical and versioned schema `$id`, root instance `$schema` constant, and `schema_version` constant to equal the exact expected identity and version. Compare packaged `LICENSE` and `NOTICE` bytes with their repository sources, require the portable approved archive mode `0644`, and record their digests in evidence.

**Rationale**: The current suffix-only `$id` check could accept an incorrect origin, and it does not independently check the instance `$schema` constant. Presence-only legal-file checks can accept stale or modified content. Exact identity and bytes are release invariants. Repository checkout modes are not a portable authority on Windows, so comparing archive modes to the working tree would produce platform-dependent proof. The governed archive rule is therefore exact repository bytes plus fixed non-executable mode `0644`.

**Alternatives considered**: Retain suffix-only and presence-only checks, which can accept counterfeit authority or legal drift; rely only on embedded byte comparison, which proves sameness but not correctness of the common schema bytes.

## Decision 5: Execute one accepted bundle on three native operating systems

**Decision**: Build and structurally verify one six-target bundle on Ubuntu, upload that exact accepted bundle, and run the same verifier with host execution against the matching amd64 artifact on hosted Linux, Windows, and macOS runners. Each smoke run writes its own evidence path and does not rebuild the candidate.

**Rationale**: Source tests on each operating system do not prove the packaged executables. Fanning out one accepted bundle proves native startup, version output, schema-version output, and emitted schema without creating divergent artifacts.

**Alternatives considered**: Rebuild independently on each runner, which yields multiple candidate sets; claim cross-compilation as execution proof, which is false; add arm64 execution without available governed runners, which overstates evidence.

## Decision 6: Preserve historical records while migrating current-state surfaces

**Decision**: Leave the v0.0.0 schema, release notes, tag links, published evidence, S012/S013 artifacts, and historical changelog statements intact. Update only current-source identities, current examples and governed fixtures, maintained current-state documentation, and current tests. Add a new S019 evidence contract rather than modifying the completed S012 contract.

**Rationale**: Historical release records are evidence. S019 must distinguish stable v1 candidate state from published v0.0.0 facts and from the still-unpublished v1 public release.

**Alternatives considered**: Global version-string replacement, which corrupts history; rewrite S012 evidence for v1, which destroys traceability; call v1 published during candidate preparation, which is false.

## Decision 7: Separate detailed history from candidate highlights

**Decision**: Move the accumulated `[Unreleased]` implementation and decision record into `## [1.0.0] - 2026-09-11`, retain a fresh empty `[Unreleased]` section, add valid comparison links, and create concise `docs/releases/v1.0.0.md` highlights ending with the exact tagged changelog link.

**Rationale**: The changelog is the complete release history, while GitHub release notes must be shorter and useful at publication time. Reviewing both now prevents publication from inventing or editing release meaning after the candidate is accepted.

**Alternatives considered**: Copy the entire changelog into release notes, which violates repository policy; defer notes until publication, which removes them from candidate review; omit the fresh `[Unreleased]` section, which breaks the continuing history model.

## Decision 8: Correct the pre-release working specification's circular public-schema gate

**Decision**: Change the working project specification's v1 completion gate so candidate readiness requires canonical, immutable repository, embedded, emitted, and packaged schema byte identity. Keep public `cueson.io` schema serving in the separately authorized post-v1 production gate, and clarify that release verification completion means issue #37 candidate proof rather than issue #38 publication.

**Rationale**: Requiring a public schema endpoint before v1 publication contradicts the constitution, release process, issue #37 scope, and the same working document's post-v1 hosting boundary. The operator has directed agents to challenge the draft where it slows or contradicts approved setup and delivery.

**Alternatives considered**: Publish production schema hosting in S019, which exceeds authority; leave the circular gate, which makes compliant v1 release impossible; silently ignore the working specification, which hides the contradiction.

## Decision 9: Bind final publication eligibility only after squash merge

**Decision**: Treat the S019 pull-request artifact proof as review evidence. After an operator-authorized squash merge, require the existing `main` push trigger to rebuild and verify the exact resulting commit before issue #38 can publish a tag or GitHub Release.

**Rationale**: The pull-request head is not the squash-merge commit. Predicting it would bind evidence to the wrong source revision.

**Alternatives considered**: Tag the pull-request head, which bypasses reviewed default-branch state; record a future commit placeholder, which is not auditable; publish during S019, which exceeds authorization.

## Decision 10: Claim semantic, not byte-for-byte, software-bill-of-material reproducibility

**Decision**: Continue verifying target association, version, source revision, catalog shape, and digest recording for every software bill of materials without claiming identical upstream document bytes across independent rebuilds. Deterministic claims apply to repository inputs, names, archive contents, source-derived timestamps, checksums, and verifier behavior.

**Rationale**: Syft includes variable timestamps and document identifiers. The candidate can still prove stable artifact and semantic source binding without making an untrue stronger claim.

**Alternatives considered**: Normalize or rewrite Syft output, which would create an unreviewed derived format; claim byte reproducibility despite known variability, which violates truthful delivery.
