package repository

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

type UploadedFileRepository interface {
	Create(ctx context.Context, info model.UploadedFile) (int64, error)
	Get(ctx context.Context, id string) (model.UploadedFile, error)
	List(ctx context.Context, limit, offset int) ([]model.UploadedFile, error)
}
