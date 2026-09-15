import { resolve4 as systemResolve4, resolve6 as systemResolve6 } from "node:dns/promises";
import { isIP } from "node:net";

const families = [
  { type: "A", version: 4, recordType: 1 },
  { type: "AAAA", version: 6, recordType: 28 },
];

function describeError(error) {
  return `${error?.code ? `${error.code}: ` : ""}${error?.message ?? String(error)}`;
}

async function verifyFamilies(hostname, resolver, query) {
  const outcomes = await Promise.allSettled(families.map(async (family) => query(family)));
  const addresses = [...new Set(outcomes.flatMap((outcome) => outcome.status === "fulfilled" ? outcome.value : []))].sort();
  if (addresses.length > 0) return addresses;
  const details = outcomes.map((outcome, index) => {
    const family = families[index];
    const reason = outcome.status === "fulfilled" ? "no addresses" : `${describeError(outcome.reason)}${outcome.reason?.cause ? ` (cause: ${describeError(outcome.reason.cause)})` : ""}`;
    return `${family.type} (IPv${family.version}): ${reason}`;
  });
  throw new Error(`${hostname}: ${resolver} returned no usable addresses; ${details.join("; ")}`);
}

export async function verifySystemDns(hostname, { resolve4 = systemResolve4, resolve6 = systemResolve6 } = {}) {
  return verifyFamilies(hostname, "system DNS", async ({ version }) => {
    const addresses = await (version === 4 ? resolve4 : resolve6)(hostname);
    if (!Array.isArray(addresses) || addresses.some((address) => typeof address !== "string" || isIP(address) !== version)) throw new Error(`invalid IPv${version} address response`);
    return addresses;
  });
}

export async function verifyDnsOverHttps(hostname, { fetch: fetchDns = globalThis.fetch } = {}) {
  return verifyFamilies(hostname, "DNS-over-HTTPS", async ({ type, version, recordType }) => {
    const response = await fetchDns(`https://dns.google/resolve?name=${encodeURIComponent(hostname)}&type=${type}`, { headers: { accept: "application/dns-json" }, signal: AbortSignal.timeout(20_000) });
    if (!response.ok) throw new Error(`HTTP status ${response.status}`);
    const result = await response.json();
    if (!result || typeof result !== "object" || Array.isArray(result) || !Number.isInteger(result.Status)) throw new Error("invalid DNS response: missing numeric Status");
    if (result.Status !== 0) throw new Error(`DNS status ${result.Status}`);
    if (result.Answer !== undefined && !Array.isArray(result.Answer)) throw new Error("invalid DNS response: Answer is not an array");
    return (result.Answer ?? []).filter((answer) => answer?.type === recordType && typeof answer.data === "string" && isIP(answer.data) === version).map((answer) => answer.data);
  });
}
