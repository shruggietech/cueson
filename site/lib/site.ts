import type { Metadata } from "next";
import contentMap from "@/content-map.json";

export const site = {
  name: "Cueson",
  description: "Universal captions and subtitles through a lossless, structured interchange layer.",
  origin: "https://cueson.io",
  github: "https://github.com/shruggietech/cueson",
  release: contentMap.release.url,
  version: contentMap.release.version,
  tag: contentMap.release.tag,
} as const;

export function canonical(pathname = "/"): string {
  return new URL(pathname, site.origin).toString();
}

export function pageMetadata(title: string, description: string, pathname: string): Metadata {
  const url = canonical(pathname);
  return {
    title,
    description,
    alternates: { canonical: url },
    openGraph: {
      type: "website",
      siteName: site.name,
      title,
      description,
      url,
      images: [{ url: "/assets/logos/cueson-social-preview-1280.png", width: 1280, height: 640, alt: "Cueson" }],
    },
    twitter: {
      card: "summary_large_image",
      title,
      description,
      images: ["/assets/logos/cueson-social-preview-1280.png"],
    },
  };
}
