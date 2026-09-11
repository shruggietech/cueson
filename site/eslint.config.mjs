import { defineConfig, globalIgnores } from "eslint/config";
import nextVitals from "eslint-config-next/core-web-vitals";
import nextTypescript from "eslint-config-next/typescript";

export default defineConfig([
  ...nextVitals,
  ...nextTypescript,
  globalIgnores([
    ".next/**",
    ".source/**",
    ".wrangler/**",
    ".wrangler-dry-run/**",
    "content/generated/**",
    "out/**",
    "playwright-report/**",
    "public/assets/**",
    "public/guides/**",
    "public/schema/**",
    "test-results/**",
  ]),
]);
