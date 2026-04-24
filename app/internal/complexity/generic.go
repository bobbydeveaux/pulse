package complexity

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/bobbydeveaux/pulse/internal/languages"
	"github.com/bobbydeveaux/pulse/internal/types"
)

// funcPattern matches common function declarations across languages.
var funcPatterns = map[string]*regexp.Regexp{
	"Python":     regexp.MustCompile(`^\s*(async\s+)?def\s+(\w+)\s*\(`),
	"JavaScript": regexp.MustCompile(`(?:function\s+(\w+)\s*\(|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s+)?(?:function|\([^)]*\)\s*=>|\w+\s*=>))`),
	"TypeScript": regexp.MustCompile(`(?:function\s+(\w+)\s*[\(<]|(?:const|let|var)\s+(\w+)\s*(?::\s*\S+\s*)?=\s*(?:async\s+)?(?:function|\([^)]*\)\s*=>|\w+\s*=>))`),
	"Java":       regexp.MustCompile(`(?:public|private|protected|static|\s)+\s+\w[\w<>\[\],\s]*\s+(\w+)\s*\(`),
	"Rust":       regexp.MustCompile(`\s*(?:pub\s+)?(?:async\s+)?fn\s+(\w+)`),
	"C":          regexp.MustCompile(`^\w[\w\s\*]*\s+(\w+)\s*\([^;]*$`),
	"C++":        regexp.MustCompile(`^\w[\w\s\*:&<>]*\s+(\w+)\s*\([^;]*$`),
	"C#":         regexp.MustCompile(`(?:public|private|protected|internal|static|\s)+\s+\w[\w<>\[\],\s]*\s+(\w+)\s*\(`),
	"Ruby":       regexp.MustCompile(`^\s*def\s+(\w+[\?!]?)`),
	"PHP":        regexp.MustCompile(`(?:public|private|protected|static|\s)*function\s+(\w+)\s*\(`),
	"Kotlin":     regexp.MustCompile(`\s*(?:fun|suspend\s+fun)\s+(\w+)`),
	"Swift":      regexp.MustCompile(`\s*(?:func)\s+(\w+)`),
	"Scala":      regexp.MustCompile(`\s*def\s+(\w+)`),
	"Lua":        regexp.MustCompile(`(?:local\s+)?function\s+(\w[\w.]*)\s*\(`),
	"Shell":      regexp.MustCompile(`(?:function\s+(\w+)|(\w+)\s*\(\s*\))`),
	"Dart":       regexp.MustCompile(`\s*(?:Future|Stream|void|int|double|String|bool|dynamic|\w+)\s+(\w+)\s*\(`),
	"Elixir":     regexp.MustCompile(`\s*(?:def|defp)\s+(\w+[\?!]?)`),
}

// AnalyzeGenericFile performs regex-based complexity analysis for non-Go languages.
func AnalyzeGenericFile(path string, lang languages.Language) ([]types.FunctionMetrics, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	funcs := findFunctions(lines, lang)
	return funcs, nil
}

type funcRange struct {
	name      string
	startLine int
	endLine   int
}

func findFunctions(lines []string, lang languages.Language) []types.FunctionMetrics {
	pattern, ok := funcPatterns[lang.Name]
	if !ok {
		return nil
	}

	var ranges []funcRange

	// Find function start lines
	for i, line := range lines {
		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		name := ""
		for _, m := range matches[1:] {
			if m != "" {
				name = m
				break
			}
		}
		if name == "" {
			continue
		}

		ranges = append(ranges, funcRange{
			name:      name,
			startLine: i + 1,
		})
	}

	// Determine function end lines using brace/indent counting
	for i := range ranges {
		if i+1 < len(ranges) {
			ranges[i].endLine = ranges[i+1].startLine - 1
		} else {
			ranges[i].endLine = len(lines)
		}
	}

	// For brace-based languages, try to find the actual end
	if isBraceLang(lang) {
		for i := range ranges {
			end := findBraceEnd(lines, ranges[i].startLine-1)
			if end > 0 {
				ranges[i].endLine = end
			}
		}
	}

	// Calculate complexity for each function
	var results []types.FunctionMetrics
	for _, fr := range ranges {
		funcLines := lines[fr.startLine-1 : fr.endLine]
		ccn := cyclomaticGeneric(funcLines, lang)
		cog := cognitiveGeneric(funcLines, lang)

		results = append(results, types.FunctionMetrics{
			Name:       fr.name,
			StartLine:  fr.startLine,
			EndLine:    fr.endLine,
			SLOC:       fr.endLine - fr.startLine + 1,
			Cyclomatic: ccn,
			Cognitive:  cog,
			Grade:      types.CyclomaticGrade(ccn),
		})
	}

	return results
}

