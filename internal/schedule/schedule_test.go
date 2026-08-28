package schedule

import (
	"labops/internal/model"
	"labops/internal/store"
	"path/filepath"
	"testing"
)

func TestAssignSlots(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	s.PutRecord(model.Record{ID: "1", AssetTag: "A", Status: "active"})
	p := New(s)
	if e := p.Assign("1", []string{"am", "pm"}); e != nil {
		t.Fatal(e)
	}
	x, _ := p.Slots("1")
	if x[1] != "pm" {
		t.Fatal(x)
	}
}
