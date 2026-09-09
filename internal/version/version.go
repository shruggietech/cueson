// Package version owns the Cueson executable version.
package version

// current is intentionally a variable so release tooling can replace it with
// -ldflags without introducing a second version source.
var current = "0.0.0"

// String returns the Cueson executable version.
func String() string {
	return current
}
