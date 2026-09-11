import { createHash } from "node:crypto";
import { spawnSync } from "node:child_process";
import { mkdir, readFile, readdir, rm, writeFile } from "node:fs/promises";
import path from "node:path";
import process from "node:process";
import { fileURLToPath } from "node:url";

const moduleDirectory = path.dirname(fileURLToPath(import.meta.url));
const defaultSiteRoot = path.resolve(moduleDirectory, "..");
const defaultRepoRoot = path.resolve(defaultSiteRoot, "..");
const gitHubBlobRoot = "https://github.com/shruggietech/cueson/blob/main/";
const generatedRoots = ["content/generated", "public/assets", "public/schema", "public/guides"];
const generatedFiles = ["public/content-manifest.json", "public/deployment.json"];

function sha256(bytes) {
  return createHash("sha256").update(bytes).digest("hex");
}

function stableJson(value) {
  return Buffer.from(`${JSON.stringify(value, null, 2)}\n`, "utf8");
}

function toPosix(value) {
  return value.split(path.sep).join("/");
}

function assertRelativePath(value, label) {
  if (typeof value !== "string" || value.length === 0 || path.isAbsolute(value) || /^[a-zA-Z]:[\\/]/.test(value)) {
    throw new Error(`${label}: unsafe generated path`);
  }
  const normalized = path.posix.normalize(value.replaceAll("\\", "/"));
  if (normalized === ".." || normalized.startsWith("../") || normalized.includes("/../")) {
    throw new Error(`${label}: unsafe generated path`);
  }
  return normalized;
}

export function safeTarget(root, relativePath) {
  const normalized = assertRelativePath(relativePath, relativePath);
  const resolvedRoot = path.resolve(root);
  const resolved = path.resolve(resolvedRoot, ...normalized.split("/"));
  if (resolved !== resolvedRoot && !resolved.startsWith(`${resolvedRoot}${path.sep}`)) {
    throw new Error(`${relativePath}: unsafe generated path`);
  }
  return resolved;
}

async function readUtf8(filePath, label) {
  const bytes = await readFile(filePath);
  try {
    return { bytes, text: new TextDecoder("utf-8", { fatal: true }).decode(bytes) };
  } catch {
    throw new Error(`${label}: source is not valid UTF-8`);
  }
}

function docsRoute(slug) {
  return slug.length === 0 ? "/docs/" : `/docs/${slug.join("/")}/`;
}

function splitFragment(target) {
  const index = target.indexOf("#");
  if (index < 0) return { pathname: target, fragment: "" };
  return { pathname: target.slice(0, index), fragment: target.slice(index) };
}

export function rewriteMarkdownLinks(markdown, sourcePath, contentMap) {
  const documents = new Map(contentMap.documents.map((document) => [path.posix.normalize(document.source), docsRoute(document.slug)]));
  const schemas = new Map(contentMap.schemas.map((schema) => [path.posix.normalize(schema.source), `/${schema.public}`]));
  const directFiles = new Map(contentMap.direct_files.map((file) => [path.posix.normalize(file.source), `/${path.posix.dirname(file.public)}/`]));

  return markdown.replace(/\]\(([^)\s]+)(?:\s+(["'][^"']*["']))?\)/g, (match, rawTarget, optionalTitle) => {
    const unwrapped = rawTarget.startsWith("<") && rawTarget.endsWith(">") ? rawTarget.slice(1, -1) : rawTarget;
    if (/^(?:[a-z][a-z0-9+.-]*:|#|\/)/i.test(unwrapped)) return match;

    const { pathname, fragment } = splitFragment(unwrapped);
    const resolved = path.posix.normalize(path.posix.join(path.posix.dirname(sourcePath), pathname));
    let replacement = documents.get(resolved);
    if (!replacement) replacement = schemas.get(resolved);
    if (!replacement) replacement = directFiles.get(resolved);
    if (!replacement) replacement = `${gitHubBlobRoot}${resolved}`;
    const title = optionalTitle ? ` ${optionalTitle}` : "";
    return `](${replacement}${fragment}${title})`;
  });
}

