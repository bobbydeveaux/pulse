package cocomo

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		sloc      int
		wantMonths float64
		wantGt    float64
	}{
		{"zero", 0, 0, -1},
		{"tiny project (1K)", 1000, 0, 0},
		{"small project (5K)", 5000, 0, 5},
		{"medium project (50K)", 50000, 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			est := Calculate(tt.sloc)

			if tt.sloc == 0 {
				if est.PersonMonths != 0 {
					t.Errorf("expected 0 person-months for 0 SLOC, got %.1f", est.PersonMonths)
				}
				return
			}

			if est.PersonMonths <= tt.wantGt {
				t.Errorf("PersonMonths = %.1f, want > %.0f", est.PersonMonths, tt.wantGt)
			}
			if est.Persons <= 0 {
				t.Error("expected Persons > 0")
			}
			if est.Duration <= 0 {
				t.Error("expected Duration > 0")
			}
			if est.Cost <= 0 {
				t.Error("expected Cost > 0")
			}
			if est.KSLOC != float64(tt.sloc)/1000.0 {
				t.Errorf("KSLOC = %.1f, want %.1f", est.KSLOC, float64(tt.sloc)/1000.0)
			}

			t.Logf("SLOC=%d: %.1f person-months, %.1f people, %.1f months, $%.0f",
				tt.sloc, est.PersonMonths, est.Persons, est.Duration, est.Cost)
		})
	}
}

func TestCalculate_Scaling(t *testing.T) {
	// Larger codebases should require more effort
	small := Calculate(1000)
	medium := Calculate(10000)
	large := Calculate(100000)

	if medium.PersonMonths <= small.PersonMonths {
		t.Errorf("10K SLOC effort (%.1f) should be > 1K SLOC (%.1f)",
			medium.PersonMonths, small.PersonMonths)
	}
	if large.PersonMonths <= medium.PersonMonths {
		t.Errorf("100K SLOC effort (%.1f) should be > 10K SLOC (%.1f)",
			large.PersonMonths, medium.PersonMonths)
	}
}
