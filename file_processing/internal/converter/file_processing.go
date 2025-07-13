package converter

import (
	"path/filepath"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/config"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToModelFromUploadRequest(req *desc.UploadCSVFileRequest, storageConfig config.StorageConfig) *model.UploadedFile {
	fileName := filepath.Base(req.GetFileName())
	if fileName == "." || fileName == ".." || fileName == "" {
		logger.Warn("invalid filename in request, using fallback", zap.String("original", req.GetFileName()))
		fileName = "unnamed.csv"
	}
	filePath := filepath.Join(storageConfig.Path(), uuid.New().String(), fileName)
	return &model.UploadedFile{
		FileName: fileName,
		FilePath: filePath,
		Size:     int64(len(req.GetContent())),
		Status:   "pending",
	}
}

func ToFileMetadata(file *model.UploadedFile, recordCount int64) *desc.FileMetadata {
	return &desc.FileMetadata{
		FileId:      file.FileID,
		FileName:    file.FileName,
		FilePath:    file.FilePath,
		UploadedAt:  timestamppb.New(file.CreatedAt),
		Status:      file.Status,
		RecordCount: recordCount,
		Size:        file.Size,
	}
}
