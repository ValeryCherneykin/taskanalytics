package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"go.uber.org/zap"
)

func (s *serv) Update(ctx context.Context, file *model.UploadedFile) error {
	logger.Info("updating file", zap.Int64("file_id", file.FileID))

	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.fileProcessingRepository.Update(ctx, file); err != nil {
			return err
		}

		_, err := s.fileProcessingRepository.Get(ctx, file.FileID)
		if err != nil {
			logger.Error("failed to update file", zap.Int64("file_id", file.FileID), zap.Error(err))
			return err
		}

		logger.Info("file updated", zap.Int64("file_id", file.FileID))
		return nil
	})
}
