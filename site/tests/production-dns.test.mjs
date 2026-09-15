import assert from "node:assert/strict";
import test from "node:test";
import { verifyDnsOverHttps, verifySystemDns } from "../scripts/production-dns.mjs";

const hostname = "www.cueson.io";
const ipv4 = "192.0.2.1";
const ipv6 = "2001:db8::1";
const dnsError = (code) => Object.assign(new Error(`query failed: ${code}`), { code });

function systemDependencies(outcomes, calls) {
  return Object.fromEntries(["A", "AAAA"].map((type) => [type === "A" ? "resolve4" : "resolve6", () => {
    calls.push(type);
    const outcome = outcomes[type];
    if (outcome instanceof Error) throw outcome;
    return Promise.resolve(outcome);
  }]));
}

function dohDependencies(outcomes, calls) {
  return { fetch: (url, options) => {
    const request = new URL(url);
    assert.equal(request.origin + request.pathname, "https://dns.google/resolve");
    assert.equal(request.searchParams.get("name"), hostname);
    assert.equal(options.headers.accept, "application/dns-json");
    assert.ok(options.signal instanceof AbortSignal);
    const type = request.searchParams.get("type");
    calls.push(type);
    const outcome = outcomes[type];
    if (outcome instanceof Error) throw outcome;
    return Promise.resolve(outcome).then((addresses) => ({
      ok: true,
      status: 200,
      json: async () => ({ Status: 0, Answer: addresses.map((data) => ({ type: type === "A" ? 1 : 28, data })) }),
    }));
  } };
}

const resolvers = [
  { name: "system DNS", verify: verifySystemDns, dependencies: systemDependencies },
  { name: "DNS-over-HTTPS", verify: verifyDnsOverHttps, dependencies: dohDependencies },
];

for (const resolver of resolvers) {
  for (const scenario of [
    { name: "IPv4-only", A: [ipv4], AAAA: [], expected: [ipv4] },
    { name: "IPv6-only", A: [], AAAA: [ipv6], expected: [ipv6] },
    { name: "dual-stack with duplicates", A: ["192.0.2.2", ipv4, ipv4], AAAA: [ipv6, ipv6], expected: [ipv4, "192.0.2.2", ipv6] },
    { name: "IPv4 no-data cannot suppress IPv6", A: dnsError("ENODATA"), AAAA: [ipv6], expected: [ipv6] },
    { name: "IPv6 failure cannot suppress IPv4", A: [ipv4], AAAA: dnsError("ETIMEOUT"), expected: [ipv4] },
  ]) {
    test(`${resolver.name}: accepts ${scenario.name} and attempts both families`, async () => {
      const calls = [];
      assert.deepEqual(await resolver.verify(hostname, resolver.dependencies(scenario, calls)), scenario.expected);
      assert.deepEqual(calls.sort(), ["A", "AAAA"]);
    });
  }

  test(`${resolver.name}: starts both families before either completes`, async () => {
    const pending = { A: Promise.withResolvers(), AAAA: Promise.withResolvers() };
    const calls = [];
    const result = resolver.verify(hostname, resolver.dependencies({ A: pending.A.promise, AAAA: pending.AAAA.promise }, calls));
    await Promise.resolve();
    assert.deepEqual(calls.sort(), ["A", "AAAA"]);
    pending.AAAA.resolve([ipv6, ipv6]);
    pending.A.resolve([ipv4]);
    assert.deepEqual(await result, [ipv4, ipv6]);
  });
}

for (const resolver of resolvers) {
  for (const scenario of [
    { name: "both empty", A: [], AAAA: [], details: ["A (IPv4): no addresses", "AAAA (IPv6): no addresses"] },
    { name: "both failed", A: dnsError("ENOTFOUND"), AAAA: dnsError("ETIMEOUT"), details: ["A (IPv4): ENOTFOUND", "AAAA (IPv6): ETIMEOUT"] },
    { name: "failure plus no addresses", A: dnsError("ESERVFAIL"), AAAA: [], details: ["A (IPv4): ESERVFAIL", "AAAA (IPv6): no addresses"] },
  ]) {
    test(`${resolver.name}: diagnoses ${scenario.name} with both family outcomes`, async () => {
      const calls = [];
      await assert.rejects(resolver.verify(hostname, resolver.dependencies(scenario, calls)), (error) => {
        assert.ok(error.message.includes(`${hostname}: ${resolver.name}`));
        for (const detail of scenario.details) assert.ok(error.message.includes(detail), error.message);
        assert.ok(error.message.indexOf("A (IPv4)") < error.message.indexOf("AAAA (IPv6)"));
        return true;
      });
      assert.deepEqual(calls.sort(), ["A", "AAAA"]);
    });
  }

  test(`${resolver.name}: an asynchronous A rejection cannot suppress completed AAAA`, async () => {
    const pending = Promise.withResolvers();
    const result = resolver.verify(hostname, resolver.dependencies({ A: pending.promise, AAAA: [ipv6] }, []));
    pending.reject(dnsError("ENODATA"));
    assert.deepEqual(await result, [ipv6]);
  });
}

