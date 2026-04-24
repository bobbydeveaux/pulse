package sloc

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobbydeveaux/pulse/internal/languages"
)

// Counts holds line counts for a single file.
type Counts struct {
	SLOC       int
	Comments   int
	Blanks     int
	TotalLines int
	Language   string
}

// CountFile counts lines in a single file.
func CountFile(path string) (Counts, error) {
	ext := filepath.Ext(path)
	lang, ok := languages.ByExtension(ext)
	if !ok {
		return Counts{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return Counts{}, err
	}
	defer f.Close()

	var c Counts
	c.Language = lang.Name

	inBlock := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		c.TotalLines++
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			c.Blanks++
			continue
		}

		// Handle block comments
		if inBlock {
			c.Comments++
			if lang.BlockEnd != "" && strings.Contains(trimmed, lang.BlockEnd) {
				inBlock = false
			}
			continue
		}

		// Check for block comment start
		if lang.BlockStart != "" && strings.Contains(trimmed, lang.BlockStart) {
			// Special case: single-line block comment /* ... */
			if lang.BlockEnd != "" && strings.Contains(trimmed, lang.BlockEnd) {
				afterStart := strings.Index(trimmed, lang.BlockStart)
				afterEnd := strings.Index(trimmed, lang.BlockEnd)
				if afterEnd > afterStart {
					// Check if there's code before or after the comment
					before := strings.TrimSpace(trimmed[:afterStart])
					after := strings.TrimSpace(trimmed[afterEnd+len(lang.BlockEnd):])
					if before == "" && after == "" {
						c.Comments++
					} else {
						c.SLOC++
					}
					continue
				}
			}
			// Multi-line block comment starts
			// Check if there's code before the comment
			idx := strings.Index(trimmed, lang.BlockStart)
			before := strings.TrimSpace(trimmed[:idx])
			if before == "" {
				c.Comments++
			} else {
				c.SLOC++
			}
			if lang.BlockEnd == "" || !strings.Contains(trimmed, lang.BlockEnd) {
				inBlock = true
			}
			continue
		}

		// Check for line comment
		if lang.LineComment != "" && strings.HasPrefix(trimmed, lang.LineComment) {
			c.Comments++
			continue
		}

		c.SLOC++
	}

	return c, scanner.Err()
}
