package api

import (
	"encoding/json"
	"labops/internal/flow075"
	"labops/internal/model"
	"labops/internal/registry"
	"labops/internal/review"
	"net/http"
)

type Server struct {
	Registry *registry.Service
	Review   *review.Service
	Flow     *flow075.Processor
}

func New(r *registry.Service, v *review.Service, f *flow075.Processor) *Server {
	return &Server{Registry: r, Review: v, Flow: f}
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/records" {
		var x model.Record
		if json.NewDecoder(r.Body).Decode(&x) != nil || s.Registry.Register(x) != nil {
			http.Error(w, "invalid", 400)
			return
		}
		w.WriteHeader(201)
		return
	}
	if r.Method == "POST" && r.URL.Path == "/review" {
		var x struct{ ID, Actor string }
		json.NewDecoder(r.Body).Decode(&x)
		if e := s.Review.Submit(x.ID, x.Actor); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		w.WriteHeader(204)
		return
	}
	http.NotFound(w, r)
}
