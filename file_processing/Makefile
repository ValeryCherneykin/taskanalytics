include .env

LOCAL_BIN := $(CURDIR)/bin
PROTO_DIR := api/file_processing_v1
OUT_DIR := pkg/file_processing_v1
VENDOR_PROTO_DIR := vendor.protogen

LOCAL_MIGRATION_DIR := $(MIGRATION_DIR)
LOCAL_MIGRATION_DSN := "host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.1
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

generate: generate-file-processing-api

generate-file-processing-api: vendor-proto
	mkdir -p $(OUT_DIR)
	mkdir -p docs/swagger
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(VENDOR_PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=file_processing \
		--plugin=protoc-gen-go=$(LOCAL_BIN)/protoc-gen-go \
		--plugin=protoc-gen-go-grpc=$(LOCAL_BIN)/protoc-gen-go-grpc \
		--plugin=protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
		--plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
		$(PROTO_DIR)/file_processing.proto

vendor-proto:
	@if [ ! -d $(VENDOR_PROTO_DIR)/google ]; then \
		echo "Cloning googleapis for google/api/*.proto..."; \
		git clone https://github.com/googleapis/googleapis $(VENDOR_PROTO_DIR)/googleapis &&\
		mkdir -p $(VENDOR_PROTO_DIR)/google &&\
		mv $(VENDOR_PROTO_DIR)/googleapis/google/api $(VENDOR_PROTO_DIR)/google &&\
		rm -rf $(VENDOR_PROTO_DIR)/googleapis ;\
	fi

local-migration-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

grpc-load-test:
	ghz \
		--proto $(PROTO_DIR)/file_processing.proto \
		--call fileprocessing_v1.FileProcessingService.GetFileMetadata \
		--data '{"file_id": 1}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051

grpc-error-load-test:
	ghz \
		--proto $(PROTO_DIR)/file_processing.proto \
		--call fileprocessing_v1.FileProcessingService.GetFileMetadata \
		--data '{"file_id": 0}' \
		--rps 100 \
		--total 3000 \
		--insecure \
		localhost:50051
