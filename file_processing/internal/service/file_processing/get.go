package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.UploadedFile, error) {
	file, err := s.fileProcessingRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return file, nil
}
