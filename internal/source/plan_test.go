package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanDestinations(t *testing.T) {
	t.Parallel()

	assets, err := prepareAssets(testDocument(t))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	output := filepath.Join(directory, "renamed.srt")
	plans, err := planDestinations(assets, RestoreOptions{Output: output, Metadata: MetadataNone})
	if err != nil || len(plans) != 1 || plans[0].destination != output {
		t.Fatalf("planDestinations() = (%#v, %v)", plans, err)
	}
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		t.Fatalf("planning created output: %v", err)
	}
	defaultPlans, err := planDestinations(assets, RestoreOptions{Metadata: MetadataNone})
	if err != nil {
		t.Fatalf("default plan error = %v", err)
	}
	wantDefault, err := filepath.Abs(assets[0].asset.FileName)
	if err != nil || defaultPlans[0].destination != wantDefault {
		t.Fatalf("default destination = %q, want %q; error = %v", defaultPlans[0].destination, wantDefault, err)
	}
	if err := os.WriteFile(output, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := planDestinations(assets, RestoreOptions{Output: output, Metadata: MetadataNone}); err == nil {
		t.Fatal("planning accepted existing output without force")
	}
	if plans, err := planDestinations(assets, RestoreOptions{Output: output, Force: true, Metadata: MetadataNone}); err != nil || plans[0].existing == nil {
		t.Fatalf("forced plan = (%#v, %v)", plans, err)
	}
}

func TestPlanRejectsModeAndTargetViolations(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	companion := document.Source.Assets[0]
	companion.ID = "asset-1"
	companion.Role = "companion"
	companion.FileName = "other.srt"
	document.Source.Assets = append(document.Source.Assets, companion)
	assets, err := prepareAssets(document)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	for _, options := range []RestoreOptions{
		{Metadata: MetadataNone},
		{Output: filepath.Join(directory, "one"), Metadata: MetadataNone},
		{Output: "one", OutputDir: directory, Metadata: MetadataNone},
	} {
		if _, err := planDestinations(assets, options); err == nil {
			t.Errorf("planDestinations(%+v) error = nil", options)
		}
	}
	nonRegular := filepath.Join(directory, document.Source.Assets[0].FileName)
	if err := os.Mkdir(nonRegular, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := planDestinations(assets[:1], RestoreOptions{OutputDir: directory, Force: true, Metadata: MetadataNone}); err == nil {
		t.Fatal("planning accepted directory destination")
	}
}

func TestPlanRejectsExistingNativeIdentityCollision(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	companion := document.Source.Assets[0]
	companion.ID = "asset-1"
	companion.Role = "companion"
	companion.FileName = "other.srt"
	document.Source.Assets = append(document.Source.Assets, companion)
	assets, err := prepareAssets(document)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	first := filepath.Join(directory, "captions.srt")
	second := filepath.Join(directory, "other.srt")
	if err := os.WriteFile(first, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(first, second); err != nil {
		t.Skipf("filesystem does not support hard-link identity test: %v", err)
	}
	if _, err := planDestinations(assets, RestoreOptions{OutputDir: directory, Force: true, Metadata: MetadataNone}); err == nil {
		t.Fatal("planDestinations() accepted two paths for one existing native file")
	}
}
