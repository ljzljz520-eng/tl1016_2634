package flow075

import (
	"labops/internal/schedule"
)

type Processor struct{ Planner *schedule.Planner }

func New(p *schedule.Planner) *Processor { return &Processor{Planner: p} }
func (p *Processor) Operate(id string, slots []string) error {
	if e := schedule.ValidateSlots(slots); e != nil {
		return e
	}
	return p.Planner.Assign(id, slots)
}
func (p *Processor) Close(id string) error               { return p.Planner.Release(id) }
func (p *Processor) Ordered(id string) ([]string, error) { return p.Planner.Slots(id) }
