package mgrpc

import (
	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type Interceptors struct {
	app *GRPCMockApp
}

type GRPCMockBuilder interface {
	Build(*GRPCMockApp) (*GRPCMock, error)
}

type GRPCMock struct {
	*foundation.BaseMock

	Reply any

	after              []matcher.OnAfterMockServed
	unaryExpectations  []*foundation.Expectation[*UnaryValueSelectorIn]
	streamExpectations []*foundation.Expectation[*StreamValueSelectorIn]
}

func newMock() *GRPCMock {
	return &GRPCMock{BaseMock: foundation.NewMock(), after: make([]matcher.OnAfterMockServed, 0)}
}

func (m *GRPCMock) getRef() string {
	if len(m.Name) == 0 {
		return m.ID
	}

	return m.Name
}
