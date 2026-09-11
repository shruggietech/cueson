package convert

import (
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
)

func TestTranslateSubRipPayloadPreservesSharedMarkupAndEscapesWebVTTSyntax(t *testing.T) {
	raw := "<B>A & B</B>\n<i><u>shared</u></i>\nliteral <v Bob> and 00:00:01,000 --> 00:00:02,000"
	got, err := translateSubRipPayload(raw)
	if err != nil {
		t.Fatal(err)
	}
	want := "<b>A &amp; B</b>\n<i><u>shared</u></i>\nliteral &lt;v Bob&gt; and 00:00:01,000 --&gt; 00:00:02,000"
	if got.Text != want || len(got.Issues) != 0 {
		t.Fatalf("translation = %#v, want text %q without issues", got, want)
	}
	parsed, err := webvtt.Parse("WEBVTT\n\n00:00.000 --> 00:03.000\n" + got.Text + "\n")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Cues[0].Payload.PlainText != subrip.PlainText(raw) {
		t.Fatalf("plain text = %q, want %q", parsed.Cues[0].Payload.PlainText, subrip.PlainText(raw))
	}
}

func TestTranslateSubRipPayloadAccountsForFontNULAndEmptyLines(t *testing.T) {
	got, err := translateSubRipPayload("<font color=red>x</font>\n\nA\x00B")
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "x\n<i></i>\nA\uFFFDB" {
		t.Fatalf("text = %q", got.Text)
	}
	if codes := issueCodes(got.Issues); strings.Join(codes, ",") != strings.Join([]string{LossCodeSubRipFontDegraded, LossCodeNULDegraded, LossCodePayloadLineDegraded}, ",") {
		t.Fatalf("issues = %#v", got.Issues)
	}
}

func TestTranslateSubRipPayloadKeepsMalformedOrUnsupportedTagsLiteral(t *testing.T) {
	raw := "<b class=x>literal</b> <i><b>crossed</i></b> <> </>"
	got, err := translateSubRipPayload(raw)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := webvtt.Parse("WEBVTT\n\n00:00.000 --> 00:01.000\n" + got.Text + "\n")
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Cues[0].Payload.PlainText != raw || len(got.Issues) != 0 {
		t.Fatalf("translation=%#v plain=%q", got, parsed.Cues[0].Payload.PlainText)
	}
}

func TestTranslateWebVTTPayloadPreservesSharedMarkupAndLinearizesNativeMarkup(t *testing.T) {
	raw := "<v Alice><b>A &amp; B</b> <c.red>C</c> <lang en>hola</lang> <ruby>漢<rt>かん</rt></ruby>\n<00:02.000>again"
	got, err := translateWebVTTPayload(raw, 1000, 3000)
	if err != nil {
		t.Fatal(err)
	}
	want := "<b>A & B</b> C hola 漢かん\nagain"
	if got.Text != want {
		t.Fatalf("text = %q, want %q", got.Text, want)
	}
	wantCodes := strings.Join([]string{LossCodeWebVTTVoiceDegraded, LossCodeWebVTTMarkupDegraded, LossCodeWebVTTMarkupDegraded, LossCodeWebVTTMarkupDegraded, LossCodeWebVTTMarkupDegraded, LossCodeWebVTTInlineTimingOmitted}, ",")
	if codes := strings.Join(issueCodes(got.Issues), ","); codes != wantCodes {
		t.Fatalf("issue codes = %q, want %q", codes, wantCodes)
	}
	if subrip.PlainText(got.Text) != "A & B C hola 漢かん\nagain" {
		t.Fatalf("SubRip plain text = %q", subrip.PlainText(got.Text))
	}
}

func TestTranslateWebVTTPayloadUsesSafePlaceholderAndRejectsReinterpretation(t *testing.T) {
	got, err := translateWebVTTPayload("<c></c>\ntext", 0, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if got.Text != "<u></u>\ntext" || len(got.Issues) != 2 {
		t.Fatalf("translation = %#v", got)
	}
	entity, err := translateWebVTTPayload("&lt;b&gt;literal&lt;/b&gt;", 0, 1000)
	if err != nil {
		t.Fatal(err)
	}
	wantEntityCodes := strings.TrimSuffix(strings.Repeat(LossCodeWebVTTEntityAmbiguous+",", 4), ",")
	if entity.Text != "&lt;b&gt;literal&lt;/b&gt;" || strings.Join(issueCodes(entity.Issues), ",") != wantEntityCodes {
		t.Fatalf("ambiguous entity translation = %#v", entity)
	}
}

func issueCodes(issues []payloadIssue) []string {
	codes := make([]string, len(issues))
	for index := range issues {
		codes[index] = issues[index].Code
	}
	return codes
}
