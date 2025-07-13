package app

import (
	"context"
	"net"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/closer"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/config"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	logger          *zap.Logger
}

func NewApp(ctx context.Context) (*App, error) {
	logger, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	a := &App{logger: logger}

	err = a.initDeps(ctx)
	if err != nil {
		logger.Error("failed to init app", zap.Error(err))
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider(a.logger)
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterFileProcessingServiceServer(a.grpcServer, a.serviceProvider.FileProcessingImpl(ctx))
	return nil
}

func (a *App) runGRPCServer() error {
	addr := a.serviceProvider.GRPCConfig().Address()

	a.logger.Info("gRPC server starting...", zap.String("address", addr))

	list, err := net.Listen("tcp", addr)
	if err != nil {
		a.logger.Error("failed to listen", zap.Error(err))
		return err
	}

	if err := a.grpcServer.Serve(list); err != nil {
		a.logger.Error("failed to serve", zap.Error(err))
		return err
	}

	return nil
}
