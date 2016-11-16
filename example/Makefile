all: build

CMD:=bin/example

deps:
	go get ./...

gen:
	go generate

fmt:
	go fmt ./...
	go vet ./...

test: gen fmt
	go test ./...
	
build: test bin/example

$(CMD):
	go build -o $@

certs: certs/server.pem

certs/server.key:
	mkdir -p certs
	openssl genrsa -out $@ 2048

certs/server.pem: certs/server.key
	@echo "For 'Common Name' use 'localhost:10000'"
	openssl req -new -x509 -key $< -out $@ -days 3650

clean:
	rm -rf bin certs swaggers

# Testing targets:
# run will start a server
# run-client-echo will start a client and send an
# echo message to the server

run: $(CMD) certs
	$(CMD) serve --key certs/server.key --crt certs/server.pem --swaggers ./swagger

run-client-echo: $(CMD) certs
	$(CMD) send --crt certs/server.pem echo hi

run-rest-echo:
	curl --insecure --cacert certs/server.pem -H 'Content-Type: application/json' -H 'Accept: application/json' -d '{"value": "string"}' 'https://localhost:10000/v1/echo'; echo
