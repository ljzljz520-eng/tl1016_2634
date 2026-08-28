package archive

import (
	"fmt"
	"labops/internal/model"
	"sort"
	"strings"
)

type Package struct {
	ID, RecordID, Reason, Operator string
	State                          string
	Items                          []string
}
type Registry struct{ packages map[string]Package }

func New() *Registry { return &Registry{packages: map[string]Package{}} }
func (r *Registry) Create(p Package) error {
	if p.ID == "" || p.RecordID == "" {
		return fmt.Errorf("identity required")
	}
	if strings.TrimSpace(p.Reason) == "" {
		return fmt.Errorf("reason required")
	}
	if _, ok := r.packages[p.ID]; ok {
		return fmt.Errorf("package exists")
	}
	p.State = "open"
	r.packages[p.ID] = p
	return nil
}
func (r *Registry) Get(id string) (Package, error) {
	p, ok := r.packages[id]
	if !ok {
		return Package{}, fmt.Errorf("package missing")
	}
	return p, nil
}
func (r *Registry) AddItem(id, item string) error {
	p, e := r.Get(id)
	if e != nil {
		return e
	}
	if p.State != "open" {
		return fmt.Errorf("package closed")
	}
	if strings.TrimSpace(item) == "" {
		return fmt.Errorf("item required")
	}
	p.Items = append(p.Items, item)
	r.packages[id] = p
	return nil
}
func (r *Registry) Seal(id string) error {
	p, e := r.Get(id)
	if e != nil {
		return e
	}
	if p.State != "open" {
		return fmt.Errorf("package not open")
	}
	if len(p.Items) == 0 {
		return fmt.Errorf("empty package")
	}
	p.State = "sealed"
	r.packages[id] = p
	return nil
}
func (r *Registry) Reopen(id string) error {
	p, e := r.Get(id)
	if e != nil {
		return e
	}
	if p.State != "sealed" {
		return fmt.Errorf("package not sealed")
	}
	p.State = "open"
	r.packages[id] = p
	return nil
}
func (r *Registry) List(state string) []Package {
	out := []Package{}
	for _, p := range r.packages {
		if state != "" && p.State != state {
			continue
		}
		p.Items = append([]string(nil), p.Items...)
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (r *Registry) Count() int       { return len(r.packages) }
func CanArchive(x model.Record) bool { return x.Status == "archived" }
func Reason(code string) string {
	switch code {
	case "retired":
		return "equipment retired"
	case "replaced":
		return "equipment replaced"
	case "duplicate":
		return "duplicate registration"
	}
	return "operational archive"
}
