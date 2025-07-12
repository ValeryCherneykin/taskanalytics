package fileprocessing

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxFileSize = 10 * 1024 * 1024

func (i *Implementation) UploadCSVFile(ctx context.Context, req *desc.UploadCSVFileRequest) (*desc.UploadCSVResponse, error) {
	if req.GetFileName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "file name cannot be empty")
	}
	content := req.GetContent()
	if len(content) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "file content cannot be empty")
	}
	if len(content) > maxFileSize {
		return nil, status.Errorf(codes.InvalidArgument, "file size exceeds limit of %d bytes", maxFileSize)
	}

	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid CSV format: %v", err)
	}
	recordCount := int64(len(records))

	file := converter.ToModelFromUploadRequest(req, i.storageConfig)

	if err := os.MkdirAll(filepath.Dir(file.FilePath), 0o755); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create directory: %v", err)
	}
	if err := os.WriteFile(file.FilePath, content, 0o644); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write file: %v", err)
	}

	id, err := i.fileProcessingService.Create(ctx, file)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save file metadata: %v", err)
	}

	log.Printf("uploaded file with id: %d, name: %s, records: %d", id, req.GetFileName(), recordCount)

	return &desc.UploadCSVResponse{
		FileId:  id,
		Message: "File uploaded successfully",
		Status:  "success",
	}, nil
}
