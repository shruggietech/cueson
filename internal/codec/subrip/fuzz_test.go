package subrip

import (
	"reflect"
	"testing"
)

const maxFuzzInputBytes = 64 << 10

func FuzzParse(f *testing.F) {
	for _, seed := range []string{
		"1\n00:00:00,000 --> 00:00:01,000\nhello\n",
		"00:00:00.1 --> 00:00:00.2\rtext\r",
		"not subtitles",
	} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > maxFuzzInputBytes {
			return
		}
		result, err := Parse(input, Options{DetectSpeakers: true})
		if err != nil {
			return
		}
		if rendered, renderErr := Render(result.Cues); renderErr == nil {
			repeated, repeatedErr := Render(result.Cues)
			if repeatedErr != nil || string(repeated) != string(rendered) {
				t.Fatalf("render is nondeterministic: %v", repeatedErr)
			}
			reparsed, parseErr := Parse(string(rendered), Options{})
			if parseErr != nil || len(reparsed.Cues) != len(result.Cues) {
				t.Fatalf("rendered output did not reparse: cues=%d/%d error=%v", len(reparsed.Cues), len(result.Cues), parseErr)
			}
			for index := range reparsed.Cues {
				if reparsed.Cues[index].Timing != result.Cues[index].Timing || !reflect.DeepEqual(reparsed.Cues[index].Payload.Lines, result.Cues[index].Payload.Lines) {
					t.Fatalf("cue %d changed during parser-renderer cycle", index)
				}
			}
		}
	})
}

func FuzzParseTimecodeLine(f *testing.F) {
	for _, seed := range []string{"00:00:00,000 --> 00:00:01,000", "1:2:3.4 --> 5:6:7.8", "bad"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, line string) {
		if len(line) <= maxFuzzInputBytes {
			_, _ = ParseTimecodeLine(line)
		}
	})
}

func FuzzPlainText(f *testing.F) {
	for _, seed := range []string{"<i>hello</i>", "<font color=red>x</font>", "<unknown>x</unknown>"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > maxFuzzInputBytes {
			return
		}
		if first, second := PlainText(input), PlainText(input); first != second {
			t.Fatal("plain-text derivation is nondeterministic")
		}
	})
}
