package analyzer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bobbydeveaux/pulse/internal/cocomo"
	"github.com/bobbydeveaux/pulse/internal/complexity"
	"github.com/bobbydeveaux/pulse/internal/duplication"
	"github.com/bobbydeveaux/pulse/internal/languages"
	"github.com/bobbydeveaux/pulse/internal/maintainability"
	"github.com/bobbydeveaux/pulse/internal/sloc"
	"github.com/bobbydeveaux/pulse/internal/types"
)

// skipDirs contains directories to always skip.
var skipDirs = map[string]bool{
	"node_modules": true, "vendor": true, ".git": true,
	"dist": true, "build": true, "__pycache__": true,
	"venv": true, ".venv": true, ".next": true, ".nuxt": true,
	"coverage": true, ".idea": true, ".vscode": true,
	"target": true, "bin": true, "obj": true,
}

// skipExts contains file extensions to skip.
var skipExts = map[string]bool{
	".min.js": true, ".min.css": true,
	".lock": true, ".sum": true,
	".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
	".ico": true, ".svg": true, ".woff": true, ".woff2": true,
	".ttf": true, ".eot": true, ".mp3": true, ".mp4": true,
	".zip": true, ".tar": true, ".gz": true,
	".pdf": true, ".exe": true, ".dll": true, ".so": true,
	".pyc": true, ".class": true, ".o": true, ".a": true,
}

// Analyze runs the full analysis pipeline on the given path.
func Analyze(root string) (*types.ProjectMetrics, error) {
	files, err := collectFiles(root)
	if err != nil {
		return nil, err
	}

	pm := &types.ProjectMetrics{}
	langMap := map[string]*types.LanguageSummary{}

	for _, path := range files {
		// SLOC counting
		counts, err := sloc.CountFile(path)
		if err != nil || counts.Language == "" {
			continue
		}

		// Complexity analysis
		funcs, err := complexity.AnalyzeFile(path)
		if err != nil {
			funcs = nil
		}

		// Set file path on each function
		for i := range funcs {
			if funcs[i].File == "" {
				funcs[i].File = path
			}
		}

		// Compute file-level maintainability from per-function MI
		var funcInputs []maintainability.FuncInput
		for _, f := range funcs {
			funcInputs = append(funcInputs, maintainability.FuncInput{
				SLOC:       f.SLOC,
				Cyclomatic: f.Cyclomatic,
			})
		}
		mi := maintainability.FileIndex(funcInputs, counts.SLOC)

		fm := types.FileMetrics{
			Path:               path,
			Language:           counts.Language,
			SLOC:               counts.SLOC,
			Comments:           counts.Comments,
			Blanks:             counts.Blanks,
			TotalLines:         counts.TotalLines,
			Functions:          funcs,
			MaintainabilityIdx: mi,
			Grade:              types.MaintainabilityGrade(mi),
		}

		pm.Files = append(pm.Files, fm)
		pm.TotalSLOC += counts.SLOC
		pm.TotalComments += counts.Comments
		pm.TotalBlanks += counts.Blanks
		pm.TotalLines += counts.TotalLines
		pm.TotalFiles++
		pm.TotalFunctions += len(funcs)

		// Aggregate per-language
		ls, ok := langMap[counts.Language]
		if !ok {
			ls = &types.LanguageSummary{Language: counts.Language}
			langMap[counts.Language] = ls
		}
		ls.Files++
		ls.SLOC += counts.SLOC
		ls.Comments += counts.Comments
		ls.Blanks += counts.Blanks
		ls.TotalLines += counts.TotalLines
	}

	// Language summaries sorted by SLOC
	for _, ls := range langMap {
		pm.LanguageSummaries = append(pm.LanguageSummaries, *ls)
	}
	sort.Slice(pm.LanguageSummaries, func(i, j int) bool {
		return pm.LanguageSummaries[i].SLOC > pm.LanguageSummaries[j].SLOC
	})

	// Aggregate function metrics
	totalCCN := 0
	totalCog := 0
	funcCount := 0
	for _, f := range pm.Files {
		for _, fn := range f.Functions {
			totalCCN += fn.Cyclomatic
			totalCog += fn.Cognitive
			funcCount++

			if fn.Cyclomatic > pm.MaxCyclomatic.Cyclomatic {
				pm.MaxCyclomatic = fn
			}
			if fn.Cognitive > pm.MaxCognitive.Cognitive {
				pm.MaxCognitive = fn
			}
		}
	}
	if funcCount > 0 {
		pm.AvgCyclomatic = float64(totalCCN) / float64(funcCount)
		pm.AvgCognitive = float64(totalCog) / float64(funcCount)
	}

	// Project-level maintainability = average of per-file MI scores
	if len(pm.Files) > 0 {
		totalMI := 0.0
		for _, f := range pm.Files {
			totalMI += f.MaintainabilityIdx
		}
		pm.MaintainabilityIdx = totalMI / float64(len(pm.Files))
	}
	pm.Grade = types.MaintainabilityGrade(pm.MaintainabilityIdx)

	// Duplication detection
	pm.Duplicates, pm.DuplicationPct = duplication.Detect(files)

	// COCOMO estimation
	est := cocomo.Calculate(pm.TotalSLOC)
	pm.COCOMOMonths = est.PersonMonths
	pm.COCOMOPersons = est.Persons
	pm.COCOMOCost = est.Cost

	return pm, nil
}

func collectFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}

		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if skipExts[ext] {
			return nil
		}

		// Check if it has a suffix like .min.js
		base := filepath.Base(path)
		if strings.Contains(base, ".min.") {
			return nil
		}

		// Only include files we know how to analyze
		if _, ok := languages.ByExtension(ext); ok {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
