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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetFileMetadata(t *testing.T) {
	t.Parallel()

	type fileProcessingServiceMockFunc func(mc *minimock.Controller) service.FileProcessingService

	type args struct {
		ctx context.Context
		req *desc.GetFileRequest
	}

	ctx := context.Background()
	fileID := gofakeit.Int64()
	if fileID <= 0 {
		fileID = 1
	}
	filename := "test.csv"
	filePath := "mock-prefix/test.csv"
	content := []byte("Column1,Column2\n1,2\n3,4")
	recordCount := int64(3)

	file := &model.UploadedFile{
		FileID:   fileID,
		FileName: filename,
		FilePath: filePath,
	}

	uploadedAt := timestamppb.Now()

	wantRes := &desc.FileMetadataResponse{
		File: &desc.FileMetadata{
			FileId:      fileID,
			FileName:    filename,
			FilePath:    filePath,
			RecordCount: recordCount,
			UploadedAt:  uploadedAt,
			Status:      "",
		},
	}

	tests := []struct {
		name          string
		args          args
		want          *desc.FileMetadataResponse
		errContains   string
		serviceMockFn fileProcessingServiceMockFunc
		uploader      *fakeUploader
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{
					FileId: fileID,
				},
			},
			want: wantRes,
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(file, nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{filePath: content},
			},
		},
		{
			name: "invalid file_id",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{
					FileId: 0,
				},
			},
			errContains: "value must be greater than 0",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				return serviceMocks.NewFileProcessingServiceMock(mc)
			},
			uploader: &fakeUploader{},
		},
		{
			name: "file not found",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{
					FileId: fileID,
				},
			},
			errContains: "file not found",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(nil, fmt.Errorf("file not found"))
				return mock
			},
			uploader: &fakeUploader{},
		},
		{
			name: "storage read error",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{
					FileId: fileID,
				},
			},
			errContains: "failed to read file from storage",
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
			name: "invalid csv format",
			args: args{
				ctx: ctx,
				req: &desc.GetFileRequest{
					FileId: fileID,
				},
			},
			errContains: "invalid CSV format",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(file, nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{filePath: []byte(`"invalid,"csv",content"`)},
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

			got, err := impl.GetFileMetadata(tt.args.ctx, tt.args.req)

			if tt.errContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want.File.FileId, got.File.FileId)
			require.Equal(t, tt.want.File.FileName, got.File.FileName)
			require.Equal(t, tt.want.File.FilePath, got.File.FilePath)
			require.Equal(t, tt.want.File.RecordCount, got.File.RecordCount)
			require.Equal(t, tt.want.File.Status, got.File.Status)
			require.NotNil(t, got.File.UploadedAt)
			require.True(t, got.File.UploadedAt.IsValid())
		})
	}
}
