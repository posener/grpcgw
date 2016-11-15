all: grpcgw/swaggerui.go certs/server.pem

fmt:
	go fmt ./...
	go vet ./...

grpcgw/swaggerui.go: swagger-ui
	mkdir -p $(dir $@)
	go-bindata -nocompress -o $@ -prefix ${PWD} -pkg grpcgw $</...

swagger-ui:
	curl -sSL https://github.com/swagger-api/swagger-ui/archive/v2.2.6.zip > /tmp/swagger-ui.zip
	mkdir -p $@
	unzip /tmp/swagger-ui.zip swagger-ui-*/dist/* -d $@
	mv swagger-ui/swagger-ui-*/dist/* swagger-ui/
	rm -r swagger-ui/swagger-ui-*

certs: certs/server.pem

certs/server.key:
	mkdir -p certs
	openssl genrsa -out $@ 2048

certs/server.pem: certs/server.key
	openssl req -new -x509 -key $< -out $@ -days 3650

clean:
	rm -r grpcgw/swaggerui.go swagger-ui
