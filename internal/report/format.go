package report

import (
	"fmt"
	"strings"
)

func Markdown(r Report) string {
	lines := []string{"# Equipment report", "", "Generated: " + r.GeneratedAt, "", FormatSummary(r.Summary), "", "## Findings"}
	for _, f := range r.Findings {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %s", f.Severity, f.RecordID, f.Message))
	}
	return strings.Join(lines, "\n")
}
func Plain(r Report) string {
	lines := []string{r.GeneratedAt, FormatSummary(r.Summary)}
	for _, f := range r.Findings {
		lines = append(lines, f.RecordID+" "+f.Code)
	}
	return strings.Join(lines, "\n")
}
func FindingCodes(r Report) []string {
	out := []string{}
	for _, f := range r.Findings {
		out = append(out, f.Code)
	}
	return out
}
func HasCode(r Report, code string) bool {
	for _, f := range r.Findings {
		if f.Code == code {
			return true
		}
	}
	return false
}
func ByRecord(r Report, id string) []Finding {
	out := []Finding{}
	for _, f := range r.Findings {
		if f.RecordID == id {
			out = append(out, f)
		}
	}
	return out
}
func Highest(r Report) string {
	rank := map[string]int{"low": 1, "medium": 2, "high": 3}
	best := 0
	out := ""
	for _, f := range r.Findings {
		if rank[f.Severity] > best {
			best = rank[f.Severity]
			out = f.Severity
		}
	}
	return out
}
func IsEmpty(r Report) bool            { return r.Summary.Total == 0 && len(r.Findings) == 0 }
func RecordCount(r Report) int         { return r.Summary.Total }
func FindingCount(r Report) int        { return len(r.Findings) }
func LocationCountReport(r Report) int { return len(r.Summary.Locations) }
func OwnerCountReport(r Report) int    { return len(r.Summary.Owners) }
func StatusRate(r Report, status string) float64 {
	if r.Summary.Total == 0 {
		return 0
	}
	switch status {
	case "draft":
		return float64(r.Summary.Draft) / float64(r.Summary.Total)
	case "pending":
		return float64(r.Summary.Pending) / float64(r.Summary.Total)
	case "active":
		return float64(r.Summary.Active) / float64(r.Summary.Total)
	case "maintenance":
		return float64(r.Summary.Maintenance) / float64(r.Summary.Total)
	case "archived":
		return float64(r.Summary.Archived) / float64(r.Summary.Total)
	}
	return 0
}
func SeverityRate(r Report, sev string) float64 {
	if r.Summary.Total == 0 {
		return 0
	}
	return float64(len(FilterFindings(r.Findings, sev))) / float64(r.Summary.Total)
}
