package grpcd

import (
	"github.com/vitorsalgado/mocha/v3/lib"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"google.golang.org/grpc/metadata"
)

var (
	_ lib.Builder[*GRPCMock, *GRPCMockApp] = (*ServerStreamMockBuilder)(nil)
	_ GRPCMockBuilder                      = (*ServerStreamMockBuilder)(nil)
)

type ServerStreamMockBuilder struct {
	m *GRPCMock
}

func ServerStreamMethod(method string) *ServerStreamMockBuilder {
	b := &ServerStreamMockBuilder{m: newMock()}
	b.m.streamExpectations = append(
		b.m.streamExpectations,
		&lib.Expectation[*StreamValueSelectorIn]{
			Target:        0,
			Key:           method,
			Matcher:       matcher.Contain(method),
			ValueSelector: streamSelectMethod,
			Weight:        10,
		},
	)

	return b
}

func (b *ServerStreamMockBuilder) Header(key string, m matcher.Matcher) *ServerStreamMockBuilder {
	b.m.streamExpectations = append(
		b.m.streamExpectations,
		&lib.Expectation[*StreamValueSelectorIn]{
			Target:        0,
			Key:           key,
			Matcher:       m,
			ValueSelector: streamSelectHeader(key),
			Weight:        3,
		},
	)

	return b
}

func (b *ServerStreamMockBuilder) Field(path string, m matcher.Matcher) *ServerStreamMockBuilder {
	b.m.streamExpectations = append(
		b.m.streamExpectations,
		&lib.Expectation[*StreamValueSelectorIn]{
			Target:        0,
			Key:           path,
			Matcher:       matcher.Field(path, m),
			ValueSelector: streamSelectBody,
			Weight:        3,
		},
	)
	return b
}

func (b *ServerStreamMockBuilder) Reply(r ServerStreamReply) *ServerStreamMockBuilder {
	b.m.Reply = r
	return b
}

func (b *ServerStreamMockBuilder) Build(_ *GRPCMockApp) (*GRPCMock, error) {
	return b.m, nil
}

func streamSelectMethod(r *StreamValueSelectorIn) any {
	return r.Info.FullMethod
}

func streamSelectHeader(k string) StreamValueSelector {
	return func(r *StreamValueSelectorIn) any {
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

func streamSelectBody(r *StreamValueSelectorIn) any {
	return r.RequestMessage
}