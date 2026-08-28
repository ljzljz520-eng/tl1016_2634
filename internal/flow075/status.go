package flow075

import (
	"fmt"
	"labops/internal/model"
)

func StatusFor(r model.Record) string {
	if r.IsArchived() {
		return "closed"
	}
	if len(r.Slots) > 0 {
		return "scheduled"
	}
	return fmt.Sprintf("%s-ready", r.Status)
}
func ValidateOperation(r model.Record) error {
	if !r.CanEdit() {
		return fmt.Errorf("record cannot be operated")
	}
	return nil
}
