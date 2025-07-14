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
	"go.uber.org/zap/zapcore"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestUploadCSVFile(t *testing.T) {
	t.Parallel()

	type fileProcessingServiceMockFunc func(mc *minimock.Controller) service.FileProcessingService

	type args struct {
		ctx context.Context
		req *desc.UploadCSVFileRequest
	}

	ctx := context.Background()
	content := "Column1,Column2\n1,2\n3,4"
	filename := "test.csv"
	basePath := "test_data"
	fileBytes := []byte(content)

	validReq := &desc.UploadCSVFileRequest{
		FileName: filename,
		Content:  fileBytes,
	}

	fileID := gofakeit.Int64()
	wantRes := &desc.UploadCSVResponse{
		FileId:  fileID,
		Message: "File uploaded successfully",
		Status:  "success",
	}

	serviceErr := fmt.Errorf("failed to save metadata")

	tests := []struct {
		name          string
		args          args
		want          *desc.UploadCSVResponse
		errContains   string
		serviceMockFn fileProcessingServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: validReq,
			},
			want: wantRes,
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.CreateMock.Set(func(ctx context.Context, file *model.UploadedFile) (int64, error) {
					if file.FileName != filename {
						return 0, fmt.Errorf("unexpected filename: %s", file.FileName)
					}
					return fileID, nil
				})
				return mock
			},
		},
		{
			name: "empty filename",
			args: args{
				ctx: ctx,
				req: &desc.UploadCSVFileRequest{
					FileName: "",
					Content:  fileBytes,
				},
			},
			errContains: "file name cannot be empty",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
		},
		{
			name: "empty content",
			args: args{
				ctx: ctx,
				req: &desc.UploadCSVFileRequest{
					FileName: filename,
					Content:  []byte{},
				},
			},
			errContains: "file content cannot be empty",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
		},
		{
			name: "invalid csv format",
			args: args{
				ctx: ctx,
				req: &desc.UploadCSVFileRequest{
					FileName: filename,
					Content:  []byte{0xff, 0xfe},
				},
			},
			errContains: "invalid CSV format",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
		},
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: validReq,
			},
			errContains: "failed to save file metadata",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.CreateMock.Set(func(_ context.Context, file *model.UploadedFile) (int64, error) {
					if file.FileName != filename {
						return 0, fmt.Errorf("unexpected filename: %s", file.FileName)
					}
					return 0, serviceErr
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

			storageCfg := &fakeStorageConfig{basePath: basePath}
			serviceMock := tt.serviceMockFn(mc)

			logger.Init(zapcore.NewNopCore())

			impl := fileprocessing.NewImplementation(serviceMock, storageCfg)

			got, err := impl.UploadCSVFile(tt.args.ctx, tt.args.req)

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
