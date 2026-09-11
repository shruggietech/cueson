package subrip

import (
	"reflect"
	"testing"
)

func TestPlainTextRemovesOnlyDocumentedTags(t *testing.T) {
	t.Parallel()
	input := `<i>Hello <b>there</b></i> <font color="#fff">friend</font> <blink>still literal</blink>`
	want := `Hello there friend <blink>still literal</blink>`
	if got := PlainText(input); got != want {
		t.Fatalf("PlainText() = %q, want %q", got, want)
	}
}

func TestPlainTextPreservesMalformedDocumentedTags(t *testing.T) {
	t.Parallel()
	input := `<i>unclosed <b>crossed</i></b> <font color="red">also unclosed`
	if got := PlainText(input); got != input {
		t.Fatalf("PlainText() = %q, want malformed tags preserved", got)
	}
}

func TestDetectSpeakerDoesNotMutateLines(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		lines []string
		want  *string
	}{
		{name: "same line", lines: []string{"Alice: Hello", "Again"}, want: stringPointer("Alice")},
		{name: "label line", lines: []string{"Speaker 1:", "Hello"}, want: stringPointer("Speaker 1")},
		{name: "time", lines: []string{"Meet at 3:30"}},
		{name: "question", lines: []string{"Who?: really"}},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			before := append([]string(nil), test.lines...)
			got := DetectSpeaker(test.lines)
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("DetectSpeaker() = %#v, want %#v", got, test.want)
			}
			if !reflect.DeepEqual(test.lines, before) {
				t.Fatalf("DetectSpeaker() mutated lines: %#v", test.lines)
			}
		})
	}
}

func stringPointer(value string) *string { return &value }
