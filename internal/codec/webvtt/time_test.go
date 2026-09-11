package webvtt

import "testing"

func TestParseTimestampStrictGrammarAndBounds(t *testing.T) {
	accepted := map[string]int64{
		"00:00.000":     0,
		"59:59.999":     3_599_999,
		"00:00:00.001":  1,
		"123:45:06.007": 445_506_007,
	}
	for input, want := range accepted {
		got, err := ParseTimestamp(input)
		if err != nil || got != want {
			t.Fatalf("ParseTimestamp(%q) = (%d, %v), want (%d, nil)", input, got, err, want)
		}
	}
	for _, input := range []string{"0:00.000", "00:0.000", "00:00,000", "00:00.00", "60:00.000", "00:60.000", "999999999999999999999:00:00.000"} {
		if _, err := ParseTimestamp(input); err == nil {
			t.Fatalf("ParseTimestamp(%q) unexpectedly succeeded", input)
		}
	}
}

func TestParseTimingLineRequiresConformingDelimiterAndIncreasingRange(t *testing.T) {
	got, err := ParseTimingLine("00:01.000 --> 00:02.500 line:20% align:start")
	if err != nil {
		t.Fatal(err)
	}
	if got.StartMilliseconds != 1000 || got.EndMilliseconds != 2500 || got.SettingsRaw != "line:20% align:start" {
		t.Fatalf("ParseTimingLine() = %#v", got)
	}
	for _, input := range []string{"00:01.000-->00:02.000", "00:02.000 --> 00:02.000", "00:03.000 --> 00:02.000", "00:01.000 -> 00:02.000"} {
		if _, err := ParseTimingLine(input); err == nil {
			t.Fatalf("ParseTimingLine(%q) unexpectedly succeeded", input)
		}
	}
}
