package review

import (
	"bytes"
	"encoding/json"
	"labops/internal/model"
)

func BuildDecision(r model.Record, action string) model.Decision {
	if action == "approve" && r.Status == "pending" {
		return model.Decision{Allowed: true}
	}
	if action == "submit" && r.Status == "draft" {
		return model.Decision{Allowed: true}
	}
	return model.Decision{Reason: "state does not permit action"}
}
func EncodeDecision(d model.Decision) ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b).Encode(d)
	return b.Bytes(), e
}
