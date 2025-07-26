FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/file_processing_server cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/file_processing_server .
COPY .env .
CMD ["./file_processing_server"]
