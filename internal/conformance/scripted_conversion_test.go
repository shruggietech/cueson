package conformance_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/convert"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
	"github.com/shruggietech/cueson/internal/source"
	"github.com/shruggietech/cueson/internal/testutil"
)

type conversionEvidence struct {
	SourceFormat string                     `json:"source_format"`
	Targets      []conversionTargetEvidence `json:"targets"`
}

type conversionTargetEvidence struct {
	TargetFormat string                   `json:"target_format"`
	Cues         []conversionCueEvidence  `json:"cues"`
	Losses       []conversionLossEvidence `json:"losses"`
}

type conversionCueEvidence struct {
	StartMilliseconds int64    `json:"start_milliseconds"`
	EndMilliseconds   int64    `json:"end_milliseconds"`
	PlainText         string   `json:"plain_text"`
	Lines             []string `json:"lines"`
}

type conversionLossEvidence struct {
	Code string `json:"code"`
	Path string `json:"path"`
}

func TestScriptedTwelveDirectionConversionConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	directions := make(map[string]bool)
	for _, format := range []string{"subrip", "webvtt", "ass", "ssa"} {
		fixture := mustFixture(t, manifest, "scripted-conversion/"+format+"-baseline")
		var expected conversionEvidence
		decodeScriptedExpected(t, readArtifact(t, root, mustArtifact(t, fixture, "expected_conversion")), &expected)
		if expected.SourceFormat != format || len(expected.Targets) != 3 {
			t.Fatalf("incomplete baseline for %s", format)
		}
		document, _ := encodeScriptedFixture(t, root, fixture)
		before, err := json.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		original := readArtifact(t, root, mustArtifact(t, fixture, "source"))
		for _, target := range expected.Targets {
			key := format + "->" + target.TargetFormat
			if directions[key] || format == target.TargetFormat {
				t.Fatalf("duplicated or same-format row %s", key)
			}
			directions[key] = true
			t.Run(key, func(t *testing.T) {
				result, err := convert.Convert(context.Background(), document, target.TargetFormat, convert.Options{})
				if err != nil {
					t.Fatalf("baseline conversion: %v", err)
				}
				want := readArtifact(t, root, mustArtifact(t, fixture, "expected_"+target.TargetFormat+"_bytes"))
				if err = testutil.CompareBytes(fixture.ID, key, want, result.Bytes); err != nil {
					t.Fatal(err)
				}
				assertConversionCues(t, target.Cues, reparseConversion(t, result.Bytes, target.TargetFormat))
				assertConversionLosses(t, document, target.Losses, result.LossReport)
				again, err := convert.Convert(context.Background(), document, target.TargetFormat, convert.Options{})
				if err != nil || !bytes.Equal(result.Bytes, again.Bytes) || !reflect.DeepEqual(result.LossReport, again.LossReport) {
					t.Fatal("conversion bytes/report changed on repetition")
				}
				strict, strictErr := convert.Convert(context.Background(), document, target.TargetFormat, convert.Options{Strict: true})
				if len(target.Losses) == 0 {
					if strictErr != nil || !bytes.Equal(strict.Bytes, want) {
						t.Fatalf("loss-free strict conversion: %v", strictErr)
					}
				} else {
					var refusal *convert.StrictLossError
					if !errors.As(strictErr, &refusal) || len(strict.Bytes) != 0 || !reflect.DeepEqual(result.LossReport, strict.LossReport) {
						t.Fatal("strict did not return complete report without payload")
					}
				}
				after, err := json.Marshal(document)
				if err != nil || !bytes.Equal(before, after) {
					t.Fatal("conversion mutated immutable source document")
				}
				if err = source.ValidateIntegrity(context.Background(), document); err != nil {
					t.Fatal(err)
				}
				path := testutil.ArtifactFile(root, mustArtifact(t, fixture, "source"))
				if !bytes.Equal(readFile(t, path), original) {
					t.Fatal("conversion changed governed source bytes")
				}
				if err = testutil.CheckNoForbidden(key, "bytes/report", append(append([]byte(nil), result.Bytes...), mustMarshalConversion(t, result.LossReport)...), root, path); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
	if len(directions) != 12 {
		t.Fatalf("verified %d directions, want twelve", len(directions))
	}
}

func TestScriptedConversionStrictFatalPublicationConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, format := range []string{"subrip", "webvtt", "ass", "ssa"} {
		fixture := mustFixture(t, manifest, "scripted-conversion/"+format+"-baseline")
		original := readArtifact(t, root, mustArtifact(t, fixture, "source"))
		for _, target := range []string{"subrip", "webvtt", "ass", "ssa"} {
			if format == target {
				continue
			}
			t.Run(format+"->"+target, func(t *testing.T) {
				if format == "ass" || format == "ssa" {
					assertConversionPublicationRefusal(t, original, format, target, true)
				} else if target == "ass" || target == "ssa" {
					imprecise := bytes.ReplaceAll(bytes.ReplaceAll(original, []byte("01,000"), []byte("01,005")), []byte("01.000"), []byte("01.005"))
					assertConversionPublicationRefusal(t, imprecise, format, target, true)
				} else {
					// Existing native losses remain strict refusals for the old pair.
					oldID := "conversion/srt-lossy"
					if format == "webvtt" {
						oldID = "conversion/webvtt-lossy"
					}
					old := mustFixture(t, manifest, oldID)
					assertConversionPublicationRefusal(t, readArtifact(t, root, mustArtifact(t, old, "source")), format, target, true)
				}
				for _, strict := range []bool{false, true} {
					assertConversionPublicationRefusal(t, []byte("{not accepted native or Cue JSON\n"), format, target, strict)
				}
			})
		}
	}
	for _, format := range []string{"ass", "ssa"} {
		for _, caseName := range []string{"zero-dialogue", "empty-overlap-drawing", "malformed-comment", "malformed-attachment", "unused-invalid-style"} {
			fixture := mustFixture(t, manifest, "scripted/"+format+"-"+caseName)
			input := readArtifact(t, root, mustArtifact(t, fixture, "source"))
			for _, target := range []string{"subrip", "webvtt"} {
				for _, strict := range []bool{false, true} {
					assertConversionPublicationRefusal(t, input, format, target, strict)
				}
			}
		}
	}
}

func TestScriptedConversionCompleteLossConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, format := range []string{"ass", "ssa"} {
		fixture := mustFixture(t, manifest, "scripted-conversion/"+format+"-native-losses")
		var expected conversionEvidence
		decodeScriptedExpected(t, readArtifact(t, root, mustArtifact(t, fixture, "expected_conversion")), &expected)
		document, _ := encodeScriptedFixture(t, root, fixture)
		before := mustMarshalConversion(t, document)
		for _, target := range expected.Targets {
			result, err := convert.Convert(context.Background(), document, target.TargetFormat, convert.Options{})
			if err != nil {
				t.Fatal(err)
			}
			assertConversionLosses(t, document, target.Losses, result.LossReport)
			assertConversionCues(t, target.Cues, reparseConversion(t, result.Bytes, target.TargetFormat))
			if err = testutil.CompareBytes(fixture.ID, "complete_native_losses", readArtifact(t, root, mustArtifact(t, fixture, "expected_"+target.TargetFormat+"_bytes")), result.Bytes); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(before, mustMarshalConversion(t, document)) {
				t.Fatal("native loss analysis changed source owners")
			}
			assertConversionPublicationRefusal(t, readArtifact(t, root, mustArtifact(t, fixture, "source")), format, target.TargetFormat, true)
		}
		// Target dialects preserve valid empty scripts and inert native content.
		targetFormat := "ass"
		if format == "ass" {
			targetFormat = "ssa"
		}
		for _, caseName := range []string{"zero-dialogue", "native-order", "attachments", "declared-fields"} {
			seed := mustFixture(t, manifest, "scripted/"+format+"-"+caseName)
			original, _ := encodeScriptedFixture(t, root, seed)
			result, err := convert.Convert(context.Background(), original, targetFormat, convert.Options{})
			if err != nil {
				t.Fatalf("retained variant %s/%s: %v", format, caseName, err)
			}
			parsed, err := scripted.Parse(context.Background(), result.Bytes, targetFormat)
			if err != nil {
				t.Fatal(err)
			}
			originalNative := original.FormatData.ASS
			if format == "ssa" {
				originalNative = original.FormatData.SSA
			}
			if len(parsed.Cues) != len(original.Cues) || len(parsed.Native.Events) != len(originalNative.Events) || len(parsed.Native.Attachments) != len(originalNative.Attachments) || len(parsed.Native.Sections) != len(originalNative.Sections) {
				t.Fatal("variant silently dropped inert native owners")
			}
			for index, attachment := range originalNative.Attachments {
				if parsed.Native.Attachments[index].Name != attachment.Name {
					t.Fatal("variant changed embedded attachment identity")
				}
			}
			for index, style := range originalNative.Styles {
				if !reflect.DeepEqual(conversionUnknownFields(style.Fields, format, "style"), conversionUnknownFields(parsed.Native.Styles[index].Fields, targetFormat, "style")) {
					t.Fatal("variant dropped declared unknown style fields")
				}
			}
			for index, event := range originalNative.Events {
				if !reflect.DeepEqual(conversionUnknownFields(event.Fields, format, "event"), conversionUnknownFields(parsed.Native.Events[index].Fields, targetFormat, "event")) {
					t.Fatal("variant dropped duplicate declared unknown event occurrences")
				}
			}
		}
	}
}

func conversionUnknownFields(fields []model.ScriptedField, format, kind string) []scriptedFieldExpected {
	values := []scriptedFieldExpected{}
	for _, field := range fields {
		if _, known := model.ScriptedFieldIdentity(field.FieldName, format, kind); !known {
			values = append(values, scriptedFieldExpected{strings.TrimSpace(field.FieldName), field.RawValue})
		}
	}
	return values
}

func TestScriptedConversionSourceIntegrityPrivacyAndHistoricalConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, format := range []string{"subrip", "webvtt", "ass", "ssa"} {
		fixture := mustFixture(t, manifest, "scripted-conversion/"+format+"-baseline")
		document, encoded := encodeScriptedFixture(t, root, fixture)
		for _, target := range []string{"subrip", "webvtt", "ass", "ssa"} {
			if format == target {
				continue
			}
			result, err := convert.Convert(context.Background(), document, target, convert.Options{})
			if err != nil {
				t.Fatal(err)
			}
			directory := t.TempDir()
			jsonPath := filepath.Join(directory, "portable.cueson.json")
			if err = os.WriteFile(jsonPath, encoded, 0o600); err != nil {
				t.Fatal(err)
			}
			var stdout, stderr bytes.Buffer
			status := cli.Run(context.Background(), []string{"convert", jsonPath, "--to", target}, strings.NewReader(""), &stdout, &stderr)
			if status != cli.ExitSuccess || !bytes.Equal(stdout.Bytes(), result.Bytes) {
				t.Fatalf("native/Cue JSON parity %s->%s: %d %s", format, target, status, stderr.Bytes())
			}
			if err = testutil.CheckNoForbidden(format+"->"+target, "parity", append(stdout.Bytes(), stderr.Bytes()...), directory, jsonPath); err != nil {
				t.Fatal(err)
			}
			corrupt := document
			corrupt.Source.Assets = append([]model.SourceAsset(nil), document.Source.Assets...)
			corrupt.Source.Assets[0].Hashes.SHA256 = strings.Repeat("0", 64)
			bad, err := convert.Convert(context.Background(), corrupt, target, convert.Options{})
			if err == nil || len(bad.Bytes) != 0 {
				t.Fatal("corrupt source integrity published conversion")
			}
			assertConversionPublicationRefusal(t, mustMarshalConversion(t, corrupt), "json", target, false)
		}
	}
	for _, format := range []string{"subrip", "webvtt"} {
		data := readFile(t, filepath.Join(root, "..", "internal", "schema", "testdata", "historical-v1.0.0-"+format+".json"))
		document, err := schema.Decode(data)
		if err != nil {
			t.Fatal(err)
		}
		before := mustMarshalConversion(t, document)
		for _, target := range []string{"ass", "ssa"} {
			result, err := convert.Convert(context.Background(), document, target, convert.Options{})
			if err != nil {
				t.Fatalf("historical %s->%s: %v", format, target, err)
			}
			if len(reparseConversion(t, result.Bytes, target)) != len(document.Cues) {
				t.Fatal("historical cue count changed")
			}
			if !bytes.Equal(before, mustMarshalConversion(t, document)) {
				t.Fatal("historical source identity/envelope changed")
			}
		}
	}
}

