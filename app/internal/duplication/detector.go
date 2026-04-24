package duplication

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"

	"github.com/bobbydeveaux/pulse/internal/types"
)

const (
	// MinLines is the minimum number of consecutive lines to consider a duplicate.
	MinLines = 6
	// MinTokens is the minimum non-whitespace characters in a block to count.
	MinTokens = 50
)

// hashEntry stores a block hash with its source location.
type hashEntry struct {
	file      string
	startLine int
	content   string
}

// Detect finds duplicated code blocks across all provided files.
func Detect(files []string) ([]types.DuplicateBlock, float64) {
	// Build hash map of all consecutive MinLines-line blocks
	hashMap := map[string][]hashEntry{}
	totalLines := 0

	for _, path := range files {
		lines, err := readLines(path)
		if err != nil {
			continue
		}
		totalLines += len(lines)

		for i := 0; i <= len(lines)-MinLines; i++ {
			block := lines[i : i+MinLines]
			normalized := normalizeBlock(block)

			// Skip blocks that are too trivial
			if len(strings.ReplaceAll(normalized, " ", "")) < MinTokens {
				continue
			}

			h := hash(normalized)
			hashMap[h] = append(hashMap[h], hashEntry{
				file:      path,
				startLine: i + 1,
				content:   strings.Join(block, "\n"),
			})
		}
	}

	// Find duplicate blocks and merge adjacent ones
	var duplicates []types.DuplicateBlock
	duplicatedLines := map[string]bool{} // "file:line" -> true

	for _, entries := range hashMap {
		if len(entries) < 2 {
			continue
		}

		// Deduplicate entries from the same file at adjacent lines
		unique := deduplicateEntries(entries)
		if len(unique) < 2 {
			continue
		}

		// Report pairwise duplicates (first occurrence vs each other)
		for i := 1; i < len(unique); i++ {
			a := unique[0]
			b := unique[i]

			// Skip if same file and overlapping
			if a.file == b.file && abs(a.startLine-b.startLine) < MinLines {
				continue
			}

			snippet := a.content
			if len(snippet) > 200 {
				snippet = snippet[:200] + "..."
			}

			duplicates = append(duplicates, types.DuplicateBlock{
				FileA:     a.file,
				StartA:    a.startLine,
				FileB:     b.file,
				StartB:    b.startLine,
				LineCount: MinLines,
				Snippet:   snippet,
			})

			// Track duplicated lines for percentage calculation
			for j := 0; j < MinLines; j++ {
				duplicatedLines[fmt.Sprintf("%s:%d", a.file, a.startLine+j)] = true
				duplicatedLines[fmt.Sprintf("%s:%d", b.file, b.startLine+j)] = true
			}
		}
	}

	// Calculate duplication percentage
	pct := 0.0
	if totalLines > 0 {
		pct = float64(len(duplicatedLines)) / float64(totalLines) * 100.0
	}

	return duplicates, pct
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// normalizeBlock removes leading/trailing whitespace from each line
// and joins them, so formatting differences don't affect matching.
func normalizeBlock(lines []string) string {
	var parts []string
	for _, l := range lines {
		parts = append(parts, strings.TrimSpace(l))
	}
	return strings.Join(parts, "\n")
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h[:16])
}

func deduplicateEntries(entries []hashEntry) []hashEntry {
	var result []hashEntry
	seen := map[string]bool{}
	for _, e := range entries {
		key := fmt.Sprintf("%s:%d", e.file, e.startLine)
		if !seen[key] {
			seen[key] = true
			result = append(result, e)
		}
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
