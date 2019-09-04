package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/kou64yama/takanawa/internal/env"
)

func TestStringEnv(t *testing.T) {
	tests := []struct {
		val string
		def string
		got string
	}{
		{val: "foo", def: "bar", got: "foo"},
		{val: "", def: "bar", got: "bar"},
		{val: "", def: "", got: ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("val=%q, def=%q, got=%q", tt.val, tt.def, tt.got), func(t *testing.T) {
			os.Setenv("TAKANAWA_TEST", tt.val)

			v := env.StringEnv("TAKANAWA_TEST", tt.def)
			if v != tt.got {
				t.Errorf("got %q, want %q", v, tt.got)
			}
		})
	}
}

func TestUintEnv(t *testing.T) {
	tests := []struct {
		val string
		def uint
		got uint
	}{
		{val: "1", def: 2, got: 1},
		{val: "0", def: 2, got: 0},
		{val: "-1", def: 2, got: 2},
		{val: "invalid", def: 2, got: 2},
		{val: "", def: 2, got: 2},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("val=%q, def=%d, got=%d", tt.val, tt.def, tt.got), func(t *testing.T) {
			os.Setenv("TAKANAWA_TEST", tt.val)

			v := env.UintEnv("TAKANAWA_TEST", tt.def)
			if v != tt.got {
				t.Errorf("got %d, want %d", v, tt.got)
			}
		})
	}
}

func TestBoolEnv(t *testing.T) {
	tests := []struct {
		val string
		def bool
		got bool
	}{
		{val: "true", def: true, got: true},
		{val: "true", def: false, got: true},
		{val: "false", def: true, got: false},
		{val: "false", def: false, got: false},
		{val: "", def: true, got: true},
		{val: "", def: false, got: false},
		{val: "invalid", def: true, got: true},
		{val: "invalid", def: false, got: false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("val=%q, def=%v, got=%v", tt.val, tt.def, tt.got), func(t *testing.T) {
			os.Setenv("TAKANAWA_TEST", tt.val)

			v := env.BoolEnv("TAKANAWA_TEST", tt.def)
			if v != tt.got {
				t.Errorf("got %v, want %v", v, tt.got)
			}
		})
	}
}
