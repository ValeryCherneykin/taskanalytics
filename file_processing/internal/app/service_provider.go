package app

import (
	"context"

	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db/pg"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/db/transaction"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/closer"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/config"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository"
	uploadFileRepository "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/repository/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	uploadFileService "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service/file_processing"
	"go.uber.org/zap"
)

type serviceProvider struct {
	logger *zap.Logger

	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	storageConfig config.StorageConfig

	dbClient                 db.Client
	txManager                db.TxManager
	fileProcessingRepository repository.UploadedFileRepository

	fileProcessingService service.FileProcessingService
	fileProcessingImpl    *fileprocessing.Implementation
}

func newServiceProvider(logger *zap.Logger) *serviceProvider {
	return &serviceProvider{logger: logger}
}

func (s *serviceProvider) PgConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			s.logger.Fatal("failed to get pg config", zap.Error(err))
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			s.logger.Fatal("failed to get grpc config: %s", zap.Error(err))
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) StorageConfig() config.StorageConfig {
	if s.storageConfig == nil {
		cfg, err := config.NewStorageConfig()
		if err != nil {
			s.logger.Fatal("failed to get storage config: %s", zap.Error(err))
		}
		s.storageConfig = cfg
	}
	return s.storageConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PgConfig().DSN())
		if err != nil {
			s.logger.Fatal("failed to create db client: %v", zap.Error(err))
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			s.logger.Fatal("ping error: %s", zap.Error(err))
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) FileProcessingRepository(ctx context.Context) repository.UploadedFileRepository {
	if s.fileProcessingRepository == nil {
		s.fileProcessingRepository = uploadFileRepository.NewRepository(s.DBClient(ctx))
	}
	return s.fileProcessingRepository
}

func (s *serviceProvider) FileProcessingService(ctx context.Context) service.FileProcessingService {
	if s.fileProcessingService == nil {
		s.fileProcessingService = uploadFileService.NewService(
			s.FileProcessingRepository(ctx),
			s.TxManager(ctx),
		)
	}
	return s.fileProcessingService
}

func (s *serviceProvider) FileProcessingImpl(ctx context.Context) *fileprocessing.Implementation {
	if s.fileProcessingImpl == nil {
		s.fileProcessingImpl = fileprocessing.NewImplementation(
			s.FileProcessingService(ctx),
			s.StorageConfig(),
		)
	}
	return s.fileProcessingImpl
}
