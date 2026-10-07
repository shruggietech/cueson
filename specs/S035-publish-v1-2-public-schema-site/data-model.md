# S035 Data Model

## Release declaration

`site/content-map.json.release` contains exact version `1.2.0`, tag `v1.2.0` and official release URL. Seven downloads have unique names/URLs and filenames matching the declared version. Historical release pages own their historical links.

## Immutable schemas

Four records retain version, source/public path, positive `byte_length` and lowercase SHA-256. Historical 0.0.0/1.0.0/1.1.0 records remain unchanged. New 1.2.0 is 191170 bytes with SHA-256 `f2661a3d52effbab4a82a4d47197b5c7fae58496dc30a397ea3f2f668358b654`.

## Document/deployment inventory

Add two document records to existing 20: release 1.2 and consumer-speakers. Total 22 docs plus landing/media guide yields 24 HTML and four schema routes. Generated deployment/content manifests bind exact source commit, route/download/schema inventories and content digest.

## Lifecycle

Preparation #88: Ready → Specced → In progress → PR review → Done after merge. Production #86: preparation dependency → separately authorized Release verification → Done only after independent live checks. Default Project Status stays unused; Slice carries S035. Milestone remains open until production acceptance.
