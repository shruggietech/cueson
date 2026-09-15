package convert

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/model"
)

func scriptedTextSource(t *testing.T, format, input string) model.Document {
	t.Helper()
	doc := conversionTestDocument(t, input, format)
	asset := &doc.Source.Assets[0]
	asset.DataBase64 = base64.StdEncoding.EncodeToString([]byte(input))
	asset.Size.Bytes = int64(len(input))
	asset.Size.Text = nil
	asset.Hashes.SHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
	return doc
}

func TestScriptedTargetsBaselineDeterminismAndHistoricalInput(t *testing.T) {
	for _, sourceFormat := range []string{"subrip", "webvtt"} {
		input := "1\n00:00:01,000 --> 00:00:02,500\nHello\n<b>world</b>\n"
		if sourceFormat == "webvtt" {
			input = "WEBVTT\n\n00:01.000 --> 00:02.500\nHello\n<b>world</b>\n"
		}
		for _, targetFormat := range []string{"ass", "ssa"} {
			t.Run(sourceFormat+"-"+targetFormat, func(t *testing.T) {
				source := scriptedTextSource(t, sourceFormat, input)
				source.Schema, source.SchemaVersion = "https://cueson.io/schema/v1.0.0/cueson.schema.json", "1.0.0"
				before, _ := json.Marshal(source)
				translation, err := translateTextToScripted(source, 0)
				if err != nil {
					t.Fatal(err)
				}
				target, losses, err := projectScriptedTarget(context.Background(), source, targetFormat, []payloadTranslation{translation})
				if err != nil {
					t.Fatal(err)
				}
				if len(losses) != 0 || len(translation.Issues) != 0 {
					t.Fatalf("baseline losses: %+v %+v", losses, translation.Issues)
				}
				if target.SchemaVersion != "1.1.0" || target.FormatSupport.Status != "stable" {
					t.Fatal("target identity")
				}
				if err := target.Validate(); err == nil {
					t.Fatal("public source grammar bypass")
				}
				result, err := scripted.RenderTarget(context.Background(), source, target)
				if err != nil {
					t.Fatal(err)
				}
				parsed, err := scripted.Parse(context.Background(), result.Bytes, targetFormat)
				if err != nil {
					t.Fatal(err)
				}
				if parsed.Cues[0].Payload.PlainText != "Hello\nworld" || parsed.Cues[0].Timing != source.Cues[0].Timing {
					t.Fatalf("target semantics: %+v", parsed.Cues[0])
				}
				if !strings.Contains(string(result.Bytes), `Hello\N{\b1}world{\b0}`) {
					t.Fatal("shared emphasis not represented")
				}
				repeated, err := scripted.RenderTarget(context.Background(), source, target)
				if err != nil || !bytes.Equal(result.Bytes, repeated.Bytes) {
					t.Fatal("nondeterminism")
				}
				after, _ := json.Marshal(source)
				if !bytes.Equal(before, after) {
					t.Fatal("input mutated")
				}
			})
		}
	}
}

func TestScriptedPrecisionEndpointsAndFatalBounds(t *testing.T) {
	for remainder := int64(0); remainder < 10; remainder++ {
		got, err := quantizeScriptedTiming(model.Timing{StartMilliseconds: 1000 + remainder, EndMilliseconds: 2000 + remainder})
		if err != nil {
			t.Fatal(err)
		}
		want := int64(1000)
		if remainder >= 5 {
			want += 10
		}
		if got.StartMilliseconds != want || got.EndMilliseconds != want+1000 {
			t.Fatalf("rounding %d: %+v", remainder, got)
		}
	}
	for _, timing := range []model.Timing{{StartMilliseconds: -1, EndMilliseconds: 100}, {StartMilliseconds: 1001, EndMilliseconds: 1004}, {StartMilliseconds: 1000, EndMilliseconds: math.MaxInt64}} {
		if _, err := quantizeScriptedTiming(timing); err == nil {
			t.Fatalf("unsafe timing accepted: %+v", timing)
		}
	}
	source := scriptedTextSource(t, "subrip", "1\n00:00:01,005 --> 00:00:02,004\nx\n")
	translation, err := translateTextToScripted(source, 0)
	if err != nil {
		t.Fatal(err)
	}
	target, losses, err := projectScriptedTarget(context.Background(), source, "ass", []payloadTranslation{translation})
	if err != nil {
		t.Fatal(err)
	}
	if len(losses) != 2 || losses[0].Code != LossCodeScriptedCentisecondQuantized || losses[1].Code != LossCodeScriptedCentisecondQuantized || target.Cues[0].Timing.StartMilliseconds != 1010 || target.Cues[0].Timing.EndMilliseconds != 2000 {
		t.Fatalf("endpoint losses: %+v %+v", losses, target.Cues)
	}
}

