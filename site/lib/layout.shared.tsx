import type { BaseLayoutProps } from "fumadocs-ui/layouts/shared";
import { Code2, Library, PackageOpen } from "lucide-react";

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      title: (
        <span className="brand-lockup">
          {/* eslint-disable-next-line @next/next/no-img-element -- byte-approved brand SVG must be served directly */}
          <img src="/assets/logos/cueson-horizontal-color.svg" alt="Cueson" />
        </span>
      ),
      url: "/",
      transparentMode: "top",
    },
    links: [
      { text: "Docs", url: "/docs/", icon: <Library aria-hidden="true" /> },
      { text: "Media guide", url: "/guides/media-formats/", icon: <PackageOpen aria-hidden="true" /> },
      { type: "main", text: "GitHub", url: "https://github.com/shruggietech/cueson", external: true, icon: <Code2 aria-hidden="true" /> },
    ],
    searchToggle: { enabled: false },
    themeSwitch: { enabled: false },
  };
}
