package cli

import (
	"context"
	"strings"
	"testing"
)

func TestScriptedConversionInvocationMatrixAndEncodingPolicy(t *testing.T) {
	aliases := map[string]string{"srt": "subrip", "vtt": "webvtt", "ass": "ass", "ssa": "ssa"}
	for source, sourceCanonical := range aliases {
		for target, targetCanonical := range aliases {
			parsed, err := parseInvocation([]string{"convert", "input.bin", "--from", source, "--to", target})
			if err != nil || parsed.convert.from != sourceCanonical || parsed.convert.target != targetCanonical {
				t.Fatalf("selection %s to %s = %#v, %v", source, target, parsed.convert, err)
			}
		}
	}
	for _, format := range []string{"ass", "ssa"} {
		for _, encoding := range []string{"utf-16le", "utf-16be", "windows-1252", "iso-8859-1"} {
			status, stdout, stderr := runForTest(context.Background(), []string{"convert", "input.bin", "--from", format, "--to", "srt", "--encoding", encoding})
			if status != ExitInvocation || stdout != "" || !strings.Contains(stderr, "require UTF-8") {
				t.Fatalf("incompatible %s source encoding %s = %d %q %q", format, encoding, status, stdout, stderr)
			}
		}
		for _, encoding := range []string{"utf-8", "utf-8-bom"} {
			if _, err := parseInvocation([]string{"convert", "input.bin", "--from", format, "--to", "srt", "--encoding", encoding}); err != nil {
				t.Fatalf("compatible selector %s %s: %v", format, encoding, err)
			}
		}
	}
}
