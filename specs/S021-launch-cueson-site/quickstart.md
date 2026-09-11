# Quickstart: Launch cueson.io

## Prerequisites

- Node.js 24.x
- Corepack with pnpm 10.28.2
- Go 1.25.x for repository checks
- Chromium installed through Playwright for browser verification
- Cloudflare credentials only for the separately owner-controlled production step, with Zone Read, DNS Read, Workers Scripts Read/Write, and Zone Rulesets Read for the declared account and zone

## Local verification

From `site/`:

```text
corepack pnpm install --frozen-lockfile
corepack pnpm generate
corepack pnpm generate:check
corepack pnpm test:unit
corepack pnpm build
corepack pnpm test:browser
corepack pnpm verify:artifact
corepack pnpm exec wrangler deploy --dry-run
```

From the repository root, run the existing CI-parity Go, documentation, brand, formatting, whitespace, UTF-8, and mojibake checks in addition to the site suite.

## Independent story checks

### US1: Public product and documentation

Serve `site/out/`, open `/`, `/docs/`, `/docs/cli/`, `/docs/schema/`, `/docs/formats/subrip/`, `/docs/formats/webvtt/`, `/docs/security/`, `/docs/contributing/`, and `/guides/media-formats/` at 360, 768, and 1440 CSS pixels, then run the browser suite for navigation, links, metadata, keyboard access, overflow, and accessibility.

### US2: Immutable schemas

Compare `site/out/schema/v0.0.0/cueson.schema.json` and `site/out/schema/v1.0.0/cueson.schema.json` byte for byte and by SHA-256 with their corresponding files under root `schema/releases/`. Confirm `site/out/schema/latest/` does not exist and a served request returns 404.

### US3: Owner-controlled deployment

Before mutation, confirm the selected full SHA is the merged `main` revision and run `corepack pnpm verify:cloudflare -- --phase before --snapshot PATH` to read and validate the live Cloudflare zone, DNS, Worker, Custom Domain, and redirect state. Deploy the verified static artifact and Worker, run the same command with `--phase after` to prove both target bindings and preservation of unrelated state, then execute the DNS, TLS, HTTP, redirect, route, download, metadata, and schema-byte probes in [deployment.md](contracts/deployment.md).

## Expected result

All repository and site checks pass, pull requests remain non-deploying, `cueson.io` serves the reviewed artifact after human merge, `www` redirects permanently without losing path or query, both canonical release schemas resolve with exact immutable bytes, and no mutable schema alias exists.
