package report

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/fatih/color"

	"github.com/bobbydeveaux/pulse/internal/types"
)

var (
	cyan    = color.New(color.FgCyan, color.Bold)
	green   = color.New(color.FgGreen)
	yellow  = color.New(color.FgYellow)
	red     = color.New(color.FgRed)
	bold    = color.New(color.Bold)
	muted   = color.New(color.FgHiBlack)
)

// PrintCheck outputs the full analysis report.
func PrintCheck(pm *types.ProjectMetrics) {
	fmt.Println()
	cyan.Println("  Pulse — Code Quality Report")
	muted.Println("  " + strings.Repeat("─", 50))
	fmt.Println()

	// Language summary table
	bold.Println("  Language         Files      SLOC  Comments    Blanks")
	muted.Println("  " + strings.Repeat("─", 54))
	for _, ls := range pm.LanguageSummaries {
		fmt.Printf("  %-16s %5d  %8s  %8s  %8s\n",
			ls.Language, ls.Files,
			fmtNum(ls.SLOC), fmtNum(ls.Comments), fmtNum(ls.Blanks))
	}
	muted.Println("  " + strings.Repeat("─", 54))
	bold.Printf("  %-16s %5d  %8s  %8s  %8s\n",
		"Total", pm.TotalFiles,
		fmtNum(pm.TotalSLOC), fmtNum(pm.TotalComments), fmtNum(pm.TotalBlanks))
	fmt.Println()

	// Top complexity hotspots
	if pm.TotalFunctions > 0 {
		cyan.Println("  Top complexity hotspots:")
		fmt.Println()

		// Collect and sort all functions by cyclomatic complexity
		var allFuncs []types.FunctionMetrics
		for _, f := range pm.Files {
			allFuncs = append(allFuncs, f.Functions...)
		}
		sort.Slice(allFuncs, func(i, j int) bool {
			return allFuncs[i].Cyclomatic > allFuncs[j].Cyclomatic
		})

		limit := 10
		if len(allFuncs) < limit {
			limit = len(allFuncs)
		}

		for _, fn := range allFuncs[:limit] {
			name := shortPath(fn.File) + ":" + fn.Name
			if len(name) > 42 {
				name = "..." + name[len(name)-39:]
			}

			ccnColor := gradeColor(types.CyclomaticGrade(fn.Cyclomatic))
			cogColor := gradeColor(types.CyclomaticGrade(fn.Cognitive))
			gradeC := gradeColor(fn.Grade)

			fmt.Printf("  %-42s ", name)
			muted.Print("CCN:")
			ccnColor.Printf(" %-4d", fn.Cyclomatic)
			muted.Print("  Cognitive:")
			cogColor.Printf(" %-4d", fn.Cognitive)
			muted.Print("  Grade:")
			gradeC.Printf(" %s\n", fn.Grade)
		}
		fmt.Println()
	}

	// Summary metrics
	muted.Println("  " + strings.Repeat("─", 50))

	if pm.TotalFunctions > 0 {
		fmt.Print("  Avg Cyclomatic:     ")
		gradeColor(types.CyclomaticGrade(int(pm.AvgCyclomatic))).
			Printf("%.1f\n", pm.AvgCyclomatic)
		fmt.Print("  Avg Cognitive:      ")
		gradeColor(types.CyclomaticGrade(int(pm.AvgCognitive))).
			Printf("%.1f\n", pm.AvgCognitive)
	}

	fmt.Print("  Duplication:        ")
	if pm.DuplicationPct < 3 {
		green.Printf("%.1f%%", pm.DuplicationPct)
	} else if pm.DuplicationPct < 10 {
		yellow.Printf("%.1f%%", pm.DuplicationPct)
	} else {
		red.Printf("%.1f%%", pm.DuplicationPct)
	}
	muted.Printf(" (%d clone(s) detected)\n", len(pm.Duplicates))

	fmt.Print("  Maintainability:    ")
	gradeColor(pm.Grade).Printf("%.1f (%s)\n", pm.MaintainabilityIdx, pm.Grade)

	fmt.Print("  Estimated effort:   ")
	bold.Printf("%.1f person-months", pm.COCOMOMonths)
	muted.Printf(" (COCOMO, $%.0f)\n", pm.COCOMOCost)

	fmt.Println()
}

