import { notFound } from "next/navigation";
import { DocsBody, DocsDescription, DocsPage, DocsTitle } from "fumadocs-ui/layouts/docs/page";
import { getMDXComponents } from "@/components/mdx";
import { pageMetadata } from "@/lib/site";
import { source } from "@/lib/source";

type PageProperties = { params: Promise<{ slug?: string[] }> };

export function generateStaticParams() { return source.generateParams(); }

export async function generateMetadata({ params }: PageProperties) {
  const { slug = [] } = await params;
  const page = source.getPage(slug);
  if (!page) notFound();
  return pageMetadata(page.data.title, page.data.description ?? page.data.title, `/docs/${slug.length ? `${slug.join("/")}/` : ""}`);
}

export default async function DocumentationPage({ params }: PageProperties) {
  const { slug = [] } = await params;
  const page = source.getPage(slug);
  if (!page) notFound();
  const Body = page.data.body;
  return (
    <DocsPage toc={page.data.toc} full={false}>
      <DocsTitle>{page.data.title}</DocsTitle>
      {page.data.description && <DocsDescription>{page.data.description}</DocsDescription>}
      <DocsBody><Body components={getMDXComponents()} /></DocsBody>
    </DocsPage>
  );
}
