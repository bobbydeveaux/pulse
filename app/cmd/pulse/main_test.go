package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bobbydeveaux/pulse/internal/analyzer"
	"github.com/bobbydeveaux/pulse/internal/types"
)

// ---------------------------------------------------------------------------
// shortPath
// ---------------------------------------------------------------------------

func TestShortPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"basename only", "main.go", "main.go"},
		{"two-component path", "cmd/main.go", "cmd/main.go"},
		{"three-component path takes dir+base", "/a/b/c.go", "b/c.go"},
		{"dot dir collapses to basename", "./main.go", "main.go"},
		{"deep path keeps last two components", "/x/y/z/main.go", "z/main.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shortPath(tt.in)
			if got != tt.want {
				t.Errorf("shortPath(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// checkGates
// ---------------------------------------------------------------------------

func TestCheckGates_AllZeroThresholdsAlwaysPass(t *testing.T) {
	pm := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Path: "x.go",
			Functions: []types.FunctionMetrics{
				{Name: "BigFn", File: "x.go", Cyclomatic: 99, Cognitive: 99, StartLine: 1},
			},
		}},
		DuplicationPct:     50.0,
		MaintainabilityIdx: 1.0,
	}
	res := checkGates(pm, types.GateConfig{})
	if !res.Passed {
		t.Fatalf("expected pass with zero thresholds, got failures: %v", res.Failures)
	}
	if len(res.Failures) != 0 {
		t.Errorf("expected no failures, got %v", res.Failures)
	}
}

func TestCheckGates_CyclomaticFail(t *testing.T) {
	pm := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Path: "x.go",
			Functions: []types.FunctionMetrics{
				{Name: "Hot", File: "src/x.go", Cyclomatic: 25, Cognitive: 5, StartLine: 12},
				{Name: "Cool", File: "src/x.go", Cyclomatic: 3, Cognitive: 2, StartLine: 42},
			},
		}},
	}
	res := checkGates(pm, types.GateConfig{MaxCyclomatic: 10})
	if res.Passed {
		t.Fatal("expected gate to fail when CCN exceeds threshold")
	}
	if len(res.Failures) != 1 {
		t.Fatalf("expected 1 failure, got %d: %v", len(res.Failures), res.Failures)
	}
	if !strings.Contains(res.Failures[0], "CCN 25 exceeds max 10") {
		t.Errorf("failure message missing CCN detail: %q", res.Failures[0])
	}
	if !strings.Contains(res.Failures[0], "Hot") {
		t.Errorf("failure should name the function: %q", res.Failures[0])
	}
}

func TestCheckGates_CognitiveFail(t *testing.T) {
	pm := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "Heavy", File: "f.go", Cognitive: 40, StartLine: 7},
			},
		}},
	}
	res := checkGates(pm, types.GateConfig{MaxCognitive: 15})
	if res.Passed {
		t.Fatal("expected gate to fail on cognitive")
	}
	if len(res.Failures) != 1 || !strings.Contains(res.Failures[0], "Cognitive 40 exceeds max 15") {
		t.Errorf("unexpected failure list: %v", res.Failures)
	}
}

func TestCheckGates_DuplicationFail(t *testing.T) {
	pm := &types.ProjectMetrics{DuplicationPct: 12.5}
	res := checkGates(pm, types.GateConfig{MaxDuplication: 5.0})
	if res.Passed {
		t.Fatal("expected gate to fail on duplication")
	}
	if len(res.Failures) != 1 || !strings.Contains(res.Failures[0], "Duplication 12.5%") {
		t.Errorf("unexpected failure list: %v", res.Failures)
	}
}

func TestCheckGates_MaintainabilityFail(t *testing.T) {
	pm := &types.ProjectMetrics{MaintainabilityIdx: 30.0}
	res := checkGates(pm, types.GateConfig{MinMaintainability: 60.0})
	if res.Passed {
		t.Fatal("expected gate to fail on maintainability")
	}
	if len(res.Failures) != 1 || !strings.Contains(res.Failures[0], "Maintainability 30.0 below min 60.0") {
		t.Errorf("unexpected failure list: %v", res.Failures)
	}
}

