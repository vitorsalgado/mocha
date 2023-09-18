package dzgrpc

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vitorsalgado/mocha/v3/dzstd"
)

type UnaryValueSelector func(ctx context.Context, in *UnaryValueSelectorIn) any

type UnaryValueSelectorIn struct {
	RequestMessage any
	Info           *grpc.UnaryServerInfo
}

type UnaryRequestValues struct {
	RequestMessage any
	Info           *grpc.UnaryServerInfo
	RawBody        any
	App            *GRPCMockApp
}

type UnaryResponse struct {
	Header  metadata.MD
	Trailer metadata.MD
	Message any
	Status  *status.Status
}

func (in *Interceptors) UnaryInterceptor(
	ctx context.Context,
	reqMsg interface{},
	info *grpc.UnaryServerInfo,
	_ grpc.UnaryHandler,
) (any, error) {
	b, err := json.Marshal(reqMsg)
	if err != nil {
		return nil, err
	}

	rawBody := string(b)
	description := dzstd.Results{Buf: make([]string, 0)}
	result, err := dzstd.FindMockForRequest(
		ctx,
		in.app.storage,
		func(m *GRPCMock) []*dzstd.Expectation[*UnaryValueSelectorIn] { return m.unaryExpectations },
		&UnaryValueSelectorIn{rawBody, info},
		&description,
		&dzstd.FindOptions{FailFast: false},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("unary: error finding mock for request: %w", err).Error())
	}

	if !result.Pass {
		return nil, interceptError("unary: request was not matched with any mock")
	}

	mock := result.Matched
	reply, ok := mock.Reply.(UnaryReply)
	if !ok {
		return nil, interceptError("unary: mock %s must implement an unary reply: got %T", mock.getRef(), mock.Reply)
	}

	res, err := reply.Build(ctx, &UnaryRequestValues{reqMsg, info, rawBody, in.app})
	if err != nil {
		return nil, err
	}

	mock.Inc()

	return res, nil
}

type UnaryReply interface {
	Build(ctx context.Context, rv *UnaryRequestValues) (any, error)
}

type BuiltInUnaryReply struct {
	response *UnaryResponse
}

func Unary() *BuiltInUnaryReply {
	return &BuiltInUnaryReply{response: &UnaryResponse{
		Header:  make(metadata.MD),
		Trailer: make(metadata.MD),
	}}
}

func (u *BuiltInUnaryReply) Status(st *status.Status) *BuiltInUnaryReply {
	u.response.Status = st
	return u
}

func (u *BuiltInUnaryReply) Header(k string, v ...string) *BuiltInUnaryReply {
	u.response.Header.Append(k, v...)
	return u
}

func (u *BuiltInUnaryReply) Trailer(k string, v ...string) *BuiltInUnaryReply {
	u.response.Trailer.Append(k, v...)
	return u
}

func (u *BuiltInUnaryReply) Message(msg any) *BuiltInUnaryReply {
	u.response.Message = msg
	return u
}

func (u *BuiltInUnaryReply) Build(ctx context.Context, _ *UnaryRequestValues) (any, error) {
	err := grpc.SendHeader(ctx, u.response.Header)
	if err != nil {
		return nil, err
	}

	err = grpc.SetTrailer(ctx, u.response.Trailer)
	if err != nil {
		return nil, err
	}

	if u.response.Status != nil {
		return nil, u.response.Status.Err()
	}

	return u.response.Message, nil
}
