.PHONY: test
test: 
	go test -race ./...

.PHONY: protos
protos:
	protoc internal/proto/v1/*/*.proto \
		--go_out=. \
		--go-grpc_out=require_unimplemented_servers=false:. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.