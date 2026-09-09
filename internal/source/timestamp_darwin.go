//go:build darwin

package source

import (
	"fmt"
	"os"
	"time"

	"github.com/shruggietech/cueson/internal/model"
	"golang.org/x/sys/unix"
)

const darwinTimestampPrecision = "1us"

func openRegularNoFollow(path string, writeAttributes bool) (*os.File, error) {
	flags := unix.O_CLOEXEC | unix.O_NOFOLLOW | unix.O_NONBLOCK
	if writeAttributes {
		flags |= unix.O_RDWR
	} else {
		flags |= unix.O_RDONLY
	}

	fd, err := unix.Open(path, flags, 0)
	if err != nil {
		return nil, fmt.Errorf("open without following links: %w", err)
	}
	file := os.NewFile(uintptr(fd), path)
	if file == nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("open without following links: invalid file descriptor")
	}

	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("inspect opened file: %w", err)
	}
	if stat.Mode&unix.S_IFMT != unix.S_IFREG {
		_ = file.Close()
		return nil, fmt.Errorf("opened path is not a regular file")
	}
	return file, nil
}

func captureNativeTimestamps(file *os.File) (model.Timestamps, error) {
	stat, err := fstatDarwin(file)
	if err != nil {
		return model.Timestamps{}, err
	}

	modified, err := modelTimestampFromDarwin(stat.Mtim)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture modification time: %w", err)
	}
	accessed, err := modelTimestampFromDarwin(stat.Atim)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture access time: %w", err)
	}

	timestamps := model.Timestamps{
		Modified:      modified,
		Accessed:      accessed,
		CreatedSource: "unavailable",
	}
	created, provenance, err := capturedDarwinCreation(stat.Btim, stat.Ctim)
	if err != nil {
		return model.Timestamps{}, err
	}
	timestamps.Created = created
	timestamps.CreatedSource = provenance
	return timestamps, nil
}

func capturedDarwinCreation(birth, change unix.Timespec) (*model.Timestamp, string, error) {
	if timespecAvailable(birth) {
		created, err := modelTimestampFromDarwin(birth)
		if err != nil {
			return nil, "", fmt.Errorf("capture birth time: %w", err)
		}
		if timespecAvailable(change) && birth.Sec == change.Sec && birth.Nsec == change.Nsec {
			return created, "ctime_fallback", nil
		}
		return created, "birthtime", nil
	}
	if timespecAvailable(change) {
		created, err := modelTimestampFromDarwin(change)
		if err != nil {
			return nil, "", fmt.Errorf("capture change-time fallback: %w", err)
		}
		return created, "ctime_fallback", nil
	}
	return nil, "unavailable", nil
}

func applyNativeTimestamps(file *os.File, timestamps model.Timestamps) []TimestampResult {
	results := []TimestampResult{
		metadataUnavailable(TimestampCreated),
		metadataUnavailable(TimestampModified),
		metadataUnavailable(TimestampAccessed),
	}
	if timestamps.Created != nil {
		results[0] = TimestampResult{
			Kind:   TimestampCreated,
			Status: TimestampUnsupported,
			Detail: "macOS creation-time restoration is not supported by the pure-Go S004 adapter",
		}
	}
	if timestamps.Modified == nil && timestamps.Accessed == nil {
		return results
	}

	stat, err := fstatDarwin(file)
	if err != nil {
		return failRequestedDarwinTimes(results, timestamps, fmt.Sprintf("inspect destination timestamps: %v", err))
	}

	accessNanos, err := timespecUnixNanoDarwin(stat.Atim)
	if err != nil {
		return failRequestedDarwinTimes(results, timestamps, fmt.Sprintf("read current access time: %v", err))
	}
	modifiedNanos, err := timespecUnixNanoDarwin(stat.Mtim)
	if err != nil {
		return failRequestedDarwinTimes(results, timestamps, fmt.Sprintf("read current modification time: %v", err))
	}
	if timestamps.Accessed != nil {
		accessNanos = timestamps.Accessed.UnixNS
	}
	if timestamps.Modified != nil {
		modifiedNanos = timestamps.Modified.UnixNS
	}
	accessTimeval, expectedAccess := timevalAtMicrosecondPrecision(accessNanos)
	modifiedTimeval, expectedModified := timevalAtMicrosecondPrecision(modifiedNanos)
	if err := unix.Futimes(int(file.Fd()), []unix.Timeval{accessTimeval, modifiedTimeval}); err != nil {
		return failRequestedDarwinTimes(results, timestamps, fmt.Sprintf("apply access and modification times: %v", err))
	}

	verified, err := fstatDarwin(file)
	if err != nil {
		return failRequestedDarwinTimes(results, timestamps, fmt.Sprintf("verify destination timestamps: %v", err))
	}
	if timestamps.Modified != nil {
		results[1] = verifiedDarwinTimestamp(TimestampModified, verified.Mtim, expectedModified)
	}
	if timestamps.Accessed != nil {
		results[2] = verifiedDarwinTimestamp(TimestampAccessed, verified.Atim, expectedAccess)
	}
	return results
}

