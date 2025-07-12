package fileprocessing

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetFileMetadata(ctx context.Context, req *desc.GetFileRequest) (*desc.FileMetadataResponse, error) {
	if req.GetFileId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "file_id must be positive")
	}

	file, err := i.fileProcessingService.Get(ctx, req.GetFileId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	content, err := os.ReadFile(file.FilePath)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read file: %v", err)
	}

	reader := csv.NewReader(strings.NewReader(string(content)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid CSV format: %v", err)
	}
	recordCount := int64(len(records))

	log.Printf("retrieved metadata for file with id: %d, name: %s, records: %d", file.FileID, file.FileName, recordCount)

	return &desc.FileMetadataResponse{
		File: converter.ToFileMetadata(file, recordCount),
	}, nil
}
