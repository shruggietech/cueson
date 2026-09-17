import AxeBuilder from "@axe-core/playwright";
import { expect, test } from "@playwright/test";
import { createHash } from "node:crypto";
import { expectMetadata, expectNoHorizontalOverflow } from "./helpers";

const routes = [
  "/", "/docs/", "/docs/architecture/", "/docs/cli/", "/docs/schema/", "/docs/formats/subrip/",
  "/docs/formats/webvtt/", "/docs/formats/ass-ssa/", "/docs/conversion/", "/docs/compatibility/", "/docs/brand/", "/docs/security/",
  "/docs/contributing/", "/docs/changelog/", "/docs/releases/v1.1.0/", "/docs/releases/v1.0.0/", "/docs/releases/v0.0.0/",
  "/docs/release-process/", "/docs/release-verification/", "/docs/project-management/", "/docs/project-specification/",
  "/guides/media-formats/",
];
const downloadFiles = [
  "cueson_1.1.0_windows_amd64.zip",
  "cueson_1.1.0_windows_arm64.zip",
  "cueson_1.1.0_darwin_amd64.tar.gz",
  "cueson_1.1.0_darwin_arm64.tar.gz",
  "cueson_1.1.0_linux_amd64.tar.gz",
  "cueson_1.1.0_linux_arm64.tar.gz",
  "cueson_1.1.0_checksums.txt",
];
const schemaIdentities = [
  ["v0.0.0", 22_263, "d15c7fa5227156109dd6be3d39b711aca3503794bb862169dfca96ee80adb975"],
  ["v1.0.0", 61_445, "1aad14567033d7e14d9beb78985e18007aefb5345095370b11b6b887df7ec541"],
  ["v1.1.0", 185_641, "223b61cbcf6337167039268576b2c739564585fff10e6cc6076e47a03526a0f7"],
] as const;

test("landing page makes installation, documentation, and release downloads obvious", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { level: 1 })).toContainText("Universal captions");
  await expect(page.getByRole("link", { name: /Read the docs/i })).toHaveAttribute("href", "/docs/");
  await expect(page.getByRole("link", { name: /Download v1.1.0/i })).toHaveAttribute("href", "https://github.com/shruggietech/cueson/releases/tag/v1.1.0");
  await expect(page.getByText("go install github.com/shruggietech/cueson/cmd/cueson@v1.1.0")).toBeVisible();
  const primaryDownloads = page.locator('section[aria-labelledby="downloads"] a[href^="https://github.com/shruggietech/cueson/releases/download/v1.1.0/"]');
  await expect(primaryDownloads).toHaveCount(7);
  expect(await primaryDownloads.evaluateAll((links) => links.map((link) => new URL((link as HTMLAnchorElement).href).pathname.split("/").at(-1)))).toEqual(downloadFiles);
  const primaryNavigation = page.getByRole("navigation", { name: "Primary" });
  for (const name of ["Documentation", "Media guide", "GitHub", "Download"]) await expect(primaryNavigation.getByRole("link", { name, exact: true })).toBeVisible();
  await expectMetadata(page, "https://cueson.io/");
});

test("standalone media guide retains primary site navigation", async ({ page }) => {
  await page.goto("/guides/media-formats/");
  const primaryNavigation = page.getByRole("navigation", { name: "Primary" });
  await expect(primaryNavigation.getByRole("link", { name: "Home", exact: true })).toBeVisible();
  await expect(primaryNavigation.getByRole("link", { name: "Documentation", exact: true })).toBeVisible();
  await expect(primaryNavigation.getByRole("link", { name: "GitHub", exact: true })).toBeVisible();
});

for (const route of routes) {
  test(`${route} is accessible, responsive, and metadata-complete`, async ({ page }) => {
    const response = await page.goto(route);
    expect(response?.status()).toBe(200);
    await expect(page.locator("html")).toHaveAttribute("lang", "en");
    await expectNoHorizontalOverflow(page);
    await expectMetadata(page, `https://cueson.io${route}`);
    const results = await new AxeBuilder({ page }).analyze();
    expect(results.violations.filter((item) => item.impact === "serious" || item.impact === "critical")).toEqual([]);
  });
}

test("documentation links remain within valid public routes", async ({ page, request }) => {
  test.skip(test.info().project.name !== "desktop-1440", "route graph needs one browser width");
  const fragments = new Set<string>();
  for (const route of routes.filter((item) => item.startsWith("/docs"))) {
    await page.goto(route);
    const links = await page.locator('main a[href^="/"]').evaluateAll((nodes) => nodes.map((node) => (node as HTMLAnchorElement).getAttribute("href") ?? ""));
    for (const link of new Set(links)) {
      const url = new URL(link, "https://cueson.io");
      expect((await request.get(url.pathname)).status(), `${route} -> ${link}`).toBeLessThan(400);
      if (url.hash) fragments.add(`${url.pathname}${url.hash}`);
    }
  }
  for (const target of fragments) {
    const url = new URL(target, "https://cueson.io");
    await page.goto(`${url.pathname}${url.hash}`);
    await expect(page.locator(`[id="${url.hash.slice(1).replaceAll('"', '\\"')}"]`), target).toHaveCount(1);
  }
});

test("indexable page titles are unique", async ({ page }) => {
  test.skip(test.info().project.name !== "desktop-1440", "metadata inventory needs one browser width");
  const titles = [];
  for (const route of routes) {
    await page.goto(route);
    titles.push(await page.title());
  }
  expect(new Set(titles).size).toBe(titles.length);
});

test("keyboard focus is visible", async ({ page }) => {
  await page.goto("/");
  await page.keyboard.press("Tab");
  const focused = page.locator(":focus-visible");
  await expect(focused).toBeVisible();
  expect(await focused.evaluate((node) => getComputedStyle(node).outlineStyle)).not.toBe("none");
});

test("all immutable schema versions have exact bytes and mutable or unknown aliases stay absent", async ({ request }) => {
  for (const [version, bytes, digest] of schemaIdentities) {
    const response = await request.get(`/schema/${version}/cueson.schema.json`);
    expect(response.status()).toBe(200);
    const body = await response.body();
    expect(body.byteLength).toBe(bytes);
    expect(createHash("sha256").update(body).digest("hex")).toBe(digest);
  }
  expect((await request.get("/schema/latest/cueson.schema.json")).status()).toBe(404);
  expect((await request.get("/schema/v9.9.9/cueson.schema.json")).status()).toBe(404);
});

test("published v1.1 release and ASS/SSA guide present stable current content", async ({ page }) => {
  test.skip(test.info().project.name !== "desktop-1440", "sidebar navigation needs one desktop width");
  await page.goto("/docs/formats/ass-ssa/");
  await expect(page.getByRole("heading", { level: 1 })).toContainText("ASS and SSA");
  await expect(page.locator("main")).toContainText("1.1.0");
  await expect(page.locator("main")).not.toContainText(/unpublished 1\.1\.0 candidate/i);
  await expect(page.getByRole("link", { name: "ASS and SSA", exact: true }).first()).toBeVisible();
  await page.goto("/docs/releases/v1.1.0/");
  await expect(page.getByRole("heading", { level: 1 })).toContainText("Cueson v1.1.0");
  await expect(page.locator("main")).toContainText(/independently verified/i);
  await page.goto("/");
  await expect(page.getByRole("link", { name: /Download v1.1.0/i })).toHaveAttribute("href", /releases\/tag\/v1\.1\.0/);
});
