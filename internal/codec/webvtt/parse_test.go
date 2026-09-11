package webvtt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func TestDetectAndDecodeUTF8Boundary(t *testing.T) {
	payload := append([]byte{0xef, 0xbb, 0xbf}, []byte("WEBVTT\r\n\r\n00:00.000 --> 00:01.000\r\nx\r\n")...)
	if !Detect(payload) || Detect([]byte("WEBVTTbad\n")) {
		t.Fatal("signature detection did not enforce its boundary")
	}
	decoded, err := DecodeUTF8(payload, "utf-8")
	if err != nil || strings.HasPrefix(decoded.Text, "\ufeff") || decoded.Observation.LineEndings != "crlf" || decoded.Observation.BOM == nil {
		t.Fatalf("DecodeUTF8() = (%#v, %v)", decoded, err)
	}
	if _, err := DecodeUTF8([]byte{0xff, 0xfe}, ""); err == nil {
		t.Fatal("DecodeUTF8() accepted invalid UTF-8")
	}
	if _, err := DecodeUTF8([]byte("WEBVTT"), "windows-1252"); err == nil {
		t.Fatal("DecodeUTF8() accepted a legacy override")
	}
}

func TestDecodeUTF8EnforcesEveryBOMSpecificAlias(t *testing.T) {
	bomless := []byte("WEBVTT\n")
	withBOM := append([]byte{0xef, 0xbb, 0xbf}, bomless...)
	for _, alias := range []string{"utf-8-bom", "utf8-bom", "utf-8-sig"} {
		if _, err := DecodeUTF8(bomless, alias); err == nil {
			t.Errorf("DecodeUTF8() accepted BOM-less input for %q", alias)
		}
		if _, err := DecodeUTF8(withBOM, alias); err != nil {
			t.Errorf("DecodeUTF8() rejected BOM input for %q: %v", alias, err)
		}
	}
	for _, alias := range []string{"utf-8", "utf8"} {
		if _, err := DecodeUTF8(bomless, alias); err != nil {
			t.Errorf("DecodeUTF8() rejected BOM-less input for %q: %v", alias, err)
		}
	}
}

func TestParsePreservesHeaderBlocksCueOrderAndRawFields(t *testing.T) {
	input := "WEBVTT sample\rX-TIMESTAMP-MAP=LOCAL:00:00:00.000,MPEGTS:0\r\rSTYLE\r::cue { color: lime }\r\rREGION\rid:captions width:80%\r\rNOTE source\rkept\r\rid-one\r00:01.000 --> 00:03.000 line:10% align:start\r<v Alice>Hello &amp; welcome</v>\r"
	result, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if result.DocumentData.SignatureLineRaw != "WEBVTT sample" || result.DocumentData.Description == nil || *result.DocumentData.Description != "sample" || len(result.DocumentData.MetadataLines) != 1 {
		t.Fatalf("document data = %#v", result.DocumentData)
	}
	if len(result.DocumentData.Blocks) != 3 || len(result.Cues) != 1 || result.Cues[0].SourceOrder != 3 {
		t.Fatalf("blocks=%#v cues=%#v", result.DocumentData.Blocks, result.Cues)
	}
	cue := result.Cues[0]
	if cue.SourceIdentifier == nil || *cue.SourceIdentifier != "id-one" || cue.Payload.RawText != "<v Alice>Hello &amp; welcome</v>" || cue.Payload.PlainText != "Hello & welcome" || len(cue.Speakers) != 1 {
		t.Fatalf("cue = %#v", cue)
	}
	if cue.FormatData.WebVTT == nil || len(cue.FormatData.WebVTT.SettingOccurrences) != 2 || cue.FormatData.WebVTT.RawPayloadLines[0] != cue.Payload.RawText {
		t.Fatalf("native cue = %#v", cue.FormatData.WebVTT)
	}
}

func TestParseToleranceIsDiagnosedWithoutSortingOrDeduplication(t *testing.T) {
	input := "WEBVTT\n\nsame\n00:02.000 --> 00:04.000 mystery:x\nroll\nsame\n00:01.000 --> 00:03.000\nroll\n\nSTYLE\n::cue { color: red }\n"
	result, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Cues) != 2 || result.Cues[0].Timing.StartMilliseconds != 2000 || result.Cues[1].Timing.StartMilliseconds != 1000 {
		t.Fatalf("cues were lost or reordered: %#v", result.Cues)
	}
	codes := make([]string, len(result.Diagnostics))
	for index := range result.Diagnostics {
		codes[index] = result.Diagnostics[index].Code
	}
	want := []string{"webvtt_setting_unknown", "webvtt_separator_missing", "webvtt_cue_identifier_duplicate", "webvtt_start_order_invalid", "webvtt_style_after_cue"}
	if strings.Join(codes, ",") != strings.Join(want, ",") {
		t.Fatalf("diagnostic codes = %v, want %v", codes, want)
	}
}

