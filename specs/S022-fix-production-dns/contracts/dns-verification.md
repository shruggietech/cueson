# S022 DNS Verification Contract

## Operational Interface

`verifySystemDns(hostname, { resolve4, resolve6 } = defaults)` and `verifyDnsOverHttps(hostname, { fetch } = defaults)` are script-local exported helpers using built-in DNS and fetch defaults. Each returns a promise for sorted unique addresses or rejects with hostname/resolver/family diagnostics. Injection is a test boundary, not a public Cueson API.

Both queries start independently, including synchronous dependency throws. System queries request string address arrays. DoH uses the existing `https://dns.google/resolve` endpoint with encoded hostname, separate `type=A` and `type=AAAA`, `accept: application/dns-json`, and the 20-second request timeout. Success requires successful HTTP, numeric DNS Status 0, and valid requested-family address records.

Missing Answer or no usable matching records is absent-address evidence. Malformed payload/status/Answer, bad JSON, unsuccessful HTTP/DNS state, and transport errors are family errors. A usable sibling permits success; complete failure reports both outcomes in A/AAAA order. Duplicates do not change results.

## Protected Boundaries

The production verifier calls both helpers for apex and www inside the existing network-identity guard. Arguments, local metadata sources, expected revision, TLS, HTTP, downloads, redirects, schemas, missing-alias behavior, and success output retain their behavior. Site validation remains unprivileged; the manual exact-main deployment workflow is unchanged.

Live acceptance targets S021 commit `46838fd5cc888b299a09b89da676db0005b00f16`; local S022 build identity does not replace production identity.
