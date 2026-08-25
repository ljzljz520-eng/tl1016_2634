package server

import (
	"os"
	"strconv"
)

type Config struct {
	Addr, DBPath string
	Debug        bool
}

func LoadConfig() Config {
	c := Config{Addr: ":8080", DBPath: "labops.db"}
	if v := os.Getenv("LABOPS_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("LABOPS_DB"); v != "" {
		c.DBPath = v
	}
	c.Debug, _ = strconv.ParseBool(os.Getenv("LABOPS_DEBUG"))
	return c
}
