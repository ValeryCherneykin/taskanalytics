package tests

import (
	"context"
	"errors"
	"io"
)

type fakeS3Config struct {
	bucket    string
	prefix    string
	endpoint  string
	accessKey string
	secretKey string
	useSSL    bool
}

func (f *fakeS3Config) Bucket() string    { return f.bucket }
func (f *fakeS3Config) Prefix() string    { return f.prefix }
func (f *fakeS3Config) Endpoint() string  { return f.endpoint }
func (f *fakeS3Config) AccessKey() string { return f.accessKey }
func (f *fakeS3Config) SecretKey() string { return f.secretKey }
func (f *fakeS3Config) UseSSL() bool      { return f.useSSL }

type fakeUploader struct {
	storage     map[string][]byte
	uploadError error
}

func (f *fakeUploader) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	if f.uploadError != nil {
		return f.uploadError
	}
	if f.storage == nil {
		f.storage = make(map[string][]byte)
	}
	data := make([]byte, size)
	n, err := io.ReadFull(reader, data)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return err
	}
	if int64(n) != size {
		return errors.New("read size mismatch")
	}
	f.storage[objectName] = data
	return nil
}

func (f *fakeUploader) Download(ctx context.Context, objectName string) ([]byte, error) {
	data, ok := f.storage[objectName]
	if !ok {
		return nil, errors.New("object not found")
	}
	return data, nil
}

func (f *fakeUploader) Delete(ctx context.Context, objectName string) error {
	if _, ok := f.storage[objectName]; !ok {
		return errors.New("object not found")
	}
	delete(f.storage, objectName)
	return nil
}
