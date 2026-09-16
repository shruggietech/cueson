# Research: Prepare v1.1.0 public schema and site

## Decision 1: Split preparation from production

**Decision**: Deliver the official PR against child #76 and leave #67 open.

**Rationale**: #67 requires post-merge Cloudflare mutation, public DNS/TLS/HTTP/content verification and milestone closure. Those facts cannot exist on a pull-request head, and production authority is separate from the user's branch/PR authorization.

**Alternatives considered**: Closing #67 from the PR would assert unperformed production acceptance. Omitting a closing reference would violate normal PR policy. Reopening #67 after an automatic close would create avoidable false lifecycle state.

## Decision 2: Use the verified release as current truth

**Decision**: Bind site content to final v1.1.0 at exact tag `v1.1.0` and official release URL.

**Rationale**: #66 independently verified the final non-draft release, exact revision, 13 assets and public package behavior. Current site and documentation claims that v1.1.0 is unpublished are now false.

**Alternatives considered**: Continuing to call v1.1.0 a candidate would mislead users. Querying GitHub dynamically during static generation would add network dependence and a mutable authority.

## Decision 3: Centralize current release identity

**Decision**: Add one validated release record to the maintained site content map and derive generated metadata/user-facing release identity from it.

**Rationale**: Current identity is duplicated across content-map downloads, site library/page copy and generator-specific v1.0.0 validation. A compact authority reduces drift without expanding architecture.

**Alternatives considered**: Replacing strings individually would preserve the present drift risk. A remote release API at runtime is unnecessary and would weaken reproducibility.

## Decision 4: Preserve all prior public contracts

**Decision**: Add v1.1.0 release/schema entries while retaining v0.0.0 and v1.0.0 routes and bytes.

**Rationale**: Canonical schema URIs and historical release pages are durable public contracts. The new release is additive.

**Alternatives considered**: Replacing v1.0.0 entries would break existing references. A mutable latest route contradicts the versioned immutable contract.

## Decision 5: Copy immutable schema bytes exactly

**Decision**: Reuse the S028-admitted `schema/releases/v1.1.0/cueson.schema.json` exact byte copy of the independently verified canonical/release schema, and verify length 185,641 plus SHA-256 `223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7`.

**Rationale**: Release schema descriptions include time-bound candidate wording but the bytes are now immutable. Editing those bytes for prose cleanup would violate release identity.

**Alternatives considered**: Recopying, regenerating or editing the schema is unnecessary and forbidden after publication. Serving the mutable embedded source directly would not provide an immutable historical source.

## Decision 6: Verify public state against local reviewed truth

**Decision**: Make production verification compare remote deployment metadata to locally generated metadata and fetch/hash the public content manifest, then probe the local declared inventories.

**Rationale**: Trusting a deployment's own route/download list allows omissions to validate themselves. The reviewed source must define completeness.

**Alternatives considered**: Checking only commit equality is insufficient because an incomplete artifact could carry the correct commit. Checking only remote internal consistency retains the omission gap.

## Decision 7: Require exact current main for deployment

**Decision**: Require selected deployment revision equality with freshly fetched `origin/main`.

**Rationale**: The existing ancestor check can deploy a stale main revision after newer changes arrive, contradicting #67's exact-main acceptance.

**Alternatives considered**: Ancestor acceptance improves rollback flexibility but silently weakens this release gate. A rollback needs its own explicit, reviewed recovery decision.

## Decision 8: Update current prose without rewriting evidence

**Decision**: Correct maintained current-state documents and the detailed changelog while leaving historical Spec Kit preparation artifacts, released schemas and earlier release records intact.

**Rationale**: Users need accurate present behavior, while frozen preparation/release evidence must continue to describe what was true when reviewed.

**Alternatives considered**: Global version-string replacement would corrupt historical claims. Updating only the homepage would leave authoritative documentation inconsistent.
