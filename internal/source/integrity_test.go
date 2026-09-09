package source

import (
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

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