// PrintTrend outputs time-series trend data.
func PrintTrend(points []types.TrendPoint) {
	fmt.Println()
	cyan.Println("  Pulse — Complexity Trend")
	muted.Println("  " + strings.Repeat("─", 50))
	fmt.Println()

	if len(points) == 0 {
		muted.Println("  No trend data available.")
		return
	}

	// ASCII chart of avg cyclomatic
	maxCCN := 0.0
	for _, p := range points {
		if p.AvgCyclomatic > maxCCN {
			maxCCN = p.AvgCyclomatic
		}
	}

	bold.Println("  Avg Cyclomatic Complexity")
	fmt.Print("  ")
	barChars := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	for _, p := range points {
		ratio := p.AvgCyclomatic / math.Max(maxCCN, 1)
		idx := int(ratio * float64(len(barChars)-1))
		if idx >= len(barChars) {
			idx = len(barChars) - 1
		}
		if idx < 0 {
			idx = 0
		}
		// Color based on value
		c := green
		if p.AvgCyclomatic > 15 {
			c = red
		} else if p.AvgCyclomatic > 10 {
			c = yellow
		}
		c.Print(barChars[idx])
	}
	fmt.Println()
	fmt.Println()

	// Table
	bold.Println("  Hash     Date         CCN    Cog   SLOC    Dup%   MI   Grade")
	muted.Println("  " + strings.Repeat("─", 66))
	for _, p := range points {
		date := p.Date
		if len(date) > 10 {
			date = date[:10]
		}
		gc := gradeColor(p.Grade)
		fmt.Printf("  %s  %s  %5.1f  %5.1f  %5d  %5.1f%%  %4.0f  ",
			p.ShortHash, date, p.AvgCyclomatic, p.AvgCognitive,
			p.TotalSLOC, p.DuplicationPct, p.Maintainability)
		gc.Println(p.Grade)
	}
	fmt.Println()

	// Delta summary
	if len(points) >= 2 {
		first := points[0]
		last := points[len(points)-1]
		delta := last.AvgCyclomatic - first.AvgCyclomatic

		if delta < 0 {
			green.Printf("  ↓ Complexity decreased by %.1f", -delta)
		} else if delta > 0 {
			red.Printf("  ↑ Complexity increased by %.1f", delta)
		} else {
			muted.Print("  → Complexity unchanged")
		}
		muted.Printf(" over %d commits\n", len(points))
		fmt.Println()
	}
}

// PrintDiff outputs a comparison between two analysis snapshots.
func PrintDiff(diff *types.DiffResult) {
	fmt.Println()
	cyan.Println("  Pulse — Complexity Diff")
	muted.Println("  " + strings.Repeat("─", 50))
	fmt.Println()

	printDelta("SLOC", float64(diff.Base.TotalSLOC), float64(diff.Head.TotalSLOC), false)
	printDelta("Avg Cyclomatic", diff.Base.AvgCyclomatic, diff.Head.AvgCyclomatic, true)
	printDelta("Avg Cognitive", diff.Base.AvgCognitive, diff.Head.AvgCognitive, true)
	printDelta("Duplication", diff.Base.DuplicationPct, diff.Head.DuplicationPct, true)
	printDelta("Maintainability", diff.Base.MaintainabilityIdx, diff.Head.MaintainabilityIdx, false)
	fmt.Println()

	if len(diff.NewHotspots) > 0 {
		yellow.Println("  New complexity hotspots:")
		for _, fn := range diff.NewHotspots {
			fmt.Printf("    %s:%s  CCN: %d  Cognitive: %d  Grade: %s\n",
				shortPath(fn.File), fn.Name, fn.Cyclomatic, fn.Cognitive, fn.Grade)
		}
		fmt.Println()
	}
}

// PrintGate outputs quality gate results.
func PrintGate(result *types.GateResult) {
	fmt.Println()
	if result.Passed {
		green.Println("  ✓ All quality gates passed")
	} else {
		red.Println("  ✗ Quality gate failures:")
		for _, f := range result.Failures {
			red.Printf("    • %s\n", f)
		}
	}
	fmt.Println()
}

func printDelta(label string, base, head float64, higherIsBad bool) {
	delta := head - base
	sign := "+"
	c := yellow
	if delta < 0 {
		sign = ""
		if higherIsBad {
			c = green
		} else {
			c = red
		}
	} else if delta > 0 {
		if higherIsBad {
			c = red
		} else {
			c = green
		}
	} else {
		c = muted
		sign = "±"
	}

	fmt.Printf("  %-22s %8.1f → %-8.1f  ", label, base, head)
	c.Printf("(%s%.1f)\n", sign, delta)
}

func gradeColor(grade string) *color.Color {
	switch grade {
	case "A":
		return green
	case "B":
		return color.New(color.FgGreen)
	case "C":
		return yellow
	case "D":
		return color.New(color.FgRed)
	default:
		return red
	}
}

func shortPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) <= 2 {
		return path
	}
	return strings.Join(parts[len(parts)-2:], "/")
}

func fmtNum(n int) string {
	s := fmt.Sprintf("%d", n)
	if n < 1000 {
		return s
	}
	// Add thousand separators
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}
