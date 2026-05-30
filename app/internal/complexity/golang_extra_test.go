package complexity

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReceiverName_ViaAST exercises receiverName by parsing real Go source
// containing several method receiver shapes (value, pointer, generic).
func TestReceiverName_ViaAST(t *testing.T) {
	src := `package x

type S struct{}
type G[T any] struct{}

func (s S) ValueMethod() {}
func (s *S) PointerMethod() {}
func (g G[T]) GenericValueMethod() {}
func (g *G[T]) GenericPointerMethod() {}

func PlainFunc() {}
`
	tmp := filepath.Join(t.TempDir(), "recv.go")
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := AnalyzeGoFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeGoFile failed: %v", err)
	}

	want := map[string]bool{
		"S.ValueMethod":          false,
		"S.PointerMethod":        false,
		"G.GenericValueMethod":   false,
		"G.GenericPointerMethod": false,
		"PlainFunc":              false,
	}
	for _, f := range funcs {
		if _, ok := want[f.Name]; ok {
			want[f.Name] = true
		}
	}
	for name, found := range want {
		if !found {
			t.Errorf("expected function name %q, not found in %+v", name, funcs)
		}
	}
}

// TestCognitiveStmt_ElseIfAndLabels exercises switch, type switch, select,
// range, labeled stmt, and labeled branch — covering the cognitiveStmt
// branches that testdata/simple.go does not.
func TestCognitiveStmt_ElseIfAndLabels(t *testing.T) {
	src := `package x

func DeepFunc(x int) int {
	switch x {
	case 1:
		return 1
	case 2:
		return 2
	}

	var v interface{} = x
	switch v.(type) {
	case int:
		_ = v
	case string:
		_ = v
	}

	ch := make(chan int)
	select {
	case <-ch:
		return 3
	default:
		return 4
	}
}

func ElseIf(x int) int {
	if x == 1 {
		return 1
	} else if x == 2 {
		return 2
	} else {
		return 3
	}
}

func LoopWithLabel() {
Outer:
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if j == 5 {
				break Outer
			}
		}
	}
}

func RangeFunc(items []int) int {
	total := 0
	for _, v := range items {
		total += v
	}
	return total
}
`
	tmp := filepath.Join(t.TempDir(), "deep.go")
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	funcs, err := AnalyzeGoFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeGoFile failed: %v", err)
	}

	if len(funcs) < 4 {
		t.Fatalf("expected >=4 functions, got %d", len(funcs))
	}

	deep := findFunc(funcs, "DeepFunc")
	if deep == nil {
		t.Fatal("DeepFunc not found")
	}
	if deep.Cognitive == 0 {
		t.Error("expected DeepFunc cognitive > 0")
	}

	ei := findFunc(funcs, "ElseIf")
	if ei == nil {
		t.Fatal("ElseIf not found")
	}
	if ei.Cognitive == 0 {
		t.Error("expected ElseIf cognitive > 0")
	}

	loop := findFunc(funcs, "LoopWithLabel")
	if loop == nil {
		t.Fatal("LoopWithLabel not found")
	}
	if loop.Cognitive == 0 {
		t.Error("expected LoopWithLabel cognitive > 0 (label penalty)")
	}

	rng := findFunc(funcs, "RangeFunc")
	if rng == nil {
		t.Fatal("RangeFunc not found")
	}
	if rng.Cyclomatic < 2 {
		t.Errorf("RangeFunc cyclomatic expected >=2 (range), got %d", rng.Cyclomatic)
	}
}

// TestCyclomaticGo_DefaultOnlyClauses ensures default cases and default-only
// select clauses don't bump cyclomatic complexity.
func TestCyclomaticGo_DefaultOnlyClauses(t *testing.T) {
	src := `package x

func OnlyDefault(x int) int {
	switch x {
	default:
		return 0
	}
}

func SelectDefaultOnly() {
	select {
	default:
		return
	}
}
`
	tmp := filepath.Join(t.TempDir(), "sw.go")
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := AnalyzeGoFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeGoFile failed: %v", err)
	}
	od := findFunc(funcs, "OnlyDefault")
	if od == nil {
		t.Fatal("OnlyDefault not found")
	}
	// CCN should remain at 1: switch alone doesn't bump, and `default:` case
	// (cc.List == nil) shouldn't bump either.
	if od.Cyclomatic != 1 {
		t.Errorf("OnlyDefault CCN expected 1 (default-only), got %d", od.Cyclomatic)
	}

	sd := findFunc(funcs, "SelectDefaultOnly")
	if sd == nil {
		t.Fatal("SelectDefaultOnly not found")
	}
	if sd.Cyclomatic != 1 {
		t.Errorf("SelectDefaultOnly CCN expected 1 (default-only), got %d", sd.Cyclomatic)
	}
}

// TestCyclomaticGo_NoBody covers the early return when fn.Body is nil
// (function declarations without bodies — e.g., //go:linkname assembly stubs).
func TestCyclomaticGo_NoBody(t *testing.T) {
	src := `package x

import _ "unsafe"

//go:linkname extern runtime.foo
func extern()
`
	tmp := filepath.Join(t.TempDir(), "nobody.go")
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := AnalyzeGoFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeGoFile failed: %v", err)
	}
	if len(funcs) != 1 {
		t.Fatalf("expected 1 function (extern), got %d", len(funcs))
	}
	if funcs[0].Cyclomatic != 1 {
		t.Errorf("extern CCN expected 1 (no body), got %d", funcs[0].Cyclomatic)
	}
	if funcs[0].Cognitive != 0 {
		t.Errorf("extern Cognitive expected 0 (no body), got %d", funcs[0].Cognitive)
	}
}

// TestCountBoolOps covers the && / || branch counting helper directly via
// an expression parsed as part of an if condition.
func TestCountBoolOps_ViaConditions(t *testing.T) {
	src := `package x

func ManyBools(a, b, c, d int) int {
	if a > 0 && b > 0 || c > 0 && d > 0 {
		return 1
	}
	return 0
}
`
	tmp := filepath.Join(t.TempDir(), "bools.go")
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := AnalyzeGoFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeGoFile failed: %v", err)
	}
	mb := findFunc(funcs, "ManyBools")
	if mb == nil {
		t.Fatal("ManyBools not found")
	}
	// CCN: 1 base + 1 if + 3 boolean ops = 5
	if mb.Cyclomatic < 5 {
		t.Errorf("ManyBools CCN expected >=5, got %d", mb.Cyclomatic)
	}
	// Cognitive includes +1 for if, +3 boolean ops (no nesting penalty on bool ops)
	if mb.Cognitive < 4 {
		t.Errorf("ManyBools Cognitive expected >=4, got %d", mb.Cognitive)
	}
}
