package store

import (
	"encoding/json"
	"io"
	"labops/internal/model"
)

func (s *Store) Export(w io.Writer) error {
	rows, e := s.ListRecords()
	if e != nil {
		return e
	}
	return json.NewEncoder(w).Encode(rows)
}
func EncodeRecord(r model.Record) ([]byte, error) { return json.Marshal(r) }
