package schedule

import (
	"fmt"
	"labops/internal/store"
)

type Planner struct{ Store *store.Store }

func New(s *store.Store) *Planner { return &Planner{Store: s} }
func (p *Planner) Assign(id string, slots []string) error {
	r, e := p.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Status != "active" && r.Status != "maintenance" {
		return fmt.Errorf("inactive")
	}
	r.Slots = append([]string(nil), slots...)
	r.Status = "maintenance"
	return p.Store.PutRecord(r)
}
func (p *Planner) Release(id string) error {
	r, e := p.Store.GetRecord(id)
	if e != nil {
		return e
	}
	r.Status = "active"
	return p.Store.PutRecord(r)
}
func (p *Planner) Slots(id string) ([]string, error) {
	r, e := p.Store.GetRecord(id)
	return r.Slots, e
}
func ValidateSlots(slots []string) error {
	seen := map[string]bool{}
	for _, v := range slots {
		if v == "" {
			return fmt.Errorf("empty slot")
		}
		if seen[v] {
			return fmt.Errorf("duplicate slot")
		}
		seen[v] = true
	}
	return nil
}
