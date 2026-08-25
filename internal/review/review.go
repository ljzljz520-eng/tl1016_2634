package review

import (
	"fmt"
	"labops/internal/model"
	"labops/internal/store"
)

type Service struct{ Store *store.Store }

func New(s *store.Store) *Service { return &Service{Store: s} }
func (s *Service) Submit(id, actor string) error {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status != "draft" {
		return fmt.Errorf("not draft")
	}
	r.Status = "pending"
	if e = s.Store.PutRecord(r); e != nil {
		return e
	}
	return s.audit(id, "submit", actor)
}
func (s *Service) Approve(id, actor string) error {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status != "pending" {
		return fmt.Errorf("not pending")
	}
	r.Status = "active"
	if e = s.Store.PutRecord(r); e != nil {
		return e
	}
	return s.audit(id, "approve", actor)
}
func (s *Service) Reject(id, actor, reason string) error {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status != "pending" {
		return fmt.Errorf("not pending")
	}
	r.Status = "draft"
	if e = s.Store.PutRecord(r); e != nil {
		return e
	}
	return s.audit(id, "reject", actor+":"+reason)
}
func (s *Service) audit(id, action, actor string) error {
	n, _ := s.Store.Count("audits")
	return s.Store.PutAudit(model.AuditEvent{ID: fmt.Sprintf("%s-%d", id, n+1), RecordID: id, Action: action, Actor: actor, Sequence: n + 1})
}
