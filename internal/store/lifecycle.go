package store

import (
	"fmt"
	"go.etcd.io/bbolt"
	"labops/internal/model"
)

func (s *Store) SaveWorkflow(w model.Workflow) error {
	if w.ID == "" || w.RecordID == "" {
		return fmt.Errorf("workflow identity required")
	}
	return s.PutWorkflow(w)
}
func (s *Store) SaveAttachment(a model.Attachment) error {
	if a.ID == "" || a.RecordID == "" {
		return fmt.Errorf("attachment identity required")
	}
	if a.Size < 0 {
		return fmt.Errorf("negative size")
	}
	return s.PutAttachment(a)
}
func (s *Store) Health() error {
	if s.db == nil {
		return fmt.Errorf("closed")
	}
	return s.db.View(func(*bbolt.Tx) error { return nil })
}
