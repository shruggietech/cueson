// Package testutil provides domain-neutral fixture and conformance-test support.
package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

var drivePrefix = regexp.MustCompile(`^[A-Za-z]:`)

// PortablePathKey validates a slash-relative repository path and returns its
// Unicode-normalized, caseless comparison key.
func PortablePathKey(value string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("path is empty")
	}
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("path is not valid UTF-8")
	}
	if strings.Contains(value, `\`) {
		return "", fmt.Errorf("path contains a backslash separator")
	}
	if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "//") || drivePrefix.MatchString(value) {
		return "", fmt.Errorf("path is absolute")
	}
	if path.Clean(value) != value {
		return "", fmt.Errorf("path is not in clean slash-relative form")
	}
	components := strings.Split(value, "/")
	for _, component := range components {
		if err := validatePortableComponent(component); err != nil {
			return "", err
		}
	}
	return norm.NFD.String(cases.Fold().String(norm.NFD.String(value))), nil
}

// CheckNoForbidden rejects explicit local identifiers and their common path
// and JSON-escaped forms without treating generic path-looking text as secret.
func CheckNoForbidden(fixtureID, surface string, payload []byte, sentinels ...string) error {
	for _, sentinel := range sentinels {
		for _, variant := range forbiddenVariants(sentinel) {
			if len(variant) > 0 && bytes.Contains(payload, variant) {
				return comparisonError(fixtureID, surface, "forbidden local identifier detected")
			}
		}
	}
	return nil
}

func forbiddenVariants(sentinel string) [][]byte {
	if sentinel == "" {
		return nil
	}
	values := []string{
		sentinel,
		filepath.ToSlash(sentinel),
		strings.ReplaceAll(sentinel, `\`, "/"),
		strings.ReplaceAll(sentinel, "/", `\`),
	}
	for _, value := range append([]string(nil), values...) {
		if encoded, err := json.Marshal(value); err == nil && len(encoded) >= 2 {
			values = append(values, string(encoded[1:len(encoded)-1]))
		}
	}
	seen := make(map[string]struct{}, len(values))
	result := make([][]byte, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, []byte(value))
	}
	return result
}

func validatePortableComponent(component string) error {
	if component == "" || component == "." || component == ".." {
		return fmt.Errorf("path contains an empty or relative component")
	}
	if len([]byte(component)) > 255 {
		return fmt.Errorf("path component exceeds 255 UTF-8 bytes")
	}
	if strings.HasSuffix(component, ".") || strings.HasSuffix(component, " ") {
		return fmt.Errorf("path component ends with a dot or space")
	}
	for _, character := range component {
		if character < 0x20 || character == 0x7f || strings.ContainsRune(`<>:"|?*`, character) {
			return fmt.Errorf("path component contains a non-portable character")
		}
	}
	stem := strings.ToUpper(strings.SplitN(component, ".", 2)[0])
	if isWindowsReservedStem(stem) {
		return fmt.Errorf("path component uses a reserved device name")
	}
	return nil
}

func isWindowsReservedStem(stem string) bool {
	switch stem {
	case "CON", "PRN", "AUX", "NUL", "COM¹", "COM²", "COM³", "LPT¹", "LPT²", "LPT³":
		return true
	}
	if len(stem) == 4 && (strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT")) && stem[3] >= '1' && stem[3] <= '9' {
		return true
	}
	return false
}