func TestCheckGates_MultipleFailures(t *testing.T) {
	pm := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "F", File: "a.go", Cyclomatic: 30, Cognitive: 30, StartLine: 1},
			},
		}},
		DuplicationPct:     20,
		MaintainabilityIdx: 10,
	}
	res := checkGates(pm, types.GateConfig{
		MaxCyclomatic:      5,
		MaxCognitive:       5,
		MaxDuplication:     5,
		MinMaintainability: 60,
	})
	if res.Passed {
		t.Fatal("expected gate to fail across multiple checks")
	}
	if len(res.Failures) != 4 {
		t.Errorf("expected 4 failures, got %d: %v", len(res.Failures), res.Failures)
	}
}

// ---------------------------------------------------------------------------
// computeDiff
// ---------------------------------------------------------------------------

func TestComputeDiff_DeltasAndHotspots(t *testing.T) {
	base := &types.ProjectMetrics{
		AvgCyclomatic:      3.0,
		AvgCognitive:       2.0,
		TotalSLOC:          100,
		DuplicationPct:     1.0,
		MaintainabilityIdx: 80.0,
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "Existing", File: "a.go", Cyclomatic: 5},
			},
		}},
	}
	head := &types.ProjectMetrics{
		AvgCyclomatic:      5.0,
		AvgCognitive:       4.0,
		TotalSLOC:          150,
		DuplicationPct:     2.5,
		MaintainabilityIdx: 75.0,
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "Existing", File: "a.go", Cyclomatic: 5},
				{Name: "NewHotspot", File: "a.go", Cyclomatic: 25},
				{Name: "AlsoHot", File: "b.go", Cyclomatic: 12},
				{Name: "ColdFunc", File: "c.go", Cyclomatic: 4},
			},
		}},
	}

	diff := computeDiff(base, head)

	if diff.DeltaCyclomatic != 2.0 {
		t.Errorf("DeltaCyclomatic = %v, want 2.0", diff.DeltaCyclomatic)
	}
	if diff.DeltaCognitive != 2.0 {
		t.Errorf("DeltaCognitive = %v, want 2.0", diff.DeltaCognitive)
	}
	if diff.DeltaSLOC != 50 {
		t.Errorf("DeltaSLOC = %d, want 50", diff.DeltaSLOC)
	}
	if diff.DeltaDuplication != 1.5 {
		t.Errorf("DeltaDuplication = %v, want 1.5", diff.DeltaDuplication)
	}
	if diff.DeltaMaintainability != -5.0 {
		t.Errorf("DeltaMaintainability = %v, want -5.0", diff.DeltaMaintainability)
	}

	if len(diff.NewHotspots) != 2 {
		t.Fatalf("expected 2 new hotspots (CCN > 10 and new/grown), got %d: %+v",
			len(diff.NewHotspots), diff.NewHotspots)
	}
	if diff.NewHotspots[0].Name != "NewHotspot" {
		t.Errorf("hotspots not sorted by CCN desc: %+v", diff.NewHotspots)
	}
}

func TestComputeDiff_CapsHotspotsAtFive(t *testing.T) {
	headFiles := []types.FileMetrics{{Functions: nil}}
	for i := 0; i < 8; i++ {
		headFiles[0].Functions = append(headFiles[0].Functions, types.FunctionMetrics{
			Name:       "F" + string(rune('A'+i)),
			File:       "x.go",
			Cyclomatic: 11 + i,
		})
	}
	diff := computeDiff(&types.ProjectMetrics{}, &types.ProjectMetrics{Files: headFiles})
	if len(diff.NewHotspots) != 5 {
		t.Errorf("expected hotspots capped at 5, got %d", len(diff.NewHotspots))
	}
}

