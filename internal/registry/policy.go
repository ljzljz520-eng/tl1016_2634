package registry

import (
	"fmt"
	"labops/internal/model"
)

func CheckTransition(from, to string) error {
	allowed := map[string][]string{"draft": {"pending", "archived"}, "pending": {"active", "draft"}, "active": {"maintenance", "archived"}, "maintenance": {"active", "archived"}}
	for _, x := range allowed[from] {
		if x == to {
			return nil
		}
	}
	return fmt.Errorf("transition %s to %s denied", from, to)
}
func CopyForUpdate(old model.Record, name, location string) model.Record {
	next := old
	if name != "" {
		next.Name = name
	}
	if location != "" {
		next.Location = location
	}
	return next
}