func fstatDarwin(file *os.File) (unix.Stat_t, error) {
	var stat unix.Stat_t
	if err := unix.Fstat(int(file.Fd()), &stat); err != nil {
		return unix.Stat_t{}, fmt.Errorf("inspect file descriptor: %w", err)
	}
	if stat.Mode&unix.S_IFMT != unix.S_IFREG {
		return unix.Stat_t{}, fmt.Errorf("opened path is not a regular file")
	}
	return stat, nil
}

func modelTimestampFromDarwin(value unix.Timespec) (*model.Timestamp, error) {
	nanos, err := timespecUnixNanoDarwin(value)
	if err != nil {
		return nil, err
	}
	instant := time.Unix(0, nanos).UTC()
	return &model.Timestamp{ISO: instant.Format(time.RFC3339Nano), UnixNS: nanos}, nil
}

func timespecUnixNanoDarwin(value unix.Timespec) (int64, error) {
	if value.Nsec < 0 || value.Nsec >= int64(time.Second) {
		return 0, fmt.Errorf("invalid nanosecond component %d", value.Nsec)
	}
	instant := time.Unix(value.Sec, value.Nsec).UTC()
	nanos := instant.UnixNano()
	if !time.Unix(0, nanos).Equal(instant) {
		return 0, fmt.Errorf("timestamp is outside signed Unix-nanosecond range")
	}
	return nanos, nil
}

func timespecAvailable(value unix.Timespec) bool {
	return value.Sec != 0 || value.Nsec != 0
}

func timevalAtMicrosecondPrecision(nanos int64) (unix.Timeval, int64) {
	effective := nanos - nanos%int64(time.Microsecond)
	seconds := effective / int64(time.Second)
	remainder := effective % int64(time.Second)
	if remainder < 0 {
		remainder += int64(time.Second)
		seconds--
	}
	return unix.Timeval{Sec: seconds, Usec: int32(remainder / int64(time.Microsecond))}, effective
}

func metadataUnavailable(kind TimestampKind) TimestampResult {
	return TimestampResult{Kind: kind, Status: TimestampUnavailable, Detail: "timestamp is unavailable in source metadata"}
}

func failRequestedDarwinTimes(results []TimestampResult, timestamps model.Timestamps, detail string) []TimestampResult {
	if timestamps.Modified != nil {
		results[1] = TimestampResult{Kind: TimestampModified, Status: TimestampFailed, Detail: detail}
	}
	if timestamps.Accessed != nil {
		results[2] = TimestampResult{Kind: TimestampAccessed, Status: TimestampFailed, Detail: detail}
	}
	return results
}

func verifiedDarwinTimestamp(kind TimestampKind, actual unix.Timespec, expected int64) TimestampResult {
	observed, err := timespecUnixNanoDarwin(actual)
	if err != nil {
		return TimestampResult{Kind: kind, Status: TimestampFailed, EffectivePrecision: darwinTimestampPrecision, Detail: fmt.Sprintf("read restored timestamp: %v", err)}
	}
	if observed != expected {
		return TimestampResult{Kind: kind, Status: TimestampFailed, EffectivePrecision: darwinTimestampPrecision, Detail: fmt.Sprintf("read back %d, want %d at microsecond precision", observed, expected)}
	}
	return TimestampResult{Kind: kind, Status: TimestampRestored, EffectivePrecision: darwinTimestampPrecision}
}
