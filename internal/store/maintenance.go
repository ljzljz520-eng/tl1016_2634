package store

import (
	"fmt"
	"go.etcd.io/bbolt"
	"labops/internal/model"
)

func (s *Store) ReplaceRecord(r model.Record) error {
	if s.db == nil {
		return fmt.Errorf("closed")
	}
	return s.PutRecord(r)
}
func (s *Store) HasRecord(id string) bool   { _, e := s.GetRecord(id); return e == nil }
func (s *Store) HasWorkflow(id string) bool { _, e := s.GetWorkflow(id); return e == nil }
func (s *Store) SaveAll(r model.Record, w model.Workflow, a model.Attachment) error {
	if e := s.PutRecord(r); e != nil {
		return e
	}
	if e := s.PutWorkflow(w); e != nil {
		return e
	}
	return s.PutAttachment(a)
}
func (s *Store) DeleteRecord(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("records")).Delete([]byte(id)) })
}
func (s *Store) DeleteWorkflow(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("workflows")).Delete([]byte(id)) })
}
func (s *Store) DeleteAttachment(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("attachments")).Delete([]byte(id)) })
}
func (s *Store) DeleteAudit(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte("audits")).Delete([]byte(id)) })
}
func (s *Store) BucketCounts() map[string]int {
	out := map[string]int{}
	for _, name := range []string{"records", "audits", "workflows", "attachments"} {
		n, _ := s.Count(name)
		out[name] = n
	}
	return out
}
func (s *Store) RecordIDs() []string {
	rows, _ := s.ListRecords()
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, r.ID)
	}
	return out
}
func (s *Store) Upsert(r model.Record) error {
	if e := r.Validate(); e != nil {
		return e
	}
	return s.PutRecord(r)
}
func (s *Store) UpdateStatus(id, status string) error {
	r, e := s.GetRecord(id)
	if e != nil {
		return e
	}
	if status == "" {
		return fmt.Errorf("status required")
	}
	r.Status = status
	return s.PutRecord(r)
}
func (s *Store) UpdateOwner(id, owner string) error {
	r, e := s.GetRecord(id)
	if e != nil {
		return e
	}
	r.Owner = owner
	return s.PutRecord(r)
}
func (s *Store) UpdateLocation(id, location string) error {
	r, e := s.GetRecord(id)
	if e != nil {
		return e
	}
	r.Location = location
	return s.PutRecord(r)
}
