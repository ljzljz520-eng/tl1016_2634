package analytics

import (
	"sort"
	"time"
)

type Event struct {
	At             time.Time
	Kind, RecordID string
	Value          int
}
type Timeline struct{ Events []Event }

func NewTimeline() *Timeline { return &Timeline{Events: []Event{}} }
func (t *Timeline) Add(e Event) {
	t.Events = append(t.Events, e)
	sort.SliceStable(t.Events, func(i, j int) bool { return t.Events[i].At.Before(t.Events[j].At) })
}
func (t *Timeline) Between(from, to time.Time) []Event {
	out := []Event{}
	for _, e := range t.Events {
		if e.At.Before(from) || e.At.After(to) {
			continue
		}
		out = append(out, e)
	}
	return out
}
func (t *Timeline) ForRecord(id string) []Event {
	out := []Event{}
	for _, e := range t.Events {
		if e.RecordID == id {
			out = append(out, e)
		}
	}
	return out
}
func (t *Timeline) Kinds() map[string]int {
	out := map[string]int{}
	for _, e := range t.Events {
		out[e.Kind]++
	}
	return out
}
func (t *Timeline) Latest(id string) (Event, bool) {
	for i := len(t.Events) - 1; i >= 0; i-- {
		if t.Events[i].RecordID == id {
			return t.Events[i], true
		}
	}
	return Event{}, false
}
func (t *Timeline) First(id string) (Event, bool) {
	for _, e := range t.Events {
		if e.RecordID == id {
			return e, true
		}
	}
	return Event{}, false
}
func (t *Timeline) Sum(kind string) int {
	n := 0
	for _, e := range t.Events {
		if e.Kind == kind {
			n += e.Value
		}
	}
	return n
}
func (t *Timeline) Count(id string) int { return len(t.ForRecord(id)) }
func (t *Timeline) Empty() bool         { return len(t.Events) == 0 }
