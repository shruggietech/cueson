package source

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestValidateIntegrityIsReadOnlyAndContextAware(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	want := testDocument(t)
	if err := ValidateIntegrity(context.Background(), document); err != nil {
		t.Fatalf("ValidateIntegrity() error = %v", err)
	}
	if !reflect.DeepEqual(document, want) {
		t.Fatal("ValidateIntegrity() mutated its input document")
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ValidateIntegrity(canceled, document); !errors.Is(err, context.Canceled) {
		t.Fatalf("ValidateIntegrity() error = %v, want context.Canceled", err)
	}
}

func TestValidateIntegrityRejectsCorruptEnvelopeWithoutOutputPlanning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*model.Document)
		want   string
	}{
		{name: "unsafe basename", mutate: func(document *model.Document) { document.Source.Assets[0].FileName = "../source.srt" }, want: "prohibited character"},
		{name: "portable collision", mutate: func(document *model.Document) {
			asset := document.Source.Assets[0]
			asset.ID = "asset-1"
			asset.Role = "companion"
			asset.FileName = strings.ToUpper(document.Source.Assets[0].FileName)
			document.Source.Assets = append(document.Source.Assets, asset)
		}, want: "colliding portable basenames"},
		{name: "corrupt source bytes", mutate: func(document *model.Document) { document.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64) }, want: "SHA-256"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := testDocument(t)
			test.mutate(&document)
			if err := ValidateIntegrity(context.Background(), document); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("ValidateIntegrity() error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestPrepareAssetsValidatesIntegrity(t *testing.T) {
	t.Parallel()

	document := testDocument(t)
	assets, err := prepareAssets(document)
	if err != nil {
		t.Fatalf("prepareAssets() error = %v", err)
	}
	if len(assets) != 1 || assets[0].asset.ID != "asset-0" {
		t.Fatalf("prepareAssets() = %#v", assets)
	}

	tests := []struct {
		name   string
		mutate func()
		want   string
	}{
		{name: "invalid base64", mutate: func() { document.Source.Assets[0].DataBase64 = "%%%%" }, want: "base64"},
		{name: "noncanonical trailing bits", mutate: func() { document.Source.Assets[0].DataBase64 = "Zh==" }, want: "base64"},
		{name: "line break", mutate: func() { document.Source.Assets[0].DataBase64 += "\n" }, want: "canonical"},
		{name: "length", mutate: func() { document.Source.Assets[0].Size.Bytes++ }, want: "decoded length"},
		{name: "digest", mutate: func() { document.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64) }, want: "SHA-256"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document = testDocument(t)
			tt.mutate()
			if _, err := prepareAssets(document); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("prepareAssets() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func testDocument(t *testing.T) model.Document {
	t.Helper()
	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	return document
}
