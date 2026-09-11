package subrip

import "testing"

func TestParseTimecodeLineVariants(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		line        string
		start       int64
		end         int64
		coordinates *Coordinates
	}{
		{name: "canonical", line: "00:00:01,250 --> 00:00:04,200", start: 1250, end: 4200},
		{name: "period and short fraction", line: "1:02:03.4 --> 1:02:04.05", start: 3723400, end: 3724050},
		{name: "coordinates", line: "00:00:01,000 --> 00:00:02,000 Y2:40 X1:10 Y1:20 X2:30", start: 1000, end: 2000, coordinates: &Coordinates{X1: 10, X2: 30, Y1: 20, Y2: 40}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseTimecodeLine(test.line)
			if err != nil {
				t.Fatalf("ParseTimecodeLine() error = %v", err)
			}
			if got.StartMilliseconds != test.start || got.EndMilliseconds != test.end {
				t.Fatalf("timing = %d..%d, want %d..%d", got.StartMilliseconds, got.EndMilliseconds, test.start, test.end)
			}
			if !equalCoordinates(got.Coordinates, test.coordinates) {
				t.Fatalf("coordinates = %#v, want %#v", got.Coordinates, test.coordinates)
			}
		})
	}
}

func TestParseTimecodeLineRejectsMalformedValues(t *testing.T) {
	t.Parallel()
	for _, line := range []string{
		"00:60:00,000 --> 00:60:01,000",
		"00:00:02,000 --> 00:00:01,000",
		"00:00:01,000 -> 00:00:02,000",
		"00:00:01,000 --> 00:00:02,000 X1:1 X2:2 Y1:3",
		"00:00:01,000 --> 00:00:02,000 X1:1 X1:2 X2:3 Y1:4 Y2:5",
	} {
		if _, err := ParseTimecodeLine(line); err == nil {
			t.Errorf("ParseTimecodeLine(%q) unexpectedly succeeded", line)
		}
	}
}

func equalCoordinates(left, right *Coordinates) bool {
	if left == nil || right == nil {
		return left == right
	}
	return *left == *right
}
