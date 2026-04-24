package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyze(t *testing.T) {
	dir := testdataDir(t)

	pm, err := Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	// Should find at least the Go, Python, and TypeScript test files
	if pm.TotalFiles < 3 {
		t.Errorf("expected at least 3 files, got %d", pm.TotalFiles)
	}

	// Should find functions
	if pm.TotalFunctions == 0 {
		t.Error("expected functions to be found")
	}

	// SLOC should be populated
	if pm.TotalSLOC == 0 {
		t.Error("expected SLOC > 0")
	}

	// Language summaries should exist
	if len(pm.LanguageSummaries) == 0 {
		t.Error("expected language summaries")
	}

	// Maintainability should be calculated
	if pm.MaintainabilityIdx <= 0 {
		t.Errorf("expected MI > 0, got %.1f", pm.MaintainabilityIdx)
	}
	if pm.Grade == "" {
		t.Error("expected grade to be set")
	}

	// COCOMO should be calculated
	if pm.COCOMOMonths <= 0 {
		t.Error("expected COCOMO months > 0")
	}

	t.Logf("Analysis results:")
	t.Logf("  Files: %d, Functions: %d", pm.TotalFiles, pm.TotalFunctions)
	t.Logf("  SLOC: %d, Comments: %d, Blanks: %d", pm.TotalSLOC, pm.TotalComments, pm.TotalBlanks)
	t.Logf("  Avg CCN: %.1f, Avg Cognitive: %.1f", pm.AvgCyclomatic, pm.AvgCognitive)
	t.Logf("  MI: %.1f (%s)", pm.MaintainabilityIdx, pm.Grade)
	t.Logf("  Duplication: %.1f%%", pm.DuplicationPct)
	t.Logf("  COCOMO: %.1f person-months", pm.COCOMOMonths)

	for _, ls := range pm.LanguageSummaries {
		t.Logf("  %s: %d files, %d SLOC", ls.Language, ls.Files, ls.SLOC)
	}
}

func TestAnalyze_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	pm, err := Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if pm.TotalFiles != 0 {
		t.Errorf("expected 0 files for empty dir, got %d", pm.TotalFiles)
	}
}

func TestAnalyze_SkipsDirs(t *testing.T) {
	dir := t.TempDir()

	// Create a file in node_modules — should be skipped
	nmDir := filepath.Join(dir, "node_modules")
	os.MkdirAll(nmDir, 0755)
	os.WriteFile(filepath.Join(nmDir, "test.js"), []byte("function x() {}"), 0644)

	// Create a real file
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0644)

	pm, err := Analyze(dir)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if pm.TotalFiles != 1 {
		t.Errorf("expected 1 file (node_modules should be skipped), got %d", pm.TotalFiles)
	}
}

func testdataDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, "..", "..", "testdata")
}
