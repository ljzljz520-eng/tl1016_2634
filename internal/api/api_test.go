package api

import (
	"labops/internal/flow075"
	"labops/internal/registry"
	"labops/internal/review"
	"labops/internal/schedule"
	"labops/internal/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	h := HealthHandler(s)
	r := httptest.NewRecorder()
	h(r, httptest.NewRequest("GET", "/health", nil))
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
	_ = New(registry.New(s), review.New(s), flow075.New(schedule.New(s)))
}
