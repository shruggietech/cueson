package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestStableInstalledCodecsPreserveLoadedCapabilityObservations(t *testing.T) {
	t.Parallel()
	for _, format := range []string{"ass", "ssa"} {
		for _, support := range []model.FormatSupport{
			{Status: "schema_only", RestoreSupported: true},
			{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
			{Status: "stable", IngestSupported: true, RenderSupported: true, RestoreSupported: true},
		} {
			t.Run(format+"/"+support.Status, func(t *testing.T) {
				fixture, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
				if err != nil {
					t.Fatal(err)
				}
				document, err := schema.Decode(fixture)
				if err != nil {
					t.Fatal(err)
				}
				document.FormatSupport = support
				document.Producer = model.Producer{Name: "independent-producer", Version: "42.7-custom"}
				payload, err := json.Marshal(document)
				if err != nil {
					t.Fatal(err)
				}
				input := filepath.Join(t.TempDir(), "misleading.vtt")
				if err := os.WriteFile(input, payload, 0o600); err != nil {
					t.Fatal(err)
				}
				status, stdout, stderr := runForTest(context.Background(), []string{"inspect", input, "--json"})
				if status != ExitSuccess {
					t.Fatalf("inspect = %d %s", status, stderr)
				}
				var report inspectionReport
				if err := json.Unmarshal([]byte(stdout), &report); err != nil {
					t.Fatal(err)
				}
				if report.Capabilities.Declared != inspectionDeclaredCapabilities(support) || !report.Capabilities.Installed.Ingest || !report.Capabilities.Installed.Render {
					t.Fatalf("loaded observation normalized by installed stable codecs: %+v", report.Capabilities)
				}
				after, err := os.ReadFile(input)
				if err != nil || !bytes.Equal(payload, after) {
					t.Fatal("inspection rewrote source document or producer")
				}
			})
		}
	}
}