function adaptMarkdown(text, document, contentMap) {
  let body = text.replaceAll("\r\n", "\n").replaceAll("\r", "\n");
  if (document.start_heading) {
    const marker = `## ${document.start_heading}`;
    const offset = body.indexOf(marker);
    if (offset < 0) throw new Error(`${document.source}: start heading not found`);
    body = body.slice(offset);
  } else {
    body = body.replace(/^# [^\n]+\n(?:\n)?/, "");
  }
  body = body.replace(/<!--[\s\S]*?-->/g, "");
  body = body.replace(/^(#{2,6}) \[([^\]]+)\](?:\([^)]+\))?/gm, "$1 $2");
  body = rewriteMarkdownLinks(body, document.source, contentMap);
  const frontmatter = [
    "---",
    `title: ${JSON.stringify(document.title)}`,
    `description: ${JSON.stringify(document.description)}`,
    `source: ${JSON.stringify(document.source)}`,
    "---",
    "",
  ].join("\n");
  return Buffer.from(`${frontmatter}${body.trimEnd()}\n`, "utf8");
}

function adaptMediaGuide(text) {
  const replacements = new Map([
    ["../brand/cueson/1.0.0/kit/icons/web/favicon.svg", "/assets/favicons/favicon.svg"],
    ["../brand/cueson/1.0.0/kit/fonts/fonts.css", "/assets/fonts/fonts.css"],
    ["../brand/cueson/1.0.0/kit/logos/svg/cueson-horizontal-color.svg", "/assets/logos/cueson-horizontal-color.svg"],
  ]);
  let output = text.replaceAll("\r\n", "\n").replaceAll("\r", "\n");
  for (const [source, target] of replacements) output = output.replaceAll(source, target);
  const metadataAnchor = '  <meta name="description" content="A plain-English guide to what Cueson unlocks and the subtitle and caption formats it is designed to understand.">';
  const publicMetadata = [
    metadataAnchor,
    '  <link rel="canonical" href="https://cueson.io/guides/media-formats/">',
    '  <meta property="og:title" content="Cueson media-format guide">',
    '  <meta property="og:description" content="A plain-English guide to what Cueson unlocks and the subtitle and caption formats it is designed to understand.">',
    '  <meta property="og:url" content="https://cueson.io/guides/media-formats/">',
    '  <meta property="og:image" content="https://cueson.io/assets/logos/cueson-social-preview-1280.png">',
    '  <meta name="twitter:card" content="summary_large_image">',
    '  <meta name="twitter:title" content="Cueson media-format guide">',
    '  <meta name="twitter:description" content="A plain-English guide to what Cueson unlocks and the subtitle and caption formats it is designed to understand.">',
    '  <meta name="twitter:image" content="https://cueson.io/assets/logos/cueson-social-preview-1280.png">',
  ].join("\n");
  if (!output.includes(metadataAnchor)) throw new Error("media guide: metadata anchor not found");
  output = output.replace(metadataAnchor, publicMetadata);
  output = output.replaceAll('<div class="table-scroll">', '<div class="table-scroll" tabindex="0" role="region" aria-label="Scrollable data table">');
  return Buffer.from(output, "utf8");
}

function resolveCommit(repoRoot, supplied) {
  const candidate = supplied ?? process.env.CUESON_SOURCE_COMMIT;
  if (candidate) {
    if (!/^[0-9a-f]{40}$/.test(candidate)) throw new Error("CUESON_SOURCE_COMMIT must be a full lowercase Git commit");
    return candidate;
  }
  const result = spawnSync("git", ["rev-parse", "HEAD"], {
    cwd: repoRoot,
    encoding: "utf8",
    windowsHide: true,
    shell: false,
  });
  const commit = result.stdout?.trim();
  if (result.status !== 0 || !/^[0-9a-f]{40}$/.test(commit)) throw new Error("unable to resolve the repository commit");
  return commit;
}

export async function loadContentMap(siteRoot = defaultSiteRoot) {
  const bytes = await readFile(path.join(siteRoot, "content-map.json"));
  const contentMap = JSON.parse(bytes.toString("utf8"));
  if (contentMap.generator_version !== 1 || contentMap.origin !== "https://cueson.io") throw new Error("content-map.json: unsupported authority metadata");
  for (const key of ["documents", "schemas", "downloads", "brand_assets", "direct_files"]) {
    if (!Array.isArray(contentMap[key])) throw new Error(`content-map.json: ${key} must be an array`);
  }
  for (const download of contentMap.downloads) {
    if (typeof download.name !== "string" || !/^https:\/\/github\.com\/shruggietech\/cueson\/releases\/download\/v1\.0\.0\//.test(download.url)) throw new Error("content-map.json: invalid release download");
  }
  const outputs = new Set();
  for (const document of contentMap.documents) {
    assertRelativePath(document.source, document.source);
    if (!Array.isArray(document.slug) || document.slug.some((segment) => !/^[a-z0-9][a-z0-9.-]*$/.test(segment))) throw new Error(`${document.source}: unsafe documentation slug`);
    const output = document.slug.length === 0 ? "content/generated/index.mdx" : `content/generated/${document.slug.join("/")}.mdx`;
    if (outputs.has(output.toLowerCase())) throw new Error(`${output}: duplicate generated path`);
    outputs.add(output.toLowerCase());
  }
  for (const record of [...contentMap.schemas, ...contentMap.brand_assets, ...contentMap.direct_files]) {
    assertRelativePath(record.source, record.source);
    assertRelativePath(record.public, record.public);
    const output = `public/${record.public}`.toLowerCase();
    if (outputs.has(output)) throw new Error(`${record.public}: duplicate generated path`);
    outputs.add(output);
  }
  return contentMap;
}

async function collectExpected(repoRoot, siteRoot, sourceCommit) {
  const contentMap = await loadContentMap(siteRoot);
  const expected = new Map();
  const documents = [];
  for (const document of [...contentMap.documents].sort((a, b) => a.order - b.order)) {
    const sourcePath = safeTarget(repoRoot, document.source);
    const source = await readUtf8(sourcePath, document.source);
    const output = adaptMarkdown(source.text, document, contentMap);
    const relative = document.slug.length === 0 ? "content/generated/index.mdx" : `content/generated/${document.slug.join("/")}.mdx`;
    expected.set(relative, output);
    documents.push({ source: document.source, route: docsRoute(document.slug), source_sha256: sha256(source.bytes), output_sha256: sha256(output) });
  }

  const rootPages = [];
  for (const document of [...contentMap.documents].sort((a, b) => a.order - b.order)) {
    const first = document.slug[0] ?? "index";
    if (!rootPages.includes(first)) rootPages.push(first);
  }
  expected.set("content/generated/meta.json", stableJson({ title: "Cueson", pages: rootPages }));
  expected.set("content/generated/formats/meta.json", stableJson({ title: "Formats", pages: ["subrip", "webvtt"] }));
  expected.set("content/generated/releases/meta.json", stableJson({ title: "Releases", pages: ["v1.0.0", "v0.0.0"] }));

  const brandManifestPath = path.join(repoRoot, "brand", "cueson", "1.0.0", "import-manifest.json");
  const brandManifest = JSON.parse((await readFile(brandManifestPath)).toString("utf8"));
  const approvedBrandFiles = new Map(brandManifest.entries.map((file) => [file.path, file.sha256]));
  const brandAssets = [];
  for (const asset of contentMap.brand_assets) {
    const sourceRelative = `brand/cueson/1.0.0/kit/${asset.source}`;
    const bytes = await readFile(safeTarget(repoRoot, sourceRelative));
    const digest = sha256(bytes);
    if (approvedBrandFiles.get(asset.source) !== digest) throw new Error(`${asset.source}: approved brand hash mismatch`);
    expected.set(`public/${asset.public}`, bytes);
    brandAssets.push({ source: asset.source, public_path: `/${asset.public}`, purpose: asset.purpose, bytes: bytes.length, sha256: digest });
  }

  const schemas = [];
  for (const schema of contentMap.schemas) {
    const bytes = await readFile(safeTarget(repoRoot, schema.source));
    const digest = sha256(bytes);
    if (digest !== schema.sha256) throw new Error(`${schema.source}: immutable schema hash mismatch`);
    expected.set(`public/${schema.public}`, bytes);
    schemas.push({ version: schema.version, source: schema.source, public_path: `/${schema.public}`, bytes: bytes.length, sha256: digest });
  }

  const directFiles = [];
  for (const directFile of contentMap.direct_files) {
    const source = await readUtf8(safeTarget(repoRoot, directFile.source), directFile.source);
    const output = directFile.transform === "media_guide" ? adaptMediaGuide(source.text) : source.bytes;
    expected.set(`public/${directFile.public}`, output);
    directFiles.push({ source: directFile.source, public_path: `/${directFile.public}`, source_sha256: sha256(source.bytes), output_sha256: sha256(output) });
  }

  const manifest = {
    generator_version: contentMap.generator_version,
    origin: contentMap.origin,
    documents,
    brand_assets: brandAssets,
    schemas,
    downloads: contentMap.downloads,
    direct_files: directFiles,
  };
  const manifestBytes = stableJson(manifest);
  expected.set("public/content-manifest.json", manifestBytes);
  expected.set("public/deployment.json", stableJson({
    commit: sourceCommit,
    content_manifest_sha256: sha256(manifestBytes),
    routes: ["/", ...documents.map((document) => document.route), "/guides/media-formats/", ...schemas.map((schema) => schema.public_path)],
    schemas: Object.fromEntries(schemas.map((schema) => [schema.version, schema.sha256])),
    downloads: contentMap.downloads,
  }));
  return { expected, summary: { commit: sourceCommit, content_manifest_sha256: sha256(manifestBytes), documents: documents.length, brand_assets: brandAssets.length, schemas: schemas.length, direct_files: directFiles.length } };
}

async function listFiles(root) {
  const files = [];
  async function walk(directory, prefix) {
    let entries;
    try {
      entries = await readdir(directory, { withFileTypes: true });
    } catch (error) {
      if (error.code === "ENOENT") return;
      throw error;
    }
    for (const entry of entries.sort((a, b) => a.name.localeCompare(b.name))) {
      const childPrefix = prefix ? `${prefix}/${entry.name}` : entry.name;
      const childPath = path.join(directory, entry.name);
      if (entry.isDirectory()) await walk(childPath, childPrefix);
      else if (entry.isFile()) files.push(childPrefix);
      else throw new Error(`${childPrefix}: generated output must be a regular file`);
    }
  }
  await walk(root, "");
  return files;
}

async function verifyExpected(siteRoot, expected) {
  const expectedPaths = new Set(expected.keys());
  for (const [relative, bytes] of expected) {
    let actual;
    try {
      actual = await readFile(safeTarget(siteRoot, relative));
    } catch (error) {
      if (error.code === "ENOENT") throw new Error(`${relative}: generated output drift (missing)`);
      throw error;
    }
    if (!actual.equals(bytes)) throw new Error(`${relative}: generated output drift (content)`);
  }
  const actualPaths = new Set(generatedFiles.filter((file) => expectedPaths.has(file)));
  for (const root of generatedRoots) {
    for (const relative of await listFiles(safeTarget(siteRoot, root))) actualPaths.add(`${root}/${toPosix(relative)}`);
  }
  const extras = [...actualPaths].filter((relative) => !expectedPaths.has(relative));
  if (extras.length > 0) throw new Error(`${extras[0]}: generated output drift (unexpected)`);
}

export async function generateSite({ repoRoot = defaultRepoRoot, siteRoot = defaultSiteRoot, sourceCommit, check = false } = {}) {
  const resolvedRepoRoot = path.resolve(repoRoot);
  const resolvedSiteRoot = path.resolve(siteRoot);
  if (path.dirname(resolvedSiteRoot) !== resolvedRepoRoot) throw new Error("site root must be an immediate child of the repository root");
  const commit = resolveCommit(resolvedRepoRoot, sourceCommit);
  const { expected, summary } = await collectExpected(resolvedRepoRoot, resolvedSiteRoot, commit);
  if (check) {
    await verifyExpected(resolvedSiteRoot, expected);
    return summary;
  }

  for (const root of generatedRoots) {
    const target = safeTarget(resolvedSiteRoot, root);
    if (target === resolvedSiteRoot) throw new Error(`${root}: refusing broad generated cleanup`);
    await rm(target, { recursive: true, force: true });
  }
  for (const file of generatedFiles) await rm(safeTarget(resolvedSiteRoot, file), { force: true });
  for (const [relative, bytes] of expected) {
    const target = safeTarget(resolvedSiteRoot, relative);
    await mkdir(path.dirname(target), { recursive: true });
    await writeFile(target, bytes);
  }
  await verifyExpected(resolvedSiteRoot, expected);
  return summary;
}

if (process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const check = process.argv.includes("--check");
  const unsupported = process.argv.slice(2).filter((argument) => argument !== "--check");
  if (unsupported.length > 0) throw new Error(`unsupported argument: ${unsupported[0]}`);
  const result = await generateSite({ check });
  process.stdout.write(`${check ? "verified" : "generated"} ${result.documents} documents, ${result.schemas} schemas, and ${result.brand_assets} brand assets for ${result.commit}\n`);
}
