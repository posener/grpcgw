all: build

deps:
	go get ./...

gen:
	go generate

fmt:
	go fmt ./...
	go vet ./...

test: gen fmt
	go test ./...
	
build: test
	go build -o bin/example
