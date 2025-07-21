package storage

import (
	"context"
	"io"
)

type MinioClient interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error
	Download(ctx context.Context, objectName string) ([]byte, error)
	Delete(ctx context.Context, objectName string) error
}
