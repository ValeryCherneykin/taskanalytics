package fileprocessing

import (
	"context"
	"log"
	"os"

	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DeleteFile(ctx context.Context, req *desc.DeleteFileRequest) (*desc.DeleteFileResponse, error) {
	if req.GetFileId() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "file_id must be positive")
	}

	file, err := i.fileProcessingService.Get(ctx, req.GetFileId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "file not found: %v", err)
	}

	if err := os.Remove(file.FilePath); err != nil && !os.IsNotExist(err) {
		return nil, status.Errorf(codes.Internal, "failed to delete file from disk: %v", err)
	}

	if err := i.fileProcessingService.Delete(ctx, req.GetFileId()); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete file metadata: %v", err)
	}

	log.Printf("deleted file with id: %d, name: %s", file.FileID, file.FileName)

	return &desc.DeleteFileResponse{
		Success: true,
		Message: "File deleted successfully",
	}, nil
}
