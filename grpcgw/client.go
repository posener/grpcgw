package grpcgw

import (
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"log"
)

type client struct {
	Address string
	CertFile string
	Insecure bool
}

func NewGRPCConnection() *grpc.ClientConn {
	var opts []grpc.DialOption
	if Client.CertFile != "" {
		cert, err := ioutil.ReadFile(Client.CertFile)
		if err != nil {
			log.Fatalf("Failed reading pem file %s", Client.CertFile)
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(cert)
		creds := credentials.NewClientTLSFromCert(certPool, Client.Address)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		if Client.Insecure {
			opts = append(opts, grpc.WithInsecure())
		}
	}
	conn, err := grpc.Dial(Client.Address, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	return conn
}
