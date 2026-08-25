package store

import (
	"labops/internal/model"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	if e = s.PutRecord(model.Record{ID: "r1", AssetTag: "A", Status: "draft"}); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	r, e := s.GetRecord("r1")
	if e != nil || r.AssetTag != "A" {
		t.Fatalf("%v %#v", e, r)
	}
}
