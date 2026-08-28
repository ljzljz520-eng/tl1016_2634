package analytics

import (
	"labops/internal/model"
	"math"
	"sort"
)

type Metrics struct {
	Total       float64
	Active      float64
	Completion  float64
	Utilization float64
	ByStatus    map[string]int
	ByLocation  map[string]int
}

func Compute(rows []model.Record) Metrics {
	m := Metrics{ByStatus: map[string]int{}, ByLocation: map[string]int{}}
	for _, r := range rows {
		m.Total++
		m.ByStatus[r.Status]++
		m.ByLocation[r.Location]++
		if r.Status == "active" || r.Status == "maintenance" {
			m.Active++
		}
		if r.IsArchived() {
			m.Completion++
		}
		if len(r.Slots) > 0 {
			m.Utilization++
		}
	}
	if m.Total > 0 {
		m.Active /= m.Total
		m.Completion /= m.Total
		m.Utilization /= m.Total
	}
	return m
}
func Round(v float64) float64 { return math.Round(v*100) / 100 }
func StatusNames(m Metrics) []string {
	k := make([]string, 0, len(m.ByStatus))
	for x := range m.ByStatus {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}
func LocationNames(m Metrics) []string {
	k := make([]string, 0, len(m.ByLocation))
	for x := range m.ByLocation {
		k = append(k, x)
	}
	sort.Strings(k)
	return k
}
func StatusCount(m Metrics, s string) int   { return m.ByStatus[s] }
func LocationCount(m Metrics, s string) int { return m.ByLocation[s] }
func IsHealthy(m Metrics) bool              { return m.Total > 0 && m.Completion >= 0 && m.Completion <= 1 }
func Risk(m Metrics) string {
	if m.Total == 0 {
		return "unknown"
	}
	if m.Completion < 0.2 {
		return "high"
	}
	if m.Completion < 0.6 {
		return "medium"
	}
	return "low"
}
func Compare(a, b Metrics) Metrics {
	return Metrics{Total: a.Total - b.Total, Active: a.Active - b.Active, Completion: a.Completion - b.Completion, Utilization: a.Utilization - b.Utilization}
}
func Merge(a, b Metrics) Metrics {
	m := Metrics{Total: a.Total + b.Total, ByStatus: map[string]int{}, ByLocation: map[string]int{}}
	for k, v := range a.ByStatus {
		m.ByStatus[k] += v
	}
	for k, v := range b.ByStatus {
		m.ByStatus[k] += v
	}
	for k, v := range a.ByLocation {
		m.ByLocation[k] += v
	}
	for k, v := range b.ByLocation {
		m.ByLocation[k] += v
	}
	if m.Total > 0 {
		m.Active = (a.Active*a.Total + b.Active*b.Total) / m.Total
		m.Completion = (a.Completion*a.Total + b.Completion*b.Total) / m.Total
		m.Utilization = (a.Utilization*a.Total + b.Utilization*b.Total) / m.Total
	}
	return m
}
