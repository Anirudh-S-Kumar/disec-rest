.PHONY: server generate build

server:
	@go run cmd/server/main.go

generate:
	@protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

build:
	@go build -o bin/server cmd/server/main.go