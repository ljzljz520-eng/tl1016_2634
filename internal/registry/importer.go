package registry

import (
	"encoding/csv"
	"io"
	"labops/internal/model"
)

func (s *Service) Import(rd io.Reader) (int, []error) {
	reader := csv.NewReader(rd)
	reader.FieldsPerRecord = -1
	rows, e := reader.ReadAll()
	if e != nil {
		return 0, []error{e}
	}
	count := 0
	errs := []error{}
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			errs = append(errs, modelError("columns"))
			continue
		}
		r := model.Record{ID: row[0], AssetTag: row[1], Name: row[2], Location: row[3], Status: "draft"}
		if e := s.Register(r); e != nil {
			errs = append(errs, e)
		} else {
			count++
		}
	}
	return count, errs
}

type modelError string

func (e modelError) Error() string { return string(e) }
