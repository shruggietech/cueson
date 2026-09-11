import { createReadStream } from "node:fs";
import { stat } from "node:fs/promises";
import { createServer } from "node:http";
import path from "node:path";
import process from "node:process";

const rootArgument = process.argv[2] ?? "out";
const port = Number.parseInt(process.argv[3] ?? "4173", 10);
const root = path.resolve(rootArgument);
const mime = new Map([
  [".css", "text/css; charset=utf-8"], [".html", "text/html; charset=utf-8"],
  [".ico", "image/x-icon"], [".js", "text/javascript; charset=utf-8"],
  [".json", "application/json; charset=utf-8"], [".map", "application/json; charset=utf-8"],
  [".png", "image/png"], [".svg", "image/svg+xml"], [".woff2", "font/woff2"],
]);

function targetFor(rawUrl) {
  const pathname = decodeURIComponent(new URL(rawUrl, "http://127.0.0.1").pathname);
  const relative = pathname.replace(/^\/+/, "");
  const candidate = path.resolve(root, relative, pathname.endsWith("/") ? "index.html" : "");
  if (candidate !== root && !candidate.startsWith(`${root}${path.sep}`)) return null;
  return candidate;
}

const server = createServer(async (request, response) => {
  let target = targetFor(request.url ?? "/");
  try {
    if (!target) throw Object.assign(new Error("unsafe path"), { code: "ENOENT" });
    let info = await stat(target);
    if (info.isDirectory()) {
      target = path.join(target, "index.html");
      info = await stat(target);
    } else if (!path.extname(target)) {
      const html = `${target}.html`;
      try { await stat(html); target = html; } catch { /* retain original */ }
    }
    response.writeHead(200, { "Content-Type": mime.get(path.extname(target)) ?? "application/octet-stream" });
    createReadStream(target).pipe(response);
  } catch {
    const fallback = path.join(root, "404.html");
    try {
      await stat(fallback);
      response.writeHead(404, { "Content-Type": "text/html; charset=utf-8" });
      createReadStream(fallback).pipe(response);
    } catch {
      response.writeHead(404, { "Content-Type": "text/plain; charset=utf-8" });
      response.end("Not found\n");
    }
  }
});

server.listen(port, "127.0.0.1", () => process.stdout.write(`serving ${root} on http://127.0.0.1:${port}\n`));
