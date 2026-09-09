//go:build linux

package source

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/shruggietech/cueson/internal/model"
	"golang.org/x/sys/unix"
)

const linuxTimestampPrecision = "1ns"

func openRegularNoFollow(path string, _ bool) (*os.File, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0)
	if err != nil {
		return nil, fmt.Errorf("open regular file without following links: %w", err)
	}

	var stat unix.Stat_t
	if err := unix.Fstat(fd, &stat); err != nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("inspect opened file: %w", err)
	}
	if stat.Mode&unix.S_IFMT != unix.S_IFREG {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("opened path is not a regular file")
	}

	file := os.NewFile(uintptr(fd), path)
	if file == nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("adopt opened file descriptor")
	}
	return file, nil
}

func captureNativeTimestamps(file *os.File) (model.Timestamps, error) {
	fd := int(file.Fd())
	var fallback unix.Stat_t
	if err := unix.Fstat(fd, &fallback); err != nil {
		return model.Timestamps{}, fmt.Errorf("capture file timestamps: %w", err)
	}
	if fallback.Mode&unix.S_IFMT != unix.S_IFREG {
		return model.Timestamps{}, fmt.Errorf("capture file timestamps: opened object is not a regular file")
	}

	accessed, err := timestampFromLinuxTimespec(fallback.Atim)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture accessed timestamp: %w", err)
	}
	modified, err := timestampFromLinuxTimespec(fallback.Mtim)
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture modified timestamp: %w", err)
	}

	timestamps := model.Timestamps{
		Modified:      modified,
		Accessed:      accessed,
		CreatedSource: "unavailable",
	}

	var statx unix.Statx_t
	err = unix.Statx(fd, "", unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT, unix.STATX_BASIC_STATS|unix.STATX_BTIME, &statx)
	if errors.Is(err, unix.ENOSYS) || errors.Is(err, unix.EPERM) {
		return timestamps, nil
	}
	if err != nil {
		return model.Timestamps{}, fmt.Errorf("capture extended file timestamps: %w", err)
	}
	if statx.Mask&unix.STATX_ATIME != 0 {
		timestamps.Accessed, err = timestampFromStatx(statx.Atime)
		if err != nil {
			return model.Timestamps{}, fmt.Errorf("capture accessed timestamp: %w", err)
		}
	}
	if statx.Mask&unix.STATX_MTIME != 0 {
		timestamps.Modified, err = timestampFromStatx(statx.Mtime)
		if err != nil {
			return model.Timestamps{}, fmt.Errorf("capture modified timestamp: %w", err)
		}
	}
	if statx.Mask&unix.STATX_BTIME != 0 {
		timestamps.Created, err = timestampFromStatx(statx.Btime)
		if err != nil {
			return model.Timestamps{}, fmt.Errorf("capture birth timestamp: %w", err)
		}
		timestamps.CreatedSource = "birthtime"
	}
	return timestamps, nil
}

func applyNativeTimestamps(file *os.File, timestamps model.Timestamps) []TimestampResult {
	results := []TimestampResult{
		linuxCreationResult(timestamps.Created),
		unavailableLinuxResult(TimestampModified, timestamps.Modified),
		unavailableLinuxResult(TimestampAccessed, timestamps.Accessed),
	}
	if timestamps.Modified == nil && timestamps.Accessed == nil {
		return results
	}

	times := []unix.Timespec{
		{Nsec: unix.UTIME_OMIT},
		{Nsec: unix.UTIME_OMIT},
	}
	if timestamps.Accessed != nil {
		times[0] = unix.NsecToTimespec(timestamps.Accessed.UnixNS)
	}
	if timestamps.Modified != nil {
		times[1] = unix.NsecToTimespec(timestamps.Modified.UnixNS)
	}

	if err := setLinuxDescriptorTimes(file, times); err != nil {
		detail := fmt.Sprintf("apply access and modification timestamps: %v", err)
		markLinuxPresentFailed(results, timestamps, detail)
		return results
	}

	var stat unix.Stat_t
	if err := unix.Fstat(int(file.Fd()), &stat); err != nil {
		detail := fmt.Sprintf("read back access and modification timestamps: %v", err)
		markLinuxPresentFailed(results, timestamps, detail)
		return results
	}
	verifyLinuxTimestamp(&results[1], timestamps.Modified, stat.Mtim)
	verifyLinuxTimestamp(&results[2], timestamps.Accessed, stat.Atim)
	return results
}

