package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/bobbydeveaux/pulse/internal/analyzer"
	"github.com/bobbydeveaux/pulse/internal/git"
	"github.com/bobbydeveaux/pulse/internal/report"
	"github.com/bobbydeveaux/pulse/internal/types"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "pulse",
		Short: "Code quality & complexity analyzer",
		Long:  "Pulse — take the pulse of your codebase.\nAnalyze complexity, duplication, maintainability, and track trends over time.",
	}

	rootCmd.AddCommand(checkCmd())
	rootCmd.AddCommand(trendCmd())
	rootCmd.AddCommand(diffCmd())
	rootCmd.AddCommand(gateCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func checkCmd() *cobra.Command {
	var noColor bool
	var skipDirs []string

	cmd := &cobra.Command{
		Use:   "check [path]",
		Short: "Analyze code quality for a directory",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			path, _ = filepath.Abs(path)

			if noColor {
				color.NoColor = true
			}

			analyzer.ExtraSkipDirs = skipDirs

			pm, err := analyzer.Analyze(path)
			if err != nil {
				return fmt.Errorf("analysis failed: %w", err)
			}

			report.PrintCheck(pm)
			return nil
		},
	}

	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	cmd.Flags().StringSliceVar(&skipDirs, "skip", nil, "Additional directories to skip (comma-separated)")
	return cmd
}

func trendCmd() *cobra.Command {
	var last int
	var noColor bool

	cmd := &cobra.Command{
		Use:   "trend [path]",
		Short: "Show complexity trends across git history",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			path, _ = filepath.Abs(path)

			if noColor {
				color.NoColor = true
			}

			// Get commits
			commits, err := git.Log(last, path)
			if err != nil {
				return fmt.Errorf("failed to read git log: %w", err)
			}
			if len(commits) == 0 {
				fmt.Println("No commits found.")
				return nil
			}

			// Save current state
			origBranch, err := git.CurrentBranch(path)
			if err != nil {
				return fmt.Errorf("failed to get current branch: %w", err)
			}

			dirty, _ := git.IsDirty(path)
			if dirty {
				if err := git.Stash(path); err != nil {
					return fmt.Errorf("failed to stash changes: %w", err)
				}
				defer git.StashPop(path)
			}
			defer git.Checkout(origBranch, path)

			// Analyze each commit (oldest first)
			var points []types.TrendPoint
			for i := len(commits) - 1; i >= 0; i-- {
				c := commits[i]
				fmt.Fprintf(os.Stderr, "\r  Analyzing commit %d/%d (%s)...",
					len(commits)-i, len(commits), c.ShortHash)

				if err := git.Checkout(c.Hash, path); err != nil {
					continue
				}

				pm, err := analyzer.Analyze(path)
				if err != nil {
					continue
				}

				points = append(points, types.TrendPoint{
					Hash:            c.Hash,
					ShortHash:       c.ShortHash,
					Author:          c.Author,
					Date:            c.Date,
					Message:         c.Message,
					AvgCyclomatic:   pm.AvgCyclomatic,
					AvgCognitive:    pm.AvgCognitive,
					TotalSLOC:       pm.TotalSLOC,
					DuplicationPct:  pm.DuplicationPct,
					Maintainability: pm.MaintainabilityIdx,
					Grade:           pm.Grade,
				})
			}
			fmt.Fprintf(os.Stderr, "\r%s\r", "                                                    ")

			report.PrintTrend(points)
			return nil
		},
	}

	cmd.Flags().IntVar(&last, "last", 20, "Number of commits to analyze")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	return cmd
}

func diffCmd() *cobra.Command {
	var base string
	var noColor bool

	cmd := &cobra.Command{
		Use:   "diff [path]",
		Short: "Compare code quality between two points",
		Long:  "Compare quality metrics between HEAD and a base ref (default: merge-base with main).",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			path, _ = filepath.Abs(path)

			if noColor {
				color.NoColor = true
			}

			// Analyze HEAD
			headMetrics, err := analyzer.Analyze(path)
			if err != nil {
				return fmt.Errorf("failed to analyze HEAD: %w", err)
			}

			// Determine base ref
			if base == "" {
				base = "main"
			}
			baseRef, err := git.MergeBase(base, "HEAD", path)
			if err != nil {
				// Fallback: use HEAD~1
				baseRef = "HEAD~1"
			}

			// Save state and checkout base
			origBranch, _ := git.CurrentBranch(path)
			dirty, _ := git.IsDirty(path)
			if dirty {
				if err := git.Stash(path); err != nil {
					return fmt.Errorf("failed to stash: %w", err)
				}
				defer git.StashPop(path)
			}

			if err := git.Checkout(baseRef, path); err != nil {
				return fmt.Errorf("failed to checkout base: %w", err)
			}

			baseMetrics, err := analyzer.Analyze(path)
			if err != nil {
				git.Checkout(origBranch, path)
				return fmt.Errorf("failed to analyze base: %w", err)
			}

			// Restore original
			git.Checkout(origBranch, path)

			// Compute diff
			diff := computeDiff(baseMetrics, headMetrics)
			report.PrintDiff(diff)
			return nil
		},
	}

	cmd.Flags().StringVar(&base, "base", "", "Base ref to compare against (default: main)")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	return cmd
}

