import assert from "node:assert/strict";
import test from "node:test";
import { handleRequest } from "./index.ts";

function environment(body: BodyInit = "asset", headers: HeadersInit = {}) {
  return {
    ASSETS: {
      fetch: async (request: Request) => new Response(body, { status: new URL(request.url).pathname.includes("latest") ? 404 : 200, headers }),
    },
  };
}

test("www permanently redirects to apex without losing path or query", async () => {
  const response = await handleRequest(new Request("https://www.cueson.io/docs/cli/?mode=full"), environment());
  assert.equal(response.status, 308);
  assert.equal(response.headers.get("location"), "https://cueson.io/docs/cli/?mode=full");
});

test("apex assets receive production security policy", async () => {
  const response = await handleRequest(new Request("https://cueson.io/docs/"), environment("page", { "content-type": "text/html" }));
  assert.equal(await response.text(), "page");
  assert.equal(response.headers.get("x-content-type-options"), "nosniff");
  assert.equal(response.headers.get("referrer-policy"), "strict-origin-when-cross-origin");
  assert.equal(response.headers.get("x-frame-options"), "DENY");
  assert.match(response.headers.get("content-security-policy") ?? "", /frame-ancestors 'none'/);
  assert.match(response.headers.get("strict-transport-security") ?? "", /max-age=/);
  assert.match(response.headers.get("cache-control") ?? "", /max-age=0/);
});

test("schemas retain bytes and receive immutable schema headers", async () => {
  const bytes = new Uint8Array([0, 1, 2, 255]);
  const response = await handleRequest(new Request("https://cueson.io/schema/v1.0.0/cueson.schema.json"), environment(bytes));
  assert.deepEqual(new Uint8Array(await response.arrayBuffer()), bytes);
  assert.equal(response.headers.get("content-type"), "application/schema+json; charset=utf-8");
  assert.equal(response.headers.get("cache-control"), "public, max-age=31536000, immutable");
});

test("deployment metadata is never cached and missing aliases stay missing", async () => {
  const deployment = await handleRequest(new Request("https://cueson.io/deployment.json"), environment("{}"));
  assert.equal(deployment.headers.get("cache-control"), "no-store, max-age=0");
  const latest = await handleRequest(new Request("https://cueson.io/schema/latest/cueson.schema.json"), environment());
  assert.equal(latest.status, 404);
  assert.equal(latest.headers.get("cache-control"), "no-store, max-age=0");
});

test("only successful content-addressed framework assets receive immutable caching", async () => {
  const frameworkAsset = await handleRequest(new Request("https://cueson.io/_next/static/chunks/app-abc123.js"), environment());
  assert.equal(frameworkAsset.headers.get("cache-control"), "public, max-age=31536000, immutable");

  const stableBrandAsset = await handleRequest(new Request("https://cueson.io/assets/logos/cueson-horizontal-color.svg"), environment());
  assert.equal(stableBrandAsset.headers.get("cache-control"), "public, max-age=300, must-revalidate");

  const missingAsset = await handleRequest(new Request("https://cueson.io/assets/latest.svg"), environment());
  assert.equal(missingAsset.status, 404);
  assert.equal(missingAsset.headers.get("cache-control"), "no-store, max-age=0");
});
