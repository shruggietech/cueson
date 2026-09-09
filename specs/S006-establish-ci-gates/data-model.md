# Data Model: Establish CI and Cross-Platform Build Gates

S006 stores no product data. Its design entities describe observable delivery state and the fixed relationships that later repository controls will consume.

## Workflow

| Field | Meaning | Validation |
|---|---|---|
| `name` | Stable workflow display name | Exactly `CI` or `CodeQL` |
| `events` | Changes that start the workflow | Pull requests to `main`; pushes to `main`; CodeQL may also run on a schedule |
| `permissions` | Workflow token authority | Read-only CI; CodeQL adds only `security-events: write` |
| `concurrency_group` | Boundary for superseded work | Includes workflow, event type, and pull-request number or ref |
| `cancel_in_progress` | Whether an older run is obsolete | True only for pull-request runs |

## Quality Gate

| Field | Meaning | Validation |
|---|---|---|
| `check_name` | GitHub-visible stable job name | Matches the check contract exactly |
| `command` | Foreground verifier | Exits non-zero on failure; no ignored error |
| `tool_version` | Analysis version where applicable | Explicit module version or immutable action commit |
| `runner` | Owning hosted environment | One explicit GA runner label |
| `evidence_kind` | Claim supported by the result | Repository, native behavior, race, security, or build portability |

## Native Platform Run

| Field | Meaning | Validation |
|---|---|---|
| `os_name` | Operating system whose behavior is exercised | Linux, Windows, or macOS |
| `runner` | Native hosted runner | `ubuntu-24.04`, `windows-2025`, or `macos-15` |
| `go_version` | Compatibility toolchain | Latest patch in Go 1.25 |
| `suite` | Executed tests | Complete root-module package suite |

## Cross-Build Target

| Field | Meaning | Validation |
|---|---|---|
| `goos` | Target operating system | `windows`, `darwin`, or `linux` |
| `goarch` | Target processor architecture | `amd64` or `arm64` |
| `cgo_enabled` | Native dependency state | Always `0` |
| `output` | Temporary binary path | Runner temporary directory only; never uploaded |

## Pinned Execution Reference

| Field | Meaning | Validation |
|---|---|---|
| `source` | Action repository or Go module | Explicit trusted source |
| `immutable_version` | Commit SHA or module version | Full action commit or exact semantic module version |
| `readable_version` | Reviewer-facing release identity | Adjacent workflow comment or command version |
| `verification_source` | How the reference was resolved | Official GitHub API or Go module checksum system |

## Controlled Failure Evidence

The state transition is append-only pull-request history:

```text
probe absent -> draft revision adds probe -> Repository text fails -> correction removes probe -> complete gate set passes -> pull request becomes ready
```

The probe file is not part of the final tree. The failed and corrected hosted runs remain attached to pull-request commits as evidence.
