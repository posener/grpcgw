package example

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type service struct{}

func NewService() *service {
	return &service{}
}

func NewClient(conn *grpc.ClientConn) EchoServiceClient {
	return NewEchoServiceClient(conn)
}

func (m *service) Echo(c context.Context, s *EchoMessage) (*EchoMessage, error) {
	log.Printf("rpc request Echo(%q)\n", s.Value)
	return s, nil
}
