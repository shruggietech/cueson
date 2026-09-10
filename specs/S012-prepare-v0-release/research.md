# Research: Prepare the v0.0.0 Release

## Decision 1: Bind the release candidate after squash merge

**Decision**: Treat pull-request proof as review evidence and automatically rerun the identical non-publishing proof on pushes to `main`. After the S012 squash merge, the exact resulting default-branch commit is the proposed v0.0.0 publication target only when its own proof passes.

**Rationale**: A pull-request head is not the commit that a squash merge creates. Recording that head as the final release target would make the release dossier stale by construction, while waiting for the `main` push provides an exact immutable revision without guessing.

**Alternatives considered**: Recording the current pre-S012 `main` commit omits the release preparation. Recording the pull-request head does not survive squash merge. Merging with a non-squash strategy would violate repository policy. Manually editing a committed SHA after merge would require another commit and move the target again.

## Decision 2: Package from the versioned release schema

**Decision**: Create `schema/releases/v0.0.0/cueson.schema.json` as an exact copy of `internal/schema/cueson.schema.json`, make GoReleaser package the versioned copy, and have the verifier reject any byte difference between the two before inspecting artifacts.

**Rationale**: The release artifact should consume the same immutable repository path that reviewers approve for that version. Comparing it to the canonical embedded source preserves the existing authoring and `go:embed` arrangement while proving all copies converge.

**Alternatives considered**: Continuing to package directly from `internal/schema` would leave the release path unused and would not prove admission. Moving the canonical schema would create unnecessary product changes. Semantic-only comparison would permit byte drift contrary to the immutable contract.

## Decision 3: Stage the final changelog with a fresh Unreleased section

**Decision**: Move all current entries into `## [0.0.0] - 2026-09-10`, create an empty `## [Unreleased]` section above it, and add Keep a Changelog comparison links for the future range and the initial tag.

**Rationale**: The reviewed tag commit needs a complete permanent version history, while the repository must retain a place for later work. The release remains publicly unreleased until tag and GitHub Release publication are separately authorized.

**Alternatives considered**: Leaving all history under `[Unreleased]` would make the intended tag incomplete. Removing `[Unreleased]` would force a follow-up edit before any later change. Claiming publication in the changelog prose would cross the state boundary.

## Decision 4: Keep release notes concise and publication-ready

**Decision**: Store highlights in `docs/releases/v0.0.0.md`, state the current envelope-only capability limit, and require the exact final tagged changelog link.

**Rationale**: The notes can be reviewed as repository content and later published byte for byte after authorization. Explicit limitations prevent the first public artifact from implying native SRT or WebVTT codecs.

**Alternatives considered**: Copying the entire changelog would violate the repository's highlights-only rule. Generating notes at publication time would bypass review. Omitting limitations would make capability claims ambiguous.

## Decision 5: Enrich deterministic evidence with schema identity

**Decision**: Add the lowercase SHA-256 of the versioned release schema to `release-evidence.json` and retain the existing exact source revision, target inventory, archive and SBOM digests, counts, host execution, and `published: false` fields.

**Rationale**: The operator's decision package must connect the permanent schema record to the ephemeral candidate artifacts. One deterministic digest closes that gap without claiming byte-reproducible SBOMs.

**Alternatives considered**: A separate manifest would duplicate the same source and target catalog. Committing generated artifact hashes is impossible before the final commit and would create a self-changing release target. Hashing only the canonical authoring path would not prove the admitted release copy.

## Decision 6: Retain the non-publishing workflow boundary

**Decision**: Add only a `main` push trigger. Keep `contents: read`, credential-free checkout, no secrets, no identity token, no tag trigger, no signing, no deployment, `release.disable: true`, snapshot mode, and `--skip=publish`.

**Rationale**: Default-branch evidence does not require publication authority. The same verifier and build configuration can prove the merge commit with the least possible workflow change.

**Alternatives considered**: A tag-triggered workflow would mix preparation with publication. Write permission or a GitHub token is unnecessary. A second workflow would duplicate the release contract and increase drift risk.

## Decision 7: Keep milestone closure after publication verification

**Decision**: Add the S012 issue to v0.0.0, close it through the merged pull request, and leave the milestone open until an authorized tag and GitHub Release pass post-publication read-back.

**Rationale**: Repository readiness and public release completion are different lifecycle facts. The milestone currently has no open issue but remains intentionally open; S012 adds the remaining preparation outcome without redefining completion.

**Alternatives considered**: Closing the milestone when the preparation PR merges would claim release completion too early. Leaving S012 outside the milestone would make the release commitment incomplete.
