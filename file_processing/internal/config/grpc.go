package config

import (
	"net"
	"os"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		logger.Error("grpc host not found", zap.String("env_var", grpcHostEnvName))
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		logger.Error("grpc port not found", zap.String("env_var", grpcPortEnvName))
		return nil, errors.New("grpc port not found")
	}

	cfg := &grpcConfig{
		host: host,
		port: port,
	}

	logger.Info("grpc config loaded", zap.String("host", host), zap.String("port", port))

	return cfg, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
