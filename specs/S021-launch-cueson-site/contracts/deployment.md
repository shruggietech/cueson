# Deployment Contract

## Pull-request boundary

The `Site` validation workflow runs on pull requests and `main` updates with `contents: read`. It installs pinned tooling, regenerates content, checks drift, builds the static export, runs unit and browser verification, verifies the deployable inventory, and performs a Wrangler dry run. It receives no Cloudflare secret and runs no deploy command without `--dry-run`.

## Owner-controlled production entry

Production deployment is available only through a manually dispatched workflow or an authenticated operator-side Cloudflare API/CLI invocation. The selected revision must be a full commit SHA on `main`. The process checks out that exact revision, verifies it equals the requested revision, regenerates all disposable inputs, runs the complete site suite, and deploys only the resulting `site/out/` with the reviewed Worker and Wrangler configuration.

## Cloudflare target

- Account: ShruggieTech (`39e3052d61e3edccea7d68269ec07182`)
- Zone: `cueson.io` (`ae7b2e3215426388c3012b5123784204`)
- Worker: `cueson-site`
- Custom Domains: `cueson.io` and `www.cueson.io`
- Apex behavior: serve static assets
- `www` behavior: 308 to the same apex path and query
- Unrelated Workers, zones, DNS records, and rulesets: untouched

## Pre-mutation checks

1. Confirm the selected SHA is the checked-out `main` revision and has green repository and site checks.
2. Read the zone status, relevant DNS records, existing `cueson-site` Worker metadata, Custom Domains, and redirect rules.
3. Stop if either hostname is owned by an unrelated resource or if a conflicting record cannot be reconciled without deleting unrelated state.
4. Verify the generated deployment record, route manifest, schema hashes, and Wrangler dry-run output.

## Mutation

Deploy Worker code and static assets atomically. Allow Wrangler or the Workers API to add or update only the two declared Custom Domains and their provider-managed DNS and certificates. Do not delete unrelated records or rules. Do not enable deployment from pull-request events.

## Read-back and public verification

1. Read the Worker version, routes or domains, both DNS records, and certificate or Custom Domain state back from Cloudflare.
2. Resolve apex and `www` through the system resolver and at least one public DNS-over-HTTPS resolver.
3. Verify a trusted TLS handshake and certificate hostname coverage.
4. Verify apex success, representative documentation routes, the media guide, deployment metadata, and every official download link.
5. Verify `www` returns 308 and preserves an arbitrary path and query.
6. Fetch both schemas without content transformation and compare raw length and SHA-256 against repository sources.
7. Verify `/schema/latest/cueson.schema.json` returns 404.
8. Record failures truthfully and retry only after identifying whether the remaining state is deployment, DNS, certificate issuance, caching, or content.

## Authority

The S021 kickoff grants one continuous authorization for its necessary GitHub administration, branch push, official pull request, Cloudflare configuration, and post-merge production deployment. The human operator alone performs the final pull-request merge. A material scope or target change requires a new decision.
