package server

import (
	"fmt"
	"labops/internal/api"
	"labops/internal/flow075"
	"labops/internal/registry"
	"labops/internal/review"
	"labops/internal/schedule"
	"labops/internal/store"
	"net/http"
	"os"
)

func Build(path string) *api.Server {
	s, e := store.Open(path)
	if e != nil {
		panic(e)
	}
	r := registry.New(s)
	v := review.New(s)
	p := schedule.New(s)
	f := flow075.New(p)
	return api.New(r, v, f)
}
func Run() {
	path := os.Getenv("LABOPS_DB")
	if path == "" {
		path = "labops.db"
	}
	fmt.Println("labops listening on :8080")
	http.ListenAndServe(":8080", Build(path))
}
