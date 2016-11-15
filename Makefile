all: grpcgw/swaggerui.go

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

clean:
	rm -r grpcgw/swaggerui.go swagger-ui
