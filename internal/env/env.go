package env

import (
	"os"
	"strconv"
)

func StringEnv(name string, def string) string {
	v := os.Getenv(name)
	if len(v) == 0 {
		return def
	}
	return v
}

func UintEnv(name string, def uint) uint {
	v := os.Getenv(name)
	if len(v) == 0 {
		return def
	}
	i, err := strconv.ParseUint(v, 10, 32)
	if err != nil {
		return def
	}
	return uint(i)
}

func BoolEnv(name string, def bool) bool {
	v := os.Getenv(name)
	if v == "true" {
		return true
	}
	if v == "false" {
		return false
	}
	return def
}
