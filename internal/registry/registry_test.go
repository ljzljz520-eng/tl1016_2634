package registry

import (
	"labops/internal/model"
	"labops/internal/store"
	"path/filepath"
	"testing"
)

func TestRegistryRegister(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := New(s)
	if e := r.Register(model.Record{ID: "1", AssetTag: "A", Status: "draft"}); e != nil {
		t.Fatal(e)
	}
	if e := r.Register(model.Record{ID: "1", AssetTag: "A", Status: "draft"}); e == nil {
		t.Fatal("duplicate accepted")
	}
}
