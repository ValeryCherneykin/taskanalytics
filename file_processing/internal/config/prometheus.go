package config

import (
	"net"
	"os"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	prometheusHostEnvName = "PROMETHEUS_HOST"
	prometheusPortEnvName = "PROMETHEUS_PORT"
)

type PrometheusConfig interface {
	Address() string
}

type prometheusConfig struct {
	host string
	port string
}

func NewPrometheusConfig() (PrometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		logger.Error("prometheus host not found", zap.String("env_var", prometheusHostEnvName))
		return nil, errors.New("prometheus host not found")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		logger.Error("prometheus port not found", zap.String("env_var", prometheusPortEnvName))
		return nil, errors.New("prometheus port not found")
	}

	cfg := &prometheusConfig{
		host: host,
		port: port,
	}

	logger.Info("prometheus config loaded", zap.String("host", host), zap.String("port", port))

	return cfg, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
