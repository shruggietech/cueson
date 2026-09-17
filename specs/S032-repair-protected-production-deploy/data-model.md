# Data Model: Protected production deployment repair

## Deployment identity

| Field | Constraint | Source of authority |
| --- | --- | --- |
| Requested revision | Full lowercase 40-character commit identifier | Manual workflow input |
| Workflow ref | Exactly `refs/heads/main` | GitHub Actions context |
| Workflow revision | Exactly requested revision | GitHub Actions context |
| Fetched default branch | Exactly requested revision at both guarded checkpoints | Fresh `origin/main` fetch |
| Public revision | Exactly requested revision | `https://cueson.io/deployment.json` |
| GitHub environment revision | Exactly requested revision | Protected environment deployment record |

Any inequality is a terminal failure before the next authority boundary.

## Protected credential contract

| Attribute | Required value |
| --- | --- |
| Secret name | `CLOUDFLARE_API_TOKEN` |
| GitHub scope | `production` environment |
| Owner | ShruggieTech organization security owner |
| Account resource | ShruggieTech Cloudflare account only |
| Zone resource | `cueson.io` only |
| Account permission | Entire Account, Developer Platform, Workers Scripts Legacy: Edit |
| Zone permissions | Specified Domains `cueson.io`: DNS: Read, Zone: Read, Zone Transform Rules: Read, Workers Routes: Read |
| Workflow exposure | Cloudflare preflight, deploy, and post-deployment read-back steps only |
| Rotation | At most 90 days between rotations |
| Emergency action | Immediate revoke, replace, preflight, and audit |

Secret value, token identifier, and other authentication material are never repository data.

## Cloudflare state snapshot

| Collection | Comparison rule |
| --- | --- |
| Zone identity | Exact expected zone/account identity |
| DNS records | Intended Cueson change allowed; unrelated records preserved |
| Worker scripts | Intended `cueson-site` content change allowed; unrelated scripts preserved |
| Custom Domains | Intended `cueson.io` and `www.cueson.io` bindings present; unrelated domains preserved |
| Redirect rulesets | Required redirect behavior present; unrelated rules preserved |

The before snapshot is an immutable run artifact consumed by post-deployment comparison.

## Issue lifecycle

| Issue | Completion point |
| --- | --- |
| S032 implementation child | Workflow, tests, documentation, and local/hosted review merge successfully |
| #78 production validation | Durable credential is provisioned and the merged default-branch run plus independent live verification pass |
