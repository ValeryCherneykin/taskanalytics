package fileprocessing

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/storage"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/config"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
)

type Implementation struct {
	desc.UnimplementedFileProcessingServiceServer
	fileProcessingService service.FileProcessingService
	storageConfig         config.S3Config
	minioClient           storage.MinioClient
}

func NewImplementation(fileProcessingService service.FileProcessingService, storageConfig config.S3Config, minioClient storage.MinioClient) *Implementation {
	return &Implementation{
		fileProcessingService: fileProcessingService,
		storageConfig:         storageConfig,
		minioClient:           minioClient,
	}
}
