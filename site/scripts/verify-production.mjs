import { createHash } from "node:crypto";
import { resolve4, resolve6 } from "node:dns/promises";
import { readFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import tls from "node:tls";
import { fileURLToPath } from "node:url";

const siteRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const args = new Map();
const argumentsList = process.argv.slice(2).filter((argument) => argument !== "--");
for (let index = 0; index < argumentsList.length; index += 1) {
  const key = argumentsList[index];
  if (key === "--skip-network-identity") args.set(key, "true");
  else {
    const value = argumentsList[index + 1];
    if (!key.startsWith("--") || !value || value.startsWith("--")) throw new Error(`invalid argument: ${key}`);
    args.set(key, value);
    index += 1;
  }
}
const origin = new URL(args.get("--origin") ?? "https://cueson.io");
const expectedCommit = args.get("--expected-commit") ?? process.env.CUESON_SOURCE_COMMIT;
const skipNetworkIdentity = args.has("--skip-network-identity");
const sha256 = (bytes) => createHash("sha256").update(bytes).digest("hex");

async function fetchSuccess(pathname, options = {}) {
  const response = await fetch(new URL(pathname, origin), { redirect: "manual", signal: AbortSignal.timeout(20_000), ...options });
  if (!response.ok) throw new Error(`${pathname}: expected success, received ${response.status}`);
  return response;
}

async function verifyTls(hostname) {
  await new Promise((resolve, reject) => {
    const socket = tls.connect({ host: hostname, port: 443, servername: hostname, rejectUnauthorized: true }, () => {
      const certificate = socket.getPeerCertificate();
      if (!certificate.subjectaltname?.split(", ").some((name) => name === `DNS:${hostname}` || (name.startsWith("DNS:*.") && hostname.endsWith(name.slice(5))))) reject(new Error(`${hostname}: certificate does not cover hostname`));
      else resolve();
      socket.end();
    });
    socket.setTimeout(20_000, () => socket.destroy(new Error(`${hostname}: TLS timeout`)));
    socket.once("error", reject);
  });
}

function verifyHtml(route, html) {
  const expectedCanonical = new URL(route, "https://cueson.io").toString();
  const requirements = [
    [/<html[^>]+lang=["']en["']/i, "English language"],
    [/<title>[^<]+<\/title>/i, "title"],
    [/<meta[^>]+name=["']description["'][^>]+content=["'][^"']+["']/i, "description"],
    [new RegExp(`<link[^>]+rel=["']canonical["'][^>]+href=["']${expectedCanonical.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")}["']`, "i"), "canonical URL"],
    [/<meta[^>]+property=["']og:title["']/i, "Open Graph title"],
    [/<meta[^>]+property=["']og:image["']/i, "Open Graph image"],
    [/<meta[^>]+name=["']twitter:card["']/i, "Twitter card"],
    [/<link[^>]+rel=["']icon["']/i, "favicon"],
  ];
  for (const [pattern, label] of requirements) if (!pattern.test(html)) throw new Error(`${route}: missing ${label}`);
}

async function verifyDnsOverHttps(hostname) {
  const response = await fetch(`https://dns.google/resolve?name=${encodeURIComponent(hostname)}&type=A`, { headers: { accept: "application/dns-json" }, signal: AbortSignal.timeout(20_000) });
  if (!response.ok) throw new Error(`${hostname}: DNS-over-HTTPS returned ${response.status}`);
  const result = await response.json();
  if (result.Status !== 0 || !Array.isArray(result.Answer) || result.Answer.length === 0) throw new Error(`${hostname}: DNS-over-HTTPS returned no A answer`);
}

async function verify() {
  const localDeployment = JSON.parse(await readFile(path.join(siteRoot, "public", "deployment.json"), "utf8"));
  const deploymentResponse = await fetchSuccess("/deployment.json");
  if (!/no-store|max-age=0/.test(deploymentResponse.headers.get("cache-control") ?? "")) throw new Error("deployment.json: unsafe cache policy");
  const deployment = await deploymentResponse.json();
  const requiredCommit = expectedCommit ?? localDeployment.commit;
  if (deployment.commit !== requiredCommit) throw new Error(`deployment.json: expected ${requiredCommit}, received ${deployment.commit}`);
  for (const route of deployment.routes.filter((route) => route === "/" || route.endsWith("/"))) {
    const response = await fetchSuccess(route);
    verifyHtml(route, await response.text());
  }

  for (const [version, expectedHash] of Object.entries(localDeployment.schemas)) {
    const route = `/schema/v${version}/cueson.schema.json`;
    const response = await fetchSuccess(route);
    const bytes = Buffer.from(await response.arrayBuffer());
    if (sha256(bytes) !== expectedHash) throw new Error(`${route}: public schema hash mismatch`);
    if (!/application\/(schema\+json|json)/.test(response.headers.get("content-type") ?? "")) throw new Error(`${route}: incorrect content type`);
    if (!/immutable/.test(response.headers.get("cache-control") ?? "")) throw new Error(`${route}: missing immutable cache policy`);
  }
  const latest = await fetch(new URL("/schema/latest/cueson.schema.json", origin), { redirect: "manual", signal: AbortSignal.timeout(20_000) });
  if (latest.status !== 404) throw new Error(`/schema/latest/cueson.schema.json: expected 404, received ${latest.status}`);

  for (const download of deployment.downloads ?? []) {
    const first = await fetch(download.url, { redirect: "manual", signal: AbortSignal.timeout(20_000) });
    if (![301, 302, 303, 307, 308].includes(first.status) || !first.headers.get("location")) throw new Error(`${download.name}: release asset did not produce its expected redirect`);
    const terminal = await fetch(first.headers.get("location"), { redirect: "follow", signal: AbortSignal.timeout(20_000) });
    if (!terminal.ok) throw new Error(`${download.name}: release asset terminal response was ${terminal.status}`);
    await terminal.body?.cancel();
  }

  if (!skipNetworkIdentity && origin.hostname === "cueson.io") {
    const addresses = [...await resolve4("cueson.io"), ...await resolve6("cueson.io").catch(() => [])];
    if (addresses.length === 0) throw new Error("cueson.io: no public DNS addresses");
    const wwwAddresses = [...await resolve4("www.cueson.io"), ...await resolve6("www.cueson.io").catch(() => [])];
    if (wwwAddresses.length === 0) throw new Error("www.cueson.io: no public DNS addresses");
    await verifyDnsOverHttps("cueson.io");
    await verifyDnsOverHttps("www.cueson.io");
    await verifyTls("cueson.io");
    await verifyTls("www.cueson.io");
    const probe = new URL("https://www.cueson.io/docs/schema/?probe=s021");
    const redirect = await fetch(probe, { redirect: "manual", signal: AbortSignal.timeout(20_000) });
    if (redirect.status !== 308 || redirect.headers.get("location") !== "https://cueson.io/docs/schema/?probe=s021") throw new Error("www: redirect does not preserve path and query");
  }
  process.stdout.write(`verified ${origin.origin} at ${deployment.commit} with ${deployment.routes.length} routes and ${Object.keys(deployment.schemas).length} schemas\n`);
}

await verify();
