package tests

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	serviceMocks "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service/mocks"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap/zapcore"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDeleteFile(t *testing.T) {
	t.Parallel()

	type fileProcessingServiceMockFunc func(mc *minimock.Controller, fileID int64, filePath string) service.FileProcessingService

	type args struct {
		ctx context.Context
		req *desc.DeleteFileRequest
	}

	ctx := context.Background()

	wantRes := &desc.DeleteFileResponse{
		Success: true,
		Message: "File deleted successfully",
	}

	serviceErr := errors.New("failed to delete metadata")

	tests := []struct {
		name          string
		args          args
		fileID        int64
		filename      string
		filePath      string
		want          *desc.DeleteFileResponse
		errContains   string
		serviceMockFn fileProcessingServiceMockFunc
		prepare       func(filePath string) error
		cleanup       func(filePath string)
	}{
		{
			name:     "success case",
			fileID:   int64(gofakeit.Uint64()) + 1,
			filename: "test.csv",
			filePath: "test_data/" + gofakeit.UUID() + "/test.csv",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{},
			},
			want: wantRes,
			serviceMockFn: func(mc *minimock.Controller, fileID int64, filePath string) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(&model.UploadedFile{
					FileID:   fileID,
					FileName: "test.csv",
					FilePath: filePath,
				}, nil)
				mock.DeleteMock.Expect(ctx, fileID).Return(nil)
				return mock
			},
			prepare: func(filePath string) error {
				if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
					return err
				}
				return os.WriteFile(filePath, []byte("test"), 0o644)
			},
			cleanup: func(filePath string) {
				os.Remove(filePath)
				os.RemoveAll(filepath.Dir(filePath))
			},
		},
		{
			name:     "invalid file_id",
			fileID:   0,
			filename: "test.csv",
			filePath: "test_data/" + gofakeit.UUID() + "/test.csv",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{
					FileId: 0,
				},
			},
			errContains: "file_id must be positive",
			serviceMockFn: func(mc *minimock.Controller, fileID int64, filePath string) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
		},
		{
			name:     "file not found",
			fileID:   int64(gofakeit.Uint64()) + 1,
			filename: "test.csv",
			filePath: "test_data/" + gofakeit.UUID() + "/test.csv",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{},
			},
			errContains: "file not found",
			serviceMockFn: func(mc *minimock.Controller, fileID int64, filePath string) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(nil, errors.New("file not found"))
				return mock
			},
		},
		{
			name:     "file does not exist on disk",
			fileID:   int64(gofakeit.Uint64()) + 1,
			filename: "test.csv",
			filePath: "test_data/" + gofakeit.UUID() + "/test.csv",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{},
			},
			want: wantRes,
			serviceMockFn: func(mc *minimock.Controller, fileID int64, filePath string) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(&model.UploadedFile{
					FileID:   fileID,
					FileName: "test.csv",
					FilePath: filePath,
				}, nil)
				mock.DeleteMock.Expect(ctx, fileID).Return(nil)
				return mock
			},
		},
		{
			name:     "service delete error",
			fileID:   int64(gofakeit.Uint64()) + 1,
			filename: "test.csv",
			filePath: "test_data/" + gofakeit.UUID() + "/test.csv",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{},
			},
			errContains: "failed to delete file metadata",
			serviceMockFn: func(mc *minimock.Controller, fileID int64, filePath string) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(&model.UploadedFile{
					FileID:   fileID,
					FileName: "test.csv",
					FilePath: filePath,
				}, nil)
				mock.DeleteMock.Expect(ctx, fileID).Return(serviceErr)
				return mock
			},
			prepare: func(filePath string) error {
				if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
					return err
				}
				return os.WriteFile(filePath, []byte("test"), 0o644)
			},
			cleanup: func(filePath string) {
				os.Remove(filePath)
				os.RemoveAll(filepath.Dir(filePath))
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)

			tt.args.req.FileId = tt.fileID

			if tt.prepare != nil {
				require.NoError(t, tt.prepare(tt.filePath))
			}
			if tt.cleanup != nil {
				t.Cleanup(func() { tt.cleanup(tt.filePath) })
			}

			logger.Init(zapcore.NewNopCore())

			storageCfg := &fakeStorageConfig{basePath: "test_data"}
			serviceMock := tt.serviceMockFn(mc, tt.fileID, tt.filePath)

			impl := fileprocessing.NewImplementation(serviceMock, storageCfg)

			got, err := impl.DeleteFile(tt.args.ctx, tt.args.req)

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
