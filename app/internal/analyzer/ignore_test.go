package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPulseignore_Missing(t *testing.T) {
	dir := t.TempDir()
	patterns, err := loadPulseignore(dir)
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if patterns != nil {
		t.Errorf("expected nil patterns for missing file, got %v", patterns)
	}
}

func TestLoadPulseignore_Empty(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".pulseignore"), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	patterns, err := loadPulseignore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patterns) != 0 {
		t.Errorf("expected 0 patterns for empty file, got %d", len(patterns))
	}
}

func TestLoadPulseignore_CommentsOnly(t *testing.T) {
	dir := t.TempDir()
	body := "# this is a comment\n\n# another comment\n   \n"
	if err := os.WriteFile(filepath.Join(dir, ".pulseignore"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	patterns, err := loadPulseignore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patterns) != 0 {
		t.Errorf("expected 0 patterns for comments-only file, got %d", len(patterns))
	}
}

func TestLoadPulseignore_DirectoryAndFilePatterns(t *testing.T) {
	dir := t.TempDir()
	body := "build/\n*.log\nfoo/bar/\n# comment\nspecific.txt\n"
	if err := os.WriteFile(filepath.Join(dir, ".pulseignore"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	patterns, err := loadPulseignore(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(patterns) != 4 {
		t.Fatalf("expected 4 patterns, got %d (%v)", len(patterns), patterns)
	}

	// Pattern order matches file order
	if patterns[0].pattern != "build" || !patterns[0].isDir {
		t.Errorf("patterns[0]: expected {build true}, got %+v", patterns[0])
	}
	if patterns[1].pattern != "*.log" || patterns[1].isDir {
		t.Errorf("patterns[1]: expected {*.log false}, got %+v", patterns[1])
	}
	if patterns[2].pattern != "foo/bar" || !patterns[2].isDir {
		t.Errorf("patterns[2]: expected {foo/bar true}, got %+v", patterns[2])
	}
	if patterns[3].pattern != "specific.txt" || patterns[3].isDir {
		t.Errorf("patterns[3]: expected {specific.txt false}, got %+v", patterns[3])
	}
}

func TestLoadPulseignore_PermissionError(t *testing.T) {
	// Open error other than NotExist: emulate with a directory at the .pulseignore path
	dir := t.TempDir()
	target := filepath.Join(dir, ".pulseignore")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err := loadPulseignore(dir)
	if err == nil {
		t.Error("expected error when .pulseignore is a directory, got nil")
	}
}

func TestMatchesIgnorePattern_DirectoryEntry(t *testing.T) {
	patterns := []ignorePattern{{pattern: "build", isDir: true}}
	info := fakeStat(t, "build", true)
	if !matchesIgnorePattern("build", info, patterns) {
		t.Error("expected directory entry 'build' to match dir pattern 'build'")
	}
}

func TestMatchesIgnorePattern_DirectoryContents(t *testing.T) {
	patterns := []ignorePattern{{pattern: "build", isDir: true}}
	fileInfo := fakeStat(t, "out.txt", false)
	if !matchesIgnorePattern(filepath.Join("build", "out.txt"), fileInfo, patterns) {
		t.Error("expected file inside build/ to match dir pattern 'build'")
	}
	// Nested deeper
	nested := filepath.Join("build", "sub", "deep.txt")
	if !matchesIgnorePattern(nested, fileInfo, patterns) {
		t.Errorf("expected nested file %s to match dir pattern", nested)
	}
}

func TestMatchesIgnorePattern_DirectoryNoMatchForSiblingPrefix(t *testing.T) {
	patterns := []ignorePattern{{pattern: "build", isDir: true}}
	// "buildkit" should not match "build" — separator boundary required
	info := fakeStat(t, "buildkit", true)
	if matchesIgnorePattern("buildkit", info, patterns) {
		t.Error("expected 'buildkit' to NOT match dir pattern 'build' (prefix boundary)")
	}
}

func TestMatchesIgnorePattern_FileGlobMatchesBase(t *testing.T) {
	patterns := []ignorePattern{{pattern: "*.log", isDir: false}}
	info := fakeStat(t, "server.log", false)
	if !matchesIgnorePattern("server.log", info, patterns) {
		t.Error("expected '*.log' glob to match 'server.log'")
	}
}

func TestMatchesIgnorePattern_FileGlobMatchesFullPath(t *testing.T) {
	patterns := []ignorePattern{{pattern: "logs/*.txt", isDir: false}}
	info := fakeStat(t, "out.txt", false)
	if !matchesIgnorePattern(filepath.Join("logs", "out.txt"), info, patterns) {
		t.Error("expected 'logs/*.txt' to match relative path 'logs/out.txt'")
	}
}

func TestMatchesIgnorePattern_NoMatch(t *testing.T) {
	patterns := []ignorePattern{
		{pattern: "build", isDir: true},
		{pattern: "*.log", isDir: false},
	}
	info := fakeStat(t, "main.go", false)
	if matchesIgnorePattern("main.go", info, patterns) {
		t.Error("expected 'main.go' to not match any pattern")
	}
}

func TestMatchesIgnorePattern_EmptyPatterns(t *testing.T) {
	info := fakeStat(t, "any.txt", false)
	if matchesIgnorePattern("any.txt", info, nil) {
		t.Error("expected empty patterns to match nothing")
	}
}

func TestMatchesIgnorePattern_DirSelfEntryFile(t *testing.T) {
	// A file whose relPath equals the directory pattern itself — should NOT match
	// (only directories match by exact name, files inside match by prefix+sep).
	patterns := []ignorePattern{{pattern: "build", isDir: true}}
	info := fakeStat(t, "build", false)
	if matchesIgnorePattern("build", info, patterns) {
		t.Error("expected a file literally named 'build' to NOT match dir pattern (no separator)")
	}
}

// fakeStat returns an os.FileInfo for the given name and dir flag by stat'ing
// a real file we create in t.TempDir. This avoids implementing FileInfo
// fully (Sys/ModTime types are tricky to satisfy across Go versions).
func fakeStat(t *testing.T, name string, isDir bool) os.FileInfo {
	t.Helper()
	dir := t.TempDir()
	full := filepath.Join(dir, name)
	if isDir {
		if err := os.Mkdir(full, 0o755); err != nil {
			t.Fatal(err)
		}
	} else {
		if err := os.WriteFile(full, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	info, err := os.Stat(full)
	if err != nil {
		t.Fatal(err)
	}
	return info
}
