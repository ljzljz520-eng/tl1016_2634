package security

import (
	"fmt"
	"sort"
	"strings"
)

type User struct {
	ID, Name string
	Roles    []string
	Active   bool
}
type Directory struct{ users map[string]User }

func New() *Directory { return &Directory{users: map[string]User{}} }
func (d *Directory) Add(u User) error {
	if u.ID == "" || u.Name == "" {
		return fmt.Errorf("identity required")
	}
	if _, ok := d.users[u.ID]; ok {
		return fmt.Errorf("user exists")
	}
	u.Active = true
	d.users[u.ID] = u
	return nil
}
func (d *Directory) Get(id string) (User, error) {
	u, ok := d.users[id]
	if !ok {
		return User{}, fmt.Errorf("user missing")
	}
	return u, nil
}
func (d *Directory) Remove(id string) error {
	if _, e := d.Get(id); e != nil {
		return e
	}
	delete(d.users, id)
	return nil
}
func (d *Directory) Activate(id string) error {
	u, e := d.Get(id)
	if e != nil {
		return e
	}
	u.Active = true
	d.users[id] = u
	return nil
}
func (d *Directory) Deactivate(id string) error {
	u, e := d.Get(id)
	if e != nil {
		return e
	}
	u.Active = false
	d.users[id] = u
	return nil
}
func (d *Directory) Grant(id, role string) error {
	u, e := d.Get(id)
	if e != nil {
		return e
	}
	if strings.TrimSpace(role) == "" {
		return fmt.Errorf("role required")
	}
	for _, r := range u.Roles {
		if r == role {
			return nil
		}
	}
	u.Roles = append(u.Roles, role)
	sort.Strings(u.Roles)
	d.users[id] = u
	return nil
}
func (d *Directory) Revoke(id, role string) error {
	u, e := d.Get(id)
	if e != nil {
		return e
	}
	for i, r := range u.Roles {
		if r == role {
			u.Roles = append(u.Roles[:i], u.Roles[i+1:]...)
			d.users[id] = u
			return nil
		}
	}
	return fmt.Errorf("role missing")
}
func (d *Directory) Has(id, role string) bool {
	u, e := d.Get(id)
	if e != nil || !u.Active {
		return false
	}
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}
func (d *Directory) List() []User {
	out := []User{}
	for _, u := range d.users {
		u.Roles = append([]string(nil), u.Roles...)
		out = append(out, u)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (d *Directory) Count() int { return len(d.users) }
