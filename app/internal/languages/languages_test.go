package languages

import (
	"testing"
)

func TestByExtension(t *testing.T) {
	tests := []struct {
		ext      string
		wantLang string
		wantOk   bool
	}{
		{".go", "Go", true},
		{".py", "Python", true},
		{".ts", "TypeScript", true},
		{".tsx", "TypeScript", true},
		{".js", "JavaScript", true},
		{".jsx", "JavaScript", true},
		{".java", "Java", true},
		{".rs", "Rust", true},
		{".rb", "Ruby", true},
		{".php", "PHP", true},
		{".kt", "Kotlin", true},
		{".swift", "Swift", true},
		{".cs", "C#", true},
		{".c", "C", true},
		{".cpp", "C++", true},
		{".sh", "Shell", true},
		{".sql", "SQL", true},
		{".tf", "Terraform", true},
		{".yaml", "YAML", true},
		{".xyz", "", false},
		{".png", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			lang, ok := ByExtension(tt.ext)
			if ok != tt.wantOk {
				t.Errorf("ByExtension(%s): ok = %v, want %v", tt.ext, ok, tt.wantOk)
			}
			if ok && lang.Name != tt.wantLang {
				t.Errorf("ByExtension(%s): name = %s, want %s", tt.ext, lang.Name, tt.wantLang)
			}
		})
	}
}

func TestSupported(t *testing.T) {
	names := Supported()
	if len(names) < 15 {
		t.Errorf("expected at least 15 supported languages, got %d", len(names))
	}
	t.Logf("Supported languages (%d): %v", len(names), names)
}
