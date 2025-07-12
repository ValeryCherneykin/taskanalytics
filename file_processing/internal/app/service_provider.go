package app

import (
	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/config"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	storageConfig config.StorageConfig

	dbClient                 db.Client
	txManager                db.TxManager
	fileprocessingRepository repository.UploadedFileRepository

	fileprocessing     service.FileProcessingService
	fileProcessingImpl *fileprocessing.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}
