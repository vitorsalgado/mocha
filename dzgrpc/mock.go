package dzgrpc

import (
	"github.com/vitorsalgado/mocha/v3/dzstd"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

type Interceptors struct {
	app *GRPCMockApp
}

type GRPCMockBuilder interface {
	Build(*GRPCMockApp) (*GRPCMock, error)
}

type GRPCMock struct {
	*dzstd.BaseMock

	Reply any

	after              []matcher.OnAfterMockServed
	unaryExpectations  []*dzstd.Expectation[*UnaryValueSelectorIn]
	streamExpectations []*dzstd.Expectation[*StreamValueSelectorIn]
}

func newMock() *GRPCMock {
	return &GRPCMock{BaseMock: dzstd.NewMock(), after: make([]matcher.OnAfterMockServed, 0)}
}

func (m *GRPCMock) getRef() string {
	if len(m.Name) == 0 {
		return m.ID
	}

	return m.Name
}

const (
	targetMethod    = "Method"
	targetHeader    = "Header"
	targetBody      = "Body"
	targetBodyField = "Field"
)

func describeTarget(target, key string) string {
	if len(key) == 0 {
		return target
	}

	return target + "(" + key + ")"
}
