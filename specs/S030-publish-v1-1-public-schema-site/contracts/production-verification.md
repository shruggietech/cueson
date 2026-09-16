# Production verification contract

## Exact revision selection

The manual deployment input MUST be a full lowercase Git commit. Immediately before proof and mutation, the workflow MUST fetch the default branch and require the selected revision to equal `origin/main`. Being an ancestor is insufficient.

The selected revision MUST receive the complete repository and site verification suite before deployment. Production credentials MUST remain limited to the protected manual deployment job and MUST never reach pull-request execution.

## Local authority and public read-back

Verification MUST load the locally generated deployment record for the selected revision as the reviewed authority. It MUST fetch the public deployment record and compare source revision, current release identity, content-manifest digest, complete route inventory, complete seven-download inventory and complete three-schema inventory with local truth.

Verification MUST fetch the public `content-manifest.json` as raw bytes and require its SHA-256 to equal the local deployment record. It MUST drive route, download and schema probes from local reviewed arrays. A remote record that omits a local entry MUST fail even if internally self-consistent.

## Public behavior

Live verification MUST retain existing trusted TLS, apex success, path/query-preserving permanent `www` redirect, DNS system and DNS-over-HTTPS, metadata, content, link and immutable schema byte/hash checks. All 25 declared routes MUST behave as specified, all seven downloads MUST resolve, all three schemas MUST match exact bytes, and every mutable latest alias probe MUST fail to return a schema.

## Cloudflare preservation

Before mutation, deployment MUST capture authenticated zone, DNS, Worker/custom-domain and relevant configuration state and reject unresolved conflicts. After mutation, it MUST read state back and prove unrelated records/settings remain unchanged. Project-owned child processes MUST remain non-interactive and hidden on Windows.

## Failure handling

Any revision mismatch, local/remote metadata difference, manifest digest difference, missing/extra route or download, schema mismatch, TLS/DNS/redirect failure or Cloudflare preservation difference blocks success. The workflow MUST retain actionable evidence and MUST NOT report #67, epic #51 or milestone completion.
