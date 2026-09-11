# Public Route Contract

## Canonical origin

The canonical origin is `https://cueson.io`. Indexable pages use apex canonical URLs. `https://www.cueson.io/{path}?{query}` returns 308 with `Location: https://cueson.io/{path}?{query}`.

## Required HTML routes

| Route | Source | Purpose |
|------|------|------|
| `/` | site landing composition from released repository facts | Product identity, current release, installation, examples, documentation, and downloads |
| `/docs/` | `README.md` | Documentation start |
| `/docs/architecture/` | `docs/architecture.md` | Architecture of record |
| `/docs/brand/` | `docs/brand.md` | Brand provenance and use |
| `/docs/cli/` | `docs/cli.md` | CLI contract |
| `/docs/compatibility/` | `docs/compatibility.md` | Compatibility contract |
| `/docs/conversion/` | `docs/conversion.md` | Conversion behavior |
| `/docs/schema/` | `docs/schema.md` | Cue JSON schema reference |
| `/docs/formats/subrip/` | `docs/formats/srt.md` | SubRip contract |
| `/docs/formats/webvtt/` | `docs/formats/webvtt.md` | WebVTT contract |
| `/docs/project-management/` | `docs/project-management.md` | Delivery governance |
| `/docs/release-process/` | `docs/release-process.md` | Release and deployment process |
| `/docs/release-verification/` | `docs/release-verification.md` | Release evidence |
| `/docs/releases/v0.0.0/` | `docs/releases/v0.0.0.md` | Historical foundation release |
| `/docs/releases/v1.0.0/` | `docs/releases/v1.0.0.md` | Stable release notes |
| `/docs/project-specification/` | `docs/Cueson-Project-Specification-v0.0.0.md` | Working product roadmap |
| `/docs/contributing/` | `CONTRIBUTING.md` | Contribution guidance |
| `/docs/security/` | `SECURITY.md` | Security policy |
| `/docs/changelog/` | `CHANGELOG.md` | Detailed history |
| `/guides/media-formats/` | `docs/cueson-media-format-guide.html` | Direct authoritative format guide |

## Required non-HTML routes

| Route | Source | Required response |
|------|------|------|
| `/schema/v0.0.0/cueson.schema.json` | `schema/releases/v0.0.0/cueson.schema.json` | 200, JSON media type, exact bytes, immutable cache policy |
| `/schema/v1.0.0/cueson.schema.json` | `schema/releases/v1.0.0/cueson.schema.json` | 200, JSON media type, exact bytes, immutable cache policy |
| `/deployment.json` | generated revision record | 200, JSON media type, no-cache policy |
| `/schema/latest/cueson.schema.json` | none | 404 and never a schema response |

## Download targets

The landing page links to the official GitHub v1.0.0 release and its Windows amd64/arm64, macOS amd64/arm64, Linux amd64/arm64, and checksum assets. Verification accepts the GitHub redirect chain only when its first response identifies the intended release asset and the terminal response succeeds.

## HTML metadata

Every indexable page has one non-empty unique title, one non-empty description, an apex canonical URL matching its route, Open Graph title/description/URL/image, Twitter card metadata, an official favicon, `lang="en"`, and responsive viewport metadata.

## Response policy

All apex responses set `X-Content-Type-Options: nosniff`, `Referrer-Policy: strict-origin-when-cross-origin`, frame-embedding protection, and a site-appropriate Content Security Policy. HTTPS production responses set HSTS. Versioned schemas and hashed static assets use immutable caching; HTML uses bounded revalidation; deployment metadata uses no-cache.
