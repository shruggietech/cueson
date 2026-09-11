package webvtt

import (
	"testing"

	"github.com/shruggietech/cueson/internal/model"
)

func FuzzParseTimestamp(f *testing.F) {
	for _, seed := range []string{"00:00.000", "00:00:00.000", "59:59.999", "999999999999999999999:00:00.000", "00:00,000"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) { _, _ = ParseTimestamp(input) })
}

func FuzzParseWebVTT(f *testing.F) {
	for _, seed := range []string{"", "WEBVTT\n\n00:00.000 --> 00:01.000\nx\n", "WEBVTT\r\rNOTE x\rkept\r", "WEBVTT\n\nREGION\nid:r width:80%\n"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) { _, _ = Parse(input) })
}

func FuzzSettingsAndMarkup(f *testing.F) {
	f.Add("line:10%", "<v A>x &amp;</v>", int64(0), int64(1000))
	f.Fuzz(func(t *testing.T, settings, payload string, start, end int64) {
		_, _, _ = parseCueSettings(settings, 0, "cue-000000")
		_, _, _ = parseRegionSettings(settings, 0)
		if start < 0 {
			start = 0
		}
		if end <= start {
			end = start + 1
		}
		_, _, _, _ = scanPayload(payload, start, end, 0, "cue-000000")
	})
}

func FuzzRenderParseCycle(f *testing.F) {
	f.Add("hello", int64(0), int64(1000))
	f.Fuzz(func(t *testing.T, payload string, start, end int64) {
		if start < 0 || end <= start || end-start > 86_400_000 || stringsContainsUnsafePayload(payload) {
			return
		}
		native := &model.WebVTTCueData{TimingLineRaw: "", SettingsRaw: "", Settings: map[string]string{}, SettingOccurrences: []model.WebVTTSettingOccurrence{}, RawPayload: payload, RawPayloadLines: []string{payload}}
		cue := model.Cue{ID: "cue-000000", SourceOrder: 0, Timing: model.Timing{StartMilliseconds: start, EndMilliseconds: end, DurationMilliseconds: end - start}, Payload: model.Payload{RawText: payload, Lines: []string{payload}}, FormatData: model.CueFormatData{WebVTT: native}}
		data := &model.WebVTTDocumentData{Signature: "WEBVTT", SignatureLineRaw: "WEBVTT", MetadataLines: []string{}, Blocks: []model.WebVTTBlock{}}
		rendered, err := Render(model.Document{Format: "webvtt", Cues: []model.Cue{cue}, FormatData: model.DocumentFormatData{WebVTT: data}}, RenderOptions{})
		if err == nil {
			_, _ = Parse(string(rendered.Bytes))
		}
	})
}

func stringsContainsUnsafePayload(value string) bool {
	for index := 0; index < len(value); index++ {
		if value[index] == '\r' || value[index] == '\n' {
			return true
		}
	}
	return value == ""
}
