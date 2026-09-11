import { readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const apiRoot = "https://api.cloudflare.com/client/v4/";
const workerName = "cueson-site";
const targetHosts = ["cueson.io", "www.cueson.io"];

function ordered(items, key) {
  return [...items].sort((left, right) => key(left).localeCompare(key(right)));
}

function same(label, expected, actual) {
  if (JSON.stringify(expected) !== JSON.stringify(actual)) throw new Error(`${label}: unrelated Cloudflare state changed during deployment`);
}

function targetDomain(snapshot, hostname) {
  return snapshot.domains.filter((domain) => domain.hostname === hostname);
}

export function validatePreMutation(snapshot) {
  if (snapshot.zone.id !== snapshot.zone_id || snapshot.zone.name !== "cueson.io" || snapshot.zone.account_id !== snapshot.account_id) throw new Error("Cloudflare zone identity does not match the declared production target");
  if (snapshot.zone.status !== "active" || snapshot.zone.paused) throw new Error("cueson.io Cloudflare zone is not active and unpaused");

  for (const hostname of targetHosts) {
    const domains = targetDomain(snapshot, hostname);
    if (domains.length > 1) throw new Error(`${hostname}: multiple Custom Domain bindings exist`);
    if (domains.some((domain) => domain.service !== workerName || domain.zone_id !== snapshot.zone_id)) throw new Error(`${hostname}: Custom Domain belongs to an unrelated resource`);
    const records = snapshot.dns_records.filter((record) => record.name === hostname);
    if (records.length > 0 && domains.length === 0) throw new Error(`${hostname}: DNS records exist without the declared Worker Custom Domain`);
  }

  for (const ruleset of snapshot.redirect_rulesets) {
    if ((ruleset.rules ?? []).some((rule) => targetHosts.some((hostname) => JSON.stringify(rule).toLowerCase().includes(hostname)))) throw new Error(`${ruleset.id}: redirect rule references a Cueson production hostname`);
  }
}

export function validatePostMutation(before, after) {
  validatePreMutation(after);
  if (!after.workers.some((worker) => worker.id === workerName)) throw new Error(`${workerName}: Worker was not present after deployment`);
  for (const hostname of targetHosts) {
    const domains = targetDomain(after, hostname);
    if (domains.length !== 1 || domains[0].service !== workerName) throw new Error(`${hostname}: expected one ${workerName} Custom Domain after deployment`);
    if (!after.dns_records.some((record) => record.name === hostname)) throw new Error(`${hostname}: provider-managed DNS record was not visible after deployment`);
  }

  same("zone", before.zone, after.zone);
  same("unrelated DNS records", before.dns_records.filter((record) => !targetHosts.includes(record.name)), after.dns_records.filter((record) => !targetHosts.includes(record.name)));
  same("unrelated Workers", before.workers.filter((worker) => worker.id !== workerName), after.workers.filter((worker) => worker.id !== workerName));
  same("unrelated Custom Domains", before.domains.filter((domain) => !targetHosts.includes(domain.hostname)), after.domains.filter((domain) => !targetHosts.includes(domain.hostname)));
  same("redirect rulesets", before.redirect_rulesets, after.redirect_rulesets);
}

function cloudflareUrl(pathname) {
  if (pathname instanceof URL) return pathname;
  return new URL(pathname.replace(/^\/+/, ""), apiRoot);
}

async function requestJson(pathname, token, fetcher) {
  const url = cloudflareUrl(pathname);
  const response = await fetcher(url, {
    headers: { Authorization: `Bearer ${token}`, Accept: "application/json" },
    signal: AbortSignal.timeout(20_000),
  });
  const payload = await response.json().catch(() => ({}));
  if (!response.ok || payload.success !== true) {
    const messages = [...(payload.errors ?? []), ...(payload.messages ?? [])].map((item) => item.message).filter(Boolean).join("; ");
    throw new Error(`Cloudflare API ${url.pathname} failed (${response.status})${messages ? `: ${messages}` : ""}`);
  }
  return payload;
}

async function listAll(pathname, token, fetcher) {
  const items = [];
  const perPage = 50;
  for (let page = 1; ; page += 1) {
    const url = cloudflareUrl(pathname);
    url.searchParams.set("page", String(page));
    url.searchParams.set("per_page", String(perPage));
    const payload = await requestJson(url, token, fetcher);
    if (!Array.isArray(payload.result)) throw new Error(`Cloudflare API ${pathname} did not return a list`);
    items.push(...payload.result);
    const resultInfo = payload.result_info;
    if (resultInfo?.total_pages && page >= resultInfo.total_pages) return items;
    if (Number.isInteger(resultInfo?.total_count) && items.length >= resultInfo.total_count) return items;
    if (payload.result.length < perPage) return items;
  }
}

export async function collectState({ accountId, zoneId, token, fetcher = fetch }) {
  const [zonePayload, dnsRecords, workers, domains, rulesets] = await Promise.all([
    requestJson(`/zones/${zoneId}`, token, fetcher),
    listAll(`/zones/${zoneId}/dns_records`, token, fetcher),
    listAll(`/accounts/${accountId}/workers/scripts`, token, fetcher),
    listAll(`/accounts/${accountId}/workers/domains?zone_id=${zoneId}`, token, fetcher),
    listAll(`/zones/${zoneId}/rulesets`, token, fetcher),
  ]);
  const redirectRulesets = [];
  for (const ruleset of rulesets.filter((item) => item.phase?.includes("redirect"))) redirectRulesets.push((await requestJson(`/zones/${zoneId}/rulesets/${ruleset.id}`, token, fetcher)).result);
  const zone = zonePayload.result;
  return {
    account_id: accountId,
    zone_id: zoneId,
    zone: { id: zone.id, name: zone.name, status: zone.status, paused: zone.paused, type: zone.type, account_id: zone.account?.id },
    dns_records: ordered(dnsRecords.map((record) => ({ id: record.id, name: record.name.toLowerCase(), type: record.type, content: record.content, proxied: record.proxied ?? null, ttl: record.ttl, comment: record.comment ?? "", tags: [...(record.tags ?? [])].sort() })), (record) => record.id),
    workers: ordered(workers.map((worker) => ({ id: worker.id, etag: worker.etag ?? "", modified_on: worker.modified_on ?? "" })), (worker) => worker.id),
    domains: ordered(domains.map((domain) => ({ id: domain.id, cert_id: domain.cert_id ?? "", hostname: domain.hostname.toLowerCase(), service: domain.service, environment: domain.environment ?? "", zone_id: domain.zone_id })), (domain) => `${domain.hostname}:${domain.id}`),
    redirect_rulesets: ordered(redirectRulesets, (ruleset) => ruleset.id),
  };
}

function parseArguments(argv) {
  const values = new Map();
  for (let index = 0; index < argv.length; index += 2) {
    const key = argv[index];
    const value = argv[index + 1];
    if (!key?.startsWith("--") || !value || value.startsWith("--")) throw new Error(`invalid argument: ${key ?? ""}`);
    values.set(key, value);
  }
  const phase = values.get("--phase");
  const snapshot = values.get("--snapshot");
  if (!(["before", "after"].includes(phase)) || !snapshot) throw new Error("usage: --phase before|after --snapshot PATH");
  return { phase, snapshot: path.resolve(snapshot) };
}

async function main() {
  const { phase, snapshot } = parseArguments(process.argv.slice(2));
  const accountId = process.env.CLOUDFLARE_ACCOUNT_ID;
  const zoneId = process.env.CLOUDFLARE_ZONE_ID;
  const token = process.env.CLOUDFLARE_API_TOKEN;
  if (!accountId || !zoneId || !token) throw new Error("CLOUDFLARE_ACCOUNT_ID, CLOUDFLARE_ZONE_ID, and CLOUDFLARE_API_TOKEN are required");
  const current = await collectState({ accountId, zoneId, token });
  if (phase === "before") {
    validatePreMutation(current);
    await writeFile(snapshot, `${JSON.stringify(current, null, 2)}\n`, { encoding: "utf8", mode: 0o600 });
    process.stdout.write(`Cloudflare preflight passed with ${current.dns_records.length} DNS records, ${current.workers.length} Workers, and ${current.domains.length} Custom Domains\n`);
    return;
  }
  const before = JSON.parse(await readFile(snapshot, "utf8"));
  validatePostMutation(before, current);
  process.stdout.write(`Cloudflare read-back passed for ${workerName} on ${targetHosts.join(" and ")}\n`);
}

if (path.resolve(process.argv[1] ?? "") === fileURLToPath(import.meta.url)) await main();
