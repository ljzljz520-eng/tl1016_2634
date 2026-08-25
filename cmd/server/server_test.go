package server

import (
	"path/filepath"
	"testing"
)

func TestBuildServer(t *testing.T) {
	s := Build(filepath.Join(t.TempDir(), "db"))
	if s == nil {
		t.Fatal("nil server")
	}
}
