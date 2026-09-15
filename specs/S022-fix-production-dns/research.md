# S022 Research

**Date**: 2026-09-15

## Independent Research

The Spec Kit plan workflow delegated bounded read-only research to `dns_research`. It confirmed that Node promise address queries reject independently, successful HTTP does not imply successful DNS status, and answers must carry the requested numeric record type and valid address value.

## Decisions

- **Decision**: Use independent A/AAAA queries with `Promise.allSettled` and async invocation wrappers. **Rationale**: Synchronous and asynchronous family failures cannot prevent the sibling query. **Alternatives**: Sequential fallbacks or shared catch blocks can skip evidence.
- **Decision**: Validate family-specific IP strings with built-in `node:net`. **Rationale**: Reject malformed values and wrong families without a dependency. **Alternatives**: Non-empty collections or custom regexes admit unsuitable evidence.
- **Decision**: Missing Answer is no-address evidence; malformed payload/status/Answer and unsuccessful HTTP/DNS state are family errors. **Rationale**: Empty answers and operational failures remain distinguishable. **Alternatives**: Swallowing errors prevents useful diagnosis.
- **Decision**: Sort and deduplicate addresses; report failed families in A/AAAA order. **Rationale**: Arrival order cannot change observable results. **Alternatives**: Completion-order inventories are unstable.
- **Decision**: Preserve Google endpoint, header, timeout, and all later probes. **Rationale**: Only address-family acceptance needs correction. **Alternatives**: Provider switching, retries, TLS rewrites, and DNS mutation expand scope.

## Primary References

- [Node promise DNS](https://nodejs.org/api/dns.html#dnspromisesresolve4hostname-options): queries yield address arrays and reject on failure.
- [Node DNS errors](https://nodejs.org/api/dns.html#error-codes): absent data and operational errors have distinguishable codes.
- [Node IP validation](https://nodejs.org/api/net.html#netisipinput): returns 4, 6, or 0 for family validation.
- [Google JSON DNS API](https://developers.google.com/speed/public-dns/docs/doh/json): A/AAAA types, numeric Status, optional Answer, and record data.
- [Google HTTP status](https://developers.google.com/speed/public-dns/docs/doh#http-status-codes): transport and DNS status are separate evidence.

All unknowns are resolved; no provider, toolchain, or public product contract changes are required.
