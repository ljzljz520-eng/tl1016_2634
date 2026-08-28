package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"labops/internal/model"
)

type Batch struct {
	Records     []model.Record
	Workflows   []model.Workflow
	Attachments []model.Attachment
}

func (s *Store) ApplyBatch(b Batch) error {
	if s.db == nil {
		return fmt.Errorf("closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, r := range b.Records {
			raw, e := json.Marshal(r)
			if e != nil {
				return e
			}
			if e = tx.Bucket([]byte("records")).Put([]byte(r.ID), raw); e != nil {
				return e
			}
		}
		for _, w := range b.Workflows {
			raw, e := json.Marshal(w)
			if e != nil {
				return e
			}
			if e = tx.Bucket([]byte("workflows")).Put([]byte(w.ID), raw); e != nil {
				return e
			}
		}
		for _, a := range b.Attachments {
			raw, e := json.Marshal(a)
			if e != nil {
				return e
			}
			if e = tx.Bucket([]byte("attachments")).Put([]byte(a.ID), raw); e != nil {
				return e
			}
		}
		return nil
	})
}
func (b Batch) Count() int  { return len(b.Records) + len(b.Workflows) + len(b.Attachments) }
func (b Batch) Empty() bool { return b.Count() == 0 }
func NewBatch() *Batch {
	return &Batch{Records: []model.Record{}, Workflows: []model.Workflow{}, Attachments: []model.Attachment{}}
}
func (b *Batch) AddRecord(r model.Record)         { b.Records = append(b.Records, r) }
func (b *Batch) AddWorkflow(w model.Workflow)     { b.Workflows = append(b.Workflows, w) }
func (b *Batch) AddAttachment(a model.Attachment) { b.Attachments = append(b.Attachments, a) }
func (b Batch) RecordIDs() []string {
	out := []string{}
	for _, r := range b.Records {
		out = append(out, r.ID)
	}
	return out
}
