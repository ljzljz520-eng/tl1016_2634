package maintenance

import (
	"fmt"
	"sort"
	"strings"
)

type Task struct {
	ID, RecordID, Kind, Technician, State string
	Priority                              int
	Notes                                 []string
}
type Queue struct{ tasks map[string]Task }

func New() *Queue { return &Queue{tasks: map[string]Task{}} }
func (q *Queue) Add(t Task) error {
	if t.ID == "" || t.RecordID == "" {
		return fmt.Errorf("task identity required")
	}
	if _, ok := q.tasks[t.ID]; ok {
		return fmt.Errorf("task exists")
	}
	if t.State == "" {
		t.State = "open"
	}
	q.tasks[t.ID] = t
	return nil
}
func (q *Queue) Get(id string) (Task, error) {
	t, ok := q.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("task missing")
	}
	return t, nil
}
func (q *Queue) Assign(id, tech string) error {
	t, e := q.Get(id)
	if e != nil {
		return e
	}
	if t.State != "open" {
		return fmt.Errorf("task not open")
	}
	if strings.TrimSpace(tech) == "" {
		return fmt.Errorf("technician required")
	}
	t.Technician = tech
	t.State = "assigned"
	q.tasks[id] = t
	return nil
}
func (q *Queue) Start(id string) error {
	t, e := q.Get(id)
	if e != nil {
		return e
	}
	if t.State != "assigned" {
		return fmt.Errorf("task not assigned")
	}
	t.State = "running"
	q.tasks[id] = t
	return nil
}
func (q *Queue) Complete(id string) error {
	t, e := q.Get(id)
	if e != nil {
		return e
	}
	if t.State != "running" {
		return fmt.Errorf("task not running")
	}
	t.State = "done"
	q.tasks[id] = t
	return nil
}
func (q *Queue) Cancel(id string) error {
	t, e := q.Get(id)
	if e != nil {
		return e
	}
	if t.State == "done" {
		return fmt.Errorf("task completed")
	}
	t.State = "cancelled"
	q.tasks[id] = t
	return nil
}
func (q *Queue) List(state string) []Task {
	out := []Task{}
	for _, t := range q.tasks {
		if state != "" && t.State != state {
			continue
		}
		t.Notes = append([]string(nil), t.Notes...)
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Priority == out[j].Priority {
			return out[i].ID < out[j].ID
		}
		return out[i].Priority > out[j].Priority
	})
	return out
}
func (q *Queue) AddNote(id, note string) error {
	t, e := q.Get(id)
	if e != nil {
		return e
	}
	if strings.TrimSpace(note) == "" {
		return fmt.Errorf("note required")
	}
	t.Notes = append(t.Notes, note)
	q.tasks[id] = t
	return nil
}
func ValidKind(k string) bool {
	switch k {
	case "inspection", "repair", "calibration", "cleaning":
		return true
	}
	return false
}
func (q *Queue) Count() int { return len(q.tasks) }
