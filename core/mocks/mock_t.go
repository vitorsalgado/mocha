package mocks

import "github.com/stretchr/testify/mock"

type FakeT struct{ mock.Mock }

func NewT() *FakeT {
	t := &FakeT{}
	t.On("Cleanup").Return()
	t.On("Helper").Return()
	t.On("Errorf", mock.Anything, mock.Anything).Return()

	return t
}

func (m *FakeT) Cleanup(_ func()) {
	m.Called()
}

func (m *FakeT) Helper() {
	m.Called()
}

func (m *FakeT) Errorf(format string, args ...any) {
	m.Called(format, args)
}
