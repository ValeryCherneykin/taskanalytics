package fileprocessing

import (
	"context"

	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
)

func (s *Implementation) ProcessCSVFile(ctx context.Context, req *desc.ProcessCSVFileRequest) (*desc.UploadCSVResponse, error) {
	return &desc.UploadCSVResponse{
		FileId:  req.FileId,
		Message: "Processed successfully",
		Status:  "success",
	}, nil
}
