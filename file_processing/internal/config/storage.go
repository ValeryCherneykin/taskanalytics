package config

import (
	"os"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	endpoint     = "MINIO_ENDPOINT"
	rootUser     = "MINIO_ROOT_USER"
	rootPassword = "MINIO_ROOT_PASSWORD"
	buckets      = "MINIO_DEFAULT_BUCKETS"
	useSSL       = "MINIO_USE_SSL"
)

type S3Config interface {
	Endpoint() string
	Bucket() string
	AccessKey() string
	SecretKey() string
	UseSSL() bool
}

type s3Config struct {
	endpoint  string
	bucket    string
	accessKey string
	secretKey string
	useSSL    bool
}

func NewS3Config() (S3Config, error) {
	endpoint := os.Getenv(endpoint)
	bucket := os.Getenv(buckets)
	accessKey := os.Getenv(rootUser)
	secretKey := os.Getenv(rootPassword)
	useSSL := os.Getenv(useSSL) == "true"

	if endpoint == "" || bucket == "" || accessKey == "" || len(secretKey) < 8 {
		logger.Error("missing required minio config")
		return nil, errors.New("invalid S3 config")
	}

	logger.Info("s3 config loaded", zap.String("endpoint", endpoint), zap.String("bucket", bucket))

	return &s3Config{
		endpoint:  endpoint,
		bucket:    bucket,
		accessKey: accessKey,
		secretKey: secretKey,
		useSSL:    useSSL,
	}, nil
}

func (c *s3Config) Endpoint() string  { return c.endpoint }
func (c *s3Config) Bucket() string    { return c.bucket }
func (c *s3Config) AccessKey() string { return c.accessKey }
func (c *s3Config) SecretKey() string { return c.secretKey }
func (c *s3Config) UseSSL() bool      { return c.useSSL }
