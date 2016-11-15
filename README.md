# grpcgw

References:

* [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
* CoreOS blog post [Take a REST with HTTP/2, Protobufs, and Swagger](https://coreos.com/blog/gRPC-protobufs-swagger.html)
* Phillip's [grpc gateway example](https://github.com/philips/grpc-gateway-example)

This is a convenience library and code-generation tool for
creating a grpc server with a REST gateway and Swagger-ui
description from `proto` files.

## grpcgw/grpcgw

## grpcgw/gen

This script easily creates from a `proto` file the following
files:

- `.pb.go`: A grpc service for the service described in the proto.

- `.pb.gw.go`: A gateway service for the REST endpoints
  described proto.

- `json.go` file: Containing the json representation of a swagger
  file which described the REST endpoints. This file is to
  be used with the swagger-ui service.

### Install

`go get -u github.com/posener/grpcgw/gen`

### Usage

`gen -swagger-out <swagger go file> <proto file> [<proto file>...]`
