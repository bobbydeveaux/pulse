package duplication

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetect_NoDuplicates(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.go"), `package main

func main() {
	fmt.Println("hello world")
	x := 1 + 2
	y := x * 3
}
`)

	writeFile(t, filepath.Join(dir, "b.go"), `package main

func other() {
	result := calculate(10, 20)
	if result > 0 {
		process(result)
	}
}
`)

	dupes, pct := Detect([]string{
		filepath.Join(dir, "a.go"),
		filepath.Join(dir, "b.go"),
	})

	if len(dupes) > 0 {
		t.Errorf("expected 0 duplicates, got %d", len(dupes))
	}
	if pct > 0 {
		t.Errorf("expected 0%% duplication, got %.1f%%", pct)
	}
}

func TestDetect_WithDuplicates(t *testing.T) {
	dir := t.TempDir()

	// Create two files with a shared block of 6+ lines
	sharedBlock := `	if err != nil {
		log.Printf("error: %v", err)
		metrics.IncrementErrors()
		notifyAdmin(err)
		return fmt.Errorf("operation failed: %w", err)
	}
	log.Println("operation succeeded")
	metrics.IncrementSuccess()`

	writeFile(t, filepath.Join(dir, "a.go"), fmt.Sprintf(`package main

func handleA() error {
	err := doA()
%s
	return nil
}
`, sharedBlock))

	writeFile(t, filepath.Join(dir, "b.go"), fmt.Sprintf(`package main

func handleB() error {
	err := doB()
%s
	return nil
}
`, sharedBlock))

	dupes, pct := Detect([]string{
		filepath.Join(dir, "a.go"),
		filepath.Join(dir, "b.go"),
	})

	if len(dupes) == 0 {
		t.Error("expected duplicates to be detected")
	}
	if pct <= 0 {
		t.Error("expected duplication percentage > 0")
	}

	t.Logf("Found %d duplicate(s), %.1f%% duplication", len(dupes), pct)
	for _, d := range dupes {
		t.Logf("  %s:%d <-> %s:%d (%d lines)",
			filepath.Base(d.FileA), d.StartA,
			filepath.Base(d.FileB), d.StartB, d.LineCount)
	}
}

func TestDetect_SkipsTrivialBlocks(t *testing.T) {
	dir := t.TempDir()

	// Create files with repeated but trivial blocks (just braces/blank)
	trivial := strings.Repeat("{\n}\n\n", 10)

	writeFile(t, filepath.Join(dir, "a.go"), trivial)
	writeFile(t, filepath.Join(dir, "b.go"), trivial)

	dupes, _ := Detect([]string{
		filepath.Join(dir, "a.go"),
		filepath.Join(dir, "b.go"),
	})

	if len(dupes) > 0 {
		t.Errorf("expected trivial blocks to be skipped, got %d duplicates", len(dupes))
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
