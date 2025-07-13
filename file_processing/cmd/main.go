package main

import (
	"context"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/app"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"go.uber.org/zap"
)

func main() {
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger.Init(zapLogger.Core())

	ctx := context.Background()

	logger.Info("initializing app")
	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to init app", zap.Error(err))
	}

	logger.Info("starting app")

	err = a.Run()
	if err != nil {
		logger.Fatal("failed to run app", zap.Error(err))
	}
}
