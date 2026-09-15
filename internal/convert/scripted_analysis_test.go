package convert

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func scriptedConversionSource(format, text string) string {
	dialect, section, first := "v4.00+", "V4+ Styles", "0"
	style := "Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,2,1,2,10,10,10,1"
	if format == "ssa" {
		dialect, section, first = "v4.00", "V4 Styles", "Marked=0"
		style = "Default,Arial,20,16777215,255,0,0,0,0,1,2,1,2,10,10,10,0,1"
	}
	return "[Script Info]\nScriptType: " + dialect + "\nWrapStyle: 0\n[" + section + "]\nFormat: " + strings.Join(model.ScriptedCanonicalFields(format, "style"), ",") + "\nStyle: " + style + "\n[Events]\nFormat: " + strings.Join(model.ScriptedCanonicalFields(format, "event"), ",") + "\nDialogue: " + first + ",0:00:01.00,0:00:03.00,Default,,0,0,0,," + text + "\n"
}

func scriptedAnalysisDocument(t *testing.T, format, input string) model.Document {
	t.Helper()
	parsed, err := scripted.Parse(context.Background(), []byte(input), format)
	if err != nil {
		t.Fatal(err)
	}
	document, err := schema.Decode(schema.Representative())
	if err != nil {
		t.Fatal(err)
	}
	document.Format = format
	document.FormatSupport = model.FormatSupport{Status: "experimental", IngestSupported: true, RenderSupported: true, RestoreSupported: true}
	document.Metadata = model.Metadata{}
	document.Cues = parsed.Cues
	document.FormatData = model.DocumentFormatData{}
	if format == "ass" {
		document.FormatData.ASS = &parsed.Native
	} else {
		document.FormatData.SSA = &parsed.Native
	}
	document.Diagnostics = parsed.Diagnostics
	for index := range document.Source.Assets {
		asset := &document.Source.Assets[index]
		if asset.ID != document.Source.PrimaryAssetID {
			continue
		}
		asset.DataBase64 = base64.StdEncoding.EncodeToString([]byte(input))
		asset.Size.Bytes = int64(len(input))
		asset.Size.Text = nil
		asset.Hashes.SHA256 = fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
		asset.Encoding = &parsed.Encoding
	}
	recomputeSummaries(&document)
	document.Stats.DiagnosticCount = len(document.Diagnostics)
	document.Stats.WarningCount, document.Stats.ErrorCount = 0, 0
	for _, diagnostic := range document.Diagnostics {
		if diagnostic.Severity == "warning" {
			document.Stats.WarningCount++
		}
		if diagnostic.Severity == "error" {
			document.Stats.ErrorCount++
		}
	}
	for _, cue := range document.Cues {
		if len(cue.Tokens) > 0 {
			document.Document.HasWordLevelTiming = true
			document.Stats.HasWordLevelTiming = true
		}
	}
	if err := document.Validate(); err != nil {
		t.Fatalf("invalid test source: %v", err)
	}
	return document
}

func TestScriptedOutboundBalancedEmphasisAndLiteralEscaping(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		for _, target := range []string{"subrip", "webvtt"} {
			t.Run(format+"_"+target, func(t *testing.T) {
				document := scriptedAnalysisDocument(t, format, scriptedConversionSource(format, `A{\b1}B{\i1}C{\b0}D{\u1}E{\r}F\Ncan't & speak`))
				before, _ := json.Marshal(document)
				analysis, err := analyzeScriptedSource(context.Background(), document, target)
				if err != nil {
					t.Fatal(err)
				}
				want := "A<b>B</b><b><i>C</i></b><i>D</i><i><u>E</u></i>F\ncan't & speak"
				if target == "webvtt" {
					want = strings.ReplaceAll(want, "&", "&amp;")
				}
				if analysis.Translations[0].Text != want {
					t.Fatalf("payload=%q want=%q", analysis.Translations[0].Text, want)
				}
				report, err := NewReport(document, analysis.Losses)
				if err != nil {
					t.Fatal(err)
				}
				if !report.HasLosses() {
					t.Fatal("native presentation losses missing")
				}
				projected, err := projectDocument(document, target, analysis.Translations)
				if err != nil {
					t.Fatal(err)
				}
				output, _, err := renderProjected(context.Background(), document, projected, target)
				if err != nil || len(output) == 0 {
					t.Fatalf("render=%q err=%v", output, err)
				}
				after, _ := json.Marshal(document)
				if !bytes.Equal(before, after) {
					t.Fatal("source document changed")
				}
			})
		}
	}
}

func TestScriptedOutboundInheritedAndNamedResetEmphasis(t *testing.T) {
	input := scriptedConversionSource("ass", `A{\i0}B{\rOther}C{\r}D`)
	input = strings.Replace(input, "Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0", "Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,1,0", 1)
	input = strings.Replace(input, "[Events]", "Style: Other,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,1,0,1,0,100,100,0,0,1,2,1,2,10,10,10,1\n[Events]", 1)
	document := scriptedAnalysisDocument(t, "ass", input)
	translation, err := translateScriptedPayload(document, 0, "subrip")
	if err != nil {
		t.Fatal(err)
	}
	if translation.Text != "<i>A</i>B<b><u>C</u></b><i>D</i>" {
		t.Fatalf("payload=%q", translation.Text)
	}
}

