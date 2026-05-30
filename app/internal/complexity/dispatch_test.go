package complexity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bobbydeveaux/pulse/internal/languages"
)

func TestAnalyzeFile_UnknownExtension(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "data.unknownext")
	if err := os.WriteFile(tmp, []byte("nothing here"), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := AnalyzeFile(tmp)
	if err != nil {
		t.Fatalf("expected nil error for unknown ext, got %v", err)
	}
	if funcs != nil {
		t.Errorf("expected nil funcs for unknown ext, got %v", funcs)
	}
}

func TestAnalyzeFile_GoDispatch(t *testing.T) {
	// Confirm .go files route through AnalyzeGoFile by giving valid Go source.
	tmp := filepath.Join(t.TempDir(), "x.go")
	src := "package x\n\nfunc Foo() int { return 1 }\n"
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := AnalyzeFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "Foo" {
		t.Errorf("expected single Foo function, got %+v", funcs)
	}
}

func TestAnalyzeFile_GenericDispatch(t *testing.T) {
	// .py routes through the generic analyzer.
	tmp := filepath.Join(t.TempDir(), "x.py")
	src := "def hello():\n    return 1\n"
	if err := os.WriteFile(tmp, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	funcs, err := AnalyzeFile(tmp)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}
	if len(funcs) != 1 || funcs[0].Name != "hello" {
		t.Errorf("expected hello function, got %+v", funcs)
	}
}

func TestAnalyzeGoFile_ParseError(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "broken.go")
	if err := os.WriteFile(tmp, []byte("package x\nfunc { not valid"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := AnalyzeGoFile(tmp)
	if err == nil {
		t.Error("expected parse error on malformed Go source")
	}
}

func TestAnalyzeGenericFile_OpenError(t *testing.T) {
	lang, ok := languages.ByExtension(".py")
	if !ok {
		t.Skip("Python not registered in languages package")
	}
	_, err := AnalyzeGenericFile(filepath.Join(t.TempDir(), "missing.py"), lang)
	if err == nil {
		t.Error("expected open error for missing file")
	}
}
