package scripted

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

// FuzzRenderParseCycle checks native semantic retention and canonical stability.
// Bounded fuzz inputs supplement the separate large-model/output regressions.
func FuzzRenderParseCycle(f *testing.F) {
	for _, format := range []string{"ass", "ssa"} {
		encoded, err := os.ReadFile("../../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			f.Fatal(err)
		}
		var d model.Document
		if err = json.Unmarshal(encoded, &d); err != nil {
			f.Fatal(err)
		}
		raw, err := base64.StdEncoding.DecodeString(d.Source.Assets[0].DataBase64)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(raw)
		f.Add(append(append([]byte{}, raw...), []byte("Inert extension\n")...))
		f.Add(append([]byte{0xef, 0xbb, 0xbf}, bytes.ReplaceAll(raw, []byte{'\n'}, []byte{'\r', '\n'})...))
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 64<<10 {
			t.Skip()
		}
		detection := Detect(raw)
		if detection.Err != nil || !detection.Candidate {
			return
		}
		if _, err := Parse(context.Background(), raw, detection.Format); err != nil {
			return
		}
		d := renderParsedFixture(t, detection.Format, raw)
		captured := d.Source.Assets[0].DataBase64
		first, err := Render(context.Background(), d, false)
		if err != nil {
			return
		} // Malformed retained records are restoration-only.
		if d.Source.Assets[0].DataBase64 != captured {
			t.Fatal("renderer modified source truth")
		}
		reparsed := renderParsedFixture(t, detection.Format, first.Bytes)
		second, err := Render(context.Background(), reparsed, false)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(first.Bytes, second.Bytes) {
			t.Fatal("canonical render/parse bytes drift")
		}
		if len(d.Cues) != len(reparsed.Cues) {
			t.Fatal("dialogue occurrence omitted")
		}
		for i := range d.Cues {
			a, b := d.Cues[i], reparsed.Cues[i]
			if a.Timing != b.Timing || !reflect.DeepEqual(a.Payload, b.Payload) || !reflect.DeepEqual(a.Speakers, b.Speakers) || !reflect.DeepEqual(a.Tokens, b.Tokens) {
				t.Fatal("common semantic data lost")
			}
		}
		a, b := renderNative(&d), renderNative(&reparsed)
		if len(a.Sections) != len(b.Sections) || len(a.Records) != len(b.Records) || len(a.Styles) != len(b.Styles) || len(a.Events) != len(b.Events) || len(a.Attachments) != len(b.Attachments) {
			t.Fatal("ordered native occurrence omitted")
		}
		for i := range a.Sections {
			if a.Sections[i].Name != b.Sections[i].Name {
				t.Fatal("section identity lost")
			}
		}
		for i := range a.Events {
			if a.Events[i].EventType != b.Events[i].EventType || a.Events[i].Text != b.Events[i].Text {
				t.Fatal("native event Text or type changed")
			}
		}
		for i := range a.Attachments {
			if a.Attachments[i].Name != b.Attachments[i].Name || a.Attachments[i].AttachmentType != b.Attachments[i].AttachmentType || a.Attachments[i].DataRecordCount != b.Attachments[i].DataRecordCount {
				t.Fatal("embedded content owner changed")
			}
		}
	})
}
