package takanawa_test

import (
	"testing"

	"github.com/kou64yama/takanawa"
)

func TestParseForwarded(t *testing.T) {
	successTests := []struct {
		input string
		fwd   []takanawa.Forwarded
	}{
		{
			input: "for=\"_gazonk\"",
			fwd: []takanawa.Forwarded{
				takanawa.Forwarded{For: "_gazonk"},
			},
		},
		{
			input: "For=\"[2001:db8:cafe::17]:4711\"",
			fwd: []takanawa.Forwarded{
				takanawa.Forwarded{For: "[2001:db8:cafe::17]:4711"},
			},
		},
		{
			input: "for=192.0.2.60;proto=http;by=203.0.113.43",
			fwd: []takanawa.Forwarded{
				takanawa.Forwarded{By: "203.0.113.43", For: "192.0.2.60", Proto: "http"},
			},
		},
		{
			input: "for=192.0.2.43, for=198.51.100.17",
			fwd: []takanawa.Forwarded{
				takanawa.Forwarded{For: "192.0.2.43"},
				takanawa.Forwarded{For: "198.51.100.17"},
			},
		},
		{
			input: "for=192.0.2.43, for=198.51.100.17;by=203.0.113.60;proto=http;host=example.com",
			fwd: []takanawa.Forwarded{
				takanawa.Forwarded{For: "192.0.2.43"},
				takanawa.Forwarded{By: "203.0.113.60", For: "198.51.100.17", Host: "example.com", Proto: "http"},
			},
		},
	}
	errorTests := []struct{ input string }{
		{input: "invalid"},
	}

	for _, tt := range successTests {
		t.Run(tt.input, func(t *testing.T) {
			t.Logf("input: %s", tt.input)
			fwd, err := takanawa.ParseForwarded(tt.input)

			t.Logf("err: %v", err)
			if err != nil {
				t.Error(err)
			}

			for i, f := range fwd {
				t.Logf("fwd[%d]: %s", i, f.String())

				if fwd[i].By != tt.fwd[i].By {
					t.Errorf("got %q, want %q", fwd[i].By, tt.fwd[i].By)
				}
				if fwd[i].For != tt.fwd[i].For {
					t.Errorf("got %q, want %q", fwd[i].For, tt.fwd[i].For)
				}
				if fwd[i].Host != tt.fwd[i].Host {
					t.Errorf("got %q, want %q", fwd[i].Host, tt.fwd[i].Host)
				}
				if fwd[i].Proto != tt.fwd[i].Proto {
					t.Errorf("got %q, want %q", fwd[i].Proto, tt.fwd[i].Proto)
				}
			}
		})
	}
	for _, tt := range errorTests {
		t.Run(tt.input, func(t *testing.T) {
			t.Logf("input: %s", tt.input)
			_, err := takanawa.ParseForwarded(tt.input)

			t.Logf("err: %v", err)
			if err == nil {
				t.Error("no error")
			}
		})
	}
}
