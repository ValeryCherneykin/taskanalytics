package converter

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	modelRepo "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository/file_processing/model"
)

func ToFileMetadataFromRepo(UploadedFile *modelRepo.UploadedFile) *model.UploadedFile {
	return &model.UploadedFile{
		FileID:    UploadedFile.FileID,
		FileName:  UploadedFile.FileName,
		FilePath:  UploadedFile.FilePath,
		Size:      UploadedFile.Size,
		Status:    UploadedFile.Status,
		CreatedAt: UploadedFile.CreatedAt,
		UpdatedAt: UploadedFile.UpdatedAt,
		DeletedAt: UploadedFile.DeletedAt,
	}
}
