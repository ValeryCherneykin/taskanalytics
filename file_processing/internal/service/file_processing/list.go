package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"go.uber.org/zap"
)

func (s *serv) List(ctx context.Context, limit, offset uint64) ([]*model.UploadedFile, error) {
	logger.Info("listing files", zap.Uint64("limit", limit), zap.Uint64("offset", offset))

	files, err := s.fileProcessingRepository.List(ctx, limit, offset)
	if err != nil {
		logger.Error("failed to list files", zap.Error(err))
		return nil, err
	}

	logger.Info("files listed", zap.Int("count", len(files)))
	return files, nil
}