function jsonResponse(payload) {
  return { ok: true, status: 200, json: async () => payload };
}

function rawDohDependencies(responses) {
  return { fetch: async (url) => {
    const type = new URL(url).searchParams.get("type");
    const response = responses[type];
    if (response instanceof Error) throw response;
    return response;
  } };
}

const emptyDoh = () => jsonResponse({ Status: 0 });
const validIpv6Doh = () => jsonResponse({ Status: 0, Answer: [{ type: 28, data: ipv6 }] });

for (const scenario of [
  { name: "absent answers", payload: { Status: 0 } },
  { name: "empty answers", payload: { Status: 0, Answer: [] } },
  { name: "alias-only answers", payload: { Status: 0, Answer: [{ type: 5, data: "target.example" }] } },
  { name: "wrong record type", payload: { Status: 0, Answer: [{ type: 28, data: ipv6 }] } },
  { name: "wrong address family", payload: { Status: 0, Answer: [{ type: 1, data: ipv6 }] } },
  { name: "invalid address value", payload: { Status: 0, Answer: [{ type: 1, data: "999.0.0.1" }] } },
  { name: "CIDR address value", payload: { Status: 0, Answer: [{ type: 1, data: "192.0.2.1/24" }] } },
  { name: "non-string address value", payload: { Status: 0, Answer: [{ type: 1, data: 123 }] } },
  { name: "malformed records", payload: { Status: 0, Answer: [null, {}, "invalid"] } },
  { name: "null payload", payload: null },
  { name: "array payload", payload: [] },
  { name: "string payload", payload: "invalid" },
  { name: "missing DNS status", payload: {} },
  { name: "string DNS status", payload: { Status: "0", Answer: [{ type: 1, data: ipv4 }] } },
  { name: "unsuccessful DNS status with addresses", payload: { Status: 2, Answer: [{ type: 1, data: ipv4 }] } },
  { name: "null answer collection", payload: { Status: 0, Answer: null } },
  { name: "object answer collection", payload: { Status: 0, Answer: {} } },
]) {
  test(`DNS-over-HTTPS: ${scenario.name} cannot prove resolution`, async () => {
    await assert.rejects(verifyDnsOverHttps(hostname, rawDohDependencies({ A: jsonResponse(scenario.payload), AAAA: emptyDoh() })), /www\.cueson\.io: DNS-over-HTTPS.*A \(IPv4\):.*AAAA \(IPv6\):/);
  });

  test(`DNS-over-HTTPS: usable IPv6 survives ${scenario.name}`, async () => {
    assert.deepEqual(await verifyDnsOverHttps(hostname, rawDohDependencies({ A: jsonResponse(scenario.payload), AAAA: validIpv6Doh() })), [ipv6]);
  });
}

test("DNS-over-HTTPS: CNAME plus actual target addresses is usable", async () => {
  assert.deepEqual(await verifyDnsOverHttps(hostname, rawDohDependencies({
    A: jsonResponse({ Status: 0, Answer: [{ type: 5, data: "target.example" }, { type: 1, data: ipv4 }] }),
    AAAA: emptyDoh(),
  })), [ipv4]);
});

for (const scenario of [
  { name: "HTTP failure", response: { ok: false, status: 503 }, detail: "HTTP status 503" },
  { name: "transport failure", response: new Error("connection refused"), detail: "connection refused" },
  { name: "fetch failure with DNS cause", response: new TypeError("fetch failed", { cause: Object.assign(new Error("getaddrinfo ENOTFOUND dns.google"), { code: "ENOTFOUND" }) }), detail: "ENOTFOUND: getaddrinfo ENOTFOUND dns.google" },
  { name: "timeout", response: new DOMException("resolver timeout", "TimeoutError"), detail: "resolver timeout" },
  { name: "invalid JSON", response: { ok: true, status: 200, json: async () => { throw new SyntaxError("invalid JSON"); } }, detail: "invalid JSON" },
]) {
  test(`DNS-over-HTTPS: ${scenario.name} is retained in complete failure`, async () => {
    await assert.rejects(verifyDnsOverHttps(hostname, rawDohDependencies({ A: scenario.response, AAAA: emptyDoh() })), (error) => {
      assert.ok(error.message.includes(hostname));
      assert.ok(error.message.includes(scenario.detail));
      assert.ok(error.message.includes("AAAA (IPv6): no addresses"));
      return true;
    });
    assert.deepEqual(await verifyDnsOverHttps(hostname, rawDohDependencies({ A: scenario.response, AAAA: validIpv6Doh() })), [ipv6]);
  });
}

for (const invalid of [null, {}, [ipv6], [123], ["invalid"]]) {
  test(`system DNS: invalid A response ${JSON.stringify(invalid)} cannot establish resolution`, async () => {
    await assert.rejects(verifySystemDns(hostname, systemDependencies({ A: invalid, AAAA: [] }, [])), /invalid IPv4 address response/);
    assert.deepEqual(await verifySystemDns(hostname, systemDependencies({ A: invalid, AAAA: [ipv6] }, [])), [ipv6]);
  });
}
