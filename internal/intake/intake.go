package intake

import (
	"fmt"
	"labops/internal/model"
	"strings"
)

type Form struct {
	AssetTag, Name, Location, Owner, Notes string
	Slots                                  []string
}
type Result struct {
	Record   model.Record
	Warnings []string
}

func Parse(f Form) (Result, error) {
	warnings := []string{}
	if strings.TrimSpace(f.AssetTag) == "" {
		return Result{}, fmt.Errorf("asset tag required")
	}
	if strings.TrimSpace(f.Name) == "" {
		return Result{}, fmt.Errorf("name required")
	}
	if strings.TrimSpace(f.Location) == "" {
		warnings = append(warnings, "location not set")
	}
	if strings.TrimSpace(f.Owner) == "" {
		warnings = append(warnings, "owner not set")
	}
	return Result{Record: model.Record{AssetTag: strings.TrimSpace(f.AssetTag), Name: strings.TrimSpace(f.Name), Location: strings.TrimSpace(f.Location), Owner: strings.TrimSpace(f.Owner), Status: "draft", Slots: append([]string(nil), f.Slots...)}, Warnings: warnings}, nil
}
func Normalize(f Form) Form {
	f.AssetTag = strings.ToUpper(strings.TrimSpace(f.AssetTag))
	f.Name = strings.TrimSpace(f.Name)
	f.Location = strings.TrimSpace(f.Location)
	f.Owner = strings.TrimSpace(f.Owner)
	f.Notes = strings.TrimSpace(f.Notes)
	return f
}
func Validate(f Form) []string {
	out := []string{}
	if f.AssetTag == "" {
		out = append(out, "asset tag required")
	}
	if f.Name == "" {
		out = append(out, "name required")
	}
	if len(f.Slots) > 8 {
		out = append(out, "too many slots")
	}
	seen := map[string]bool{}
	for _, s := range f.Slots {
		if seen[s] {
			out = append(out, "duplicate slot")
		}
		seen[s] = true
	}
	return out
}
func MakeID(tag string, sequence int) string {
	return fmt.Sprintf("%s-%04d", strings.ToLower(strings.ReplaceAll(tag, " ", "-")), sequence)
}
func Merge(base, patch Form) Form {
	if patch.AssetTag != "" {
		base.AssetTag = patch.AssetTag
	}
	if patch.Name != "" {
		base.Name = patch.Name
	}
	if patch.Location != "" {
		base.Location = patch.Location
	}
	if patch.Owner != "" {
		base.Owner = patch.Owner
	}
	if patch.Notes != "" {
		base.Notes = patch.Notes
	}
	if patch.Slots != nil {
		base.Slots = append([]string(nil), patch.Slots...)
	}
	return base
}
func IsComplete(f Form) bool        { return len(Validate(f)) == 0 && f.Location != "" && f.Owner != "" }
func WarningText(w []string) string { return strings.Join(w, "; ") }
func Clone(f Form) Form             { f.Slots = append([]string(nil), f.Slots...); return f }
