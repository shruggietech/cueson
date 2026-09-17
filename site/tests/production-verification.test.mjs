import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import test from "node:test";

import {
  assertContentManifestDigest,
  assertDeploymentMatches,
} from "../scripts/verify-production.mjs";

const manifestBytes = Buffer.from('{"release":{"version":"1.1.0"}}\n', "utf8");
const release = {
  version: "1.1.0",
  tag: "v1.1.0",
  url: "https://github.com/shruggietech/cueson/releases/tag/v1.1.0",
};
const downloads = [
  "cueson_1.1.0_windows_amd64.zip",
  "cueson_1.1.0_windows_arm64.zip",
  "cueson_1.1.0_darwin_amd64.tar.gz",
  "cueson_1.1.0_darwin_arm64.tar.gz",
  "cueson_1.1.0_linux_amd64.tar.gz",
  "cueson_1.1.0_linux_arm64.tar.gz",
  "cueson_1.1.0_checksums.txt",
].map((filename) => ({
  name: filename,
  url: `https://github.com/shruggietech/cueson/releases/download/v1.1.0/${filename}`,
}));
const schemas = [
  ["0.0.0", 22_263, "d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975"],
  ["1.0.0", 61_445, "1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541"],
  ["1.1.0", 185_641, "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7"],
].map(([version, byte_length, sha256]) => ({
  version,
  public_path: `/schema/v${version}/cueson.schema.json`,
  byte_length,
  sha256,
}));
const localDeployment = {
  commit: "0123456789abcdef0123456789abcdef01234567",
  release,
  content_manifest_sha256: createHash("sha256").update(manifestBytes).digest("hex"),
  routes: [
    "/",
    "/docs/",
    "/docs/releases/v1.1.0/",
    "/schema/v0.0.0/cueson.schema.json",
    "/schema/v1.0.0/cueson.schema.json",
    "/schema/v1.1.0/cueson.schema.json",
  ],
  downloads,
  schemas,
};

const copy = (value) => structuredClone(value);

test("public deployment metadata must exactly match the locally reviewed authority", () => {
  assert.doesNotThrow(() => assertDeploymentMatches(localDeployment, copy(localDeployment)));

  const cases = [
    ["missing route", (remote) => remote.routes.pop()],
    ["extra download", (remote) => remote.downloads.push({ name: "unexpected", url: "https://example.invalid/unexpected" })],
    ["reordered inventory", (remote) => remote.routes.reverse()],
    ["changed release", (remote) => { remote.release.version = "1.0.0"; }],
    ["missing schema", (remote) => remote.schemas.shift()],
  ];
  for (const [label, mutate] of cases) {
    const remote = copy(localDeployment);
    mutate(remote);
    assert.throws(() => assertDeploymentMatches(localDeployment, remote), /deployment\.json/i, label);
  }
});

test("public content manifest bytes must match the locally reviewed digest", () => {
  assert.doesNotThrow(() => assertContentManifestDigest(localDeployment, manifestBytes));
  assert.throws(
    () => assertContentManifestDigest(localDeployment, Buffer.concat([manifestBytes, Buffer.from(" ")])),
    /content-manifest\.json.*digest/i,
  );
});
