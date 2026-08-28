package policy

import (
	"fmt"
	"labops/internal/model"
	"strings"
)

type Rule struct {
	Code, Description string
	Check             func(model.Record) bool
}
type Engine struct{ rules []Rule }

func New() *Engine { return &Engine{rules: []Rule{}} }
func (e *Engine) Add(r Rule) error {
	if r.Code == "" || r.Check == nil {
		return fmt.Errorf("invalid rule")
	}
	for _, x := range e.rules {
		if x.Code == r.Code {
			return fmt.Errorf("rule exists")
		}
	}
	e.rules = append(e.rules, r)
	return nil
}
func (e *Engine) Remove(code string) bool {
	for i, r := range e.rules {
		if r.Code == code {
			e.rules = append(e.rules[:i], e.rules[i+1:]...)
			return true
		}
	}
	return false
}
func (e *Engine) Evaluate(r model.Record) []string {
	out := []string{}
	for _, x := range e.rules {
		if !x.Check(r) {
			out = append(out, x.Code)
		}
	}
	return out
}
func (e *Engine) Allowed(r model.Record) bool { return len(e.Evaluate(r)) == 0 }
func (e *Engine) Explain(r model.Record) string {
	failed := e.Evaluate(r)
	if len(failed) == 0 {
		return "allowed"
	}
	return "blocked:" + strings.Join(failed, ",")
}
func Default() *Engine {
	e := New()
	e.Add(Rule{"owner-required", "owner required", func(r model.Record) bool { return strings.TrimSpace(r.Owner) != "" }})
	e.Add(Rule{"known-status", "known status", func(r model.Record) bool {
		switch r.Status {
		case "draft", "pending", "active", "maintenance", "archived":
			return true
		}
		return false
	}})
	e.Add(Rule{"version-positive", "version nonnegative", func(r model.Record) bool { return r.Version >= 0 }})
	return e
}
func (e *Engine) Codes() []string {
	out := []string{}
	for _, r := range e.rules {
		out = append(out, r.Code)
	}
	return out
}
