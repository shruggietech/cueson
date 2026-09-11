import { createMDX } from "fumadocs-mdx/next";

const sourceCommit = process.env.CUESON_SOURCE_COMMIT ?? "development";

/** @type {import('next').NextConfig} */
const config = {
  output: "export",
  trailingSlash: true,
  reactStrictMode: true,
  images: {
    unoptimized: true,
  },
  generateBuildId: async () => sourceCommit,
};

export default createMDX()(config);
