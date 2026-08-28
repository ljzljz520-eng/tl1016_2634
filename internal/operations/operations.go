package operations

import (
	"fmt"
	"labops/internal/model"
	"strings"
)

type Checklist struct {
	Items []string
	Done  map[string]bool
}

func NewChecklist(items []string) *Checklist {
	c := &Checklist{Items: append([]string(nil), items...), Done: map[string]bool{}}
	for _, x := range items {
		c.Done[x] = false
	}
	return c
}
func (c *Checklist) Mark(item string) error {
	if _, ok := c.Done[item]; !ok {
		return fmt.Errorf("unknown item")
	}
	c.Done[item] = true
	return nil
}
func (c *Checklist) Unmark(item string) error {
	if _, ok := c.Done[item]; !ok {
		return fmt.Errorf("unknown item")
	}
	c.Done[item] = false
	return nil
}
func (c *Checklist) Complete() bool {
	for _, x := range c.Items {
		if !c.Done[x] {
			return false
		}
	}
	return true
}
func (c *Checklist) Remaining() []string {
	out := []string{}
	for _, x := range c.Items {
		if !c.Done[x] {
			out = append(out, x)
		}
	}
	return out
}
func (c *Checklist) Progress() float64 {
	if len(c.Items) == 0 {
		return 1
	}
	n := 0
	for _, x := range c.Items {
		if c.Done[x] {
			n++
		}
	}
	return float64(n) / float64(len(c.Items))
}
func (c *Checklist) Reset() {
	for x := range c.Done {
		c.Done[x] = false
	}
}
func (c *Checklist) Add(item string) error {
	if strings.TrimSpace(item) == "" {
		return fmt.Errorf("item required")
	}
	if _, ok := c.Done[item]; ok {
		return fmt.Errorf("item exists")
	}
	c.Items = append(c.Items, item)
	c.Done[item] = false
	return nil
}
func (c *Checklist) Remove(item string) error {
	if _, ok := c.Done[item]; !ok {
		return fmt.Errorf("unknown item")
	}
	delete(c.Done, item)
	for i, x := range c.Items {
		if x == item {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			break
		}
	}
	return nil
}
func StatusLabel(r model.Record) string {
	switch r.Status {
	case "draft":
		return "待登记"
	case "pending":
		return "待审核"
	case "active":
		return "运行中"
	case "maintenance":
		return "维护中"
	case "archived":
		return "已归档"
	}
	return "未知"
}
func CanSchedule(r model.Record) bool { return r.Status == "active" || r.Status == "maintenance" }
