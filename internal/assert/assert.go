package assert

import "testing"

type Assertions struct {
	t *testing.T
}

func NewAssertions(t *testing.T) *Assertions {
	return &Assertions{t: t}
}

func (ass *Assertions) AssertEquals(got, want interface{}) bool {
	b := got == want
	if !b {
		ass.t.Errorf("got %v, want %v", got, want)
	}
	return b
}

func (ass *Assertions) AssertTrue(got bool) bool {
	if !got {
		ass.t.Errorf("got %v, want true", got)
	}
	return got
}
