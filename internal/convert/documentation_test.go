package convert

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestConversionDocumentationListsEveryStableLossCodeExactlyOnce(t *testing.T) {
	t.Parallel()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate conversion package")
	}
	path := filepath.Join(filepath.Dir(file), "..", "..", "docs", "conversion.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	document := string(data)
	for code := range knownLossCodes {
		if count := strings.Count(document, "`"+code+"`"); count != 1 {
			t.Errorf("loss code %q appears %d times, want exactly once", code, count)
		}
	}
	rows := regexp.MustCompile(`(?m)^\| \x60(conversion_[a-z0-9_]+)\x60 \|`).FindAllStringSubmatch(document, -1)
	if len(rows) != len(knownLossCodes) {
		t.Fatalf("documented loss rows = %d, want %d", len(rows), len(knownLossCodes))
	}
	for _, row := range rows {
		if !IsKnownLossCode(row[1]) {
			t.Errorf("documentation contains unknown loss code %q", row[1])
		}
	}
	for _, phrase := range []string{"--strict", "writes no target payload", "8,192", "caller paths", "cueson restore"} {
		if !strings.Contains(document, phrase) {
			t.Errorf("conversion documentation is missing %q", phrase)
		}
	}
}