func TestScriptedLiteralAndNestedEmphasisPolicies(t *testing.T) {
	for _, literal := range []string{"{literal}", `literal \N`, `literal \n`, `literal \h`, "}\n"} {
		if _, err := encodeScriptedLiteral(literal); err == nil {
			t.Fatalf("control literal accepted %q", literal)
		}
	}
	for _, text := range []string{`ordinary \x and comma,`, "line1\n\nline3\u00a0tail", "<b>outer<b>inner</b>tail</b>"} {
		source := scriptedTextSource(t, "subrip", "1\n00:00:01,000 --> 00:00:02,000\n"+strings.ReplaceAll(text, "\n\n", "\n<i></i>\n")+"\n")
		translation, err := translateTextToScripted(source, 0)
		if err != nil {
			t.Fatal(err)
		}
		facts, err := model.ProjectScriptedText(translation.Text, 0, 1000, 2000)
		if err != nil || len(facts.DiagnosticCodes) != 0 {
			t.Fatalf("literal facts: %+v %v", facts, err)
		}
		if strings.Contains(text, "outer") && translation.Text != `{\b1}outerinnertail{\b0}` {
			t.Fatalf("nested state lost: %s", translation.Text)
		}
	}
	source := scriptedTextSource(t, "webvtt", "WEBVTT\n\n00:01.000 --> 00:02.000\n&lt;b&gt;&amp; &nbsp;\n")
	translation, err := translateTextToScripted(source, 0)
	if err != nil {
		t.Fatal(err)
	}
	facts, err := model.ProjectScriptedText(translation.Text, 0, 1000, 2000)
	if err != nil || strings.Join(facts.Lines, "\n") != "<b>& \u00a0" {
		t.Fatalf("entity literal reinterpreted: %+v %v", facts, err)
	}
}

func TestScriptedVariantsPreserveNativeOwnersAndAccountFields(t *testing.T) {
	for _, sourceFormat := range []string{"ass", "ssa"} {
		t.Run(sourceFormat, func(t *testing.T) {
			targetFormat := "ass"
			if sourceFormat == "ass" {
				targetFormat = "ssa"
			}
			input := scriptedConversionSource(sourceFormat, `Hello{\i1} world{\i0}`) + "; retained comment\n[Fonts]\nfontname: embedded.ttf\n!!!!\n\n"
			source := scriptedAnalysisDocument(t, sourceFormat, input)
			before, _ := json.Marshal(source)
			target, losses, err := projectScriptedVariant(context.Background(), source, targetFormat)
			if err != nil {
				t.Fatal(err)
			}
			if len(losses) == 0 {
				t.Fatal("variant field losses missing")
			}
			result, err := scripted.RenderTarget(context.Background(), source, target)
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := scripted.Parse(context.Background(), result.Bytes, targetFormat)
			if err != nil {
				t.Fatal(err)
			}
			if parsed.Cues[0].Timing != source.Cues[0].Timing || parsed.Cues[0].Payload.RawText != source.Cues[0].Payload.RawText || len(parsed.Native.Attachments) != 1 || !strings.Contains(string(result.Bytes), "; retained comment\n") {
				t.Fatal("native owner lost")
			}
			if !strings.HasSuffix(string(result.Bytes), "!!!!\n\n") {
				t.Fatal("trailing blank source record lost")
			}
			after, _ := json.Marshal(source)
			if !bytes.Equal(before, after) {
				t.Fatal("variant mutated source")
			}
		})
	}
}

func TestScriptedVariantAlignmentAndColorRoles(t *testing.T) {
	mapping := map[string]string{"1": "1", "2": "2", "3": "3", "5": "7", "6": "8", "7": "9", "9": "4", "10": "5", "11": "6"}
	for ssa, ass := range mapping {
		source := scriptedAnalysisDocument(t, "ssa", scriptedConversionSource("ssa", "x"))
		style := source.FormatData.SSA.Styles[0]
		for index := range style.Fields {
			if strings.EqualFold(style.Fields[index].FieldName, "Alignment") {
				style.Fields[index].RawValue = ssa
				style.Fields[index].TypedValue = nil
			}
		}
		losses := []Loss{}
		fields, err := scriptedVariantStyle(source, "ass", 0, style, 5, &losses)
		if err != nil {
			t.Fatal(err)
		}
		for _, field := range fields {
			if field.FieldName == "Alignment" && field.RawValue != ass {
				t.Fatalf("alignment %s: %s", ssa, field.RawValue)
			}
			if field.FieldName == "BackColour" && field.RawValue != "&H80000000" {
				t.Fatal("SSA shadow alpha not mapped")
			}
		}
	}
}

