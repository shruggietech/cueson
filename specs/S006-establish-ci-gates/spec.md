# Feature Specification: Establish CI and Cross-Platform Build Gates

**Feature Branch**: `S006-establish-ci-gates`

**Created**: 2026-09-09

**Status**: Draft

**Input**: User description: "Deliver issue #8 as S006: establish least-privilege continuous-integration, security-analysis, native-platform, and pure-Go cross-build gates; publish the work through the authorized autopilot pull-request and bounded two-round review protocol; exclude release publication, repository protection settings, pull-request review automation, and unimplemented codec claims."

## User Scenarios & Testing

### User Story 1 - Trust Every Proposed Change (Priority: P1)

As a maintainer, I can see a deterministic set of automated quality gates on every pull request and default-branch update so regressions in implemented Cueson behavior cannot pass silently.

**Why this priority**: Repeatable pull-request evidence is the prerequisite for every later repository rule, feature slice, and release decision.

**Independent Test**: Open a controlled pull request, observe the documented stable checks, prove that an intentional repository failure makes the workflow fail, remove the failure, and observe the same checks pass without manual background polling.

**Acceptance Scenarios**:

1. **Given** a pull request whose repository state satisfies every implemented contract, **When** automation completes, **Then** formatting, vetting, tests, race detection, static analysis, vulnerability analysis, schema integrity, fixture integrity, and repository-text checks report success under stable names.
2. **Given** a pull request containing a controlled failing condition, **When** automation runs, **Then** at least one required quality gate reports failure and the overall workflow cannot report success.
3. **Given** an obsolete workflow run for the same branch, **When** a newer commit starts, **Then** the obsolete run is cancelled without cancelling work for another branch.

---

### User Story 2 - Prove Portable Native Behavior and Builds (Priority: P2)

As a maintainer, I can distinguish native behavioral evidence from cross-compilation evidence across the supported operating systems and processor architectures.

**Why this priority**: Timestamp behavior is platform-specific, while release portability requires pure-Go builds for every supported target. Treating cross-compilation as native proof would overstate compatibility.

**Independent Test**: Inspect one completed automation run and verify native tests execute separately on Windows, macOS, and Linux while pure-Go binary builds cover every supported operating-system and architecture pair.

**Acceptance Scenarios**:

1. **Given** platform-selectable timestamp tests, **When** continuous integration runs, **Then** Windows, macOS, and Linux each execute their own native behavior rather than relying on cross-compilation alone.
2. **Given** the supported target matrix, **When** cross-build validation runs, **Then** Windows, macOS, and Linux each build for amd64 and arm64 with native dependencies disabled.
3. **Given** Linux cannot restore a birth timestamp, **When** Linux-native tests run, **Then** the documented unsupported result is tested without simulating support.

---

### User Story 3 - Receive Independent Security Evidence (Priority: P3)

As a maintainer or security reviewer, I receive independently reported code-scanning and dependency-vulnerability results without granting ordinary build jobs write authority.

**Why this priority**: Security evidence must be visible and trustworthy, but it must not unnecessarily expand workflow permissions or become coupled to the general build matrix.

**Independent Test**: Inspect workflow declarations and a completed pull-request run to verify ordinary jobs are read-only, code scanning is independently named, tool and action inputs are pinned, and no release or repository-setting mutation occurs.

**Acceptance Scenarios**:

1. **Given** a normal pull request, **When** security workflows run, **Then** static analysis, dependency vulnerability analysis, and independent code scanning produce visible results.
2. **Given** workflow permission declarations, **When** they are audited, **Then** ordinary jobs have read-only repository access and only the independent scanning workflow receives its narrowly required reporting permission.
3. **Given** a dependency or action reference used by automation, **When** the workflow is reviewed, **Then** its executable version is reproducibly pinned and its human-readable release identity remains discoverable.

### Edge Cases

- Documentation-only changes still receive the stable gate set so required-check presence does not depend on changed paths.
- A superseded run is cancelled only within the same workflow and change reference; default-branch work is not cancelled by an unrelated pull request.
- Fork-originated pull requests do not receive write-capable ordinary jobs or access to repository secrets.
- A security reporting permission unavailable to a fork does not widen permissions for general CI.
- Tool download or module-proxy failure produces a failed job rather than a false success or a silently skipped gate.
- Cross-built binaries remain temporary verification outputs and are never uploaded, released, or published by S006.
- Checks for native codecs, cross-format conversion, Codex-review automation, and release packaging remain absent until their owning issues implement those surfaces.

## Requirements

### Functional Requirements

