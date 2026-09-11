import Link from "next/link";
import type { ReactNode } from "react";
import { site } from "@/lib/site";

export default function HomeLayout({ children }: { children: ReactNode }) {
  return (
    <>
      <header className="site-header">
        <nav className="site-nav" aria-label="Primary">
          <Link href="/" aria-label="Cueson home">
            {/* eslint-disable-next-line @next/next/no-img-element -- byte-approved brand SVG must be served directly */}
            <img src="/assets/logos/cueson-horizontal-color.svg" alt="" />
          </Link>
          <div className="site-links">
            <Link href="/docs/">Documentation</Link>
            <Link href="/guides/media-formats/">Media guide</Link>
            <a href={site.github}>GitHub</a>
            <a className="nav-download" href={site.release}>Download</a>
          </div>
        </nav>
      </header>
      {children}
      <footer className="site-footer">
        <div className="site-footer-inner"><span>Cueson is a ShruggieTech project.</span><a href={site.github}>Apache-2.0 on GitHub</a></div>
      </footer>
    </>
  );
}
