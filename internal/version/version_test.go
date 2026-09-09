package version

import "testing"

func TestString(t *testing.T) {
	if got, want := String(), "0.0.0"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestStringUsesReleaseOverride(t *testing.T) {
	previous := releaseOverride
	releaseOverride = releaseMarkerPrefix + "1.2.3"
	t.Cleanup(func() { releaseOverride = previous })

	if got, want := String(), "1.2.3"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
