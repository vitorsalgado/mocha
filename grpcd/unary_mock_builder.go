package grpcd

import (
	"github.com/vitorsalgado/mocha/v3/foundation"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"google.golang.org/grpc/metadata"
)

var (
	_ foundation.Builder[*GRPCMock, *GRPCMockApp] = (*UnaryMockBuilder)(nil)
	_ GRPCMockBuilder                             = (*UnaryMockBuilder)(nil)
)

type UnaryMockBuilder struct {
	m *GRPCMock
}

func UnaryMethod(method string) *UnaryMockBuilder {
	b := &UnaryMockBuilder{m: newMock()}
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&foundation.Expectation[*UnaryValueSelectorIn]{
			Target:        0,
			Key:           method,
			Matcher:       matcher.Contain(method),
			ValueSelector: unarySelectMethod,
			Weight:        10,
		},
	)

	return b
}

func (b *UnaryMockBuilder) Header(key string, m matcher.Matcher) *UnaryMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&foundation.Expectation[*UnaryValueSelectorIn]{
			Target:        0,
			Key:           key,
			Matcher:       m,
			ValueSelector: unarySelectHeader(key),
			Weight:        3,
		},
	)

	return b
}

func (b *UnaryMockBuilder) Field(path string, m matcher.Matcher) *UnaryMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&foundation.Expectation[*UnaryValueSelectorIn]{
			Target:        0,
			Key:           path,
			Matcher:       matcher.Field(path, m),
			ValueSelector: unarySelectBody,
			Weight:        3,
		},
	)
	return b
}

func (b *UnaryMockBuilder) Reply(r UnaryReply) *UnaryMockBuilder {
	b.m.Reply = r
	return b
}

func (b *UnaryMockBuilder) Build(_ *GRPCMockApp) (*GRPCMock, error) {
	return b.m, nil
}

func unarySelectMethod(r *UnaryValueSelectorIn) any {
	return r.Info.FullMethod
}

func unarySelectHeader(k string) UnaryValueSelector {
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

func unarySelectBody(r *UnaryValueSelectorIn) any {
	return r.RequestMessage
}
