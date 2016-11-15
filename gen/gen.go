package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	_ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"
	"github.com/posener/grpcgw/gen/swagger"
)

func main() {
	swaggerOut := flag.String("swagger-out", "src/swagger/swagger.go", "Output go file for swagger json content")
	flag.Parse()
	protos := flag.Args()

	if len(protos) == 0 {
		log.Fatal("No proto files provided")
	}

	includes := getIncludes()

	cmds := []*exec.Cmd{}

	for _, proto := range protos {
		log.Printf("Generating code for %s", proto)
		cmds = append(cmds, generateSwagger(proto, includes, *swaggerOut))
		cmds = append(cmds, generateGRPC(proto, includes))
		cmds = append(cmds, generateGateway(proto, includes))
	}

	log.Print("Waiting for generating code...")
	for _, cmd := range cmds {
		cmd.Wait()
	}

	log.Print("Creating go file from swagger jsons...")
	swagger.GenerateSwaggerGoFile(*swaggerOut)

	log.Print("Finished successfully")
}

// Generate a swagger json file in the directory of the swaggerOutFile
func generateSwagger(proto string, includes []string, swaggerOutFile string) *exec.Cmd {
	swaggerOutDir := filepath.Dir(swaggerOutFile)
	out := fmt.Sprintf("--swagger_out=logtostderr=true:%s", swaggerOutDir)
	cmd := exec.Command("protoc", append(includes, out, proto)...)
	err := cmd.Start()
	if err != nil {
		log.Panicf("Failed running protoc swagger %s", err)
	}
	return cmd
}

// Generate grpc server code
func generateGRPC(proto string, includes []string) *exec.Cmd {
	out := "--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:."
	cmd := exec.Command("protoc", append(includes, out, proto)...)
	err := cmd.Start()
	if err != nil {
		log.Panicf("Failed running protoc grpc: %s", err)
	}
	return cmd
}

// Generate gateway server code
func generateGateway(proto string, includes []string) *exec.Cmd {
	out := "--grpc-gateway_out=logtostderr=true:."
	cmd := exec.Command("protoc", append(includes, out, proto)...)
	err := cmd.Start()
	if err != nil {
		log.Panicf("Failed running protoc gateway: %s", err)
	}
	return cmd
}

func getIncludes() []string {
	goPath, ok := os.LookupEnv("GOPATH")
	if !ok {
		log.Panic("GOPATH must be defined!")
	}
	includes := []string{
		"/usr/local/include",
		".",
		filepath.Join(goPath, "src"),
		filepath.Join(goPath, "src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis"),
	}
	for i := range includes {
		includes[i] = "-I" + includes[i]
	}
	return includes
}
