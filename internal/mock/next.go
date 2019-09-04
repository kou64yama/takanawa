package mock

type NextMock struct {
	CalledN int
}

func NewNextMock() *NextMock {
	return &NextMock{}
}

func (mock *NextMock) Mock() func() {
	return func() {
		mock.CalledN++
	}
}
