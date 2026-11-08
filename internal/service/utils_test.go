package service

import (
	"context"
	"testing"
)

func TestHasNextPage(t *testing.T) {
	service := &svc{limit: 10}
	tests := []struct {
		totalRows int
		page      int
		want      bool
	}{
		{totalRows: 0, page: 1, want: false},
		{totalRows: 10, page: 1, want: false},
		{totalRows: 11, page: 1, want: true},
		{totalRows: 15, page: 1, want: true},
		{totalRows: 20, page: 2, want: false},
		{totalRows: 21, page: 2, want: true},
	}
	for _, tt := range tests {
		if got := service.HasNextPage(context.Background(), tt.totalRows, tt.page); got != tt.want {
			t.Errorf("HasNextPage(total=%d, page=%d) = %v, want %v", tt.totalRows, tt.page, got, tt.want)
		}
	}
}

func TestTotalPages(t *testing.T) {
	service := &svc{limit: 10}
	tests := map[int]int{0: 0, 1: 1, 10: 1, 11: 2, 20: 2, 21: 3}
	for totalRows, want := range tests {
		if got := service.TotalPages(totalRows); got != want {
			t.Errorf("TotalPages(%d) = %d, want %d", totalRows, got, want)
		}
	}
}
