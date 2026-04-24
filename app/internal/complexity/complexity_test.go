package complexity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bobbydeveaux/pulse/internal/types"
)

func TestAnalyzeGoFile(t *testing.T) {
	path := filepath.Join(testdataDir(t), "simple.go")
	funcs, err := AnalyzeGoFile(path)
	if err != nil {
		t.Fatalf("AnalyzeGoFile failed: %v", err)
	}

	if len(funcs) != 3 {
		t.Fatalf("expected 3 functions, got %d", len(funcs))
	}

	// SimpleFunc: should have CCN 1 (no branches)
	simple := findFunc(funcs, "SimpleFunc")
	if simple == nil {
		t.Fatal("SimpleFunc not found")
	}
	if simple.Cyclomatic != 1 {
		t.Errorf("SimpleFunc CCN: expected 1, got %d", simple.Cyclomatic)
	}
	if simple.Cognitive != 0 {
		t.Errorf("SimpleFunc Cognitive: expected 0, got %d", simple.Cognitive)
	}
	if simple.Grade != "A" {
		t.Errorf("SimpleFunc Grade: expected A, got %s", simple.Grade)
	}

	// MediumFunc: has if/else if/for/if chains
	medium := findFunc(funcs, "MediumFunc")
	if medium == nil {
		t.Fatal("MediumFunc not found")
	}
	if medium.Cyclomatic < 4 {
		t.Errorf("MediumFunc CCN: expected >= 4, got %d", medium.Cyclomatic)
	}
	if medium.Cognitive < 3 {
		t.Errorf("MediumFunc Cognitive: expected >= 3, got %d", medium.Cognitive)
	}

	// ComplexFunc: has nested if/switch/for/&& /||
	complex := findFunc(funcs, "ComplexFunc")
	if complex == nil {
		t.Fatal("ComplexFunc not found")
	}
	if complex.Cyclomatic < 10 {
		t.Errorf("ComplexFunc CCN: expected >= 10, got %d", complex.Cyclomatic)
	}
	if complex.Cognitive < 5 {
		t.Errorf("ComplexFunc Cognitive: expected >= 5, got %d", complex.Cognitive)
	}

	t.Logf("Results:")
	for _, f := range funcs {
		t.Logf("  %s: CCN=%d Cognitive=%d Grade=%s",
			f.Name, f.Cyclomatic, f.Cognitive, f.Grade)
	}
}

func TestAnalyzePythonFile(t *testing.T) {
	path := filepath.Join(testdataDir(t), "simple.py")
	funcs, err := AnalyzeFile(path)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if len(funcs) < 2 {
		t.Fatalf("expected at least 2 functions, got %d", len(funcs))
	}

	// hello: should be simple
	hello := findFunc(funcs, "hello")
	if hello == nil {
		t.Fatal("hello not found")
	}
	if hello.Cyclomatic != 1 {
		t.Errorf("hello CCN: expected 1, got %d", hello.Cyclomatic)
	}

	// calculate: has if/elif chain
	calc := findFunc(funcs, "calculate")
	if calc == nil {
		t.Fatal("calculate not found")
	}
	if calc.Cyclomatic < 4 {
		t.Errorf("calculate CCN: expected >= 4, got %d", calc.Cyclomatic)
	}

	t.Logf("Python results:")
	for _, f := range funcs {
		t.Logf("  %s: CCN=%d Cognitive=%d Grade=%s",
			f.Name, f.Cyclomatic, f.Cognitive, f.Grade)
	}
}

func TestAnalyzeTypeScriptFile(t *testing.T) {
	path := filepath.Join(testdataDir(t), "simple.ts")
	funcs, err := AnalyzeFile(path)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if len(funcs) < 2 {
		t.Fatalf("expected at least 2 functions, got %d", len(funcs))
	}

	// greet: should be simple
	greet := findFunc(funcs, "greet")
	if greet == nil {
		t.Fatal("greet not found")
	}
	if greet.Cyclomatic != 1 {
		t.Errorf("greet CCN: expected 1, got %d", greet.Cyclomatic)
	}

	t.Logf("TypeScript results:")
	for _, f := range funcs {
		t.Logf("  %s: CCN=%d Cognitive=%d Grade=%s",
			f.Name, f.Cyclomatic, f.Cognitive, f.Grade)
	}
}

func findFunc(funcs []types.FunctionMetrics, name string) *types.FunctionMetrics {
	for i := range funcs {
		if funcs[i].Name == name {
			return &funcs[i]
		}
	}
	return nil
}

func testdataDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, "..", "..", "testdata")
}
