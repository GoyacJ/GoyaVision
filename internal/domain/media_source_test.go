package domain

import (
	"regexp"
	"testing"
)

func TestGeneratePathName(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"ascii", "camera1"},
		{"with spaces", "摄像头 1"},
		{"special chars", "a/b?c#d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GeneratePathName(tt.in)
			matched, _ := regexp.MatchString(`^live/[a-zA-Z0-9_-]+-[a-f0-9]{8}$`, got)
			if !matched {
				t.Errorf("GeneratePathName(%q) = %q, want match live/{slug}-{8hex}", tt.in, got)
			}
		})
	}
}

func TestGeneratePathName_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		s := GeneratePathName("same")
		if seen[s] {
			t.Errorf("GeneratePathName produced duplicate %q", s)
		}
		seen[s] = true
	}
}
