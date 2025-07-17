package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	serviceMocks "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service/mocks"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestGetFileMetadata(t *testing.T) {
	t.Parallel()

	type fileProcessingServiceMockFunc func(mc *minimock.Controller) service.FileProcessingService

	type args struct {
		ctx context.Context
		req *desc.GetFileRequest
	}

	ctx := context.Background()
	basePath := "test_data"
	fileID := int64(gofakeit.Number(1, 1000))
	filename := "sample.csv"
	content := "col1,col2\n1,2\n3,4"
	recordCount := int64(3)

	filePath := filepath.Join(basePath, filename)

	createdAt := time.Now().UTC()

	fileModel := &model.UploadedFile{
		FileID:    fileID,
		FileName:  filename,
		FilePath:  filePath,
		Size:      int64(len(content)),
		Status:    "success",
		CreatedAt: createdAt,
	}

	want := &desc.FileMetadataResponse{
		File: converter.ToFileMetadata(fileModel, recordCount),
	}

	require.NoError(t, os.MkdirAll(basePath, 0o755))
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o644))
	t.Cleanup(func() {
		_ = os.RemoveAll(basePath)
	})

	tests := []struct {
		name          string
		args          args
		want          *desc.FileMetadataResponse
		errContains   string
		serviceMockFn fileProcessingServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{FileId: fileID},
			},
			want: want,
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					require.Equal(t, fileID, id)
					return fileModel, nil
				})
				return mock
			},
		},
		{
			name: "invalid file id",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{FileId: 0},
			},
			errContains: "file_id must be positive",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
		},
		{
			name: "file not found in service",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{FileId: fileID},
			},
			errContains: "file not found",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return nil, fmt.Errorf("not found")
				})
				return mock
			},
		},
		{
			name: "read file error",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{FileId: fileID},
			},
			errContains: "failed to read file",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				badPath := filepath.Join(basePath, "missing.csv")
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return &model.UploadedFile{
						FileID:    fileID,
						FileName:  "missing.csv",
						FilePath:  badPath,
						Status:    "pending",
						CreatedAt: time.Now(),
					}, nil
				})
				return mock
			},
		},
		{
			name: "invalid CSV format",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{FileId: fileID},
			},
			errContains: "invalid CSV format",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				badCSV := "col1,col2\n1,2\n3"
				badPath := filepath.Join(basePath, "bad.csv")
				require.NoError(t, os.WriteFile(badPath, []byte(badCSV), 0o644))

				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return &model.UploadedFile{
						FileID:    fileID,
						FileName:  "bad.csv",
						FilePath:  badPath,
						Status:    "pending",
						CreatedAt: time.Now(),
					}, nil
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			logger.Init(zapcore.NewNopCore())

			storageCfg := &fakeStorageConfig{basePath: basePath}
			serviceMock := tt.serviceMockFn(mc)
			impl := fileprocessing.NewImplementation(serviceMock, storageCfg)

			got, err := impl.GetFileMetadata(tt.args.ctx, tt.args.req)

			if tt.errContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
