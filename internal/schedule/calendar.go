package schedule

import (
	"labops/internal/model"
	"sort"
)

type Calendar struct{ Days map[string][]string }

func NewCalendar() *Calendar { return &Calendar{Days: map[string][]string{}} }
func (c *Calendar) Add(day, record string) {
	c.Days[day] = append(c.Days[day], record)
	sort.Strings(c.Days[day])
}
func (c *Calendar) Records(day string) []string { return append([]string(nil), c.Days[day]...) }
func Summarize(r model.Record) string           { return r.AssetTag + " @ " + r.Location + " [" + r.Status + "]" }
