package fileprocessing

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) UpdateCSVFile(ctx context.Context, req *desc.UpdateCSVFileRequest) (*desc.UploadCSVResponse, error) {
	if req.GetFileId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "file_id must be positive")
	}
	newContent := req.GetNewContent()
	if len(newContent) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "new_content cannot be empty")
	}
	if len(newContent) > maxFileSize {
		return nil, status.Errorf(codes.InvalidArgument, "file size exceeds limit of %d bytes", maxFileSize)
	}

	reader := csv.NewReader(strings.NewReader(string(newContent)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid CSV format: %v", err)
	}
	recordCount := int64(len(records))

	file, err := i.fileProcessingService.Get(ctx, req.GetFileId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	newFileName := file.FileName
	if req.GetFileName() != "" {
		newFileName = filepath.Base(req.GetFileName())
		if newFileName == "." || newFileName == ".." || newFileName == "" {
			newFileName = "unnamed.csv"
		}
	}

	updatedFile := &model.UploadedFile{
		FileID:    file.FileID,
		FileName:  newFileName,
		FilePath:  file.FilePath,
		Size:      int64(len(newContent)),
		Status:    "updated",
		CreatedAt: file.CreatedAt,
	}

	if err := os.WriteFile(file.FilePath, newContent, 0o644); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write file: %v", err)
	}

	if err := i.fileProcessingService.Update(ctx, updatedFile); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update file metadata: %v", err)
	}

	logger.Info("updated file",
		zap.Int64("file_id", file.FileID),
		zap.String("file_name", newFileName),
		zap.Int64("record_count", recordCount),
	)

	return &desc.UploadCSVResponse{
		FileId:  file.FileID,
		Message: "File updated successfully",
		Status:  "success",
	}, nil
}
