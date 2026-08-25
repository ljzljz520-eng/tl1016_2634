package catalog

import (
	"fmt"
	"sort"
	"strings"
)

type Entry struct {
	Code, Category, Manufacturer, Model string
	Active                              bool
	Tags                                []string
}
type Catalog struct {
	entries map[string]Entry
	aliases map[string]string
}

func New() *Catalog { return &Catalog{entries: map[string]Entry{}, aliases: map[string]string{}} }
func (c *Catalog) Add(e Entry) error {
	if strings.TrimSpace(e.Code) == "" {
		return fmt.Errorf("code required")
	}
	if _, ok := c.entries[e.Code]; ok {
		return fmt.Errorf("code exists")
	}
	if e.Category == "" {
		e.Category = "general"
	}
	c.entries[e.Code] = e
	return nil
}
func (c *Catalog) Update(e Entry) error {
	if _, ok := c.entries[e.Code]; !ok {
		return fmt.Errorf("code missing")
	}
	c.entries[e.Code] = e
	return nil
}
func (c *Catalog) Remove(code string) error {
	if _, ok := c.entries[code]; !ok {
		return fmt.Errorf("code missing")
	}
	delete(c.entries, code)
	for a, v := range c.aliases {
		if v == code {
			delete(c.aliases, a)
		}
	}
	return nil
}
func (c *Catalog) Get(code string) (Entry, error) {
	if x, ok := c.entries[code]; ok {
		return x, nil
	}
	if real, ok := c.aliases[code]; ok {
		return c.entries[real], nil
	}
	return Entry{}, fmt.Errorf("not found")
}
func (c *Catalog) Alias(alias, code string) error {
	if _, e := c.Get(code); e != nil {
		return e
	}
	if alias == "" {
		return fmt.Errorf("alias required")
	}
	c.aliases[alias] = code
	return nil
}
func (c *Catalog) Search(text, category string) []Entry {
	out := []Entry{}
	needle := strings.ToLower(text)
	for _, e := range c.entries {
		if !e.Active {
			continue
		}
		if category != "" && e.Category != category {
			continue
		}
		hay := strings.ToLower(e.Code + " " + e.Manufacturer + " " + e.Model + " " + strings.Join(e.Tags, " "))
		if needle != "" && !strings.Contains(hay, needle) {
			continue
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}
func (c *Catalog) Categories() []string {
	seen := map[string]bool{}
	for _, e := range c.entries {
		seen[e.Category] = true
	}
	out := make([]string, 0, len(seen))
	for x := range seen {
		out = append(out, x)
	}
	sort.Strings(out)
	return out
}
func (c *Catalog) Activate(code string) error {
	e, x := c.Get(code)
	if x != nil {
		return x
	}
	e.Active = true
	return c.Update(e)
}
func (c *Catalog) Deactivate(code string) error {
	e, x := c.Get(code)
	if x != nil {
		return x
	}
	e.Active = false
	return c.Update(e)
}
func (c *Catalog) Count() int { return len(c.entries) }
func (c *Catalog) Snapshot() []Entry {
	out := make([]Entry, 0, len(c.entries))
	for _, e := range c.entries {
		e.Tags = append([]string(nil), e.Tags...)
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}
