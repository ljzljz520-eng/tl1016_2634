package model

import "strings"

func (r Record) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return err("id required")
	}
	if strings.TrimSpace(r.AssetTag) == "" {
		return err("asset tag required")
	}
	if r.Status == "" {
		return err("status required")
	}
	return nil
}
func (r Record) IsArchived() bool { return r.Status == "archived" }
func (r Record) CanEdit() bool {
	return r.Status == "draft" || r.Status == "active" || r.Status == "maintenance"
}
func NormalizeStatus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return "draft"
	}
	return s
}

type validationError string

func (e validationError) Error() string { return string(e) }
func err(s string) error                { return validationError(s) }
