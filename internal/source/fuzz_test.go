package source

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"
)

const (
	maxEncodedFuzzBytes  = 64 << 10
	maxBasenameFuzzBytes = 1 << 10
)

func FuzzInspectEncoded(f *testing.F) {
	for _, seed := range []string{
		"",
		"Zg==",
		"Zm8=",
		"Zm9v",
		"AP+A",
		"Zg=",
		"Zg===",
		"Z g==",
		"Zg==\n",
		"-w==",
		"%%%%",
		"Zg",
		"Zh==",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, encoded string) {
		if len(encoded) > maxEncodedFuzzBytes {
			return
		}

		count, digest, err := inspectEncoded(encoded)
		repeatedCount, repeatedDigest, repeatedErr := inspectEncoded(encoded)
		if (err == nil) != (repeatedErr == nil) {
			t.Fatalf("inspectEncoded() acceptance changed between calls: first = %v, second = %v", err, repeatedErr)
		}
		if err != nil {
			if err.Error() != repeatedErr.Error() {
				t.Fatalf("inspectEncoded() rejection changed between calls: first = %q, second = %q", err, repeatedErr)
			}
			return
		}
		if count != repeatedCount || digest != repeatedDigest {
			t.Fatalf("inspectEncoded() result changed between calls: first = (%d, %q), second = (%d, %q)", count, digest, repeatedCount, repeatedDigest)
		}

		decoded, err := base64.StdEncoding.Strict().DecodeString(encoded)
		if err != nil {
			t.Fatalf("inspectEncoded() accepted strict-decode rejection: %v", err)
		}
		wantDigest := fmt.Sprintf("%x", sha256.Sum256(decoded))
		if count != int64(len(decoded)) || digest != wantDigest {
			t.Fatalf("inspectEncoded() = (%d, %q), want (%d, %q)", count, digest, len(decoded), wantDigest)
		}
	})
}

func FuzzValidateSafeBasename(f *testing.F) {
	for _, seed := range []string{
		"captions.srt",
		"résumé.vtt",
		"re\u0301sume\u0301.vtt",
		strings.Repeat("a", 255),
		strings.Repeat("a", 256),
		strings.Repeat("é", 126) + ".s",
		strings.Repeat("é", 127) + ".s",
		".",
		"..",
		"../captions.srt",
		"/tmp/captions.srt",
		`C:\Users\cueson-test\captions.srt`,
		`\\server\share\captions.srt`,
		"file:///tmp/captions.srt",
		"dir/captions.srt",
		`dir\captions.srt`,
		"captions.srt:stream",
		"bad\x00name.srt",
		"captions.srt ",
		"captions.srt.",
		"CON",
		"con.txt",
		"COM¹.srt",
		"LPT³.vtt",
		string([]byte{0xff, 'x'}),
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, name string) {
		if len(name) > maxBasenameFuzzBytes {
			return
		}

		err := validateSafeBasename(name)
		repeatedErr := validateSafeBasename(name)
		if (err == nil) != (repeatedErr == nil) {
			t.Fatalf("validateSafeBasename() acceptance changed between calls: first = %v, second = %v", err, repeatedErr)
		}
		if err != nil {
			if err.Error() != repeatedErr.Error() {
				t.Fatalf("validateSafeBasename() rejection changed between calls: first = %q, second = %q", err, repeatedErr)
			}
			return
		}
		if !utf8.ValidString(name) {
			t.Fatal("validateSafeBasename() accepted invalid UTF-8")
		}
	})
}
