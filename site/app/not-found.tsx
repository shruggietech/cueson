import Link from "next/link";

export default function NotFound() {
  return (
    <main className="not-found">
      <p className="eyebrow">404</p>
      <h1>That cue is not on the timeline.</h1>
      <p>The requested page does not exist.</p>
      <Link className="button" href="/docs/">Open the documentation</Link>
    </main>
  );
}
