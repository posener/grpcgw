package grpcgw

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Service interface {
	RegisterGRPC(*grpc.Server)
	RegisterGatewayEndpoints(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
