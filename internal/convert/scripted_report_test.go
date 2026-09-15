package convert

import (
	"os"
	"testing"

	"github.com/shruggietech/cueson/internal/schema"
)

func TestScriptedLossReferencesResolveNativeOrders(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			data, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
			if err != nil {
				t.Fatal(err)
			}
			document, err := schema.Decode(data)
			if err != nil {
				t.Fatal(err)
			}
			native := document.FormatData.ASS
			if format == "ssa" {
				native = document.FormatData.SSA
			}
			order := native.Records[0].SourceOrder
			loss := newLoss(document, "subrip", LossCodeScriptedMetadataOmitted, KindOmitted, "native metadata field is omitted", "/format_data/"+format+"/records/0/fields/0", &order, nil, nil)
			report, err := NewReport(document, []Loss{loss})
			if err != nil || len(report.Losses) != 1 {
				t.Fatalf("native reference rejected: %v", err)
			}
			fabricated := len(native.Sections) + len(native.Records) + 1
			loss.SourceOrder = &fabricated
			if _, err := NewReport(document, []Loss{loss}); err == nil {
				t.Fatal("fabricated native reference accepted")
			}
			header := native.Sections[0].SourceOrder
			loss.Code, loss.Path, loss.SourceOrder = LossCodeScriptedSectionOmitted, "/format_data/"+format+"/sections/0", &header
			if _, err := NewReport(document, []Loss{loss}); err != nil {
				t.Fatalf("native header reference rejected: %v", err)
			}
		})
	}
}
