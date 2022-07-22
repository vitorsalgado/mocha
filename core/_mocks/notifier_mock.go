package mocks

import "github.com/stretchr/testify/mock"

type FakeNotifier struct{ mock.Mock }

func NewFakeNotifier() *FakeNotifier {
	t := &FakeNotifier{}
	t.On("Helper").Return()
	t.On("Logf", mock.Anything, mock.Anything).Return()
	t.On("Errorf", mock.Anything, mock.Anything).Return()
	t.On("FailNow").Return()

	return t
}

func (m *FakeNotifier) Helper() {
	m.Called()
}

func (m *FakeNotifier) Logf(format string, args ...any) {
	m.Called(format, args)
}

func (m *FakeNotifier) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *FakeNotifier) FailNow() {
	m.Called()
}
