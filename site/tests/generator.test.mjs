import assert from "node:assert/strict";
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

test("generated schemas retain immutable repository bytes and latest is absent", async () => {
  const commit = "0123456789abcdef0123456789abcdef01234567";
  await generateSite({ repoRoot, siteRoot, sourceCommit: commit });
  for (const version of ["0.0.0", "1.0.0"]) {
    const source = await readFile(path.join(repoRoot, "schema", "releases", `v${version}`, "cueson.schema.json"));
    const generated = await readFile(path.join(siteRoot, "public", "schema", `v${version}`, "cueson.schema.json"));
    assert.deepEqual(generated, source);
  }
  await assert.rejects(readFile(path.join(siteRoot, "public", "schema", "latest", "cueson.schema.json")), { code: "ENOENT" });
  await assert.rejects(readFile(path.join(siteRoot, "public", "schema", "v1.1.0", "cueson.schema.json")), { code: "ENOENT" });
});

test("candidate format navigation derives from authored routes while public releases stay unchanged", async () => {
  const contentMap = await loadContentMap(siteRoot);
  await generateSite({ repoRoot, siteRoot });
  const formatDocs = [...contentMap.documents].sort((a, b) => a.order - b.order)
    .filter((document) => document.slug.length === 2 && document.slug[0] === "formats");
  assert.ok(formatDocs.some((document) => document.source === "docs/formats/ass-ssa.md" && document.slug[1] === "ass-ssa"));
  const meta = JSON.parse(await readFile(path.join(siteRoot, "content", "generated", "formats", "meta.json"), "utf8"));
  assert.deepEqual(meta.pages, formatDocs.map((document) => document.slug[1]));
  const source = "See [ASS/SSA](formats/ass-ssa.md#encoding-timing-and-source-authority).";
  assert.equal(rewriteMarkdownLinks(source, "docs/architecture.md", contentMap), "See [ASS/SSA](/docs/formats/ass-ssa/#encoding-timing-and-source-authority).");
  assert.deepEqual(contentMap.schemas.map((schema) => schema.version), ["0.0.0", "1.0.0"]);
  assert.ok(contentMap.downloads.every((download) => download.url.includes("/download/v1.0.0/")));
  const rendered = await readFile(path.join(siteRoot, "content", "generated", "formats", "ass-ssa.mdx"), "utf8");
  assert.match(rendered, /unpublished 1\.1\.0 candidate/);
  assert.doesNotMatch(rendered, /Unavailable for scripted formats|conversion and the complete later CLI freeze remain deferred/);
  const guide = await readFile(path.join(siteRoot, "public", "guides", "media-formats", "index.html"), "utf8");
  for (const dialect of ["ASS", "SSA"]) {
    const row = guide.match(new RegExp(`<div class="format-name">${dialect}</div>[\\s\\S]*?</tr>`));
    assert.ok(row, `${dialect} support row is absent`);
    assert.match(row[0], /Candidate 1\.1\.0/);
    assert.doesNotMatch(row[0], /status-future">Future/);
  }
});

test("deployment workflow is manual-only and binds a full main revision", async () => {
  const workflow = await readFile(path.join(repoRoot, ".github", "workflows", "site-deploy.yml"), "utf8");
  assert.match(workflow, /^on:\s*\n\s*workflow_dispatch:/m);
  assert.doesNotMatch(workflow, /pull_request:|\n\s*push:/);
  assert.ok(workflow.includes("^[0-9a-f]{40}$"));
  assert.match(workflow, /git merge-base --is-ancestor/);
  assert.match(workflow, /wrangler deploy/);
  assert.match(workflow, /verify:cloudflare -- --phase before/);
  assert.match(workflow, /verify:cloudflare -- --phase after/);
  assert.match(workflow, /verify:production -- --expected-commit "\$REVISION"/);
  assert.ok(workflow.indexOf("--phase before") < workflow.indexOf("wrangler deploy"));
  assert.ok(workflow.indexOf("wrangler deploy") < workflow.indexOf("--phase after"));
  assert.ok(workflow.indexOf("wrangler deploy") < workflow.indexOf("verify:production"));
});
