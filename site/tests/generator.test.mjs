import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readFile, writeFile } from "node:fs/promises";
import path from "node:path";
import test, { after } from "node:test";

import {
  generateSite,
  loadContentMap,
  rewriteMarkdownLinks,
  safeTarget,
} from "../scripts/generate.mjs";

const siteRoot = path.resolve(import.meta.dirname, "..");
const repoRoot = path.resolve(siteRoot, "..");
const release = {
  version: "1.1.0",
  tag: "v1.1.0",
  url: "https://github.com/shruggietech/cueson/releases/tag/v1.1.0",
};
const downloads = [
  ["Windows x86-64", "cueson_1.1.0_windows_amd64.zip"],
  ["Windows Arm64", "cueson_1.1.0_windows_arm64.zip"],
  ["macOS Intel", "cueson_1.1.0_darwin_amd64.tar.gz"],
  ["macOS Apple silicon", "cueson_1.1.0_darwin_arm64.tar.gz"],
  ["Linux x86-64", "cueson_1.1.0_linux_amd64.tar.gz"],
  ["Linux Arm64", "cueson_1.1.0_linux_arm64.tar.gz"],
  ["SHA-256 checksums", "cueson_1.1.0_checksums.txt"],
].map(([name, filename]) => ({
  name,
  url: `https://github.com/shruggietech/cueson/releases/download/v1.1.0/${filename}`,
}));
const schemas = [
  { version: "0.0.0", byte_length: 22_263, sha256: "d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975" },
  { version: "1.0.0", byte_length: 61_445, sha256: "1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541" },
  { version: "1.1.0", byte_length: 185_641, sha256: "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7" },
];
const sha256 = (bytes) => createHash("sha256").update(bytes).digest("hex");

after(async () => {
  await generateSite({ repoRoot, siteRoot });
});

test("safeTarget confines generated paths to their declared root", () => {
  const root = path.join(siteRoot, "public");
  assert.equal(safeTarget(root, "schema/v1.0.0/cueson.schema.json"), path.join(root, "schema/v1.0.0/cueson.schema.json"));
  assert.throws(() => safeTarget(root, "../outside"), /unsafe generated path/);
  assert.throws(() => safeTarget(root, "C:/outside"), /unsafe generated path/);
});

test("repository documentation links become stable public routes", async () => {
  const contentMap = await loadContentMap(siteRoot);
  const source = "See [schema](schema.md#root-object), [WebVTT](formats/webvtt.md), [source](../internal/schema/schema.go), and [verifier](../scripts/brand-verify/).";
  const actual = rewriteMarkdownLinks(source, "docs/architecture.md", contentMap);
  assert.equal(actual, "See [schema](/docs/schema/#root-object), [WebVTT](/docs/formats/webvtt/), [source](https://github.com/shruggietech/cueson/blob/main/internal/schema/schema.go), and [verifier](https://github.com/shruggietech/cueson/tree/main/scripts/brand-verify/)." );
});

test("generation is deterministic and check mode detects drift", async () => {
  const commit = "0123456789abcdef0123456789abcdef01234567";
  const first = await generateSite({ repoRoot, siteRoot, sourceCommit: commit });
  const second = await generateSite({ repoRoot, siteRoot, sourceCommit: commit, check: true });
  assert.deepEqual(second, first);

  const deploymentPath = path.join(siteRoot, "public", "deployment.json");
  const original = await readFile(deploymentPath);
  await writeFile(deploymentPath, Buffer.concat([original, Buffer.from("\n")]));
  await assert.rejects(generateSite({ repoRoot, siteRoot, sourceCommit: commit, check: true }), /generated output drift/);
  await generateSite({ repoRoot, siteRoot, sourceCommit: commit });
});

test("generated schemas retain all three immutable repository identities and latest is absent", async () => {
  const commit = "0123456789abcdef0123456789abcdef01234567";
  await generateSite({ repoRoot, siteRoot, sourceCommit: commit });
  for (const schema of schemas) {
    const { version } = schema;
    const source = await readFile(path.join(repoRoot, "schema", "releases", `v${version}`, "cueson.schema.json"));
    const generated = await readFile(path.join(siteRoot, "public", "schema", `v${version}`, "cueson.schema.json"));
    assert.deepEqual(generated, source);
    assert.equal(generated.length, schema.byte_length);
    assert.equal(sha256(generated), schema.sha256);
  }
  await assert.rejects(readFile(path.join(siteRoot, "public", "schema", "latest", "cueson.schema.json")), { code: "ENOENT" });
  await assert.rejects(readFile(path.join(siteRoot, "public", "schema", "v9.9.9", "cueson.schema.json")), { code: "ENOENT" });
});

