# S035 Production Handoff Contract

## Preparation acceptance

The reviewed artifact must include 28 declared routes, 24 HTML, four immutable schemas and seven primary v1.2.0 downloads. The new release and speaker routes must participate in existing link/metadata/browser/accessibility checks. Historical schemas and release routes remain intact.

## Later deployment entry

1. Operator confirms this specific PR merge; verify actual squash merge and clean current main.
2. Fetch/prune and bind a fresh artifact/proof to full lowercase actual main revision, with current remote main rechecked.
3. Review exact generated deployment/content-manifest/schema/download identities; bind the recorded S035 production authorization to this concrete reviewed artifact/configuration. If material state changes invalidate this reviewed outcome, obtain a fresh decision.
4. Dispatch existing manual `site-deploy.yml` on main with `revision=<full-current-main>`. Retain protected environment and step-scoped credential boundaries; never print credentials.
5. Preserve workflow freshness checks before proof and immediately before mutation, infrastructure preflight/post-readback and unrelated-state preservation.

## Independent completion

- Verify deployment/environment/public revision identifies actual accepted main.
- Verify all 28 routes, seven official download bindings, public content-manifest digest and four immutable schema bytes/sizes/hashes.
- Verify HTTPS schema content type, immutable caching, nosniff, absence of mutable latest route, old routes and redirect behavior using established verifiers plus explicit header probes where the existing verifier does not assert those headers.
- Separately demonstrate public HTTPS connectivity over IPv4 and IPv6, including authoritative/public DNS evidence. DNS resolution by itself is insufficient; one successful address family cannot hide a failed other.

For each of `cueson.io` and `www.cueson.io`, run hidden `curl.exe -4 --fail --silent --show-error --max-time 20 --output <evidence-file> --write-out <status-and-remote-ip> https://<host>/deployment.json` and the equivalent `-6` command (use the apex destination after checking the www redirect without following). Require successful TLS verification, expected HTTP status/redirect, address family and exact deployment identity. Capture all four schema response headers with hidden curl `--dump-header`, compare bytes to reviewed hashes and require `application/schema+json`, `immutable` and `x-content-type-options: nosniff`. These explicit probes complement `pnpm verify:production --expected-commit <actual-main>` and the S022 DNS checks. If the operator host lacks a usable family, execute equivalent probes from an authorized dual-stack runner and retain that runner's evidence; never mark an untested family passed.
- Retain artifact/run/revision/authority/configuration/public evidence; reconcile #86 acceptance, dependencies, Project Stage and unused Status. Close milestone only when no committed outcome remains incomplete.
- Any stale main, wrong bytes/metadata, partial deployment or failed live check prevents completion and requires fresh decision/evidence. Do not move published tags, rewrite release assets or historical schemas as recovery.

## Authority at preparation time

Current kickoff explicitly authorizes push and official PR, all review findings and at most two Codex rounds. The owner subsequently instructed: "Modify the instructions then. Don't leave shit hanging" in direct response to deferring deployment. This authorizes completing S035 production deployment/live verification immediately after their final specific PR merge, without another kickoff or approval. Final PR merge remains the requested pause. No production mutation occurs before it.
