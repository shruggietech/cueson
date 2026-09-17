# Research: Repair protected production deployment

## Decision 1: Split reviewable implementation from post-merge proof

**Decision**: Create one implementation issue for the workflow, tests, and durable credential contract. The pull request closes that issue and references #78. Issue #78 remains open until the merged default-branch workflow completes the authorized production validation.

**Rationale**: A pull request cannot prove a workflow definition that exists only after merge. Closing #78 at merge would incorrectly claim the protected production path was operational before the default-branch run and Cloudflare read-back.

**Alternatives considered**: Close #78 from the pull request and reopen it on a deployment failure. Rejected because lifecycle state would knowingly become false. Keep one issue with no closing reference. Rejected because normal pull requests require a complete closing reference and the implementation is independently testable.

## Decision 2: Bind workflow execution identity to the deployment revision

**Decision**: Require `GITHUB_REF` to equal `refs/heads/main` and `GITHUB_SHA` to equal the requested revision before checkout, then preserve the two freshly fetched `origin/main` equality checks before artifact proof and mutation.

**Rationale**: GitHub environment deployment records use the workflow execution ref and SHA. A recovery-branch workflow can deploy exact main while GitHub records the recovery commit. Binding dispatch identity makes the GitHub deployment record and public record describe the same reviewed commit.

**Alternatives considered**: Continue accepting another workflow-definition branch when its input equals main. Rejected because it produces contradictory deployment evidence. Infer identity only from checkout state. Rejected because the environment record is created from workflow execution metadata.

## Decision 3: Remove the package-script separator

**Decision**: Invoke `pnpm verify:cloudflare --phase ...` and `pnpm verify:production --expected-commit ...` directly.

**Rationale**: The current pnpm invocation forwards a literal `--` to the Node script. The Cloudflare parser rejects it, causing the protected workflow to stop before mutation. Direct named arguments are the declared package-script contract.

**Alternatives considered**: Teach both scripts to ignore a literal separator. Rejected because it masks a malformed workflow boundary and makes parser behavior inconsistent. Invoke Node scripts directly. Rejected because tests and operators should exercise the published package-script interface used by CI.

## Decision 4: Exercise the real command boundary

**Decision**: Add focused Node tests that spawn the platform-specific Corepack executable without a shell, from the site workspace, with hidden Windows process behavior. The tests assert each command passes argument parsing and reaches the next deterministic validation failure.

**Rationale**: Static text assertions allowed the invalid separator to ship. An actual subprocess test covers pnpm forwarding and script parsing together while remaining credential-free and non-mutating.

**Alternatives considered**: Keep only static workflow assertions. Rejected because they cannot prove package-manager forwarding. Unit-test parser helpers only. Rejected because they omit the package-script layer where the defect occurred.

## Decision 5: Use an account-owned restricted Cloudflare API token

**Decision**: Store an account-owned token as `CLOUDFLARE_API_TOKEN` in the protected GitHub `production` environment. Grant the entire ShruggieTech account Developer Platform `Workers Scripts Legacy: Edit`. Grant only the specified `cueson.io` domain DNS & Zones `DNS: Read` and `Zone: Read`, Rules & Configuration `Zone Transform Rules: Read`, and Developer Platform `Workers Routes: Read`.

**Rationale**: Worker upload, Worker inventory, and Custom Domain operations require Workers Scripts Legacy Edit at account scope. Wrangler 4.131.1 also reads `GET /zones/{zone_id}/workers/routes` to detect route conflicts before applying configured Custom Domains, which requires Workers Routes Read on the zone. The verifier reads zone identity, DNS records, Workers, Custom Domains, and redirect rulesets. The selected permissions cover those exact operations without Workers Routes Edit or broad zone administration.

**Alternatives considered**: Reuse a user OAuth session. Rejected because it is temporary, user-bound, and unsuitable for unattended protected Actions. Use a broad account token. Rejected because it violates least privilege. Grant Workers Routes Edit. Rejected because this configuration performs only the route-list read for conflict detection; Custom Domain mutation remains covered by Workers Scripts Legacy Edit.

## Decision 6: Limit secret exposure to Cloudflare-facing steps

**Decision**: Remove the token from job-level environment variables and add it only to preflight, deploy, and post-deployment read-back steps.

**Rationale**: Job-level exposure makes the token available to checkout, tool setup, dependency installation, tests, and local artifact proof. Step scope limits the blast radius while supporting every necessary Cloudflare call.

## Decision 7: Document a bounded credential lifecycle

**Decision**: The organization security owner creates the token, records a non-secret inventory entry, uses an explicit expiry where available, rotates at least every 90 days, validates replacement through preflight before retiring the prior token, and revokes immediately after suspected exposure or ownership change.

**Rationale**: A durable automation credential still needs time-bounded ownership and a tested replacement procedure. GitHub secret read-back can prove presence and update time but never reveal the value.

## Primary references

- Cloudflare Workers Scripts API: https://developers.cloudflare.com/api/resources/workers/subresources/scripts/methods/list/
- Cloudflare Workers Custom Domains API: https://developers.cloudflare.com/api/resources/workers/subresources/domains/methods/update/
- Cloudflare Workers Routes API: https://developers.cloudflare.com/api/resources/workers/subresources/routes/methods/list/
- Cloudflare DNS Records API: https://developers.cloudflare.com/api/resources/dns/subresources/records/methods/list/
- Cloudflare Zones API: https://developers.cloudflare.com/api/resources/zones/methods/list/
- Cloudflare Rulesets API: https://developers.cloudflare.com/api/resources/rulesets/methods/list/
