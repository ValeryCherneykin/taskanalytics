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

func TestListFiles(t *testing.T) {
	t.Parallel()

	type fileProcessingServiceMockFunc func(mc *minimock.Controller) service.FileProcessingService

	type args struct {
		ctx context.Context
		req *desc.ListFilesRequest
	}

	ctx := context.Background()
	fileID1 := gofakeit.Int64()
	if fileID1 <= 0 {
		fileID1 = 1
	}
	fileID2 := fileID1 + 1
	filename1 := "test1.csv"
	filename2 := "test2.csv"
	filePath1 := "mock-prefix/test1.csv"
	filePath2 := "mock-prefix/test2.csv"
	content := []byte("Column1,Column2\n1,2\n3,4")
	recordCount := int64(3)

	file1 := &model.UploadedFile{
		FileID:   fileID1,
		FileName: filename1,
		FilePath: filePath1,
	}
	file2 := &model.UploadedFile{
		FileID:   fileID2,
		FileName: filename2,
		FilePath: filePath2,
	}

	uploadedAt := timestamppb.Now()

	wantRes := &desc.ListFilesResponse{
		Files: []*desc.FileMetadata{
			{
				FileId:      fileID1,
				FileName:    filename1,
				FilePath:    filePath1,
				RecordCount: recordCount,
				UploadedAt:  uploadedAt,
				Status:      "",
			},
			{
				FileId:      fileID2,
				FileName:    filename2,
				FilePath:    filePath2,
				RecordCount: recordCount,
				UploadedAt:  uploadedAt,
			},
		},
	}

	tests := []struct {
		name          string
		args          args
		want          *desc.ListFilesResponse
		errContains   string
		serviceMockFn fileProcessingServiceMockFunc
		uploader      *fakeUploader
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &desc.ListFilesRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			want: wantRes,
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.ListMock.Expect(ctx, uint64(10), uint64(0)).Return([]*model.UploadedFile{file1, file2}, nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{
					filePath1: content,
					filePath2: content,
				},
			},
		},
		{
			name: "service error",
			args: args{
				ctx: ctx,
				req: &desc.ListFilesRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			errContains: "failed to list files",
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.ListMock.Expect(ctx, uint64(10), uint64(0)).Return(nil, fmt.Errorf("service error"))
				return mock
			},
			uploader: &fakeUploader{},
		},
		{
			name: "partial download failure",
			args: args{
				ctx: ctx,
				req: &desc.ListFilesRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			want: &desc.ListFilesResponse{
				Files: []*desc.FileMetadata{
					{
						FileId:      fileID1,
						FileName:    filename1,
						FilePath:    filePath1,
						RecordCount: recordCount,
						UploadedAt:  uploadedAt,
						Status:      "",
					},
				},
			},
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.ListMock.Expect(ctx, uint64(10), uint64(0)).Return([]*model.UploadedFile{file1, file2}, nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{
					filePath1: content,
					// filePath2 отсутствует, что вызовет ошибку скачивания
				},
			},
		},
		{
			name: "partial invalid csv",
			args: args{
				ctx: ctx,
				req: &desc.ListFilesRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			want: &desc.ListFilesResponse{
				Files: []*desc.FileMetadata{
					{
						FileId:      fileID1,
						FileName:    filename1,
						FilePath:    filePath1,
						RecordCount: recordCount,
						UploadedAt:  uploadedAt,
						Status:      "",
					},
				},
			},
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.ListMock.Expect(ctx, uint64(10), uint64(0)).Return([]*model.UploadedFile{file1, file2}, nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: map[string][]byte{
					filePath1: content,
					filePath2: []byte(`"invalid,"csv",content"`), // Несбалансированные кавычки для ошибки CSV
				},
			},
		},
		{
			name: "empty file list",
			args: args{
				ctx: ctx,
				req: &desc.ListFilesRequest{
					Limit:  10,
					Offset: 0,
				},
			},
			want: &desc.ListFilesResponse{
				Files: []*desc.FileMetadata{},
			},
			serviceMockFn: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.ListMock.Expect(ctx, uint64(10), uint64(0)).Return([]*model.UploadedFile{}, nil)
				return mock
			},
			uploader: &fakeUploader{},
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

			got, err := impl.ListFiles(tt.args.ctx, tt.args.req)

			if tt.errContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tt.want.Files), len(got.Files))
			for i, wantFile := range tt.want.Files {
				require.Equal(t, wantFile.FileId, got.Files[i].FileId)
				require.Equal(t, wantFile.FileName, got.Files[i].FileName)
				require.Equal(t, wantFile.FilePath, got.Files[i].FilePath)
				require.Equal(t, wantFile.RecordCount, got.Files[i].RecordCount)
				require.Equal(t, wantFile.Status, got.Files[i].Status)
				require.NotNil(t, got.Files[i].UploadedAt)
				require.True(t, got.Files[i].UploadedAt.IsValid())
			}
		})
	}
}
