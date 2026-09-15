package convert_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/shruggietech/cueson/internal/cli"
	"github.com/shruggietech/cueson/internal/codec/scripted"
	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/convert"
	"github.com/shruggietech/cueson/internal/model"
	"github.com/shruggietech/cueson/internal/schema"
)

func FuzzScriptedConversionCycle(f *testing.F) {
	f.Add([]byte("1\n00:00:01,000 --> 00:00:02,005\n<b>Hello</b>\n"), uint8(0))
	f.Add([]byte("WEBVTT\n\n00:01.000 --> 00:02.000\n<i>Hello &amp; world</i>\n"), uint8(1))
	for index, format := range []string{"ass", "ssa"} {
		fixture, err := os.ReadFile("../schema/testdata/scripted-" + format + ".cueson.json")
		if err != nil {
			f.Fatal(err)
		}
		document, err := schema.Decode(fixture)
		if err != nil {
			f.Fatal(err)
		}
		raw, err := base64.StdEncoding.DecodeString(document.Source.Assets[0].DataBase64)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(raw, uint8(index+2))
	}
	f.Fuzz(func(t *testing.T, raw []byte, selection uint8) {
		if len(raw) > 64<<10 {
			t.Skip()
		}
		formats := []string{"subrip", "webvtt", "ass", "ssa"}
		format := formats[int(selection)%len(formats)]
		document, accepted := encodeConversionFuzzSource(t, raw, format)
		if !accepted {
			return
		}
		before, err := json.Marshal(document)
		if err != nil {
			t.Fatal(err)
		}
		for _, target := range formats {
			if target == format {
				continue
			}
			permissive, conversionErr := convert.Convert(context.Background(), document, target, convert.Options{})
			if conversionErr != nil {
				if len(permissive.Bytes) != 0 {
					t.Fatal("fatal conversion returned partial bytes")
				}
				strict, strictErr := convert.Convert(context.Background(), document, target, convert.Options{Strict: true})
				if strictErr == nil || len(strict.Bytes) != 0 {
					t.Fatal("fatal boundary was accepted in strict mode")
				}
				continue
			}
			if err := permissive.LossReport.Validate(document); err != nil {
				t.Fatal(err)
			}
			parseConversionFuzzTarget(t, permissive.Bytes, target)
			repeated, err := convert.Convert(context.Background(), document, target, convert.Options{})
			if err != nil || !bytes.Equal(repeated.Bytes, permissive.Bytes) || !reflect.DeepEqual(repeated.LossReport, permissive.LossReport) {
				t.Fatalf("nondeterministic conversion to %s: %v", target, err)
			}
			strict, strictErr := convert.Convert(context.Background(), document, target, convert.Options{Strict: true})
			if permissive.LossReport.HasLosses() {
				var lossErr *convert.StrictLossError
				if !errors.As(strictErr, &lossErr) || len(strict.Bytes) != 0 || !reflect.DeepEqual(lossErr.Report, permissive.LossReport) {
					t.Fatalf("incomplete strict policy to %s: %v", target, strictErr)
				}
			} else if strictErr != nil || !bytes.Equal(strict.Bytes, permissive.Bytes) {
				t.Fatalf("loss-free strict conversion to %s: %v", target, strictErr)
			}
		}
		after, err := json.Marshal(document)
		if err != nil || !bytes.Equal(before, after) {
			t.Fatal("conversion mutated original source model")
		}
	})
}

func encodeConversionFuzzSource(t *testing.T, raw []byte, format string) (model.Document, bool) {
	t.Helper()
	input := filepath.Join(t.TempDir(), "input.bin")
	if err := os.WriteFile(input, raw, 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	status := cli.Run(context.Background(), []string{"encode", input, "--format", format, "--stdout", "--no-speaker-detection"}, nil, &stdout, &stderr)
	if status != cli.ExitSuccess {
		return model.Document{}, false
	}
	document, err := schema.Decode(stdout.Bytes())
	if err != nil {
		t.Fatalf("accepted input did not decode: %v", err)
	}
	return document, true
}

func parseConversionFuzzTarget(t *testing.T, raw []byte, format string) {
	t.Helper()
	var err error
	switch format {
	case "subrip":
		_, err = subrip.Parse(string(raw), subrip.Options{})
	case "webvtt":
		_, err = webvtt.Parse(string(raw))
	default:
		_, err = scripted.Parse(context.Background(), raw, format)
	}
	if err != nil {
		t.Fatalf("target %s rejected complete conversion: %v", format, err)
	}
}

func TestConversionRejectsNilContextCancellationAndCorruptSource(t *testing.T) {
	for _, format := range []string{"subrip", "webvtt"} {
		raw := []byte("1\n00:00:01,000 --> 00:00:02,000\nx\n")
		if format == "webvtt" {
			raw = []byte("WEBVTT\n\n00:01.000 --> 00:02.000\nx\n")
		}
		document, accepted := encodeConversionFuzzSource(t, raw, format)
		if !accepted {
			t.Fatal("baseline rejected")
		}
		for _, target := range []string{"ass", "ssa"} {
			//lint:ignore SA1012 Verify the explicitly supported nil-context refusal boundary.
			result, err := convert.Convert(nil, document, target, convert.Options{})
			if err == nil || len(result.Bytes) != 0 {
				t.Fatal("nil context did not refuse safely")
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			result, err = convert.Convert(ctx, document, target, convert.Options{})
			if !errors.Is(err, context.Canceled) || len(result.Bytes) != 0 {
				t.Fatal("cancelled conversion did not refuse safely")
			}
			corrupt := document
			corrupt.Source.Assets = append([]model.SourceAsset(nil), document.Source.Assets...)
			corrupt.Source.Assets[0].Hashes.SHA256 = "0000000000000000000000000000000000000000000000000000000000000000"
			result, err = convert.Convert(context.Background(), corrupt, target, convert.Options{})
			if err == nil || len(result.Bytes) != 0 {
				t.Fatal("corrupt source did not refuse safely")
			}
		}
	}
}