func TestScriptedVariantAliasAndUnsupportedControls(t *testing.T) {
	input := strings.ReplaceAll(scriptedConversionSource("ass", "x"), ",Name,", ",Actor,")
	input = strings.ReplaceAll(input, "Default,,0,0,0,,x", "Default,Narrator,0,0,0,,x")
	source := scriptedAnalysisDocument(t, "ass", input)
	target, _, err := projectScriptedVariant(context.Background(), source, "ssa")
	if err != nil {
		t.Fatal(err)
	}
	if target.Cues[0].FormatData.SSA.Projection.SpeakerFieldName == nil || *target.Cues[0].FormatData.SSA.Projection.SpeakerFieldName != "Name" || target.Cues[0].Speakers[0].Name != "Narrator" {
		t.Fatal("actor alias linkage was not normalized to constructed target field")
	}
	unknown := scriptedAnalysisDocument(t, "ass", scriptedConversionSource("ass", `unknown{\foo1}control`))
	if _, _, err := projectScriptedVariant(context.Background(), unknown, "ssa"); err == nil {
		t.Fatal("unknown control target meaning was assumed")
	}
}

func TestScriptedVariantEditedCommonPrecisionPolicy(t *testing.T) {
	source := scriptedAnalysisDocument(t, "ass", scriptedConversionSource("ass", "x"))
	source.Cues[0].Timing = model.Timing{StartMilliseconds: 1005, EndMilliseconds: 3004, DurationMilliseconds: 1999}
	recomputeSummaries(&source)
	if err := source.Validate(); err != nil {
		t.Fatal(err)
	}
	target, losses, err := projectScriptedVariant(context.Background(), source, "ssa")
	if err != nil {
		t.Fatal(err)
	}
	quantized := 0
	for _, loss := range losses {
		if loss.Code == LossCodeScriptedCentisecondQuantized {
			quantized++
		}
	}
	if quantized != 2 || target.Cues[0].Timing.StartMilliseconds != 1010 || target.Cues[0].Timing.EndMilliseconds != 3000 {
		t.Fatalf("edited precision: %+v %+v", target.Cues[0].Timing, losses)
	}
	called := false
	renderer := func(context.Context, model.Document, model.Document, string) ([]byte, []model.Diagnostic, error) {
		called = true
		return []byte("unexpected"), nil, nil
	}
	result, err := convertWithRenderer(context.Background(), source, "ssa", Options{Strict: true}, renderer)
	if err == nil || called || len(result.Bytes) != 0 {
		t.Fatal("strict edited variant invoked renderer")
	}
	source.Cues[0].Timing = model.Timing{StartMilliseconds: 1001, EndMilliseconds: 1004, DurationMilliseconds: 3}
	recomputeSummaries(&source)
	if _, _, err := projectScriptedVariant(context.Background(), source, "ssa"); err == nil {
		t.Fatal("edited interval collapsed silently")
	}
}

func TestCaptureFreeScriptedTargetAcceptsPreservedUTF16Source(t *testing.T) {
	input := "1\n00:00:01,000 --> 00:00:02,000\nHello\n"
	source := scriptedTextSource(t, "subrip", input)
	encoded := []byte{0xff, 0xfe}
	for _, unit := range utf16.Encode([]rune(input)) {
		encoded = append(encoded, byte(unit), byte(unit>>8))
	}
	asset := &source.Source.Assets[0]
	asset.DataBase64 = base64.StdEncoding.EncodeToString(encoded)
	asset.Size.Bytes = int64(len(encoded))
	asset.Hashes.SHA256 = fmt.Sprintf("%x", sha256.Sum256(encoded))
	asset.Encoding = nil
	translation, err := translateTextToScripted(source, 0)
	if err != nil {
		t.Fatal(err)
	}
	target, _, err := projectScriptedTarget(context.Background(), source, "ass", []payloadTranslation{translation})
	if err != nil {
		t.Fatal(err)
	}
	result, err := scripted.RenderTarget(context.Background(), source, target)
	if err != nil || !strings.Contains(string(result.Bytes), ",Hello\n") {
		t.Fatalf("preserved UTF16 source rejected: %v", err)
	}
	if source.Source.Assets[0].DataBase64 != base64.StdEncoding.EncodeToString(encoded) {
		t.Fatal("source bytes changed")
	}
}

