# Implementation Plan: Launch cueson.io

**Branch**: `codex/S021-launch-cueson-site` | **Date**: 2026-09-11 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `/specs/S021-launch-cueson-site/spec.md`

## Summary

Add a deterministic static documentation and product site under `site/`, generate its documentation inputs and selected web assets from the authoritative root documents and approved brand kit, publish immutable release-schema bytes, and verify the complete artifact with unit and browser tests. Deploy the reviewed artifact after human merge through Cloudflare Workers Static Assets using exact custom domains for `cueson.io` and `www.cueson.io`; a minimal Worker serves apex assets, applies response policy, and returns a path/query-preserving permanent redirect for `www`.

## Technical Context

**Language/Version**: TypeScript 5.9.3 and ECMAScript modules on Node.js 24.x; existing product remains Go 1.25.x

**Primary Dependencies**: Next.js 16.3.5, React and React DOM 19.3.0, Fumadocs Core/UI 16.15.9, Fumadocs MDX 15.4.0, Tailwind CSS 4.3.3, next-themes 0.4.6, Wrangler 4.131.1

**Storage**: Version-controlled root documents, immutable schema files, and brand kit inputs; generated static files only, with no database or mutable runtime storage

**Testing**: Node's built-in test runner for generators and Worker behavior; Playwright 1.63.0 with `@axe-core/playwright` 4.13.0 for route, link, accessibility, responsive, and metadata checks; Wrangler dry-run artifact verification; existing Go and repository checks

**Target Platform**: Standards-compliant desktop and mobile browsers served globally through Cloudflare Workers Static Assets

**Project Type**: Static web application plus a minimal edge request policy and repository automation

**Performance Goals**: Pre-render every page; avoid runtime content fetches; keep first-party landing-page JavaScript below 200 KiB compressed where practical; serve immutable schemas and hashed assets with long-lived cache policy

**Constraints**: Root documents remain authored authority; generated adaptations are build-only; released schema bytes cannot change; brand inputs remain byte-protected; no pull-request deployment credentials; production changes occur only after human merge; `www` redirects without losing path or query

**Scale/Scope**: One landing page, approximately seventeen generated documentation pages, two immutable schema routes, one static Worker, two custom domains, one hosted site-validation workflow, and one owner-controlled deployment workflow

## Constitution Check

*GATE: Passed before Phase 0 research and passed again after Phase 1 design.*

- **I. Lossless Source Preservation**: Site work does not modify Cue JSON or source envelopes. Generated copies preserve authoritative inputs, and published schemas use byte copies rather than parse/serialize cycles.
- **II. Schema and Official Software Discipline**: Only already released v0.0.0 and v1.0.0 schemas are published at immutable versioned paths. No `latest` alias is created.
- **III. Common Model Plus Native Fidelity**: Documentation explains the maintained common-model and native-fidelity boundaries without changing them.
- **IV. No Silent Loss**: The adaptation generator fails on unmapped sources, broken mappings, unsafe paths, stale output, or copy drift instead of silently dropping content.
- **V. Test-First Format Work**: No codec behavior changes. Site generator, Worker, route, link, metadata, responsive, accessibility, and schema-byte behavior receive focused automated tests.
- **VI. Portable and Secure Operation**: The shipped Go product is unchanged. The static site treats documentation as build input, emits no machine-local paths, uses local fonts, and applies bounded browser-facing response policy.
- **VII. Documentation and Delivery Authority**: Root documentation remains authoritative, S021 uses the full Spec Kit and analysis path, issue #47 and Project item S021 own delivery, GitHub publication bodies are formatted and read back, and the human retains final merge authority.
- **Protected actions**: The operator explicitly authorized S021 kickoff administration, push, pull-request creation, GitHub configuration, Cloudflare configuration, and post-merge production deployment. This authorization does not authorize the AI to merge the pull request.
- **Post-design re-check**: The route, generation, deployment, and verification contracts preserve every gate above. No constitutional exception or complexity waiver is required.

## Project Structure

### Documentation (this feature)

```text
specs/S021-launch-cueson-site/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── content-generation.md
│   ├── deployment.md
│   └── routes.md
├── checklists/
│   └── requirements.md
└── tasks.md
```

### Source Code (repository root)

```text
site/
├── app/
│   ├── (home)/
│   ├── docs/[[...slug]]/
│   ├── layout.tsx
│   ├── not-found.tsx
│   └── global.css
├── components/
├── content/
│   └── generated/                 # build-only, ignored
├── lib/
├── public/
│   ├── assets/                    # build-only approved copies, ignored
│   ├── schema/                    # build-only byte copies, ignored
│   └── deployment.json            # build-only revision record, ignored
├── scripts/
│   ├── generate.mjs
│   ├── verify-artifact.mjs
│   └── serve.mjs
├── tests/
│   ├── generator.test.mjs
│   └── site.spec.ts
├── worker/
│   ├── index.ts
│   └── index.test.ts
├── content-map.json
├── next.config.mjs
├── playwright.config.ts
├── source.config.ts
├── tsconfig.json
├── wrangler.jsonc
├── package.json
└── pnpm-lock.yaml

.github/workflows/
├── site.yml
└── site-deploy.yml

docs/
├── architecture.md
├── release-process.md
└── schema.md

CHANGELOG.md
.gitignore
.gitattributes
```

**Structure Decision**: `site/` is a separately managed pnpm package so Node dependencies and generated outputs do not enter the pure-Go product module. Root documents and brand files stay in place; the generator produces disposable MDX, web assets, schema paths, and revision metadata inside ignored `site/` subtrees before each build. Cloudflare configuration is colocated with the artifact it deploys.

## Complexity Tracking

No constitution violation or unjustified complexity is introduced.
