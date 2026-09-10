# Feature Specification: v0.0.0 Documentation and Milestone Verification

**Feature Branch**: `S010-complete-v0-docs`

**Created**: 2026-09-09

**Status**: Draft

**Input**: User description: "Use Spec Kit and the repository autopilot protocol to complete the v0.0.0 documentation and milestone-verification boundary, automatically publish the official pull request, complete no more than two Codex review rounds, and stop for the operator's final review and merge ritual."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Trust the v0.0.0 Documentation (Priority: P1)

A maintainer can compare every release-specific document with the implemented executable, schema, source-restoration behavior, automation, repository controls, and non-publishing release proof and find no stale or unsupported claim.

**Why this priority**: The foundation cannot be considered release-ready while its documentation contradicts the behavior and controls already delivered.

**Independent Test**: Audit every required repository document against the current executable, canonical schema, architecture of record, closed issue evidence, and hosted checks, then confirm that every current-state claim is supported and every future capability is labeled as planned.

**Acceptance Scenarios**:

1. **Given** the merged v0.0.0 foundation, **When** a maintainer follows current-behavior links and examples, **Then** every referenced document and command exists and agrees with implemented behavior.
2. **Given** a statement about SRT, WebVTT, rendering, conversion, or publication, **When** it is compared with the current capability boundary, **Then** it distinguishes envelope-only support from future native support and does not imply a public release.
3. **Given** historical delivery records, **When** they are audited after S009, **Then** transient pre-merge statements are either updated to durable final evidence or explicitly identified as historical context.

---

### User Story 2 - Navigate the Required Documentation Set (Priority: P1)

A contributor can locate dedicated SRT, WebVTT, and release-process documentation from the repository entry points and understand both the current v0.0.0 boundary and the work planned for v1.0.0.

**Why this priority**: Missing format and release-process pages prevent the documented repository contract from being complete even when the underlying implementation is sound.

**Independent Test**: Start from the README and documentation entry points, follow every local link, and verify that the required format and release pages exist, have distinct authority boundaries, and identify current versus planned behavior without ambiguity.

**Acceptance Scenarios**:

1. **Given** a contributor investigating SubRip, **When** they open the dedicated SRT page, **Then** they can identify the current envelope-only capability and the separately planned v1 grammar, fidelity, diagnostic, and rendering obligations.
2. **Given** a contributor investigating WebVTT, **When** they open the dedicated WebVTT page, **Then** they can identify the current envelope-only capability and the separately planned v1 structure, rolling-caption, fidelity, diagnostic, and rendering obligations.
3. **Given** an operator preparing a future release, **When** they open the release-process page, **Then** they can distinguish candidate verification from separately authorized tag, GitHub Release, immutable schema-copy, and production publication actions.
4. **Given** any repository-authored Markdown link in the audited documentation scope, **When** documentation verification runs, **Then** every local target and anchor resolves or the verification fails visibly.

---

### User Story 3 - Review the Completed Foundation Milestone (Priority: P2)

An operator can review one coherent evidence record showing that all atomic v0.0.0 foundation outcomes are complete and that the milestone is ready for a separate release decision without S010 performing that release.

**Why this priority**: Closing the foundation epic requires complete and truthful evidence, but release authority remains a separate human decision.

**Independent Test**: Review issues #1 through #12, the parent-child relationship, Project membership and fields, the milestone inventory, version lockstep, current-head checks, and release-proof evidence, then confirm that the S010 pull request closes issues #12 and #2 only on merge.

**Acceptance Scenarios**:

1. **Given** issues #1 through #11 are closed and issue #12 is the final open child, **When** the S010 evidence is prepared, **Then** every child acceptance outcome is traceable and the pull request carries complete closing references for issues #12 and #2.
2. **Given** the v0.0.0 milestone and Delivery Project, **When** their state is audited, **Then** every in-scope issue appears exactly once, completed issues use Stage `Done`, S010 issues use Stage `PR review` during review, Slice is `S010`, and the default Status field is unused.
3. **Given** the S010 pull request is ready for the operator, **When** final evidence is read back, **Then** all CI and security checks are green, all authorized review rounds are satisfied, and no tag, release, schema publication, milestone closure, or production mutation has occurred.

### Edge Cases

