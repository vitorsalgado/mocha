package dzgrpc

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func interceptError(format string, a ...any) error {
	return status.Error(codes.Internal, fmt.Errorf(format, a...).Error())
}
