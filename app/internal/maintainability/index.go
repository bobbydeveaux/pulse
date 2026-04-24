package maintainability

import (
	"math"
)

// Index computes the Maintainability Index using the Microsoft Visual Studio formula:
// MI = max(0, (171 - 5.2*ln(HV) - 0.23*CC - 16.2*ln(SLOC)) * 100 / 171)
//
// Where:
//   - HV = Halstead Volume
//   - CC = average cyclomatic complexity
//   - SLOC = source lines of code
func Index(sloc int, avgCyclomatic float64, halsteadVolume float64) float64 {
	if sloc <= 0 {
		return 100.0
	}

	lnSLOC := math.Log(float64(sloc))
	lnHV := math.Log(math.Max(halsteadVolume, 1.0))

	mi := 171.0 - 5.2*lnHV - 0.23*avgCyclomatic - 16.2*lnSLOC
	mi = mi * 100.0 / 171.0

	return math.Max(0, math.Min(100, mi))
}

// HalsteadVolume provides a simplified Halstead volume estimate.
// Approximation: V ≈ SLOC * log2(SLOC), which matches empirical
// observations for typical function-sized code blocks (5-50 SLOC).
func HalsteadVolume(sloc int) float64 {
	if sloc <= 0 {
		return 0
	}
	s := float64(sloc)
	return s * math.Log2(math.Max(s, 2.0))
}

// FileIndex computes the maintainability index for a file by averaging
// per-function MI scores. This is more accurate than computing MI on
// total file SLOC, since the formula was designed for function-level metrics.
func FileIndex(functions []FuncInput, fileSLOC int) float64 {
	if len(functions) == 0 {
		// No functions: compute directly on file SLOC with CCN=1
		hv := HalsteadVolume(fileSLOC)
		return Index(fileSLOC, 1, hv)
	}

	total := 0.0
	for _, f := range functions {
		hv := HalsteadVolume(f.SLOC)
		mi := Index(f.SLOC, float64(f.Cyclomatic), hv)
		total += mi
	}
	return total / float64(len(functions))
}

// FuncInput holds the minimum data needed to compute MI for a function.
type FuncInput struct {
	SLOC       int
	Cyclomatic int
}
