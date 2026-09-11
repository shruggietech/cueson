import defaultComponents from "fumadocs-ui/mdx";
import type { MDXComponents } from "mdx/types";
import type { InputHTMLAttributes, TableHTMLAttributes } from "react";

function AccessibleInput(properties: InputHTMLAttributes<HTMLInputElement>) {
  return <input {...properties} aria-label={properties["aria-label"] ?? (properties.checked ? "Completed item" : "Incomplete item")} />;
}

function AccessibleTable(properties: TableHTMLAttributes<HTMLTableElement>) {
  return <div className="relative overflow-auto prose-no-margin my-6" tabIndex={0} role="region" aria-label="Scrollable data table"><table {...properties} /></div>;
}

export function getMDXComponents(components?: MDXComponents): MDXComponents {
  return {
    ...defaultComponents,
    input: AccessibleInput,
    table: AccessibleTable,
    ...components,
  };
}
