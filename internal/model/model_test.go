package model

import "testing"

func TestRecordValidation(t *testing.T) {
	if (Record{ID: "x", AssetTag: "a", Status: "draft"}).Validate() != nil {
		t.Fatal("valid record rejected")
	}
	if (Record{}).Validate() == nil {
		t.Fatal("invalid record accepted")
	}
}
