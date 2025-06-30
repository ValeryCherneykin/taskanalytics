package converter

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToFileMetadata(file model.UploadedFile, recordCount int64) *desc.FileMetadata {
	return &desc.FileMetadata{
		FileId:      file.FileID,
		FileName:    file.FileName,
		FilePath:    file.FilePath,
		Size:        file.Size,
		Status:      file.Status,
		UploadedAt:  timestamppb.New(file.CreatedAt),
		RecordCount: recordCount,
	}
}