func TestScriptedConversionPrecisionAndDefaultsConformance(t *testing.T) {
	t.Parallel()
	root := fixtureRoot(t)
	manifest := verifiedManifest(t, root)
	for _, format := range []string{"subrip", "webvtt"} {
		fixture := mustFixture(t, manifest, "scripted-conversion/"+format+"-baseline")
		baseline := readArtifact(t, root, mustArtifact(t, fixture, "source"))
		for _, target := range []string{"ass", "ssa"} {
			for _, point := range []struct {
				lexeme string
				want   int64
				losses int
			}{{"000", 1000, 0}, {"004", 1000, 1}, {"005", 1010, 1}, {"006", 1010, 1}, {"010", 1010, 0}} {
				t.Run(format+"->"+target+"/"+point.lexeme, func(t *testing.T) {
					input := bytes.ReplaceAll(bytes.ReplaceAll(baseline, []byte("01,000"), []byte("01,"+point.lexeme)), []byte("01.000"), []byte("01."+point.lexeme))
					directory := t.TempDir()
					path := filepath.Join(directory, "precision."+format)
					if err := os.WriteFile(path, input, 0o600); err != nil {
						t.Fatal(err)
					}
					var encoded, diagnostics bytes.Buffer
					if status := cli.Run(context.Background(), []string{"encode", path, "--stdout"}, strings.NewReader(""), &encoded, &diagnostics); status != cli.ExitSuccess {
						t.Fatalf("precision encode %d: %s", status, diagnostics.Bytes())
					}
					document, err := schema.Decode(encoded.Bytes())
					if err != nil {
						t.Fatal(err)
					}
					result, err := convert.Convert(context.Background(), document, target, convert.Options{})
					if err != nil {
						t.Fatal(err)
					}
					cues := reparseConversion(t, result.Bytes, target)
					if len(cues) != 1 || cues[0].Timing.StartMilliseconds != point.want || cues[0].Timing.EndMilliseconds != 2500 {
						t.Fatalf("precision reparsed %v", cues)
					}
					if len(result.LossReport.Losses) != point.losses {
						t.Fatalf("endpoint losses=%v", result.LossReport.Losses)
					}
					if point.losses != 0 && (result.LossReport.Losses[0].Code != "conversion_scripted_centisecond_quantized" || result.LossReport.Losses[0].Path != "/cues/0/timing/start_milliseconds") {
						t.Fatal("quantization lacks atomic endpoint reference")
					}
					if point.losses != 0 {
						assertConversionPublicationRefusal(t, input, format, target, true)
					}
				})
			}
			for _, payload := range []string{"literal {braces}", `literal \N`, `literal \n`, `literal \h`} {
				input := bytes.ReplaceAll(baseline, []byte("Hello\nworld"), []byte(payload))
				for _, strict := range []bool{false, true} {
					assertConversionPublicationRefusal(t, input, format, target, strict)
				}
			}
			separator := ","
			if format == "webvtt" {
				separator = "."
			}
			input := []byte("1\n00:00:01" + separator + "001 --> 00:00:01" + separator + "004\ncollapse\n")
			if format == "webvtt" {
				input = []byte("WEBVTT\n\n00:00:01.001 --> 00:00:01.004\ncollapse\n")
			}
			for _, strict := range []bool{false, true} {
				assertConversionPublicationRefusal(t, input, format, target, strict)
			}
		}
	}
}

