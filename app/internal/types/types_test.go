package types

import "testing"

func TestMaintainabilityGrade(t *testing.T) {
	tests := []struct {
		mi   float64
		want string
	}{
		{100, "A"},
		{80, "A"},
		{79.9, "B"},
		{60, "B"},
		{59.9, "C"},
		{40, "C"},
		{39.9, "D"},
		{20, "D"},
		{19.9, "F"},
		{0, "F"},
	}

	for _, tt := range tests {
		got := MaintainabilityGrade(tt.mi)
		if got != tt.want {
			t.Errorf("MaintainabilityGrade(%.1f) = %s, want %s", tt.mi, got, tt.want)
		}
	}
}

func TestCyclomaticGrade(t *testing.T) {
	tests := []struct {
		ccn  int
		want string
	}{
		{1, "A"},
		{5, "A"},
		{6, "B"},
		{10, "B"},
		{11, "C"},
		{20, "C"},
		{21, "D"},
		{30, "D"},
		{31, "F"},
		{50, "F"},
	}

	for _, tt := range tests {
		got := CyclomaticGrade(tt.ccn)
		if got != tt.want {
			t.Errorf("CyclomaticGrade(%d) = %s, want %s", tt.ccn, got, tt.want)
		}
	}
}
