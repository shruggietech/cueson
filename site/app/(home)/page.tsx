import { ArrowRight, CheckCircle2, FileJson2, Languages, ShieldCheck } from "lucide-react";
import Link from "next/link";
import { pageMetadata, site } from "@/lib/site";
import contentMap from "@/content-map.json";

export const metadata = pageMetadata("Universal captions and subtitles", site.description, "/");

const install = "go install github.com/shruggietech/cueson/cmd/cueson@v1.0.0";

export default function HomePage() {
  return (
    <main>
      <section className="hero">
        <div className="hero-copy">
          <p className="eyebrow">Cueson v1.0.0</p>
          <h1>Universal captions and subtitles.</h1>
          <p>A lossless, structured interchange layer for subtitle and caption content. Parse, preserve, inspect, validate, render, and convert SubRip and WebVTT through one stable Cue JSON surface.</p>
          <div className="hero-actions">
            <Link className="button" href="/docs/">Read the docs <ArrowRight size={18} aria-hidden="true" /></Link>
            <a className="button secondary" href={site.release}>Download v1.0.0</a>
          </div>
          <code className="install" tabIndex={0} aria-label="Go installation command">{install}</code>
        </div>
        <div className="cue-card" aria-label="Caption timeline illustration">
          <p className="eyebrow">Cue JSON</p><div className="cue-rail" aria-hidden="true" />
          <div className="cue"><small>00:00:01.250 → 00:00:03.800</small><strong>The words stay addressable.</strong></div>
          <div className="cue"><small>source · preserved · sha256</small><strong>The original carrier stays intact.</strong></div>
          <div className="cue"><small>SRT ⇄ VTT · explicit loss report</small><strong>Conversion stays honest.</strong></div>
        </div>
      </section>
      <section className="section" aria-labelledby="capabilities">
        <p className="eyebrow">Stable v1 contract</p><h2 id="capabilities">One precise surface for caption content.</h2>
        <p className="section-lead">Cueson keeps source evidence separate from derived structure, giving tools a predictable model without discarding what arrived.</p>
        <div className="feature-grid">
          <article className="feature"><FileJson2 aria-hidden="true" /><h3>Structured</h3><p>Stable Cue JSON exposes timing, text, native format details, diagnostics, and source provenance.</p></article>
          <article className="feature"><ShieldCheck aria-hidden="true" /><h3>Lossless</h3><p>Restore original source bytes exactly while deterministic rendering remains a separate, explicit operation.</p></article>
          <article className="feature"><Languages aria-hidden="true" /><h3>Format-literate</h3><p>Convert SubRip and WebVTT with complete loss reporting and strict rejection when fidelity matters.</p></article>
        </div>
      </section>
      <section className="section" aria-labelledby="quick-start">
        <p className="eyebrow">First command</p><h2 id="quick-start">Turn subtitles into Cue JSON.</h2>
        <p className="section-lead">Encode a SubRip or WebVTT file into the stable structured contract. The source envelope preserves the original bytes for exact restoration.</p>
        <pre className="command-example" tabIndex={0} aria-label="Cueson encode example"><code>cueson encode --pretty --output document.cueson.json subtitles.srt</code></pre>
        <p><Link href="/docs/cli/">Continue with the complete CLI guide →</Link></p>
      </section>
      <section className="section" aria-labelledby="resources">
        <p className="eyebrow">Use the contract</p><h2 id="resources">Start with the source of truth.</h2>
        <div className="feature-grid">
          <article className="feature"><CheckCircle2 aria-hidden="true" /><h3>Documentation</h3><p>Read the complete CLI, schema, format, conversion, compatibility, and architecture contracts.</p><Link href="/docs/">Open documentation →</Link></article>
          <article className="feature"><FileJson2 aria-hidden="true" /><h3>Versioned schema</h3><p>Resolve the immutable released Cue JSON schema without a mutable latest alias.</p><a href="/schema/v1.0.0/cueson.schema.json">Open v1 schema →</a></article>
          <article className="feature"><Languages aria-hidden="true" /><h3>Media-format guide</h3><p>Explore the broader caption, subtitle, transcript, and timed-text landscape.</p><a href="/guides/media-formats/">Open the guide →</a></article>
        </div>
      </section>
      <section className="section" aria-labelledby="downloads">
        <p className="eyebrow">Native release</p><h2 id="downloads">Download Cueson v1.0.0.</h2>
        <p className="section-lead">Choose the archive for your platform, then verify it with the published SHA-256 checksum manifest.</p>
        <div className="download-grid">
          {contentMap.downloads.map((download) => <a key={download.url} href={download.url}>{download.name} <ArrowRight size={16} aria-hidden="true" /></a>)}
        </div>
      </section>
    </main>
  );
}
