package source

import (
	"fmt"

	"github.com/shruggietech/cueson/internal/model"
)

func timestampResults(filePath string, timestamps model.Timestamps, mode MetadataMode) ([]TimestampResult, []string, error) {
	if mode == MetadataNone {
		return nil, nil, nil
	}
	file, err := openRegularNoFollow(filePath, true)
	if err != nil {
		return nil, nil, fmt.Errorf("open restored destination metadata handle: %w", err)
	}
	results := applyNativeTimestamps(file, timestamps)
	closeErr := file.Close()
	warnings := make([]string, 0, len(results))
	for _, result := range results {
		switch result.Status {
		case TimestampFailed:
			return results, warnings, fmt.Errorf("restore %s timestamp: %s", result.Kind, result.Detail)
		case TimestampUnsupported:
			if mode == MetadataStrict {
				return results, warnings, fmt.Errorf("restore %s timestamp: %s", result.Kind, result.Detail)
			}
			warnings = append(warnings, fmt.Sprintf("%s timestamp for %q is unsupported: %s", result.Kind, filePath, result.Detail))
		}
	}
	if closeErr != nil {
		return results, warnings, fmt.Errorf("close restored destination metadata handle: %w", closeErr)
	}
	return results, warnings, nil
}

//lint:ignore U1000 Used by non-Windows adapters selected through build constraints.
func unavailableResult(kind TimestampKind) TimestampResult {
	return TimestampResult{Kind: kind, Status: TimestampUnavailable, Detail: "source timestamp was unavailable"}
}

//lint:ignore U1000 Used by non-Windows adapters selected through build constraints.
func unsupportedResult(kind TimestampKind, detail string) TimestampResult {
	return TimestampResult{Kind: kind, Status: TimestampUnsupported, Detail: detail}
}