func reparseConversion(t *testing.T, data []byte, format string) []model.Cue {
	t.Helper()
	switch format {
	case "subrip":
		parsed, err := subrip.Parse(string(data), subrip.Options{})
		if err != nil {
			t.Fatal(err)
		}
		return parsed.Cues
	case "webvtt":
		parsed, err := webvtt.Parse(string(data))
		if err != nil {
			t.Fatal(err)
		}
		return parsed.Cues
	case "ass", "ssa":
		parsed, err := scripted.Parse(context.Background(), data, format)
		if err != nil {
			t.Fatal(err)
		}
		return parsed.Cues
	default:
		t.Fatalf("unknown target %s", format)
		return nil
	}
}

func assertConversionCues(t *testing.T, expected []conversionCueEvidence, cues []model.Cue) {
	t.Helper()
	actual := make([]conversionCueEvidence, 0, len(cues))
	for _, cue := range cues {
		actual = append(actual, conversionCueEvidence{cue.Timing.StartMilliseconds, cue.Timing.EndMilliseconds, cue.Payload.PlainText, cue.Payload.Lines})
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("reparsed cues = %#v, want %#v", actual, expected)
	}
}

func assertConversionLosses(t *testing.T, document model.Document, expected []conversionLossEvidence, report convert.Report) {
	t.Helper()
	if err := report.Validate(document); err != nil {
		t.Fatal(err)
	}
	actual := make([]conversionLossEvidence, 0, len(report.Losses))
	for _, loss := range report.Losses {
		actual = append(actual, conversionLossEvidence{loss.Code, loss.Path})
		if !convert.IsKnownLossCode(loss.Code) || loss.Message == "" || loss.SourceFormat != document.Format {
			t.Fatal("loss escaped governed vocabulary")
		}
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("complete ordered losses = %#v, want %#v", actual, expected)
	}
}

