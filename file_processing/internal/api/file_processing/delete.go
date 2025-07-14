package fileprocessing

import (
	"context"
	"os"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DeleteFile(ctx context.Context, req *desc.DeleteFileRequest) (*desc.DeleteFileResponse, error) {
	if req.GetFileId() <= 0 {
		logger.Error("Invalid file_id", zap.Int64("file_id", req.GetFileId()))
		return nil, status.Errorf(codes.InvalidArgument, "file_id must be positive")
	}

	file, err := i.fileProcessingService.Get(ctx, req.GetFileId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		logger.Error("Failed to delete file", zap.String("path", file.FilePath), zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete file from disk: %v", err)
	}

	if err := i.fileProcessingService.Delete(ctx, req.GetFileId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete file metadata: %v", err)
	}

	logger.Info("deleted file",
		zap.Int64("file_id", file.FileID),
		zap.String("file_name", file.FileName),
	)

	return &desc.DeleteFileResponse{
		Success: true,
		Message: "File deleted successfully",
	}, nil
}
