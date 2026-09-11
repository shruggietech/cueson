package subrip

import "testing"

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		"1\n00:00:00,000 --> 00:00:01,000\nhello\n",
		"00:00:00.1 --> 00:00:00.2\rtext\r",
		"not subtitles",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		result, err := Parse(input, Options{DetectSpeakers: true})
		if err != nil {
			return
		}
		if rendered, renderErr := Render(result.Cues); renderErr == nil {
			_, _ = Parse(string(rendered), Options{})
		}
	})
}

func FuzzParseTimecodeLine(f *testing.F) {
	for _, seed := range []string{"00:00:00,000 --> 00:00:01,000", "1:2:3.4 --> 5:6:7.8", "bad"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, line string) { _, _ = ParseTimecodeLine(line) })
}

func FuzzPlainText(f *testing.F) {
	for _, seed := range []string{"<i>hello</i>", "<font color=red>x</font>", "<unknown>x</unknown>"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) { _ = PlainText(input) })
}
