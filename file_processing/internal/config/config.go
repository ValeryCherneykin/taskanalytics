package config

import (
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		logger.Error("failed to load env file", zap.String("path", path), zap.Error(err))
		return err
	}

	logger.Info("env file loaded", zap.String("path", path))
	return nil
}
