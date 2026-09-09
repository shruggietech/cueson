# Quickstart: S004 Verification

Run commands from the repository root on a clean checkout of `S004-establish-source-integrity`.

## Focused Package Verification

```powershell
go test ./internal/schema ./internal/model ./internal/source ./internal/cli
```

## Complete Local Verification

```powershell
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/cueson
```

## Platform Selection Verification

Run native Windows timestamp tests in the current environment, then compile the Linux and macOS test binaries with CGO disabled. Cross-compilation verifies build selection only and is not recorded as native behavioral proof.

```powershell
go test ./internal/source
$env:CGO_ENABLED = '0'
$env:GOOS = 'linux'
$env:GOARCH = 'amd64'
go test -c ./internal/source
$env:GOOS = 'darwin'
go test -c ./internal/source
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
Remove-Item Env:CGO_ENABLED
```

## Single-Asset Smoke Test

```powershell
go run ./cmd/cueson restore --no-metadata --output .\restored.bin .\internal\schema\testdata\representative.cueson.json
```

Verify that stdout is empty and the restored file matches the representative source asset's declared size and SHA-256. Remove the smoke-test output after verification.

## Multi-Asset and Safety Smoke Tests

Use a temporary existing directory and a valid multi-asset test document:

```powershell
go run ./cmd/cueson restore --no-metadata --output-dir .\restore-output .\path\to\multi-asset.cueson.json
```

Repeat without `--force` against existing regular destinations to verify refusal, then with `--force` to verify approved replacement. Also exercise one corrupt integrity declaration and confirm that no destination changes and no staging files remain.

The reusable repository-wide fixture, malformed-corpus, and fuzz foundation remains issue #7. Hosted native operating-system execution remains issue #8.
