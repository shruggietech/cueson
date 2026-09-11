import assert from "node:assert/strict";
import test from "node:test";
import { validatePostMutation, validatePreMutation } from "../scripts/verify-cloudflare-state.mjs";

function state() {
  return {
    account_id: "account",
    zone_id: "zone",
    zone: { id: "zone", name: "cueson.io", status: "active", paused: false, type: "full", account_id: "account" },
    dns_records: [{ id: "unrelated-dns", name: "api.cueson.io", type: "A", content: "192.0.2.1" }],
    workers: [{ id: "unrelated-worker", etag: "one", modified_on: "unchanged" }],
    domains: [{ id: "unrelated-domain", hostname: "api.cueson.io", service: "unrelated-worker", zone_id: "zone" }],
    redirect_rulesets: [],
  };
}

test("preflight rejects target DNS that is not owned by the declared Worker", () => {
  const snapshot = state();
  snapshot.dns_records.push({ id: "conflict", name: "cueson.io", type: "A", content: "192.0.2.2" });
  assert.throws(() => validatePreMutation(snapshot), /DNS records exist without the declared Worker Custom Domain/);
});

test("preflight accepts an existing declared Worker binding", () => {
  const snapshot = state();
  snapshot.dns_records.push({ id: "apex", name: "cueson.io", type: "AAAA", content: "100::" });
  snapshot.domains.push({ id: "apex-domain", hostname: "cueson.io", service: "cueson-site", zone_id: "zone" });
  assert.doesNotThrow(() => validatePreMutation(snapshot));
});

test("post-mutation read-back proves both target bindings and preserves unrelated state", () => {
  const before = state();
  const after = structuredClone(before);
  after.workers.push({ id: "cueson-site", etag: "site", modified_on: "now" });
  for (const [index, hostname] of ["cueson.io", "www.cueson.io"].entries()) {
    after.dns_records.push({ id: `target-dns-${index}`, name: hostname, type: "AAAA", content: "100::" });
    after.domains.push({ id: `target-domain-${index}`, hostname, service: "cueson-site", zone_id: "zone" });
  }
  assert.doesNotThrow(() => validatePostMutation(before, after));
  after.workers[0].etag = "changed";
  assert.throws(() => validatePostMutation(before, after), /unrelated Workers/);
});
