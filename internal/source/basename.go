package source

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

var portableFolder = cases.Fold()

func validateSafeBasename(name string) error {
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("file_name %q is not a safe basename", name)
	}
	if !utf8.ValidString(name) || len(name) > 255 {
		return fmt.Errorf("file_name %q is not a valid portable basename", name)
	}
	for _, character := range name {
		if character < 0x20 || character == 0x7f || strings.ContainsRune(`<>:"/\|?*`, character) {
			return fmt.Errorf("file_name %q contains a prohibited character", name)
		}
	}
	if strings.HasSuffix(name, " ") || strings.HasSuffix(name, ".") {
		return fmt.Errorf("file_name %q has a prohibited trailing character", name)
	}
	stem := name
	if dot := strings.IndexByte(stem, '.'); dot >= 0 {
		stem = stem[:dot]
	}
	switch strings.ToUpper(stem) {
	case "CON", "PRN", "AUX", "NUL", "CLOCK$", "CONIN$", "CONOUT$",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9":
		return fmt.Errorf("file_name %q uses a reserved Windows device stem", name)
	}
	return nil
}

func portableIdentity(value string) string {
	decomposed := norm.NFD.String(value)
	return norm.NFD.String(portableFolder.String(decomposed))
}
