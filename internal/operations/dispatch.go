package operations

import (
	"fmt"
	"sort"
)

type Dispatch struct {
	ID, RecordID, Assignee, State string
	Priority                      int
}
type Board struct{ rows map[string]Dispatch }

func NewBoard() *Board { return &Board{rows: map[string]Dispatch{}} }
func (b *Board) Add(d Dispatch) error {
	if d.ID == "" || d.RecordID == "" {
		return fmt.Errorf("identity required")
	}
	if _, ok := b.rows[d.ID]; ok {
		return fmt.Errorf("dispatch exists")
	}
	if d.State == "" {
		d.State = "queued"
	}
	b.rows[d.ID] = d
	return nil
}
func (b *Board) Get(id string) (Dispatch, error) {
	d, ok := b.rows[id]
	if !ok {
		return Dispatch{}, fmt.Errorf("dispatch missing")
	}
	return d, nil
}
func (b *Board) Assign(id, who string) error {
	d, e := b.Get(id)
	if e != nil {
		return e
	}
	if d.State != "queued" {
		return fmt.Errorf("not queued")
	}
	if who == "" {
		return fmt.Errorf("assignee required")
	}
	d.Assignee = who
	d.State = "assigned"
	b.rows[id] = d
	return nil
}
func (b *Board) Begin(id string) error {
	d, e := b.Get(id)
	if e != nil {
		return e
	}
	if d.State != "assigned" {
		return fmt.Errorf("not assigned")
	}
	d.State = "running"
	b.rows[id] = d
	return nil
}
func (b *Board) Finish(id string) error {
	d, e := b.Get(id)
	if e != nil {
		return e
	}
	if d.State != "running" {
		return fmt.Errorf("not running")
	}
	d.State = "finished"
	b.rows[id] = d
	return nil
}
func (b *Board) Cancel(id string) error {
	d, e := b.Get(id)
	if e != nil {
		return e
	}
	if d.State == "finished" {
		return fmt.Errorf("already finished")
	}
	d.State = "cancelled"
	b.rows[id] = d
	return nil
}
func (b *Board) List(state string) []Dispatch {
	out := []Dispatch{}
	for _, d := range b.rows {
		if state != "" && d.State != state {
			continue
		}
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Priority > out[j].Priority })
	return out
}
func (b *Board) Count() int { return len(b.rows) }
