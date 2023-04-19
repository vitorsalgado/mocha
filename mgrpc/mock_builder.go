package mgrpc

import (
	"google.golang.org/grpc/metadata"

	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

var _ foundation.Builder[*GRPCMock, *GRPCMockApp] = (*GRPCMockBuilder)(nil)

type GRPCMockBuilder struct {
	m *GRPCMock
}

type baseGRPCMockBuilder[T any] struct {
}

func ForMethod(method string) *GRPCMockBuilder {
	b := &GRPCMockBuilder{m: newMock()}
	return b.Method(method)
}

func (b *GRPCMockBuilder) Method(method string) *GRPCMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&foundation.Expectation[*UnaryValueSelectorIn]{
			Target:        0,
			Key:           method,
			Matcher:       matcher.Contain(method),
			ValueSelector: selectMethod,
			Weight:        10,
		},
	)

	return b
}

func (b *GRPCMockBuilder) Header(key string, m matcher.Matcher) *GRPCMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&foundation.Expectation[*UnaryValueSelectorIn]{
			Target:        0,
			Key:           key,
			Matcher:       m,
			ValueSelector: selectHeader(key),
			Weight:        3,
		},
	)

	return b
}

func (b *GRPCMockBuilder) Field(path string, m matcher.Matcher) *GRPCMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&foundation.Expectation[*UnaryValueSelectorIn]{
			Target:        0,
			Key:           path,
			Matcher:       matcher.Field(path, m),
			ValueSelector: selectBody,
			Weight:        3,
		},
	)
	return b
}

func (b *GRPCMockBuilder) Reply(r UnaryReply) *GRPCMockBuilder {
	b.m.Reply = r
	return b
}

func (b *GRPCMockBuilder) Stream(r ServerStreamReply) *GRPCMockBuilder {
	b.m.Reply = r
	return b
}

func (b *GRPCMockBuilder) Build(_ *GRPCMockApp) (*GRPCMock, error) {
	return b.m, nil
}

func selectMethod(r *UnaryValueSelectorIn) any {
	return r.Info.FullMethod
}

func selectHeader(k string) UnaryValueSelector {
	return func(r *UnaryValueSelectorIn) any {
		md, ok := metadata.FromIncomingContext(r.Context)
		if !ok {
			return nil
		}

		v := md.Get(k)
		if len(v) == 1 {
			return v[0]
		}

		return v
	}
}

func selectBody(r *UnaryValueSelectorIn) any {
	return r.RequestMessage
}
