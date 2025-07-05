package tests

import (
	"testing"
	"time"

	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/converter"
	"github.com/ValeryCherneykin/taskanalytics/file_processing/internal/model"
	desc "github.com/ValeryCherneykin/taskanalytics/file_processing/pkg/file_processing_v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestToFileMetadata(t *testing.T) {
	now := time.Now()
	file := model.UploadedFile{
		FileID:    1,
		FileName:  "data.csv",
		FilePath:  "/files/data.csv",
		Size:      1024,
		Status:    "completed",
		CreatedAt: now,
	}

	expected := &desc.FileMetadata{
		FileId:      1,
		FileName:    "data.csv",
		FilePath:    "/files/data.csv",
		Size:        1024,
		Status:      "completed",
		UploadedAt:  timestamppb.New(now),
		RecordCount: 42,
	}

	result := converter.ToFileMetadata(file, 42)

	assert.Equal(t, expected, result)
}
