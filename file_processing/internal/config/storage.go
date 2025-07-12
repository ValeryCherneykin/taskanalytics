package config

import (
	"os"

	"github.com/pkg/errors"
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
		return nil, errors.New("storage path not found")
	}
	return &storageConfig{
		path: path,
	}, nil
}

func (cfg *storageConfig) Path() string {
	return cfg.path
}