func TestScriptedOutboundWrapSpacesAndWeightedEmphasis(t *testing.T) {
	document := scriptedAnalysisDocument(t, "ass", scriptedConversionSource("ass", `A\nB{\q2}C\nD\hE{\b400}F{\b700}G`))
	analysis, err := analyzeScriptedSource(context.Background(), document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	if analysis.Translations[0].Text != "A BC\nD\u00a0EF<b>G</b>" {
		t.Fatalf("projection=%q", analysis.Translations[0].Text)
	}
	degraded := 0
	for _, loss := range analysis.Losses {
		if loss.Code == LossCodeScriptedOverrideDegraded {
			degraded++
		}
	}
	if degraded != 2 {
		t.Fatalf("weighted emphasis losses=%d", degraded)
	}
}

func TestScriptedOutboundCompleteAtomicNativeLosses(t *testing.T) {
	input := scriptedConversionSource("ass", `{comment\k10}A{\p1}m 0 0 l 1 1{\p0}B{\pos(1,2)}C`)
	input += "Comment: 0,0:00:01.00,0:00:03.00,Default,,0,0,0,,Note\nkept inert extension\n[Extra]\n; inert\n[Fonts]\nfontname: bundle.ttf\n!!!!\n"
	document := scriptedAnalysisDocument(t, "ass", input)
	analysis, err := analyzeScriptedSource(context.Background(), document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	report, err := NewReport(document, analysis.Losses)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{}
	for _, loss := range report.Losses {
		counts[loss.Code]++
		if loss.SourceOrder == nil || strings.Contains(loss.Message, "bundle") || strings.Contains(loss.Message, "inert") {
			t.Fatalf("unsafe/unresolved native loss: %+v", loss)
		}
	}
	wantCounts := map[string]int{LossCodeScriptedMetadataOmitted: 1, LossCodeScriptedStyleFieldOmitted: 20, LossCodeScriptedEventFieldOmitted: 7, LossCodeScriptedRecordOmitted: 3, LossCodeScriptedSectionOmitted: 1, LossCodeScriptedAttachmentOmitted: 1, LossCodeScriptedOverrideOmitted: 4, LossCodeScriptedOverrideCommentOmitted: 1, LossCodeScriptedDrawingOmitted: 1, LossCodeTokenTimingOmitted: 1}
	for code, want := range wantCounts {
		if counts[code] != want {
			t.Fatalf("%s count=%d want=%d full=%+v", code, counts[code], want, report)
		}
	}
	if analysis.Translations[0].Text != "ABC" {
		t.Fatalf("mixed drawing invented text: %q", analysis.Translations[0].Text)
	}
	if len(report.Losses) != 40 {
		t.Fatalf("complete report count=%d", len(report.Losses))
	}
	repeated, err := analyzeScriptedSource(context.Background(), document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	repeatedReport, err := NewReport(document, repeated.Losses)
	if err != nil || !reportsEqual(report, repeatedReport) {
		t.Fatal("report not deterministic")
	}
}

func TestScriptedOutboundFatalTextBoundaries(t *testing.T) {
	for _, text := range []string{"", `{\p1}m 0 0 l 1 1{\p0}`, `{\i2}A`, `{\b(1)}A`, `A{unclosed`} {
		t.Run(text, func(t *testing.T) {
			document := scriptedAnalysisDocument(t, "ass", scriptedConversionSource("ass", text))
			for _, target := range []string{"subrip", "webvtt"} {
				analysis, err := analyzeScriptedSource(context.Background(), document, target)
				if err == nil || len(analysis.Translations) != 0 {
					t.Fatalf("unsafe payload=%q target=%s result=%+v err=%v", text, target, analysis, err)
				}
			}
		})
	}
	document := scriptedAnalysisDocument(t, "ass", scriptedConversionSource("ass", `literal <b>text</b>`))
	if _, err := analyzeScriptedSource(context.Background(), document, "subrip"); err == nil {
		t.Fatal("SubRip literal reinterpretation accepted")
	}
	analysis, err := analyzeScriptedSource(context.Background(), document, "webvtt")
	if err != nil || analysis.Translations[0].Text != "literal &lt;b&gt;text&lt;/b&gt;" {
		t.Fatalf("WebVTT escaping=%+v err=%v", analysis, err)
	}
}

func TestScriptedOutboundUnknownIdentityIsNotGuessedAsKnownTag(t *testing.T) {
	document := scriptedAnalysisDocument(t, "ass", scriptedConversionSource("ass", `{\rMissing\bno}A`))
	analysis, err := analyzeScriptedSource(context.Background(), document, "webvtt")
	if err != nil {
		t.Fatal(err)
	}
	if analysis.Translations[0].Text != "A" {
		t.Fatalf("unknown identity changed dialogue: %q", analysis.Translations[0].Text)
	}
	count := 0
	for _, loss := range analysis.Losses {
		if loss.Code == LossCodeScriptedOverrideOmitted {
			count++
		}
	}
	if count != 2 {
		t.Fatalf("unknown identities losses=%d", count)
	}
}

func TestScriptedOutboundLossBoundAndCancellation(t *testing.T) {
	input := scriptedConversionSource("ass", "A")
	for index := 0; index < MaxLosses/7; index++ {
		input += "Dialogue: 0,0:00:03.00,0:00:04.00,Default,,0,0,0,,A\n"
	}
	document := scriptedAnalysisDocument(t, "ass", input)
	if analysis, err := analyzeScriptedSource(context.Background(), document, "webvtt"); err == nil || len(analysis.Losses) != 0 {
		t.Fatalf("loss limit failed result=%+v err=%v", analysis, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := analyzeScriptedSource(ctx, document, "webvtt"); err != context.Canceled {
		t.Fatalf("cancellation=%v", err)
	}
	//lint:ignore SA1012 Intentional invalid-context boundary regression.
	if _, err := analyzeScriptedSource(nil, document, "webvtt"); err == nil {
		t.Fatal("nil context accepted")
	}
}
