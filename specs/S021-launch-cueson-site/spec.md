# Feature Specification: Launch cueson.io

**Feature Branch**: `codex/S021-launch-cueson-site`

**Created**: 2026-09-11

**Status**: Approved for implementation

**Input**: User description: "Make `cueson.io` the public home for Cueson in one post-v1 work slice, publish the authoritative documentation and immutable released schemas, use approved brand assets, configure Cloudflare, and drive the work through a green and fully reviewed pull request before the human merge ritual."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Learn and use Cueson from its public home (Priority: P1)

A visitor can understand what Cueson does, install the current release, find practical command examples, and navigate the maintained product, CLI, format, compatibility, security, contribution, and release documentation from `cueson.io` on desktop or mobile.

**Why this priority**: The released software currently has no public home beyond its source repository. A coherent product and documentation experience is the primary value of domain activation.

**Independent Test**: Build the public artifact, serve it locally, open every planned route at desktop and mobile viewports, and verify navigation, content, links, metadata, keyboard access, and automated accessibility results.

**Acceptance Scenarios**:

1. **Given** a visitor opens the apex domain, **When** the landing page loads, **Then** the visitor sees the official Cueson identity, a concise product explanation, the current stable release, installation options, and clear documentation and download actions.
2. **Given** a visitor follows documentation navigation, **When** they open CLI, schema, SubRip, WebVTT, conversion, compatibility, security, contribution, changelog, or release material, **Then** the corresponding maintained repository content is available through a stable public route with working internal links.
3. **Given** a visitor uses a narrow mobile viewport or keyboard-only navigation, **When** they browse the site, **Then** the content remains readable, navigation remains operable, and no page requires horizontal scrolling at the supported viewport widths.

---

### User Story 2 - Resolve an immutable released schema (Priority: P1)

A consumer can resolve the canonical schema URI emitted by Cueson v0.0.0 or v1.0.0 and receive exactly the immutable schema bytes published with that release.

**Why this priority**: Canonical schema identifiers are already embedded in released Cue JSON. Making those identifiers resolvable completes the public contract promised after domain activation.

**Independent Test**: Fetch both versioned public schema URLs, compare their raw bytes and SHA-256 values with the immutable repository copies, verify their content type, and verify that no mutable `latest` schema route is published.

**Acceptance Scenarios**:

1. **Given** a consumer requests `/schema/v0.0.0/cueson.schema.json`, **When** the response succeeds, **Then** its body is byte-identical to the immutable v0.0.0 repository schema and identifies v0.0.0.
2. **Given** a consumer requests `/schema/v1.0.0/cueson.schema.json`, **When** the response succeeds, **Then** its body is byte-identical to the immutable v1.0.0 repository schema and identifies v1.0.0.
3. **Given** a consumer requests a mutable `latest` schema path, **When** the public host handles the request, **Then** no schema document is returned from that path.

---

### User Story 3 - Publish and verify an owner-controlled site (Priority: P1)

The maintainer can deploy the exact reviewed default-branch artifact to the managed Cloudflare zone without granting arbitrary pull requests production authority, then prove the resulting DNS, TLS, redirects, routes, and schema identities independently.

**Why this priority**: A site that only builds locally does not give Cueson real estate on the Internet. The deployment and public verification are part of this slice's outcome.

**Independent Test**: Confirm no pull-request event can deploy, select the reviewed default-branch revision, deploy its verified artifact with owner-controlled credentials, read Cloudflare state back, and run public DNS, TLS, HTTP, redirect, route, content, and schema-byte probes.

**Acceptance Scenarios**:

1. **Given** an arbitrary pull request, **When** hosted validation runs, **Then** it can build and test the site but cannot receive production credentials or mutate the production deployment.
2. **Given** the reviewed pull request has been merged by the human operator, **When** the maintainer deploys the selected default-branch revision, **Then** the deployed artifact is traceable to that revision and Cloudflare reports the intended apex-domain configuration without deleting unrelated zone records.
3. **Given** production activation completed, **When** a visitor opens `https://www.cueson.io`, **Then** the request redirects permanently to the equivalent `https://cueson.io` location without losing its path or query.
4. **Given** production activation completed, **When** independent verification probes the apex, representative documentation, download targets, and schemas, **Then** DNS, TLS, status, redirect, metadata, content, and byte-identity checks all pass.

### Edge Cases

