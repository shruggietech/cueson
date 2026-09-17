import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { existsSync, mkdtempSync, readFileSync } from "node:fs";
import os from "node:os";
import path from "node:path";
import process from "node:process";
import test from "node:test";
import { fileURLToPath } from "node:url";

const siteRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
function resolveWindowsPnpmCli() {
  for (const directory of (process.env.Path ?? process.env.PATH ?? "").split(path.delimiter)) {
    const wrapper = path.join(directory, "pnpm.cmd");
    if (!existsSync(wrapper)) continue;
    const match = readFileSync(wrapper, "utf8").match(/"%~dp0([^"]*pnpm\.(?:cjs|mjs))"/i);
    if (match) return path.resolve(path.dirname(wrapper), match[1].replaceAll("\\", path.sep));
  }
  throw new Error("pnpm.cmd with a direct JavaScript CLI target was not found on PATH");
}

function runPackageScript(argumentsList, environment = process.env) {
  const command = process.platform === "win32" ? process.execPath : "corepack";
  const commandArguments = process.platform === "win32" ? [resolveWindowsPnpmCli(), ...argumentsList] : ["pnpm", ...argumentsList];
  return spawnSync(command, commandArguments, {
    cwd: siteRoot,
    encoding: "utf8",
    env: environment,
    shell: false,
    windowsHide: true,
  });
}

function output(result) {
  return `${result.stdout ?? ""}\n${result.stderr ?? ""}\n${result.error?.message ?? ""}`;
}

function cloudflareEnvironment() {
  const environment = { ...process.env };
  delete environment.CLOUDFLARE_ACCOUNT_ID;
  delete environment.CLOUDFLARE_ZONE_ID;
  delete environment.CLOUDFLARE_API_TOKEN;
  return environment;
}

test("Cloudflare package script forwards the workflow's named arguments", () => {
  const snapshot = path.join(mkdtempSync(path.join(os.tmpdir(), "cueson-cloudflare-")), "before.json");
  const result = runPackageScript(["verify:cloudflare", "--phase", "before", "--snapshot", snapshot], cloudflareEnvironment());
  assert.equal(result.error, undefined);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /CLOUDFLARE_ACCOUNT_ID, CLOUDFLARE_ZONE_ID, and CLOUDFLARE_API_TOKEN are required/);
  assert.doesNotMatch(output(result), /invalid argument/);
});

test("Cloudflare package script rejects a standalone separator", () => {
  const snapshot = path.join(mkdtempSync(path.join(os.tmpdir(), "cueson-cloudflare-")), "before.json");
  const result = runPackageScript(["verify:cloudflare", "--", "--phase", "before", "--snapshot", snapshot], cloudflareEnvironment());
  assert.equal(result.error, undefined);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /invalid argument: --/);
});

test("production package script forwards the workflow's named arguments", () => {
  const result = runPackageScript([
    "verify:production",
    "--expected-commit",
    "0000000000000000000000000000000000000000",
    "--origin",
    "http://127.0.0.1:1",
    "--skip-network-identity",
  ]);
  assert.equal(result.error, undefined);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /fetch failed/);
  assert.doesNotMatch(output(result), /invalid argument/);
});

test("production package script rejects a standalone separator", () => {
  const result = runPackageScript([
    "verify:production",
    "--",
    "--expected-commit",
    "0000000000000000000000000000000000000000",
    "--origin",
    "http://127.0.0.1:1",
    "--skip-network-identity",
  ]);
  assert.equal(result.error, undefined);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /invalid argument: --/);
});
