//go:build windows

package source

import (
	"errors"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/shruggietech/cueson/internal/model"
	"golang.org/x/sys/windows"
)

const (
	filetimeEpochTicks uint64 = 116444736000000000
	filetimeQuantum           = int64(100)
	windowsPrecision          = "100ns"
)

// openRegularNoFollow opens one regular file without resolving a reparse point.
// Omitting write and delete sharing keeps the named object stable while the
// caller captures, reads, applies, or verifies metadata.
func openRegularNoFollow(path string, writeAttributes bool) (*os.File, error) {
	pathPointer, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, fmt.Errorf("encode Windows path: %w", err)
	}

	access := uint32(windows.GENERIC_READ)
	if writeAttributes {
		access |= windows.FILE_WRITE_ATTRIBUTES
	}
	handle, err := windows.CreateFile(
		pathPointer,
		access,
		windows.FILE_SHARE_READ,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_FLAG_OPEN_REPARSE_POINT|windows.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("open regular file without following links: %w", err)
	}

	file := os.NewFile(uintptr(handle), path)
	if file == nil {
		_ = windows.CloseHandle(handle)
		return nil, fmt.Errorf("open regular file without following links: invalid Windows handle")
	}
	if _, err := nativeFileInformation(file); err != nil {
		_ = file.Close()
		return nil, err
	}
	return file, nil
}

func captureNativeTimestamps(file *os.File) (model.Timestamps, error) {
	information, err := nativeFileInformation(file)
	if err != nil {
		return model.Timestamps{}, err
	}

	created, createdAvailable, err := timestampFromFiletime(information.CreationTime)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture creation timestamp: %w", err)
	}
	modified, modifiedAvailable, err := timestampFromFiletime(information.LastWriteTime)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture modification timestamp: %w", err)
	}
	accessed, accessedAvailable, err := timestampFromFiletime(information.LastAccessTime)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture access timestamp: %w", err)
	}

	result := model.Timestamps{CreatedSource: "unavailable"}
	if createdAvailable {
		result.Created = created
		result.CreatedSource = "windows_creation_time"
	}
	if modifiedAvailable {
		result.Modified = modified
	}
	if accessedAvailable {
		result.Accessed = accessed
	}
	return result, nil
}

func applyNativeTimestamps(file *os.File, timestamps model.Timestamps) []TimestampResult {
	results := []TimestampResult{
		unavailableTimestampResult(TimestampCreated),
		unavailableTimestampResult(TimestampModified),
		unavailableTimestampResult(TimestampAccessed),
	}

	var creationTime *windows.Filetime
	var modificationTime *windows.Filetime
	var accessTime *windows.Filetime
	attempted := [3]bool{}

	if timestamps.Created != nil {
		switch timestamps.CreatedSource {
		case "windows_creation_time", "birthtime":
			if converted, result := prepareWindowsTimestamp(TimestampCreated, timestamps.Created); converted != nil {
				creationTime = converted
				attempted[0] = true
				results[0] = result
			} else {
				results[0] = result
			}
		case "ctime_fallback":
			results[0] = TimestampResult{Kind: TimestampCreated, Status: TimestampUnsupported, EffectivePrecision: windowsPrecision, Detail: "ctime fallback is not a restorable creation timestamp"}
		case "unavailable":
			results[0] = TimestampResult{Kind: TimestampCreated, Status: TimestampFailed, Detail: "non-null creation timestamp has unavailable provenance"}
		default:
			results[0] = TimestampResult{Kind: TimestampCreated, Status: TimestampFailed, Detail: fmt.Sprintf("unrecognized creation timestamp provenance %q", timestamps.CreatedSource)}
		}
	}
	if timestamps.Modified != nil {
		if converted, result := prepareWindowsTimestamp(TimestampModified, timestamps.Modified); converted != nil {
			modificationTime = converted
			attempted[1] = true
			results[1] = result
		} else {
			results[1] = result
		}
	}
	if timestamps.Accessed != nil {
		if converted, result := prepareWindowsTimestamp(TimestampAccessed, timestamps.Accessed); converted != nil {
			accessTime = converted
			attempted[2] = true
			results[2] = result
		} else {
			results[2] = result
		}
	}

	if !attempted[0] && !attempted[1] && !attempted[2] {
		return results
	}
	anchor, err := nativeFileInformation(file)
	if err != nil {
		return failAttemptedTimestamps(results, attempted, err)
	}

	setter, err := openRegularNoFollow(file.Name(), true)
	if err != nil {
		return failAttemptedTimestamps(results, attempted, err)
	}
	setterInformation, err := nativeFileInformation(setter)
	if err == nil && !sameNativeFile(anchor, setterInformation) {
		err = fmt.Errorf("timestamp destination identity changed before application")
	}
	if err == nil {
		err = withWindowsHandle(setter, func(handle windows.Handle) error {
			return windows.SetFileTime(handle, creationTime, accessTime, modificationTime)
		})
	}
	if closeErr := setter.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf("close timestamp setter: %w", closeErr)
	}
	if err != nil {
		if isUnsupportedWindowsTimestampError(err) {
			return markAttemptedTimestamps(results, attempted, TimestampUnsupported, err)
		}
		return failAttemptedTimestamps(results, attempted, err)
	}

	verifier, err := openRegularNoFollow(file.Name(), false)
	if err != nil {
		return failAttemptedTimestamps(results, attempted, err)
	}
	verified, err := nativeFileInformation(verifier)
	if err == nil && !sameNativeFile(anchor, verified) {
		err = fmt.Errorf("timestamp destination identity changed before verification")
	}
	if closeErr := verifier.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf("close timestamp verifier: %w", closeErr)
	}
	if err != nil {
		return failAttemptedTimestamps(results, attempted, err)
	}

	actual := [3]windows.Filetime{verified.CreationTime, verified.LastWriteTime, verified.LastAccessTime}
	wanted := [3]*windows.Filetime{creationTime, modificationTime, accessTime}
	for index := range results {
		if !attempted[index] {
			continue
		}
		if equalFiletime(actual[index], *wanted[index]) {
			results[index].Status = TimestampRestored
			results[index].Detail = "timestamp restored and verified"
			continue
		}
		results[index].Status = TimestampUnsupported
		results[index].Detail = "destination filesystem did not preserve the exact 100ns timestamp"
	}
	return results
}

