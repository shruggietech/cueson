package cli

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/schema"
)

func TestScriptedNativeEndToEndAndSourceSeparation(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		t.Run(format, func(t *testing.T) {
			fixture, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
			if err != nil {
				t.Fatal(err)
			}
			seed, err := schema.Decode(fixture)
			if err != nil {
				t.Fatal(err)
			}
			sourceBytes, err := base64.StdEncoding.DecodeString(seed.Source.Assets[0].DataBase64)
			if err != nil {
				t.Fatal(err)
			}
			for _, extension := range []string{".srt", ".vtt", ".json", "." + format} {
				t.Run(extension, func(t *testing.T) {
					directory := t.TempDir()
					input := filepath.Join(directory, "misleading"+extension)
					if err := os.WriteFile(input, sourceBytes, 0600); err != nil {
						t.Fatal(err)
					}
					status, stdout, stderr := runForTest(context.Background(), []string{"encode", input, "--stdout", "--no-speaker-detection"})
					if status != ExitSuccess {
						t.Fatalf("encode = %d %s", status, stderr)
					}
					doc, err := schema.Decode([]byte(stdout))
					if err != nil {
						t.Fatal(err)
					}
					if doc.Format != format || doc.FormatSupport.Status != "stable" || !doc.FormatSupport.IngestSupported || !doc.FormatSupport.RenderSupported || doc.Source.Assets[0].Encoding == nil {
						t.Fatalf("native contract = %#v", doc.FormatSupport)
					}
					if len(doc.Cues) != len(seed.Cues) || !reflect.DeepEqual(doc.Cues[0].Payload, seed.Cues[0].Payload) || !reflect.DeepEqual(doc.Cues[0].Speakers, seed.Cues[0].Speakers) {
						t.Fatal("native semantics/actor provenance changed")
					}
					encoded := filepath.Join(directory, "model.cueson.json")
					if err := os.WriteFile(encoded, []byte(stdout), 0600); err != nil {
						t.Fatal(err)
					}
					for _, path := range []string{input, encoded} {
						for _, command := range []string{"validate", "inspect"} {
							args := []string{command, path}
							if command == "inspect" {
								args = append(args, "--json")
							}
							status, out, errOut := runForTest(context.Background(), args)
							if status != ExitSuccess {
								t.Fatalf("%v = %d %s", args, status, errOut)
							}
							if command == "inspect" && (strings.Contains(out, directory) || strings.Contains(out, filepath.ToSlash(directory))) {
								t.Fatal("inspection path leak")
							}
						}
					}
					restored := filepath.Join(directory, "restored."+format)
					status, _, stderr = runForTest(context.Background(), []string{"restore", encoded, "--no-metadata", "--output", restored})
					if status != ExitSuccess {
						t.Fatalf("restore = %d %s", status, stderr)
					}
					got, err := os.ReadFile(restored)
					if err != nil || !bytes.Equal(got, sourceBytes) {
						t.Fatalf("exact restore mismatch: %v", err)
					}
					status, native, stderr := runForTest(context.Background(), []string{"render", encoded, "--to", format})
					if status != ExitSuccess {
						t.Fatalf("render = %d %s", status, stderr)
					}
					strictStatus, strictNative, strictDiagnostics := runForTest(context.Background(), []string{"render", encoded, "--to", format, "--strict"})
					if strictStatus != ExitSuccess || strictNative != native {
						t.Fatalf("strict render of conforming %s source under %s = %d %s", format, extension, strictStatus, strictDiagnostics)
					}
					parsed, err := scripted.Parse(context.Background(), []byte(native), format)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(parsed.Cues[0].Payload, doc.Cues[0].Payload) || parsed.Cues[0].Timing != doc.Cues[0].Timing {
						t.Fatal("render/reparse projection changed")
					}
					for index, cue := range parsed.Cues {
						if !reflect.DeepEqual(cue.Speakers, doc.Cues[index].Speakers) {
							t.Fatal("actor provenance changed on render")
						}
					}
					status, again, stderr := runForTest(context.Background(), []string{"render", encoded, "--to", format})
					if status != ExitSuccess || again != native {
						t.Fatalf("nondeterministic render %d %s", status, stderr)
					}
				})
			}
		})
	}
}