func isBraceLang(lang languages.Language) bool {
	switch lang.Name {
	case "Python", "Ruby", "Elixir", "Lua", "Shell", "YAML", "TOML":
		return false
	}
	return true
}

func findBraceEnd(lines []string, startIdx int) int {
	depth := 0
	started := false
	for i := startIdx; i < len(lines); i++ {
		for _, ch := range lines[i] {
			if ch == '{' {
				depth++
				started = true
			} else if ch == '}' {
				depth--
				if started && depth == 0 {
					return i + 1 // 1-indexed
				}
			}
		}
	}
	return 0
}

// cyclomaticGeneric counts decision points using keyword matching.
func cyclomaticGeneric(lines []string, lang languages.Language) int {
	ccn := 1
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip comments
		if lang.BlockStart != "" && strings.Contains(trimmed, lang.BlockStart) {
			inBlock = true
		}
		if inBlock {
			if lang.BlockEnd != "" && strings.Contains(trimmed, lang.BlockEnd) {
				inBlock = false
			}
			continue
		}
		if lang.LineComment != "" && strings.HasPrefix(trimmed, lang.LineComment) {
			continue
		}

		for _, kw := range lang.BranchKeywords {
			if kw == "&&" || kw == "||" || kw == "??" || kw == "and" || kw == "or" {
				ccn += strings.Count(trimmed, kw)
			} else {
				ccn += countKeyword(trimmed, kw)
			}
		}
	}

	return ccn
}

// cognitiveGeneric estimates cognitive complexity using nesting depth.
func cognitiveGeneric(lines []string, lang languages.Language) int {
	cog := 0
	nesting := 0
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if lang.BlockStart != "" && strings.Contains(trimmed, lang.BlockStart) {
			inBlock = true
		}
		if inBlock {
			if lang.BlockEnd != "" && strings.Contains(trimmed, lang.BlockEnd) {
				inBlock = false
			}
			continue
		}
		if lang.LineComment != "" && strings.HasPrefix(trimmed, lang.LineComment) {
			continue
		}

		// Count nesting-increasing keywords
		for _, kw := range lang.NestKeywords {
			count := countKeyword(trimmed, kw)
			if count > 0 {
				cog += count * (1 + nesting)
				nesting += count
			}
		}

		// Boolean operators add +1 each (no nesting penalty)
		for _, op := range []string{"&&", "||", "??", " and ", " or "} {
			cog += strings.Count(trimmed, op)
		}

		// Track brace-based nesting
		if isBraceLang(lang) {
			for _, ch := range trimmed {
				if ch == '{' {
					// We already incremented for keywords; braces just confirm
				} else if ch == '}' {
					if nesting > 0 {
						nesting--
					}
				}
			}
		}
	}

	return cog
}

// countKeyword counts occurrences of a keyword as a whole word in a line.
func countKeyword(line, keyword string) int {
	count := 0
	idx := 0
	for {
		pos := strings.Index(line[idx:], keyword)
		if pos == -1 {
			break
		}
		absPos := idx + pos
		end := absPos + len(keyword)

		// Check word boundaries
		beforeOk := absPos == 0 || !isIdentChar(line[absPos-1])
		afterOk := end >= len(line) || !isIdentChar(line[end])

		if beforeOk && afterOk {
			count++
		}
		idx = end
	}
	return count
}

func isIdentChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}