func nativeFileInformation(file *os.File) (windows.ByHandleFileInformation, error) {
	var information windows.ByHandleFileInformation
	err := withWindowsHandle(file, func(handle windows.Handle) error {
		return windows.GetFileInformationByHandle(handle, &information)
	})
	if err != nil {
		return windows.ByHandleFileInformation{}, fmt.Errorf("query Windows file information: %w", err)
	}
	if information.FileAttributes&(windows.FILE_ATTRIBUTE_REPARSE_POINT|windows.FILE_ATTRIBUTE_DIRECTORY|windows.FILE_ATTRIBUTE_DEVICE) != 0 {
		return windows.ByHandleFileInformation{}, fmt.Errorf("path does not identify a no-follow regular file")
	}
	return information, nil
}

func withWindowsHandle(file *os.File, operation func(windows.Handle) error) error {
	connection, err := file.SyscallConn()
	if err != nil {
		return err
	}
	var operationErr error
	if err := connection.Control(func(handle uintptr) {
		operationErr = operation(windows.Handle(handle))
	}); err != nil {
		return err
	}
	return operationErr
}

func timestampFromFiletime(filetime windows.Filetime) (*model.Timestamp, bool, error) {
	raw := rawFiletime(filetime)
	if raw == 0 {
		return nil, false, nil
	}
	unixNS, err := filetimeToUnixNS(raw)
	if err != nil {
		return nil, false, err
	}
	return &model.Timestamp{ISO: time.Unix(0, unixNS).UTC().Format(time.RFC3339Nano), UnixNS: unixNS}, true, nil
}

func filetimeToUnixNS(raw uint64) (int64, error) {
	const unixTickLimit = uint64(math.MaxInt64 / filetimeQuantum)
	minimum := filetimeEpochTicks - unixTickLimit
	maximum := filetimeEpochTicks + unixTickLimit
	if raw < minimum || raw > maximum {
		return 0, fmt.Errorf("FILETIME value %d is outside signed Unix nanoseconds", raw)
	}
	return (int64(raw) - int64(filetimeEpochTicks)) * filetimeQuantum, nil
}

func unixNSToFiletime(unixNS int64) (windows.Filetime, error) {
	if unixNS%filetimeQuantum != 0 {
		return windows.Filetime{}, fmt.Errorf("unix nanoseconds %d are not representable at 100ns precision", unixNS)
	}
	raw := uint64(unixNS/filetimeQuantum + int64(filetimeEpochTicks))
	return windows.Filetime{LowDateTime: uint32(raw), HighDateTime: uint32(raw >> 32)}, nil
}

func rawFiletime(filetime windows.Filetime) uint64 {
	return uint64(filetime.HighDateTime)<<32 | uint64(filetime.LowDateTime)
}

func equalFiletime(left, right windows.Filetime) bool {
	return left.LowDateTime == right.LowDateTime && left.HighDateTime == right.HighDateTime
}

func sameNativeFile(left, right windows.ByHandleFileInformation) bool {
	return left.VolumeSerialNumber == right.VolumeSerialNumber && left.FileIndexHigh == right.FileIndexHigh && left.FileIndexLow == right.FileIndexLow
}

func unavailableTimestampResult(kind TimestampKind) TimestampResult {
	return TimestampResult{Kind: kind, Status: TimestampUnavailable, Detail: "source timestamp is unavailable"}
}

func prepareWindowsTimestamp(kind TimestampKind, timestamp *model.Timestamp) (*windows.Filetime, TimestampResult) {
	converted, err := unixNSToFiletime(timestamp.UnixNS)
	if err != nil {
		return nil, TimestampResult{Kind: kind, Status: TimestampUnsupported, EffectivePrecision: windowsPrecision, Detail: err.Error()}
	}
	return &converted, TimestampResult{Kind: kind, EffectivePrecision: windowsPrecision}
}

func failAttemptedTimestamps(results []TimestampResult, attempted [3]bool, err error) []TimestampResult {
	return markAttemptedTimestamps(results, attempted, TimestampFailed, err)
}

func markAttemptedTimestamps(results []TimestampResult, attempted [3]bool, status TimestampStatus, err error) []TimestampResult {
	for index := range results {
		if attempted[index] {
			results[index].Status = status
			results[index].Detail = err.Error()
		}
	}
	return results
}

func isUnsupportedWindowsTimestampError(err error) bool {
	return errors.Is(err, windows.ERROR_NOT_SUPPORTED) || errors.Is(err, windows.ERROR_INVALID_FUNCTION) || errors.Is(err, windows.ERROR_CALL_NOT_IMPLEMENTED)
}
