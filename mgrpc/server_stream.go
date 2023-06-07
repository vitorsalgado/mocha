package mgrpc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	"github.com/vitorsalgado/mocha/v3/foundation"
)

type GRPCStreamMock struct {
	*GRPCMock
}

func (s *GRPCStreamMock) GetExpectations() []*foundation.Expectation[*StreamValueSelectorIn] {
	return s.streamExpectations
}

type StreamValueSelector func(in *StreamValueSelectorIn) any

type StreamValueSelectorIn struct {
	Context        context.Context
	RequestMessage any
	Info           *grpc.StreamServerInfo
}

type StreamRequestValues struct {
	Context        context.Context
	RequestMessage any
	ServerStream   grpc.ServerStream
	Info           *grpc.StreamServerInfo
	App            *GRPCMockApp
}

type StreamType int

const (
	StreamTypeJSON StreamType = iota
	StreamTypeText
)

type StreamResponse struct {
	Header     metadata.MD
	Trailer    metadata.MD
	MsgType    proto.Message
	Stream     any
	StreamType StreamType
}

type ServerStreamReply interface {
	Build(values *StreamRequestValues) error
}

type BuiltInServerStreamReply struct {
	response *StreamResponse
}

func (r *BuiltInServerStreamReply) Messages(arr []proto.Message) *BuiltInServerStreamReply {
	r.response.Stream = arr
	return r
}

func (r *BuiltInServerStreamReply) AnyMessages(arr []any) *BuiltInServerStreamReply {
	r.response.Stream = arr
	return r
}

func (r *BuiltInServerStreamReply) Text(reader io.Reader) *BuiltInServerStreamReply {
	r.response.Stream = reader
	r.response.StreamType = StreamTypeText
	return r
}

func (r *BuiltInServerStreamReply) JSON(reader io.Reader) *BuiltInServerStreamReply {
	r.response.Stream = reader
	r.response.StreamType = StreamTypeJSON
	return r
}

func (r *BuiltInServerStreamReply) Build(values *StreamRequestValues) error {
	err := grpc.SendHeader(values.Context, r.response.Header)
	if err != nil {
		return err
	}

	err = grpc.SetTrailer(values.Context, r.response.Trailer)
	if err != nil {
		return err
	}

	switch s := r.response.Stream.(type) {
	case io.Reader:
		scan := bufio.NewScanner(s)
		msgType := reflect.New(reflect.TypeOf(r.response.MsgType).Elem())

		for scan.Scan() {
			msg := msgType.Interface().(proto.Message)
			err := r.decode(scan.Bytes(), msg)
			if err != nil {
				return err
			}

			err = values.ServerStream.SendMsg(msg)
			if err != nil {
				return err
			}
		}

	default:
		t := reflect.TypeOf(s)
		switch t.Kind() {
		case reflect.Array, reflect.Slice:
			v := reflect.ValueOf(s)
			for i := 0; i < v.Len(); i++ {
				err := values.ServerStream.SendMsg(v.Index(i).Interface())
				if err != nil {
					return err
				}
			}

		default:
			return status.Errorf(
				codes.Internal,
				"stream: reply stream type %T is not allowed. expected: io.Reader | array/slice",
				s,
			)
		}
	}

	return nil
}

func (r *BuiltInServerStreamReply) decode(b []byte, m proto.Message) error {
	switch r.response.StreamType {
	case StreamTypeText:
		return prototext.Unmarshal(b, m)
	case StreamTypeJSON:
		return protojson.Unmarshal(b, m)
	}

	return fmt.Errorf("stream: unexpected stream type %d", r.response.StreamType)
}

func ServerStreamT[T proto.Message]() *BuiltInServerStreamReply {
	var msgType T

	return &BuiltInServerStreamReply{response: &StreamResponse{
		Header:  make(metadata.MD),
		Trailer: make(metadata.MD),
		MsgType: msgType,
	}}
}

func ServerStream(messageType proto.Message) *BuiltInServerStreamReply {
	return &BuiltInServerStreamReply{response: &StreamResponse{
		Header:  make(metadata.MD),
		Trailer: make(metadata.MD),
		MsgType: messageType,
	}}
}

func (in *Interceptors) StreamInterceptor(
	reqMsg interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	_ grpc.StreamHandler,
) error {
	b, err := json.Marshal(reqMsg)
	if err != nil {
		return err
	}

	rawBody := string(b)
	mocks := in.app.storage.GetAll()
	wrappedMocks := make([]*GRPCStreamMock, len(mocks))
	for i, v := range mocks {
		wrappedMocks[i] = &GRPCStreamMock{v}
	}

	result := foundation.FindMockForRequest(wrappedMocks, &StreamValueSelectorIn{
		Context:        stream.Context(),
		RequestMessage: rawBody,
		Info:           info,
	})

	if !result.Pass {
		return status.Error(codes.NotFound, "stream: request was not matched with any mock")
	}

	mock := result.Matched
	reply, ok := mock.Reply.(ServerStreamReply)
	if !ok {
		return status.Errorf(
			codes.Unknown,
			"stream: mock %s must implement an server stream reply: got %T",
			mock.getRef(),
			mock.Reply,
		)
	}

	err = reply.Build(&StreamRequestValues{stream.Context(), rawBody, stream, info, in.app})
	if err != nil {
		return fmt.Errorf("stream: failed to reply: %w", err)
	}

	mock.Inc()

	return nil
}