func TestScriptedWebVTTEncodedMarkupRemainsLiteral(t *testing.T) {
	for _, targetFormat := range []string{"ass", "ssa"} {
		for _, text := range []string{"&lt;b&gt;literal&lt;/b&gt;", "unterminated &amp"} {
			source := scriptedTextSource(t, "webvtt", "WEBVTT\n\n00:01.000 --> 00:02.000\n"+text+"\n")
			translation, err := translateTextToScripted(source, 0)
			if err != nil || len(translation.Issues) != 0 {
				t.Fatalf("false target ambiguity: %+v %v", translation, err)
			}
			result, err := Convert(context.Background(), source, targetFormat, Options{Strict: true})
			if err != nil || result.LossReport.HasLosses() {
				t.Fatalf("literal conversion: %+v %v", result.LossReport, err)
			}
			parsed, err := scripted.Parse(context.Background(), result.Bytes, targetFormat)
			if err != nil || parsed.Cues[0].Payload.PlainText != source.Cues[0].Payload.PlainText {
				t.Fatalf("literal changed: %+v %v", parsed.Cues, err)
			}
		}
	}
}

func TestScriptedVariantNativeInertAlphaObservationIsAccounted(t *testing.T) {
	input := strings.ReplaceAll(scriptedConversionSource("ssa", "x"), "Style: Default,Arial,20,16777215,255,0,0,", "Style: Default,Arial,20,553648127,570425599,0,587202560,")
	source := scriptedAnalysisDocument(t, "ssa", input)
	before, _ := json.Marshal(source)
	target, losses, err := projectScriptedVariant(context.Background(), source, "ass")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, loss := range losses {
		if loss.Code == LossCodeScriptedVariantFieldDegraded {
			count++
		}
	}
	if count != 3 {
		t.Fatalf("inert alpha observations disappeared: %+v", losses)
	}
	for _, field := range target.FormatData.ASS.Styles[0].Fields {
		if field.FieldName == "PrimaryColour" && field.RawValue != "&H00FFFFFF" {
			t.Fatal("active global alpha mapping incorrect")
		}
		if field.FieldName == "BackColour" && field.RawValue != "&H80000000" {
			t.Fatal("active fixed shadow alpha mapping incorrect")
		}
	}
	after, _ := json.Marshal(source)
	if !bytes.Equal(before, after) {
		t.Fatal("native alpha observations mutated")
	}
}

func TestScriptedVariantIdentifierOmissionKeepsSourceImmutable(t *testing.T) {
	for _, sourceFormat := range []string{"ass", "ssa"} {
		targetFormat := "ssa"
		if sourceFormat == "ssa" {
			targetFormat = "ass"
		}
		source := scriptedAnalysisDocument(t, sourceFormat, scriptedConversionSource(sourceFormat, "x"))
		identifier := "operator-assigned-label"
		source.Cues[0].SourceIdentifier = &identifier
		before, _ := json.Marshal(source)
		target, losses, err := projectScriptedVariant(context.Background(), source, targetFormat)
		if err != nil {
			t.Fatal(err)
		}
		count := 0
		for _, loss := range losses {
			if loss.Code == LossCodeSourceIdentifierOmitted && loss.Path == "/cues/0/source_identifier" {
				count++
			}
		}
		if count != 1 || target.Cues[0].SourceIdentifier != nil {
			t.Fatalf("identifier omission incomplete: %+v", losses)
		}
		called := false
		renderer := func(context.Context, model.Document, model.Document, string) ([]byte, []model.Diagnostic, error) {
			called = true
			return []byte("unexpected"), nil, nil
		}
		result, err := convertWithRenderer(context.Background(), source, targetFormat, Options{Strict: true}, renderer)
		if err == nil || called || len(result.Bytes) != 0 {
			t.Fatal("strict identifier conversion published output")
		}
		count = 0
		for _, loss := range result.LossReport.Losses {
			if loss.Code == LossCodeSourceIdentifierOmitted {
				count++
			}
		}
		if count != 1 {
			t.Fatal("strict report lost identifier omission")
		}
		after, _ := json.Marshal(source)
		if !bytes.Equal(before, after) {
			t.Fatal("identifier annotation changed in original source")
		}
	}
}