- **FR-001**: S006 MUST add pull-request and default-branch automation with stable workflow and job names suitable for later repository protection.
- **FR-002**: General continuous-integration automation MUST declare read-only repository permissions, and any workflow requiring broader reporting authority MUST declare only the narrow permission needed at the narrowest practical scope.
- **FR-003**: Third-party and platform action executions MUST be pinned to immutable commit identifiers with a nearby human-readable release tag or version comment.
- **FR-004**: Repeated updates to the same change MUST cancel obsolete runs without cancelling runs for other changes or default-branch updates.
- **FR-005**: The quality gates MUST cover Go formatting, repository Markdown formatting, vetting, root-module tests, publication-formatter module tests, schema and software-version lockstep, schema key style, source-path prohibition, fixture manifest and byte integrity, golden comparisons, and repository text hygiene.
- **FR-006**: Race detection MUST run on one supported native target that provides reliable race-detector execution, without implying that one race run is native behavioral proof for every platform.
- **FR-007**: Static analysis and dependency vulnerability analysis MUST use explicit reproducible tool versions and fail closed when analysis cannot complete; reachable findings MUST be remediated or explicitly rejected through a separately approved security decision rather than suppressed for compatibility convenience.
- **FR-008**: Code scanning MUST run as an independently named security workflow and report results using only its required security-reporting permission.
- **FR-009**: Windows, macOS, and Linux MUST each run the root-module test suite natively so platform-selected timestamp behavior is executed on its owning operating system.
- **FR-010**: Pure-Go cross-build proof MUST cover Windows, macOS, and Linux on both amd64 and arm64 with `CGO_ENABLED=0` and MUST NOT publish the resulting binaries.
- **FR-011**: Repository text verification MUST reject invalid UTF-8, BOM-prefixed repository-authored text, mojibake indicators, line endings that violate `.gitattributes`, hard-wrapped Markdown prose where governed, and whitespace errors without inspecting byte-preserved fixture payloads as authored text.
- **FR-012**: The workflow MUST include a controlled, reviewable means to prove that an intentional failing repository state blocks success, and the official pull request MUST retain evidence of both the failing and corrected runs.
- **FR-013**: All verification MUST execute in the foreground of its job and propagate non-zero status without background polling shortcuts or success-on-error behavior.
- **FR-014**: The README MUST replace the temporary planned-CI badge with the stable live workflow badge in the same change that creates the workflow; final rendering on the default branch remains a post-merge verification item.
- **FR-015**: S006 MUST update architecture and unreleased changelog documentation to describe the implemented automation boundary and stable check contract.
- **FR-016**: S006 MUST NOT implement pull-request issue-link or Codex-review automation from issue #9, repository rules or security-setting mutations from issue #10, release packaging from issue #11, native subtitle codecs, cross-format conversion, tags, releases, or production-domain changes.
- **FR-017**: Checks named in the broader draft roadmap but owned by unimplemented codec, conversion, review-automation, or release surfaces MUST remain explicitly deferred rather than represented by empty or falsely passing jobs.

### Key Entities

- **Workflow**: An event-triggered automation definition with a stable name, declared permissions, concurrency boundary, and one or more visible jobs.
- **Quality Gate**: A stable, fail-closed verification result covering one coherent repository obligation.
- **Native Platform Run**: A test execution performed on the operating system whose behavior it claims to verify.
- **Cross-Build Target**: One supported operating-system and architecture pair built with native dependencies disabled as portability evidence only.
- **Pinned Execution Reference**: An immutable action or tool version paired with a readable release identity for auditability.
- **Controlled Failure Evidence**: A pull-request run intentionally made red through a bounded repository condition and followed by a corrected green run.

## Success Criteria

### Measurable Outcomes

- **SC-001**: Every pull request and default-branch update receives the same documented stable quality-gate names, with zero path-based omissions.
- **SC-002**: One automation run supplies native test evidence for all three supported operating systems and build evidence for all six supported operating-system and architecture pairs.
- **SC-003**: One controlled pull-request run fails because of the documented probe, and the next corrected run completes with every required S006 gate green.
- **SC-004**: One hundred percent of executable action references and analysis-tool versions are reproducibly pinned and auditable to a readable release identity.
- **SC-005**: Ordinary continuous-integration jobs request zero write permissions, and the independent scanning workflow requests no write permission beyond security-result publication.
- **SC-006**: A healthy pull-request run completes within 20 minutes under normal hosted-runner availability.
- **SC-007**: S006 produces zero published artifacts, releases, tags, repository-rule mutations, or production-domain changes.
- **SC-008**: The specification, workflows, README, architecture, changelog, issue acceptance criteria, and Project tracking agree on the S006 boundary with no unsupported capability claim.

## Assumptions

- GitHub-hosted Windows, macOS, and Linux runners remain available for this public repository.
- Go 1.25.0 is the minimum supported toolchain after vulnerability analysis proved that the prior Go 1.24 floor and its compatible `golang.org/x/text` release could not satisfy the security gate.
- Existing root-module tests are the authority for schema, version, source, fixture, and conformance behavior; S006 wires them into hosted automation rather than duplicating their semantics.
- The repository's current organization-level workflow defaults may remain broader until issue #10; each S006 workflow compensates with explicit least-privilege declarations.
- Race detection on one Linux amd64 hosted runner is sufficient for the S006 race gate, while native functional tests on all three operating systems supply platform behavior evidence.
- A dynamic CI badge may not produce a successful image until the workflow exists and runs on `main`; S006 verifies its target structure before merge and post-merge housekeeping verifies the live image.
