package service

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

type fileProcessingService interface {
	Create(ctx context.Context, file *model.UploadedFile) (int64, error)
	Get(ctx context.Context, id int64) (*model.UploadedFile, error)
}
