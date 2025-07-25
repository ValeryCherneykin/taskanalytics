package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	fileprocessing "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/api/file_processing"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/logger"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service"
	serviceMocks "github.com/ValeryCherneykin/taskanalytics/file_processing/internal/service/mocks"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestUpdateCSVFile(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	fileID := int64(1)
	existingFile := &model.UploadedFile{
		FileID:    fileID,
		FileName:  "old.csv",
		FilePath:  "mock-prefix/old.csv",
		CreatedAt: time.Time{},
	}
	validCSV := []byte("col1,col2\n1,2\n3,4\n")
	invalidCSV := []byte(`"invalid,"csv",content"`)

	tests := []struct {
		name       string
		req        *desc.UpdateCSVFileRequest
		setupMocks func(mc *minimock.Controller) service.FileProcessingService
		uploader   *fakeUploader
		wantErr    string
	}{
		{
			name: "success update with new filename",
			req: &desc.UpdateCSVFileRequest{
				FileId:     fileID,
				FileName:   "newfile.csv",
				NewContent: validCSV,
			},
			setupMocks: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(existingFile, nil)
				mock.UpdateMock.Expect(ctx, &model.UploadedFile{
					FileID:    existingFile.FileID,
					FileName:  "newfile.csv",
					FilePath:  existingFile.FilePath,
					Size:      int64(len(validCSV)),
					Status:    "updated",
					CreatedAt: existingFile.CreatedAt,
				}).Return(nil)
				return mock
			},
			uploader: &fakeUploader{
				storage: make(map[string][]byte),
			},
			wantErr: "",
		},
		{
			name: "file not found",
			req: &desc.UpdateCSVFileRequest{
				FileId:     fileID,
				NewContent: validCSV,
			},
			setupMocks: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(nil, errors.New("not found"))
				return mock
			},
			uploader: &fakeUploader{},
			wantErr:  "file not found",
		},
		{
			name: "invalid csv format",
			req: &desc.UpdateCSVFileRequest{
				FileId:     fileID,
				NewContent: invalidCSV,
			},
			setupMocks: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				return mock
			},
			uploader: &fakeUploader{},
			wantErr:  "invalid CSV format",
		},
		{
			name: "upload error",
			req: &desc.UpdateCSVFileRequest{
				FileId:     fileID,
				NewContent: validCSV,
			},
			setupMocks: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(existingFile, nil)
				return mock
			},
			uploader: &fakeUploader{
				uploadError: errors.New("upload failed"),
			},
			wantErr: "failed to upload updated file",
		},
		{
			name: "update metadata error",
			req: &desc.UpdateCSVFileRequest{
				FileId:     fileID,
				NewContent: validCSV,
			},
			setupMocks: func(mc *minimock.Controller) service.FileProcessingService {
				mock := serviceMocks.NewFileProcessingServiceMock(mc)
				mock.GetMock.Expect(ctx, fileID).Return(existingFile, nil)
				mock.UpdateMock.Expect(ctx, &model.UploadedFile{
					FileID:    existingFile.FileID,
					FileName:  existingFile.FileName,
					FilePath:  existingFile.FilePath,
					Size:      int64(len(validCSV)),
					Status:    "updated",
					CreatedAt: existingFile.CreatedAt,
				}).Return(errors.New("update error"))
				return mock
			},
			uploader: &fakeUploader{
				storage: make(map[string][]byte),
			},
			wantErr: "failed to update file metadata",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			defer mc.Finish()

			logger.Init(zapcore.NewNopCore())

			impl := fileprocessing.NewImplementation(tt.setupMocks(mc), nil, tt.uploader)

			got, err := impl.UpdateCSVFile(ctx, tt.req)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			require.Equal(t, fileID, got.FileId)
			require.Equal(t, "success", got.Status)
		})
	}
}
