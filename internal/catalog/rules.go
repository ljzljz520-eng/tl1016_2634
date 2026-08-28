package catalog

import (
	"fmt"
	"strings"
)

type FieldRule struct {
	Name     string
	Required bool
	Max      int
}
type Validator struct{ Rules []FieldRule }

func NewValidator() *Validator { return &Validator{Rules: []FieldRule{}} }
func (v *Validator) Add(r FieldRule) error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name required")
	}
	if r.Max < 0 {
		return fmt.Errorf("max invalid")
	}
	v.Rules = append(v.Rules, r)
	return nil
}
func (v *Validator) Check(values map[string]string) []string {
	out := []string{}
	for _, r := range v.Rules {
		x := strings.TrimSpace(values[r.Name])
		if r.Required && x == "" {
			out = append(out, r.Name+" required")
		}
		if r.Max > 0 && len(x) > r.Max {
			out = append(out, r.Name+" too long")
		}
	}
	return out
}
func (v *Validator) Valid(values map[string]string) bool { return len(v.Check(values)) == 0 }
func (v *Validator) Names() []string {
	out := []string{}
	for _, r := range v.Rules {
		out = append(out, r.Name)
	}
	return out
}
func (v *Validator) Remove(name string) bool {
	for i, r := range v.Rules {
		if r.Name == name {
			v.Rules = append(v.Rules[:i], v.Rules[i+1:]...)
			return true
		}
	}
	return false
}
func (v *Validator) Required() []string {
	out := []string{}
	for _, r := range v.Rules {
		if r.Required {
			out = append(out, r.Name)
		}
	}
	return out
}
func NormalizeTags(tags []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, t := range tags {
		x := strings.ToLower(strings.TrimSpace(t))
		if x != "" && !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
func MatchTags(want, have []string) bool {
	hs := map[string]bool{}
	for _, x := range NormalizeTags(have) {
		hs[x] = true
	}
	for _, x := range NormalizeTags(want) {
		if !hs[x] {
			return false
		}
	}
	return true
}
func CategoryCode(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "" {
		return "GEN"
	}
	return s[:min(3, len(s))]
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