func mustMarshalConversion(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertConversionPublicationRefusal(t *testing.T, input []byte, format, target string, strict bool) {
	t.Helper()
	directory := t.TempDir()
	inputPath := filepath.Join(directory, "captions."+format)
	if err := os.WriteFile(inputPath, input, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, route := range []string{"default_stdout", "explicit_stdout", "new_destination", "forced_destination"} {
		args := []string{"convert", inputPath, "--to", target}
		if strict {
			args = append(args, "--strict")
		}
		outputPath := filepath.Join(directory, "target."+target)
		preserved := []byte("existing destination must remain byte-identical\n")
		switch route {
		case "explicit_stdout":
			args = append(args, "--output", "-")
		case "new_destination":
			args = append(args, "--output", outputPath)
		case "forced_destination":
			if err := os.WriteFile(outputPath, preserved, 0o600); err != nil {
				t.Fatal(err)
			}
			args = append(args, "--output", outputPath, "--force")
		}
		before, err := os.ReadDir(directory)
		if err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		status := cli.Run(context.Background(), args, strings.NewReader(""), &stdout, &stderr)
		if status != cli.ExitRuntimeFailure || stdout.Len() != 0 || stderr.Len() == 0 {
			t.Fatalf("%s->%s strict=%t %s refusal = %d, %q, %q", format, target, strict, route, status, stdout.Bytes(), stderr.Bytes())
		}
		if route == "forced_destination" && !bytes.Equal(readFile(t, outputPath), preserved) {
			t.Fatal("refusal changed forced destination")
		}
		if route == "new_destination" {
			if _, err = os.Lstat(outputPath); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("refusal created destination: %v", err)
			}
		}
		after, err := os.ReadDir(directory)
		if err != nil || len(before) != len(after) {
			t.Fatal("refusal left partial staging artifact")
		}
		if !bytes.Equal(readFile(t, inputPath), input) {
			t.Fatal("refusal changed source bytes")
		}
		if err = testutil.CheckNoForbidden("conversion/refusal", "diagnostics", stderr.Bytes(), directory, inputPath, outputPath); err != nil {
			t.Fatal(err)
		}
	}
}
