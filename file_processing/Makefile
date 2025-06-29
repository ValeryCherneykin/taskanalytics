LOCAL_BIN=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-file-processing-api

generate-file-processing-api:
	mkdir -p pkg/file_processing_v1
	protoc --proto_path=api/file_processing_v1 \
		--go_out=pkg/file_processing_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--go-grpc_out=pkg/file_processing_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		api/file_processing_v1/file_processing.proto

grpc-load-test:
	ghz \
		--proto api/file_processing_v1/file_processing.proto \
		--call file_processing_v1.FileProcessingService.GetFileMetadata \
		--data '{"file_id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto api/file_processing_v1/file_processing.proto \
		--call file_processing_v1.FileProcessingService.GetFileMetadata \
		--data '{"file_id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051
