// Package version owns the Cueson executable version.
package version

import "strings"

const releaseMarkerPrefix = "cueson-release-version:"

// current is the current source-build contract version.
var current = "1.1.0"

// releaseOverride remains empty in source builds. Release tooling injects
// the version with a distinctive prefix so artifact verification can prove the
// public version of foreign-target binaries without executing them.
var releaseOverride string

// String returns the Cueson executable version.
func String() string {
	if strings.HasPrefix(releaseOverride, releaseMarkerPrefix) {
		return strings.TrimPrefix(releaseOverride, releaseMarkerPrefix)
	}
	return current
}
