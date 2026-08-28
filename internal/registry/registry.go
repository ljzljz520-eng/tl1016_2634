package registry

import (
	"fmt"
	"labops/internal/model"
	"labops/internal/store"
)

type Service struct{ Store *store.Store }

func New(s *store.Store) *Service { return &Service{Store: s} }
func (s *Service) Register(r model.Record) error {
	r.Status = model.NormalizeStatus(r.Status)
	if e := r.Validate(); e != nil {
		return e
	}
	if _, e := s.Store.GetRecord(r.ID); e == nil {
		return fmt.Errorf("record exists")
	}
	return s.Store.PutRecord(r)
}
func (s *Service) Get(id string) (model.Record, error) { return s.Store.GetRecord(id) }
func (s *Service) Update(r model.Record) error {
	if e := r.Validate(); e != nil {
		return e
	}
	old, e := s.Get(r.ID)
	if e != nil {
		return e
	}
	if old.IsArchived() {
		return fmt.Errorf("archived record")
	}
	r.Version = old.Version + 1
	return s.Store.PutRecord(r)
}
func (s *Service) Archive(id string) error {
	r, e := s.Get(id)
	if e != nil {
		return e
	}
	r.Status = "archived"
	return s.Store.PutRecord(r)
}
