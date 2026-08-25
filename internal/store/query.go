package store

import (
	"go.etcd.io/bbolt"
	"labops/internal/model"
	"sort"
	"strings"
)

func (s *Store) Search(q model.Query) (model.Result, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return model.Result{}, e
	}
	filtered := make([]model.Record, 0)
	for _, r := range rows {
		if q.Status != "" && r.Status != q.Status {
			continue
		}
		if q.Location != "" && r.Location != q.Location {
			continue
		}
		needle := strings.ToLower(q.Text)
		if needle != "" && !strings.Contains(strings.ToLower(r.Name+" "+r.AssetTag), needle) {
			continue
		}
		filtered = append(filtered, r)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].ID < filtered[j].ID })
	limit := q.Limit
	if limit <= 0 || limit > len(filtered) {
		limit = len(filtered)
	}
	return model.Result{Records: filtered[:limit], Total: len(filtered)}, nil
}
func (s *Store) Count(bucket string) (int, error) {
	n := 0
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		n = b.Stats().KeyN
		return nil
	})
	return n, e
}
