package grpcd

import (
	"github.com/vitorsalgado/mocha/v3/lib"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type Interceptors struct {
	app *GRPCMockApp
}

type GRPCMockBuilder interface {
	Build(*GRPCMockApp) (*GRPCMock, error)
}

type GRPCMock struct {
	*lib.BaseMock

	Reply any

	after              []matcher.OnAfterMockServed
	unaryExpectations  []*lib.Expectation[*UnaryValueSelectorIn]
	streamExpectations []*lib.Expectation[*StreamValueSelectorIn]
}

func newMock() *GRPCMock {
	return &GRPCMock{BaseMock: lib.NewMock(), after: make([]matcher.OnAfterMockServed, 0)}
}

func (m *GRPCMock) getRef() string {
	if len(m.Name) == 0 {
		return m.ID
	}

	return m.Name
}
