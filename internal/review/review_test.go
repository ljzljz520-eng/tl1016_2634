package review

import (
	"labops/internal/model"
	"labops/internal/store"
	"path/filepath"
	"testing"
)

func TestReviewTransitions(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	s.PutRecord(model.Record{ID: "1", AssetTag: "A", Status: "draft"})
	v := New(s)
	if e := v.Submit("1", "u"); e != nil {
		t.Fatal(e)
	}
	if e := v.Approve("1", "u"); e != nil {
		t.Fatal(e)
	}
}