func gateCmd() *cobra.Command {
	var (
		maxCCN        int
		maxCognitive  int
		maxDup        float64
		minMI         float64
		noColor       bool
	)

	cmd := &cobra.Command{
		Use:   "gate [path]",
		Short: "Check code against quality gates",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			path, _ = filepath.Abs(path)

			if noColor {
				color.NoColor = true
			}

			pm, err := analyzer.Analyze(path)
			if err != nil {
				return fmt.Errorf("analysis failed: %w", err)
			}

			result := checkGates(pm, types.GateConfig{
				MaxCyclomatic:      maxCCN,
				MaxCognitive:       maxCognitive,
				MaxDuplication:     maxDup,
				MinMaintainability: minMI,
			})

			// Print brief summary first
			report.PrintCheck(pm)
			report.PrintGate(result)

			if !result.Passed {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().IntVar(&maxCCN, "max-ccn", 0, "Maximum cyclomatic complexity per function (0 = no limit)")
	cmd.Flags().IntVar(&maxCognitive, "max-cognitive", 0, "Maximum cognitive complexity per function (0 = no limit)")
	cmd.Flags().Float64Var(&maxDup, "max-duplication", 0, "Maximum duplication percentage (0 = no limit)")
	cmd.Flags().Float64Var(&minMI, "min-maintainability", 0, "Minimum maintainability index (0 = no limit)")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	return cmd
}

func checkGates(pm *types.ProjectMetrics, cfg types.GateConfig) *types.GateResult {
	result := &types.GateResult{Passed: true}

	// Check per-function cyclomatic complexity
	if cfg.MaxCyclomatic > 0 {
		for _, f := range pm.Files {
			for _, fn := range f.Functions {
				if fn.Cyclomatic > cfg.MaxCyclomatic {
					result.Passed = false
					result.Failures = append(result.Failures,
						fmt.Sprintf("CCN %d exceeds max %d: %s:%s (line %d)",
							fn.Cyclomatic, cfg.MaxCyclomatic,
							shortPath(fn.File), fn.Name, fn.StartLine))
				}
			}
		}
	}

	// Check per-function cognitive complexity
	if cfg.MaxCognitive > 0 {
		for _, f := range pm.Files {
			for _, fn := range f.Functions {
				if fn.Cognitive > cfg.MaxCognitive {
					result.Passed = false
					result.Failures = append(result.Failures,
						fmt.Sprintf("Cognitive %d exceeds max %d: %s:%s (line %d)",
							fn.Cognitive, cfg.MaxCognitive,
							shortPath(fn.File), fn.Name, fn.StartLine))
				}
			}
		}
	}

	// Check duplication
	if cfg.MaxDuplication > 0 && pm.DuplicationPct > cfg.MaxDuplication {
		result.Passed = false
		result.Failures = append(result.Failures,
			fmt.Sprintf("Duplication %.1f%% exceeds max %.1f%%",
				pm.DuplicationPct, cfg.MaxDuplication))
	}

	// Check maintainability
	if cfg.MinMaintainability > 0 && pm.MaintainabilityIdx < cfg.MinMaintainability {
		result.Passed = false
		result.Failures = append(result.Failures,
			fmt.Sprintf("Maintainability %.1f below min %.1f",
				pm.MaintainabilityIdx, cfg.MinMaintainability))
	}

	return result
}

func computeDiff(base, head *types.ProjectMetrics) *types.DiffResult {
	diff := &types.DiffResult{
		Base:                 *base,
		Head:                 *head,
		DeltaCyclomatic:      head.AvgCyclomatic - base.AvgCyclomatic,
		DeltaCognitive:       head.AvgCognitive - base.AvgCognitive,
		DeltaSLOC:            head.TotalSLOC - base.TotalSLOC,
		DeltaDuplication:     head.DuplicationPct - base.DuplicationPct,
		DeltaMaintainability: head.MaintainabilityIdx - base.MaintainabilityIdx,
	}

	// Find new hotspots (functions in head with CCN > 10 that didn't exist in base)
	baseFuncs := map[string]int{}
	for _, f := range base.Files {
		for _, fn := range f.Functions {
			baseFuncs[fn.File+":"+fn.Name] = fn.Cyclomatic
		}
	}

	for _, f := range head.Files {
		for _, fn := range f.Functions {
			key := fn.File + ":" + fn.Name
			baseCCN, existed := baseFuncs[key]
			if fn.Cyclomatic > 10 && (!existed || fn.Cyclomatic > baseCCN+5) {
				diff.NewHotspots = append(diff.NewHotspots, fn)
			}
		}
	}

	sort.Slice(diff.NewHotspots, func(i, j int) bool {
		return diff.NewHotspots[i].Cyclomatic > diff.NewHotspots[j].Cyclomatic
	})

	if len(diff.NewHotspots) > 5 {
		diff.NewHotspots = diff.NewHotspots[:5]
	}

	return diff
}

func shortPath(path string) string {
	parts := filepath.SplitList(path)
	if len(parts) == 0 {
		return path
	}
	base := filepath.Base(path)
	dir := filepath.Base(filepath.Dir(path))
	if dir == "." || dir == "/" {
		return base
	}
	return dir + "/" + base
}