- A maintained source document contains repository-relative links, heading fragments, raw HTML, fenced code, or tables that require a deterministic web adaptation.
- A generated adaptation is stale, missing, or differs from what the current authoritative source would produce.
- A requested documentation path is unknown or differs only by a trailing slash.
- A release download is unavailable or renamed upstream after the static artifact was built.
- A versioned schema source is missing, modified, normalized, or copied to the wrong public route.
- An existing Cloudflare DNS record or zone setting conflicts with the desired apex or `www` route.
- Cloudflare deployment succeeds but certificate issuance, DNS propagation, or the custom-domain binding is not yet publicly usable.
- Production verification receives cached, compressed, redirected, or error content instead of the expected raw schema bytes.
- JavaScript is unavailable, prefers-reduced-motion is enabled, or the viewport is narrow.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The public artifact MUST provide an official Cueson landing page with product identity, value, current stable version, installation guidance, primary command examples, documentation entry points, and release downloads.
- **FR-002**: The public artifact MUST provide stable routes for the maintained architecture, CLI, schema, SubRip, WebVTT, conversion, compatibility, security, contribution, changelog, brand, release-process, and v1.0.0 release material.
- **FR-003**: Public documentation content MUST originate from maintained root repository documents, and the build MUST reject stale or manually divergent generated adaptations.
- **FR-004**: Generated adaptations MUST preserve meaningful headings, paragraphs, lists, tables, code fences, links, and anchors from their authoritative sources or report an actionable build failure.
- **FR-005**: The public experience MUST use the approved repository Cueson brand kit and MUST NOT invent replacement logos, colors, typography, or slogans.
- **FR-006**: The site MUST provide coherent desktop and mobile navigation, a visible current-page state, keyboard-operable controls, visible focus, and a reduced-motion-safe experience.
- **FR-007**: Every public page MUST provide a unique title, useful description, canonical URL, social-preview metadata, favicon, language declaration, and responsive viewport metadata.
- **FR-008**: The public artifact MUST expose the v0.0.0 and v1.0.0 schemas at their exact canonical versioned paths with byte-for-byte identity to the corresponding immutable repository files.
- **FR-009**: The public artifact MUST NOT expose a mutable `latest` schema document or rewrite a released schema.
- **FR-010**: Schema responses MUST use an appropriate JSON media type and remain verifiable from raw response bytes regardless of transfer compression.
- **FR-011**: Site verification MUST check the static build, required routes, internal links and anchors, accessibility, representative desktop and mobile layouts, metadata, schema bytes, and production-equivalent deployable artifact.
- **FR-012**: Pull-request validation MUST run without production deployment authority, production credentials, or a production mutation path.
- **FR-013**: Production deployment MUST require an owner-controlled invocation selecting an exact default-branch revision whose repository and site checks have passed.
- **FR-014**: The deployment process MUST use the reviewed static artifact associated with the selected revision and MUST record enough evidence to identify that revision after deployment.
- **FR-015**: Domain activation MUST inspect existing Cloudflare zone and DNS state before mutation, make only the smallest necessary changes, preserve unrelated records, and read the resulting state back.
- **FR-016**: The apex domain MUST serve the public artifact over valid HTTPS, and `www.cueson.io` MUST permanently redirect to the equivalent apex path and query.
- **FR-017**: Post-deployment verification MUST independently check public DNS, TLS, apex and `www` HTTP behavior, representative site routes, release downloads, metadata, and both schema hashes.
- **FR-018**: Deployment or verification failures MUST remain visible, MUST NOT be reported as successful activation, and MUST leave repeatable diagnostic evidence for retry.
- **FR-019**: Repository documentation MUST record the site architecture, source-to-generated boundary, deployment authority, Cloudflare configuration, public route contract, and verification commands.
- **FR-020**: The detailed changelog MUST record the public-site architecture and production-domain activation decision without altering published release history.

### Key Entities

- **Authoritative document**: A maintained Markdown or HTML source in root `docs/` or a maintained root repository document that supplies public content.
- **Generated adaptation**: A deterministic build-only representation of one authoritative document, with source identity and public-route metadata but no independent editing authority.
- **Public route**: A stable path, page purpose, source relationship, expected metadata, and expected HTTP behavior in the deployable artifact.
- **Released schema**: An immutable version, canonical public path, source repository path, exact byte length, SHA-256 identity, and JSON media type.
- **Deployment artifact**: A verified static output tied to one repository revision and suitable for owner-controlled production publication.
- **Domain binding**: The Cloudflare-managed apex custom domain, `www` redirect, DNS and certificate state, and retained unrelated records.
- **Verification record**: The selected revision and results of build, route, link, accessibility, responsive, metadata, DNS, TLS, redirect, and schema-byte checks.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A first-time visitor can reach installation instructions, one working encode example, and the current release downloads from the landing page in no more than two navigation actions each.
- **SC-002**: All required public routes return their specified successful or redirect status, and automated internal-link and fragment checks report zero broken targets.
- **SC-003**: Automated accessibility checks report zero serious or critical violations across the landing page and every distinct documentation layout at desktop and mobile viewports.
- **SC-004**: Representative pages at widths of 360, 768, and 1440 CSS pixels have no unintended horizontal overflow, obscured primary content, or inaccessible navigation.
- **SC-005**: Every indexable page passes the metadata contract with one unique title, one non-empty description, one apex canonical URL, and complete social-preview metadata.
- **SC-006**: Re-running source adaptation without source changes produces no diff, while changing an authoritative document without regenerating causes the drift check to fail.
- **SC-007**: Public v0.0.0 and v1.0.0 schema responses match their immutable repository sources byte for byte and by SHA-256, while the mutable `latest` path does not return a schema.
- **SC-008**: Pull-request workflows expose zero production credentials and contain zero event paths capable of deploying the production site.
- **SC-009**: Public verification after deployment confirms DNS answers, a trusted TLS certificate, apex success, path-preserving permanent `www` redirection, successful representative routes, and exact schema identity from at least two independent network clients or resolvers where applicable.
- **SC-010**: The deployment record identifies one merged default-branch commit, and a clean checkout of that commit reproduces the verified deployable artifact and its route and schema inventory.

## Assumptions

- Cueson v1.0.0 and both immutable repository schema copies remain published and unchanged.
- The existing approved Cueson brand kit is the only visual identity source required for launch.
- Root `docs/`, `README.md`, `CONTRIBUTING.md`, `SECURITY.md`, and `CHANGELOG.md` are maintained repository authorities; web adaptations are generated build inputs rather than new authored copies.
- The Cloudflare account already owns the active `cueson.io` zone and the operator-provided authorization covers the necessary API-backed configuration and post-merge production deployment for this slice.
- The apex domain is the canonical public origin, and `www` is redirect-only.
- Production deployment follows human merge so the public artifact can be reproduced from the exact reviewed default-branch revision.
- Authentication, analytics, comments, search indexing integrations, server-side application behavior, and mutable schema aliases are outside this launch slice.
