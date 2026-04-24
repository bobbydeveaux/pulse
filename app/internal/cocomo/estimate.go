package cocomo

import "math"

// Basic COCOMO model coefficients for organic projects.
const (
	coeffA = 2.4    // effort multiplier
	coeffB = 1.05   // effort exponent
	coeffC = 2.5    // schedule multiplier
	coeffD = 0.38   // schedule exponent
	costPerMonth = 8000.0 // default developer cost per month (USD)
)

// Estimate holds COCOMO estimation results.
type Estimate struct {
	KSLOC        float64 // thousands of source lines of code
	PersonMonths float64 // estimated effort in person-months
	Persons      float64 // average number of developers
	Duration     float64 // estimated calendar months
	Cost         float64 // estimated cost in USD
}

// Calculate runs the Basic COCOMO model for organic projects.
//
// The Basic COCOMO model:
//   Effort = a * (KSLOC)^b person-months
//   Duration = c * (Effort)^d months
//   Persons = Effort / Duration
func Calculate(sloc int) Estimate {
	if sloc <= 0 {
		return Estimate{}
	}

	ksloc := float64(sloc) / 1000.0
	effort := coeffA * math.Pow(ksloc, coeffB)
	duration := coeffC * math.Pow(effort, coeffD)
	persons := effort / duration
	cost := effort * costPerMonth

	return Estimate{
		KSLOC:        ksloc,
		PersonMonths: math.Round(effort*10) / 10,
		Persons:      math.Round(persons*10) / 10,
		Duration:     math.Round(duration*10) / 10,
		Cost:         math.Round(cost),
	}
}