func setLinuxDescriptorTimes(file *os.File, times []unix.Timespec) error {
	fd := int(file.Fd())
	err := unix.UtimesNanoAt(fd, "", times, unix.AT_EMPTY_PATH)
	if !errors.Is(err, unix.EINVAL) && !errors.Is(err, unix.ENOENT) {
		return err
	}

	// AT_EMPTY_PATH for utimensat requires Linux 5.8. The proc descriptor path
	// still identifies the already-open file on older kernels without exposing
	// timestamp mutation to a caller-controlled pathname.
	return unix.UtimesNano("/proc/self/fd/"+strconv.Itoa(fd), times)
}

func linuxCreationResult(timestamp *model.Timestamp) TimestampResult {
	if timestamp == nil {
		return TimestampResult{Kind: TimestampCreated, Status: TimestampUnavailable, Detail: "source timestamp is unavailable"}
	}
	return TimestampResult{Kind: TimestampCreated, Status: TimestampUnsupported, Detail: "Linux does not support setting filesystem creation or birth time"}
}

func unavailableLinuxResult(kind TimestampKind, timestamp *model.Timestamp) TimestampResult {
	if timestamp == nil {
		return TimestampResult{Kind: kind, Status: TimestampUnavailable, Detail: "source timestamp is unavailable"}
	}
	return TimestampResult{Kind: kind, EffectivePrecision: linuxTimestampPrecision}
}

func markLinuxPresentFailed(results []TimestampResult, timestamps model.Timestamps, detail string) {
	if timestamps.Modified != nil {
		results[1].Status = TimestampFailed
		results[1].Detail = detail
	}
	if timestamps.Accessed != nil {
		results[2].Status = TimestampFailed
		results[2].Detail = detail
	}
}

func verifyLinuxTimestamp(result *TimestampResult, requested *model.Timestamp, actual unix.Timespec) {
	if requested == nil {
		return
	}
	actualNS, err := linuxUnixNS(int64(actual.Sec), int64(actual.Nsec))
	if err != nil {
		result.Status = TimestampFailed
		result.Detail = fmt.Sprintf("readback timestamp is not representable as unix_ns: %v", err)
		return
	}
	if actualNS != requested.UnixNS {
		result.Status = TimestampFailed
		result.Detail = fmt.Sprintf("readback timestamp is %d, want %d", actualNS, requested.UnixNS)
		return
	}
	result.Status = TimestampRestored
	result.Detail = "verified by descriptor readback"
}

func timestampFromLinuxTimespec(timestamp unix.Timespec) (*model.Timestamp, error) {
	return timestampFromLinuxParts(int64(timestamp.Sec), int64(timestamp.Nsec))
}

func timestampFromStatx(timestamp unix.StatxTimestamp) (*model.Timestamp, error) {
	return timestampFromLinuxParts(timestamp.Sec, int64(timestamp.Nsec))
}

func timestampFromLinuxParts(seconds, nanoseconds int64) (*model.Timestamp, error) {
	unixNS, err := linuxUnixNS(seconds, nanoseconds)
	if err != nil {
		return nil, err
	}
	instant := time.Unix(seconds, nanoseconds).UTC()
	return &model.Timestamp{ISO: instant.Format(time.RFC3339Nano), UnixNS: unixNS}, nil
}

func linuxUnixNS(seconds, nanoseconds int64) (int64, error) {
	if nanoseconds < 0 || nanoseconds >= int64(time.Second) {
		return 0, fmt.Errorf("nanoseconds %d are outside [0, 999999999]", nanoseconds)
	}

	const (
		maximumSeconds     = int64(math.MaxInt64 / int64(time.Second))
		maximumNanoseconds = int64(math.MaxInt64 % int64(time.Second))
		minimumSeconds     = int64(math.MinInt64/int64(time.Second) - 1)
		minimumNanoseconds = int64(time.Second) + int64(math.MinInt64%int64(time.Second))
	)
	if seconds < minimumSeconds || seconds > maximumSeconds ||
		(seconds == minimumSeconds && nanoseconds < minimumNanoseconds) ||
		(seconds == maximumSeconds && nanoseconds > maximumNanoseconds) {
		return 0, fmt.Errorf("timestamp %d.%09d is outside signed unix_ns range", seconds, nanoseconds)
	}
	if seconds == minimumSeconds {
		return math.MinInt64 + (nanoseconds - minimumNanoseconds), nil
	}
	return seconds*int64(time.Second) + nanoseconds, nil
}
