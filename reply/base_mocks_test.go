package reply

import "github.com/stretchr/testify/mock"

type mMock struct {
	mock.Mock
}

func (m *mMock) Hits() int {
	args := m.Called()
	return args.Int(0)
}
