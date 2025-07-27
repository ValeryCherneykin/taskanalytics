LOCAL_BIN := $(CURDIR)/bin
PROTO_DIR := api/task_distribution_v1
OUT_DIR := pkg/task_distribution_v1

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-task-distribution-api

generate-task-distribution-api:
	mkdir -p pkg/task_distribution_v1
	protoc --proto_path $(PROTO_DIR) \
	--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/task_distribution_v1/task_distribution.proto


