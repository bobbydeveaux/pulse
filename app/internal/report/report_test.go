package report

import (
	"testing"
)

func TestShortPath(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"main.go", "main.go"},
		{"pkg/main.go", "pkg/main.go"},
		{"internal/foo/main.go", "foo/main.go"},
		{"a/b/c/d/main.go", "d/main.go"},
		{"", ""},
	}
	for _, c := range cases {
		if got := shortPath(c.in); got != c.want {
			t.Errorf("shortPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFmtNum(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{0, "0"},
		{42, "42"},
		{999, "999"},
		{1000, "1,000"},
		{12345, "12,345"},
		{1000000, "1,000,000"},
	}
	for _, c := range cases {
		if got := fmtNum(c.in); got != c.want {
			t.Errorf("fmtNum(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGradeColor(t *testing.T) {
	// gradeColor returns *color.Color — we don't assert on the colour itself
	// (it's stateful and varies with TTY), just that we get a non-nil value
	// for every input grade including the default.
	cases := []string{"A", "B", "C", "D", "F", "", "Z"}
	for _, g := range cases {
		if got := gradeColor(g); got == nil {
			t.Errorf("gradeColor(%q) returned nil", g)
		}
	}
}

func TestPrintDelta_AllBranches(t *testing.T) {
	// Calling printDelta writes to stdout; the test purpose here is purely
	// coverage of the sign/colour decision logic. We exercise all 6 paths:
	// (delta>0 + higherIsBad), (delta>0 + !higherIsBad),
	// (delta<0 + higherIsBad), (delta<0 + !higherIsBad),
	// (delta==0), and verify the function does not panic.
	cases := []struct {
		base, head    float64
		higherIsBad   bool
	}{
		{10, 20, true},
		{10, 20, false},
		{20, 10, true},
		{20, 10, false},
		{10, 10, true},
		{10, 10, false},
	}
	for _, c := range cases {
		// Should not panic; output is discarded by the test runner unless -v.
		printDelta("test", c.base, c.head, c.higherIsBad)
	}
}
