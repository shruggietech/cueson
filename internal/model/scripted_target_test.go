package model

import (
	"strings"
	"testing"
)

func captureFreeScriptedExample(t *testing.T, format string) Document {
	t.Helper()
	doc := scriptedExample(t, format)
	native := doc.FormatData.ASS
	if format == "ssa" {
		native = doc.FormatData.SSA
	}
	for index := range native.Sections {
		native.Sections[index].RawHeader = nil
	}
	for index := range native.Records {
		native.Records[index].RawLine = nil
	}
	for index := range doc.Cues {
		data := doc.Cues[index].FormatData.ASS
		if format == "ssa" {
			data = doc.Cues[index].FormatData.SSA
		}
		data.StartTimestampRaw, data.EndTimestampRaw = nil, nil
	}
	return doc
}

func TestScriptedTargetKeepsPublicSourceGrammarGate(t *testing.T) {
	source := scriptedExample(t, "ssa")
	target := captureFreeScriptedExample(t, "ass")
	target.Source = source.Source
	if err := target.Validate(); err == nil {
		t.Fatal("public validation accepted mismatched original dialect")
	}
	if err := ValidateScriptedTarget(source, target); err != nil {
		t.Fatal(err)
	}
	target.Cues[0].Payload.PlainText = "invented"
	if err := ValidateScriptedTarget(source, target); err == nil {
		t.Fatal("constructed target skipped projection validation")
	}
}

func TestScriptedTargetRejectsChangedEnvelopeAndCaptures(t *testing.T) {
	source := scriptedExample(t, "ssa")
	target := captureFreeScriptedExample(t, "ass")
	target.Source = source.Source
	target.Source.Assets = append([]SourceAsset{}, source.Source.Assets...)
	target.Source.Assets[0].FileName = "different.ssa"
	if err := ValidateScriptedTarget(source, target); err == nil || !strings.Contains(err.Error(), "envelope_changed") {
		t.Fatalf("changed envelope: %v", err)
	}
	target.Source = source.Source
	wrong := "not an original capture"
	target.FormatData.ASS.Records[0].RawLine = &wrong
	if err := ValidateScriptedTarget(source, target); err == nil || !strings.Contains(err.Error(), "inconsistent_source_capture") {
		t.Fatalf("changed capture: %v", err)
	}
}
