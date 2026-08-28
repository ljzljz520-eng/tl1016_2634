package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"labops/internal/model"
	"path/filepath"
	"sync"
)

var buckets = [][]byte{[]byte("records"), []byte("audits"), []byte("workflows"), []byte("attachments")}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, e := bbolt.Open(filepath.Clean(path), 0600, nil)
	if e != nil {
		return nil, e
	}
	s := &Store{db: db}
	e = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, x := tx.CreateBucketIfNotExists(b); x != nil {
				return x
			}
		}
		return nil
	})
	if e != nil {
		db.Close()
		return nil, e
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	e := s.db.Close()
	s.db = nil
	return e
}
func (s *Store) put(bucket, key string, v any) error {
	raw, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), raw) })
}
func (s *Store) get(bucket, key string, out any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if v == nil {
			return fmt.Errorf("not found")
		}
		return json.Unmarshal(v, out)
	})
}
func (s *Store) PutRecord(r model.Record) error { return s.put("records", r.ID, r) }
func (s *Store) GetRecord(id string) (model.Record, error) {
	var r model.Record
	e := s.get("records", id, &r)
	return r, e
}
func (s *Store) PutWorkflow(w model.Workflow) error { return s.put("workflows", w.ID, w) }
func (s *Store) GetWorkflow(id string) (model.Workflow, error) {
	var w model.Workflow
	e := s.get("workflows", id, &w)
	return w, e
}
func (s *Store) PutAttachment(a model.Attachment) error { return s.put("attachments", a.ID, a) }
func (s *Store) PutAudit(a model.AuditEvent) error      { return s.put("audits", a.ID, a) }
func (s *Store) ListRecords() ([]model.Record, error) {
	out := []model.Record{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r model.Record
			if x := json.Unmarshal(v, &r); x != nil {
				return x
			}
			out = append(out, r)
			return nil
		})
	})
	return out, e
}
