package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

func (s *serv) Update(ctx context.Context, file *model.UploadedFile) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		if err := s.fileProcessingRepository.Update(ctx, file); err != nil {
			return err
		}

		_, err := s.fileProcessingRepository.Get(ctx, file.FileID)
		if err != nil {
			return err
		}

		return nil
	})
}
