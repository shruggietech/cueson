# Feature Specification: Prepare reviewed v1.1.0 publication

**Feature Branch**: `codex/S029-publish-verify-v1-1`

**Created**: 2026-09-15

**Status**: Specified

**Input**: Operator S029 autopilot, explicit branch push/official PR authorization, after merged S028 main `6a55b48c654e14f5baaf89a9add51fa4f5a4d04a`. Preparation issue #74 is a sub-issue of publication #66. No tag/release/production or final-merge authority is supplied.

## Clarifications

### Session 2026-09-15

- Q: Can the preparation PR close publication #66? A: No. Close independently testable #74 and reference #66, which stays open until authorized public release verification. The finalized changelog must reach reviewed main before selecting the tag target.
- Q: Can a future squash-merge revision or bundle be predicted? A: No. Preparation/PR proof is review evidence. Only fresh accepted proof of the actual post-merge main can populate the exact publication decision package.
- Q: Which release notes are approved for public use? A: Preserve prospective candidate guidance and separately freeze concise publication notes without an unpublished-candidate banner or implicit removal transformation. The future decision package binds their exact formatted body digest.
- Q: Does a dated changelog assert availability? A: No. It is prepared tagged metadata; current publication state remains explicitly unpublished, public downloads remain v1.0.0, and public schema hosting remains #67.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Review complete tagged history (Priority: P1)

The operator and downstream consumers can review one complete prepared v1.1.0 history and accurate concise release notes before any release transaction.

**Why this priority**: The final tag must include the reviewed changelog and notes without rewriting historical releases or claiming early public availability.

**Independent Test**: Audit prepared history and notes, valid dates/order/categories, exact final changelog links, bounded capability disclosures and current unpublished state.

**Acceptance Scenarios**:

1. **Given** accumulated unreleased history, **When** preparation is reviewed, **Then** one dated v1.1.0 section contains the complete history, one fresh Unreleased section remains, duplicate categories are consolidated and dated decisions stay chronological.
2. **Given** prepared public highlights and candidate guidance, **When** compared, **Then** capabilities/compatibility/limitations agree, the public body is exact and substantially shorter than the changelog, and candidate guidance still states unpublished availability.
3. **Given** older released history and schema/source assets, **When** preparation is applied, **Then** their bytes remain unchanged.

### User Story 2 - Make an exact release decision (Priority: P1)

The operator receives a bounded contract for reviewing the actual post-merge source and accepted assets before authorizing publication.

**Why this priority**: A predicted merge target, stale bundle or implicit notes transformation cannot safely define a public release.

**Independent Test**: Inspect the normative decision contract and repository checks for source/proof/asset/notes/authority bindings and material-change refresh rules.

**Acceptance Scenarios**:

1. **Given** a reviewed preparation PR, **When** it merges, **Then** the future decision package selects the actual verified main revision and requires fresh green default-branch proof instead of a PR head or superseded baseline.
2. **Given** an accepted bundle, **When** the decision package is populated, **Then** it binds six archives, six target-specific SPDX SBOMs and one exact checksum manifest, schema/legal identity, every name/size/digest, and three same-bundle native evidence records.
3. **Given** a requested protected action, **When** its authority is evaluated, **Then** exact tag creation/push and release/asset publication each need explicit approval; changes to source, tools, assets, schema, notes or checks require a refreshed decision.

### User Story 3 - Receive one fully reviewed preparation PR (Priority: P2)

The operator receives an official verified PR closing preparation #74 while publication and hosting outcomes retain their independent gates.

**Why this priority**: Completion must reflect executed preparation evidence without prematurely closing public-release acceptance.

**Independent Test**: Read back issue/native Project relationships, publication bodies, exact-head check results and all terminal review states.

**Acceptance Scenarios**:

1. **Given** the authorized slice branch, **When** the official PR is published, **Then** it closes #74, references #66, renders correctly and has green exact-head checks with every actionable review handled using at most one requested second Codex round.
2. **Given** preparation completion, **When** lifecycle is reconciled, **Then** #66/#67 and epic/milestone remain open, every issue has unique Project membership, and the default Status field remains unused.

