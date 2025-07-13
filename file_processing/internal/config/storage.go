package config

import (
	"os"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	storagePathEnvName = "STORAGE_PATH"
)

type StorageConfig interface {
	Path() string
}

type storageConfig struct {
	path string
}

func NewStorageConfig() (StorageConfig, error) {
	path := os.Getenv(storagePathEnvName)
	if len(path) == 0 {
		logger.Error("storage path not found", zap.String("env_var", storagePathEnvName))
		return nil, errors.New("storage path not found")
	}

	logger.Info("storage config loaded", zap.String("path", path))

	return &storageConfig{
		path: path,
	}, nil
}

func (cfg *storageConfig) Path() string {
	return cfg.path
}