func TestComputeDiff_ExistingFunctionNotFlaggedUnlessJumpsBy6(t *testing.T) {
	base := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "F", File: "x.go", Cyclomatic: 11},
			},
		}},
	}
	// CCN went from 11 to 14 (+3) — should NOT be flagged (needs +5 jump per code).
	headSmall := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "F", File: "x.go", Cyclomatic: 14},
			},
		}},
	}
	if got := len(computeDiff(base, headSmall).NewHotspots); got != 0 {
		t.Errorf("expected 0 hotspots for +3 jump, got %d", got)
	}

	// CCN went from 11 to 17 (+6) — SHOULD be flagged.
	headBig := &types.ProjectMetrics{
		Files: []types.FileMetrics{{
			Functions: []types.FunctionMetrics{
				{Name: "F", File: "x.go", Cyclomatic: 17},
			},
		}},
	}
	if got := len(computeDiff(base, headBig).NewHotspots); got != 1 {
		t.Errorf("expected 1 hotspot for +6 jump, got %d", got)
	}
}

// ---------------------------------------------------------------------------
// checkCmd / gateCmd — full CLI invocation against a temp dir
// ---------------------------------------------------------------------------

func writeGoFixture(t *testing.T, dir string) {
	t.Helper()
	src := `package fixture

func Add(a, b int) int { return a + b }

func Branch(x int) string {
	if x > 0 {
		return "pos"
	}
	if x < 0 {
		return "neg"
	}
	return "zero"
}
`
	if err := os.WriteFile(filepath.Join(dir, "fixture.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckCmd_Smoke(t *testing.T) {
	dir := t.TempDir()
	writeGoFixture(t, dir)

	resetExtraSkipDirs(t)
	cmd := checkCmd()
	cmd.SetArgs([]string{"--no-color", "--skip", "vendor", dir})
	// Sink output so tests don't pollute stdout. cobra's SetOut/SetErr handles
	// errors; we silence stdout via os.Stdout swap for the duration.
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("checkCmd execute failed: %v", err)
	}
}

func TestCheckCmd_RelativePathFromCwd(t *testing.T) {
	dir := t.TempDir()
	writeGoFixture(t, dir)

	// Run with no positional arg → defaults to "." → use Chdir.
	origWD, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	resetExtraSkipDirs(t)
	cmd := checkCmd()
	cmd.SetArgs([]string{"--no-color"})
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("checkCmd execute (cwd) failed: %v", err)
	}
}

func TestGateCmd_PassesWhenNoThresholdsSet(t *testing.T) {
	dir := t.TempDir()
	writeGoFixture(t, dir)

	exited := -1
	swapExitFunc(t, func(code int) { exited = code })
	resetExtraSkipDirs(t)

	cmd := gateCmd()
	cmd.SetArgs([]string{"--no-color", dir})
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gateCmd execute failed: %v", err)
	}
	if exited != -1 {
		t.Errorf("expected exitFunc NOT to fire on pass, got exit=%d", exited)
	}
}

func TestGateCmd_FailsWhenThresholdTripped(t *testing.T) {
	dir := t.TempDir()
	// A heavy switch to force CCN > 1.
	src := `package fixture

func Branchy(x int) string {
	if x == 1 { return "a" }
	if x == 2 { return "b" }
	if x == 3 { return "c" }
	if x == 4 { return "d" }
	return "z"
}
`
	if err := os.WriteFile(filepath.Join(dir, "fix.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	exited := -1
	swapExitFunc(t, func(code int) { exited = code })
	resetExtraSkipDirs(t)

	cmd := gateCmd()
	// max-ccn=1 will definitely trip Branchy.
	cmd.SetArgs([]string{"--no-color", "--max-ccn", "1", dir})
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("gateCmd execute returned error: %v", err)
	}
	if exited != 1 {
		t.Errorf("expected exitFunc(1) to fire on gate fail, got exit=%d", exited)
	}
}

// ---------------------------------------------------------------------------
// trendCmd / diffCmd — git-backed
// ---------------------------------------------------------------------------

func initRepoWithCommits(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not available — skipping")
	}
	dir := t.TempDir()

	bootstrap := [][]string{
		{"git", "init", "--quiet", "--initial-branch=main"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "config", "commit.gpgsign", "false"},
	}
	for _, c := range bootstrap {
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			if strings.Contains(string(out), "initial-branch") {
				retry := exec.Command("git", "init", "--quiet")
				retry.Dir = dir
				if err2 := retry.Run(); err2 != nil {
					t.Fatalf("git init fallback: %v", err2)
				}
				continue
			}
			t.Fatalf("%v failed: %v\n%s", c, err, out)
		}
	}

	writeGoFixture(t, dir)
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "--quiet", "-m", "first")

	// Second commit with extra fn
	more := `package fixture

func Two() int { return 2 }
`
	if err := os.WriteFile(filepath.Join(dir, "two.go"), []byte(more), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "--quiet", "-m", "second")

	return dir
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func TestTrendCmd_Smoke(t *testing.T) {
	dir := initRepoWithCommits(t)

	cmd := trendCmd()
	cmd.SetArgs([]string{"--no-color", "--last", "5", dir})
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("trendCmd execute failed: %v", err)
	}
}

