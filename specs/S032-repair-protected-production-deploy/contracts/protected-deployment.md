# Protected Production Deployment Contract

## Entry conditions

1. The workflow is manually dispatched from `main`.
2. The `revision` input is one full lowercase commit identifier.
3. `GITHUB_REF` is `refs/heads/main`.
4. `GITHUB_SHA` equals `revision`.
5. A fresh fetch proves `origin/main` equals `revision` before checkout-controlled proof.
6. The GitHub `production` environment has approved the run and supplies the protected credential only to Cloudflare-facing steps.

## Command boundary

The workflow invokes:

```text
corepack pnpm verify:cloudflare --phase before --snapshot <path>
corepack pnpm verify:cloudflare --phase after --snapshot <path>
corepack pnpm verify:production --expected-commit <revision>
```

No standalone `--` is passed to a verifier.

## Credential boundary

`CLOUDFLARE_API_TOKEN` is absent from job scope and all source checkout, setup, install, test, build, and local verification steps. It is present only for Cloudflare preflight, Wrangler deploy, and Cloudflare post-deployment read-back.

## Mutation gate

Immediately before deployment, another fresh fetch proves `origin/main` still equals `revision`. Failure stops the run without production mutation.

## Completion evidence

1. Preflight proves target account/zone identity and captures DNS, Workers, Custom Domains, and redirect rulesets.
2. Wrangler deploy succeeds for `cueson-site`.
3. Post-deployment read-back proves intended bindings and unrelated-state preservation.
4. Public verification proves exact deployment metadata, 25 routes, seven downloads, three immutable schemas, redirect behavior, and content-manifest integrity.
5. GitHub environment deployment metadata and public deployment metadata identify the same exact main revision.
