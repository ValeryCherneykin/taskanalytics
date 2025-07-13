package config

import (
	"errors"
	"os"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"go.uber.org/zap"
)

const (
	dsnEnvName = "PG_DSN"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		logger.Error("postgres DSN not found", zap.String("env_var", dsnEnvName))
		return nil, errors.New("pg dsn not found")
	}

	logger.Info("postgres config loaded", zap.String("dsn_prefix", dsnPrefix(dsn)))

	return &pgConfig{
		dsn: dsn,
	}, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}

func dsnPrefix(dsn string) string {
	if len(dsn) > 20 {
		return dsn[:20] + "..."
	}
	return dsn
}
