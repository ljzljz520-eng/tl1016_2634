package flow075

import (
	"labops/internal/model"
	"labops/internal/schedule"
	"labops/internal/store"
	"path/filepath"
	"testing"
)

func Test1016BusinessRegression(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	s.PutRecord(model.Record{ID: "1", AssetTag: "A", Status: "active"})
	p := New(schedule.New(s))
	if e := p.Operate("1", []string{"first", "second"}); e != nil {
		t.Fatal(e)
	}
	got, _ := p.Ordered("1")
	if got[0] != "first" || got[1] != "second" {
		t.Fatalf("unexpected order: %v", got)
	}
}
