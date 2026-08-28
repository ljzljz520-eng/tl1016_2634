package integration

import (
	"bytes"
	"labops/internal/flow075"
	"labops/internal/model"
	"labops/internal/registry"
	"labops/internal/review"
	"labops/internal/schedule"
	"labops/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := registry.New(s)
	v := review.New(s)
	r.Register(model.Record{ID: "r", AssetTag: "A", Status: "draft"})
	v.Submit("r", "auditor")
	v.Approve("r", "auditor")
	r.Archive("r")
	x, _ := r.Get("r")
	if !x.IsArchived() {
		t.Fatal(x.Status)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	r := registry.New(s)
	r.Register(model.Record{ID: "r", AssetTag: "A", Name: "centrifuge", Status: "active"})
	p := schedule.New(s)
	f := flow075.New(p)
	f.Operate("r", []string{"08:00", "12:00"})
	q, _ := s.Search(model.Query{Text: "centrifuge"})
	if len(q.Records) != 1 {
		t.Fatal(q)
	}
}
func TestWorkflowImportReport(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "db"))
	defer s.Close()
	n, errs := registry.New(s).Import(bytes.NewBufferString("id,tag,name,location\n1,A,scope,lab\n2,B\n"))
	if n != 1 || len(errs) != 1 {
		t.Fatalf("%d %v", n, errs)
	}
}
