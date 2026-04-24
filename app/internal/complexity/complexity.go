package complexity

import (
	"path/filepath"

	"github.com/bobbydeveaux/pulse/internal/languages"
	"github.com/bobbydeveaux/pulse/internal/types"
)

// AnalyzeFile dispatches to the appropriate analyzer based on file extension.
func AnalyzeFile(path string) ([]types.FunctionMetrics, error) {
	ext := filepath.Ext(path)
	lang, ok := languages.ByExtension(ext)
	if !ok {
		return nil, nil
	}

	if lang.Name == "Go" {
		return AnalyzeGoFile(path)
	}

	return AnalyzeGenericFile(path, lang)
}
