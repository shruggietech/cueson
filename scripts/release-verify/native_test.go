package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPackagedNativeWorkflowConformance(t *testing.T) {
	repository, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "cueson")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}
	command := exec.Command("go", "build", "-trimpath", "-o", binary, "./cmd/cueson")
	command.Dir = repository
	configureProcess(command)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build native test binary: %v: %s", err, output)
	}
	if _, err := verifyNativeWorkflows(context.Background(), binary, repository, t.TempDir(), "1.1.0"); err != nil {
		t.Fatal(err)
	}
	if count, err := verifyHistoricalNative(context.Background(), binary, repository, t.TempDir()); err != nil || count != 2 {
		t.Fatalf("historical packaged proof: count=%d error=%v", count, err)
	}
	// The normal unit suite never needs network access. Explicit proof execution
	// uses the actual public release asset, not a source rebuild substitute.
	if os.Getenv("CUESON_VERIFY_PUBLISHED_CONSUMER") == "1" {
		binaryBytes, err := os.ReadFile(binary)
		if err != nil {
			t.Fatal(err)
		}
		proof, err := verifyStableHostCLI(context.Background(), binaryBytes, filepath.Base(binary), repository, "1.1.0")
		if err != nil {
			t.Fatal(err)
		}
		if proof.OldConsumerRefusalCount != 32 || proof.OldConsumerIdentityProbeCount != 32 {
			t.Fatalf("old consumer refusal count = %d", proof.OldConsumerRefusalCount)
		}
	}
}

func TestNativePayloadContractRejectsIdentityAndSourceDrift(t *testing.T) {
	valid := []byte(`{"$schema":"https://cueson.io/schema/v1.1.0/cueson.schema.json","schema_version":"1.1.0","format":"subrip","format_support":{"status":"stable","ingest_supported":true,"render_supported":true,"restore_supported":true,"ocr_required_for_semantic_output":false},"producer":{"name":"cueson","version":"1.1.0"},"source":{"assets":[{"data_base64":"SGVsbG8=","size":{"bytes":5},"hashes":{"sha256":"185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969"}}]}}`)
	for name, payload := range map[string][]byte{
		"valid":        valid,
		"development":  bytes.ReplaceAll(valid, []byte("1.1.0"), []byte("1.1.0-dev")),
		"experimental": bytes.ReplaceAll(valid, []byte(`"stable"`), []byte(`"experimental"`)),
		"capability":   bytes.ReplaceAll(valid, []byte(`"render_supported":true`), []byte(`"render_supported":false`)),
		"ocr":          bytes.ReplaceAll(valid, []byte(`"ocr_required_for_semantic_output":false`), []byte(`"ocr_required_for_semantic_output":true`)),
		"producer":     bytes.ReplaceAll(valid, []byte(`"name":"cueson"`), []byte(`"name":"third-party"`)),
		"format":       bytes.ReplaceAll(valid, []byte(`"format":"subrip"`), []byte(`"format":"ass"`)),
		"bytes":        bytes.ReplaceAll(valid, []byte("SGVsbG8="), []byte("V29ybGQ=")),
		"hash":         bytes.ReplaceAll(valid, []byte("185f8db3"), []byte("085f8db3")),
		"length":       bytes.ReplaceAll(valid, []byte(`"bytes":5`), []byte(`"bytes":4`)),
		"path":         bytes.ReplaceAll(valid, []byte(`"version":"1.1.0"`), []byte(`"version":"1.1.0","path":"C:\\private\\source"`)),
	} {
		t.Run(name, func(t *testing.T) {
			err := verifyNativePayload(payload, []byte("Hello"), "1.1.0", "subrip", nil)
			if (err == nil) != (name == "valid") {
				t.Fatalf("contract result = %v", err)
			}
		})
	}
}

func TestOldConsumerRefusalRequiresExactIdentityDiagnostic(t *testing.T) {
	for _, test := range []struct {
		diagnostic string
		accept     bool
	}{
		{`schema_version: value must be "1.0.0"`, true},
		{`$schema: value must be "https://cueson.io/schema/v1.0.0/cueson.schema.json"`, true},
		{"source hash mismatch", false}, {"cannot read input", false},
		{"schema_version: malformed value", false}, {"1.0.0: conversion refused", false},
	} {
		if (verifyOldIdentityDiagnostic([]byte(test.diagnostic)) == nil) != test.accept {
			t.Fatalf("identity diagnostic acceptance differs for %q", test.diagnostic)
		}
	}
}

func TestRefusalRequiresRuntimeFailureAndNoPublication(t *testing.T) {
	for _, test := range []struct {
		name           string
		status         int
		stdout, stderr []byte
		change         bool
		accept         bool
	}{
		{"valid", 1, nil, []byte("unsupported identity"), false, true},
		{"success", 0, nil, nil, false, false},
		{"usage", 2, nil, []byte("usage"), false, false},
		{"payload", 1, []byte("partial"), []byte("error"), false, false},
		{"silent", 1, nil, nil, false, false},
		{"replaced", 1, nil, []byte("error"), true, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "output")
			data := []byte("sentinel")
			if test.change {
				data = []byte("changed")
			}
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
			err := verifyRefusal(test.status, test.stdout, test.stderr, path, nil)
			if (err == nil) != test.accept {
				t.Fatalf("refusal result = %v", err)
			}
		})
	}
}
