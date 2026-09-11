# Research: Launch cueson.io

## Decision 1: Use the approved current house stack with exact pins

**Decision**: Pin Node.js 24.x, pnpm 10.28.2, TypeScript 5.9.3, ESLint 9.39.5, Next.js 16.3.5, React and React DOM 19.3.0, Fumadocs Core/UI 16.15.9, Fumadocs MDX 15.4.0, Tailwind CSS 4.3.3, next-themes 0.4.6, Playwright 1.63.0, `@axe-core/playwright` 4.13.0, and Wrangler 4.131.1.

**Rationale**: The application packages are the current registry releases on 2026-09-11. Fumadocs Core/UI target Next.js 16 and React 19.2 or newer, while Fumadocs MDX accepts Next.js 15.3 or 16 and Fumadocs Core 16.15.3 or newer. Wrangler requires Node.js 22 or newer, satisfied by Node.js 24. TypeScript 7.0.2 and ESLint 10.10.0 are newer registry tags but violate the declared peer ranges of the current Next.js lint dependency graph, so the newest supported TypeScript 5.9.3 and ESLint 9.39.5 releases are pinned instead.

**Alternatives considered**: Older pre-v1 house pins would add immediate upgrade debt. Floating ranges would make clean-checkout output and review evidence unstable.

## Decision 2: Generate Fumadocs MDX from root authority

**Decision**: Maintain a version-controlled `content-map.json` and deterministic generator that reads root documentation, prepends controlled front matter, rewrites known repository-relative documentation links to public routes, preserves source bodies, emits a source/output hash manifest, and rejects stale or unmapped results. Generated MDX remains ignored build output.

**Rationale**: Fumadocs MDX expects a content directory, but the repository contract makes root `docs/` authoritative. A reproducible adapter gives Fumadocs its build input without asking maintainers to edit a second copy.

**Alternatives considered**: Hand-maintained MDX duplicates documentation authority. Moving root docs into the site would break repository links and established tooling. Rendering Markdown at request time conflicts with static export and adds runtime code.

## Decision 3: Use Next.js static export

**Decision**: Configure `output: "export"`, explicit trailing-slash behavior, and a revision-derived build ID. Generate all documentation route parameters at build time and exclude server-only Next.js features.

**Rationale**: Current Next.js guidance produces one HTML file per route under `out/`, supports App Router Server Components at build time, and can be hosted by any static server. A stable build ID removes avoidable cross-build randomness.

**Alternatives considered**: A Node.js server adds an origin and runtime maintenance without product value. A client-only single-page application weakens direct-route, metadata, and no-JavaScript behavior.

**Primary source**: https://nextjs.org/docs/app/guides/static-exports

## Decision 4: Use Cloudflare Workers Static Assets as the origin

**Decision**: Deploy `site/out/` with Wrangler's static-assets configuration, `404-page` not-found handling, automatic HTML handling, and a small TypeScript Worker using an `ASSETS` binding.

**Rationale**: Cloudflare's current Workers guidance recommends Static Assets for new static projects. Assets and Worker code deploy atomically, while the Worker can implement the narrow redirect and response-header policy.

**Alternatives considered**: Pages is viable but is no longer the preferred new-project path for this owner-controlled static origin. An external host would add another vendor and DNS integration surface.

**Primary sources**: https://developers.cloudflare.com/workers/static-assets/ and https://developers.cloudflare.com/workers/static-assets/routing/static-site-generation/

## Decision 5: Attach both hostnames as exact Custom Domains

**Decision**: Declare exact Custom Domains for `cueson.io` and `www.cueson.io`. Run the Worker before asset serving, return an HTTP 308 from `www` to the equivalent apex URL, and pass apex requests to the asset binding.

**Rationale**: Custom Domains are intended when the Worker is the origin and create DNS records and certificates automatically. Attaching both exact hostnames keeps redirect behavior in reviewed code and avoids a separately mutable redirect-ruleset dependency. The 308 changes only the hostname, so path and query remain intact.

**Alternatives considered**: A Worker route needs a pre-existing proxied DNS record and implies another origin. A Bulk or Single Redirect rule also works, but it adds an independent ruleset plus a placeholder DNS record for a behavior the Worker already owns.

**Primary sources**: https://developers.cloudflare.com/workers/configuration/routing/custom-domains/ and https://developers.cloudflare.com/workers/best-practices/workers-best-practices/

## Decision 6: Keep production deployment owner-controlled and post-merge

**Decision**: Pull requests run build, tests, browser verification, and `wrangler deploy --dry-run` only. A separate manual workflow and the authenticated Cloudflare connector can deploy only a selected `main` revision after the human merge. The deployment regenerates and re-verifies the artifact before mutation.

**Rationale**: This prevents fork or branch code from receiving production credentials, preserves the human merge boundary, and binds the public artifact to reviewed default-branch history.

**Alternatives considered**: Automatic deployment on every push to `main` is fast but expands mutation authority beyond the explicitly selected launch. Preview deployment from pull requests is unnecessary for this first static launch and would expose another public surface.

## Decision 7: Verify the public result independently

**Decision**: Record the Cloudflare pre-mutation inventory, read zone/DNS/Worker/custom-domain state after deployment, resolve DNS through system and public resolvers, inspect TLS and HTTP behavior, crawl representative pages, verify download targets, and compare raw public schema bytes and SHA-256 values against repository sources.

**Rationale**: An accepted API mutation or successful upload does not prove public DNS, certificate readiness, redirects, cached content, or exact schema bytes.

**Alternatives considered**: Dashboard screenshots and deployment command success are not reproducible evidence.

## Live-state findings

- ShruggieTech account ID: `39e3052d61e3edccea7d68269ec07182`.
- `cueson.io` zone ID: `ae7b2e3215426388c3012b5123784204`.
- Zone state on 2026-09-11: active, full setup, not paused, Free Website plan.
- Authoritative nameservers: `norah.ns.cloudflare.com` and `rex.ns.cloudflare.com`.
- Existing DNS records: zero.
- Existing Cueson Worker: none; four unrelated account Workers exist and remain outside scope.
- Existing zone rulesets: Cloudflare-managed normalization, managed-free, and DDoS rules only; no custom redirect ruleset exists.
- The local research pass was performed by the coordinating agent because this turn did not authorize sub-agent delegation.
