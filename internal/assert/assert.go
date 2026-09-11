package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {
	// Use the t.Helper() method to tell the test suite that this is a
	// helper function. The test suite will then ignore any lines in this
	// function when reporting errors, and instead show the line number of the
	// caller.
	t.Helper()

	if actual != expected {
		t.Errorf("got %v, want %v", actual, expected)
	}
}

func StringContains(t *testing.T, actual, expectedSubstring string) {
	t.Helper()
	if !strings.Contains(actual, expectedSubstring) {
		t.Errorf("got: %q; expected to contain: %q", actual, expectedSubstring)
	}
}
