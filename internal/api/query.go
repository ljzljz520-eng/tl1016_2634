package api

import (
	"encoding/json"
	"labops/internal/model"
	"labops/internal/store"
	"net/http"
)

func QueryHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := model.Query{Text: r.URL.Query().Get("q"), Status: r.URL.Query().Get("status"), Location: r.URL.Query().Get("location")}
		out, e := s.Search(q)
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(out)
	}
}
