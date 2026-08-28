package api

import (
	"labops/internal/store"
	"net/http"
)

func HealthHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if e := s.Health(); e != nil {
			http.Error(w, "unhealthy", 503)
			return
		}
		w.WriteHeader(200)
	}
}
