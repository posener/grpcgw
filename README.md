# grpcgw

References:

* [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
* CoreOS blog post [Take a REST with HTTP/2, Protobufs, and Swagger](https://coreos.com/blog/gRPC-protobufs-swagger.html)
* Phillip's [grpc gateway example](https://github.com/philips/grpc-gateway-example)

This is a convenience library and code-generation tool for
creating a grpc server with a REST gateway and Swagger-ui
description from `proto` files.

## grpcgw/grpcgw

grpcgw is a library extracted and modified from the example in
Phillip's [github page](https://github.com/philips/grpc-gateway-example).
It enables the user to write a service without effort and minimal
boiler plate code.

### Install

`go get -u github.com/posener/grpcgw/gen`

### grpcgw/example

This directory contains a basic example of echo server.
Lets examine the code:

* `example/service`

  This directory contains the service main business.

  - `service.proto`: The RPC/Rest API of the service, represented
    with protobuf and the grpc-gw extension. There you could find
    the definition of the `EchoMessage` message, the `EchoService`
    interface, and the `Echo` method of that interface.

  - `service.go`: The business logic of the service. There you could
    find the `service` struct, implementing the `EchoService` interface,
    and the `Echo` method, which is the implementation of the interface.

    You could also find there functions that creates new service and
    new client.

  - `register.go`: Methods implementing the `grpcgw.Service` interface.
    Those methods are used to register the `service` when the server
    starts.

  - After invoking the `go genenrate` command in the `example` project,
    two more files will appear in this folder: `service.pb.go` which
    is the implementation of the grpc server, and `service.pb.gw.go`,
    which is the implementation of the REST gateway.
    In the `example` directory, a `swagger` directory will appear
    with auto-generated swagger json.

* `generate.go`: A file containing the auto-generation script.

* `main.go` and `cmd/*` are auto generated with
  (cobra)[https://github.com/spf13/cobra].

  - grpcgw provides basic commands to start and stop the client.
    In the `main.go` we call
    `grpcgw.AddCommands(cmd.RootCmd, example.NewService())`, which
    customize the `serve` command to use the `example` service.

  - `cmd/echo.go` is the echo sub-command, notice that it uses
    `grpcgw.NewGRPCConnection()` to get connection to the defined
    server.

* `Makefile`: can give an impression about how to create your own service
  and how to run it. You can play with the `run` target, which will also
  create certificates for your server and run it. The `run-client-echo`
  target will send an RPC to the running server. The `run-rest-echo`
  target will send the same request through REST.

* When the server is running, browse to
  [https://localhost:10000/swagger-ui](https://localhost:10000/swagger-ui]),
  then, enter in the text box:
  [https://localhost:10000/swaggers/service.swagger.json](https://localhost:10000/swaggers/service.swagger.json).

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
