package version

import "testing"

func TestString(t *testing.T) {
	t.Parallel()

	if got, want := String(), "0.0.0"; got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}
