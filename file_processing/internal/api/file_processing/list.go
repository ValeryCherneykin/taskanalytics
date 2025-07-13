package fileprocessing

import (
	"context"
	"encoding/csv"
	"os"
	"strings"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListFiles(ctx context.Context, req *desc.ListFilesRequest) (*desc.ListFilesResponse, error) {
	files, err := i.fileProcessingService.List(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list files: %v", err)
	}

	var result []*desc.FileMetadata
	for _, file := range files {
		content, err := os.ReadFile(file.FilePath)
		if err != nil {
			continue
		}
		reader := csv.NewReader(strings.NewReader(string(content)))
		records, err := reader.ReadAll()
		if err != nil {
			continue
		}
		recordCount := int64(len(records))

		result = append(result, converter.ToFileMetadata(file, recordCount))
	}

	logger.Info("listed files",
		zap.Int("files_count", len(result)),
		zap.Uint64("limit", req.GetLimit()),
		zap.Uint64("offset", req.GetOffset()),
	)

	return &desc.ListFilesResponse{
		Files: result,
	}, nil
}
