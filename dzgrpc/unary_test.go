package dzgrpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/vitorsalgado/mocha/v3/dzgrpc/internal/protobuf"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestUnary(t *testing.T) {
	m := NewGRPCWithT(t)
	pb.RegisterTestServiceServer(m.Server(), pb.UnimplementedTestServiceServer{})
	m.MustStart()
	m.MustMock(UnaryMethod("Greetings").
		Header("h1", Eq("v1")).
		Field("message", Eq("hi")).
		Reply(Unary().
			Header("k1", "v1", "v2").
			Trailer("tk1", "tv1", "tv2").
			Message(&pb.HiRes{Message: "bye"})))

	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "h1", "v1")
	conn, err := grpc.Dial(m.Addr(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	require.NoError(t, err)

	client := pb.NewTestServiceClient(conn)
	header := make(metadata.MD)
	trailer := make(metadata.MD)
	res, err := client.Greetings(ctx, &pb.HiReq{Message: "hi"}, grpc.Header(&header), grpc.Trailer(&trailer))

	require.NoError(t, err)
	require.Equal(t, 2, header.Len()) // includes the content-type
	require.Equal(t, []string{"v1", "v2"}, header.Get("k1"))
	require.Equal(t, 1, trailer.Len())
	require.Equal(t, []string{"tv1", "tv2"}, trailer.Get("tk1"))
	require.Equal(t, "bye", res.Message)
}
