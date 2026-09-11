package subrip

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParsePreservesNativeAndDerivedViews(t *testing.T) {
	t.Parallel()
	input := "7\r\n00:00:01.2 --> 00:00:02,34 X1:10 X2:20 Y1:30 Y2:40\r\n<i>Alice: Hello</i>  \r\nsecond line\r\n\r\n00:00:03,000 --> 00:00:04,000\r\nNo sequence\r\n"
	result, err := Parse(input, Options{DetectSpeakers: true})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Cues) != 2 {
		t.Fatalf("cue count = %d, want 2", len(result.Cues))
	}
	first := result.Cues[0]
	if first.SourceIdentifier == nil || *first.SourceIdentifier != "7" {
		t.Fatalf("source identifier = %#v, want 7", first.SourceIdentifier)
	}
	if first.FormatData.SubRip.SequenceLineRaw != "7" || first.FormatData.SubRip.TimingLineRaw != "00:00:01.2 --> 00:00:02,34 X1:10 X2:20 Y1:30 Y2:40" {
		t.Fatalf("native data = %#v", first.FormatData.SubRip)
	}
	if first.Payload.RawText != "<i>Alice: Hello</i>  \nsecond line" {
		t.Fatalf("raw text = %q", first.Payload.RawText)
	}
	if first.Payload.PlainText != "Alice: Hello  \nsecond line" {
		t.Fatalf("plain text = %q", first.Payload.PlainText)
	}
	if !reflect.DeepEqual(first.Payload.Lines, []string{"<i>Alice: Hello</i>  ", "second line"}) {
		t.Fatalf("lines = %#v", first.Payload.Lines)
	}
	if len(first.Speakers) != 0 {
		t.Fatalf("tag-wrapped speaker should not be inferred: %#v", first.Speakers)
	}
	second := result.Cues[1]
	if second.SourceIdentifier != nil || second.FormatData.SubRip.SequenceLineRaw != "" {
		t.Fatalf("missing sequence was not preserved: %#v", second)
	}
	if !hasDiagnostic(result.Diagnostics, "subrip_sequence_irregular") || !hasDiagnostic(result.Diagnostics, "subrip_sequence_missing") {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func TestParsePreservesEmptyPayloadLineAndRecoversMissingSeparator(t *testing.T) {
	t.Parallel()
	input := "1\n00:00:00,000 --> 00:00:01,000\nline one\n\nline three\n2\n00:00:01,000 --> 00:00:02,000\nBob: hi\n"
	result, err := Parse(input, Options{DetectSpeakers: true})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if got, want := result.Cues[0].Payload.Lines, []string{"line one", "", "line three"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("first lines = %#v, want %#v", got, want)
	}
	if len(result.Cues[1].Speakers) != 1 || result.Cues[1].Speakers[0].Name != "Bob" || result.Cues[1].Payload.RawText != "Bob: hi" {
		t.Fatalf("speaker cue = %#v", result.Cues[1])
	}
	if !hasDiagnostic(result.Diagnostics, "subrip_separator_missing") {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func TestParsePreservesInternalBlankPayloadLineBeforeSeparator(t *testing.T) {
	t.Parallel()
	input := "1\n00:00:00,000 --> 00:00:01,000\nfirst\n\nthird\n\n2\n00:00:01,000 --> 00:00:02,000\nsecond\n"
	result, err := Parse(input, Options{})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Cues) != 2 {
		t.Fatalf("cue count = %d, want 2", len(result.Cues))
	}
	if got, want := result.Cues[0].Payload.Lines, []string{"first", "", "third"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("first lines = %#v, want %#v", got, want)
	}
	rendered, err := Render(result.Cues)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	reparsed, err := Parse(string(rendered), Options{})
	if err != nil {
		t.Fatalf("Parse(rendered) error = %v", err)
	}
	if got, want := reparsed.Cues[0].Payload.Lines, result.Cues[0].Payload.Lines; !reflect.DeepEqual(got, want) {
		t.Fatalf("reparsed first lines = %#v, want %#v", got, want)
	}
}

func TestParsePreservesLeadingUFEFFAsContent(t *testing.T) {
	t.Parallel()
	input := "\ufeff1\n00:00:00,000 --> 00:00:01,000\ntext\n"
	result, err := Parse(input, Options{})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if result.Cues[0].SourceIdentifier == nil || *result.Cues[0].SourceIdentifier != "\ufeff1" {
		t.Fatalf("source identifier = %#v, want literal U+FEFF followed by 1", result.Cues[0].SourceIdentifier)
	}
	if result.Cues[0].FormatData.SubRip.SequenceLineRaw != "\ufeff1" {
		t.Fatalf("sequence line = %q, want literal U+FEFF followed by 1", result.Cues[0].FormatData.SubRip.SequenceLineRaw)
	}
}

func TestParseAccountsForUnrecognizedBlocks(t *testing.T) {
	t.Parallel()
	input := "unmodeled header\n\n1\n00:00:01,000 --> 00:00:02,000\nhello\n\ninterstitial block\n\n2\n00:00:02,000 --> 00:00:03,000\nworld\n"
	result, err := Parse(input, Options{})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Cues) != 2 || result.Cues[0].Payload.RawText != "hello\n\ninterstitial block" || !hasDiagnostic(result.Diagnostics, "subrip_block_unrecognized") {
		t.Fatalf("result = %#v", result)
	}
}

func TestParsePreservesIrregularSequenceForms(t *testing.T) {
	t.Parallel()
	input := "alpha\r00:00:00,000 --> 00:00:01,000\rfirst\r\r2\n00:00:01,000 --> 00:00:02,000\nsecond\n\n2\r\n00:00:02,000 --> 00:00:03,000\r\nthird\r\n"
	result, err := Parse(input, Options{})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(result.Cues) != 3 {
		t.Fatalf("cue count = %d, want 3", len(result.Cues))
	}
	if got := result.Cues[0].FormatData.SubRip.SequenceLineRaw; got != "alpha" {
		t.Fatalf("non-integer sequence = %q", got)
	}
	if !hasDiagnostic(result.Diagnostics, "subrip_sequence_non_integer") || !hasDiagnostic(result.Diagnostics, "subrip_sequence_irregular") {
		t.Fatalf("diagnostics = %#v", result.Diagnostics)
	}
}

func TestParseRejectsMalformedTimingBlock(t *testing.T) {
	t.Parallel()
	input := "1\n00:00:01,000 --> 00:00:02,000 X1:1 X2:2 Y1:3\ntext\n"
	if _, err := Parse(input, Options{}); err == nil {
		t.Fatal("Parse() unexpectedly accepted an incomplete coordinate group")
	}
}

func TestParseRejectsInputWithoutValidCue(t *testing.T) {
	t.Parallel()
	for _, input := range []string{"", " \r\n\t", "not subtitles\n"} {
		if _, err := Parse(input, Options{}); err == nil {
			t.Errorf("Parse(%q) unexpectedly succeeded", input)
		}
	}
}

func TestGovernedSubRipFixtures(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..", "..", "testdata")
	source, err := os.ReadFile(filepath.Join(root, "fixtures", "subrip", "grammar-variants", "source", "variants.srt"))
	if err != nil {
		t.Fatal(err)
	}
	result, err := Parse(string(source), Options{DetectSpeakers: true})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	diagnosticBytes, err := os.ReadFile(filepath.Join(root, "fixtures", "subrip", "grammar-variants", "expected", "diagnostics.json"))
	if err != nil {
		t.Fatal(err)
	}
	var wantDiagnostics []Diagnostic
	if err := json.Unmarshal(diagnosticBytes, &wantDiagnostics); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(result.Diagnostics, wantDiagnostics) {
		t.Fatalf("diagnostics = %#v, want %#v", result.Diagnostics, wantDiagnostics)
	}
	rendered, err := Render(result.Cues)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	want, err := os.ReadFile(filepath.Join(root, "fixtures", "subrip", "grammar-variants", "expected", "bytes", "rendered.srt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(rendered) != string(want) {
		t.Fatalf("rendered fixture = %q, want %q", rendered, want)
	}
	malformed, err := os.ReadFile(filepath.Join(root, "malformed", "subrip", "reversed-time", "input", "reversed.srt"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(string(malformed), Options{}); err == nil {
		t.Fatal("reversed-time fixture unexpectedly succeeded")
	}
}

func hasDiagnostic(diagnostics []Diagnostic, code string) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return true
		}
	}
	return false
}
