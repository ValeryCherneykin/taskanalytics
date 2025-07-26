package fileprocessing

import (
	"context"
	"encoding/csv"
	"strings"
	"unicode/utf8"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxFileSize = 10 * 1024 * 1024

func (i *Implementation) UploadCSVFile(ctx context.Context, req *desc.UploadCSVFileRequest) (*desc.UploadCSVResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	content := req.GetContent()

	if len(content) > maxFileSize {
		return nil, status.Errorf(codes.InvalidArgument, "file size exceeds limit of %d bytes", maxFileSize)
	}

	if !utf8.Valid(content) {
		logger.Error("invalid CSV content", zap.ByteString("content", content))
		return nil, status.Errorf(codes.InvalidArgument, "invalid CSV format: content is not valid UTF-8")
	}

	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid CSV format: %v", err)
	}

	if len(records) <= 1 {
		return nil, status.Errorf(codes.InvalidArgument, "CSV must contain header and at least one row")
	}

	expectedColumns := len(records[0])
	for i, row := range records[1:] {
		if len(row) != expectedColumns {
			return nil, status.Errorf(codes.InvalidArgument, "inconsistent column count in row %d", i+2)
		}
	}

	file := converter.ToModelFromUploadRequest(req, i.storageConfig)

	err = i.minioClient.Upload(ctx, file.FilePath, strings.NewReader(string(content)), file.Size, "text/csv")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload file to storage: %v", err)
	}

	id, err := i.fileProcessingService.Create(ctx, file)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file metadata: %v", err)
	}

	logger.Info("uploaded file",
		zap.Int64("file_id", id),
		zap.String("name", req.GetFileName()),
		zap.Int64("records", int64(len(records)-1)),
	)

	return &desc.UploadCSVResponse{
		FileId:  id,
		Message: "File uploaded successfully",
		Status:  "success",
	}, nil
}
