package maintainability

import (
	"testing"
)

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		sloc     int
		ccn      float64
		hv       float64
		wantMin  float64
		wantMax  float64
	}{
		{"empty file", 0, 0, 0, 100, 100},
		{"tiny simple func", 5, 1, HalsteadVolume(5), 60, 100},
		{"small func low complexity", 20, 2, HalsteadVolume(20), 40, 80},
		{"medium func moderate complexity", 50, 8, HalsteadVolume(50), 20, 60},
		{"large func high complexity", 200, 25, HalsteadVolume(200), 0, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mi := Index(tt.sloc, tt.ccn, tt.hv)
			if mi < tt.wantMin || mi > tt.wantMax {
				t.Errorf("MI = %.1f, want between %.0f and %.0f", mi, tt.wantMin, tt.wantMax)
			}
			t.Logf("MI = %.1f (SLOC=%d, CCN=%.0f)", mi, tt.sloc, tt.ccn)
		})
	}
}

func TestHalsteadVolume(t *testing.T) {
	tests := []struct {
		sloc    int
		wantGt  float64
	}{
		{0, -1},
		{1, 0},
		{10, 10},
		{100, 100},
	}

	for _, tt := range tests {
		hv := HalsteadVolume(tt.sloc)
		if hv <= tt.wantGt {
			t.Errorf("HalsteadVolume(%d) = %.1f, want > %.0f", tt.sloc, hv, tt.wantGt)
		}
		t.Logf("HalsteadVolume(%d) = %.1f", tt.sloc, hv)
	}
}

func TestFileIndex(t *testing.T) {
	// File with simple functions should score well
	simple := []FuncInput{
		{SLOC: 5, Cyclomatic: 1},
		{SLOC: 10, Cyclomatic: 2},
		{SLOC: 8, Cyclomatic: 1},
	}
	mi := FileIndex(simple, 30)
	if mi < 50 {
		t.Errorf("simple functions MI = %.1f, expected >= 50", mi)
	}
	t.Logf("Simple functions file MI = %.1f", mi)

	// File with complex functions should score poorly
	complex := []FuncInput{
		{SLOC: 100, Cyclomatic: 20},
		{SLOC: 80, Cyclomatic: 15},
	}
	miComplex := FileIndex(complex, 200)
	if miComplex > mi {
		t.Errorf("complex file MI (%.1f) should be lower than simple (%.1f)", miComplex, mi)
	}
	t.Logf("Complex functions file MI = %.1f", miComplex)
}
