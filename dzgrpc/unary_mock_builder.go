package dzgrpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/vitorsalgado/mocha/v3/dzstd"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

var (
	_ dzstd.Builder[*GRPCMock, *GRPCMockApp] = (*UnaryMockBuilder)(nil)
	_ GRPCMockBuilder                        = (*UnaryMockBuilder)(nil)
)

type UnaryMockBuilder struct {
	m *GRPCMock
}

func UnaryMethod(method string) *UnaryMockBuilder {
	b := &UnaryMockBuilder{m: newMock()}
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&dzstd.Expectation[*UnaryValueSelectorIn]{
			TargetDescription: describeTarget(targetMethod, method),
			Matcher:           matcher.Contain(method),
			ValueSelector:     unarySelectMethod,
			Weight:            10,
		},
	)

	return b
}

func (b *UnaryMockBuilder) Header(key string, m matcher.Matcher) *UnaryMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&dzstd.Expectation[*UnaryValueSelectorIn]{
			TargetDescription: describeTarget(targetHeader, key),
			Matcher:           m,
			ValueSelector:     unarySelectHeader(key),
			Weight:            3,
		},
	)

	return b
}

func (b *UnaryMockBuilder) Field(path string, m matcher.Matcher) *UnaryMockBuilder {
	b.m.unaryExpectations = append(
		b.m.unaryExpectations,
		&dzstd.Expectation[*UnaryValueSelectorIn]{
			TargetDescription: describeTarget(targetBodyField, path),
			Matcher:           matcher.Field(path, m),
			ValueSelector:     unarySelectBody,
			Weight:            3,
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

func unarySelectMethod(_ context.Context, r *UnaryValueSelectorIn) any {
	return r.Info.FullMethod
}

func unarySelectHeader(k string) UnaryValueSelector {
	return func(ctx context.Context, r *UnaryValueSelectorIn) any {
		md, ok := metadata.FromIncomingContext(ctx)
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

func unarySelectBody(_ context.Context, r *UnaryValueSelectorIn) any {
	return r.RequestMessage
}