func TestParseRejectsFatalSignatureAndTiming(t *testing.T) {
	for _, input := range []string{"NOTVTT\n", "WEBVTT bad --> header\n", "WEBVTT\n\n00:00.000 --> 00:00.000\nx\n", "WEBVTT\n\n00:00.000 --> 00:01.00\nx\n"} {
		if _, err := Parse(input); err == nil {
			t.Fatalf("Parse(%q) unexpectedly succeeded", input)
		}
	}
}

func TestGovernedToleratedFixturePreservesEveryDiagnosedConstruct(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "fixtures", "webvtt", "tolerated", "source", "tolerated.vtt")
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeUTF8(source, "")
	if err != nil {
		t.Fatal(err)
	}
	result, err := Parse(decoded.Text)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Cues) != 2 || len(result.DocumentData.Blocks) != 3 {
		t.Fatalf("cues=%d blocks=%d", len(result.Cues), len(result.DocumentData.Blocks))
	}
	want := []string{
		"webvtt_setting_invalid",
		"webvtt_setting_duplicate",
		"webvtt_setting_unknown",
		"webvtt_entity_invalid",
		"webvtt_markup_unbalanced",
		"webvtt_separator_missing",
		"webvtt_cue_identifier_duplicate",
		"webvtt_start_order_invalid",
		"webvtt_nul_replaced",
		"webvtt_style_after_cue",
		"webvtt_region_after_cue",
		"webvtt_block_unrecognized",
	}
	for _, code := range want {
		found := false
		for _, diagnostic := range result.Diagnostics {
			if diagnostic.Code == code {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing diagnostic %s: %#v", code, result.Diagnostics)
		}
	}
	if !strings.ContainsRune(result.Cues[1].Payload.RawText, '\x00') || !strings.Contains(result.Cues[1].Payload.PlainText, "\ufffd") || !strings.ContainsRune(string(source), '\x00') {
		t.Fatalf("NUL source and semantic replacement were not kept distinct")
	}
}

func TestParsePreservesNULInNativeFieldsAndReplacesOnlyDerivedText(t *testing.T) {
	input := "WEBVTT desc\x00value\n\nNOTE raw\x00note\nkept\n\nid\x00raw\n00:00.000 --> 00:01.000\npayload\x00text\n"
	result, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.ContainsRune(result.DocumentData.SignatureLineRaw, '\x00') || result.DocumentData.Description == nil || !strings.Contains(*result.DocumentData.Description, "\ufffd") {
		t.Fatalf("signature=%q description=%v", result.DocumentData.SignatureLineRaw, result.DocumentData.Description)
	}
	if !strings.ContainsRune(result.DocumentData.Blocks[0].Raw, '\x00') || !strings.ContainsRune(*result.Cues[0].FormatData.WebVTT.IdentifierRaw, '\x00') || !strings.ContainsRune(result.Cues[0].Payload.RawText, '\x00') {
		t.Fatalf("native fields lost NUL: %#v %#v", result.DocumentData.Blocks[0], result.Cues[0])
	}
	if strings.ContainsRune(result.Cues[0].Payload.PlainText, '\x00') || !strings.Contains(result.Cues[0].Payload.PlainText, "\ufffd") {
		t.Fatalf("plain_text = %q", result.Cues[0].Payload.PlainText)
	}
}

func TestParseRejectsOccurrenceAndDiagnosticAmplification(t *testing.T) {
	t.Parallel()

	settings := strings.Repeat("x:y ", model.MaxItemOccurrences+1)
	input := "WEBVTT\n\n00:00.000 --> 00:01.000 " + settings + "\nx\n"
	if _, err := Parse(input); err == nil || !strings.Contains(err.Error(), "webvtt_occurrence_limit_exceeded") {
		t.Fatalf("setting limit error = %v", err)
	}

	input = "WEBVTT\n\n" + strings.Repeat("unknown\n\n", model.MaxDiagnostics+1) + "00:00.000 --> 00:01.000\nx\n"
	if _, err := Parse(input); err == nil || !strings.Contains(err.Error(), "webvtt_diagnostic_limit_exceeded") {
		t.Fatalf("diagnostic limit error = %v", err)
	}
}
