package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Log returns the last N commit hashes from the current branch.
func Log(n int, dir string) ([]Commit, error) {
	cmd := exec.Command("git", "log", "--format=%H|%h|%an|%ai|%s", "-n", itoa(n))
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []Commit
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}
		commits = append(commits, Commit{
			Hash:      parts[0],
			ShortHash: parts[1],
			Author:    parts[2],
			Date:      parts[3],
			Message:   parts[4],
		})
	}
	return commits, nil
}

// Commit represents a git commit.
type Commit struct {
	Hash      string
	ShortHash string
	Author    string
	Date      string
	Message   string
}

// Checkout switches to a given ref.
func Checkout(ref, dir string) error {
	cmd := exec.Command("git", "checkout", ref, "--quiet")
	cmd.Dir = dir
	return cmd.Run()
}

// CurrentRef returns the current HEAD ref.
func CurrentRef(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// CurrentBranch returns the current branch name.
func CurrentBranch(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// MergeBase returns the merge base between two refs.
func MergeBase(a, b, dir string) (string, error) {
	cmd := exec.Command("git", "merge-base", a, b)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// Stash saves current changes (including untracked).
func Stash(dir string) error {
	cmd := exec.Command("git", "stash", "push", "--include-untracked", "-m", "pulse-temp")
	cmd.Dir = dir
	return cmd.Run()
}

// StashPop restores stashed changes.
func StashPop(dir string) error {
	cmd := exec.Command("git", "stash", "pop")
	cmd.Dir = dir
	return cmd.Run()
}

// IsDirty returns true if there are uncommitted changes.
func IsDirty(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
