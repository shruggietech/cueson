//go:build !windows && !linux && !darwin

package source

import (
	"fmt"
	"os"

	"github.com/shruggietech/cueson/internal/model"
)

func openRegularNoFollow(path string, writeAttributes bool) (*os.File, error) {
	if writeAttributes {
		return os.OpenFile(path, os.O_RDWR, 0)
	}
	return os.Open(path)
}

func captureNativeTimestamps(file *os.File) (model.Timestamps, error) {
	return model.Timestamps{CreatedSource: "unavailable"}, nil
}

func applyNativeTimestamps(file *os.File, timestamps model.Timestamps) []TimestampResult {
	return []TimestampResult{
		unsupportedOrUnavailable(TimestampCreated, timestamps.Created),
		unsupportedOrUnavailable(TimestampModified, timestamps.Modified),
		unsupportedOrUnavailable(TimestampAccessed, timestamps.Accessed),
	}
}

func unsupportedOrUnavailable(kind TimestampKind, value *model.Timestamp) TimestampResult {
	if value == nil {
		return unavailableResult(kind)
	}
	return unsupportedResult(kind, fmt.Sprintf("timestamp restoration is unsupported on %s", filePlatform()))
}

func filePlatform() string { return "this platform" }
