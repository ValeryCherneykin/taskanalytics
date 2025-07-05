package repository

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

type UploadedFileRepository interface {
	Create(ctx context.Context, file *model.UploadedFile) (int64, error)
	Get(ctx context.Context, id int64) (*model.UploadedFile, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, file model.UploadedFile) error
	List(ctx context.Context, limit, offset uint64) ([]model.UploadedFile, error)
}
