package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	serviceMocks "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service/mocks"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"go.uber.org/zap/zapcore"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestDeleteFile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	basePath := "test_data"
	testFileID := int64(123)
	testFilePath := basePath + "/test.csv"
	testFileName := "test.csv"

	testFile := &model.UploadedFile{
		FileID:   testFileID,
		FilePath: testFilePath,
		FileName: testFileName,
	}

	tests := []struct {
		name          string
		req           *desc.DeleteFileRequest
		setupMock     func(mc *minimock.Controller) service.FileProcessingService
		wantErr       bool
		errContains   string
		expectedReply *desc.DeleteFileResponse
		setupFile     bool
	}{
		{
			name: "successfully delete file",
			req: &desc.DeleteFileRequest{
				FileId: testFileID,
			},
			setupFile: true,
			setupMock: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return testFile, nil
				})
				mock.DeleteMock.Set(func(ctx context.Context, id int64) error {
					return nil
				})
				return mock
			},
			expectedReply: &desc.DeleteFileResponse{
				Success: true,
				Message: "File deleted successfully",
			},
		},
		{
			name: "file_id is invalid",
			req: &desc.DeleteFileRequest{
				FileId: -1,
			},
			setupMock: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
			wantErr:     true,
			errContains: "file_id must be positive",
		},
		{
			name: "file not found in metadata",
			req: &desc.DeleteFileRequest{
				FileId: testFileID,
			},
			setupMock: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return nil, fmt.Errorf("not found")
				})
				return mock
			},
			wantErr:     true,
			errContains: "file not found",
		},
		{
			name: "os.Remove fails",
			req: &desc.DeleteFileRequest{
				FileId: testFileID,
			},
			setupMock: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return testFile, nil
				})
				mock.DeleteMock.Set(func(ctx context.Context, id int64) error {
					return nil
				})
				return mock
			},
			setupFile: false,
			wantErr:   false,
			expectedReply: &desc.DeleteFileResponse{
				Success: true,
				Message: "File deleted successfully",
			},
		},
		{
			name: "error during metadata delete",
			req: &desc.DeleteFileRequest{
				FileId: testFileID,
			},
			setupFile: true,
			setupMock: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Set(func(ctx context.Context, id int64) (*model.UploadedFile, error) {
					return testFile, nil
				})
				mock.DeleteMock.Set(func(ctx context.Context, id int64) error {
					return fmt.Errorf("db error")
				})
				return mock
			},
			wantErr:     true,
			errContains: "failed to delete file metadata",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.setupFile {
				_ = os.MkdirAll(basePath, 0o755)
				err := os.WriteFile(testFilePath, []byte("test data"), 0o644)
				require.NoError(t, err)
				t.Cleanup(func() { os.Remove(testFilePath) })
			}

			mc := minimock.NewController(t)
			t.Cleanup(mc.Finish)

			logger.Init(zapcore.NewNopCore())

			storageCfg := &fakeStorageConfig{basePath: basePath}
			serviceMock := tt.setupMock(mc)

			impl := fileprocessing.NewImplementation(serviceMock, storageCfg)

			resp, err := impl.DeleteFile(ctx, tt.req)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedReply, resp)
		})
	}
}