func TestDiffCmd_AgainstMergeBaseFallback(t *testing.T) {
	dir := initRepoWithCommits(t)

	cmd := diffCmd()
	// No --base flag → defaults to "main" → MergeBase succeeds for main itself.
	cmd.SetArgs([]string{"--no-color", dir})
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("diffCmd execute failed: %v", err)
	}
}

func TestDiffCmd_ExplicitBase(t *testing.T) {
	dir := initRepoWithCommits(t)

	// Use HEAD~1 as the base ref — guarantees a real ref exists.
	cmd := diffCmd()
	cmd.SetArgs([]string{"--no-color", "--base", "HEAD~1", dir})
	silenceStdout(t)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("diffCmd with --base HEAD~1 failed: %v", err)
	}
}

// ---------------------------------------------------------------------------
// newRootCmd / main wiring
// ---------------------------------------------------------------------------

func TestNewRootCmd_HasAllSubcommands(t *testing.T) {
	root := newRootCmd()
	want := map[string]bool{
		"check": false,
		"trend": false,
		"diff":  false,
		"gate":  false,
	}
	for _, c := range root.Commands() {
		if _, ok := want[c.Name()]; ok {
			want[c.Name()] = true
		}
	}
	for name, present := range want {
		if !present {
			t.Errorf("expected subcommand %q to be registered", name)
		}
	}
}

func TestNewRootCmd_HelpRuns(t *testing.T) {
	root := newRootCmd()
	root.SetArgs([]string{"--help"})
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	if err := root.Execute(); err != nil {
		t.Fatalf("root --help failed: %v", err)
	}
	if !strings.Contains(out.String(), "pulse") {
		t.Errorf("expected help output to mention 'pulse', got:\n%s", out.String())
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// silenceStdout redirects os.Stdout to a temp file for the duration of the
// test, restoring it on cleanup. We do this because report.Print* writes to
// os.Stdout directly via the `color` package and would clutter `go test -v`.
func silenceStdout(t *testing.T) {
	t.Helper()
	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Skip("cannot open /dev/null on this platform — skipping silencing")
	}
	orig := os.Stdout
	os.Stdout = devnull
	t.Cleanup(func() {
		os.Stdout = orig
		_ = devnull.Close()
	})
}

// swapExitFunc installs a stub exitFunc and restores the original on cleanup.
func swapExitFunc(t *testing.T, fn func(int)) {
	t.Helper()
	orig := exitFunc
	exitFunc = fn
	t.Cleanup(func() { exitFunc = orig })
}

// resetExtraSkipDirs snapshots analyzer.ExtraSkipDirs (a package-level
// global mutated by checkCmd and gateCmd) and restores it on cleanup so
// tests don't leak state into each other.
func resetExtraSkipDirs(t *testing.T) {
	t.Helper()
	orig := analyzer.ExtraSkipDirs
	t.Cleanup(func() { analyzer.ExtraSkipDirs = orig })
}