func TestScriptedNativeInvalidInputsNeverPublish(t *testing.T) {
	b, err := os.ReadFile("../schema/testdata/scripted-ass.cueson.json")
	if err != nil {
		t.Fatal(err)
	}
	seed, err := schema.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(seed.Source.Assets[0].DataBase64)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		source []byte
		format string
	}{
		{"mixed dialect", []byte(strings.Replace(string(raw), "v4.00+", "v4.00", 1)), "auto"},
		{"unsafe relative resource", []byte(strings.Replace(string(raw), "ScriptType:", "Video File: assets/private.mkv\nScriptType:", 1)), "auto"},
		{"explicit mismatched dialect", raw, "ssa"},
		{"explicit other codec escape", raw, "srt"},
	} {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "source.srt")
			output := filepath.Join(dir, "selected.json")
			if err := os.WriteFile(input, test.source, 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
				t.Fatal(err)
			}
			args := []string{"encode", input, "--format", test.format, "--output", output, "--force"}
			status, stdout, stderr := runForTest(context.Background(), args)
			if status != ExitRuntimeFailure || stdout != "" || strings.Contains(stderr, "assets/private.mkv") {
				t.Fatalf("invalid source = %d %q %q", status, stdout, stderr)
			}
			after, err := os.ReadFile(output)
			if err != nil || string(after) != "preserve" {
				t.Fatal("fatal encode published output")
			}
			status, stdout, _ = runForTest(context.Background(), []string{"encode", input, "--format", test.format, "--stdout"})
			if status != ExitRuntimeFailure || stdout != "" {
				t.Fatal("fatal encode published stdout")
			}
			after, err = os.ReadFile(input)
			if err != nil || !bytes.Equal(after, test.source) {
				t.Fatal("source bytes modified")
			}
		})
	}
}

func TestScriptedSelectedEncodingRequiresMatchingProfile(t *testing.T) {
	for _, format := range []string{"ass", "ssa"} {
		fixture, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			t.Fatal(err)
		}
		doc, err := schema.Decode(fixture)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := base64.StdEncoding.DecodeString(doc.Source.Assets[0].DataBase64)
		if err != nil {
			t.Fatal(err)
		}
		raw = bytes.TrimPrefix(raw, []byte{0xef, 0xbb, 0xbf})
		for _, test := range []struct {
			encoding string
			bom      bool
			success  bool
		}{
			{"utf-8-bom", false, false}, {"utf8-bom", true, true},
			{"utf-8", false, true}, {"utf-8", true, true}, {"windows-1252", false, false},
		} {
			input := filepath.Join(t.TempDir(), "input."+format)
			payload := raw
			if test.bom {
				payload = append([]byte{0xef, 0xbb, 0xbf}, raw...)
			}
			if err := os.WriteFile(input, payload, 0600); err != nil {
				t.Fatal(err)
			}
			status, stdout, stderr := runForTest(context.Background(), []string{"encode", input, "--stdout", "--encoding", test.encoding})
			if test.success && status != ExitSuccess {
				t.Fatalf("%s %s BOM=%t: %d %s", format, test.encoding, test.bom, status, stderr)
			}
			if !test.success && (status != ExitRuntimeFailure || stdout != "") {
				t.Fatalf("incompatible profile published: %d %q", status, stdout)
			}
		}
	}
}

func TestScriptedRenderPrecisionRefusesBeforeReplacement(t *testing.T) {
	payload, err := os.ReadFile("../schema/testdata/scripted-ass.cueson.json")
	if err != nil {
		t.Fatal(err)
	}
	doc, err := schema.Decode(payload)
	if err != nil {
		t.Fatal(err)
	}
	doc.Cues[0].Timing.StartMilliseconds++
	doc.Cues[0].Timing.DurationMilliseconds--
	*doc.Document.MediaStartMilliseconds++
	*doc.Document.MediaSpanMilliseconds--
	*doc.Stats.MediaSpanMilliseconds--
	payload, err = json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	input := filepath.Join(directory, "edited.json")
	output := filepath.Join(directory, "selected.ass")
	if err := os.WriteFile(input, payload, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, []byte("preserve"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, strict := range []bool{false, true} {
		args := []string{"render", input, "--to", "ass", "--output", output, "--force"}
		if strict {
			args = append(args, "--strict")
		}
		status, stdout, stderr := runForTest(context.Background(), args)
		if status != ExitRuntimeFailure || stdout != "" || !strings.Contains(stderr, "unrepresentable_centiseconds") {
			t.Fatalf("precision = %d %q %q", status, stdout, stderr)
		}
		after, err := os.ReadFile(output)
		if err != nil || string(after) != "preserve" {
			t.Fatal("precision failure replaced destination")
		}
	}
}