test("published v1.1 release, downloads, and format navigation derive from maintained authority", async () => {
  const contentMap = await loadContentMap(siteRoot);
  await generateSite({ repoRoot, siteRoot });
  assert.deepEqual(contentMap.release, release);
  assert.deepEqual(contentMap.downloads, downloads);
  assert.equal(new Set(contentMap.downloads.map((download) => download.name)).size, 7);
  assert.equal(new Set(contentMap.downloads.map((download) => download.url)).size, 7);
  assert.deepEqual(contentMap.documents
    .filter((document) => document.slug[0] === "releases")
    .map((document) => document.slug[1]), ["v1.1.0", "v1.0.0", "v0.0.0"]);
  const formatDocs = [...contentMap.documents].sort((a, b) => a.order - b.order)
    .filter((document) => document.slug.length === 2 && document.slug[0] === "formats");
  const scriptedDocument = formatDocs.find((document) => document.source === "docs/formats/ass-ssa.md" && document.slug[1] === "ass-ssa");
  assert.ok(scriptedDocument);
  assert.match(scriptedDocument.description, /stable/i);
  assert.doesNotMatch(scriptedDocument.description, /candidate|unpublished/i);
  const meta = JSON.parse(await readFile(path.join(siteRoot, "content", "generated", "formats", "meta.json"), "utf8"));
  assert.deepEqual(meta.pages, formatDocs.map((document) => document.slug[1]));
  const source = "See [ASS/SSA](formats/ass-ssa.md#encoding-timing-and-source-authority).";
  assert.equal(rewriteMarkdownLinks(source, "docs/architecture.md", contentMap), "See [ASS/SSA](/docs/formats/ass-ssa/#encoding-timing-and-source-authority).");
  assert.deepEqual(contentMap.schemas.map(({ version, byte_length, sha256: digest }) => ({ version, byte_length, sha256: digest })), schemas);
  const rendered = await readFile(path.join(siteRoot, "content", "generated", "formats", "ass-ssa.mdx"), "utf8");
  assert.match(rendered, /1\.1\.0/);
  assert.doesNotMatch(rendered, /unpublished 1\.1\.0 candidate/i);
  assert.doesNotMatch(rendered, /Unavailable for scripted formats|conversion and the complete later CLI freeze remain deferred/);
  const currentGeneratedDocuments = await Promise.all([
    readFile(path.join(siteRoot, "content", "generated", "contributing.mdx"), "utf8"),
    readFile(path.join(siteRoot, "content", "generated", "security.mdx"), "utf8"),
  ]);
  assert.doesNotMatch(currentGeneratedDocuments.join("\n"), /published v1\.0\.0 release is the current stable line|unpublished exact 1\.1\.0|Current source targets the unpublished exact 1\.1\.0/i);
  const releasePage = await readFile(path.join(siteRoot, "content", "generated", "releases", "v1.1.0.mdx"), "utf8");
  assert.match(releasePage, /Cueson v1\.1\.0/);
  assert.match(releasePage, /independently verified/i);
  const guide = await readFile(path.join(siteRoot, "public", "guides", "media-formats", "index.html"), "utf8");
  for (const dialect of ["ASS", "SSA"]) {
    const row = guide.match(new RegExp(`<div class="format-name">${dialect}</div>[\\s\\S]*?</tr>`));
    assert.ok(row, `${dialect} support row is absent`);
    assert.match(row[0], /Stable 1\.1\.0/);
    assert.doesNotMatch(row[0], /status-future">Future/);
  }
});

test("release asset filenames must match the declared release version", async () => {
  const contentMapPath = path.join(siteRoot, "content-map.json");
  const original = await readFile(contentMapPath, "utf8");
  const changed = JSON.parse(original);
  changed.release = {
    version: "1.2.0",
    tag: "v1.2.0",
    url: "https://github.com/shruggietech/cueson/releases/tag/v1.2.0",
  };
  changed.downloads = changed.downloads.map((download) => ({
    ...download,
    url: download.url.replace("/download/v1.1.0/", "/download/v1.2.0/"),
  }));
  try {
    await writeFile(contentMapPath, `${JSON.stringify(changed, null, 2)}\n`);
    await assert.rejects(loadContentMap(siteRoot), /release download inventory mismatch/);
  } finally {
    await writeFile(contentMapPath, original);
  }
});

test("generated deployment declares exactly 22 HTML routes, three schemas, and seven downloads", async () => {
  const commit = "0123456789abcdef0123456789abcdef01234567";
  await generateSite({ repoRoot, siteRoot, sourceCommit: commit });
  const deployment = JSON.parse(await readFile(path.join(siteRoot, "public", "deployment.json"), "utf8"));
  const manifest = JSON.parse(await readFile(path.join(siteRoot, "public", "content-manifest.json"), "utf8"));
  const htmlRoutes = deployment.routes.filter((route) => route === "/" || route.endsWith("/"));
  const schemaRoutes = deployment.routes.filter((route) => route.endsWith("/cueson.schema.json"));
  assert.equal(deployment.routes.length, 25);
  assert.equal(new Set(deployment.routes).size, 25);
  assert.equal(htmlRoutes.length, 22);
  assert.deepEqual(schemaRoutes, schemas.map((schema) => `/schema/v${schema.version}/cueson.schema.json`));
  assert.deepEqual(deployment.release, release);
  assert.deepEqual(deployment.downloads, downloads);
  assert.deepEqual(manifest.release, release);
  assert.deepEqual(manifest.downloads, downloads);
  assert.deepEqual(manifest.schemas.map(({ version, public_path, byte_length, sha256: digest }) => ({ version, public_path, byte_length, sha256: digest })), schemas.map((schema) => ({
    version: schema.version,
    public_path: `/schema/v${schema.version}/cueson.schema.json`,
    byte_length: schema.byte_length,
    sha256: schema.sha256,
  })));
});

test("deployment workflow is manual-only and requires the exact current main revision", async () => {
  const workflow = await readFile(path.join(repoRoot, ".github", "workflows", "site-deploy.yml"), "utf8");
  assert.match(workflow, /^on:\s*\n\s*workflow_dispatch:/m);
  assert.doesNotMatch(workflow, /pull_request:|\n\s*push:/);
  assert.ok(workflow.includes("^[0-9a-f]{40}$"));
  assert.match(workflow, /git fetch origin main/);
  assert.equal(workflow.match(/test\s+"\$REVISION"\s+=\s+"\$\(git rev-parse origin\/main\)"/g)?.length, 3);
  assert.doesNotMatch(workflow, /git merge-base --is-ancestor/);
  assert.match(workflow, /wrangler deploy/);
  assert.match(workflow, /verify:cloudflare -- --phase before/);
  assert.match(workflow, /verify:cloudflare -- --phase after/);
  assert.match(workflow, /verify:production -- --expected-commit "\$REVISION"/);
  const bindIndex = workflow.indexOf("Bind checkout to current main");
  const installIndex = workflow.indexOf("corepack pnpm install");
  const proofConfirmIndex = workflow.indexOf("Confirm exact current main revision before proof");
  const proofIndex = workflow.indexOf("corepack pnpm test");
  const preflightIndex = workflow.indexOf("--phase before");
  const reconfirmIndex = workflow.indexOf("Reconfirm exact current main revision before deployment");
  const deployIndex = workflow.indexOf("wrangler deploy");
  assert.ok(bindIndex < installIndex);
  const bindBlock = workflow.slice(bindIndex, installIndex);
  assert.match(bindBlock, /git fetch origin main/);
  assert.match(bindBlock, /test\s+"\$REVISION"\s+=\s+"\$\(git rev-parse origin\/main\)"/);
  assert.ok(installIndex < proofConfirmIndex);
  assert.ok(proofConfirmIndex < proofIndex);
  assert.ok(proofIndex < preflightIndex);
  assert.ok(preflightIndex < reconfirmIndex);
  assert.ok(reconfirmIndex < deployIndex);
  assert.ok(deployIndex < workflow.indexOf("--phase after"));
  assert.ok(deployIndex < workflow.indexOf("verify:production"));
});
