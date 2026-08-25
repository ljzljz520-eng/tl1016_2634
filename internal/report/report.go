package report

import (
	"fmt"
	"labops/internal/model"
	"sort"
	"strings"
	"time"
)

type Summary struct {
	Total, Draft, Pending, Active, Maintenance, Archived int
	Locations                                            map[string]int
	Owners                                               map[string]int
}
type Finding struct{ Code, Severity, Message, RecordID string }
type Report struct {
	Summary     Summary
	Findings    []Finding
	GeneratedAt string
}

func NewSummary(rows []model.Record) Summary {
	s := Summary{Locations: map[string]int{}, Owners: map[string]int{}}
	for _, r := range rows {
		s.Total++
		s.Locations[r.Location]++
		s.Owners[r.Owner]++
		switch r.Status {
		case "draft":
			s.Draft++
		case "pending":
			s.Pending++
		case "active":
			s.Active++
		case "maintenance":
			s.Maintenance++
		case "archived":
			s.Archived++
		}
	}
	return s
}
func Build(rows []model.Record, stamp time.Time) Report {
	r := Report{Summary: NewSummary(rows), GeneratedAt: stamp.UTC().Format(time.RFC3339)}
	for _, x := range rows {
		r.Findings = append(r.Findings, FindingsFor(x)...)
	}
	sort.Slice(r.Findings, func(i, j int) bool { return r.Findings[i].RecordID < r.Findings[j].RecordID })
	return r
}
func FindingsFor(r model.Record) []Finding {
	out := []Finding{}
	if r.AssetTag == "" {
		out = append(out, Finding{"MISSING_TAG", "high", "asset tag is missing", r.ID})
	}
	if r.Location == "" {
		out = append(out, Finding{"MISSING_LOCATION", "medium", "location is missing", r.ID})
	}
	if r.Status == "pending" {
		out = append(out, Finding{"WAITING_REVIEW", "low", "record awaits review", r.ID})
	}
	if r.IsArchived() && len(r.Slots) > 0 {
		out = append(out, Finding{"STALE_SLOTS", "medium", "archived record has slots", r.ID})
	}
	return out
}
func FormatSummary(s Summary) string {
	return fmt.Sprintf("total=%d draft=%d pending=%d active=%d maintenance=%d archived=%d", s.Total, s.Draft, s.Pending, s.Active, s.Maintenance, s.Archived)
}
func CSV(rows []model.Record) string {
	lines := []string{"id,asset_tag,name,location,status,owner,version"}
	for _, r := range rows {
		lines = append(lines, strings.Join([]string{r.ID, r.AssetTag, r.Name, r.Location, r.Status, r.Owner, fmt.Sprint(r.Version)}, ","))
	}
	return strings.Join(lines, "\n")
}
func FilterFindings(in []Finding, severity string) []Finding {
	out := []Finding{}
	for _, f := range in {
		if severity == "" || f.Severity == severity {
			out = append(out, f)
		}
	}
	return out
}
func CountSeverity(in []Finding) map[string]int {
	out := map[string]int{}
	for _, f := range in {
		out[f.Severity]++
	}
	return out
}
func Merge(a, b Report) Report {
	r := a
	r.Summary.Total += b.Summary.Total
	r.Summary.Draft += b.Summary.Draft
	r.Summary.Pending += b.Summary.Pending
	r.Summary.Active += b.Summary.Active
	r.Summary.Maintenance += b.Summary.Maintenance
	r.Summary.Archived += b.Summary.Archived
	r.Findings = append(r.Findings, b.Findings...)
	for k, v := range b.Summary.Locations {
		r.Summary.Locations[k] += v
	}
	for k, v := range b.Summary.Owners {
		r.Summary.Owners[k] += v
	}
	return r
}
func LocationRows(s Summary) []string {
	keys := make([]string, 0, len(s.Locations))
	for k := range s.Locations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []string{}
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s:%d", k, s.Locations[k]))
	}
	return out
}
func OwnerRows(s Summary) []string {
	keys := make([]string, 0, len(s.Owners))
	for k := range s.Owners {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []string{}
	for _, k := range keys {
		out = append(out, fmt.Sprintf("%s:%d", k, s.Owners[k]))
	}
	return out
}
