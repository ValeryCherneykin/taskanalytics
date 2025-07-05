package fileprocessing

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
)

func (s *serv) List(ctx context.Context, limit, offset uint64) ([]*model.UploadedFile, error) {
	files, err := s.fileProcessingRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return files, nil
}
