package example

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func (s *service) RegisterGRPC(grpcServer *grpc.Server) {
	RegisterEchoServiceServer(grpcServer, s)
}
func (s *service) RegisterGatewayEndpoints(ctx context.Context, gwmux *runtime.ServeMux, grpcEndpointAddr string, opts []grpc.DialOption) error {
	return RegisterEchoServiceHandlerFromEndpoint(ctx, gwmux, grpcEndpointAddr, opts)
}
