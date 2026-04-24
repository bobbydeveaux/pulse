package sloc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCountFile_Go(t *testing.T) {
	path := filepath.Join(testdataDir(t), "simple.go")
	c, err := CountFile(path)
	if err != nil {
		t.Fatalf("CountFile failed: %v", err)
	}

	if c.Language != "Go" {
		t.Errorf("expected language Go, got %s", c.Language)
	}
	if c.SLOC == 0 {
		t.Error("expected SLOC > 0")
	}
	if c.Comments == 0 {
		t.Error("expected Comments > 0 (file has doc comments)")
	}
	if c.Blanks == 0 {
		t.Error("expected Blanks > 0")
	}
	if c.TotalLines != c.SLOC+c.Comments+c.Blanks {
		t.Errorf("total lines %d != SLOC(%d) + Comments(%d) + Blanks(%d)",
			c.TotalLines, c.SLOC, c.Comments, c.Blanks)
	}

	t.Logf("Go: SLOC=%d Comments=%d Blanks=%d Total=%d",
		c.SLOC, c.Comments, c.Blanks, c.TotalLines)
}

func TestCountFile_Python(t *testing.T) {
	path := filepath.Join(testdataDir(t), "simple.py")
	c, err := CountFile(path)
	if err != nil {
		t.Fatalf("CountFile failed: %v", err)
	}

	if c.Language != "Python" {
		t.Errorf("expected language Python, got %s", c.Language)
	}
	if c.SLOC == 0 {
		t.Error("expected SLOC > 0")
	}

	t.Logf("Python: SLOC=%d Comments=%d Blanks=%d Total=%d",
		c.SLOC, c.Comments, c.Blanks, c.TotalLines)
}

func TestCountFile_TypeScript(t *testing.T) {
	path := filepath.Join(testdataDir(t), "simple.ts")
	c, err := CountFile(path)
	if err != nil {
		t.Fatalf("CountFile failed: %v", err)
	}

	if c.Language != "TypeScript" {
		t.Errorf("expected language TypeScript, got %s", c.Language)
	}
	if c.SLOC == 0 {
		t.Error("expected SLOC > 0")
	}

	t.Logf("TypeScript: SLOC=%d Comments=%d Blanks=%d Total=%d",
		c.SLOC, c.Comments, c.Blanks, c.TotalLines)
}

func TestCountFile_UnknownExtension(t *testing.T) {
	// Create a temp file with unknown extension
	tmp := filepath.Join(t.TempDir(), "test.xyz")
	os.WriteFile(tmp, []byte("hello\nworld\n"), 0644)

	c, err := CountFile(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Language != "" {
		t.Errorf("expected empty language for unknown extension, got %s", c.Language)
	}
}

func testdataDir(t *testing.T) string {
	t.Helper()
	// Navigate from internal/sloc/ to testdata/
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(dir, "..", "..", "testdata")
}
