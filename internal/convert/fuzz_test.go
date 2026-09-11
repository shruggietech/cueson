package convert

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/shruggietech/cueson/internal/codec/subrip"
	"github.com/shruggietech/cueson/internal/codec/webvtt"
	"github.com/shruggietech/cueson/internal/model"
)

const maxFuzzInputBytes = 64 << 10

func FuzzSubRipConversionCycle(f *testing.F) {
	f.Add("1\n00:00:01,000 --> 00:00:02,000\n<b>x & y</b>\n")
	f.Add("7\n00:00:01,000 --> 00:00:02,000 X1:1 X2:2 Y1:3 Y2:4\n<font>x</font>\n")
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > maxFuzzInputBytes {
			t.Skip()
		}
		if _, err := subrip.Parse(input, subrip.Options{}); err != nil {
			return
		}
		document := conversionTestDocument(t, input, "subrip")
		assertFuzzConversionPolicy(t, document, "webvtt")
	})
}

func FuzzWebVTTConversionCycle(f *testing.F) {
	f.Add("WEBVTT\n\n00:00.000 --> 00:01.000\n<b>x &amp; y</b>\n")
	f.Add("WEBVTT note\nKind: captions\n\nNOTE x\n\n00:00.000 --> 00:01.000 align:start\n<v A>x\n")
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > maxFuzzInputBytes {
			t.Skip()
		}
		if _, err := webvtt.Parse(input); err != nil {
			return
		}
		document := conversionTestDocument(t, input, "webvtt")
		assertFuzzConversionPolicy(t, document, "subrip")
	})
}

func FuzzStructuredConversion(f *testing.F) {
	f.Add("<b>text</b>", uint8(0))
	f.Add(`C:\Users\operator\secret.srt`, uint8(1))
	f.Add("text\n\nmore", uint8(2))
	f.Fuzz(func(t *testing.T, payload string, mutation uint8) {
		if len(payload) > maxFuzzInputBytes || strings.ContainsRune(payload, '\r') {
			t.Skip()
		}
		document := conversionTestDocument(t, "1\n00:00:01,000 --> 00:00:02,000\ntext\n", "subrip")
		document.Cues[0].Payload = model.Payload{RawText: payload, PlainText: subrip.PlainText(payload), Lines: strings.Split(payload, "\n")}
		pathSentinel := `C:\Users\operator\secret.srt`
		switch mutation % 4 {
		case 1:
			document.Metadata.Title = &pathSentinel
		case 2:
			document.Cues[0].Ordinal = 7
		case 3:
			document.Cues[0].Timing.DurationMilliseconds++
		}

		validationErr := document.Validate()
		result, err := Convert(context.Background(), document, "webvtt", Options{})
		if validationErr != nil {
			if err == nil || len(result.Bytes) != 0 {
				t.Fatalf("invalid structured document conversion = %#v, %v; validation = %v", result, err, validationErr)
			}
			return
		}
		if err != nil {
			var projection *ProjectionError
			if !errors.As(err, &projection) {
				t.Fatalf("unexpected conversion error: %v", err)
			}
			return
		}
		if err := result.LossReport.Validate(document); err != nil {
			t.Fatalf("loss report validation: %v", err)
		}
		canonical, err := json.Marshal(result.LossReport)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(canonical, []byte(pathSentinel)) {
			t.Fatalf("loss report leaked structured path sentinel: %s", canonical)
		}
		assertFuzzConversionPolicy(t, document, "webvtt")
	})
}

func assertFuzzConversionPolicy(t *testing.T, document model.Document, target string) {
	t.Helper()
	permissive, err := Convert(context.Background(), document, target, Options{})
	if err != nil {
		var projection *ProjectionError
		if !errors.As(err, &projection) {
			t.Fatalf("unexpected conversion error: %v", err)
		}
		return
	}
	if err := validateTargetParser(permissive.Bytes, target); err != nil {
		t.Fatalf("target parser rejected output: %v", err)
	}
	repeated, err := Convert(context.Background(), document, target, Options{})
	if err != nil || string(repeated.Bytes) != string(permissive.Bytes) || !reportsEqual(repeated.LossReport, permissive.LossReport) {
		t.Fatalf("conversion is nondeterministic: %v", err)
	}
	strict, strictErr := Convert(context.Background(), document, target, Options{Strict: true})
	if permissive.LossReport.HasLosses() {
		var lossErr *StrictLossError
		if !errors.As(strictErr, &lossErr) || len(strict.Bytes) != 0 || !reportsEqual(lossErr.Report, permissive.LossReport) {
			t.Fatalf("strict loss policy mismatch: %#v, %v", strict, strictErr)
		}
	} else if strictErr != nil || string(strict.Bytes) != string(permissive.Bytes) {
		t.Fatalf("strict loss-free conversion mismatch: %#v, %v", strict, strictErr)
	}
}
