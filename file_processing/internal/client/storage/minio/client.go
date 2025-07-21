package minio

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/client/storage"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/config"
)

var _ storage.MinioClient = (*client)(nil)

type handler func(ctx context.Context, client *minio.Client) error

type client struct {
	minio  *minio.Client
	config config.S3Config
}

func NewClient(cfg config.S3Config) (*client, error) {
	minioClient, err := minio.New(cfg.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey(), cfg.SecretKey(), ""),
		Secure: cfg.UseSSL(),
	})
	if err != nil {
		return nil, err
	}

	return &client{
		minio:  minioClient,
		config: cfg,
	}, nil
}

func (c *client) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	return c.execute(ctx, func(ctx context.Context, cli *minio.Client) error {
		_, err := cli.PutObject(ctx, c.config.Bucket(), objectName, reader, size, minio.PutObjectOptions{
			ContentType: contentType,
		})
		return err
	})
}

func (c *client) Download(ctx context.Context, objectName string) ([]byte, error) {
	var buf bytes.Buffer

	err := c.execute(ctx, func(ctx context.Context, cli *minio.Client) error {
		object, err := cli.GetObject(ctx, c.config.Bucket(), objectName, minio.GetObjectOptions{})
		if err != nil {
			return err
		}
		defer object.Close()

		_, err = io.Copy(&buf, object)
		return err
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *client) Delete(ctx context.Context, objectName string) error {
	return c.execute(ctx, func(ctx context.Context, cli *minio.Client) error {
		return cli.RemoveObject(ctx, c.config.Bucket(), objectName, minio.RemoveObjectOptions{})
	})
}

func (c *client) execute(ctx context.Context, h handler) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if c.minio == nil {
		return errors.New("minio client is nil")
	}

	err := h(timeoutCtx, c.minio)
	if err != nil {
		log.Printf("minio error: %v", err)
		return err
	}

	return nil
}

