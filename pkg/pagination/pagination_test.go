package pagination

import "testing"

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    Pagination
		wantPage int
		wantSize int
	}{
		{"defaults", Pagination{0, 0}, DefaultPage, DefaultPageSize},
		{"negative page", Pagination{-1, 10}, DefaultPage, 10},
		{"negative size", Pagination{2, -5}, 2, DefaultPageSize},
		{"exceeds max", Pagination{1, 2000}, 1, MaxPageSize},
		{"valid", Pagination{3, 50}, 3, 50},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Normalize()
			if tt.input.Page != tt.wantPage {
				t.Errorf("page: got %d, want %d", tt.input.Page, tt.wantPage)
			}
			if tt.input.PageSize != tt.wantSize {
				t.Errorf("pageSize: got %d, want %d", tt.input.PageSize, tt.wantSize)
			}
		})
	}
}

func TestOffset(t *testing.T) {
	p := Pagination{Page: 3, PageSize: 20}
	if p.Offset() != 40 {
		t.Errorf("expected offset 40, got %d", p.Offset())
	}
}

func TestNewPagedResult(t *testing.T) {
	items := []string{"a", "b"}
	r := NewPagedResult(items, 100, 1, 20)
	if len(r.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(r.Items))
	}
	if r.Total != 100 {
		t.Errorf("expected total 100, got %d", r.Total)
	}
}

func TestNewPagedResultNilItems(t *testing.T) {
	r := NewPagedResult[string](nil, 0, 1, 20)
	if r.Items == nil {
		t.Error("items should not be nil")
	}
	if len(r.Items) != 0 {
		t.Error("items should be empty")
	}
}
