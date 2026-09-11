interface AssetsBinding {
  fetch(request: Request): Promise<Response>;
}

export interface Environment {
  ASSETS: AssetsBinding;
}

const contentSecurityPolicy = [
  "default-src 'self'",
  "base-uri 'self'",
  "connect-src 'self'",
  "font-src 'self'",
  "form-action 'none'",
  "frame-ancestors 'none'",
  "img-src 'self' data:",
  "object-src 'none'",
  "script-src 'self' 'unsafe-inline'",
  "style-src 'self' 'unsafe-inline'",
  "upgrade-insecure-requests",
].join("; ");

function applyPolicy(request: Request, response: Response): Response {
  const url = new URL(request.url);
  const headers = new Headers(response.headers);
  headers.set("Content-Security-Policy", contentSecurityPolicy);
  headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
  headers.set("X-Content-Type-Options", "nosniff");
  headers.set("X-Frame-Options", "DENY");
  if (url.protocol === "https:") headers.set("Strict-Transport-Security", "max-age=31536000; includeSubDomains");

  if (/^\/schema\/v\d+\.\d+\.\d+\/cueson\.schema\.json$/.test(url.pathname) && response.ok) {
    headers.set("Content-Type", "application/schema+json; charset=utf-8");
    headers.set("Cache-Control", "public, max-age=31536000, immutable");
  } else if (url.pathname === "/deployment.json") {
    headers.set("Cache-Control", "no-store, max-age=0");
  } else if ((headers.get("content-type") ?? "").startsWith("text/html")) {
    headers.set("Cache-Control", "public, max-age=0, must-revalidate");
  } else if (url.pathname.startsWith("/_next/static/") || url.pathname.startsWith("/assets/")) {
    headers.set("Cache-Control", "public, max-age=31536000, immutable");
  } else {
    headers.set("Cache-Control", "public, max-age=300, must-revalidate");
  }
  return new Response(response.body, { status: response.status, statusText: response.statusText, headers });
}

export async function handleRequest(request: Request, environment: Environment): Promise<Response> {
  const url = new URL(request.url);
  if (url.hostname.toLowerCase() === "www.cueson.io") {
    url.protocol = "https:";
    url.hostname = "cueson.io";
    url.port = "";
    return applyPolicy(request, Response.redirect(url.toString(), 308));
  }
  return applyPolicy(request, await environment.ASSETS.fetch(request));
}

export default { fetch: handleRequest } satisfies ExportedHandler<Environment>;
