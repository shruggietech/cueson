# Quickstart: Establish CI and Cross-Platform Build Gates

## Prerequisites

- A clean checkout of `S006-establish-ci-gates`.
- Go 1.25 or newer for product verification.
- PowerShell 7 on Windows, macOS, or Linux for the cross-platform local sequence below.
- Network access for pinned Staticcheck, govulncheck, and actionlint modules.
- GitHub CLI authentication for publication and hosted-run inspection.
- Every intended final path staged before the whole-tree whitespace check, so new files are included without changing their content.

## Local parity

Run repository and product verification from the repository root in PowerShell 7. The explicit Go 1.26.8 analysis toolchain prevents an older installed launcher from compiling the pinned analyzers with an incompatible standard library. Product tests and builds continue to prove the Go 1.25 compatibility floor.

```powershell
function Assert-NativeSuccess([string]$Step) {
    if ($LASTEXITCODE -ne 0) {
        throw "$Step failed with exit code $LASTEXITCODE"
    }
}

$unformatted = @(gofmt -l cmd internal scripts)
if ($LASTEXITCODE -ne 0 -or $unformatted.Count -ne 0) {
    $unformatted | Write-Error
    throw "gofmt reported repository drift"
}

go run ./scripts/github-format/main.go .
Assert-NativeSuccess "repository text check"
go -C scripts/github-format test -count=1 ./...
Assert-NativeSuccess "publication formatter tests"
go vet ./...
Assert-NativeSuccess "go vet"
go test -count=1 ./...
Assert-NativeSuccess "root tests"
go test -race -count=1 ./...
Assert-NativeSuccess "race tests"
go test -count=1 ./internal/schema ./internal/conformance ./internal/testutil
Assert-NativeSuccess "schema and conformance tests"

$previousToolchain = $env:GOTOOLCHAIN
try {
    $env:GOTOOLCHAIN = "go1.26.8"
    go run honnef.co/go/tools/cmd/staticcheck@v0.7.0 ./...
    Assert-NativeSuccess "Staticcheck"
    go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
    Assert-NativeSuccess "govulncheck"
    go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
    Assert-NativeSuccess "actionlint"
} finally {
    $env:GOTOOLCHAIN = $previousToolchain
}

$emptyTree = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
git diff --cached --check $emptyTree
Assert-NativeSuccess "staged whole-tree whitespace check"
git diff --check
Assert-NativeSuccess "unstaged whitespace check"
```

Cross-build all six supported targets with `CGO_ENABLED=0`. This closed PowerShell sequence sends local outputs to the operating-system null device, so compilation is proven without creating or publishing an artifact. The hosted matrix separately verifies non-empty binaries beneath the runner temporary directory.

```powershell
$targets = @(
    @{ Name = "linux-amd64"; GOOS = "linux"; GOARCH = "amd64" },
    @{ Name = "linux-arm64"; GOOS = "linux"; GOARCH = "arm64" },
    @{ Name = "windows-amd64"; GOOS = "windows"; GOARCH = "amd64" },
    @{ Name = "windows-arm64"; GOOS = "windows"; GOARCH = "arm64" },
    @{ Name = "darwin-amd64"; GOOS = "darwin"; GOARCH = "amd64" },
    @{ Name = "darwin-arm64"; GOOS = "darwin"; GOARCH = "arm64" }
)

$previousCgo = $env:CGO_ENABLED
$previousGoos = $env:GOOS
$previousGoarch = $env:GOARCH
$previousToolchain = $env:GOTOOLCHAIN
try {
    $env:CGO_ENABLED = "0"
    $env:GOTOOLCHAIN = "go1.25.0"
    $nullDevice = if ($IsWindows) { "NUL" } else { "/dev/null" }
    foreach ($target in $targets) {
        $env:GOOS = $target.GOOS
        $env:GOARCH = $target.GOARCH
        go build -trimpath -o $nullDevice ./cmd/cueson
        if ($LASTEXITCODE -ne 0) {
            throw "cross-build failed for $($target.Name)"
        }
    }
} finally {
    $env:CGO_ENABLED = $previousCgo
    $env:GOOS = $previousGoos
    $env:GOARCH = $previousGoarch
    $env:GOTOOLCHAIN = $previousToolchain
}
```

The repository formatter covers UTF-8, BOM, mojibake, and governed line-ending checks. Confirm the final `git status --short` contains only intended S006 changes and no build output.

## Hosted proof

Follow [failure-proof.md](contracts/failure-proof.md) exactly. The first draft revision must show the controlled `Repository text` failure. The corrected revision must show every name in [check-contract.md](contracts/check-contract.md) green before the pull request is marked ready.

Inspect workflow permissions and triggers against [permissions-and-events.md](contracts/permissions-and-events.md). Read each workflow and hosted run back from GitHub rather than treating a successful command response as proof.

## Completion

Confirm the final tree has no `.github/ci-failure-probe`, no uploaded artifacts, no release or tag, no repository-setting change, and no check that claims an unimplemented codec, conversion, review-automation, or release capability. Post-merge housekeeping verifies the live CI badge and default-branch runs.
