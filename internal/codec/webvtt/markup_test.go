package webvtt

import "testing"

func TestPayloadScannerDerivesMarkupEntitiesSpeakerAndInlineTiming(t *testing.T) {
	raw := "<v Roger><c.green>Hello &amp; <b>world</b></c>\n<00:02.000><ruby>漢<rt>かん</rt></ruby></v>"
	plain, speakers, tokens, diagnostics := scanPayload(raw, 1000, 3000, 0, "cue-000000")
	if plain != "Hello & world\n漢かん" {
		t.Fatalf("plain = %q", plain)
	}
	if len(speakers) != 1 || speakers[0].Name != "Roger" || speakers[0].Origin != "native" {
		t.Fatalf("speakers = %#v", speakers)
	}
	if len(tokens) != 2 || tokens[0].StartMilliseconds != 1000 || tokens[0].EndMilliseconds != 2000 || tokens[1].StartMilliseconds != 2000 || tokens[1].EndMilliseconds != 3000 {
		t.Fatalf("tokens = %#v", tokens)
	}
	if len(diagnostics) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
}

func TestPayloadScannerPreservesUnknownAndInvalidSyntax(t *testing.T) {
	raw := "<blink>x</blink> &bogus; <00:00.500>"
	plain, _, tokens, diagnostics := scanPayload(raw, 1000, 2000, 0, "cue-000000")
	if plain != raw || len(tokens) != 0 {
		t.Fatalf("plain=%q tokens=%#v", plain, tokens)
	}
	if len(diagnostics) != 4 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	if diagnostics[3].Code != "webvtt_inline_timestamp_invalid" {
		t.Fatalf("last diagnostic = %#v", diagnostics[3])
	}
}

func TestPayloadScannerPreservesUnterminatedEntity(t *testing.T) {
	plain, _, _, diagnostics := scanPayload("A &amp B", 0, 1000, 0, "cue-000000")
	if plain != "A &amp B" || len(diagnostics) != 1 || diagnostics[0].Code != "webvtt_entity_invalid" {
		t.Fatalf("plain=%q diagnostics=%#v", plain, diagnostics)
	}
}

func TestPayloadScannerAcceptsSoleVoiceSpanWithOmittedEndTag(t *testing.T) {
	plain, speakers, _, diagnostics := scanPayload("<v Roger>Voice text", 0, 1000, 0, "cue-000000")
	if plain != "Voice text" || len(speakers) != 1 || speakers[0].Name != "Roger" || len(diagnostics) != 0 {
		t.Fatalf("plain=%q speakers=%#v diagnostics=%#v", plain, speakers, diagnostics)
	}
}

func TestPayloadScannerUsesWebVTTCharacterReferenceSubset(t *testing.T) {
	plain, _, _, diagnostics := scanPayload("&lt;&gt;&lrm;&rlm;&nbsp;&#65;&copy;", 0, 1000, 0, "cue-000000")
	if plain != "<>\u200e\u200f\u00a0A&copy;" || len(diagnostics) != 1 || diagnostics[0].Code != "webvtt_entity_invalid" {
		t.Fatalf("plain=%q diagnostics=%#v", plain, diagnostics)
	}
}
