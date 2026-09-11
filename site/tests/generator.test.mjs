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
  const source = "See [schema](schema.md#root-object), [WebVTT](formats/webvtt.md), and [source](../internal/schema/schema.go).";
  const actual = rewriteMarkdownLinks(source, "docs/architecture.md", contentMap);
  assert.equal(actual, "See [schema](/docs/schema/#root-object), [WebVTT](/docs/formats/webvtt/), and [source](https://github.com/shruggietech/cueson/blob/main/internal/schema/schema.go)." );
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
});

test("deployment workflow is manual-only and binds a full main revision", async () => {
  const workflow = await readFile(path.join(repoRoot, ".github", "workflows", "site-deploy.yml"), "utf8");
  assert.match(workflow, /^on:\s*\n\s*workflow_dispatch:/m);
  assert.doesNotMatch(workflow, /pull_request:|\n\s*push:/);
  assert.ok(workflow.includes("^[0-9a-f]{40}$"));
  assert.match(workflow, /git merge-base --is-ancestor/);
  assert.match(workflow, /wrangler deploy/);
  assert.match(workflow, /verify:production -- --expected-commit "\$REVISION"/);
  assert.ok(workflow.indexOf("wrangler deploy") < workflow.indexOf("verify:production"));
});
