package codec

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/shruggietech/cueson/internal/model"
)

// Stable detection diagnostic codes.
const DiagnosticExtensionDisagreement = "format_extension_disagreement"

// Evidence is one detector's content-based format claim.
type Evidence struct {
	Matched    bool
	Confidence int
	Reason     string
}

// Selection records the chosen registration and the evidence used to choose it.
type Selection struct {
	Format          Format
	Registration    Registration
	Explicit        bool
	ContentFormat   *Format
	ExtensionFormat *Format
	Diagnostics     []model.Diagnostic
}

// AmbiguousFormatError reports equal strongest content claims.
type AmbiguousFormatError struct {
	Formats []Format
}

func (err *AmbiguousFormatError) Error() string {
	values := make([]string, len(err.Formats))
	for index := range err.Formats {
		values[index] = string(err.Formats[index])
	}
	return fmt.Sprintf("source format is ambiguous between %s; provide --format", strings.Join(values, ", "))
}

// Select resolves an explicit selector or deterministic content and extension evidence.
func (registry *Registry) Select(data []byte, fileName, requested string) (Selection, error) {
	normalized := normalizeToken(requested)
	if normalized != "" && normalized != string(FormatAuto) {
		registration, exists := registry.Lookup(normalized)
		if !exists {
			return Selection{}, &UnknownFormatError{Selector: requested}
		}
		return Selection{Format: registration.Format, Registration: registration, Explicit: true}, nil
	}

	contentIndex, contentFormat, err := registry.selectContent(data)
	if err != nil {
		return Selection{}, err
	}
	extensionIndex, extensionFormat := registry.selectExtension(fileName)

	selectedIndex := contentIndex
	if selectedIndex < 0 {
		selectedIndex = extensionIndex
	}
	if selectedIndex < 0 {
		return Selection{}, &UnknownFormatError{}
	}
	registration := cloneRegistration(registry.registrations[selectedIndex])
	selection := Selection{Format: registration.Format, Registration: registration, ContentFormat: contentFormat, ExtensionFormat: extensionFormat}
	if contentFormat != nil && extensionFormat != nil && *contentFormat != *extensionFormat {
		selection.Diagnostics = append(selection.Diagnostics, model.Diagnostic{
			Severity: "warning",
			Code:     DiagnosticExtensionDisagreement,
			Message:  fmt.Sprintf("content identifies format %q but file extension identifies %q; content evidence was used", *contentFormat, *extensionFormat),
		})
	}
	return selection, nil
}

func (registry *Registry) selectContent(data []byte) (int, *Format, error) {
	type candidate struct {
		index      int
		format     Format
		confidence int
	}
	candidates := make([]candidate, 0, len(registry.registrations))
	for index := range registry.registrations {
		detector := registry.registrations[index].Detect
		if detector == nil {
			continue
		}
		evidence := detector(data)
		if !evidence.Matched || evidence.Confidence <= 0 {
			continue
		}
		candidates = append(candidates, candidate{index: index, format: registry.registrations[index].Format, confidence: evidence.Confidence})
	}
	if len(candidates) == 0 {
		return -1, nil, nil
	}
	sort.Slice(candidates, func(left, right int) bool {
		if candidates[left].confidence != candidates[right].confidence {
			return candidates[left].confidence > candidates[right].confidence
		}
		return candidates[left].format < candidates[right].format
	})
	strongest := candidates[0].confidence
	tied := []Format{candidates[0].format}
	for index := 1; index < len(candidates) && candidates[index].confidence == strongest; index++ {
		tied = append(tied, candidates[index].format)
	}
	if len(tied) > 1 {
		return -1, nil, &AmbiguousFormatError{Formats: tied}
	}
	format := candidates[0].format
	return candidates[0].index, &format, nil
}

func (registry *Registry) selectExtension(fileName string) (int, *Format) {
	extension := normalizeExtension(filepath.Ext(fileName))
	if extension == "" {
		return -1, nil
	}
	index, exists := registry.extensions[extension]
	if !exists {
		return -1, nil
	}
	format := registry.registrations[index].Format
	return index, &format
}
