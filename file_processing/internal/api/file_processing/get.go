package fileprocessing

import (
	"context"
	"encoding/csv"
	"strings"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetFileMetadata(ctx context.Context, req *desc.GetFileRequest) (*desc.FileMetadataResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	if req.GetFileId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "file_id must be positive")
	}

	file, err := i.fileProcessingService.Get(ctx, req.GetFileId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	content, err := i.minioClient.Download(ctx, file.FilePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read file from storage: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid CSV format: %v", err)
	}
	recordCount := int64(len(records))

	logger.Info("retrieved file metadata",
		zap.Int64("file_id", file.FileID),
		zap.String("file_name", file.FileName),
		zap.Int64("record_count", recordCount),
	)

	return &desc.FileMetadataResponse{
		File: converter.ToFileMetadata(file, recordCount),
	}, nil
}
