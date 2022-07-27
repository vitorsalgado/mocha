package reply

import "github.com/stretchr/testify/mock"

type mmock struct {
	mock.Mock
}

func (m *mmock) Hits() int {
	args := m.Called()
	return args.Int(0)
}
