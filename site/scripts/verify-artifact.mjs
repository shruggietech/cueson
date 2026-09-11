import { createHash } from "node:crypto";
import { lstat, readFile, readdir } from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const siteRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
const repoRoot = path.resolve(siteRoot, "..");
const outputRoot = path.join(siteRoot, "out");

function sha256(bytes) { return createHash("sha256").update(bytes).digest("hex"); }

function target(relative) {
  const result = path.resolve(outputRoot, relative);
  if (result !== outputRoot && !result.startsWith(`${outputRoot}${path.sep}`)) throw new Error(`${relative}: unsafe artifact path`);
  return result;
}

async function assertRegular(relative) {
  const info = await lstat(target(relative));
  if (!info.isFile() || info.isSymbolicLink()) throw new Error(`${relative}: expected a regular artifact file`);
}

async function verifyArtifact() {
  const deploymentBytes = await readFile(target("deployment.json"));
  const deployment = JSON.parse(deploymentBytes.toString("utf8"));
  if (!/^[0-9a-f]{40}$/.test(deployment.commit)) throw new Error("deployment.json: invalid commit identity");
  if (process.env.CUESON_SOURCE_COMMIT && deployment.commit !== process.env.CUESON_SOURCE_COMMIT) throw new Error("deployment.json: commit does not match CUESON_SOURCE_COMMIT");
  const manifest = await readFile(target("content-manifest.json"));
  if (sha256(manifest) !== deployment.content_manifest_sha256) throw new Error("content-manifest.json: deployment hash mismatch");

  for (const route of deployment.routes) {
    const relative = route === "/" ? "index.html" : route.endsWith("/") ? `${route.slice(1)}index.html` : route.slice(1);
    await assertRegular(relative);
  }
  for (const [version, expectedHash] of Object.entries(deployment.schemas)) {
    const relative = `schema/v${version}/cueson.schema.json`;
    const artifact = await readFile(target(relative));
    const source = await readFile(path.join(repoRoot, "schema", "releases", `v${version}`, "cueson.schema.json"));
    if (!artifact.equals(source) || sha256(artifact) !== expectedHash) throw new Error(`${relative}: immutable schema identity mismatch`);
  }
  try {
    await lstat(target("schema/latest"));
    throw new Error("schema/latest: mutable alias must not exist");
  } catch (error) {
    if (error.code !== "ENOENT") throw error;
  }
  const entries = await readdir(outputRoot);
  if (!entries.includes("404.html")) throw new Error("404.html: missing static not-found surface");
  process.stdout.write(`verified deployable artifact for ${deployment.commit}: ${deployment.routes.length} declared routes and ${Object.keys(deployment.schemas).length} immutable schemas\n`);
}

await verifyArtifact();
