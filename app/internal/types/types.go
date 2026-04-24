package types

// FunctionMetrics holds complexity metrics for a single function/method.
type FunctionMetrics struct {
	Name       string
	File       string
	StartLine  int
	EndLine    int
	SLOC       int
	Cyclomatic int
	Cognitive  int
	Grade      string // A-F based on maintainability
}

// FileMetrics holds aggregated metrics for a single file.
type FileMetrics struct {
	Path               string
	Language           string
	SLOC               int
	Comments           int
	Blanks             int
	TotalLines         int
	Functions          []FunctionMetrics
	MaintainabilityIdx float64
	Grade              string
}

// DuplicateBlock represents a pair of duplicated code regions.
type DuplicateBlock struct {
	FileA     string
	StartA    int
	FileB     string
	StartB    int
	LineCount int
	Snippet   string
}

// ProjectMetrics holds aggregated metrics for the whole analysis.
type ProjectMetrics struct {
	Files             []FileMetrics
	LanguageSummaries []LanguageSummary
	TotalSLOC         int
	TotalComments     int
	TotalBlanks       int
	TotalLines        int
	TotalFiles        int
	TotalFunctions    int
	AvgCyclomatic     float64
	AvgCognitive      float64
	MaxCyclomatic     FunctionMetrics
	MaxCognitive      FunctionMetrics
	MaintainabilityIdx float64
	Grade             string
	DuplicationPct    float64
	Duplicates        []DuplicateBlock
	COCOMOMonths      float64
	COCOMOPersons     float64
	COCOMOCost        float64
}

// LanguageSummary holds per-language aggregated stats.
type LanguageSummary struct {
	Language   string
	Files      int
	SLOC       int
	Comments   int
	Blanks     int
	TotalLines int
}

// TrendPoint holds metrics for a single commit in a time-series.
type TrendPoint struct {
	Hash          string
	ShortHash     string
	Author        string
	Date          string
	Message       string
	AvgCyclomatic float64
	AvgCognitive  float64
	TotalSLOC     int
	DuplicationPct float64
	Maintainability float64
	Grade         string
}

// DiffResult holds the comparison between two analysis snapshots.
type DiffResult struct {
	Base    ProjectMetrics
	Head    ProjectMetrics
	DeltaCyclomatic float64
	DeltaCognitive  float64
	DeltaSLOC       int
	DeltaDuplication float64
	DeltaMaintainability float64
	NewHotspots     []FunctionMetrics
	ImprovedFuncs   []FunctionMetrics
}

// GateConfig defines quality gate thresholds.
type GateConfig struct {
	MaxCyclomatic   int
	MaxCognitive    int
	MaxDuplication  float64
	MinMaintainability float64
	MinGrade        string
}

// GateResult holds the outcome of a quality gate check.
type GateResult struct {
	Passed  bool
	Failures []string
}

// MaintainabilityGrade returns a letter grade for a maintainability index score.
func MaintainabilityGrade(mi float64) string {
	switch {
	case mi >= 80:
		return "A"
	case mi >= 60:
		return "B"
	case mi >= 40:
		return "C"
	case mi >= 20:
		return "D"
	default:
		return "F"
	}
}

// CyclomaticGrade returns a letter grade for a cyclomatic complexity value.
func CyclomaticGrade(ccn int) string {
	switch {
	case ccn <= 5:
		return "A"
	case ccn <= 10:
		return "B"
	case ccn <= 20:
		return "C"
	case ccn <= 30:
		return "D"
	default:
		return "F"
	}
}
