package convert

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func scriptedEmphasisLossSource(t *testing.T, format, field, value, text string, other, secondCue bool) string {
	t.Helper()
	source := scriptedConversionSource(format, text)
	lines := strings.Split(source, "\n")
	for index, line := range lines {
		if !strings.HasPrefix(line, "Style: ") {
			continue
		}
		values := strings.Split(strings.TrimPrefix(line, "Style: "), ",")
		for fieldIndex, name := range model.ScriptedCanonicalFields(format, "style") {
			if strings.EqualFold(name, field) {
				values[fieldIndex] = value
			}
		}
		lines[index] = "Style: " + strings.Join(values, ",")
		if other {
			values[0] = "Other"
			lines[index] += "\nStyle: " + strings.Join(values, ",")
		}
	}
	source = strings.Join(lines, "\n")
	if secondCue {
		first := "0"
		if format == "ssa" {
			first = "Marked=0"
		}
		source += "Dialogue: " + first + ",0:00:04.00,0:00:05.00,Default,,0,0,0,,Another\n"
	}
	return source
}

func TestScriptedOutboundStyleEmphasisRequiresReadableRepresentation(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		fields := []string{"Bold", "Italic"}
		if format == "ass" {
			fields = append(fields, "Underline")
		}
		for _, field := range fields {
			tag := map[string]string{"Bold": "b", "Italic": "i", "Underline": "u"}[field]
			for _, target := range []string{"subrip", "webvtt"} {
				cases := []struct {
					name, value, text string
					other, secondCue  bool
					omitted           []bool
				}{
					{"inherited_true", "-1", "Hello", false, false, []bool{false}},
					{"neutral_baseline", "0", "Hello", false, false, []bool{false}},
					{"true_overridden_before_text", "-1", "{\\" + tag + "0}Hello", false, false, []bool{true}},
					{"false_overridden_before_text", "0", "{\\" + tag + "1}Hello", false, false, []bool{true}},
					{"matching_run_before_override", "-1", "A{\\" + tag + "0}Hello", false, false, []bool{false}},
					{"reset_restores_matching_run", "-1", "{\\" + tag + "0\\r}Hello", false, false, []bool{false}},
					{"reset_overridden_before_text", "-1", "{\\r\\" + tag + "0}Hello", false, false, []bool{true}},
					{"matching_other_cue", "-1", "{\\" + tag + "0}Hello", false, true, []bool{false}},
					{"drawing_is_not_readable", "-1", "{\\p1}m 0 0 l 1 1{\\p0\\" + tag + "0}Hello", false, false, []bool{true}},
					{"whitespace_is_not_readable", "-1", " \\N{\\" + tag + "0}Hello", false, false, []bool{true}},
					{"named_reset_has_readable_run", "-1", "{\\rOther}Hello", true, false, []bool{true, false}},
					{"named_reset_overridden_before_text", "-1", "{\\rOther\\" + tag + "0}Hello", true, false, []bool{true, true}},
					{"named_reset_without_readable_run", "-1", "{\\rOther\\r}Hello", true, false, []bool{false, true}},
				}
				for _, test := range cases {
					t.Run(format+"/"+field+"/"+target+"/"+test.name, func(t *testing.T) {
						source := scriptedAnalysisDocument(t, format, scriptedEmphasisLossSource(t, format, field, test.value, test.text, test.other, test.secondCue))
						before, _ := json.Marshal(source)
						result, err := Convert(context.Background(), source, target, Options{})
						if err != nil {
							t.Fatal(err)
						}
						native := source.FormatData.ASS
						if format == "ssa" {
							native = source.FormatData.SSA
						}
						for styleIndex, want := range test.omitted {
							styleOrder := -1
							for _, record := range native.Records {
								if record.RecordID == native.Styles[styleIndex].RecordID {
									styleOrder = record.SourceOrder
								}
							}
							fieldIndex := -1
							for index, nativeField := range native.Styles[styleIndex].Fields {
								if strings.EqualFold(nativeField.FieldName, field) {
									fieldIndex = index
								}
							}
							path := fmt.Sprintf("/format_data/%s/styles/%d/fields/%d", format, styleIndex, fieldIndex)
							found := 0
							for _, loss := range result.LossReport.Losses {
								if loss.Code == LossCodeScriptedStyleFieldOmitted && loss.Path == path {
									found++
									if loss.SourceOrder == nil || *loss.SourceOrder != styleOrder || loss.CueID != nil {
										t.Fatalf("style loss lost its native reference: %+v", loss)
									}
								}
							}
							if found > 1 || (found == 1) != want {
								t.Fatalf("%s omission count=%d want omission=%v", path, found, want)
							}
						}
						called := false
						renderer := func(context.Context, model.Document, model.Document, string) ([]byte, []model.Diagnostic, error) {
							called = true
							return []byte("unexpected"), nil, nil
						}
						strict, strictErr := convertWithRenderer(context.Background(), source, target, Options{Strict: true}, renderer)
						var refusal *StrictLossError
						if !errors.As(strictErr, &refusal) || called || len(strict.Bytes) != 0 || !reportsEqual(result.LossReport, strict.LossReport) {
							t.Fatalf("strict did not refuse with complete identical report: %+v %v", strict, strictErr)
						}
						after, _ := json.Marshal(source)
						if string(before) != string(after) {
							t.Fatal("source changed while tracking style representation")
						}
					})
				}
			}
		}
	}
}
