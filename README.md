# grpcgw/gen tool

A convinience tool, wrapper of
[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway).
This script easily creates from a `proto` file the following
files:

- `.pb.go`: A grpc service for the service described in the proto.

- `.pb.gw.go`: A gateway service for the REST endpoints
  described proto.

- `json.go` file: Containing the json representation of a swagger
  file which described the REST endpoints. This file is to
  be used with the swagger-ui service.

## Install

`go get -u github.com/posener/grpcgw/gen`

## Usage

`gen -swagger-out <swagger go file> <proto file> [<proto file>...]`
