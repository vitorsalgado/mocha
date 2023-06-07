package mgrpc

import (
	"context"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	pb "github.com/vitorsalgado/mocha/v3/mgrpc/internal/protobuf"
)

func TestServerStreaming(t *testing.T) {
	m := NewGRPCWithT(t, Setup().Service(&pb.TestService_ServiceDesc, &pb.UnimplementedTestServiceServer{}))
	m.MustStart()

	ctx := context.Background()
	list := []proto.Message{&pb.ListItem{Key: "a"}, &pb.ListItem{Key: "b"}}
	scope := m.MustMock(
		ServerStreamMethod("List").Reply(
			ServerStream(&pb.ListItem{}).
				Messages(list)))

	conn, err := grpc.Dial(m.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)

	client := pb.NewTestServiceClient(conn)
	res, err := client.List(ctx, &pb.ListReq{})

	require.NoError(t, err)
	require.NotNil(t, res)

	mu := sync.Mutex{}
	items := make([]*pb.ListItem, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			item, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					cancel()
					return
				}

				require.Fail(t, err.Error())
				cancel()
			}

			mu.Lock()
			items = append(items, item)
			mu.Unlock()
		}
	}()

	<-ctx.Done()

	require.Len(t, items, 2)
	require.Equal(t, items[0].String(), list[0].(*pb.ListItem).String())
	require.Equal(t, items[1].String(), list[1].(*pb.ListItem).String())
	require.True(t, scope.AssertCalled(t))
}

func TestServerStream_MessageReader(t *testing.T) {
	m := NewGRPCWithT(t, Setup().Service(&pb.TestService_ServiceDesc, &pb.UnimplementedTestServiceServer{}))
	m.MustStart()

	buf := strings.Builder{}
	m1 := &pb.ListItem{Key: "\ntest\n"}
	m1b, _ := prototext.Marshal(m1)
	m2 := &pb.ListItem{Key: "\ndev"}
	m2b, _ := prototext.Marshal(m2)

	buf.WriteString(string(m1b))
	buf.WriteString("\n")
	buf.WriteString(string(m2b))

	ctx := context.Background()
	scope := m.MustMock(ServerStreamMethod("List").Reply(
		ServerStreamT[*pb.ListItem]().
			Text(strings.NewReader(buf.String()))))

	conn, err := grpc.Dial(m.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)

	client := pb.NewTestServiceClient(conn)
	res, err := client.List(ctx, &pb.ListReq{})

	require.NoError(t, err)
	require.NotNil(t, res)

	mu := sync.Mutex{}
	items := make([]*pb.ListItem, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			item, err := res.Recv()
			if err != nil {
				if err == io.EOF {
					cancel()
					return
				}

				require.Fail(t, err.Error())
				cancel()
			}

			mu.Lock()
			items = append(items, item)
			mu.Unlock()
		}
	}()

	<-ctx.Done()

	require.Len(t, items, 2)
	require.Equal(t, m1, items[0])
	require.Equal(t, m2, items[1])
	require.True(t, scope.AssertCalled(t))
}
