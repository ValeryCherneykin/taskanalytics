package fileprocessing

import (
	"context"
	"fmt"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"go.uber.org/zap"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	logger.Info("deleting file", zap.Int64("file_id", id))

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.fileProcessingRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		_, err := s.fileProcessingRepository.Get(ctx, id)
		if err == nil {
			return fmt.Errorf("file with id %d still exists after delete", id)
		}

		return nil
	})
	if err != nil {
		logger.Error("failed to delete file", zap.Int64("file_id", id), zap.Error(err))
		return err
	}

	logger.Info("file deleted", zap.Int64("file_id", id))
	return nil
}