- A local Markdown link can point to an existing file while its fragment identifies no heading.
- Generated badges can be technically reachable while communicating a stale or false project state.
- A roadmap section can use mandatory language for v1 behavior and accidentally read as a v0.0.0 capability claim when the time boundary is omitted.
- Historical pull-request evidence can become false when copied into an enduring current-state section after merge.
- The current release-verification guide and the new release-process guide can duplicate authority unless candidate proof and publication procedure are separated explicitly.
- The changelog can imply that v0.0.0 was released merely because a version heading exists.
- Closing the documentation issue before merge would make the epic appear complete without the official changes on the default branch.
- GitHub's default Project Status automation can repopulate a field the project contract deliberately leaves unused.
- A documentation verifier can mistake external URLs, URL-encoded fragments, images, code examples, or generated references for broken local links.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: S010 MUST audit every repository document required by the project specification against the merged v0.0.0 implementation and recorded GitHub evidence.
- **FR-002**: The repository MUST contain dedicated SRT and WebVTT format pages and a dedicated release-process page at the documented canonical paths.
- **FR-003**: The SRT and WebVTT pages MUST describe v0.0.0 as envelope-only and MUST separate current source preservation and restoration from planned v1 native ingest, semantic extraction, model-driven rendering, and conversion.
- **FR-004**: The format pages MUST summarize their planned v1 grammar and fidelity obligations without claiming that fixtures, parsers, renderers, or conversion behavior already exist.
- **FR-005**: The release-process page MUST distinguish the existing non-publishing candidate proof from future tag creation, immutable release-schema copying, GitHub Release publication, release-asset verification, milestone closure, and production-domain publication.
- **FR-006**: The release-process page MUST state which release actions require specific operator authorization and MUST NOT introduce an executable publishing path in S010.
- **FR-007**: README badges, status text, examples, documentation links, release statements, and support claims MUST match the current default-branch behavior and public GitHub state.
- **FR-008**: Architecture, schema, CLI, project-management, repository-control, and release-verification documents MUST contain durable current-state language and MUST remove or clearly contextualize stale pre-merge statements.
- **FR-009**: The working project specification MUST record completion evidence for every satisfied v0.0.0 foundation gate while leaving v1 and public-release gates unchecked.
- **FR-010**: `CHANGELOG.md` MUST account for all v0.0.0 foundation work, retain a truthful unpublished state, and distinguish accumulated candidate changes from an actual released version.
- **FR-011**: Documentation verification MUST fail visibly for a missing required document, unresolved local file link, unresolved local heading fragment, malformed local target, or forbidden hard-wrapped or corrupted repository text in the audited scope.
- **FR-012**: Documentation verification MUST accept valid relative links, repository-root links, URL-encoded paths and fragments, images, external URLs, and intentional links to GitHub resources without network-dependent success claims.
- **FR-013**: Documentation and verification changes MUST preserve UTF-8 without BOM, repository line-ending policy, and the one-line-per-paragraph Markdown rule.
- **FR-014**: S010 MUST verify executable and schema version lockstep, the complete existing repository quality gate, and the non-publishing release contract without changing shipped product behavior.
- **FR-015**: S010 MUST audit issues #1 through #12, the v0.0.0 milestone, native parent and sub-issue relationships, dependencies, Project membership, Stage, Slice, and unused Status fields.
- **FR-016**: The official pull request MUST contain complete closing references for issues #12 and #2 so the final child and its parent epic close together only when the operator merges the reviewed change.
- **FR-017**: The final milestone evidence MUST identify any state that cannot become final until merge and MUST avoid marking a release or milestone as published or closed prematurely.
- **FR-018**: S010 MUST complete no more than two Codex review rounds, address every review and security finding, leave all review threads resolved, and stop for the operator's final review and merge ritual.
- **FR-019**: S010 MUST NOT merge or auto-merge its pull request, create or move a tag, publish a release or schema, create an immutable release-schema copy, close the milestone, or mutate production configuration.

### Key Entities

- **Documentation Contract**: One required repository document with a canonical path, authority boundary, current-state claims, future-state claims, and verifiable local links.
- **Capability Claim**: A statement about behavior or support classified as implemented, envelope-only, planned, release-gated, or post-release.
- **Documentation Link**: A local file or heading reference whose target must resolve deterministically, or an external reference whose syntax is validated without network dependence.
- **Foundation Gate**: One v0.0.0 completion criterion paired with repository or GitHub evidence.
- **Milestone Evidence Record**: The combined issue, Project, version, automation, release-proof, and publication-boundary evidence used for operator review.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of the required repository documentation paths exist, including the dedicated SRT, WebVTT, and release-process pages.
- **SC-002**: Deterministic verification resolves 100% of local files and heading fragments in the audited Markdown scope with zero network dependency.
- **SC-003**: An audit of current-state documentation finds zero unsupported claims of native SRT or WebVTT ingest, model-driven rendering, conversion, public release, schema publication, or production activation.
- **SC-004**: Every v0.0.0 foundation completion criterion has explicit evidence, while 100% of v1, public-release, milestone-closure, and production gates remain visibly incomplete where they depend on later authorization or implementation.
- **SC-005**: All eleven child issues of epic #2 have independently traceable acceptance evidence, and the S010 pull request is configured to close the final child and parent epic together on merge.
- **SC-006**: 100% of in-scope issues appear exactly once in `cueson Delivery`, use only the governed Stage and Slice planning fields, and leave the default Status field empty after reconciliation.
- **SC-007**: The full local and hosted quality gates complete successfully with zero unresolved review threads and no more than two Codex review rounds.
- **SC-008**: S010 creates zero tags, GitHub Releases, immutable release-schema copies, published schemas, milestone closures, or production-domain mutations.

## Assumptions

- Issues #12 and #2 form one coherent S010 closeout because the documentation issue is the epic's only remaining open child and both share one milestone-evidence boundary.
- Changelog content remains in an unpublished candidate state until a separately authorized release action supplies a tag and release date.
- External links are not fetched during deterministic documentation verification; only their syntax and classification are in scope.
- The existing release-verification document remains the authority for running the non-publishing artifact proof, while the new release-process document owns the broader human-authorized release sequence.
- Documentation verification may use repository-owned standalone tooling outside the shipped Cueson product module.
- Current issue, Project, milestone, check, and pull-request state must be read directly from GitHub before publication because it can change independently of the working tree.
