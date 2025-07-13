package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"go.uber.org/zap"
)

func (s *serv) Create(ctx context.Context, file *model.UploadedFile) (int64, error) {
	logger.Info("creating file", zap.String("file_name", file.FileName))

	var id int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.fileProcessingRepository.Create(ctx, file)
		if errTx != nil {
			return errTx
		}
		_, errTx = s.fileProcessingRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		logger.Error("failed to create file", zap.String("file_name", file.FileName), zap.Error(err))
		return 0, err
	}

	logger.Info("file created", zap.Int64("file_id", id))
	return id, nil
}
