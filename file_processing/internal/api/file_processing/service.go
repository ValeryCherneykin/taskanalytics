package fileprocessing

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
)

type Implementation struct {
	desc.UnimplementedFileProcessingServiceServer
	fileProcessingService service.FileProcessingService
}

func NewInplementation(fileProcessingService service.FileProcessingService) *Implementation {
	return &Implementation{
		fileProcessingService: fileProcessingService,
	}
}
