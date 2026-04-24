package languages

// Language defines how to parse a source file for a given language.
type Language struct {
	Name           string
	Extensions     []string
	LineComment    string   // e.g. "//"
	BlockStart     string   // e.g. "/*"
	BlockEnd       string   // e.g. "*/"
	FuncKeywords   []string // e.g. ["func"] for Go, ["def"] for Python
	BranchKeywords []string // keywords that add cyclomatic complexity
	NestKeywords   []string // keywords that increase cognitive nesting
}

var registry = map[string]Language{}

func init() {
	for _, lang := range All {
		for _, ext := range lang.Extensions {
			registry[ext] = lang
		}
	}
}

// ByExtension returns the language for a file extension (including the dot).
func ByExtension(ext string) (Language, bool) {
	l, ok := registry[ext]
	return l, ok
}

// All is the complete list of supported languages.
var All = []Language{
	{
		Name:           "Go",
		Extensions:     []string{".go"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"func"},
		BranchKeywords: []string{"if", "else if", "for", "case", "&&", "||"},
		NestKeywords:   []string{"if", "for", "switch", "select"},
	},
	{
		Name:           "Python",
		Extensions:     []string{".py", ".pyw"},
		LineComment:    "#",
		BlockStart:     `"""`,
		BlockEnd:       `"""`,
		FuncKeywords:   []string{"def", "class"},
		BranchKeywords: []string{"if", "elif", "for", "while", "except", "and", "or"},
		NestKeywords:   []string{"if", "for", "while", "try", "with"},
	},
	{
		Name:           "JavaScript",
		Extensions:     []string{".js", ".jsx", ".mjs", ".cjs"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"function", "=>"},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||", "??"},
		NestKeywords:   []string{"if", "for", "while", "switch", "try"},
	},
	{
		Name:           "TypeScript",
		Extensions:     []string{".ts", ".tsx"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"function", "=>"},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||", "??"},
		NestKeywords:   []string{"if", "for", "while", "switch", "try"},
	},
	{
		Name:           "Java",
		Extensions:     []string{".java"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"void", "public", "private", "protected", "static"},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "switch", "try"},
	},
	{
		Name:           "Rust",
		Extensions:     []string{".rs"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"fn"},
		BranchKeywords: []string{"if", "else if", "for", "while", "match", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "match", "loop"},
	},
	{
		Name:           "C",
		Extensions:     []string{".c", ".h"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "switch"},
	},
	{
		Name:           "C++",
		Extensions:     []string{".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "switch", "try"},
	},
	{
		Name:           "C#",
		Extensions:     []string{".cs"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"void", "public", "private", "protected", "static", "async"},
		BranchKeywords: []string{"if", "else if", "for", "foreach", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "foreach", "while", "switch", "try"},
	},
	{
		Name:           "Ruby",
		Extensions:     []string{".rb"},
		LineComment:    "#",
		BlockStart:     "=begin",
		BlockEnd:       "=end",
		FuncKeywords:   []string{"def", "class"},
		BranchKeywords: []string{"if", "elsif", "unless", "when", "while", "until", "rescue", "&&", "||"},
		NestKeywords:   []string{"if", "unless", "while", "until", "begin"},
	},
	{
		Name:           "PHP",
		Extensions:     []string{".php"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"function"},
		BranchKeywords: []string{"if", "elseif", "for", "foreach", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "foreach", "while", "switch", "try"},
	},
	{
		Name:           "Kotlin",
		Extensions:     []string{".kt", ".kts"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"fun"},
		BranchKeywords: []string{"if", "else if", "for", "while", "when", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "when", "try"},
	},
	{
		Name:           "Swift",
		Extensions:     []string{".swift"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"func"},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "switch", "do"},
	},
	{
		Name:           "Scala",
		Extensions:     []string{".scala"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"def"},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "match", "try"},
	},
	{
		Name:           "Lua",
		Extensions:     []string{".lua"},
		LineComment:    "--",
		BlockStart:     "--[[",
		BlockEnd:       "]]",
		FuncKeywords:   []string{"function"},
		BranchKeywords: []string{"if", "elseif", "for", "while", "and", "or"},
		NestKeywords:   []string{"if", "for", "while", "repeat"},
	},
	{
		Name:           "Shell",
		Extensions:     []string{".sh", ".bash", ".zsh"},
		LineComment:    "#",
		BlockStart:     "",
		BlockEnd:       "",
		FuncKeywords:   []string{"function"},
		BranchKeywords: []string{"if", "elif", "for", "while", "case", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "case"},
	},
	{
		Name:           "YAML",
		Extensions:     []string{".yaml", ".yml"},
		LineComment:    "#",
		BlockStart:     "",
		BlockEnd:       "",
		FuncKeywords:   []string{},
		BranchKeywords: []string{},
		NestKeywords:   []string{},
	},
	{
		Name:           "TOML",
		Extensions:     []string{".toml"},
		LineComment:    "#",
		BlockStart:     "",
		BlockEnd:       "",
		FuncKeywords:   []string{},
		BranchKeywords: []string{},
		NestKeywords:   []string{},
	},
	{
		Name:           "SQL",
		Extensions:     []string{".sql"},
		LineComment:    "--",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{},
		BranchKeywords: []string{"CASE", "WHEN", "IF", "ELSEIF"},
		NestKeywords:   []string{"CASE", "IF", "BEGIN"},
	},
	{
		Name:           "Terraform",
		Extensions:     []string{".tf"},
		LineComment:    "#",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{},
		BranchKeywords: []string{},
		NestKeywords:   []string{},
	},
	{
		Name:           "Dart",
		Extensions:     []string{".dart"},
		LineComment:    "//",
		BlockStart:     "/*",
		BlockEnd:       "*/",
		FuncKeywords:   []string{"void", "Future", "Stream"},
		BranchKeywords: []string{"if", "else if", "for", "while", "case", "catch", "&&", "||"},
		NestKeywords:   []string{"if", "for", "while", "switch", "try"},
	},
	{
		Name:           "Elixir",
		Extensions:     []string{".ex", ".exs"},
		LineComment:    "#",
		BlockStart:     "",
		BlockEnd:       "",
		FuncKeywords:   []string{"def", "defp", "defmodule"},
		BranchKeywords: []string{"if", "cond", "case", "with", "and", "or"},
		NestKeywords:   []string{"if", "cond", "case", "with", "try"},
	},
}

// Supported returns a sorted list of supported language names.
func Supported() []string {
	seen := map[string]bool{}
	var names []string
	for _, l := range All {
		if !seen[l.Name] {
			seen[l.Name] = true
			names = append(names, l.Name)
		}
	}
	return names
}
