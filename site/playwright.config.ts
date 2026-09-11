import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./tests",
  testMatch: "site.spec.ts",
  fullyParallel: true,
  forbidOnly: Boolean(process.env.CI),
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "github" : "list",
  use: {
    baseURL: "http://127.0.0.1:4173",
    trace: "retain-on-failure",
  },
  projects: [
    { name: "mobile-360", use: { viewport: { width: 360, height: 800 } } },
    { name: "tablet-768", use: { viewport: { width: 768, height: 1024 } } },
    { name: "desktop-1440", use: { viewport: { width: 1440, height: 1000 } } },
  ],
  webServer: {
    command: "node scripts/serve.mjs out 4173",
    port: 4173,
    reuseExistingServer: !process.env.CI,
  },
});
