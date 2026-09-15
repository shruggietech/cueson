package model

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const scriptedOwnershipMutations = 18

// FuzzScriptedDocumentOwnership mutates an accepted ownership graph, not random
// JSON that almost always stops at syntax. Fixtures are acquired once; every
// mutation and repeated validation operates on independent in-memory copies.
func FuzzScriptedDocumentOwnership(f *testing.F) {
	fixtures := make([][]byte, 2)
	for i, format := range []string{"ass", "ssa"} {
		var err error
		fixtures[i], err = os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			f.Fatal(err)
		}
		for mutation := 0; mutation < scriptedOwnershipMutations; mutation++ {
			f.Add("edited caption", uint8(mutation), uint8(i))
		}
	}
	f.Fuzz(func(t *testing.T, payload string, mutation, dialect uint8) {
		if len(payload) > 4096 {
			t.Skip()
		}
		var d Document
		if err := json.Unmarshal(fixtures[int(dialect)%len(fixtures)], &d); err != nil {
			t.Fatal(err)
		}
		if err := d.Validate(); err != nil {
			t.Fatalf("baseline rejected: %v", err)
		}
		wantAccepted := mutateScriptedOwnership(&d, payload, int(mutation)%scriptedOwnershipMutations)
		before, err := json.Marshal(d)
		if err != nil {
			t.Fatal(err)
		}
		first, second := d.Validate(), d.Validate()
		if (first == nil) != wantAccepted || (second == nil) != wantAccepted {
			t.Fatalf("mutation %d acceptance: %v / %v, want %v", int(mutation)%scriptedOwnershipMutations, first, second, wantAccepted)
		}
		if first != nil && first.Error() != second.Error() {
			t.Fatal("nondeterministic owning-boundary rejection")
		}
		if first != nil && strings.Contains(first.Error(), "C:/private-ownership-sentinel/") {
			t.Fatal("owning-boundary diagnostic exposed private metadata")
		}
		after, err := json.Marshal(d)
		if err != nil || !bytes.Equal(before, after) {
			t.Fatal("validation mutated owners or source")
		}
	})
}

func mutateScriptedOwnership(d *Document, payload string, mutation int) bool {
	n := d.FormatData.ASS
	c := d.Cues[0].FormatData.ASS
	if d.Format == "ssa" {
		n, c = d.FormatData.SSA, d.Cues[0].FormatData.SSA
	}
	missing := "missing-owner"
	private := "C:/private-ownership-sentinel/" + payload
	switch mutation {
	case 0:
		// Accepted control: capability observations can truthfully change from
		// schema-only example recognition to the installed native codec.
		d.FormatSupport = FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true}
		return true
	case 1:
		// Common timing is intentionally editable independently of source truth.
		d.Cues[0].Timing.StartMilliseconds++
		d.Cues[0].Timing.DurationMilliseconds--
		*d.Document.MediaStartMilliseconds++
		*d.Document.MediaSpanMilliseconds--
		*d.Stats.MediaSpanMilliseconds--
		return true
	case 2:
		n.Records[0].SectionID = missing
	case 3:
		n.Records[0].SourceOrder = len(n.Sections) + len(n.Records) + 1
	case 4:
		n.Events[0].RecordID = missing
	case 5:
		n.Events[0].DeclarationID = &missing
	case 6:
		c.StyleID = missing
	case 7:
		c.Projection.SourceEventID = missing
	case 8:
		n.Events[0].Spans[0].EndScalar++
	case 9:
		d.Cues[0].Payload.PlainText = "forged derived content " + payload
	case 10:
		n.Sections[0].Name = private
	case 11:
		d.Metadata.Title = &private
	case 12:
		d.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64)
	case 13:
		d.Source.Assets[0].Size.Bytes++
	case 14:
		d.Source.Assets[0].DataBase64 = "%%%" + payload
	case 15:
		n.Sections[0].RawHeader = &missing
	case 16:
		n.Styles[0].Fields[2].RawValue = "NaN"
	case 17:
		n.Dialect = "v4.00++"
	}
	return false
}

func TestScriptedOwnershipMutationControlsAndRejections(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for mutation := 0; mutation < scriptedOwnershipMutations; mutation++ {
			d := scriptedExample(t, format)
			wantAccepted := mutateScriptedOwnership(&d, "字幕 😀 C:/nested/private", mutation)
			if err := d.Validate(); (err == nil) != wantAccepted {
				t.Fatalf("%s mutation %d: %v, want accepted %v", format, mutation, err, wantAccepted)
			}
		}
	}
}
