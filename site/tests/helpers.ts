import { expect, type Page } from "@playwright/test";

export async function expectNoHorizontalOverflow(page: Page): Promise<void> {
  const dimensions = await page.evaluate(() => ({ viewport: document.documentElement.clientWidth, content: document.documentElement.scrollWidth }));
  expect(dimensions.content).toBeLessThanOrEqual(dimensions.viewport);
}

export async function expectMetadata(page: Page, canonical: string): Promise<void> {
  await expect(page).toHaveTitle(/Cueson/);
  await expect(page.locator('meta[name="description"]')).toHaveAttribute("content", /.+/);
  await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", canonical);
  await expect(page.locator('meta[property="og:title"]')).toHaveAttribute("content", /.+/);
  await expect(page.locator('meta[name="twitter:card"]')).toHaveAttribute("content", "summary_large_image");
  expect(await page.locator('link[rel~="icon"]').count()).toBeGreaterThanOrEqual(1);
}
