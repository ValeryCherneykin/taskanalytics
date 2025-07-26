package tests

import (
	"context"
	"fmt"
	"testing"

	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
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

func TestDeleteFile(t *testing.T) {
	t.Parallel()

	type fileProcessingServiceMockFunc func(mc *minimock.Controller) service.FileProcessingService

	type args struct {
		ctx context.Context
		req *desc.DeleteFileRequest
	}

	ctx := context.Background()
	fileID := gofakeit.Int64()
	if fileID <= 0 {
		fileID = 1
	}
	filename := "test.csv"
	filePath := "mock-prefix/test.csv"
	file := &model.UploadedFile{
		FileID:   fileID,
		FileName: filename,
		FilePath: filePath,
	}

	wantRes := &desc.DeleteFileResponse{
		Success: true,
		Message: "File deleted successfully",
	}

	serviceErr := fmt.Errorf("failed to retrieve or delete file")

	tests := []struct {
		name          string
		args          args
		want          *desc.DeleteFileResponse
		errContains   string
		serviceMockFn fileProcessingServiceMockFunc
		uploader      *fakeUploader
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{
					FileId: fileID,
				},
			},
			want: wantRes,
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(file, nil)
				mock.DeleteMock.Expect(ctx, fileID).Return(nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{filePath: {}},
			},
		},
		{
			name: "invalid file_id",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{
					FileId: 0,
				},
			},
			errContains: "value must be greater than 0", // обновлено здесь
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
			uploader: &fakeUploader{},
		},
		{
			name: "file not found",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{
					FileId: fileID,
				},
			},
			errContains: "file not found",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(nil, serviceErr)
				return mock
			},
			uploader: &fakeUploader{},
		},
		{
			name: "storage deletion error",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{
					FileId: fileID,
				},
			},
			errContains: "failed to delete file from storage",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(file, nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{},
			},
		},
		{
			name: "metadata deletion error",
			args: args{
				ctx: ctx,
				req: &desc.DeleteFileRequest{
					FileId: fileID,
				},
			},
			errContains: "failed to delete file metadata",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(file, nil)
				mock.DeleteMock.Expect(ctx, fileID).Return(serviceErr)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{filePath: {}},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Finish()

			logger.Init(zapcore.NewNopCore())

			s3Cfg := &fakeS3Config{
				bucket:    "mock-bucket",
				prefix:    "mock-prefix/",
				endpoint:  "mock-endpoint",
				accessKey: "mock-access",
				secretKey: "mock-secret",
				useSSL:    false,
			}

			impl := fileprocessing.NewImplementation(tt.serviceMockFn(mc), s3Cfg, tt.uploader)

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
