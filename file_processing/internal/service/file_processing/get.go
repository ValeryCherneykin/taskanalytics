package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"go.uber.org/zap"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.UploadedFile, error) {
	logger.Info("getting file", zap.Int64("file_id", id))

	file, err := s.fileProcessingRepository.Get(ctx, id)
	if err != nil {
		logger.Error("failed to get file", zap.Int64("file_id", id), zap.Error(err))
		return nil, err
	}

	logger.Info("file loaded", zap.Int64("file_id", file.FileID))
	return file, nil
}