### Edge Cases

- Duplicate, missing, malformed or reordered release headings must not pass metadata checks.
- Candidate notes containing an unpublished banner are not silently repurposed as the final public body.
- A changed default branch, expired artifact, stale notes digest or partial public release blocks publication and dependent hosting.
- Tag/release absence is a preflight observation, not authority to create either.
- Arm64 build/structural acceptance is not foreign native execution.
- Earlier dated release history and immutable schemas cannot be rewritten during finalization.
- A release that has not been independently verified cannot close #66 or unblock #67.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Prepare exactly one valid dated v1.1.0 changelog section and retain one fresh Unreleased section, consolidating duplicate categories without losing accumulated history or chronological decisions.
- **FR-002**: Preserve historical released changelog, schemas, legal/source assets and public inventories unchanged.
- **FR-003**: Maintain truthful candidate/public availability distinctions and exact current/historical capability and compatibility disclosures.
- **FR-004**: Freeze a separate concise public release-note body with intentional formatting and exact versioned full-changelog suffix; bind its future reviewed digest without implicit editing.
- **FR-005**: Specify a normative publication decision contract requiring the actual post-merge verified main revision and fresh accepted default-branch proof; no predicted target or populated approval may be substituted.
- **FR-006**: Require exact thirteen public assets, six target tuples, names/sizes/digests/checksum bijection, schema/legal identity and three same-bundle governed native proofs.
- **FR-007**: Require distinct explicit approvals for exact tag creation/push and release/asset publication, fail-closed partial-publication handling and refresh after any material change or artifact expiry.
- **FR-008**: Add meaningful metadata/decision-contract checks that reject invalid release preparation while retaining all established publication-disabled/read-only/security checks.
- **FR-009**: Complete installed Spec Kit design, blocking analysis, implementation, independent convergence and appropriate foreground verification before publication readiness.
- **FR-010**: Publish the explicitly authorized official PR closing #74 and referencing #66, verify formatted body read-back and exact-head checks, and handle every finding with at most two Codex rounds total.
- **FR-011**: Preserve native issue hierarchy/dependencies/milestone and unique Project membership; preparation tracks S029 and default Status stays unused.
- **FR-012**: Perform no tag/release/asset/production/public-schema/final-merge transition; leave #66/#67 and epic/milestone open until their independently authorized verified outcomes.

### Key Entities

- **Prepared release history**: Dated detailed candidate history with chronological decisions and immutable earlier release sections.
- **Public notes body**: Exact short highlights and compatibility disclosures selected for later authorized publication.
- **Decision contract**: Source, proof, asset, notes, authority and refresh obligations for a future exact publication package.
- **Preparation evidence**: Independent #74 acceptance, local/hosted verification, exact-head reviews and native delivery lifecycle.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All six #74 acceptance criteria have explicit coverage with zero material analysis/convergence findings.
- **SC-002**: One valid dated current release section, one Unreleased section and unchanged historical history pass all positive/negative metadata checks.
- **SC-003**: Every governed target/asset/proof/notes/authority/refresh obligation is explicit in the decision contract with no predicted revision or approval.
- **SC-004**: Official PR bodies render correctly, exact-head checks are green, and no actionable review remains within the two-round limit.
- **SC-005**: Zero protected publication/production/final-merge transitions occur; #66/#67 and epic/milestone remain open with truthful Project/native lifecycle.

## Assumptions

- Clean S028 post-merge main is the preparation baseline; all four default-branch workflows and three package native smokes passed.
- Fresh kickoff found only #51/#66/#67 open, with no operator-created arrivals; #74 is the narrowly scoped preparation outcome created during S029.
- Branch/PR authority is explicit and overrides the normal skill pre-push pause; release/tag and production authority remain separate.
- No new shipped runtime, schema, publisher, signing, installer or Site route is necessary for this metadata/decision slice.
- A future exact machine publication record can be populated only after human preparation merge and fresh accepted main proof, without changing this reviewed source to predict its own identity.
