package util_test

import (
	"github.com/kou64yama/takanawa/internal/util"
	"os"
	"testing"
)

func TestDefaultHost(t *testing.T) {
	t.Run("HOST=", func(t *testing.T) {
		t.Helper()

		p := os.Getenv("HOST")
		defer os.Setenv("HOST", p)

		os.Setenv("HOST", "")
		got := util.DefaultHost()
		if got != "127.0.0.1" && got != "[::1]" {
			t.Errorf("got %q, want \"127.0.0.1\" or \"[::1]\"", got)
		}
	})
	t.Run("HOST=0.0.0.0", func(t *testing.T) {
		t.Helper()

		p := os.Getenv("HOST")
		defer os.Setenv("HOST", p)

		os.Setenv("HOST", "0.0.0.0")
		want := "0.0.0.0"
		got := util.DefaultHost()
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestDefaultPort(t *testing.T) {
	tests := []struct {
		port string
		want uint
	}{
		{port: "", want: 5000},
		{port: "3000", want: 3000},
		{port: "-1", want: 5000},
		{port: "0", want: 0},
		{port: "65535", want: 65535},
		{port: "65536", want: 5000},
		{port: "abc", want: 5000},
	}

	for _, tt := range tests {
		t.Run("PORT="+tt.port, func(t *testing.T) {
			t.Helper()

			p := os.Getenv("PORT")
			defer os.Setenv("PORT", p)

			os.Setenv("PORT", tt.port)
			got := util.DefaultPort()
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
