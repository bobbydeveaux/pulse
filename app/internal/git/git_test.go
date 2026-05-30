package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// initRepo builds a small git repo in a temp directory and returns its path.
// Skips the calling test if `git` is not on PATH.
func initRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not available — skipping")
	}
	dir := t.TempDir()
	cmds := [][]string{
		{"git", "init", "--quiet", "--initial-branch=main"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test"},
		{"git", "config", "commit.gpgsign", "false"},
	}
	for _, c := range cmds {
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			// Fall back without --initial-branch flag (older git)
			if strings.Contains(string(out), "initial-branch") {
				retry := exec.Command("git", "init", "--quiet")
				retry.Dir = dir
				if err2 := retry.Run(); err2 != nil {
					t.Fatalf("git init fallback failed: %v", err2)
				}
				continue
			}
			t.Fatalf("%v failed: %v\n%s", c, err, out)
		}
	}
	// First commit
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "git", "add", "a.txt")
	mustRun(t, dir, "git", "commit", "--quiet", "-m", "first commit")
	return dir
}

func mustRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v failed: %v\n%s", args, err, out)
	}
}

func TestLog(t *testing.T) {
	dir := initRepo(t)

	// Add a second commit so we have two
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "git", "add", "b.txt")
	mustRun(t, dir, "git", "commit", "--quiet", "-m", "second commit")

	commits, err := Log(5, dir)
	if err != nil {
		t.Fatalf("Log failed: %v", err)
	}
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}
	if commits[0].Message != "second commit" {
		t.Errorf("expected newest commit message 'second commit', got %q", commits[0].Message)
	}
	if commits[0].Hash == "" || commits[0].ShortHash == "" {
		t.Errorf("commit hashes should be populated: %+v", commits[0])
	}
	if commits[0].Author == "" || commits[0].Date == "" {
		t.Errorf("commit author/date should be populated: %+v", commits[0])
	}
}

func TestLog_Error(t *testing.T) {
	if _, err := Log(1, t.TempDir()); err == nil {
		t.Error("expected error from non-git directory, got nil")
	}
}

func TestCurrentRef(t *testing.T) {
	dir := initRepo(t)
	ref, err := CurrentRef(dir)
	if err != nil {
		t.Fatalf("CurrentRef failed: %v", err)
	}
	if len(ref) != 40 {
		t.Errorf("expected 40-char SHA, got %q (%d)", ref, len(ref))
	}
}

func TestCurrentRef_Error(t *testing.T) {
	if _, err := CurrentRef(t.TempDir()); err == nil {
		t.Error("expected error from non-git directory")
	}
}

func TestCurrentBranch(t *testing.T) {
	dir := initRepo(t)
	branch, err := CurrentBranch(dir)
	if err != nil {
		t.Fatalf("CurrentBranch failed: %v", err)
	}
	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestCurrentBranch_Error(t *testing.T) {
	if _, err := CurrentBranch(t.TempDir()); err == nil {
		t.Error("expected error from non-git directory")
	}
}

func TestMergeBase(t *testing.T) {
	dir := initRepo(t)
	ref, err := CurrentRef(dir)
	if err != nil {
		t.Fatal(err)
	}
	mb, err := MergeBase(ref, ref, dir)
	if err != nil {
		t.Fatalf("MergeBase failed: %v", err)
	}
	if mb != ref {
		t.Errorf("merge-base of ref with itself should be the ref; got %q vs %q", mb, ref)
	}
}

func TestMergeBase_Error(t *testing.T) {
	if _, err := MergeBase("a", "b", t.TempDir()); err == nil {
		t.Error("expected error from non-git directory")
	}
}

func TestIsDirty_Clean(t *testing.T) {
	dir := initRepo(t)
	dirty, err := IsDirty(dir)
	if err != nil {
		t.Fatalf("IsDirty failed: %v", err)
	}
	if dirty {
		t.Error("expected clean repo to report not dirty")
	}
}

func TestIsDirty_Modified(t *testing.T) {
	dir := initRepo(t)
	if err := os.WriteFile(filepath.Join(dir, "c.txt"), []byte("change\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	dirty, err := IsDirty(dir)
	if err != nil {
		t.Fatalf("IsDirty failed: %v", err)
	}
	if !dirty {
		t.Error("expected modified repo to report dirty")
	}
}

func TestIsDirty_Error(t *testing.T) {
	if _, err := IsDirty(t.TempDir()); err == nil {
		t.Error("expected error from non-git directory")
	}
}

func TestStashAndPop(t *testing.T) {
	dir := initRepo(t)
	// Create an untracked file
	if err := os.WriteFile(filepath.Join(dir, "untracked.txt"), []byte("temp\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Stash(dir); err != nil {
		t.Fatalf("Stash failed: %v", err)
	}

	// After stash --include-untracked the file should be gone
	if _, err := os.Stat(filepath.Join(dir, "untracked.txt")); !os.IsNotExist(err) {
		t.Errorf("expected stashed file to be removed, got err=%v", err)
	}

	if err := StashPop(dir); err != nil {
		t.Fatalf("StashPop failed: %v", err)
	}

	// File should be back
	if _, err := os.Stat(filepath.Join(dir, "untracked.txt")); err != nil {
		t.Errorf("expected file restored after StashPop, got %v", err)
	}
}

func TestStash_Error(t *testing.T) {
	if err := Stash(t.TempDir()); err == nil {
		t.Error("expected error stashing in non-git directory")
	}
}

func TestStashPop_Error(t *testing.T) {
	if err := StashPop(t.TempDir()); err == nil {
		t.Error("expected error popping in non-git directory")
	}
}

func TestCheckout(t *testing.T) {
	dir := initRepo(t)

	// Make a second commit, then a new branch pointing at the first
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("second\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "git", "add", "b.txt")
	mustRun(t, dir, "git", "commit", "--quiet", "-m", "second")

	// Find the first commit hash
	commits, err := Log(5, dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) < 2 {
		t.Fatalf("need 2 commits, got %d", len(commits))
	}
	firstHash := commits[len(commits)-1].Hash

	if err := Checkout(firstHash, dir); err != nil {
		t.Fatalf("Checkout failed: %v", err)
	}

	// After checkout, b.txt should not exist
	if _, err := os.Stat(filepath.Join(dir, "b.txt")); !os.IsNotExist(err) {
		t.Errorf("expected b.txt absent after checkout to first commit, got err=%v", err)
	}
}

func TestCheckout_Error(t *testing.T) {
	if err := Checkout("nonexistent", t.TempDir()); err == nil {
		t.Error("expected error checking out in non-git dir")
	}
}

func TestItoa(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{
		{0, "0"},
		{1, "1"},
		{42, "42"},
		{-1, "-1"},
	}
	for _, c := range cases {
		if got := itoa(c.in); got != c.want {
			t.Errorf("itoa(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestLog_EmptyRepo(t *testing.T) {
	// `git log` on a repo with no commits returns a non-zero exit — we hit
	// the err branch of Log.
	dir := t.TempDir()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git binary not available")
	}
	cmd := exec.Command("git", "init", "--quiet")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	_, err := Log(5, dir)
	if err == nil {
		t.Error("expected error from git log on empty repo")
	}
}
