import { RootProvider } from "fumadocs-ui/provider/next";
import type { Metadata, Viewport } from "next";
import type { ReactNode } from "react";
import { site } from "@/lib/site";
import "./global.css";

export const metadata: Metadata = {
  metadataBase: new URL(site.origin),
  title: { default: "Cueson | Universal captions and subtitles", template: "%s | Cueson" },
  description: site.description,
  icons: {
    icon: [{ url: "/assets/favicons/favicon.ico" }, { url: "/assets/favicons/favicon.svg", type: "image/svg+xml" }],
    apple: "/assets/favicons/apple-touch-icon.png",
  },
};

export const viewport: Viewport = { width: "device-width", initialScale: 1, colorScheme: "dark" };

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <body>
        <RootProvider search={{ enabled: false }} theme={{ attribute: "class", defaultTheme: "dark", enableSystem: false, forcedTheme: "dark" }}>
          {children}
        </RootProvider>
      </body>
    </html>
  );
}
