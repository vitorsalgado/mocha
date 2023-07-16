package testmock

import "github.com/stretchr/testify/mock"

type FakeT struct{ mock.Mock }

func NewFakeT() *FakeT {
	t := &FakeT{}
	t.On("Helper").Return()
	t.On("Logf", mock.Anything, mock.Anything).Return()
	t.On("Errorf", mock.Anything, mock.Anything).Return()
	t.On("FailNow").Return()

	return t
}

func (m *FakeT) Helper() {
	m.Called()
}

func (m *FakeT) Logf(format string, args ...any) {
	m.Called(format, args)
}

func (m *FakeT) Error(args ...any) {
	m.Called(args...)
}

func (m *FakeT) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *FakeT) Cleanup(_ func()) { /* not necessary */ }
