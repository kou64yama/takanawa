package util_test

import "testing"

import "github.com/kou64yama/takanawa/internal/util"

func TestSplitAndTrimSpace(t *testing.T) {
	tests := []struct {
		str  string
		want []string
	}{
		{str: "foo,bar", want: []string{"foo", "bar"}},
		{str: "foo, bar", want: []string{"foo", "bar"}},
		{str: "foo", want: []string{"foo"}},
		{str: " ", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			t.Helper()

			got := util.SplitAndTrimSpace(tt.str, ",")
			if len(got) != len(tt.want) {
				t.Errorf("got %d, want %d", len(got), len(tt.want))
				return
			}
			for i := 0; i < len(got); i++ {
				if got[i] != tt.want[i] {
					t.Errorf("got %q want %q", got[i], tt.want[i])
				}
			}
		})
	}
}
