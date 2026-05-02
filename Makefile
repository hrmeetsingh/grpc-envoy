.PHONY: proto test build up down

proto:
	protoc --go_out=services/user/internal/adapter/grpcserver/pb --go_opt=paths=source_relative \
		--go-grpc_out=services/user/internal/adapter/grpcserver/pb --go-grpc_opt=paths=source_relative \
		-I proto proto/user/v1/user.proto
	protoc --go_out=services/greeter/internal/adapter/grpcserver/pb --go_opt=paths=source_relative \
		--go-grpc_out=services/greeter/internal/adapter/grpcserver/pb --go-grpc_opt=paths=source_relative \
		-I proto proto/greeter/v1/greeter.proto

test:
	cd services/user && go test ./...
	cd services/greeter && go test ./...
	cd services/auth && go test ./...

build:
	cd services/user && go build -o ../../bin/user-service ./cmd/server
	cd services/greeter && go build -o ../../bin/greeter-service ./cmd/server
	cd services/auth && go build -o ../../bin/auth-service ./cmd/server

up:
	docker compose up --build -d

down:
	docker compose down
