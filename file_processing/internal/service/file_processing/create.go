package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

func (s *serv) Create(ctx context.Context, file *model.UploadedFile) (int64, error) {
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
		return 0, err
	}
	return id, nil
}
