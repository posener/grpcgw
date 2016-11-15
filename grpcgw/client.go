package grpcgw

import (
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func NewGRPCConnection(addr string, credsPool *x509.CertPool) *grpc.ClientConn {
	var opts []grpc.DialOption
	creds := credentials.NewClientTLSFromCert(credsPool, addr)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	return conn
}
